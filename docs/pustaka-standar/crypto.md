---
outline: deep
---

# Cryptography Library (`crypto` / `kripto`)

The **`crypto`** module (official Indonesian name **`kripto`**) provides a modern, production-grade cryptography ecosystem in Celang. Built with pure Go standard library primitives (zero external dependencies and no CGo bindings), it is designed around the **safe-by-default** philosophy and inspired by modern high-level cryptographic libraries like Python's `cryptography.io`.

The library provides comprehensive cryptographic functionality ranging from authenticated encryption and GPU-resistant password hashing to asymmetric Ed25519 signatures, standard digests, modern hashes, and post-quantum cryptography (PQC):

1. **Authenticated Symmetric Encryption (AEAD AES-256-GCM & ChaCha20-Poly1305)**: Safe symmetric encryption with automatic nonce generation and built-in tamper detection (`generate_key`, `encrypt`, `decrypt`, `generate_chacha_key`, `encrypt_chacha`, `decrypt_chacha`).
2. **Password Hashing & Key Derivation (Argon2id, PBKDF2, HKDF)**: Memory-hard Argon2id and standard PBKDF2 password storage, plus RFC 5869 HKDF (*extract-and-expand*) (`hash_password_argon2`, `verify_password_argon2`, `hkdf`, `hkdf_extract`, `hkdf_expand`).
3. **Asymmetric Cryptography & Key Exchange (Ed25519, X25519, RSA)**: Modern Ed25519 digital signatures, high-speed X25519 ECDH key exchange, and enterprise RSA 2048/4096-bit (OAEP encryption & PSS signatures).
4. **Modern Hashes & Message Authentication (SHA-3, SHAKE, BLAKE2, BLAKE3, Poly1305)**: FIPS 202 SHA-3, SHAKE extensible-output functions, BLAKE2b/s, ultra-fast BLAKE3 tree hashing, and Poly1305 one-time authenticator.
5. **Post-Quantum Cryptography & Hybrid KEM (ML-KEM-768 & Hybrid X-Wing)**: NIST FIPS 203 ML-KEM-768 (Kyber) quantum-resistant KEM and future-proof Hybrid X-Wing (X25519 + ML-KEM-768).
6. **Cryptographic Randomness & Identifiers**: Operating-system-backed CSPRNG entropy for session tokens and salt generation (`random_hex`, `random_bytes`, `uuid`).
7. **Anti-Timing Attack Protection**: Constant-time byte and string comparison to defend against timing side-channel attacks (`timing_safe_equal`).

> [!TIP]
> Looking for the Indonesian documentation? Read **[Pustaka Kriptografi (`kripto`)](/pustaka-standar/kripto)**.

---

## 🌍 Universal Mode (`"use universal"`) Compatibility

Celang provides a strict `"use universal"` directive designed for international teams, global open-source libraries, and automated CI/CD pipelines.

When `"use universal"` is placed at the very first line of a `.ce` file:
- **Module Name Requirement**: You **must** import or reference the module as `"crypto"` (e.g., `import { encrypt, decrypt } from "crypto"` or `crypto.encrypt(...)`).
- **Empathetic Diagnostic**: Referencing the Indonesian module name `kripto` is strictly intercepted by Celang's compiler and reports an empathetic error advising the developer to use `crypto`.
- **English Language Syntax**: Code written in Universal Mode uses standard English keywords (`let`, `const`, `if`, `else`, `fn` / `def`, `return`, `try`, `catch`, `finally`, `true`, `false`, `null`) along with English standard library aliases.

```celang
"use universal"
import { generate_key, encrypt, decrypt } from "crypto"

let key = generate_key()
let ciphertext = encrypt("Confidential payload", key)
let recovered = decrypt(ciphertext, key)
```

---

## 📑 Complete API Reference

### 1. Symmetric Encryption Recipes (AES-256-GCM AEAD)

| Function (English) | Indonesian Alias | Arguments | Return Type | Description |
| :--- | :--- | :--- | :--- | :--- |
| `crypto.generate_key()` | `kripto.buat_kunci()` | `-` | `String` (Base64) | Generates a cryptographically secure 256-bit (32-byte) key encoded in Base64 |
| `crypto.encrypt(data, key)` | `kripto.enkripsi(data, kunci)` | `data`: `String` / `Bytes`<br>`key`: `String` / `Bytes` | `String` (Base64) | Encrypts data using AES-256-GCM with an automatic 12-byte CSPRNG nonce and 16-byte authentication tag |
| `crypto.decrypt(ciphertext, key)` | `kripto.dekripsi(ciphertext, kunci)` | `ciphertext`: `String` / `Bytes`<br>`key`: `String` / `Bytes` | `String` | Decrypts ciphertext, verifies authentication tag, or safely raises a friendly error on tampering / wrong key |

### 2. Password Hashing & Key Derivation (PBKDF2-HMAC-SHA256)

