# Fungsi Bawaan Global (Built-in Functions)

Fungsi bawaan (*built-in functions*) adalah fungsi-fungsi standar di Celang yang **selalu siap digunakan langsung** di seluruh bagian program Anda tanpa memerlukan perintah impor (`impor`) berkas atau modul apa pun.

Seluruh fungsi bawaan mendukung **nama Bahasa Indonesia yang ramah** serta **nama padanan internasional (Bahasa Inggris)**, dan terintegrasi penuh dengan sistem pesan empati serta Mode Universal (`"use universal"`).

---

## ⚡ Navigasi Kategori Tipe Data

Pilih kategori fungsi bawaan berdasarkan jenis data yang ingin Anda kelola:

- [1. Konversi Tipe Data & Meta-Objek](#1-konversi-tipe-data--meta-objek) (`tipe`, `Angka`, `Teks`, `Logika`, `salin`)
- [2. Manipulasi Teks / String](#2-manipulasi-teks--string) (`panjang`, `potong_spasi`, `kapitalisasi`)
- [3. Format Mata Uang & Angka Terbilang](#3-format-mata-uang--angka-terbilang) (`rupiah`, `terbilang`, `format_usd`, `in_words`, `dollars_in_words`)
- [4. Koleksi & Larik (Array & Iterasi)](#4-koleksi--larik-array--iterasi) (`rentang`, `saring`, `petakan`, `panjang`)
- [5. Masukan & Keluaran Konsol (I/O)](#5-masukan--keluaran-konsol-io) (`cetak`, `masukan`)
- [6. Konkurensi & Tugas Asinkron (Async / CSP)](#6-konkurensi--tugas-asinkron-async--csp) (`Saluran`, `jemput`, `jemput_semua`, `jemput_pertama`)

---

## 📋 Tabel Referensi Cepat

| Kategori | Perintah (ID) | Padanan Global (EN) | Argumen | Nilai Balik | Kegunaan Utama |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **Meta & Tipe** | `tipe(x)` | `type` | `x: ApaPun` | `Teks` | Mengetahui jenis data suatu nilai |
| **Meta & Tipe** | `Angka(x)` | `Number`, `number` | `x: ApaPun` | `Angka`/`Pecahan` | Mengubah teks atau nilai lain menjadi angka |
| **Meta & Tipe** | `Teks(x)` | `String`, `string` | `x: ApaPun` | `Teks` | Mengubah nilai apa pun menjadi tulisan teks |
| **Meta & Tipe** | `Logika(x)` | `Boolean`, `boolean` | `x: ApaPun` | `Logika` | Mengevaluasi status kebenaran boolean |
| **Meta & Tipe** | `salin(x)` | `clone` | `x: ApaPun` | Objek Baru | Membuat salinan utuh (*deep copy*) mandiri |
| **Teks** | `panjang(s)` | `len`, `length` | `s: Teks` | `Angka` | Menghitung jumlah karakter teks (Unicode) |
| **Teks** | `potong_spasi(s)`| `trim` | `s: Teks` | `Teks` | Memangkas spasi di ujung kiri dan kanan |
| **Teks** | `kapitalisasi(s)`| `uppercase`, `capitalize`| `s: Teks` | `Teks` | Mengubah seluruh huruf menjadi kapital |
| **Mata Uang & Kata**| `rupiah(n)` | `format_rupiah` | `n: Numerik` | `Teks` | Format nominal Rupiah (`Rp 150.000,00`) |
| **Mata Uang & Kata**| `terbilang(n)` | - | `n: Numerik` | `Teks` | Mengeja angka ke teks Bahasa Indonesia |
| **Mata Uang & Kata**| `format_usd(n)`| `format_currency` | `n: Numerik` | `Teks` | Format mata uang Dolar AS (`$150.00`) |
| **Mata Uang & Kata**| `in_words(n)` | `spelled_out` | `n: Numerik` | `Teks` | Mengeja angka ke teks Bahasa Inggris |
| **Mata Uang & Kata**| `dollars_in_words(n)`| - | `n: Numerik` | `Teks` | Mengeja nominal USD dalam Bahasa Inggris |
| **Koleksi** | `rentang(a, b)` | `range` | `[awal], akhir` | `Array` | Membuat urutan angka otomatis |
| **Koleksi** | `saring(arr, fn)`| `filter` | `arr, fungsi(el)` | `Array` | Menyaring elemen yang memenuhi kriteria |
| **Koleksi** | `petakan(arr, fn)`| `map` | `arr, fungsi(el)` | `Array` | Mengubah setiap elemen dengan fungsi |
| **Koleksi** | `panjang(arr/map)`| `len`, `length` | `arr: Array / Map` | `Angka` | Menghitung banyak elemen atau pasangan |
| **I/O Terminal** | `cetak(...val)` | `print` | `...nilai: ApaPun` | `kosong` | Menampilkan data ke layar terminal |
| **I/O Terminal** | `masukan(tanya)`| `input` | `[tanya: Teks]` | `Teks` | Menanyakan dan membaca masukan pengguna |
| **Konkurensi** | `Saluran(kap)` | `Channel`, `channel` | `[kapasitas: Angka]`| `Saluran` | Membuat saluran komunikasi CSP |
| **Konkurensi** | `jemput(tiket)` | `await` | `tiket, [timeout]` | Objek Hasil | Menunggu penyelesaian tugas latar belakang |
| **Konkurensi** | `jemput_semua(t)`| `await_all` | `daftarTiket: Array`| `Array` | Menunggu seluruh tiket selesai sekaligus |
| **Konkurensi** | `jemput_pertama(t)`| `await_race`, `await_first`| `daftarTiket: Array`| Objek Hasil | Mengambil hasil tiket yang paling awal selesai |

---

## 1. Konversi Tipe Data & Meta-Objek

Kelompok fungsi ini digunakan untuk memeriksa tipe data, melakukan konversi (*casting*), serta menduplikasi struktur data kompleks di memori.

### `tipe(nilai)` / `type`
Mengembalikan nama jenis data asli dari objek yang diberikan dalam bentuk `Teks`:

```celang
cetak(tipe(42))            // "Angka"
cetak(tipe(3.14))          // "Pecahan"
cetak(tipe("Halo"))        // "Teks"
cetak(tipe(benar))         // "Logika"
cetak(tipe([1, 2, 3]))     // "Array"
cetak(tipe({"a": 1}))      // "Map"
cetak(tipe(Saluran(5)))    // "Saluran"
cetak(tipe(kosong))        // "KOSONG"
```

### `Angka(nilai)` / `Number`
Mengubah teks angka atau nilai logika menjadi bilangan utuh (`Angka`) atau bilangan berkoma (`Pecahan`):

```celang
buat a = Angka("150000")   // 150000 (Angka)
buat b = Angka("3.14159")  // 3.14159 (Pecahan)
buat c = Angka(benar)      // 1
buat d = Angka(salah)      // 0

// Mendukung tulisan angka berformat mata uang secara cerdas:
buat harga = Angka("Rp 25.000") // 25000
```

### `Teks(nilai)` / `String`
Mengubah objek apa pun menjadi representasi tulisan (`Teks`):

```celang
buat t1 = Teks(123)        // "123"
buat t2 = Teks(benar)      // "benar"
buat t3 = Teks([1, 2, 3])  // "[1, 2, 3]"
```

### `Logika(nilai)` / `Boolean`
Mengevaluasi status nilai kebenaran (*truthiness*):
- Bernilai `salah` (*falsy*): angka `0`, teks kosong `""`, array kosong `[]`, atau nilai `kosong`.
- Bernilai `benar` (*truthy*): semua nilai lainnya yang memiliki isi.

```celang
cetak(Logika(1))          // benar
cetak(Logika(0))          // salah
cetak(Logika("Celang"))   // benar
cetak(Logika(""))         // salah
cetak(Logika([]))         // salah
```

### `salin(objek)` / `clone`
Melakukan penggandaan utuh (*deep copy*) pada memori terhadap Array, Map, atau Bentuk (Struct). Mengubah salinan data tidak akan pernah memengaruhi data aslinya:

```celang
buat asli = [10, 20, 30]
buat tiruan = salin(asli)

tiruan.push(40)
cetak(asli)   // [10, 20, 30] (asli tetap aman!)
cetak(tiruan) // [10, 20, 30, 40]
```

---

## 2. Manipulasi Teks / String

Fungsi bawaan untuk mengukur panjang dan memanipulasi susunan huruf pada tipe data `Teks`.

### `panjang(teks)` / `len` / `length`
Menghitung jumlah karakter dalam teks. Fungsi ini menghitung karakter Unicode secara akurat (*rune-aware*), sehingga emoji atau karakter multi-byte dihitung sebagai 1 karakter utuh:

```celang
buat judul = "Celang"
cetak(panjang(judul)) // 6

buat sapaan = "Halo 👋"
cetak(panjang(sapaan)) // 6 (termasuk 1 karakter emoji)
```

### `potong_spasi(teks)` / `trim`
Menghapus spasi kosong (*whitespace*), karakter tab, dan baris baru di ujung awal dan ujung akhir teks:

```celang
buat masukanKotor = "   \t selamat pagi!   \n "
buat bersih = potong_spasi(masukanKotor)
cetak(bersih) // "selamat pagi!"
```

### `kapitalisasi(teks)` / `uppercase` / `capitalize`
Mengubah semua huruf di dalam tulisan menjadi huruf besar (kapital):

```celang
buat status = "sukses"
cetak(kapitalisasi(status)) // "SUKSES"
```

> 💡 **Manipulasi Teks Lanjutan:** Untuk operasi string lanjutan seperti penggantian teks, pemisahan (`split`), pengecekan awalan/akhiran, gunakan modul bawaan [`teks`](/pustaka-standar/teks) atau [`strings`](/pustaka-standar/strings).

---

## 3. Format Mata Uang & Angka Terbilang

Celang memiliki dukungan bawaan tingkat bahasa untuk penanganan format moneter lokal (IDR) dan internasional (USD), serta konversi angka menjadi kata-kata (*number to words*).

### `rupiah(angka)` / `format_rupiah`
Memformat bilangan angka menjadi format standar mata uang Rupiah Indonesia (IDR):

```celang
cetak(rupiah(150000))     // "Rp 150.000,00"
cetak(rupiah(2500000.5))  // "Rp 2.500.000,50"
```

### `terbilang(angka)`
Mengubah nominal bilangan menjadi kalimat ejaan Bahasa Indonesia yang baku dan alami:

```celang
cetak(terbilang(150000))  // "seratus lima puluh ribu"
cetak(terbilang(1234567)) // "satu juta dua ratus tiga puluh empat ribu lima ratus enam puluh tujuh"
```

### `format_usd(angka)` / `format_currency`
Memformat bilangan angka menjadi format standar mata uang Dolar Amerika Serikat (USD):

```celang
cetak(format_usd(1250.75)) // "$1,250.75"
cetak(format_usd(50))      // "$50.00"
```

### `in_words(angka)` / `spelled_out`
Mengeja nominal angka menjadi kata-kata dalam Bahasa Inggris:

```celang
cetak(in_words(150))      // "one hundred fifty"
cetak(in_words(2024))     // "two thousand twenty-four"
```

### `dollars_in_words(angka)`
Mengeja nominal moneter dolar lengkap dengan sen dalam Bahasa Inggris standar:

```celang
cetak(dollars_in_words(150.25)) // "one hundred fifty dollars and twenty-five cents"
cetak(dollars_in_words(1.00))   // "one dollar"
```

> ⚠️ **Catatan Mode Universal:** Saat direktif `"use universal"` aktif di baris pertama file program, fungsi kultural lokal seperti `rupiah()` dan `terbilang()` dinonaktifkan demi kepatuhan standar internasional, dan Celang akan mengarahkan Anda menggunakan `format_usd()` atau `in_words()`.

---

## 4. Koleksi & Larik (Array & Iterasi)

Fungsi bawaan untuk membuat deretan nilai dan melakukan iterasi fungsional pada data array.

### `rentang(awal?, akhir)` / `range`
Menghasilkan array angka berurutan:
- Jika diberikan 1 argumen `rentang(akhir)`: menghasilkan deretan dari `0` sampai `akhir`.
- Jika diberikan 2 argumen `rentang(awal, akhir)`: menghasilkan deretan dari `awal` sampai `akhir` (bisa naik maupun turun).

```celang
buat deret1 = rentang(4)      // [0, 1, 2, 3, 4]
buat deret2 = rentang(1, 5)   // [1, 2, 3, 4, 5]
buat deret3 = rentang(5, 1)   // [5, 4, 3, 2, 1] (urutan menurun)
```

### `saring(daftar, rumus)` / `filter`
Menyaring elemen array berdasarkan fungsi predikat yang mengembalikan `benar` atau `salah`:

```celang
buat angka = [10, 15, 20, 25, 30]
buat kelipatanSepuluh = saring(angka, x => x % 10 == 0)
cetak(kelipatanSepuluh) // [10, 20, 30]
```

### `petakan(daftar, rumus)` / `map`
Mentransformasi setiap elemen di dalam array menjadi bentuk baru menggunakan fungsi pemetaan:

```celang
buat namaKecil = ["budi", "siti", "eko"]
buat namaBesar = petakan(namaKecil, s => kapitalisasi(s))
cetak(namaBesar) // ["BUDI", "SITI", "EKO"]
```

### `panjang(wadah)` / `len` / `length`
Menghitung jumlah elemen yang tersimpan di dalam array atau kamus map:

```celang
buat buah = ["Apel", "Mangga", "Jeruk"]
cetak(panjang(buah)) // 3

buat data = {"id": 1, "status": "aktif"}
cetak(panjang(data)) // 2
```

> 📖 **Metode Manipulasi Koleksi Lengkap:** Pelajari metode manipulasi modern berbasis objek (seperti `.push()`, `.pop()`, `.shift()`, `.insert()`, `.sort()`, `.keys()`, `.has()`, `.get()`, dll) pada bab khusus [Manipulasi Koleksi (Array & Map)](/fitur-khas/koleksi-array-map).

---

## 5. Masukan & Keluaran Konsol (I/O)

Perintah bawaan untuk berinteraksi dengan terminal atau konsol pengguna.

### `cetak(...nilai)` / `print`
Mencetak satu atau lebih nilai ke layar terminal secara otomatis dipisahkan spasi dan diakhiri dengan baris baru (*newline*):

```celang
cetak("Halo!")
cetak("Skor:", 95, "| Status: Lulus")
// Output: Skor: 95 | Status: Lulus
```

### `masukan(pertanyaan?)` / `input`
Menampilkan pesan pertanyaan di layar dan membaca teks yang diketikkan pengguna hingga tombol Enter ditekan:

```celang
buat nama = masukan("Masukkan nama Anda: ")
cetak("Selamat datang,", nama)
```

---

## 6. Konkurensi & Tugas Asinkron (Async / CSP)

Celang menyediakan fungsi bawaan untuk membuat saluran komunikasi antar-pekerjaan (*Communicating Sequential Processes / CSP*) dan mengelola tugas latar belakang (*background tickets*).

### `Saluran(kapasitas?)` / `Channel` / `channel`
Membuat saluran data (*channel*) untuk berkirim pesan antar-proses asinkron (`sambil`). Saluran dapat bersifat *unbuffered* (kapasitas 0) atau *buffered* (kapasitas > 0):

```celang
// Membuat saluran dengan kapasitas buffer 2 pesan:
buat ch = Saluran(2)

sambil {
    ch.kirim("pesan pertama")
    ch.kirim("pesan kedua")
    ch.tutup()
}

cetak(ch.terima()) // "pesan pertama"
cetak(ch.terima()) // "pesan kedua"
```

### `jemput(tiket, timeout?)` / `await`
Menunggu tugas latar belakang yang dibuat dengan kata kunci `titip` sampai selesai, dengan batas waktu (*timeout*) opsional:

```celang
buat tugas = titip {
    // Simulasi komputasi berat
    kembalikan 42 * 10
}

// Menunggu hasil tugas selesai:
buat hasil = jemput(tugas)
cetak("Hasil komputasi:", hasil) // 420
```

### `jemput_semua(daftarTiket)` / `await_all`
Menerima array berisi banyak tiket dan menunggu semuanya selesai secara paralel. Mengembalikan array berisi seluruh hasil eksekusi sesuai urutan tiket aslinya:

```celang
buat t1 = titip { 10 + 10 }
buat t2 = titip { 20 * 2 }
buat t3 = titip { 100 - 30 }

buat semuaHasil = jemput_semua([t1, t2, t3])
cetak(semuaHasil) // [20, 40, 70]
```

### `jemput_pertama(daftarTiket)` / `await_race` / `await_first`
Menunggu beberapa tiket tugas latar belakang dan mengembalikan hasil dari tiket yang **paling cepat selesai** pertama kali:

```celang
buat sumberA = titip { "Data dari Server A" }
buat sumberB = titip { "Data dari Server B" }

buat yangTercepat = jemput_pertama([sumberA, sumberB])
cetak("Server tercepat:", yangTercepat)
```

---

## Ringkasan Praktik Terbaik

1. **Gunakan Padanan Konsisten**: Anda bebas menggunakan nama Bahasa Indonesia (misal `cetak`, `panjang`, `rentang`) untuk nuansa lokal atau padanan Bahasa Inggris (`print`, `len`, `range`) untuk nuansa global.
2. **Prioritaskan Metode Objek untuk Koleksi**: Untuk manipulasi array dan map, metode berantai seperti `arr.push(x)` dan `map.get(k)` lebih ergonomis dibandingkan fungsi prosedural.
3. **Patuhi Mode Universal**: Saat bekerja dalam tim internasional dengan direktif `"use universal"`, gunakan fungsi global berbahasa Inggris (`format_usd`, `print`, `await`) agar kode terbebas dari kesalahan sintaks.
