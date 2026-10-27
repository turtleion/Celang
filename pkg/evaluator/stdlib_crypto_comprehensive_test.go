package evaluator_test

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
	"testing"

	"celang/pkg/evaluator"
	"golang.org/x/crypto/chacha20poly1305"
)

func TestCrypto_Pillar1_HashesAndMAC(t *testing.T) {
	script := `
	buat pesan = "Kriptografi Modern Celang 2026"

	// 1. SHA-3
	buat h_sha3_256 = kripto.sha3_256(pesan)
	buat h_sha3_512 = kripto.sha3_512(pesan)
	buat h_shake128 = kripto.shake128(pesan, 32)
	buat h_shake256 = kripto.shake256(pesan, 64)

	// 2. BLAKE2b & BLAKE2s
	buat h_blake2b = kripto.blake2b(pesan)
	buat h_blake2s = kripto.blake2s(pesan)

	buat kunciB2 = "01234567890123456789012345678901"
	buat h_blake2b_kunci = kripto.blake2b(pesan, kunciB2)

	// 3. BLAKE3
	buat h_blake3 = kripto.blake3(pesan)
	buat h_blake3_kunci = kripto.blake3(pesan, kunciB2)

	// 4. Poly1305
	buat kunciPoly = kripto.buat_kunci()
	buat tagPoly = kripto.poly1305_buat_mac(pesan, kunciPoly)
	buat cekPolyOk = kripto.poly1305_verifikasi(pesan, tagPoly, kunciPoly)
	buat cekPolyPalsu = kripto.poly1305_verifikasi("pesan palsu", tagPoly, kunciPoly)

	[h_sha3_256, h_sha3_512, h_shake128, h_shake256, h_blake2b, h_blake2s, h_blake3, cekPolyOk, cekPolyPalsu]
	`
	res := runCryptoScript(script)
	if res == nil || res.Type() == evaluator.ERROR_OBJ {
		t.Fatalf("Script error: %v", res)
	}

	arr, ok := res.(*evaluator.Array)
	if !ok {
		t.Fatalf("Expected Array, got %v", res)
	}

	// Verify lengths of hashes
	if len(arr.Elements[0].(*evaluator.String).Value) != 64 { // SHA3-256 = 32 bytes = 64 hex
		t.Errorf("SHA3-256 length mismatch: %s", arr.Elements[0].Inspect())
	}
	if len(arr.Elements[1].(*evaluator.String).Value) != 128 { // SHA3-512 = 64 bytes = 128 hex
		t.Errorf("SHA3-512 length mismatch: %s", arr.Elements[1].Inspect())
	}
	if len(arr.Elements[2].(*evaluator.String).Value) != 64 { // SHAKE128 (32 bytes) = 64 hex
		t.Errorf("SHAKE128 length mismatch: %s", arr.Elements[2].Inspect())
	}
	if len(arr.Elements[4].(*evaluator.String).Value) != 128 { // BLAKE2b (64 bytes) = 128 hex
		t.Errorf("BLAKE2b length mismatch: %s", arr.Elements[4].Inspect())
	}
	if len(arr.Elements[5].(*evaluator.String).Value) != 64 { // BLAKE2s (32 bytes) = 64 hex
		t.Errorf("BLAKE2s length mismatch: %s", arr.Elements[5].Inspect())
	}
	if len(arr.Elements[6].(*evaluator.String).Value) != 64 { // BLAKE3 (32 bytes) = 64 hex
		t.Errorf("BLAKE3 length mismatch: %s", arr.Elements[6].Inspect())
	}

	// Poly1305 verification
	if arr.Elements[7] != evaluator.TRUE {
		t.Errorf("Poly1305 verification expected true, got %v", arr.Elements[7])
	}
	if arr.Elements[8] != evaluator.FALSE {
		t.Errorf("Poly1305 verification expected false for tampered data, got %v", arr.Elements[8])
	}
}

