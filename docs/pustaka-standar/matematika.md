# Pustaka Perhitungan Matematika (`matematika` / `math`)

Modul **`matematika`** (dapat juga dipanggil dengan alias internasional **`math`**) menyediakan ekosistem terlengkap untuk kalkulasi numerik, statistik data, trigonometri, hingga utilitas finansial & lokal Indonesia di Celang.

Selain rumus-rumus kalkulus dan aljabar standar, modul `matematika` di Celang dirancang dengan **fitur-fitur cerdas** yang human-centric:
- **Ejaan Angka Otomatis (`matematika.terbilang`)**: Mengubah bilangan menjadi kata ejaan bahasa Indonesia secara otomatis (sangat cocok untuk kwitansi, faktur, nota, dan aplikasi kasir).
- **Format Rupiah Cepat (`matematika.rupiah`)**: Memformat nominal angka ke mata uang Rupiah standar Indonesia lengkap dengan pemisah ribuan titik (`.`) dan opsi desimal koma (`,00`).
- **Kalkulasi Persentase Ramah Pengguna (`matematika.persen` & `matematika.dari_persen`)**: Menghitung persentase proporsi dan nilai diskon/pajak tanpa perlu mengingat rumus pembagian/perkalian 100.
- **Statistik Kumpulan Data Instan**: Menghitung total (`jumlah`), rata-rata (`rata_rata`), nilai tengah (`median`), nilai terbanyak (`modus`), dan variabilitas (`standar_deviasi`) langsung dari sebuah Daftar (`Array`) maupun deretan argumen.
- **Interpolasi & Pemetaan Rentang (`lerp` & `petakan_rentang`)**: Memudahkan animasi, transisi visual, simulasi, dan pengembangan game.
- **Penghitung Jarak Euclidean 2D & 3D (`matematika.jarak`)**: Menghitung jarak lurus antar koordinat titik instan.
- **Aritmatika Khusus**: Pembulatan presisi desimal kustom, pembatasan rentang (`jepit` / `clamp`), faktor persekutuan terbesar (`fpb` / `gcd`), dan kelipatan persekutuan terkecil (`kpk` / `lcm`).

---

## Cara Memanggil

Modul ini langsung aktif secara global maupun dapat diimpor secara modular:

```celang
// 1. Langsung panggil via objek global:
buat bayar = matematika.rupiah(1500000)
buat eja = math.terbilang(125000)

// 2. Atau menggunakan sintaks impor terstruktur:
impor { terbilang, rupiah, persen, rata_rata } dari "matematika"
```

---

## 🌟 Fitur Cerdas & Human-Centric Celang

### 1. Ejaan Terbilang Bahasa Indonesia (`terbilang` / `spell_number`)
Mengonversi angka menjadi ejaan kalimat bahasa Indonesia yang akurat (mendukung bilangan bulat positif, minus, hingga bilangan desimal/koma):

```celang
buat totalBelanja = 1250500
cetak(matematika.terbilang(totalBelanja))
// Hasil: "satu juta dua ratus lima puluh ribu lima ratus"

// Mendukung bilangan negatif dan koma desimal:
cetak(matematika.terbilang(-50))   // "minus lima puluh"
cetak(matematika.terbilang(15.5))  // "lima belas koma lima"
```

---

### 2. Format Mata Uang Rupiah (`rupiah` / `format_rupiah`)
Mengubah angka mentah menjadi teks mata uang resmi Indonesia dengan format titik ribuan dan koma desimal:

```celang
// Format standar tanpa desimal:
buat harga = 2500000
cetak(matematika.rupiah(harga)) // "Rp 2.500.000"

// Format dengan 2 digit sen desimal (opsi kedua bernilai benar):
cetak(matematika.rupiah(harga, benar)) // "Rp 2.500.000,00"

// Mendukung nominal negatif (misal saldo minus / tagihan):
cetak(matematika.rupiah(-75000)) // "-Rp 75.000"
```

---

### 3. Persentase & Diskon Cepat (`persen` & `dari_persen`)
Menghilangkan kerumitan kalkulasi persentase dan diskon/pajak:

```celang
// A. Hitung persentase: "Berapa persen kah 25 dari total 200?"
buat skor = matematika.persen(25, 200)
cetak("Pencapaian:", skor, "%") // "Pencapaian: 12.5 %"

// B. Hitung nilai potongan: "Berapa nilai diskon 15% dari total belanja Rp200.000?"
buat diskon = matematika.dari_persen(15, 200000)
cetak("Potongan harga:", matematika.rupiah(diskon)) // "Potongan harga: Rp 30.000"
```

---

### 4. Interpolasi Linear & Pemetaan Rentang (`lerp` & `petakan_rentang`)
Sangat berguna untuk transisi halus pada game dev, UI, dan grafik:

```celang
// Linear interpolation (lerp): Mencari nilai antara awal & akhir berdasarkan t (0.0 .. 1.0)
buat tengah = matematika.lerp(0, 100, 0.5)
cetak("Posisi tengah:", tengah) // 50

// Petakan rentang (re-mapping):
// Mengubah input rentang [0, 100] menjadi rentang opasitas [0.0, 1.0]
buat opasitas = matematika.petakan_rentang(75, 0, 100, 0.0, 1.0)
cetak("Opasitas layar:", opasitas) // 0.75
```

