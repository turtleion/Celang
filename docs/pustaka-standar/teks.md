# Pustaka Pengolah Teks (`teks` / `strings`)

Modul **`teks`** (dapat juga dipanggil dengan alias internasional **`strings`**) menyediakan ekosistem terlengkap untuk memanipulasi, mencari, mengubah format, dan menganalisis kalimat di Celang.

Selain operasi pengolahan teks standar seperti memotong spasi dan memecah kalimat, modul `teks` di Celang dirancang dengan **fitur-fitur cerdas** bawaan:
- **Pembuat Slug URL Otomatis (`teks.slug`)**: Mengonversi judul artikel atau kalimat menjadi tautan URL ramah web.
- **Penyamaran Data Sensitif (`teks.samarkan`)**: Sensor otomatis untuk email, nomor ponsel, NIK, dan teks pribadi.
- **Toleransi Salah Ketik & Kemiripan Levenshtein (`teks.kemiripan` & `teks.terdekat`)**: Mengukur kesamaan kata dan mencari kata paling mirip dari kumpulan pilihan.
- **Format Templat Serbaguna (`teks.format`)**: Format teks dinamis berbasis pemetaan kamus (`{kunci}`) maupun indeks posisi (`{0}`, `{1}`).
- **Konversi Gaya Huruf Beragam**: Mendukung *camelCase* (`teks.unta`), *snake_case* (`teks.ular`), *kebab-case* (`teks.kebab`), *Title Case* (`teks.judul`), hingga kapital awal (`teks.kapital_awal`).
- **Pola & Ekspresi Reguler (Regex)**: Validasi pola (`teks.sesuai`), pencarian banyak kecocokan (`teks.temukan_semua`), dan substitusi pola (`teks.ganti_pola`).
- **Perataan Teks & Padding**: Meratakan teks ke kiri, kanan, atau tengah tabel terminal (`teks.isi_kiri`, `teks.isi_kanan`, `teks.tengah`).

---

## Cara Memanggil

Modul ini aktif secara global dan dapat dipanggil menggunakan nama lokal maupun alias internasional:

```celang
// 1. Menggunakan nama Bahasa Indonesia:
buat judulWeb = teks.slug("Belajar Celang 2.0: Cepat & Mudah!")
cetak(judulWeb) // "belajar-celang-2-0-cepat-dan-mudah"

// 2. Menggunakan alias internasional:
buat hasil = strings.trim_left("===DATA===", "=")
cetak(hasil) // "DATA==="

// 3. Atau menggunakan sintaks impor terstruktur:
impor { slug, samarkan, kemiripan, terdekat } dari "teks"
```

---

## 🌟 Fitur Cerdas & Unik Celang

### 1. Pembuat Slug URL (`slug`)
Mengubah kalimat bebas menjadi format URL slug yang aman, bersih, dan ramah SEO. Karakter khusus seperti tanda dan (`&`) otomatis diubah menjadi kata penghubung bahasa Indonesia (`dan`).

```celang
buat judulArtikel = "Panduan Belajar Celang 2.0 & Arsitektur U-AST!"
buat tautan = teks.slug(judulArtikel)

cetak(tautan)
// Hasil: "panduan-belajar-celang-2-0-dan-arsitektur-u-ast"
```

---

### 2. Penyamar Data Sensitif (`samarkan` / `mask`)
Mencegah kebocoran data pribadi (PII) pada tampilan log atau profil pengguna secara otomatis:
- **Email**: Otomatis mendeteksi alamat email dan menyamarkan bagian tengah nama pengguna (misal: `b***o@domain.com`).
- **Nomor Telepon / NIK / Kartu**: Otomatis menampilkan beberapa digit awal dan akhir, menyamarkan bagian tengahnya.
- **Karakter Kustom**: Karakter sensor dapat diganti sesuai kebutuhan (default: `*`).

