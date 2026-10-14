# Keamanan Nilai Kosong & Kenyamanan Menulis

Salah satu penyebab paling sering sebuah program komputer mendadak macet adalah **mencari data yang ternyata kosong**. 

Celang menyediakan beberapa jaring pengaman alami agar kode Anda tetap berjalan tenang dan tidak mudah terhenti:

1. **Pemeriksaan Aman:** `mungkin`
2. **Nilai Cadangan:** `kalau kosong` (atau `??`)
3. **Pembersihan Otomatis:** `saat_selesai`
4. **Rentang Angka & Pemotongan:** `..`

---

## 1. Menjaga Agar Tidak Macet dengan `mungkin`

Bayangkan Anda ingin mencari nama kota dari data pengguna:
```celang
buat kota = mungkin pengguna.profil.alamat.kota
```

### Bagaimana Cara Kerjanya?
* Jika data `pengguna`, `profil`, atau `alamat` ternyata belum diisi atau bernilai kosong, Celang **tidak akan macet atau mengeluarkan pesan galat**.
* Celang akan berhenti mencari dengan tenang dan memberikan hasil `kosong` (null).
* Jika semua datanya lengkap, nama kota akan diambil seperti biasa.

> **Tips Singkat:** Anda juga bisa menggunakan tanda `?.` (contoh: `pengguna?.profil?.kota`) jika menyukai penulisan cepat.

---

## 2. Menyiapkan Nilai Cadangan dengan `kalau kosong`

Jika data yang dicari ternyata kosong, Anda bisa langsung menyiapkan jawaban pengganti:

```celang
buat namaTampilan = inputNama kalau kosong "Tamu Anonim"
buat nomorPort = konfigurasi.port kalau kosong 8080
```

### Menggabungkan `mungkin` dan `kalau kosong`:
Keduanya bisa dibaca seperti kalimat bahasa Indonesia biasa yang sangat wajar:

```celang
buat kota = mungkin pengguna.profil.alamat.kota kalau kosong "Jakarta"
```
Artinya: *Ambil kota jika ada, tetapi kalau kosong, gunakan "Jakarta".*

> **Catatan Penting:**  
> Nilai `salah` (false) dan angka `0` bukanlah nilai kosong! Nilai cadangan hanya akan dipakai jika datanya benar-benar tidak ada (`kosong` / null).

---

## 3. Pembersihan Otomatis dengan `saat_selesai`

Terkadang kita perlu memastikan bahwa suatu tindakan perapian selalu dikerjakan sebelum sebuah fungsi selesai (misalnya menutup dokumen atau melepaskan sambungan):

```celang
fungsi bacaCatatan(namaBerkas) {
    buat dokumen = bukaBerkas(namaBerkas)
    saat_selesai tutupBerkas(dokumen) // DIJAMIN selalu dijalankan sebelum fungsi selesai!

    jika dokumenKosong(dokumen) {
        kembalikan kosong // dokumen tetap ditutup dengan aman!
    }

    kembalikan olahIsi(dokumen)
}
```

### Urutan Pelaksanaan
Jika ada lebih dari satu perintah `saat_selesai`, perintah tersebut akan dijalankan dengan urutan terbalik (yang terakhir didaftarkan akan dikerjakan lebih dulu, seperti menumpuk piring lalu mengambil piring teratas terlebih dahulu).

---

## 4. Rentang Angka & Memotong Bagian (`..`)

### A. Membuat Deretan Angka Cepat (`awal..akhir`)
Anda dapat membuat urutan angka langsung menggunakan tanda dua titik `..`:

```celang
buat angka = 1..5
cetak(angka) // [1, 2, 3, 4, 5]
```

Dapat langsung dipadukan dengan perulangan:
```celang
untuk i dalam 1..5 {
    cetak("Langkah ke-", i)
}
```

### B. Mengambil Sebagian Isi Daftar
Anda dapat memotong sebagian daftar dengan menyebutkan rentangnya:

```celang
buat data = [10, 20, 30, 40, 50]

buat tengah = data[1..4]  // Mengambil dari nomor 1 sampai sebelum 4 -> [20, 30, 40]
buat awalan = data[..2]   // Dari awal sampai nomor 2 -> [10, 20]
buat akhiran = data[3..]  // Dari nomor 3 sampai habis -> [40, 50]
```

### C. Mengambil Sebagian Kata dalam Kalimat
Cara yang sama juga berlaku untuk kalimat teks:

```celang
buat nama = "Bahasa Celang"
buat cuplikan = nama[7..13] // "Celang"
cetak(cuplikan)
```
Celang menangani pemotongan teks secara cerdas, sehingga simbol dan karakter emoji tetap terpotong rapi tanpa merusak susunan hurufnya.
