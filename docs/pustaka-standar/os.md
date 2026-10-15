---
outline: deep
---

# Pustaka `os` (`sistem`)

> [!NOTE]
> **Modul `os`** (*operating system*) adalah alias internasional resmi untuk modul **`sistem`** di Celang. Seluruh fungsi, parameter, konstanta, dan perilakunya 100% identik.
> 
> Dokumentasi lengkap dan panduan menyeluruh dapat dibaca di:
> 👉 **[Pustaka Sistem Operasi (`sistem` / `os`)](/pustaka-standar/sistem)**

## Contoh Penggunaan Cepat

```celang
// 1. Eksekusi Perintah Terminal (Process Execution):
let hasil = os.exec("git status")
if hasil.success {
    print(hasil.stdout)
}

// 2. Variabel Lingkungan (Environment Variables):
let port = os.getenv("PORT", "8080")
os.setenv("APP_ENV", "production")

// 3. Argumen Baris Perintah (CLI Arguments):
let args = os.args()
let target = os.args(0, "default_target")

// 4. Informasi Sistem & Platform:
print("Platform :", os.platform())
print("Arch     :", os.arch())
print("CPUs     :", os.cpu_count())
print("CWD      :", os.cwd())
print("Home     :", os.homedir())

// 5. Memori & Proses:
let mem = os.memory()
print("Allocated:", mem.alloc)
```

Untuk tabel fungsi lengkap dan penjelasan detail setiap kemampuan, silakan buka **[Dokumentasi Lengkap Sistem](/pustaka-standar/sistem)**.
