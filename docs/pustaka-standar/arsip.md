---
outline: deep
---

# Pustaka Arsip & Kompresi (`arsip` / `archive`)

Modul **`arsip`** (dengan alias **`archive`**, **`kompresi`**, **`compression`**, **`zip`**, dan **`tar`**) menyediakan kapabilitas pengemasan, ekstraksi, inspeksi, dan kompresi data berstandar industri pada bahasa pemrograman Celang. 

Diimplementasikan murni menggunakan pustaka standar Go (*pure Go*, tanpa dependensi shell eksternal seperti `tar`, `unzip`, atau `curl`), modul ini menjamin portabilitas 100% di berbagai sistem operasi (Linux, macOS, dan Windows) serta dirancang secara **aman secara bawaan (*safe-by-default*)** dengan perlindungan bawaan terhadap kerentanan **Zip Slip** dan **Tar Slip**.

---

## 🎯 Fitur Utama

1. **Format ZIP (`buat_zip`, `ekstrak_zip`, `daftar_zip`, `baca_berkas_zip`)**:
   - Mengemas folder rekursif atau daftar berkas menjadi file `.zip`.
   - Mengekstrak isi ZIP ke direktori tujuan dengan proteksi keamanan *path traversal*.
   - Menginspeksi metadata berkas (nama, ukuran terkompresi/asli, tanggal modifikasi) tanpa mengekstrak.
   - Membaca satu berkas langsung dari dalam ZIP sebagai objek `Bita`.
2. **Format TAR.GZ (`buat_tar_gz`, `ekstrak_tar_gz`, `daftar_tar_gz`, `baca_berkas_tar_gz`)**:
   - Standar pengemasan tarball yang digunakan untuk distribusi paket pada **Package Manager Celang**.
   - Preservasi izin berkas (*file permissions* `0755` untuk skrip/biner vs `0644` untuk teks).
3. **Format TAR Polos (`buat_tar`, `ekstrak_tar`, `daftar_tar`)**:
   - Pengarsipan tanpa kompresi untuk pemrosesan berkecepatan tinggi.
4. **Kompresi Bita Dalam-Memori (*In-Memory Compression*)**:
   - **Gzip**: `gzip(data)` dan `gunzip(bita)`.
   - **Zlib / Deflate**: `zlib_kompres(data)` dan `zlib_dekompres(bita)`.
5. **Pertahanan Keamanan Ketat**:
   - **Pencegahan Zip Slip / Tar Slip**: Otomatis menggagalkan ekstraksi berkas berbahaya yang mencoba menulis di luar direktori tujuan (`../../etc/passwd` atau `..\..\Windows`).
   - Penanganan berkas terpotong (*truncated*) atau rusak secara elegan tanpa *crash*.

---

## 🚀 Cara Menggunakan

Modul `arsip` dapat langsung diakses melalui objek global atau diimpor melalui sistem modul:

```celang
// 1. Akses langsung melalui objek global:
buat ok = arsip.buat_zip("./projek", "projek.zip")
arsip.ekstrak_zip("projek.zip", "./tujuan")

// 2. Menggunakan alias kompresi / archive:
buat terkompresi = kompresi.gzip("Teks panjang yang ingin dihemat...")
buat teksAsli = kompresi.gunzip(terkompresi).ke_teks()

// 3. Impor spesifik melalui sistem modul:
ambil { buat_tar_gz, ekstrak_tar_gz } dari "arsip"
```

### 🌍 Dukungan Mode Universal (`"use universal"`)

Pada berkas yang mengaktifkan direktif `"use universal"`, gunakan nama modul bahasa Inggris `archive` atau `compression`:

```celang
"use universal"
import { create_zip, extract_zip, gzip, gunzip } from "archive"

const ok = create_zip("./src", "bundle.zip")
const total = extract_zip("bundle.zip", "./dest")
```

---

## 📚 Referensi API Lengkap

### 1. Arsip ZIP

#### `arsip.buat_zip(sumber, tujuanZip, [opsi])`
Mengemas berkas atau folder menjadi file ZIP.
- **Parameter**:
  - `sumber`: `Teks` (jalur folder atau berkas) atau `Daftar` teks (kumpulan jalur berkas).
  - `tujuanZip`: `Teks` (jalur berkas ZIP keluaran).
  - `opsi` (opsional): `Peta` opsi:
    - `"metode"` / `"method"`: `"deflate"` (bawaan, kompresi standar) atau `"store"` / `"simpan"` (tanpa kompresi).
