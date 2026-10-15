# Saluran Komunikasi Proses (Saluran / Channel CSP)

Ketika Anda memiliki dua bagian program yang bekerja sendiri-sendiri, bagaimana cara keduanya saling bertukar informasi secara teratur tanpa saling mengganggu?

Celang menyediakan **Saluran** (`saluran`), yaitu pipa pengantar pesan khusus yang menghubungkan dua bagian program agar bisa saling mengirim dan menerima data dengan aman.

---

## 🏭 Analogi Nyata: Pipa Pengantar Surat atau Ban Berjalan

Bayangkan sebuah ban berjalan di restoran atau pipa pengantar surat di kantor:
* **Pengirim**: Orang yang memasukkan pesan atau barang ke dalam pipa (`.kirim()`).
* **Penerima**: Orang di ujung pipa yang mengambil barang yang keluar (`.terima()`).
* **Kapasitas Pipa**: Pipa bisa memiliki ruang tampung antrean. Selama pipa belum penuh, pengirim bisa terus memasukkan pesan tanpa harus menunggu penerima mengambilnya.

```
 [ Pengirim ] ──► .kirim(pesan) ──► [ Saluran Pipa ] ──► .terima() ──► [ Penerima ]
```

---

## 1. Membuat Saluran

Anda dapat membuat saluran biasa atau saluran yang memiliki kapasitas antrean:

```celang
// 1. Saluran langsung (kapasitas 0)
// Pengirim baru akan melangkah setelah penerima siap mengambilnya
buat pipaLangsung = Saluran()

// 2. Saluran dengan antrean (kapasitas 3 pesan)
// Pengirim bisa memasukkan hingga 3 pesan sekaligus ke dalam antrean
buat antrean = Saluran(3)
```

### Menentukan Jenis Pesan yang Boleh Masuk
Anda bisa membatasi saluran agar hanya menerima jenis data tertentu (misalnya hanya menerima angka):

```celang
buat saluranAngka = Saluran<Angka>(5)

saluranAngka.kirim(100) // ✅ Sah!
// saluranAngka.kirim("halo") // ❌ Celang akan mengingatkan bahwa isi harus berupa Angka
```

---

## 2. Cara Mengirim dan Menerima Pesan

```celang
buat antrean = Saluran(2)

// 1. Memasukkan pesan ke dalam saluran
antrean.kirim("Tugas A")
antrean.kirim("Tugas B")

// 2. Memeriksa isi antrean
cetak("Banyak pesan di dalam:", antrean.panjang)     // 2
cetak("Kapasitas maksimal    :", antrean.kapasitas) // 2

// 3. Mengambil pesan dari saluran
buat tugas1 = antrean.terima()
buat tugas2 = antrean.terima()
cetak(tugas1, tugas2) // "Tugas A" "Tugas B"

// 4. Menutup saluran setelah semua selesai
antrean.tutup()
```

---

## 3. Membaca Pesan Mengalir (`untuk .. dalam`)

Salah satu kenyamanan terbaik di Celang adalah Anda bisa menyedot pesan yang mengalir terus-menerus menggunakan perulangan biasa `untuk .. dalam`. Perulangan ini akan membaca pesan satu per satu dan **otomatis berhenti dengan sendirinya** begitu saluran ditutup:

```celang
buat antreanPesanan = Saluran(3)

// Di latar belakang: Koki menyiapkan pesanan
sambil {
    buat daftarMenu = ["Nasi Goreng", "Sate Ayam", "Es Teh"]
    untuk menu dalam daftarMenu {
        cetak("Koki memasak:", menu)
        antreanPesanan.kirim(menu)
    }
    antreanPesanan.tutup() // Beritahu bahwa semua pesanan sudah selesai
}

// Di bagian depan: Pelayan menyajikan makanan saat tiba
untuk makanan dalam antreanPesanan {
    cetak("Pelayan menyajikan:", makanan)
}
```

---

## 🎯 Ringkasan Perintah Saluran

| Perintah | Kegunaan |
| :--- | :--- |
| `Saluran(kapasitas)` | Membuat pipa saluran dengan kapasitas antrean tertentu |
| `saluran.kirim(data)` | Memasukkan pesan ke dalam saluran |
| `saluran.terima()` | Mengambil pesan berikutnya dari saluran |
| `saluran.tutup()` | Menutup saluran sebagai tanda bahwa pengiriman sudah selesai |
| `saluran.panjang` | Mengetahui berapa banyak pesan yang sedang mengantre |
| `saluran.kapasitas` | Mengetahui daya tampung maksimal saluran |
| `untuk pesan dalam saluran` | Mengambil pesan secara otomatis sampai saluran ditutup |
