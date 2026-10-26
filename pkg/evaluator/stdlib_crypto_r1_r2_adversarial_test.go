package evaluator_test

import (
	"celang/pkg/evaluator"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
	"testing"
)

// ============================================================================
// ADVERSARIAL TEST SUITE: R1 (AES-256-GCM AEAD Symmetric Encryption)
// ============================================================================

// TestAdversarial_AESGCM_BitTampering systematically flips single bits across
// different regions (nonce, ciphertext, authentication tag) and verifies that
// decryption fails safely with an empathetic error and NEVER panics or leaks data.
func TestAdversarial_AESGCM_BitTampering(t *testing.T) {
	initScript := `
	buat kunci = kripto.buat_kunci()
	buat pesan = "Kerahasiaan data tingkat tinggi Celang 2026!"
	buat terenkripsi = kripto.enkripsi(pesan, kunci)
	[kunci, terenkripsi]
	`
	resInit := runCryptoScript(initScript)
	arrInit, ok := resInit.(*evaluator.Array)
	if !ok {
		t.Fatalf("inisialisasi AES-GCM gagal: %v", resInit)
	}
	kunci := arrInit.Elements[0].Inspect()
	b64Ciphertext := arrInit.Elements[1].Inspect()

	payload, err := base64.StdEncoding.DecodeString(b64Ciphertext)
	if err != nil {
		t.Fatalf("gagal mendekode base64 ciphertext: %v", err)
	}

	// Payload structure: 12 bytes nonce + len(plaintext) ciphertext + 16 bytes tag
	totalBytes := len(payload)
	if totalBytes < 28 {
		t.Fatalf("panjang payload tidak valid: %d (minimal 28)", totalBytes)
	}

	regions := []struct {
		name      string
		startByte int
		endByte   int
	}{
		{"Nonce (0..11)", 0, 11},
		{"Nonce-Ciphertext Boundary (11..12)", 11, 12},
		{"Ciphertext Region", 12, totalBytes - 17},
		{"Ciphertext-Tag Boundary", totalBytes - 17, totalBytes - 16},
		{"Tag Region (last 16 bytes)", totalBytes - 16, totalBytes - 1},
	}

	tamperCount := 0
	for _, reg := range regions {
		t.Run("Region_"+reg.name, func(t *testing.T) {
			for byteIdx := reg.startByte; byteIdx <= reg.endByte && byteIdx < totalBytes; byteIdx++ {
				for _, bitShift := range []uint{0, 3, 7} {
					corrupted := make([]byte, len(payload))
					copy(corrupted, payload)
					corrupted[byteIdx] ^= (1 << bitShift)

					corruptedB64 := base64.StdEncoding.EncodeToString(corrupted)
					script := fmt.Sprintf(`
					buat hasil = ""
					coba {
						kripto.dekripsi("%s", "%s")
					} tangkap (e) {
						hasil = e
					}
					hasil
					`, corruptedB64, kunci)

					res := runCryptoScript(script)
					errStr := res.Inspect()

					if !strings.Contains(errStr, "dekripsi gagal") && !strings.Contains(errStr, "otentikasi gagal") {
						t.Errorf("byte %d (bit %d) tampered payload did not yield safe empathetic error: got %s",
							byteIdx, bitShift, errStr)
					}
					tamperCount++
				}
			}
		})
	}

	// Also verify bit tampering on empty plaintext payload (exactly 28 bytes: 12 nonce + 16 tag)
	t.Run("EmptyPlaintext_28Bytes_AllBits", func(t *testing.T) {
		emptyEncScript := fmt.Sprintf(`kripto.enkripsi("", "%s")`, kunci)
		resEmpty := runCryptoScript(emptyEncScript)
		b64Empty := resEmpty.Inspect()
		emptyPayload, err := base64.StdEncoding.DecodeString(b64Empty)
		if err != nil {
			t.Fatalf("decode empty ciphertext failed: %v", err)
		}
		if len(emptyPayload) != 28 {
			t.Fatalf("empty payload length expected 28, got %d", len(emptyPayload))
		}

		for i := 0; i < len(emptyPayload); i++ {
			for _, bit := range []uint{0, 7} {
				corrupted := make([]byte, len(emptyPayload))
				copy(corrupted, emptyPayload)
				corrupted[i] ^= (1 << bit)
				corruptedB64 := base64.StdEncoding.EncodeToString(corrupted)

				script := fmt.Sprintf(`
				buat hasil = ""
				coba {
					kripto.dekripsi("%s", "%s")
				} tangkap (e) {
					hasil = e
				}
				hasil
				`, corruptedB64, kunci)

				res := runCryptoScript(script)
				errStr := res.Inspect()
				if !strings.Contains(errStr, "dekripsi gagal") {
					t.Errorf("empty payload byte %d bit %d tamper failed: %s", i, bit, errStr)
				}
			}
		}
	})

	t.Logf("Total 1-bit tampering variations tested and verified: %d", tamperCount)
}