- **Kembalian**: `Logika` (`benar` jika sukses).

```celang
// Mengemas seluruh isi folder secara rekursif:
arsip.buat_zip("./projek", "projek.zip")

// Mengemas tanpa kompresi (store mode):
arsip.buat_zip("./video", "video.zip", { "metode": "store" })

// Mengemas daftar berkas tertentu:
arsip.buat_zip(["./berkas1.txt", "./berkas2.png"], "aset.zip")
```

#### `arsip.buat_zip_memori(petaKonten, tujuanZip, [opsi])`
Mengemas kumpulan data string atau `Bita` langsung dari memori menjadi berkas ZIP di disk tanpa perlu membuat berkas sementara. Dilengkapi proteksi otomatis terhadap upaya *path traversal* (Zip Slip).
- **Parameter**:
  - `petaKonten`: `Peta` dengan pasangan `nama_berkas: data` (data dapat berupa `Teks` atau `Bita`).
  - `tujuanZip`: `Teks` (jalur berkas ZIP keluaran).
  - `opsi` (opsional): `Peta` opsi (`"metode"`: `"deflate"` atau `"store"`).
- **Kembalian**: `Logika` (`benar` jika sukses).

```celang
buat berkasMemori = {
    "catatan.txt": "Dibuat langsung dari memori Celang!",
    "aset/logo.png": berkas.baca_bita("./logo.png"),
    "data/konfig.json": "{\"versi\": 1}"
}
arsip.buat_zip_memori(berkasMemori, "paket_memori.zip")
```

#### `arsip.uji_zip(jalur)`
Memverifikasi integritas checksum CRC-32 dan header seluruh anggota di dalam arsip ZIP tanpa perlu mengekstrak berkas ke disk (mirip `ZipFile.testzip()` di Python).
- **Parameter**:
  - `jalur`: `Teks` (jalur berkas ZIP).
- **Kembalian**:
  - `kosong` (`null`): Jika seluruh berkas di dalam arsip sehat dan lolos verifikasi CRC-32.
  - `Teks`: Nama berkas pertama yang terdeteksi rusak atau mengalami korupsi data.

```celang
buat hasil = arsip.uji_zip("unduhan.zip")
jika hasil == kosong {
    cetak("✓ Seluruh arsip ZIP sehat dan lolos uji integritas CRC-32.")
} jika tidak {
    cetak("⚠ Terdeteksi korupsi data pada berkas:", hasil)
}
```

#### `arsip.info_berkas_zip(berkasZip, namaBerkas)`
Mengambil metadata detail satu entri berkas di dalam arsip ZIP tanpa mengekstraknya.
- **Parameter**:
  - `berkasZip`: `Teks` (jalur berkas ZIP).
  - `namaBerkas`: `Teks` (jalur berkas di dalam arsip).
- **Kembalian**: `Peta` metadata yang memuat:
  - `nama`: Nama entri berkas.
  - `ukuran`: Ukuran asli sebelum kompresi (*uncompressed size*).
  - `ukuran_terkompresi`: Ukuran setelah kompresi (*compressed size*).
  - `crc32`: Nilai checksum CRC-32 (`Angka`).
  - `metode_kompresi` / `compression_method`: Teks metode kompresi (`"deflate"` atau `"store"`).
  - `apakah_direktori`: `benar` jika entri folder.
  - `tipe`: `"berkas"` atau `"direktori"`.
  - `waktu_modifikasi`: String RFC 3339.

```celang
buat info = arsip.info_berkas_zip("paket.zip", "data.csv")
cetak("Nama:", info.nama)
cetak("Ukuran asli:", info.ukuran, "bita")
cetak("CRC32:", info.crc32)
cetak("Metode:", info.metode_kompresi)
```

