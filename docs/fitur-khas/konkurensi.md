# Konkurensi Alami (`sambil` & `bersamaan`)

Celang dibangun di atas runtime Go yang terkenal dengan keandalan goroutine dan primitif konkurensinya. Namun, alih-alih memaksa pengembang mengelola `sync.WaitGroup`, channel kompleks, atau syntax threading yang rumit, Celang menghadirkan **Natural Concurrency** menggunakan kata kunci alami bahasa Indonesia: **`sambil`** dan **`bersamaan`**.

---

## 1. Asinkron di Latar Belakang: `sambil`

Kata kunci `sambil` memungkinkan Anda mengeksekusi blok kode secara asinkron di latar belakang tanpa memblokir baris kode utama berikutnya.

Secara teknis, `sambil` memunculkan Goroutine terpisah yang berjalan secara non-blocking.

### Sintaks

```celang
sambil {
    // Pekerjaan berat atau tugas latar belakang
    // Dijalankan secara independen
}

// Baris ini langsung dieksekusi tanpa menunggu blok sambil selesai
```

### Contoh Nyata

```celang
cetak("Memulai proses...")

sambil {
    cetak("[Latar] Mengunduh data pembaruan...")
    waktu.tidur(100)
    cetak("[Latar] Pembaruan berhasil diunduh!")
}

cetak("Program utama melanjutkan interaksi dengan pengguna...")
waktu.tidur(150)
cetak("Selesai!")
```

### Output

```text
Memulai proses...
Program utama melanjutkan interaksi dengan pengguna...
[Latar] Mengunduh data pembaruan...
[Latar] Pembaruan berhasil diunduh!
Selesai!
```

---

## 2. Paralelisme Terkoordinasi: `bersamaan`

Jika `sambil` digunakan untuk tugas latar belakang yang tidak perlu ditunggu hasilnya, `bersamaan` digunakan ketika Anda memiliki beberapa tugas yang ingin diproses **secara paralel** sekaligus, dan Anda **harus menunggu semuanya selesai** sebelum melanjutkan ke langkah berikutnya.

Di balik layar, Celang secara otomatis mengelola `sync.WaitGroup` untuk Anda: menambahkan counter untuk tiap proses dan menahan eksekusi utama sampai semua cabang selesai.

### Sintaks

```celang
bersamaan {
    tugasA()
    tugasB()
    tugasC()
}

// Blok ini baru akan dieksekusi setelah tugasA, tugasB, dan tugasC semuanya selesai
cetak("Seluruh tugas selesai!")
```

### Contoh Kasus: Memproses Berkas Multitasking

```celang
buat daftarLaporan = ["laporan_januari.csv", "laporan_februari.csv", "laporan_maret.csv"]

fungsi olahLaporan(namaBerkas) {
    cetak("Mulai memproses:", namaBerkas)
    waktu.tidur(50)
    cetak("Selesai memproses:", namaBerkas)
}

cetak("--- Mulai Pemrosesan Paralel ---")

bersamaan {
    olahLaporan("laporan_januari.csv")
    olahLaporan("laporan_februari.csv")
    olahLaporan("laporan_maret.csv")
}

cetak("--- Semua 3 berkas laporan telah selesai diproses! Melanjutkan ekspor data... ---")
```

---

## 3. Keamanan Memori & Variabel Bersama

Ketika beberapa thread atau goroutine mengakses variabel yang sama, Celang menyediakan lingkungan aman. Namun, disarankan mengikuti prinsip Go:
> *"Jangan berkomunikasi dengan berbagi memori; alih-alih, bagilah memori dengan berkomunikasi."*

Dalam Celang, jika Anda mengakses variabel time-travel di dalam thread konkurensi:
- Tiap frame pembaruan state dicatat secara atomik.
- Hindari *race condition* penulisan variabel bersama tanpa sinkronisasi terarah.

---

## Perbandingan Dialek

::: code-group

```celang [Bahasa Indonesia]
sambil {
    unduhData()
}

bersamaan {
    prosesA()
    prosesB()
}
```

```celang [JavaScript Dialect]
async {
    unduhData()
}

parallel {
    prosesA()
    prosesB()
}
```

```celang [Sunda Dialect]
bari {
    unduhData()
}

babarengan {
    prosesA()
    prosesB()
}
```

:::
