---
outline: deep
---

# Pustaka Kriptografi (`kripto` / `crypto`)

Modul **`kripto`** (atau dengan alias internasional **`crypto`**) menyediakan ekosistem kriptografi berstandar industri modern di Celang. Diimplementasikan murni menggunakan pustaka kriptografi Go (*pure Go*, tanpa dependensi C/eksternal), modul ini dirancang dengan prinsip **aman secara bawaan (*safe-by-default*)** dan terinspirasi dari pola resep tingkat tinggi pada pustaka modern seperti Python `cryptography.io`.

Modul ini mencakup solusi komprehensif mulai dari resep enkripsi terotentikasi, hashing sandi tahan serangan brute-force, tanda tangan digital asimetris, hingga fungsi hash modern dan kriptografi pasca-kuantum (PQC):

1. **Enkripsi Simetris Terotentikasi (AEAD AES-256-GCM & ChaCha20-Poly1305)**: Enkripsi dan dekripsi data aman dengan perlindungan manipulasi otomatis (`buat_kunci`, `enkripsi`, `dekripsi`, `buat_kunci_chacha`, `enkripsi_chacha`, `dekripsi_chacha`).
2. **Password Hashing & Derivasi Kunci Lanjutan (Argon2id, PBKDF2, HKDF)**: Pengamanan kata sandi berstandar industri dengan Argon2id dan PBKDF2, serta RFC 5869 HKDF (*extract-and-expand*) (`hash_sandi_argon2`, `verifikasi_sandi_argon2`, `hkdf`, `hkdf_ekstrak`, `hkdf_ekspansi`).
3. **Kriptografi Asimetris & Pertukaran Kunci (Ed25519, X25519, RSA)**: Tanda tangan digital modern Ed25519, pertukaran kunci ECDH X25519, dan RSA 2048/4096-bit (OAEP enkripsi & PSS tanda tangan).
4. **Fungsi Hash Modern & MAC Kinerja Tinggi (SHA-3, SHAKE, BLAKE2, BLAKE3, Poly1305)**: Keccak SHA-3, SHAKE extensible output, BLAKE2b/s, BLAKE3 paralel, dan Poly1305 MAC satu-kali.
5. **Kriptografi Pasca-Kuantum (Post-Quantum Cryptography - PQC & Hybrid KEM)**: Standar NIST FIPS 203 ML-KEM-768 (Kyber) dan konstruksi hibrida tahan masa depan X-Wing (X25519 + ML-KEM-768).
6. **Entropi & Keacakan Kriptografis**: Mengambil keacakan aman tingkat sistem operasi untuk token sesi dan salt (`acak_heks`, `acak_bita`, `uuid`).
7. **Perlindungan Anti-Timing Attack**: Perbandingan waktu konstan untuk mencegah kebocoran informasi melalui perbedaan latensi CPU (`konstanta_sama`).

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

### 5. Hash Modern & Message Authentication Code (SHA-3, SHAKE, BLAKE2, BLAKE3, Poly1305)

