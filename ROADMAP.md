# Peta Jalan Fitur Mendatang Celang (`ROADMAP.md`)

Dokumen ini mencatat rencana pengembangan fitur lanjutan untuk bahasa pemrograman Celang setelah implementasi sistem Visibilitas Modul, Runtime Error Handling, Operator Pipeline (`|>`), dan In-Code Styling (`gaya "..." { }`).

---

## 📌 Prioritas Saat Ini: Standard Library Bawaan (*Batteries Included*)
Fokus utama saat ini adalah menyediakan pustaka standar bawaan agar Celang dapat digunakan untuk membuat program dunia nyata (I/O berkas, matematika, manipulasi teks, waktu, dan format JSON) tanpa bergantung pada pustaka eksternal.

> Rincian pendekatan teknis dan arsitektur pengikatan host runtime (Go) dibahas secara lengkap di berkas terpisah: [`STDLIB_APPROACH.md`](./STDLIB_APPROACH.md).

---

## 🗺️ Rangkuman Fitur-Fitur Masa Depan

### 1. 📦 Standard Library Bawaan (✅ Selesai di `feature/stdlib`)
Modul-modul inti yang terikat langsung ke sistem operasi melalui backend runtime Go:
- **`berkas` / `fs`**: I/O berkas dan folder (`baca`, `tulis`, `tambah`, `ada`, `hapus`).
- **`matematika` / `math`**: Operasi numerik presisi (`akar`, `pangkat`, `mutlak`, `bulat`, `lantai`, `atap`, `acak`, `maks`, `min`, konstanta `PI`, `E`).
- **`teks` / `strings`**: Pemrosesan string (`pisah`, `gabung`, `ganti`, `berisi`, `kapital`, `kecil`, `potong`, `dimulai_dengan`, `diakhiri_dengan`).
- **`waktu` / `time`**: Penanganan waktu & penundaan (`sekarang`, `sekarang_detik`, `format`, `tidur`).
- **`json`**: Serialisasi & parsing objek JSON rekursif (`uraikan` / `parse`, `format` / `stringify`).

---

### 2. ✨ String Interpolation (✅ Selesai)
Menggantikan penggabungan string bertingkat (`"Halo " + nama + ", umur " + umur`) dengan sintaks templat yang bersih dan ekspresif.

#### Contoh Sintaks:
```celang
buat nama = "Budi"
buat saldo = 50000

// Gaya Indonesia:
buat pesan = "Halo {nama}, sisa saldo kamu: Rp {saldo}."

// Gaya JavaScript:
let pesan = `Halo ${nama}, sisa saldo kamu: Rp ${saldo}.`
```

---

### 3. 🧩 Destructuring Assignment (✅ Selesai)
Memungkinkan pemecahan elemen array atau properti map langsung ke variabel mandiri dalam satu baris deklarasi.

#### Contoh Pembongkaran Array:
```celang
buat koordinat = [100, 250]
buat [x, y] = koordinat
// x = 100, y = 250
```

#### Contoh Pembongkaran Map (Object):
```celang
buat profil = { "nama": "Andi", "peran": "Admin" }
buat { nama, peran } = profil
// nama = "Andi", peran = "Admin"
```

---

### 4. 🏗️ Tipe Data Kustom / Struct (`bentuk` / `struct` / `tipe`) (✅ Selesai)
Mendukung pendefinisian cetak biru objek dengan tipe properti yang jelas serta kemampuan menyematkan fungsi (method).

#### Contoh Sintaks:
```celang
bentuk RekeningBank {
    pemilik: Teks,
    saldo itu Angka,

    fungsi setor(jumlah: Angka) {
        diri.saldo += jumlah
    }
}

buat akun = RekeningBank("Budi", 100000)
akun.setor(50000)
cetak(akun.saldo) // 150000
```

---

### 5. 💻 REPL Interaktif (`celang repl` / `celang`)
Menyediakan terminal interaktif baris perintah di mana programmer dapat:
- Mengetik baris kode Celang dan langsung melihat output hasil evaluasinya.
- Menguji pipa ekspresi (`|>`) secara instan.
- Melihat konversi dialek sintaks (*Live Styling Preview*) langsung di terminal.

```bash
$ celang
🌟 Celang REPL v1.0 (Ketik 'keluar' atau 'exit' untuk selesai)
celang> buat x = 10
celang> x |> saring(...) |> cetak
```

---

### 6. 🛠️ Tooling Formatter In-Place (`celang format -w`)
Menambahkan flag `-w` / `--write` pada CLI pemformatan (mirip `gofmt -w` atau `prettier --write`):
```bash
# Menimpa file di tempat dengan kode yang sudah terformat rapi
celang format program.ce --style=indonesian -w
```
Serta integrasi filter otomatis saat check-in ke sistem Git (*clean/smudge filter*).