| Function (English) | Indonesian Alias | Arguments | Return Type | Description |
| :--- | :--- | :--- | :--- | :--- |
| `crypto.hash_password(password)` | `kripto.hash_sandi(sandi)` | `password`: `String` / `Bytes` | `String` (PHC) | Hashes password with 100,000 PBKDF2-HMAC-SHA256 iterations and random 16-byte salt in standard PHC format |
| `crypto.verify_password(password, hash)` | `kripto.verifikasi_sandi(sandi, hash)` | `password`: `String` / `Bytes`<br>`hash`: `String` | `Boolean` (`true` / `false`) | Verifies password against PHC string in constant time (resistant to timing attacks) |
| `crypto.derive_key(password, salt, len, [format])` | `kripto.turunkan_kunci(sandi, salt, panjang, [format])` | `password`: `String` / `Bytes`<br>`salt`: `String` / `Bytes`<br>`len`: `Integer`<br>`format`: `String` (optional) | `Bytes` / `String` | Derives symmetric key from passphrase and salt. Output format: `"bytes"` (default), `"hex"`, or `"base64"` |

### 3. Asymmetric Digital Signatures (Ed25519 - RFC 8032)

| Function (English) | Indonesian Alias | Arguments | Return Type | Description |
| :--- | :--- | :--- | :--- | :--- |
| `crypto.generate_keypair()` | `kripto.buat_pasangan_kunci()` | `-` | `Map` | Generates Ed25519 keypair containing `public_key` (64 hex) and `private_key` (64 hex seed) |
| `crypto.sign(data, private_key)` | `kripto.tanda_tangan(data, kunci_privat)` | `data`: `String` / `Bytes`<br>`private_key`: `String` / `Bytes` | `String` (hex) | Computes a 64-byte Ed25519 digital signature (128 hex characters) |
| `crypto.verify(data, signature, public_key)` | `kripto.verifikasi_tanda_tangan(data, signature, kunci_publik)` | `data`: `String` / `Bytes`<br>`signature`: `String` / `Bytes`<br>`public_key`: `String` / `Bytes` | `Boolean` (`true` / `false`) | Verifies digital signature authenticity against sender's public key (alias: `verify_signature`) |

### 4. Hashes, HMAC, Randomness, & Timing-Safe Utilities

| Function (English) | Indonesian Alias | Arguments | Return Type | Description |
| :--- | :--- | :--- | :--- | :--- |
| `crypto.sha256(data)` | `kripto.sha256(data)` | `data`: `String` / `Bytes` | `String` (hex) | Computes SHA-256 digest (64 hex characters) |
| `crypto.sha256_file(path)` | `kripto.sha256_berkas(jalur)` | `path`: `String` | `String` (hex) | Computes SHA-256 directly from file stream with zero memory bloat |
| `crypto.sha512(data)` | `kripto.sha512(data)` | `data`: `String` / `Bytes` | `String` (hex) | Computes SHA-512 digest (128 hex characters) |
| `crypto.sha512_file(path)` | `kripto.sha512_berkas(jalur)` | `path`: `String` | `String` (hex) | Computes SHA-512 directly from file stream |
| `crypto.sha1(data)` | `kripto.sha1(data)` | `data`: `String` / `Bytes` | `String` (hex) | Computes SHA-1 digest (legacy compatibility, 40 hex characters) |
| `crypto.md5(data)` | `kripto.md5(data)` | `data`: `String` / `Bytes` | `String` (hex) | Computes MD5 digest (legacy compatibility, 32 hex characters) |
| `crypto.md5_file(path)` | `kripto.md5_berkas(jalur)` | `path`: `String` | `String` (hex) | Computes MD5 directly from file stream |
| `crypto.crc32(data)` | `kripto.crc32(data)` | `data`: `String` / `Bytes` | `Integer` (UInt32) | Computes fast IEEE CRC-32 checksum (ZIP and PNG standard) |
| `crypto.crc32_file(path)` | `kripto.crc32_berkas(jalur)` | `path`: `String` | `Integer` (UInt32) | Computes IEEE CRC-32 directly from file stream |
| `crypto.hmac(key, data)` | `kripto.hmac_sha256(kunci, data)` | `key`: `String` / `Bytes`<br>`data`: `String` / `Bytes` | `String` (hex) | Computes HMAC-SHA256 message authentication code (alias: `hmac_sha256`) |
| `crypto.random_hex(bytes)` | `kripto.acak_heks(panjang)` | `bytes`: `Integer` | `String` (hex) | Generates CSPRNG random hex string (`2 * bytes` hex characters) |
| `crypto.random_bytes(bytes)` | `kripto.acak_bita(panjang)` | `bytes`: `Integer` | `Bytes` | Generates raw `Bytes` buffer containing CSPRNG random entropy |
| `crypto.uuid()` | `kripto.uuid()` | `-` | `String` (UUID v4) | Generates RFC 4122 Version 4 random UUID (aliases: `uuid_v4`, `uuid4`) |
| `crypto.timing_safe_equal(a, b)` | `kripto.konstanta_sama(a, b)` | `a`, `b`: `String` / `Bytes` | `Boolean` (`true` / `false`) | Compares two secrets in constant CPU time to prevent timing attacks |