| Fungsi (Bahasa Indonesia) | Alias Universal (Inggris) | Parameter Masukan | Tipe Kembalian | Deskripsi Singkat |
| :--- | :--- | :--- | :--- | :--- |
| `kripto.sha3_256(data)` | `crypto.sha3_256(data)` | `data`: `Teks` / `Bita` | `Teks` (heksadesimal) | Menghitung digest SHA3-256 (FIPS 202 Keccak, 64 heks) |
| `kripto.sha3_512(data)` | `crypto.sha3_512(data)` | `data`: `Teks` / `Bita` | `Teks` (heksadesimal) | Menghitung digest SHA3-512 (FIPS 202 Keccak, 128 heks) |
| `kripto.shake128(data, panjang)` | `crypto.shake128(data, length)` | `data`: `Teks` / `Bita`<br>`panjang`: `Angka` (byte) | `Teks` (heksadesimal) | XOF (Extendable-Output Function) SHAKE128 dengan panjang dinamis |
| `kripto.shake256(data, panjang)` | `crypto.shake256(data, length)` | `data`: `Teks` / `Bita`<br>`panjang`: `Angka` (byte) | `Teks` (heksadesimal) | XOF (Extendable-Output Function) SHAKE256 dengan panjang dinamis |
| `kripto.blake2b(data, [kunci])` | `crypto.blake2b(data, [key])` | `data`: `Teks` / `Bita`<br>`kunci`: `Teks` / `Bita` (opsional) | `Teks` (heksadesimal) | Hash kecepatan tinggi 64-byte (128 heks) RFC 7693 dengan opsi keyed MAC |
| `kripto.blake2s(data, [kunci])` | `crypto.blake2s(data, [key])` | `data`: `Teks` / `Bita`<br>`kunci`: `Teks` / `Bita` (opsional) | `Teks` (heksadesimal) | Hash 32-bit teroptimasi 32-byte (64 heks) RFC 7693 dengan opsi keyed MAC |
| `kripto.blake3(data, [kunci])` | `crypto.blake3(data, [key])` | `data`: `Teks` / `Bita`<br>`kunci`: `Teks` / `Bita` (opsional) | `Teks` (heksadesimal) | Hash pohon Merkle modern super cepat dengan dukungan keyed MAC |
| `kripto.poly1305_buat_mac(pesan, kunci)` | `crypto.poly1305_auth(msg, key)` | `pesan`: `Teks` / `Bita`<br>`kunci`: `Teks` / `Bita` (32 bita) | `Teks` (heksadesimal) | Menghasilkan autentikator Poly1305 16-byte (32 heks) satu-kali (RFC 8439) |
| `kripto.poly1305_verifikasi(pesan, tag, kunci)` | `crypto.poly1305_verify(msg, tag, key)` | `pesan`: `Teks` / `Bita`<br>`tag`: `Teks` / `Bita`<br>`kunci`: `Teks` / `Bita` | `Benar` / `Salah` | Memverifikasi tag autentikasi Poly1305 dalam waktu konstan |

### 6. Derivasi Kunci Lanjutan & Password Hashing Modern (HKDF & Argon2id)

| Fungsi (Bahasa Indonesia) | Alias Universal (Inggris) | Parameter Masukan | Tipe Kembalian | Deskripsi Singkat |
| :--- | :--- | :--- | :--- | :--- |
| `kripto.hkdf(secret, salt, info, panjang)` | `crypto.hkdf(...)` | `secret`, `salt`, `info`: `Teks`/`Bita`<br>`panjang`: `Angka` | `Bita` | Menurunkan kunci biner berkualitas tinggi menggunakan RFC 5869 HKDF |
| `kripto.hkdf_ekstrak(secret, salt)` | `crypto.hkdf_extract(...)` | `secret`, `salt`: `Teks` / `Bita` | `Bita` | Tahap ekstraksi HKDF menghasilkan Pseudorandom Key (PRK) 32-byte |
| `kripto.hkdf_ekspansi(prk, info, panjang)` | `crypto.hkdf_expand(...)` | `prk`, `info`: `Teks`/`Bita`<br>`panjang`: `Angka` | `Bita` | Tahap ekspansi HKDF memperluas PRK menjadi material kunci target |
| `kripto.hash_sandi_argon2(sandi, [time, mem, threads])` | `crypto.hash_password_argon2(...)` | `sandi`: `Teks`/`Bita`<br>opsional cost parameters | `Teks` (PHC) | Hash sandi standar industri Argon2id tahan ASIC/GPU memory-hard |
| `kripto.verifikasi_sandi_argon2(sandi, hash)` | `crypto.verify_password_argon2(...)` | `sandi`: `Teks`/`Bita`<br>`hash`: `Teks` | `Benar` / `Salah` | Memverifikasi kata sandi terhadap string hash format PHC `$argon2id$` |
| `kripto.argon2id_turunkan_kunci(sandi, salt, time, mem, threads, panjang)` | `crypto.argon2id_derive_key(...)` | `sandi`, `salt`<br>`time`, `mem`, `threads`, `panjang` | `Bita` | Menurunkan kunci enkripsi dari sandi menggunakan Argon2id memory-hard |

### 7. Enkripsi Simetris ChaCha20-Poly1305 & Mode Blok AES

