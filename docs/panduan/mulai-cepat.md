# Mulai Cepat & Instalasi

Panduan ini akan membantumu memasang runtime Celang dan menjalankan program pertamamu dalam hitungan menit.

---

## 🛠️ Prasyarat

Untuk membangun runtime Celang dari kode sumber:
- **Go** versi 1.20 atau yang lebih baru.
- Sistem Operasi: Linux, macOS, atau Windows.

---

## 📥 Pemasangan (Build dari Sumber)

1. **Kloning repositori**:
   ```bash
   git clone https://github.com/celang/celang.git
   cd celang
   ```

2. **Kompilasi binary runtime**:
   ```bash
   go build -o celang ./cmd/celang
   ```

3. *(Opsional)* **Pindahkan binary ke PATH sistem**:
   ```bash
   # Di Linux / macOS:
   sudo mv celang /usr/local/bin/
   ```

4. **Verifikasi instalasi**:
   ```bash
   celang versi
   ```
   *Output:*
   ```text
   🌟 Celang Runtime v1.0.0 (linux/amd64)
      Sistem Bahasa Pemrograman Human-Centric & Syntax-Styling
      Bytecode Compiler & Virtual Machine (VM) Terintegrasi
      Lisensi: MIT
   ```

---

## 🚀 Menjalankan Program Pertama

1. Buat berkas baru bernama `halo.ce`:
   ```celang
   // halo.ce
   buat salam = "Halo, Dunia!"
   cetak(salam)

   buat nama = ["Budi", "Siti", "Andi"]
   untuk orang dalam nama {
       cetak("Selamat datang,", orang)
   }
   ```

2. Jalankan program dengan perintah `run`:
   ```bash
   celang run halo.ce
   ```

   *Output:*
   ```text
   Halo, Dunia!
   Selamat datang, Budi
   Selamat datang, Siti
   Selamat datang, Andi
   ```

---

## ⚡ Mode Eksekusi

Secara *default*, perintah `celang run` mengeksekusi program menggunakan **Bytecode Virtual Machine (VM)** berkecepatan tinggi.

Jika kamu ingin membandingkannya dengan mesin AST interpreter lama untuk keperluan analisis/debug:
```bash
# Menjalankan via AST Evaluator lama (legacy tree-walk)
celang run halo.ce --ast
```

Untuk membongkar instruksi bytecode di balik programmu:
```bash
# Melihat kode rakitan bytecode (Disassembler)
celang disasm halo.ce
```
