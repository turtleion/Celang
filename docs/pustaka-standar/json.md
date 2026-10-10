# Modul JSON (`json`)

Modul `json` menyediakan fungsi serialisasi dan deserialisasi data format JSON secara cepat dan aman, memungkinkan pertukaran data API dan penyimpanan dokumen konfigurasi terstruktur.

---

## Memuat Modul

Tersedia secara global tanpa konfigurasi tambahan:

```celang
// Mengurai teks JSON
buat data = json.uraikan("{\"nama\": \"Celang\", \"versi\": 1}")

// Memformat objek Celang ke string JSON
buat teksJSON = json.format(data, 2)
```

---

## Daftar Fungsi

### 1. `uraikan(teksJSON)` / `parse(teksJSON)`
Mengonversi teks berformat JSON menjadi objek, peta (*map*), array, angka, atau nilai Celang yang sesuai.

```celang
buat payload = "{\"id\": 101, \"nama\": \"Andi\", \"hobi\": [\"membaca\", \"koding\"], \"aktif\": true}"

buat pengguna = json.uraikan(payload)

cetak("ID Pengguna:", pengguna.id)
cetak("Nama Pengguna:", pengguna.nama)
cetak("Hobi pertama:", pengguna.hobi[0])
```

::: tip Penanganan Kesalahan
Jika teks JSON tidak valid (kurang tanda kutip, kurung tidak berpasangan, atau koma berlebih), Celang akan menampilkan pesan diagnostik yang ramah:
```text
Kesalahan: Gagal mengurai teks JSON: invalid character 'a' after object key
Saran: Periksa kembali sintaks JSON Anda (tanda kutip ganda, koma, kurung kurawal).
```
:::

### 2. `format(objek, spasi?)` / `stringify(objek, spasi?)`
Mengonversi objek atau struktur data Celang kembali ke bentuk representasi string JSON.
- Argumen kedua bersifat opsional untuk memberikan indentasi spasi (*pretty-print*).

```celang
buat konfigurasi = {
    "server": "localhost",
    "port": 8080,
    "debug": benar,
    "fitur": ["uast", "bytecode-vm", "syntax-styling"]
}

// Format ringkas satu baris (compact)
buat compact = json.format(konfigurasi)
cetak(compact)

// Format rapi dengan indentasi 2 spasi (pretty-print)
buat rapi = json.format(konfigurasi, 2)
cetak("JSON Berindentasi:\n" + rapi)
```
