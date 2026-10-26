# 🗺️ Rencana Prasyarat & Fondasi Sebelum Package Manager Celang

Dokumen ini memetakan rencana kerja strategis untuk melengkapi fitur-fitur pustaka standar, ketahanan runtime, dan ergonomi bahasa Celang sebelum memasuki subsistem **Package Manager** (`celang pasang` / `celang get`).

---

## 🧭 Mengapa Fitur-Fitur Ini Diperlukan Lebih Dulu?

Sebelum Celang dapat mengunduh, mengekstrak, dan memverifikasi modul pihak ketiga dari Git/Registry secara andal, Celang membutuhkan kapabilitas bawaan (*pure Go*, tanpa dependensi shell eksternal seperti `tar`, `unzip`, `curl`):
1. **Portabilitas Lintas Platform**: Mencegah kegagalan perintah terminal di Windows yang tidak memiliki `tar` atau `curl` bawaan.
2. **Efisiensi Memori**: Mengalirkan unduhan paket tarball/zip besar langsung ke disk (*streaming I/O*) tanpa membebani memori RAM.
3. **Keamanan Rantai Pasok (*Supply Chain Security*)**: Verifikasi checksum SHA-256 pada berkas pengunci (`celang.kunci`) untuk memastikan pustaka tidak dirusak di tengah jalan.
4. **Ergonomi Pengembang**: Memberikan kemampuan pemrosesan koleksi dan manipulasi bita yang nyaman bagi pengembang pembuat pustaka di ekosistem Celang.

---

## 📊 Peta Tahapan Pengembangan (Roadmap)

```
[ Fase A: Ergonomi Koleksi & Bita ] ➔ [ Fase B: Modul Standar Kripto, Kompresi, HTTP ] ➔ [ Fase C: SemVer & Parser celang.mod ] ➔ [ Fase D: Package Manager CLI ]
```

---

## 🚀 Fase A: Ekstensi Koleksi Array & Decoder Bita

Fase ini berfokus pada kenyamanan ekspresi data fungsional dan pengolahan data bita mentah:

### 1. Chaining Metode Fungsional pada `Array`
Menghadirkan metode fungsional langsung pada objek array yang dapat dirangkai (*method chaining*):
* `.petakan(fungsi)` / `.map(fn)`: Mentransformasi setiap elemen.
* `.saring(fungsi)` / `.filter(fn)`: Menyaring elemen berdasarkan predikat boolean.
* `.lipat(fungsi, nilai_awal)` / `.reduksi(...)` / `.reduce(fn, initial)`: Mereduksi array menjadi nilai tunggal akumulatif.
* `.temukan(fungsi)` / `.find(fn)`: Mengembalikan elemen pertama yang cocok dengan predikat.
* `.ada_yang(fungsi)` / `.some(fn)` & `.semuanya(fungsi)` / `.every(fn)`: Pengujian kondisi elemen.

### 2. Decoder & Konstruktor Lanjutan `Bita` (`Bytes`)
* `bita.dari_base64(teks)` / `bytes.from_base64(str)`: Mengubah string Base64 menjadi objek `Bita`.
* `bita.dari_heks(teks)` / `bytes.from_hex(str)`: Mengubah string heksadesimal menjadi objek `Bita`.
* Konstruktor fleksibel `bita(teks, "base64" | "hex")`.

---

## 🔒 Fase B: Pustaka Standar Fondasi (`kripto`, `kompresi`, `http.unduh`)

### 1. Modul `kompresi` (`compression`)
* `kompresi.ekstrak_tar_gz(arsip, target_folder)` / `compression.extract_tar_gz(...)`
* `kompresi.buat_tar_gz(sumber_folder, berkas_keluar)` / `compression.create_tar_gz(...)`
* `kompresi.ekstrak_zip(arsip, target_folder)` / `compression.extract_zip(...)`
* `kompresi.buat_zip(daftar_berkas, berkas_keluar)` / `compression.create_zip(...)`
* `kompresi.gzip(bita)` & `kompresi.gunzip(bita)`

### 2. Modul `kripto` (`crypto`)
* `kripto.sha256(data)` / `crypto.sha256(data)`
* `kripto.sha512(data)`
* `kripto.md5(data)`
* `kripto.hmac_sha256(kunci, data)`
* `kripto.acak_heks(panjang)`
* `kripto.uuid()` (Format UUID v4 standar)

### 3. Streaming `http.unduh` (`http.download`)
* Mengalirkan respon HTTP `GET` langsung ke berkas di disk secara atomik.
* Opsi batas waktu (`batas_waktu`), tajuk kustom (`tajuk`), dan penanganan pengalihan (*redirect*).
* Callback progres unduhan: `pada_progres: fungsi(terunduh, total)`.

---

## 📦 Fase C: SemVer & Parser Manifest (`celang.mod` & `celang.kunci`)

### 1. Pustaka Pembanding Versi Semantik (`semver`)
* Penguraian spesifikasi SemVer 2.0.0 (`mayor.minor.tambalan[-pratayang]`).
* Operator rentang dan batasan: `^1.2.0`, `~1.2.0`, `>=1.0.0, <2.0.0`.
* Fungsi pembanding dan pemilihan versi tertinggi yang kompatibel.

### 2. Spesifikasi Dokumen `celang.mod`
```celang
modul github.com/pengguna/projek-ku
versi 1.2.0
celang >= 1.0.0

butuh (
    github.com/celang/pustaka-web v0.3.0
    github.com/celang/kripto-jwt ^1.0.0
    github.com/organisasi/logger ~2.1.0
)

ganti (
    github.com/celang/pustaka-web => ../lokal_pustaka_web
)
```

### 3. Berkas Pengunci Integritas (`celang.kunci` / `celang.lock`)
Mencatat pohon dependensi terperinci, versi pasti (*resolved version*), dan hash SHA-256 tarball untuk verifikasi deterministik.

---

## 🛠️ Fase D: Package Manager CLI (`celang pasang` / `celang get`)

Implementasi perintah CLI terintegrasi:
* `celang pasang <url_repo>[@versi]` / `celang get ...`
* `celang pasang` (tanpa argumen: mengunduh seluruh dependensi dari `celang.mod`)
* `celang mutakhirkan [paket]` / `celang update [pkg]`
* `celang kemas` / `celang pack` (membuat paket `.tar.gz` siap distribusi)
