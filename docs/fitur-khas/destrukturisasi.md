# Destrukturisasi (Destructuring)

Destrukturisasi memungkinkanmu mengekstrak elemen dari `Array` atau nilai dari `Map` secara langsung ke dalam variabel terpisah dalam satu baris kode yang ringkas.

---

## 📦 Destrukturisasi Array

Ekstrak elemen array berdasarkan urutan posisinya:

```celang
buat koordinat = [120, 250, 480]

// Deklarasi destrukturisasi array
buat [x, y, z] = koordinat

cetak("X:", x) // 120
cetak("Y:", y) // 250
cetak("Z:", z) // 480
```

### Penugasan Ulang (*Re-assignment*)
Kamu juga bisa mengubah variabel yang sudah ada tanpa kata kunci `buat`:

```celang
[x, y] = [999, 888]
cetak("X baru:", x, "Y baru:", y) // 999, 888
```

---

## 🗂️ Destrukturisasi Map / Object

Ekstrak properti kamus berdasarkan nama kuncinya:

```celang
buat dataPegawai = {
    "nama": "Siti Rahma",
    "divisi": "Teknologi",
    "gaji": 15000000
}

// Ekstrak nama dan divisi langsung
buat { nama, divisi } = dataPegawai
cetak(nama, "-", divisi) // Siti Rahma - Teknologi
```

### Mengganti Nama Variabel (*Aliasing*)
Jika ingin menyimpan properti ke nama variabel yang berbeda:

```celang
buat { gaji: nominalGaji } = dataPegawai
cetak("Nominal gaji:", nominalGaji) // 15000000
```
