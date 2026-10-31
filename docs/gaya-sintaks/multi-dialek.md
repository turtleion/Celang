# Menulis dengan Berbagai Gaya Bahasa (Multi-Dialek)

Pernahkah Anda berdebat dengan teman tentang mana yang lebih enak dibaca: gaya penulisan kurung kurawal atau gaya indentasi, kata kunci Bahasa Indonesia atau Bahasa Inggris?

Di Celang, Anda tidak perlu lagi berdebat. Celang memisahkan antara **makna logika program** dan **gaya tulisan yang tampak di mata**. Anda bebas menulis dan membaca program yang sama dalam gaya Bahasa Indonesia, Basa Jawa, JavaScript, atau Python.

---

## 1. Pilihan Gaya Bahasa Bawaan

Celang menyediakan beberapa pilihan gaya bahasa siap pakai:

### A. Gaya Bahasa Indonesia (`indonesian` - Gaya Utama)
Dirancang agar terasa wajar dan langsung dipahami oleh siapa saja yang berbahasa Indonesia:
```celang
buat nama = "Budi"
jika nama == "Budi" {
    cetak("Halo Budi!")
} selain {
    cetak("Halo tamu!")
}
```

### B. Gaya JavaScript (`javascript` atau `js`)
Bagi pembuat program web yang sudah terbiasa dengan gaya JavaScript:
```celang
let nama = "Budi"
if (nama == "Budi") {
    print("Halo Budi!")
} else {
    print("Halo tamu!")
}
```

### C. Gaya Python (`pythonic` atau `py`)
Bagi yang menyukai gaya penulisan ringkas ala Python:
```celang
nama = "Budi"
if nama == "Budi" {
    print("Halo Budi!")
} else {
    print("Halo tamu!")
}
```

### D. Gaya Bahasa Daerah (Basa Jawa via `jawa` atau `jawa.style`)
Bisa langsung dipanggil bawaan (`--style=jawa`) atau dimuat dari berkas gaya:
```celang
gawe nama = "Budi"
yen nama == "Budi" {
    cetak("Sugeng rawuh Budi!")
} yen ora {
    cetak("Sugeng rawuh tamu!")
}
```

---

## 2. Tabel Persamaan Kata Kunci

Berikut perbandingan kata kunci di berbagai gaya:

| Maksud Perintah | Bahasa Indonesia | JavaScript | Pythonic | Basa Jawa |
| :--- | :--- | :--- | :--- | :--- |
| **Menyiapkan Variabel** | `buat` | `let` / `var` | *(langsung nama)* | `gawe` |
| **Nilai Tetap (Konstanta)** | `konstanta` | `const` | `const` | `tetep` |
| **Pemeriksaan Kondisi** | `jika` | `if` | `if` | `yen` |
| **Pilihan Kedua** | `selain jika` | `else if` | `elif` | `utawa yen` |
| **Pilihan Terakhir** | `selain` | `else` | `else` | `yen ora` |
| **Perulangan Selama** | `selama` | `while` | `while` | `sasuwene` |
| **Perulangan Daftar** | `untuk ... dalam` | `for ... of` | `for ... in` | `kanggo ... ing` |
| **Membuat Rumus/Fungsi** | `fungsi` | `function` | `def` | `fungsi` |
| **Mengembalikan Hasil** | `kembalikan` | `return` | `return` | `balekna` |
| **Tugas Latar Belakang** | `sambil` | `async` | `async` | `sinambi` |
| **Tugas Serentak** | `bersamaan` | `parallel` | `parallel` | `bebarengan` |
| **Nilai Ya / Benar** | `benar` | `true` | `True` | `bener` |
| **Nilai Tidak / Salah** | `salah` | `false` | `False` | `salah` |
| **Nilai Kosong** | `kosong` | `null` | `None` | `suwung` |

---

## 3. Mengubah Tampilan Kode Antar-Gaya (`celang format`)

Anda dapat menyuruh Celang menerjemahkan tampilan berkas kode ke gaya lain secara instan tanpa khawatir merusak alur perhitungan:

```bash
# Menampilkan kode dalam gaya Python di layar:
./celang format program.ce --style=pythonic

# Menampilkan kode dalam gaya JavaScript di layar:
./celang format program.ce --style=javascript

# Menampilkan kode dalam dialek Basa Jawa bawaan:
./celang format program.ce --style=jawa

# Menampilkan kode menggunakan berkas gaya kustom:
./celang format program.ce --style=examples/jawa.style

# Langsung mengubah isi berkas secara permanen (tambahkan opsi -w):
./celang format -w program.ce --style=javascript
```

Dengan fitur ini, seluruh anggota tim dapat bekerja sama dengan tenang, membaca kode dalam gaya yang paling mereka sukai, dan tetap menghasilkan program yang sama persis!
