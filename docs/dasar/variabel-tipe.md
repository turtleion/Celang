# Variabel & Jenis Data

Di Celang, menyimpan dan mengelola data dibuat sealami mungkin, seperti Anda mencatat hal-hal penting di buku catatan.

---

## 1. Membuat Catatan Nilai (`buat`)

Gunakan kata kunci `buat` untuk menyiapkan sebuah wadah catatan nilai yang isinya bisa diubah di kemudian hari:

```celang
buat nama = "Budi Santoso"
buat umur = 24
buat aktif = benar

// Mengubah isi variabel dengan nilai baru:
umur = 25
cetak("Umur sekarang:", umur) // 25
```

> **Nilai Awal Kosong:** Jika Anda membuat variabel tanpa langsung mengisinya, Celang secara otomatis mengisinya dengan nilai `kosong` (null):
> ```celang
> buat alamat
> cetak(alamat) // kosong
> ```

---

## 2. Nilai Tetap yang Tidak Boleh Berubah (`konstanta`)

Jika ada nilai yang harus selalu tetap dan tidak boleh diubah secara tidak sengaja oleh bagian program lain, gunakan kata kunci `konstanta`:

```celang
konstanta PAJAK = 0.11
konstanta NAMA_APLIKASI = "TokoKita"

// Jika Anda mencoba mengubah nilai konstanta:
// PAJAK = 0.15 
// ❌ Celang akan menolak dengan ramah: nilai konstanta tidak boleh diubah!
```

---

## 3. Batas Lingkup Catatan (`{ ... }`)

Variabel yang Anda buat di dalam kurung kurawal `{ ... }` hanya akan berlaku di dalam ruangan tersebut. Hal ini menjaga agar nama variabel tidak saling bertabrakan:

```celang
buat judul = "Buku Utama"

{
    // Catatan di dalam kurung kurawal hanya dikenal di sini:
    buat catatanKecil = "Catatan Bab 1"
    cetak(judul)        // "Buku Utama" (bisa membaca dari luar)
    cetak(catatanKecil) // "Catatan Bab 1"
}

// Di luar kurung kurawal, 'catatanKecil' sudah tidak ada:
// cetak(catatanKecil) ❌ (Variabel tidak ditemukan)
cetak(judul) // "Buku Utama" (tetap aman)
```

---

## 4. Jenis-Jenis Nilai Dasar

Celang menyediakan beberapa jenis nilai yang sering kita jumpai dalam kehidupan sehari-hari:

| Jenis Nilai | Nama Jenis | Contoh | Keterangan |
| :--- | :--- | :--- | :--- |
| **Bilangan Bulat** | `Angka` | `42`, `-10`, `1_000_000` | Angka utuh tanpa koma (tanda garis bawah `_` boleh dipakai sebagai pemisah ribuan agar mudah dibaca) |
| **Bilangan Desimal** | `Angka` / `Pecahan` | `3.14`, `-0.05` | Angka dengan tanda titik sebagai koma desimal |
| **Tulisan / Teks** | `Teks` | `"Halo"`, `'Dunia'` | Tulisan huruf, kata, atau kalimat |
| **Pilihan Logika** | `Logika` | `benar`, `salah` | Jawaban ya atau tidak |
| **Nilai Kosong** | - | `kosong` (nihil) | Menandakan belum ada data atau tidak ada nilai |

---

## 5. Kumpulan Data: Daftar dan Kamus

### A. Daftar / Kumpulan Berurutan (`[...]`)
Daftar digunakan untuk menyimpan banyak barang secara berurutan dalam satu wadah:

```celang
buat buah = ["Apel", "Mangga", "Jeruk", "Pisang"]

// Mengambil barang berdasarkan urutan nomor (dimulai dari nomor 0):
cetak("Buah pertama:", buah[0]) // "Apel"

// Menghitung berapa banyak buah di dalam keranjang:
cetak("Total buah  :", panjang(buah)) // 4

// Mengambil potongan tertentu:
buat pilihan = buah[1..3] // ["Mangga", "Jeruk"]
```

### B. Kamus Data Bertanda Label (`{...}`)
Kamus digunakan ketika setiap data memiliki nama label tersendiri (seperti data formulir):

```celang
buat pengguna = {
    "nama": "Siti Rahma",
    "divisi": "Teknologi",
    "gaji": 15000000,
    "aktif": benar
}

// Membaca data menggunakan tanda titik atau kurung siku:
cetak(pengguna.nama)      // "Siti Rahma"
cetak(pengguna["divisi"]) // "Teknologi"

// Jika mencari kolom yang belum ada, hasilnya adalah 'kosong' tanpa membuat program macet:
cetak(pengguna.alamat)    // kosong
```

---

## 6. Mengubah Jenis Data (Konversi)

Jika Anda ingin mengubah tulisan angka `"150"` menjadi angka sungguhan, atau sebaliknya, Celang menyediakan alat pengubah jenis yang sangat mudah:

```celang
buat teksAngka = "25000"
buat angkaMurni = Angka(teksAngka) // Mengubah teks menjadi angka 25000

buat total = 1000
buat kalimat = Teks(total)         // Mengubah angka 1000 menjadi teks "1000"
```
Celang juga secara cerdas mengenali tulisan berformat mata uang (misalnya `"Rp 150000" - 25000` akan otomatis menghasilkan angka `125000`).
