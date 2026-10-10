# Spesifikasi Sintaks Lanjutan Celang (`ADVANCED_SYNTAX.md`)

> **Catatan untuk Kamu:**  
> Dokumen ini berisi draf rancangan sintaks untuk fitur-fitur lanjutan Celang:  
> **Modul & Impor/Ekspor, Error Handling (Try-Catch), Pipeline Operator, In-Code Styling Blocks, dan Standard Library.**  
> Silakan review, ubah keyword, edit format, atau hapus bagian yang kurang cocok. Nanti setelah kamu selesai mengedit, beri tahu aku dan kita akan lanjutkan eksekusi di branch baru!

---

## 1. Sistem Modul & Impor / Ekspor

Memungkinkan kode dipecah ke dalam beberapa berkas terpisah untuk kerapian arsitektur proyek.

### 1.1 Visibilitas

Secara default, semua hal yang ada di suatu file, visibilitasnya itu privat untuk satu file itu sendiri. Jika file itu di dalam suato modul, maka secara default, semua hal yang ada di modul itu privat untuk modul itu sendiri (hanya yang didalam modul itu yang dapat mengakses variabel/fungsi/dll yang didalam modul itu)

Sebelum itu, keywords visibilitas itu:

- **publik!** Untuk visibilitas publik
- **privat-modul!** Membatasi visibilitas untuk modul
- **privat!** Membatasi visibilitas untuk file itu sendiri dan ini adalah properti default untuk semua.
- _[SPESIAL]_ **privat-[FilePath/Modul]**! Allowlist untuk file tertentu yang bisa ambil suatu nilai dari file tersebut

### 1.1.1 Non module, single file program.

#### Code/program1.ce

```celang
buat nama = "Kayasa"
```

#### Code/program2.ce

```celang
ambil { nama } dari "./program1.ce" <- ERROR! Nama tidak ditemukan
```

Bisa dihandle dengan..

#### Code/program1.ce

```celang
publik! buat nama = "Kayasa"
```

#### Code/program2.ce

```celang
ambil { nama } dari "./program1.ce" <- Nama = "Kayasa"
```

### 1.1.2 Module scoped.

#### Transport/lib/avt.ce

```celang
fungsi komunikasikan(ip: Teks) mengembalikan Kebenaran
```

#### Transport/core.ce

```celang
ambil { komunikasikan } dari "./lib/avt.ce" <- ERROR! komunikasikan gaada
```

Bisa dihandle dengan..

#### Transport/lib/avt.ce

```celang
privat-modul! fungsi komunikasikan(ip: Teks) mengembalikan Kebenaran
```

#### Transport/core.ce

```celang
ambil { komunikasikan } dari "./lib/avt.ce" <- komunikasikan = fungsi komunikasikan di avt.ce
```

Tapi... di scope lain!

#### Code/program3.ce

```celang
ambil { komunikasikan } dari "../Transport/lib/avt.ce" <- ERROR! komunikasikan gaada
```

### 1.1.3 privat-[FilePath/Modul], inovasi terbaru

#### Code/program4.ce

```celang
privat-["./program5.ce"]! buat nama = "abdi"
```

#### Code/program5.ce

```celang
ambil { nama } dari "./program4.ce" <- nama = "Abdi"
```

#### Code/program6.ce

```celang
ambil { nama } dari "./program4.ce" <- ERROR! nama gaada
```

#### 1.1.4 Publik

Sepertinya aku gausa jelasin ini. Intinya publik mengekspos ke siapapun

```celang
// File: matematika.ce

// Cara 1: Langsung di depan deklarasi
publik! fungsi tambah(a, b) {
    kembalikan a + b
}

privat-["./Modul-A"] konstanta PI = 3.14159

// Cara 2: Di akhir berkas
fungsi kali(a, b) {
    kembalikan a * b
}
publik! { kali }
```

### 1.2 Impor (`ambil ... dari ...`)

Mengambil fungsi atau nilai dari berkas lain.

#### Gaya Indonesia (Default):

