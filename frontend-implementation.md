# Frontend Implementation Plan

Dokumen ini adalah panduan implementasi frontend berdasarkan endpoint backend yang tersedia pada proyek ini, termasuk perbaikan perilaku error terbaru.

## 1. Tujuan Produk Frontend

- Menyediakan halaman publik untuk menampilkan daftar portfolio.
- Menyediakan panel admin untuk login, tambah/edit/hapus portfolio, dan manajemen user.
- Menjaga pengalaman penggunaan yang aman, cepat, dan konsisten dengan kontrak API backend.

## 2. Ringkasan Endpoint Backend (Sumber Kebenaran)

### Public

- `POST /api/login`
  - Body:
    - `username: string` (required)
    - `password: string` (required)
  - Success `200`: `{ "token": "..." }`
  - Error: `400`, `401`

- `GET /api/portfolios?page=1&limit=10&search=&industry=&type=`
  - Success `200`: `{ data: Portfolio[], page: number, limit: number, total: number }`
  - Error: `400` untuk query `page/limit` tidak valid.

### Private (Header `Authorization: Bearer <token>`)

- `POST /api/admin/portfolios`
- `PUT /api/admin/portfolios/:id`
  - Body: `multipart/form-data` (Bukan JSON lagi)
  - Fields:
    - `title`, `description`, `role`, `type`, `industry`, `tech_stack`, `project_link`, `is_published`: string/text
    - `start_date`, `end_date`: string (format `YYYY-MM-DD`)
  - Files:
    - `image`: file (Single, untuk gambar utama)
    - `images`: file[] (Multiple, untuk galeri gambar)
  - Error tambahan: `404` jika portfolio tidak ditemukan, `400` jika file bukan gambar atau mengandung konten berbahaya.
- `DELETE /api/admin/portfolios/:id`
  - Error tambahan: `404` jika portfolio tidak ditemukan.
- `PUT /api/admin/users/:id`
  - Error tambahan: `404` jika user tidak ditemukan.
- `DELETE /api/admin/users/:id`
  - Error tambahan: `404` jika user tidak ditemukan.

## 3. Kontrak Data Frontend

Gunakan tipe data yang konsisten:

```ts
export type PortfolioImage = {
  id: number;
  portfolio_id: number;
  image_url: string;
  created_at: string;
  updated_at: string;
};

export type Portfolio = {
  id: number;
  title: string;
  description: string;
  role: string;
  type: string;
  industry: string;
  tech_stack: string;
  project_link: string;
  image_url: string; // URL gambar utama
  images: PortfolioImage[]; // Daftar gambar galeri
  start_date: string; // ISO datetime
  end_date: string;   // ISO datetime
  is_published: boolean;
  created_at: string;
  updated_at: string;
};

export type PaginatedPortfolios = {
  data: Portfolio[];
  page: number;
  limit: number;
  total: number;
};
```

Catatan:
- Kirim tanggal ke backend dalam format ISO.
- Saat menampilkan tanggal ke user, format di sisi frontend (mis. `MMM YYYY`).

## 4. Struktur Halaman yang Disarankan

- `/` (Public Portfolio List)
  - Search input.
  - Filter `industry` dan `type`.
  - Pagination.
  - Empty state + loading skeleton.

- `/admin/login`
  - Form username/password.
  - Simpan token jika login sukses.
  - Redirect ke dashboard admin.

- `/admin/portfolios`
  - Tabel/list portfolio.
  - Tombol tambah, edit, hapus.
  - Konfirmasi sebelum hapus.

- `/admin/portfolios/new`
- `/admin/portfolios/:id/edit`
  - Form lengkap field portfolio.
  - Validasi client-side minimal.

- `/admin/users/:id/edit` (opsional tahap 2)
  - Form update username/password.
  - Password boleh kosong jika tidak mau diubah.

## 5. Arsitektur Frontend yang Disarankan

- Framework: React + TypeScript (boleh Next.js atau Vite).
- HTTP client: `fetch` wrapper atau `axios` dengan interceptor.
- State server: TanStack Query/SWR untuk caching, retry, invalidation.
- Form: React Hook Form + schema validator (Zod/Yup).
- UI: komponen reusable (button, input, modal, table, pagination).

Struktur folder minimal:

- `src/api/` untuk fungsi endpoint.
- `src/types/` untuk model TypeScript.
- `src/features/public-portfolios/`
- `src/features/admin-auth/`
- `src/features/admin-portfolios/`
- `src/components/ui/`
- `src/lib/` (token storage, config, helpers).

## 6. API Client Plan

