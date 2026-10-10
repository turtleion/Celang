# Modul Waktu (`waktu` / `time`)

Modul `waktu` (alias `time`) menyediakan kapabilitas untuk mengukur durasi, mengambil waktu saat ini, memformat penanggalan, dan menunda eksekusi (*sleep*).

---

## Memuat Modul

Tersedia langsung di seluruh skrip Celang:

```celang
buat mulai = waktu.sekarang()
waktu.tidur(100)
buat selesai = time.now()
cetak("Waktu eksekusi:", selesai - mulai, "ms")
```

---

## Daftar Fungsi

### 1. `sekarang()` / `now()`
Mengembalikan stempel waktu (*timestamp*) saat ini dalam satuan **milidetik** (*epoch milliseconds*). Sangat cocok untuk *benchmarking* dan pengukuran performa.

```celang
buat ms = waktu.sekarang()
cetak("Milidetik saat ini:", ms)
```

### 2. `sekarang_detik()` / `now_seconds()`
Mengembalikan stempel waktu Unix standar dalam satuan **detik** (*epoch seconds*).

```celang
buat detik = waktu.sekarang_detik()
cetak("Unix timestamp detik:", detik)
```

### 3. `format(timestamp_ms, layout?)`
Memformat stempel milidetik ke representasi tanggal dan waktu yang dapat dibaca manusia.
- Format bawaan jika `layout` tidak diberikan adalah `"2006-01-02 15:04:05"`.
- Layout mengikuti standar Go reference time (`2006` tahun, `01` bulan, `02` hari, `15` jam, `04` menit, `05` detik).

```celang
buat saatIni = waktu.sekarang()

// Format standar YYYY-MM-DD HH:MM:SS
cetak(waktu.format(saatIni)) // Contoh: 2026-10-11 04:30:00

// Format kustom (hanya jam dan menit)
cetak(waktu.format(saatIni, "15:04")) // Contoh: 04:30

// Format tanggal Indonesia
cetak(waktu.format(saatIni, "02/01/2006")) // Contoh: 11/10/2026
```

### 4. `tidur(milidetik)` / `sleep(milidetik)`
Menghentikan sementara eksekusi thread/program saat ini selama durasi milidetik yang ditentukan.

```celang
cetak("Menghitung mundur...")
cetak("3...")
waktu.tidur(1000) // Tunggu 1 detik (1000 ms)
cetak("2...")
waktu.tidur(1000)
cetak("1...")
waktu.tidur(1000)
cetak("Mulai!")
```
