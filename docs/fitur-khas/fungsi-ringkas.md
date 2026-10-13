# 🏹 Fungsi Ringkas Panah (Arrow Functions `=>`)

Celang menyediakan sintaks **Fungsi Panah (`=>`)** ala JavaScript/TypeScript modern yang sangat ringkas untuk membuat fungsi anonim (*lambda*) tanpa perlu menulis kata kunci `fungsi { ... }`.

---

## 💡 Mengapa Menggunakan Fungsi Panah?

Bandingkan penulisan fungsi konvensional dengan fungsi panah:

```celang
// Gaya konvensional:
buat kaliDuaLama = fungsi(x) { kembalikan x * 2 }

// Gaya panah ringkas:
buat kaliDua = x => x * 2
```

Lebih bersih, ekspresif, dan hemat baris kode!

---

## 1. Variasi Bentuk Sintaks

### Parameter Tunggal (Tanpa Tanda Kurung)
Jika fungsi hanya membutuhkan satu parameter, Anda tidak perlu menuliskan tanda kurung:

```celang
buat kuadrat = x => x * x
cetak(kuadrat(7)) // 49
```

### Banyak Parameter (Dengan Tanda Kurung)
Jika ada dua parameter atau lebih, bungkus daftar parameter di dalam `( ... )`:

```celang
buat tambah = (a, b) => a + b
buat hitungLuas = (panjang, lebar) => panjang * lebar

cetak(tambah(10, 25))       // 35
cetak(hitungLuas(5, 4))     // 20
```

### Tanpa Parameter
Gunakan tanda kurung kosong `()` jika fungsi tidak menerima argumen:

```celang
buat salam = () => "Halo, selamat datang di Celang!"
cetak(salam())
```

### Parameter dengan Tipe Data Eksplisit
Anda tetap bisa menentukan tipe data parameter untuk keamanan kode:

```celang
buat tambahSatu = (x: Angka) => x + 1
cetak(tambahSatu(99)) // 100
```

### Tubuh Fungsi Berbentuk Blok (`{ ... }`)
Jika fungsi membutuhkan beberapa langkah pemrosesan atau baris jamak, bungkus tubuh fungsi dengan kurung kurawal `{ ... }`:

```celang
buat sapaNama = nama => {
    buat pesan = "Halo, " + nama + "!"
    cetak("Mencatat ke log...")
    kembalikan pesan
}

cetak(sapaNama("Andi"))
```

---

## 2. Sangat Cocok untuk Functional Programming & Pipeline

Kekuatan terbesar fungsi panah terasa saat dikombinasikan dengan fungsi koleksi seperti `saring()` (*filter*), `petakan()` (*map*), dan **Pipeline Operator (`|>`)**:

### Pemrosesan Koleksi Ringkas
```celang
buat nilaiSiswa = [45, 78, 90, 62, 85, 30]

// Ambil nilai yang lulus (>= 60), lalu tambahkan bonus 5 poin
buat nilaiLulusBonus = nilaiSiswa
    |> saring(skor => skor >= 60)
    |> petakan(skor => skor + 5)

cetak(nilaiLulusBonus) // [83, 95, 67, 90]
```

### Menyaring Objek Kompleks
```celang
buat daftarProduk = [
    { "nama": "Buku", "harga": 25000 },
    { "nama": "Laptop", "harga": 8000000 },
    { "nama": "Pena", "harga": 5000 }
]

buat produkMurah = daftarProduk |> saring(item => item.harga < 50000)
cetak(produkMurah)
```

---

## 🎯 Rangkuman Bentuk Sintaks

| Bentuk Sintaks | Penjelasan |
| :--- | :--- |
| `x => x * 2` | Parameter tunggal, satu ekspresi otomatis dikembalikan |
| `(a, b) => a + b` | Banyak parameter |
| `() => "selesai"` | Tanpa parameter |
| `(a: Angka, b: Angka) => a + b` | Parameter bertipe data eksplisit |
| `x => { baris1; baris2; }` | Tubuh fungsi multi-baris |
