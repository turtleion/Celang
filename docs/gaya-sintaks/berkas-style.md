# Membuat Berkas .style Kustom

Celang tidak membatasi Anda hanya pada dialek bawaan. Anda dapat membuat dialek Anda sendiri — baik itu bahasa daerah (Sunda, Jawa, Minang, Bali), bahasa internasional lainnya, atau istilah teknis spesifik tim Anda — menggunakan berkas konfigurasi berekstensi `.style`.

Format `.style` di Celang dirancang secara **terstruktur per bab/bagian logis**, sehingga sangat mudah dibaca, dirawat, dan dikembangkan seiring bertambahnya fitur bahasa.

---

## 1. Struktur Berkas `.style` (Format Per Bab)

Berkas `.style` adalah dokumen JSON yang memetakan konsep semantik U-AST ke kata kunci pilihan Anda. Konfigurasi dibagi ke dalam beberapa bab kategori:

Berikut adalah contoh lengkap berkas [examples/sunda.style](file:///home/yosatsu/Coding/Celang/examples/sunda.style):

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

## 2. Rincian Bab & Kunci Konfigurasi

### 2.1 `metadata` (Informasi Profil)
| Kunci | Deskripsi | Contoh |
| :--- | :--- | :--- |
| `name` | Pengidentifikasi unik gaya (huruf kecil tanpa spasi). | `"sunda"` |
| `display_name` | Nama resmi gaya untuk log terminal dan format. | `"Basa Sunda"` |
| `version` | Versi dialek/gaya sintaks. | `"1.0.0"` |
| `author` | Penulis atau pemelihara dialek. | `"Komunitas Celang"` |
| `description` | Deskripsi singkat kegunaan dialek. | `"Dialek basa Sunda"` |

### 2.2 `deklarasi` (Variabel, Tipe & Struktur)
| Kunci | Makna Semantik | Bawaan (ID) |
| :--- | :--- | :--- |
| `var` | Deklarasi variabel yang dapat diubah. | `"buat"` |
| `const` | Deklarasi konstanta nilai tetap. | `"konstanta"` |
| `type_assign` | Penetapan tipe data eksplisit. | `"itu"` (atau `:`) |
| `struct` | Deklarasi bentuk/struktur objek. | `"bentuk"` |
| `ref` | Deklarasi parameter rujukan (*pass-by-reference*). | `"rujuk"` |
| `as` | Konversi tipe eksplisit (*type casting*). | `"sebagai"` |

### 2.3 `percabangan` (Kondisi & Kasus)
| Kunci | Makna Semantik | Bawaan (ID) |
| :--- | :--- | :--- |
| `if` | Kondisi percabangan utama. | `"jika"` |
| `else_if` | Cabang kondisi berikutnya. | `"selain jika"` |
| `else` | Cabang kondisi terakhir. | `"selain"` |
| `then` | Kata hubung opsional penegas blok. | `"maka"` |
| `switch` | Pencocokan pola kasus. | `"cocok"` |
| `case` | Klausul nilai pencocokan kasus. | `"kasus"` |
| `default` | Klausul bawaan jika tidak ada kasus cocok. | `"lainnya"` |

### 2.4 `perulangan` (Iterasi & Koleksi)
| Kunci | Makna Semantik | Bawaan (ID) |
| :--- | :--- | :--- |
| `while` | Perulangan selama kondisi bernilai benar. | `"selama"` |
| `for` | Perulangan koleksi/rentang. | `"untuk"` |
| `in` | Kata hubung anggota iterasi. | `"dalam"` |

### 2.5 `fungsi` (Subrutin & Eksekusi Tertunda)
| Kunci | Makna Semantik | Bawaan (ID) |
| :--- | :--- | :--- |
| `def` | Deklarasi pembuatan fungsi. | `"fungsi"` |
| `return` | Pengembalian nilai dari fungsi. | `"kembalikan"` |
| `return_type` | Anotasi tipe kembalian fungsi. | `"mengembalikan"` (atau `->`) |
| `defer` | Eksekusi tertunda tepat sebelum fungsi keluar. | `"saat_selesai"` |

### 2.6 `penanganan_galat` (Try-Catch-Finally-Throw)
| Kunci | Makna Semantik | Bawaan (ID) |
| :--- | :--- | :--- |
| `try` | Blok percobaan eksekusi yang rawan galat. | `"coba"` |
| `catch` | Blok penangkap objek galat. | `"tangkap"` |
| `finally` | Blok pembersihan yang dijamin selalu berjalan. | `"akhirnya"` |
| `throw` | Melempar galat kustom secara sengaja. | `"lempar"` |

### 2.7 `operator_logika` & `literal`
| Kategori | Kunci | Makna Semantik | Bawaan (ID) |
| :--- | :--- | :--- | :--- |
| **Operator** | `and`, `or`, `not` | Logika kebenaran. | `"dan"`, `"atau"`, `"bukan"` |
| | `null_coalesce` | Nilai fallback jika kosong. | `"kalau kosong"` (atau `??`) |
| | `maybe` | Navigasi aman null (*safe navigation*). | `"mungkin"` (atau `?.`) |
| **Literal** | `true`, `false` | Nilai kebenaran boolean. | `"benar"`, `"salah"` |
| | `null` | Ketiadaan nilai. | `"kosong"` |

### 2.8 `konkurensi` & `format`
| Kategori | Kunci | Makna Semantik | Bawaan (ID) |
| :--- | :--- | :--- | :--- |
| **Konkurensi** | `async` | Eksekusi asinkron latar belakang. | `"sambil"` |
| | `parallel` | Eksekusi tugas serentak. | `"bersamaan"` |
| **Format** | `block_style` | Gaya blok kode: `"brace"` (`{}`) atau `"keyword"`. | `"brace"` |
| | `indent_width` | Jumlah spasi per tingkat indentasi. | `4` |
| | `quote_style` | Gaya kutipan teks: `"double"` atau `"single"`. | `"double"` |

::: tip Kompatibilitas Mundur (*Backward Compatibility*)
Celang mendukung format berbab ini secara native, namun tetap dapat membaca format lama (*flat JSON*) secara otomatis jika Anda memiliki berkas `.style` versi sebelumnya.
:::

---

## 3. Cara Menggunakan Berkas `.style`

### A. Menjalankan Kode dengan Dialek Kustom
Cukup arahkan bendera `--style` ke lokasi berkas `.style`:

```bash
celang run --style ./examples/sunda.style program_sunda.ce
```

### B. Memformat dan Mengonversi Kode Antar-Dialek
Anda dapat mengonversi kode yang ada dari bahasa Indonesia ke dialek kustom Anda secara otomatis menggunakan perintah `format`:

```bash
# Menampilkan hasil format ke terminal
celang format --style ./examples/sunda.style program_indo.ce

# Mengonversi dan menyimpan langsung ke file baru
celang format --style ./examples/sunda.style program_indo.ce -w
```
