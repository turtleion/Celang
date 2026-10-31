# Panduan Perintah Terminal & Jendela Interaktif

Celang menyediakan dua cara praktis untuk berinteraksi dengan komputer: melalui **perintah terminal** untuk memproses berkas, dan melalui **jendela interaktif** untuk mencoba ketikan kode secara langsung.

---

## 💻 Daftar Perintah di Terminal

Jalankan `./celang bantuan` di terminal Anda untuk melihat daftar perintah yang tersedia:

| Perintah | Contoh Penggunaan | Kegunaan |
| :--- | :--- | :--- |
| `celang` | `./celang` | Langsung membuka jendela interaktif untuk mencoba kode seketika. |
| `celang repl` | `./celang repl --style=jawa` | Membuka jendela interaktif dengan gaya bahasa tertentu. |
| `celang run` | `./celang run program.ce` | Menjalankan program berkas `.ce` atau `.ast`. |
| `celang bangun` | `./celang bangun game.ce -o GameKeren.exe` | Membangun aplikasi mandiri (*standalone executable*) siap pakai tanpa butuh instalasi Celang/Go. |
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
- `:gaya <nama>` — Mengganti gaya bahasa aktif seketika (contoh: `:gaya jawa` atau `:gaya javascript`).
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

---

## 📦 Membangun Aplikasi Mandiri (`celang bangun` / `celang build`)

Salah satu fitur paling praktis dan membanggakan di Celang adalah kemampuannya mengubah berkas kode program Anda menjadi **aplikasi biner mandiri (*standalone executable*)** dalam hitungan detik.

### Contoh Membangun Aplikasi
```bash
# Menggunakan perintah bahasa Indonesia:
./celang bangun game.ce -o Monster.exe

# Atau menggunakan alias internasional:
./celang build game.ce -o Monster.exe
```

### Mengapa Fitur Ini Luar Biasa?

1. **100% Bebas Prasyarat (*Zero Dependencies*)**:
   Berkas biner yang dihasilkan sudah mengemas seluruh mesin pengeksekusi Celang beserta kode program Anda. Pengguna akhir **sama sekali tidak perlu memasang Celang ataupun Go** di komputer mereka.
2. **Praktis Dibagikan (*Klik Dua Kali Langsung Jalan*)**:
   Kirimkan berkas `.exe` (Windows) atau biner (Linux/macOS) ke teman, guru, atau keluarga lewat flashdisk atau obrolan online. Mereka dapat langsung mengklik dua kali berkas tersebut untuk menjalankan karya Anda!
3. **Mendukung Argumen Baris Perintah**:
   Aplikasi biner mandiri yang Anda bangun tetap dapat menerima argumen terminal melalui variabel bawaan `argumen` (atau `args`).
4. **Bisa Kompilasi ke Windows Langsung dari Linux/macOS (*Cross-Compilation*)**:
   Anda tidak perlu menyalakan komputer Windows hanya untuk membuat file `.exe`. Cukup jalankan perintah berikut dari Linux atau macOS:
   ```bash
   # Otomatis mendeteksi target Windows dari akhiran .exe
   ./celang bangun game.ce -o GameMonster.exe

   # Atau tentukan target secara eksplisit
   ./celang bangun game.ce --target=windows
   ```
   Celang akan otomatis menyiapkan dan menyimpan template pelari Windows ke direktori `~/.celang/runners/windows-amd64.exe`. Pembuatan file `.exe` Windows berikutnya akan selesai dalam **0,02 detik**!
5. **Kustomisasi Pelari (*Custom Runner*)**:
   Jika Anda memiliki biner Windows kustom, Anda juga dapat menentukannya secara langsung lewat opsi `--runner`:
   ```bash
   ./celang bangun game.ce -o Game.exe --runner=path/ke/celang-windows.exe
   ```