| Fungsi (Bahasa Indonesia) | Alias Universal (Inggris) | Parameter Masukan | Tipe Kembalian | Deskripsi Singkat |
| :--- | :--- | :--- | :--- | :--- |
| `kripto.buat_kunci_chacha()` | `crypto.generate_chacha_key()` | `-` | `Teks` (Base64) | Menghasilkan kunci 32-byte acak format Base64 untuk ChaCha20 |
| `kripto.enkripsi_chacha(data, kunci, [aad])` | `crypto.encrypt_chacha(...)` | `data`, `kunci`: `Teks`/`Bita`<br>`aad`: `Teks`/`Bita` (opsional) | `Teks` (Base64) | Enkripsi terotentikasi RFC 8439 ChaCha20-Poly1305 dengan nonce 12-byte otomatis |
| `kripto.dekripsi_chacha(ciphertext, kunci, [aad])` | `crypto.decrypt_chacha(...)` | `ciphertext`, `kunci`<br>`aad` (opsional) | `Teks` | Mendekripsi dan memverifikasi ciphertext ChaCha20-Poly1305 |
| `kripto.aes_ctr(data, kunci, iv)` | `crypto.aes_ctr(data, key, iv)` | `data`: `Teks`/`Bita`<br>`kunci`: 16/24/32 bita<br>`iv`: 16 bita | `Bita` | Enkripsi/dekripsi simetris AES mode Counter (CTR) stream cipher |
| `kripto.aes_cbc_enkripsi(data, kunci, iv)` | `crypto.aes_cbc_encrypt(...)` | `data`: `Teks`/`Bita`<br>`kunci`: 16/24/32 bita<br>`iv`: 16 bita | `Bita` | Enkripsi AES mode CBC dengan padding otomatis standar PKCS#7 |
| `kripto.aes_cbc_dekripsi(ciphertext, kunci, iv)` | `crypto.aes_cbc_decrypt(...)` | `ciphertext`, `kunci`, `iv` | `Bita` | Dekripsi AES mode CBC dengan validasi dan pemotongan padding PKCS#7 |

### 8. Kriptografi Asimetris & Pertukaran Kunci (X25519 & RSA)

| Fungsi (Bahasa Indonesia) | Alias Universal (Inggris) | Parameter Masukan | Tipe Kembalian | Deskripsi Singkat |
| :--- | :--- | :--- | :--- | :--- |
| `kripto.buat_kunci_x25519()` | `crypto.generate_x25519_keypair()` | `-` | `Peta` (`Map`) | Menghasilkan pasangan kunci ECDH X25519 (`kunci_publik`, `kunci_privat`) |
| `kripto.x25519_hitung_rahasia(kunci_privat, kunci_publik_lawan)` | `crypto.x25519_shared_secret(...)` | `kunci_privat`, `kunci_publik_lawan` | `Bita` | Menghitung shared secret Diffie-Hellman 32-byte berkecepatan tinggi |
| `kripto.rsa_buat_kunci([bits])` | `crypto.rsa_generate_keypair([bits])` | `bits`: `Angka` (opsional, default 2048) | `Peta` (`Map`) | Menghasilkan pasangan kunci RSA dalam format standar PEM PKCS#1 / PKIX |
| `kripto.rsa_enkripsi_oaep(data, kunci_publik_pem)` | `crypto.rsa_encrypt_oaep(...)` | `data`: `Teks`/`Bita`<br>`kunci_publik_pem`: `Teks` | `Teks` (Base64) | Mengenkripsi data menggunakan RSA-OAEP (SHA-256 MGF1) |
| `kripto.rsa_dekripsi_oaep(ciphertext, kunci_privat_pem)` | `crypto.rsa_decrypt_oaep(...)` | `ciphertext`: `Teks`<br>`kunci_privat_pem`: `Teks` | `Teks` | Mendekripsi data menggunakan RSA-OAEP private key |
| `kripto.rsa_tanda_tangan_pss(data, kunci_privat_pem)` | `crypto.rsa_sign_pss(...)` | `data`: `Teks`/`Bita`<br>`kunci_privat_pem`: `Teks` | `Teks` (Base64) | Menghasilkan tanda tangan digital probabilistik standar RSA-PSS (SHA-256) |
| `kripto.rsa_verifikasi_pss(data, signature, kunci_publik_pem)` | `crypto.rsa_verify_pss(...)` | `data`, `signature`<br>`kunci_publik_pem`: `Teks` | `Benar` / `Salah` | Memvalidasi keabsahan tanda tangan digital RSA-PSS |

### 9. Kriptografi Pasca-Kuantum & Hybrid KEM (NIST FIPS 203 & X-Wing)

