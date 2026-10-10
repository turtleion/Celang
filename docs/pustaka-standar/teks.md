# Modul Teks (`teks` / `strings`)

Modul `teks` (atau alias `strings`) menyediakan fungsi manipulasi teks dan string yang lengkap, mendukung encoding UTF-8, pemotongan spasi, pemisahan kata, dan pencarian pola.

---

## Memuat Modul

Dapat digunakan secara langsung tanpa instalasi eksternal:

```celang
// Menggunakan identifier Bahasa Indonesia
buat kataKata = teks.pisah("apel,mangga,jeruk", ",")

// Atau menggunakan identifier internasional
buat hasil = strings.join(["A", "B", "C"], " - ")
```

---

## Daftar Fungsi

### 1. Pemisahan & Penggabungan
- **`pisah(teks, pemisah)` / `split(str, sep)`**: Memecah string menjadi array berdasarkan pemisah tertentu.
- **`gabung(array, pemisah)` / `join(arr, sep)`**: Menggabungkan elemen-elemen array menjadi satu string.

```celang
buat data = "kopi,teh,susu"
buat daftar = teks.pisah(data, ",")
cetak(daftar) // ["kopi", "teh", "susu"]

buat kalimat = teks.gabung(daftar, " & ")
cetak(kalimat) // "kopi & teh & susu"
```

### 2. Transformasi Huruf
- **`kapital(teks)` / `upper(str)`**: Mengubah seluruh teks menjadi huruf besar/kapital.
- **`kecil(teks)` / `lower(str)`**: Mengubah seluruh teks menjadi huruf kecil.

```celang
cetak(teks.kapital("celang bahasa keren")) // "CELANG BAHASA KEREN"
cetak(teks.kecil("KATA KAPITAL"))         // "kata kapital"
```

### 3. Pembersihan Spasi
- **`potong(teks)` / `trim(str)`**: Menghilangkan spasi dan karakter whitespace di awal dan akhir teks.

```celang
buat inputPengguna = "   admin@domain.com   "
buat emailBersih = teks.potong(inputPengguna)
cetak(emailBersih) // "admin@domain.com"
```

### 4. Penggantian Teks
- **`ganti(teks, kata_lama, kata_baru)` / `replace(str, old, new)`**: Mengganti semua kemunculan `kata_lama` dengan `kata_baru`.

```celang
buat template = "Halo [NAMA], selamat datang di [NAMA] Corp!"
buat pesan = teks.ganti(template, "[NAMA]", "Budi")
cetak(pesan) // "Halo Budi, selamat datang di Budi Corp!"
```

### 5. Pencarian & Validasi Pola
- **`berisi(teks, kata_cari)` / `contains(str, substr)`**: Memeriksa apakah teks mengandung substring tertentu (mengembalikan `benar` atau `salah`).
- **`dimulai_dengan(teks, awalan)` / `starts_with(str, prefix)`**: Memeriksa apakah teks diawali dengan awalan tertentu.
- **`diakhiri_dengan(teks, akhiran)` / `ends_with(str, suffix)`**: Memeriksa apakah teks diakhiri dengan akhiran tertentu.

```celang
buat namaBerkas = "dokumen_keuangan.pdf"

cetak(teks.berisi(namaBerkas, "keuangan"))      // benar
cetak(teks.dimulai_dengan(namaBerkas, "dokumen")) // benar
cetak(teks.diakhiri_dengan(namaBerkas, ".pdf"))   // benar
```
