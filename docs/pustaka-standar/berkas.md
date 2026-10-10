# Modul Berkas (`berkas` / `fs`)

Modul `berkas` (alias `fs`) menyediakan antarmuka I/O berkas dan sistem penyimpanan lokal yang aman dan mudah digunakan untuk membaca, menulis, menambahkan, memeriksa, serta menghapus berkas.

---

## Memuat Modul

Tersedia secara langsung sebagai pustaka bawaan:

```celang
// Menggunakan identifier Bahasa Indonesia
berkas.tulis("./laporan.txt", "Isi berkas...")

// Atau menggunakan identifier internasional
buat isi = fs.read("./laporan.txt")
```

---

## Daftar Fungsi

### 1. `baca(path)` / `read(path)`
Membaca seluruh isi berkas teks dari disk ke dalam variabel string.

```celang
buat isiCatatan = berkas.baca("./catatan.txt")
cetak("Isi berkas:", isiCatatan)
```

::: tip Pesan Error Empatik
Jika berkas tidak ditemukan, Celang memberikan pesan kesalahan yang mendetail beserta saran perbaikan alur path.
:::

### 2. `tulis(path, data)` / `write(path, data)`
Menulis data string ke berkas tujuan. Jika berkas belum ada, berkas baru akan dibuat. Jika direktori induk belum ada, direktori tersebut akan dibuat secara otomatis (*auto-mkdir*).

```celang
buat konten = "Konfigurasi Sistem Celang v1.0\nMode: Produksi\n"
berkas.tulis("./konfig/app.conf", konten)
cetak("Berkas konfigurasi berhasil ditulis!")
```

### 3. `tambah(path, data)` / `append(path, data)`
Menambahkan data teks di akhir berkas (*append mode*) tanpa menimpa isi yang sudah ada sebelumnya.

```celang
buat logBaru = "[{waktu.sekarang()}] Pengguna berhasil masuk.\n"
berkas.tambah("./logs/app.log", logBaru)
```

### 4. `ada(path)` / `exists(path)`
Memeriksa keberadaan suatu berkas atau direktori di disk. Mengembalikan `benar` jika ada, dan `salah` jika tidak ditemukan.

```celang
jika berkas.ada("./kunci_rahasia.key") {
    cetak("Kunci otentikasi ditemukan.")
} selain {
    cetak("Peringatan: Kunci otentikasi belum dibuat!")
}
```

### 5. `hapus(path)` / `remove(path)`
Menghapus berkas dari disk.

```celang
jika berkas.ada("./temp_file.txt") {
    berkas.hapus("./temp_file.txt")
    cetak("Berkas sementara berhasil dibersihkan.")
}
```