// TestAdversarial_AESGCM_TruncatedCiphertext tests truncated payloads of all
// sizes from 0 to 27 bytes (below the 28-byte minimum) and corrupted base64 strings.
func TestAdversarial_AESGCM_TruncatedCiphertext(t *testing.T) {
	keyScript := `kripto.buat_kunci()`
	kunci := runCryptoScript(keyScript).Inspect()

	// 1. Raw byte lengths from 0 to 27 bytes encoded as Base64
	for length := 0; length < 28; length++ {
		rawBytes := make([]byte, length)
		b64 := base64.StdEncoding.EncodeToString(rawBytes)

		script := fmt.Sprintf(`
		buat err = ""
		coba {
			kripto.dekripsi("%s", "%s")
		} tangkap (e) {
			err = e
		}
		err
		`, b64, kunci)

		res := runCryptoScript(script)
		errStr := res.Inspect()
		if !strings.Contains(errStr, "format data tersandi tidak valid") && !strings.Contains(errStr, "terlalu pendek") {
			t.Errorf("length %d truncated ciphertext did not yield empathetic error: got %s", length, errStr)
		}
	}

	// 2. Truncated Bita objects from 0 to 27 bytes
	for length := 0; length < 28; length++ {
		script := fmt.Sprintf(`
		buat err = ""
		coba {
			kripto.dekripsi(kripto.acak_bita(%d), "%s")
		} tangkap (e) {
			err = e
		}
		err
		`, length, kunci)

		res := runCryptoScript(script)
		errStr := res.Inspect()
		if !strings.Contains(errStr, "format data tersandi tidak valid") && !strings.Contains(errStr, "terlalu pendek") {
			t.Errorf("Bita length %d truncated ciphertext did not yield empathetic error: got %s", length, errStr)
		}
	}

	// 3. Corrupted base64 inputs (malformed characters, incomplete paddings)
	malformedInputs := []string{
		"!",
		"===",
		"abc",
		"dGVz!",
		"????",
		"dGVzdA=====",
		"    \t\n   ",
	}
	for _, badB64 := range malformedInputs {
		script := fmt.Sprintf(`
		buat err = ""
		coba {
			kripto.dekripsi("%s", "%s")
		} tangkap (e) {
			err = e
		}
		err
		`, badB64, kunci)

		res := runCryptoScript(script)
		errStr := res.Inspect()
		if !strings.Contains(errStr, "format") {
			t.Errorf("malformed base64 '%s' did not return empathetic format error: got %s", badB64, errStr)
		}
	}
}

// TestAdversarial_AESGCM_CorruptedNonces tests nonces replaced with all zeros,
// all 0xFF, or nonces transplanted from a different encryption session.
func TestAdversarial_AESGCM_CorruptedNonces(t *testing.T) {
	initScript := `
	buat kunci = kripto.buat_kunci()
	buat enc1 = kripto.enkripsi("Pesan rahasia 1", kunci)
	buat enc2 = kripto.enkripsi("Pesan rahasia 2", kunci)
	[kunci, enc1, enc2]
	`
	resInit := runCryptoScript(initScript)
	arrInit := resInit.(*evaluator.Array)
	kunci := arrInit.Elements[0].Inspect()
	p1, _ := base64.StdEncoding.DecodeString(arrInit.Elements[1].Inspect())
	p2, _ := base64.StdEncoding.DecodeString(arrInit.Elements[2].Inspect())

	// Nonce transplant: swap 12-byte nonce of enc1 into enc2
	transplanted := make([]byte, len(p2))
	copy(transplanted, p2)
	copy(transplanted[:12], p1[:12])
	b64Trans := base64.StdEncoding.EncodeToString(transplanted)

	// All zero nonce
	allZeros := make([]byte, len(p1))
	copy(allZeros, p1)
	for i := 0; i < 12; i++ {
		allZeros[i] = 0x00
	}
	b64Zeros := base64.StdEncoding.EncodeToString(allZeros)

	// All 0xFF nonce
	allFF := make([]byte, len(p1))
	copy(allFF, p1)
	for i := 0; i < 12; i++ {
		allFF[i] = 0xff
	}
	b64FF := base64.StdEncoding.EncodeToString(allFF)

	cases := map[string]string{
		"TransplantedNonce": b64Trans,
		"AllZerosNonce":     b64Zeros,
		"AllFFNonce":        b64FF,
	}

	for name, b64Payload := range cases {
		script := fmt.Sprintf(`
		buat err = ""
		coba {
			kripto.dekripsi("%s", "%s")
		} tangkap (e) {
			err = e
		}
		err
		`, b64Payload, kunci)

		res := runCryptoScript(script)
		errStr := res.Inspect()
		if !strings.Contains(errStr, "dekripsi gagal") {
			t.Errorf("corrupted nonce scenario '%s' failed to reject with safe empathetic error: got %s", name, errStr)
		}
	}
}

