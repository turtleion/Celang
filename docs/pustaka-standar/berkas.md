# Pustaka Pengelola Berkas (`berkas` / `fs`)

Modul **`berkas`** (dapat juga dipanggil dengan nama alias internasional **`fs`**, singkatan dari *file system*) menyediakan ekosistem lengkap untuk berinteraksi dengan sistem penyimpanan komputer Anda.

Selain fungsi baca-tulis standar, modul `berkas` di Celang dirancang dengan **fitur-fitur khas** yang tidak ada di bahasa lain: aliran streaming CSP berbasis Saluran, penyimpanan atomik anti-rusak (*atomic write*), pemuatan cerdas JSON/Teks otomatis, dan deteksi salah ketik nama berkas (*typo detector*).

---

## Cara Memanggil

Modul ini sudah terpasang bawaan dan siap digunakan secara langsung melalui objek global maupun sintaks impor:

```celang
// 1. Langsung panggil via objek global:
berkas.tulis("./laporan.txt", "Isi catatan...")
buat data = fs.read("./laporan.txt")

// 2. Atau impor fungsi tertentu:
impor { muat, simpan, alir } dari "berkas"
```

---

## 🌟 Fitur-Fitur Khas & Unik Celang

### 1. Aliran Streaming CSP (`alir` / `stream`)
Membaca berkas raksasa (misalnya log berukuran 10 GB) tanpa takut kehabisan memori RAM. Fungsi ini mengembalikan **Saluran (`Saluran`)** yang membaca berkas baris demi baris di latar belakang:

```celang
// Efisien & hemat RAM: memproses baris demi baris saat tiba
untuk baris dalam berkas.alir("./server.log") {
    jika baris.mengandung("ERROR") {
        cetak("Peringatan log:", baris)
    }
}
```

---

### 2. Muat & Simpan Otomatis (`muat` / `load` & `simpan` / `save`)
Tidak perlu lagi mengurai JSON secara manual!
- Jika nama berkas berakhiran `.json`, **`muat()`** otomatis mengubah isi berkas menjadi Kamus (`Map`) atau Daftar (`Array`) Celang.
- Sebaliknya, **`simpan()`** otomatis mengubah data objek Celang menjadi format JSON yang rapi (*pretty printed*) dan menyimpannya secara atomik.

```celang
// 1. Otomatis membaca JSON menjadi Objek:
buat setelan = berkas.muat("./setelan.json")
cetak("Nama Pengguna:", setelan.pengguna)

// 2. Mengubah nilai dan langsung menyimpannya kembali:
setelan.tema = "gelap"
berkas.simpan("./setelan.json", setelan)
cetak("Setelan berhasil diperbarui!")
```

---

### 3. Penyimpanan Atomik (*Atomic Write*)
Fungsi `tulis()` dan `simpan()` menerapkan prinsip **Atomic Write** secara transparan:
Data ditulis terlebih dahulu ke berkas sementara di direktori yang sama, disinkronkan ke disk, lalu ditukar secara instan menggantikan berkas target. 
> Berkas Anda **tidak akan pernah rusak (*corrupt*) atau menjadi 0 byte** meskipun komputer mati lampu atau program terhenti mendadak saat proses tulis berlangsung.

---

### 4. Diagnostik Empati & *Typo Detector*
Jika Anda salah mengetikkan nama berkas, Celang tidak sekadar melempar galat mentah, melainkan memeriksa direktori sekitar dan memberikan saran cerdas:

```celang
buat isi = berkas.baca("./catatan_rapat.txt")
```
Jika Anda tidak sengaja mengetik `./catatan_rapatt.txt`, terminal Celang akan memberi tahu:
```text
❌ Kesalahan Saat Program Berjalan (Runtime)
│ Berkas './catatan_rapatt.txt' tidak ditemukan.
│ 💡 Saran Solusi:
│    Maksud Anda mungkin './catatan_rapat.txt'? Berkas serupa ditemukan di direktori target.
```

---

## Operasi Dasar Berkas

### Membaca Teks (`baca` / `read`)
Membaca seluruh isi berkas teks ke dalam sebuah variabel:
```celang
buat isi = berkas.baca("./dokumen.txt")
cetak(isi)
```

