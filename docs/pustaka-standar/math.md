---
outline: deep
---

# Pustaka `math` (`matematika`)

> [!NOTE]
> **Modul `math`** adalah alias internasional resmi untuk modul **`matematika`** di Celang. Seluruh fungsi, parameter, konstanta, dan perilakunya 100% identik.
> 
> Dokumentasi lengkap dan panduan menyeluruh dapat dibaca di:
> 👉 **[Pustaka Perhitungan Matematika (`matematika` / `math`)](/pustaka-standar/matematika)**

## Contoh Penggunaan Cepat

```celang
// 1. Fitur Cerdas Finansial & Ejaan:
buat harga = math.rupiah(1500000)
buat ejaan = math.terbilang(1500000)
cetak(harga) // "Rp 1.500.000"
cetak(ejaan) // "satu juta lima ratus ribu"

// 2. Operasi Matematika & Statistik:
buat akar = math.sqrt(144)         // 12
buat pangkat = math.pow(2, 8)       // 256
buat rata = math.average([10, 20, 30]) // 20
buat jarak = math.distance(0, 0, 3, 4) // 5

// 3. Trigonometri & Konstanta:
buat sin90 = math.sin(math.to_radians(90))
cetak("Pi:", math.PI)
cetak("Euler e:", math.E)
```

Untuk tabel fungsi lengkap dan penjelasan detail setiap fungsi, silakan buka **[Dokumentasi Lengkap Matematika](/pustaka-standar/matematika)**.