// TestAdversarial_AESGCM_BadKeys tests various invalid keys (short, long, wrong format, wrong type).
func TestAdversarial_AESGCM_BadKeys(t *testing.T) {
	validKey := runCryptoScript(`kripto.buat_kunci()`).Inspect()
	ciphertext := runCryptoScript(fmt.Sprintf(`kripto.enkripsi("halo", "%s")`, validKey)).Inspect()

	// 1. Raw string keys with invalid lengths
	invalidStrKeys := []string{
		"",
		"a",
		"12345678",
		"123456789012345",                   // 15 bytes
		"1234567890123456",                  // 16 bytes (AES-128, should reject because AES-256 requires 32 bytes)
		"123456789012345678901234",          // 24 bytes (AES-192, should reject)
		"1234567890123456789012345678901",   // 31 bytes
		"123456789012345678901234567890123", // 33 bytes
		strings.Repeat("X", 63),             // 63 chars
		strings.Repeat("X", 65),             // 65 chars
		strings.Repeat("X", 128),
		strings.Repeat("X", 1024),
	}

	for _, badKey := range invalidStrKeys {
		script := fmt.Sprintf(`
		buat errEnc = ""
		buat errDec = ""
		coba {
			kripto.enkripsi("data", "%s")
		} tangkap (e) {
			errEnc = e
		}
		coba {
			kripto.dekripsi("%s", "%s")
		} tangkap (e) {
			errDec = e
		}
		[errEnc, errDec]
		`, badKey, ciphertext, badKey)

		res := runCryptoScript(script)
		arr, ok := res.(*evaluator.Array)
		if !ok {
			t.Fatalf("script evaluation failed for bad key len %d: %v", len(badKey), res)
		}
		errEnc := arr.Elements[0].Inspect()
		errDec := arr.Elements[1].Inspect()

		if !strings.Contains(errEnc, "format atau panjang kunci tidak valid") {
			t.Errorf("enkripsi with key len %d did not return empathetic format error: got %s", len(badKey), errEnc)
		}
		if !strings.Contains(errDec, "format atau panjang kunci tidak valid") {
			t.Errorf("dekripsi with key len %d did not return empathetic format error: got %s", len(badKey), errDec)
		}
	}

	// 2. Non-hex characters of length 64
	badHexKey := strings.Repeat("z", 64)
	scriptBadHex := fmt.Sprintf(`
	buat err = ""
	coba {
		kripto.enkripsi("data", "%s")
	} tangkap (e) {
		err = e
	}
	err
	`, badHexKey)
	resBadHex := runCryptoScript(scriptBadHex).Inspect()
	if !strings.Contains(resBadHex, "format atau panjang kunci tidak valid") {
		t.Errorf("64-char non-hex key did not return empathetic error: got %s", resBadHex)
	}

	// 3. Invalid types for key (Integer, Boolean, Array, Map)
	badTypeKeys := []string{
		"12345",
		"benar",
		"salah",
		"[1, 2, 3]",
		`{"k": "v"}`,
	}
	for _, typeKey := range badTypeKeys {
		script := fmt.Sprintf(`
		buat err = ""
		coba {
			kripto.enkripsi("data", %s)
		} tangkap (e) {
			err = e
		}
		err
		`, typeKey)

		res := runCryptoScript(script).Inspect()
		if !strings.Contains(res, "harus berupa teks atau Bita") {
			t.Errorf("bad key type %s did not return type error: got %s", typeKey, res)
		}
	}

	// 4. Bita keys of invalid lengths
	for _, badLen := range []int{0, 1, 16, 24, 31, 33, 64} {
		script := fmt.Sprintf(`
		buat err = ""
		coba {
			kripto.enkripsi("data", kripto.acak_bita(%d))
		} tangkap (e) {
			err = e
		}
		err
		`, badLen)

		res := runCryptoScript(script).Inspect()
		if !strings.Contains(res, "panjang kunci Bita tidak valid") {
			t.Errorf("Bita key of len %d did not return length error: got %s", badLen, res)
		}
	}
}

// TestAdversarial_AESGCM_TypeConfusion tests passing invalid types for ciphertext or data.
func TestAdversarial_AESGCM_TypeConfusion(t *testing.T) {
	key := runCryptoScript(`kripto.buat_kunci()`).Inspect()

	invalidArgs := []string{
		"12345",
		"benar",
		"[1, 2, 3]",
		`{"teks": "rahasia"}`,
	}

	for _, arg := range invalidArgs {
		script := fmt.Sprintf(`
		buat errData = ""
		buat errCipher = ""
		coba {
			kripto.enkripsi(%s, "%s")
		} tangkap (e) {
			errData = e
		}
		coba {
			kripto.dekripsi(%s, "%s")
		} tangkap (e) {
			errCipher = e
		}
		[errData, errCipher]
		`, arg, key, arg, key)

		res := runCryptoScript(script)
		arr, ok := res.(*evaluator.Array)
		if !ok {
			t.Fatalf("type confusion test failed: %v", res)
		}
		errData := arr.Elements[0].Inspect()
		errCipher := arr.Elements[1].Inspect()

		if !strings.Contains(errData, "harus berupa teks atau Bita") {
			t.Errorf("enkripsi with data %s did not return type error: got %s", arg, errData)
		}
		if !strings.Contains(errCipher, "harus berupa teks atau Bita") {
			t.Errorf("dekripsi with ciphertext %s did not return type error: got %s", arg, errCipher)
		}
	}
}

// ============================================================================
// ADVERSARIAL TEST SUITE: R2 (Password Hashing PBKDF2 & Key Derivation)
// ============================================================================

// TestAdversarial_PasswordHashing_SaltUniqueness verifies salt uniqueness and
// randomness across 30 consecutive generations of password hashing.
func TestAdversarial_PasswordHashing_SaltUniqueness(t *testing.T) {
	const iterations = 30
	salts := make(map[string]bool)
	hashes := make(map[string]bool)

	for i := 0; i < iterations; i++ {
		script := `kripto.hash_sandi("SamaSemuaPassword123!")`
		res := runCryptoScript(script)
		phc := res.Inspect()

		if hashes[phc] {
			t.Fatalf("iteration %d: duplicate full PHC string produced: %s", i, phc)
		}
		hashes[phc] = true

		// Parse PHC: $pbkdf2-sha256$i=100000$<salt>$<hash>
		parts := strings.Split(phc, "$")
		if len(parts) != 5 || parts[1] != "pbkdf2-sha256" {
			t.Fatalf("iteration %d: invalid PHC format: %s", i, phc)
		}

		saltB64 := parts[3]
		saltBytes, err := base64.RawStdEncoding.DecodeString(saltB64)
		if err != nil {
			t.Fatalf("iteration %d: invalid base64 salt '%s': %v", i, saltB64, err)
		}
		if len(saltBytes) != 16 {
			t.Fatalf("iteration %d: salt length expected 16 bytes, got %d", i, len(saltBytes))
		}

		saltHex := hex.EncodeToString(saltBytes)
		if salts[saltHex] {
			t.Fatalf("iteration %d: duplicate salt collision detected: %s", i, saltHex)
		}
		salts[saltHex] = true
	}

	if len(salts) != iterations {
		t.Errorf("expected %d unique salts, got %d", iterations, len(salts))
	}
	t.Logf("Successfully verified %d unique 16-byte random salts across consecutive hash_sandi calls", iterations)
}

