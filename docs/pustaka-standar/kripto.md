---
outline: deep
---

# Pustaka Kriptografi (`kripto` / `crypto`)

Modul **`kripto`** (atau dengan alias internasional **`crypto`**) menyediakan ekosistem kriptografi berstandar industri modern di Celang. Diimplementasikan murni menggunakan pustaka kriptografi Go (*pure Go*, tanpa dependensi C/eksternal), modul ini dirancang dengan prinsip **aman secara bawaan (*safe-by-default*)** dan terinspirasi dari pola resep tingkat tinggi pada pustaka modern seperti Python `cryptography.io`.

Modul ini mencakup solusi komprehensif mulai dari resep enkripsi terotentikasi, hashing sandi tahan serangan brute-force, tanda tangan digital asimetris Ed25519, hingga fungsi hash dan perbandingan waktu konstan:

1. **Enkripsi Simetris Terotentikasi (AEAD AES-256-GCM)**: Enkripsi dan dekripsi data aman dengan perlindungan manipulasi otomatis (`buat_kunci`, `enkripsi`, `dekripsi`).
2. **Password Hashing & Derivasi Kunci (PBKDF2-HMAC-SHA256)**: Pengamanan kata sandi berstandar format PHC dengan *constant-time comparison* serta pembuatan kunci simetris dari sandi (`hash_sandi`, `verifikasi_sandi`, `turunkan_kunci`).
3. **Tanda Tangan Digital Asimetris (Ed25519 - RFC 8032)**: Penandatanganan dan verifikasi pesan digital berkecepatan tinggi dengan ukuran kunci ringkas 32-byte (`buat_pasangan_kunci`, `tanda_tangan`, `verifikasi_tanda_tangan`).
4. **Verifikasi Integritas & Checksum (Hash Digest)**: SHA-256, SHA-512, SHA-1, dan MD5.
5. **Otentikasi Pesan Berbasis Kunci (HMAC)**: HMAC-SHA256 untuk tanda tangan webhook dan integritas payload API.
6. **Entropi & Keacakan Kriptografis**: Mengambil keacakan aman tingkat sistem operasi untuk token sesi dan salt (`acak_heks`, `acak_bita`).
7. **Pengenal Unik Global (UUID v4)**: Pembangkit UUID versi 4 standar RFC 4122 (`uuid`).
8. **Perlindungan Anti-Timing Attack**: Perbandingan waktu konstan untuk mencegah kebocoran informasi melalui perbedaan latensi CPU (`konstanta_sama`).

---

## 🚀 Cara Menggunakan

Modul `kripto` dapat diakses langsung sebagai objek global maupun diimpor secara modular:

```celang
// 1. Akses langsung melalui objek global:
buat kunci = kripto.buat_kunci()
buat rahasia = kripto.enkripsi("Data rahasia", kunci)
buat asli = kripto.dekripsi(rahasia, kunci)

// 2. Menggunakan alias bahasa Inggris / Universal:
buat token = crypto.random_hex(16)
buat id = crypto.uuid()

// 3. Impor spesifik melalui sistem modul:
ambil { buat_kunci, enkripsi, dekripsi, hash_sandi, verifikasi_sandi } dari "kripto"
```

### 🌍 Dukungan Mode Universal (`"use universal"`)

Celang menyediakan direktif `"use universal"` untuk basis kode internasional. Saat mode universal diaktifkan di baris paling awal berkas:

- Wajib menggunakan modul dengan nama bahasa Inggris: `crypto` (misal: `import { encrypt, decrypt } from "crypto"`).
- Pemanggilan nama modul bahasa Indonesia `kripto` akan memicu diagnostik empati Celang yang ramah, memandu pengembang untuk beralih ke nama `crypto`.
- Semua alias fungsi bahasa Inggris (`crypto.generate_key()`, `crypto.encrypt()`, `crypto.decrypt()`, `crypto.hash_password()`, `crypto.verify()`, dll.) siap digunakan berdampingan dengan kata kunci bahasa Inggris (`let`, `const`, `if`, `def`, `try`, dll.).

