# Fungsi & Fungsi Anonim

Fungsi di Celang adalah *first-class citizen* — Anda dapat mendeklarasikan fungsi bernama, membuat fungsi anonim (*closures/lambda*), menyimpannya ke dalam variabel, serta meneruskannya sebagai argumen ke fungsi lain.

---

## 1. Deklarasi Fungsi Bernama

Gunakan kata kunci `fungsi` diikuti nama fungsi, parameter di dalam tanda kurung `(...)`, dan blok kode di dalam kurung kurawal `{...}`.

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

## 3. Parameter dengan Tipe Data & Nilai Bawaan

Anda dapat menentukan tipe data parameter serta nilai default untuk meningkatkan keterbacaan:

```celang
fungsi sambut(nama: Teks, sapaan: Teks = "Halo") {
    kembalikan "{sapaan}, {nama}!"
}

cetak(sambut("Rian"))                   // Halo, Rian!
cetak(sambut("Siti", "Selamat Siang"))  // Selamat Siang, Siti!
```

---

## 4. Fungsi Anonim & Closure

Fungsi anonim adalah fungsi tanpa nama yang dapat disimpan langsung di dalam variabel atau dikirimkan sebagai argumen fungsi tingkat tinggi (*higher-order functions*).

```celang
// Menyimpan fungsi ke variabel
buat kuadratkan = fungsi(x) {
    kembalikan x * x
}

cetak("4 dikuadratkan =", kuadratkan(4)) // 16
```

### Kemampuan Penutupan (*Closure*)

Fungsi anonim dapat mengingat dan mengakses variabel dari lingkungan (*scope*) tempat ia diciptakan:

```celang
fungsi buatPenghitung(nilaiAwal) {
    buat hitungan = nilaiAwal

    kembalikan fungsi() {
        hitungan += 1
        kembalikan hitungan
    }
}

buat counterA = buatPenghitung(0)
cetak(counterA()) // 1
cetak(counterA()) // 2
cetak(counterA()) // 3
```

---

## 5. Fungsi Rekursif

Celang mendukung fungsi rekursif, yaitu fungsi yang memanggil dirinya sendiri:

```celang
fungsi faktorial(n) {
    jika n <= 1 {
        kembalikan 1
    }
    kembalikan n * faktorial(n - 1)
}

cetak("Faktorial dari 5 adalah:", faktorial(5)) // 120
```
