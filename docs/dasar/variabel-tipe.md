# Tipe Data & Variabel

Setiap informasi yang Anda olah di dalam Celang memiliki jenis atau **tipe data**. Halaman ini menyajikan seluruh tipe data bawaan Celang terlebih dahulu agar Anda dapat langsung mengenali dan menggunakannya, diikuti cara menyimpan data ke dalam variabel dan konstanta.

---

## 1. Ikhtisar Tipe Data

Celang menyediakan tipe data lengkap untuk kebutuhan komputasi sehari-hari hingga pengolahan data biner tingkat lanjut:

| Tipe Data | Nama Tipe | Contoh Nilai | Keterangan |
| :--- | :--- | :--- | :--- |
| **Bilangan Bulat** | `Angka` | `42`, `-10`, `1_000_000` | Bilangan bulat utuh tanpa koma (garis bawah `_` sebagai pemisah ribuan) |
| **Bilangan Desimal** | `Pecahan` / `Desimal` | `3.14`, `-0.05` | Bilangan desimal dengan tanda titik |
| **Tulisan / Teks** | `Teks` | `"Halo"`, `'Dunia'`, `"${nama}"` | Tulisan huruf, kata, kalimat, atau teks Unicode |
| **Nilai Kebenaran** | `Kebenaran` | `benar`, `salah` | Nilai logika kebenaran (*boolean*: ya atau tidak) |
| **Koleksi Bita Biner** | `Bita` / `Bytes` | `bita("Halo")`, `bita(8)` | Deretan byte mentah (0–255) untuk berkas, gambar, dan jaringan |
| **Daftar / Larik** | `Array` | `["Apel", "Mangga", "Jeruk"]` | Kumpulan data berurutan dalam satu wadah |
| **Kamus Data** | `Map` | `{"nama": "Budi", "umur": 24}` | Pasangan label (kunci) dan nilai |
| **Nilai Kosong** | `kosong` | `kosong` | Menandakan nilai nihil atau ketiadaan data |

---

## 2. Mengenal Masing-Masing Tipe Data

### A. Bilangan Bulat (`Angka`) & Desimal (`Pecahan`)

Celang membedakan angka bulat dan angka berkoma secara otomatis:

```celang
buat jumlah = 42
buat jarak = -10
buat populasi = 275_000_000 // Tanda '_' memudahkan pembacaan angka besar

buat pi = 3.14159
buat suhu = -0.5
```

Celang juga mendukung penulisan bilangan heksadesimal (`0xFF`), biner (`0b1010`), dan oktal (`0o77`).

### B. Tulisan Teks (`Teks`)

Teks dapat diapit tanda kutip ganda (`"..."`) atau kutip tunggal (`'...'`). Anda juga dapat menyisipkan variabel langsung ke dalam kalimat menggunakan *f-string* (awalan `f`) dengan notasi `${...}` atau `{...}`:

```celang
buat salam = "Halo Dunia!"
buat kota = 'Jakarta'

buat nama = "Budi"
buat pesan = f"Selamat datang, ${nama}! Semoga harimu menyenangkan."
cetak(pesan) // "Selamat datang, Budi! Semoga harimu menyenangkan."
```

### C. Nilai Kebenaran (`Kebenaran`)

Tipe `Kebenaran` hanya memiliki dua nilai: `benar` (*true*) atau `salah` (*false*). Tipe ini dihasilkan oleh operasi perbandingan atau kondisi logika:

```celang
buat aktif = benar
buat lulus = salah

buat nilai = 85
buat memenuhiSyarat = nilai >= 75 // bernilai benar
```

### D. Koleksi Bita Biner (`Bita` / `Bytes`)

`Bita` menyimpan deretan bita mentah (angka `0` hingga `255`). Berbeda dengan `Teks` yang memproses karakter Unicode, `Bita` dirancang khusus untuk memanipulasi berkas biner (gambar, berkas ZIP, audio), protokol jaringan, dan enkripsi:

```celang
// 1. Membuat Bita:
buat b1 = bita("Halo")               // dari teks: bita([72, 97, 108, 111])
buat b2 = bita(4)                    // alokasi buffer 4 byte nol: bita([0, 0, 0, 0])
buat b3 = bita([0x43, 0x65, 0x6C])   // dari daftar kode byte ASCII: bita([67, 101, 108])
buat b4 = bita("SmFsbw==", "base64") // dari teks Base64
buat b5 = bita("4a616c6f", "hex")    // dari teks heksadesimal

// 2. Baca dan ubah byte langsung berdasarkan indeks:
cetak(b1[0])       // 72 (kode ASCII huruf 'H')
b1[0] = 74         // Ubah 'H' menjadi 'J' (ASCII 74)
cetak(b1.ke_teks()) // "Jalo"

// 3. Konversi format biner:
cetak(b1.ke_heks())   // "4a616c6f"
cetak(b1.ke_base64()) // "SmFsbw=="
cetak(b1.ke_daftar()) // [74, 97, 108, 111] (Array Angka)

// 4. Dekode dari Base64 dan Heksadesimal:
buat bDariB64 = bita.dari_base64("SmFsbw==")
buat bDariHex = bita.dari_heks("4a616c6f")
cetak(bDariB64.ke_teks()) // "Jalo"

// 5. Mutasi, Pemotongan, & Penggabungan:
b1.tambah(33)                 // Menambah byte ASCII '!' (nilai 33)
b1.tambah(" Dunia")           // Menambah teks langsung ke bita
b1.tambah(bita([46, 46, 46])) // Menambah bita lain
cetak(b1.ke_teks())           // "Jalo! Dunia..."

buat subBita = b1.potong(0, 4) // Potong 4 byte pertama
cetak(subBita.ke_teks())      // "Jalo"

buat salinan = b1.salin()     // Duplikasi mandiri (deep clone)
cetak("Ukuran byte:", b1.panjang) // atau b1.length
```

