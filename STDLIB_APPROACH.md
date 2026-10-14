# Cara Celang Berkomunikasi dengan Komputer (`STDLIB_APPROACH.md`)

Pernahkah Anda bertanya:
> *"Bagaimana sebuah bahasa pemrograman bisa membaca berkas di harddisk, mengetahui jam dan tanggal, atau mencetak tulisan ke layar komputer?"*

Dokumen ini menjelaskan secara sederhana bagaimana **Pustaka Bawaan (Standard Library)** Celang berkomunikasi dengan sistem operasi komputer Anda.

---

## 1. Menghubungi Sistem Operasi Melalui Jembatan Terpercaya

Program komputer tingkat tinggi (seperti Celang, Python, atau JavaScript) tidak bisa sembarangan mengutak-atik perangkat keras komputer secara langsung. Sistem operasi (seperti Linux, Windows, atau macOS) menjaga hal tersebut demi keamanan data Anda.

Ketika program Celang Anda menulis:
```celang
buat isi = berkas.baca("./data.txt")
```

Berikut alur yang terjadi di balik layar:

```
┌────────────────────────────────────────────────────────┐
│ 1. Kode Celang Anda (.ce)                              │
│    "berkas.baca('./data.txt')"                         │
└───────────────────────────┬────────────────────────────┘
                            │
┌───────────────────────────▼────────────────────────────┐
│ 2. Mesin Evaluator Celang                              │
│    Mengenali bahwa pengguna ingin membaca dokumen      │
└───────────────────────────┬────────────────────────────┘
                            │
┌───────────────────────────▼────────────────────────────┐
│ 3. Jembatan Penghubung (Go Engine)                     │
│    Memanggil perintah resmi pembaca berkas komputer    │
└───────────────────────────┬────────────────────────────┘
                            │
┌───────────────────────────▼────────────────────────────┐
│ 4. Sistem Operasi Komputer (Linux / Mac / Windows)     │
│    Mengambil isi dokumen dari penyimpanan dan          │
│    memberikannya kembali ke Celang                     │
└────────────────────────────────────────────────────────┘
```

---

## 2. Mengapa Menggunakan Jembatan Ini?

1. **Aman dan Stabil**: Jembatan ini memanfaatkan kemampuan bahasa Go yang sudah terbukti sangat stabil, cepat, dan aman dari kerusakan memori.
2. **Bisa Berjalan di Mana Saja**: Perintah yang Anda tulis di Celang akan bekerja persis sama baik di komputer Linux, komputer Mac, maupun laptop Windows.
3. **Pesan Kesalahan yang Ramah**: Jika dokumen yang dicari tidak ada atau Anda tidak memiliki izin untuk membukanya, jembatan ini tidak membiarkan program mati mendadak. Jembatan ini membungkus masalah tersebut ke dalam kotak petunjuk yang ramah dan memberi tahu letak kesalahannya dengan jelas.

---

## 3. Lima Modul Utama yang Selalu Siap Pakai

Dengan arsitektur ini, Celang menyediakan lima kelompok alat bantu bawaan yang langsung aktif tanpa perlu mengunduh aplikasi lain:

1. **`berkas` / `fs`**: Mengelola dokumen di penyimpanan komputer (baca, tulis, tambah, hapus, cek keberadaan).
2. **`matematika` / `math`**: Melakukan perhitungan matematika presisi tinggi (akar, pangkat, nilai mutlak, pembulatan, angka acak, konstanta pi).
3. **`teks` / `strings`**: Memproses tulisan (memotong spasi, memecah kata, menyambung kalimat, mengubah huruf besar/kecil).
4. **`waktu` / `time`**: Mengukur durasi, mengambil waktu saat ini, memformat tanggal, dan menunda jeda sejenak.
5. **`json`**: Membaca dan menyusun data berformat standar JSON untuk bertukar informasi dengan aplikasi lain.
