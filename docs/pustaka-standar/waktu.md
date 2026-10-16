# Pustaka Penunjuk Waktu (`waktu` / `time`)

Modul **`waktu`** (dapat juga dipanggil dengan alias internasional **`time`**) menyediakan ekosistem lengkap untuk membaca, memanipulasi, mengukur, dan memformat penanggalan waktu di Celang.

Selain operasi stempel waktu dasar, modul `waktu` di Celang dirancang dengan **fitur-fitur cerdas** yang sangat ramah pengembang (*human-centric*):
- **Waktu Relatif Bahasa Indonesia (`waktu.lalu`)**: Menghasilkan teks waktu relatif otomatis (*"baru saja"*, *"5 menit yang lalu"*, *"2 jam yang lalu"*, *"kemarin"*, *"3 hari yang lalu"*). Sangat dibutuhkan untuk aplikasi media sosial, komentar, riwayat aktivitas, dan notifikasi.
- **Format Penanggalan Indonesia Resmi (`waktu.format_indonesia`)**: Memformat tanggal dengan nama hari dan bulan bahasa Indonesia baku sesuai PUEBI/KBBI (*"Rabu, 15 Oktober 2026, 17:30 WIB"*).
- **Pola Format Intuitif Modern**: Mendukung format token standar industri (`YYYY-MM-DD HH:mm:ss`, `DD/MM/YYYY`, dll.) tanpa perlu menghafal tata letak penanggalan yang rumit.
- **Pengukur Durasi Eksekusi Otomatis (`waktu.ukur` / `benchmark`)**: Mengukur berapa lama fungsi atau algoritma berjalan secara presisi hingga ke pecahan milidetik/mikrodetik.
- **Pembongkar Komponen Penanggalan (`waktu.info`)**: Mengekstrak tahun, bulan, hari, nama hari, jam, menit, hingga zona waktu ke dalam bentuk Kamus (`Map`).
- **Aritmatika Penanggalan Praktis**: Menambah atau mengurangi tanggal menggunakan unit intuitif (`{"hari": 3, "jam": 5}`) serta menghitung selisih rentang waktu (`waktu.selisih`).
- **Pengurai String Tanggal Fleksibel (`waktu.urai` / `parse`)**: Otomatis mendeteksi string ISO 8601, RFC3339, format tanggal database SQL, maupun pola kustom.
- **Generator & Format Durasi**: Konverter cepat `waktu.detik(5)`, `waktu.menit(10)`, `waktu.jam(2)`, serta pengubah milidetik ke kalimat durasi (`waktu.format_durasi`).

---

## Cara Memanggil

Modul ini aktif secara global dan dapat dipanggil menggunakan nama lokal maupun alias internasional:

```celang
// 1. Menggunakan nama Bahasa Indonesia:
buat sekarang = waktu.sekarang()
buat teksIndo = waktu.format_indonesia(sekarang)
cetak(teksIndo) // "Rabu, 15 Oktober 2026, 17:30 WIB"

// 2. Menggunakan alias internasional:
buat ms = time.now()
time.sleep(100) // istirahat 100 milidetik

// 3. Atau impor fungsi tertentu (tanpa tanda kutip):
ambil { lalu, format_indonesia, ukur, tambah, selisih } dari waktu
```

---

## 🌟 Fitur Cerdas & Human-Centric Celang

### 1. Waktu Relatif Humanis (`lalu` / `relatif` / `time_ago`)
Mengubah selisih waktu menjadi kalimat percakapan bahasa Indonesia yang alami:

```celang
buat t_sekarang = waktu.sekarang()

// Simulasi postingan di media sosial:
buat t_komentar = t_sekarang - waktu.menit(15)
cetak(waktu.lalu(t_komentar)) // "15 menit yang lalu"

buat t_kemarin = t_sekarang - waktu.hari(1)
cetak(waktu.lalu(t_kemarin))  // "kemarin"

buat t_baru = t_sekarang - waktu.detik(10)
cetak(waktu.lalu(t_baru))     // "baru saja"

// Mendukung juga waktu di masa depan:
buat t_jadwal = t_sekarang + waktu.jam(2)
cetak(waktu.lalu(t_jadwal))   // "2 jam lagi"
```

---

### 2. Format Penanggalan Indonesia Resmi (`format_indonesia` / `format_id`)
Menghasilkan penanggalan resmi berbahasa Indonesia lengkap dengan zona waktu lokal (WIB / WITA / WIT):