### 5. Modern Hashes & Message Authentication Codes (SHA-3, SHAKE, BLAKE2, BLAKE3, Poly1305)

| Function (English) | Indonesian Alias | Arguments | Return Type | Description |
| :--- | :--- | :--- | :--- | :--- |
| `crypto.sha3_256(data)` | `kripto.sha3_256(data)` | `data`: `String` / `Bytes` | `String` (hex) | Computes SHA3-256 digest (FIPS 202 Keccak, 64 hex characters) |
| `crypto.sha3_512(data)` | `kripto.sha3_512(data)` | `data`: `String` / `Bytes` | `String` (hex) | Computes SHA3-512 digest (FIPS 202 Keccak, 128 hex characters) |
| `crypto.shake128(data, length)` | `kripto.shake128(data, panjang)` | `data`: `String` / `Bytes`<br>`length`: `Integer` (bytes) | `String` (hex) | Extendable-Output Function (XOF) SHAKE128 with custom output length |
| `crypto.shake256(data, length)` | `kripto.shake256(data, panjang)` | `data`: `String` / `Bytes`<br>`length`: `Integer` (bytes) | `String` (hex) | Extendable-Output Function (XOF) SHAKE256 with custom output length |
| `crypto.blake2b(data, [key])` | `kripto.blake2b(data, [kunci])` | `data`: `String` / `Bytes`<br>`key`: `String` / `Bytes` (optional) | `String` (hex) | High-speed 64-byte (128 hex) hash (RFC 7693) with optional keyed MAC |
| `crypto.blake2s(data, [key])` | `kripto.blake2s(data, [kunci])` | `data`: `String` / `Bytes`<br>`key`: `String` / `Bytes` (optional) | `String` (hex) | 32-bit optimized 32-byte (64 hex) hash (RFC 7693) with optional keyed MAC |
| `crypto.blake3(data, [key])` | `kripto.blake3(data, [kunci])` | `data`: `String` / `Bytes`<br>`key`: `String` / `Bytes` (optional) | `String` (hex) | Ultra-fast parallel Merkle-tree hash with keyed MAC support |
| `crypto.poly1305_auth(msg, key)` | `kripto.poly1305_buat_mac(pesan, kunci)` | `msg`: `String` / `Bytes`<br>`key`: `String` / `Bytes` (32 bytes) | `String` (hex) | Generates a 16-byte (32 hex) one-time Poly1305 authenticator (RFC 8439) |
| `crypto.poly1305_verify(msg, tag, key)` | `kripto.poly1305_verifikasi(...)` | `msg`, `tag`, `key` | `Boolean` (`true` / `false`) | Verifies Poly1305 authentication tag in constant CPU time |

### 6. Advanced Key Derivation & Modern Password Hashing (HKDF & Argon2id)

| Function (English) | Indonesian Alias | Arguments | Return Type | Description |
| :--- | :--- | :--- | :--- | :--- |
| `crypto.hkdf(secret, salt, info, length)` | `kripto.hkdf(...)` | `secret`, `salt`, `info`: `String`/`Bytes`<br>`length`: `Integer` | `Bytes` | Derives cryptographic keys using RFC 5869 HKDF (*extract-and-expand*) |
| `crypto.hkdf_extract(secret, salt)` | `kripto.hkdf_ekstrak(...)` | `secret`, `salt`: `String` / `Bytes` | `Bytes` | HKDF extract step producing a 32-byte Pseudorandom Key (PRK) |
| `crypto.hkdf_expand(prk, info, length)` | `kripto.hkdf_ekspansi(...)` | `prk`, `info`: `String`/`Bytes`<br>`length`: `Integer` | `Bytes` | HKDF expand step producing derived key bytes of requested length |
| `crypto.hash_password_argon2(password, ...)` | `kripto.hash_sandi_argon2(...)` | `password`: `String`/`Bytes`<br>optional cost parameters | `String` (PHC) | Industry-standard memory-hard Argon2id password hash ($argon2id$ format) |
| `crypto.verify_password_argon2(pwd, hash)` | `kripto.verifikasi_sandi_argon2(...)` | `pwd`: `String`/`Bytes`<br>`hash`: `String` | `Boolean` (`true` / `false`) | Verifies password against Argon2id PHC string in constant time |
| `crypto.argon2id_derive_key(...)` | `kripto.argon2id_turunkan_kunci(...)` | `password`, `salt`<br>`time`, `mem`, `threads`, `length` | `Bytes` | Derives symmetric key using memory-hard Argon2id KDF |

