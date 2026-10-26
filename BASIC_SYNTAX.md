# Panduan Ringkas Sintaks Dasar Celang (`BASIC_SYNTAX.md`)

Dokumen ini adalah ringkasan aturan penulisan kode dasar dalam bahasa pemrograman Celang. Seluruh sintaks dirancang agar mudah dibaca seperti kalimat sehari-hari.

---

## 1. Komentar (Catatan Penjelas)

Komentar adalah catatan yang Anda tulis untuk dibaca oleh manusia. Celang akan mengabaikan tulisan ini saat program dijalankan:

```celang
// Ini adalah catatan satu baris menggunakan garis miring ganda
# Ini juga catatan satu baris menggunakan tanda pagar

/*
   Ini adalah catatan panjang
   yang memakan beberapa baris sekaligus
*/
```

---

## 2. Variabel & Tipe Data

### 2.1 Menyiapkan Variabel (`buat`)
Gunakan kata `buat` untuk menyiapkan wadah catatan nilai yang isinya bisa diubah kapan saja:

```celang
buat nama = "Budi"
buat umur = 24

// Mengubah isi variabel:
umur = 25
cetak("Umur Budi sekarang:", umur) // 25
```

### 2.2 Nilai Tetap / Konstanta (`konstanta`)
Gunakan kata `konstanta` untuk nilai yang tidak boleh diubah selama program berjalan:

```celang
konstanta PAJAK = 0.11
konstanta KOTA = "Bandung"
```

### 2.3 Jenis Nilai yang Tersedia
- **Angka Bulat (`Angka`)**: `10`, `-5`, `1_000_000`
- **Angka Pecahan (`Pecahan` / `Angka`)**: `3.14`, `-0.5`
- **Teks (`Teks`)**: `"Halo Dunia"`, `'Selamat Pagi'`
- **Pilihan Logika (`Logika`)**: `benar` (true) atau `salah` (false)
- **Nilai Kosong**: `kosong` (nihil / null)

---

## 3. Tanda Hitung & Logika

### 3.1 Perhitungan Matematika
- Penjumlahan & Sambung Kata: `+`
- Pengurangan: `-`
- Perkalian: `*`
- Pembagian: `/`
- Sisa Bagi (Modulo): `%`

### 3.2 Pembanding Nilai
- Sama persis: `==`
- Tidak sama: `!=`
- Lebih besar: `>`
- Lebih kecil: `<`
- Lebih besar atau sama: `>=`
- Lebih kecil atau sama: `<=`

### 3.3 Penghubung Kalimat
- Kedua syarat harus benar: `dan` (atau `&&`)
- Salah satu syarat saja cukup: `atau` (atau `||`)
- Membalikkan keadaan: `bukan` atau `tidak` (atau `!`)

---

## 4. Pengambilan Keputusan (`jika ... selain`)

```celang
buat nilai = 85

jika nilai >= 90 {
    cetak("Istimewa")
} selain jika nilai >= 75 {
    cetak("Bagus")
} selain {
    cetak("Perlu ditingkatkan")
}
```

Bisa juga menggunakan kata `maka` agar lebih alami:
```celang
jika umur >= 17 maka {
    cetak("Boleh membuat KTP")
}
```

---

## 5. Memilih Berdasarkan Pilihan (`cocok ... kasus`)

```celang
buat lampu = "merah"

cocok lampu {
    kasus "merah":
        cetak("Berhenti!")
    kasus "kuning":
        cetak("Hati-hati!")
    kasus "hijau":
        cetak("Jalan terus!")
    lainnya:
        cetak("Warna lampu tidak dikenal.")
}
```

---

## 6. Perulangan

### 6.1 Selama Syarat Terpenuhi (`selama`)
```celang
buat hitungan = 1
selama hitungan <= 3 {
    cetak("Putaran ke-", hitungan)
    hitungan += 1
}
```

### 6.2 Membaca Isi Daftar (`untuk ... dalam`)
```celang
buat daftarNama = ["Andi", "Budi", "Cici"]

untuk nama dalam daftarNama {
    cetak("Halo,", nama)
}

// Menggunakan rentang angka:
untuk angka dalam 1..5 {
    cetak("Angka:", angka)
}
```

---

## 7. Fungsi / Rumus (`fungsi`)

```celang
fungsi hitungLuas(panjang, lebar) {
    kembalikan panjang * lebar
}

buat hasil = hitungLuas(8, 5)
cetak("Luas:", hasil) // 40
```

### Rumus Cepat Panah (`=>`)
```celang
buat kuadrat = x => x * x
cetak(kuadrat(6)) // 36
```

---

## 8. Kumpulan Data: Daftar dan Kamus

```celang
// 1. Daftar berurutan (Array):
buat buah = ["Apel", "Mangga", "Jeruk"]
cetak(buah[0]) // "Apel"

// Menambah atau manipulasi:
buah.tambah("Pisang")

// Transformasi fungsional (menghasilkan array baru / imutabel):
buat hasil = [1, 2, 3, 4]
    .saring(x => x % 2 == 0)          // [2, 4]
    .petakan(x => x * 10)             // [20, 40]
    .reduksi((acc, x) => acc + x, 0)  // 60

// 2. Kamus data berlabel (Map):
buat profil = {
    "nama": "Siti",
    "kota": "Surabaya"
}
cetak(profil.nama) // "Siti"
```

---

## 9. Tiga Cara Menjalankan Perintah (Arsitektur Pemanggilan)

Celang mengelompokkan pemanggilan fungsi menjadi 3 gaya agar rapi dan tidak saling tabrakan:

1. **Fungsi Bawaan Global**: Langsung dipanggil tanpa awalan (misal: `cetak()`, `panjang()`, `tipe()`, `bita()`).
2. **Metode Objek (`.xxx()`)**: Menempel pada data tertentu (misal: `daftar.saring()`, `daftar.petakan()`, `teks.huruf_besar()`).
3. **Modul Ber-namespace**: Dikelompokkan dalam domain pustaka standar yang otomatis tersedia secara global (misal: `berkas.baca()`, `matematika.akar()`, `json.urai()`).

