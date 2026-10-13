# Kontrol Alur & Penanganan Galat

Celang menyediakan struktur kontrol alur yang dirancang secara ergonomis dan alami dalam bahasa Indonesia. Seluruh struktur logika—mulai dari percabangan bersyarat, pencocokan pola kasus, perulangan koleksi, hingga penanganan galat terstruktur—dapat ditulis dan dibaca layaknya kalimat lugas sehari-hari.

---

## 1. Percabangan Bersyarat (`jika`, `selain jika`, `selain`)

Pernyataan kondisi mengevaluasi ekspresi nilai kebenaran (*boolean*). Di Celang, tanda kurung `(...)` di sekeliling kondisi bersifat **opsional**.

### Sintaks Standar

```celang
buat nilaiUjian = 85

jika nilaiUjian >= 90 {
    cetak("Predikat: A (Istimewa)")
} selain jika nilaiUjian >= 80 {
    cetak("Predikat: B (Sangat Baik)")
} selain jika nilaiUjian >= 70 {
    cetak("Predikat: C (Cukup)")
} selain {
    cetak("Predikat: D (Perlu Perbaikan)")
}
```

### Gaya Ekspresif dengan Kata Kunci `maka`

Untuk mempertegas gaya bahasa manusia, Celang mengizinkan penyertaan kata kunci opsional `maka` sebelum blok pembuka `{`:

```celang
buat usia = 18

jika usia >= 17 maka {
    cetak("KTP dapat diterbitkan.")
} selain {
    cetak("Belum memenuhi syarat usia.")
}
```

### Kondisi Gabungan (Operator Logika & Pembanding)

Anda dapat menggabungkan beberapa kondisi menggunakan operator logika kata (`dan`, `atau`, `bukan`) maupun simbol (`&&`, `||`, `!`):

```celang
buat punyaTiket = benar
buat usiaPenonton = 16
buat didampingiOrtu = benar

jika punyaTiket dan (usiaPenonton >= 17 atau didampingiOrtu) {
    cetak("Akses masuk bioskop diizinkan.")
} selain {
    cetak("Akses masuk ditolak.")
}
```

::: tip Operator Logika Bersahabat
Gunakan `dan`, `atau`, dan `bukan` untuk meningkatkan kenyamanan membaca kode dalam tim tanpa mengurangi presisi logika.
:::

---

## 2. Pencocokan Pola / Kasus (`cocok ... kasus ... lainnya`)

Ketika menghadapi banyak percabangan nilai yang spesifik, pernyataan `cocok` (*switch / pattern matching*) menghasilkan kode yang jauh lebih ringkas dan tertata dibanding rantai `selain jika` yang panjang.

### Sintaks Dasar

```celang
buat peranPengguna = "editor"

cocok peranPengguna {
    kasus "admin":
        cetak("Akses penuh: Dashboard, Pengaturan, dan Database.")
    kasus "editor":
        cetak("Akses terbatas: Publikasi dan Penyuntingan Artikel.")
    kasus "penulis":
        cetak("Akses dasar: Menulis Draf Artikel.")
    lainnya:
        cetak("Peran tidak dikenali. Akses publik hanya-baca.")
}
```

### Karakteristik `cocok`:

1. **Evaluasi Nilai Presisi**: Mencocokkan nilai ekspresi target terhadap nilai yang tertera pada tiap klausul `kasus`.
2. **Klausul Default (`lainnya`)**: Dijalankan apabila tidak ada satupun klausul `kasus` yang cocok dengan nilai target.
3. **Mendukung Beragam Tipe Data**: Dapat digunakan dengan `Angka`, `Teks`, `Logika`, maupun ekspresi kompleks.

```celang
buat kodeStatus = 404

cocok kodeStatus {
    kasus 200:
        cetak("Sukses (OK)")
    kasus 400:
        cetak("Permintaan Tidak Valid (Bad Request)")
    kasus 404:
        cetak("Halaman Tidak Ditemukan (Not Found)")
    kasus 500:
        cetak("Kesalahan Server Internal")
    lainnya:
        cetak("Kode status tidak terdaftar:", kodeStatus)
}
```

---

## 3. Perulangan Kondisional (`selama`)

Gunakan kata kunci `selama` (*while loop*) untuk mengeksekusi blok kode berulang kali selama kondisi evaluasi masih bernilai `benar`.

```celang
buat hitungan = 1

selama hitungan <= 4 {
    cetak("Iterasi ke-", hitungan)
    hitungan += 1
}

cetak("Perulangan selesai!")
```

### Output:
```text
Iterasi ke- 1
Iterasi ke- 2
Iterasi ke- 3
Iterasi ke- 4
Perulangan selesai!
```