```celang
"use universal"
import { generate_key, encrypt, decrypt } from "crypto"

let key = generate_key()
let ciphertext = encrypt("Confidential payload", key)
let recovered = decrypt(ciphertext, key)
```

Untuk panduan sintaks lengkap dalam bahasa Inggris, silakan lihat **[Dokumentasi Pustaka `crypto` (English)](/pustaka-standar/crypto)**.

---

## 📑 Daftar Fungsi Lengkap

### 1. Resep Enkripsi Simetris (AEAD AES-256-GCM)

| Fungsi (Bahasa Indonesia) | Alias Universal (Inggris) | Parameter Masukan | Tipe Kembalian | Deskripsi Singkat |
| :--- | :--- | :--- | :--- | :--- |
| `kripto.buat_kunci()` | `crypto.generate_key()` | `-` | `Teks` (Base64) | Menghasilkan kunci acak 32-byte (256-bit) aman format Base64 |
| `kripto.enkripsi(data, kunci)` | `crypto.encrypt(data, key)` | `data`: `Teks` / `Bita`<br>`kunci`: `Teks` / `Bita` | `Teks` (Base64) | Mengenkripsi data dengan AES-256-GCM, nonce 12-byte otomatis, dan tag autentikasi 16-byte |
| `kripto.dekripsi(ciphertext, kunci)` | `crypto.decrypt(ciphertext, key)` | `ciphertext`: `Teks` / `Bita`<br>`kunci`: `Teks` / `Bita` | `Teks` | Mendekripsi ciphertext, memverifikasi integritas, atau melempar galat empati jika rusak/kunci salah |

### 2. Password Hashing & Derivasi Kunci (PBKDF2-HMAC-SHA256)

| Fungsi (Bahasa Indonesia) | Alias Universal (Inggris) | Parameter Masukan | Tipe Kembalian | Deskripsi Singkat |
| :--- | :--- | :--- | :--- | :--- |
| `kripto.hash_sandi(sandi)` | `crypto.hash_password(password)` | `sandi`: `Teks` / `Bita` | `Teks` (PHC) | Menghasilkan hash sandi tahan brute-force GPU (100.000 iterasi, salt acak 16-byte, format PHC) |
| `kripto.verifikasi_sandi(sandi, hash)` | `crypto.verify_password(password, hash)` | `sandi`: `Teks` / `Bita`<br>`hash`: `Teks` | `Benar` / `Salah` | Memverifikasi kecocokan kata sandi dengan hash menggunakan perbandingan waktu konstan |
| `kripto.turunkan_kunci(sandi, salt, panjang, [format])` | `crypto.derive_key(password, salt, len, [format])` | `sandi`: `Teks` / `Bita`<br>`salt`: `Teks` / `Bita`<br>`panjang`: `Angka`<br>`format`: `Teks` (opsional) | `Bita` / `Teks` | Menurunkan kunci kriptografis dari sandi dan salt. Format: `"bita"` (bawaan), `"heks"`, `"base64"` |

### 3. Tanda Tangan Digital Asimetris (Ed25519 - RFC 8032)

| Fungsi (Bahasa Indonesia) | Alias Universal (Inggris) | Parameter Masukan | Tipe Kembalian | Deskripsi Singkat |
| :--- | :--- | :--- | :--- | :--- |
| `kripto.buat_pasangan_kunci()` | `crypto.generate_keypair()` | `-` | `Peta` (`Map`) | Menghasilkan pasangan kunci Ed25519: `kunci_publik` (64 heks) dan `kunci_privat` (64 heks seed) |
| `kripto.tanda_tangan(data, kunci_privat)` | `crypto.sign(data, private_key)` | `data`: `Teks` / `Bita`<br>`kunci_privat`: `Teks` / `Bita` | `Teks` (heksadesimal) | Menghasilkan tanda tangan digital Ed25519 sepanjang 64-byte (128 karakter heksadesimal) |
| `kripto.verifikasi_tanda_tangan(data, signature, kunci_publik)` | `crypto.verify(data, signature, public_key)` | `data`: `Teks` / `Bita`<br>`signature`: `Teks` / `Bita`<br>`kunci_publik`: `Teks` / `Bita` | `Benar` / `Salah` | Memvalidasi keaslian tanda tangan digital terhadap kunci publik pengirim |

