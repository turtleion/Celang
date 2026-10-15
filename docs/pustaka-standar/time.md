---
outline: deep
---

# Pustaka `time` (`waktu`)

> [!NOTE]
> **Modul `time`** adalah alias internasional resmi untuk modul **`waktu`** di Celang. Seluruh fungsi, parameter, konstanta, dan perilakunya 100% identik.
> 
> Dokumentasi lengkap dan panduan menyeluruh dapat dibaca di:
> 👉 **[Pustaka Penunjuk Waktu (`waktu` / `time`)](/pustaka-standar/waktu)**

## Contoh Penggunaan Cepat

```celang
// 1. Waktu Saat Ini & Waktu Relatif:
buat sekarang = time.now()
buat relative = time.time_ago(sekarang - time.minutes(5))
cetak(relative) // "5 menit yang lalu"

// 2. Format Penanggalan:
buat formatted = time.format(sekarang, "YYYY-MM-DD HH:mm:ss")
cetak(formatted) // "2026-10-15 17:30:00"

// 3. Profiler Eksekusi Kode:
buat hasilUji = time.benchmark(fungsi() {
    buat x = 100 * 200
    kembalikan x
})
cetak("Waktu eksekusi:", hasilUji.durasi_teks)

// 4. Aritmatika Tanggal:
buat besok = time.add(sekarang, { "days": 1 })
buat selisih = time.diff(sekarang, besok)
cetak("Selisih:", selisih.hari, "hari")
```

Untuk tabel fungsi lengkap dan penjelasan detail setiap fungsi, silakan buka **[Dokumentasi Lengkap Waktu](/pustaka-standar/waktu)**.
