# Ergonomi & Null Safety di Celang

Celang mengusung filosofi *human-centric*: setiap operasi yang sering digunakan dibuat sejelas dan senyaman mungkin bagi manusia tanpa mengorbankan keamanan sistem.

Fitur-fitur ergonomi modern berikut hadir tanpa simbol-simbol kriptik yang membingungkan:
1. **Navigasi Aman (*Null Safety*):** `mungkin`
2. **Nilai Cadangan (*Null Coalescing*):** `kalau kosong` (atau `??`)
3. **Pembersihan Otomatis (*Defer*):** `saat_selesai`
4. **Sintaks Rentang & Pemotongan (*Range & Slicing*):** `..`

---

## 1. Navigasi Aman dengan `mungkin`

Dalam bahasa seperti JavaScript, C#, atau Kotlin, pengecekan properti bersarang menggunakan simbol `?.` (*optional chaining*). Di Celang, Anda dapat menulisnya secara ekspresif menggunakan kata kunci `mungkin`:

```celang
buat kota = mungkin pengguna.profil.alamat.kota
```

### Cara Kerja:
* Jika salah satu bagian dari rantai objek (`pengguna`, `profil`, atau `alamat`) bernilai `kosong` (`nihil` / `null`), Celang **tidak akan melempar error / crash**.
* Evaluasi rantai akan langsung berhenti dan mengembalikan nilai `kosong`.
* Jika seluruh rantai valid, nilai properti tujuan akan dikembalikan seperti biasa.

> **Catatan:** Operator singkat `?.` (contoh: `pengguna?.profil?.kota`) juga tetap didukung bagi yang menyukai penulisan cepat.

---

## 2. Nilai Cadangan dengan `kalau kosong`

Operasi ini menentukan nilai pengganti (*fallback*) jika ekspresi di sebelah kiri bernilai `kosong`.

```celang
buat nama_tampilan = input_nama kalau kosong "Anonim"
buat port = konfigurasi.port kalau kosong 8080
```

### Menggabungkan `mungkin` dan `kalau kosong`:
Keduanya berpadu secara sangat alami seperti kalimat bahasa Indonesia biasa:

```celang
buat kota = mungkin pengguna.profil.alamat.kota kalau kosong "Jakarta"
```
Jika `pengguna` tidak ada atau alamat belum diisi, variabel `kota` langsung berisi `"Jakarta"`!

> **Penting Mengenai Boolean dan Angka Nol:**  
> Nilai boolean `salah` (false) dan angka `0` **bukanlah** nilai kosong. Jadi:
> ```celang
> buat aktif = salah kalau kosong benar   // Hasil tetap: salah
> buat skor = 0 kalau kosong 100          // Hasil tetap: 0
> ```
> Nilai cadangan hanya dipicu jika nilai sebelah kiri benar-benar bernilai `kosong` / `null`.

> **Catatan:** Operator simbol `??` (contoh: `a ?? b`) juga didukung sebagai alternatif ringkas.

---

## 3. Pembersihan Otomatis dengan `saat_selesai` (*Defer*)

Terinspirasi dari fitur `defer` di bahasa Go, Celang menyediakan kata kunci `saat_selesai`. Kata kunci ini menjadwalkan suatu perintah atau pembersihan sumber daya (*cleanup*) agar dieksekusi **tepat sebelum fungsi selesai**:

```celang
fungsi olah_berkas(nama_berkas) {
    buat berkas = buka_berkas(nama_berkas)
    saat_selesai tutup_berkas(berkas) // DIJAMIN dieksekusi sebelum fungsi keluar

    jika berkas_kosong(berkas) {
        kembalikan kosong // tutup_berkas tetap terpanggil!
    }

    kembalikan proses(berkas)
}
```

### Urutan Eksekusi (LIFO - Last In, First Out)
Jika ada lebih dari satu `saat_selesai` dalam satu fungsi, mereka akan dieksekusi dalam urutan terbalik:

```celang
fungsi urutan() {
    saat_selesai cetak("Langkah 3")
    saat_selesai cetak("Langkah 2")
    cetak("Langkah 1")
}

urutan()
// Output:
// Langkah 1
// Langkah 2
// Langkah 3
```

---

## 4. Sintaks Rentang (`..`) & Pemotongan Koleksi (*Slicing*)

### A. Rentang Angka Native (`awal..akhir`)
Anda dapat membuat array angka berurutan secara instan menggunakan operator rentang `..`:

```celang
buat angka = 1..5
cetak(angka) // [1, 2, 3, 4, 5]
```

Dapat langsung digunakan pada perulangan `untuk ... dalam`:
```celang
untuk i dalam 1..5 {
    cetak("Iterasi ke-", i)
}
```

### B. Pemotongan Array (*Array Slicing*)
Mengambil sebagian elemen array menggunakan notasi indeks rentang:

```celang
buat data = [10, 20, 30, 40, 50]

buat tengah = data[1..4]  // Ambil indeks 1 sampai sebelum 4 -> [20, 30, 40]
buat awalan = data[..2]   // Dari awal sampai indeks 2 -> [10, 20]
buat akhiran = data[3..]  // Dari indeks 3 sampai akhir -> [40, 50]
```

### C. Pemotongan Teks (*String Slicing*)
Berlaku juga untuk tipe data teks (`Teks`):

```celang
buat nama = "Bahasa Celang"
buat sub = nama[7..13] // "Celang"
cetak(sub)
```
Celang menangani pemotongan string berbasis karakter Unicode (*rune-safe*), sehingga karakter multi-byte dan emoji dipotong dengan benar tanpa merusak byte string.