### 4. Hash Digest, HMAC, Entropi Acak, & Keamanan Waktu

| Fungsi (Bahasa Indonesia) | Alias Universal (Inggris) | Parameter Masukan | Tipe Kembalian | Deskripsi Singkat |
| :--- | :--- | :--- | :--- | :--- |
| `kripto.sha256(data)` | `crypto.sha256(data)` | `data`: `Teks` / `Bita` | `Teks` (heksadesimal) | Menghitung digest SHA-256 (64 karakter heksadesimal) |
| `kripto.sha256_berkas(jalur)` | `crypto.sha256_file(path)` | `jalur`: `Teks` | `Teks` (heksadesimal) | Menghitung SHA-256 berkas langsung via *buffered streaming* (hemat RAM) |
| `kripto.sha512(data)` | `crypto.sha512(data)` | `data`: `Teks` / `Bita` | `Teks` (heksadesimal) | Menghitung digest SHA-512 (128 karakter heksadesimal) |
| `kripto.sha512_berkas(jalur)` | `crypto.sha512_file(path)` | `jalur`: `Teks` | `Teks` (heksadesimal) | Menghitung SHA-512 berkas langsung via *buffered streaming* |
| `kripto.sha1(data)` | `crypto.sha1(data)` | `data`: `Teks` / `Bita` | `Teks` (heksadesimal) | Menghitung digest SHA-1 (kompatibilitas warisan, 40 karakter) |
| `kripto.md5(data)` | `crypto.md5(data)` | `data`: `Teks` / `Bita` | `Teks` (heksadesimal) | Menghitung digest MD5 (kompatibilitas warisan, 32 karakter) |
| `kripto.md5_berkas(jalur)` | `crypto.md5_file(path)` | `jalur`: `Teks` | `Teks` (heksadesimal) | Menghitung MD5 berkas langsung via *buffered streaming* |
| `kripto.crc32(data)` | `crypto.crc32(data)` | `data`: `Teks` / `Bita` | `Angka` (UInt32) | Menghitung checksum cepat CRC-32 IEEE (standar berkas ZIP & PNG) |
| `kripto.crc32_berkas(jalur)` | `crypto.crc32_file(path)` | `jalur`: `Teks` | `Angka` (UInt32) | Menghitung CRC-32 berkas langsung via *buffered streaming* |
| `kripto.hmac_sha256(kunci, data)` | `crypto.hmac_sha256(key, data)` | `kunci`: `Teks` / `Bita`<br>`data`: `Teks` / `Bita` | `Teks` (heksadesimal) | Menghitung Kode Autentikasi Pesan HMAC berbasis SHA-256 |
| `kripto.acak_heks(panjang)` | `crypto.random_hex(panjang)` | `panjang`: `Angka` (byte) | `Teks` (heksadesimal) | Menghasilkan string heksadesimal acak CSPRNG (`2 * panjang` heks) |
| `kripto.acak_bita(panjang)` | `crypto.random_bytes(panjang)` | `panjang`: `Angka` (byte) | `Bita` | Menghasilkan buffer objek `Bita` biner acak dari CSPRNG |
| `kripto.uuid()` | `crypto.uuid()` | `-` | `Teks` (UUID v4) | Menghasilkan ID acak unik standar RFC 4122 versi 4 |
| `kripto.konstanta_sama(a, b)` | `crypto.timing_safe_equal(a, b)` | `a`, `b`: `Teks` / `Bita` | `Benar` / `Salah` | Membandingkan dua nilai rahasia dalam waktu konstan anti-timing attack |

---

## 🛡️ Resep Kriptografi Tingkat Tinggi (High-Level Recipes)

