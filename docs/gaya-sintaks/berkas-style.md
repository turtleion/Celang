# Membuat Gaya Bahasa Sendiri (Berkas .style)

Celang tidak membatasi Anda hanya pada gaya bahasa bawaan. Anda dapat membuat gaya bahasa atau dialek Anda sendiri—baik itu bahasa daerah Nusantara (seperti Basa Sunda, Jawa, Minang, Batak, Bali, Bugis), bahasa internasional lainnya, maupun gaya santai tim internal Anda—menggunakan berkas konfigurasi berekstensi `.style`.

Berkas `.style` adalah dokumen konfigurasi berformat JSON yang memetakan kata kunci khusus Anda ke struktur abstrak bahasa Celang (*Universal Abstract Syntax Tree* / U-AST).

---

## 1. Template Siap Pakai (`template.style`)

Untuk mempermudah pembuatan dialek baru, Celang menyediakan berkas template resmi di repositori: [examples/template.style](file:///home/yosatsu/Coding/Celang/examples/template.style).

Anda dapat langsung menduplikasinya ke proyek Anda:

```bash
cp examples/template.style dialek_saya.style
```

Atau salin langsung template JSON di bawah ini (seluruh nilai awal telah diisi dengan kata kunci standar Bahasa Indonesia):

```json
{
  "metadata": {
    "name": "nama_dialek",
    "display_name": "Nama Tampilan Dialek",
    "version": "1.0.0",
    "author": "Nama Anda / Komunitas Anda",
    "description": "Deskripsi singkat mengenai dialek atau gaya bahasa ini"
  },
  "deklarasi": {
    "var": "buat",
    "const": "konstanta",
    "type_assign": "itu",
    "struct": "bentuk",
    "ref": "rujuk",
    "as": "sebagai"
  },
  "percabangan": {
    "if": "jika",
    "else_if": "selain jika",
    "else": "selain",
    "then": "maka",
    "switch": "pilih",
    "case": "bila",
    "default": "lainnya"
  },
  "perulangan": {
    "while": "selama",
    "for": "untuk",
    "in": "dalam"
  },
  "fungsi": {
    "def": "fungsi",
    "return": "kembalikan",
    "return_type": "mengembalikan",
    "defer": "saat_selesai"
  },
  "penanganan_galat": {
    "try": "coba",
    "catch": "tangkap",
    "finally": "akhirnya",
    "throw": "lempar"
  },
  "operator_logika": {
    "and": "dan",
    "or": "atau",
    "not": "bukan",
    "null_coalesce": "kalau kosong",
    "maybe": "mungkin"
  },
  "literal": {
    "true": "benar",
    "false": "salah",
    "null": "kosong"
  },
  "konkurensi": {
    "async": "sambil",
    "parallel": "bersamaan"
  },
  "format": {
    "block_style": "brace",
    "indent_width": 4,
    "quote_style": "double"
  }
}
```

---

## 2. Contoh Nyata: Dialek Basa Sunda (`sunda.style`)

Berikut adalah contoh implementasi nyata dialek Basa Sunda dari berkas [examples/sunda.style](file:///home/yosatsu/Coding/Celang/examples/sunda.style):

```json
{
  "metadata": {
    "name": "sunda",
    "display_name": "Basa Sunda",
    "version": "1.0.0",
    "author": "Komunitas Celang",
    "description": "Dialek basa Sunda pikeun basa pamrograman Celang"
  },
  "deklarasi": {
    "var": "jieun",
    "const": "tetep",
    "type_assign": "nyaeta",
    "struct": "wangun",
    "ref": "rujuk",
    "as": "salaku"
  },
  "percabangan": {
    "if": "lamun",
    "else_if": "atawa lamun",
    "else": "lamun henteu",
    "then": "maka",
    "switch": "cocogkeun",
    "case": "bisi",
    "default": "lianna"
  },
  "perulangan": {
    "while": "salami",
    "for": "pikeun",
    "in": "dina"
  },
  "fungsi": {
    "def": "fungsi",
    "return": "balikkeun",
    "return_type": "ngahasilkeun",
    "defer": "pas_rengse"
  },
  "penanganan_galat": {
    "try": "cobi",
    "catch": "tewak",
    "finally": "pamungkas",
    "throw": "alungkeun"
  },
  "operator_logika": {
    "and": "sareng",
    "or": "atawa",
    "not": "sanes",
    "null_coalesce": "mun_kosong",
    "maybe": "meureun"
  },
  "literal": {
    "true": "bener",
    "false": "salah",
    "null": "kosong"
  },
  "konkurensi": {
    "async": "bari",
    "parallel": "babarengan"
  },
  "format": {
    "block_style": "brace",
    "indent_width": 4,
    "quote_style": "double"
  }
}
```

---

## 3. Tabel Referensi Kata Kunci Lengkap

Setiap kunci dalam berkas `.style` mewakili elemen sintaksis bahasa Celang:

| Kategori JSON | Kunci JSON | Nilai Bawaan Celang | Contoh Basa Sunda | Keterangan & Penggunaan |
|---|---|---|---|---|
| **deklarasi** | `var` | `buat` | `jieun` | Deklarasi variabel yang dapat diubah nilainya |
| | `const` | `konstanta` | `tetep` | Deklarasi nilai konstan yang tidak dapat diubah |
| | `type_assign` | `itu` | `nyaeta` | Anotasi tipe data statis (`x itu Angka`) |
| | `struct` | `bentuk` | `wangun` | Pendefinisian struktur data objek / tipe bentukan |
| | `ref` | `rujuk` | `rujuk` | Kata kunci referensi / pointer |
| | `as` | `sebagai` | `salaku` | Alias impor atau casting tipe |
| **percabangan** | `if` | `jika` | `lamun` | Kondisi percabangan utama |
| | `else_if` | `selain jika` | `atawa lamun` | Cabang alternatif bersyarat |
| | `else` | `selain` | `lamun henteu` | Cabang alternatif terakhir |
| | `then` | `maka` | `maka` | Kata kunci pembuka aksi kondisi |
| | `switch` | `pilih` | `cocogkeun` | Pemilihan multi-kasus |
| | `case` | `bila` | `bisi` | Klausa pencocokan kasus |
| | `default` | `lainnya` | `lianna` | Kasus default bila tidak ada yang cocok |
| **perulangan** | `while` | `selama` | `salami` | Perulangan berbasis kondisi benar |
| | `for` | `untuk` | `pikeun` | Perulangan iterasi koleksi atau rentang |
| | `in` | `dalam` | `dina` | Operator keanggotaan koleksi perulangan |
| **fungsi** | `def` / `func` | `fungsi` | `fungsi` | Deklarasi fungsi / prosedur |
| | `return` | `kembalikan` | `balikkeun` | Mengembalikan nilai hasil pemrosesan fungsi |
| | `return_type` | `mengembalikan` | `ngahasilkeun` | Anotasi tipe data nilai kembalian fungsi |
| | `defer` | `saat_selesai` | `pas_rengse` | Eksekusi tertunda di akhir cakupan blok |
| **penanganan_galat** | `try` | `coba` | `cobi` | Memulai blok penanganan pengecualian |
| | `catch` | `tangkap` | `tewak` | Menangkap galat yang terjadi dalam blok coba |
| | `finally` | `akhirnya` | `pamungkas` | Blok pembersihan yang selalu dijalankan |
| | `throw` | `lempar` | `alungkeun` | Memicu lemparan galat secara eksplisit |
| **operator_logika** | `and` | `dan` | `sareng` | Operasi logika konjungsi (AND) |
| | `or` | `atau` | `atawa` | Operasi logika disjungsi (OR) |
| | `not` | `bukan` | `sanes` | Operasi logika negasi (NOT) |
| | `null_coalesce` | `kalau kosong` | `mun_kosong` | Operator fallback nilai kosong (`a kalau kosong b`) |
| | `maybe` | `mungkin` | `meureun` | Operator chaining aman (safe navigation) |
| **literal** | `true` | `benar` | `bener` | Nilai kebenaran positif |
| | `false` | `salah` | `salah` | Nilai kebenaran negatif |
| | `null` | `kosong` | `kosong` | Nilai ketiadaan objek / null pointer |
| **konkurensi** | `async` | `sambil` | `bari` | Menjalankan blok tugas secara asinkron |
| | `parallel` | `bersamaan` | `babarengan` | Menjalankan kumpulan tugas secara paralel |
| **format** | `block_style` | `"brace"` | `"brace"` | Gaya pembungkus blok (`brace`, `python`, `ruby`) |
| | `indent_width` | `4` | `4` | Lebar indentasi spasi saat kode ditata ulang |
| | `quote_style` | `"double"` | `"double"` | Gaya tanda kutip teks (`double` atau `single`) |

---

## 4. Format Datar (*Flat JSON*)

Selain format berjenjang (*categorized JSON*), Celang juga mendukung **Format Datar (*Flat JSON*)** jika Anda ingin menulis konfigurasi yang lebih ringkas atau hanya ingin menimpa sebagian kecil kata kunci:

```json
{
  "name": "dialek_ringkas",
  "display_name": "Dialek Ringkas",
  "var_decl": "gawe",
  "const_decl": "tetep",
  "if": "yen",
  "else": "liya",
  "while": "saksuwene",
  "func": "fungsi",
  "return": "balekno",
  "true": "bener",
  "false": "kliru",
  "null": "suwung",
  "block_style": "brace",
  "indent_width": 2
}
```

> [!NOTE]
> Celang secara otomatis mendeteksi apakah berkas `.style` Anda menggunakan format kategorikal atau format datar. Keduanya valid dan didukung penuh.

---

## 5. Pengaturan Gaya Blok (`block_style`)

Properti `block_style` pada bagian `format` menentukan bagaimana struktur blok kode akan diformat ulang:

- `"brace"`: Menggunakan kurung kurawal `{ ... }` ala C, JavaScript, Go, dan Rust (Gaya bawaan Celang).
- `"python"`: Menggunakan titik dua `:` dan indentasi spasi murni tanpa kurung kurawal.
- `"ruby"`: Menggunakan kata kunci penutup `end`.

---

## 6. Cara Menggunakan Berkas `.style` Anda

Setelah Anda menyimpan berkas gaya Anda (misalnya `jawa.style` atau `dialek_saya.style`):

### A. Membuka Jendela Interaktif (REPL) dengan Gaya Anda
```bash
./celang repl --style=dialek_saya.style
```
Jendela interaktif akan menyambut Anda langsung dengan aturan kata kunci gaya baru tersebut.

### B. Menerjemahkan dan Menata Ulang Kode
Anda dapat menata ulang berkas kode yang ada menjadi gaya baru Anda:

```bash
# Menampilkan hasil terjemahan di layar terminal:
./celang format program.ce --style=dialek_saya.style

# Langsung menimpa dan mengubah berkas kode secara permanen:
./celang format -w program.ce --style=dialek_saya.style
```

### C. Menjalankan Kode
Kode yang sudah ditulis menggunakan dialek tersebut dapat langsung dieksekusi:
```bash
./celang run program.ce
```
Celang secara cerdas mengenali kata kunci yang didaftarkan melalui berkas konfigurasi maupun direktif in-code styling.
