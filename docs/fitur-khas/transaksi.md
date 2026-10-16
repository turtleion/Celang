# 🛡️ Transaksi & Rollback Otomatis (`transaksi`)

Celang menghadirkan fitur revolusioner yang umumnya hanya ditemukan di database tingkat enterprise (seperti PostgreSQL) ke level bahasa pemrograman: **Transaksi Status & Rollback Otomatis**.

Dengan blok `transaksi`, Anda dapat menjalankan serangkaian mutasi variabel yang dijamin bersifat atomik: jika terjadi kegagalan atau dipanggil perintah `batalkan`, seluruh mutasi variabel dikembalikan (*rollback*) ke kondisi awal persis sebelum transaksi dimulai.

---

## 📌 Mengapa Fitur Ini Eksklusif?

Di sebagian besar bahasa pemrograman (JavaScript, Python, Go, Java):
- Jika Anda mengubah 3 variabel berturut-turut, lalu operasi ke-4 melempar error (*exception*), 3 variabel pertama **tetap berada dalam kondisi rusak/termodifikasi separuh jalan** (*dirty state*).
- Pengembang terpaksa menulis kode *cleanup* manual yang panjang, rawan bug, dan sulit diuji.

Di Celang, runtime secara otomatis membuat **Savepoint Ringan O(1)**. Jika ada kesalahan atau pembatalan eksplisit, semua mutasi dibatalkan otomatis!

---

## 🚀 Sintaks Dasar

```celang
transaksi {
    // Serangkaian mutasi variabel
    // Jika ada kegagalan, panggil:
    batalkan "alasan pembatalan"
} pulihkan (err) {
    // Blok penanganan kesalahan (opsional)
}
```

---

## 💡 Contoh 1: Transfer Saldo Antar Akun

Bayangkan transaksi keuangan di mana saldo pengirim dikurangi, lalu saldo penerima ditambah. Jika saldo pengirim tidak cukup atau ada gangguan jaringan, transaksi dibatalkan:

```celang
buat saldo_andi = 100.ribu
buat saldo_budi = 50.ribu
buat riwayat_status = "belum diproses"

transaksi {
    saldo_andi = saldo_andi - 30.ribu
    saldo_budi = saldo_budi + 30.ribu

    jika saldo_andi < 0 {
        batalkan "Saldo Andi tidak mencukupi"
    }
} pulihkan (alasan) {
    riwayat_status = "Gagal: " + alasan
}

cetak(saldo_andi.rupiah) // Rp 70.000
cetak(saldo_budi.rupiah) // Rp 80.000
```

Jika dicoba dengan nominal yang melebihi saldo:
```celang
transaksi {
    saldo_andi = saldo_andi - 200.ribu
    saldo_budi = saldo_budi + 200.ribu

    jika saldo_andi < 0 {
        batalkan "Saldo Andi minus!"
    }
} pulihkan (alasan) {
    riwayat_status = "Transaksi dibatalkan: " + alasan
}

// Saldo Andi dan Budi otomatis kembali ke awal persis!
cetak(saldo_andi.rupiah) // Rp 70.000
cetak(saldo_budi.rupiah) // Rp 80.000
cetak(riwayat_status)    // Transaksi dibatalkan: Saldo Andi minus!
```

---

## ⚡ Contoh 2: Rollback Otomatis Saat Terjadi Error / Eksepsi

Bahkan tanpa memanggil `batalkan`, jika terjadi *runtime exception* atau pemanggilan `lempar` (*throw*), rollback berlangsung otomatis:

```celang
buat data_penting = 999

transaksi {
    data_penting = 123
    lempar "Koneksi database terputus di tengah jalan"
} pulihkan (err) {
    cetak("Log error:", err)
}

// data_penting tidak pernah terkontaminasi oleh mutasi '123'
cetak(data_penting) // 999
```

---

## 🔄 Transaksi Bersarang (*Nested Transactions*)

Celang mendukung transaksi bertingkat dengan semantik *Savepoints*:
- Transaksi dalam (*inner*) yang dibatalkan hanya mengembalikan perubahan di dalam bloknya sendiri tanpa membatalkan transaksi luar.
- Jika transaksi luar (*outer*) dibatalkan, semua mutasi di seluruh rantai transaksi akan dikembalikan ke kondisi awal.

```celang
buat x = 10
buat y = 20

transaksi {
    x = 100
    transaksi {
        y = 200
        batalkan "hanya batalkan inner"
    }
}

cetak(x) // 100 (tetap berkomitmen)
cetak(y) // 20  (kembali ke nilai semula)
```
