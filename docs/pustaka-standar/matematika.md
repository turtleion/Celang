# Alat Bantu Perhitungan (`matematika` / `math`)

Modul `matematika` (bisa juga dipanggil dengan nama `math`) menyediakan berbagai rumus perhitungan angka siap pakai, mulai dari akar, pangkat, pembulatan, angka acak, hingga nilai tertinggi dan terendah.

---

## Cara Memanggil

Alat ini sudah terpasang secara otomatis di Celang. Anda bisa langsung memanggilnya menggunakan Bahasa Indonesia maupun bahasa internasional:

```celang
// Menggunakan nama Bahasa Indonesia:
buat r = matematika.akar(144)

// Menggunakan nama internasional:
buat p = math.pow(2, 8)
```

---

## Daftar Rumus dan Contoh Penggunaan

### 1. Mencari Akar (`akar` / `sqrt`)
Menghitung akar kuadrat dari suatu bilangan positif:
```celang
cetak(matematika.akar(25)) // 5
cetak(math.sqrt(144))      // 12
```

### 2. Menghitung Pangkat (`pangkat` / `pow`)
Mengalikan angka dengan dirinya sendiri sebanyak pangkat yang diminta:
```celang
cetak(matematika.pangkat(3, 4)) // 81 (artinya 3 x 3 x 3 x 3)
cetak(math.pow(2, 8))           // 256
```

### 3. Mengubah ke Nilai Positif (`mutlak` / `abs`)
Mengubah angka minus menjadi angka positif:
```celang
cetak(matematika.mutlak(-42)) // 42
cetak(math.abs(-3.14))        // 3.14
```

### 4. Pembulatan Angka Pecahan
Celang menyediakan 3 pilihan pembulatan sesuai kebutuhan Anda:
- **`bulat(x)` / `round(x)`**: Membulatkan ke angka bulat terdekat (di atas 0.5 naik ke atas, di bawah 0.5 turun ke bawah).
- **`lantai(x)` / `floor(x)`**: Selalu membulatkan ke bawah.
- **`atap(x)` / `ceil(x)`**: Selalu membulatkan ke atas.

```celang
cetak(matematika.bulat(4.6))  // 5
cetak(matematika.lantai(4.9)) // 4
cetak(matematika.atap(4.1))   // 5
```

### 5. Mengambil Angka Acak (`acak` / `random`)
Sangat berguna untuk permainan tebak angka atau undian:
- Jika dipanggil tanpa angka masukan `matematika.acak()`: Menghasilkan pecahan acak antara 0.0 sampai 1.0.
- Jika dipanggil dengan dua angka batas `matematika.acak(awal, akhir)`: Menghasilkan angka bulat acak di antara kedua batas tersebut.

```celang
// Melempar dadu (angka acak dari 1 sampai 6):
buat dadu = matematika.acak(1, 6)
cetak("Angka dadu yang keluar:", dadu)
```

### 6. Mencari Nilai Terbesar dan Terkecil (`maks` & `min`)
Membandingkan sekumpulan angka dan mengambil yang paling tinggi atau paling rendah:
```celang
cetak(matematika.maks(12, 45, 9, 88, 23)) // 88 (paling besar)
cetak(matematika.min(12, 45, 9, 88, 23))  // 9  (paling kecil)
```

---

## Angka Tetap (Konstanta)

- **`matematika.PI`** / **`math.PI`**: Nilai pi lingkaran ($\approx 3.141592653589793$)
- **`matematika.E`** / **`math.E`**: Nilai bilangan Euler ($\approx 2.718281828459045$)

```celang
buat kelilingLingkaran = 2 * matematika.PI * 7
cetak("Keliling lingkaran:", kelilingLingkaran)
```
