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

<div class="vm-flow-container">
  <div class="vm-flow-step">
    <div class="vm-step-icon">🌳</div>
    <div class="vm-step-title">Pohon Logika (U-AST)</div>
    <div class="vm-step-desc">Struktur data program hasil parsing</div>
  </div>
  <div class="vm-flow-arrow">➔</div>
  <div class="vm-flow-step">
    <div class="vm-step-icon">⚙️</div>
    <div class="vm-step-title">Kompiler Bytecode</div>
    <div class="vm-step-desc">Menyusun instruksi kecil & padat</div>
  </div>
  <div class="vm-flow-arrow">➔</div>
  <div class="vm-flow-step">
    <div class="vm-step-icon">📜</div>
    <div class="vm-step-title">Instruksi Chunk</div>
    <div class="vm-step-desc">Deretan opcode ringkas di memori</div>
  </div>
  <div class="vm-flow-arrow">➔</div>
  <div class="vm-flow-step">
    <div class="vm-step-icon">⚡</div>
    <div class="vm-step-title">Mesin Virtual (VM)</div>
    <div class="vm-step-desc">Eksekusi stack super cepat</div>
  </div>
  <div class="vm-flow-arrow">➔</div>
  <div class="vm-flow-step">
    <div class="vm-step-icon">🖥️</div>
    <div class="vm-step-title">Hasil Output</div>
    <div class="vm-step-desc">Tampil langsung di terminal</div>
  </div>
</div>

<style>
.vm-flow-container {
  display: flex;
  align-items: stretch;
  justify-content: space-between;
  gap: 10px;
  margin: 20px 0;
  padding: 16px;
  background: var(--vp-c-bg-soft);
  border: 1px solid var(--vp-c-divider);
  border-radius: 12px;
  overflow-x: auto;
}
.vm-flow-step {
  flex: 1;
  min-width: 130px;
  background: var(--vp-c-bg);
  border: 1px solid var(--vp-c-divider);
  border-radius: 8px;
  padding: 12px 10px;
  text-align: center;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.03);
}
.vm-step-icon {
  font-size: 22px;
  margin-bottom: 4px;
}
.vm-step-title {
  font-weight: 700;
  font-size: 13px;
  color: var(--vp-c-text-1);
  margin-bottom: 4px;
}
.vm-step-desc {
  font-size: 11px;
  color: var(--vp-c-text-2);
  line-height: 1.3;
}
.vm-flow-arrow {
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  color: var(--vp-c-brand-1);
  font-weight: 700;
  opacity: 0.7;
}
@media (max-width: 768px) {
  .vm-flow-container {
    flex-direction: column;
    align-items: center;
  }
  .vm-flow-step {
    width: 100%;
    min-width: 0;
  }
  .vm-flow-arrow {
    transform: rotate(90deg);
    padding: 4px 0;
  }
}
</style>

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