Celang mengadopsi filosofi *safe-by-default*: pengembang aplikasi tidak perlu merancang konstruksi kriptografi tingkat rendah sendiri yang rentan kesalahan fatal (seperti penggunaan mode CBC tanpa HMAC, penggunaan IV statis, atau penggunaan hash cepat untuk kata sandi). Seluruh algoritma yang dipilih telah memenuhi standar modern.

---

### 1. Enkripsi Simetris Terotentikasi (AES-256-GCM AEAD)

Resep enkripsi simetris Celang menggunakan **AES-256** dalam mode **GCM** (*Galois/Counter Mode*), yaitu algoritma AEAD (*Authenticated Encryption with Associated Data*).

#### Keunggulan & Jaminan Keamanan:
- **Kerahasiaan & Integritas Bersama**: Selain mengacak data asli, GCM menghitung *authentication tag* sepanjang 16-byte (128-bit) yang memastikan data tidak dapat diubah di tengah jalan tanpa terdeteksi.
- **Nonce Acak Otomatis**: Setiap pemanggilan `kripto.enkripsi()` otomatis menghasilkan *nonce* acak unik sepanjang 12-byte dari CSPRNG sistem (`crypto/rand`). Dengan demikian, mengenkripsi teks yang sama berulang kali akan menghasilkan *ciphertext* yang selalu berbeda, mencegah penyerang menebak pola data.
- **Format Payload Standar**: Ciphertext keluaran merupakan string Base64 yang menggabungkan:
  $$\text{Payload} = \text{Base64}(\text{Nonce 12-byte} \parallel \text{Ciphertext} \parallel \text{Tag Autentikasi 16-byte})$$
- **Deteksi Manipulasi Data (*Tamper-Evident*)**: Jika ciphertext diubah meskipun hanya 1 bit, atau jika didekripsi dengan kunci yang salah, proses dekripsi **gagal secara aman** dan mengembalikan galat empati Celang tanpa menyebabkan aplikasi berhenti mendadak (*crash*).

#### Format Kunci yang Didukung:
Argumen `kunci` pada `enkripsi` dan `dekripsi` sangat fleksibel:
1. Kunci Base64 standar 44-karakter (hasil dari `kripto.buat_kunci()`).
2. Kunci heksadesimal 64-karakter (32 byte murni).
3. Objek `Bita` tepat 32 bita (misal dari `kripto.turunkan_kunci()`).
4. Teks mentah sepanjang tepat 32 karakter byte.

#### Contoh Kode Lengkap: Enkripsi, Dekripsi, & Deteksi Manipulasi

```celang
// 1. Buat kunci acak 256-bit (32 byte)
buat kunci = kripto.buat_kunci()
cetak("Kunci AES-256:", kunci)

// 2. Enkripsi pesan rahasia
buat pesanAsli = "Dokumen transaksi bernilai Rp 150.000.000"
buat terenkripsi = kripto.enkripsi(pesanAsli, kunci)
cetak("Ciphertext (Base64):", terenkripsi)

// 3. Dekripsi dan pulihkan data asli
buat hasilDekripsi = kripto.dekripsi(terenkripsi, kunci)
cetak("Data berhasil dipulihkan:", hasilDekripsi)

// 4. Uji Nonce Acak Unik:
// Enkripsi kedua terhadap pesan yang sama menghasilkan ciphertext berbeda!
buat terenkripsiKedua = kripto.enkripsi(pesanAsli, kunci)
jika (terenkripsi != terenkripsiKedua) {
    cetak("✅ Setiap pemanggilan menghasilkan ciphertext unik berkat nonce otomatis.")
}

// 5. Uji Deteksi Modifikasi (Tamper Detection):
// Jika ada pihak ketiga memodifikasi ciphertext atau kunci keliru, tangkap secara elegan:
buat kunciPalsu = kripto.buat_kunci()

coba {
    // Mencoba mendekripsi dengan kunci palsu
    kripto.dekripsi(terenkripsi, kunciPalsu)
} tangkap (galat) {
    cetak("🛡️ Berhasil dicegat:", galat)
    // Galat: kripto.dekripsi: dekripsi gagal (otentikasi gagal, data tersandi telah dimodifikasi atau kunci salah)
}
```