| Fungsi (Bahasa Indonesia) | Alias Universal (Inggris) | Parameter Masukan | Tipe Kembalian | Deskripsi Singkat |
| :--- | :--- | :--- | :--- | :--- |
| `kripto.pqc_buat_kunci()` | `crypto.pqc_generate_keypair()` | `-` | `Peta` (`Map`) | Pembangkit kunci NIST FIPS 203 ML-KEM-768 (`kunci_enkapsulasi`, `kunci_dekapsulasi`) |
| `kripto.pqc_enkapsulasi(kunci_enkapsulasi)` | `crypto.pqc_encapsulate(...)` | `kunci_enkapsulasi`: `Teks`/`Bita` | `Peta` (`Map`) | Mengenkapsulasi kunci bersama pasca-kuantum: mengembalikan `{ ciphertext, kunci_bersama }` |
| `kripto.pqc_dekapsulasi(ciphertext, kunci_dekapsulasi)` | `crypto.pqc_decapsulate(...)` | `ciphertext`, `kunci_dekapsulasi` | `Bita` | Mendekapsulasi shared secret 32-byte dari ciphertext ML-KEM-768 |
| `kripto.xwing_buat_kunci()` | `crypto.xwing_generate_keypair()` | `-` | `Peta` (`Map`) | Pembangkit kunci hibrida X-Wing (kombinasi X25519 klasik + ML-KEM-768 kuantum) |
| `kripto.xwing_enkapsulasi(kunci_publik)` | `crypto.xwing_encapsulate(...)` | `kunci_publik`: `Teks`/`Bita` | `Peta` (`Map`) | Enkapsulasi kunci hibrida tahan kuantum: mengembalikan `{ ciphertext, kunci_bersama }` |
| `kripto.xwing_dekapsulasi(ciphertext, kunci_privat)` | `crypto.xwing_decapsulate(...)` | `ciphertext`, `kunci_privat` | `Bita` | Dekapsulasi shared secret hibrida 32-byte dari ciphertext X-Wing |


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

### 4. Hash Kinerja Tinggi Modern (SHA-3, BLAKE2, BLAKE3, dan Poly1305)

Celang menyediakan fungsi hash modern dan message authentication code dengan perlindungan terhadap serangan length extension:

```celang
buat payload = "Pesan berharga untuk diverifikasi"

// 1. Keccak SHA-3 & SHAKE XOF
buat h3_256 = kripto.sha3_256(payload)
buat h3_512 = kripto.sha3_512(payload)
buat shake = kripto.shake128(payload, 32) // Output 32 bita dinamis

// 2. BLAKE2 & BLAKE3
buat b2b = kripto.blake2b(payload)
buat b3 = kripto.blake3(payload)

// Keyed MAC menggunakan BLAKE3 (32 bita)
buat kunciMac = "kunci-rahasia-32-karakter-tepat!"
buat macB3 = kripto.blake3(payload, kunciMac)

// 3. Poly1305 MAC satu-kali
buat kunciPoly = kripto.buat_kunci()
buat tagPoly = kripto.poly1305_buat_mac(payload, kunciPoly)
buat tagCocok = kripto.poly1305_verifikasi(payload, tagPoly, kunciPoly)
```

### 5. Password Hashing Modern (Argon2id) & Key Derivation (HKDF)

Untuk pengamanan kata sandi berstandar ketat yang tahan serangan komputasi GPU/ASIC:

```celang
// 1. Argon2id Password Hashing (Format Standar PHC)
buat sandi = "KatasandiKuat#2026!"
buat hashArgon = kripto.hash_sandi_argon2(sandi)
// Output: "$argon2id$v=19$m=65536,t=3,p=4$..."

// Opsi kustom dapat diberikan via Peta (Map):
buat hashKustom = kripto.hash_sandi_argon2(sandi, {"waktu": 2, "memori": 65536, "threads": 2})

buat sukses = kripto.verifikasi_sandi_argon2(sandi, hashArgon)
cetak("Verifikasi sandi cocok:", sukses) // benar

// 2. HKDF Key Derivation (RFC 5869 Extract-and-Expand)
// Sesuai RFC 5869 (SHA-256), panjang kunci turunan dibatasi hingga 255 * 32 = 8.160 bita.
buat masterEntropy = "seed-entropy-sangat-rahasia"
buat salt = "garam-unik-sesi-123"
buat appKey = kripto.hkdf(masterEntropy, salt, "celang-app-v1", 32)
```