```celang
buat saatIni = waktu.sekarang()

// 1. Gaya Lengkap (default):
cetak(waktu.format_indonesia(saatIni))
// Hasil: "Rabu, 15 Oktober 2026, 17:30 WIB"

// 2. Gaya Tanggal Saja:
cetak(waktu.format_indonesia(saatIni, "tanggal"))
// Hasil: "15 Oktober 2026"

// 3. Gaya Hari dan Tanggal:
cetak(waktu.format_indonesia(saatIni, "hari_tanggal"))
// Hasil: "Rabu, 15 Oktober 2026"

// 4. Gaya Jam dan Waktu:
cetak(waktu.format_indonesia(saatIni, "waktu"))
// Hasil: "17:30 WIB"
```

---

### 3. Format Penanggalan Intuitif (`format`)
Mendukung pola token modern (`YYYY`, `MM`, `DD`, `HH`, `mm`, `ss`):

```celang
buat t = waktu.sekarang()

// Format standar database SQL:
cetak(waktu.format(t, "YYYY-MM-DD HH:mm:ss"))
// Hasil: "2026-10-15 17:30:00"

// Format tanggal garis miring:
cetak(waktu.format(t, "DD/MM/YYYY"))
// Hasil: "15/10/2026"

// Format jam dan menit:
cetak(waktu.format(t, "HH:mm"))
// Hasil: "17:30"
```

---

### 4. Pengukur Durasi Eksekusi / Profiler (`ukur` / `benchmark`)
Mengukur performa sebuah fungsi tanpa perlu mencatat `start` dan `stop` manual:

```celang
buat profil = waktu.ukur(fungsi() {
    buat total = 0
    untuk i dalam rentang(1, 100000) {
        total = total + i
    }
    kembalikan total
})

cetak("Hasil komputasi :", profil.hasil)
cetak("Waktu eksekusi  :", profil.durasi_teks) // misal: "1.45 ms"
cetak("Angka milidetik :", profil.durasi_ms)   // 1.45
```

---

## 🗓️ Rincian & Komponen Waktu (`info` / `rincian`)

Mendapatkan seluruh informasi tanggal dan jam dalam bentuk Kamus (`Map`) tanpa perlu parsing teks:

```celang
buat r = waktu.info()

cetak("Tahun             :", r.tahun)             // 2026
cetak("Bulan (angka)     :", r.bulan)             // 10
cetak("Bulan (nama)      :", r.nama_bulan)        // "Oktober"
cetak("Tanggal           :", r.hari)              // 15
cetak("Hari (nama)       :", r.nama_hari)         // "Rabu"
cetak("Hari dalam minggu :", r.hari_dalam_minggu) // 3
cetak("Jam               :", r.jam)               // 17
cetak("Menit             :", r.menit)             // 30
cetak("Detik             :", r.detik)             // 45
cetak("Zona Waktu        :", r.zona_waktu)        // "WIB"
```

Tersedia juga fungsi pembaca cepat:
- `waktu.tahun([ms])`
- `waktu.bulan([ms])`
- `waktu.hari_ke([ms])`

---

## ⏳ Aritmatika Penanggalan & Manipulasi

### 1. Menambah dan Mengurangi Waktu (`tambah` & `kurang`)
Gunakan kamus penanggalan yang intuitif:

```celang
buat sekarang = waktu.sekarang()

// Tambah 3 hari dan 2 jam:
buat tenggatWaktu = waktu.tambah(sekarang, { "hari": 3, "jam": 2 })
cetak("Batas akhir:", waktu.format(tenggatWaktu, "YYYY-MM-DD HH:mm"))

// Kurang 1 minggu (7 hari):
buat mingguLalu = waktu.kurang(sekarang, { "hari": 7 })
cetak("Satu minggu lalu:", waktu.format_indonesia(mingguLalu, "tanggal"))
```

### 2. Menghitung Rentang Selisih Waktu (`selisih` / `diff`)
Mencari jarak selisih antara dua waktu:

```celang
buat mulai = waktu.urai("2026-10-15 08:00:00")
buat selesai = waktu.urai("2026-10-17 11:30:15")

buat jarak = waktu.selisih(mulai, selesai)
cetak("Selisih:", jarak.hari, "hari,", jarak.jam, "jam,", jarak.menit, "menit")
// "Selisih: 2 hari, 3 jam, 30 menit"
```

### 3. Awal & Akhir Hari (`awal_hari` & `akhir_hari`)
Reset waktu ke jam `00:00:00.000` atau `23:59:59.999`:

```celang
buat hariIniMulai = waktu.awal_hari()
buat hariIniSelesai = waktu.akhir_hari()

cetak("Awal hari :", waktu.format(hariIniMulai, "YYYY-MM-DD HH:mm:ss"))
cetak("Akhir hari:", waktu.format(hariIniSelesai, "YYYY-MM-DD HH:mm:ss"))
```

