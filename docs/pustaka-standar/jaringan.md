---
outline: deep
---

# Pustaka `jaringan` (`http`)

> [!NOTE]
> **Modul `jaringan`** (dan alias `net`) adalah alias resmi untuk modul **`http`** di Celang. Seluruh fungsi, parameter, dan perilakunya 100% identik.
> 
> Dokumentasi lengkap dan panduan menyeluruh dapat dibaca di:
> 👉 **[Pustaka Jaringan & Protokol HTTP (`http` / `jaringan`)](/pustaka-standar/http)**

## Contoh Penggunaan Cepat

```celang
// 1. Klien HTTP:
buat res = jaringan.ambil("https://jsonplaceholder.typicode.com/posts/1")
cetak(res.status) // 200
cetak(res.json().title)

// 2. Server HTTP Sederhana:
buat server = jaringan.buat_server()
server.ambil("/", fungsi(req, res) {
    res.kirim("Halo dari Celang HTTP Server!")
})
server.ambil("/api/waktu", fungsi(req, res) {
    res.json({"waktu": waktu.sekarang()})
})

cetak("Menjalankan server di port 8080...")
server.layani(8080)
```

Untuk tabel fungsi lengkap dan penjelasan detail setiap fungsi, silakan buka **[Dokumentasi Lengkap HTTP](/pustaka-standar/http)**.