---

### 2. Password Hashing & Derivasi Kunci (PBKDF2-HMAC-SHA256)

Jangan pernah menyimpan kata sandi pengguna dalam bentuk teks mentah ataupun menggunakan fungsi hash cepat biasa seperti MD5 atau SHA-256 (karena mudah dibongkar menggunakan kartu grafis / GPU farm).

Celang menyediakan fungsi *password hashing* berbasis **PBKDF2-HMAC-SHA256** dengan konfigurasi aman:
- **100.000 Iterasi (*Work Factor*)**: Memberikan beban komputasi yang memperlambat serangan *brute-force* dan *dictionary attack*.
- **Salt Acak Unik 16-byte**: Dibuat otomatis untuk setiap pemanggilan fungsi `hash_sandi`, mencegah serangan *rainbow table*.
- **Format Standar PHC (Password Hashing Competition)**:
  Format string yang dihasilkan mengikuti standar industri modern:
  `$pbkdf2-sha256$i=100000$<salt_base64>$<hash_base64>`
- **Verifikasi Waktu Konstan**: `kripto.verifikasi_sandi` membandingkan digest menggunakan `subtle.ConstantTimeCompare`, melindungi sistem dari serangan analisis waktu CPU (*timing attacks*).

#### Menurunkan Kunci Enkripsi (`kripto.turunkan_kunci`)

Jika Anda ingin mengenkripsi data menggunakan kata sandi/frasa sandi yang dimasukkan pengguna (*passphrase*), gunakan `kripto.turunkan_kunci()` untuk mengubah kata sandi bervariasi menjadi kunci biner berkualitas tinggi yang siap digunakan pada `kripto.enkripsi()`.

Pilihan format keluaran pada argumen ke-4 (opsional):
- `"bita"` (bawaan): Mengembalikan objek `Bita`.
- `"heks"` / `"hex"`: Mengembalikan teks heksadesimal.
- `"base64"` / `"b64"`: Mengembalikan teks Base64.

#### Contoh Kode Lengkap: Otentikasi Pengguna & KDF

```celang
// -------------------------------------------------------------
// Skenario A: Pendaftaran & Otentikasi Kata Sandi Pengguna
// -------------------------------------------------------------

buat sandiPengguna = "KataSandiRahasia2026!"

// 1. Simpan hash ini ke basis data saat pendaftaran pengguna:
buat hashDisimpan = kripto.hash_sandi(sandiPengguna)
cetak("Hash PHC tersimpan:", hashDisimpan)
// Contoh: $pbkdf2-sha256$i=100000$AbCdEf123456...$XyZ987...

// 2. Saat pengguna login, verifikasi kata sandi yang dimasukkan:
buat inputLogin = "KataSandiRahasia2026!"
buat loginSukses = kripto.verifikasi_sandi(inputLogin, hashDisimpan)

jika (loginSukses) {
    cetak("✅ Login berhasil! Sandi cocok.")
} selain {
    cetak("❌ Login ditolak: Sandi keliru.")
}

// Uji percobaan dengan kata sandi salah:
buat inputKeliru = "SandiTebakanSembarangan"
buat hasilCek = kripto.verifikasi_sandi(inputKeliru, hashDisimpan)
cetak("Verifikasi sandi salah bernilai:", hasilCek) // salah


// -------------------------------------------------------------
// Skenario B: Menurunkan Kunci Enkripsi dari Kata Sandi (KDF)
// -------------------------------------------------------------

buat passphrase = "kunci-master-dompet-digital"
// Buat salt acak unik (minimal 16 bita)
buat salt = kripto.acak_bita(16)

// Turunkan kunci simetris AES-256 sepanjang 32-byte berformat Base64
buat kunciTurunan = kripto.turunkan_kunci(passphrase, salt, 32, "base64")

// Kunci turunan dapat langsung dipakai mengenkripsi data sensitif:
buat catatan = "Catatan rahasia yang dilindungi passphrase"
buat arsipTersandi = kripto.enkripsi(catatan, kunciTurunan)
buat catatanPulih = kripto.dekripsi(arsipTersandi, kunciTurunan)

cetak("Catatan berhasil dipulihkan dari derived key:", catatanPulih)
```

