# Kinerja dan Pengelolaan Memori di Celang (`PERFORMANCE_AND_MEMORY.md`)

Dokumen ini menjelaskan bagaimana Celang dirancang agar tetap **berjalan cepat** dan **hemat memori di komputer Anda**, meskipun memiliki fitur-fitur canggih seperti variabel penjelajah waktu (*time-travel*) dan gaya bahasa dinamis.

---

## 1. Menjaga Memori pada Variabel Penjelajah Waktu

Salah satu pertanyaan yang sering muncul adalah:
> *"Jika setiap variabel mengingat riwayat perubahannya di masa lalu, apakah memori komputer tidak akan cepat penuh jika variabel diubah jutaan kali?"*

Jawabannya: **Tidak, komputer Anda tetap aman.**

Celang menerapkan aturan penyimpanan yang sangat terukur:
1. **Nilai Mula-mula (`.awal`)**: Disimpan di catatan awal tersendiri yang tidak pernah hilang.
2. **Nilai Sebelumnya (`.sebelumnya`)**: Selalu siap diambil seketika tepat satu langkah sebelum nilai aktif.
3. **Buku Riwayat (`.histori`)**: Dibatasi maksimal **50 perubahan terakhir**. Jika variabel diubah untuk ke-51 kalinya, catatan yang paling lama akan digeser keluar secara teratur.

Dengan cara ini, sebanyak apa pun perulangan yang Anda lakukan (bahkan hingga 10 juta kali), ruang memori yang dipakai oleh suatu variabel **tidak akan membengkak**, dan komputer Anda tetap bekerja dengan ringan dan lancar.

---

## 2. Mengapa Mesin Virtual Celang Sangat Cepat?

Secara bawaan, Celang tidak menjalankan kode dengan cara membolak-balik teks program secara lambat. Celang bekerja melalui dua tahap cepat:

1. **Peringkasan Instruksi (Kompilasi ke Bytecode)**:
   Program Anda terlebih dahulu diringkas menjadi deretan instruksi kecil yang sangat padat.
2. **Mesin Tumpukan Sederhana (Virtual Machine)**:
   Mesin virtual Celang membaca instruksi kecil tersebut satu per satu menggunakan model tumpukan sederhana (seperti menaruh dan mengambil piring di meja).

Karena instruksi yang dibaca berukuran sangat kecil dan proses tumpukan nilainya sangat langsung, komputer tidak perlu berpikir berbelit-belit, sehingga eksekusi program berlangsung dalam sekejap mata.

---

## 3. Menghindari Pemborosan Data dengan `salin`

Ketika Anda bekerja dengan data yang besar (seperti daftar atau bentuk data), Celang secara cerdas membagikan data tersebut secara aman. Namun, jika Anda secara sengaja ingin membuat duplikat terpisah yang benar-benar mandiri:

```celang
buat dataAsli = [10, 20, 30]
buat dataTiruan = salin(dataAsli)

dataTiruan.tambah(40)
cetak(dataAsli)   // [10, 20, 30] (tetap aman, tidak terpengaruh)
cetak(dataTiruan) // [10, 20, 30, 40]
```
Perintah `salin` menduplikasi data sampai ke bagian terdalamnya, sehingga kedua data hidup secara mandiri tanpa saling mengganggu.
