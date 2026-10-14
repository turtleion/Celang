# Alat Bantu Pertukaran Data (`json`)

Format **JSON** adalah standar tulisan yang paling populer di dunia untuk bertukar data antar-aplikasi dan situs web.

Modul `json` di Celang menyediakan alat untuk membaca data teks JSON menjadi data yang bisa langsung diproses oleh Celang, serta mengubah data Celang kembali menjadi teks JSON yang rapi.

---

## Cara Memanggil

Modul ini sudah siap digunakan secara langsung:

```celang
// 1. Membaca teks JSON menjadi data Celang:
buat data = json.uraikan("{\"nama\": \"Celang\", \"versi\": 1}")

// 2. Mengubah data Celang menjadi teks JSON yang rapi:
buat teksJSON = json.format(data, 2)
cetak(teksJSON)
```

---

## Daftar Alat Bantu

### 1. Mengurai Teks JSON (`uraikan` / `parse`)
Mengubah kalimat teks JSON menjadi data Celang (seperti kamus, daftar, angka, atau teks) agar mudah diakses:

```celang
buat teksMentah = "{\"id\": 101, \"nama\": \"Andi\", \"hobi\": [\"membaca\", \"memasak\"]}"

buat pengguna = json.uraikan(teksMentah)

cetak("Nomor ID    :", pengguna.id)
cetak("Nama        :", pengguna.nama)
cetak("Hobi pertama:", pengguna.hobi[0])
```

> **Petunjuk Ramah:** Jika teks JSON yang Anda masukkan kurang tanda kutip atau kurungnya tidak lengkap, Celang akan memberi tahu secara lembut letak kesalahan ketik tersebut agar mudah Anda betulkan.

---

### 2. Menyusun Teks JSON (`format` / `stringify`)
Mengubah data yang ada di Celang menjadi teks berformat JSON agar bisa dikirim ke aplikasi lain atau disimpan ke berkas:

```celang
buat dataPengaturan = {
    "server": "localhost",
    "port": 8080,
    "aktif": benar,
    "fitur": ["mudah", "cepat", "ramah"]
}

// 1. Format ringkas satu baris:
buat satuBaris = json.format(dataPengaturan)
cetak(satuBaris)

// 2. Format rapi berjarak 2 spasi agar nyaman dibaca mata:
buat teksRapi = json.format(dataPengaturan, 2)
cetak(teksRapi)
```
Tampilan yang dihasilkan akan berjejer rapi ke bawah dengan kurung kurawal yang tertata indah.