> [!NOTE]
> **Perlindungan DoS pada Verifikasi Argon2**:
> Saat memverifikasi hash dari input luar, `kripto.verifikasi_sandi_argon2` membatasi alokasi memori maksimal 256 MB (`256 * 1024` KB) dan waktu iterasi $t \le 10$ guna mencegah kelumpuhan server akibat serangan *Denial of Service* (DoS / OOM).

### 6. Enkripsi Simetris XChaCha20-Poly1305 & Mode Blok AES

Pilihan enkripsi simetris modern tanpa akselerasi perangkat keras khusus serta mode standar:

```celang
buat pesan = "Payload rahasia protokol Celang"

// 1. XChaCha20-Poly1305 AEAD (Nonce 24-byte / 192-bit)
// Menggunakan nonce acak 24-byte untuk mengeliminasi risiko tabrakan nonce (birthday bound).
// Fungsi dekripsi bersifat adaptif: mendukung ciphertext 24-byte baru dan kompatibel mundur dengan 12-byte lama.
buat kunciChaCha = kripto.buat_kunci_chacha()
buat aad = "header-metadata-id:101"
buat ciphertext = kripto.enkripsi_chacha(pesan, kunciChaCha, aad)
buat plaintext = kripto.dekripsi_chacha(ciphertext, kunciChaCha, aad)

// 2. Mode Blok AES (CTR & CBC PKCS#7)
// AES-CBC Celang menggunakan unpadding PKCS#7 dalam waktu konstan (constant-time)
// untuk memitigasi serangan padding oracle.
buat kunciAES = "01234567890123456789012345678901" // 32 bita (AES-256)
buat iv = "0123456789012345"                     // 16 bita
buat encCTR = kripto.aes_ctr(pesan, kunciAES, iv)
buat decCTR = kripto.aes_ctr(encCTR, kunciAES, iv)

buat encCBC = kripto.aes_cbc_enkripsi(pesan, kunciAES, iv)
buat decCBC = kripto.aes_cbc_dekripsi(encCBC, kunciAES, iv)
```

### 7. Pertukaran Kunci Asimetris (X25519 & RSA)

```celang
// 1. Pertukaran Kunci X25519 ECDH (Alice & Bob)
buat alice = kripto.buat_kunci_x25519()
buat bob = kripto.buat_kunci_x25519()

buat rahasiaAlice = kripto.x25519_hitung_rahasia(alice.kunci_privat, bob.kunci_publik)
buat rahasiaBob = kripto.x25519_hitung_rahasia(bob.kunci_privat, alice.kunci_publik)
buat rahasiaCocok = kripto.konstanta_sama(rahasiaAlice, rahasiaBob) // benar!

// 2. RSA (2048/4096-bit) Enkripsi OAEP & Tanda Tangan PSS
buat pasanganRSA = kripto.rsa_buat_kunci(2048)
buat encRSA = kripto.rsa_enkripsi_oaep("Pesan Penting", pasanganRSA.kunci_publik)
buat decRSA = kripto.rsa_dekripsi_oaep(encRSA, pasanganRSA.kunci_privat)

// Penandatanganan PSS menggunakan SaltLength deterministik sesuai hash (SHA-256).
// Verifikasi PSS mendukung SaltLengthAuto (menerima signature hex maupun Base64 dari OpenSSL, Java, Python).
buat sigRSA = kripto.rsa_tanda_tangan_pss("Dokumen Sah", pasanganRSA.kunci_privat)
buat validRSA = kripto.rsa_verifikasi_pss("Dokumen Sah", sigRSA, pasanganRSA.kunci_publik)
```

### 8. Kriptografi Pasca-Kuantum (Post-Quantum ML-KEM-768 & Hybrid X-Wing)

Melindungi transmisi data jangka panjang dari ancaman komputer kuantum masa depan (*Harvest Now, Decrypt Later*):

