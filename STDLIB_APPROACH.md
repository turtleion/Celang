# Pendekatan Arsitektur Standard Library Celang (`STDLIB_APPROACH.md`)

> **Pertanyaan Mendasar dari Pengembang:**  
> *"Ga mungkin kan kita buat Standard Library murni dari Celang (misal berkas, kan harus gunain third party sebagai eksekutor, kayak Go/C/Dll.)? Gimana kita bakal nerapin ini?"*

Jawaban singkatnya: **Tepat sekali!**  
Bahasa pemrograman apa pun di dunia (Python, Node.js, Ruby, Go, Rust) **tidak pernah** bisa mengakses perangkat keras atau sistem operasi (seperti membaca harddisk, membuka port jaringan, atau mengakses jam CPU) hanya menggunakan kode tingkat tinggi murni tanpa jembatan (*bridge*) ke sistem di bawahnya.

Dokumen ini menjelaskan secara gamblang **arsitektur, strategi pengikatan host (Go runtime), dan mekanisme modul bawaan** untuk Celang Standard Library.

---

## 1. Realitas Fundamental: Batas Sistem Operasi (*Syscall Boundary*)

Di sistem operasi (Linux, Windows, macOS), operasi berkas dan perangkat keras dilindungi oleh Kernel. Program hanya bisa mengaksesnya melalui *System Calls* (`sys_read`, `sys_write`, `sys_open`, dll.).

```
┌────────────────────────────────────────────────────────┐
│             Kode Pengguna Celang (.ce)                 │
│         buat data = berkas.baca("./data.txt")          │
└───────────────────────────┬────────────────────────────┘
                            │ (1) U-AST / CallExpression
┌───────────────────────────▼────────────────────────────┐
│              Celang Evaluator Runtime                  │
└───────────────────────────┬────────────────────────────┘
                            │ (2) Marshaling Bridge
┌───────────────────────────▼────────────────────────────┐
│        Go Host Engine Backend (`pkg/stdlib/`)          │
│          os.ReadFile("./data.txt")                     │
└───────────────────────────┬────────────────────────────┘
                            │ (3) Go Runtime Syscall
┌───────────────────────────▼────────────────────────────┐
│                 Kernel Sistem Operasi                  │
│               sys_openat -> sys_read                   │
└────────────────────────────────────────────────────────┘
```

### Bagaimana Bahasa Lain Menanganinya?
1. **Python (CPython):** Menggunakan bahasa C. Ketika kamu menulis `import os` atau `open()`, Python memanggil kode C di berkas `posixmodule.c` atau `_iomodule.c`.
2. **Node.js:** Menggunakan C++ dan pustaka `libuv`. Ketika kamu menulis `fs.readFileSync()`, Node.js memanggil C++ binding ke fungsi POSIX `read()`.
3. **Celang:** Menggunakan **Go Engine**. Karena interpreter Celang dibangun di atas bahasa Go, Go bertindak sebagai *Host Execution Engine* dan *Foreign Function Interface (FFI)* bawaan kita. Pustaka standar Go (`os`, `io`, `math`, `strings`, `time`, `encoding/json`) sangat cepat, aman dari memory corruption, dan bekerja secara *cross-platform* di Linux, macOS, dan Windows.

---

## 2. Arsitektur 3 Lapis Standard Library Celang

Untuk memastikan kode Celang tetap bersih dan tidak tercemar oleh tipe data internal Go, kita menerapkan pemisahan 3 lapis:

### Lapis 1: Native Host Executor (Go Backend)
Berkas Go di dalam direktori `pkg/stdlib/` yang bertugas mengeksekusi operasi sistem riil.
- Contoh: `pkg/stdlib/berkas.go`, `pkg/stdlib/matematika.go`, `pkg/stdlib/waktu.go`, dll.