---

### 5. Jarak Antar Titik 2D & 3D (`jarak` / `distance` / `hypot`)
Menghitung jarak lurus Euclidean antara dua koordinat titik:

```celang
// Jarak 2D: (x1, y1) ke (x2, y2)
buat jarak2D = matematika.jarak(0, 0, 3, 4)
cetak("Jarak 2D:", jarak2D) // 5

// Jarak 3D: (x1, y1, z1) ke (x2, y2, z2)
buat jarak3D = matematika.jarak(0, 0, 0, 2, 3, 6)
cetak("Jarak 3D:", jarak3D) // 7
```

---

## 📊 Statistik & Analisis Kumpulan Data

Semua fungsi statistik dapat menerima argumen berupa sebuah **Daftar (`Array`)** maupun **deretan angka langsung (`...args`)**:

```celang
buat nilaiUjian = [70, 85, 90, 85, 60, 95, 85, 100]

// 1. Total penjumlahan seluruh elemen (sum):
cetak("Total nilai :", matematika.jumlah(nilaiUjian)) // 670

// 2. Rata-rata nilai (average / mean):
cetak("Rata-rata   :", matematika.rata_rata(nilaiUjian)) // 83.75

// 3. Nilai tengah (median):
cetak("Median      :", matematika.median(nilaiUjian)) // 85

// 4. Nilai yang paling sering muncul (modus / mode):
cetak("Modus       :", matematika.modus(nilaiUjian)) // 85

// 5. Standar deviasi (penyebaran data):
cetak("Std Deviasi :", matematika.bulat(matematika.standar_deviasi(nilaiUjian), 2))
```

---

## 📐 Trigonometri & Sudut

Semua fungsi trigonometri dasar bekerja dalam satuan **radian**. Gunakan `ke_derajat` atau `ke_radian` untuk konversi:

```celang
// 1. Konversi sudut:
buat sudutRadian = matematika.ke_radian(90) // ~1.57079 rad
buat sudutDerajat = matematika.ke_derajat(matematika.PI) // 180 derajat

// 2. Sinus, Kosinus, Tangen:
cetak(matematika.sin(sudutRadian)) // 1
cetak(matematika.cos(sudutRadian)) // 0 (atau mendekati 0 karena floating point)
cetak(matematika.tan(0))           // 0

// 3. Arcus (Invers Trigonometri):
cetak(matematika.asin(1)) // 1.57079... (Pi / 2)
cetak(matematika.atan2(4, 3)) // Sudut koordinat y=4, x=3
```

---

## 📈 Logaritma & Eksponensial

```celang
// 1. Logaritma Natural (ln / basis e):
cetak(matematika.log(matematika.E)) // 1

// 2. Logaritma Basis 10:
cetak(matematika.log10(1000)) // 3

// 3. Logaritma Basis 2:
cetak(matematika.log2(256)) // 8

// 4. Eksponensial (e^x):
cetak(matematika.eksponen(1)) // ~2.71828 (nilai e)
```

---

## ⚙️ Aritmatika Presisi & Utilitas

### 1. Pembulatan Fleksibel (`bulat`, `potong_desimal`, `lantai`, `atap`)
```celang
buat pi = 3.14159265

// Membulatkan dengan presisi digit desimal:
cetak(matematika.bulat(pi, 2))          // 3.14
cetak(matematika.bulat(pi, 4))          // 3.1416
cetak(matematika.potong_desimal(pi, 3)) // 3.141 (memotong tanpa membulatkan)

// Pembulatan ke integer terdekat, bawah, atau atas:
cetak(matematika.bulat(4.6))  // 5
cetak(matematika.lantai(4.9)) // 4
cetak(matematika.atap(4.1))   // 5
```

### 2. Pembatas Nilai (`jepit` / `clamp`)
Menjaga sebuah nilai agar tidak keluar dari batas bawah dan batas atas:
```celang
// Contoh: Membatasi volume audio antara 0 hingga 100
buat masukanVolume = 140
buat volumeAman = matematika.jepit(masukanVolume, 0, 100)
cetak("Volume:", volumeAman) // 100

buat volumeMinus = -25
cetak("Volume:", matematika.jepit(volumeMinus, 0, 100)) // 0
```

### 3. Tanda Bilangan (`tanda` / `sign`)
Menghasilkan `1` jika positif, `-1` jika negatif, dan `0` jika nol:
```celang
cetak(matematika.tanda(42))  // 1
cetak(matematika.tanda(-15)) // -1
cetak(matematika.tanda(0))   // 0
```

### 4. Faktorial, FPB, dan KPK
```celang
// Faktorial (n!):
cetak(matematika.faktorial(5)) // 120 (5 x 4 x 3 x 2 x 1)

// Faktor Persekutuan Terbesar (FPB / GCD):
cetak(matematika.fpb(24, 36)) // 12

// Kelipatan Persekutuan Terkecil (KPK / LCM):
cetak(matematika.kpk(4, 6)) // 12
```

