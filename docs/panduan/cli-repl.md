# CLI & Interactive REPL

Celang dilengkapi dengan antarmuka baris perintah (*Command Line Interface*) yang kaya fitur serta REPL interaktif terintegrasi.

---

## 💻 Panduan Perintah CLI

Cukup jalankan `celang bantuan` atau `celang -h` untuk melihat seluruh opsi yang tersedia:

| Perintah | Deskripsi |
| :--- | :--- |
| `celang` | Langsung masuk ke sesi Interactive REPL |
| `celang repl [--style=...]` | Membuka REPL dengan gaya sintaks awal tertentu |
| `celang run <file> [--ast] [-- args...]` | Menjalankan program Celang (`.ce`) atau AST (`.ast`) |
| `celang disasm <file.ce>` | Menampilkan instruksi rakitan Bytecode (Disassembler) |
| `celang format [-w] <file> [--style=...]` | Memformat kode ke gaya tertentu (`-w` untuk menimpa file langsung) |
| `celang compile <file.ce> [-o out.ast]` | Mengompilasi kode menjadi berkas U-AST JSON portabel |
| `celang ast <file.ce>` | Mencetak pohon struktur Universal AST ke terminal |
| `celang init [nama_proyek]` | Menginisialisasi proyek Celang baru (`celang.mod` & `utama.ce`) |
| `celang versi` | Menampilkan informasi versi runtime Celang |

---

## 🎮 Interactive REPL

REPL (*Read-Eval-Print Loop*) memungkinkanmu mengeksplorasi sintaks Celang secara langsung dan instan.

### Memulai Sesi
```bash
celang
```

Tampilan sambutan:
```text
🌟 Celang Interactive REPL v1.0.0 (Bytecode VM Engine)
Gaya aktif: Bahasa Indonesia (Baku) | Ketik ':bantuan' untuk petunjuk, ':keluar' untuk selesai.

celang> buat a = 10
celang> buat b = 20
celang> a + b
=> 30
```

### Perintah Khusus (*Meta-Commands*)

Di dalam REPL, awali perintah dengan titik dua (`:`):

- `:bantuan` atau `:?` — Menampilkan panduan perintah REPL.
- `:lingkungan` atau `:vars` — Melihat daftar variabel dan nilainya yang saat ini aktif di memori.
- `:gaya <nama|path>` — Mengganti profil gaya sintaks aktif (misal `:gaya sunda` atau `:gaya javascript`).
- `:bc <ekspresi>` — Membongkar ekspresi menjadi instruksi Bytecode VM.
- `:ast <ekspresi>` — Melihat pohon Universal AST JSON dari baris kode.
- `:bersih` atau `:cls` — Membersihkan layar terminal.
- `:reset` atau `:ulang` — Mengosongkan memori variabel kembali ke kondisi awal.
- `:keluar` atau `:q` — Keluar dari sesi REPL.

### Input Multi-Baris Cerdas
REPL otomatis mendeteksi tanda kurung kurawal `{ ... }` yang belum tertutup dan menyediakan prompt multi-baris (`......>`):

```text
celang> fungsi luas_persegi(sisi) {
......>     kembalikan sisi * sisi
......> }
celang> luas_persegi(8)
=> 64
```
