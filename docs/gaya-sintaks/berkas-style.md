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

## 2. Contoh Nyata: Dialek Basa Jawa (`jawa.style`)

Berikut adalah contoh implementasi nyata dialek Basa Jawa dari berkas [examples/jawa.style](file:///home/yosatsu/Coding/Celang/examples/jawa.style):

```json
{
  "metadata": {
    "name": "jawa",
    "display_name": "Basa Jawa",
    "version": "1.0.0",
    "author": "Komunitas Celang",
    "description": "Dialek basa Jawa kanggo basa pamrograman Celang"
  },
  "deklarasi": {
    "var": "gawe",
    "const": "tetep",
    "type_assign": "yaiku",
    "struct": "wujud",
    "ref": "rujuk",
    "as": "minangka"
  },
  "percabangan": {
    "if": "yen",
    "else_if": "utawa yen",
    "else": "yen ora",
    "then": "mula",
    "switch": "cocogna",
    "case": "nalika",
    "default": "liyane"
  },
  "perulangan": {
    "while": "sasuwene",
    "for": "kanggo",
    "in": "ing"
  },
  "fungsi": {
    "def": "fungsi",
    "return": "balekna",
    "return_type": "ngasilakeun",
    "defer": "pas_bubar"
  },
  "penanganan_galat": {
    "try": "jajal",
    "catch": "cekel",
    "finally": "pungkasan",
    "throw": "uncalna"
  },
  "operator_logika": {
    "and": "lan",
    "or": "utawa",
    "not": "dudu",
    "null_coalesce": "yen_suwung",
    "maybe": "mbokmenawa"
  },
  "literal": {
    "true": "bener",
    "false": "salah",
    "null": "suwung"
  },
  "konkurensi": {
    "async": "sinambi",
    "parallel": "bebarengan"
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

| Kategori JSON | Kunci JSON | Nilai Bawaan Celang | Contoh Basa Jawa | Keterangan & Penggunaan |
|---|---|---|---|---|
| **deklarasi** | `var` | `buat` | `gawe` | Deklarasi variabel yang dapat diubah nilainya |
| | `const` | `konstanta` | `tetep` | Deklarasi nilai konstan yang tidak dapat diubah |
| | `type_assign` | `itu` | `yaiku` | Anotasi tipe data statis (`x itu Angka`) |
| | `struct` | `bentuk` | `wujud` | Pendefinisian struktur data objek / tipe bentukan |
| | `ref` | `rujuk` | `rujuk` | Kata kunci referensi / pointer |
| | `as` | `sebagai` | `minangka` | Alias impor atau casting tipe |
| **percabangan** | `if` | `jika` | `yen` | Kondisi percabangan utama |
| | `else_if` | `selain jika` | `utawa yen` | Cabang alternatif bersyarat |
| | `else` | `selain` | `yen ora` | Cabang alternatif terakhir |
| | `then` | `maka` | `mula` | Kata kunci pembuka aksi kondisi |
| | `switch` | `pilih` | `cocogna` | Pemilihan multi-kasus |
| | `case` | `bila` | `nalika` | Klausa pencocokan kasus |
| | `default` | `lainnya` | `liyane` | Kasus default bila tidak ada yang cocok |
| **perulangan** | `while` | `selama` | `sasuwene` | Perulangan berbasis kondisi benar |
| | `for` | `untuk` | `kanggo` | Perulangan iterasi koleksi atau rentang |
| | `in` | `dalam` | `ing` | Operator keanggotaan koleksi perulangan |
| **fungsi** | `def` / `func` | `fungsi` | `fungsi` | Deklarasi fungsi / prosedur |
| | `return` | `kembalikan` | `balekna` | Mengembalikan nilai hasil pemrosesan fungsi |
| | `return_type` | `mengembalikan` | `ngasilakeun` | Anotasi tipe data nilai kembalian fungsi |
| | `defer` | `saat_selesai` | `pas_bubar` | Eksekusi tertunda di akhir cakupan blok |
| **penanganan_galat** | `try` | `coba` | `jajal` | Memulai blok penanganan pengecualian |
| | `catch` | `tangkap` | `cekel` | Menangkap galat yang terjadi dalam blok coba |
| | `finally` | `akhirnya` | `pungkasan` | Blok pembersihan yang selalu dijalankan |
| | `throw` | `lempar` | `uncalna` | Memicu lemparan galat secara eksplisit |
| **operator_logika** | `and` | `dan` | `lan` | Operasi logika konjungsi (AND) |
| | `or` | `atau` | `utawa` | Operasi logika disjungsi (OR) |
| | `not` | `bukan` | `dudu` | Operasi logika negasi (NOT) |
| | `null_coalesce` | `kalau kosong` | `yen_suwung` | Operator fallback nilai kosong (`a kalau kosong b`) |
| | `maybe` | `mungkin` | `mbokmenawa` | Operator chaining aman (safe navigation) |
| **literal** | `true` | `benar` | `bener` | Nilai kebenaran positif |
| | `false` | `salah` | `salah` | Nilai kebenaran negatif |
| | `null` | `kosong` | `suwung` | Nilai ketiadaan objek / null pointer |
| **konkurensi** | `async` | `sambil` | `sinambi` | Menjalankan blok tugas secara asinkron |
| | `parallel` | `bersamaan` | `bebarengan` | Menjalankan kumpulan tugas secara paralel |
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