### 5. Angka Acak & Nilai Ekstrem
```celang
// Acak float antara 0.0 - 1.0:
buat acakPecahan = matematika.acak()

// Acak bilangan bulat dalam rentang batas [min, max]:
buat dadu = matematika.acak(1, 6)

// Nilai terbesar dan terkecil:
cetak(matematika.maks([10, 45, 92, 18])) // 92
cetak(matematika.min(10, 45, 92, 18))    // 10
```

---

## 💎 Konstanta Matematika

| Nama Konstanta | Nilai | Penjelasan |
| :--- | :--- | :--- |
| **`PI`** | `3.141592653589793` | Nilai konstanta lingkaran $\pi$ |
| **`E`** | `2.718281828459045` | Nilai bilangan Euler $e$ |
| **`PHI`** / **`RASIO_EMAS`** | `1.618033988749895` | Rasio Emas (*Golden Ratio*) $\phi$ |
| **`INF`** / **`TAK_HINGGA`** | `+Infinity` | Nilai tak terhingga positif |

```celang
buat luasLingkaran = matematika.PI * matematika.pangkat(7, 2)
cetak("Luas lingkaran r=7:", matematika.bulat(luasLingkaran, 2))
```

---

## 📋 Tabel Rangkuman Fungsi

| Fungsi Bahasa Indonesia | Alias Internasional | Keterangan Singkat |
| :--- | :--- | :--- |
| **`terbilang(angka)`** | `spell_number` | Mengubah angka menjadi ejaan kata bahasa Indonesia |
| **`rupiah(angka, [desimal])`** | `format_rupiah` | Format mata uang Rupiah Indonesia (`Rp 1.500.000`) |
| **`persen(bagian, total)`** | `percent` | Menghitung persentase proporsi |
| **`dari_persen(persen, total)`**| `from_percent` | Menghitung nominal diskon atau pajak |
| **`lerp(a, b, t)`** | `lerp` | *Linear interpolation* antara titik $a$ dan $b$ |
| **`petakan_rentang(val, ...)`** | `map_range` | Mengonversi nilai dari rentang asal ke rentang target |
| **`jarak(x1, y1, [z1], ...)`** | `distance` / `hypot` | Menghitung jarak Euclidean koordinat 2D atau 3D |
| **`jumlah(daftar)`** | `sum` | Menghitung total seluruh angka dalam daftar/argumen |
| **`rata_rata(daftar)`** | `average` / `mean` | Menghitung rata-rata aritmatika |
| **`median(daftar)`** | `median` | Mengambil nilai tengah kumpulan data |
| **`modus(daftar)`** | `mode` | Mengambil nilai yang frekuensi kemunculannya tertinggi |
| **`standar_deviasi(daftar)`** | `std_dev` | Menghitung standar deviasi data |
| **`jepit(nilai, min, max)`** | `clamp` | Membatasi angka agar selalu di rentang [min, max] |
| **`tanda(angka)`** | `sign` | Menghasilkan `-1`, `0`, atau `1` |
| **`faktorial(n)`** | `factorial` | Menghitung faktorial $n!$ |
| **`fpb(a, b)`** | `gcd` | Faktor Persekutuan Terbesar |
| **`kpk(a, b)`** | `lcm` | Kelipatan Persekutuan Terkecil |
| **`bulat(x, [desimal])`** | `round` | Membulatkan ke integer atau ke $N$ digit desimal |
| **`potong_desimal(x, [desimal])`**| `trunc` | Memotong digit desimal tanpa pembulatan |
| **`lantai(x)`** | `floor` | Membulatkan ke bawah (*floor*) |
| **`atap(x)`** | `ceil` | Membulatkan ke atas (*ceiling*) |
| **`akar(x)`** | `sqrt` | Akar kuadrat $\sqrt{x}$ |
| **`pangkat(basis, eksponen)`**| `pow` | Pemangkatan $\text{basis}^{\text{eksponen}}$ |
| **`mutlak(x)`** | `abs` | Nilai absolut $\|x\|$ |
| **`acak([min, max])`** | `random` | Menghasilkan angka acak pecahan atau integer |
| **`maks(...angka / daftar)`** | `max` | Nilai terbesar |
| **`min(...angka / daftar)`** | `min` | Nilai terkecil |
| **`sin(rad)` / `cos` / `tan`** | `sin` / `cos` / `tan` | Nilai trigonometri sudut radian |
| **`asin(x)` / `acos` / `atan`**| `asin` / `acos` / `atan`| Invers trigonometri sudut |
| **`atan2(y, x)`** | `atan2` | Sudut arc-tangent 2 koordinat |
| **`ke_derajat(radian)`** | `to_degrees` | Konversi radian ke derajat |
| **`ke_radian(derajat)`** | `to_radians` | Konversi derajat ke radian |
| **`log(x)`** | `log` | Logaritma natural $\ln(x)$ |
| **`log10(x)`** | `log10` | Logaritma basis 10 |
| **`log2(x)`** | `log2` | Logaritma basis 2 |
| **`eksponen(x)`** | `exp` | Menghitung eksponensial $e^x$ |
