# Desain Optimasi Memori & Arsitektur Bytecode VM (`PERFORMANCE_AND_MEMORY.md`)

> **Catatan untuk Kamu:**  
> Dokumen ini membahas dua topik besar:  
> 1. **Filosofi Trade-off** antara Kepraktisan/Keamanan vs Kecepatan Mesin dalam perancangan bahasa.  
> 2. **Implementasi Fix Memory Leak** pada *Time-Travel Variables* (menggunakan *Anchor Snapshot + Ring Buffer*).  
> 3. **Rancangan & Cara Kerja Bytecode Virtual Machine (VM)** untuk melipatgandakan kecepatan eksekusi Celang di masa depan.

---

## 1. ⚖️ Apakah Benar Trade-off Kepraktisan & Keamanan Adalah Kecepatan?

**Jawabannya: Ya, secara fundamental itu adalah hukum alam di dunia *Computer Science*.**

Dalam perancangan bahasa pemrograman (*Language Design*), ada konsep yang dikenal sebagai **Segitiga Trade-off**:

```
                 Keamanan & Kenyamanan
               (Safety & Ergonomics)
                     ▲
                    / \
                   /   \
                  /     \
                 /       \
  Kecepatan Mesin ◀───────▶ Kesederhanaan Pembuatan
   (Raw Speed)               (Simplicity)
```

### Mengapa Kepraktisan Mengorbankan Kecepatan?
1. **Garbage Collection (GC) vs Manual Memory:**
   - **Bahasa Praktis (Python, JS, Go, Celang):** Memori dibersihkan otomatis. Programmer tidak perlu pusing `free()` atau takut crash. *Trade-off:* Komputer harus meluangkan siklus CPU untuk memindai dan membersihkan RAM.
   - **Bahasa C / C++:** Super cepat karena programmer mengatur alamat RAM manual, tapi rawan *memory leak*, *buffer overflow*, dan celah keamanan siber (hacker bisa mengeksploitasi pointer liar).
2. **Dynamic Typing & Smart Fallback:**
   - Celang punya fitur humanis: `"Rp 50000" - 10000` otomatis menghasilkan `40000`. Ini sangat praktis untuk manusia, tetapi di balik layar mesin harus melakukan regex/parsing teks setiap operasi.
3. **Time-Travel Variables:**
   - Sangat praktis untuk debugging (`x.sebelumnya`, `x.awal`), tapi mesin harus menyisihkan sedikit memori untuk mencatat jejak masa lalu.

### Apakah Trade-off Ini "Takdir Mati"?
**Tidak! Di sinilah peran rekayasa perangkat lunak tingkat lanjut (*Engineering Optimization*):**
- **Generasi 1 (AST Interpreter):** Paling lambat, tapi paling fleksibel dan mudah dibuat (posisi Celang saat ini).
- **Generasi 2 (Bytecode VM):** **5x s/d 20x lebih cepat.** Mengubah pohon AST menjadi baris instruksi biner linear (seperti Python, Lua, Java).
- **Generasi 3 (JIT - Just-In-Time Compiler):** Mendekati kecepatan C/Go. Mengompilasi kode yang sering dijalankan (*hot loops*) langsung ke bahasa mesin CPU saat program berjalan (seperti V8 di Node.js atau PyPy).

---

## 2. 🛡️ Bagian 1: Solusi Memory Leak pada Time-Travel Variables

### 2.1 Masalah Saat Ini (*Unbounded Slice Growth*)
Saat ini di `pkg/evaluator/environment.go`:
```go
entry.History = append(entry.History, val)
```
Setiap kali variabel ditugaskan ulang (`counter = counter + 1`), nilai lama dimasukkan ke slice `History` tanpa batas.

Jika ada perulangan 10 juta kali:
```celang
buat i = 0
selama i < 10000000 {
    i = i + 1
}
```
Slice `i.History` akan menampung **10.000.000 objek di RAM** ($\approx$ ratusan Megabyte hingga Gigabyte), dan Garbage Collector Go tidak bisa membersihkannya karena referensinya masih menempel di variabel `i`.

---

### 2.2 Desain Solusi: *Anchor Snapshot + Bounded Ring Buffer*

Alih-alih menyimpan array tanpa batas, kita rancang struktur data cerdas untuk `VariableEntry`:

```go
type VariableEntry struct {
    Name       string
    InitialVal Object       // Pointer ke nilai PERTAMA kali (abadi untuk .awal)
    PreviousVal Object      // Pointer ke nilai TEPAT SEBELUM saat ini (untuk .sebelumnya)
    CurrentVal  Object      // Pointer ke nilai saat ini
    
    // Ring buffer terbatas untuk properti .histori (maks misal 50 state terakhir)
    HistoryRing []Object
    HistoryCap  int          // Default: 50
    HistoryHead int
    IsFull      bool
}
```

```
[ Deklarasi Awal: x = 10 ]
InitialVal  ──▶ 10   (Tersimpan abadi di O(1))
CurrentVal  ──▶ 10
PreviousVal ──▶ nil
HistoryRing ──▶ [10]

[ Diubah: x = 20 ]
PreviousVal ──▶ 10
CurrentVal  ──▶ 20
HistoryRing ──▶ [10, 20]

... setelah 10.000.000 kali loop ...
InitialVal  ──▶ 10          (Tetap 10!)
PreviousVal ──▶ 9999999     (Tetap nilai tepat sebelum selesai!)
CurrentVal  ──▶ 10000000
HistoryRing ──▶ [50 nilai terakhir saja]
```