::: warning Hindari Perulangan Tak Berujung (*Infinite Loop*)
Pastikan variabel kondisi selalu diperbarui di dalam badan perulangan agar kondisi terminasi tercapai dan program tidak mengalami kebuntuan (*freeze*).
:::

---

## 4. Perulangan Koleksi & Rentang (`untuk ... dalam`)

Kata kunci `untuk` (*for-in loop*) adalah cara utama di Celang untuk mengiterasi urutan elemen.

### 4.1 Iterasi Rentang Angka (`..` dan `rentang`)

Anda dapat menulis perulangan rentang angka berurutan menggunakan operator rentang alami `awal..akhir` atau fungsi `rentang(awal, akhir)`:

```celang
// Menggunakan operator rentang inklusif ..
untuk i dalam 1..5 {
    cetak("Angka:", i)
}

// Menghitung mundur juga didukung secara alami
untuk mundur dalam 3..1 {
    cetak("Hitung mundur:", mundur)
}
```

### 4.2 Iterasi Elemen Array

```celang
buat daftarKota = ["Jakarta", "Bandung", "Surabaya", "Yogyakarta"]

untuk kota dalam daftarKota {
    cetak("Kunjungan ke kota:", kota)
}
```

### 4.3 Iterasi Karakter Teks (*Rune by Rune*)

String di Celang dapat diiterasi langsung karakter per karakter secara aman (*Unicode rune-safe*):

```celang
buat salam = "Kopi ☕"

untuk huruf dalam salam {
    cetak("Karakter:", huruf)
}
```

---

## 5. Penanganan Galat Terstruktur (`coba`, `tangkap`, `akhirnya`, `lempar`)

Celang menyediakan mekanisme penanganan galat tingkat industri menggunakan konstruksi `coba ... tangkap ... akhirnya` (*try-catch-finally*), serta kata kunci `lempar` (*throw*) untuk memicu pengecualian sengaja.

### Menangkap Kesalahan Runtime

Ketika terjadi kesalahan saat eksekusi (seperti pembagian dengan nol), alur program tidak akan langsung mati (*crash*), melainkan dialihkan ke blok `tangkap`:

```celang
buat penyebut = 0

coba {
    buat hasil = 100 / penyebut
    cetak("Hasil pembagian:", hasil)
} tangkap (galat) {
    cetak("Terjadi kesalahan:", galat)
} akhirnya {
    cetak("Operasi hitung selesai diproses.")
}
```

### Melempar Kesalahan Mandiri (`lempar`)

Gunakan kata kunci `lempar` untuk menghentikan alur jika ditemukan kondisi abnormal dalam logika bisnis:

```celang
fungsi tarikSaldo(saldo: Angka, jumlah: Angka) {
    jika jumlah > saldo {
        lempar "Penarikan gagal: Saldo tidak mencukupi!"
    }
    kembalikan saldo - jumlah
}

coba {
    buat sisaSaldo = tarikSaldo(50000, 100000)
    cetak("Sisa saldo:", sisaSaldo)
} tangkap (pesanGalat) {
    cetak("Pemberitahuan Sistem:", pesanGalat)
}
```

### Peran Blok `akhirnya` (*Finally*)

Blok `akhirnya` **selalu dijamin berjalan**, baik blok `coba` berhasil tanpa galat maupun terjadi galat yang tertangkap:

```celang
buat statusKoneksi = "Tersambung"

coba {
    cetak("Mengunduh data...")
    // Simulasi kegagalan
    lempar "Batas waktu koneksi habis (Timeout)"
} tangkap (err) {
    cetak("Log kegagalan:", err)
} akhirnya {
    // Selalu dieksekusi untuk pembersihan sumber daya
    statusKoneksi = "Terputus"
    cetak("Koneksi ditutup secara aman.")
}
```

---

## Ringkasan Alur Kontrol

| Konstruksi | Kata Kunci | Kegunaan Utama |
| :--- | :--- | :--- |
| **Percabangan** | `jika`, `selain jika`, `selain`, `maka` | Evaluasi kondisi boolean dengan blok logika fleksibel. |
| **Pencocokan Kasus** | `cocok`, `kasus`, `lainnya` | Pencocokan nilai eksplisit multikondisi yang rapi. |
| **Perulangan Kondisi**| `selama` | Iterasi berbasis evaluasi kondisi boolean. |
| **Perulangan Koleksi**| `untuk ... dalam` | Iterasi array, rentang `awal..akhir`, string, dan koleksi data. |
| **Penanganan Galat** | `coba`, `tangkap`, `akhirnya`, `lempar` | Menjaga stabilitas aplikasi dari gangguan runtime dan pengecualian. |