```celang
// Sensor email:
buat emailTersamar = teks.samarkan("budi.santoso@perusahaan.com")
cetak(emailTersamar) // "b***o@perusahaan.com"

// Sensor nomor kontak:
buat kontakTersamar = teks.samarkan("081234567890")
cetak(kontakTersamar) // "0812****7890"

// Sensor dengan karakter kustom (misal tanda pagar '#'):
buat pinTersamar = teks.samarkan("123456", "#")
cetak(pinTersamar) // "1####6"
```

---

### 3. Kemiripan Kata / Levenshtein Distance (`kemiripan` / `similarity`)
Menghitung derajat kesamaan antara dua teks dengan skala angka desimal `0.0` (sama sekali berbeda) hingga `1.0` (persis identik):

```celang
buat skor1 = teks.kemiripan("jakarta", "jakartaa")
cetak("Skor kemiripan:", skor1) // ~0.875 (sangat mirip)

buat skor2 = teks.kemiripan("apel", "kucing")
cetak("Skor kemiripan:", skor2) // 0.0 (tidak mirip)
```

---

### 4. Pencari Kata Terdekat / Fuzzy Matcher (`terdekat` / `closest`)
Membantu menemukan kata yang paling cocok dari sekumpulan daftar pilihan kata. Fitur ini sangat bermanfaat untuk rekomendasi salah ketik (*did you mean?*) atau pencarian toleran typo:

```celang
buat kamusBuah = ["apel", "jeruk", "mangga", "pisang", "semangka"]

// Pengguna salah mengetik 'aple':
buat saran = teks.terdekat("aple", kamusBuah)
cetak("Maksud Anda mungkin:", saran) // "apel"

// Pengguna mencari 'jruk':
buat saran2 = teks.terdekat("jruk", kamusBuah)
cetak("Maksud Anda mungkin:", saran2) // "jeruk"
```

---

### 5. Format Templat Serbaguna (`format`)
Menyusun kalimat dinamis menggunakan placeholder berbasis nama kamus (`Map`) atau posisi indeks urutan (`Array` / argumen langsung):

```celang
// A. Format berbasis Kamus data (Map):
buat templateKamus = "Halo \{nama}, status akun Anda: \{status}."
buat hasil1 = teks.format(templateKamus, {
    "nama": "Andi",
    "status": "Aktif"
})
cetak(hasil1) // "Halo Andi, status akun Anda: Aktif."

// B. Format berbasis Indeks Posisi ({0}, {1}, dst):
buat templatePos = "Barang \{0} berjumlah \{1} buah, total: Rp\{2}."
buat hasil2 = teks.format(templatePos, "Buku", 3, 45000)
cetak(hasil2) // "Barang Buku berjumlah 3 buah, total: Rp45000."
```

> [!NOTE]
> Di Celang, tanda `{variabel}` di dalam kutip teks merupakan fitur interpolasi bawaan. Jika Anda ingin menuliskan template teks secara harfiah, beri tanda garis miring terbalik di depan kurung kurawal pembuka: `\{kunci}` atau baca template tersebut langsung dari berkas konfigurasi.

---

## 🔠 Transformasi Gaya Penulisan Huruf (Case Styles)

Celang menyediakan fungsi transformasi gaya penulisan yang umum dipakai dalam rekayasa perangkat lunak:

```celang
// 1. Kapital huruf pertama kalimat/kata (Capitalize):
cetak(teks.kapital_awal("celang")) // "Celang"

// 2. Gaya Judul Buku (Title Case):
cetak(teks.judul("panduan bahasa celang")) // "Panduan Bahasa Celang"

// 3. Gaya Unta (camelCase):
cetak(teks.unta("nama_pengguna_aktif")) // "namaPenggunaAktif"
cetak(teks.unta("Halo Dunia"))         // "haloDunia"

// 4. Gaya Ular (snake_case):
cetak(teks.ular("namaPenggunaAktif")) // "nama_pengguna_aktif"
cetak(teks.ular("Halo Dunia"))        // "halo_dunia"

// 5. Gaya Kebab (kebab-case):
cetak(teks.kebab("Halo Dunia Celang")) // "halo-dunia-celang"
cetak(teks.kebab("data_pelanggan"))    // "data-pelanggan"
```

