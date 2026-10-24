---
outline: deep
---

# Pustaka Sistem Operasi (`sistem` / `os`)

Modul **`sistem`** (dengan alias internasional **`os`**) menyediakan akses lengkap ke lingkungan sistem operasi, interaksi terminal, eksekusi proses eksternal, variabel lingkungan, inspeksi memori, hingga informasi perangkat keras komputer.

Modul ini sangat berguna saat Anda membangun:
- Skrip otomasi (*automation scripts*) & pembuat *build tool* (seperti `make` atau *task runner*).
- Utilitas baris perintah (*CLI tools*) yang menerima argumen terminal.
- Aplikasi *backend* & layanan peladen (*server*) yang membutuhkan konfigurasi lingkungan (*environment variables*).
- Pemantau kesehatan proses & konsumsi memori (*system health monitors*).

---

## Cara Memanggil

Modul ini langsung aktif secara global maupun dapat diimpor secara spesifik:

```celang
// 1. Langsung panggil via objek global:
buat info = sistem.jalankan("git status")
buat home = os.homedir()

// 2. Atau impor fungsi tertentu (tanpa tanda kutip):
ambil { jalankan, dapatkan_env, direktori_kerja } dari sistem

// 3. Atau impor modul dengan alias:
impor sistem sebagai os
```

---

## Ringkasan Fungsi & Padanan Istilah

| Fungsi Bahasa Indonesia | Alias Bahasa Inggris | Keterangan |
| :--- | :--- | :--- |
| `sistem.jalankan(cmd, [opsi])` | `os.exec(cmd, [options])` | Menjalankan perintah terminal / proses eksternal (**Fitur Unggulan**) |
| `sistem.cari_program(nama)` | `os.which(nama)` / `os.look_path` | Menemukan lokasi biner di `PATH` sistem operasi |
| `sistem.dapatkan_env(kunci, [bawaan])` | `os.getenv(kunci, [bawaan])` | Membaca variabel lingkungan sistem |
| `sistem.setel_env(kunci, nilai)` | `os.setenv(kunci, nilai)` | Menyetel nilai variabel lingkungan |
| `sistem.hapus_env(kunci)` | `os.unsetenv(kunci)` | Menghapus variabel lingkungan |
| `sistem.ada_env(kunci)` | `os.has_env(kunci)` | Memeriksa apakah variabel lingkungan aktif |
| `sistem.semua_env()` | `os.environ()` / `os.env()` | Mengambil seluruh variabel lingkungan dalam bentuk Map |
| `sistem.argumen([indeks, [bawaan]])` | `os.args([index, [default]])` | Mengambil argumen CLI skrip (semua atau per indeks) |
| `sistem.argumen_mentah()` | `os.raw_args()` | Mengambil seluruh argumen terminal termasuk biner pelaksana |
| `sistem.direktori_kerja()` | `os.cwd()` | Mendapatkan lokasi direktori kerja saat ini |
| `sistem.ubah_direktori(jalur)` | `os.chdir(path)` | Berpindah ke direktori tertentu |
| `sistem.os()` | `os.platform()` | Nama sistem operasi (`"linux"`, `"darwin"`, `"windows"`) |
| `sistem.arsitektur()` | `os.arch()` | Arsitektur CPU (`"amd64"`, `"arm64"`, `"386"`) |
| `sistem.jumlah_cpu()` | `os.cpu_count()` / `os.cpus()` | Jumlah inti prosesor (CPU cores) |
| `sistem.nama_host()` | `os.hostname()` | Nama komputer / nama *host* |
| `sistem.folder_home()` | `os.homedir()` | Jalur absolut direktori beranda pengguna saat ini |
| `sistem.folder_sementara()` | `os.tmpdir()` | Jalur direktori penyimpanan berkas sementara (*temp*) |
| `sistem.pid()` | `os.pid()` | Nomor identitas proses (*Process ID*) Celang |
| `sistem.waktu_aktif([satuan])` | `os.uptime([unit])` | Lama waktu proses berjalan dalam detik atau milidetik |
| `sistem.apakah_admin()` | `os.is_admin()` / `os.is_root()` | Mengecek apakah skrip berjalan dengan hak akses root/administrator |
| `sistem.memori()` | `os.memory()` | Statistik penggunaan memori RAM aplikasi |
| `sistem.gc()` | `os.gc()` | Memicu pembersihan memori (*Garbage Collector*) secara manual |
| `sistem.keluar([kode])` | `os.exit([kode])` | Menghentikan eksekusi skrip dengan kode keluar (*exit code*) |

