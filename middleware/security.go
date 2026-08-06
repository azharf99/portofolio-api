package middleware

import (
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// SecurityHeaders menambahkan header keamanan standar industri
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		// KEAMANAN: Content-Security-Policy sebagai lapisan pertahanan terhadap XSS.
		// API ini hanya melayani JSON + file statis di /uploads, jadi kebijakan ketat aman dipakai.
		c.Header("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'")
		c.Next()
	}
}

// SetupCORS membatasi domain mana saja yang boleh mengakses API ini
func SetupCORS() gin.HandlerFunc {
	// KEAMANAN: Konfigurasi CORS Dinamis
	allowedOriginsEnv := os.Getenv("ALLOWED_ORIGINS")
	var allowedOrigins []string
	if allowedOriginsEnv == "" {
		allowedOrigins = []string{"http://localhost:5173"} // Fallback aman
	} else {
		allowedOrigins = strings.Split(allowedOriginsEnv, ",")
	}
	return cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true, // dibutuhkan agar cookie httpOnly auth_token ikut terkirim
		MaxAge:           12 * time.Hour,
	})
}

// --- Rate Limiter In-Memory ---
// Melindungi server dari spam request (misal: frontend stuck di infinite loop refresh token,
// atau percobaan enumerasi endpoint publik seperti transaction history).

type limiterEntry struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type limiterStore struct {
	mu       sync.Mutex
	visitors map[string]*limiterEntry
	r        rate.Limit
	burst    int
}

func newLimiterStore(r rate.Limit, burst int) *limiterStore {
	s := &limiterStore{
		visitors: make(map[string]*limiterEntry),
		r:        r,
		burst:    burst,
	}
	go s.cleanupLoop()
	return s
}

// cleanupLoop mencegah map visitors tumbuh tanpa batas (memory leak) dengan
// membuang entry yang sudah tidak aktif selama >10 menit.
func (s *limiterStore) cleanupLoop() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		s.mu.Lock()
		for ip, entry := range s.visitors {
			if time.Since(entry.lastSeen) > 10*time.Minute {
				delete(s.visitors, ip)
			}
		}
		s.mu.Unlock()
	}
}

func (s *limiterStore) allow(key string) bool {
	s.mu.Lock()
	entry, exists := s.visitors[key]
	if !exists {
		entry = &limiterEntry{limiter: rate.NewLimiter(s.r, s.burst)}
		s.visitors[key] = entry
	}
	entry.lastSeen = time.Now()
	limiter := entry.limiter
	s.mu.Unlock()

	return limiter.Allow()
}

// Limiter global: 3 request/detik, burst 6 — melindungi seluruh API dari spam umum.
var globalLimiter = newLimiterStore(3, 6)

func RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !globalLimiter.allow(c.ClientIP()) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Terlalu banyak permintaan, silakan coba beberapa saat lagi.",
			})
			return
		}
		c.Next()
	}
}

// Limiter ketat khusus endpoint yang rawan enumerasi (mis. lookup transaksi via email):
// 1 request setiap 2 menit dengan burst 3, per IP.
var historyLimiter = newLimiterStore(rate.Every(2*time.Minute), 3)

func HistoryRateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !historyLimiter.allow(c.ClientIP()) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Terlalu banyak permintaan pengecekan riwayat. Silakan coba lagi dalam beberapa menit.",
			})
			return
		}
		c.Next()
	}
}
