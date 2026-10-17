# Celang Language Support untuk Visual Studio Code

Ekstensi resmi Visual Studio Code untuk bahasa pemrograman **Celang** — bahasa berkonsep *Human-Centric*, *Universal Mode*, *Transaksi ACID Memory*, *Variabel Reaktif*, *Time-Travel Variables*, dan *Dual Runtime (Bytecode VM & AST)*.

---

## ✨ Fitur Unggulan

- **🎨 Syntax Highlighting Lengkap & Modern**:
  - **Kata Kunci Indonesia & Universal Mode**: Mewarnai sintaks alami Indonesia (`buat`, `jika`, `fungsi`, `transaksi`, `reaktif`, `pantau`) dan standar global (`"use universal"`, `let`, `if`, `def`, `transaction`, `reactive`, `watch`).
  - **Mata Uang & Satuan Budaya**: Pewarnaan khusus untuk properti nominal dan mata uang (`.rupiah`, `.usd`, `.k`, `.million`, `.billion`, `.words`, `.dollars_in_words`).
  - **Time-Travel Variables**: Penandaan properti riwayat variabel (`.sebelumnya`, `.awal`, `.histori`, `.previous`, `.initial`, `.history`).
  - **Modern Ergonomics**: Penyorotan operator pipa (`|>`), fungsi panah ringkas (`=>`), dan null-safety (`?.`, `??`).
- **💡 Cuplikan Cepat (Code Snippets)**:
  - Ketik `universal` $\rightarrow$ langsung memasukkan direktif `"use universal"`.
  - Ketik `transaksi` atau `transaction` $\rightarrow$ kerangka transaksi lengkap dengan `batalkan`/`rollback` dan `pulihkan`/`recover`.
  - Ketik `reaktif` atau `reactive` $\rightarrow$ deklarasi variabel reaktif otomatis.
  - Ketik `pantau` atau `watch` $\rightarrow$ blok pengamat perubahan nilai variabel target.
  - Ketik `saluran` atau `channel` $\rightarrow$ deklarasi komunikasi saluran multi-pekerjaan.
  - Ketik `buat`, `let`, `fungsi`, `def`, `bentuk`, `jika`, `if`, `sambil`, `bersamaan`, `rupiah`, `usd`, atau `pipe`.
- **🧱 Auto-Closing & Smart Indentation**: Pasangan kurung kurawal `{}`, kurung siku `[]`, tanda kurung `()`, serta tanda kutip otomatis tertutup rapi dengan lekukan (*indentation*) cerdas.
- **📜 Komentar Cepat**: Pintasan `Ctrl + /` untuk menyisipkan komentar satu baris `//` atau blok `/* ... */`.

---

## 🚀 Cara Instalasi

### Cara 1: Tautan Lokal (Langsung Aktif di VS Code)

Cukup buat tautan simbolik (*symlink*) direktori ekstensi ini ke folder ekstensi VS Code Anda:

```bash
# Untuk Linux / macOS:
ln -s "$(pwd)/editors/vscode" ~/.vscode/extensions/celang

# Untuk Windows (PowerShell Administrator):
# New-Item -ItemType SymbolicLink -Path "$HOME\.vscode\extensions\celang" -Target "$PWD\editors\vscode"
```

Setelah itu, muat ulang VS Code (`Ctrl+Shift+P` -> *Developer: Reload Window*). Buka berkas berekstensi `.ce`, dan sintaks Celang akan langsung berwarna indah!

### Cara 2: Pasang Paket Berkas `.vsix`

Ekstensi ini menyediakan bundel biner `.vsix` siap pasang.
Di VS Code:
1. Buka menu Extensions (`Ctrl+Shift+X`).
2. Klik tombol opsi titik tiga `...` di sudut kanan atas panel Extensions.
3. Pilih **Install from VSIX...**.
4. Arahkan ke berkas `editors/vscode/celang-1.1.0.vsix`.
