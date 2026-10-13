# 📦 Generik (Type Parameters & Constraints)

Generik di Celang memungkinkan Anda menulis **bentuk (struct)** dan **fungsi** yang fleksibel untuk berbagai tipe data, namun tetap mempertahankan **keamanan tipe data (type safety)** secara ketat.

---

## 💡 Mengapa Butuh Generik?

### Analogi Nyata: Kotak Kado Serbaguna
Bayangkan Anda memiliki sebuah kotak kado:
* Tanpa generik, Anda harus membuat `KotakAngka`, `KotakTeks`, dan `KotakKebenaran` secara terpisah walaupun cara kerja kotaknya persis sama.
* Dengan **Generik**, Anda cukup membuat satu cetak biru `Kotak<T>`. Tanda `<T>` adalah penampung tipe (*type parameter*) yang akan ditentukan saat kotak tersebut digunakan.

---

## 1. Generik pada Bentuk (`bentuk`)

### Bentuk Generik Dasar (`<T>`)
```celang
bentuk Kotak<T> {
    isi: T
}

// Menggunakan tipe konkret:
buat kotakAngka = Kotak<Angka>(42)
buat kotakTeks = Kotak<Teks>("Surat Rahasia")

cetak(kotakAngka.isi) // 42
cetak(kotakTeks.isi)  // "Surat Rahasia"
```

### Multiple Type Parameters (`<K, V>`)
Anda dapat menggunakan lebih dari satu parameter tipe, seperti pasangan Kunci-Nilai (*Key-Value*):

```celang
bentuk Pasangan<K, V> {
    kunci: K,
    nilai: V
}

buat skorPemain = Pasangan<Teks, Angka>("Budi", 950)
buat statusAktif = Pasangan<Angka, Kebenaran>(101, benar)

cetak(skorPemain.kunci, skorPemain.nilai)  // "Budi" 950
```

---

## 2. Batasan Tipe (Type Constraints)

Terkadang Anda ingin membuat tipe generik, tetapi **membatasi** tipe data apa saja yang diperbolehkan masuk. Celang menyediakan sintaks yang sangat manusiawi:

### Opsi 1: Menggunakan Titik Dua (`Nilai: Angka`)
```celang
bentuk Dompet<Nilai: Angka> {
    saldo: Nilai
}

buat dompetA = Dompet<Angka>(50000) // ✅ Sah! Angka memenuhi batasan Angka
```

### Opsi 2: Menggunakan Kata `harus` (`Nilai harus Angka`)
Bagi yang menyukai keterbacaan bahasa alami khas Celang:

```celang
bentuk Tabungan<Nilai harus Angka> {
    jumlah: Nilai
}

buat tabunganB = Tabungan<Angka>(100000) // ✅ Sah!
```

### Mengizinkan Bebas (`Sembarang`)
Jika ingin menegaskan bahwa tipe ini bebas menerima apa pun tanpa batasan:

```celang
bentuk Loker<Barang: Sembarang> {
    barang: Barang
}

buat loker1 = Loker<Teks>("Helm")
buat loker2 = Loker<Angka>(12345)
```

### Apa yang Terjadi Jika Melanggar Batasan?
Jika Anda mencoba memasukkan tipe yang tidak sesuai batasan, Celang akan memberikan pesan galat ramah pemula:

```celang
// ❌ Error!
buat dompetSalah = Dompet<Teks>("Uang Kertas")
// Kesalahan: Argumen tipe 'Teks' tidak memenuhi batasan 'Angka' pada bentuk 'Dompet'
// 💡 Saran: Gunakan tipe yang memenuhi batasan Angka.
```

---

## 3. Generik pada Fungsi

Fungsi juga dapat menggunakan parameter tipe generik untuk menciptakan algoritma serbaguna:

```celang
// Fungsi penukar dua elemen
fungsi tukar<T>(a: T, b: T) {
    kembalikan [b, a]
}

buat angkaDitukar = tukar<Angka>(10, 20)
cetak(angkaDitukar) // [20, 10]

buat teksDitukar = tukar<Teks>("Halo", "Dunia")
cetak(teksDitukar)  // ["Dunia", "Halo"]
```

Dengan batasan tipe pada fungsi:
```celang
fungsi kaliDua<Nilai harus Angka>(x: Nilai) {
    kembalikan x * 2
}

cetak(kaliDua<Angka>(25)) // 50
```

---

## 4. Inferensi Tipe Otomatis

Di Celang, Anda bahkan tidak wajib selalu menulis tanda kurung sudut `<...>` saat memanggil fungsi atau membuat objek jika tipenya sudah jelas dari argumen yang diberikan:

```celang
bentuk Kotak<T> {
    isi: T
}

// Celang otomatis mendeteksi bahwa T adalah Angka!
buat k = Kotak(100)
cetak(k.isi) // 100
```

---

## 🎯 Rangkuman Sintaks

| Bentuk Sintaks | Arti / Kegunaan |
| :--- | :--- |
| `bentuk Box<T> { ... }` | Bentuk generik dengan 1 parameter tipe |
| `bentuk Peta<K, V> { ... }` | Bentuk generik dengan 2 parameter tipe |
| `bentuk Box<Nilai: Angka> { ... }` | Parameter tipe wajib memenuhi tipe `Angka` |
| `bentuk Box<Nilai harus Angka> { ... }` | Gaya bahasa alami untuk pembatasan tipe |
| `bentuk Box<T: Sembarang> { ... }` | Parameter tipe bebas tanpa batasan |
| `fungsi proses<T>(x: T) { ... }` | Fungsi generik serbaguna |
| `Box<Angka>(42)` | Pemanggilan dengan argumen tipe eksplisit |