Buat `apiClient` dengan behavior berikut:

- Base URL dari env (`VITE_API_BASE_URL` atau `NEXT_PUBLIC_API_BASE_URL`).
- Otomatis menambahkan header `Authorization` jika token tersedia.
- Parsing error response backend ke format seragam:
  - `400`: validasi input/query.
  - `401`: token invalid/expired -> logout + redirect login.
  - `404`: data tidak ditemukan (update/delete).
  - `429`: tampilkan pesan rate limit.
  - `500`: fallback pesan umum.

Contoh fungsi API:

- `login(payload)`
- `getPortfolios(params)`
- `createPortfolio(payload, token)`
- `updatePortfolio(id, payload, token)`
- `deletePortfolio(id, token)`
- `updateUser(id, payload, token)`
- `deleteUser(id, token)`

## 7. Validasi Form yang Wajib

- Login:
  - username required
  - password required

- Portfolio form:
  - `title`, `description`, `role`, `type`, `industry` required.
  - `project_link` valid URL (jika diisi).
  - `image` (Main Image) wajib saat Create (opsional saat Update).
  - `images` (Gallery) opsional, dukung upload banyak file sekaligus.
  - Tipe file wajib gambar (`jpg`, `jpeg`, `png`, `webp`).
  - `end_date >= start_date`.

- Query list:
  - pastikan `page` dan `limit` selalu angka positif sebelum request.

## 8. Auth & Security Frontend

- Simpan token di memory + fallback `localStorage` (praktis) atau sessionStorage.
- Jangan render route admin tanpa token valid.
- Tambahkan logout manual dan auto-logout saat `401`.
- Jangan log token ke console.

## 9. UX States yang Harus Ada

Setiap halaman data harus punya:

- loading state,
- empty state,
- error state,
- success feedback (toast/snackbar) untuk create/update/delete.

## 10. Tahapan Implementasi (Roadmap)

### Phase 1 - Fondasi

- Setup project frontend (TS, router, query client, UI base).
- Buat `.env` dan konfigurasi `API_BASE_URL`.
- Implement API client + handler error global.

### Phase 2 - Public Portfolio

- Implement halaman list portfolio.
- Integrasikan search/filter/pagination ke endpoint `GET /api/portfolios`.
- Tambahkan debouncing untuk search.

### Phase 3 - Admin Auth

- Implement halaman login.
- Simpan token, proteksi route admin, auto-redirect.
- Implement logout.

### Phase 4 - Admin Portfolio CRUD

- Implement list admin portfolio.
- Implement create/update/delete.
- Handle error `404` dan `401` sesuai backend.

### Phase 5 - Admin User Management

- Implement edit user dan delete user.
- Handle update password opsional.

### Phase 6 - Hardening & Quality

- Tambah unit/integration test untuk API layer dan form.
- Tambah e2e test skenario utama (login, CRUD, filter).
- Audit aksesibilitas dasar (label, keyboard, focus state).

## 11. Checklist Implementasi untuk Programmer/AI

- [ ] Base URL env berjalan untuk dev/staging/prod.
- [ ] Semua endpoint backend terintegrasi.
- [ ] Error status `400/401/404/429/500` memiliki UI message yang tepat.
- [ ] Query `page/limit` selalu valid (>0).
- [ ] Route admin diproteksi.
- [ ] Token terhapus saat logout dan saat `401`.
- [ ] CRUD portfolio lengkap + optimistic/refresh invalidate.
- [ ] Form validation aktif dan user-friendly.
- [ ] README frontend berisi cara setup, env, dan run.

## 12. Catatan Backend yang Sudah Diperbaiki (Dampak ke Frontend)

- Middleware auth kini memvalidasi format `Bearer` lebih ketat.
- Update/delete user dan portfolio sekarang mengembalikan `404` jika ID tidak ada.
- Query `page` dan `limit` pada endpoint public sekarang mengembalikan `400` jika invalid.
- Dokumentasi endpoint user admin diseragamkan ke `/api/admin/users/:id`.
- **Dukungan Galeri Gambar**: CRUD Portofolio kini menggunakan `multipart/form-data` untuk mendukung upload gambar (field `image` dan `images`).
- **Keamanan File**: Backend sekarang memvalidasi tipe MIME dan ekstensi file gambar, serta memberikan nama unik (UUID) untuk mencegah konflik.
- **Static Files**: Semua file yang di-upload dapat diakses melalui prefix URL `/uploads/`.

Frontend wajib menyesuaikan handling error dan metode pengiriman data (Form Data) terhadap perubahan ini.
