---
title: Pengujian Unit Otomatis (celang uji)
description: Panduan lengkap pengujian unit bawaan dan assertion empati di Celang tanpa dependensi luar.
---

# Pengujian Unit Otomatis (`celang uji` / `celang test`)

> *"Kode yang tidak diuji adalah kode yang belum tentu berfungsi di masa depan."*
> 
> Celang percaya bahwa pengujian perangkat lunak bukanlah beban tambahan atau pustaka pihak ketiga yang rumit disiapkan. Pengujian adalah pilar keandalan utama (*reliability*) yang tertanam langsung di inti bahasa (*first-class citizen*).

Dengan **Celang Test Runner**, Anda tidak memerlukan npm, pip, go test wrapper, atau file konfigurasi JSON rumit. Cukup ketik `celang uji` di terminal Anda!

---

## 🌟 Mengapa `celang uji` Begitu Istimewa?

1. **Nol Konfigurasi & Dependensi Luar**: Langsung jalan dari biner tunggal Celang di Linux, macOS, maupun Windows.
2. **Dua Gaya Sintaks Elegan**: Bebas memilih antara gaya deklaratif blok (`uji "..." { }`) atau gaya callback (`uji("...", () => { })`).
3. **Pesan Galat Penuh Empati**: Jika pengujian gagal, terminal menampilkan perbedaan (*diff*) yang jelas antara ekspektasi dan nilai aktual, lengkap dengan baris dan kolom kode yang bermasalah.
4. **Siap CI/CD & Pipeline**: Otomatis mengembalikan *exit code* `0` jika seluruh tes lulus dan *exit code* `1` jika terdapat kegagalan, sangat cocok untuk GitHub Actions, GitLab CI, atau Docker.
5. **Dukungan Dwibahasa Penuh**: Seluruh asersi dapat ditulis dalam Bahasa Indonesia alami maupun standar internasional (Inggris).

---

## 🚀 Dua Gaya Menulis Pengujian

### 1. Gaya Blok Pernyataan (`uji "nama tes" { ... }`)
Gaya paling bersih, deklaratif, dan alami di Celang. Setiap blok pengujian memiliki isolasi cakupan (*scope*) lokal sendiri.

```celang
fungsi hitung_diskon(total, persen) {
    kembalikan total - (total * persen / 100)
}

uji "perhitungan diskon promo 20%" {
    buat hargaAwal = 100000
    buat hargaAkhir = hitung_diskon(hargaAwal, 20)
    
    harus_sama(hargaAkhir, 80000, "Harga akhir harus terpotong 20%")
    harus_beda(hargaAkhir, hargaAwal)
    harus_benar(hargaAkhir > 0)
}
```

### 2. Gaya Callback Fungsi (`uji("nama tes", () => { ... })`)
Gaya fungsional modern yang akrab bagi pengembang dari ekosistem Jest, Vitest, atau Mocha.

```celang
uji("diskon member VIP", () => {
    buat memberVIP = benar
    buat total = 200000
    buat diskon = memberVIP ? 50000 : 0
    
    harus_sama(total - diskon, 150000)
})

// Dalam mode Bahasa Inggris (Universal)
test("vip member discount calculation", fn() {
    let isVip = true
    let total = 200000
    let discount = isVip ? 50000 : 0
    
    assert_eq(total - discount, 150000)
})
```

---

## 📋 Pustaka Asersi (*Assertion Library*) Bawaan

Celang menyediakan pustaka asersi komprehensif yang siap digunakan langsung tanpa perlu mengimpor modul apapun:

| Bahasa Indonesia | Standar Universal (Inggris) | Penjelasan & Contoh |
| :--- | :--- | :--- |
| `harus(kondisi, [pesan])` | `assert(cond, [msg])` | Memastikan kondisi bernilai benar (*truthy*). |
| `harus_sama(aktual, target, [pesan])` | `assert_eq(act, exp, [msg])` | Memastikan kesetaraan mendalam (*deep equality*) antar nilai, larik, atau peta. |
| `harus_beda(aktual, target, [pesan])` | `assert_ne(act, exp, [msg])` | Memastikan kedua nilai bernilai tidak sama. |
| `harus_benar(kondisi, [pesan])` | `assert_true(cond, [msg])` | Memastikan nilai adalah boolean `benar` (`true`). |
| `harus_salah(kondisi, [pesan])` | `assert_false(cond, [msg])` | Memastikan nilai adalah boolean `salah` (`false`). |
| `harus_kosong(aktual, [pesan])` | `assert_null(act, [msg])` | Memastikan nilai bernilai `kosong` (`null`). |
| `harus_tidak_kosong(aktual, [pesan])` | `assert_not_null(act, [msg])` | Memastikan nilai tidak `kosong` (memiliki isi). |
| `harus_berisi(wadah, elemen, [pesan])` | `assert_contains(cnt, el, [msg])` | Memeriksa apakah teks, larik, atau peta memuat elemen tertentu. |
| `harus_galat(fungsi, [pesan_galat])` | `assert_throws(fn, [msg])` | Memastikan pemanggilan fungsi melempar galat (*throw/error*). |

