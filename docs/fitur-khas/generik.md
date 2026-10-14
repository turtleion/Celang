# Wadah Fleksibel & Pembatas Jenis Data (Generik)

Saat membuat program, sering kali kita ingin membuat sebuah wadah atau fungsi yang dapat digunakan untuk berbagai jenis data—baik angka, teks, maupun daftar—tanpa harus menulis ulang kode yang sama berulang kali.

Fitur **Generik** di Celang memungkinkan Anda membuat satu cetak biru serbaguna yang fleksibel, namun tetap memastikan data yang dimasukkan sesuai dengan yang diharapkan.

---

## 💡 Analogi Sederhana: Kotak Serbaguna

Bayangkan Anda memiliki sebuah kotak penyimpanan:
* Tanpa fitur ini, Anda harus membuat `KotakKhususAngka`, `KotakKhususTeks`, dan `KotakKhususSurat` secara terpisah, padahal bentuk dan cara kerja kotaknya persis sama.
* Dengan fitur ini, Anda cukup membuat satu cetak biru `Kotak<Jenis>` saja. Tanda `<Jenis>` adalah penanda jenis barang yang baru akan Anda tentukan saat kotak tersebut benar-benar dipakai.

---

## 1. Membuat Bentuk Serbaguna (`<T>`)

Anda dapat menyematkan tanda `<T>` (atau nama apa pun yang Anda sukai) di belakang nama bentuk:

```celang
bentuk Kotak<T> {
    isi: T
}

// Saat digunakan, kita tentukan jenis datanya:
buat kotakAngka = Kotak<Angka>(42)
buat kotakTeks = Kotak<Teks>("Surat Penting")

cetak(kotakAngka.isi) // 42
cetak(kotakTeks.isi)  // "Surat Penting"
```

### Menggunakan Lebih dari Satu Jenis Data (`<K, V>`)
Anda juga dapat menggunakan dua jenis data sekaligus, seperti pasangan nama dan nilai:

```celang
bentuk Pasangan<Kunci, Nilai> {
    kunci: Kunci,
    nilai: Nilai
}

buat skor = Pasangan<Teks, Angka>("Budi", 950)
cetak(skor.kunci, ":", skor.nilai) // Budi : 950
```

---

## 2. Membatasi Jenis Data yang Boleh Masuk (`harus`)

Bagaimana jika Anda ingin membuat wadah serbaguna, tetapi **hanya boleh diisi angka** (misalnya dompet atau celengan)?

Celang menyediakan cara penulisan yang sangat alami menggunakan kata `harus`:

```celang
// Menegaskan bahwa nilai di dalam tabungan HARUS berupa angka
bentuk Tabungan<Nilai harus Angka> {
    jumlah: Nilai
}

buat celengan = Tabungan<Angka>(100000) // ✅ Berhasil!
```

Jika ada seseorang yang sengaja atau tidak sengaja memasukkan jenis data yang salah:
```celang
// ❌ Coba memasukkan tulisan ke celengan angka:
// buat tabunganSalah = Tabungan<Teks>("Uang Kertas")
```
Celang akan memberi tahu secara ramah:
*Jenis data 'Teks' tidak memenuhi syarat 'Angka' pada bentuk 'Tabungan'. Gunakan jenis data yang sesuai dengan angka.*

### Mengizinkan Jenis Data Apa Saja (`Sembarang`)
Jika wadah Anda benar-benar bebas diisi apa saja:
```celang
bentuk Loker<Barang: Sembarang> {
    barang: Barang
}

buat loker1 = Loker<Teks>("Helm")
buat loker2 = Loker<Angka>(999)
```

---

## 3. Fungsi Serbaguna

Selain bentuk data, fungsi juga bisa dibuat serbaguna:

```celang
// Fungsi untuk menukar posisi dua barang apa saja
fungsi tukar<T>(a: T, b: T) {
    kembalikan [b, a]
}

buat hasil1 = tukar<Angka>(10, 20)
cetak(hasil1) // [20, 10]

buat hasil2 = tukar<Teks>("Selamat", "Pagi")
cetak(hasil2) // ["Pagi", "Selamat"]
```

Dengan cara ini, Anda cukup menulis satu rumus logika untuk melayani berbagai jenis data dengan rapi dan aman.
