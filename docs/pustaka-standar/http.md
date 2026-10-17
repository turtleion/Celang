# Pustaka Jaringan & Protokol HTTP (`http` / `jaringan`)

Modul **`http`** (dapat juga dipanggil dengan alias Bahasa Indonesia **`jaringan`** atau alias singkat **`net`**) menyediakan fungsionalitas lengkap untuk komunikasi jaringan web di Celang:
- **Klien HTTP (*HTTP Client*)**: Pengambilan data REST API, download berkas biner, dan pengiriman payload JSON otomatis.
- **Server HTTP & Router (*Micro-Framework*)**: Pembuatan REST API dan web server modern dengan perutean parameter (`:id`), *middleware*, dan *response chaining*.
- **Utilitas URL**: Enkoding, dekoding, dan penguraian struktur alamat web.

---

## Cara Memanggil

Modul `http` langsung aktif secara global maupun dapat diimpor secara selektif:

```celang
// 1. Langsung panggil via objek global:
buat res = http.get("https://api.github.com/users/octocat")
buat server = http.server()

// 2. Atau menggunakan alias lokal Bahasa Indonesia:
buat res = jaringan.ambil("https://api.github.com/users/octocat")

// 3. Atau impor fungsi tertentu secara spesifik:
ambil { get, post, server } dari http
ambil { ambil, kirim, buat_server } dari jaringan
```

---

## 🚀 1. Klien HTTP (HTTP Client)

