# 🔄 Saluran (Channels & CSP Concurrency)

**Saluran** (*Channel*) adalah pipa komunikasi antar-proses bergaya **CSP (Communicating Sequential Processes)** seperti yang dipopulerkan oleh bahasa Go. Saluran memungkinkan goroutine untuk saling bertukar data dengan aman tanpa saling berebut memori (*thread-safe*).

---

## 🏭 Analogi Nyata: Ban Berjalan di Pabrik

Bayangkan sebuah ban berjalan (*conveyor belt*) di sebuah pabrik:
* **Produsen (Pengirim)**: Pekerja yang meletakkan produk jadi ke ban berjalan (`.kirim()`).
* **Konsumen (Penerima)**: Pekerja di ujung ban yang mengambil produk tersebut (`.terima()`).
* **Kapasitas Penampung**: Jika ban berjalan memiliki tempat tunggu (*buffer*), pekerja pengirim tidak perlu menunggu sampai pekerja penerima selesai mengambil barangnya.

```
 [ Produsen ] ──► .kirim(data) ──► [ Saluran / Ban Berjalan ] ──► .terima() ──► [ Konsumen ]
```

---

## 1. Membuat Saluran

Celang mendukung saluran tanpa penyangga (*unbuffered*) maupun berpanyangga (*buffered*):

```celang
// 1. Saluran tanpa penyangga (kapasitas 0)
// Pengirim akan menunggu sampai ada penerima yang siap
buat pipaLangsung = Saluran()

// 2. Saluran dengan penyangga (kapasitas 3)
// Pengirim bisa mengirim hingga 3 item sebelum tertahan
buat antrean = Saluran(3)
```

### Saluran Bertipe (*Typed Channel*)
Anda dapat memadukan saluran dengan fitur Generik untuk menjamin hanya tipe data tertentu yang bisa dikirim:

```celang
buat saluranAngka = Saluran<Angka>(5)

saluranAngka.kirim(100) // ✅ Sah!
// saluranAngka.kirim("halo") // ❌ Error! Tipe data tidak cocok
```

---

## 2. Operasi Saluran

Saluran memiliki metode bawaan yang ringkas:

```celang
buat antrean = Saluran(2)

// 1. Mengirim data
antrean.kirim("Tugas A")
antrean.kirim("Tugas B")

// 2. Memeriksa isi dan kapasitas
cetak("Banyak item:", antrean.panjang)     // 2
cetak("Kapasitas total:", antrean.kapasitas) // 2

// 3. Menerima data
buat tugas1 = antrean.terima()
buat tugas2 = antrean.terima()
cetak(tugas1, tugas2) // "Tugas A" "Tugas B"

// 4. Menutup saluran
antrean.tutup()
```

---

## 3. Streaming dengan Perulangan (`untuk .. dalam`)

Celang memungkinkan Anda mengonsumsi data yang dialirkan terus-menerus lewat perulangan `untuk .. dalam`. Perulangan ini akan otomatis membaca data satu per satu dan **otomatis berhenti** begitu saluran ditutup (`.tutup()`):

```celang
buat pipaPesan = Saluran(5)

// Produsen: Kirim data dari latar belakang
sambil {
    untuk angka dalam [10, 20, 30] {
        pipaPesan.kirim(angka)
    }
    pipaPesan.tutup() // Jangan lupa tutup setelah selesai mengirim!
}

// Konsumen: Baca data secara streaming
cetak("Menerima data secara streaming:")
untuk nilai dalam pipaPesan {
    cetak("Data masuk:", nilai)
}
// Output:
// Data masuk: 10
// Data masuk: 20
// Data masuk: 30
```

---

## 4. Pola Produsen & Konsumen (*Producer-Consumer*)

Berikut adalah contoh lengkap komunikasi asinkron antara produsen dan konsumen:

```celang
buat antreanPesanan = Saluran(3)

// Produsen: Menyiapkan pesanan
sambil {
    buat daftarMenu = ["Nasi Goreng", "Sate Ayam", "Es Teh"]
    untuk menu dalam daftarMenu {
        cetak("Koki memasak:", menu)
        antreanPesanan.kirim(menu)
    }
    antreanPesanan.tutup()
}

// Konsumen: Pelayan menyajikan pesanan
untuk makanan dalam antreanPesanan {
    cetak("Pelayan menyajikan:", makanan)
}
```

---

## 🎯 Rangkuman Metode Saluran

| Metode / Properti | Penjelasan |
| :--- | :--- |
| `Saluran(kapasitas)` | Membuat saluran dengan kapasitas buffer tertentu |
| `Saluran<T>(kapasitas)` | Saluran dengan batasan tipe data generik |
| `saluran.kirim(data)` | Mengirim data ke dalam saluran |
| `saluran.terima()` | Mengambil data berikutnya dari saluran |
| `saluran.tutup()` | Menutup saluran dan memberi sinyal selesai |
| `saluran.panjang` | Jumlah elemen yang saat ini ada di antrean |
| `saluran.kapasitas` | Kapasitas maksimal penampung saluran |
| `untuk x dalam saluran` | Iterasi otomatis sampai saluran ditutup |