#### Ringkasan Metode & Properti `Bita`

| Metode / Properti (Indonesia) | Alias Universal / Inggris | Keterangan |
| :--- | :--- | :--- |
| `bita(ukuran)` | `bytes(size)` | Mengalokasikan buffer nol sebesar $N$ byte (maks 100 MB) |
| `bita(teks, format)` | `bytes(str, enc)` | Membuat bita dari teks dengan encoding `"base64"` atau `"hex"` |
| `bita.dari_base64(teks)` | `bytes.from_base64(str)` | Mengurai teks Base64 menjadi objek `Bita` |
| `bita.dari_heks(teks)` | `bytes.from_hex(str)` | Mengurai teks heksadesimal menjadi objek `Bita` |
| `.ke_teks()` | `.to_string()` / `.to_text()` | Mengonversi deretan bita menjadi teks (*string*) UTF-8 |
| `.ke_heks()` / `.heks()` | `.to_hex()` / `.hex()` | Mengonversi deretan bita menjadi teks heksadesimal |
| `.ke_base64()` | `.to_base64()` | Mengonversi deretan bita menjadi string Base64 standar |
| `.ke_daftar()` | `.to_array()` | Mengonversi bita menjadi daftar array angka `[0..255]` |
| `.tambah(...)` | `.append(...)` | Menambahkan byte angka, teks, atau bita lain ke akhir buffer |
| `.potong(awal, [akhir])` | `.slice(start, [end])` | Mengambil potongan bita (mendukung indeks negatif) |
| `.salin()` | `.clone()` | Mengkloning bita secara independen (*deep copy*) |
| `.panjang` | `.length` / `.len` | Properti pembaca jumlah total bita |

> 🌍 **Mode Universal:** Saat menggunakan `"use universal"`, panggil fungsi dengan nama bahasa Inggris (`bytes`, `to_string`, `to_hex`, `to_base64`, `to_array`, `append`, `slice`, `clone`, `length`). Celang akan memberikan saran empati jika nama lokal terpanggil.

### E. Kumpulan Data: Daftar (`Array`) & Kamus (`Map`)

#### 1. Daftar / Array (`[...]`)
Daftar digunakan untuk menyimpan banyak nilai berurutan. Indeks nomor urut selalu dimulai dari `0`:

```celang
buat buah = ["Apel", "Mangga", "Jeruk", "Pisang"]

cetak(buah[0])       // "Apel" (elemen pertama)
cetak(panjang(buah)) // 4 (jumlah elemen)

// Mengambil potongan tertentu (rentang):
buat pilihan = buah[1..3] // ["Mangga", "Jeruk"]

// Manipulasi daftar bawaan:
buah.push("Melon")        // Menambah di belakang
buat akhir = buah.pop()   // Mengambil elemen terakhir ("Melon")
buah.reverse()            // Membalikkan urutan
```

#### 2. Kamus Data / Map (`{...}`)
Kamus digunakan ketika setiap data memiliki nama label tersendiri (seperti data formulir atau profil):

```celang
buat pengguna = {
    "nama": "Siti Rahma",
    "divisi": "Teknologi",
    "gaji": 15000000,
    "aktif": benar
}

// Membaca data dengan titik atau kurung siku:
cetak(pengguna.nama)      // "Siti Rahma"
cetak(pengguna["divisi"]) // "Teknologi"

// Jika mencari kunci yang belum ada, hasilnya adalah 'kosong':
cetak(pengguna.alamat)    // kosong

// Manipulasi kamus:
buat daftarKunci = pengguna.keys() // ["nama", "divisi", "gaji", "aktif"]
pengguna.set("lokasi", "Jakarta")  // Menambah atau mengubah nilai
```

> 💡 **Panduan Lengkap Koleksi:** Pelajari seluruh metode manipulasi array dan map pada bab [Manipulasi Koleksi (Array & Map)](/fitur-khas/koleksi-array-map).

### F. Nilai Kosong (`kosong`)

`kosong` (*null* / *nil*) menandakan ketiadaan nilai atau data yang belum diisi. Operasi aman (*null-safety*) didukung penuh di Celang agar aplikasi tidak mogok saat membaca data kosong.

---

## 3. Deklarasi Variabel & Nilai Tetap

