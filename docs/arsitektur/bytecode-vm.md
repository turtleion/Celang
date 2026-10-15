# Kompiler Bytecode & Mesin Virtual (VM)

Agar program Anda dapat berjalan secepat kilat tanpa membebani komputer, Celang menggunakan **Mesin Virtual Berbasis Instruksi Ringkas (Bytecode VM)** sebagai mesin penggerak utamanya.

---

## 1. Dua Cara Menjalankan Program

Bayangkan Anda ingin memasak sup:
* **Cara Pertama (Membaca Buku Tebal Berulang Kali)**: Setiap kali ingin memotong wortel atau menabur garam, Anda membuka lembar demi lembar buku resep yang tebal dari awal. Ini membutuhkan waktu lebih lama untuk membolak-balik halaman.
* **Cara Kedua (Menyalin Catatan Ringkas ke Kertas Selembar)**: Sebelum memasak, Anda merangkum seluruh langkah ke selembar kertas kecil berisi instruksi singkat: `1. Potong wortel`, `2. Rebus air`, `3. Tabur garam`. Sambil memasak, Anda membaca catatan kecil ini dalam sekejap mata.

Celang bekerja dengan **Cara Kedua**:
1. Pohon logika program Anda diringkas menjadi deretan instruksi kecil dan padat.
2. Mesin virtual Celang mengeksekusi instruksi kecil tersebut satu per satu dengan sangat cepat.

```mermaid
flowchart LR
    A[Pohon Logika Program] --> B[Penyusun Instruksi / Compiler]
    B --> C[Daftar Instruksi Ringkas]
    C --> D[Mesin Virtual Celang / VM]
    D --> E[Hasil Muncul di Layar]
```

---

## 2. Bagaimana Mesin Menghitung Nilai?

Mesin virtual Celang menggunakan konsep **Tumpukan Nilai** (seperti tumpukan piring di meja):

1. Misalkan ada perhitungan `10 + 25`.
2. Mesin menaruh angka `10` di meja.
3. Mesin menaruh angka `25` di atas angka `10`.
4. Saat perintah `Tambah` tiba, mesin mengambil dua angka teratas tersebut, menjumlahkannya menjadi `35`, lalu menaruh kembali angka `35` di meja.

Karena proses ini hanya melibatkan tumpukan nilai sederhana, komputer Anda bisa melakukannya dalam hitungan sekejap mata tanpa perlu berpikir rumit.

---

## 3. Pilihan Mode di Terminal

Secara bawaan (*default*), Celang **selalu menggunakan Mesin Virtual Cepat** setiap kali Anda menjalankan perintah:

```bash
# Menjalankan program (secara otomatis menggunakan mesin virtual cepat):
./celang run program.ce
```

Jika untuk keperluan belajar atau penelitian Anda ingin melihat bagaimana program dijalankan menggunakan penafsir pohon logika langsung (AST Evaluator):
```bash
# Menjalankan dengan mode penafsir pohon langsung:
./celang run program.ce --ast
```

Jika Anda ingin melihat daftar instruksi ringkas yang dibaca oleh mesin di balik program Anda:
```bash
# Melihat daftar instruksi ringkas di layar:
./celang disasm program.ce
```

---

## 4. Ringkasan Perbedaan Dua Mesin

| Hal yang Dilihat | Mesin Virtual Cepat (Bawaan) | Penafsir Pohon Langsung (`--ast`) |
| :--- | :--- | :--- |
| **Status Pemakaian** | Mode utama dan bawaan | Pilihan alternatif untuk analisis |
| **Cara Menjalankan** | Membaca daftar instruksi ringkas | Menelusuri bagan logika langkah demi langkah |
| **Beban Memori Komputer** | Sangat ringan | Sedikit lebih berat |
| **Kecepatan Perhitungan** | Sangat cepat | Kecepatan standar |
