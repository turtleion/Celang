# Operator & Ekspresi

Operator adalah simbol atau kata kunci khusus yang digunakan untuk melakukan komputasi, perbandingan, manipulasi data, atau pengaliran logika di dalam Celang.

Celang mendukung penulisan simbolik standar maupun kata penghubung bahasa manusia (seperti `dan`, `atau`, `bukan`, `kalau kosong`).

---

## 1. Operator Aritmatika

Digunakan untuk operasi matematika dasar pada tipe `Angka` dan `Pecahan`:

| Operator | Nama | Contoh | Keterangan |
| :--- | :--- | :--- | :--- |
| `+` | Penjumlahan / Penggabungan | `a + b` | Menjumlahkan angka atau menggabungkan teks |
| `-` | Pengurangan / Negasi | `a - b`, `-x` | Mengurangi angka atau membalikkan tanda |
| `*` | Perkalian | `a * b` | Mengalikan dua angka |
| `/` | Pembagian | `a / b` | Membagi angka |
| `%` | Modulo (Sisa Bagi) | `a % b` | Menghitung sisa hasil pembagian bulat |

```celang
buat x = 15
buat y = 4

cetak(x + y) // 19
cetak(x - y) // 11
cetak(x * y) // 60
cetak(x / y) // 3.75
cetak(x % y) // 3
```

::: tip Penggabungan Teks dengan `+`
Operator `+` juga otomatis menggabungkan string teks dengan tipe data lain:
```celang
buat skor = 95
cetak("Skor akhir: " + skor) // "Skor akhir: 95"
```
:::

---

## 2. Operator Penugasan (*Assignment*)

Digunakan untuk menyimpan atau memodifikasi nilai dalam variabel:

| Operator | Contoh | Ekuivalen Dengan | Keterangan |
| :--- | :--- | :--- | :--- |
| `=` | `x = 10` | `x = 10` | Penugasan nilai dasar |
| `+=` | `x += 5` | `x = x + 5` | Tambah lalu tugaskan |
| `-=` | `x -= 3` | `x = x - 3` | Kurang lalu tugaskan |
| `*=` | `x *= 2` | `x = x * 2` | Kali lalu tugaskan |
| `/=` | `x /= 4` | `x = x / 4` | Bagi lalu tugaskan |
| `%=` | `x %= 3` | `x = x % 3` | Modulo lalu tugaskan |

```celang
buat tabungan = 100000
tabungan += 50000  // tabungan = 150000
tabungan -= 20000  // tabungan = 130000
tabungan *= 2      // tabungan = 260000
```

---

## 3. Operator Perbandingan (*Comparison*)

Menghasilkan nilai boolean `Logika` (`benar` atau `salah`):

| Operator | Nama | Contoh | Hasil jika `a = 10, b = 20` |
| :--- | :--- | :--- | :--- |
| `==` | Sama dengan | `a == b` | `salah` |
| `!=` | Tidak sama dengan | `a != b` | `benar` |
| `>` | Lebih besar dari | `a > b` | `salah` |
| `<` | Lebih kecil dari | `a < b` | `benar` |
| `>=` | Lebih besar atau sama | `a >= 10` | `benar` |
| `<=` | Lebih kecil atau sama | `a <= 20` | `benar` |

```celang
buat batasUsia = 17
buat usiaUser = 19

jika usiaUser >= batasUsia {
    cetak("Memenuhi syarat umur.")
}
```

---

## 4. Operator Logika Boolean

Celang mendukung simbol konvensional (`&&`, `||`, `!`) dan kata penghubung bahasa Indonesia (`dan`, `atau`, `bukan` / `tidak`):

| Operator Simbol | Kata Kunci Celang | Keterangan | Contoh |
| :--- | :--- | :--- | :--- |
| `&&` | `dan` | Bernilai `benar` jika kedua sisi benar (*Short-circuit AND*) | `adaKTP dan cukupUmur` |
| `\|\|` | `atau` | Bernilai `benar` jika salah satu sisi benar (*Short-circuit OR*) | `punyaSIM atau punyaPaspor` |
| `!` | `bukan`, `tidak` | Membalikkan nilai logika (*Negation*) | `bukan sudahSelesai` |

```celang
buat punyaAkun = benar
buat saldoCukup = salah

// Menggunakan kata kunci bahasa Indonesia
jika punyaAkun dan saldoCukup {
    cetak("Transaksi diproses.")
} selain {
    cetak("Transaksi ditolak.")
}

// Negasi boolean
buat sedangMemuat = salah
jika bukan sedangMemuat {
    cetak("Halaman siap ditampilkan.")
}
```

