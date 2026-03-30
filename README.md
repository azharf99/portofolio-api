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