func TestCrypto_Pillar2_KDFAndArgon2(t *testing.T) {
	script := `
	// 1. HKDF
	buat secret = "master-secret-key-material-2026"
	buat salt = "random-salt-1234"
	buat derivedKey = kripto.hkdf(secret, salt, "celang-session", 32)
	buat prk = kripto.hkdf_ekstrak(secret, salt)
	buat expanded = kripto.hkdf_ekspansi(prk, "celang-session", 32)

	// 2. Argon2id Password Hashing
	buat sandi = "SandiRahasiaSuperAman123!"
	buat hashArgon = kripto.hash_sandi_argon2(sandi)
	buat cekArgonBenar = kripto.verifikasi_sandi_argon2(sandi, hashArgon)
	buat cekArgonSalah = kripto.verifikasi_sandi_argon2("SandiSalah", hashArgon)

	[derivedKey.panjang, expanded.panjang, cekArgonBenar, cekArgonSalah, hashArgon]
	`
	res := runCryptoScript(script)
	if res == nil || res.Type() == evaluator.ERROR_OBJ {
		t.Fatalf("Script error: %v", res)
	}

	arr := res.(*evaluator.Array)
	if arr.Elements[0].(*evaluator.Integer).Value != 32 {
		t.Errorf("HKDF key length expected 32, got %v", arr.Elements[0])
	}
	if arr.Elements[1].(*evaluator.Integer).Value != 32 {
		t.Errorf("HKDF expanded length expected 32, got %v", arr.Elements[1])
	}
	if arr.Elements[2] != evaluator.TRUE {
		t.Errorf("Argon2id verify expected true for correct password")
	}
	if arr.Elements[3] != evaluator.FALSE {
		t.Errorf("Argon2id verify expected false for wrong password")
	}
	phc := arr.Elements[4].(*evaluator.String).Value
	if !strings.HasPrefix(phc, "$argon2id$v=19$") {
		t.Errorf("Argon2id string prefix invalid: %s", phc)
	}
}

func TestCrypto_Pillar3_SymmetricChaChaAndAES(t *testing.T) {
	script := `
	// 1. ChaCha20-Poly1305 AEAD
	buat kunciChaCha = kripto.buat_kunci_chacha()
	buat plaintext = "Pesan rahasia dengan enkripsi ChaCha20-Poly1305"
	buat aad = "header-metadata-id-42"
	buat cipherChaCha = kripto.enkripsi_chacha(plaintext, kunciChaCha, aad)
	buat pulihChaCha = kripto.dekripsi_chacha(cipherChaCha, kunciChaCha, aad)

	// 2. AES-CTR
	buat kunciAES = "01234567890123456789012345678901"
	buat iv = "0123456789012345"
	buat encCTR = kripto.aes_ctr(plaintext, kunciAES, iv)
	buat decCTR = kripto.aes_ctr(encCTR, kunciAES, iv)

	// 3. AES-CBC with PKCS#7
	buat encCBC = kripto.aes_cbc_enkripsi(plaintext, kunciAES, iv)
	buat decCBC = kripto.aes_cbc_dekripsi(encCBC, kunciAES, iv)

	[pulihChaCha, decCTR.ke_teks(), decCBC.ke_teks()]
	`
	res := runCryptoScript(script)
	if res == nil || res.Type() == evaluator.ERROR_OBJ {
		t.Fatalf("Script error: %v", res)
	}

	arr := res.(*evaluator.Array)
	expected := "Pesan rahasia dengan enkripsi ChaCha20-Poly1305"
	if arr.Elements[0].(*evaluator.String).Value != expected {
		t.Errorf("ChaCha20 mismatch: got %s", arr.Elements[0].(*evaluator.String).Value)
	}
	if arr.Elements[1].(*evaluator.String).Value != expected {
		t.Errorf("AES-CTR mismatch: got %s", arr.Elements[1].(*evaluator.String).Value)
	}
	if arr.Elements[2].(*evaluator.String).Value != expected {
		t.Errorf("AES-CBC mismatch: got %s", arr.Elements[2].(*evaluator.String).Value)
	}
}

