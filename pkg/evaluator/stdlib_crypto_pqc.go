package evaluator

import (
	"crypto/ecdh"
	"crypto/hkdf"
	"crypto/mlkem"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

// registerPQCCrypto mendaftarkan primitif Post-Quantum Cryptography (Pilar 5):
// ML-KEM-768 (NIST FIPS 203 / Kyber) dan Hybrid KEM (X-Wing: X25519 + ML-KEM-768).
func registerPQCCrypto(env *Environment) {
	// =========================================================================
	// 1. ML-KEM-768 (NIST FIPS 203 / Kyber)
	// =========================================================================

	// pqc_buat_kunci() / pqc_generate_keypair
	pqcBuatKunciFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 0 {
				return newError(nil, "kripto.pqc_buat_kunci tidak memerlukan argumen", "Panggil: kripto.pqc_buat_kunci() atau crypto.pqc_generate_keypair()")
			}
			dk, err := mlkem.GenerateKey768()
			if err != nil {
				return newError(nil, fmt.Sprintf("kripto.pqc_buat_kunci gagal: %v", err), "")
			}
			ek := dk.EncapsulationKey()

			ekHex := hex.EncodeToString(ek.Bytes())
			dkHex := hex.EncodeToString(dk.Bytes())

			pairs := map[string]Object{
				"kunci_enkapsulasi": &String{Value: ekHex},
				"encapsulation_key": &String{Value: ekHex},
				"kunci_dekapsulasi": &String{Value: dkHex},
				"decapsulation_key": &String{Value: dkHex},
				"algoritma":         &String{Value: "ML-KEM-768"},
				"algorithm":         &String{Value: "ML-KEM-768"},
				"standar":           &String{Value: "NIST FIPS 203"},
			}
			order := []string{"kunci_enkapsulasi", "encapsulation_key", "kunci_dekapsulasi", "decapsulation_key", "algoritma", "algorithm", "standar"}
			return &Map{Pairs: pairs, Order: order}
		},
	}
	definePublicFunction(env, "pqc_buat_kunci", pqcBuatKunciFn)
	definePublicFunction(env, "pqc_generate_keypair", pqcBuatKunciFn)

	// pqc_enkapsulasi(kunci_enkapsulasi) / pqc_encapsulate
	pqcEnkapsulasiFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "kripto.pqc_enkapsulasi membutuhkan 1 argumen: (kunci_enkapsulasi_hex)", "")
			}
			ekBytes, errKey := parsePQCBytes(args[0], 1184, "kripto.pqc_enkapsulasi (kunci enkapsulasi)")
			if errKey != nil {
				return errKey
			}

			ek, err := mlkem.NewEncapsulationKey768(ekBytes)
			if err != nil {
				return newError(nil, fmt.Sprintf("kripto.pqc_enkapsulasi: kunci enkapsulasi ML-KEM-768 tidak valid: %v", err), "")
			}

			sharedSecret, ciphertext := ek.Encapsulate()

			pairs := map[string]Object{
				"ciphertext":    &String{Value: hex.EncodeToString(ciphertext)},
				"kunci_bersama": &Bytes{Elements: sharedSecret},
				"shared_secret": &Bytes{Elements: sharedSecret},
			}
			order := []string{"ciphertext", "kunci_bersama", "shared_secret"}
			return &Map{Pairs: pairs, Order: order}
		},
	}
	definePublicFunction(env, "pqc_enkapsulasi", pqcEnkapsulasiFn)
	definePublicFunction(env, "pqc_encapsulate", pqcEnkapsulasiFn)

	// pqc_dekapsulasi(ciphertext, kunci_dekapsulasi) / pqc_decapsulate
	pqcDekapsulasiFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "kripto.pqc_dekapsulasi membutuhkan 2 argumen: (ciphertext_hex, kunci_dekapsulasi_hex)", "")
			}
			cipherBytes, errC := parsePQCBytes(args[0], 1088, "kripto.pqc_dekapsulasi (ciphertext)")
			if errC != nil {
				return errC
			}
			dkBytes, errKey := parsePQCBytes(args[1], 64, "kripto.pqc_dekapsulasi (kunci dekapsulasi)")
			if errKey != nil {
				return errKey
			}

			dk, err := mlkem.NewDecapsulationKey768(dkBytes)
			if err != nil {
				return newError(nil, fmt.Sprintf("kripto.pqc_dekapsulasi: kunci dekapsulasi ML-KEM-768 tidak valid: %v", err), "")
			}

			sharedSecret, err := dk.Decapsulate(cipherBytes)
			if err != nil {
				return newError(nil, fmt.Sprintf("kripto.pqc_dekapsulasi gagal: %v", err), "")
			}
			return &Bytes{Elements: sharedSecret}
		},
	}
	definePublicFunction(env, "pqc_dekapsulasi", pqcDekapsulasiFn)
	definePublicFunction(env, "pqc_decapsulate", pqcDekapsulasiFn)

	// =========================================================================
	// 2. HYBRID KEM (X-WING: X25519 + ML-KEM-768)
	// =========================================================================

	// xwing_buat_kunci() / xwing_generate_keypair
	xwingBuatKunciFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 0 {
				return newError(nil, "kripto.xwing_buat_kunci tidak memerlukan argumen", "")
			}
			// 1. Classical X25519 keypair
			curve := ecdh.X25519()
			xPriv, err := curve.GenerateKey(rand.Reader)
			if err != nil {
				return newError(nil, fmt.Sprintf("kripto.xwing_buat_kunci gagal membuat kunci X25519: %v", err), "")
			}
			xPub := xPriv.PublicKey()

			// 2. Post-Quantum ML-KEM-768 keypair
			mlDk, err := mlkem.GenerateKey768()
			if err != nil {
				return newError(nil, fmt.Sprintf("kripto.xwing_buat_kunci gagal membuat kunci ML-KEM: %v", err), "")
			}
			mlEk := mlDk.EncapsulationKey()

			// Concatenate public key: 32 bytes X25519 + 1184 bytes ML-KEM-768 = 1216 bytes
			jointPub := append(xPub.Bytes(), mlEk.Bytes()...)
			// Concatenate private key: 32 bytes X25519 + 64 bytes ML-KEM-768 seed = 96 bytes
			jointPriv := append(xPriv.Bytes(), mlDk.Bytes()...)

			pubHex := hex.EncodeToString(jointPub)
			privHex := hex.EncodeToString(jointPriv)

			pairs := map[string]Object{
				"kunci_publik": &String{Value: pubHex},
				"public_key":   &String{Value: pubHex},
				"kunci_privat": &String{Value: privHex},
				"private_key":  &String{Value: privHex},
				"algoritma":    &String{Value: "Hybrid X-Wing (X25519 + ML-KEM-768)"},
				"algorithm":    &String{Value: "Hybrid X-Wing (X25519 + ML-KEM-768)"},
			}
			order := []string{"kunci_publik", "public_key", "kunci_privat", "private_key", "algoritma", "algorithm"}
			return &Map{Pairs: pairs, Order: order}
		},
	}
	definePublicFunction(env, "xwing_buat_kunci", xwingBuatKunciFn)
	definePublicFunction(env, "xwing_generate_keypair", xwingBuatKunciFn)

	// xwing_enkapsulasi(kunci_publik_xwing) / xwing_encapsulate
	xwingEnkapsulasiFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "kripto.xwing_enkapsulasi membutuhkan 1 argumen: (kunci_publik_xwing_hex)", "")
			}
			pubBytes, errKey := parsePQCBytes(args[0], 1216, "kripto.xwing_enkapsulasi (kunci publik)")
			if errKey != nil {
				return errKey
			}

			peerXPubBytes := pubBytes[:32]
			peerMLEkBytes := pubBytes[32:]

			// 1. Classical Ephemeral X25519 ECDH
			curve := ecdh.X25519()
			peerXPub, err := curve.NewPublicKey(peerXPubBytes)
			if err != nil {
				return newError(nil, fmt.Sprintf("kripto.xwing_enkapsulasi: kunci publik X25519 lawan tidak valid: %v", err), "")
			}
			ephPriv, err := curve.GenerateKey(rand.Reader)
			if err != nil {
				return newError(nil, fmt.Sprintf("kripto.xwing_enkapsulasi gagal membuat ephemeral X25519: %v", err), "")
			}
			ephPub := ephPriv.PublicKey()
			xSharedSecret, err := ephPriv.ECDH(peerXPub)
			if err != nil {
				return newError(nil, fmt.Sprintf("kripto.xwing_enkapsulasi: ECDH gagal: %v", err), "")
			}

			// 2. Post-Quantum ML-KEM-768 Encapsulation
			mlEk, err := mlkem.NewEncapsulationKey768(peerMLEkBytes)
			if err != nil {
				return newError(nil, fmt.Sprintf("kripto.xwing_enkapsulasi: kunci ML-KEM lawan tidak valid: %v", err), "")
			}
			mlSharedSecret, mlCiphertext := mlEk.Encapsulate()

			// 3. Combine both shared secrets via HKDF-SHA256
			combinedSecret := append(xSharedSecret, mlSharedSecret...)
			jointKey, err := hkdf.Key(sha256.New, combinedSecret, nil, "X-Wing-v1", 32)
			if err != nil {
				return newError(nil, fmt.Sprintf("kripto.xwing_enkapsulasi: HKDF derivation gagal: %v", err), "")
			}

			// Ciphertext = ephPub (32 bytes) + mlCiphertext (1088 bytes) = 1120 bytes
			jointCiphertext := append(ephPub.Bytes(), mlCiphertext...)

			pairs := map[string]Object{
				"ciphertext":    &String{Value: hex.EncodeToString(jointCiphertext)},
				"kunci_bersama": &Bytes{Elements: jointKey},
				"shared_secret": &Bytes{Elements: jointKey},
			}
			order := []string{"ciphertext", "kunci_bersama", "shared_secret"}
			return &Map{Pairs: pairs, Order: order}
		},
	}
	definePublicFunction(env, "xwing_enkapsulasi", xwingEnkapsulasiFn)
	definePublicFunction(env, "xwing_encapsulate", xwingEnkapsulasiFn)

	// xwing_dekapsulasi(ciphertext, kunci_privat_xwing) / xwing_decapsulate
	xwingDekapsulasiFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "kripto.xwing_dekapsulasi membutuhkan 2 argumen: (ciphertext_hex, kunci_privat_xwing_hex)", "")
			}
			cipherBytes, errC := parsePQCBytes(args[0], 1120, "kripto.xwing_dekapsulasi (ciphertext)")
			if errC != nil {
				return errC
			}
			privBytes, errKey := parsePQCBytes(args[1], 96, "kripto.xwing_dekapsulasi (kunci privat)")
			if errKey != nil {
				return errKey
			}

			ephPubBytes := cipherBytes[:32]
			mlCiphertext := cipherBytes[32:]

			xPrivBytes := privBytes[:32]
			mlDkBytes := privBytes[32:]

			// 1. Classical X25519 ECDH
			curve := ecdh.X25519()
			xPriv, err := curve.NewPrivateKey(xPrivBytes)
			if err != nil {
				return newError(nil, fmt.Sprintf("kripto.xwing_dekapsulasi: kunci privat X25519 tidak valid: %v", err), "")
			}
			ephPub, err := curve.NewPublicKey(ephPubBytes)
			if err != nil {
				return newError(nil, fmt.Sprintf("kripto.xwing_dekapsulasi: ephemeral public key X25519 tidak valid: %v", err), "")
			}
			xSharedSecret, err := xPriv.ECDH(ephPub)
			if err != nil {
				return newError(nil, fmt.Sprintf("kripto.xwing_dekapsulasi: ECDH gagal: %v", err), "")
			}

			// 2. Post-Quantum ML-KEM-768 Decapsulation
			mlDk, err := mlkem.NewDecapsulationKey768(mlDkBytes)
			if err != nil {
				return newError(nil, fmt.Sprintf("kripto.xwing_dekapsulasi: kunci dekapsulasi ML-KEM tidak valid: %v", err), "")
			}
			mlSharedSecret, err := mlDk.Decapsulate(mlCiphertext)
			if err != nil {
				return newError(nil, fmt.Sprintf("kripto.xwing_dekapsulasi: ML-KEM dekapsulasi gagal: %v", err), "")
			}

			// 3. Combine both shared secrets via HKDF-SHA256
			combinedSecret := append(xSharedSecret, mlSharedSecret...)
			jointKey, err := hkdf.Key(sha256.New, combinedSecret, nil, "X-Wing-v1", 32)
			if err != nil {
				return newError(nil, fmt.Sprintf("kripto.xwing_dekapsulasi: HKDF derivation gagal: %v", err), "")
			}

			return &Bytes{Elements: jointKey}
		},
	}
	definePublicFunction(env, "xwing_dekapsulasi", xwingDekapsulasiFn)
	definePublicFunction(env, "xwing_decapsulate", xwingDekapsulasiFn)
}

func parsePQCBytes(arg Object, expectedLen int, funcName string) ([]byte, Object) {
	switch v := arg.(type) {
	case *Bytes:
		if len(v.Elements) != expectedLen {
			return nil, newError(nil, fmt.Sprintf("%s: panjang bita harus %d bita, menerima %d bita", funcName, expectedLen, len(v.Elements)), "")
		}
		return v.Elements, nil
	case *String:
		s := strings.TrimSpace(v.Value)
		if len(s) == expectedLen*2 {
			if hx, err := hex.DecodeString(s); err == nil && len(hx) == expectedLen {
				return hx, nil
			}
		}
		if b64, err := decodeBase64Flex(s); err == nil && len(b64) == expectedLen {
			return b64, nil
		}
		return nil, newError(nil, fmt.Sprintf("%s: format kunci/data tidak valid (diharapkan heksadesimal %d karakter atau Base64 %d bita)", funcName, expectedLen*2, expectedLen), "")
	default:
		return nil, newError(nil, fmt.Sprintf("%s: argumen harus berupa teks atau Bita", funcName), "")
	}
}