### 2.3 Keuntungan Arsitektur Ini:
1. `x.awal`: **100% selalu tepat**, tanpa memakan memori berulang.
2. `x.sebelumnya`: **100% selalu tepat**, hanya memakan 1 pointer memori.
3. `x.histori`: Mengembalikan snapshot riwayat terkini (50–100 langkah terakhir), bukan jutaan langkah usang.
4. **Kompleksitas Memori:** Turun dari $O(N)$ (tak terbatas) menjadi **$O(1)$ (konstan)**. Loop 100 juta kali sekalipun memorinya tetap dingin dan hemat (< 1 KB per variabel)!

---

## 3. ⚡ Bagian 2: Cara Kerja Bytecode Virtual Machine (VM)

Untuk melipatgandakan performa Celang hingga mendekati kecepatan Lua/Python, kita bisa menambahkan backend eksekusi berupa **Bytecode VM**.

### 3.1 Mengapa AST Interpreter Terasa Berat?
Di AST Interpreter saat ini:
- Untuk menghitung `1 + 2 * 3`, evaluator harus membuat node pointer Go:
  `BinaryExpression(BinaryExpression(2, "*", 3), "+", 1)`
- Saat runtime menelusuri pohon ini (*tree-walking*), CPU harus sering melompat ke alamat memori yang terpencar-pencar (*pointer indirection & cache miss*).

### 3.2 Apa Itu Bytecode?
**Bytecode adalah deretan instruksi biner angka satu dimensi** yang sangat padat, mirip bahasa mesin (assembly), tetapi berjalan di atas mesin virtual (*software processor*).

Contoh: Ekspresi `a = 10 + 20` diubah menjadi instruksi biner sederhana:
```text
Bytecode Offset | OpCode           | Parameter / Operand
-----------------------------------------------------------
0000            | OP_CONSTANT      | Index 0 (Konstanta 10)
0002            | OP_CONSTANT      | Index 1 (Konstanta 20)
0004            | OP_ADD           | -
0005            | OP_SET_GLOBAL    | Index 2 ("a")
```

---

### 3.3 Anatomi Mesin Celang VM

Celang VM dirancang menggunakan model **Stack-Based VM** (model yang sama dengan JVM Java, Python VM, dan WebAssembly):

```
       [ Universal AST (U-AST) ]
                   │
                   ▼
       [ Celang Bytecode Compiler ]  ──▶ Mengubah AST jadi Chunk Bytecode
                   │
                   ▼
       [ Celang Virtual Machine ]
         ┌─────────────────────┐
         │ Instruction Pointer │ ──▶ Menunjuk baris bytecode saat ini
         ├─────────────────────┤
         │ Value Stack         │ ──▶ [ 10 ] [ 20 ] ──▶ (OP_ADD) ──▶ [ 30 ]
         ├─────────────────────┤
         │ Global / Local Env  │ ──▶ Tempat penyimpanan variabel
         └─────────────────────┘
```

#### Loop Eksekusi (*Dispatch Loop*):
Di dalam kode Go, VM hanyalah sebuah perulangan super cepat:
```go
for {
    switch readOpcode() {
    case OP_CONSTANT:
        push(readConstant())
    case OP_ADD:
        b := pop()
        a := pop()
        push(a + b)
    case OP_PRINT:
        fmt.Println(pop())
    case OP_RETURN:
        return
    }
}
```
Karena eksekusi hanya membaca array biner linear di memori, CPU komputer bisa memanfaatkannya dengan **Cache L1/L2 super cepat** tanpa ada overhead penelusuran pohon AST!

---

### 3.4 Hubungan U-AST dan Bytecode VM di Celang

Apakah U-AST kita menjadi sia-sia jika ada Bytecode VM? **Sama sekali tidak!**  
U-AST tetap menjadi **jantung utama ekosistem Celang**:

1. **U-AST:** Menangani *Syntax Styling* (Indonesian, JS, Pythonic, Sunda), Formatter, dan Linter.
2. **Bytecode VM:** Menjadi *Execution Engine* berkecepatan tinggi yang menerima input U-AST dan menerjemahkannya ke biner VM.

```
Kode (.ce) ──▶ Parser ──▶ U-AST ──┬──▶ Formatter / Stylist (Human-Friendly)
                                   └──▶ Bytecode Compiler ──▶ Celang VM (Super Fast)
```

---

## 4. 🚀 Rencana Eksekusi Bertahap

1. **Fase 1 (Saat ini):** Terapkan perbaikan **Anchor Snapshot + Bounded Ring Buffer** pada `VariableEntry` di branch ini. Ini langsung mematikan 100% potensi memory leak pada time-travel.
2. **Fase 2 (Berikutnya):** Tambahkan benchmark test (`go test -bench`) untuk mengukur efisiensi memori loop panjang.
3. **Fase 3 (Masa Depan):** Bangun prototipe paket `pkg/vm` dan `pkg/compiler` sebagai opsi eksekusi alternatif berkecepatan tinggi (`celang run --vm app.ce`).