#### `arsip.ekstrak_zip(berkasZip, folderTujuan, [opsi])`
Mengekstrak seluruh berkas dari arsip ZIP ke folder tujuan secara aman.
- **Parameter**:
  - `berkasZip`: `Teks` (jalur berkas ZIP).
  - `folderTujuan`: `Teks` (direktori tempat berkas diekstrak).
  - `opsi` (opsional): `Peta` opsi:
    - `"timpa"` / `"overwrite"`: `Logika` (bawaan `benar`).
    - `"saring"` / `"filter"`: `Fungsi` penyaring kustom `fungsi(info) -> Logika`.
    - `"batas_ukuran"` / `"max_size"`: `Angka` (batas total kuota ekstraksi dalam bita, bawaan `4 GB` / 4.294.967.296 bita, perlindungan anti *decompression bomb*).
    - `"batas_entri"` / `"max_entries"`: `Angka` (batas total jumlah entri arsip yang diproses, bawaan `100.000`).
- **Kembalian**: `Angka` (jumlah berkas yang berhasil diekstrak).

```celang
// Ekstraksi aman dengan batas kuota dan saringan filter:
buat total = arsip.ekstrak_zip("paket.zip", "./distribusi", {
    "batas_ukuran": 50 * 1024 * 1024, // Maksimal 50 MB
    "saring": fungsi(info) {
        kembalikan teks.diakhiri_dengan(info.nama, ".ce")
    }
})
cetak("Berhasil mengekstrak", total, "berkas Celang.")
```

#### `arsip.ekstrak_berkas_zip(berkasZip, namaBerkas, tujuan, [opsi])`
Mengekstrak satu berkas spesifik dari dalam ZIP langsung ke disk (folder tujuan atau nama berkas baru) tanpa harus mengekstrak seisi arsip.

```celang
// Ekstrak ke dalam folder tujuan:
arsip.ekstrak_berkas_zip("arsip.zip", "docs/panduan.md", "./keluaran/")

// Ekstrak dengan mengganti nama berkas tujuan:
arsip.ekstrak_berkas_zip("arsip.zip", "docs/panduan.md", "./keluaran/baca_saya.md")
```

#### `arsip.apakah_zip(jalur)`
Memeriksa apakah berkas merupakan arsip ZIP yang valid berdasarkan *magic bytes* dan direktori sentral ZIP tanpa perlu mengekstrak atau melempar runtime error.
- **Kembalian**: `Logika` (`benar` jika valid ZIP, `salah` jika bukan).

```celang
jika arsip.apakah_zip("paket.unduhan") {
    arsip.ekstrak_zip("paket.unduhan", "./tujuan")
}
```

#### `arsip.daftar_zip(berkasZip)`
Melihat isi berkas ZIP tanpa melakukan ekstraksi ke disk.
- **Kembalian**: `Daftar` berisi `Peta` metadata untuk setiap berkas:
  - `nama`: Nama entri berkas.
  - `ukuran`: Ukuran berkas asli sebelum kompresi (dalam bita).
  - `ukuran_terkompresi`: Ukuran setelah kompresi.
  - `crc32`: Checksum CRC-32.
  - `metode_kompresi` / `compression_method`: `"deflate"` atau `"store"`.
  - `apakah_direktori`: `benar` jika entri merupakan folder.
  - `apakah_tautan`: `benar` jika entri merupakan tautan simbolik (*symlink*).
  - `tipe`: Teks jenis entri (`"berkas"`, `"direktori"`, atau `"tautan"`).
  - `mode`: Nilai izin berkas POSIX (misal: `420` untuk `0644`).
  - `waktu_modifikasi`: String waktu RFC 3339.

```celang
buat daftar = arsip.daftar_zip("paket.zip")
untuk item dalam daftar {
    jika item.tipe == "berkas" {
        cetak("- Berkas:", item.nama, "(", item.ukuran, "bita, metode:", item.metode_kompresi, ")")
    }
}
```

#### `arsip.baca_berkas_zip(berkasZip, namaBerkas)`
Membaca langsung isi berkas tertentu dari dalam ZIP tanpa perlu mengekstrak seluruh arsip.
- **Kembalian**: `Bita` (*Bytes*).

```celang
buat isi = arsip.baca_berkas_zip("projek.zip", "package.json")
buat jsonTeks = isi.ke_teks()
```

---

### 2. Arsip TAR.GZ (Tarball Gzip)

