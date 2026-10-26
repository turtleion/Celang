# Manipulasi Koleksi (Array & Map)

Celang menyediakan metode manipulasi bawaan (*built-in collection methods*) yang lengkap dan ekspresif untuk **Array** (Daftar Berurutan) dan **Map** (Kamus Pasangan Kunci-Nilai).

Metode-metode ini mendukung **pemanggilan berantai (*fluent chaining*)**, **penamaan dwibahasa (Bahasa Indonesia & Universal/Inggris)**, serta kepatuhan ketat pada **Mode Universal (`"use universal"`)**.

---

## 1. Manipulasi Array (Daftar)

Array di Celang merepresentasikan kumpulan data berurutan. Anda dapat memanipulasi elemen array secara langsung menggunakan metode bawaan:

### Ringkasan Metode Array

| Metode (Universal) | Metode (Indonesia) | Argumen | Keterangan | Nilai Kembali |
| :--- | :--- | :--- | :--- | :--- |
| `.push(...)` | `.tambah(...)` | `elemen, ...` | Menambahkan satu atau lebih elemen ke ujung belakang array | Array (memungkinkan *chaining*) |
| `.pop()` | `.hapus_terakhir()` | *(tidak ada)* | Mengambil dan menghapus elemen terakhir dari array | Elemen yang terhapus / `kosong` |
| `.shift()` | `.hapus_awal()` | *(tidak ada)* | Mengambil dan menghapus elemen pertama dari array | Elemen yang terhapus / `kosong` |
| `.unshift(...)` | `.sisipkan_depan(...)` | `elemen, ...` | Menyisipkan satu atau lebih elemen di awal array | Array (memungkinkan *chaining*) |
| `.insert(indeks, el)`| `.sisipkan(indeks, el)` | `indeks, elemen` | Menyisipkan elemen pada indeks tertentu (mendukung indeks negatif) | Array (memungkinkan *chaining*) |
| `.remove_at(indeks)` | `.hapus_indeks(indeks)` | `indeks` | Menghapus elemen pada indeks tertentu | Elemen yang terhapus |
| `.reverse()` | `.balik()` | *(tidak ada)* | Membalikkan urutan elemen array secara langsung (*in-place*) | Array (memungkinkan *chaining*) |
| `.sort([pembanding])`| `.urut([pembanding])` | `[fungsi(a, b)]` | Mengurutkan elemen (alami atau menggunakan fungsi perbandingan kustom) | Array (memungkinkan *chaining*) |
| `.join([pemisah])` | `.gabung([pemisah])` | `[pemisah]` | Menggabungkan seluruh elemen menjadi satu teks (bawaan: `", "`) | `Teks` (*String*) |
| `.contains(el)` / `.includes(el)` | `.berisi(el)` | `elemen` | Memeriksa apakah suatu elemen ada di dalam array | `Kebenaran` (`benar`/`salah`) |
| `.index_of(el)` | `.indeks_dari(el)` | `elemen` | Mencari posisi indeks pertama dari elemen (mengembalikan `-1` jika tidak ada) | `Angka` (*Integer*) |
| `.map(fungsi)` | `.petakan(fungsi)` | `fungsi(el, [i])` | Mengubah setiap elemen dengan fungsi (tidak mengubah array asal) | Salinan Array baru (Imutabel) |
| `.filter(fungsi)` | `.saring(fungsi)` | `fungsi(el, [i])` | Menyaring elemen yang memenuhi kondisi predikat (tidak mengubah array asal) | Salinan Array baru (Imutabel) |
| `.reduce(fungsi, [awal])` | `.reduksi(...)` / `.lipat(...)` | `fungsi(acc, el, [i]), [awal]` | Mereduksi seluruh elemen menjadi satu nilai akumulasi | Nilai akumulasi |
| `.find(fungsi)` | `.temukan(fungsi)` | `fungsi(el, [i])` | Mencari elemen pertama yang memenuhi kriteria (*null-safe*) | Elemen pertama / `kosong` |
| `.some(fungsi)` | `.ada_yang(fungsi)` | `fungsi(el, [i])` | Memeriksa apakah minimal satu elemen memenuhi kriteria | `Kebenaran` (`benar`/`salah`) |
| `.every(fungsi)` | `.semuanya(fungsi)` | `fungsi(el, [i])` | Memeriksa apakah seluruh elemen memenuhi kriteria | `Kebenaran` (`benar`/`salah`) |
| `.slice([mulai, [akhir]])` | `.potong([mulai, [akhir]])`| `[mulai, [akhir]]` | Mengambil salinan potongan array tanpa mengubah array asli | Salinan Array baru |
| `.concat(...)` | `.gabungkan(...)` | `array_atau_elemen, ...`| Menggabungkan array dengan elemen atau array lain | Array baru |
| `.clone()` | `.salin()` | *(tidak ada)* | Menduplikasi array secara mandiri (*shallow copy*) | Salinan Array baru |
| `.clear()` | `.bersihkan()` | *(tidak ada)* | Mengosongkan seluruh isi array | Array (kosong `[]`) |
| `.length` / `.len` | `.panjang` | *(properti)* | Mengambil jumlah total elemen | `Angka` (*Integer*) |
| `.size()` | `.ukuran()` | *(tidak ada)* | Mengambil jumlah total elemen sebagai pemanggilan fungsi | `Angka` (*Integer*) |

