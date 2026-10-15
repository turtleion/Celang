# Fungsi Bawaan Global (Built-in Functions)

Celang menyediakan sejumlah perintah dan fungsi bawaan yang **selalu siap digunakan langsung** di mana saja di dalam program Anda, tanpa perlu memuat berkas atau modul tambahan.

Fungsi-fungsi ini mendukung nama Bahasa Indonesia yang nyaman serta nama padanan internasionalnya.

---

## 📋 Daftar Perintah Bawaan

| Kelompok | Perintah (Bahasa Indonesia) | Padanan Internasional | Kegunaan |
| :--- | :--- | :--- | :--- |
| **Layar & Masukan** | `cetak(...nilai)` | `print` | Menampilkan tulisan atau hasil hitungan ke layar terminal |
| | `masukan(pertanyaan?)` | `input` | Menanyakan dan menerima ketikan jawaban dari pengguna |
| **Pemeriksaan Data** | `panjang(wadah)` | `len` | Menghitung banyak huruf di kalimat atau banyak barang di daftar |
| | `tipe(nilai)` | `type` | Mengetahui jenis data (seperti `Angka`, `Teks`, atau `Array`) |
| **Pengubah Jenis Data** | `Angka(nilai)` | `angka` | Mengubah tulisan atau data menjadi angka murni |
| | `Teks(nilai)` | `teks` | Mengubah angka atau data apa pun menjadi tulisan teks |
| | `Logika(nilai)` | `logika` | Memeriksa apakah suatu nilai bernilai ada (`benar`) atau kosong (`salah`) |
| **Daftar & Kumpulan** | `rentang(awal?, akhir)` | `range` | Menghasilkan urutan angka dari awal sampai akhir |
| | `saring(daftar, rumus)` | `filter` | Menyaring dan mengambil hanya barang yang lolos kriteria |
| | `petakan(daftar, rumus)` | `map` | Mengubah setiap barang di daftar menggunakan rumus tertentu |
| **Pembersih Teks** | `potong_spasi(teks)` | `trim` | Membuang spasi kosong di ujung kiri dan kanan tulisan |
| | `kapitalisasi(teks)` | `capitalize`, `uppercase` | Mengubah seluruh tulisan menjadi huruf besar semua |
| **Duplikasi Data** | `salin(objek)` | `clone` | Membuat salinan data mandiri agar perubahan tidak merusak data asli |

---

## 1. Menampilkan ke Layar & Menerima Jawaban

### `cetak(...nilai)` / `print`
Menampilkan tulisan atau angka ke layar. Jika Anda memasukkan beberapa hal sekaligus dengan koma, Celang akan memisahkannya dengan spasi:

```celang
cetak("Halo Dunia!")
cetak("Nama:", "Budi", "| Usia:", 25)
// Hasil di layar: Nama: Budi | Usia: 25
```

### `masukan(pertanyaan?)` / `input`
Mengajukan pertanyaan di layar dan menunggu sampai pengguna mengetik jawaban lalu menekan tombol Enter:

```celang
buat nama = masukan("Siapa namamu? ")
cetak("Senang berkenalan denganmu,", nama)
```

---

## 2. Memeriksa Isi & Panjang Data

### `panjang(wadah)` / `len`
Menghitung banyak karakter dalam tulisan, banyak barang dalam daftar, atau banyak kolom dalam kamus:

```celang
cetak(panjang("Celang"))         // 6 (ada 6 huruf)
cetak(panjang([10, 20, 30]))     // 3 (ada 3 angka)
cetak(panjang({"a": 1, "b": 2})) // 2 (ada 2 pasang data)
```

### `tipe(nilai)` / `type`
Mengetahui apa jenis data dari suatu nilai:

```celang
cetak(tipe(100))        // "Angka"
cetak(tipe(3.14))       // "Pecahan"
cetak(tipe("Halo"))     // "Teks"
cetak(tipe(benar))      // "Logika"
cetak(tipe([1, 2]))     // "Array"
cetak(tipe({"x": 10}))  // "Map"
```

---

## 3. Mengubah Jenis Data

### `Angka(nilai)`
Mengubah tulisan angka menjadi angka sungguhan agar bisa dihitung matematika:
```celang
buat x = Angka("125000") // angka 125000
buat y = Angka(benar)    // angka 1
```

### `Teks(nilai)`
Mengubah nilai apa pun menjadi tulisan teks:
```celang
buat s1 = Teks(42)       // "42"
buat s2 = Teks(benar)    // "benar"
```

### `Logika(nilai)`
Memeriksa apakah data tersebut bernilai ada atau kosong:
- Angka `0`, tulisan kosong `""`, daftar kosong `[]`, atau nilai `kosong` bernilai `salah`.
- Data yang berisi bernilai `benar`.

```celang
cetak(Logika(1))        // benar
cetak(Logika(0))        // salah
cetak(Logika("Celang")) // benar
cetak(Logika(""))       // salah
```

---

## 4. Mengolah Kumpulan Daftar

### `rentang(awal?, akhir)` / `range`
Membuat deretan angka berurutan secara cepat:
```celang
buat r1 = rentang(4)    // [0, 1, 2, 3, 4]
buat r2 = rentang(1, 5) // [1, 2, 3, 4, 5]
buat r3 = rentang(3, 1) // [3, 2, 1]
```

### `saring(daftar, rumus)` / `filter`
Menyaring barang di dalam daftar, mengambil hanya yang memenuhi syarat:
```celang
buat angka = [1, 2, 3, 4, 5, 6]
buat angkaGenap = saring(angka, x => x % 2 == 0)
cetak(angkaGenap) // [2, 4, 6]
```

### `petakan(daftar, rumus)` / `map`
Mengubah semua isi daftar menggunakan rumus yang sama:
```celang
buat angka = [1, 2, 3]
buat kuadrat = petakan(angka, x => x * x)
cetak(kuadrat) // [1, 4, 9]
```

---

## 5. Memoles Teks

### `potong_spasi(teks)` / `trim`
Membuang spasi kosong di awal atau akhir tulisan:
```celang
buat pesan = "   Halo Celang!   "
cetak(potong_spasi(pesan)) // "Halo Celang!"
```

### `kapitalisasi(teks)` / `uppercase`
Mengubah seluruh tulisan menjadi huruf besar:
```celang
cetak(kapitalisasi("selamat pagi")) // "SELAMAT PAGI"
```

---

## 6. Menduplikasi Data Secara Aman (`salin`)

Membuat salinan utuh dari data yang kompleks sehingga perubahan pada salinan tidak akan mengubah data aslinya:

```celang
buat asli = [1, 2, 3]
buat tiruan = salin(asli)

tiruan.tambah(4)
cetak(asli)   // [1, 2, 3] (tetap aman!)
cetak(tiruan) // [1, 2, 3, 4]
```