```celang
// 1. NIST FIPS 203 ML-KEM-768 (Kyber)
buat kunciPQC = kripto.pqc_buat_kunci()

// Pengirim mengenkapsulasi kunci bersama menggunakan kunci publik penerima
buat hasilEnkapsulasi = kripto.pqc_enkapsulasi(kunciPQC.kunci_enkapsulasi)
buat kunciBersamaPengirim = hasilEnkapsulasi.kunci_bersama
buat ciphertextKuantum = hasilEnkapsulasi.ciphertext

// Penerima mendekapsulasi ciphertext menggunakan kunci dekapsulasi privat
buat kunciBersamaPenerima = kripto.pqc_dekapsulasi(ciphertextKuantum, kunciPQC.kunci_dekapsulasi)
buat cocokPQC = kripto.konstanta_sama(kunciBersamaPengirim, kunciBersamaPenerima) // benar!

// 2. Hybrid KEM: X-Wing (X25519 Klasik + ML-KEM-768 Pasca-Kuantum)
// Rekomendasi terkuat: Menggabungkan keamanan kriptografi klasik dan pasca-kuantum secara simultan
buat kunciXWing = kripto.xwing_buat_kunci()
buat resXW = kripto.xwing_enkapsulasi(kunciXWing.kunci_publik)
buat kunciXWPenerima = kripto.xwing_dekapsulasi(resXW.ciphertext, kunciXWing.kunci_privat)
buat cocokXW = kripto.konstanta_sama(resXW.kunci_bersama, kunciXWPenerima) // benar!
```

---

## 🔒 Praktik Terbaik & Rekomendasi Keamanan

1. **Gunakan `konstanta_sama` untuk Membandingkan Rahasia**:
   Jangan gunakan operator perbandingan biasa `==` saat memverifikasi API key, token rahasia, atau signature HMAC. Operator `==` melakukan *early return* saat menemukan bita pertama yang berbeda, membuka celah kerentanan *timing side-channel*. Gunakan selalu `kripto.konstanta_sama(a, b)` atau `crypto.timing_safe_equal(a, b)`.

2. **Gunakan `kripto.enkripsi` untuk Data yang Disimpan**:
   Jangan membuat skema enkripsi sendiri dengan cipher stream mentah. `kripto.enkripsi` mengimplementasikan AES-256-GCM terotentikasi dengan nonce acak otomatis dan tag verifikasi integritas yang menjamin perlindungan terhadap serangan manipulasi ciphertext.

3. **Gunakan `kripto.hash_sandi` untuk Kredensial Pengguna**:
   Jangan pernah menggunakan fungsi hash digest seperti `kripto.sha256` atau `kripto.md5` untuk menyimpan kata sandi. Gunakan selalu `kripto.hash_sandi` yang menggunakan PBKDF2-HMAC-SHA256 dengan 100.000 iterasi dan salt acak unik.

4. **Simpan Kunci Privat Ed25519 & RSA di Lingkungan yang Aman**:
   Kunci privat yang dihasilkan dari `kripto.buat_pasangan_kunci()` atau `kripto.rsa_buat_kunci()` tidak boleh disimpan di repositori publik atau dibagikan kepada klien. Cukup bagikan `kunci_publik` untuk proses verifikasi.

5. **Utamakan Enkripsi Terotentikasi (AEAD) Dibandingkan Mode Mentah (CBC / CTR)**:
   Mode `aes_cbc_enkripsi` dan `aes_ctr` tidak menyertakan kode autentikasi pesan bawaan. Penggunaan CBC mentah di lingkungan jaringan tanpa otentikasi dapat membuka risiko serangan *padding oracle*. Selalu utamakan `kripto.enkripsi` (AES-256-GCM) atau `kripto.enkripsi_chacha` (ChaCha20-Poly1305). Jika menggunakan CBC/CTR, selalu pasangkan dengan pola *Encrypt-then-MAC* (HMAC-SHA256).

6. **Kunci Poly1305 Bersifat Sekali-Pakai (*One-Time Key*)**:
   Jangan pernah menggunakan kunci 32-byte yang sama untuk membuat tag Poly1305 pada lebih dari satu pesan. Menggunakan kembali kunci pada Poly1305 mematahkan integritas otentikasi. Jika membutuhkan MAC dengan kunci statis berulang kali, gunakan `kripto.hmac_sha256` atau keyed `kripto.blake3(pesan, kunci)`.

7. **Turunkan Kunci dari `x25519_hitung_rahasia` Melalui HKDF**:
   Hasil *shared secret* mentah dari pertukaran kunci Diffie-Hellman X25519 sebaiknya tidak langsung digunakan sebagai kunci enkripsi simetris. Gunakan fungsi derivasi kunci `kripto.hkdf(sharedSecret, salt, "tujuan-protokol", 32)` untuk mengekstraksi material kunci enkripsi simetris yang seragam.
