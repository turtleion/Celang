# Panduan Sintaks Lanjutan Celang (`ADVANCED_SYNTAX.md`)

Dokumen ini menjelaskan fitur-fitur lanjutan di bahasa pemrograman Celang untuk membangun aplikasi yang terstruktur, aman, dan mudah dipelihara.

---

## 1. Sistem Modul & Hak Akses Berkas (Visibilitas)

Ketika program Anda bertambah besar dan dipecah ke dalam beberapa berkas, Celang menjaga agar data tidak saling bertabrakan menggunakan aturan hak akses:

### 1.1 Aturan Bawaan: Rahasia untuk Berkas Sendiri (`privat!`)
Secara bawaan, semua variabel, fungsi, dan bentuk data yang Anda buat di dalam suatu berkas **hanya bisa digunakan oleh berkas itu sendiri**. Berkas lain tidak dapat melihatnya kecuali Anda secara sengaja mengizinkannya.

### 1.2 Membuka untuk Umum (`publik!` atau `export`)
Sematkan kata `publik!` jika Anda ingin suatu fungsi atau data dapat diimpor dan digunakan oleh berkas mana pun:

```celang
// berkas: matematika.ce
publik! fungsi tambah(a, b) {
    kembalikan a + b
}

// Fungsi ini tetap rahasia, hanya untuk berkas matematika.ce:
fungsi rumusRahasia() {
    kembalikan 42
}
```

### 1.3 Terbuka Hanya untuk Modul yang Sama (`privat-modul!`)
Data hanya boleh diakses oleh berkas-berkas yang berada di dalam folder modul yang sama (folder yang memiliki berkas penanda `celang.mod`):

```celang
privat-modul! konstanta KUNCI_MODUL = "SEC-123"
```

### 1.4 Memuat Berkas Lain (`impor` / `ambil`)
```celang
// Memuat fungsi dari berkas lain:
ambil tambah dari "./matematika.ce"

buat hasil = tambah(10, 20)
cetak(hasil) // 30
```

---

## 2. Bentuk Data Kustom (`bentuk`) & Hak Akses Privat

Anda dapat mengelompokkan data dan tindakan ke dalam cetak biru yang rapi menggunakan kata `bentuk`:

```celang
bentuk RekeningBank {
    pemilik: Teks,

    // Data rahasia yang tidak boleh diubah sembarangan dari luar:
    privat {
        saldo: Angka = 0,
        pin: Angka = 1234
    }

    setor(jumlah: Angka) {
        @saldo += jumlah
        kembalikan @saldo
    }

    tarik(jumlah: Angka, inputPin: Angka) {
        jika inputPin != @pin {
            lempar "PIN salah!"
        }
        jika jumlah > @saldo {
            lempar "Saldo tidak mencukupi!"
        }
        @saldo -= jumlah
        kembalikan @saldo
    }
}

buat akun = RekeningBank("Budi")
akun.setor(50000)
```
Tanda `@` di dalam metode bentuk adalah jalan pintas untuk menyebut data milik dirinya sendiri (misalnya `@saldo`).

---

## 3. Parameter Rujukan Langsung (`rujuk`)

Secara umum, angka yang dikirim ke fungsi hanyalah salinan. Gunakan kata `rujuk` jika Anda ingin fungsi **langsung mengubah isi variabel asli** di luar:

```celang
fungsi tambahPoin(rujuk skor, tambahan) {
    skor += tambahan
}

buat skorPemain = 100
tambahPoin(skorPemain, 50)
cetak(skorPemain) // 150 (variabel asli langsung bertambah)
```

---

## 4. Penanganan Kendala (`coba`, `tangkap`, `akhirnya`, `lempar`)

Mengamankan program agar tidak mendadak mati saat terjadi situasi yang tidak diinginkan:

```celang
fungsi bagiAngka(a, b) {
    jika b == 0 {
        lempar "Tidak bisa membagi dengan angka nol!"
    }
    kembalikan a / b
}

coba {
    buat hasil = bagiAngka(100, 0)
    cetak(hasil)
} tangkap (kendala) {
    cetak("Pemberitahuan kendala:", kendala)
} akhirnya {
    cetak("Pemeriksaan selesai.")
}
```

---

## 5. Pembersihan Terjadwal (`saat_selesai` / `defer`)

Menjamin sebuah perintah selalu dijalankan tepat sebelum fungsi selesai (misalnya menutup dokumen atau membereskan sambungan):

```celang
fungsi olahData() {
    cetak("1. Membuka berkas...")
    saat_selesai {
        cetak("3. Menutup berkas dengan aman.")
    }
    cetak("2. Memproses isi berkas...")
}

olahData()
```

---

## 6. Penyaluran Berantai (`|>`)

Mengalirkan hasil perhitungan dari langkah sebelumnya ke langkah berikutnya seperti ban berjalan:

```celang
buat angka = [1, 2, 3, 4, 5, 6]

buat hasil = angka
    |> saring(x => x % 2 == 0)
    |> petakan(x => x * 10)

cetak(hasil) // [20, 40, 60]
```

---

## 7. Membongkar Paket Data (Destrukturisasi)

Mengambil data dari daftar atau kamus langsung ke variabel terpisah dalam satu baris:

```celang
// 1. Membongkar Daftar (Array):
buat [pertama, kedua] = [10, 20]
cetak(pertama, kedua) // 10 20

// 2. Membongkar Kamus (Map):
buat data = { "nama": "Dewi", "divisi": "Desain" }
buat { nama, divisi } = data
cetak(nama, "bekerja di", divisi) // Dewi bekerja di Desain
```

---

## 8. Mengubah Gaya Bahasa di Tengah Kode (`gaya "..." { }`)

Menyisipkan beberapa baris kode dengan gaya bahasa lain di tengah program utama:

```celang
buat nama = "Budi"

// Meminjam gaya JavaScript untuk satu blok kecil:
gaya "javascript" {
    let statusAktif = true
    function cek(user) {
        return user == "Budi" && statusAktif
    }
    let hasilCek = cek(nama)
}

// Kembali ke gaya Bahasa Indonesia:
jika hasilCek {
    cetak("Selamat datang,", nama)
}
```