---

## 🔍 Pola & Ekspresi Reguler (Regex)

Untuk validasi formulir dan ekstraksi pola teks tingkat lanjut:

### 1. `sesuai` / `apakah_cocok` / `match`
Memeriksa apakah seluruh atau sebagian teks memenuhi kriteria pola regex:

```celang
buat polaEmail = "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"

jika teks.sesuai("pengguna@celang.id", polaEmail) {
    cetak("Alamat surel valid!")
}
```

### 2. `temukan_semua` / `find_all`
Mengekstrak seluruh potongan kalimat yang sesuai dengan pola regex menjadi sebuah Daftar (`Array`):

```celang
buat catatan = "Pesanan #101 dan #105 selesai diproses, sedangkan #209 menunggu pengiriman."
buat nomorPesanan = teks.temukan_semua(catatan, "#[0-9]+")

cetak(nomorPesanan) // ["#101", "#105", "#209"]
```

### 3. `ganti_pola` / `replace_regex`
Mengganti teks berdasarkan pencocokan pola:

```celang
buat teksAcak = "Harga: 50000 rupiah dan 25000 rupiah"
buat disensor = teks.ganti_pola(teksAcak, "[0-9]+", "XXX")

cetak(disensor) // "Harga: XXX rupiah dan XXX rupiah"
```

---

## 📏 Perataan & Pengisian Karakter (Padding & Alignment)

Sangat berguna untuk menyusun tabel laporan rapi di layar terminal / konsol:

```celang
// 1. Pengisian dari kiri (pad_start):
cetak(teks.isi_kiri("42", 5, "0"))      // "00042"
cetak(teks.isi_kiri("Rp500", 10, " "))  // "     Rp500"

// 2. Pengisian dari kanan (pad_end):
cetak(teks.isi_kanan("Barang", 10, ".")) // "Barang...."

// 3. Meratakan ke tengah (center):
cetak(teks.tengah("MENU UTAMA", 20, "=")) 
// "=====MENU UTAMA====="
```

---

## 🛠️ Manipulasi & Pencarian Dasar

```celang
buat kalimat = "  Halo Celang Dunia  "

// 1. Memotong spasi / karakter pinggir:
buat bersih = teks.potong(kalimat)         // "Halo Celang Dunia"
buat potongKiri = teks.potong_kiri("---teks---", "-")  // "teks---"
buat potongKanan = teks.potong_kanan("---teks---", "-") // "---teks"

// 2. Mengubah kapital / huruf kecil:
cetak(teks.kapital("celang")) // "CELANG"
cetak(teks.kecil("CELANG"))   // "celang"

// 3. Mencari posisi indeks & frekuensi kemunculan:
buat pesan = "kopi kopi kopi"
cetak(teks.indeks_dari(pesan, "kopi"))     // 0
cetak(teks.indeks_terakhir(pesan, "kopi")) // 10
cetak(teks.hitung(pesan, "kopi"))          // 3

// 4. Mengulang teks:
cetak(teks.ulangi("Abc ", 3)) // "Abc Abc Abc "

// 5. Memecah dan Menggabung:
buat daftarKata = teks.pisah("apel,jeruk,mangga", ",")
cetak(daftarKata) // ["apel", "jeruk", "mangga"]

buat teksGabung = teks.gabung(daftarKata, " - ")
cetak(teksGabung) // "apel - jeruk - mangga"

// 6. Memeriksa keberadaan teks:
cetak(teks.berisi("halo dunia", "dunia"))           // benar
cetak(teks.dimulai_dengan("celang.go", "celang"))   // benar
cetak(teks.diakhiri_dengan("dokumen.pdf", ".pdf"))  // benar

// 7. Penggantian teks biasa:
cetak(teks.ganti("kucing merah", "merah", "putih")) // "kucing putih"
```

