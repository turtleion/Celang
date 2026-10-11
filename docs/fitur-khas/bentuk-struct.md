# Bentuk (Struct) & Metode Berorientasi Objek

Celang mendukung pemodelan data berorientasi objek yang bersih menggunakan kata kunci `bentuk` (atau `struct`).

---

## 🏗️ Mendefinisikan Bentuk

Bentuk dapat memiliki properti (*fields*) dengan tipe data opsional, nilai default, dan fungsi metode yang terikat.

```celang
bentuk RekeningBank {
    pemilik: Teks,
    saldo itu Angka = 0,

    fungsi setor(jumlah: Angka) {
        diri.saldo += jumlah
        kembalikan diri.saldo
    }

    fungsi tarik(jumlah: Angka) {
        jika jumlah > diri.saldo {
            lempar "Saldo tidak mencukupi!"
        }
        diri.saldo -= jumlah
        kembalikan diri.saldo
    }
}
```

---

## 🎯 Instansiasi Objek

Instansiasi dapat dilakukan melalui konstruktor posisi terurut atau melalui kamus (*Map*):

```celang
// 1. Instansiasi dengan argumen posisi
buat akun1 = RekeningBank("Budi Santoso", 100000)

// 2. Instansiasi dengan kamus (Map properties)
buat akun2 = RekeningBank({
    "pemilik": "Siti Rahma",
    "saldo": 250000
})

cetak(akun1) // RekeningBank{pemilik: Budi Santoso, saldo: 100000}
```

---

## 🧭 Memanggil Metode & Referensi Diri (@properti)

Di dalam metode bentuk, kamu dapat menggunakan:
- `@properti` atau `@metode()` sebagai **shorthand elegan** untuk mengakses anggota instance internal.
- `diri` (idiomatik Indonesia), `this`, atau `self` ketika membutuhkan referensi objek utuh (contoh: `buat ref = diri`).

```celang
bentuk Dompet {
    pemilik: Teks,
    saldo: Angka = 0,

    setor(jumlah: Angka) {
        @saldo += jumlah
        kembalikan @saldo
    }
}
```

---

## 🔒 Enkapsulasi & Visibilitas (`privat`)

Secara default, anggota bentuk bersifat **publik** (*loose mode*). Jika ingin mengamankan field atau metode sensitif dari manipulasi luar, gunakan kata kunci `privat`.

Celang mendukung dua gaya penulisan visibilitas:
1. **Per-baris**: `privat saldo: Angka = 0`
2. **Blok terkelompok**: `privat { ... }`

```celang
bentuk RekeningBank {
    pemilik: Teks, // publik secara default

    privat {
        saldo: Angka = 0,
        pin: Angka = 1234,

        validasiPin(inputPin: Angka) {
            kembalikan inputPin == @pin
        }
    }

    setor(jumlah: Angka) {
        @saldo += jumlah
        kembalikan @saldo
    }

    ambil(jumlah: Angka, inputPin: Angka) {
        jika !@validasiPin(inputPin) {
            lempar "PIN salah!"
        }
        jika jumlah > @saldo {
            lempar "Saldo tidak mencukupi!"
        }
        @saldo -= jumlah
        kembalikan @saldo
    }
}

buat akun = RekeningBank("Budi")
akun.setor(50000)

// ❌ Akses langsung dari luar akan diblokir dengan pesan empati:
// akun.saldo = 999999 // Kesalahan: Properti 'saldo' bersifat privat pada bentuk 'RekeningBank'
```

---

## 🧩 Komposisi Bentuk (Struct Embedding)

Celang mengutamakan prinsip *composition over inheritance*. Kamu dapat menanamkan (*embed*) satu atau lebih bentuk lain ke dalam sebuah bentuk cukup dengan menuliskan nama tipenya:

```celang
bentuk Identitas {
    nama: Teks,
    umur: Angka,

    sapa() {
        kembalikan "Halo, aku " + @nama
    }
}

bentuk Karyawan {
    Identitas, // Di-embed langsung
    gaji: Angka = 5000000,
    divisi: Teks = "Teknologi"
}

buat staf = Karyawan("Dewi", 25, 8000000, "Desain")
cetak(staf.nama)    // "Dewi" (diwarisi dari Identitas)
cetak(staf.sapa())   // "Halo, aku Dewi"
cetak(staf.divisi) // "Desain"
```

---

## 🐑 Kloning Objek (Deep Memory Clone: `salin` / `clone`)

Untuk menduplikasi instance objek atau koleksi secara mendalam tanpa berbagi referensi memori, gunakan fungsi bawaan `salin()` atau `clone()`:

```celang
buat asli = RekeningBank("Ali", 100000)
buat tiruan = salin(asli)

tiruan.pemilik = "Budi"
cetak(asli.pemilik)   // "Ali" (tidak terpengaruh)
cetak(tiruan.pemilik) // "Budi"
```