func TestCrypto_Pillar4_X25519AndRSA(t *testing.T) {
	script := `
	// 1. X25519 Diffie-Hellman Key Exchange (Alice & Bob)
	buat alice = kripto.buat_kunci_x25519()
	buat bob = kripto.buat_kunci_x25519()

	// Alice menghitung shared secret dengan kunci publik Bob
	buat rahasiaAlice = kripto.x25519_hitung_rahasia(alice.kunci_privat, bob.kunci_publik)
	// Bob menghitung shared secret dengan kunci publik Alice
	buat rahasiaBob = kripto.x25519_hitung_rahasia(bob.kunci_privat, alice.kunci_publik)

	buat ecdhSama = kripto.konstanta_sama(rahasiaAlice, rahasiaBob)

	// 2. RSA (2048-bit) Key Generation, OAEP Encrypt, PSS Sign
	buat rsaKeys = kripto.rsa_buat_kunci(2048)
	buat teksAsli = "Dokumen penting bertanda tangan digital RSA-PSS"

	// Enkripsi & Dekripsi OAEP
	buat rsaCipher = kripto.rsa_enkripsi_oaep(teksAsli, rsaKeys.kunci_publik)
	buat rsaPulih = kripto.rsa_dekripsi_oaep(rsaCipher, rsaKeys.kunci_privat)

	// Tanda Tangan & Verifikasi PSS
	buat ttdRSA = kripto.rsa_tanda_tangan_pss(teksAsli, rsaKeys.kunci_privat)
	buat validPSS = kripto.rsa_verifikasi_pss(teksAsli, ttdRSA, rsaKeys.kunci_publik)
	buat palsuPSS = kripto.rsa_verifikasi_pss("Dokumen dipalsukan!", ttdRSA, rsaKeys.kunci_publik)

	[ecdhSama, rsaPulih, validPSS, palsuPSS]
	`
	res := runCryptoScript(script)
	if res == nil || res.Type() == evaluator.ERROR_OBJ {
		t.Fatalf("Script error: %v", res)
	}

	arr := res.(*evaluator.Array)
	if arr.Elements[0] != evaluator.TRUE {
		t.Errorf("X25519 Alice and Bob shared secret mismatch!")
	}
	if arr.Elements[1].(*evaluator.String).Value != "Dokumen penting bertanda tangan digital RSA-PSS" {
		t.Errorf("RSA OAEP decrypted plaintext mismatch")
	}
	if arr.Elements[2] != evaluator.TRUE {
		t.Errorf("RSA PSS signature verification expected true")
	}
	if arr.Elements[3] != evaluator.FALSE {
		t.Errorf("RSA PSS signature verification expected false for tampered doc")
	}
}

func TestCrypto_Pillar5_PostQuantumAndHybrid(t *testing.T) {
	script := `
	// 1. Post-Quantum ML-KEM-768 (NIST FIPS 203 / Kyber)
	buat pqcKunci = kripto.pqc_buat_kunci()
	buat resEnkapsulasi = kripto.pqc_enkapsulasi(pqcKunci.kunci_enkapsulasi)
	buat sharedSecretSender = resEnkapsulasi.kunci_bersama
	buat ciphertext = resEnkapsulasi.ciphertext

	// Penerima mendekapsulasi menggunakan kunci dekapsulasi
	buat sharedSecretReceiver = kripto.pqc_dekapsulasi(ciphertext, pqcKunci.kunci_dekapsulasi)
	buat pqcCocok = kripto.konstanta_sama(sharedSecretSender, sharedSecretReceiver)

	// 2. Hybrid KEM (X-Wing: X25519 + ML-KEM-768)
	buat xwingKunci = kripto.xwing_buat_kunci()
	buat resXwing = kripto.xwing_enkapsulasi(xwingKunci.kunci_publik)
	buat xwingSharedSender = resXwing.kunci_bersama
	buat xwingCipher = resXwing.ciphertext

	buat xwingSharedReceiver = kripto.xwing_dekapsulasi(xwingCipher, xwingKunci.kunci_privat)
	buat xwingCocok = kripto.konstanta_sama(xwingSharedSender, xwingSharedReceiver)

	[pqcCocok, xwingCocok, sharedSecretSender.panjang, xwingSharedSender.panjang]
	`
	res := runCryptoScript(script)
	if res == nil || res.Type() == evaluator.ERROR_OBJ {
		t.Fatalf("Script error: %v", res)
	}

	arr := res.(*evaluator.Array)
	if arr.Elements[0] != evaluator.TRUE {
		t.Errorf("ML-KEM-768 sender and receiver shared secret mismatch!")
	}
	if arr.Elements[1] != evaluator.TRUE {
		t.Errorf("Hybrid X-Wing sender and receiver shared secret mismatch!")
	}
	if arr.Elements[2].(*evaluator.Integer).Value != 32 {
		t.Errorf("ML-KEM shared secret length expected 32, got %v", arr.Elements[2])
	}
	if arr.Elements[3].(*evaluator.Integer).Value != 32 {
		t.Errorf("Hybrid X-Wing shared secret length expected 32, got %v", arr.Elements[3])
	}
}