### 7. Symmetric ChaCha20-Poly1305 AEAD & AES Block Modes

| Function (English) | Indonesian Alias | Arguments | Return Type | Description |
| :--- | :--- | :--- | :--- | :--- |
| `crypto.generate_chacha_key()` | `kripto.buat_kunci_chacha()` | `-` | `String` (Base64) | Generates a cryptographically secure 32-byte Base64 key for ChaCha20 |
| `crypto.encrypt_chacha(data, key, [aad])` | `kripto.enkripsi_chacha(...)` | `data`, `key`: `String`/`Bytes`<br>`aad`: `String`/`Bytes` (optional) | `String` (Base64) | Authenticated encryption using RFC 8439 ChaCha20-Poly1305 with automatic nonce |
| `crypto.decrypt_chacha(ciphertext, key, [aad])` | `kripto.dekripsi_chacha(...)` | `ciphertext`, `key`<br>`aad` (optional) | `String` | Decrypts and authenticates ChaCha20-Poly1305 ciphertext |
| `crypto.aes_ctr(data, key, iv)` | `kripto.aes_ctr(data, kunci, iv)` | `data`: `String`/`Bytes`<br>`key`: 16/24/32 bytes<br>`iv`: 16 bytes | `Bytes` | AES Counter (CTR) mode stream encryption / decryption |
| `crypto.aes_cbc_encrypt(data, key, iv)` | `kripto.aes_cbc_enkripsi(...)` | `data`: `String`/`Bytes`<br>`key`: 16/24/32 bytes<br>`iv`: 16 bytes | `Bytes` | AES CBC encryption with automatic standard PKCS#7 padding |
| `crypto.aes_cbc_decrypt(ciphertext, key, iv)` | `kripto.aes_cbc_dekripsi(...)` | `ciphertext`, `key`, `iv` | `Bytes` | AES CBC decryption with PKCS#7 padding validation and strip |

### 8. Asymmetric Cryptography & Key Exchange (X25519 & RSA)

| Function (English) | Indonesian Alias | Arguments | Return Type | Description |
| :--- | :--- | :--- | :--- | :--- |
| `crypto.generate_x25519_keypair()` | `kripto.buat_kunci_x25519()` | `-` | `Map` | Generates X25519 ECDH keypair (`public_key`, `private_key`) |
| `crypto.x25519_shared_secret(priv, pub)` | `kripto.x25519_hitung_rahasia(...)` | `priv`: private key<br>`pub`: peer public key | `Bytes` | Computes 32-byte Diffie-Hellman shared secret over Curve25519 |
| `crypto.rsa_generate_keypair([bits])` | `kripto.rsa_buat_kunci([bits])` | `bits`: `Integer` (optional, default 2048) | `Map` | Generates RSA keypair in standard PKCS#1 / PKIX PEM format |
| `crypto.rsa_encrypt_oaep(data, pub_pem)` | `kripto.rsa_enkripsi_oaep(...)` | `data`: `String`/`Bytes`<br>`pub_pem`: `String` | `String` (Base64) | Encrypts data using RSA-OAEP with SHA-256 and MGF1 |
| `crypto.rsa_decrypt_oaep(cipher, priv_pem)` | `kripto.rsa_dekripsi_oaep(...)` | `cipher`: `String`<br>`priv_pem`: `String` | `String` | Decrypts data using RSA-OAEP private key |
| `crypto.rsa_sign_pss(data, priv_pem)` | `kripto.rsa_tanda_tangan_pss(...)` | `data`: `String`/`Bytes`<br>`priv_pem`: `String` | `String` (Base64) | Computes probabilistic RSA-PSS digital signature (SHA-256) |
| `crypto.rsa_verify_pss(data, sig, pub_pem)` | `kripto.rsa_verifikasi_pss(...)` | `data`, `sig`<br>`pub_pem`: `String` | `Boolean` (`true` / `false`) | Verifies RSA-PSS digital signature against sender's public key |

### 9. Post-Quantum Cryptography & Hybrid KEM (NIST FIPS 203 & X-Wing)

| Function (English) | Indonesian Alias | Arguments | Return Type | Description |
| :--- | :--- | :--- | :--- | :--- |
| `crypto.pqc_generate_keypair()` | `kripto.pqc_buat_kunci()` | `-` | `Map` | Generates NIST FIPS 203 ML-KEM-768 keypair (`encapsulation_key`, `decapsulation_key`) |
| `crypto.pqc_encapsulate(ek)` | `kripto.pqc_enkapsulasi(...)` | `ek`: encapsulation key | `Map` | Encapsulates quantum-resistant shared secret: returns `{ ciphertext, shared_secret }` |
| `crypto.pqc_decapsulate(ciphertext, dk)` | `kripto.pqc_dekapsulasi(...)` | `ciphertext`, `dk`: decapsulation key | `Bytes` | Decapsulates 32-byte shared secret from ML-KEM-768 ciphertext |
| `crypto.xwing_generate_keypair()` | `kripto.xwing_buat_kunci()` | `-` | `Map` | Generates Hybrid X-Wing keypair (combining X25519 classical + ML-KEM-768 quantum) |
| `crypto.xwing_encapsulate(public_key)` | `kripto.xwing_enkapsulasi(...)` | `public_key`: `String`/`Bytes` | `Map` | Encapsulates hybrid shared secret: returns `{ ciphertext, shared_secret }` |
| `crypto.xwing_decapsulate(ciphertext, priv_key)` | `kripto.xwing_dekapsulasi(...)` | `ciphertext`, `priv_key` | `Bytes` | Decapsulates 32-byte hybrid shared secret from X-Wing ciphertext |

