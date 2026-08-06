# Portofolio API

[![Go Report Card](https://goreportcard.com/badge/github.com/azharf99/portofolio-api)](https://goreportcard.com/report/github.com/azharf99/portofolio-api)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

Portofolio API adalah backend service yang dibangun menggunakan bahasa pemrograman Go dengan prinsip **Clean Architecture**. API ini dirancang untuk mengelola data portofolio profesional, lengkap dengan fitur autentikasi, manajemen gambar, dan dukungan multibahasa (i18n).

## 🚀 Fitur Utama

- **Clean Architecture**: Pemisahan tanggung jawab yang jelas antara Domain, Usecase, Repository, dan Delivery.
- **Autentikasi JWT**: Keamanan akses admin menggunakan JSON Web Token.
- **Multibahasa (i18n)**: Dukungan pesan respons dalam bahasa Inggris (en), Indonesia (id), dan Rusia (ru).
- **Manajemen Portofolio**: CRUD lengkap termasuk fitur unggah gambar dan manajemen tech stack.
- **Keamanan**:
    - **Rate Limiting**: Perlindungan dari serangan brute-force atau spamming.
    - **Security Headers**: Implementasi header keamanan standar (XSS Protection, HSTS, dll).
    - **CORS**: Pengaturan Cross-Origin Resource Sharing yang fleksibel.
- **Database**: PostgreSQL dengan GORM untuk manajemen database yang handal.
- **Docker Ready**: Siap dijalankan di lingkungan kontainer menggunakan Docker & Docker Compose.
- **CI/CD**: Workflow GitHub Actions untuk pengujian otomatis.

## 🛠️ Tech Stack

- **Language**: Go (Golang) 1.26.1
- **Web Framework**: [Gin Gonic](https://gin-gonic.com/)
- **ORM**: [GORM](https://gorm.io/)
- **Database**: PostgreSQL
- **Auth**: JWT (JSON Web Token)
- **i18n**: go-i18n
- **Tools**: Docker, Docker Compose

## 📋 Struktur Proyek

```text
├── config/             # Konfigurasi database dan env
├── delivery/           # Layer delivery (HTTP Handlers)
├── domain/             # Layer entitas dan interface
├── locales/            # File translasi i18n (JSON)
├── middleware/         # Middleware (Auth, i18n, Security)
├── pkg/                # Package utilitas (i18n helper)
├── repository/         # Layer akses database
├── routes/             # Definisi rute API
├── usecase/            # Layer logika bisnis
└── main.go             # Entry point aplikasi
```

## ⚙️ Instalasi & Penggunaan

### Menggunakan Docker (Rekomendasi)

1. Clone repository:
   ```bash
   git clone https://github.com/azharf99/portofolio-api.git
   cd portofolio-api
   ```

2. Salin file `.env.example` ke `.env` dan sesuaikan konfigurasinya:
   ```bash
   cp .env.example .env
   ```

3. Jalankan aplikasi menggunakan Docker Compose:
   ```bash
   docker-compose up -d --build
   ```

4. API akan berjalan di `http://localhost:8080`.

### Tanpa Docker

1. Pastikan Anda memiliki Go 1.26.1+ dan PostgreSQL yang sedang berjalan.
2. Buat database di PostgreSQL.
3. Atur environment variables di file `.env`.
4. Jalankan aplikasi:
   ```bash
   go run main.go
   ```

## 🔐 Akun Admin Default

Saat pertama kali dijalankan (database kosong), sistem akan membuat SATU akun admin dari
env var `ADMIN_USERNAME` dan `ADMIN_PASSWORD`. **Kedua env var ini wajib diisi** —
aplikasi akan gagal start jika kosong, sama seperti `JWT_SECRET`. Password minimal 12
karakter. Tidak ada default hardcoded lagi — pilih kredensial Anda sendiri sebelum
menjalankan API di lingkungan manapun yang bisa diakses publik.

## 🛣️ API Endpoints

### Public Routes
- `POST /api/login` - Login admin. Token JWT dikirim sebagai cookie httpOnly (`auth_token`), bukan di body response.
- `POST /api/logout` - Menghapus cookie auth.
- `GET /api/portfolios` - Mengambil daftar portofolio yang dipublikasikan.
- `GET /uploads/*` - Mengakses file gambar statis.

### Admin Routes (Protected by JWT)
- `GET /api/admin/portfolios` - Daftar semua portofolio (termasuk yang belum dipublikasikan).
- `POST /api/admin/portfolios` - Menambah portofolio baru.
- `PUT /api/admin/portfolios/:id` - Memperbarui data portofolio.
- `DELETE /api/admin/portfolios/:id` - Menghapus portofolio.
- `PUT /api/admin/users/:id` - Memperbarui data user/admin.
- `DELETE /api/admin/users/:id` - Menghapus user/admin.

## 📜 Lisensi & Atribusi

Proyek ini dilisensikan di bawah **Apache License 2.0**.

### Syarat Penggunaan (Atribusi)
Sesuai dengan keinginan pengembang asli, **setiap individu atau organisasi yang menggunakan, memodifikasi, atau mendistribusikan ulang kode ini WAJIB mencantumkan atribusi kepada pengembang asli**:

**Azhar Faturohman Ahidin**  
GitHub: [azharf99](https://github.com/azharf99)

---
Dibuat dengan ❤️ oleh [Azhar Faturohman Ahidin](https://github.com/azharf99)