---

## 🌟 Eksekusi Terminal & Proses Eksternal

Fitur utama modul `sistem` adalah kemampuan menjalankan perintah terminal secara fleksibel, aman, dan mudah dianalisis.

### 1. Menjalankan Perintah Teks (*Shell Command*)
Perintah teks otomatis dieksekusi melalui *shell* bawaan sistem (`sh -c` pada Linux/macOS, `cmd.exe /C` pada Windows):

```celang
buat hasil = sistem.jalankan("git status")

jika hasil.sukses {
    cetak("Keluaran git:\n" + hasil.keluaran)
} selain {
    cetak("Perintah gagal dengan kode:", hasil.kode_keluar)
    cetak("Galat:", hasil.galat)
}
```

Struktur Map yang dikembalikan oleh `sistem.jalankan` mendukung properti dwibahasa:
- `hasil.keluaran` / `hasil.stdout` (*Teks*): Teks keluaran standar.
- `hasil.galat` / `hasil.stderr` (*Teks*): Pesan kesalahan standar.
- `hasil.kode_keluar` / `hasil.exit_code` (*Angka*): Kode keluar proses (0 jika berhasil).
- `hasil.sukses` / `hasil.success` (*Kebenaran*): `benar` jika `kode_keluar == 0`.
- `hasil.durasi` / `hasil.duration` (*Angka*): Waktu eksekusi dalam milidetik.

### 2. Menjalankan Daftar Argumen Langsung (*Direct Execution*)
Untuk menghindari masalah pengutipan (*escaping*) spasi pada argumen atau serangan *shell injection*, Anda dapat meneruskan daftar teks:

```celang
buat hasil = os.exec(["git", "commit", "-m", "fitur: sistem modern"])
cetak("Waktu eksekusi:", hasil.duration, "ms")
```

### 3. Mengatur Batas Waktu (*Timeout*), Direktori Kerja, & Input
Anda dapat menentukan opsi tambahan melalui parameter kedua:

```celang
buat opsi = {
    "dir": "/home/pengguna/proyek",        // Jalankan di folder tertentu
    "timeout": 5000,                       // Hentikan jika lebih dari 5 detik (5000 ms)
    "env": { "NODE_ENV": "production" },   // Tambahan variabel lingkungan
    "stdin": "data masukan baris 1\n",     // Alirkan teks ke input proses
    "gabung_keluar": benar                 // Gabungkan stderr ke stdout
}

buat hasil = sistem.jalankan("npm test", opsi)
jika !hasil.sukses {
    cetak("Pengujian gagal atau melebihi batas waktu!")
}
```

### 4. Mencari Lokasi Program (*Which / Look Path*)
Mengecek apakah suatu perkakas terinstal di komputer sebelum mengeksekusinya:

```celang
buat jalurGit = sistem.cari_program("git")

jika jalurGit != kosong {
    cetak("Git ditemukan di:", jalurGit)
} selain {
    cetak("Harap instal Git terlebih dahulu!")
}
```

---

## 🔐 Variabel Lingkungan (*Environment Variables*)

Celang menyediakan operasi lengkap dan aman untuk membaca konfigurasi aplikasi dari variabel lingkungan:

```celang
// 1. Membaca variabel dengan nilai cadangan (default) jika belum disetel
buat port = sistem.dapatkan_env("PORT", "3000")
buat mode = os.getenv("APP_MODE", "development")

// 2. Mengecek keberadaan variabel lingkungan
jika sistem.ada_env("DATABASE_URL") {
    cetak("Koneksi basis data terkonfigurasi.")
}

// 3. Menyetel dan menghapus variabel lingkungan saat runtime
sistem.setel_env("KUNCI_SESI", "abc123secret")
sistem.hapus_env("KUNCI_SESI")

// 4. Membaca seluruh variabel lingkungan komputer
buat semua = sistem.semua_env()
cetak("Jalur PATH komputer:", semua["PATH"])
```

