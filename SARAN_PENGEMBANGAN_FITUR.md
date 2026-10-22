# 🧭 Rekomendasi & Saran Pengembangan Fitur Celang

Dokumen ini memetakan evaluasi kesiapan ekosistem Celang saat ini serta rekomendasi langkah strategis berikutnya untuk membawa Celang dari bahasa pemrograman yang **siap pakai** menjadi **platform ekosistem berskala industri (*production-grade*)**.

---

## 📊 1. Evaluasi Kesiapan Celang Saat Ini (Current State)

Secara fungsional, Celang **sudah sangat siap dan menyenangkan digunakan** untuk kebutuhan:
* 🛠️ **Skrip Otomatisasi & Perkakas CLI** (Didukung kompilasi mandiri `.exe` / ELF).
* ⚙️ **Integrasi Pustaka C Lintas Platform** (Modul `ffi` berbasis *Dynamic ABI Invocation* tanpa butuh GCC/CGO).
* 🧪 **Pengembangan Teruji (*Test-Driven Development*)** (Test runner bawaan `celang uji` / `celang test` dengan asersi lengkap).
* 📚 **Edukasi & Pembelajaran Pemrograman** (Sintaks dwibahasa alami Indonesia/Inggris, diagnostik empati, tur interaktif WebAssembly).

---

## 🚀 2. Empat Pilar Fitur Strategis Berikutnya

Untuk memperluas adopsi pengembang di dunia industri dan komunitas *open-source*, berikut adalah 4 pilar fitur yang sangat direkomendasikan untuk diimplementasikan:

### Pilar 1: Package Manager & Manajemen Dependensi Git (`celang pasang` / `celang get`)
* **Tujuan**: Memungkinkan pengembang mengunduh, menggunakan, dan membagikan modul/pustaka pihak ketiga dari GitHub/Git tanpa perlu salin-tempel manual.
* **Fitur Utama**:
  1. Perintah CLI:
     ```bash
     celang pasang github.com/pengguna/pustaka-jwt
     # Alias internasional:
     celang get github.com/user/jwt-lib
     ```
  2. Resolusi berkas konfigurasi `celang.mod` dan berkas pengunci versi `celang.kunci` (*lockfile*).
  3. Struktur folder dependensi lokal (`celang_modules/`) atau cache terpusat global (`~/.celang/pkg/`).

### Pilar 2: Language Server Protocol (LSP) untuk Editor (`celang lsp`)
* **Tujuan**: Menghadirkan pengalaman koding kelas satu (*first-class developer experience*) di Visual Studio Code, Neovim, dan JetBrains.
* **Fitur Utama**:
  1. **Autocomplete & IntelliSense**: Saran kata kunci, nama variabel, properti bentuk (*struct*), dan fungsi modul saat mengetik.
  2. **Hover Documentation**: Menampilkan dokumentasi ringkas, tipe data, dan petunjuk fungsi saat kursor berada di atas kode.
  3. **Diagnostik Waktu-Nyata (*Real-Time Linting*)**: Garis bawah merah saat ada kesalahan sintaks tanpa harus menjalankan program terlebih dahulu di terminal.
  4. **Go-to-Definition**: Lompat langsung ke lokasi deklarasi fungsi atau variabel (F12 di VS Code).

### Pilar 3: Generator Dokumentasi Proyek Otomatis (`celang dok` / `celang doc`)
* **Tujuan**: Menjamin keterbacaan proyek besar dengan mendokumentasikan kode secara otomatis.
* **Fitur Utama**:
  1. Memindai komentar dokumentasi berformat `///` atau blok `/** ... */` di atas fungsi, bentuk, dan konstanta.
  2. Perintah CLI:
     ```bash
     celang dok [--format=html|markdown] [-o docs/]
     ```
  3. Menghasilkan situs dokumentasi statis yang dapat langsung diunggah ke GitHub Pages atau dibaca secara lokal.

### Pilar 4: Ekosistem Web & Microservices Siap Pakai (Pustaka Web Terintegrasi)
* **Tujuan**: Mempermudah pembuatan API RESTful dan layanan mikro backend.
* **Fitur Utama**:
  1. Penyempurnaan pustaka bawaan `http` dengan *Router* dinamis yang mendukung parameter rute (contoh: `/pengguna/:id`).
  2. Lapisan *Middleware* siap pakai:
     - CORS (*Cross-Origin Resource Sharing*)
     - JSON Body Parser & Serializer otomatis
     - Logger permintaan HTTP
     - Penanganan token otentikasi (JWT)

---

## 💡 3. Fitur Tambahan Pelengkap (*Nice-to-Have*)

1. **Benchmark Runner Bawaan (`celang tolok` / `celang bench`)**:
   - Pengujian performa fungsi berulang dengan output alokasi memori dan durasi rata-rata per operasi.
2. **Pemeriksa Tipe Statis AOT Opsional (`celang periksa` / `celang check`)**:
   - Analisis statis tanpa menjalankan kode untuk memvalidasi kompatibilitas tipe generik dan keberadaan variabel.
3. **Debug Adapter Protocol (DAP)**:
   - Integrasi titik henti (*breakpoints*), *step into*, dan inspeksi variabel langsung di panel debugger VS Code.

---

## 🎯 Kesimpulan & Rekomendasi Urutan Pengerjaan

| Urutan | Fitur | Dampak terhadap Pengembang | Tingkat Kerumitan |
| :---: | :--- | :--- | :---: |
| **1** | **Package Manager (`celang pasang`)** | ⭐⭐⭐⭐⭐ (Sangat Tinggi - Kolaborasi Komunitas) | Sedang |
| **2** | **Language Server Protocol (`celang lsp`)** | ⭐⭐⭐⭐⭐ (Sangat Tinggi - Kesenangan Koding Editor) | Sedang - Tinggi |
| **3** | **Generator Dokumentasi (`celang dok`)** | ⭐⭐⭐⭐ (Tinggi - Pemeliharaan Proyek) | Rendah - Sedang |
| **4** | **Penyempurnaan Web/HTTP Router** | ⭐⭐⭐⭐ (Tinggi - Kasus Penggunaan Riil Backend) | Sedang |