---

## 🛡️ High-Level Cryptographic Recipes

### 1. Authenticated Symmetric Encryption (AES-256-GCM AEAD)

Celang provides authenticated encryption using **AES-256** in **Galois/Counter Mode (GCM)**.

#### Core Guarantees:
- **Confidentiality and Authenticity Combined**: GCM protects the privacy of the plaintext and generates a 16-byte authentication tag that detects any unauthorized modification.
- **Automatic 12-Byte CSPRNG Nonce**: Every call to `crypto.encrypt()` generates a fresh 12-byte random nonce from the operating system's cryptographic entropy source (`crypto/rand`). Encrypting the exact same message multiple times yields distinct ciphertexts each time.
- **Combined Payload**: The returned Base64 string packages the components into a single transportable token:
  $$\text{Payload} = \text{Base64}(\text{12-byte Nonce} \parallel \text{Ciphertext} \parallel \text{16-byte Tag})$$
- **Tamper Detection**: If even a single bit of the ciphertext or tag is modified, or if an incorrect key is provided, `crypto.decrypt()` fails safely and raises an empathetic Celang runtime error without crashing the process.

#### Supported Key Formats:
- Standard 44-character Base64 string (produced by `crypto.generate_key()`).
- 64-character hexadecimal string (32 raw bytes).
- Raw 32-byte `Bytes` object (e.g., from `crypto.derive_key()`).
- Raw 32-character string.

#### Runnable Example (Universal Mode):

```celang
"use universal"
import { generate_key, encrypt, decrypt } from "crypto"

// 1. Generate a secure 256-bit (32-byte) key
let key = generate_key()
print("AES-256 Key:", key)

// 2. Encrypt sensitive data
let secretMessage = "Confidential business contract details"
let ciphertext = encrypt(secretMessage, key)
print("Ciphertext (Base64):", ciphertext)

// 3. Decrypt and verify
let recovered = decrypt(ciphertext, key)
print("Decrypted plaintext:", recovered)

// 4. Nonce uniqueness guarantee
let ciphertext2 = encrypt(secretMessage, key)
if (ciphertext != ciphertext2) {
    print("✅ Nonces are generated automatically per encryption call.")
}

// 5. Tamper detection & safe error handling
let wrongKey = generate_key()

try {
    decrypt(ciphertext, wrongKey)
} catch (err) {
    print("🛡️ Tamper/invalid key intercepted successfully:", err)
    // Error: crypto.decrypt: decryption failed (authentication failed, data modified or invalid key)
}
```

---

### 2. Password Hashing & Key Derivation (PBKDF2-HMAC-SHA256)

Plain hash functions like SHA-256 or MD5 should never be used for user passwords due to high vulnerability to GPU clusters and rainbow table attacks.

Celang implements **PBKDF2-HMAC-SHA256** with conservative security parameters:
- **100,000 Iterations**: Imposes computational cost (*work factor*) that slows down brute-force and dictionary attacks.
- **Automatic 16-Byte Random Salt**: Every call to `crypto.hash_password()` samples 16 fresh random bytes from CSPRNG, making pre-computed rainbow tables ineffective.
- **PHC String Format**: Formatted cleanly as:
  `$pbkdf2-sha256$i=100000$<salt_base64>$<hash_base64>`
- **Constant-Time Verification**: `crypto.verify_password()` performs constant-time comparison using `crypto/subtle`, preventing timing side-channel attacks.

#### Key Derivation (`crypto.derive_key`)

When deriving symmetric encryption keys from human-chosen passphrases, use `crypto.derive_key()`. It converts arbitrary-length passphrases into uniform, high-entropy cryptographic keys.

Supported output formats for the optional 4th argument:
- `"bytes"` / `"raw"` (default): Returns a `Bytes` object.
- `"hex"`: Returns a hexadecimal `String`.
- `"base64"`: Returns a standard Base64 `String` (directly compatible with `crypto.encrypt()`).

#### Runnable Example (Universal Mode):