---

## 🎯 Argumen Baris Perintah (*CLI Arguments*)

Saat menjalankan skrip melalui terminal (`celang skrip.ce arg1 arg2`), argumen dapat diakses dengan sangat ringkas:

```bash
celang aplikasi.ce buat-laporan --format=pdf
```

```celang
// Ambil seluruh argumen skrip sebagai daftar:
buat args = sistem.argumen()
cetak("Jumlah argumen:", panjang(args)) // 2

// Ambil berdasarkan indeks secara aman:
buat perintah = sistem.argumen(0) // "buat-laporan"
buat format = sistem.argumen(1, "--format=json") // Nilai bawaan jika argumen ke-2 tidak ada

// Akses biner pelaksana mentah (termasuk "celang" dan path skrip):
buat mentah = os.raw_args()
```

---

## 📂 Direktori Kerja & Navigasi

```celang
// 1. Dapatkan folder kerja saat ini:
buat lokasiSekarang = sistem.direktori_kerja()
cetak("Direktori aktif:", lokasiSekarang)

// 2. Berpindah folder:
sistem.ubah_direktori("/tmp")
cetak("Sekarang berada di:", os.cwd())

// Kembali ke folder awal:
sistem.ubah_direktori(lokasiSekarang)
```

---

## 💻 Informasi Perangkat Keras & Host

```celang
cetak("Sistem Operasi :", sistem.os())               // "linux", "darwin", "windows"
cetak("Arsitektur CPU :", sistem.arsitektur())       // "amd64", "arm64"
cetak("Jumlah Inti CPU:", sistem.jumlah_cpu())       // misal: 8 atau 16
cetak("Nama Komputer  :", sistem.nama_host())        // misal: "laptop-dev"
cetak("Folder Home    :", sistem.folder_home())      // misal: "/home/budi"
cetak("Folder Temp    :", sistem.folder_sementara()) // misal: "/tmp"
cetak("ID Proses (PID):", sistem.pid())              // misal: 14230
cetak("Waktu Aktif    :", sistem.waktu_aktif(), "detik")
cetak("Hak Akses Root :", sistem.apakah_admin())     // benar / salah
```

---

## 📊 Pemantauan Memori & Garbage Collector

Pantau konsumsi RAM proses Celang Anda saat memproses tugas besar:

```celang
// Dapatkan statistik memori dalam satuan bytes
buat infoMem = sistem.memori()

cetak("Memori aktif teralokasi:", infoMem.alokasi / 1024, "KB")
cetak("Total alokasi kumulatif:", infoMem.total_alokasi / 1024, "KB")
cetak("Memori dari OS         :", infoMem.sistem / 1024, "KB")
cetak("Jumlah siklus GC       :", infoMem.jumlah_gc)

// Memicu siklus Garbage Collector secara langsung jika diperlukan
sistem.gc()
```

---

## 🚪 Siklus Hidup Proses (*Process Lifecycle*)

Gunakan `sistem.keluar` untuk menghentikan program Celang seketika dengan kode status tertentu:

```celang
jika !berkas.ada("./kunci-rahasia.txt") {
    cetak("Galat kritis: berkas kunci tidak ditemukan!")
    sistem.keluar(1) // Keluar dengan status gagal (exit code 1)
}

// Keluar sukses (exit code 0):
sistem.keluar(0)
```

---

## 📌 Konstanta Sistem

Modul `sistem` juga menyediakan konstanta penanda jalur dan pembatas lintas platform:

```celang
cetak("Pemisah berkas :", sistem.PEMISAH_BERKAS) // "/" di Linux/macOS, "\" di Windows
cetak("Pemisah jalur  :", sistem.PEMISAH_JALUR)  // ":" di Linux/macOS, ";" di Windows
cetak("Dev Null       :", sistem.DEV_NULL)       // "/dev/null" atau "NUL"
cetak("Karakter EOL   :", sistem.GARIS_BARU)     // "\n" atau "\r\n"
```
