# Rancangan Desain & Spesifikasi Celang CLI & REPL (`CLI_REPL.md`)

> **Catatan untuk Kamu:**  
> Dokumen ini adalah rancangan arsitektur dan spesifikasi teknis untuk membuat antarmuka baris perintah (**Celang CLI**) dan terminal interaktif (**Celang REPL**) menjadi alat pengembang yang profesional, intuitif, dan nyaman digunakan.  
> Silakan review poin-poin di bawah ini. Jika ada nama perintah, flag, atau perilaku REPL yang ingin diubah, kabari aku!

---

## 1. 🎯 Visi & Filosofi Tooling Celang

1. **Human-Centric & Welcoming**: Pesan terminal tidak dingin atau kaku; menggunakan pewarnaan yang elegan (*ANSI colors/emojis*), pesan bantuan yang mudah dipahami, dan diagnostik error yang empatik.
2. **Batteries Included**: Satu biner `celang` sudah memuat compiler, interpreter, in-place code formatter, terminal REPL interaktif, dan project initializer tanpa perlu toolchain terpisah.
3. **Syntax-Agnostic First**: REPL dan CLI mendukung penuh konversi dialek instan (*live style preview*).

---

## 2. 📋 Struktur & Daftar Perintah CLI (`celang <command>`)

Format umum:
```bash
celang [perintah] [opsi/flag] [argumen...]
```

| Perintah | Contoh Penggunaan | Penjelasan |
| :--- | :--- | :--- |
| *(tanpa argumen)* | `celang` | Otomatis masuk ke **Interactive REPL**. |
| `repl` | `celang repl` | Membuka sesi Interactive REPL secara eksplisit. |
| `run` | `celang run app.ce` | Menjalankan file kode sumber (`.ce`) atau Universal AST (`.ast`). |
| `run` *(dg argumen)* | `celang run app.ce -- foo bar` | Menjalankan skrip dan meneruskan argumen ke variabel global `argumen`. |
| `format` | `celang format app.ce --style=jawa.style` | Memformat kode menggunakan gaya bawaan atau file style lokal kustom (`.style`). |
| `format -w` | `celang format -w app.ce --style=jawa.style` | **In-place formatting**: Langsung menimpa berkas sumber di disk dengan gaya yang dipilih. |
| `compile` / `parse` | `celang compile app.ce -o app.ast` | Mengubah kode sumber teks menjadi berkas U-AST JSON. |
| `ast` | `celang ast app.ce` | Menginspeksi struktur pohon U-AST secara langsung di terminal dengan pewarnaan. |
| `init` | `celang init [nama_proyek]` | Menginisialisasi proyek modular baru (`celang.mod` + `utama.ce`). |
| `versi` / `--version` / `-v` | `celang --version` | Menampilkan nomor versi runtime, OS/Arch, dan ASCII logo. |
| `bantuan` / `--help` / `-h` | `celang help` | Menampilkan daftar lengkap petunjuk perintah CLI. |

---

## 3. 💻 Interactive REPL (Read-Eval-Print-Loop)

REPL memungkinkan programmer mencoba ekspresi, menguji algoritma, atau mendemonstrasikan sintaks Celang secara langsung di terminal tanpa harus membuat berkas `.ce`.

### 3.1 Tampilan Sambutan & Prompt Terminal

Saat pengguna menjalankan perintah `celang`:

```text
🌟 Celang v1.0.0 (Linux/amd64) - Human-Centric & Syntax-Styling Language
Ketik kode Celang langsung, atau ':bantuan' untuk daftar perintah khusus.
Tekan Ctrl+C atau ketik ':keluar' untuk mengakhiri sesi.

celang> buat nama = "Kayasa"
celang> "Halo " + nama + "!"
=> "Halo Kayasa!"

celang> [1, 2, 3] |> petakan(fungsi(x) { kembalikan x * 10 })
=> [10, 20, 30]
```

### 3.2 Auto-Printing Hasil Evaluasi
Di dalam REPL, programmer **tidak wajib** mengetik `cetak(...)` untuk melihat hasil:
- Setiap baris ekspresi yang menghasilkan nilai (bukan penugasan/statement kosong) akan otomatis dicetak dengan prefix tanda panah `=> <hasil>`.
- Jika baris berupa statement deklarasi (`buat x = 10`), REPL tetap tenang tanpa output bising, siap menerima baris berikutnya.

### 3.3 Penanganan Blok Multi-Baris (*Multi-line Input*)
Jika pengguna membuka kurung kurawal `{`, kurung siku `[`, atau tanda kurung `(` yang belum ditutup, REPL otomatis berganti prompt menjadi titik-titik kelanjutan (`......> `) hingga seluruh blok selesai:

```text
celang> fungsi hitung_diskon(total) {
......>     jika total > 50000 {
......>         kembalikan total * 0.9
......>     } selain {
......>         kembalikan total
......>     }
......> }
celang> hitung_diskon(100000)
=> 90000
```

### 3.4 Perintah Khusus REPL (Meta-Commands dengan awalan `:`)