---

### Contoh Penggunaan Array

#### Menambah, Mengambil, dan Menyisipkan Elemen
```celang
buat angka = [10, 20, 30]

// Menambah elemen ke ujung belakang (bisa berantai)
angka.push(40, 50)
cetak(angka) // [10, 20, 30, 40, 50]

// Menghapus dan mengambil elemen terakhir & pertama
buat terakhir = angka.hapus_terakhir() // 50 (alias: .pop() / .ambil_terakhir())
buat pertama  = angka.hapus_awal()     // 10 (alias: .shift() / .geser())
cetak(terakhir, pertama)              // 50 10
cetak(angka)                          // [20, 30, 40]

// Sisipkan di depan atau pada indeks spesifik
angka.sisipkan_depan(5)               // sisip di depan: [5, 20, 30, 40]
angka.sisipkan(2, 25)                 // sisipkan pada indeks 2: [5, 20, 25, 30, 40]
buat dihapus = angka.hapus_indeks(1)  // hapus elemen pada indeks 1 (nilai 20)
cetak(angka)                          // [5, 25, 30, 40]
```

#### Pengurutan (*Sorting*) & Pembalikan (*Reverse*)
```celang
buat kata = ["mangga", "apel", "jeruk", "durian"]

// Urutkan secara alfabet
kata.sort()
cetak(kata) // [apel, durian, jeruk, mangga]

// Balikkan urutannya
kata.reverse()
cetak(kata) // [mangga, jeruk, durian, apel]

// Urutkan angka dengan pembanding kustom (misal: urutan menurun / descending)
buat skor = [85, 92, 78, 99, 88]
skor.sort((a, b) => b - a)
cetak(skor) // [99, 92, 88, 85, 78]
```

#### Pencarian, Pemotongan, dan Penggabungan
```celang
buat hewan = ["kucing", "anjing", "kelinci", "panda"]

cetak(hewan.contains("kucing"))   // benar
cetak(hewan.index_of("kelinci"))  // 2
cetak(hewan.index_of("singa"))    // -1

// Menggabungkan menjadi teks string
buat daftarTeks = hewan.join(" | ")
cetak(daftarTeks) // "kucing | anjing | kelinci | panda"

// Potong sebagian array (slice)
buat favorit = hewan.slice(1, 3)
cetak(favorit) // [anjing, kelinci]
```

#### Transformasi Fungsional & Imutabilitas (*Functional Pipelines*)

Celang mendukung metode manipulasi fungsional modern yang bersifat **imutabel (tidak memutasi array asli)**. Setiap pemanggilan `.saring()` atau `.petakan()` mengembalikan objek `Array` baru di memori:

```celang
buat angka = [1, 2, 3, 4, 5, 6]

// 1. .petakan() / .map() — Mengalikan setiap angka dengan 10
buat sepuluhKali = angka.petakan(x => x * 10)
cetak(sepuluhKali) // [10, 20, 30, 40, 50, 60]
cetak(angka)       // [1, 2, 3, 4, 5, 6] (Array asal tetap utuh!)

// 2. .saring() / .filter() — Mengambil hanya bilangan genap
buat genap = angka.saring(x => x % 2 == 0)
cetak(genap)       // [2, 4, 6]

// 3. .reduksi() / .lipat() / .reduce() — Menjumlahkan total
buat total = angka.reduksi((acc, x) => acc + x, 0)
cetak(total)       // 21

// 4. .temukan() / .find() — Menemukan elemen pertama (null-safe)
buat pertamaBesar = angka.temukan(x => x > 4) // 5
buat nihil = angka.temukan(x => x > 100)       // kosong

// 5. .ada_yang() / .some() & .semuanya() / .every()
cetak(angka.ada_yang(x => x > 5))  // benar
cetak(angka.semuanya(x => x > 0))  // benar

// 6. Rantai Proses (Fluent Pipeline Chaining)
buat ringkasan = angka
    .saring(x => x % 2 != 0)          // [1, 3, 5]
    .petakan(x => x * x)              // [1, 9, 25]
    .reduksi((acc, x) => acc + x, 0)  // 35
cetak("Total kuadrat ganjil:", ringkasan) // 35
```

> 💡 **Mengapa `.saring()` dan `.petakan()` Bersifat Imutabel? Apakah Memakai `const`?**
>
> 1. **Bukan karena `const` (`tetap`)**: Kata kunci `const` hanya mengunci penugasan (*variable binding*) agar variabel tidak bisa di-assign ulang, bukan mengunci data di dalam memori.
> 2. **Imutabilitas Struktural (Fungsi Murni)**: Di balik layar mesin Celang, `.saring()` dan `.petakan()` mengalokasikan irisan memori baru (`new Array`) dan mengembalikannya ke pemanggil. Array asli sama sekali tidak disentuh, sehingga bebas dari efek samping (*zero side-effects*).
> 3. **Pola Imutabel Tanpa `const`**: Untuk menjaga data tidak berubah tanpa bergantung pada `const`, biasakan memakai metode murni di atas, atau duplikasi data terlebih dahulu dengan `.salin()` sebelum memanggil metode mutasi seperti `.tambah()`.

---

## 2. Manipulasi Map (Kamus Pasangan Kunci-Nilai)

Map di Celang menyimpan pasangan kunci dan nilai terurut berdasarkan urutan penambahan (*insertion order*).

### Ringkasan Metode Map

| Metode (Universal) | Metode (Indonesia) | Argumen | Keterangan | Nilai Kembali |
| :--- | :--- | :--- | :--- | :--- |
| `.keys()` | `.kunci()` | *(tidak ada)* | Mengambil seluruh kunci map dalam urutan penambahan | `Array` berisi teks kunci |
| `.values()` | `.nilai()` | *(tidak ada)* | Mengambil seluruh nilai map dalam urutan penambahan | `Array` berisi nilai |
| `.entries()` | `.entri()` | *(tidak ada)* | Mengambil seluruh pasangan sebagai array `[[kunci, nilai], ...]` | `Array` 2-dimensi |
| `.has(kunci)` / `.contains(kunci)` | `.memiliki(kunci)` | `kunci` | Memeriksa apakah kunci tertentu ada di dalam map | `Kebenaran` (`benar`/`salah`) |
| `.get(kunci, [bawaan])` | `.ambil(kunci, [bawaan])` | `kunci, [bawaan]`| Mengambil nilai berdasarkan kunci, dengan nilai cadangan opsional jika kunci tidak ditemukan | Nilai objek / `bawaan` / `kosong` |
| `.set(kunci, nilai)` | `.atur(kunci, nilai)` | `kunci, nilai` | Menyimpan pasangan kunci-nilai ke dalam map | Map (memungkinkan *chaining*) |
| `.delete(kunci)` / `.remove(kunci)` | `.hapus(kunci)` | `kunci` | Menghapus kunci dari map | `Kebenaran` (`benar` jika dihapus, `salah` jika tidak ada) |
| `.merge(mapLain)` | `.gabung(mapLain)` | `mapLain` | Menggabungkan seluruh pasangan dari map lain ke dalam map ini | Map saat ini |
| `.clear()` | `.bersihkan()` | *(tidak ada)* | Mengosongkan seluruh isi map | Map saat ini (kosong `{}`) |
| `.length` / `.len` | `.panjang` | *(properti)* | Jumlah kunci yang tersimpan | `Angka` (*Integer*) |
| `.size()` | `.ukuran()` | *(tidak ada)* | Jumlah kunci sebagai pemanggilan fungsi | `Angka` (*Integer*) |