### Lapis 2: Marshaling Bridge & Empathetic Error Wrapper
Jembatan penerjemah antara tipe data Go dan objek U-AST Celang:
| Tipe Data Asli Go | Objek Celang (`evaluator.Object`) | Catatan Marshaling |
|---|---|---|
| `string` / `[]byte` | `*evaluator.String` | Diubah menjadi teks Celang |
| `int64` / `int` | `*evaluator.Integer` | Bilangan bulat 64-bit |
| `float64` | `*evaluator.Float` | Bilangan desimal presisi ganda |
| `bool` | `*evaluator.Boolean` | `benar` / `salah` / `true` / `false` |
| `[]any` / slice | `*evaluator.Array` | Array dinamis Celang |
| `map[string]any` | `*evaluator.Map` | Map/Dictionary Celang |
| `nil` | `*evaluator.Null` | `kosong` / `nihil` / `null` |
| `error` | `*evaluator.Error` | Dibungkus ke dalam format **Empathetic Diagnostic** |

> **Prinsip Empati Celang pada Error OS:**  
> Jika berkas tidak ditemukan, alih-alih melempar error teknis mentah seperti `ENOENT: no such file or directory`, jembatan kita membungkusnya menjadi:  
> `Pesan: berkas './laporan.txt' tidak ditemukan`  
> `Solusi: Pastikan path berkas sudah benar dan berkas tersebut sudah dibuat di disk.`

### Lapis 3: Celang Module Namespace
Representasi yang dilihat oleh programmer Celang: modul dapat diakses langsung sebagai objek global bawaan (`berkas.baca(...)`) ATAU diimpor melalui modul virtual (`ambil { baca } dari "berkas"`).

---

## 3. Mekanisme Akses: Virtual Module Interception