---

### 3. Tanda Tangan Digital Modern Asimetris (Ed25519 - RFC 8032)

Untuk skenario di mana dua pihak ingin memastikan keaslian data tanpa berbagi kunci rahasia yang sama (seperti verifikasi pembaruan perangkat lunak, *package manager*, atau pengesahan dokumen), Celang menyediakan **Ed25519** (*Edwards-curve Digital Signature Algorithm* berbasis Curve25519) sesuai spesifikasi **RFC 8032**.

#### Keunggulan Ed25519:
- **Kunci Publik Ringkas**: Hanya 32 byte (64 karakter heksadesimal).
- **Keamanan Ekstra Tinggi**: Menyediakan tingkat keamanan setara RSA 3072-bit namun dengan performa jauh lebih cepat.
- **Tahan Serangan Side-Channel**: Desain matematis bebas cabang dan tabel acuan memori.
- **Deterministik**: Penandatanganan pesan yang sama dengan kunci yang sama menghasilkan tanda tangan yang identik dan dapat diprediksi tanpa risiko kehabisan entropi acak saat menandatangani.

#### Struktur Nilai Kembalian `buat_pasangan_kunci()`:
Mengembalikan objek `Peta` (`Map`) yang memuat kunci publik dan kunci privat (seed):
- `pasangan.kunci_publik` (alias `pasangan.public_key`): String heksadesimal 64 karakter.
- `pasangan.kunci_privat` (alias `pasangan.private_key`): String heksadesimal 64 karakter (seed 32-byte).

#### Fleksibilitas Format Kunci & Tanda Tangan:
Fungsi `tanda_tangan` dan `verifikasi_tanda_tangan` menerima kunci privat, kunci publik, serta tanda tangan baik dalam format string **heksadesimal**, string **Base64**, teks bita mentah, maupun objek **`Bita`**.

#### Contoh Kode Lengkap: Penandatanganan & Verifikasi Rilis Paket

```celang
// 1. Pengembang menghasilkan pasangan kunci Ed25519
buat pasangan = kripto.buat_pasangan_kunci()
buat kunciPublik = pasangan.kunci_publik
buat kunciPrivat = pasangan.kunci_privat

cetak("Kunci Publik (bagikan ke pengguna):", kunciPublik)
cetak("Kunci Privat (jaga kerahasiaannya):", kunciPrivat)

// 2. Penulis paket menandatangani manifes rilis
buat manifesPaket = "nama=celang-http&versi=2.1.0&sha256=a1b2c3d4e5f6..."
buat signature = kripto.tanda_tangan(manifesPaket, kunciPrivat)
cetak("Signature Ed25519 (128 heks):", signature)

// 3. Pengguna / CLI Celang memverifikasi manifes paket sebelum instalasi
buat valid = kripto.verifikasi_tanda_tangan(manifesPaket, signature, kunciPublik)

jika (valid) {
    cetak("✅ Paket otentik! Manifes ditandatangani oleh pemegang kunci privat yang sah.")
} selain {
    cetak("❌ Verifikasi gagal: Paket mencurigakan atau kunci tidak cocok.")
}

// 4. Deteksi Pemalsuan Paket (Tamper Verification):
// Jika isi manifes diubah oleh peretas:
buat manifesPalsu = "nama=celang-http&versi=2.1.0-malware&sha256=99999..."
buat cekPalsu = kripto.verifikasi_tanda_tangan(manifesPalsu, signature, kunciPublik)
cetak("Verifikasi berkas yang dimodifikasi:", cekPalsu) // salah
```

---

## 💡 Contoh Praktis Fungsi Dasar & Utilitas

