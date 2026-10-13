# Fungsi & Fungsi Anonim

Fungsi di Celang adalah *first-class citizen* — Anda dapat mendeklarasikan fungsi bernama, membuat fungsi anonim (*closures/lambda*), menyimpannya ke dalam variabel, meneruskannya sebagai argumen, hingga merangkaikannya secara berurutan menggunakan pipeline operator (`|>`).

---

## 1. Deklarasi Fungsi Bernama

Gunakan kata kunci `fungsi` diikuti nama fungsi, daftar parameter di dalam tanda kurung `(...)`, dan blok kode di dalam kurung kurawal `{...}`:

```celang
fungsi sapa(nama) {
    cetak("Halo, " + nama + "! Selamat datang di Celang.")
}

// Memanggil fungsi
sapa("Andi")
sapa("Dewi")
```

---

## 2. Mengembalikan Nilai (`kembalikan`)

Gunakan kata kunci `kembalikan` (*return*) untuk mengirimkan hasil kalkulasi atau data keluar dari fungsi:

```celang
fungsi hitungLuasPersegi(panjang, lebar) {
    kembalikan panjang * lebar
}

buat luasRuangan = hitungLuasPersegi(10, 5)
cetak("Luas ruangan:", luasRuangan, "m²")
```

::: tip Return Dini (*Early Return*)
Eksekusi di dalam fungsi akan langsung berhenti begitu pernyataan `kembalikan` dieksekusi.
```celang
fungsi cekKelulusan(nilai) {
    jika nilai < 70 {
        kembalikan "Belum Lulus"
    }
    kembalikan "Lulus"
}
```
:::

---

## 3. Anotasi Tipe Parameter & Nilai Kembalian

Untuk meningkatkan dokumentasi kode, keterbacaan, dan keamanan tipe, Celang mendukung anotasi tipe eksplisit pada parameter maupun nilai kembalian fungsi.

### Tipe Kembalian (`mengembalikan` atau `->`)

Anda dapat menuliskan tipe kembalian menggunakan kata kunci alami `mengembalikan` atau simbol panah `->`:

```celang
// Menggunakan kata kunci alami 'mengembalikan'
fungsi tambah(a: Angka, b: Angka) mengembalikan Angka {
    kembalikan a + b
}

// Menggunakan simbol panah '->'
fungsi gandakan(x: Angka) -> Angka {
    kembalikan x * 2
}

cetak("Hasil tambah:", tambah(15, 25))  // 40
cetak("Hasil ganda:", gandakan(7))      // 14
```

---

## 4. Parameter dengan Nilai Bawaan (*Default Values*)

Parameter dapat memiliki nilai bawaan (*default value*) yang akan otomatis digunakan apabila pemanggil tidak menyertakan argumen tersebut:

```celang
fungsi sambut(nama: Teks, sapaan: Teks = "Halo") {
    kembalikan sapaan + ", " + nama + "!"
}

cetak(sambut("Rian"))                   // Halo, Rian!
cetak(sambut("Siti", "Selamat Siang"))  // Selamat Siang, Siti!
```

---

## 5. Parameter Rujukan (`rujuk` / Pass-by-Reference)

Secara bawaan, Celang melewatkan argumen tipe primitif berdasarkan nilai (*pass-by-value*). Jika Anda ingin fungsi memodifikasi variabel asli di luar *scope*, gunakan kata kunci `rujuk` (*pass-by-reference*):

```celang
fungsi tambahBonus(rujuk saldo: Angka, bonus: Angka) {
    saldo += bonus
}

buat saldoTabungan = 100000

// Melewatkan variabel ke parameter rujukan
tambahBonus(saldoTabungan, 25000)

cetak("Saldo setelah bonus:", saldoTabungan) // 125000
```

::: warning Aturan Parameter Rujukan
Argumen yang diberikan ke parameter `rujuk` **harus berupa nama variabel** (misalnya `saldoTabungan`), bukan berupa nilai langsung atau ekspresi hitungan (seperti `100` atau `saldo + 5`).
:::

---

## 6. Eksekusi Tertunda (`saat_selesai` / `defer`)

Kata kunci `saat_selesai` (atau `defer`) menjadwalkan suatu instruksi untuk dieksekusi **tepat sebelum fungsi selesai berjalan**, baik fungsi keluar secara normal maupun melalui `kembalikan` dini:

```celang
fungsi prosesTransaksi(nomorAkun: Teks) {
    cetak("1. Membuka koneksi database untuk akun:", nomorAkun)

    // Dijamin berjalan saat fungsi selesai
    saat_selesai {
        cetak("4. Menutup koneksi database dengan aman.")
    }

    cetak("2. Memverifikasi saldo...")
    cetak("3. Melakukan debit saldo...")
}

prosesTransaksi("ACC-9921")
```