Saat ini, modul resolver di [`pkg/evaluator/module.go`](file:///home/yosatsu/Coding/Celang/pkg/evaluator/module.go) mencari berkas di disk:
```celang
ambil { hitung } dari "./rumus.ce" // Mencari berkas fisik di disk
```

Untuk Standard Library, kita membuat **Virtual Module Registry (Pencegat Modul Virtual)**:
Saat evaluator membaca perintah impor:
```celang
ambil { baca, tulis } dari "berkas"
// atau
import { read, write } from "fs"
```

1. Module resolver memeriksa: *Apakah nama impor ada dalam daftar Standard Library (`"berkas"`, `"fs"`, `"matematika"`, `"math"`, `"teks"`, `"strings"`, `"waktu"`, `"time"`, `"json"`)?*
2. **Jika YA:** Evaluator **TIDAK** mencari berkas di harddisk. Evaluator langsung mengambil *pre-compiled environment* modul Go yang sudah didaftarkan di memori.
3. **Jika BUKAN:** Evaluator mencari berkas fisik seperti biasa.

Ini persis seperti cara kerja `import "fmt"` di Go, `import fs from 'fs'` di Node.js, atau `import math` di Python!

---

## 4. Konsep Dwibahasa (*Bilingual Identity*)

Mengikuti filosofi Syntax Styling Celang, seluruh modul dan fungsi Standard Library memiliki **identitas kembar** (Bahasa Indonesia & Bahasa Inggris) yang menunjuk ke fungsi memori yang persis sama:

```
           ┌───────────────────────────────────────┐
           │        Go Native Function:            │
           │       os.ReadFile(filePath)           │
           └──────────────────┬────────────────────┘
                              │
               ┌──────────────┴──────────────┐
               ▼                             ▼
       Dialek Indonesia:             Dialek JavaScript / English:
       `berkas.baca(path)`           `fs.read(path)`
```

### Matriks Pemetaan Modul & Fungsi:

#### 1. Modul Berkas (`berkas` $\leftrightarrow$ `fs`)
| Fungsi Indonesia | Fungsi English / JS | Keterangan |
|---|---|---|
| `berkas.baca(path)` | `fs.read(path)` | Membaca seluruh isi file teks |
| `berkas.tulis(path, data)` | `fs.write(path, data)` | Menimpa/menulis file baru |
| `berkas.tambah(path, data)`| `fs.append(path, data)` | Menambahkan teks ke akhir file |
| `berkas.ada(path)` | `fs.exists(path)` | Mengecek keberadaan file (`benar`/`salah`)|
| `berkas.hapus(path)` | `fs.remove(path)` | Menghapus file |

#### 2. Modul Matematika (`matematika` $\leftrightarrow$ `math`)
| Fungsi Indonesia | Fungsi English / JS | Keterangan |
|---|---|---|
| `matematika.akar(x)` | `math.sqrt(x)` | Akar kuadrat ($\sqrt{x}$) |
| `matematika.pangkat(a, b)`| `math.pow(a, b)` | Perpangkatan ($a^b$) |
| `matematika.mutlak(x)` | `math.abs(x)` | Nilai absolut ($|x|$) |
| `matematika.bulat(x)` | `math.round(x)` | Pembulatan terdekat |
| `matematika.lantai(x)` | `math.floor(x)` | Pembulatan ke bawah |
| `matematika.atap(x)` | `math.ceil(x)` | Pembulatan ke atas |
| `matematika.acak(min, max)`| `math.random(min, max)`| Bilangan acak integer |
| `matematika.PI` | `math.PI` | Konstanta $\pi$ ($3.14159...$) |

#### 3. Modul Teks (`teks` $\leftrightarrow$ `strings`)
| Fungsi Indonesia | Fungsi English / JS | Keterangan |
|---|---|---|
| `teks.pisah(str, pemisah)` | `strings.split(str, sep)` | Memecah teks ke Array |
| `teks.gabung(arr, pemisah)`| `strings.join(arr, sep)` | Menggabungkan Array ke teks |
| `teks.kapital(str)` | `strings.upper(str)` | Huruf besar semua |
| `teks.kecil(str)` | `strings.lower(str)` | Huruf kecil semua |
| `teks.potong(str)` | `strings.trim(str)` | Memotong spasi di ujung |
| `teks.ganti(str, pola, baru)`| `strings.replace(str, old, new)` | Mengganti bagian teks |
| `teks.berisi(str, cari)` | `strings.contains(str, substr)` | Mengecek substring |

#### 4. Modul Waktu (`waktu` $\leftrightarrow$ `time`)
| Fungsi Indonesia | Fungsi English / JS | Keterangan |
|---|---|---|
| `waktu.sekarang()` | `time.now()` | Timestamp milidetik sekarang |
| `waktu.tidur(ms)` | `time.sleep(ms)` | Menunda eksekusi sejumlah milidetik |

#### 5. Modul JSON (`json` $\leftrightarrow$ `json`)
| Fungsi Indonesia | Fungsi English / JS | Keterangan |
|---|---|---|
| `json.uraikan(str)` | `json.parse(str)` | Mengubah string JSON ke Map/Array |
| `json.format(obj)` | `json.stringify(obj)` | Mengubah Map/Array ke string JSON |

---

## 5. Rencana Langkah Implementasi Teknis (Tahapan Eksekusi)

Berikut rencana kerja bertahap yang akan kita eksekusi di branch `feature/stdlib`:

1. **Langkah 1: Pembuatan Paket `pkg/stdlib`**:
   - Membuat file registrasi modul standar dan tipe helper untuk registrasi fungsi host Go.
2. **Langkah 2: Integrasi ke Module Loader & Evaluator**:
   - Memodifikasi `loadModule()` di `pkg/evaluator/module.go` agar mengintersep nama modul standar bawaan.
   - Mendaftarkan objek modul ke global environment saat evaluasi dimulai sehingga bisa diakses langsung (misal `berkas.baca`) tanpa wajib `ambil` jika diinginkan.
3. **Langkah 3: Implementasi Modul per Modul**:
   - `berkas.go`: Menggunakan `os` dan `path/filepath`.
   - `matematika.go`: Menggunakan `math` dan `math/rand`.
   - `teks.go`: Menggunakan `strings`.
   - `waktu.go`: Menggunakan `time`.
   - `json.go`: Menggunakan `encoding/json` dengan konversi rekursif ke `evaluator.Map` dan `evaluator.Array`.
4. **Langkah 4: Unit Testing & Verifikasi**:
   - Menulis unit test untuk setiap modul di `pkg/stdlib/stdlib_test.go`.
   - Membuat berkas contoh program di `examples/demo_stdlib.ce`.
   - Menjalankan pengujian CLI dan memastikan 100% lulus.
