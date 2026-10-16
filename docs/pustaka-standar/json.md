# Pustaka Pertukaran Data (`json`)

Modul **`json`** menyediakan alat lengkap untuk membaca, menyusun, memvalidasi, memanipulasi, dan menyamarkan data berformat **JSON** (*JavaScript Object Notation*) di Celang.

Selain operasi penguraian (*parse*) dan penataan teks (*stringify*) standar, modul `json` di Celang dirancang dengan **fitur-fitur cerdas** yang sangat memudahkan pengolahan data bersarang (*nested data*) dan integrasi API modern:
- **Akses Data Bersarang Aman (`json.dapatkan` / `get`)**: Mengambil nilai bertingkat dalam (*deep path*) seperti `"pengguna.alamat.kota"` atau `"items[0].nama"` tanpa takut galat *null/undefined*.
- **Modifikasi Bersarang Instan (`json.setel` / `set` & `json.hapus` / `delete`)**: Menulis atau menghapus nilai pada jalur bersarang secara otomatis.
- **Penggabungan Mendalam (*Deep Merge*) (`json.gabung` / `merge`)**: Menggabungkan dua atau lebih Kamus JSON secara rekursif (sangat cocok untuk menggabungkan konfigurasi *default* dengan pengaturan pengguna).
- **Sensor Otomatis Data Sensitif (`json.samarkan` / `mask`)**: Otomatis mendeteksi dan menyensor field sensitif (*password*, *token*, *kartu_kredit*, *api_key*, *pin*, dll.) menjadi `"***"` agar aman saat dicetak ke file log server.
- **Dukungan Streaming JSON Lines / NDJSON (`json.uraikan_baris` & `json.format_baris`)**: Standar industri untuk membaca dan menulis *streaming logs* dan *big data* (satu JSON per baris).
- **Ratakan & Susun (*Flatten & Unflatten*) (`json.ratakan` & `json.susun`)**: Mengubah struktur bersarang menjadi pasangan kunci datar bertitik, dan sebaliknya.
- **Pemeriksaan & Pemadatan**: Validasi JSON tanpa *crash* (`json.apakah_valid`), pembuatan JSON rapi berindentasi (`json.rapi`), dan pemadatan tanpa spasi (`json.padatkan`).

---

## Cara Memanggil

Modul ini aktif secara global dan dapat langsung dipanggil:

```celang
// 1. Membaca teks JSON menjadi objek Celang:
buat data = json.uraikan("{\"nama\": \"Celang\", \"versi\": 2}")
cetak(data.nama) // "Celang"

// 2. Mengubah objek Celang menjadi teks JSON rapi:
buat teksJSON = json.rapi(data)
cetak(teksJSON)

// 3. Atau impor fungsi tertentu (tanpa tanda kutip):
ambil { dapatkan, setel, gabung, samarkan } dari json
```

---

## 🌟 Fitur Cerdas & Manipulasi Jalur (*JSONPath*)

### 1. Mengambil Nilai Bersarang Aman (`dapatkan` / `petik` / `get`)
Mengakses data di dalam struktur objek atau daftar yang bertingkat dalam tanpa khawatir terjadi galat *null reference*:

```celang
buat responsAPI = {
    "status": 200,
    "data": {
        "pengguna": [
            {
                "id": 101,
                "nama": "Andi",
                "alamat": {
                    "kota": "Jakarta",
                    "kode_pos": 10110
                }
            }
        ]
    }
}

// Akses properti objek bersarang:
buat kota = json.dapatkan(responsAPI, "data.pengguna[0].alamat.kota")
cetak("Kota domisili:", kota) // "Jakarta"

// Akses dengan nilai bawaan (fallback) jika jalur tidak ditemukan:
buat nomorTelepon = json.dapatkan(responsAPI, "data.pengguna[0].kontak.telepon", "Belum Diisi")
cetak("Telepon:", nomorTelepon) // "Belum Diisi"
```

---

### 2. Menulis & Menghapus Jalur Bersarang (`setel` & `hapus`)
Menuliskan data ke jalur bersarang yang belum ada (Celang otomatis membuatkan Kamus perantara di sepanjang jalur):

```celang
buat setelan = {}

// 1. Menulis nilai ke path bertingkat dalam:
json.setel(setelan, "keamanan.autentikasi.dua_faktor", benar)
json.setel(setelan, "tampilan.tema", "gelap")

cetak(json.rapi(setelan))
// Hasil:
// {
//   "keamanan": {
//     "autentikasi": {
//       "dua_faktor": true
//     }
//   },
//   "tampilan": {
//     "tema": "gelap"
//   }
// }

// 2. Menghapus field tertentu:
json.hapus(setelan, "tampilan.tema")
```

---

### 3. Penggabungan Mendalam (*Deep Merge*) (`gabung` / `merge`)
Menggabungkan dua Kamus secara rekursif hingga ke lapisan terdalam:

```celang
buat konfigurasiBawaan = {
    "tema": "terang",
    "server": {
        "host": "localhost",
        "port": 8080,
        "ssl": salah
    }
}

buat konfigurasiPengguna = {
    "tema": "gelap",
    "server": {
        "ssl": benar
    }
}

buat konfigurasiFinal = json.gabung(konfigurasiBawaan, konfigurasiPengguna)
cetak(json.rapi(konfigurasiFinal))
// Hasil:
// {
//   "tema": "gelap",
//   "server": {
//     "host": "localhost",
//     "port": 8080,
//     "ssl": true
//   }
// }
```

