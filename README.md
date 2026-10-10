# Celang - Syntax Styling & Human-Centric Language Core

Celang adalah bahasa pemrograman eksperimental berbasis **Syntax Styling (Sintaks Dinamis)** dan pendekatan **Human-Centric**. Semua kode disimpan dalam format data Universal AST (U-AST), bukan teks statis. Teks kode hanyalah "gaya visual/styling" yang dapat diubah sesuai preferensi programmer.

## 🚀 Implementasi Fase 1: Fondasi Arsitektur Inti

### 1.1 Universal AST (U-AST) (`pkg/uast`)
- **Struktur Node, Statement, Expression, dan Block** dalam interface Go yang bersih.
- **Dukungan Metadata & Serialisasi JSON Dua Arah**: U-AST dapat disimpan, dibaca, dan diserialisasi ke JSON tanpa kehilangan detail tipe konkret.
- Node yang didukung:
  - `Literal` (`int`, `float`, `string`, `bool`, `null`)
  - `Identifier`, `BinaryExpression`, `UnaryExpression`, `CallExpression`, `MemberExpression`, `ArrayLiteral`
  - `VariableDeclaration`, `AssignmentStatement`, `BlockStatement`, `IfStatement`, `WhileStatement`, `ForStatement`, `FunctionDeclaration`, `ReturnStatement`, `ConcurrencyStatement`

### 1.2 Core Evaluator / Interpreter (`pkg/evaluator`)
- **Tree-Walking Interpreter** di Go.
- **Tipe Data Internal**: `INTEGER`, `FLOAT`, `BOOLEAN`, `STRING`, `NULL`, `ARRAY`, `FUNCTION`, `BUILTIN`.
- **Time-Travel Variables**: Setiap variabel dalam `Environment` menyimpan timeline riwayat nilainya (`History []Object`). Properti bawaan seperti `.sebelumnya`, `.awal`, dan `.histori` dapat langsung diakses.

### 1.3 Smart Fallback System (Tipe Adaptif) (`pkg/evaluator/coercion.go`)
- **Type Coercion Engine**: Mengekstrak nilai numerik dari string secara cerdas (misal: `"100" - 20` -> `80`, `"Rp 150000" - 25000` -> `125000`).
- Tetap mempertahankan penggabungan string yang natural (misal: `"Total: " + 50` -> `"Total: 50"`).

### Groundwork Fase 3: Natural Concurrency
- `ConcurrencyStatement` bawaan dengan mode:
  - `sambil` (eksekusi di Goroutine latar belakang).
  - `bersamaan` (eksekusi paralel terkoordinasi dengan `sync.WaitGroup`).

---

## 🛠️ Menjalankan dan Menguji

### 1. Jalankan Unit Test
```bash
go test -v ./...
```

### 2. Jalankan Program Contoh U-AST
```bash
# Contoh 1: Time-Travel Variables (.sebelumnya, .awal, .histori)
go run ./cmd/celang run examples/timetravel.ast

# Contoh 2: Smart Fallback Adaptive Coercion
go run ./cmd/celang run examples/smart_fallback.ast
```
