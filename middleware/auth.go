package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthCookieName adalah nama cookie httpOnly tempat JWT admin disimpan.
// KEAMANAN: Token disimpan sebagai cookie httpOnly (bukan di localStorage/body response)
// supaya tidak bisa dibaca oleh JavaScript — ini menutup celah pencurian token via XSS.
const AuthCookieName = "auth_token"

const authCookieMaxAgeSeconds = 24 * 60 * 60 // 24 jam, samakan dengan exp klaim JWT

// isSecureEnv menentukan apakah cookie harus diberi flag Secure (hanya dikirim via HTTPS).
// Default aman: true, kecuali eksplisit dinonaktifkan untuk pengembangan lokal via http://localhost.
func isSecureEnv() bool {
	return os.Getenv("COOKIE_INSECURE") != "true"
}

// SetAuthCookie menaruh JWT sebagai cookie httpOnly + Secure + SameSite=Strict.
func SetAuthCookie(c *gin.Context, token string) {
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie(AuthCookieName, token, authCookieMaxAgeSeconds, "/", "", isSecureEnv(), true)
}

// ClearAuthCookie menghapus cookie auth saat logout.
func ClearAuthCookie(c *gin.Context) {
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie(AuthCookieName, "", -1, "/", "", isSecureEnv(), true)
}

// AuthMiddleware mengecek dan memvalidasi JWT Token.
// Mendukung dua sumber token (cookie httpOnly diprioritaskan, header Authorization
// sebagai fallback untuk kebutuhan testing manual/API client seperti curl/Postman).
func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := extractToken(c)
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token tidak ditemukan"})
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("metode signing tidak valid")
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
			return
		}

		userID, exists := claims["sub"]
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
			return
		}

		c.Set("user_id", userID)
		c.Set("jwt_claims", claims)
		c.Next()
	}
}

func extractToken(c *gin.Context) string {
	if cookie, err := c.Cookie(AuthCookieName); err == nil && cookie != "" {
		return cookie
	}

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return ""
	}
	return strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
}
