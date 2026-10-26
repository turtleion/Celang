# Project: Celang High-Level Cryptography Recipes & Bilingual Stdlib

## Architecture
- **Language**: Celang (Go implementation, pure standard library).
- **Core Package**: `pkg/evaluator` (`stdlib_crypto.go`, `stdlib.go`, `stdlib_crypto_test.go`).
- **Standard Library Modules**:
  - `kripto`: Indonesian module name with empathetic diagnostics.
  - `crypto`: English / Universal module name.
- **Data Flow**: Celang AST -> Evaluator Builtins -> Go stdlib crypto (`crypto/aes`, `crypto/cipher`, `crypto/ed25519`, `crypto/hmac`, `crypto/sha256`, `crypto/rand`, `crypto/subtle`) -> Celang Object representations (`*object.String`, `*object.Boolean`, `*object.Map`, `*object.Error`).

## Code Layout
- `pkg/evaluator/stdlib_crypto.go`: Core implementation of symmetric encryption, KDF, password hashing, Ed25519, and bilingual bindings.
- `pkg/evaluator/stdlib_crypto_test.go`: Go unit tests for crypto builtins.
- `tests/kripto_uji.ce`: Native Celang test suite covering end-user script behavior.
- `docs/pustaka-standar/kripto.md`: Official Indonesian documentation.
- `docs/pustaka-standar/crypto.md`: Official English documentation.
- CRITICAL: `tes.ce` is strictly untouched and must never be staged or committed.

## Feature Inventory
| # | Feature | Description | Milestone | Source |
|---|---------|-------------|-----------|--------|
| 1 | R1: Symmetric Encryption | AES-256-GCM AEAD (`buat_kunci`/`generate_key`, `enkripsi`/`encrypt`, `dekripsi`/`decrypt`) with automatic 12-byte nonce and tamper detection | M1 | ORIGINAL_REQUEST §R1 |
| 2 | R2: Password Hashing & KDF | PBKDF2-HMAC-SHA256 (`hash_sandi`/`hash_password`, `verifikasi_sandi`/`verify_password`, `turunkan_kunci`/`derive_key`) with PHC string format and constant-time verification | M1 | ORIGINAL_REQUEST §R2 |
| 3 | R3: Ed25519 Signatures | Modern asymmetric digital signatures (`buat_pasangan_kunci`/`generate_keypair`, `tanda_tangan`/`sign`, `verifikasi_tanda_tangan`/`verify`) with RFC 8032 test vectors | M1 | ORIGINAL_REQUEST §R3 |
| 4 | R4: Bilingual & Universal Mode | Dual module exposure (`kripto` and `crypto`), "use universal" enforcement | M1 | ORIGINAL_REQUEST §R4 |
| 5 | R5: Comprehensive Test Suite | Go unit tests, Celang native test scripts, round-trip, tamper detection, RFC 8032, 100% pass | M2 | ORIGINAL_REQUEST §Acceptance |
| 6 | R6: Documentation | Up-to-date documentation in `kripto.md` and `crypto.md` passing VitePress build | M3 | ORIGINAL_REQUEST §Acceptance |

## Milestones
| # | Name | Scope | Dependencies | Status |
|---|------|-------|-------------|--------|
| M1 | Crypto Stdlib Implementation | Implement R1, R2, R3, R4 in `pkg/evaluator/stdlib_crypto.go` | none | DONE |
| M2 | E2E & Cryptographic Test Suite | Comprehensive unit, integration, and Celang native test suite; publish `TEST_READY.md` | M1 | DONE |
| M3 | Documentation Update | Update `docs/pustaka-standar/kripto.md` and `docs/pustaka-standar/crypto.md` | M1 | DONE |
| M4 | Final Milestone & Hardening | Phase 1: 100% test pass; Phase 2: Tier 5 adversarial testing & Forensic Audit | M1, M2, M3 | DONE |

## Interface Contracts
### R1: Symmetric Encryption (AES-256-GCM)
- `kripto.buat_kunci()` / `crypto.generate_key()` -> `String` (32 random bytes encoded in standard Base64).
- `kripto.enkripsi(data: String|Bita, kunci: String)` / `crypto.encrypt(data, key)` -> `String` (Base64 string of: 12-byte nonce + ciphertext + 16-byte authentication tag).
- `kripto.dekripsi(ciphertext: String, kunci: String)` / `crypto.decrypt(ciphertext, key)` -> `String` (Decrypted plaintext string) or `*object.Error` on tampering/invalid key.

### R2: Password Hashing & KDF (PBKDF2-HMAC-SHA256)
- `kripto.hash_sandi(sandi: String)` / `crypto.hash_password(password)` -> `String` (PHC format: `$pbkdf2-sha256$i=100000$<salt_b64>$<hash_b64>`).
- `kripto.verifikasi_sandi(sandi: String, hash: String)` / `crypto.verify_password(password, hash)` -> `Boolean` (`TRUE` or `FALSE`) using `subtle.ConstantTimeCompare`.
- `kripto.turunkan_kunci(sandi: String, salt: String|Bita, panjang: Integer)` / `crypto.derive_key(password, salt, len)` -> `String` (hex or base64 derived key) or `Bita`.

### R3: Ed25519 Signatures
- `kripto.buat_pasangan_kunci()` / `crypto.generate_keypair()` -> `Map` containing:
  - `"kunci_publik"` / `"public_key"`: 32 bytes hex string (64 characters).
  - `"kunci_privat"` / `"private_key"`: 64 bytes hex string (128 characters, seed+public).
- `kripto.tanda_tangan(data: String|Bita, kunci_privat: String)` / `crypto.sign(data, private_key)` -> `String` (64 bytes signature in hex string, 128 characters).
- `kripto.verifikasi_tanda_tangan(data: String|Bita, signature: String, kunci_publik: String)` / `crypto.verify(data, signature, public_key)` -> `Boolean` (`TRUE` or `FALSE`).
