# Operator Penyaluran (Pipeline Operator `|>`)

Ketika Anda mengolah data melalui beberapa tahap secara berturut-turut, menulis fungsi di dalam fungsi yang bertumpuk-tumpuk sering kali membuat kode sulit dibaca karena harus dibaca dari tengah ke luar:
```celang
// Sulit dibaca karena bertumpuk dari dalam ke luar:
buat hasil = kapitalisasi(potong_spasi("  belajar celang  "))
```

Tanda penyaluran berantai `|>` memecahkan masalah ini. Tanda ini mengalirkan hasil perhitungan dari langkah sebelumnya langsung ke langkah berikutnya, seperti air mengalir di pipa atau ban berjalan di pabrik.

---

## 🌊 Contoh Memoles Teks

Dengan tanda `|>`, Anda dapat membaca alur kerja secara wajar dari atas ke bawah:

```celang
buat hasil = "  belajar celang  "
    |> potong_spasi()
    |> kapitalisasi()

cetak(hasil) // "BELAJAR CELANG"
```

1. Mulai dari teks mentah `"  belajar celang  "`.
2. Alirkan ke `potong_spasi()` untuk membuang spasi di ujungnya.
3. Alirkan hasilnya ke `kapitalisasi()` untuk mengubah semua huruf menjadi besar.

---

## 🔄 Mengolah Daftar Angka Bertahap

Penyaluran ini sangat nyaman dipadukan dengan alat penyaring (`saring`) dan pemeta (`petakan`):

```celang
buat deret = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]

buat hasil = deret
    |> saring(fungsi(x) => x % 2 == 0)  // Ambil hanya angka genap: [2, 4, 6, 8, 10]
    |> petakan(fungsi(x) => x * 10)     // Kalikan masing-masing dengan 10

cetak(hasil) // [20, 40, 60, 80, 100]
```

---

## 🏷️ Jika Fungsi Membutuhkan Keterangan Tambahan

Jika fungsi tujuan memerlukan lebih dari satu masukan, nilai yang dialirkan dari pipa sebelah kiri akan otomatis menjadi **masukan yang pertama**:

```celang
fungsi bungkusTag(isi, tag) {
    kembalikan "<" + tag + ">" + isi + "</" + tag + ">"
}

buat teksBungkus = "Selamat Datang"
    |> bungkusTag("b")
    |> bungkusTag("div")

cetak(teksBungkus) // <div><b>Selamat Datang</b></div>
```
Kalimat `"Selamat Datang"` masuk sebagai `isi`, sedangkan `"b"` menjadi `tag`-nya. Hasilnya kemudian dialirkan lagi ke langkah berikutnya.