#### `arsip.buat_tar_gz(sumber, tujuanTarGz)`
Mengemas berkas atau folder menjadi file `.tar.gz`.

```celang
arsip.buat_tar_gz("./modul_pengguna", "modul-1.0.0.tar.gz")
```

#### `arsip.ekstrak_tar_gz(berkasTarGz, folderTujuan, [opsi])`
Mengekstrak arsip `.tar.gz` ke direktori tujuan dengan penjagaan keamanan Tar Slip dan dukungan fungsi saring (`"saring"` / `"filter"`).

```celang
arsip.ekstrak_tar_gz("modul-1.0.0.tar.gz", "./celang_modules/modul_pengguna")
```

#### `arsip.ekstrak_berkas_tar_gz(berkasTarGz, namaBerkas, tujuan, [opsi])`
Mengekstrak satu berkas tertentu dari `.tar.gz` langsung ke disk.

```celang
arsip.ekstrak_berkas_tar_gz("modul.tar.gz", "celang.mod", "./tujuan/")
```

#### `arsip.apakah_tar_gz(jalur)`
Memeriksa apakah berkas merupakan arsip `.tar.gz` yang valid.

#### `arsip.daftar_tar_gz(berkasTarGz)`
Membaca daftar entri di dalam file `.tar.gz`.

#### `arsip.baca_berkas_tar_gz(berkasTarGz, namaBerkas)`
Membaca isi satu berkas dari dalam `.tar.gz` sebagai `Bita`.

---

### 3. Arsip TAR Polos (Uncompressed)

- `arsip.buat_tar(sumber, tujuanTar)`: Membuat arsip `.tar` uncompressed.
- `arsip.ekstrak_tar(berkasTar, folderTujuan, [opsi])`: Mengekstrak arsip `.tar` (mendukung filter callback).
- `arsip.ekstrak_berkas_tar(berkasTar, namaBerkas, tujuan, [opsi])`: Mengekstrak satu berkas dari TAR polos ke disk.
- `arsip.baca_berkas_tar(berkasTar, namaBerkas)`: Membaca satu berkas dari TAR polos ke memori sebagai `Bita`.
- `arsip.daftar_tar(berkasTar)`: Membaca daftar isi berkas `.tar`.
- `arsip.apakah_tar(jalur)`: Memeriksa apakah berkas merupakan arsip `.tar` polos yang valid.

---

### 4. Kompresi Memori (Gzip & Zlib)

#### `arsip.gzip(data)` & `arsip.gunzip(bitaTerkompresi)`
Mengompresi dan mendekompresi data menggunakan algoritma Gzip (RFC 1952).
- `data`: dapat berupa `Teks` atau `Bita`.
- `gunzip`: mengembalikan `Bita` yang dapat dikonversi dengan `.ke_teks()`.

```celang
buat teks = "Data besar yang dikirim melalui jaringan..."
buat comp = arsip.gzip(teks)
buat pulih = arsip.gunzip(comp).ke_teks()

harus_sama(pulih, teks)
```

#### `arsip.zlib_kompres(data)` & `arsip.zlib_dekompres(bita)`
Mengompresi dan mendekompresi data menggunakan format Zlib Deflate (RFC 1950).

```celang
buat bitaZlib = arsip.zlib_kompres("Pesan rahasia")
buat teksAsli = arsip.zlib_dekompres(bitaZlib).ke_teks()
```

---

## 🛡️ Keamanan & Mitigasi Kerentanan

> [!IMPORTANT]
> **Pencegahan Otomatis Zip Slip / Tar Slip**:
> Penyerang dapat membuat berkas arsip berbahaya yang memuat entri seperti `../../../../etc/shadow` atau `..\..\Windows\System32\cmd.exe`.
> 
> Modul `arsip` Celang secara ketat memvalidasi setiap jalur ekstraksi menggunakan resolusi jalur absolut kanonikal. Jika terdeteksi berkas yang mencoba meloloskan diri dari `folderTujuan`, operasi ekstraksi langsung dibatalkan dengan galat empati:
> 
> ```
> Kesalahan: arsip.ekstrak_zip gagal: terdeteksi ancaman path traversal (Zip/Tar Slip) pada entri: '../../pwned.txt'
> ```