---

## 🔍 Menguji Penanganan Galat (`harus_galat`)

Salah satu bagian terpenting dari kode yang andal adalah memastikan program melempar galat yang tepat saat menerima masukan tidak valid. `harus_galat` mempermudah pengujian ini:

```celang
fungsi transfer_dana(saldo, jumlah) {
    jika jumlah <= 0 {
        lempar "nominal transfer harus lebih besar dari nol"
    }
    jika jumlah > saldo {
        lempar "saldo tidak mencukupi untuk transfer"
    }
    kembalikan saldo - jumlah
}

uji "validasi transfer dana gagal" {
    buat saldoSekarang = 50000

    // 1. Memastikan galat dilempar jika nominal negatif
    harus_galat(() => {
        transfer_dana(saldoSekarang, -1000)
    }, "lebih besar dari nol")

    // 2. Memastikan galat dilempar jika saldo kurang
    harus_galat(() => {
        transfer_dana(saldoSekarang, 100000)
    }, "saldo tidak mencukupi")
}
```

---

## 💻 Menjalankan Pengujian di Terminal

### 1. Pencarian Otomatis Seluruh Proyek
Cukup jalankan perintah tanpa argumen:

```bash
celang uji
# atau
celang test
```

Runner akan memindai folder proyek dan secara otomatis mendeteksi berkas dengan konvensi:
- Berakhiran `*_uji.ce` atau `*_test.ce` (contoh: `kalkulator_uji.ce`, `auth_test.ce`)
- Berakhiran `*.uji.ce` atau `*.test.ce` (contoh: `transaksi.uji.ce`)
- Seluruh berkas `.ce` yang berada di dalam folder `uji/` atau `tests/`

### 2. Menguji Berkas Tertentu
```bash
celang uji uji/matematika_uji.ce
```

### 3. Menyaring Berdasarkan Nama Tes (`-f` / `--filter` / `--saring`)
Hanya menjalankan pengujian yang namanya mengandung kata kunci tertentu:

```bash
celang uji -f "diskon"
# atau
celang test --filter="transfer"
```

### 4. Mode Verbose & Tanpa Warna
```bash
# Menampilkan rincian setiap pengujian
celang uji -v

# Mode monokrom (tanpa ANSI escape color, cocok untuk log server)
celang uji --no-color
```

---

## 📊 Contoh Output Terminal

### Tampilan Saat Pengujian Lulus (Sukses)

```text
🔍 Menjalankan Celang Test Runner (2 berkas)

✓ uji/aritmatika_uji.ce (1.20ms)
  ✓ operasi penjumlahan dua angka (180.2µs)
  ✓ pembagian desimal valid (25.1µs)
  ✓ validasi galat pembagian nol (42.0µs)

✓ uji/koleksi_uji.ce (850.4µs)
  ✓ pencarian elemen dalam larik (30.1µs)
  ✓ filter teks dan kapitalisasi (22.8µs)

----------------------------------------------------------------------
✓ SUKSES: 5 lulus dari 5 pengujian (2 berkas) [2.15ms]
----------------------------------------------------------------------
```

### Tampilan Saat Terjadi Kegagalan (Empathetic Diff)

Jika terdapat asersi yang tidak sesuai, runner akan segera memberikan konteks yang jelas:

```text
🔍 Menjalankan Celang Test Runner (1 berkas)

✕ uji/keuangan_uji.ce (1.45ms)
  ✓ cek saldo awal akun (110.0µs)
  ✕ kalkulasi bunga tabungan (320.5µs)
    baris 28, kolom 5: Nilai aktual tidak sama dengan ekspektasi
    Ekspektasi: 1050000
    Aktual:     1025000

----------------------------------------------------------------------
✕ GAGAL: 1 gagal, 1 lulus dari 2 pengujian (1 berkas) [1.52ms]
----------------------------------------------------------------------
```

---

## 🛠 Integrasi ke CI/CD (GitHub Actions)

Karena `celang uji` mengembalikan status *exit code* standar sistem operasi (`0` untuk sukses, `1` untuk gagal), Anda dapat langsung memasukkannya ke dalam workflow otomatis:

```yaml
name: Celang Automated Test
on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Jalankan Unit Test Celang
        run: |
          ./celang uji
```

Dengan fitur pengujian bawaan ini, Celang menjamin setiap baris kode Anda siap dipakai di dunia nyata dengan keyakinan penuh!
