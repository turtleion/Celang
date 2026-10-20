# Perulangan & Iterasi (Looping)

Komputer sangat unggul dalam **mengulangi pekerjaan** yang sama ratusan, ribuan, hingga jutaan kali tanpa merasa lelah atau membuat kesalahan ceroboh. 

Di Celang, Anda memiliki dua cara utama untuk melakukan perulangan:
1. **`selama`** (*while loop*): Mengulang selama suatu syarat masih terpenuhi.
2. **`untuk ... dalam`** (*for-in loop*): Menjelajahi setiap elemen dalam sebuah kumpulan data atau rentang angka.

---

## 1. Perulangan Bersyarat (`selama`)

Gunakan kata kunci `selama` jika Anda ingin mengulang suatu tindakan selama kondisi yang ditentukan bernilai `benar`:

```celang
buat hitungan = 1

selama hitungan <= 4 {
    cetak("Langkah ke-", hitungan)
    hitungan += 1 // Menambah hitungan agar tidak berputar selamanya
}

cetak("Perulangan selesai!")
```

Keluaran:
```text
Langkah ke- 1
Langkah ke- 2
Langkah ke- 3
Langkah ke- 4
Perulangan selesai!
```

> ⚠️ **Penting:** Pastikan kondisi di dalam `selama` selalu diperbarui di dalam badan perulangan agar syaratnya dapat berhenti di waktu tertentu, sehingga program tidak mengalami putaran tanpa akhir (*infinite loop*).

---

## 2. Perulangan Koleksi & Rentang (`untuk ... dalam`)

Perulangan `untuk ... dalam` digunakan ketika Anda memiliki kumpulan data (seperti daftar/array, kamus/map, atau rentang angka) dan ingin memproses elemen-elemennya satu per satu.

### A. Membaca Isi Daftar (Array)

```celang
buat daftarKota = ["Bandung", "Surabaya", "Medan", "Makassar"]

untuk kota dalam daftarKota {
    cetak("Mengunjungi kota:", kota)
}
```

### B. Menggunakan Rentang Angka Cepat (`..`)

Anda dapat menghasilkan deretan angka secara instan menggunakan operator rentang dua titik `..`:

```celang
untuk nomor dalam 1..5 {
    cetak("Antrean nomor:", nomor)
}
```

Keluaran:
```text
Antrean nomor: 1
Antrean nomor: 2
Antrean nomor: 3
Antrean nomor: 4
Antrean nomor: 5
```

Anda juga dapat menggunakan fungsi bawaan `rentang(awal, akhir)`:

```celang
untuk angka dalam rentang(0, 10) {
    cetak("Nilai:", angka)
}
```

### C. Membaca Kunci Kamus (Map)

Ketika Anda melakukan perulangan pada sebuah `Map`, Celang akan mengiterasikan setiap kunci (*keys*) yang ada di dalam map secara teratur:

```celang
buat profil = {
    "nama": "Siti Rahma",
    "divisi": "Teknologi",
    "kota": "Jakarta"
}

untuk kunci dalam profil {
    cetak(kunci, "=>", profil[kunci])
}
```

### D. Mengalirkan Data dari Saluran (Channel)

Celang mendukung pembacaan saluran inter-proses konkurensi (CSP) secara alami menggunakan loop `untuk`:

```celang
buat ch = Saluran(3)
ch.kirim("Pesan 1")
ch.kirim("Pesan 2")
ch.tutup()

untuk pesan dalam ch {
    cetak("Menerima:", pesan)
}
```

---

## 3. Pengendali Putaran: `berhenti` & `lanjut`

Terkadang Anda ingin keluar dari perulangan lebih awal atau melompati bagian tertentu saat kondisi khusus terjadi:

| Perintah | Padanan Global | Maksud |
| :--- | :--- | :--- |
| **`berhenti`** | `break` | Langsung menyudahi dan keluar dari seluruh proses perulangan. |
| **`lanjut`** | `continue` | Melompati sisa perintah pada putaran saat ini dan langsung melompat ke putaran berikutnya. |

```celang
untuk angka dalam 1..10 {
    // Lewati angka 3, jangan dicetak
    jika angka == 3 {
        lanjut
    }

    // Berhenti total saat mencapai angka 7
    jika angka == 7 {
        berhenti
    }

    cetak("Angka:", angka)
}
```

Keluaran:
```text
Angka: 1
Angka: 2
Angka: 4
Angka: 5
Angka: 6
```
*(Perhatikan bahwa angka 3 dilompati oleh `lanjut`, dan perulangan selesai di angka 6 sebelum angka 7 dicetak).*

---

## 4. Kapan Menggunakan `selama` vs `untuk`?

- Gunakan **`untuk ... dalam`** saat Anda sudah mengetahui jumlah atau kumpulan data yang ingin dijelajahi (seperti elemen array, kunci map, atau rentang angka pasti). Ini adalah perulangan yang paling sering digunakan dan paling aman.
- Gunakan **`selama`** saat Anda perlu mengulang sampai suatu kondisi eksternal terpenuhi (misalnya menunggu pengguna menekan tombol tertentu, polling status jaringan, atau game loop).

---

## 🧭 Langkah Selanjutnya

- Ingin mengamankan aplikasi dari kesalahan tak terduga saat perulangan atau pengolahan data? Pelajari [**Penanganan Galat (Try-Catch)**](/dasar/penanganan-galat).
- Ingin mengolah array dan map lebih ringkas dengan method fungsional (`.filter()`, `.map()`, `.reduce()`)? Pelajari [**Manipulasi Koleksi (Array & Map)**](/fitur-khas/koleksi-array-map).