---

### Contoh Penggunaan Map

#### Membaca Kunci, Nilai, dan Entri
```celang
buat produk = {
    "kode": "PRD-01",
    "nama": "Kopi Arabika",
    "harga": 45000
}

buat semuaKunci = produk.keys()
buat semuaNilai = produk.values()
buat daftarEntri = produk.entries()

cetak("Kunci :", semuaKunci)   // [kode, nama, harga]
cetak("Nilai :", semuaNilai)   // [PRD-01, Kopi Arabika, 45000]
cetak("Entri :", daftarEntri)  // [[kode, PRD-01], [nama, Kopi Arabika], [harga, 45000]]
cetak("Ukuran:", produk.length) // 3
```

#### Pengambilan Aman dengan Cadangan (`.get()`) & Pemeriksaan (`.has()`)
```celang
buat config = {
    "port": 8080,
    "host": "localhost"
}

// Memeriksa keberadaan kunci
cetak(config.has("port"))    // benar
cetak(config.has("timeout")) // salah

// Mengambil nilai dengan nilai bawaan (fallback) jika kunci tidak ada
buat timeout = config.get("timeout", 5000)
cetak("Timeout:", timeout) // 5000 (menggunakan nilai bawaan)

buat port = config.get("port", 3000)
cetak("Port   :", port)    // 8080 (menggunakan nilai yang ada)
```

#### Modifikasi, Penggabungan (*Merge*), dan Pembersihan
```celang
buat akun = {"pengguna": "yosatsu"}

// Menyimpan nilai dengan chaining
akun.set("peran", "admin").set("aktif", benar)

// Menggabungkan dengan map preferensi lain
buat preferensi = {"tema": "gelap", "notifikasi": benar}
akun.merge(preferensi)

cetak(akun.pengguna)   // "yosatsu"
cetak(akun.tema)       // "gelap"

// Menghapus kunci
buat berhasilHapus = akun.delete("aktif")
cetak("Terhapus:", berhasilHapus) // benar

// Membersihkan seluruh isi map
akun.clear()
cetak("Jumlah setelah dibersihkan:", akun.length) // 0
```

---

## 3. Manipulasi Koleksi Bita Biner (`Bita` / `Bytes`)

Koleksi bita biner `Bita` di Celang dirancang sebagai wadah byte efisien (`[]byte`) untuk memanipulasi berkas biner, soket jaringan, dan enkripsi. Objek `Bita` memiliki serangkaian metode bawaan yang mendukung nama Bahasa Indonesia dan nama internasional (Universal).

### Ringkasan Metode Bita