```celang
"use universal"
import { hash_password, verify_password, derive_key, encrypt, decrypt, random_bytes } from "crypto"

// -------------------------------------------------------------
// Part A: User Password Storage & Verification
// -------------------------------------------------------------

let userPassword = "MySuperSecretPassword#2026"

// 1. Hash during user registration and save to database:
let storedHash = hash_password(userPassword)
print("Stored PHC Hash:", storedHash)
// e.g. $pbkdf2-sha256$i=100000$c2FsdF9leGFtcGxl...$aGFzaF9leGFtcGxl...

// 2. Verify during user login:
let attempt1 = "MySuperSecretPassword#2026"
let success = verify_password(attempt1, storedHash)

if (success) {
    print("✅ Login successful: Password matches!")
} else {
    print("❌ Login rejected: Incorrect password.")
}

// 3. Test wrong password attempt:
let attempt2 = "WrongPasswordGuess"
let failed = verify_password(attempt2, storedHash)
print("Invalid attempt verification result:", failed) // false


// -------------------------------------------------------------
// Part B: Deriving Encryption Keys from a Passphrase
// -------------------------------------------------------------

let masterPassphrase = "master-wallet-seed-phrase"
let salt = random_bytes(16)

// Derive a 32-byte AES-256 key in Base64 format:
let derivedKey = derive_key(masterPassphrase, salt, 32, "base64")

// Encrypt data directly using the derived key:
let secureNotes = "Backup seeds and recovery credentials"
let encryptedNotes = encrypt(secureNotes, derivedKey)
let decryptedNotes = decrypt(encryptedNotes, derivedKey)

print("Decrypted with derived key:", decryptedNotes)
```

---

### 3. Modern Asymmetric Digital Signatures (Ed25519 - RFC 8032)

For scenarios requiring asymmetric authentication without sharing secret keys (e.g., package manager verification, code signing, cryptographic licenses), Celang provides **Ed25519** (*Edwards-curve Digital Signature Algorithm* on Curve25519) according to **RFC 8032**.

#### Advantages of Ed25519:
- **Compact Keys**: 32-byte public keys (64 hex characters) and 32-byte seed private keys (64 hex characters).
- **High Security**: Equivalent to ~3072-bit RSA with superior performance and small signature sizes (64 bytes / 128 hex characters).
- **Side-Channel Resilience**: Constant-time curve arithmetic and memory-access patterns.
- **Deterministic**: Signing identical messages with the same private key always produces the exact same signature, avoiding RNG failures during signature creation.

#### Return Value of `generate_keypair()`:
Returns a `Map` with both English and Indonesian property names:
- `pair.public_key` (alias `pair.kunci_publik`): 64-character hex string (32 bytes).
- `pair.private_key` (alias `pair.kunci_privat`): 64-character hex string (32-byte seed).

#### Format Flexibility:
`crypto.sign()` and `crypto.verify()` accept private keys, public keys, and signatures in **hexadecimal**, **Base64**, raw strings, or `Bytes` objects.

#### Runnable Example (Universal Mode):

```celang
"use universal"
import { generate_keypair, sign, verify } from "crypto"

// 1. Package publisher generates Ed25519 keypair
let pair = generate_keypair()
let publicKey = pair.public_key
let privateKey = pair.private_key

print("Public Key (distribute to clients):", publicKey)
print("Private Key (keep secret on build server):", privateKey)

// 2. Publisher signs a package release manifest
let manifest = "package=celang-net&version=1.4.0&checksum=7f83b165..."
let signature = sign(manifest, privateKey)
print("Ed25519 Signature (128 hex chars):", signature)

// 3. Client / Package Manager verifies the signature before installation
let isValid = verify(manifest, signature, publicKey)

if (isValid) {
    print("✅ Package manifest is authentic and untampered.")
} else {
    print("❌ Verification failed: Untrusted publisher or corrupt payload.")
}

// 4. Tamper detection:
// Modifying a single character of the manifest causes verification to fail:
let tamperedManifest = "package=celang-net&version=1.4.0-malware&checksum=9999..."
let isTamperedValid = verify(tamperedManifest, signature, publicKey)
print("Tampered payload verification:", isTamperedValid) // false
```

---

## 💡 Practical Examples: Hashes, HMAC, & Utilities

### 1. Computing SHA-256 Checksums on Strings & Raw Bytes

```celang
"use universal"
import { sha256 } from "crypto"

// Hash string payload
let textHash = sha256("Celang Universal Cryptography")
print("SHA-256 text:", textHash)

// Hash raw byte array
let byteBuffer = bytes([67, 101, 108, 97, 110, 103]) // "Celang"
let byteHash = sha256(byteBuffer)
print("SHA-256 bytes:", byteHash)
```

### 2. Webhook Authentication with HMAC-SHA256 & Timing-Safe Comparison

