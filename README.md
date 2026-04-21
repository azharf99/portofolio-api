# 🚀 Portofolio API

Sebuah RESTful API tangguh dan aman yang dibangun menggunakan **Golang** untuk mengelola data portofolio profesional. Proyek ini menerapkan prinsip **Clean Architecture** dan praktik terbaik *Cybersecurity* untuk memastikan skalabilitas, kemudahan *testing*, dan keamanan data.

## 🌟 Fitur Utama

- **Clean Architecture**: Kode dipisahkan menjadi layer `Domain`, `Repository`, `Usecase`, dan `Delivery` agar mudah di- *maintenance* dan diuji.
- **Autentikasi Aman**: Menggunakan **JWT (JSON Web Tokens)** untuk perlindungan *route* *private* dan **Bcrypt** untuk *hashing password*.
- **Security First**: Dilengkapi dengan *Middleware Security Headers* (mencegah XSS, Clickjacking) dan konfigurasi CORS yang ketat.
- **Fitur Lengkap CRUD Portofolio**: Mendukung penambahan, pembaruan, penghapusan, dan pengambilan data.
- **Advanced Fetching**: Mendukung fitur *Searching*, *Filtering* (berdasarkan tipe dan industri), serta *Pagination* bawaan pada *endpoint public*.
- **Database Relasional**: Menggunakan **PostgreSQL** yang diintegrasikan dengan **GORM** (mencegah SQL Injection).
- **Dockerized**: Siap untuk di-*deploy* menggunakan Docker dan Docker Compose.
- **Unit Testing**: Dilengkapi dengan *unit test* standar industri menggunakan *Mocking* (`testify`) untuk mencapai *code coverage* yang tinggi.

## 🛠️ Tech Stack

- **Bahasa**: Go (Golang)
- **Framework Web**: [Gin Gonic](https://gin-gonic.com/)
- **ORM**: [GORM](https://gorm.io/)
- **Database**: PostgreSQL
- **Security**: `golang-jwt`, `bcrypt`
- **Testing**: `testify`
- **Containerization**: Docker, Docker Compose

## 📁 Struktur Direktori (Clean Architecture)

```text
.
├── config/           # Setup database & koneksi
├── delivery/         # Layer HTTP Controller (Gin Handlers)
├── domain/           # Entitas inti (Struct) dan Interface (Kontrak)
│   └── mocks/        # File mock untuk Unit Testing
├── middleware/       # JWT Auth, CORS, dan Security Headers
├── repository/       # Interaksi langsung ke database (GORM)
├── routes/           # Registrasi semua endpoint API
├── usecase/          # Layer logika bisnis
├── Dockerfile        # Konfigurasi image aplikasi
├── docker-compose.yml# Orkestrasi Database dan API
└── main.go           # Titik masuk aplikasi (Entry Point)
```

## 🚀 Cara Menjalankan Aplikasi
- Persyaratan SistemGo (v1.21 atau lebih baru)
- Docker & Docker Compose

# 1. Menjalankan via Docker (Direkomendasikan)
Cara termudah untuk menjalankan API beserta databasenya adalah menggunakan Docker Compose.

**Clone repository**
```bash
git clone https://github.com/azharf99/portofolio-api.git
cd portofolio-api
```

**Jalankan container di background**
```bash
docker-compose up -d --build
```

**Cek log aplikasi untuk memastikan berjalan lancar**
```bash
docker-compose logs -f api
```

**Aplikasi akan berjalan di http://localhost:8080.**

# 2. Menjalankan Unit Test

Proyek ini mengutamakan kualitas kode. Untuk menjalankan unit test dan melihat code coverage:

**Download semua dependency**
```bash
go mod tidy
```

**Jalankan seluruh test**
```bash
go test ./... -cover
```

## 📡 Dokumentasi API (Endpoints)

# Public Routes (Tidak butuh token)

| Method | Endpoint | Deskripsi | Query Params (Opsional) |
| POST | /api/login | Login admin (Mendapatkan JWT) | - |
| GET | /api/portfolios | Mengambil data portofolio | page, limit, search, industry, type |

# Private Routes (Butuh Header Authorization: Bearer <token>)

| Method | Endpoint | Deskripsi |
| POST | /api/admin/portfolios | Menambahkan portofolio baru |
| PUT | /api/admin/portfolios/:id | Memperbarui data portofolio |
| DELETE | /api/admin/portfolios/:id | Menghapus portofolio |
| PUT | /api/admin/users/:id | Memperbarui data user |
| DELETE | /api/admin/users/:id | Menghapus user |

## 🔒 Variabel Lingkungan (Environment Variables)

Aplikasi ini membaca konfigurasi dari environment variables. Saat menggunakan Docker Compose, ini sudah diatur secara otomatis. Jika menjalankan secara lokal, Anda bisa mengatur:
`DB_HOST`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_PORT`, `JWT_SECRET` (Sangat disarankan untuk diganti di lingkungan produksi), `GIN_MODE` (Set ke release saat deployment)


**Dibuat dengan ❤️ oleh azharf99**