| Metode (Universal) | Metode (Indonesia) | Argumen | Keterangan | Nilai Kembali |
| :--- | :--- | :--- | :--- | :--- |
| `.to_string()` | `.ke_teks()` | *(tidak ada)* | Mengonversi isi bita menjadi tulisan string teks UTF-8 | `Teks` |
| `.to_hex()` | `.ke_heks()` | *(tidak ada)* | Mengonversi deretan bita menjadi string heksadesimal | `Teks` |
| `.to_base64()` | `.ke_base64()` | *(tidak ada)* | Mengonversi deretan bita menjadi string Base64 | `Teks` |
| `.to_array()` | `.ke_daftar()` | *(tidak ada)* | Mengonversi bita menjadi array berisi angka bulat `0..255` | `Array` |
| `.append(...)` | `.tambah(...)` | `data_bita` | Menambahkan bita, teks, atau array integer ke akhir bita | Bita (memungkinkan *chaining*) |
| `.slice(mulai, [akhir])` | `.potong(mulai, [akhir])` | `mulai, [akhir]` | Mengambil cuplikan rentang bita | Objek `Bita` baru |
| `.clone()` | `.salin()` | *(tidak ada)* | Membuat salinan utuh (*deep copy*) mandiri | Objek `Bita` baru |
| `.length` / `.len` | `.panjang` | *(properti)* | Mengambil ukuran total bita | `Angka` (*Integer*) |

---

### Contoh Penggunaan Bita

```celang
// Membuat bita dari teks
buat b = bita("Celang")

// Menambah bita secara fluent chaining
b.tambah(" v2").tambah([33]) // kode ASCII 33 adalah '!'
cetak(b.ke_teks())           // "Celang v2!"

// Mengakses dan mengubah nilai byte langsung
cetak(b[0])                  // 67 (huruf 'C')
b[0] = 71                    // ganti dengan 'G' (ASCII 71)
cetak(b.ke_teks())           // "Gelang v2!"

// Enkoding heksadesimal dan Base64
cetak(b.ke_heks())           // "47656c616e6720763221"
cetak(b.ke_base64())         // "R2VsYW5nIHYyIQ=="

// Dekode Base64 & Heksadesimal menjadi Bita
buat bBase64 = bita.dari_base64("R2VsYW5nIHYyIQ==")
buat bHeks   = bita.dari_heks("47656c616e6720763221")
cetak(bBase64.ke_teks())     // "Gelang v2!"
cetak(bHeks.ke_teks())       // "Gelang v2!"

// Atau menggunakan konstruktor dengan argumen format:
buat b3 = bita("R2VsYW5n", "base64")
buat b4 = bita("43656c616e67", "heks") // "Celang"

// Operasi potong (slice)
buat potong = b[0..6]        // bita([71, 101, 108, 97, 110, 103])
cetak(potong.ke_teks())      // "Gelang"

// Menggandakan bita (operator perkalian)
buat ulang = bita("ABC") * 3
cetak(ulang.ke_teks())       // "ABCABCABC"
```

---

## 4. Dukungan Mode Universal (`"use universal"`)

Di bawah direktif `"use universal"` pada baris pertama file program, Celang memberlakukan standar metode internasional dalam Bahasa Inggris.

Apabila nama metode Bahasa Indonesia dipanggil secara tidak sengaja pada Mode Universal, Celang akan memberikan pesan empati yang mengarahkan Anda ke metode standar yang sesuai:

```celang
"use universal"

let data = [3, 1, 2]
data.push(4) // ✅ Valid
data.sort() // ✅ Valid

// Memanggil nama metode Indonesia pada Mode Universal:
// data.tambah(5)
// ❌ Kesalahan: Indonesian method 'tambah' is not allowed in Universal Mode.
//    Saran: Use '.push(...)' instead.

let profil = {"name": "Alex"}
profil.set("role", "developer") // ✅ Valid
let k = profil.keys()           // ✅ Valid

// profil.kunci()
// ❌ Kesalahan: Indonesian method 'kunci' is not allowed in Universal Mode.
//    Saran: Use '.keys()' instead.
```

Dengan fitur ini, kode Celang Anda dapat ditulis secara fleksibel dalam Bahasa Indonesia untuk kenyamanan lokal, maupun secara konsisten dalam Bahasa Inggris untuk proyek kolaborasi global.
