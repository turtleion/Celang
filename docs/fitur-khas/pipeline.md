# Pipeline Operator (`|>`)

Operator pipeline `|>` mengalirkan hasil ekspresi di sebelah kiri sebagai argumen pertama ke fungsi di sebelah kanan. Fitur ini menghilangkan sarang tanda kurung yang rumit (*nested calls*) dan membuat alur data mudah dibaca dari atas ke bawah.

---

## 🌊 Contoh Pemrosesan String

Tanpa pipeline:
```celang
buat hasil = kapitalisasi(potong_spasi("  belajar celang  "))
```

Dengan pipeline:
```celang
buat hasil = "  belajar celang  "
    |> potong_spasi()
    |> kapitalisasi()

cetak(hasil) // "BELAJAR CELANG"
```

---

## 🔄 Pemrosesan Koleksi (Filter & Map)

Pipeline sangat bertenaga saat dikombinasikan dengan fungsi transformasi array tingkat tinggi seperti `saring` (*filter*) dan `petakan` (*map*):

```celang
buat angka = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]

buat hasil = angka
    |> saring(fungsi(x) { kembalikan x % 2 == 0 })  // Ambil genap: [2, 4, 6, 8, 10]
    |> petakan(fungsi(x) { kembalikan x * 10 })     // Kalikan 10: [20, 40, 60, 80, 100]

cetak(hasil) // [20, 40, 60, 80, 100]
```

---

## 🏷️ Fungsi dengan Argumen Tambahan

Jika fungsi tujuan menerima lebih dari satu argumen, nilai dari sebelah kiri pipa selalu disuntikkan sebagai **argumen pertama**:

```celang
fungsi bungkus_tag(isi, tag) {
    kembalikan "<" + tag + ">" + isi + "</" + tag + ">"
}

buat html = "Selamat Datang"
    |> bungkus_tag("b")
    |> bungkus_tag("div")

cetak(html) // <div><b>Selamat Datang</b></div>
```
