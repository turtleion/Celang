# Multi-Dialek (ID, Sunda, JS, Pythonic)

Salah satu inovasi terbesar Celang adalah pemisahan antara **Struktur Semantik Program (Universal AST)** dan **Tampilan Sintaks (Syntax Style)**. 

Di bahasa pemrograman konvensional, sintaks bersifat kaku. Di Celang, Anda dapat menulis logika yang sama persis menggunakan kata kunci Bahasa Indonesia, gaya JavaScript, gaya Python, atau bahkan bahasa daerah seperti Bahasa Sunda.

---

## 1. Dialek Bawaan (*Built-in Profiles*)

Celang menyediakan profil bawaan siap pakai:

### A. Dialek Bahasa Indonesia (`indonesian` - Default)
Didesain agar mudah dipahami secara alami oleh siapa saja yang berbahasa Indonesia:
```celang
buat nama = "Budi"
jika nama == "Budi" {
    cetak("Halo Budi!")
} selain {
    cetak("Halo tamu!")
}
```

### B. Dialek JavaScript (`javascript`)
Bagi pengembang web atau pengguna Node.js yang sudah terbiasa dengan sintaks JS modern:
```celang
let nama = "Budi"
if (nama === "Budi") {
    print("Halo Budi!")
} else {
    print("Halo tamu!")
}
```

### C. Dialek Pythonic (`pythonic`)
Bagi pengembang yang menyukai gaya penulisan Python:
```celang
nama = "Budi"
if nama == "Budi" {
    print("Halo Budi!")
} else {
    print("Halo tamu!")
}
```

### D. Dialek Bahasa Daerah (Contoh: `sunda`)
Dapat dimuat dari berkas kustom seperti `sunda.style`:
```celang
jieun ngaran = "Budi"
lamun ngaran == "Budi" {
    cetak("Sampurasun Budi!")
} lamun henteu {
    cetak("Sampurasun semah!")
}
```

---

## 2. Tabel Perbandingan Kata Kunci

| Konsep | Indonesia (`indonesian`) | JavaScript (`javascript`) | Pythonic (`pythonic`) | Sunda (`sunda.style`) |
| :--- | :--- | :--- | :--- | :--- |
| Deklarasi Variabel | `buat` | `let` / `var` | *(langsung)* | `jieun` |
| Konstanta | `konstanta` | `const` | `const` | `tetep` |
| Percabangan | `jika` | `if` | `if` | `lamun` |
| Kondisi Kedua | `selain jika` | `else if` | `elif` | `atawa lamun` |
| Kondisi Terakhir | `selain` | `else` | `else` | `lamun henteu` |
| Perulangan While | `selama` | `while` | `while` | `salami` |
| Perulangan For | `untuk ... dalam` | `for ... of` | `for ... in` | `pikeun ... dina` |
| Fungsi | `fungsi` | `function` | `def` | `fungsi` |
| Kembalian | `kembalikan` | `return` | `return` | `balikkeun` |
| Asinkron | `sambil` | `async` | `async` | `bari` |
| Paralel | `bersamaan` | `parallel` | `parallel` | `babarengan` |
| Nilai Boolean | `benar`, `salah` | `true`, `false` | `True`, `False` | `bener`, `salah` |
| Nilai Kosong | `kosong` | `null` | `None` | `kosong` |

---

## 3. Menjalankan Kode dengan Dialek Tertentu

Anda dapat menginstruksikan Celang untuk mem-parsing kode dengan style tertentu melalui bendera `--style`:

```bash
# Menjalankan berkas bertema JavaScript
celang run --style javascript kode_web.ce

# Menjalankan berkas bertema Sunda dengan berkas style eksternal
celang run --style examples/sunda.style program_sunda.ce
```

---

## 4. Konversi Otomatis Antar-Gaya (`celang format`)

Karena Celang mengonversi kode ke Universal AST (U-AST), Celang dapat **menerjemahkan gaya kode secara instan dan tanpa kehilangan logika**:

```bash
# Mengubah kode bergaya JavaScript ke Bahasa Indonesia murni
celang format --to indonesian app_js.ce -w

# Mengubah kode Bahasa Indonesia ke dialek Sunda
celang format --to examples/sunda.style program.ce -o program_sunda.ce
```

Hal ini memungkinkan kolaborasi tim di mana setiap pengembang dapat membaca atau menulis dalam dialek yang paling nyaman bagi mereka!
