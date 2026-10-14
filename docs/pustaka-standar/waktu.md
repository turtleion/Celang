# Alat Bantu Penunjuk Waktu (`waktu` / `time`)

Modul `waktu` (bisa juga dipanggil dengan nama `time`) menyediakan alat untuk mengetahui jam dan tanggal saat ini, menghitung lama waktu yang dibutuhkan suatu proses, dan memberikan jeda sejenak (*istirahat/tidur*).

---

## Cara Memanggil

Modul ini siap dipakai secara langsung:

```celang
buat mulai = waktu.sekarang()
waktu.tidur(100) // menunggu 100 milidetik
buat selesai = waktu.sekarang()
cetak("Lama proses:", selesai - mulai, "milidetik")
```

---

## Daftar Alat Bantu

### 1. Mengetahui Waktu Saat Ini (`sekarang` / `now`)
Memberikan angka stempel waktu saat ini dalam satuan **milidetik** (1 detik = 1000 milidetik). Sangat cocok untuk mengukur seberapa cepat program Anda menyelesaikan suatu tugas:

```celang
buat ms = waktu.sekarang()
cetak("Stempel waktu milidetik:", ms)
```

Jika Anda hanya membutuhkan satuan **detik**, gunakan `waktu.sekarang_detik()`:
```celang
buat detik = waktu.sekarang_detik()
cetak("Stempel waktu detik:", detik)
```

---

### 2. Mengubah Waktu Menjadi Teks yang Mudah Dibaca (`format`)
Mengubah angka stempel waktu menjadi teks penanggalan yang rapi dan nyaman dibaca:

```celang
buat saatIni = waktu.sekarang()

// Format standar (Tahun-Bulan-Hari Jam:Menit:Detik):
cetak(waktu.format(saatIni)) // Contoh: 2026-10-14 00:30:00

// Hanya menampilkan jam dan menit:
cetak(waktu.format(saatIni, "15:04")) // Contoh: 00:30

// Format tanggal Indonesia (Hari/Bulan/Tahun):
cetak(waktu.format(saatIni, "02/01/2006")) // Contoh: 14/10/2026
```

---

### 3. Memberi Jeda Waktu (`tidur` / `sleep`)
Menghentikan sementara program Anda selama beberapa saat sebelum melanjutkan ke baris berikutnya:

```celang
cetak("Menghitung mundur...")
cetak("3...")
waktu.tidur(1000) // Berhenti sejenak selama 1 detik (1000 milidetik)
cetak("2...")
waktu.tidur(1000)
cetak("1...")
waktu.tidur(1000)
cetak("Mulai!")
```

### 4. Satuan Waktu Praktis
Selain menggunakan angka milidetik mentah, Celang juga mendukung penulisan satuan waktu langsung pada angka:
- `3.detik` (artinya 3000 milidetik)
- `500.milidetik` (artinya 500 milidetik)
- `1.menit` (artinya 60000 milidetik)
