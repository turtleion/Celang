# Panduan Lengkap: Gaya Bahasa & Ekosistem Celang (`STYLING-AND-MORE.md`)

Dokumen ini menjelaskan konsep **Gaya Bahasa Fleksibel (Syntax Styling)** di Celang: bagaimana satu program yang sama dapat dilihat dan ditulis dalam berbagai gaya bahasa tanpa mengubah cara kerjanya.

---

## 1. Filosofi Inti: Memisahkan Ide Logika dari Tampilan Tulisan

Di kebanyakan bahasa pemrograman konvensional, teks tulisan dan ide program menyatu secara kaku. Jika anggota tim berdebat antara gaya tanda kurung kurawal vs indentasi, atau kata kunci bahasa Inggris vs bahasa Indonesia, biasanya salah satu pihak terpaksa mengalah.

Celang memisahkan keduanya:
- **Pohon Logika Murni (U-AST)**: Bagan pemahaman murni dari program Anda (kebenaran tunggal yang disimpan dan dijalankan).
- **Profil Gaya (.style)**: Pilihan tampilan visual bagaimana bagan logika tersebut ditampilkan sebagai tulisan yang dibaca oleh mata manusia.

```
                      ┌────────────────────────┐
                      │   Pohon Logika Murni   │
                      │       (U-AST)          │
                      └───────────┬────────────┘
                                  │
         ┌────────────────────────┼────────────────────────┐
         ▼                        ▼                        ▼
┌──────────────────┐    ┌──────────────────┐    ┌──────────────────┐
│ Gaya Indonesia   │    │  Gaya Pythonic   │    │  Gaya JavaScript │
│   (buat x = 10)  │    │     (x = 10)     │    │   (let x = 10)   │
└──────────────────┘    └──────────────────┘    └──────────────────┘
```

---

## 2. Mengubah Gaya Melalui Perintah Terminal (`celang format`)

Anda dapat menerjemahkan tampilan berkas kode ke gaya lain secara instan:

```bash
# 1. Menampilkan berkas dalam gaya Bahasa Indonesia:
./celang format program.ce --style=indonesian

# 2. Menampilkan berkas dalam gaya Python:
./celang format program.ce --style=pythonic

# 3. Menampilkan berkas dalam gaya JavaScript:
./celang format program.ce --style=javascript

# 4. Menampilkan berkas dalam dialek Basa Sunda kustom:
./celang format program.ce --style=examples/sunda.style

# 5. Langsung menimpa isi berkas dengan gaya baru secara permanen (opsi -w):
./celang format -w program.ce --style=javascript
```

---

## 3. Menentukan Gaya Bahasa di Dalam Berkas Kode

### 3.1 Penanda di Baris Pertama (Header Pragma)
Anda dapat menyematkan catatan di baris pertama agar Celang otomatis mengetahui gaya bahasa yang diinginkan untuk berkas tersebut:

```celang
// @gaya: indonesian
buat nama = "Budi"
cetak(nama)
```

Atau dalam berkas bergaya Python:
```celang
# @style: pythonic
nama = "Budi"
print(nama)
```

### 3.2 Mengubah Gaya di Tengah-Tengah Kode (`gaya "..." { }`)
Anda dapat meminjam gaya bahasa lain untuk beberapa baris kode tertentu di tengah program utama:

```celang
// Bagian utama: Bahasa Indonesia
buat saldo = 1000000

// Meminjam gaya Python untuk rumus tertentu:
gaya "pythonic" {
    def hitungBonus(pokok, persen): {
        return pokok * persen / 100
    }
    bonus = hitungBonus(saldo, 10)
}

cetak("Bonus yang diterima:", bonus) // 100000
```
Seluruh variabel yang dibuat di dalam blok tersebut tetap hidup dan dapat digunakan oleh bagian program lainnya.

---

## 4. Tabel Pemetaan Kata Kunci Antar-Gaya

| Maksud Perintah | Bahasa Indonesia | Pythonic | JavaScript | Basa Sunda |
| :--- | :--- | :--- | :--- | :--- |
| **Membuat Variabel** | `buat` | *(langsung nama)* | `let` / `var` | `jieun` |
| **Nilai Tetap** | `konstanta` | `const` | `const` | `tetep` |
| **Pemeriksaan Syarat** | `jika` | `if` | `if` | `lamun` |
| **Pilihan Syarat Kedua**| `selain jika` | `elif` | `else if` | `atawa lamun` |
| **Pilihan Terakhir** | `selain` | `else` | `else` | `lamun henteu` |
| **Perulangan Selama** | `selama` | `while` | `while` | `salami` |
| **Perulangan Daftar** | `untuk ... dalam` | `for ... in` | `for ... of` | `pikeun ... dina` |
| **Membuat Rumus/Fungsi**| `fungsi` | `def` | `function` | `fungsi` |
| **Mengembalikan Nilai** | `kembalikan` | `return` | `return` | `balikkeun` |
| **Tugas Latar Belakang**| `sambil` | `async` | `async` | `bari` |
| **Tugas Serentak** | `bersamaan` | `parallel` | `parallel` | `babarengan` |
| **Nilai Benar** | `benar` | `True` | `true` | `bener` |
| **Nilai Salah** | `salah` | `False` | `false` | `salah` |
| **Nilai Kosong** | `kosong` | `None` | `null` | `kosong` |

---

## 5. Rencana Ekosistem Masa Depan

Konsep Pohon Logika Murni (U-AST) membuka peluang kolaborasi yang sangat menarik untuk masa depan:
- **Kolaborasi Tanpa Konflik Gaya**: Anggota tim dapat menyimpan program dalam format data U-AST murni (`.ast`), sementara masing-masing pengembang melihat kodenya di editor sesuai gaya favorit pribadi.
- **Penerjemahan Bahasa Daerah Lainnya**: Komunitas dapat dengan mudah membuat berkas `.style` baru untuk Bahasa Jawa, Minang, Bali, dan bahasa lainnya cukup dengan menerjemahkan kata kunci di dokumen JSON sederhana.
