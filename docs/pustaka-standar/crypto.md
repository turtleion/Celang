---
outline: deep
---

# Cryptography Library (`crypto` / `kripto`)

The **`crypto`** module (official Indonesian name **`kripto`**) provides a modern, production-grade cryptography ecosystem in Celang. Built with pure Go standard library primitives (zero external dependencies and no CGo bindings), it is designed around the **safe-by-default** philosophy and inspired by modern high-level cryptographic libraries like Python's `cryptography.io`.

The library provides comprehensive cryptographic functionality ranging from authenticated encryption and GPU-resistant password hashing to asymmetric Ed25519 signatures, standard digests, and constant-time utilities:

1. **Authenticated Symmetric Encryption (AES-256-GCM AEAD)**: Safe, foolproof symmetric encryption with automatic 12-byte random nonce generation, 16-byte authentication tag, and built-in tamper detection (`generate_key`, `encrypt`, `decrypt`).
2. **Password Hashing & Key Derivation (PBKDF2-HMAC-SHA256)**: GPU brute-force resistant password storage adhering to the PHC string format, timing-safe verification, and cryptographic key derivation (`hash_password`, `verify_password`, `derive_key`).
3. **Modern Asymmetric Digital Signatures (Ed25519 - RFC 8032)**: High-performance digital signing and verification with compact 32-byte keys for package manager verification, software updates, and API contracts (`generate_keypair`, `sign`, `verify`).
4. **Cryptographic Hashes & Checksums**: Fast and secure digests (SHA-256, SHA-512, SHA-1, MD5).
5. **Keyed-Hash Message Authentication Codes (HMAC)**: HMAC-SHA256 for webhook signing and API request verification.
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

---

## 🔒 Security Best Practices

1. **Always Use `timing_safe_equal` for Secret Verification**:
   Standard string equality (`==`) terminates evaluation early upon encountering the first mismatched character. This creates microsecond variations exploitable via timing attacks. Always use `crypto.timing_safe_equal()` or `kripto.konstanta_sama()`.

2. **Use `encrypt` / `decrypt` (AES-256-GCM) for Stored Data**:
   Never implement custom encryption using unauthenticated cipher modes (such as CBC or ECB). Celang's `crypto.encrypt` combines AES-256-GCM with automatic CSPRNG nonces and authentication tags to guarantee both confidentiality and integrity.

3. **Always Hash User Passwords with `hash_password`**:
   Never store user passwords using simple cryptographic digests like SHA-256 or MD5. Always rely on `crypto.hash_password()`, which enforces 100,000 PBKDF2-HMAC-SHA256 iterations and random 16-byte salt generation.

4. **Safeguard Private Keys**:
   Never commit Ed25519 private keys to source control. Store them in secure key vaults or environment variables, and distribute only the public key to clients.
