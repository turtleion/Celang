# Membongkar Isi Paket Data (Destrukturisasi)

Bayangkan Anda baru pulang belanja dan membawa sekantong belanjaan berisi apel, jeruk, dan mangga. Daripada mengambil barang satu per satu dengan tangan kanan berkali-kali, Anda langsung membongkar kantong tersebut dan menaruh masing-masing buah ke mangkuknya masing-masing dalam satu gerakan.

Fitur **Destrukturisasi** memungkinkan Anda memecah isi daftar (*Array*) atau kamus data (*Map*) langsung ke dalam beberapa variabel terpisah dalam satu baris kode yang singkat.

---

## 📦 Membongkar Daftar Berdasarkan Urutan

Jika Anda memiliki daftar data berurutan, Anda dapat memecahnya sesuai urutan posisinya:

```celang
buat titikKoordinat = [120, 250, 480]

// Membongkar posisi pertama ke x, kedua ke y, dan ketiga ke z:
buat [x, y, z] = titikKoordinat

cetak("Titik X:", x) // 120
cetak("Titik Y:", y) // 250
cetak("Titik Z:", z) // 480
```

### Mengisi Ulang Variabel yang Sudah Ada
Anda juga dapat memperbarui variabel yang sudah ada tanpa perlu menulis kata `buat` lagi:

```celang
[x, y] = [999, 888]
cetak("X baru:", x, "Y baru:", y) // 999 888
```

---

## 🗂️ Membongkar Kamus Berdasarkan Nama Kolom

Jika Anda memiliki sekumpulan data yang memiliki nama label (seperti data kartu nama), Anda dapat mengambil bagian yang Anda butuhkan dengan menyebutkan nama labelnya:

```celang
buat dataKaryawan = {
    "nama": "Siti Rahma",
    "divisi": "Teknologi",
    "gaji": 15000000
}

// Mengambil nama dan divisi secara langsung:
buat { nama, divisi } = dataKaryawan

cetak(nama, "bekerja di divisi", divisi) // Siti Rahma bekerja di divisi Teknologi
```

### Menggunakan Nama Panggilan Baru
Jika nama kolom dari data asal terasa terlalu panjang atau bentrok dengan variabel lain, Anda bisa mengganti namanya saat membongkar:

```celang
buat { gaji: uangBulanan } = dataKaryawan
cetak("Uang bulanan:", uangBulanan) // 15000000
```