---

## 5. Operator Null Safety & Nilai Cadangan

Mencegah error *null pointer* dan memberikan nilai default secara instan:

### A. Navigasi Aman (`mungkin` / `?.`)
Mengakses rantai properti objek yang belum pasti ada tanpa memicu crash:
```celang
// Jika salah satu dari profil atau alamat bernilai kosong, hasilnya langsung kosong
buat kota = mungkin user.profil.alamat.kota
```

### B. Nilai Cadangan / Coalescing (`kalau kosong` / `??`)
Memberikan nilai pengganti jika sisi kiri bernilai `kosong`:
```celang
buat nama = inputNama kalau kosong "Tamu Anonim"
buat port = konfigurasi.port ?? 8080

// Dipadukan bersama 'mungkin':
buat kotaTinggal = mungkin user.alamat.kota kalau kosong "Jakarta"
```

> **Catatan Penting:** Boolean `salah` dan angka `0` **bukanlah** nilai kosong:
> ```celang
> buat aktif = salah kalau kosong benar // Hasil: salah (tidak diganti)
> buat skor = 0 kalau kosong 100        // Hasil: 0 (tidak diganti)
> ```

---

## 6. Operator Rentang (*Range*) & Slicing (`..`)

Digunakan untuk membuat deretan angka dan memotong koleksi:

```celang
// 1. Rentang Angka
buat angka = 1..5      // [1, 2, 3, 4, 5]

// 2. Loop Rentang
untuk i dalam 1..3 {
    cetak("Hitungan:", i)
}

// 3. Slicing Koleksi
buat daftar = [10, 20, 30, 40, 50]
buat sub = daftar[1..4]   // [20, 30, 40]
buat awal = daftar[..2]   // [10, 20]
buat akhir = daftar[3..]  // [40, 50]

buat teks = "Celang"[0..2] // "Ce"
```

---

## 7. Operator Aliran Pipeline (`|>`)

Mengalirkan hasil ekspresi sebelah kiri sebagai argumen pertama ke fungsi sebelah kanan:

```celang
buat teksMentah = "   belajar celang   "

buat hasil = teksMentah
    |> potong_spasi()
    |> kapitalisasi()

cetak(hasil) // "BELAJAR CELANG"
```

---

## 8. Operator Pointer Rendah (`&` dan `*`)

Selain referensi transparan (`rujuk`), Celang juga menyediakan akses tingkat rendah ke alamat memori variabel:

* **`&` (Address-of):** Mengambil alamat referensi memori variabel.
* **`*` (Dereference):** Membaca atau memutasi nilai langsung pada alamat memori tersebut.

```celang
buat angkaAsli = 50
buat ptr = &angkaAsli

// Mutasi nilai langsung di memori melalui pointer:
*ptr = 99
cetak(angkaAsli) // 99
```

---

## 9. Tabel Tingkat Presedensi Operator

Urutan prioritas evaluasi operator dari yang paling kuat (dieksekusi lebih dulu) hingga terendah:

| Prioritas | Kategori | Operator | Asosiatif |
| :---: | :--- | :--- | :---: |
| **1 (Tertinggi)**| Akses Member & Panggilan | `.`, `?.`, `[]`, `()` | Kiri ke Kanan |
| **2** | Unary / Prefix | `-` (negatif), `!`, `bukan`, `tidak`, `&`, `*`, `mungkin` | Kanan ke Kiri |
| **3** | Type Casting | `sebagai`, `as` | Kiri ke Kanan |
| **4** | Perkalian & Pembagian | `*`, `/`, `%` | Kiri ke Kanan |
| **5** | Penjumlahan & Pengurangan | `+`, `-` | Kiri ke Kanan |
| **6** | Rentang (*Range*) | `..` | Kiri ke Kanan |
| **7** | Perbandingan Relasional | `<`, `<=`, `>`, `>=` | Kiri ke Kanan |
| **8** | Perbandingan Kesetaraan | `==`, `!=` | Kiri ke Kanan |
| **9** | Logika DAN | `dan`, `&&` | Kiri ke Kanan |
| **10** | Logika ATAU | `atau`, `\|\|` | Kiri ke Kanan |
| **11** | Nilai Cadangan (*Coalescing*) | `kalau kosong`, `bila kosong`, `??` | Kiri ke Kanan |
| **12** | Aliran Pipeline | `\|>` | Kiri ke Kanan |
| **13 (Terendah)**| Penugasan (*Assignment*) | `=`, `+=`, `-=`, `*=`, `/=`, `%=` | Kanan ke Kiri |