| Perintah Meta | Alias | Deskripsi |
| :--- | :--- | :--- |
| `:keluar` | `:exit`, `:quit`, `:q` | Keluar dari sesi REPL. |
| `:bersih` | `:clear`, `:cls` | Membersihkan layar terminal. |
| `:lingkungan` | `:env`, `:vars` | Menampilkan seluruh variabel, konstanta, dan bentuk yang aktif di memori sesi. |
| `:gaya <nama|path>` | `:style <nama|path>` | Mengganti profil gaya aktif (*live styling*) di REPL. Bisa nama bawaan (`indonesian`, `javascript`, `pythonic`) maupun berkas lokal (`./jawa.style`, `examples/sunda.style`). |
| `:ast <kode>` | `:uast` | Menampilkan representasi U-AST JSON dari baris kode yang diketik. |
| `:reset` | `:ulang` | Mengosongkan seluruh memori environment kembali ke sesi awal. |
| `:bantuan` | `:help`, `:?` | Menampilkan panduan bantuan perintah khusus REPL. |

---

## 4. 📥 Fitur Keyboard Input: `masukan(...)` / `input(...)`

Fungsi bawaan ini melengkapi ekosistem CLI sehingga programmer Celang dapat membuat program interaktif berbasis teks:

### Sintaks & Contoh:
```celang
buat nama = masukan("Masukkan nama Anda: ")
buat umur = masukan("Masukkan umur Anda: ")

cetak("Halo {nama}, tahun depan umurmu adalah {umur + 1} tahun!")
```

### Karakteristik:
1. Menerima 1 argumen opsional berupa teks prompt (*string*).
2. Menampilkan prompt tanpa baris baru (*no newline*), lalu menunggu pengguna mengetik hingga tombol Enter ditekan.
3. Mengembalikan nilai bertipe `Teks` (`String`) dengan karakter newline (`\r\n`) yang otomatis dibersihkan (*trimmed*).

---

## 5. 🚩 Penerusan Argumen Baris Perintah (`argumen` / `args`)

Saat membuat script otomatisasi atau tool CLI dengan Celang, pengguna dapat membaca argumen terminal yang dikirimkan:

### Cara Menjalankan:
```bash
celang run kalkulator.ce -- tambah 100 50
```

### Cara Membaca di Celang:
Pemisah tanda `--` menandai dimulainya argumen skrip. Di dalam kode Celang, runtime menginjeksi array bawaan `argumen` (alias `args`):

```celang
// File: kalkulator.ce
jika panjang(argumen) < 3 {
    cetak("Penggunaan: celang run kalkulator.ce -- <operasi> <a b>")
    kembalikan kosong
}

buat operasi = argumen[0]
buat a = argumen[1]
buat b = argumen[2]

jika operasi == "tambah" {
    cetak("Hasil: {a + b}")
}
```

---

## 6. 🛠️ In-Place Code Formatter (`celang format -w`)

Melengkapi perintah `celang format` yang sebelumnya hanya menampilkan hasil ke stdout:

```bash
# 1. Menampilkan ke terminal (dry-run):
celang format program.ce --style=javascript

# 2. In-place write (langsung menimpa berkas sumber):
celang format -w program.ce

# 3. In-place write dengan transformasi gaya bahasa bawaan:
celang format -w program.ce --style=pythonic

# 4. In-place write dengan file style lokal kustom (.style):
celang format -w program.ce --style=./jawa.style
```

---

## 7. 📦 Project Initializer (`celang init`)

Membuat boilerplate struktur proyek Celang yang rapi dan siap dijalankan:

### Penggunaan:
```bash
celang init projek_baru
```

### Berkas yang Dihasilkan:
```text
projek_baru/
├── celang.mod          # Penanda modul root proyek
├── utama.ce            # Berkas entrypoint utama
└── README.md           # Dokumentasi awal proyek
```

Isi `celang.mod`:
```ini
modul projek_baru
versi 1.0.0
```

Isi `utama.ce`:
```celang
// Entrypoint utama aplikasi Celang
fungsi utama() {
    cetak("Halo dari Celang! Selamat datang di projek_baru 🚀")
}

utama()
```

---

## 8. 🔍 Pemeriksa Pohon U-AST (`celang ast`)

Alat diagnostik bagi programmer dan pembuat compiler untuk melihat representasi pohon U-AST dari suatu berkas:

```bash
celang ast program.ce
```
Output mencetak U-AST berformat JSON dengan indentasi rapi dan sintaks yang terbaca.

---

## 9. 💬 Pertanyaan & Review Sebelum Eksekusi

Silakan beri masukan untuk hal-hal berikut:
1. **Nama Perintah REPL:** Apakah kamu lebih suka REPL berjalan otomatis saat mengetik `celang` tanpa argumen, atau harus selalu mengetik `celang repl`? (Rancanganku: dua-duanya bisa).
2. **Karakter Prompt REPL:** Apakah `celang> ` dan kelanjutan `......> ` sudah nyaman di mata?
3. **Nama Variabel Argumen CLI:** Apakah `argumen` dan `args` sudah pas?
4. **Nama Fungsi Input Keyboard:** Apakah pasangan `masukan(pesan)` dan `input(pesan)` sudah ideal?