func TestCrypto_Pillars_UniversalMode(t *testing.T) {
	script := `"use universal"
	import "crypto"

	// Pillar 1: Hash & MAC
	let h3 = crypto.sha3_256("test")
	let b3 = crypto.blake3("test")
	let polyKey = crypto.generate_key()
	let tag = crypto.poly1305_mac("test", polyKey)
	let tagOk = crypto.poly1305_verify("test", tag, polyKey)

	// Pillar 2: HKDF & Argon2
	let salt = "salt-universal-1234"
	let k = crypto.hkdf("secret", salt, "info", 32)
	let pwdHash = crypto.hash_password_argon2("secretpass")
	let pwdOk = crypto.verify_password_argon2("secretpass", pwdHash)

	// Pillar 3: ChaCha20 & AES
	let chachaKey = crypto.generate_chacha_key()
	let ctChaCha = crypto.encrypt_chacha("plaintext", chachaKey, "aad")
	let ptChaCha = crypto.decrypt_chacha(ctChaCha, chachaKey, "aad")

	// Pillar 4: X25519 & RSA
	let alice = crypto.generate_x25519_keypair()
	let bob = crypto.generate_x25519_keypair()
	let ssAlice = crypto.x25519_compute_secret(alice.private_key, bob.public_key)
	let ssBob = crypto.x25519_compute_secret(bob.private_key, alice.public_key)
	let ecdhMatch = crypto.constant_time_compare(ssAlice, ssBob)

	let rsa = crypto.rsa_generate_key(2048)
	let encRSA = crypto.rsa_encrypt_oaep("rsa secret", rsa.public_key)
	let decRSA = crypto.rsa_decrypt_oaep(encRSA, rsa.private_key)

	// Pillar 5: PQC & Hybrid
	let pqc = crypto.pqc_generate_keypair()
	let encPQC = crypto.pqc_encapsulate(pqc.encapsulation_key)
	let decPQCSecret = crypto.pqc_decapsulate(encPQC.ciphertext, pqc.decapsulation_key)
	let pqcMatch = crypto.constant_time_compare(encPQC.shared_secret, decPQCSecret)

	let xwing = crypto.xwing_generate_keypair()
	let encXW = crypto.xwing_encapsulate(xwing.public_key)
	let decXWSecret = crypto.xwing_decapsulate(encXW.ciphertext, xwing.private_key)
	let xwMatch = crypto.constant_time_compare(encXW.shared_secret, decXWSecret)

	[tagOk, pwdOk, ptChaCha, ecdhMatch, decRSA, pqcMatch, xwMatch]
	`
	res := runCryptoScript(script)
	if res == nil || res.Type() == evaluator.ERROR_OBJ {
		t.Fatalf("Universal Mode script error: %v", res)
	}

	arr := res.(*evaluator.Array)
	if arr.Elements[0] != evaluator.TRUE {
		t.Errorf("Universal mode tagOk expected true, got %v", arr.Elements[0])
	}
	if arr.Elements[1] != evaluator.TRUE {
		t.Errorf("Universal mode pwdOk expected true, got %v", arr.Elements[1])
	}
	if arr.Elements[2].(*evaluator.String).Value != "plaintext" {
		t.Errorf("Universal mode ptChaCha expected plaintext, got %v", arr.Elements[2])
	}
	if arr.Elements[3] != evaluator.TRUE {
		t.Errorf("Universal mode ecdhMatch expected true, got %v", arr.Elements[3])
	}
	if arr.Elements[4].(*evaluator.String).Value != "rsa secret" {
		t.Errorf("Universal mode decRSA expected 'rsa secret', got %v", arr.Elements[4])
	}
	if arr.Elements[5] != evaluator.TRUE {
		t.Errorf("Universal mode pqcMatch expected true, got %v", arr.Elements[5])
	}
	if arr.Elements[6] != evaluator.TRUE {
		t.Errorf("Universal mode xwMatch expected true, got %v", arr.Elements[6])
	}
}