### Urutan Eksekusi LIFO (*Last-In, First-Out*)

Jika terdapat lebih dari satu `saat_selesai` dalam satu fungsi, eksekusi dilakukan secara terbalik (instruksi yang didaftarkan terakhir akan dijalankan lebih dulu):

```celang
fungsi ujiUrutanDefer() {
    saat_selesai cetak("Pembersihan Tahap 1")
    saat_selesai cetak("Pembersihan Tahap 2")
    cetak("Inti Operasi")
}

ujiUrutanDefer()
```

### Output:
```text
Inti Operasi
Pembersihan Tahap 2
Pembersihan Tahap 1
```

---

## 7. Multi-Return Nilai & Destrukturisasi

Sebuah fungsi dapat menghasilkan lebih dari satu nilai sekaligus dengan mengembalikan array atau struktur data:

```celang
fungsi bagiDanSisa(angka: Angka, pembagi: Angka) {
    buat hasilBagi = angka / pembagi
    buat sisaBagi = angka % pembagi
    kembalikan [hasilBagi, sisaBagi]
}

// Membongkar hasil kembalian langsung ke dalam variabel
buat [hasil, sisa] = bagiDanSisa(14, 4)

cetak("Hasil bagi:", hasil) // 3
cetak("Sisa bagi:", sisa)   // 2
```

---

## 8. Fungsi Anonim, Closure & Pipeline

### Fungsi Anonim

Fungsi dapat dideklarasikan tanpa nama (*anonymous function* / *lambda*) dan disimpan langsung ke dalam variabel:

```celang
buat kuadratkan = fungsi(x: Angka) -> Angka {
    kembalikan x * x
}

cetak("5 dikuadratkan =", kuadratkan(5)) // 25
```

### Penutupan Lingkungan (*Closure*)

Fungsi anonim mampu mengingat dan mengakses variabel dari lingkungan (*lexical scope*) induknya:

```celang
fungsi buatPencatat(awal: Angka) {
    buat hitungan = awal

    kembalikan fungsi() {
        hitungan += 1
        kembalikan hitungan
    }
}

buat counter = buatPencatat(10)
cetak(counter()) // 11
cetak(counter()) // 12
cetak(counter()) // 13
```

### Integrasi dengan Pipeline Operator (`|>`)

Fungsi yang menerima satu argumen dapat dirangkaikan secara berurutan dengan sangat elegan menggunakan operator `|>`:

```celang
fungsi tambahDua(n) { kembalikan n + 2 }
fungsi kalikanTiga(n) { kembalikan n * 3 }
fungsi jadikanTeks(n) { kembalikan "Hasil kalkulasi: " + n }

buat hasilAkhir = 4 |> tambahDua |> kalikanTiga |> jadikanTeks
cetak(hasilAkhir) // Hasil kalkulasi: 18
```

---

## 9. Fungsi Rekursif

Fungsi yang memanggil dirinya sendiri didukung penuh untuk menyelesaikan komputasi matematis dan struktur data berulang:

```celang
fungsi faktorial(n: Angka) mengembalikan Angka {
    jika n <= 1 {
        kembalikan 1
    }
    kembalikan n * faktorial(n - 1)
}

cetak("Faktorial dari 5:", faktorial(5)) // 120
```

---

## Ringkasan Fitur Fungsi

| Fitur | Sintaks Contoh | Manfaat |
| :--- | :--- | :--- |
| **Deklarasi Dasar** | `fungsi sapa(nama) { ... }` | Pembuatan subrutin modular yang dapat dipanggil kembali. |
| **Nilai Kembali** | `kembalikan nilai` | Mengirim data hasil eksekusi ke pemanggil. |
| **Anotasi Tipe** | `(x: Angka) mengembalikan Angka` atau `-> Teks` | Kejelasan tipe dan dokumentasi kode mandiri. |
| **Nilai Bawaan** | `(sapaan: Teks = "Halo")` | Argumen opsional yang fleksibel. |
| **Parameter Rujukan** | `(rujuk x: Angka)` | Modifikasi langsung variabel luar (*pass-by-reference*). |
| **Cleanup Tertunda** | `saat_selesai { ... }` | Penjaminan pembersihan sumber daya (LIFO). |
| **Multi-Return** | `buat [a, b] = bagi(x, y)` | Pengembalian banyak nilai sekaligus lewat destrukturisasi. |
| **Pipeline Ready** | `data \|> filter \|> format` | Aliran pemrosesan data linier yang bersih. |