// TestAdversarial_PasswordHashing_SubtlyWrongPasswords tests that verify_password
// returns FALSE for all subtly mutated passwords (1-char edits, case, whitespace, nulls).
func TestAdversarial_PasswordHashing_SubtlyWrongPasswords(t *testing.T) {
	originalPassword := "C3l@ng-S3cur3-P@ssw0rd!"

	initScript := fmt.Sprintf(`kripto.hash_sandi("%s")`, originalPassword)
	phc := runCryptoScript(initScript).Inspect()

	// Mutated password test cases (string literals without raw null bytes)
	mutations := []struct {
		name     string
		password string
	}{
		{"LowerInitial", "c3l@ng-S3cur3-P@ssw0rd!"},
		{"UpperSubstitution", "C3L@ng-S3cur3-P@ssw0rd!"},
		{"MiddleDigitToLetter", "C3l@ng-Secur3-P@ssw0rd!"},
		{"LastCharMissing", "C3l@ng-S3cur3-P@ssw0rd"},
		{"LastCharDifferent", "C3l@ng-S3cur3-P@ssw0rd?"},
		{"FirstCharInserted", "CC3l@ng-S3cur3-P@ssw0rd!"},
		{"FirstCharDeleted", "3l@ng-S3cur3-P@ssw0rd!"},
		{"TransposedAdjacent", "3Cl@ng-S3cur3-P@ssw0rd!"},
		{"TrailingSpace", "C3l@ng-S3cur3-P@ssw0rd! "},
		{"LeadingSpace", " C3l@ng-S3cur3-P@ssw0rd!"},
		{"InternalSpace", "C3l@ng-S3cur3- P@ssw0rd!"},
		{"ZeroWidthSpace", "C3l@ng-S3cur3-P@ssw0rd!\u200B"},
		{"AllUpperCase", "C3L@NG-S3CUR3-P@SSW0RD!"},
		{"AllLowerCase", "c3l@ng-s3cur3-p@ssw0rd!"},
		{"CompletelyWrong", "PasswordRahasiaSalahTotal"},
		{"SingleChar", "C"},
		{"EmptyString", ""},
	}

	// 1. Verify original matches
	okScript := fmt.Sprintf(`kripto.verifikasi_sandi("%s", "%s")`, originalPassword, phc)
	if runCryptoScript(okScript) != evaluator.TRUE {
		t.Fatalf("original password must verify as true")
	}

	// 2. Verify all string mutations fail safely
	for _, mut := range mutations {
		t.Run("Mutation_"+mut.name, func(t *testing.T) {
			script := fmt.Sprintf(`
			buat hasil = kripto.verifikasi_sandi("%s", "%s")
			buat hasilInggris = crypto.verify_password("%s", "%s")
			[hasil, hasilInggris]
			`, mut.password, phc, mut.password, phc)

			res := runCryptoScript(script)
			arr, ok := res.(*evaluator.Array)
			if !ok {
				t.Fatalf("verification script failed for %s: %v", mut.name, res)
			}

			if arr.Elements[0] != evaluator.FALSE || arr.Elements[1] != evaluator.FALSE {
				t.Errorf("mutation '%s' ('%s') incorrectly verified as TRUE!", mut.name, mut.password)
			}
		})
	}

	// 3. Verify null-byte injection attacks using Bita representations
	bitaNullTests := []struct {
		name        string
		data        []byte
		targetPhc   string
		expectMatch bool
	}{
		{"NullBytePrepended", append([]byte{0x00}, []byte(originalPassword)...), phc, false},
		{"NullByteEmbedded", append(append([]byte(originalPassword[:8]), 0x00), []byte(originalPassword[8:])...), phc, false},
	}

	for _, tc := range bitaNullTests {
		t.Run("Mutation_"+tc.name, func(t *testing.T) {
			bitaElems := make([]string, len(tc.data))
			for i, b := range tc.data {
				bitaElems[i] = fmt.Sprintf("%d", b)
			}
			script := fmt.Sprintf(`
			buat pwBita = bita([%s])
			buat hasil = kripto.verifikasi_sandi(pwBita, "%s")
			buat hasilInggris = crypto.verify_password(pwBita, "%s")
			[hasil, hasilInggris]
			`, strings.Join(bitaElems, ", "), tc.targetPhc, tc.targetPhc)

			res := runCryptoScript(script)
			arr, ok := res.(*evaluator.Array)
			if !ok {
				t.Fatalf("bita null verification script failed for %s: %v", tc.name, res)
			}

			expectedBool := evaluator.FALSE
			if tc.expectMatch {
				expectedBool = evaluator.TRUE
			}
			if arr.Elements[0] != expectedBool || arr.Elements[1] != expectedBool {
				t.Errorf("bita null mutation '%s' expected %v, got [0]=%s, [1]=%s",
					tc.name, tc.expectMatch, arr.Elements[0].Inspect(), arr.Elements[1].Inspect())
			}
		})
	}

	// 4. Test 64-byte password boundary where HMAC pre-hashes the key with SHA-256 (RFC 2104 Section 2)
	t.Run("Mutation_64ByteBoundary_AppendedNull", func(t *testing.T) {
		longPw := strings.Repeat("A", 64) // exactly 64 bytes
		longPhc := runCryptoScript(fmt.Sprintf(`kripto.hash_sandi("%s")`, longPw)).Inspect()

		longPwBytesWithNull := append([]byte(longPw), 0x00)
		bitaElems := make([]string, len(longPwBytesWithNull))
		for i, b := range longPwBytesWithNull {
			bitaElems[i] = fmt.Sprintf("%d", b)
		}
		script := fmt.Sprintf(`
		buat pwBita = bita([%s])
		kripto.verifikasi_sandi(pwBita, "%s")
		`, strings.Join(bitaElems, ", "), longPhc)

		res := runCryptoScript(script)
		if res != evaluator.FALSE {
			t.Errorf("64-byte password with appended null byte must return FALSE: got %v", res)
		}
	})
}

