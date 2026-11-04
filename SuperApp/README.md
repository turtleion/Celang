# SuperApp (Celang Showcase Application)

SuperApp adalah aplikasi web konsol aman dan terdistribusi yang mendemonstrasikan **seluruh fitur bahasa Celang** dan **pustaka standarnya (stdlib)** dalam satu arsitektur nyata:
- **Server HTTP REST** dengan middleware keamanan (CSP, X-Frame-Options, nosniff, dsb.)
- **Konsol Berbasis Web (Allowlist)** untuk administrasi sistem aman (menolak perintah berbahaya OS secara ketat)
- **Kriptografi Mutakhir**: AES-256-GCM, PBKDF2 (PHC format), Ed25519 digital signature, HMAC-SHA256 bearer tokens, dan UUID v4
- **Penyimpanan Terenkripsi & Transaksional**: Catatan pengguna terenkripsi otomatis dengan jaminan ACID (`transaksi` & `pulihkan`)
- **Pencadangan Arsip ZIP**: Backup otomatis seluruh direktori data ke dalam format `.zip`
- **Konkurensi Modern**: Saluran komunikasi CSP (`Saluran`), pemrosesan latar belakang goroutine (`sambil`), dan evaluasi paralel (`titip` & `jemput`)
- **Fitur Bahasa Lengkap**: Generik (`Bentuk<T>`), Destrukturisasi, Operator Pipeline (`|>`), String Interpolation, Rujukan (`rujuk`), Variabel Reaktif (`reaktif` & `pantau`), Time-Travel History (`.histori` & `.sebelumnya`), dan FFI bridge.

---

## 📁 Struktur Berkas

```
SuperApp/
├── celang.mod            # Deklarasi modul Celang
├── utama.ce              # Titik masuk utama server (entrypoint)
├── lib/
│   ├── keamanan.ce       # Modul kriptografi, token sesi, dan tanda tangan digital
│   ├── konsol.ce         # Runner perintah konsol aman & demonstrasi fitur bahasa
│   ├── pekerja.ce        # Antrean CSP (Saluran) dan pekerja latar belakang (sambil)
│   ├── penyimpanan.ce    # Penyimpanan JSON, catatan terenkripsi, transaksi, & arsip ZIP
│   └── rute.ce           # Pengarah rute HTTP REST, middleware, & handler statis
├── web/
│   ├── index.html        # Antarmuka web konsol
│   ├── app.js            # Logika klien web (autentikasi & eksekusi API)
│   └── app.css           # Tema visual dark modern konsol
└── uji/
    └── superapp_uji.ce   # 13 pengujian menyeluruh (unit & integrasi)
```

---

## 🚀 Cara Menjalankan

### 1. Menjalankan Server Utama
Jalankan perintah berikut di direktori utama repositori Celang:

```bash
./celang run SuperApp/utama.ce
```

Saat pertama kali dijalankan, SuperApp akan:
1. Menyiapkan direktori data (`SuperApp/data` atau sesuai `SUPERAPP_DATA`).
2. Membuat akun `admin` otomatis dengan kata sandi acak yang dicetak di terminal.
3. Menjalankan server HTTP di `http://127.0.0.1:8080`.

Buka browser Anda dan akses `http://127.0.0.1:8080`, lalu masuk menggunakan kredensial admin tersebut.

### 2. Opsi Variabel Lingkungan (Environment Variables)
Anda dapat mengonfigurasi SuperApp menggunakan variabel lingkungan:
- `SUPERAPP_PORT`: Port HTTP (default: `8080`)
- `SUPERAPP_ADMIN_PASSWORD`: Kata sandi tetap untuk akun admin
- `SUPERAPP_SECRET`: Rahasia HMAC untuk penandatanganan token sesi
- `SUPERAPP_KEY`: Kunci heksadesimal AES-256 untuk enkripsi data catatan
- `SUPERAPP_DATA`: Jalur direktori penyimpanan data (default: `./SuperApp/data`)
- `SUPERAPP_BACKUP`: Jalur direktori cadangan arsip ZIP (default: `./SuperApp/backup`)

Contoh:
```bash
SUPERAPP_PORT=9000 SUPERAPP_ADMIN_PASSWORD="SandiRahasia2026!" ./celang run SuperApp/utama.ce
```

---

## 🧪 Menjalankan Pengujian (Test Suite)

Jalankan pengujian otomatis untuk memvalidasi seluruh fungsi:

```bash
./celang uji SuperApp/uji/superapp_uji.ce
```

Seluruh 13 pengujian mencakup:
1. `Kripto: hash sandi dan verifikasi` (PBKDF2 timing-safe)
2. `Kripto: AES-256-GCM round-trip dan deteksi manipulasi` (AEAD authentication)
3. `Kripto: token sesi bertanda tangan HMAC`
4. `Kripto: tanda tangan Ed25519`
5. `Penyimpanan: simpan dan muat pengguna`
6. `Penyimpanan: catatan terenkripsi`
7. `Konsol: allowlist menolak perintah OS` (keamanan webshell aman)
8. `Konsol: demonstrasi fitur bahasa` (generik, destrukturisasi, pipeline, rujukan, dsb.)
9. `Arsip: cadangan ZIP data aplikasi`
10. `Pekerja: antrean CSP memproses hash` (`Saluran` + `sambil`)
11. `HTTP: kesehatan, login, konsol, catatan, dan 401`
12. `FFI: pointer NULL aman`
13. `Stdlib: matematika, teks, JSON, waktu`

---

## 🛡️ Keamanan & Allowlist Konsol

SuperApp menerapkan filosofi **Safe-by-Default**:
- Bukan webshell sembarangan yang membiarkan sembarang `exec/sh/bash` berjalan.
- Menggunakan **Allowlist Perintah**: `bantuan`, `status`, `waktu`, `echo`, `hash`, `enkripsi`, `dekripsi`, `statistik`, `siapa`, `cadangkan`, `catatan`, `fitur`, `hitung`.
- Perintah eksternal OS berbahaya seperti `rm`, `cat /etc/passwd`, `shutdown`, dsb. ditolak secara otomatis dan dicatat ke log audit.