func TestCrypto_SecurityHardening_Issues(t *testing.T) {
	// 1. Test RSA minimum key size and non-integer argument rejection
	scriptRSA1024 := `kripto.rsa_buat_kunci(1024)`
	resRSA1024 := runCryptoScript(scriptRSA1024)
	if resRSA1024 == nil || resRSA1024.Type() != evaluator.ERROR_OBJ {
		t.Errorf("Expected rsa_buat_kunci(1024) to be rejected as insecure, got %v", resRSA1024)
	}

	scriptRSABadArg := `kripto.rsa_buat_kunci("2048")`
	resRSABadArg := runCryptoScript(scriptRSABadArg)
	if resRSABadArg == nil || resRSABadArg.Type() != evaluator.ERROR_OBJ {
		t.Errorf("Expected rsa_buat_kunci with non-integer arg to return error, got %v", resRSABadArg)
	}

	// 2. Test Argon2 DoS with malicious huge parameters
	scriptArgonDoS := `
	buat hashJahat = "$argon2id$v=19$m=4000000000,t=1000,p=64$c2FsdHNhbHRzYWx0$aGFzaGhhc2hoYXNo"
	kripto.verifikasi_sandi_argon2("password", hashJahat)
	`
	resArgonDoS := runCryptoScript(scriptArgonDoS)
	if resArgonDoS != evaluator.FALSE {
		t.Errorf("Expected verifikasi_sandi_argon2 with malicious m to return FALSE, got %v", resArgonDoS)
	}

	// 3. Test HKDF info argument type check
	scriptHKDFBadInfo := `kripto.hkdf("secret", "salt", 12345, 32)`
	resHKDF := runCryptoScript(scriptHKDFBadInfo)
	if resHKDF == nil || resHKDF.Type() != evaluator.ERROR_OBJ {
		t.Errorf("Expected hkdf with numeric info to return error, got %v", resHKDF)
	}

	// 4. Test AES-CBC backing array non-mutation (aliasing prevention)
	raw := make([]byte, 5, 32)
	copy(raw, []byte{10, 20, 30, 40, 50})
	bObj := &evaluator.Bytes{Elements: raw}
	keyObj := &evaluator.String{Value: "01234567890123456789012345678901"}
	ivObj := &evaluator.String{Value: "0123456789012345"}
	env := evaluator.NewEnvironment()
	evaluator.RegisterStdlibGlobals(env)
	kriptoObj, _ := env.Get("kripto")
	kriptoMap := kriptoObj.(*evaluator.Map)
	cbcEnc := kriptoMap.Pairs["aes_cbc_enkripsi"].(*evaluator.Builtin)
	resEnc := cbcEnc.Fn(bObj, keyObj, ivObj)
	if resEnc.Type() == evaluator.ERROR_OBJ {
		t.Fatalf("CBC encryption failed: %v", resEnc)
	}
	if len(bObj.Elements) != 5 || bObj.Elements[0] != 10 || bObj.Elements[4] != 50 {
		t.Errorf("AES-CBC mutated caller's original Bytes elements: %v", bObj.Elements)
	}
	if raw[:32][5] != 0 {
		t.Errorf("AES-CBC leaked padding bytes into caller's backing array slice capacity: %v", raw[:16])
	}
}