```celang
"use universal"
import { hmac, timing_safe_equal } from "crypto"

let webhookSecret = "super_secure_webhook_signing_secret"
let requestBody = "{\"event\":\"order.created\",\"amount\":250.00}"

// Compute HMAC-SHA256 digest:
let computedSignature = hmac(webhookSecret, requestBody)
print("Calculated X-Signature:", computedSignature)

// Compare against incoming HTTP header using constant-time check:
let incomingSignature = computedSignature // received in request header
let isAuthentic = timing_safe_equal(computedSignature, incomingSignature)

if (isAuthentic) {
    print("✅ Webhook payload verified successfully.")
} else {
    print("❌ Invalid webhook signature!")
}
```

### 3. Session Tokens & Random UUID v4

```celang
"use universal"
import { random_hex, random_bytes, uuid } from "crypto"

// Generate 16 random bytes (32 hex characters) for session token:
let sessionToken = random_hex(16)
print("Session Token:", sessionToken)

// Generate 32 raw random bytes for secret seeds:
let rawEntropy = random_bytes(32)
print("Entropy size:", len(rawEntropy))

// Standard RFC 4122 UUID v4:
let accountId = uuid()
print("Account UUID:", accountId)
// e.g. "a98835e0-82a8-48b4-8777-62f79f2ff049"
```

### 4. High-Performance Modern Hashes & Poly1305 MAC

```celang
"use universal"
import "crypto"

let message = "Confidential API payload"

// 1. Keccak SHA-3 & SHAKE XOF
let h3_256 = crypto.sha3_256(message)
let h3_512 = crypto.sha3_512(message)
let shake = crypto.shake128(message, 32) // 32-byte dynamic length

// 2. High-speed BLAKE2 & BLAKE3
let b2 = crypto.blake2b(message)
let b3 = crypto.blake3(message)

// Keyed MAC using BLAKE3 (32 bytes)
let macKey = "exact-32-byte-secret-key-material!"
let mac = crypto.blake3(message, macKey)

// 3. One-time Poly1305 MAC (RFC 8439)
let polyKey = crypto.generate_key()
let tag = crypto.poly1305_auth(message, polyKey)
let isTagValid = crypto.poly1305_verify(message, tag, polyKey)
```

### 5. Memory-Hard Password Hashing (Argon2id) & HKDF (RFC 5869)

```celang
"use universal"
import "crypto"

// 1. Argon2id Password Storage (PHC String Format)
let password = "SuperSecretPassword#2026!"
let argonHash = crypto.hash_password_argon2(password)
// Output: "$argon2id$v=19$m=65536,t=3,p=4$..."

// Custom options can be passed via Map:
let customHash = crypto.hash_password_argon2(password, {"time": 2, "memory": 65536, "threads": 2})

let isPasswordMatch = crypto.verify_password_argon2(password, argonHash)
print("Password valid:", isPasswordMatch) // true

// 2. HKDF Key Derivation (RFC 5869 Extract-and-Expand)
// Per RFC 5869 with SHA-256, max derived key length is 255 * 32 = 8,160 bytes.
let masterSecret = "high-entropy-master-secret"
let salt = "app-salt-session-1234"
let derivedKey = crypto.hkdf(masterSecret, salt, "celang-service-v1", 32)
```

> [!NOTE]
> **Argon2 Verification DoS Protection**:
> When verifying password hashes from external inputs, `crypto.verify_password_argon2` enforces strict upper bounds: maximum memory $\le 256$ MB (`256 * 1024` KB) and iterations $t \le 10$ to protect servers against Denial of Service (DoS / OOM) attacks.

### 6. Symmetric XChaCha20-Poly1305 AEAD & AES Block Modes

```celang
"use universal"
import "crypto"

let plaintext = "Sensitive data transmission over untrusted network"

// 1. XChaCha20-Poly1305 AEAD with Extended 24-byte Nonce
// Uses a 192-bit (24-byte) random nonce to eliminate nonce collision risk (birthday bound).
// Decryption is adaptive: transparently supports both 24-byte nonces and legacy 12-byte nonces.
let chachaKey = crypto.generate_chacha_key()
let aad = "user-id:1001;org:celang"
let ciphertext = crypto.encrypt_chacha(plaintext, chachaKey, aad)
let recovered = crypto.decrypt_chacha(ciphertext, chachaKey, aad)

// 2. Standard AES Modes (CTR & CBC with PKCS#7)
// AES-CBC incorporates constant-time PKCS#7 unpadding to mitigate padding oracle attacks.
let aesKey = "01234567890123456789012345678901" // 32 bytes (AES-256)
let iv = "0123456789012345"                     // 16 bytes
let encCTR = crypto.aes_ctr(plaintext, aesKey, iv)
let decCTR = crypto.aes_ctr(encCTR, aesKey, iv)

let encCBC = crypto.aes_cbc_encrypt(plaintext, aesKey, iv)
let decCBC = crypto.aes_cbc_decrypt(encCBC, aesKey, iv)
```

### 7. Asymmetric Key Exchange & Signatures (X25519 & RSA)

