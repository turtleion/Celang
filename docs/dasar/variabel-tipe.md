# Variabel & Tipe Data

Celang dirancang agar deklarasi variabel terasa alami, intuitif, dan aman. Bahasa ini mendukung pendekatan dinamis yang ringkas (*dynamic typing*) sekaligus deklarasi tipe bertahap (*gradual typing*) untuk kejelasan arsitektur proyek skala besar.

---

## 1. Deklarasi Variabel (`buat`)

Gunakan kata kunci `buat` untuk mendeklarasikan variabel baru yang nilainya dapat diubah (*mutable*):

```celang
buat nama = "Budi Santoso"
buat umur = 24
buat aktif = benar

// Mengubah nilai variabel
umur = 25
cetak("Umur sekarang:", umur) // 25
```

::: tip Nilai Awal Default
Jika suatu variabel dideklarasikan tanpa memberikan nilai awal, nilainya secara otomatis adalah `kosong` (*null*):
```celang
buat alamat
cetak(alamat) // Output: kosong
```
:::

---

## 2. Konstanta (`konstanta`)

Gunakan kata kunci `konstanta` untuk nilai tetap yang tidak boleh diubah selama eksekusi program:

```celang
konstanta PAJAK = 0.11
konstanta NAMA_APLIKASI = "CelangApp"

// Mencoba mengubah nilai konstanta akan memicu error:
// PAJAK = 0.15 ❌ (Galat: tidak dapat menetapkan ulang nilai konstanta)
```

---

## 3. Ruang Lingkup Variabel (*Scope & Shadowing*)

Celang menerapkan aturan ruang lingkup leksikal (*lexical block scope*):

```celang
buat x = "Global"

{
    // Variabel dalam blok terisolasi dari luar
    buat y = "Lokal Blok"
    cetak(x) // "Global" (bisa membaca lingkup luar)
    cetak(y) // "Lokal Blok"

    // Shadowing: mendeklarasikan ulang variabel dengan nama yang sama di dalam blok
    buat x = "Shadowed"
    cetak(x) // "Shadowed"
}

// cetak(y) ❌ (Galat: variabel 'y' tidak ditemukan di luar blok)
cetak(x) // "Global" (nilai luar tetap aman)
```

---

## 4. Tipe Data Primitif

Celang menyediakan tipe data dasar yang teroptimasi:

| Tipe Data | Kata Kunci | Contoh Nilai | Keterangan |
| :--- | :--- | :--- | :--- |
| **Bilangan Bulat** | `Angka` | `42`, `-10`, `1_000_000` | Integer 64-bit (mendukung pemisah visual `_`) |
| **Bilangan Desimal**| `Pecahan` / `Angka` | `3.14159`, `-0.05` | Float 64-bit berpresisi tinggi |
| **Teks** | `Teks` | `"Halo"`, `'Dunia'`, `` `Format` `` | Deretan karakter UTF-8 |
| **Logika** | `Logika` | `benar`, `salah` | Nilai kebenaran boolean |
| **Ketiadaan Nilai**| - | `kosong` / `nihil` | Representasi objek null / kosong |

### Penulisan Teks & Escape Karakter
Teks di Celang dapat ditulis menggunakan tanda kutip ganda (`"..."`), kutip tunggal (`'...'`), atau *backtick* (`` `...` ``):

```celang
buat teksBiasa = "Baris pertama\nBaris kedua\t(dengan tab)"
buat kutipDalam = 'Dia berkata: "Celang sangat mudah!"'
```

### Pengindeksan & Pemotongan Teks (*String Slicing*)
Teks dapat diindeks per karakter atau dipotong menggunakan operator rentang `..`:
```celang
buat bahasa = "Celang"
cetak(bahasa[0])      // "C"
cetak(bahasa[0..2])    // "Ce" (indeks 0 sampai sebelum 2)
cetak(bahasa[2..])     // "lang" (dari indeks 2 ke akhir)
```
*(Pemotongan teks di Celang bersifat rune-safe / mendukung karakter Unicode dan emoji).*

---

## 5. Tipe Data Koleksi

### A. Daftar / Array (`[...]`)
Array di Celang dinamis, fleksibel, dan dapat menampung elemen dengan tipe data campuran:

```celang
buat buah = ["Apel", "Mangga", "Jeruk", "Pisang"]

// Akses indeks (dimulai dari 0)
cetak("Buah pertama:", buah[0]) // "Apel"

// Mengetahui panjang array
cetak("Total item:", panjang(buah)) // 4

// Menambahkan elemen baru
buah.tambah("Durian") // atau buah.push("Durian")

// Pemotongan Array (Slicing)
buat pilihan = buah[1..3]  // ["Mangga", "Jeruk"]
buat awal = buah[..2]      // ["Apel", "Mangga"]
```

### B. Kamus / Map (`{...}`)
Map menyimpan data dalam pasangan kunci (*key*) dan nilai (*value*). Celang mendukung format standar (`"kunci": nilai`) maupun format deskriptif khas Celang (`"kunci" itu nilai`):

```celang
buat pengguna = {
    "nama": "Siti Rahma",
    "divisi" itu "Teknologi", // Gaya khas Celang
    "gaji": 15000000,
    "aktif": benar
}

// Mengakses nilai via notasi titik atau bracket
cetak(pengguna.nama)        // "Siti Rahma"
cetak(pengguna["divisi"])   // "Teknologi"

// Mengakses kunci yang tidak ada menghasilkan 'kosong' tanpa crash:
cetak(pengguna.alamat)      // kosong
```

---

## 6. Anotasi Tipe Data (*Gradual Typing*)

Anda dapat menambahkan keterangan tipe eksplisit menggunakan simbol titik dua (`:`) atau kata hubung `itu`:

```celang
// Menggunakan simbol titik dua (:)
buat namaLengkap: Teks = "Ahmad Dani"
buat saldo: Angka = 500000

// Atau menggunakan kata hubung 'itu'
buat batasKredit itu Angka = 10000000
buat statusVerifikasi itu Logika = benar
```

---

## 7. Referensi Transparan (`rujuk` / `ref`)

Secara *default*, penetapan variabel primitif di Celang menyalin nilai (*pass-by-value*). Jika Anda ingin dua variabel merujuk pada alamat memori yang sama tanpa kerumitan pointer, gunakan kata kunci `rujuk`:

```celang
buat angkaUtama = 10
rujuk aliasAngka = angkaUtama // aliasAngka terhubung ke angkaUtama

// Mengubah alias otomatis mengubah variabel sumber:
aliasAngka = 25
cetak(angkaUtama) // 25
```

Bekerja secara sangat mulus pada parameter fungsi:
```celang
fungsi kaliDua(rujuk x: Angka) {
    x *= 2
}

buat skor = 50
kaliDua(skor)
cetak(skor) // 100
```

---

## 8. Duplikasi Mendalam (*Deep Memory Clone*) dengan `salin`

Jika Anda bekerja dengan objek kompleks (Array bersarang, Map, atau Struct) dan ingin menyalinnya secara independen tanpa berbagi referensi memori, gunakan fungsi bawaan `salin` (atau `clone`):

```celang
buat daftarAsli = [1, 2, [10, 20]]
buat daftarKloning = salin(daftarAsli)

daftarKloning[0] = 99
daftarKloning[2][0] = 999

cetak(daftarAsli[0])       // 1 (tidak terpengaruh)
cetak(daftarAsli[2][0])    // 10 (elemen bersarang tetap aman!)
```

---

## 9. Alias Tipe (`tipe`)

Membuat nama alias semantik untuk tipe data agar kode lebih deskriptif:

```celang
tipe IDPengguna = Angka
tipe Surel = Teks

buat id: IDPengguna = 1001
buat kontak: Surel = "pengguna@celang.id"
```

Setiap alias tipe dapat langsung dipanggil sebagai fungsi konstruktor/casting:
```celang
buat idBaru = IDPengguna("2026") // Mengonversi teks "2026" -> Angka 2026
```

---

## 10. Konversi Tipe Eksplisit (*Type Casting*)

Celang mendukung konversi tipe eksplisit melalui dua pendekatan:

### A. Operator `sebagai` (atau `as`)
```celang
buat bilangan = "123" sebagai Angka    // 123
buat bulat = 3.99 sebagai Angka         // 3 (pemotongan desimal)
buat teks = 450 sebagai Teks            // "450"
buat id = "789" sebagai IDPengguna      // 789
```

### B. Konstruktor Tipe Bawaan
```celang
buat a = Angka("100")   // 100
buat t = Teks(500)       // "500"
buat b = Logika(1)       // benar
```
