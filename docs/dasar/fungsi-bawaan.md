# Fungsi Bawaan (Built-in Functions)

Celang menyediakan sejumlah fungsi bawaan global yang selalu tersedia secara langsung di seluruh cakupan (*scope*) program tanpa perlu mengimpor modul pustaka standar terlebih dahulu.

Fungsi-fungsi ini mendukung penamaan idiomatik Bahasa Indonesia dan alias internasional yang ramah bagi programmer dari berbagai latar belakang.

---

## 📋 Ringkasan Daftar Fungsi

| Kategori | Fungsi (Indonesia) | Alias Internasional | Keterangan |
| :--- | :--- | :--- | :--- |
| **I/O & Terminal** | `cetak(...nilai)` | `print` | Mencetak nilai ke terminal dengan spasi & baris baru otomatis |
| | `masukan(prompt?)` | `input` | Membaca input baris teks dari terminal pengguna |
| **Inspeksi Data** | `panjang(koleksi)` | `len` | Menghitung jumlah elemen teks, array, atau map |
| | `tipe(nilai)` | `type` | Mengembalikan nama tipe data objek saat runtime |
| **Konversi Tipe** | `Angka(nilai)` | `angka` | Mengonversi nilai teks/logika menjadi bilangan bulat/pecahan |
| | `Teks(nilai)` | `teks` | Mengonversi nilai apa pun menjadi string teks |
| | `Logika(nilai)` | `logika` | Mengevaluasi nilai kebenaran boolean (`benar` / `salah`) |
| **Koleksi & Array** | `rentang(awal?, akhir)` | `range` | Menghasilkan daftar array bilangan terurut |
| | `saring(array, fn)` | `filter` | Menyaring elemen array menggunakan predikat fungsi |
| | `petakan(array, fn)` | `map` | Mentransformasikan setiap elemen array dengan fungsi callback |
| **Manipulasi Teks** | `potong_spasi(teks)` | `trim` | Menghapus spasi putih di awal dan akhir teks |
| | `kapitalisasi(teks)` | `capitalize`, `uppercase` | Mengubah seluruh karakter teks menjadi huruf besar |
| **Manajemen Memori**| `salin(obj)` | `clone` | Melakukan duplikasi mendalam (*deep clone*) tanpa berbagi memori |

---

## 1. I/O & Interaksi Terminal

### `cetak(...nilai)` / `print`
Mencetak satu atau beberapa nilai ke output terminal. Nilai yang dipisahkan tanda koma akan dicetak dengan spasi pemisah, dan diakhiri dengan baris baru (*newline*).

```celang
cetak("Halo Dunia!")
cetak("Nilai:", 42, "Status:", benar)
// Output: Nilai: 42 Status: benar
```

### `masukan(prompt?)` / `input`
Membaca satu baris masukan teks dari pengguna melalui terminal. Parameter prompt bersifat opsional.

```celang
buat nama = masukan("Masukkan nama Anda: ")
cetak("Selamat datang,", nama)
```

---

## 2. Inspeksi Data & Panjang

### `panjang(koleksi)` / `len`
Mengembalikan jumlah elemen dari `Teks` (panjang karakter), `Array` (jumlah item), atau `Map` (jumlah entri pasangan kunci-nilai).

```celang
cetak(panjang("Celang"))       // 6
cetak(panjang([10, 20, 30]))   // 3
cetak(panjang({"a": 1, "b": 2})) // 2
```

### `tipe(nilai)` / `type`
Mengecek nama tipe data objek saat runtime. Mengembalikan objek `Teks`.

```celang
cetak(tipe(100))        // "Angka"
cetak(tipe(3.14))       // "Pecahan"
cetak(tipe("Halo"))     // "Teks"
cetak(tipe(benar))      // "Logika"
cetak(tipe([1, 2]))     // "Array"
cetak(tipe({"x": 10}))  // "Map"
```

---

## 3. Konversi Tipe Data (*Type Casting*)

Fungsi-fungsi ini bekerja sebagai konstruktor konversi eksplisit:

### `Angka(nilai)` / `angka`
Mengonversi teks berformat angka atau boolean menjadi `Angka`.
```celang
buat x = Angka("125000") // 125000
buat y = Angka(3.95)     // 3 (pemotongan bulat)
buat z = Angka(benar)    // 1
```