```celang
"use universal"
import "crypto"

// 1. Curve25519 ECDH Key Exchange (Alice & Bob)
let alice = crypto.generate_x25519_keypair()
let bob = crypto.generate_x25519_keypair()

let aliceSecret = crypto.x25519_shared_secret(alice.private_key, bob.public_key)
let bobSecret = crypto.x25519_shared_secret(bob.private_key, alice.public_key)
let ecdhMatches = crypto.timing_safe_equal(aliceSecret, bobSecret) // true

// 2. RSA 2048/4096-bit Encryption (OAEP) & Signatures (PSS)
let rsa = crypto.rsa_generate_keypair(2048)
let rsaCipher = crypto.rsa_encrypt_oaep("Classified Note", rsa.public_key)
let rsaPlain = crypto.rsa_decrypt_oaep(rsaCipher, rsa.private_key)

// RSA-PSS signing uses deterministic SaltLengthEqualsHash.
// RSA-PSS verification uses SaltLengthAuto (accepting both hex and Base64 signatures from OpenSSL, Java, Python).
let rsaSig = crypto.rsa_sign_pss("Legal Agreement", rsa.private_key)
let isSigValid = crypto.rsa_verify_pss("Legal Agreement", rsaSig, rsa.public_key)
```

### 8. Post-Quantum Cryptography & Hybrid KEM (ML-KEM-768 & Hybrid X-Wing)

Guarantees long-term forward secrecy against future quantum attacks (*Harvest Now, Decrypt Later*):

```celang
"use universal"
import "crypto"

// 1. NIST FIPS 203 ML-KEM-768 (Kyber)
let pqc = crypto.pqc_generate_keypair()

// Sender encapsulates shared secret with recipient's encapsulation key
let encResult = crypto.pqc_encapsulate(pqc.encapsulation_key)
let senderSecret = encResult.shared_secret
let pqcCiphertext = encResult.ciphertext

// Recipient decapsulates using private decapsulation key
let recipientSecret = crypto.pqc_decapsulate(pqcCiphertext, pqc.decapsulation_key)
let pqcMatches = crypto.timing_safe_equal(senderSecret, recipientSecret) // true

// 2. Hybrid KEM: X-Wing (X25519 Classical + ML-KEM-768 Quantum)
// Recommended for production: resists both classical cryptanalysis and quantum attacks
let xwing = crypto.xwing_generate_keypair()
let xwEnc = crypto.xwing_encapsulate(xwing.public_key)
let xwSecret = crypto.xwing_decapsulate(xwEnc.ciphertext, xwing.private_key)
let xwMatches = crypto.timing_safe_equal(xwEnc.shared_secret, xwSecret) // true
```

---

## 🔒 Security Best Practices

1. **Always Use `timing_safe_equal` for Secret Verification**:
   Standard string equality (`==`) terminates evaluation early upon encountering the first mismatched character. This creates microsecond variations exploitable via timing attacks. Always use `crypto.timing_safe_equal()` or `kripto.konstanta_sama()`.

2. **Use `encrypt` / `decrypt` (AES-256-GCM) for Stored Data**:
   Never implement custom encryption using unauthenticated cipher modes (such as CBC or ECB). Celang's `crypto.encrypt` combines AES-256-GCM with automatic CSPRNG nonces and authentication tags to guarantee both confidentiality and integrity.

3. **Always Hash User Passwords with `hash_password`**:
   Never store user passwords using simple cryptographic digests like SHA-256 or MD5. Always rely on `crypto.hash_password()`, which enforces 100,000 PBKDF2-HMAC-SHA256 iterations and random 16-byte salt generation.

4. **Safeguard Private Keys**:
   Never commit Ed25519 or RSA private keys to source control. Store them in secure key vaults or environment variables, and distribute only the public key to clients.

5. **Prefer Authenticated Encryption (AEAD) over Raw Cipher Modes (CBC / CTR)**:
   `crypto.aes_cbc_encrypt` and `crypto.aes_ctr` do not include built-in authenticity tags. Using unauthenticated CBC over a network can expose your service to padding oracle vulnerabilities. Always prefer `crypto.encrypt` (AES-256-GCM) or `crypto.encrypt_chacha` (ChaCha20-Poly1305). If using CBC or CTR, always combine it with HMAC-SHA256 (*Encrypt-then-MAC*).

6. **Poly1305 Requires a One-Time Key**:
   Never reuse the same 32-byte key to authenticate multiple messages with Poly1305. Key reuse catastrophically compromises the authenticator. If you need a reusable key for message authentication, use `crypto.hmac` (HMAC-SHA256) or keyed `crypto.blake3(data, key)`.

7. **Derive Symmetric Keys from `x25519_shared_secret` using HKDF**:
   Do not use the raw 32-byte shared secret from X25519 directly as a symmetric cipher key. Pass the shared secret through HKDF (`crypto.hkdf(sharedSecret, salt, "protocol-purpose", 32)`) to obtain uniform key material.