```celang
// 1. Impor entitas tertentu (Named Import)
ambil { tambah, PI } dari "./matematika.ce"

// 2. Impor seluruh modul ke dalam satu objek namespace
ambil * sebagai math dari "./matematika.ce"

// 3. Impor dengan alias
ambil { tambah sebagai penjumlahan } dari "./matematika.ce"

// Penggunaan:
cetak(tambah(10, 20))
cetak(math.kali(5, 5))
```

#### Gaya JavaScript:

```javascript
import { tambah, PI } from "./matematika.ce";
import * as math from "./matematika.ce";
import { tambah as penjumlahan } from "./matematika.ce";
```

#### Gaya Pythonic:

```python
from "./matematika.ce" import tambah, PI
import "./matematika.ce" as math
from "./matematika.ce" import tambah as penjumlahan
```

---

## 2. Penanganan Error Runtime (`coba ... tangkap ... akhirnya`)

Menangani kesalahan saat program berjalan (_runtime exceptions_) dengan aman tanpa mematikan program secara kasar.

### 2.1 Sintaks Percabangan Error

#### Gaya Indonesia:

```celang
coba {
    buat hasil = 100 / 0
    cetak("Hasil:", hasil)
} tangkap (kesalahan) {
    cetak("Ups, terjadi masalah:", kesalahan)
} akhirnya {
    cetak("Operasi selesai, membersihkan memori.")
}
```

- Bagian `akhirnya` bersifat opsional (boleh ditiadakan).
- Variabel penampung kesalahan (misal: `kesalahan`) menyimpan objek pesan error.

#### Gaya JavaScript:

```javascript
try {
    let hasil = 100 / 0;
    cetak("Hasil:", hasil);
} catch (error) {
    cetak("Ups, terjadi masalah:", error);
} finally {
    cetak("Operasi selesai, membersihkan memori.");
}
```

#### Gaya Pythonic:

```python
try:
    hasil = 100 / 0
    print("Hasil:", hasil)
except error:
    print("Ups, terjadi masalah:", error)
finally:
    print("Operasi selesai, membersihkan memori.")
```

### 2.2 Melempar Kesalahan Manual (`lempar` / `throw` / `raise`)

Untuk memicu error validasi secara sengaja:

```celang
// Gaya Indonesia:
fungsi set_umur(u) {
    jika (u < 0) {
        lempar "Umur tidak boleh negatif!"
    }
    kembalikan u
}

// Gaya JS:
throw "Umur tidak boleh negatif!"

// Gaya Pythonic:
raise "Umur tidak boleh negatif!"
```

---

## 3. Operator Pipeline (`|>`) / Pengaliran Data

Memungkinkan pemrosesan data berantai (_chaining_) dari kiri ke kanan secara bersih, alami, dan mudah dibaca tanpa sarang tanda kurung yang rumit.

Nilai di sisi kiri `|>` otomatis disisipkan sebagai **argumen pertama** ke fungsi di sisi kanan.

```celang
// Bentuk konvensional (sulit dibaca / banyak kurung bertingkat):
cetak(kapitalisasi(potong_spasi("  halo dunia  ")))

// Bentuk Pipeline (Elegan & Mengalir):
"  halo dunia  "
    |> potong_spasi()
    |> kapitalisasi()
    |> cetak()
```

### Contoh Pemrosesan Koleksi (Transformasi Array):

```celang
buat angka = [1, 2, 3, 4, 5, 6]

buat hasil = angka
    |> saring(fungsi(x) { kembalikan x % 2 == 0 })
    |> petakan(fungsi(x) { kembalikan x * 10 })

cetak(hasil) // Output: [20, 40, 60]
```

_(Catatan: Operator `|>` otomatis dipertahankan atau ditransformasikan ke pemanggilan berantai/nested sesuai profile)._

---

## 4. In-Code Styling Blocks (`gaya "..." { }`)

Memungkinkan programmer mencampur dialek dalam satu berkas bila blok tertentu lebih enak ditulis dengan gaya lain (misal: penulisan algoritma dalam gaya Pythonic atau penulisan manipulasi web dalam gaya JS).

```celang
// Kode utama dalam gaya Indonesia
buat nama_pengguna = "Budi"

// Blok dengan gaya JavaScript
gaya "javascript" {
    let status_login = true
    function cek_akses(user) {
        return user === "Budi" && status_login
    }
    let diizinkan = cek_akses(nama_pengguna)
}

// Kembali ke konteks global
jika (diizinkan) {
    cetak("Selamat datang, " + nama_pengguna)
}
```

