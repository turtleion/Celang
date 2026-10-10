# Panduan Lengkap: Syntax Styling & Ekosistem Celang (`STYLING-AND-MORE.md`)

Dokumen ini mengupas tuntas arsitektur **Syntax Styling (Sintaks Dinamis)** pada Celang, mulai dari cara berpindah gaya di luar kode (CLI & Editor/Git), di dalam kode (Pragmas & Scope Styling), hingga eksplorasi mendalam mengenai sistem keyword dan pembuatan profil gaya kustom (*Syntax Definition Language*).

---

## 1. Filosofi Inti: Pemisahan Logika (Data) dan Tampilan (Presentasi)

Dalam bahasa pemrograman konvensional (Go, Python, C), **teks kode dan logika program menyatu secara kaku**. Jika tim berdebat antara gaya kurung kurawal vs indentasi, atau keyword bahasa Inggris vs Indonesia, salah satu pihak harus mengalah.

Celang memisahkan keduanya seperti arsitektur web modern:
- **Universal AST (U-AST)** = **HTML / Data Model** (Kebenaran tunggal yang disimpan di disk dan Git).
- **Style Profile (.style)** = **CSS** (Aturan visual bagaimana node AST dirender menjadi teks yang dibaca manusia).

```
                      ┌────────────────────────┐
                      │   Universal AST (Data) │
                      │      (file.ast)        │
                      └───────────┬────────────┘
                                  │
         ┌────────────────────────┼────────────────────────┐
         ▼                        ▼                        ▼
┌──────────────────┐    ┌──────────────────┐    ┌──────────────────┐
│ Profil Indonesia │    │  Profil Pythonic │    │  Profil Classic  │
│      (.ce)       │    │      (.pyce)     │    │      (.jsce)     │
└──────────────────┘    └──────────────────┘    └──────────────────┘
```

---

## 2. Mengubah Style di Luar Kode (Eksternal / Tooling)

Pengubahan style di luar kode terjadi tanpa mengubah satu baris pun logika di dalam U-AST.

### 2.1 Melalui CLI (Command Line Interface)

CLI Celang menyediakan fitur format dan restyling instan:

```bash
# 1. Menampilkan berkas dengan gaya Bahasa Indonesia
celang format program.ce --style=indonesian

# 2. Merender berkas yang sama menjadi gaya Pythonic
celang format program.ce --style=pythonic

# 3. Merender berkas menjadi gaya C / JavaScript
celang format program.ce --style=classic

# 4. Mengubah langsung file sumber ke style baru
celang format program.ce --style=pythonic -o program_pythonic.ce
```

### 2.2 Konfigurasi Profil Global & Lokal (`celang.style`)

Pengembang dapat menyetel preferensi gaya di tingkat proyek atau pengguna:

1. **File Konfigurasi Proyek (`.celangrc` atau `celang.style.json`)**:
   Ditaruh di root repositori proyek:
   ```json
   {
     "default_style": "indonesian",
     "tab_width": 4,
     "use_parentheses_in_if": false,
     "map_delimiter": "itu"
   }
   ```
2. **Global Config User (`~/.celang/config.json`)**:
   Setiap programmer di tim yang sama bisa memiliki preferensi lokal yang berbeda di komputernya masing-masing.

### 2.3 Git Integration: Custom Diff Driver & Clean/Smudge Filter

Karena file kanonikal di Git adalah U-AST (`.ast`), tim programmer dapat berkolaborasi tanpa pernah mengalami *merge conflict* akibat selera penulisan sintaks:

1. **Smudge Filter (Saat Git Checkout):**
   U-AST otomatis dirender ke teks dengan style favorit programmer lokal (misal: Programmer A melihat Pythonic, Programmer B melihat Gaya Indonesia).
2. **Clean Filter (Saat Git Commit):**
   Teks lokal dikonversi kembali ke format U-AST kanonikal sebelum masuk ke commit Git.
3. **Custom Git Diff Driver:**
   Saat menjalankan `git diff`, Celang menampilkan perubahan dalam bentuk sintaks teks yang manusiawi, bukan perbedaan JSON AST mentah.

### 2.4 Language Server Protocol (LSP) & Editor