func TestCrypto_SecurityHardening_Points5And6(t *testing.T) {
	// 1. HKDF bound 8160 bita vs 8161 bita (RFC 5869)
	scriptHKDF8160 := `kripto.hkdf("secret", "salt", "info", 8160)`
	resHKDF8160 := runCryptoScript(scriptHKDF8160)
	if resHKDF8160 == nil || resHKDF8160.Type() != evaluator.BYTES_OBJ {
		t.Fatalf("Expected hkdf with 8160 bytes to succeed, got %v", resHKDF8160)
	}
	if len(resHKDF8160.(*evaluator.Bytes).Elements) != 8160 {
		t.Errorf("Expected 8160 bytes, got %d", len(resHKDF8160.(*evaluator.Bytes).Elements))
	}

	scriptHKDF8161 := `kripto.hkdf("secret", "salt", "info", 8161)`
	resHKDF8161 := runCryptoScript(scriptHKDF8161)
	if resHKDF8161 == nil || resHKDF8161.Type() != evaluator.ERROR_OBJ {
		t.Fatalf("Expected hkdf with 8161 bytes to fail, got %v", resHKDF8161)
	}
	if !strings.Contains(resHKDF8161.(*evaluator.Error).Message, "8160") {
		t.Errorf("Expected error message to mention 8160, got: %s", resHKDF8161.(*evaluator.Error).Message)
	}

	// 2. hash_sandi_argon2 strictly validates second argument as Map
	scriptArgonBadOpt := `kripto.hash_sandi_argon2("password", 12345)`
	resArgonBadOpt := runCryptoScript(scriptArgonBadOpt)
	if resArgonBadOpt == nil || resArgonBadOpt.Type() != evaluator.ERROR_OBJ {
		t.Fatalf("Expected hash_sandi_argon2 with non-map option to return error, got %v", resArgonBadOpt)
	}

	scriptArgonGoodOpt := `
	buat hashCustom = kripto.hash_sandi_argon2("sandi", {"waktu": 2, "memori": 65536, "threads": 2})
	kripto.verifikasi_sandi_argon2("sandi", hashCustom)
	`
	resArgonGoodOpt := runCryptoScript(scriptArgonGoodOpt)
	if resArgonGoodOpt != evaluator.TRUE {
		t.Errorf("Expected custom Argon2 options hash verification to be TRUE, got %v", resArgonGoodOpt)
	}

	// 3. verifikasi_sandi_argon2 DoS protections (m > 256MB, t > 10, invalid salt length)
	// m = 300MB = 307200 KB (> 256*1024)
	resArgonMemOver := runCryptoScript(`kripto.verifikasi_sandi_argon2("sandi", "$argon2id$v=19$m=307200,t=2,p=1$c2FsdHNhbHRzYWx0$aGFzaGhhc2hoYXNo")`)
	if resArgonMemOver != evaluator.FALSE {
		t.Errorf("Expected Argon2 verification with >256MB to return FALSE, got %v", resArgonMemOver)
	}

	// t = 11 (> 10)
	resArgonTimeOver := runCryptoScript(`kripto.verifikasi_sandi_argon2("sandi", "$argon2id$v=19$m=65536,t=11,p=1$c2FsdHNhbHRzYWx0$aGFzaGhhc2hoYXNo")`)
	if resArgonTimeOver != evaluator.FALSE {
		t.Errorf("Expected Argon2 verification with >10 iterations to return FALSE, got %v", resArgonTimeOver)
	}

	// 4. XChaCha20-Poly1305 (24-byte nonce) & Backward Compatibility with 12-byte nonce
	scriptXChaCha := `
	buat k = kripto.buat_kunci_chacha()
	buat cipher = kripto.enkripsi_chacha("Rahasia XChaCha20", k, "metadata")
	buat pulih = kripto.dekripsi_chacha(cipher, k, "metadata")
	[cipher, pulih, k]
	`
	resXChaCha := runCryptoScript(scriptXChaCha)
	if resXChaCha == nil || resXChaCha.Type() != evaluator.ARRAY_OBJ {
		t.Fatalf("XChaCha script failed: %v", resXChaCha)
	}
	arrX := resXChaCha.(*evaluator.Array)
	cipherB64 := arrX.Elements[0].(*evaluator.String).Value
	recovered := arrX.Elements[1].(*evaluator.String).Value
	keyB64 := arrX.Elements[2].(*evaluator.String).Value

	if recovered != "Rahasia XChaCha20" {
		t.Errorf("Recovered plaintext mismatch: %s", recovered)
	}
	rawCipher, err := base64.StdEncoding.DecodeString(cipherB64)
	if err != nil {
		t.Fatalf("Base64 decode failed: %v", err)
	}
	// XChaCha20 Nonce is 24 bytes, Poly1305 tag is 16 bytes. Total = 24 + len(data) + 16
	expectedLen := 24 + len("Rahasia XChaCha20") + 16
	if len(rawCipher) != expectedLen {
		t.Errorf("Expected XChaCha20 ciphertext length %d (24-byte nonce + data + 16-byte tag), got %d", expectedLen, len(rawCipher))
	}

	// Backward Compatibility: Legacy 12-byte nonce ChaCha20-Poly1305 decryption in Celang
	keyRaw, err := base64.StdEncoding.DecodeString(keyB64)
	if err != nil {
		t.Fatalf("Base64 decode key failed: %v", err)
	}
	legacyAead, err := chacha20poly1305.New(keyRaw)
	if err != nil {
		t.Fatalf("Legacy AEAD creation failed: %v", err)
	}
	legacyNonce := make([]byte, 12)
	rand.Read(legacyNonce)
	legacyPlain := []byte("Pesan Kompatibel ChaCha Legacy 12-byte")
	legacyAAD := []byte("legacy-aad")
	legacySealed := legacyAead.Seal(nil, legacyNonce, legacyPlain, legacyAAD)
	legacyPayload := append(legacyNonce, legacySealed...)
	legacyPayloadB64 := base64.StdEncoding.EncodeToString(legacyPayload)

	scriptLegacyDec := fmt.Sprintf(`kripto.dekripsi_chacha("%s", "%s", "%s")`, legacyPayloadB64, keyB64, string(legacyAAD))
	resLegacyDec := runCryptoScript(scriptLegacyDec)
	if resLegacyDec == nil || resLegacyDec.Type() != evaluator.STRING_OBJ {
		t.Fatalf("Legacy ChaCha decryption failed: %v", resLegacyDec)
	}
	if resLegacyDec.(*evaluator.String).Value != string(legacyPlain) {
		t.Errorf("Expected legacy plain %s, got %s", string(legacyPlain), resLegacyDec.(*evaluator.String).Value)
	}

	// 5. Constant-time PKCS#7 unpadding & AES-CBC error
	scriptCBCErr := `
	buat k = "01234567890123456789012345678901"
	buat iv = "0123456789012345"
	buat enc = kripto.aes_cbc_enkripsi("pesan", k, iv)
	buat ivSalah = "0000000000000000"
	kripto.aes_cbc_dekripsi(enc, k, ivSalah)
	`
	resCBCErr := runCryptoScript(scriptCBCErr)
	if resCBCErr == nil || resCBCErr.Type() != evaluator.ERROR_OBJ {
		t.Fatalf("Expected AES-CBC decryption with wrong IV to fail, got %v", resCBCErr)
	}
	if !strings.Contains(resCBCErr.(*evaluator.Error).Message, "dekripsi gagal") {
		t.Errorf("Expected unified safe error message, got: %s", resCBCErr.(*evaluator.Error).Message)
	}

	// 6. RSA-PSS Verification with SaltLengthAuto (interoperability with different salt lengths)
	scriptRSASetup := `
	buat keys = kripto.rsa_buat_kunci(2048)
	buat sig = kripto.rsa_tanda_tangan_pss("data uji", keys.kunci_privat)
	buat vHex = kripto.rsa_verifikasi_pss("data uji", sig, keys.kunci_publik)
	[keys.kunci_publik, keys.kunci_privat, sig, vHex]
	`
	resRSASetup := runCryptoScript(scriptRSASetup)
	if resRSASetup == nil || resRSASetup.Type() != evaluator.ARRAY_OBJ {
		t.Fatalf("RSA setup failed: %v", resRSASetup)
	}
	arrRSA := resRSASetup.(*evaluator.Array)
	pubPEM := arrRSA.Elements[0].(*evaluator.String).Value
	sigHex := arrRSA.Elements[2].(*evaluator.String).Value
	vHex := arrRSA.Elements[3]
	if vHex != evaluator.TRUE {
		t.Errorf("RSA PSS hex verify failed")
	}

	// Verify Base64 version of the same signature
	sigRaw, err := hex.DecodeString(sigHex)
	if err != nil {
		t.Fatalf("Hex decode failed: %v", err)
	}
	sigB64 := base64.StdEncoding.EncodeToString(sigRaw)
	scriptRSAB64 := fmt.Sprintf(`kripto.rsa_verifikasi_pss("data uji", "%s", %q)`, sigB64, pubPEM)
	resRSAB64 := runCryptoScript(scriptRSAB64)
	if resRSAB64 != evaluator.TRUE {
		t.Errorf("RSA PSS base64 verify failed: %v", resRSAB64)
	}

	// Verify Base64URL (RawURLEncoding with - and _) of the same signature
	sigB64URL := base64.RawURLEncoding.EncodeToString(sigRaw)
	scriptRSAB64URL := fmt.Sprintf(`kripto.rsa_verifikasi_pss("data uji", "%s", %q)`, sigB64URL, pubPEM)
	resRSAB64URL := runCryptoScript(scriptRSAB64URL)
	if resRSAB64URL != evaluator.TRUE {
		t.Errorf("RSA PSS base64url verify failed: %v", resRSAB64URL)
	}
}



