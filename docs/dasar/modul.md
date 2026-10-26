---
outline: deep
---

# Sistem Modul & Impor (`ambil`, `impor`, `dari`)

Ketika program Anda bertambah besar, membagi kode ke dalam beberapa berkas terpisah (modul) sangat penting agar kode tetap rapi, mudah dirawat, dan dapat digunakan kembali.

Celang memiliki sistem modul modern dengan aturan impor yang sangat bersih (*clean import syntax*) serta kontrol visibilitas ekspor yang presisi.

---

## 💡 Filosofi Impor Bersih (*Clean Imports*)

Celang membedakan antara modul pustaka/paket dengan berkas lokal secara visual:

| Jenis Impor | Sintaks | Contoh |
| :--- | :--- | :--- |
| **Pustaka Standar (*Builtin*)** | **Tanpa tanda kutip** | `ambil { akar, PI } dari matematika` |
| **Paket Eksternal (*Package*)** | **Tanpa tanda kutip** | `ambil { router } dari express` |
| **Berkas Jalur Lokal (*File Path*)** | **Menggunakan tanda kutip `""`** | `ambil { helper } dari "./helper.ce"` |

> 💡 **Catatan Kompatibilitas:** Impor modul pustaka menggunakan tanda kutip (`dari "matematika"`) tetap didukung penuh agar kode versi terdahulu tidak rusak.

---

## 📥 Sintaks Impor

Celang mendukung beragam variasi impor yang ergonomis sesuai kebiasaan dan gaya bahasa pengembang:

### 1. Impor Simbol Spesifik (*Named Import*)
Mengambil hanya fungsi, konstanta, atau tipe yang dibutuhkan:

```celang
// Dari pustaka bawaan atau paket (tanpa kutip):
ambil { akar, PI } dari matematika
ambil { uraikan, rapi } dari json
ambil { get, post, server } dari http

// Dari berkas lokal (menggunakan kutip):
ambil { formatTanggal, namaBulan } dari "./utilitas.ce"
```

### 2. Menggunakan Nama Alias (*Import with Alias*)
Jika nama simbol terlalu panjang atau terjadi konflik nama dengan variabel lokal:

```celang
ambil { tambah sebagai penjumlahan } dari "./kalkulator.ce"

buat hasil = penjumlahan(10, 20)
```

### 3. Impor Seluruh Modul (*Direct Module Import*)
Mengimpor seluruh isi modul menjadi satu objek namespace:

```celang
// Memberikan alias namespace:
impor matematika sebagai math
impor sistem sebagai os

buat luas = math.PI * math.pangkat(r, 2)
buat cwd = os.direktori_kerja()

// Atau impor modul langsung:
impor matematika
buat nilai = matematika.akar(144)
```

Atau menggunakan sintaks bintang:
```celang
ambil * sebagai math dari matematika
```

### 4. Gaya Bahasa Pythonic (`dari ... ambil`)
Bagi pengembang yang menyukai gaya sintaks Python (`from ... import ...`):

```celang
dari matematika ambil { akar, sinus }
dari json ambil uraikan
dari matematika ambil *

from "./utilitas.ce" import cetak_baris, cetak_kolom as kolom
```

> 🌟 **Pustaka Standar Tersedia Otomatis (*Auto-Injected Globals*)**
>
> Seluruh modul standar inti di Celang (`berkas`, `matematika`, `teks`, `waktu`, `json`, `sistem`, `http`, `ffi`) sudah **otomatis terdaftar sebagai objek global**.
> Anda dapat langsung menulis `berkas.baca("catatan.txt")` atau `matematika.akar(16)` **tanpa harus memanggil perintah `impor` terlebih dahulu**.
>
> Kata kunci `impor` / `ambil` digunakan ketika:
> 1. Mengimpor berkas kode lokal buatan Anda sendiri (misal: `ambil { formatData } dari "./util.ce"`).
> 2. Mengambil simbol tertentu secara langsung ke ruang lingkup lokal (misal: `ambil { akar, PI } dari matematika`).
> 3. Memberikan alias namespace pendek (misal: `impor matematika sebagai math`).
> 4. Menggunakan modul/paket eksternal yang diinstal melalui Package Manager.

---

## 🔒 Kontrol Visibilitas Ekspor

Di Celang, secara bawaan (*default*), semua variabel, fungsi, dan tipe bersifat **privat** untuk berkas tempat mereka dideklarasikan. Agar dapat diakses dari berkas lain, Anda perlu menambahkan pengubah visibilitas:

### 1. `publik!` (Dapat Diakses Bebas)
Dapat diimpor oleh berkas atau modul mana pun di dalam proyek:

```celang
// Diekspor langsung saat deklarasi:
publik! fungsi tambah(a, b) {
    kembalikan a + b
}

publik! buat VERSI = "2.0.0"

// Atau diekspor secara terpusat di akhir berkas:
fungsi kali(a, b) { kembalikan a * b }
fungsi bagi(a, b) { kembalikan a / b }

publik! { kali, bagi }
```

### 2. `privat-modul!` (Khusus Modul yang Sama)
Hanya dapat diakses oleh berkas yang berada di dalam lingkup modul yang sama (direktori yang memiliki penanda `celang.mod` atau `modul.ce`):

```celang
// Hanya boleh digunakan oleh berkas-berkas internal modul:
privat-modul! fungsi inisialisasiDatabase() {
    // ...
}
```

Jika berkas di luar modul mencoba mengimpor fungsi ini, Celang akan menolaknya saat evaluasi untuk menjaga arsitektur aplikasi tetap bersih.

### 3. `privat-["./berkas.ce"]!` (*Allowlist* Khusus)
Membatasi akses simbol hanya untuk berkas atau jalur tertentu:

```celang
// Hanya berkas "./pengujian.ce" yang boleh mengimpor fungsi rahasia ini:
privat-["./pengujian.ce"]! buat kunciRahasia = "TOKEN_12345"
```

---

## 📦 Resolusi Paket (`celang_modules/`)

Ketika Anda mengimpor nama paket tanpa awalan jalur (misal: `ambil { sesuatu } dari nama_paket`), Celang secara cerdas mencari paket tersebut di dalam folder `celang_modules/` dengan menelusuri folder saat ini dan folder-folder induk di atasnya:

1. `celang_modules/<paket>/utama.ce`
2. `celang_modules/<paket>/main.ce`
3. `celang_modules/<paket>/modul.ce`
4. `celang_modules/<paket>/<paket>.ce`
5. `celang_modules/<paket>/index.ce`
6. `celang_modules/<paket>.ce`

---

## 📂 Resolusi Impor Folder / Direktori

Celang mendukung impor langsung ke direktori tanpa perlu menulis nama berkas entrypoint secara manual. Ini sangat memudahkan pengorganisasian komponen atau sub-paket.

```celang
// Mengimpor dari folder "./komponen"
ambil { Tombol, InputTeks } dari "./komponen"
```

Ketika jalur yang dituju adalah sebuah folder, Celang akan mencari berkas entrypoint dengan urutan prioritas:
1. `modul.ce` *(Konvensi utama Celang)*
2. `utama.ce`
3. `index.ce`
4. `<nama_folder>.ce` *(misal: `komponen/komponen.ce`)*
5. Varian `.ast` (U-AST terkompilasi)

---

## 🔄 Re-Ekspor & Pola Barrel (*Barrel Exports*)

Pola *barrel* memungkinkan sebuah berkas menyatukan dan mengekspor kembali simbol-simbol dari sub-berkas lain, sehingga pengguna modul hanya perlu mengimpor dari satu pintu masuk utama.

### 1. Re-ekspor Simbol Tertentu (*Named Re-export*)
```celang
// Di dalam komponen/modul.ce:
ekspor { Tombol, Modal } dari "./tombol.ce"
ekspor { Input as KolomInput } dari "./input.ce"

// Mode universal:
export { Button, Modal } from "./button.ce"
```

### 2. Re-ekspor Seluruh Simbol (*Wildcard Re-export*)
```celang
// Mengekspor semua simbol publik dari sub-modul:
ekspor * dari "./utilitas.ce"

// Mode universal:
export * from "./utils.ce"
```

> 🛡️ **Penjaga Privasi:** Simbol privat atau bertanda `privat-modul!` pada sub-modul tidak akan bocor keluar melalui *barrel export*, memastikan enkapsulasi tetap terlindungi 100%.

---

## ⚡ Impor Dinamis (*Dynamic Import*)

Selain deklarasi impor statis di awal berkas, Celang menyediakan fungsi bawaan `impor_dinamis(path)` (atau `import_dynamic(path)` pada Mode Universal) untuk memuat modul secara kondisional atau saat runtime:

```celang
buat konfigurasiPlugin = "modul_tambahan.ce"

// Memuat modul saat runtime:
buat plugin = impor_dinamis("./plugins/" + konfigurasiPlugin)

// Hasil impor dinamis adalah Map berisi seluruh simbol publik modul:
buat hasil = plugin.jalankan()
buat daftarFungsi = plugin.kunci()

// Aman dibungkus dalam blok coba/tangkap:
coba {
    buat driver = impor_dinamis("./drivers/" + namaDriver)
    driver.hubungkan()
} tangkap (e) {
    cetak("Gagal memuat driver:", e)
}
```

Struktur ini dirancang agar siap menyambut *package manager* ekosistem Celang di masa mendatang!

