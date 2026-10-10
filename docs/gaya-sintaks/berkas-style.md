# Membuat Berkas .style Kustom

Celang tidak membatasi Anda hanya pada dialek bawaan. Anda dapat membuat dialek Anda sendiri — baik itu bahasa daerah (Jawa, Sunda, Minang, Bali), bahasa internasional lainnya, atau istilah teknis spesifik tim Anda — menggunakan berkas konfigurasi berekstensi `.style`.

---

## 1. Struktur Berkas `.style`

Berkas `.style` adalah dokumen JSON yang memetakan konsep semantik bahasa ke kata kunci yang Anda inginkan.

Berikut adalah contoh lengkap berkas [examples/sunda.style](file:///home/yosatsu/Coding/Celang/examples/sunda.style):

```json
{
  "name": "sunda",
  "display_name": "Bahasa Sunda",
  "var_decl": "jieun",
  "const_decl": "tetep",
  "type_assign": "nyaeta",
  "if": "lamun",
  "else_if": "atawa lamun",
  "else": "lamun henteu",
  "switch": "cocogkeun",
  "case": "bisi",
  "default": "lianna",
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
  "return_arrow": "ngahasilkeun",
  "block_style": "brace",
  "indent_width": 4
}
```

---

## 2. Penjelasan Kunci Konfigurasi

| Kunci (*Key*) | Deskripsi | Nilai Default (ID) |
| :--- | :--- | :--- |
| `name` | Pengidentifikasi unik profil gaya (string singkat). | `"indonesian"` |
| `display_name` | Nama yang ditampilkan pada log dan dokumentasi. | `"Bahasa Indonesia"` |
| `var_decl` | Kata kunci pembuatan variabel baru. | `"buat"` |
| `const_decl` | Kata kunci deklarasi konstanta. | `"konstanta"` |
| `type_assign` | Kata kunci / simbol penetapan tipe eksplisit. | `"itu"` |
| `if` | Kata kunci kondisi percabangan utama. | `"jika"` |
| `else_if` | Kata kunci cabang kondisi lanjutan. | `"selain jika"` |
| `else` | Kata kunci cabang kondisi penutup. | `"selain"` |
| `while` | Kata kunci perulangan selama kondisi benar. | `"selama"` |
| `for` | Kata kunci perulangan koleksi/rentang. | `"untuk"` |
| `in` | Kata hubung anggota iterasi perulangan. | `"dalam"` |
| `func` | Kata kunci deklarasi fungsi. | `"fungsi"` |
| `return` | Kata kunci pengembalian nilai dari fungsi. | `"kembalikan"` |
| `async` | Kata kunci konkurensi non-blocking di latar. | `"sambil"` |
| `parallel` | Kata kunci eksekusi banyak tugas bersamaan. | `"bersamaan"` |
| `true` / `false` | Literal nilai kebenaran logika (*boolean*). | `"benar"` / `"salah"` |
| `null` | Literal ketiadaan nilai (*null*). | `"kosong"` |
| `and` / `or` / `not`| Operator logika deskriptif. | `"dan"` / `"atau"` / `"bukan"` |

---

## 3. Cara Menggunakan Berkas `.style`

### A. Menjalankan Kode dengan Style Kustom
Gunakan bendera `--style` dan arahkan ke path berkas `.style`:

```bash
celang run --style ./sunda.style program_sunda.ce
```

### B. Memformat dan Mengonversi Kode
Anda dapat mengonversi kode yang sudah ada ke style baru:

```bash
celang format --to ./sunda.style program_utama.ce -o program_versi_sunda.ce
```

### C. Mendaftarkan Secara Global
Simpan berkas style Anda di direktori profil Celang atau sertakan dalam repositori proyek untuk dibagikan ke seluruh anggota tim.
