# 🌍 Mode Universal & Internasional (`"use universal"`)

Celang dirancang ramah untuk pengembang Indonesia dengan kata kunci alami dan konvensi lokal. Namun, ketika kode ingin dipublikasikan sebagai pustaka *open-source* global, dibagikan ke tim internasional, atau dijalankan dalam standar internasional yang ketat, Celang menyediakan **Mode Universal** (*Universal Mode*).

Cukup tambahkan direktif `"use universal"` (atau `"use english"`, `"use standard"`) di baris paling awal berkas `.ce`:

```celang
"use universal"

let balance = 100.k
const taxRate = 0.05
reactive grandTotal = balance + (balance * taxRate)

print("Total USD:", grandTotal.usd)
// Output: Total USD: $105,000.00
```

---

## 📌 Mengapa Mode Universal?

1. **Kolaborasi Lintas Batas**: Kode Celang Anda dapat dibaca dan dipelihara oleh pengembang di seluruh dunia tanpa hambatan bahasa.
2. **Ketat & Konsisten (*English Only*)**: Mode Universal menegakkan sintaks bahasa Inggris standar secara menyeluruh. Kata kunci dan pembantu Indonesia dilarang dan langsung dilaporkan dengan pesan empati berbahasa Inggris beserta saran perbaikan.
3. **Mata Uang & Skala Global**: Mengganti format lokal Rupiah (`.rupiah`, `.terbilang`, `.ribu`, `.juta`) menjadi standar internasional Dolar AS (`.usd`, `.k`, `.million`, `.billion`, `.words`, `.dollars_in_words`).
4. **Semua Fitur Modern Tetap Tersedia**: Keunggulan revolusioner Celang seperti **Time-Travel Variables**, **Transaksi & Rollback**, **Variabel Reaktif**, serta **Saluran Asinkron** tetap berfungsi penuh menggunakan kata kunci bahasa Inggris.

---

## 🔒 Penegakan Sintaks Ketat (Strict Enforcement)

Dalam Mode Universal, parser Celang menolak kata kunci Indonesia dan memberikan pesan perbaikan instan:

| Kata Kunci Indonesia Dilarang | Wajib Digunakan di Universal Mode | Catatan |
| :--- | :--- | :--- |
| `buat` | `let` | Deklarasi variabel mutabel |
| `konstanta`, `tetap` | `const` | Deklarasi variabel konstan |
| `jika`, `kalau`, `bila` | `if` | Percabangan kondisi |
| `selain jika`, `lain jika` | `else if` atau `elif` | Percabangan alternatif |
| `selain`, `lain` | `else` | Blok kondisi fallback |
| `fungsi` | `def`, `func`, atau `fn` | Deklarasi fungsi |
| `kembalikan`, `mengembalikan` | `return` | Pengembalian nilai fungsi |
| `transaksi` | `transaction` | Blok atomik ACID memory |
| `batalkan` | `rollback` | Sinyal pembatalan transaksi |
| `pulihkan` | `recover` | Penanganan kesalahan transaksi |
| `pantau` | `watch` | Pengamat mutasi reaktif |
| `reaktif` | `reactive` | Variabel turunan reaktif |
| `benar`, `salah`, `nihil` | `true`, `false`, `null` | Nilai literal standar |
| `dan`, `atau`, `bukan` | `and`, `or`, `not` (atau `&&`, `\|\|`, `!`) | Operator logika |

### Contoh Diagnostik Empati di Universal Mode

Jika tidak sengaja menulis kata kunci Indonesia dalam berkas Universal Mode:

```celang
"use universal"
buat balance = 50.k
```

Terminal akan menampilkan diagnostik dalam bahasa Inggris:

```text
╭─ 🛑 Syntax Error ──────────────────────────────────────────
│ 📍 Location: Line 2, Column 1
│ 💬 Explanation: Indonesian keyword 'buat' is not allowed in Universal Mode. Use 'let' instead.
│
│   1 | "use universal"
│   2 | buat balance = 50.k
│        ^ here
│   3 | 
│
│ 💡 Suggested Solution:
│    Universal Mode requires standard English syntax and keywords.
╰────────────────────────────────────────────────────────────
```

---

## 💵 Satuan Nominal & Mata Uang Internasional

Dalam Mode Universal, seluruh pembantu nominal dan mata uang dialihkan ke standar internasional:

### 1. Pengali Nominal (*Multipliers*)

```celang
"use universal"

let users = 50.k             // 50,000
let views = 10.thousand      // 10,000
let downloads = 2.5.million  // 2,500,000
let population = 8.billion   // 8,000,000,000
let debt = 1.trillion        // 1,000,000,000,000
```

### 2. Format Mata Uang USD (`.usd` / `format_usd`)

```celang
"use universal"

let price = 1250.usd          // "$1,250.00"
let smallFee = 25.50.usd      // "$25.50"
let formatted = format_usd(99.9) // "$99.90"
```

### 3. Angka Terbilang Bahasa Inggris (`.words` / `.dollars_in_words`)

```celang
"use universal"

print(42.words)
// Output: forty-two

print(1500.in_words)
// Output: one thousand five hundred

print(25.50.dollars_in_words)
// Output: twenty-five dollars and fifty cents

print(10.dollars_in_words)
// Output: ten dollars
```

### 4. Satuan Waktu & Durasi

```celang
"use universal"

let cacheTtl = 5.sec          // 5,000 ms
let timeout = 2.min           // 120,000 ms
let sessionDuration = 1.hour  // 3,600,000 ms
let backupInterval = 1.day    // 86,400,000 ms
```

> [!NOTE]
> Properti lokal Indonesia seperti `.rupiah`, `.terbilang`, `.ribu`, `.juta`, `.detik`, dan modul `berkas`/`matematika` **otomatis dilarang** di Universal Mode dengan petunjuk ekuivalen bahasa Inggrisnya.

---

## ⚡ Fitur Modern di Universal Mode

### 1. Transaksi & Rollback

```celang
"use universal"

let sender = 200
let receiver = 50

transaction {
    sender = sender - 300
    receiver = receiver + 300
    if sender < 0 {
        rollback "Insufficient balance"
    }
} recover (reason) {
    print("Transaction rolled back:", reason)
}

print("Sender safely preserved:", sender) // 200
```

### 2. Variabel Reaktif & Pengamat (*Watchers*)

```celang
"use universal"

let price = 100
reactive tax = price * 0.1
reactive grandTotal = price + tax

watch price (newVal, oldVal) {
    print("Price adjusted from", oldVal, "to", newVal)
}

price = 150
print("New Grand Total:", grandTotal.usd) // $165.00
```

### 3. Time-Travel Variables

```celang
"use universal"

let score = 10
score = 25
score = 50

print("Current:", score)          // 50
print("Previous:", score.previous) // 25 (atau score.prev)
print("Initial:", score.initial)   // 10
print("History:", score.history)   // [10, 25, 50]
```

### 4. Pustaka Standar Internasional

Di Universal Mode, impor modul pustaka standar menggunakan nama global:

```celang
"use universal"

import "math"
import "fs"
import "strings"
import "time"
import "os"
import "http"
```

---

## 🔄 Kompatibilitas Mundur

Jika direktif `"use universal"` **tidak dituliskan**, Celang tetap beroperasi dalam mode santai dan fleksibel (*default permissive mode*):
- Anda bebas menggunakan kata kunci Indonesia (`buat`, `jika`, `transaksi`, `50.ribu`, `.rupiah`), kata kunci bahasa Inggris (`let`, `if`), maupun campuran keduanya.
- Seluruh kode warisan (*legacy code*) tetap berjalan 100% tanpa perubahan apa pun.
