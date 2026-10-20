# Penanganan Galat (Try-Catch)

Dalam dunia nyata, program komputer terkadang berhadapan dengan situasi yang tidak terduga: pembagian dengan angka nol, berkas yang tidak ditemukan di disk, atau jaringan internet yang tiba-tiba terputus.

Daripada membiarkan program Anda berhenti mendadak (*crash*), Celang menyediakan mekanisme penanganan galat yang tangguh, ramah, dan alami.

---

## 1. Menangkap Kendala (`coba` & `tangkap`)

Gunakan kata kunci `coba` untuk menjalankan perintah yang berisiko mengalami kesalahan, dan siapkan blok `tangkap` untuk mengantisipasi masalah tersebut jika benar-benar terjadi:

```celang
fungsi bagiUang(total, jumlahOrang) {
    jika jumlahOrang <= 0 {
        lempar "Jumlah orang tidak boleh nol atau negatif!"
    }
    kembalikan total / jumlahOrang
}

coba {
    buat bagian = bagiUang(100000, 0)
    cetak("Masing-masing mendapat:", bagian)
} tangkap (kendala) {
    cetak("Terjadi kendala perhitungan:", kendala)
}
```

Keluaran:
```text
Terjadi kendala perhitungan: Jumlah orang tidak boleh nol atau negatif!
```

Program tetap berjalan dengan lancar hingga baris terakhir tanpa terhenti paksa!

---

## 2. Jaminan Pelaksanaan Pasti (`akhirnya`)

Terkadang ada proses penting yang **harus selalu dikerjakan**, baik saat perhitungan berhasil maupun saat terjadi kesalahan (misalnya: menutup sambungan berkas atau mematikan animasi loading). Gunakan blok `akhirnya`:

```celang
coba {
    cetak("1. Membuka koneksi data...")
    buat hasil = 10 / 0 // Memicu galat pembagian nol
    cetak("2. Hasil:", hasil)
} tangkap (galat) {
    cetak("⚠️ Terjadi masalah:", galat)
} akhirnya {
    cetak("3. Menutup koneksi data. (Pasti dieksekusi)")
}
```

Keluaran:
```text
1. Membuka koneksi data...
⚠️ Terjadi masalah: pembagian dengan nol (zero division)
3. Menutup koneksi data. (Pasti dieksekusi)
```

---

## 3. Melempar Kesalahan Mandiri (`lempar`)

Jika fungsi yang Anda buat menerima masukan yang tidak masuk akal atau melanggar aturan bisnis, Anda dapat sengaja memunculkan galat dengan kata kunci `lempar`:

```celang
fungsi daftarkanPengguna(nama, umur) {
    jika panjang(nama) == 0 {
        lempar "Nama pengguna wajib diisi!"
    }
    jika umur < 13 {
        lempar "Usia minimal pendaftaran adalah 13 tahun."
    }

    cetak("Pengguna ${nama} berhasil didaftarkan.")
}

coba {
    daftarkanPengguna("Budi", 10)
} tangkap (pesan) {
    cetak("Pendaftaran ditolak:", pesan)
}
```

---

## 4. Penanganan Galat Ringkas & Pencegahan Dini

Selain blok `coba...tangkap`, Celang juga menyediakan fitur ergonomis untuk mencegah terjadinya galat sejak awal:

### A. Navigasi Aman (`mungkin`)
Cegah kesalahan saat membaca properti bertingkat yang berpotensi `kosong` (*null*):

```celang
buat pengguna = {"nama": "Budi"}

// Tanpa 'mungkin', mengakses pengguna.alamat.kota bisa membuat program gagal
buat kota = mungkin pengguna.alamat.kota
cetak(kota) // kosong (tetap aman tanpa error!)
```

### B. Nilai Cadangan Bawaan (`kalau kosong`)
Siapkan nilai pengganti jika data bernilai `kosong`:

```celang
buat konfigurasiPort = kosong
buat portAktif = konfigurasiPort kalau kosong 8080
cetak(portAktif) // 8080
```

### C. Pembersihan Otomatis (`saat_selesai`)
Pastikan perintah pembersihan selalu dieksekusi sebelum sebuah fungsi selesai:

```celang
fungsi prosesBerkas() {
    cetak("Buka berkas")
    saat_selesai {
        cetak("Tutup berkas (dijamin dipanggil saat fungsi selesai)")
    }

    cetak("Membaca isi berkas...")
}
```

---

## 📋 Ringkasan Kata Kunci

| Kata Kunci | Padanan Global | Peran |
| :--- | :--- | :--- |
| **`coba`** | `try` | Membungkus perintah yang rawan mengalami galat |
| **`tangkap (err)`** | `catch (err)` | Menangkap dan menangani pesan galat jika terjadi masalah |
| **`akhirnya`** | `finally` | Blok penutup yang dijamin pasti selalu dijalankan |
| **`lempar`** | `throw` | Menerbitkan atau melempar galat secara sengaja |
| **`saat_selesai`** | `defer` | Menunda eksekusi perintah pembersihan hingga fungsi selesai |
| **`mungkin`** | `?.` | Membaca properti bertingkat secara aman dari nilai kosong |
