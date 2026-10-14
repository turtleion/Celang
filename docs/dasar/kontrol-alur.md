# Pengambilan Keputusan, Perulangan, & Kendala

Sebuah program komputer menjadi pintar karena ia mampu **mengambil keputusan** berdasarkan situasi yang terjadi, **mengulangi pekerjaan** yang sama berkali-kali tanpa lelah, dan **menangani kendala** jika ada sesuatu yang tidak terduga.

Di Celang, seluruh alur logika ini ditulis menggunakan kalimat bahasa Indonesia yang mengalir alami.

---

## 1. Mengambil Keputusan (`jika`, `selain jika`, `selain`)

Gunakan kata kunci `jika` untuk menjalankan perintah hanya ketika syarat tertentu terpenuhi:

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
Anda boleh menyisipkan kata `maka` sebelum kurung kurawal agar kalimatnya terbaca lebih mantap:
```celang
buat usia = 18

jika usia >= 17 maka {
    cetak("Kartu Tanda Penduduk dapat diterbitkan.")
} selain {
    cetak("Belum memenuhi syarat usia.")
}
```

---

## 2. Memilih Berdasarkan Kasus (`cocok ... kasus ... lainnya`)

Ketika Anda memiliki banyak pilihan yang pasti (misalnya menu pilihan, peran pengguna, atau nomor kode), pernyataan `cocok` jauh lebih rapi dibanding menulis banyak `selain jika`:

```celang
buat peranPengguna = "editor"

cocok peranPengguna {
    kasus "admin":
        cetak("Akses penuh: Semua menu terbuka.")
    kasus "editor":
        cetak("Akses penyunting: Boleh menyunting artikel.")
    kasus "penulis":
        cetak("Akses pembuat: Boleh menulis draf baru.")
    lainnya:
        cetak("Peran tidak dikenali. Akses hanya untuk membaca.")
}
```

---

## 3. Mengulangi Selama Syarat Terpenuhi (`selama`)

Gunakan kata kunci `selama` jika Anda ingin mengulangi suatu tindakan selama kondisinya masih bernilai `benar`:

```celang
buat hitungan = 1

selama hitungan <= 4 {
    cetak("Langkah ke-", hitungan)
    hitungan += 1 // Menambah hitungan agar tidak berputar selamanya
}

cetak("Perulangan selesai!")
```

> **Perhatian:** Pastikan angka hitungan selalu ditambah di dalam perulangan agar syaratnya bisa tercapai dan program tidak berputar tanpa henti.

---

## 4. Membaca Seluruh Isi Daftar (`untuk ... dalam`)

Jika Anda memiliki sebuah daftar barang dan ingin memeriksa isinya satu per satu:

```celang
buat daftarKota = ["Bandung", "Surabaya", "Medan"]

untuk kota dalam daftarKota {
    cetak("Mengunjungi kota:", kota)
}
```

Bisa juga dipadukan dengan rentang angka cepat:
```celang
untuk nomor dalam 1..5 {
    cetak("Antrean nomor:", nomor)
}
```

### Mengatur Alur: `berhenti` dan `lanjut`
- **`berhenti`**: Langsung menghentikan perulangan saat itu juga.
- **`lanjut`**: Melompati putaran saat ini dan langsung maju ke putaran berikutnya.

```celang
untuk angka dalam 1..10 {
    jika angka == 3 {
        lanjut // Lewati angka 3, jangan cetak
    }
    jika angka == 7 {
        berhenti // Berhenti total saat mencapai angka 7
    }
    cetak("Angka:", angka)
}
```

---

## 5. Menangani Kendala (`coba`, `tangkap`, `akhirnya`, `lempar`)

Terkadang program kita menghadapi situasi yang tidak terduga, misalnya pembagian dengan angka nol atau sambungan yang putus. Daripada membiarkan program langsung macet, Anda dapat menangkap kendala tersebut dengan tenang:

```celang
fungsi bagiUang(total, orang) {
    jika orang <= 0 {
        lempar "Jumlah orang tidak boleh nol atau minus!"
    }
    kembalikan total / orang
}

coba {
    buat hasil = bagiUang(100000, 0)
    cetak("Masing-masing mendapat:", hasil)
} tangkap (kendala) {
    cetak("Terjadi kendala perhitungan:", kendala)
} akhirnya {
    cetak("Pemeriksaan selesai dikerjakan.")
}
```

- **`coba`**: Mencoba menjalankan perintah yang rawan masalah.
- **`tangkap`**: Menangkap pesan masalah jika terjadi kesalahan, sehingga program tetap bisa berjalan normal.
- **`akhirnya`**: Bagian yang dijamin **pasti dijalankan**, baik perhitungannya berhasil maupun gagal.
- **`lempar`**: Sengaja mengeluarkan peringatan masalah jika ada data yang tidak masuk akal.