Klien HTTP Celang dirancang sangat ringkas, otomatis memproses JSON, dan mendukung unduhan biner langsung ke koleksi [**`Bita`**](/dasar/variabel-tipe#c-koleksi-bita-biner-bita--bytes).

### A. Permintaan GET Sederhana (`get` / `ambil`)
```celang
buat res = http.get("https://jsonplaceholder.typicode.com/posts/1")

cetak("Status Code :", res.status)  // 200
cetak("Apakah Sukses:", res.sukses)  // benar
cetak("Content-Type :", res.headers["content-type"])

// Parsing JSON otomatis ke Map Celang:
buat postingan = res.json()
cetak("Judul Artikel:", postingan.title)
```

### B. Permintaan POST dengan JSON Otomatis (`post` / `kirim`)
Jika data yang dikirim berupa `Map` atau `Array`, Celang otomatis mengubahnya menjadi teks JSON dan menyetel *header* `Content-Type: application/json`:

```celang
buat penggunaBaru = {
    "nama": "Siti Rahma",
    "divisi": "Teknologi",
    "email": "siti@kantor.id"
}

buat res = http.post("https://api.example.com/users", penggunaBaru)
cetak("Status:", res.status)
cetak("Balasan:", res.json())
```

### C. Mengunduh Berkas Biner Langsung ke `Bita`
Metode `res.bita()` / `res.bytes()` membaca respons mentah ke objek `Bita` tanpa risiko korupsi encoding:

```celang
buat res = http.get("https://celang.org/logo.png")
jika res.sukses {
    // Simpan bita gambar secara atomik ke disk:
    berkas.tulis_bita("./logo.png", res.bita())
    cetak("Gambar berhasil diunduh! Ukuran bita:", res.bita().panjang)
}
```

### D. Permintaan Kustom Penuh (`request` / `minta`)
Untuk kendali penuh atas *headers*, *timeout* (tenggat waktu), atau *query parameters*:

```celang
buat res = http.request({
    "method": "POST",
    "url": "https://api.pembayaran.id/transaksi",
    "headers": {
        "Authorization": "Bearer rahasia_token_123",
        "X-Partner-ID": "CELANG-01"
    },
    "params": {
        "mode": "uji_coba"
    },
    "json": {
        "nominal": 250000,
        "keterangan": "Langganan Layanan"
    },
    "timeout": 5000 // dalam milidetik (5 detik)
})
```

---

## 🌐 2. Server HTTP & Router (Micro-Framework)

Celang menyediakan mesin router bawaan yang sangat ringan dan intuitif untuk membangun REST API dan layanan web.

### A. Membuat Server dan Rute REST API
```celang
buat app = http.server()

// 1. Rute Beranda Sederhana:
app.get("/", fungsi(req, res) {
    res.kirim("Selamat datang di Server Celang!")
})

// 2. Rute dengan Parameter Jalur (:id):
app.get("/api/pengguna/:id", fungsi(req, res) {
    buat id = req.params.id
    res.json({
        "id": id,
        "nama": "Budi Santoso",
        "status": "aktif"
    })
})

// 3. Rute POST Pembacaan Payload JSON:
app.post("/api/produk", fungsi(req, res) {
    buat data = req.json()
    // Kembalikan respons status 201 Created:
    res.status(201).json({
        "status": "sukses",
        "pesan": "Produk berhasil ditambahkan",
        "produk": data
    })
})

// 4. Rute HTML:
app.get("/tentang", fungsi(req, res) {
    res.html("<h1>Tentang Celang</h1><p>Bahasa pemrograman human-centric.</p>")
})

// 5. Menjalankan Server:
cetak("Server aktif di http://localhost:8080")
app.listen(8080) // alias: app.layani(8080)
```

### B. Middleware
Gunakan `app.use()` / `app.gunakan()` untuk logging, autentikasi, atau manipulasi *headers* CORS:

```celang
buat app = http.server()

// Middleware global:
app.use(fungsi(req, res, next) {
    cetak("Permintaan masuk:", req.method, req.path)
    res.header("X-Powered-By", "Celang-HTTP")
    res.header("Access-Control-Allow-Origin", "*")
    next() // Lanjutkan ke rute berikutnya
})

// Middleware khusus awalan rute tertentu:
app.use("/admin", fungsi(req, res, next) {
    jika req.headers["authorization"] == kosong {
        res.status(401).json({"error": "Akses ditolak: token tidak ditemukan"})
        kembalikan // Berhenti, jangan panggil next()
    }
    next()
})
```

### C. Ringkasan Objek `req` (Permintaan Masuk)
| Properti / Metode | Tipe | Keterangan |
| :--- | :--- | :--- |
| `req.method` / `req.metode` | `Teks` | Metode HTTP (`"GET"`, `"POST"`, `"PUT"`, dll) |
| `req.path` / `req.jalur` | `Teks` | Jalur URL yang diminta (misal: `"/api/user"`) |
| `req.url` | `Teks` | URL lengkap beserta *query string* |
| `req.params` / `req.parameter`| `Map` | Parameter rute tertangkap (misal: `req.params.id`) |
| `req.query` / `req.kueri` | `Map` | Parameter *query string* (misal: `?kategori=buku`) |
| `req.headers` / `req.kepala` | `Map` | *Header* permintaan masuk (*lowercase keys*) |
| `req.ip` / `req.alamat` | `Teks` | Alamat IP pemanggil |
| `req.body` / `req.badan` | `Teks` | Isi teks mentah *request body* |
| `req.json()` | `Map` / `Array` | Mengurai *request body* sebagai JSON |
| `req.bytes()` / `req.bita()` | `Bita` | Mengambil *request body* dalam bentuk bita biner |

### D. Ringkasan Objek `res` (Respons Keluar)
| Metode | Keterangan | Nilai Balik |
| :--- | :--- | :--- |
| `res.status(code)` | Menyetel kode status HTTP (default `200`) | `res` (*chainable*) |
| `res.header(k, v)` / `res.kepala(k, v)` | Menyetel *header* respons HTTP | `res` (*chainable*) |
| `res.send(data)` / `res.kirim(data)` | Mengirim teks atau data biner `Bita` | `res` |
| `res.json(data, [code])` | Mengirim data terformat JSON secara otomatis | `res` |
| `res.html(str, [code])` | Mengirim dokumen HTML (`text/html; charset=utf-8`)| `res` |
| `res.redirect(url, [code])` | Mengalihkan halaman (*redirect*, default status `302`) | `res` |

---

## 🛠️ 3. Utilitas URL

Untuk memanipulasi dan membedah alamat tautan:

```celang
buat tautan = "https://celang.org:8080/cari?kata=belajar&urut=baru#bab2"
buat info = http.parse_url(tautan)

cetak(info.scheme)    // "https"
cetak(info.host)      // "celang.org"
cetak(info.port)      // 8080
cetak(info.path)      // "/cari"
cetak(info.params)    // {"kata": "belajar", "urut": "baru"}
cetak(info.fragment)  // "bab2"

// Enkoding parameter URL:
buat aman = http.encode_uri_component("kopi & teh + susu")
cetak(aman)           // "kopi%20%26%20teh%20%2B%20susu"
cetak(http.decode_uri_component(aman)) // "kopi & teh + susu"
```

---

## 📋 Ringkasan Daftar Fungsi Global Modul `http`

### Klien HTTP
| Fungsi Indonesia | Alias Internasional | Keterangan |
| :--- | :--- | :--- |
| `ambil(url, [opsi])` | `get` | Melakukan permintaan HTTP GET |
| `kirim(url, [data], [opsi])` | `post` | Melakukan permintaan HTTP POST |
| `taruh(url, [data], [opsi])` | `put` | Melakukan permintaan HTTP PUT |
| `tambal(url, [data], [opsi])` | `patch` | Melakukan permintaan HTTP PATCH |
| `hapus(url, [opsi])` | `delete` | Melakukan permintaan HTTP DELETE |
| `kepala(url, [opsi])` | `head` | Melakukan permintaan HTTP HEAD |
| `minta(opsi)` | `request` | Permintaan HTTP kustom dengan objek opsi |

### Server & Router
| Fungsi Indonesia | Alias Internasional | Keterangan |
| :--- | :--- | :--- |
| `buat_server()` | `server` | Membuat instansi router server HTTP baru |
| `layani(port, handler)` | `serve`, `listen` | Menjalankan server tunggal sederhana (memblokir) |
| `mulai(port, handler)` | `start` | Menjalankan server tunggal di latar belakang (*non-blocking*) |

### Utilitas URL
| Fungsi Indonesia | Alias Internasional | Keterangan |
| :--- | :--- | :--- |
| `urai_url(str)` | `parse_url` | Membedah komponen URL menjadi Map terstruktur |
| `enkode_uri(str)` | `encode_uri` | Enkoding jalur URI |
| `dekode_uri(str)` | `decode_uri` | Dekoding jalur URI |
| `enkode_komponen_uri(str)`| `encode_uri_component` | Enkoding komponen parameter query string |
| `dekode_komponen_uri(str)`| `decode_uri_component` | Dekoding komponen parameter query string |
