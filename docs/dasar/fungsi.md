# Membuat Rumus & Tindakan Sendiri (Fungsi)

Fungsi adalah sekumpulan perintah yang diberi nama khusus agar bisa Anda panggil dan gunakan berulang kali kapan pun dibutuhkan tanpa harus mengetik ulang seluruh baris kodenya.

---

## 1. Membuat Fungsi Sederhana (`fungsi`)

Gunakan kata kunci `fungsi`, diikuti nama yang Anda inginkan, tanda kurung untuk menerima masukan, dan kurung kurawal `{ ... }` untuk membungkus perintahnya:

```celang
fungsi sapa(nama) {
    cetak("Halo, " + nama + "! Senang bertemu denganmu.")
}

// Memanggil fungsi:
sapa("Andi") // Halo, Andi! Senang bertemu denganmu.
sapa("Dewi") // Halo, Dewi! Senang bertemu denganmu.
```

---

## 2. Mengirimkan Hasil Kembali (`kembalikan`)

Jika fungsi Anda melakukan perhitungan dan ingin memberikan hasilnya kembali ke pemanggil, gunakan kata `kembalikan`:

```celang
fungsi hitungLuasPersegi(panjang, lebar) {
    kembalikan panjang * lebar
}

buat luasKamar = hitungLuasPersegi(10, 5)
cetak("Luas kamar:", luasKamar, "meter") // 50 meter
```

> **Berhenti Seketika:** Begitu kata `kembalikan` dijalankan, fungsi akan langsung selesai saat itu juga dan perintah di bawahnya tidak akan dijalankan lagi.

---

## 3. Masukan dengan Nilai Cadangan Bawaan

Anda dapat menyiapkan nilai cadangan jika pemanggil tidak memberikan masukan:

```celang
fungsi sambut(nama, sapaan = "Halo") {
    kembalikan sapaan + ", " + nama + "!"
}

cetak(sambut("Rian"))                  // Halo, Rian!
cetak(sambut("Siti", "Selamat Pagi"))  // Selamat Pagi, Siti!
```

---

## 4. Mengubah Variabel Asli secara Langsung (`rujuk`)

Secara umum, saat Anda mengirimkan angka ke dalam fungsi, fungsi hanya menerima salinan angka tersebut. Jika Anda ingin fungsi **langsung mengubah isi variabel asli** di luar, gunakan kata `rujuk`:

```celang
fungsi beriBonus(rujuk saldo, bonus) {
    saldo += bonus // Langsung mengubah variabel asli di luar
}

buat saldoTabungan = 100000
beriBonus(saldoTabungan, 25000)

cetak("Saldo setelah bonus:", saldoTabungan) // 125000
```

> **Perhatian:** Masukan untuk parameter `rujuk` harus berupa nama variabel (seperti `saldoTabungan`), bukan angka langsung seperti `100000`.

---

## 5. Menitipkan Pesan Pembersihan (`saat_selesai`)

Anda dapat memastikan bahwa suatu tindakan perapian selalu dikerjakan sebelum fungsi selesai (misalnya menutup dokumen atau koneksi):

```celang
fungsi layaniPelanggan(nama) {
    cetak("1. Membuka loket untuk:", nama)
    saat_selesai {
        cetak("3. Menutup loket kembali (selalu dijalankan!)")
    }
    cetak("2. Memproses transaksi...")
}

layaniPelanggan("Andi")
```

---

## 6. Rumus Singkat Panah (`=>`)

Untuk rumus kecil yang hanya satu baris, Anda bisa menyingkatnya dengan tanda panah:

```celang
buat kuadrat = x => x * x
cetak(kuadrat(9)) // 81
```
