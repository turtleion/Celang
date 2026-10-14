# Variabel yang Mengingat Masa Lalu (Time-Travel)

Salah satu keistimewaan paling unik di Celang adalah **Variabel Penjelajah Waktu**. 

Di hampir semua bahasa pemrograman lain, saat Anda mengubah isi sebuah variabel, nilai yang lama langsung terhapus dan hilang selamanya. Tetapi di Celang, variabel secara alami mengingat masa lalunya.

---

## ⏳ Tiga Hal yang Diingat oleh Variabel

Setiap variabel yang Anda buat menggunakan kata kunci `buat` secara otomatis memiliki 3 catatan waktu:

| Catatan Waktu | Panggilan Alternatif | Keterangan |
| :--- | :--- | :--- |
| `.sebelumnya` | `.prev`, `.previous` | Memberikan nilai tepat satu langkah sebelum nilai saat ini. Jika belum pernah diubah, nilainya adalah `kosong` (null). |
| `.awal` | `.initial`, `.first` | Memberikan nilai mula-mula saat variabel pertama kali diciptakan. |
| `.histori` | `.history` | Memberikan daftar urutan perubahan nilai variabel dari waktu ke waktu. |

---

## 📝 Contoh Nyata: Menelusuri Kenaikan Harga

Bayangkan Anda sedang mencatat harga barang yang mengalami beberapa kali kenaikan:

```celang
// 1. Harga saat pertama kali dicatat
buat harga = 50000

// 2. Harga mengalami beberapa kali kenaikan
harga = 75000
harga = 90000
harga = 120000

// 3. Mari kita periksa ingatan variabel harga:
cetak("Harga sekarang   :", harga)             // 120000
cetak("Harga sebelumnya :", harga.sebelumnya)  // 90000
cetak("Harga mula-mula  :", harga.awal)        // 50000
cetak("Histori kenaikan :", harga.histori)     // [50000, 75000, 90000, 120000]
```

### Kapan Fitur Ini Sangat Membantu?
- **Menghitung Selisih**: Mengetahui selisih kenaikan harga antara hari ini dan kemarin tanpa perlu membuat variabel `harga_kemarin`.
- **Fitur Batal (Undo)**: Mengembalikan nilai ke kondisi sebelumnya jika pengguna membatalkan tindakan.
- **Audit dan Riwayat**: Mengetahui apakah suatu nilai sudah pernah diubah dari nilai aslinya.

---

## 🛡️ Apakah Komputer Tidak Keberatan Menyimpan Riwayat?

Mungkin Anda bertanya:
> *"Jika variabel diubah ribuan atau jutaan kali di dalam perulangan, apakah memori komputer tidak akan kepenuhan?"*

Jawabannya: **Tidak, komputer Anda tetap aman dan ringan.**

Celang mengelola ingatan variabel dengan sangat bijak:
1. **Nilai Awal (`.awal`)**: Disimpan secara terpisah dalam satu catatan tetap, sehingga selalu aman.
2. **Nilai Terakhir (`.sebelumnya`)**: Selalu siap diambil seketika.
3. **Buku Riwayat (`.histori`)**: Membatasi catatan maksimal hingga **50 perubahan terakhir**. Jika ada perubahan ke-51, catatan paling lama akan digeser keluar secara teratur.

Dengan cara ini, meskipun Anda mengubah isi variabel sebanyak sepuluh juta kali di dalam perulangan, penggunaan memori komputer Anda **tetap stabil dan ringan**, tanpa membuang-buang ruang penyimpanan.
