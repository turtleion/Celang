# E2E Test Infra: Celang High-Level Cryptography

## Test Philosophy
- Opaque-box, requirement-driven, mathematically rigorous.
- Methodology: Category-Partition + BVA + Pairwise + Real-World Workload Testing.
- RFC 8032 test vector verification for Ed25519.
- Timing-attack safety verification via constant-time comparison.
- Empathetic error validation on tampered payloads and malformed keys.

## Feature Inventory
| # | Feature | Source | Tier 1 | Tier 2 | Tier 3 | Tier 4 |
|---|---------|--------|:------:|:------:|:------:|:------:|
| 1 | R1: Symmetric Encryption (AES-256-GCM) | ORIGINAL_REQUEST §R1 | 5 | 5 | ✓ | ✓ |
| 2 | R2: Password Hashing & KDF | ORIGINAL_REQUEST §R2 | 5 | 5 | ✓ | ✓ |
| 3 | R3: Ed25519 Signatures | ORIGINAL_REQUEST §R3 | 5 | 5 | ✓ | ✓ |
| 4 | R4: Bilingual & Universal Mode | ORIGINAL_REQUEST §R4 | 5 | 5 | ✓ | ✓ |

## Test Architecture
- **Go Unit & Builtin Test Runner**: `go test -v -count=1 ./pkg/evaluator -run TestStdlibCrypto`
- **Full Celang Suite**: `go test ./...`
- **Celang Script Runner**: `./celang uji tests/` or `./celang run <script>`
- **Pass/Fail Semantics**: 0 exit code, 100% test pass, zero regressions, `tes.ce` untouched.

## Test Tiers
- **Tier 1: Feature Coverage (≥5 per feature)**: Basic happy-path encryption/decryption, password hashing, verification, keypair generation, signature, and verification.
- **Tier 2: Boundary & Corner Cases (≥5 per feature)**: Empty data, large data, 1-bit modified ciphertext (tamper detection), wrong key, malformed hash string, corrupted Ed25519 signatures, invalid public keys.
- **Tier 3: Cross-Feature Interactions**: Derive key with KDF -> Encrypt with AES-GCM -> Decrypt with derived key. Sign derived ciphertext with Ed25519 -> Verify signature -> Decrypt.
- **Tier 4: Real-World Scenarios**: Package manager manifest signing & verification, secure user authentication vault, encrypted credential storage.
- **Tier 5: Adversarial Hardening**: White-box edge cases, coverage audit, mutation tests.
