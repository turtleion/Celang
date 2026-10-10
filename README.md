# Celang - Syntax Styling & Human-Centric Language Core

Celang adalah bahasa pemrograman eksperimental berbasis **Syntax Styling (Sintaks Dinamis)** dan pendekatan **Human-Centric**. Semua kode disimpan dalam format data Universal AST (U-AST), bukan teks statis. Teks kode hanyalah "gaya visual/styling" yang dapat diubah sesuai preferensi programmer.

---

## 🚀 Fitur & Komponen yang Telah Diimplementasikan

### 1. Fondasi Arsitektur Inti (Fase 1)
- **Universal AST (U-AST)** (`pkg/uast`):
  - Interface `Node`, `Statement`, `Expression`, serta serialisasi polimorfik JSON lossless dua arah.
  - Node: `Literal`, `Identifier`, `BinaryExpression`, `UnaryExpression`, `CallExpression`, `MemberExpression`, `ArrayLiteral`, `MapLiteral`, `VariableDeclaration`, `AssignmentStatement`, `BlockStatement`, `IfStatement`, `WhileStatement`, `ForStatement`, `ForEachStatement`, `FunctionDeclaration`, `ReturnStatement`, `ConcurrencyStatement`.
- **Core Tree-Walking Evaluator** (`pkg/evaluator`):
  - Model objek: `Integer`, `Float`, `String`, `Boolean`, `Null`, `Array`, `Map`, `Function`, `Builtin`, `Error`.
  - **Time-Travel Variables**: Pelacakan riwayat mutasi variabel (`.sebelumnya`, `.awal`, `.histori`).
  - **Smart Fallback System**: Type coercion cerdas berbasis regex (misal: `"Rp 150000" - 25000` -> `125000`).
  - **Natural Concurrency**: `sambil` (goroutine latar belakang) dan `bersamaan` (sinkronisasi paralel).

### 2. Mesin "Syntax Styling" & Two-Way Parser (Fase 2)
- **Style Profiles** (`pkg/syntax/profile.go`):
  - **Bahasa Indonesia (Default)**: `buat`, `konstanta`, `jika ... selain jika ... selain`, `selama`, `untuk ... dalam`, `fungsi`, `kembalikan`, `sambil`, `bersamaan`.
  - **Pythonic Style**: `x = val`, `def`, `return`, `if ... elif ... else`, `for ... in`, `async`.
  - **C / JavaScript Style**: `let`, `const`, `fn`, `return`, `if ... else if ... else`.
- **Forward Parser (Teks -> U-AST)** (`pkg/parser`):
  - Lexer dan recursive-descent/Pratt parser yang mengubah teks kode Celang menjadi struktur data U-AST.
- **Reverse Parser / Formatter (U-AST -> Teks)** (`pkg/syntax/formatter.go`):
  - Mengubah U-AST kembali menjadi teks kode dalam profil gaya mana pun secara instan (*on-the-fly syntax restyling*).

---

## 🛠️ Panduan Penggunaan CLI

### 1. Menjalankan Unit Test
```bash
go test -v ./...
```

### 2. Menjalankan Berkas Kode Teks Celang (`.ce`)
```bash
go run ./cmd/celang run examples/program_indo.ce
```

### 3. Mengompilasi Kode Teks ke Berkas Data U-AST (`.ast`)
```bash
go run ./cmd/celang compile examples/program_indo.ce -o examples/program_indo.ast
```

### 4. Menjalankan Berkas U-AST Langsung
```bash
go run ./cmd/celang run examples/program_indo.ast
```

### 5. Syntax Styling: Mengubah Tampilan Kode Antar Gaya
```bash
# Render kode dalam gaya Pythonic:
go run ./cmd/celang format examples/program_indo.ce --style=pythonic

# Render kode dalam gaya C / JavaScript:
go run ./cmd/celang format examples/program_indo.ce --style=classic

# Render kembali ke gaya Bahasa Indonesia:
go run ./cmd/celang format examples/program_indo.ast --style=indonesian
```
