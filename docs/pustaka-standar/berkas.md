# Pustaka Pengelola Berkas (`berkas` / `fs`)

Modul **`berkas`** (dapat juga dipanggil dengan alias internasional **`fs`**, singkatan dari *file system*) menyediakan ekosistem terlengkap untuk mengelola dokumen, folder, jalur penyimpanan, dan pengawasan berkas di komputer Anda.

Selain fungsi baca-tulis standar, modul `berkas` di Celang dirancang dengan **fitur-fitur khas** yang menjadikannya sangat modern dan nyaman digunakan:
- Aliran streaming baris CSP berbasis `Saluran` (`berkas.alir`).
- Pemantau perubahan berkas (*file watcher*) bawaan via `Saluran` (`berkas.pantau`).
- Pemuatan & penyimpanan cerdas data JSON/Teks (`berkas.muat` & `berkas.simpan`).
- Pencarian pola jalur (*globbing*) `**/*` bawaan (`berkas.cari`).
- Penggantian teks langsung di dalam berkas (*in-place find & replace*) (`berkas.ganti_teks`).
- Pembaca awal & akhir dokumen (*head & tail*) (`berkas.kepala` & `berkas.ekor`).
- Penghitung sidik jari berkas (*checksum / hash*) (`berkas.hash`).
- Penelusur folder rekursif (*tree walker*) (`berkas.telusuri`).
- Penyimpanan atomik anti-rusak (*atomic write*).
- Diagnostik ramah salah ketik nama berkas (*typo detector*).

---

## Cara Memanggil

Modul ini langsung aktif secara global maupun dapat diimpor secara spesifik:

```celang
// 1. Langsung panggil via objek global:
berkas.tulis("./laporan.txt", "Isi catatan...")
buat data = fs.read("./laporan.txt")

// 2. Atau impor fungsi tertentu (tanpa tanda kutip):
ambil { muat, simpan, alir, cari, pantau } dari berkas
```

---

## 🌟 Fitur-Fitur Khas & Unik Celang

### 1. Aliran Streaming CSP (`alir` / `stream`)
Membaca berkas raksasa (misalnya log berukuran puluhan GB) tanpa membebani memori RAM. Fungsi ini mengembalikan **`Saluran`** yang membaca baris demi baris di latar belakang:

```celang
untuk baris dalam berkas.alir("./server.log") {
    jika baris.mengandung("ERROR") {
        cetak("Peringatan log:", baris)
    }
}
```

---

### 2. Pemantau Berkas Bawaan (`pantau` / `watch`)
Memonitor perubahan berkas dan direktori secara *real-time* menggunakan aliran `Saluran` atau *callback*, tanpa butuh pustaka eksternal:

```celang
cetak("Memulai pemantauan folder proyek...")

buat pengawas = berkas.pantau("./src")
untuk kejadian dalam pengawas {
    cetak("Berkas berubah :", kejadian.path)
    cetak("Aksi perubahan :", kejadian.aksi) // "buat", "ubah", "hapus"
}

// Untuk menghentikan pengawasan:
// pengawas.tutup()
```

---

### 3. Pencari Berkas Wildcard (`cari` / `glob`)
Mencari berkas menggunakan pola pencarian (*globbing*) termasuk penelusuran subfolder rekursif `**`:

```celang
// Temukan semua berkas kode Celang di folder src dan seluruh subfoldernya:
buat berkasKode = berkas.cari("./src/**/*.ce")
cetak("Ditemukan berkas kode:", panjang(berkasKode))

// Temukan semua gambar PNG:
buat foto = berkas.cari("./assets/*.png")
```

---

### 4. Pemuat & Penyimpan Cerdas (`muat` / `load` & `simpan` / `save`)
Bebas *boilerplate*!
- Jika nama berkas berakhiran `.json`, **`muat()`** otomatis mengubah isi berkas menjadi Kamus (`Map`) atau Daftar (`Array`) Celang.
- Sebaliknya, **`simpan()`** otomatis mengubah data objek Celang menjadi format JSON rapi (*pretty printed*) secara atomik.

```celang
// 1. Otomatis membaca JSON menjadi Objek Celang:
buat setelan = berkas.muat("./setelan.json")
cetak("Tema saat ini:", setelan.tema)

// 2. Mengubah nilai dan langsung menyimpannya kembali:
setelan.tema = "gelap"
berkas.simpan("./setelan.json", setelan)
cetak("Setelan berhasil disimpan!")
```

---

### 5. Pengganti Teks Langsung (`ganti_teks` / `replace_text`)
Mengganti teks di dalam dokumen secara instan (*find & replace in-place*) dengan jaminan penulisan atomik:

```celang
// Ganti pengaturan konfigurasi dalam 1 baris kode:
berkas.ganti_teks("./config.env", "MODE=pengembangan", "MODE=produksi")
```

---

### 6. Pembaca Awal & Akhir Dokumen (`kepala` / `head` & `ekor` / `tail`)
Mengintip $N$ baris pertama atau $N$ baris terakhir dari berkas log besar tanpa perlu membuka seluruh berkas:

```celang
// Ambil 3 baris pertama:
buat awal = berkas.kepala("./app.log", 3)

// Ambil 5 baris terakhir (sangat berguna untuk melihat galat terbaru):
buat akhir = berkas.ekor("./app.log", 5)
```

---

### 7. Sidik Jari / Checksum Berkas (`hash` / `checksum`)
Menghitung kode hash isi berkas (default: SHA-256, atau MD5, SHA-1) secara efisien:

```celang
buat sha = berkas.hash("./dokumen.pdf")
buat md5 = berkas.hash("./dokumen.pdf", "md5")
cetak("SHA-256:", sha)
```

---

### 8. Penelusur Folder Rekursif (`telusuri` / `walk`)
Menjelajahi seluruh pohon folder dan berkas dengan fungsi *callback*:

```celang
berkas.telusuri("./proyek", fungsi (jalur, info) {
    jika info.ekstensi == ".json" {
        cetak("Menemukan berkas konfigurasi:", jalur, "Ukuran:", info.ukuran)
    }
})
```

---

### 9. Penyimpanan Atomik (*Atomic Write*)
Fungsi `tulis()`, `simpan()`, dan `ganti_teks()` menerapkan **Atomic Write** bawaan:
Data ditulis terlebih dahulu ke berkas temporer di direktori yang sama, disinkronkan ke disk, lalu ditukar secara instan menggantikan berkas target. 
> Berkas Anda **tidak akan pernah rusak (*corrupt*) atau menjadi 0 byte** meskipun komputer mati mendadak saat proses tulis berlangsung.

---

### 10. Diagnostik Empati (*Typo Detector*)
Jika Anda salah mengetikkan nama berkas, Celang memeriksa direktori sekitar menggunakan algoritma jarak Levenshtein dan memberikan saran:

```text
❌ Kesalahan Saat Program Berjalan (Runtime)
│ Berkas './catatan_rapatt.txt' tidak ditemukan.
│ 💡 Saran Solusi:
│    Maksud Anda mungkin './catatan_rapat.txt'? Berkas serupa ditemukan di direktori target.
```

---

## Ringkasan Daftar Fungsi

### Operasi Berkas & Folder
| Fungsi Indonesia | Alias Internasional | Keterangan |
| :--- | :--- | :--- |
| `baca(path)` | `read` | Baca seluruh isi teks berkas |
| `tulis(path, teks)` | `write` | Tulis berkas (otomatis buat folder & atomik) |
| `tambah(path, teks)` | `append` | Tambah teks di akhir berkas |
| `buat_folder(path)` | `mkdir` | Buat folder secara rekursif |
| `daftar(path?)` | `list` / `ls` | Ambil daftar nama isi folder |
| `hapus(path)` | `remove` | Hapus 1 berkas / folder kosong |
| `hapus_semua(path)` | `remove_all` | Hapus folder dan seluruh isinya secara rekursif |
| `salin(asal, tujuan)` | `copy` | Salin berkas atau folder |
| `pindah(asal, tujuan)` | `move` / `rename` | Pindah atau ubah nama berkas/folder |

### Pemeriksaan & Metadata
| Fungsi Indonesia | Alias Internasional | Keterangan |
| :--- | :--- | :--- |
| `ada(path)` | `exists` | Cek keberadaan (`benar` / `salah`) |
| `apakah_folder(path)` | `is_dir` | Cek apakah direktori |
| `apakah_berkas(path)` | `is_file` | Cek apakah berkas biasa |
| `ukuran(path)` | `size` | Ukuran berkas dalam byte |
| `info(path)` | `stat` | Kamus metadata (nama, ukuran, tanggal, izin) |

### Utilitas Jalur (*Path Helpers*)
| Fungsi Indonesia | Alias Internasional | Keterangan |
| :--- | :--- | :--- |
| `gabung(...paths)` | `join` | Gabung bagian jalur secara aman lintas OS |
| `nama(path)` | `basename` | Ambil nama berkas saja |
| `induk(path)` | `dirname` | Ambil folder induk |
| `ekstensi(path)` | `extname` | Ambil ekstensi berkas (`.txt`, `.json`) |
| `jalur_mutlak(path)`| `absolute` | Ambil path absolut lengkap |
| `bersihkan_jalur(path)` | `clean` | Normalisasi penulisan path |
