# EventHub Backend

[![Go](https://img.shields.io/badge/Go-1.26.5-00ADD8?logo=go)](https://go.dev/)
[![Gin](https://img.shields.io/badge/Gin_Gonic-1.12.0-008ECF?logo=gin&logoColor=white)](https://gin-gonic.com/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16.15-4169E1?logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![Redis](https://img.shields.io/badge/Redis-7.0.15-DC382D?logo=redis)](https://redis.io/)

EventHub Backend adalah REST API untuk mendukung platform penemuan dan pengelolaan acara serta komunitas. API ini menyediakan satu layanan bagi peserta, penyelenggara, dan admin untuk mengelola akun, menemukan acara dan komunitas, berpartisipasi, serta mengakses dashboard sesuai peran.

## Features

- Registrasi, login, logout, permintaan reset kata sandi, dan reset kata sandi.
- Autentikasi menggunakan JWT Bearer Token; Redis digunakan untuk memeriksa blacklist token saat logout.
- Pencarian dan filter acara berdasarkan kata kunci, kategori, dan lokasi; detail acara dan daftar acara mendatang.
- Fitur bergabung/keluar dari acara, menyimpan acara, serta melihat acara milik pengguna.
- Pembuatan dan pengeditan acara untuk akun berperan penyelenggara.
- Penelusuran komunitas, detail, anggota, komunitas populer, serta fitur bergabung/keluar.
- Profil pengguna, pembaruan profil dan avatar, informasi header, serta penggantian kata sandi.
- Notifikasi pengguna dan testimoni.
- Dashboard khusus penyelenggara dan admin dengan pembatasan akses berdasarkan peran.
- Dokumentasi API interaktif melalui Swagger.

## Usage Instruction

1. **Clone repository :**

   ```bash
   git clone https://github.com/MGT06/EventHub_Backend.git
   cd EventHub_Backend
   ```

2. **Siapkan PostgreSQL dan Redis.** 

3. **Buat file `.env`di root proyek.**

   ```dotenv
   DBUSER={YOUR_DB_USER}
   DBPASS={YOUR_DB_PASS}
   DBHOST={YOUR_DB_HOST}
   DBPORT={YOUR_DB_PORT}
   DBNAME={YOUR_DB_NAME}

   RDBUSER={YOUR_RDB_USER}
   RDBPASS={YOUR_RDB_PASS}
   RDBHOST={YOUR_RDB_HOST}
   RDBPORT={YOUR_RDB_PORT}

   JWT_KEY={YOUR_SECRET_KEY}
   JWT_ISSUER={YOUR_ISSUER}
   ```

   Sesuaikan seluruh nilai dengan lingkungan masing-masing.

4. **Jalankan migrasi database.** Dari root proyek, gunakan target yang tersedia :

   ```bash
   make migrate-up
   ```

5. **Unduh dependensi dan jalankan server:**

   ```bash
   go mod download
   go run .
   ```

6. **Jalankan Backend :**

   ```bash
   go run main.go ## dapat juga menggunakan 'air'
   ```

## Endpoint API Utama

### Autentikasi

| Metode | Endpoint | Keterangan |
| --- | --- | --- |
| `POST` | `/auth/register` | Registrasi akun |
| `POST` | `/auth/login` | Login |
| `POST` | `/auth/logout` | Logout (**Bearer**) |
| `POST` | `/auth/forgot-password` | Meminta proses reset kata sandi |
| `POST` | `/auth/reset-password` | Mengatur ulang kata sandi |

### Event

| Metode | Endpoint | Keterangan |
| --- | --- | --- |
| `GET` | `/event` | List event; Dapat melakukan search dan filter juga jika ditambahkan query param |
| `GET` | `/event/detail/:id` | Detail event berdasarkan ID. |
| `GET` | `/event/upcoming` | List event mendatang. |
| `GET` | `/event/myevent` | event yang user ikuti (**Bearer**). |
| `POST` | `/event/join/:id` | Join atau Leave event (**Bearer**). |
| `POST` | `/event/save/:id` | Save atau UnSave event (**Bearer**). |
| `POST` | `/event/create` | Buat event (**Bearer**, peran `organizer`). |
| `PATCH` | `/event/:id/edit` | Edit event (**Bearer**, peran `organizer`). |

### Community

| Metode | Endpoint | Keterangan |
| --- | --- | --- |
| `GET` | `/community` | List community dengan pencarian/filter |
| `GET` | `/community/detail/:id` | Detail community |
| `GET` | `/community/members/:id` | List member community |
| `GET` | `/community/popular` | List 4 community yang populer |
| `POST` | `/community/join/:id` | Join atau Leave community (**Bearer**) |

### User, notifikasi, dashboard, dan testimoni

| Metode | Endpoint | Keterangan |
| --- | --- | --- |
| `GET` | `/user` | Profil user (**Bearer**). |
| `PATCH` | `/user/edit` | Edit user (**Bearer**) |
| `POST` | `/user/change-password` | Ganti kata sandi (**Bearer**). |
| `GET` | `/user/headerinfo` | Informasi ringkas header user (**Bearer**). |
| `GET` | `/notification` | Notifikasi user (**Bearer**). |
| `GET` | `/organizer` | Dashboard organizer (**Bearer**, peran `organizer`). |
| `GET` | `/admin` | Dashboard admin (**Bearer**, peran `admin`). |
| `GET` | `/testimony` | List testimoni. |
| `POST` | `/testimony/send-testimony` | Kirim testimoni (**Bearer**). |

Dokumentasi Swagger interaktif tersedia di `/api/docs/index.html`. Periksa spesifikasi tersebut untuk detail skema request dan response.

## License
This project is licensed under the MIT License

## Related Project

[Frontend](https://github.com/MGT06/EventHub)