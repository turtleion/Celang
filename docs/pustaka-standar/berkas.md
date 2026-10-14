# Alat Bantu Pengelola Berkas (`berkas` / `fs`)

Modul `berkas` (bisa juga dipanggil dengan nama `fs`, singkatan dari *file system*) menyediakan alat untuk mengelola dokumen dan berkas di komputer Anda: membaca isi berkas, membuat catatan baru, menambahkan tulisan, memeriksa keberadaan berkas, dan menghapus berkas.

---

## Cara Memanggil

Modul ini sudah terpasang dan siap digunakan:

```celang
// Menggunakan nama Bahasa Indonesia:
berkas.tulis("./laporan.txt", "Catatan rapat hari ini...")

// Menggunakan nama internasional:
buat isi = fs.read("./laporan.txt")
cetak(isi)
```

---

## Daftar Alat Bantu

### 1. Membaca Dokumen (`baca` / `read`)
Membuka dan membaca seluruh isi dokumen teks di komputer Anda ke dalam variabel:

```celang
buat isiCatatan = berkas.baca("./catatan.txt")
cetak("Isi dokumen:", isiCatatan)
```

> **Petunjuk Ramah:** Jika berkas yang Anda cari ternyata tidak ada di folder tersebut, Celang tidak akan panik. Celang akan memberi tahu secara jelas nama berkas yang dicari dan menyarankan Anda memeriksa kembali lokasinya.

---

### 2. Menulis Berkas Baru (`tulis` / `write`)
Menyimpan tulisan ke dalam berkas di komputer. Jika berkas atau foldernya belum ada, Celang akan membuatnya secara otomatis:

```celang
buat teks = "Daftar Belanja:\n1. Beras\n2. Telur\n"
berkas.tulis("./catatan_belanja.txt", teks)
cetak("Berkas berhasil disimpan!")
```
*(Catatan: Perintah ini akan menggantikan seluruh isi berkas lama jika berkas tersebut sudah ada sebelumnya).*

---

### 3. Menambahkan Tulisan di Akhir Berkas (`tambah` / `append`)
Menambahkan baris baru di bagian paling bawah berkas tanpa menghapus tulisan yang sudah ada sebelumnya (sangat cocok untuk buku harian atau catatan riwayat):

```celang
buat barisBaru = "Catatan tambahan: Jangan lupa beli garam.\n"
berkas.tambah("./catatan_belanja.txt", barisBaru)
```

---

### 4. Memeriksa Apakah Berkas Ada (`ada` / `exists`)
Memeriksa apakah suatu berkas benar-benar ada di penyimpanan komputer. Memberikan jawaban `benar` jika ada, dan `salah` jika tidak ditemukan:

```celang
jika berkas.ada("./kunci.txt") {
    cetak("Berkas kunci ditemukan!")
} selain {
    cetak("Berkas kunci belum ada, silakan buat terlebih dahulu.")
}
```

---

### 5. Menghapus Berkas (`hapus` / `remove`)
Menghapus dokumen dari komputer saat sudah tidak dibutuhkan lagi:

```celang
jika berkas.ada("./sementara.txt") {
    berkas.hapus("./sementara.txt")
    cetak("Berkas sementara berhasil dibersihkan.")
}
```
