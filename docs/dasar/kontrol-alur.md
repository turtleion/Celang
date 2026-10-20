# Percabangan & Kontrol Alur

Sebuah program komputer dapat merespons situasi yang berbeda-beda karena ia mampu **mengambil keputusan**. Di Celang, Anda dapat mengarahkan jalannya program menggunakan pernyataan kondisional berbahasa Indonesia yang mengalir alami dan mudah dipahami.

---

## 1. Percabangan Bersyarat (`jika`, `selain jika`, `selain`)

Gunakan kata kunci `jika` untuk menjalankan blok kode hanya ketika kondisi tertentu terpenuhi (`benar`):

```celang
buat nilaiUjian = 85

jika nilaiUjian >= 90 {
    cetak("Hasil: Istimewa (A)")
} selain jika nilaiUjian >= 80 {
    cetak("Hasil: Sangat Baik (B)")
} selain jika nilaiUjian >= 70 {
    cetak("Hasil: Cukup (C)")
} selain {
    cetak("Hasil: Perlu Belajar Lagi (D)")
}
```

### Menggunakan Kata `maka` agar Lebih Alami

Untuk meningkatkan keterbacaan, Anda dapat menyisipkan kata `maka` sebelum kurung kurawal pembuka `{`:

```celang
buat usia = 18

jika usia >= 17 maka {
    cetak("Kartu Tanda Penduduk dapat diterbitkan.")
} selain {
    cetak("Belum memenuhi syarat usia.")
}
```

### Merangkai Banyak Syarat (`dan`, `atau`, `bukan`)

Anda dapat memadukan beberapa kondisi menggunakan operator logika:

```celang
buat punyaTiket = benar
buat usia = 20

jika punyaTiket dan usia >= 18 maka {
    cetak("Silakan masuk ke area bioskop.")
} selain {
    cetak("Akses ditolak.")
}
```

---

## 2. Percabangan Multi-Pilihan (`pilih ... bila ... lainnya`)

Ketika Anda perlu mencocokkan suatu nilai terhadap banyak kemungkinan tetap (seperti menu pilihan, peran pengguna, kode status, atau tipe data), pernyataan `pilih ... bila` jauh lebih rapi, ekspresif, dan enak dibaca dibandingkan menulis banyak `selain jika`:

```celang
buat peranPengguna = "editor"

pilih peranPengguna {
    bila "admin":
        cetak("Akses penuh: Semua menu terbuka.")
    bila "editor":
        cetak("Akses penyunting: Boleh menyunting artikel.")
    bila "penulis":
        cetak("Akses pembuat: Boleh menulis draf baru.")
    lainnya:
        cetak("Peran tidak dikenali. Akses hanya untuk membaca.")
}
```

### Karakteristik `pilih ... bila` di Celang:
- **Otomatis Berhenti:** Celang tidak memerlukan pernyataan `break` manual di setiap cabang (*no fallthrough*). Setelah satu cabang `bila` yang cocok dieksekusi, evaluasi blok `pilih` langsung selesai.
- **Cabang Bawaan (`lainnya`):** Menangani seluruh kemungkinan yang tidak tercakup oleh cabang-cabang `bila` sebelumnya.
- **Dukungan Berhenti Awal (`berhenti`):** Jika Anda ingin keluar dari blok `bila` lebih cepat, Anda dapat menggunakan kata kunci `berhenti`.

---

## 🧭 Langkah Selanjutnya

- Ingin mengulangi pekerjaan berkali-kali? Pelajari [**Perulangan & Iterasi (Looping)**](/dasar/perulangan).
- Ingin mengamankan program dari kesalahan tak terduga? Pelajari [**Penanganan Galat (Try-Catch)**](/dasar/penanganan-galat).