### 4. Pemeriksaan Tahun Kabisat (`apakah_kabisat` / `is_leap_year`)
```celang
cetak(waktu.apakah_kabisat(2024)) // benar
cetak(waktu.apakah_kabisat(2025)) // salah
```

---

## 🔍 Pengurai String ke Waktu (`urai` / `parse`)

Membaca string penanggalan menjadi angka stempel waktu milidetik:

```celang
// 1. Deteksi otomatis format umum (ISO 8601, SQL, Tanggal biasa):
buat ms1 = waktu.urai("2026-10-15 17:30:00")
buat ms2 = waktu.urai("2026-10-15T10:30:00Z")
buat ms3 = waktu.urai("15/10/2026")

// 2. Menggunakan pola kustom eksplisit:
buat ms4 = waktu.urai("15-10-2026 17.30", "DD-MM-YYYY HH.mm")
```

---

## ⏱️ Generator & Format Durasi

```celang
// 1. Pembuat angka milidetik:
buat jeda = waktu.detik(5)    // 5000 ms
buat dur1 = waktu.menit(10)   // 600000 ms
buat dur2 = waktu.jam(2)      // 7200000 ms
buat dur3 = waktu.hari(1)     // 86400000 ms

// 2. Format milidetik ke kalimat durasi (format_durasi):
buat totalMs = waktu.jam(2) + waktu.menit(15) + waktu.detik(30)
cetak(waktu.format_durasi(totalMs))
// Hasil: "2 jam 15 menit 30 detik"
```

---

## 📋 Tabel Rangkuman Fungsi

| Fungsi Bahasa Indonesia | Alias Internasional | Keterangan Singkat |
| :--- | :--- | :--- |
| **`sekarang()`** | `now()` | Stempel waktu saat ini dalam milidetik |
| **`sekarang_detik()`** | `now_seconds()` | Stempel waktu saat ini dalam detik |
| **`sekarang_nano()`** | `now_nanos()` | Stempel waktu saat ini dalam nanodetik |
| **`tidur(ms)`** | `sleep(ms)` | Menjeda eksekusi program selama $N$ milidetik |
| **`format([ms], [pola])`** | `format` | Format waktu dengan token modern (`YYYY-MM-DD`) |
| **`format_indonesia([ms], [gaya])`** | `format_id` | Format penanggalan resmi Indonesia (*PUEBI*) |
| **`lalu(target_ms, [base_ms])`** | `relatif` / `time_ago` | Format waktu relatif (*"5 menit yang lalu"*) |
| **`ukur(fungsi)`** | `benchmark` | Profiler pengukur waktu eksekusi kode |
| **`info([ms])`** | `rincian` | Kamus rincian komponen waktu lengkap |
| **`tahun([ms])`** | `year` | Mengambil angka tahun |
| **`bulan([ms])`** | `month` | Mengambil angka bulan (1–12) |
| **`hari_ke([ms])`** | `day_of_month` | Mengambil tanggal hari dalam bulan (1–31) |
| **`tambah(ms, durasi)`** | `add` | Menambah waktu dengan kamus durasi |
| **`kurang(ms, durasi)`** | `sub` | Mengurangi waktu dengan kamus durasi |
| **`selisih(msA, msB)`** | `diff` | Menghitung jarak selisih dua waktu |
| **`awal_hari([ms])`** | `start_of_day` | Reset waktu ke `00:00:00.000` |
| **`akhir_hari([ms])`** | `end_of_day` | Reset waktu ke `23:59:59.999` |
| **`apakah_kabisat(tahun)`**| `is_leap_year` | Memeriksa tahun kabisat (boolean) |
| **`urai(teks, [pola])`** | `parse` | Mengurai teks menjadi stempel waktu angka |
| **`format_durasi(ms)`** | `format_duration` | Mengubah milidetik ke kalimat durasi |
| **`detik(n)`** | `seconds(n)` | Mengalikan $n$ ke satuan milidetik ($n \times 1.000$) |
| **`menit(n)`** | `minutes(n)` | Mengalikan $n$ ke satuan milidetik ($n \times 60.000$) |
| **`jam(n)`** | `hours(n)` | Mengalikan $n$ ke satuan milidetik ($n \times 3.600.000$) |
| **`hari(n)`** | `days(n)` | Mengalikan $n$ ke satuan milidetik ($n \times 86.400.000$) |
| **`minggu(n)`** | `weeks(n)` | Mengalikan $n$ ke satuan milidetik ($n \times 604.800.000$) |
