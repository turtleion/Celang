# Modul Matematika (`matematika` / `math`)

Modul `matematika` (tersedia juga dengan alias `math`) menyediakan fungsi kalkulasi matematika, trigonometri, pembulatan, dan bilangan acak yang efisien dan dwibasa.

---

## Memuat Modul

Modul ini terdaftar sebagai pustaka standar bawaan (*batteries included*), sehingga dapat diakses secara langsung:

```celang
// Menggunakan identifier Bahasa Indonesia
buat r = matematika.akar(144)

// Atau menggunakan identifier internasional
buat p = math.pow(2, 8)
```

---

## Daftar Fungsi

### 1. `akar(x)` / `sqrt(x)`
Menghitung akar kuadrat dari bilangan non-negatif.
```celang
cetak(matematika.akar(25)) // 5
cetak(math.sqrt(2))        // 1.4142135623730951
```

### 2. `pangkat(basis, eksponen)` / `pow(basis, eksponen)`
Menghitung perpangkatan $basis^{eksponen}$.
```celang
cetak(matematika.pangkat(3, 4)) // 81
cetak(math.pow(2, 10))          // 1024
```

### 3. `mutlak(x)` / `abs(x)`
Mengembalikan nilai absolut/mutlak (selalu bernilai positif atau nol).
```celang
cetak(matematika.mutlak(-42)) // 42
cetak(math.abs(-3.14))        // 3.14
```

### 4. Pembulatan Angka
- **`bulat(x)` / `round(x)`**: Membulatkan ke bilangan bulat terdekat.
- **`lantai(x)` / `floor(x)`**: Membulatkan ke bawah (*floor*).
- **`atap(x)` / `ceil(x)`**: Membulatkan ke atas (*ceiling*).

```celang
cetak(matematika.bulat(4.6))  // 5
cetak(matematika.lantai(4.9)) // 4
cetak(matematika.atap(4.1))   // 5
```

### 5. `acak([min, max])` / `random([min, max])`
- Tanpa argumen: Menghasilkan bilangan riil acak antara `0.0` hingga `1.0`.
- Dengan 2 argumen `acak(min, max)`: Menghasilkan bilangan bulat acak inklusif antara `min` dan `max`.

```celang
// Nilai desimal acak 0.0 - 1.0
buat rasio = matematika.acak()

// Dadu acak 1 sampai 6
buat dadu = matematika.acak(1, 6)
cetak("Angka dadu:", dadu)
```

### 6. `maks(...angka)` / `min(...angka)`
Mencari nilai terbesar (*maksimum*) atau terkecil (*minimum*) dari sekumpulan argumen.
```celang
cetak(matematika.maks(12, 45, 9, 88, 23)) // 88
cetak(matematika.min(12, 45, 9, 88, 23))  // 9
```

---

## Konstanta

- **`matematika.PI`** / **`math.PI`**: $\pi \approx 3.141592653589793$
- **`matematika.E`** / **`math.E`**: $e \approx 2.718281828459045$

```celang
buat kelilingLingkaran = 2 * math.PI * 7
cetak("Keliling lingkaran:", kelilingLingkaran)
```
