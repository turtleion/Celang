# Membuat Gaya Bahasa Sendiri (Berkas .style)

Celang tidak membatasi Anda hanya pada gaya bahasa yang sudah ada. Anda dapat membuat gaya bahasa Anda sendiri—baik itu bahasa daerah (seperti Basa Sunda, Jawa, Minang, Batak, Bali), bahasa internasional lainnya, maupun istilah akrab tim Anda—menggunakan berkas konfigurasi berekstensi `.style`.

Berkas `.style` adalah dokumen sederhana (berformat JSON) yang berisi daftar padanan kata: kata apa yang ingin Anda gunakan untuk menggantikan kata kunci standar Celang.

---

## 1. Contoh Berkas Gaya Basa Sunda (`sunda.style`)

Berikut adalah contoh isi berkas [examples/sunda.style](file:///home/yosatsu/Coding/Celang/examples/sunda.style):

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

## 2. Bagian-Bagian yang Perlu Anda Isi

Anda cukup mengisi kata-kata yang sesuai dengan bahasa pilihan Anda:

1. **`metadata`**: Keterangan nama gaya dan siapa pembuatnya.
2. **`deklarasi`**: Kata untuk membuat variabel (`var`), nilai tetap (`const`), dan bentuk data (`struct`).
3. **`percabangan`**: Kata untuk mengambil keputusan, seperti `jika` (`if`), `selain` (`else`), dan pilihan kasus (`case`).
4. **`perulangan`**: Kata untuk pengulangan, seperti `selama` (`while`) dan `untuk ... dalam` (`for ... in`).
5. **`fungsi`**: Kata untuk membuat rumus/fungsi (`def`) dan mengembalikan hasil (`return`).
6. **`penanganan_galat`**: Kata untuk mencoba perintah rawan masalah (`try`, `catch`, `finally`, `throw`).
7. **`operator_logika` & `literal`**: Kata untuk logika `dan`, `atau`, `bukan`, serta nilai kebenaran `benar`, `salah`, dan `kosong`.
8. **`format`**: Pengaturan kerapian tampilan, seperti jarak spasi indentasi (biasanya 4 spasi).

---

## 3. Cara Menggunakan Berkas `.style` Anda

Setelah Anda menyimpan berkas gaya Anda (misalnya `jawa.style` atau `sunda.style`):

### A. Membuka Jendela Interaktif dengan Gaya Anda
```bash
./celang repl --style=examples/sunda.style
```
Jendela interaktif akan menyambut Anda langsung dengan aturan kata kunci gaya baru tersebut.

### B. Menerjemahkan Tampilan Kode yang Ada
Anda dapat menata ulang berkas kode berbahasa Indonesia menjadi gaya baru Anda:

```bash
# Menampilkan hasil terjemahan di layar terminal:
./celang format program.ce --style=examples/sunda.style

# Langsung mengubah dan menimpa isi berkas program:
./celang format -w program.ce --style=examples/sunda.style
```

### C. Menjalankan Kode
Anda dapat langsung menjalankan berkas kode yang sudah ditulis dalam dialek tersebut seperti biasa:
```bash
./celang run program.ce
```
Celang secara cerdas mengenali kata kunci yang telah didaftarkan dan langsung menjalankannya tanpa hambatan.