---

### 4. Sensor Otomatis Data Sensitif (`samarkan` / `mask`)
Mencegah bocornya data rahasia (*password*, *token*, *pin*, *kartu_kredit*, dll.) saat mencetak payload ke file log:

```celang
buat permintaanMasuk = {
    "pengguna": "budi_santoso",
    "kata_sandi": "kucing_oren_123",
    "kartu_kredit": "4111-2222-3333-4444",
    "metadata": {
        "api_key": "sec_live_987654321"
    }
}

// Otomatis menyensor field yang mengandung kata kunci sensitif:
buat payloadAman = json.samarkan(permintaanMasuk)
cetak(json.rapi(payloadAman))
// Hasil:
// {
//   "pengguna": "budi_santoso",
//   "kata_sandi": "***",
//   "kartu_kredit": "***",
//   "metadata": {
//     "api_key": "***"
//   }
// }

// Anda juga bisa menambahkan kata kunci kustom:
buat kustom = json.samarkan(permintaanMasuk, ["pengguna"])
```

---

### 5. Dukungan JSON Lines / NDJSON (`uraikan_baris` & `format_baris`)
Format streaming populer di mana setiap baris adalah satu objek JSON mandiri:

```celang
// 1. Mengurai teks JSONL menjadi Daftar (Array) objek:
buat logMentah = "{\"waktu\": \"08:00\", \"level\": \"info\"}\n{\"waktu\": \"08:01\", \"level\": \"error\"}\n"
buat daftarLog = json.uraikan_baris(logMentah)

cetak("Total baris log:", panjang(daftarLog)) // 2
cetak("Status baris 2 :", daftarLog[1].level)   // "error"

// 2. Mengubah Daftar objek menjadi teks JSONL:
buat stringJSONL = json.format_baris(daftarLog)
```

---

### 6. Ratakan & Susun (*Flatten & Unflatten*) (`ratakan` & `susun`)
Mengonversi struktur pohon menjadi pasangan kunci datar bertitik untuk serialisasi tabel/formulir, dan merekonstruksinya kembali:

```celang
buat dataPengguna = {
    "profil": {
        "nama": "Andi",
        "usia": 25
    }
}

// Ratakan (Flatten):
buat datar = json.ratakan(dataPengguna)
cetak(datar) // {"profil.nama": "Andi", "profil.usia": 25}

// Susun kembali (Unflatten):
buat bersarang = json.susun(datar)
cetak(bersarang.profil.nama) // "Andi"
```

---

## 🛠️ Validasi, Format, & Pemadatan

### 1. Pemeriksaan Validitas JSON (`apakah_valid` / `is_valid`)
Memeriksa keabsahan sintaks JSON sebelum diolah tanpa khawatir memicu *runtime error*:

```celang
jika json.apakah_valid("{\"ok\": benar}") {
    cetak("Data JSON valid!")
}

jika !json.apakah_valid("teks sembarang") {
    cetak("Bukan format JSON yang sah.")
}
```

### 2. Format Rapi & Pemadatan (`rapi` & `padatkan`)
```celang
buat data = {"a": 1, "b": [2, 3]}

// Format rapi berjarak 2 spasi:
cetak(json.rapi(data))

// Padatkan ke 1 baris tanpa spasi (hemat ukuran byte):
buat minim = json.padatkan(data)
cetak(minim) // "{"a":1,"b":[2,3]}"
```

---

## 📋 Tabel Rangkuman Fungsi

| Fungsi Bahasa Indonesia | Alias Internasional | Keterangan Singkat |
| :--- | :--- | :--- |
| **`uraikan(str)`** | `parse` | Mengurai teks JSON menjadi objek Celang |
| **`format(obj, [spasi])`** | `stringify` | Menyusun objek Celang menjadi teks JSON |
| **`apakah_valid(str)`** | `is_valid` | Memeriksa validitas teks JSON (boolean) |
| **`dapatkan(obj, path, [def])`**| `get` / `petik` | Mengambil data bertingkat dalam via path |
| **`setel(obj, path, nilai)`** | `set` | Menulis data ke path bertingkat secara otomatis |
| **`hapus(obj, path)`** | `delete` | Menghapus properti pada path bertingkat |
| **`gabung(...objs)`** | `merge` | Menggabungkan Kamus JSON secara rekursif (*deep merge*) |
| **`samarkan(obj, [kunci])`** | `mask` | Sensor otomatis field sensitif (*password*, *token*, dll.) |
| **`uraikan_baris(str)`** | `parse_lines` | Mengurai teks JSONL / NDJSON menjadi Daftar objek |
| **`format_baris(daftar)`** | `stringify_lines`| Menyusun Daftar objek menjadi format JSONL / NDJSON |
| **`ratakan(obj)`** | `flatten` | Meratakan objek bertingkat menjadi kunci bertitik |
| **`susun(objDatar)`** | `unflatten` | Menyusun kembali kunci bertitik menjadi objek bersarang |
| **`rapi(obj/str, [spasi])`**| `pretty` | Memformat JSON dengan indentasi rapi |
| **`padatkan(obj/str)`** | `minify` | Memadatkan JSON menjadi 1 baris tanpa spasi |
