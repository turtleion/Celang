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
| `.pop()` | `.ambil_terakhir()` | *(tidak ada)* | Mengambil dan menghapus elemen terakhir dari array | Elemen yang terhapus / `kosong` |
| `.shift()` | `.geser()` | *(tidak ada)* | Mengambil dan menghapus elemen pertama dari array | Elemen yang terhapus / `kosong` |
| `.unshift(...)` | `.sisip_depan(...)` | `elemen, ...` | Menyisipkan satu atau lebih elemen di awal array | Array (memungkinkan *chaining*) |
| `.insert(indeks, el)`| `.sisip(indeks, el)` | `indeks, elemen` | Menyisipkan elemen pada indeks tertentu (mendukung indeks negatif) | Array (memungkinkan *chaining*) |
| `.remove_at(indeks)` | `.hapus_indeks(indeks)` | `indeks` | Menghapus elemen pada indeks tertentu | Elemen yang terhapus |
| `.reverse()` | `.balik()` | *(tidak ada)* | Membalikkan urutan elemen array secara langsung (*in-place*) | Array (memungkinkan *chaining*) |
| `.sort([pembanding])`| `.urut([pembanding])` | `[fungsi(a, b)]` | Mengurutkan elemen (alami atau menggunakan fungsi perbandingan kustom) | Array (memungkinkan *chaining*) |
| `.join([pemisah])` | `.gabung([pemisah])` | `[pemisah]` | Menggabungkan seluruh elemen menjadi satu teks (bawaan: `", "`) | `Teks` (*String*) |
| `.contains(el)` / `.includes(el)` | `.berisi(el)` | `elemen` | Memeriksa apakah suatu elemen ada di dalam array | `Logika` (`benar`/`salah`) |
| `.index_of(el)` | `.indeks_dari(el)` | `elemen` | Mencari posisi indeks pertama dari elemen (mengembalikan `-1` jika tidak ada) | `Angka` (*Integer*) |
| `.slice([mulai, [akhir]])` | `.potong([mulai, [akhir]])`| `[mulai, [akhir]]` | Mengambil salinan potongan array tanpa mengubah array asli | Salinan Array baru |
| `.concat(...)` | `.gabungkan(...)` | `array_atau_elemen, ...`| Menggabungkan array dengan elemen atau array lain | Array baru |
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

// Mengambil elemen terakhir & pertama
buat terakhir = angka.pop()   // 50
buat pertama  = angka.shift() // 10
cetak(terakhir, pertama)     // 50 10
cetak(angka)                 // [20, 30, 40]

// Sisip di depan atau pada indeks spesifik
angka.unshift(5)             // sisip di depan: [5, 20, 30, 40]
angka.insert(2, 25)          // sisip pada indeks 2: [5, 20, 25, 30, 40]
buat dihapus = angka.remove_at(1) // hapus elemen pada indeks 1 (nilai 20)
cetak(angka)                 // [5, 25, 30, 40]
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

---

## 2. Manipulasi Map (Kamus Pasangan Kunci-Nilai)

Map di Celang menyimpan pasangan kunci dan nilai terurut berdasarkan urutan penambahan (*insertion order*).

### Ringkasan Metode Map

| Metode (Universal) | Metode (Indonesia) | Argumen | Keterangan | Nilai Kembali |
| :--- | :--- | :--- | :--- | :--- |
| `.keys()` | `.kunci()` | *(tidak ada)* | Mengambil seluruh kunci map dalam urutan penambahan | `Array` berisi teks kunci |
| `.values()` | `.nilai()` | *(tidak ada)* | Mengambil seluruh nilai map dalam urutan penambahan | `Array` berisi nilai |
| `.entries()` | `.entri()` | *(tidak ada)* | Mengambil seluruh pasangan sebagai array `[[kunci, nilai], ...]` | `Array` 2-dimensi |
| `.has(kunci)` / `.contains(kunci)` | `.memiliki(kunci)` | `kunci` | Memeriksa apakah kunci tertentu ada di dalam map | `Logika` (`benar`/`salah`) |
| `.get(kunci, [bawaan])` | `.ambil(kunci, [bawaan])` | `kunci, [bawaan]`| Mengambil nilai berdasarkan kunci, dengan nilai cadangan opsional jika kunci tidak ditemukan | Nilai objek / `bawaan` / `kosong` |
| `.set(kunci, nilai)` | `.atur(kunci, nilai)` | `kunci, nilai` | Menyimpan pasangan kunci-nilai ke dalam map | Map (memungkinkan *chaining*) |
| `.delete(kunci)` / `.remove(kunci)` | `.hapus(kunci)` | `kunci` | Menghapus kunci dari map | `Logika` (`benar` jika dihapus, `salah` jika tidak ada) |
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

## 3. Dukungan Mode Universal (`"use universal"`)

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
