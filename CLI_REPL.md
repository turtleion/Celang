# Panduan Perintah Terminal & Jendela Interaktif (`CLI_REPL.md`)

Dokumen ini adalah panduan lengkap cara berinteraksi dengan Celang melalui perintah di terminal (*Command Line Interface*) dan melalui jendela interaktif (*REPL*).

---

## 1. Perintah Baris di Terminal (`celang <perintah>`)

Format umum penulisan perintah:
```bash
./celang [perintah] [opsi] [berkas / argumen]
```

### Tabel Perintah Lengkap

| Perintah | Contoh Perintah | Kegunaan |
| :--- | :--- | :--- |
| `celang` (tanpa perintah) | `./celang` | Langsung masuk ke jendela interaktif (REPL) |
| `celang repl` | `./celang repl --style=jawa` | Membuka jendela interaktif dengan gaya bahasa tertentu |
| `celang run` | `./celang run program.ce` | Menjalankan program Celang menggunakan mesin virtual cepat (bawaan) |
| `celang run ... --ast` | `./celang run program.ce --ast` | Menjalankan program menggunakan penafsir pohon logika langsung |
| `celang run ... -- args` | `./celang run program.ce -- a b` | Meneruskan argumen tambahan ke variabel `argumen` di dalam program |
| `celang format` | `./celang format program.ce --style=pythonic` | Menata ulang tampilan kode ke gaya bahasa lain di layar terminal |
| `celang format -w` | `./celang format -w program.ce --style=javascript` | Menata ulang dan langsung menimpa isi berkas program secara permanen |
| `celang compile` | `./celang compile program.ce -o program.ast` | Mengubah kode teks menjadi berkas data logika U-AST berformat JSON |
| `celang ast` | `./celang ast program.ce` | Menampilkan susunan pohon data logika kode ke terminal |
| `celang disasm` | `./celang disasm program.ce` | Menampilkan daftar instruksi mesin rakitan di balik program |
| `celang init` | `./celang init nama_proyek` | Membuat folder proyek baru dengan berkas `celang.mod` dan `utama.ce` |
| `celang versi` | `./celang versi` | Memeriksa nomor versi Celang dan sistem operasi |
| `celang bantuan` | `./celang bantuan` | Menampilkan daftar bantuan seluruh perintah |

---

## 2. Jendela Interaktif (REPL)

Jendela interaktif adalah tempat Anda mengetik sebaris instruksi dan komputer langsung memberikan jawaban seketika.

### Memulai Sesi
Cukup ketik `./celang` di terminal:
```text
🌟 Celang Interactive REPL v1.0.0 (Bytecode VM Engine)
Gaya aktif: Bahasa Indonesia | Ketik ':bantuan' untuk petunjuk, ':keluar' untuk selesai.

celang> buat a = 10
celang> buat b = 25
celang> a + b
=> 35
```

---

### Perintah Pengaturan di Dalam REPL (Diawali Titik Dua `:`)

| Perintah | Panggilan Cepat | Kegunaan |
| :--- | :--- | :--- |
| `:keluar` | `:q`, `:exit` | Menutup sesi interaktif dan kembali ke terminal |
| `:bersih` | `:cls`, `:clear` | Membersihkan layar terminal agar kembali rapi |
| `:lingkungan` | `:vars`, `:env` | Melihat daftar variabel yang sedang diingat komputer |
| `:gaya <nama>` | `:style <nama>` | Mengganti gaya bahasa aktif seketika (misal `:gaya jawa` atau `:gaya javascript`) |
| `:reset` | `:ulang` | Mengosongkan ingatan komputer kembali ke kondisi bersih |
| `:ast <kode>` | `:uast <kode>` | Melihat bentuk pohon logika U-AST JSON dari baris kode |
| `:bc <kode>` | `:disasm <kode>` | Melihat instruksi rakitan mesin dari baris kode |
| `:bantuan` | `:?`, `:help` | Menampilkan daftar bantuan di jendela interaktif |

---

### Mengetik Perintah Lebih dari Satu Baris
Jika Anda membuka kurung kurawal `{`, kurung biasa `(`, atau kurung siku `[`, Celang akan otomatis menunggu sampai Anda menutup kurung tersebut sebelum menjalankan seluruh kodenya bersama-sama:

```text
celang> fungsi luasPersegi(sisi) {
......>     kembalikan sisi * sisi
......> }
celang> luasPersegi(7)
=> 49
```
Tanda `......>` menunjukkan bahwa komputer sedang menunggu kelanjutan baris kode Anda.