Pada VS Code atau editor lain:
- LSP Celang membaca file `.ast` atau `.ce`.
- LSP merender virtual document / dekorator sintaks di layar editor sesuai preferensi pengguna saat itu juga.
- Programmer A mengetik `buat nama = "Budi"`, tetapi programmer B yang sedang *live pair-programming* melihatnya sebagai `nama = "Budi"`.

---

## 3. Mengubah Style di Dalam Kode (In-Code Styling)

Selain pengaturan eksternal, Celang mendukung instruksi gaya langsung di dalam berkas kode.

### 3.1 Berkas Header Pragma (`# @gaya:` atau `// @style:`)

Diletakkan di baris pertama sebuah file untuk menentukan gaya default file tersebut:

```celang
// @gaya: indonesian
buat pesan = "Halo dari Celang"
cetak(pesan)
```

Atau dalam file bergaya Pythonic:
```celang
# @style: pythonic
pesan = "Halo dari Celang"
print(pesan)
```

### 3.2 Block-Level Styling (`gaya "..." { }`)

Fitur inovatif Celang di mana dalam **satu program yang sama**, blok-blok tertentu dapat ditulis menggunakan dialek yang berbeda jika dirasa lebih ekspresif:

```celang
// Bagian Utama: Logika Bisnis (Gaya Indonesia)
buat saldo = 1000000

// Blok Khusus: Rumus Matematis / Data Science (Gaya Pythonic)
gaya "pythonic" {
    def hitung_proyeksi(pokok, bunga, tahun):
        total = pokok
        for t in rentang(1, tahun):
            total = total + (total * bunga / 100)
        return total
    
    saldo_akhir = hitung_proyeksi(saldo, 5, 3)
}

cetak("Saldo Proyeksi:", saldo_akhir)
```

> **Bagaimana Mesin Menanganinya?**  
> Saat Lexer menemui blok `gaya "<nama>"`, parser sementara mengalihkan kamus keyword dan aturan penutup blok ke profil yang diminta, lalu menyatukan node-node hasilnya ke dalam satu pohon U-AST yang sama!

---

## 4. Pembahasan Mendalam: Sistem Keyword Celang

Celang mengusung konsep **Kamus Dwibahasa / Poliglot Terbuka**.

### 4.1 Matriks Pemetaan Keyword Antar Profil

| Konsep Semantik | Gaya Default (Indonesia) | Gaya Pythonic | Gaya Classic (C/JS) | Catatan Khusus |
|---|---|---|---|---|
| **Deklarasi Variabel** | `buat` | *(implisit)* | `let` / `var` | Menandai variabel mutable |
| **Konstanta** | `konstanta` | `const` | `const` | Immutable |
| **Penugasan Tipe** | `itu` / `:` | `:` | `:` | Anotasi tipe data |
| **Percabangan** | `jika` | `if` | `if` | Kurung `( )` opsional |
| **Alternatif Jika** | `selain jika` | `elif` | `else if` | Dapat berulang |
| **Alternatif Akhir** | `selain` | `else` | `else` | Blok penutup kondisi |
| **Kondisi Selesai** | `selesai` / `}` | `}` | `}` | Mode kurung atau keyword |
| **Perulangan While** | `selama` | `while` | `while` | Perulangan berdasar kondisi |
| **Perulangan Iterasi** | `untuk ... dalam` | `for ... in` | `for (let item of ...)` | 100% gaya JavaScript |
| **Pencocokan Nilai** | `cocok ... kasus ... lainnya` | `match ... case` | `switch ... case ... default` | Pattern / Switch matching |
| **Deklarasi Fungsi** | `fungsi` | `def` | `function` | 100% JavaScript keyword |
| **Nilai Kembali** | `kembalikan` | `return` | `return` | Return statement |
| **Panah Tipe Return** | `mengembalikan` / `->` | `->` | `->` | Anotasi tipe keluaran |
| **Pointer & Alamat** | `&` (alamat), `*` (deref) | `&` / `*` | `&` / `*` | Manipulasi referensi memori |
| **Asinkron Latar** | `sambil` | `async` | `async` | Goroutine mandiri |
| **Paralel Sinkron** | `bersamaan` | `parallel` | `parallel` | WaitGroup otomatis |
| **Nilai Benar** | `benar` | `True` | `true` | Boolean truthy |
| **Nilai Salah** | `salah` | `False` | `false` | Boolean falsy |
| **Ketiadaan Nilai** | `kosong` / `nihil` | `None` | `null` | Null literal |
| **Logika DAN** | `dan` / `&&` | `and` / `&&` | `&&` | Operator boolean |
| **Logika ATAU** | `atau` / `\|\|` | `or` / `\|\|` | `\|\|` | Operator boolean |
| **Logika BUKAN** | `bukan` / `tidak` / `!` | `not` / `!` | `!` | Negasi boolean |