### Menulis Teks (`tulis` / `write`)
Menulis tulisan ke berkas. Jika folder induknya belum ada, Celang otomatis membuatnya:
```celang
berkas.tulis("./catatan/hari_ini.txt", "Halo dari Celang!")
```

### Menambahkan Teks di Akhir (`tambah` / `append`)
Menambahkan teks baru ke baris paling bawah berkas tanpa menimpa isi sebelumnya:
```celang
berkas.tambah("./riwayat.txt", "Pengguna masuk pada 20:00\n")
```

### Memeriksa Keberadaan (`ada` / `exists`)
Mengecek apakah berkas atau folder ada di disk:
```celang
jika berkas.ada("./kunci.txt") {
    cetak("Berkas kunci tersedia.")
}
```

---

## Operasi Folder & Direktori

| Fungsi Indonesia | Alias Internasional | Keterangan |
| :--- | :--- | :--- |
| `buat_folder(path)` | `mkdir(path)` | Membuat direktori (termasuk folder induknya) |
| `daftar(path?)` | `list(path?)` / `ls` | Mengembalikan daftar nama berkas/folder |
| `hapus(path)` | `remove(path)` | Menghapus 1 berkas atau folder kosong |
| `hapus_semua(path)` | `remove_all(path)` | Menghapus folder beserta seluruh isinya secara rekursif |
| `salin(asal, tujuan)` | `copy(src, dst)` | Menyalin berkas atau folder |
| `pindah(asal, tujuan)` | `move(src, dst)` | Memindahkan atau mengubah nama (*rename*) |

Contoh penggunaan:
```celang
// Buat folder proyek
berkas.buat_folder("./proyek/sumber")

// Lihat isi folder saat ini
buat daftarBerkas = berkas.daftar("./proyek")
untuk nama dalam daftarBerkas {
    cetak("Item:", nama)
}

// Salin berkas
berkas.salin("./proyek/a.txt", "./proyek/a_cadangan.txt")

// Pindah atau ubah nama
berkas.pindah("./proyek/a_cadangan.txt", "./proyek/b.txt")

// Hapus folder proyek dan seluruh isinya
berkas.hapus_semua("./proyek")
```

---

## Pemeriksaan Tipe & Metadata

### `apakah_folder` & `apakah_berkas`
```celang
jika berkas.apakah_folder("./data") {
    cetak("Ini adalah folder.")
}

jika berkas.apakah_berkas("./data/catatan.txt") {
    cetak("Ini adalah berkas.")
}
```

### `ukuran` (`size`)
Mengambil ukuran berkas dalam satuan byte:
```celang
buat bytes = berkas.ukuran("./laporan.pdf")
cetak("Ukuran dokumen:", bytes, "byte")
```

### `info` (`stat`)
Mengambil rincian metadata lengkap dalam bentuk Kamus (`Map`):
```celang
buat detail = berkas.info("./gambar.png")

cetak("Nama berkas  :", detail.nama)
cetak("Ukuran       :", detail.ukuran, "byte")
cetak("Apakah berkas:", detail.apakah_berkas)
cetak("Diubah pada  :", detail.diubah_pada)
cetak("Hak izin     :", detail.izin)
```

---

## Utilitas Jalur (*Path Helpers*)

Modul `berkas` juga menyertakan alat pengolah jalur berkas lintas platform (kompatibel dengan format pemisah Linux/macOS `/` maupun Windows `\`):

```celang
// 1. Menggabungkan jalur secara aman
buat jalur = berkas.gabung("pengguna", "dokumen", "laporan.txt")
// Hasil: "pengguna/dokumen/laporan.txt"

// 2. Mengambil nama berkas
buat nama = berkas.nama(jalur)
// Hasil: "laporan.txt"

// 3. Mengambil ekstensi
buat ext = berkas.ekstensi(jalur)
// Hasil: ".txt"

// 4. Mengambil folder induk
buat folder = berkas.induk(jalur)
// Hasil: "pengguna/dokumen"

// 5. Mengambil jalur absolut lengkap
buat mutlak = berkas.jalur_mutlak(jalur)
```