---

## 📋 Tabel Rangkuman Fungsi

| Fungsi Bahasa Indonesia | Alias Internasional | Keterangan Singkat |
| :--- | :--- | :--- |
| **`slug(teks)`** | `slug` | Mengubah kalimat menjadi URL slug SEO yang bersih |
| **`samarkan(teks, [simbol])`** | `mask` | Sensor otomatis email / nomor telepon / teks sensitif |
| **`kemiripan(teksA, teksB)`** | `similarity` | Skor kemiripan teks (0.0 sampai 1.0) |
| **`terdekat(target, daftar)`** | `closest` | Mencari kata paling mirip dari kumpulan opsi |
| **`format(template, data)`** | `format` | Format teks dinamis via kamus atau indeks posisi |
| **`kapital_awal(teks)`** | `capitalize` | Mengubah huruf pertama teks menjadi huruf besar |
| **`judul(teks)`** | `title` | Mengubah teks ke gaya Title Case |
| **`unta(teks)`** | `camel` | Mengonversi format kalimat ke `camelCase` |
| **`ular(teks)`** | `snake` | Mengonversi format kalimat ke `snake_case` |
| **`kebab(teks)`** | `kebab` | Mengonversi format kalimat ke `kebab-case` |
| **`sesuai(teks, pola)`** | `match` / `apakah_cocok` | Memeriksa kecocokan pola Regex (boolean) |
| **`temukan_semua(teks, pola)`** | `find_all` | Mengambil seluruh kecocokan Regex ke dalam daftar |
| **`ganti_pola(teks, pola, baru)`**| `replace_regex` | Substitusi teks dengan aturan Regex |
| **`isi_kiri(teks, len, [pad])`** | `pad_start` / `pad_left` | Menambahkan karakter pengisi di sisi kiri |
| **`isi_kanan(teks, len, [pad])`**| `pad_end` / `pad_right` | Menambahkan karakter pengisi di sisi kanan |
| **`tengah(teks, len, [pad])`** | `center` | Memposisikan teks di tengah dengan karakter pengisi |
| **`indeks_dari(teks, sub)`** | `index_of` | Mencari indeks pertama kemunculan substring (-1 bila tiada) |
| **`indeks_terakhir(teks, sub)`**| `last_index_of` | Mencari indeks terakhir kemunculan substring |
| **`hitung(teks, sub)`** | `count` | Menghitung banyaknya kemunculan substring |
| **`ulangi(teks, n)`** | `repeat` | Menduplikasi teks sebanyak $n$ kali |
| **`potong(teks, [cut])`** | `trim` | Memotong spasi atau karakter dari kedua ujung |
| **`potong_kiri(teks, [cut])`** | `trim_left` | Memotong spasi atau karakter dari ujung kiri |
| **`potong_kanan(teks, [cut])`**| `trim_right` | Memotong spasi atau karakter dari ujung kanan |
| **`pisah(teks, pemisah)`** | `split` | Memecah teks menjadi potongan daftar kata |
| **`gabung(daftar, pemisah)`**| `join` | Menggabungkan daftar kata menjadi kalimat utuh |
| **`kapital(teks)`** | `upper` | Mengubah seluruh huruf menjadi huruf besar |
| **`kecil(teks)`** | `lower` | Mengubah seluruh huruf menjadi huruf kecil |
| **`ganti(teks, lama, baru)`** | `replace` | Mengganti kemunculan teks lama dengan teks baru |
| **`berisi(teks, sub)`** | `contains` | Memeriksa apakah teks memuat kata tertentu |
| **`dimulai_dengan(teks, awal)`**| `starts_with` | Memeriksa apakah teks diawali substring tertentu |
| **`diakhiri_dengan(teks, akhir)`**| `ends_with` | Memeriksa apakah teks diakhiri substring tertentu |