// TestAdversarial_PasswordHashing_MalformedPHCString tests that verifikasi_sandi
// safely handles corrupted, truncated, oversized, or syntactically invalid PHC strings.
func TestAdversarial_PasswordHashing_MalformedPHCString(t *testing.T) {
	validSaltB64 := base64.RawStdEncoding.EncodeToString([]byte("1234567890123456"))
	validHashB64 := base64.RawStdEncoding.EncodeToString(make([]byte, 32))

	malformedStrings := []struct {
		name string
		phc  string
	}{
		{"EmptyString", ""},
		{"JustDollar", "$"},
		{"TwoDollars", "$$"},
		{"ThreeDollars", "$$$"},
		{"FourDollars", "$$$$"},
		{"FiveDollars", "$$$$$"},
		{"WrongAlgorithm_Bcrypt", "$2a$12$R9h/cIPz0gi.URNNXRkh2OPST9/PgBkqquzi.Ss7KIUgO2t0jWMUW"},
		{"WrongAlgorithm_Argon2id", "$argon2id$v=19$m=65536,t=3,p=4$c29tZXNhbHQ$i90..."},
		{"WrongAlgorithm_PBKDF2_SHA512", "$pbkdf2-sha512$i=100000$" + validSaltB64 + "$" + validHashB64},
		{"IncompletePrefix", "$pbkdf2$i=1000$" + validSaltB64 + "$" + validHashB64},
		{"MissingHashField", "$pbkdf2-sha256$i=100000$" + validSaltB64},
		{"ExtraTrailingField", "$pbkdf2-sha256$i=100000$" + validSaltB64 + "$" + validHashB64 + "$extra"},
		{"NonNumericIterations", "$pbkdf2-sha256$i=notanumber$" + validSaltB64 + "$" + validHashB64},
		{"ZeroIterations", "$pbkdf2-sha256$i=0$" + validSaltB64 + "$" + validHashB64},
		{"NegativeIterations", "$pbkdf2-sha256$i=-5000$" + validSaltB64 + "$" + validHashB64},
		{"HugeOverflowIterations", "$pbkdf2-sha256$i=9999999999999999999999999999999999999999$" + validSaltB64 + "$" + validHashB64},
		{"ExceedsDosIterationCap", "$pbkdf2-sha256$i=20000000$" + validSaltB64 + "$" + validHashB64},
		{"MissingIterKey", "$pbkdf2-sha256$cost=100000$" + validSaltB64 + "$" + validHashB64},
		{"EmptyParamsField", "$pbkdf2-sha256$$" + validSaltB64 + "$" + validHashB64},
		{"EmptySaltField", "$pbkdf2-sha256$i=100000$$" + validHashB64},
		{"InvalidBase64Salt", "$pbkdf2-sha256$i=100000$???salt???$" + validHashB64},
		{"EmptyHashField", "$pbkdf2-sha256$i=100000$" + validSaltB64 + "$"},
		{"InvalidBase64Hash", "$pbkdf2-sha256$i=100000$" + validSaltB64 + "$!@#$%^&*()"},
		{"PlainTextWithoutDollars", "ini-bukan-format-phc-sama-sekali"},
		{"LeadingWhitespace", "  $pbkdf2-sha256$i=100000$" + validSaltB64 + "$" + validHashB64},
		{"TrailingWhitespace", "$pbkdf2-sha256$i=100000$" + validSaltB64 + "$" + validHashB64 + "  \n"},
	}

	for _, tc := range malformedStrings {
		t.Run(tc.name, func(t *testing.T) {
			script := fmt.Sprintf(`
			buat hasil = kripto.verifikasi_sandi("kata-sandi", "%s")
			hasil
			`, strings.ReplaceAll(tc.phc, `"`, `\"`))

			res := runCryptoScript(script)
			if res != evaluator.FALSE {
				t.Errorf("malformed PHC '%s' did not return FALSE safely: got %v", tc.name, res)
			}
		})
	}

	// Non-string arguments to verifikasi_sandi
	badTypeArgs := []string{
		`kripto.verifikasi_sandi("sandi", 12345)`,
		`kripto.verifikasi_sandi("sandi", benar)`,
		`kripto.verifikasi_sandi("sandi", [1, 2, 3])`,
		`kripto.verifikasi_sandi("sandi", {"hash": "abc"})`,
		`kripto.verifikasi_sandi(12345, "$pbkdf2-sha256$i=1000$salt$hash")`,
	}

	for i, expr := range badTypeArgs {
		script := fmt.Sprintf(`
		buat err = ""
		coba {
			%s
		} tangkap (e) {
			err = e
		}
		err
		`, expr)

		res := runCryptoScript(script).Inspect()
		if !strings.Contains(res, "harus berupa teks") {
			t.Errorf("case %d (%s) did not return type error: got %s", i, expr, res)
		}
	}
}