### 1. Menghitung Checksum SHA-256 pada Teks & Berkas Biner

Semua fungsi hashing Celang menerima masukan berupa teks biasa (`Teks`) maupun objek biner mentah (`Bita`):

```celang
impor berkas

// Hash teks biasa
buat hashTeks = kripto.sha256("Celang Programming Language")
cetak("SHA-256 teks:", hashTeks)

// Hash berkas biner mentah
buat dataBita = berkas.baca_bita("./arsip.tar.gz")
buat hashBerkas = kripto.sha256(dataBita)
cetak("SHA-256 berkas:", hashBerkas)
```

### 2. Otentikasi Webhook API dengan HMAC-SHA256

Verifikasi tanda tangan permintaan HTTP yang dikirim oleh penyedia layanan pihak ketiga (misal: payment gateway):

```celang
buat rahasiaWebhook = "kunci_rahasia_webhook_mitra_32char"
buat payload = "{\"event\": \"pembayaran_sukses\", \"nominal\": 500000}"

// Buat signature HMAC-SHA256
buat signatureDikirim = kripto.hmac_sha256(rahasiaWebhook, payload)
cetak("Header X-Signature:", signatureDikirim)

// Verifikasi di sisi penerima menggunakan perbandingan waktu konstan:
buat signatureDiterima = signatureDikirim
buat valid = kripto.konstanta_sama(signatureDikirim, signatureDiterima)

jika (valid) {
    cetak("✅ Payload webhook valid dan terotentikasi.")
} selain {
    cetak("❌ Signature tidak cocok!")
}
```

### 3. Pembuatan Token Sesi & Pengenal Acak Kriptografis

Gunakan `kripto.acak_heks` dan `kripto.acak_bita` untuk menghasilkan token sesi, OTP, dan token reset yang tidak dapat ditebak:

```celang
// 16 byte acak = 32 karakter heksadesimal
buat tokenSesi = kripto.acak_heks(16)
cetak("Token sesi:", tokenSesi)

// Objek Bita acak sepanjang 32 bita murni
buat rawBytes = kripto.acak_bita(32)
cetak("Bita acak:", rawBytes)

// Pengenal Unik Global RFC 4122 UUID v4
buat userId = kripto.uuid()
cetak("User ID (UUID v4):", userId)
// Contoh: "4b848092-23c8-47a8-a53b-e85d996d9d1c"
```

---

## 🔒 Praktik Terbaik & Rekomendasi Keamanan

1. **Gunakan `konstanta_sama` untuk Membandingkan Rahasia**:
   Jangan gunakan operator perbandingan biasa `==` saat memverifikasi API key, token rahasia, atau signature HMAC. Operator `==` melakukan *early return* saat menemukan bita pertama yang berbeda, membuka celah kerentanan *timing side-channel*. Gunakan selalu `kripto.konstanta_sama(a, b)` atau `crypto.timing_safe_equal(a, b)`.

2. **Gunakan `kripto.enkripsi` untuk Data yang Disimpan**:
   Jangan membuat skema enkripsi sendiri dengan cipher stream mentah. `kripto.enkripsi` mengimplementasikan AES-256-GCM terotentikasi dengan nonce acak otomatis dan tag verifikasi integritas yang menjamin perlindungan terhadap serangan manipulasi ciphertext.

3. **Gunakan `kripto.hash_sandi` untuk Kredensial Pengguna**:
   Jangan pernah menggunakan fungsi hash digest seperti `kripto.sha256` atau `kripto.md5` untuk menyimpan kata sandi. Gunakan selalu `kripto.hash_sandi` yang menggunakan PBKDF2-HMAC-SHA256 dengan 100.000 iterasi dan salt acak unik.

4. **Simpan Kunci Privat Ed25519 di Lingkungan yang Aman**:
   Kunci privat Ed25519 yang dihasilkan dari `kripto.buat_pasangan_kunci()` tidak boleh disimpan di repositori publik atau dibagikan kepada klien. Cukup bagikan `kunci_publik` untuk proses verifikasi.
