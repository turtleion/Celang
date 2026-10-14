# Alat Bantu Pengolah Teks (`teks` / `strings`)

Modul `teks` (bisa juga dipanggil dengan nama `strings`) menyediakan berbagai alat praktis untuk mengolah kalimat, mulai dari memotong spasi, memecah kata, menyambung daftar, hingga mencari kata tertentu di dalam kalimat.

---

## Cara Memanggil

Modul ini sudah siap digunakan tanpa perlu pengaturan apa pun:

```celang
// Menggunakan nama Bahasa Indonesia:
buat buah = teks.pisah("apel,mangga,jeruk", ",")

// Atau menggunakan nama internasional:
buat hasil = strings.join(["A", "B", "C"], " - ")
```

---

## Daftar Alat Bantu

### 1. Memecah dan Menyambung Kalimat
- **`pisah(kalimat, pemisah)` / `split(str, sep)`**: Memecah kalimat panjang menjadi potongan daftar kata berdasarkan tanda pemisah tertentu.
- **`gabung(daftar, pemisah)` / `join(arr, sep)`**: Menyambung potongan-potongan kata menjadi satu kalimat utuh dengan tanda penghubung.

```celang
buat belanjaan = "kopi,teh,susu"
buat daftar = teks.pisah(belanjaan, ",")
cetak(daftar) // ["kopi", "teh", "susu"]

buat kalimat = teks.gabung(daftar, " & ")
cetak(kalimat) // "kopi & teh & susu"
```

### 2. Mengubah Huruf Besar dan Kecil
- **`kapital(kalimat)` / `upper(str)`**: Mengubah seluruh huruf menjadi huruf besar (kapital).
- **`kecil(kalimat)` / `lower(str)`**: Mengubah seluruh huruf menjadi huruf kecil.

```celang
cetak(teks.kapital("belajar celang itu menyenangkan")) // "BELAJAR CELANG ITU MENYENANGKAN"
cetak(teks.kecil("HURUF BESAR"))                     // "huruf besar"
```

### 3. Merapikan Spasi di Ujung Kata
- **`potong(kalimat)` / `trim(str)`**: Membuang spasi kosong yang tidak sengaja terketik di awal atau akhir tulisan.

```celang
buat masukanKetik = "   budi@email.com   "
buat emailBersih = teks.potong(masukanKetik)
cetak(emailBersih) // "budi@email.com"
```

### 4. Mengganti Kata Tertentu
- **`ganti(kalimat, kataLama, kataBaru)` / `replace(str, old, new)`**: Menemukan kata lama di dalam kalimat dan menggantinya dengan kata baru.

```celang
buat surat = "Halo [NAMA], selamat datang!"
buat hasilSurat = teks.ganti(surat, "[NAMA]", "Budi")
cetak(hasilSurat) // "Halo Budi, selamat datang!"
```

### 5. Memeriksa Isi Kalimat
- **`berisi(kalimat, kataCari)` / `contains(str, substr)`**: Memeriksa apakah kata yang dicari ada di dalam kalimat (memberikan jawaban `benar` atau `salah`).
- **`dimulai_dengan(kalimat, awalan)` / `starts_with(...)`**: Memeriksa apakah kalimat diawali dengan kata tertentu.
- **`diakhiri_dengan(kalimat, akhiran)` / `ends_with(...)`**: Memeriksa apakah kalimat diakhiri dengan kata tertentu.

```celang
buat namaDokumen = "laporan_keuangan.pdf"

cetak(teks.berisi(namaDokumen, "keuangan"))      // benar
cetak(teks.dimulai_dengan(namaDokumen, "laporan")) // benar
cetak(teks.diakhiri_dengan(namaDokumen, ".pdf"))   // benar
```
