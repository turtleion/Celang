# Pohon Sintaks Universal (Universal AST / U-AST)

Pernahkah Anda memikirkan sebuah ide, lalu Anda bisa menceritakan ide tersebut dalam Bahasa Indonesia, Bahasa Inggris, atau Basa Jawa? 

Idenya tetap sama, yang berganti hanyalah kata-kata yang keluar dari mulut Anda.

Itulah cara kerja **Universal Abstract Syntax Tree (U-AST)** atau **Pohon Logika Murni** di Celang.

---

## 1. Ide Pikiran vs Kata-Kata yang Terlihat

Di kebanyakan bahasa pemrograman konvensional, teks tulisan dan logika program menempel secara kaku. Jika Anda mengubah kata kuncinya, komputernya langsung bingung.

Celang memisahkan keduanya:
- **Pohon Logika (U-AST)** ibarat **isi pikiran murni** dari program Anda.
- **Teks Kode (.ce)** ibarat **pakaian atau bahasa** yang digunakan untuk membaca pikiran tersebut ke mata manusia.

<div class="uast-flow-container">
  <div class="uast-flow-col">
    <div class="uast-node-header">Dialek Masukan</div>
    <div class="uast-node input-node">🇮🇩 Bahasa Indonesia</div>
    <div class="uast-node input-node">🌐 JavaScript Style</div>
    <div class="uast-node input-node">🌾 Basa Jawa</div>
  </div>
  <div class="uast-flow-arrow">➔</div>
  <div class="uast-flow-col center-col">
    <div class="uast-node-header">Inti Celang</div>
    <div class="uast-node parser-node">⚙️ Parser & Normalizer</div>
    <div class="uast-flow-arrow down-arrow">⬇</div>
    <div class="uast-node core-node">🌳 U-AST<br/><span style="font-size: 11px; opacity: 0.8;">Pohon Logika Universal</span></div>
  </div>
  <div class="uast-flow-arrow">➔</div>
  <div class="uast-flow-col">
    <div class="uast-node-header">Bentuk Keluaran</div>
    <div class="uast-node out-node">⚡ Eksekusi VM Langsung</div>
    <div class="uast-node out-node">🎨 Format ke Gaya Apa Pun</div>
    <div class="uast-node out-node">📄 Ekspor Berkas .ast JSON</div>
  </div>
</div>

<style>
.uast-flow-container {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin: 24px 0;
  padding: 20px;
  background: var(--vp-c-bg-soft);
  border: 1px solid var(--vp-c-divider);
  border-radius: 12px;
  overflow-x: auto;
}
.uast-flow-col {
  display: flex;
  flex-direction: column;
  gap: 8px;
  flex: 1;
  min-width: 140px;
}
.uast-node-header {
  font-size: 11px;
  font-weight: 700;
  text-transform: uppercase;
  color: var(--vp-c-text-2);
  margin-bottom: 2px;
  text-align: center;
  letter-spacing: 0.5px;
}
.uast-node {
  background: var(--vp-c-bg);
  border: 1px solid var(--vp-c-divider);
  border-radius: 8px;
  padding: 10px 12px;
  font-size: 12px;
  font-weight: 600;
  text-align: center;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.03);
}
.input-node {
  border-left: 3px solid #3b82f6;
}
.parser-node {
  border: 1px solid #f59e0b;
  background: rgba(245, 158, 11, 0.08);
}
.core-node {
  border: 2px solid var(--vp-c-brand-1);
  background: rgba(62, 175, 124, 0.1);
  color: var(--vp-c-brand-1);
}
.out-node {
  border-left: 3px solid var(--vp-c-brand-1);
}
.uast-flow-arrow {
  font-size: 18px;
  color: var(--vp-c-brand-1);
  font-weight: 700;
  opacity: 0.7;
}
.down-arrow {
  text-align: center;
}
@media (max-width: 768px) {
  .uast-flow-container {
    flex-direction: column;
    align-items: stretch;
  }
  .uast-flow-arrow:not(.down-arrow) {
    text-align: center;
    transform: rotate(90deg);
    padding: 4px 0;
  }
}
</style>

### Apa Manfaatnya bagi Pengguna?
1. **Tidak Ada Perdebatan Format**: Jika rekan kerja Anda lebih suka gaya JavaScript dan Anda lebih suka Bahasa Indonesia, Anda berdua bisa melihat logika yang sama persis sesuai kenyamanan mata masing-masing.
2. **Kekebalan dari Kesalahan Ejaan**: Karena yang disimpan dan dijalankan adalah inti logikanya, mengganti tampilan tidak akan merusak alur perhitungan program.
3. **Bisa Disimpan sebagai Data**: Seluruh pohon logika ini dapat disimpan ke berkas data berekstensi `.ast` dan dijalankan langsung kapan saja.

---

## 2. Mengenal Bagian-Bagian dalam Pohon Logika

Di dalam ingatan komputer, program Anda disusun seperti dahan pohon keluarga yang rapi:

1. **Pernyataan Tindakan**: Perintah yang melakukan sesuatu tetapi tidak menghasilkan nilai balasan langsung.
   - Membuat wadah nilai baru (`buat nama = "Budi"`).
   - Mengambil keputusan (`jika nilai > 75 ...`).
   - Melakukan pengulangan (`untuk orang dalam sahabat ...`).
2. **Ungkapan Nilai**: Bagian yang menghasilkan sebuah jawaban atau angka.
   - Perhitungan hitungan (`10 + 20`).
   - Pengambilan keputusan ya/tidak (`umur >= 17`).
   - Pemanggilan tindakan (`cetak("Halo")`).

---

## 3. Menyimpan dan Menjalankan Berkas Logika (.ast)

Selain berkas teks biasa (`.ce`), Celang dapat mengubah program Anda langsung menjadi berkas data logika berformat JSON (`.ast`):

```bash
# Mengubah berkas kode menjadi berkas data logika:
./celang compile program.ce -o program.ast
```

Berkas `.ast` yang dihasilkan berisi catatan struktur logika murni. Anda dapat menjalankannya langsung di terminal:

```bash
./celang run program.ast
```

Karena formatnya adalah data biasa, alat bantu lain seperti editor kode, pemeriksa otomatis, atau kecerdasan buatan dapat membaca dan memahami program Anda dengan sangat mudah tanpa harus meraba-raba teks mentah dari awal.
