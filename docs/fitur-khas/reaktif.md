# ⚡ Variabel Reaktif & Pengamat (`reaktif` & `pantau`)

Celang mengintegrasikan paradigma **Reaktivitas Tingkat Bahasa** (*First-Class Language-Level Reactivity*) tanpa memerlukan pustaka pihak ketiga seperti RxJS, MobX, atau Vue Reactivity.

Dengan kata kunci `reaktif` dan `pantau`, variabel Celang dapat saling terhubung seperti sel pada spreadsheet atau mengalirkan sinyal pemantauan secara otomatis saat nilai berubah.

---

## 🚀 1. Variabel Reaktif (`reaktif`)

Variabel yang dideklarasikan dengan `reaktif` adalah *derived state* (nilai turunan terkomputasi). Variabel ini merekam secara otomatis variabel apa saja yang dibacanya. Setiap kali variabel penyusunnya berubah, nilainya diperbarui seketika!

### Sintaks Dasar

```celang
reaktif <nama> = <ekspresi>
```

### Contoh Penggunaan

```celang
buat harga = 50.ribu
buat jumlah = 2

// total otomatis terikat pada 'harga' dan 'jumlah'
reaktif total = harga * jumlah

cetak(total) // 100000

// Saat harga atau jumlah diubah, total terupdate otomatis!
harga = 75.ribu
cetak(total) // 150000

jumlah = 4
cetak(total) // 300000
```

---

## 👁️ 2. Pengamat Nilai (`pantau`)

Blok `pantau` memungkinkan Anda mendaftarkan pengamat (*observer/watcher*) terhadap variabel mana pun. Kapan pun variabel tersebut menerima nilai baru (baik melalui penugasan biasa atau pembaruan reaktif), pengamat langsung dieksekusi.

### Sintaks Dasar

```celang
pantau <variabel_target> (baru, lama) {
    // Dipanggil otomatis saat <variabel_target> berubah nilainya
}
```

Jika nama parameter dihilangkan, variabel bawaan `baru` dan `lama` tetap dapat diakses secara implisit.

### Contoh Penggunaan

```celang
buat saldo = 500.ribu

pantau saldo (baru, lama) {
    cetak("Notifikasi Bank: Saldo Anda berubah dari " + lama.rupiah + " menjadi " + baru.rupiah)
}

saldo = 350.ribu
// Output: Notifikasi Bank: Saldo Anda berubah dari Rp 500.000 menjadi Rp 350.000

saldo = 1.juta
// Output: Notifikasi Bank: Saldo Anda berubah dari Rp 350.000 menjadi Rp 1.000.000
```

---

## 🔮 3. Sinergi: Reaktif + Pantau + Transaksi

Ketiga fitur eksklusif Celang saling melengkapi secara mulus:
1. Anda membuat formula reaktif: `reaktif pajak = total * 0.11`
2. Anda memasang pengamat: `pantau pajak (baru, lama) { ... }`
3. Jika mutasi terjadi di dalam `transaksi` yang kemudian di-`batalkan`, nilai reaktif otomatis kembali ke semula dan pengamat dipicu dengan nilai pemulihan tersebut!

```celang
buat harga = 100.ribu
reaktif ppn = harga * 0.11

pantau ppn (baru, lama) {
    cetak("Perubahan PPN:", lama, "->", baru)
}

transaksi {
    harga = 200.ribu // memicu: Perubahan PPN: 11000 -> 22000
    batalkan "promo dibatalkan" // memicu rollback: Perubahan PPN: 22000 -> 11000
}

cetak(harga.rupiah) // Rp 100.000
cetak(ppn)          // 11000
```
