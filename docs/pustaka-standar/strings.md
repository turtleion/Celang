---
outline: deep
---

# Pustaka `strings` (`teks`)

> [!NOTE]
> **Modul `strings`** adalah alias internasional resmi untuk modul **`teks`** di Celang. Seluruh fungsi, parameter, dan perilakunya 100% identik.
> 
> Dokumentasi lengkap dan panduan menyeluruh dapat dibaca di:
> 👉 **[Pustaka Pengolah Teks (`teks` / `strings`)](/pustaka-standar/teks)**

## Contoh Penggunaan Cepat

```celang
// 1. Fitur Cerdas:
buat tautan = strings.slug("Panduan Belajar Celang 2.0!")
buat disensor = strings.mask("budi@gmail.com")
buat mirip = strings.similarity("jakarta", "jakartaa")
buat terdekat = strings.closest("aple", ["apel", "jeruk", "mangga"])

// 2. Format Gaya Penulisan (Case Styles):
cetak(strings.camel("nama_pengguna")) // "namaPengguna"
cetak(strings.snake("namaPengguna")) // "nama_pengguna"
cetak(strings.kebab("Halo Dunia"))    // "halo-dunia"

// 3. Pola & Pencarian Reguler (Regex):
jika strings.match("user@mail.com", "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$") {
    cetak("Email valid!")
}
```

Untuk tabel fungsi lengkap dan penjelasan detail setiap fungsi, silakan buka **[Dokumentasi Lengkap Teks](/pustaka-standar/teks)**.
