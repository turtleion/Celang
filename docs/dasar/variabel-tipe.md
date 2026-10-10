# Variabel & Tipe Data

Celang dirancang agar deklarasi variabel terasa natural, manusiawi, dan fleksibel. Anda dapat menulis kode secara dinamis (*dynamically-typed*) atau menyertakan tipe data opsional (*gradual typing*) saat menginginkan kejelasan arsitektur.

---

## 1. Deklarasi Variabel

Gunakan kata kunci `buat` untuk mendeklarasikan variabel baru yang nilainya dapat diubah (*mutable*).

```celang
buat nama = "Budi Santoso"
buat umur = 24
buat aktif = benar

// Mengubah nilai variabel
umur = 25
cetak("Umur Budi sekarang:", umur)
```

::: tip Nilai Awal Default
Jika variabel dideklarasikan tanpa nilai inisialisasi, nilainya adalah `kosong` (*null*).
```celang
buat alamat
cetak(alamat) // Output: kosong
```
:::

---

## 2. Konstanta

Jika Anda memiliki nilai tetap yang tidak boleh diubah selama eksekusi program, gunakan kata kunci `konstanta`.

```celang
konstanta PAJAK = 0.11
konstanta KOTA_ASAL = "Bandung"

// Mencoba mengubah konstanta akan memunculkan error waktu kompilasi/eksekusi
// PAJAK = 0.15 ❌ (Error: tidak dapat menetapkan ulang nilai konstanta)
```

---

## 3. Tipe Data Primitif

Celang memiliki tipe data dasar berikut:

| Tipe Data | Kata Kunci | Contoh Nilai | Keterangan |
| :--- | :--- | :--- | :--- |
| **Angka** | `Angka` | `42`, `-10`, `3.14159` | Bilangan bulat maupun pecahan |
| **Teks** | `Teks` | `"Halo"`, `'Dunia'`, `` `Format` `` | Deretan karakter UTF-8 |
| **Logika** | `Logika` | `benar`, `salah` | Boolean |
| **Kosong** | - | `kosong` | Representasi ketiadaan nilai (*null*) |

```celang
buat bilanganBulat = 100
buat bilanganDesimal = 99.95
buat statusAktif = benar
buat pesan = "Selamat pagi!"
buat belumDiisi = kosong
```

---

## 4. Tipe Data Koleksi

### Daftar / Array

Array di Celang dinamis dan dapat menampung elemen dengan berbagai tipe data:

```celang
buat buah = ["Apel", "Mangga", "Jeruk", "Pisang"]

// Mengakses elemen berdasarkan indeks (mulai dari 0)
cetak("Buah pertama:", buah[0]) // Apel

// Mengetahui panjang array
cetak("Jumlah buah:", panjang(buah)) // 4

// Menambahkan elemen baru
buah.tambah("Durian")
```

### Kamus / Map / Objek

Map menyimpan data dalam pasangan kunci (*key*) dan nilai (*value*):

```celang
buat dataPengguna = {
    "nama": "Siti Rahma",
    "divisi": "Teknologi",
    "gaji": 15000000,
    "aktif": benar
}

// Mengakses data menggunakan bracket atau notasi titik
cetak(dataPengguna["nama"]) // Siti Rahma
cetak(dataPengguna.divisi)   // Teknologi
```

---

## 5. Anotasi Tipe Data (Opsional)

Celang mendukung *type hints* agar kode Anda lebih mudah dibaca oleh tim dan editor:

```celang
// Menggunakan tanda titik dua (:)
buat namaLengkap: Teks = "Ahmad Dani"
buat saldo: Angka = 500000

// Atau menggunakan kata sambung 'itu'
buat batasKredit itu Angka = 10000000
buat statusVerifikasi itu Logika = benar
```

Jika terjadi ketidaksesuaian tipe saat penetapan nilai statis, Celang akan memberikan pesan peringatan atau kesalahan yang jelas dan ramah pengguna.

---

## 6. Pemeriksaan Tipe Data Bawaan

Gunakan fungsi bawaan `tipe(nilai)` untuk mengecek tipe objek saat runtime:

```celang
cetak(tipe(100))          // Output: Angka
cetak(tipe("Celang"))     // Output: Teks
cetak(tipe(benar))        // Output: Logika
cetak(tipe([1, 2, 3]))    // Output: Array
cetak(tipe({"a": 1}))     // Output: Map
```