### 4.2 Toleransi Dwibahasa (*Bilingual Tolerance*)

Lexer Celang dirancang **ramah manusia**:
- Di dalam kode bergaya Indonesia, programmer yang tidak sengaja mengetik `return` alih-alih `kembalikan`, atau `if` alih-alih `jika`, **tetap dikenali secara valid**.
- Di dalam pemisah Map: programmer boleh menulis `{"nama": "Budi"}` maupun `{"nama" itu "Budi"}`. Keduanya menghasilkan node `MapLiteral` yang persis sama di U-AST.

---

## 5. Merancang Gaya Kustom: Syntax Definition Language (SDL)

Pengembang dapat membuat profil gaya mereka sendiri menggunakan format definisi sintaks `.style`:

### Contoh Berkas: `sunda.style`
```json
{
  "name": "sunda",
  "display_name": "Basa Sunda",
  "var_decl": "jieun",
  "const_decl": "tetep",
  "type_assign": "nyaeta",
  "if": "lamun",
  "else_if": "atawa lamun",
  "else": "lamun henteu",
  "while": "salami",
  "for": "pikeun",
  "in": "dina",
  "func": "fungsi",
  "return": "balikkeun",
  "async": "bari",
  "parallel": "babarengan",
  "true": "bener",
  "false": "salah",
  "null": "kosong",
  "and": "sareng",
  "or": "atawa",
  "not": "sanes",
  "return_arrow": "ngahasilkeun"
}
```

Dengan mendaftarkan file ini, kode program dapat langsung diformat atau ditulis dalam dialek Sunda:
```celang
jieun ngaran = "Asep"
lamun (ngaran == "Asep") {
    balikkeun "Sampurasun!"
}
```

---

## 6. Fitur Canggih yang Telah Terimplementasi

Fitur-fitur inovatif yang telah siap pakai:

1. **Fitur Pointer & Referensi (`&` dan `*`)**:
   - Untuk efisiensi dan mutasi memori langsung: `&var` untuk mengambil alamat memori variabel, dan `*ptr = nilai_baru` untuk dereferensi dan mutasi langsung dalam environment.
2. **Pattern Matching / Switch (`cocok ... kasus ... lainnya`)**:
   - Percabangan multi-kasus yang bersih. Terformat otomatis menjadi `switch ... case ... default` di JavaScript dan `match ... case` di Pythonic style.
3. **Empathetic Error Handler (Pesan Kesalahan Humanis)**:
   - Jika pengguna salah mengetik sintaks atau tipe data, compiler menyajikan visual snippet baris kode dengan penunjuk kolom (`^`), serta saran (*Tips*) solusi perbaikan ramah manusia.
4. **Custom SDL (.style) Loader Dinamis**:
   - Memungkinkan pemuatan profil dialek baru (misal: `examples/sunda.style`) secara dinamis melalui flag `--style=path/to.style`.
5. **Gaya Classic 100% JavaScript**:
   - Menggunakan `function`, `let`, `for (let item of array)`, `switch`, dan format idiomatik JavaScript murni.

---

## 7. Rangkuman Alur Kerja Sintaks Celang

```
 Programmer Mengetik Kode (.ce)
              │
              ▼
   [ Lexer & Parser Celang ]
              │
              ▼
   Universal AST (.ast / JSON)  ◄── Disimpan ke Git / Disk
              │
      ┌───────┴───────┐
      ▼               ▼
 [ Evaluator ]   [ Formatter / Reverse Parser ]
 (Dijalankan)    (Dirender ke Style Indonesia / Python / C / Custom)
```

Dengan sistem ini, **Celang adalah bahasa pertama yang tidak membatasi programmer pada satu dialek teks**, melainkan membebaskan setiap orang untuk mengekspresikan logika mereka dengan sintaks yang paling nyaman bagi mereka!