### `Teks(nilai)` / `teks`
Mengonversi tipe data apa pun menjadi `Teks`.
```celang
buat s1 = Teks(42)       // "42"
buat s2 = Teks(benar)    // "benar"
buat s3 = Teks([1, 2])   // "[1, 2]"
```

### `Logika(nilai)` / `logika`
Mengevaluasi kebenaran (*truthiness*) suatu nilai. Nilai `0`, string kosong `""`, array kosong `[]`, atau `kosong` bernilai `salah`, sisanya bernilai `benar`.
```celang
cetak(Logika(1))         // benar
cetak(Logika(0))         // salah
cetak(Logika("Celang"))  // benar
cetak(Logika(""))        // salah
```

---

## 4. Manipulasi Koleksi & Array

### `rentang(awal?, akhir)` / `range`
Menghasilkan array angka terurut. Jika hanya diberikan 1 argumen `rentang(akhir)`, rentang dimulai dari `0` sampai `akhir`. Jika diberikan 2 argumen `rentang(awal, akhir)`, rentang mencakup dari `awal` hingga `akhir` (bisa naik atau turun).

```celang
buat r1 = rentang(4)       // [0, 1, 2, 3, 4]
buat r2 = rentang(1, 5)    // [1, 2, 3, 4, 5]
buat r3 = rentang(3, 1)    // [3, 2, 1]
```

### `saring(array, fungsi)` / `filter`
Menyaring elemen array berdasarkan fungsi predikat yang mengembalikan `benar` atau `salah`.

```celang
buat angka = [1, 2, 3, 4, 5, 6]
buat genap = saring(angka, fungsi(n) {
    kembalikan n % 2 == 0
})
cetak(genap) // [2, 4, 6]
```

### `petakan(array, fungsi)` / `map`
Mentransformasikan setiap elemen array ke dalam bentuk baru melalui fungsi callback.

```celang
buat angka = [1, 2, 3]
buat kuadrat = petakan(angka, fungsi(n) {
    kembalikan n * n
})
cetak(kuadrat) // [1, 4, 9]
```

---

## 5. Manipulasi Teks

### `potong_spasi(teks)` / `trim`
Menghapus karakter spasi putih (*whitespace*) di awal dan di akhir teks.

```celang
buat pesan = "   Halo Celang!   "
cetak(potong_spasi(pesan)) // "Halo Celang!"
```

### `kapitalisasi(teks)` / `capitalize` / `uppercase`
Mengubah semua huruf dalam teks menjadi huruf besar (*uppercase*).

```celang
buat judul = "bahasa pemrograman indonesia"
cetak(kapitalisasi(judul)) // "BAHASA PEMROGRAMAN INDONESIA"
```

---

## 6. Manajemen Memori & Kloning

### `salin(obj)` / `clone`
Melakukan *deep memory clone* pada objek, array bersarang, map, atau bentuk (*struct*). Objek hasil kloning berdiri sendiri secara independen, sehingga mutasi pada hasil salinan tidak akan mempengaruhi objek asli.

```celang
bentuk Pengguna {
    nama: Teks,
    poin: Angka
}

buat p1 = Pengguna("Ahmad", 100)
buat p2 = salin(p1) // kloning independen

p2.nama = "Budi"
p2.poin = 200

cetak(p1.nama, p1.poin) // "Ahmad", 100 (tidak berubah)
cetak(p2.nama, p2.poin) // "Budi", 200
```

---

## 🌊 Integrasi Bersih dengan Pipeline Operator (`|>`)

Semua fungsi bawaan Celang dirancang agar selaras dengan **Pipeline Operator** (`|>`), memungkinkan penulisan alur data bergaya fungsional yang mudah dibaca:

```celang
buat dataMentah = "   celang itu keren   "

dataMentah
    |> potong_spasi
    |> kapitalisasi
    |> cetak
// Output: CELANG ITU KEREN

// Mengolah koleksi angka dengan pipeline:
rentang(1, 10)
    |> saring(fungsi(x) { kembalikan x % 2 != 0 })
    |> petakan(fungsi(x) { kembalikan x * 10 })
    |> cetak
// Output: [10, 30, 50, 70, 90]
```
