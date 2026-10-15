# Tugas Latar Belakang (Titip & Jemput / Future & Promise)

Pernahkah Anda memesan makanan di kedai, lalu kasir memberi Anda nomor antrean atau tiket? Anda tidak perlu berdiri mematung di depan dapur; Anda bisa duduk santai membaca buku. Ketika pesanan selesai dipanggil, barulah Anda menukar nomor antrean tersebut untuk mengambil makanan Anda.

Inilah cara kerja **Titip & Jemput** di Celang. Fitur ini memungkinkan program menjalankan tugas yang memakan waktu lama di latar belakang tanpa membuat aktivitas utama terhenti.

---

## 🧺 Dua Langkah Sederhana

```
[ Program Utama ]
       │
       ├─► titip hitungLaporan() ──► Langsung mendapat <Tiket>!
       │                                │
       │                                ▼
       ├─► Melakukan hal lain...     [ Tugas dikerjakan di latar belakang ]
       │                                │
       ▼                                ▼
   jemput tiket  ◄─────────────────  [ Selesai! Mengembalikan hasil ]
```

1. **Titip (`titip`)**: Anda menyerahkan suatu pekerjaan untuk dikerjakan di latar belakang. Anda langsung menerima sebuah **Tiket** sebagai bukti.
2. **Jemput (`jemput`)**: Kapan pun Anda membutuhkan hasil dari pekerjaan tersebut, gunakan tiket tadi untuk mengambil hasilnya.

---

## 1. Cara Menggunakan Dasar

```celang
fungsi unduhLaporan() {
    waktu.tidur(1000) // pura-pura proses butuh 1 detik
    kembalikan "Laporan Keuangan 2026.pdf"
}

// 1. Titipkan pekerjaan -> program langsung lanjut tanpa terhenti
buat tiket = titip unduhLaporan()
cetak("Tiket diterima, proses sedang berjalan di belakang...")

// 2. Lakukan pekerjaan lain
cetak("Melayani pelanggan lain di meja kasir...")

// 3. Ambil hasil saat siap
buat dokumen = jemput tiket
cetak("Dokumen berhasil diterima:", dokumen)
```

---

## 2. Batas Waktu Tunggu

Jika Anda khawatir proses latar belakang macet atau memakan waktu terlalu lama, Anda dapat menetapkan batas waktu tunggu:

```celang
buat tiketProses = titip prosesKomputasiBerat()

// Menunggu maksimal 3 detik
buat hasil = jemput tiketProses dalam 3.detik
```

> **Petunjuk Satuan Waktu:** Anda bisa menulis durasi dengan `3.detik`, `500.milidetik`, atau angka langsung seperti `jemput tiket dalam 3`. Jika waktu habis sebelum tugas selesai, Celang akan memberikan pemberitahuan yang ramah.

---

## 3. Mengelola Banyak Tugas Sekaligus

Ketika Anda menitipkan beberapa pekerjaan secara bersamaan, Celang menyediakan dua cara mudah untuk mengambil hasilnya:

### Menunggu Semua Selesai (`jemput_semua`)
Menunggu semua tiket selesai dikerjakan, lalu mengumpulkan semua hasilnya ke dalam sebuah daftar:

```celang
buat t1 = titip ambilDataPengguna()
buat t2 = titip ambilDataProduk()
buat t3 = titip ambilDataDiskon()

// Menunggu ketiganya selesai:
buat [pengguna, produk, diskon] = jemput_semua([t1, t2, t3])

cetak("Semua data berhasil dimuat!")
```

### Mengambil yang Paling Cepat Selesai (`jemput_pertama`)
Jika Anda memanggil beberapa jalur cadangan dan hanya membutuhkan satu jawaban dari yang paling cepat selesai:

```celang
buat jalur1 = titip panggilServerJakarta()
buat jalur2 = titip panggilServerSingapura()

// Ambil jawaban dari jalur mana pun yang paling cepat merespons:
buat responCepat = jemput_pertama([jalur1, jalur2])
cetak("Jawaban dari server tercepat:", responCepat)
```

---

## 4. Bagaimana Jika Terjadi Masalah pada Tugas Latar Belakang?

Jika tugas yang dititipkan mengalami kegagalan, masalah tersebut tidak akan hilang atau membungkam program diam-diam. Masalah tersebut akan disampaikan saat Anda memanggil `jemput`, sehingga dapat Anda tangani dengan aman menggunakan blok `coba ... tangkap`:

```celang
fungsi simpanData() {
    lempar "Koneksi internet terputus!"
}

buat tiket = titip simpanData()

coba {
    buat hasil = jemput tiket
    cetak("Sukses:", hasil)
} tangkap (masalah) {
    cetak("Pemberitahuan:", masalah)
    // Tampilan: Pemberitahuan: Koneksi internet terputus!
}
```