// TestAdversarial_ConstantTimeCompare_TimingSafe tests constant-time comparison
// across identical, different length, and 1-bit differing inputs.
func TestAdversarial_ConstantTimeCompare_TimingSafe(t *testing.T) {
	script := `
	buat s1 = "KunciAutentikasiPalingRahasia2026"
	buat s2 = "KunciAutentikasiPalingRahasia2026"
	buat s3 = "kunciAutentikasiPalingRahasia2026" // 1-bit flip at start
	buat s4 = "KunciAutentikasiPalingRahasia2027" // 1-bit flip at end
	buat s5 = "KunciAutentikasi"                  // shorter

	buat r1 = kripto.konstanta_sama(s1, s2)
	buat r2 = kripto.konstanta_sama(s1, s3)
	buat r3 = kripto.konstanta_sama(s1, s4)
	buat r4 = kripto.konstanta_sama(s1, s5)
	buat r5 = crypto.timing_safe_equal("", "")
	buat r6 = crypto.timing_safe_equal("", "a")

	// Bita comparison
	buat r7 = kripto.banding_aman(bita(s1), bita(s2))
	buat r8 = kripto.banding_aman(bita(s1), bita(s3))

	[r1, r2, r3, r4, r5, r6, r7, r8]
	`
	res := runCryptoScript(script)
	arr, ok := res.(*evaluator.Array)
	if !ok {
		t.Fatalf("constant time compare script failed: %v", res)
	}

	expected := []evaluator.Object{
		evaluator.TRUE,
		evaluator.FALSE,
		evaluator.FALSE,
		evaluator.FALSE,
		evaluator.TRUE,
		evaluator.FALSE,
		evaluator.TRUE,
		evaluator.FALSE,
	}

	for i, exp := range expected {
		if arr.Elements[i] != exp {
			t.Errorf("assertion %d expected %v, got %v", i, exp, arr.Elements[i])
		}
	}
}

// TestAdversarial_DeriveKey_Stress tests key derivation with extreme lengths,
// empty inputs, and invalid formats.
func TestAdversarial_DeriveKey_Stress(t *testing.T) {
	script := `
	buat salt = kripto.acak_bita(16)

	// Valid varying lengths
	buat k16 = crypto.derive_key("password", salt, 16, "hex")
	buat k32 = crypto.derive_key("password", salt, 32, "hex")
	buat k64 = crypto.derive_key("password", salt, 64, "hex")

	// Empty password error
	buat errPass = ""
	coba {
		kripto.turunkan_kunci("", salt, 32)
	} tangkap (e) {
		errPass = e
	}

	// Empty salt error
	buat errSalt = ""
	coba {
		kripto.turunkan_kunci("password", "", 32)
	} tangkap (e) {
		errSalt = e
	}

	// Negative length error
	buat errNeg = ""
	coba {
		kripto.turunkan_kunci("password", salt, -5)
	} tangkap (e) {
		errNeg = e
	}

	// Length exceeding 1MB safe cap
	buat errHuge = ""
	coba {
		kripto.turunkan_kunci("password", salt, 2000000)
	} tangkap (e) {
		errHuge = e
	}

	[panjang(k16), panjang(k32), panjang(k64), errPass, errSalt, errNeg, errHuge]
	`
	res := runCryptoScript(script)
	arr, ok := res.(*evaluator.Array)
	if !ok {
		t.Fatalf("derive_key stress failed: %v", res)
	}

	if arr.Elements[0].(*evaluator.Integer).Value != 32 { // 16 bytes = 32 hex chars
		t.Errorf("k16 hex length expected 32, got %v", arr.Elements[0])
	}
	if arr.Elements[1].(*evaluator.Integer).Value != 64 { // 32 bytes = 64 hex chars
		t.Errorf("k32 hex length expected 64, got %v", arr.Elements[1])
	}
	if arr.Elements[2].(*evaluator.Integer).Value != 128 { // 64 bytes = 128 hex chars
		t.Errorf("k64 hex length expected 128, got %v", arr.Elements[2])
	}

	errPass := arr.Elements[3].Inspect()
	if !strings.Contains(errPass, "kata sandi tidak boleh kosong") {
		t.Errorf("empty password error unexpected: %s", errPass)
	}

	errSalt := arr.Elements[4].Inspect()
	if !strings.Contains(errSalt, "salt tidak boleh kosong") {
		t.Errorf("empty salt error unexpected: %s", errSalt)
	}

	errNeg := arr.Elements[5].Inspect()
	if !strings.Contains(errNeg, "harus lebih besar dari 0") {
		t.Errorf("negative length error unexpected: %s", errNeg)
	}

	errHuge := arr.Elements[6].Inspect()
	if !strings.Contains(errHuge, "melebihi batas aman") {
		t.Errorf("huge length error unexpected: %s", errHuge)
	}
}

// TestAdversarial_AESGCM_NonceUniqueness_Stress tests 1,000 consecutive AES-GCM
// encryptions of the same plaintext with the same key, ensuring that all 12-byte
// nonces are uniquely generated without a single collision.
func TestAdversarial_AESGCM_NonceUniqueness_Stress(t *testing.T) {
	key := runCryptoScript(`kripto.buat_kunci()`).Inspect()
	const count = 1000
	nonces := make(map[string]bool, count)

	for i := 0; i < count; i++ {
		script := fmt.Sprintf(`kripto.enkripsi("Data rahasia nonce stress", "%s")`, key)
		res := runCryptoScript(script)
		ctB64, ok := res.(*evaluator.String)
		if !ok || ctB64.Value == "" {
			t.Fatalf("enkripsi failed at iteration %d: %v", i, res)
		}

		payload, err := base64.StdEncoding.DecodeString(ctB64.Value)
		if err != nil {
			t.Fatalf("base64 decode failed at iteration %d: %v", i, err)
		}
		if len(payload) < 28 {
			t.Fatalf("payload too short at iteration %d: %d", i, len(payload))
		}

		nonce := hex.EncodeToString(payload[:12])
		if nonces[nonce] {
			t.Fatalf("CRITICAL: Nonce collision detected at iteration %d! Nonce: %s", i, nonce)
		}
		nonces[nonce] = true
	}

	if len(nonces) != count {
		t.Errorf("expected %d unique nonces, got %d", count, len(nonces))
	}
	t.Logf("Successfully verified %d unique 12-byte nonces (0 collisions)", count)
}

