# 🇮🇩 Satuan Budaya & Terbilang Bawaan

Celang adalah bahasa pemrograman pertama yang dirancang untuk benar-benar memahami cara berhitung dan bertransaksi masyarakat Indonesia: dari singkatan bilangan sehari-hari (`ribu`, `juta`, `miliar`), konversi ejaan terbilang bahasa Indonesia (`terbilang`), hingga format mata uang Rupiah (`rupiah`) langsung di tingkat sintaks!

---

## 💎 1. Literal Satuan Budaya & Angka Ergonomis

Alih-alih mengetik angka dengan deretan nol yang rawan kesalahan ketik (seperti `1000000000`):

```celang
buat uang_saku = 50.ribu        // 50000
buat gaji_pokok = 2.5.juta      // 2500000
buat dana_apbd = 1.miliar       // 1000000000
buat investasi = 500.juta       // 500000000
```

### Satuan Digital
Celang juga mendukung satuan komputasi modern:
```celang
buat memori = 16.gb             // 17179869184 byte
buat berkas_size = 512.mb       // 536870912 byte
buat transfer = 64.kb           // 65536 byte
```

> **Optimasi Kompiler:** Literal satuan seperti `50.ribu` dioptimasi saat kompilasi (*constant folding*) langsung menjadi nilai integer utuh sehingga tidak membebani performa runtime.

---

## 🗣️ 2. Konversi Terbilang Bawaan (`.terbilang` & `terbilang()`)

Di industri perbankan, akuntansi, dan invoice, mencetak angka dalam bentuk ejaan kata bahasa Indonesia adalah persyaratan mutlak (misalnya *"dua juta lima ratus ribu rupiah"*). Di bahasa lain, pengembang harus menginstal paket tambahan. Di Celang, kemampuan ini hadir langsung di dalam inti bahasa!

### Properti Objek `.terbilang`
```celang
buat harga = 50.ribu
cetak(harga.terbilang) 
// Output: lima puluh ribu rupiah

buat gaji = 2.5.juta
cetak(gaji.terbilang)
// Output: dua juta lima ratus ribu rupiah

buat angka = 1000
cetak(angka.terbilang)
// Output: seribu rupiah
```

### Properti `.kata` (Tanpa Akhiran Rupiah)
Jika Anda hanya ingin mengeja bilangan tanpa kata "rupiah":
```celang
buat nilai = 105
cetak(nilai.kata)
// Output: seratus lima
```

### Fungsi Bawaan Global `terbilang(angka)`
```celang
buat teks = terbilang(1500000)
cetak(teks)
// Output: satu juta lima ratus ribu rupiah
```

---

## 💰 3. Pemformatan Rupiah (`.rupiah` & `format_rupiah()`)

Mencetak angka dengan pemisah ribuan standar Indonesia (`Rp 50.000`):

```celang
buat tagihan = 1250000
cetak(tagihan.rupiah)
// Output: Rp 1.250.000

buat total = 75.ribu
cetak(format_rupiah(total))
// Output: Rp 75.000
```
