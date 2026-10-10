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

## 🧭 Memanggil Metode & Referensi Diri

Di dalam metode bentuk, kamu dapat menggunakan kata kunci:
- `diri` (idiomatik Indonesia)
- `this` (idiomatik JavaScript/C++)
- `self` (idiomatik Python/Rust)

Ketiganya merujuk pada instance objek yang sama.

```celang
akun1.setor(50000)
cetak("Saldo setelah setor:", akun1.saldo) // 150000

akun1.tarik(20000)
cetak("Saldo setelah tarik:", akun1.saldo) // 130000
```
