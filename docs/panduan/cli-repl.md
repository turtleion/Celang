# Panduan Perintah Terminal & Jendela Interaktif

Celang menyediakan dua cara praktis untuk berinteraksi dengan komputer: melalui **perintah terminal** untuk memproses berkas, dan melalui **jendela interaktif** untuk mencoba ketikan kode secara langsung.

---

## 💻 Daftar Perintah di Terminal

Jalankan `./celang bantuan` di terminal Anda untuk melihat daftar perintah yang tersedia:

| Perintah | Contoh Penggunaan | Kegunaan |
| :--- | :--- | :--- |
| `celang` | `./celang` | Langsung membuka jendela interaktif untuk mencoba kode seketika. |
| `celang repl` | `./celang repl --style=sunda` | Membuka jendela interaktif dengan gaya bahasa tertentu. |
| `celang run` | `./celang run program.ce` | Menjalankan program berkas `.ce` atau `.ast`. |
| `celang format` | `./celang format program.ce --style=javascript` | Menata ulang tampilan kode ke gaya lain di layar (tambahkan `-w` untuk langsung menimpa berkas). |
| `celang compile` | `./celang compile program.ce -o program.ast` | Mengubah kode teks menjadi berkas data logika U-AST berformat JSON. |
| `celang init` | `./celang init toko_buku` | Membuat folder proyek baru lengkap dengan berkas pembukanya. |
| `celang ast` | `./celang ast program.ce` | Menampilkan susunan pohon data logika kode ke layar. |
| `celang disasm` | `./celang disasm program.ce` | Menampilkan daftar instruksi mesin rakitan di balik program. |
| `celang versi` | `./celang versi` | Memeriksa nomor versi Celang yang terpasang. |

---

## 🎮 Jendela Interaktif (REPL)

Jendela interaktif ibarat ruang obrolan langsung antara Anda dan Celang. Setiap kali Anda mengetik satu baris perintah, Celang akan langsung menjawab hasilnya.

### Membuka Jendela Interaktif
Cukup ketik:
```bash
./celang
```

Layar Anda akan menampilkan salam pembuka:
```text
🌟 Celang Interactive REPL v1.0.0 (Bytecode VM Engine)
Gaya aktif: Bahasa Indonesia | Ketik ':bantuan' untuk petunjuk, ':keluar' untuk selesai.

celang> buat a = 10
celang> buat b = 20
celang> a + b
=> 30
```

---

### Perintah Khusus di Jendela Interaktif

Anda dapat mengatur sesi Anda menggunakan perintah yang diawali dengan tanda titik dua (`:`):

- `:bantuan` (atau `:?`) — Menampilkan petunjuk perintah di jendela interaktif.
- `:lingkungan` (atau `:vars`) — Melihat daftar semua variabel dan nilainya yang sedang diingat komputer saat ini.
- `:gaya <nama>` — Mengganti gaya bahasa aktif seketika (contoh: `:gaya sunda` atau `:gaya javascript`).
- `:bersih` (atau `:cls`) — Membersihkan layar terminal agar kembali rapi.
- `:reset` (atau `:ulang`) — Mengosongkan ingatan komputer kembali ke kondisi bersih seperti saat baru dibuka.
- `:keluar` (atau `:q`) — Menutup jendela interaktif dan kembali ke terminal biasa.

---

### Mengetik Perintah Lebih dari Satu Baris

Anda tidak perlu khawatir jika ingin membuat fungsi atau perintah yang panjang. Celang secara otomatis mengetahui jika Anda membuka kurung kurawal `{`:

```text
celang> fungsi luasPersegi(sisi) {
......>     kembalikan sisi * sisi
......> }
celang> luasPersegi(8)
=> 64
```
Tanda titik-titik (`......>`) menandakan bahwa Celang sedang sabar menunggu sampai Anda menutup kurung kurawal `}` untuk menjalankan seluruh blok kode tersebut bersama-sama.
