# Bentuk Data Kustom (`bentuk`)

Ketika Anda membuat program yang menangani data dunia nyata—seperti data rekening bank, data murid, atau barang belanjaan—menggunakan variabel terpisah-pisah bisa membuat kode berantakan.

Di Celang, Anda dapat mengelompokkan informasi yang saling berkaitan ke dalam satu cetak biru yang rapi menggunakan kata kunci `bentuk`.

---

## 🏗️ Membuat Cetak Biru Bentuk

Sebuah `bentuk` dapat memiliki daftar informasi (properti) dan daftar tindakan (fungsi/metode) yang bisa dilakukannya:

```celang
bentuk RekeningBank {
    pemilik: Teks,
    saldo: Angka = 0,

    fungsi setor(jumlah: Angka) {
        @saldo += jumlah
        kembalikan @saldo
    }

    fungsi tarik(jumlah: Angka) {
        jika jumlah > @saldo {
            lempar "Saldo tidak mencukupi!"
        }
        @saldo -= jumlah
        kembalikan @saldo
    }
}
```

---

## 🎯 Mengisi Data ke dalam Bentuk

Setelah cetak birunya siap, Anda dapat membuat objek data nyata berdasarkan cetak biru tersebut:

```celang
// Cara 1: Mengisi data berdasarkan urutan
buat akun1 = RekeningBank("Budi Santoso", 100000)

// Cara 2: Mengisi data dengan menyebutkan nama kolomnya
buat akun2 = RekeningBank({
    "pemilik": "Siti Rahma",
    "saldo": 250000
})

cetak(akun1) // RekeningBank{pemilik: Budi Santoso, saldo: 100000}
```

---

## 🧭 Memanggil Tindakan, Tanda `@`, dan Kata Kunci `diri`

Di dalam fungsi yang ada di dalam bentuk, Anda dapat menggunakan tanda `@` di depan nama data sebagai jalan pintas untuk menyebut data milik dirinya sendiri:

- `@saldo` artinya sama dengan `diri.saldo` (data saldo milik akun ini).
- `@setor(10000)` artinya memanggil tindakan `setor` milik akun ini (`diri.setor(10000)`).

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

### 🛡️ Aturan Anti-Ambiguitas: Kapan Menggunakan `diri` vs `@`?

Tanda `@` didesain khusus sebagai jalan pintas (*syntactic sugar*) untuk mengakses **properti atau metode** milik objek ini (`@nama`, `@saldo`, `@hitung()`). Karena itu, tanda `@` **wajib** diikuti nama anggota dan **tidak boleh berdiri sendirian**.

Jika Anda ingin merujuk, mengembalikan, atau mengirim **seluruh objek ini** ke fungsi lain (seperti `this` atau `self` pada bahasa lain), Anda **wajib menggunakan kata kunci `diri`**, bukan tanda `@` sendirian!

```celang
fungsi catatAktivitas(objekRekening, namaNasabah) {
    cetak("Mencatat aktivitas untuk:", namaNasabah)
}

bentuk RekeningBank {
    pemilik: Teks,
    saldo: Angka = 0,

    fungsi proses() {
        // ✅ BENAR: Mengirim objek diri seutuhnya menggunakan kata kunci 'diri'
        catatAktivitas(diri, @pemilik)

        // ❌ SALAH / Ditolak Celang:
        // catatAktivitas(@, @pemilik) -> Celang akan menolak tanda '@' tanpa nama anggota!
    }

    fungsi ambilObjekIni() {
        kembalikan diri // ✅ BENAR, bukan 'kembalikan @'
    }
}
```

---

## 🔒 Menjaga Data Sensitif (`privat`)

Secara umum, data di dalam bentuk dapat dilihat dan diubah oleh siapa saja. Namun, jika ada informasi sensitif (seperti nomor PIN atau saldo rahasia) yang tidak boleh diubah sembarangan dari luar, Anda dapat melindunginya dengan kata `privat`.

```celang
bentuk RekeningBank {
    pemilik: Teks, // Bisa dilihat siapa saja

    privat {
        saldo: Angka = 0, // Rahasia, hanya bisa diubah lewat tindakan resmi
        pin: Angka = 1234,

        validasiPin(inputPin: Angka) {
            kembalikan inputPin == @pin
        }
    }

    setor(jumlah: Angka) {
        @saldo += jumlah
        kembalikan @saldo
    }

    tarik(jumlah: Angka, inputPin: Angka) {
        jika !@validasiPin(inputPin) {
            lempar "Nomor PIN salah!"
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

// ❌ Jika seseorang mencoba mengubah saldo secara paksa dari luar:
// akun.saldo = 9999999 
// Celang akan menolak dengan ramah: 'saldo' bersifat privat!
```

---

## 🧩 Menggabungkan Bentuk Lain (Penyisipan Bentuk)

Anda tidak perlu mengetik ulang data yang sama jika sudah ada bentuk lain yang memilikinya. Anda cukup menyisipkan nama bentuk lain tersebut ke dalam bentuk baru:

```celang
bentuk Identitas {
    nama: Teks,
    umur: Angka,

    sapa() {
        kembalikan "Halo, nama saya " + @nama
    }
}

bentuk Karyawan {
    Identitas, // Langsung menyertakan nama, umur, dan kemampuan sapa
    gaji: Angka = 5000000,
    divisi: Teks = "Teknologi"
}

buat staf = Karyawan("Dewi", 25, 8000000, "Desain")
cetak(staf.nama)   // "Dewi"
cetak(staf.sapa())  // "Halo, nama saya Dewi"
cetak(staf.divisi) // "Desain"
```

---

## 🐑 Membuat Duplikat Mandiri (`salin`)

Jika Anda ingin membuat salinan persis dari suatu data tanpa khawatir perubahan pada salinan akan ikut merusak data aslinya, gunakan perintah `salin`:

```celang
buat asli = RekeningBank("Ali", 100000)
buat tiruan = salin(asli)

tiruan.pemilik = "Budi"
cetak(asli.pemilik)   // "Ali" (tetap aman, tidak berubah)
cetak(tiruan.pemilik) // "Budi"
```
