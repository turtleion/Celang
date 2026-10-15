# Fungsi Panah Ringkas (Arrow Functions `=>`)

Terkadang kita hanya ingin membuat rumus perhitungan sederhana tanpa perlu mengetik kata `fungsi`, kurung kurawal, dan kata `kembalikan` yang panjang.

Celang menyediakan tanda **Panah Ringkas (`=>`)** untuk menulis rumus kecil dalam satu tarikan napas.

---

## 💡 Perbandingan Menulis Rumus

Bandingkan dua cara menulis rumus pengali dua berikut:

```celang
// Cara biasa yang panjang:
buat kaliDuaPanjang = fungsi(x) { kembalikan x * 2 }

// Cara cepat dengan tanda panah:
buat kaliDua = x => x * 2
```

Tanda panah `=>` artinya: *"terima masukan di sebelah kiri, lalu langsung berikan hasil perhitungan di sebelah kanan"*. Lebih ringkas, bersih, dan hemat waktu!

---

## 1. Berbagai Bentuk Penulisan

### Satu Masukan (Tanpa Kurung)
Jika rumus hanya menerima satu nilai, Anda tidak perlu memakai tanda kurung:

```celang
buat kuadrat = x => x * x
cetak(kuadrat(7)) // 49
```

### Dua Masukan atau Lebih (Dengan Kurung)
Jika ada dua masukan atau lebih, kumpulkan di dalam tanda kurung `( ... )`:

```celang
buat tambah = (a, b) => a + b
buat hitungLuas = (panjang, lebar) => panjang * lebar

cetak(tambah(10, 25))   // 35
cetak(hitungLuas(5, 4)) // 20
```

### Tanpa Masukan
Jika rumus hanya menghasilkan salam tetap tanpa perlu masukan:

```celang
buat salam = () => "Selamat pagi, kawan!"
cetak(salam()) // "Selamat pagi, kawan!"
```

### Rumus Lebih dari Satu Langkah (`{ ... }`)
Jika Anda membutuhkan beberapa langkah sebelum memberikan hasil akhir, gunakan kurung kurawal `{ ... }` dan kata `kembalikan`:

```celang
buat sapaPengguna = nama => {
    buat ucapan = "Halo, " + nama + "!"
    cetak("Mencatat kunjungan...")
    kembalikan ucapan
}

cetak(sapaPengguna("Andi"))
```

---

## 2. Sangat Nyaman untuk Menyaring dan Mengubah Data

Rumus panah ini terasa paling menyenangkan saat dipadukan dengan alat penyaring data (`saring`) dan alat pengubah data (`petakan`):

```celang
buat nilaiSiswa = [45, 78, 90, 62, 85, 30]

// Ambil nilai yang lulus (>= 60), lalu tambahkan nilai bonus 5 poin:
buat nilaiLulus = nilaiSiswa
    |> saring(skor => skor >= 60)
    |> petakan(skor => skor + 5)

cetak(nilaiLulus) // [83, 95, 67, 90]
```
Kodenya terbaca sangat alami: *saring mereka yang skornya 60 ke atas, lalu petakan masing-masing skor dengan menambahkan 5 poin*.
