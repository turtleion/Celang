---
outline: deep
---

# Pustaka `fs` (`berkas`)

> [!NOTE]
> **Modul `fs`** (*file system*) adalah alias internasional resmi untuk modul **`berkas`** di Celang. Seluruh fungsi, parameter, dan perilakunya 100% identik.
> 
> Dokumentasi lengkap dan panduan menyeluruh dapat dibaca di:
> 👉 **[Pustaka Pengelola Berkas (`berkas` / `fs`)](/pustaka-standar/berkas)**

## Contoh Penggunaan Cepat

```celang
// 1. Baca, Tulis, & Operasi Atomik:
fs.write("./output.txt", "Halo dari Celang!")
buat konten = fs.read("./output.txt")
cetak(konten)

// 2. Aliran Streaming CSP (Hemat Memori):
untuk baris dalam fs.stream("./server.log") {
    jika baris.contains("ERROR") {
        cetak("Galat ditemukan:", baris)
    }
}

// 3. Pencarian Berkas Wildcard (Glob):
buat daftarBerkas = fs.glob("./src/**/*.ce")
cetak("Total berkas kode:", len(daftarBerkas))

// 4. Pemuatan & Penyimpanan Cerdas JSON:
buat config = fs.load("./setelan.json")
config.mode = "produksi"
fs.save("./setelan.json", config)

// 5. Pemantau Perubahan Berkas (Watcher):
buat watcher = fs.watch("./src")
untuk event dalam watcher {
    cetak("Perubahan:", event.path, event.action)
}

// 6. Baca & Tulis Data Biner Mentah (Bytes):
buat fotoBytes = fs.read_bytes("./foto.png")
fs.write_bytes("./foto_salinan.png", fotoBytes)
```

Untuk tabel fungsi lengkap dan penjelasan detail setiap fungsi, silakan buka **[Dokumentasi Lengkap Berkas](/pustaka-standar/berkas)**.