---

## 5. Standard Library Bawaan (Baterai Siap Pakai)

Dapat diakses langsung atau melalui modul bawaan.

### 5.1 Modul Berkas (`berkas` / `fs`)

Manipulasi file dan direktori di komputer:

- `berkas.baca(path)`: Membaca seluruh isi file teks.
- `berkas.tulis(path, data)`: Menulis atau menimpa isi file.
- `berkas.tambah(path, data)`: Menambahkan teks ke akhir file (_append_).
- `berkas.ada(path)`: Mengecek apakah file/folder ada (`benar` / `salah`).
- `berkas.hapus(path)`: Menghapus file.

### 5.2 Modul Matematika (`matematika` / `math`)

- `matematika.akar(16)` -> `4`
- `matematika.pangkat(2, 3)` -> `8`
- `matematika.mutlak(-5)` -> `5` (Absolute)
- `matematika.bulat(3.7)` -> `4` (Round)
- `matematika.lantai(3.9)` -> `3` (Floor)
- `matematika.atap(3.1)` -> `4` (Ceil)
- `matematika.acak(1, 100)` -> Angka integer acak antara 1 s/d 100.
- `matematika.PI` -> `3.141592653589793`

### 5.3 Modul Teks & String (`teks` / `strings`)

- `teks.pisah(kalimat, pemisah)`: Memecah string menjadi array.
- `teks.gabung(array, pemisah)`: Menggabungkan array string menjadi satu kalimat.
- `teks.kapital(str)` / `teks.besar(str)`: Huruf kapital semua (_UPPERCASE_).
- `teks.kecil(str)`: Huruf kecil semua (_lowercase_).
- `teks.potong(str)`: Menghapus spasi di awal dan akhir (_trim_).
- `teks.ganti(str, pola, baru)`: Mengganti substring.
- `teks.berisi(str, cari)`: Mengecek apakah string mengandung teks tertentu (`benar` / `salah`).

### 5.4 Modul Waktu (`waktu` / `time`)

- `waktu.sekarang()`: Mendapatkan timestamp atau tanggal sekarang.
- `waktu.tidur(milidetik)`: Menunda eksekusi (sleep).

### 5.5 Modul Format JSON (`json`)

- `json.uraikan(teks_json)`: Mengubah string JSON menjadi Map/Array Celang.
- `json.format(objek)`: Mengubah Map/Array Celang menjadi string JSON.

---

## 6. Pertanyaan & Area Penyesuaian Sebelum Kita Mulai Coding

Silakan ubah langsung di file ini atau beri tahu aku preferensimu:

1. **Keyword Impor & Ekspor:**
    - Apakah kata `ambil ... dari ...` dan `ekspor` sudah pas? Belum. Aku ubah total ekspor jadi visibilitas.
    - Atau kamu lebih suka `impor ... dari ...` / `masukkan` / `muat`? Biarkan impor ... dari ...
2. **Keyword Try-Catch:**
    - Apakah pasangan `coba { } tangkap (e) { } akhirnya { }` sudah nyaman dan manusiawi? Iya
    - Untuk melempar error: apakah `lempar "error"` sudah bagus, atau ada ide lain? Iya
3. **Pipeline Operator:**
    - Apakah simbol `|>` sudah oke (standar industri di Elixir, F#, proposal JS), atau mau dibuat keyword bahasa seperti `alirkan ke`? Itu aja.
4. **Nama-nama Modul Standard Library:**
    - Apakah kamu lebih suka nama bahasa Indonesia (`berkas`, `matematika`, `teks`, `waktu`) dengan alias bahasa Inggris (`fs`, `math`, `strings`, `time`), atau bagaimana? Iya
5. **Urutan Pengerjaan:**
    - Di antara fitur-fitur di atas, fitur mana yang ingin diprioritaskan terlebih dahulu di branch baru? Visibilitas dan impor paling wajib, lalu error handling dan pipeline, dan incode styling. Kalau modul bawaan nanti aja