// TestAdversarial_AESGCM_PayloadBoundarySizes verifies round-trip encryption
// and decryption across exact payload boundary conditions:
// 0B, 1B, 15B, 16B (AES block), 17B, 31B, 32B, 33B, 64B, 1024B, 65536B (64KB).
func TestAdversarial_AESGCM_PayloadBoundarySizes(t *testing.T) {
	key := runCryptoScript(`kripto.buat_kunci()`).Inspect()

	sizes := []int{0, 1, 15, 16, 17, 31, 32, 33, 63, 64, 65, 127, 128, 256, 1024, 8192, 65536}

	for _, size := range sizes {
		t.Run(fmt.Sprintf("Size_%d_Bytes", size), func(t *testing.T) {
			raw := make([]byte, size)
			for i := range raw {
				raw[i] = byte((i * 31) % 256)
			}

			// Test via Bita object
			elems := make([]string, size)
			for i, b := range raw {
				elems[i] = fmt.Sprintf("%d", b)
			}
			var dataExpr string
			if size == 0 {
				dataExpr = `bita([])`
			} else {
				dataExpr = fmt.Sprintf(`bita([%s])`, strings.Join(elems, ","))
			}

			script := fmt.Sprintf(`
			buat data = %s
			buat enc = kripto.enkripsi(data, "%s")
			buat dec = kripto.dekripsi(enc, "%s")
			[panjang(enc) > 0, dec]
			`, dataExpr, key, key)

			res := runCryptoScript(script)
			arr, ok := res.(*evaluator.Array)
			if !ok {
				t.Fatalf("script execution failed for size %d: %v", size, res)
			}
			if arr.Elements[0] != evaluator.TRUE {
				t.Errorf("size %d encryption returned empty string", size)
			}
			decStr := arr.Elements[1].Inspect()
			if decStr != string(raw) {
				t.Errorf("size %d decryption mismatch: expected %d bytes, got %d bytes", size, len(raw), len(decStr))
			}
		})
	}
}

// TestAdversarial_AESGCM_BitFlipping_AllBoundaryBits tests single-bit and double-bit
// tampering explicitly on the exact boundary transitions:
// - byte 0 (nonce start)
// - byte 11 (nonce end)
// - byte 12 (data start / nonce boundary)
// - byte 11+N (data end)
// - byte 12+N (tag start / data boundary)
// - byte 27+N (tag end)
func TestAdversarial_AESGCM_BitFlipping_AllBoundaryBits(t *testing.T) {
	key := runCryptoScript(`kripto.buat_kunci()`).Inspect()
	// 5-byte payload: total = 12 nonce + 5 data + 16 tag = 33 bytes
	// Nonce: 0..11
	// Data: 12..16
	// Tag: 17..32
	msg := "Celang"
	scriptEnc := fmt.Sprintf(`kripto.enkripsi("%s", "%s")`, msg, key)
	ctB64 := runCryptoScript(scriptEnc).Inspect()
	payload, err := base64.StdEncoding.DecodeString(ctB64)
	if err != nil {
		t.Fatalf("decode base64 failed: %v", err)
	}

	criticalBytes := []struct {
		name    string
		bytePos int
	}{
		{"NonceStart", 0},
		{"NonceEnd", 11},
		{"DataStart", 12},
		{"DataEnd", 11 + len(msg)},
		{"TagStart", 12 + len(msg)},
		{"TagEnd", len(payload) - 1},
	}

	for _, cb := range criticalBytes {
		t.Run(cb.name, func(t *testing.T) {
			for bit := uint(0); bit < 8; bit++ {
				corrupted := make([]byte, len(payload))
				copy(corrupted, payload)
				corrupted[cb.bytePos] ^= (1 << bit)

				corruptedB64 := base64.StdEncoding.EncodeToString(corrupted)
				script := fmt.Sprintf(`
				buat err = ""
				coba {
					kripto.dekripsi("%s", "%s")
				} tangkap (e) {
					err = e
				}
				err
				`, corruptedB64, key)

				res := runCryptoScript(script).Inspect()
				if !strings.Contains(res, "dekripsi gagal") {
					t.Errorf("%s bit %d flip failed: expected empathetic error, got %s", cb.name, bit, res)
				}
			}
		})
	}

	// Cross-boundary dual bit flips
	t.Run("CrossBoundaryDualFlips", func(t *testing.T) {
		boundaryPairs := [][2]int{
			{11, 12},                        // Nonce/Data boundary
			{11 + len(msg), 12 + len(msg)}, // Data/Tag boundary
			{0, len(payload) - 1},          // Nonce start and Tag end
		}

		for _, pair := range boundaryPairs {
			corrupted := make([]byte, len(payload))
			copy(corrupted, payload)
			corrupted[pair[0]] ^= 0x01
			corrupted[pair[1]] ^= 0x80

			corruptedB64 := base64.StdEncoding.EncodeToString(corrupted)
			script := fmt.Sprintf(`
			buat err = ""
			coba {
				kripto.dekripsi("%s", "%s")
			} tangkap (e) {
				err = e
			}
			err
			`, corruptedB64, key)

			res := runCryptoScript(script).Inspect()
			if !strings.Contains(res, "dekripsi gagal") {
				t.Errorf("dual boundary flip (%d, %d) failed: expected empathetic error, got %s", pair[0], pair[1], res)
			}
		}
	})
}

