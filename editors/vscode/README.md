# Celang Language Support untuk Visual Studio Code

Ekstensi resmi Visual Studio Code untuk bahasa pemrograman **Celang** — bahasa berkonsep *Human-Centric*, *Syntax-Styling*, *Time-Travel Variables*, dan *Dual Runtime (Bytecode VM & AST)*.

---

## ✨ Fitur Unggulan

- **🎨 Syntax Highlighting Lengkap**: Mewarnai kata kunci Bahasa Indonesia (`buat`, `jika`, `selama`, `fungsi`, `sambil`, `bersamaan`), dialek alternatif (`let`, `def`, `async`), properti time-travel (`.sebelumnya`, `.awal`), serta operator pipa (`|>`).
- **💡 Cuplikan Cepat (Code Snippets)**: Ketik `buat`, `fungsi`, `bentuk`, `jika`, `sambil`, `bersamaan`, atau `gaya` dan tekan `Tab` untuk menghasilkan kerangka kode otomatis.
- **🧱 Auto-Closing & Bracket Matching**: Pasangan kurung kurawal `{}`, kurung siku `[]`, tanda kurung `()`, serta tanda kutip otomatis tertutup rapi.
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

### Cara 2: Paket Berkas `.vsix`

Jika Anda memiliki `vsce` terinstal (`npm install -g @vscode/vsce`):

```bash
cd editors/vscode
vsce package
```

Lalu pasang berkas `.vsix` yang dihasilkan via VS Code:
- Klik menu Extensions -> tanda titik tiga `...` di pojok kanan atas -> **Install from VSIX...**.
