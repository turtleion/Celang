# Panduan Mulai Cepat & Langkah Pertama

Panduan ini akan membantu Anda memasang Celang di komputer Anda dan menjalankan program pertama dalam beberapa menit.

---

## 🛠️ Persiapan Awal

Jika Anda ingin membuat aplikasi Celang langsung dari kode sumbernya, pastikan komputer Anda telah memiliki:
- **Go** (versi 1.20 atau yang lebih baru).
- Sistem Operasi: Linux, macOS, atau Windows.

> **Catatan:** Jika Anda belum ingin memasang apa pun di komputer saat ini, Anda bisa langsung mencoba Celang di peramban web Anda melalui menu [Playground](/playground).

---

## 📥 Cara Memasang di Komputer

1. **Unduh kode proyek**:
   ```bash
   git clone https://github.com/celang/celang.git
   cd celang
   ```

2. **Bangun aplikasi Celang**:
   ```bash
   go build -o celang ./cmd/celang
   ```

3. *(Opsional)* **Pindahkan ke folder sistem agar bisa dipanggil dari mana saja**:
   ```bash
   # Di Linux atau macOS:
   sudo mv celang /usr/local/bin/
   ```

4. **Periksa hasil pemasangan**:
   ```bash
   ./celang versi
   ```
   Komputer akan menampilkan pesan sambutan dan nomor versi Celang yang aktif.

---

## 🚀 Menjalankan Program Pertama

1. Buat berkas baru bernama `halo.ce` di folder Anda:
   ```celang
   // halo.ce
   buat salam = "Halo, Dunia!"
   cetak(salam)

   buat sahabat = ["Budi", "Siti", "Andi"]
   untuk orang dalam sahabat {
       cetak("Selamat datang,", orang)
   }
   ```

2. Jalankan program Anda dengan perintah `run`:
   ```bash
   ./celang run halo.ce
   ```

3. Hasil yang muncul di layar:
   ```text
   Halo, Dunia!
   Selamat datang, Budi
   Selamat datang, Siti
   Selamat datang, Andi
   ```

---

## ⚡ Pilihan Cara Menjalankan

Secara bawaan (*default*), perintah `./celang run` langsung menjalankan program menggunakan mesin virtual cepat.

Jika di kemudian hari Anda ingin membandingkannya dengan cara kerja penafsir lama (misalnya saat menganalisis alur kode):
```bash
# Menjalankan melalui penafsir pohon logika langsung (AST Evaluator):
./celang run halo.ce --ast
```

Jika Anda penasaran dan ingin melihat bagaimana komputer menyusun instruksi rakitan di balik program Anda:
```bash
# Melihat daftar instruksi mesin rakitan:
./celang disasm halo.ce
```