// TestAdversarial_PBKDF2_IterationCaps_Boundary tests the exact boundaries of
// iteration count enforcement:
// - i=0 -> FALSE
// - i=-1 -> FALSE
// - i=10_000_001 (1 over 10M cap) -> FALSE immediately
// - i=50_000_000 -> FALSE immediately
func TestAdversarial_PBKDF2_IterationCaps_Boundary(t *testing.T) {
	saltB64 := base64.RawStdEncoding.EncodeToString([]byte("1234567890123456"))
	hashB64 := base64.RawStdEncoding.EncodeToString(make([]byte, 32))

	caps := []struct {
		name string
		iter string
	}{
		{"Zero", "0"},
		{"NegativeOne", "-1"},
		{"NegativeTenThousand", "-10000"},
		{"CapPlusOne", "10000001"},
		{"FifteenMillion", "15000000"},
		{"OneHundredMillion", "100000000"},
		{"MaxInt64Overflow", "9223372036854775808"},
	}

	for _, c := range caps {
		t.Run(c.name, func(t *testing.T) {
			phc := fmt.Sprintf("$pbkdf2-sha256$i=%s$%s$%s", c.iter, saltB64, hashB64)
			script := fmt.Sprintf(`kripto.verifikasi_sandi("mypassword", "%s")`, phc)
			res := runCryptoScript(script)
			if res != evaluator.FALSE {
				t.Errorf("iteration %s did not return FALSE: got %v", c.iter, res)
			}
		})
	}
}

// TestAdversarial_PBKDF2_TimingSafe_VariedPasswordLengths tests verify_password
// against candidate passwords of widely varying lengths (0 to 4096 chars)
// against a target hash, verifying constant-time compare integrity, absence of
// panics, and safe FALSE return.
func TestAdversarial_PBKDF2_TimingSafe_VariedPasswordLengths(t *testing.T) {
	targetPw := "StandardUserPassword123!"
	targetPhc := runCryptoScript(fmt.Sprintf(`kripto.hash_sandi("%s")`, targetPw)).Inspect()

	lengths := []int{0, 1, 5, 10, 23, 24, 25, 32, 64, 128, 256, 512, 1024, 4096}

	for _, length := range lengths {
		t.Run(fmt.Sprintf("Len_%d", length), func(t *testing.T) {
			candidate := strings.Repeat("x", length)
			script := fmt.Sprintf(`
			buat r1 = kripto.verifikasi_sandi("%s", "%s")
			buat r2 = crypto.verify_password("%s", "%s")
			[r1, r2]
			`, candidate, targetPhc, candidate, targetPhc)

			res := runCryptoScript(script)
			arr, ok := res.(*evaluator.Array)
			if !ok {
				t.Fatalf("verification script failed for length %d: %v", length, res)
			}
			if arr.Elements[0] != evaluator.FALSE || arr.Elements[1] != evaluator.FALSE {
				t.Errorf("candidate of length %d wrongly verified as TRUE", length)
			}
		})
	}

	// Verify exact match still succeeds
	t.Run("ExactMatch_True", func(t *testing.T) {
		script := fmt.Sprintf(`kripto.verifikasi_sandi("%s", "%s")`, targetPw, targetPhc)
		res := runCryptoScript(script)
		if res != evaluator.TRUE {
			t.Errorf("exact match failed to return TRUE: got %v", res)
		}
	})
}

// TestAdversarial_PBKDF2_SaltLengthBoundaries tests salt lengths:
// 0B, 1B, 15B (all < 16, must return FALSE), and 16B, 32B (>= 16, valid).
func TestAdversarial_PBKDF2_SaltLengthBoundaries(t *testing.T) {
	pw := "SecretPassword"
	dummyHash := base64.RawStdEncoding.EncodeToString(make([]byte, 32))

	// Salts below 16 bytes: must return FALSE
	for saltLen := 0; saltLen < 16; saltLen++ {
		t.Run(fmt.Sprintf("SaltBelow16_Len%d", saltLen), func(t *testing.T) {
			salt := make([]byte, saltLen)
			saltB64 := base64.RawStdEncoding.EncodeToString(salt)
			phc := fmt.Sprintf("$pbkdf2-sha256$i=1000$%s$%s", saltB64, dummyHash)

			script := fmt.Sprintf(`kripto.verifikasi_sandi("%s", "%s")`, pw, phc)
			res := runCryptoScript(script)
			if res != evaluator.FALSE {
				t.Errorf("salt length %d expected FALSE, got %v", saltLen, res)
			}
		})
	}

	// Salts 16 bytes and 32 bytes: compute actual valid PBKDF2 hash and verify match
	for _, saltLen := range []int{16, 32} {
		t.Run(fmt.Sprintf("ValidSalt_Len%d", saltLen), func(t *testing.T) {
			salt := make([]byte, saltLen)
			for i := range salt {
				salt[i] = byte(i + 42)
			}
			derived := evaluator.PBKDF2SHA256([]byte(pw), salt, 1000, 32)
			saltB64 := base64.RawStdEncoding.EncodeToString(salt)
			hashB64 := base64.RawStdEncoding.EncodeToString(derived)
			phc := fmt.Sprintf("$pbkdf2-sha256$i=1000$%s$%s", saltB64, hashB64)

			script := fmt.Sprintf(`kripto.verifikasi_sandi("%s", "%s")`, pw, phc)
			res := runCryptoScript(script)
			if res != evaluator.TRUE {
				t.Errorf("valid salt length %d failed to verify: got %v", saltLen, res)
			}
		})
	}
}



