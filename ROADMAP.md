# Peta Perkembangan Proyek Celang (`ROADMAP.md`)

Dokumen ini mencatat pencapaian fitur yang telah berhasil diimplementasikan di dalam Celang serta rencana arah pengembangan untuk masa depan.

---

## 🏆 Fitur yang Telah Selesai Diimplementasikan

Seluruh fitur di bawah ini telah teruji penuh dan siap digunakan di versi Celang saat ini:

### 1. Fondasi Inti & Bahasa
- [x] **Pohon Logika Murni (Universal AST)**: Representasi kode independen dari gaya bahasa, mendukung ekspor dan impor berkas JSON lossless dua arah (`.ast`).
- [x] **Variabel Penjelajah Waktu**: Kemampuan variabel mengingat masa lalunya (`.sebelumnya`, `.awal`, dan `.histori`) dengan pembatas memori aman (50 catatan).
- [x] **Pengenalan Angka Alami (Smart Fallback Coercion)**: Kemampuan mengenali angka di dalam teks sehari-hari seperti `"Rp 150000" - 25000`.
- [x] **Penyaluran Data Berantai (`|>`)**: Mengalirkan data dari kiri ke kanan secara alami.
- [x] **Interpolasi Teks**: Menyisipkan nilai langsung ke dalam kalimat menggunakan tanda `{...}` atau `${...}`.
- [x] **Membongkar Paket Data (Destrukturisasi)**: Membongkar isi array dan map langsung ke dalam variabel terpisah.
- [x] **Jaring Pengaman Nilai Kosong**: Kata kunci `mungkin`, `?.`, dan nilai cadangan `kalau kosong` (atau `??`).
- [x] **Pembersihan Otomatis**: Kata kunci `saat_selesai` / `defer` untuk tindakan perapian terjamin.

### 2. Bentuk Data & Fitur Modern
- [x] **Bentuk Data Kustom (`bentuk`)**: Pemodelan data dengan metode, hak akses `privat`, dan jalan pintas `@`.
- [x] **Wadah Fleksibel (Generik)**: Dukungan parameter jenis data `<T>` dengan pembatas alami `harus` (contoh `<T harus Angka>`).
- [x] **Saluran Pengiriman Pesan (`saluran`)**: Pipa komunikasi aman dengan antrean, pengiriman (`.kirim`), penerimaan (`.terima`), penutupan (`.tutup`), dan perulangan `untuk .. dalam`.
- [x] **Titip & Jemput**: Menitipkan tugas di latar belakang (`titip`), menukar tiket untuk mengambil hasil (`jemput`), serta kombinasi `jemput_semua` dan `jemput_pertama`.
- [x] **Rumus Cepat Panah (`=>`)**: Penulisan rumus satu baris yang ringkas.

### 3. Mesin Eksekusi & Kinerja
- [x] **Kompiler Bytecode & Mesin Virtual (VM)**: Kompilasi pohon logika menjadi instruksi ringkas dan eksekusi tumpukan cepat sebagai mode eksekusi bawaan (*default*).
- [x] **Penafsir Pohon Logika Alternatif**: Opsi `--ast` untuk analisis langsung pohon logika.
- [x] **Diagnostik Empati Manusiawi**: Pesan kesalahan yang menampilkan kotak petunjuk ramah dan saran perbaikan praktis.

### 4. Ekosistem Gaya Bahasa & Perkakas
- [x] **Gaya Bahasa Fleksibel (Multi-Dialek)**: Profil bawaan Bahasa Indonesia, Basa Sunda, JavaScript, dan Pythonic.
- [x] **Format Ulang Kode Otomatis**: Perintah `celang format` untuk menerjemahkan tampilan kode antar gaya.
- [x] **Mengubah Gaya di Tengah Kode**: Blok `gaya "..." { ... }`.
- [x] **Jendela Interaktif (REPL)**: Sesi interaktif terminal dengan riwayat, dukungan banyak baris, dan perintah khusus.
- [x] **Pustaka Standar Bawaan**: Modul `matematika` (`math`), `teks` (`strings`), `waktu` (`time`), `berkas` (`fs`), dan `json`.
- [x] **WebAssembly Playground**: Menjalankan Celang langsung di peramban web tanpa perlu instalasi.
- [x] **Ekstensi Resmi Visual Studio Code**: Pewarnaan kode, cuplikan cepat (snippets), dan penutupan kurung otomatis.

---

## 🚀 Rencana Masa Depan (Eksplorasi Lanjutan)

Berikut beberapa ide yang dipertimbangkan untuk pengembangan berikutnya:

1. **Pengelola Paket Komunitas (Package Manager)**:
   - Kemudahan mengunduh dan berbagi modul buatan komunitas langsung lewat perintah terminal.
2. **Kompilasi Biner Mandiri**:
   - Kemampuan membungkus program Celang menjadi satu berkas biner mandiri yang bisa dijalankan langsung di komputer orang lain tanpa perlu memasang Celang terlebih dahulu.
3. **Penyempurnaan Integrasi Editor (Language Server Protocol / LSP)**:
   - Menyediakan fitur pelengkapan kata otomatis (*auto-complete*) dan pemeriksaan kesalahan saat mengetik di editor.
