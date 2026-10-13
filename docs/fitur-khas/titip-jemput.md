# 🎟️ Titip & Jemput (Non-blocking Asynchronous)

**Titip & Jemput** adalah terobosan Celang untuk pemrograman asinkron (*asynchronous*) yang menggantikan kerumitan konsep `Promise`, `async`, dan `await` di JavaScript.

---

## 🧺 Analogi Nyata: Karcis Laundry

Bayangkan Anda membawa pakaian kotor ke tempat laundry:
1. **Titip (`titip`)**: Anda menyerahkan pakaian ke kasir. Kasir langsung memberikan **Karcis/Tiket Bukti Titip**. Anda tidak berdiri melongo menunggu mesin cuci selesai berputar berjam-jam; Anda langsung bebas pergi melakukan hal lain.
2. **Jemput (`jemput`)**: Nanti saat Anda butuh pakaian bersih tersebut, Anda datang membawa karcis untuk menjemput cucian. Jika sudah selesai, Anda menerima pakaian bersihnya. Jika mesin cuci rusak, kasir mengabari Anda (galat/error).

```
[ Program Utama ]
       │
       ├─► titip hitungData()  ──► Dapat <Tiket> langsung!
       │                             │
       │                             ▼
       ├─► Lanjut kerja lain...   [ Goroutine latar belakang berjalan... ]
       │                             │
       ▼                             ▼
   jemput tiket  ◄────────────── [ Selesai! Mengembalikan hasil ]
```

---

## 1. Perintah Dasar: `titip` dan `jemput`

### Menitipkan Pekerjaan ke Latar Belakang
Gunakan kata kunci `titip` di depan ekspresi atau fungsi yang ingin dijalankan di latar belakang:

```celang
fungsi unduhLaporan() {
    waktu.tidur(1000) // pura-pura proses lama 1 detik
    kembalikan "Laporan Tahunan 2026.pdf"
}

// 1. Titip pekerjaan -> langsung dapat tiket tanpa berhenti menunggu
buat tiket = titip unduhLaporan()
cetak("Tiket berhasil didapat, proses berjalan di latar belakang!")

// 2. Jemput hasilnya ketika dibutuhkan
buat dokumen = jemput tiket
cetak("Dokumen diterima:", dokumen)
```

---

## 2. Batas Waktu Tunggu (*Timeout*)

Tidak ingin menunggu selamanya jika proses latar belakang macet? Gunakan klausa `dalam <waktu>`:

```celang
buat tiketProses = titip prosesKomputasiBerat()

// Menunggu maksimal 3 detik
buat hasil = jemput tiketProses dalam 3.detik
```

> [!NOTE]
> Anda bisa menulis durasi dengan `3.detik`, `500.milidetik`, atau angka langsung seperti `jemput tiket dalam 3`. Jika waktu habis sebelum tugas selesai, Celang akan melempar pesan kesalahan yang ramah.

---

## 3. Kombinator Lanjutan: `jemput_semua` & `jemput_pertama`

Ketika Anda memiliki beberapa tugas yang dititipkan sekaligus, Celang menyediakan dua fungsi bawaan yang sangat intuitif:

### `jemput_semua([t1, t2, ...])` (Mirip `Promise.all`)
Menunggu semua tiket selesai dikerjakan secara paralel, lalu mengembalikan sebuah daftar (*Array*) berisi seluruh hasilnya:

```celang
buat t1 = titip ambilDataPengguna()
buat t2 = titip ambilDataProduk()
buat t3 = titip ambilDataDiskon()

// Menunggu ketiganya selesai secara bersamaan:
buat [pengguna, produk, diskon] = jemput_semua([t1, t2, t3])

cetak("Semua data berhasil dimuat!")
```

### `jemput_pertama([t1, t2, ...])` (Mirip `Promise.race`)
Mengambil hasil dari tiket mana pun yang **paling cepat selesai**:

```celang
buat serverA = titip panggilServerJakarta()
buat serverB = titip panggilServerSingapura()

// Ambil respon tercepat:
buat responTercepat = jemput_pertama([serverA, serverB])
cetak("Respon dari server tercepat:", responTercepat)
```

---

## 4. Penanganan Galat (*Error Handling*)

Jika tugas yang dititipkan mengalami kegagalan, galat tersebut tidak akan hilang atau membungkam program diam-diam. Galat akan otomatis dilempar saat Anda memanggil `jemput`, sehingga dapat ditangkap secara elegan menggunakan blok `coba ... tangkap`:

```celang
fungsi simpanKeServer() {
    lempar "Koneksi internet terputus!"
}

buat tiket = titip simpanKeServer()

coba {
    buat hasil = jemput tiket
    cetak("Sukses:", hasil)
} tangkap (kesalahan) {
    cetak("Peringatan:", kesalahan)
    // Output: Peringatan: Koneksi internet terputus!
}
```

---

## 🎯 Perbandingan dengan Bahasa Lain

| Konsep | JavaScript | Go | Celang |
| :--- | :--- | :--- | :--- |
| **Pemicu Async** | `new Promise(...)` / `async` | `go routine()` | `titip ekspresi` |
| **Menunggu Hasil** | `await promise` | `<-channel` | `jemput tiket` |
| **Tunggu Semua** | `Promise.all(...)` | `sync.WaitGroup` | `jemput_semua(...)` |
| **Tunggu Tercepat**| `Promise.race(...)` | `select` | `jemput_pertama(...)` |
| **Batas Waktu** | `AbortController` | `context.WithTimeout` | `jemput tiket dalam 3.detik` |
