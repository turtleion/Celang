# Mengubah Gaya di Tengah Kode (`gaya "..." { }`)

Selain menentukan gaya bahasa untuk satu berkas utuh, Celang memiliki fitur yang sangat unik: **Mengubah Gaya Bahasa di Tengah-Tengah Program**.

Fitur ini memungkinkan Anda berpindah gaya bahasa untuk beberapa baris kode tertentu menggunakan blok `gaya "<nama_gaya>" { ... }`.

---

## 💡 Kapan Ini Bermanfaat?

Bayangkan situasi berikut:
1. Anda sedang menulis program utama menggunakan Bahasa Indonesia yang nyaman.
2. Anda menemukan cuplikan rumus yang bagus dari internet yang ditulis dalam gaya JavaScript atau Python.
3. Alih-alih harus menerjemahkan kata demi kata secara manual, Anda cukup membungkus rumus tersebut di dalam blok `gaya "javascript"` atau `gaya "pythonic"`.
4. Variabel dan fungsi yang Anda buat di dalam blok tersebut **tetap bisa dipakai langsung** oleh bagian program lainnya!

---

## 📝 Contoh Nyata

Berikut contoh program yang memadukan Bahasa Indonesia, JavaScript, dan Python dalam satu berkas:

```celang
// 1. Bagian Utama Ditulis dalam Bahasa Indonesia
buat namaPengguna = "Budi"
cetak("Nama Pengguna:", namaPengguna)

// 2. Meminjam Gaya JavaScript untuk Memeriksa Akses
gaya "javascript" {
    let statusLogin = true
    function cekAkses(user) {
        return user == "Budi" && statusLogin
    }
    let diizinkan = cekAkses(namaPengguna)
}

// 3. Kembali ke Bahasa Indonesia (Variabel 'diizinkan' tetap bisa dipakai!)
jika diizinkan {
    cetak("Akses Diberikan: Selamat datang,", namaPengguna)
} selain {
    cetak("Akses Ditolak!")
}

// 4. Meminjam Gaya Python untuk Rumus Bunga
gaya "pythonic" {
    def hitungBunga(pokok, persen): {
        return pokok * persen / 100
    }
    bonus = hitungBunga(1000000, 10)
}

cetak("Bonus Akhir:", bonus) // 100000
```

---

## ⚙️ Cara Kerja di Balik Layar

Ketika Celang menemukan perintah `gaya "nama_gaya" { ... }`:

1. **Berganti Kacamata Baca**: Celang sementara waktu mengganti kamus kata kuncinya menjadi kamus gaya yang diminta (misalnya mengenali kata `let`, `function`, atau `def`).
2. **Menyimpan ke Logika Murni yang Sama**: Kata `let` atau `def` disimpan ke dalam pohon logika yang sama persis seperti `buat` dan `fungsi`.
3. **Kembali Seperti Semula**: Begitu kurung kurawal penutup `}` tercapai, Celang kembali menggunakan kacamata baca bahasa utama.
4. **Variabel Saling Terhubung**: Semua variabel yang dibuat di dalam blok tersebut tetap hidup dan dapat dibaca oleh bagian kode lainnya tanpa terisolasi.