Setelah memahami jenis-jenis data di atas, Anda dapat menyimpannya ke dalam wadah menggunakan kata kunci berikut:

### A. Variabel yang Bisa Diubah (`buat`)

Gunakan `buat` untuk membuat wadah variabel yang nilainya dapat diperbarui di kemudian hari:

```celang
buat skor = 10
cetak(skor) // 10

// Memperbarui nilai:
skor = 25
cetak(skor) // 25
```

> **Nilai Awal Otomatis Kosong:** Jika Anda membuat variabel tanpa langsung mengisinya, Celang secara otomatis mengisinya dengan nilai `kosong`:
> ```celang
> buat alamat
> cetak(alamat) // kosong
> ```

### B. Nilai Tetap yang Tidak Boleh Berubah (`konstanta`)

Jika ada nilai patokan yang harus selalu tetap dan tidak boleh diubah oleh bagian program lain, gunakan kata kunci `konstanta`:

```celang
konstanta PAJAK = 0.11
konstanta NAMA_APLIKASI = "TokoKita"

// Jika mencoba mengubah nilai konstanta:
// PAJAK = 0.15 
// ❌ Celang akan menolak dengan ramah: nilai konstanta tidak boleh diubah!
```

> 💡 **Konsep Immutability: Binding vs Nilai Objek**
>
> Perlu dibedakan antara **Binding Immutability** dan **Value Immutability**:
> * `konstanta` / `tetap` hanya mengunci **label variabel** (*binding*) agar tidak bisa di-assign ulang ke objek lain.
> * Jika konstanta menampung sebuah Array, pemanggilan metode mutasi internal seperti `daftar.tambah(4)` tetap dapat memodifikasi isi memori array tersebut.
> * Untuk menciptakan **Immutability Nilai (Bebas Efek Samping)** tanpa `konstanta`, gunakan metode transformasi fungsional seperti `.saring()`, `.petakan()`, atau salin array dengan `.salin()` sebelum memodifikasinya.

### C. Anotasi Tipe Statis Opsional (Type Hinting)

Celang mendukung penentuan tipe secara eksplisit untuk menjaga kejelasan kode Anda:

```celang
buat nama: Teks = "Budi Santoso"
buat umur: Angka = 25
buat saldo: Pecahan = 175000.50
buat aktif: Kebenaran = benar
buat buffer: Bita = bita(16)
```

### D. Batas Lingkup Catatan / Ruang Lingkup (`{ ... }`)

Variabel yang dibuat di dalam kurung kurawal `{ ... }` hanya berlaku di dalam blok tersebut:

```celang
buat judul = "Buku Utama"

{
    // Catatan di dalam kurung kurawal hanya dikenal di sini:
    buat catatanKecil = "Catatan Bab 1"
    cetak(judul)        // "Buku Utama" (dapat membaca variabel luar)
    cetak(catatanKecil) // "Catatan Bab 1"
}

// Di luar kurung kurawal, 'catatanKecil' sudah tidak ada:
// cetak(catatanKecil) ❌ (Variabel tidak ditemukan)
cetak(judul) // "Buku Utama" (tetap aman)
```

---

## 4. Konversi Tipe Data (Casting) & Pemeriksaan Tipe

### A. Memeriksa Jenis Data (`tipe`)

Gunakan fungsi bawaan `tipe(nilai)` untuk mengetahui jenis data dari suatu nilai:

```celang
cetak(tipe(42))            // "INTEGER"
cetak(tipe(3.14))          // "FLOAT"
cetak(tipe("Halo"))        // "STRING"
cetak(tipe(benar))         // "BOOLEAN"
cetak(tipe(bita("data")))  // "BITA"
cetak(tipe([1, 2, 3]))     // "ARRAY"
cetak(tipe({"a": 1}))      // "MAP"
cetak(tipe(kosong))        // "NULL"
```

### B. Fungsi Pengubah Jenis Data (Casting)

Celang menyediakan fungsi konversi eksplisit yang ramah dan konsisten:

```celang
// 1. Mengubah teks ke Angka atau Pecahan:
buat angka1 = Angka("25000")       // 25000
buat desimal1 = Pecahan("3.14")    // 3.14

// 2. Mengubah nilai ke Teks:
buat teksAngka = Teks(1000)        // "1000"
buat teksBool = Teks(benar)        // "benar"

// 3. Konversi nilai ke Kebenaran (Boolean):
buat b1 = Kebenaran(1)             // benar
buat b2 = Kebenaran(0)             // salah
buat b3 = Kebenaran("")            // salah (teks kosong adalah falsy)
buat b4 = Kebenaran("Celang")      // benar
buat b5 = Kebenaran([])            // salah (daftar kosong adalah falsy)

// 4. Konversi antara Teks dan Bita:
buat dataBiner = Bita("Celang")    // Bita([67, 101, 108, 97, 110, 103])
buat teksBalik = Teks(dataBiner)   // "Celang" (langsung diterjemahkan ke UTF-8)
```

Celang juga secara cerdas mengenali format mata uang saat operasi matematika berlangsung (misalnya `"Rp 150000" - 25000` akan otomatis menghasilkan angka `125000`).
