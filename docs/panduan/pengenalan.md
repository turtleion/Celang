# Pengenalan & Filosofi Celang

**Celang** adalah bahasa pemrograman eksperimental modern yang dibangun dengan fokus pada **kemanusiaan (Human-Centric)**, **fleksibilitas gaya (Syntax Styling)**, dan **kecepatan eksekusi (Bytecode Virtual Machine)**.

---

## 💡 Mengapa Celang Diciptakan?

Sebagian besar bahasa pemrograman modern memaksa manusia untuk tunduk pada satu sintaks kaku yang sering kali berbahasa Inggris. Namun bagi sebagian besar orang di dunia—termasuk di Indonesia—belajar logika komputasi dalam bahasa ibu memberikan intuisi dan pemahaman mendasar yang jauh lebih cepat.

Celang lahir untuk menjawab pertanyaan:
> *"Bagaimana jika kode program tidak terikat pada satu bahasa manusia saja, melainkan berupa pohon semantik abstrak yang bisa 'didandani' (styled) sesuai kenyamanan pembacanya?"*

---

## 🌟 Pilar Utama Celang

### 1. Human-Centric (Berpusat pada Manusia)
Kode bukan sekadar instruksi untuk mesin, melainkan alat komunikasi antarsesama manusia. Celang menyediakan:
- Diagnostik galat empatik yang ramah dan memberi saran perbaikan konkret.
- Kosakata alami yang intuitif seperti `buat`, `jika`, `selama`, `kembalikan`.
- Coercion numerik cerdas (*adaptive coercion*) yang memahami format angka dunia nyata seperti `"Rp 150000"`.

### 2. Universal AST (U-AST) & Syntax Styling
Celang memisahkan antara **makna logika program** dan **tampilan sintaks**:
- Kode yang ditulis dalam bahasa Indonesia dapat diformat secara otomatis menjadi bahasa Sunda, gaya JavaScript, atau Pythonic.
- Tim programmer yang berbeda gaya bisa bekerja sama pada basis kode yang sama tanpa perdebatan gaya kode.

### 3. Variabel Time-Travel (Penjelajah Waktu)
Di Celang, variabel bukan hanya menyimpan nilai saat ini, melainkan ingatannya sendiri:
- Akses nilai sebelumnya lewat `.sebelumnya`.
- Akses nilai pertama saat variabel lahir lewat `.awal`.
- Telusuri histori mutasi lewat `.histori`.
- Semua histori ini dioptimasi dengan arsitektur **Anchor Snapshot + Bounded Ring Buffer**, menjamin konsumsi memori selalu konstan **$O(1)$** tanpa kebocoran memori (*zero memory leak*).

### 4. Bytecode Compiler & Virtual Machine (VM)
Celang tidak lambat! Setiap program Celang dikompilasi menjadi instruksi bytecode padat dan dieksekusi oleh mesin virtual (*stack-based VM*) berkecepatan tinggi:
- **2x lebih cepat** dibanding interpreter tree-walk AST konvensional.
- **47% lebih hemat memori**.
- Terintegrasi penuh dengan seluruh pustaka standar Go.
