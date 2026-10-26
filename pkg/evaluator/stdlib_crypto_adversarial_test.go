package evaluator_test

import (
	"celang/pkg/evaluator"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
	"testing"
)

// ============================================================================
// ADVERSARIAL TEST SUITE: R3 (Ed25519 Digital Signatures - RFC 8032)
// ============================================================================

// TestAdversarial_Ed25519_RFC8032_Comprehensive tests RFC 8032 Test Vectors 1, 2, 3
// with both seed and expanded 64-byte private keys, as well as multiple encodings (hex, base64, raw, bita).
func TestAdversarial_Ed25519_RFC8032_Comprehensive(t *testing.T) {
	type rfcVector struct {
		name       string
		seedHex    string
		pubHex     string
		msg        string
		isByteMsg  bool
		msgBytes   []byte
		sigHex     string
	}

	vectors := []rfcVector{
		{
			name:      "RFC 8032 Test 1 (empty message)",
			seedHex:   "9d61b19deffd5a60ba844af492ec2cc44449c5697b326919703bac031cae7f60",
			pubHex:    "d75a980182b10ab7d54bfed3c964073a0ee172f3daa62325af021a68f707511a",
			msg:       "",
			isByteMsg: false,
			sigHex:    "e5564300c360ac729086e2cc806e828a84877f1eb8e5d974d873e065224901555fb8821590a33bacc61e39701cf9b46bd25bf5f0595bbe24655141438e7a100b",
		},
		{
			name:      "RFC 8032 Test 2 (1-byte message 'r')",
			seedHex:   "4ccd089b28ff96da9db6c346ec114e0f5b8a319f35aba624da8cf6ed4fb8a6fb",
			pubHex:    "3d4017c3e843895a92b70aa74d1b7ebc9c982ccf2ec4968cc0cd55f12af4660c",
			msg:       "r",
			isByteMsg: false,
			sigHex:    "92a009a9f0d4cab8720e820b5f642540a2b27b5416503f8fb3762223ebdb69da085ac1e43e15996e458f3613d0f11d8c387b2eaeb4302aeeb00d291612bb0c00",
		},
		{
			name:      "RFC 8032 Test 3 (2-byte message 0xaf 0x82)",
			seedHex:   "c5aa8df43f9f837bedb7442f31dcb7b166d38535076f094b85ce3a2e0b4458f7",
			pubHex:    "fc51cd8e6218a1a38da47ed00230f0580816ed13ba3303ac5deb911548908025",
			isByteMsg: true,
			msgBytes:  []byte{0xaf, 0x82},
			sigHex:    "6291d657deec24024827e69c3abe01a30ce548a284743a445e3680d7db5ac3ac18ff9b538d16f290ae67f760984dc6594a7c15e9716ed28dc027beceea1ec40a",
		},
	}

	for _, vec := range vectors {
		t.Run(vec.name, func(t *testing.T) {
			// Expanded 64-byte private key is seed (32B) + pubkey (32B)
			seedBytes, _ := hex.DecodeString(vec.seedHex)
			pubBytes, _ := hex.DecodeString(vec.pubHex)
			expandedPrivBytes := append(append([]byte{}, seedBytes...), pubBytes...)
			expandedPrivHex := hex.EncodeToString(expandedPrivBytes)

			// 1. Sign using 32-byte hex seed
			var scriptSignSeed string
			if vec.isByteMsg {
				scriptSignSeed = fmt.Sprintf(`
					buat data = bita([%d, %d])
					buat sig = kripto.tanda_tangan(data, "%s")
					buat ver = kripto.verifikasi_tanda_tangan(data, sig, "%s")
					[sig, ver]
				`, vec.msgBytes[0], vec.msgBytes[1], vec.seedHex, vec.pubHex)
			} else {
				scriptSignSeed = fmt.Sprintf(`
					buat data = "%s"
					buat sig = kripto.tanda_tangan(data, "%s")
					buat ver = kripto.verifikasi_tanda_tangan(data, sig, "%s")
					[sig, ver]
				`, vec.msg, vec.seedHex, vec.pubHex)
			}

			resSeed := runCryptoScript(scriptSignSeed)
			arrSeed, ok := resSeed.(*evaluator.Array)
			if !ok {
				t.Fatalf("%s (seed) execution failed: %v", vec.name, resSeed)
			}
			if arrSeed.Elements[0].Inspect() != vec.sigHex {
				t.Errorf("%s signature mismatch with seed:\nexpected: %s\ngot:      %s", vec.name, vec.sigHex, arrSeed.Elements[0].Inspect())
			}
			if arrSeed.Elements[1] != evaluator.TRUE {
				t.Errorf("%s verification failed with seed", vec.name)
			}

			// 2. Sign using 64-byte expanded private key
			var scriptSignExpanded string
			if vec.isByteMsg {
				scriptSignExpanded = fmt.Sprintf(`
					buat data = bita([%d, %d])
					buat sig = crypto.sign(data, "%s")
					buat ver = crypto.verify(data, sig, "%s")
					[sig, ver]
				`, vec.msgBytes[0], vec.msgBytes[1], expandedPrivHex, vec.pubHex)
			} else {
				scriptSignExpanded = fmt.Sprintf(`
					buat data = "%s"
					buat sig = crypto.sign(data, "%s")
					buat ver = crypto.verify(data, sig, "%s")
					[sig, ver]
				`, vec.msg, expandedPrivHex, vec.pubHex)
			}

			resExp := runCryptoScript(scriptSignExpanded)
			arrExp, ok := resExp.(*evaluator.Array)
			if !ok {
				t.Fatalf("%s (expanded) execution failed: %v", vec.name, resExp)
			}
			if arrExp.Elements[0].Inspect() != vec.sigHex {
				t.Errorf("%s signature mismatch with expanded key:\nexpected: %s\ngot:      %s", vec.name, vec.sigHex, arrExp.Elements[0].Inspect())
			}
			if arrExp.Elements[1] != evaluator.TRUE {
				t.Errorf("%s verification failed with expanded key", vec.name)
			}

			// 3. Verify with Base64-encoded signature and Base64 public key
			sigBytes, _ := hex.DecodeString(vec.sigHex)
			sigB64 := base64.StdEncoding.EncodeToString(sigBytes)
			pubB64 := base64.StdEncoding.EncodeToString(pubBytes)

			var scriptB64 string
			if vec.isByteMsg {
				scriptB64 = fmt.Sprintf(`
					buat data = bita([%d, %d])
					crypto.verify(data, "%s", "%s")
				`, vec.msgBytes[0], vec.msgBytes[1], sigB64, pubB64)
			} else {
				scriptB64 = fmt.Sprintf(`
					crypto.verify("%s", "%s", "%s")
				`, vec.msg, sigB64, pubB64)
			}
			resB64 := runCryptoScript(scriptB64)
			if resB64 != evaluator.TRUE {
				t.Errorf("%s verification failed with Base64 signature and pubkey: %v", vec.name, resB64)
			}
		})
	}
}

// TestAdversarial_Ed25519_BitFlippingAndMalleability verifies that ANY bit flip
// across the 64-byte signature safely results in FALSE (never true, never panics).
func TestAdversarial_Ed25519_BitFlippingAndMalleability(t *testing.T) {
	// Standard Test 1 vector
	msg := "adversarial-signature-verification-payload"
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate ed25519 key: %v", err)
	}

	origSig := ed25519.Sign(priv, []byte(msg))
	pubHex := hex.EncodeToString(pub)

	// Critical byte offsets to mutate:
	// - 0 (start of R)
	// - 15 (middle of R)
	// - 31 (end of R)
	// - 32 (start of S)
	// - 47 (middle of S)
	// - 63 (end of S / high bits)
	targetOffsets := []int{0, 1, 15, 31, 32, 33, 47, 62, 63}

	for _, offset := range targetOffsets {
		for bit := 0; bit < 8; bit++ {
			mutatedSig := make([]byte, len(origSig))
			copy(mutatedSig, origSig)
			mutatedSig[offset] ^= 1 << bit

			mutatedHex := hex.EncodeToString(mutatedSig)

			script := fmt.Sprintf(`
				buat v1 = kripto.verifikasi_tanda_tangan("%s", "%s", "%s")
				buat v2 = crypto.verify("%s", "%s", "%s")
				[v1, v2]
			`, msg, mutatedHex, pubHex, msg, mutatedHex, pubHex)

			res := runCryptoScript(script)
			arr, ok := res.(*evaluator.Array)
			if !ok {
				t.Fatalf("eval failed on bit flip at offset %d bit %d: %v", offset, bit, res)
			}

			if arr.Elements[0] != evaluator.FALSE || arr.Elements[1] != evaluator.FALSE {
				t.Errorf("CRITICAL BUG: verification succeeded on bit-flipped signature! offset=%d, bit=%d, v1=%v, v2=%v",
					offset, bit, arr.Elements[0], arr.Elements[1])
			}
		}
	}
}

// TestAdversarial_Ed25519_Malleability_ScalarOrder tests non-canonical S values (S >= L).
// RFC 8032 section 5.1.7 requires that S < L; any S >= L must fail verification.
func TestAdversarial_Ed25519_Malleability_ScalarOrder(t *testing.T) {
	msg := "test-malleability-scalar"
	pub, priv, _ := ed25519.GenerateKey(rand.Reader)
	sig := ed25519.Sign(priv, []byte(msg))
	pubHex := hex.EncodeToString(pub)

	// L in little-endian bytes:
	// L = 2^252 + 27742317777372353535851937790883648493
	// L in bytes: 0xed, 0xd3, 0xf5, 0x5c, 0x1a, 0x63, 0x12, 0x58, 0xd6, 0x9c, 0xf7, 0xa2, 0xde, 0xf9, 0xde, 0x14, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10
	lBytes := []byte{
		0xed, 0xd3, 0xf5, 0x5c, 0x1a, 0x63, 0x12, 0x58,
		0xd6, 0x9c, 0xf7, 0xa2, 0xde, 0xf9, 0xde, 0x14,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10,
	}

	// 1. Set S to exactly L (S = L >= L, must be rejected)
	sigWithL := make([]byte, 64)
	copy(sigWithL[:32], sig[:32])
	copy(sigWithL[32:], lBytes)

	// 2. Set S with highest bit set (0x80 on byte 63 -> definitely S > L)
	sigWithHighS := make([]byte, 64)
	copy(sigWithHighS, sig)
	sigWithHighS[63] |= 0x80

	// 3. Set S to 0xFF * 32
	sigWithMaxS := make([]byte, 64)
	copy(sigWithMaxS[:32], sig[:32])
	for i := 32; i < 64; i++ {
		sigWithMaxS[i] = 0xff
	}

	testCases := []struct {
		name string
		sig  []byte
	}{
		{"S = L", sigWithL},
		{"S high bit set", sigWithHighS},
		{"S = 0xFF..FF", sigWithMaxS},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sHex := hex.EncodeToString(tc.sig)
			script := fmt.Sprintf(`kripto.verifikasi_tanda_tangan("%s", "%s", "%s")`, msg, sHex, pubHex)
			res := runCryptoScript(script)
			if res != evaluator.FALSE {
				t.Errorf("malleable signature with %s must return FALSE, got: %v", tc.name, res)
			}
		})
	}
}

// TestAdversarial_Ed25519_CorruptedInputsSafeFailure ensures verification safely returns FALSE
// without panicking on corrupted signatures, malformed hex, invalid sizes, or bad public keys.
func TestAdversarial_Ed25519_CorruptedInputsSafeFailure(t *testing.T) {
	validMsg := "valid message to verify"
	pub, priv, _ := ed25519.GenerateKey(rand.Reader)
	validSigHex := hex.EncodeToString(ed25519.Sign(priv, []byte(validMsg)))
	validPubHex := hex.EncodeToString(pub)

	corruptedSignatures := []struct {
		name string
		sig  string
	}{
		{"Empty string signature", `""`},
		{"1 byte hex (odd length)", `"a"`},
		{"Odd length 127 hex characters", `"` + strings.Repeat("a", 127) + `"`},
		{"Odd length 129 hex characters", `"` + strings.Repeat("a", 129) + `"`},
		{"Truncated 32 bytes hex (64 chars)", `"` + strings.Repeat("0", 64) + `"`},
		{"Truncated 63 bytes hex (126 chars)", `"` + strings.Repeat("0", 126) + `"`},
		{"Oversized 65 bytes hex (130 chars)", `"` + strings.Repeat("0", 130) + `"`},
		{"Non-hex characters in 128 chars", `"` + strings.Repeat("z", 128) + `"`},
		{"Invalid hex char inside valid length", `"` + validSigHex[:120] + "xyz!@#" + `"`},
		{"Garbage Base64 padding", `"not-a-valid-base64-signature====="`},
		{"Whitespace signature", `"   "`},
	}

	for _, tc := range corruptedSignatures {
		t.Run("Sig_"+tc.name, func(t *testing.T) {
			script := fmt.Sprintf(`
				buat res = kripto.verifikasi_tanda_tangan("%s", %s, "%s")
				res
			`, validMsg, tc.sig, validPubHex)

			res := runCryptoScript(script)
			if res != evaluator.FALSE {
				t.Errorf("corrupted signature '%s' did not return FALSE safely, got: %v", tc.name, res)
			}
		})
	}

	corruptedPubKeys := []struct {
		name string
		pub  string
	}{
		{"Empty pubkey", `""`},
		{"Odd length 63 hex chars", `"` + strings.Repeat("f", 63) + `"`},
		{"Odd length 65 hex chars", `"` + strings.Repeat("f", 65) + `"`},
		{"Truncated 16 bytes (32 hex chars)", `"` + strings.Repeat("1", 32) + `"`},
		{"Non-hex pubkey 64 chars", `"` + strings.Repeat("x", 64) + `"`},
		{"All zeroes pubkey 64 hex chars (0x00*32)", `"` + strings.Repeat("0", 64) + `"`},
		{"Off-curve point y >= 2^255 - 19 (e.g. 0xFF*32)", `"` + strings.Repeat("f", 64) + `"`},
		{"Whitespace pubkey", `"   "`},
	}

	for _, tc := range corruptedPubKeys {
		t.Run("Pub_"+tc.name, func(t *testing.T) {
			script := fmt.Sprintf(`
				buat res = kripto.verifikasi_tanda_tangan("%s", "%s", %s)
				res
			`, validMsg, validSigHex, tc.pub)

			res := runCryptoScript(script)
			if res != evaluator.FALSE {
				t.Errorf("corrupted pubkey '%s' did not return FALSE safely, got: %v", tc.name, res)
			}
		})
	}

	// Non-string / non-bita signature or pubkey arguments
	nonStringNonBita := []struct {
		name string
		val  string
	}{
		{"Integer", "12345"},
		{"Boolean", "benar"},
		{"Array", "[1, 2, 3]"},
		{"Map", `{"kunci": "nilai"}`},
		{"Nihil", "nihil"},
	}

	for _, tc := range nonStringNonBita {
		t.Run("NonStringType_Sig_"+tc.name, func(t *testing.T) {
			script := fmt.Sprintf(`kripto.verifikasi_tanda_tangan("%s", %s, "%s")`, validMsg, tc.val, validPubHex)
			res := runCryptoScript(script)
			if res != evaluator.FALSE {
				t.Errorf("non-string signature type '%s' must safely evaluate to FALSE, got: %v", tc.name, res)
			}
		})

		t.Run("NonStringType_Pub_"+tc.name, func(t *testing.T) {
			script := fmt.Sprintf(`kripto.verifikasi_tanda_tangan("%s", "%s", %s)`, validMsg, validSigHex, tc.val)
			res := runCryptoScript(script)
			if res != evaluator.FALSE {
				t.Errorf("non-string pubkey type '%s' must safely evaluate to FALSE, got: %v", tc.name, res)
			}
		})
	}
}

// TestAdversarial_Ed25519_CorruptedPrivateKeyInSign tests that passing invalid private keys
// to kripto.tanda_tangan / crypto.sign returns an empathetic Celang Error without panicking.
func TestAdversarial_Ed25519_CorruptedPrivateKeyInSign(t *testing.T) {
	badKeys := []struct {
		name string
		key  string
	}{
		{"Empty string", `""`},
		{"1 byte hex", `"aa"`},
		{"31 bytes hex (62 chars)", `"` + strings.Repeat("a", 62) + `"`},
		{"33 bytes hex (66 chars)", `"` + strings.Repeat("b", 66) + `"`},
		{"63 bytes hex (126 chars)", `"` + strings.Repeat("c", 126) + `"`},
		{"65 bytes hex (130 chars)", `"` + strings.Repeat("d", 130) + `"`},
		{"Integer type", "12345"},
		{"Boolean type", "salah"},
	}

	for _, tc := range badKeys {
		t.Run(tc.name, func(t *testing.T) {
			script := fmt.Sprintf(`
				buat errMsg = ""
				coba {
					kripto.tanda_tangan("halo", %s)
				} tangkap (e) {
					errMsg = e
				}
				errMsg
			`, tc.key)

			res := runCryptoScript(script)
			str, ok := res.(*evaluator.String)
			if !ok || str.Value == "" {
				t.Errorf("tanda_tangan with bad key '%s' should produce an error, got: %v", tc.name, res)
			}
		})
	}
}

// TestAdversarial_Ed25519_Corrupted64CharHexKey_Vulnerability empirically demonstrates
// that parseEd25519PrivateKey silently accepts a 64-character string containing non-hex characters
// as a raw 64-byte private key instead of rejecting it as a corrupted hex seed.
func TestAdversarial_Ed25519_Corrupted64CharHexKey_Vulnerability(t *testing.T) {
	// A standard 64-char hex seed from RFC 8032 Test 1 with a single typo (last char '0' -> 'z')
	corruptedSeed := "9d61b19deffd5a60ba844af492ec2cc44449c5697b326919703bac031cae7f6z"
	pubHex := "d75a980182b10ab7d54bfed3c964073a0ee172f3daa62325af021a68f707511a"

	// When a user makes a typo in their 64-char hex seed:
	scriptSign := fmt.Sprintf(`kripto.tanda_tangan("message", "%s")`, corruptedSeed)
	resSign := runCryptoScript(scriptSign)

	// Observation: Instead of returning an Error object, Celang returns a String signature!
	sigStr, isString := resSign.(*evaluator.String)
	if isString {
		t.Logf("[CONFIRMED VULNERABILITY] kripto.tanda_tangan silently accepted corrupted 64-char hex key and generated signature: %s", sigStr.Value)

		// Now attempt to verify this signature with the genuine public key:
		scriptVerify := fmt.Sprintf(`kripto.verifikasi_tanda_tangan("message", "%s", "%s")`, sigStr.Value, pubHex)
		resVerify := runCryptoScript(scriptVerify)
		if resVerify == evaluator.FALSE {
			t.Logf("[CONSEQUENCE] Verification with intended public key returned FALSE: the signature is permanently invalid/unverifiable.")
		}
	} else if errObj, isErr := resSign.(*evaluator.Error); isErr {
		t.Logf("kripto.tanda_tangan properly rejected corrupted 64-char key: %s", errObj.Message)
	}
}

// TestAdversarial_Ed25519_DifferentialOracle compares Celang sign and verify against Go's standard library.
func TestAdversarial_Ed25519_DifferentialOracle(t *testing.T) {
	for i := 0; i < 10; i++ {
		pub, priv, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			t.Fatalf("failed keygen: %v", err)
		}
		seedHex := hex.EncodeToString(priv.Seed())
		pubHex := hex.EncodeToString(pub)

		msg := fmt.Sprintf("random-differential-message-#%d", i)
		goSig := ed25519.Sign(priv, []byte(msg))
		goSigHex := hex.EncodeToString(goSig)

		// 1. Sign in Celang and verify with Go oracle
		scriptSign := fmt.Sprintf(`crypto.sign("%s", "%s")`, msg, seedHex)
		celangSigObj := runCryptoScript(scriptSign)
		celangSigStr, ok := celangSigObj.(*evaluator.String)
		if !ok {
			t.Fatalf("Celang sign failed: %v", celangSigObj)
		}
		if celangSigStr.Value != goSigHex {
			t.Fatalf("Differential mismatch in sign! Iteration %d\nGo:     %s\nCelang: %s", i, goSigHex, celangSigStr.Value)
		}

		// 2. Verify Go signature in Celang
		scriptVerify := fmt.Sprintf(`crypto.verify("%s", "%s", "%s")`, msg, goSigHex, pubHex)
		resVerify := runCryptoScript(scriptVerify)
		if resVerify != evaluator.TRUE {
			t.Fatalf("Celang failed to verify Go signature in iteration %d: %v", i, resVerify)
		}
	}
}

// ============================================================================
// ADVERSARIAL TEST SUITE: R4 (Bilingual Parity & Universal Mode Enforcement)
// ============================================================================

// TestAdversarial_Bilingual_FullParity verifies that all cryptographic functions
// are present and functional under BOTH 'kripto' and 'crypto' namespaces in normal mode.
func TestAdversarial_Bilingual_FullParity(t *testing.T) {
	script := `
	// 1. Key generation parity
	buat k1 = kripto.buat_kunci()
	buat k2 = crypto.generate_key()
	buat p1 = kripto.buat_pasangan_kunci()
	buat p2 = crypto.generate_keypair()

	// Check map keys
	buat p1_keys_ok = (p1.kunci_publik != nihil) dan (p1.kunci_privat != nihil) dan (p1.public_key != nihil) dan (p1.private_key != nihil)
	buat p2_keys_ok = (p2.kunci_publik != nihil) dan (p2.kunci_privat != nihil) dan (p2.public_key != nihil) dan (p2.private_key != nihil)

	// 2. Symmetric encryption cross-talk
	buat enc_id = kripto.enkripsi("pesan rahasia", k1)
	buat dec_en = crypto.decrypt(enc_id, k1)

	buat enc_en = crypto.encrypt("secret message", k2)
	buat dec_id = kripto.dekripsi(enc_en, k2)

	// 3. Password hashing cross-talk
	buat pw = "sandi-rahasia-123"
	buat h_id = kripto.hash_sandi(pw)
	buat v_en = crypto.verify_password(pw, h_id)

	buat h_en = crypto.hash_password(pw)
	buat v_id = kripto.verifikasi_sandi(pw, h_en)

	// 4. Ed25519 signing cross-talk
	buat data = "data transaksi celang"
	buat sig_id = kripto.tanda_tangan(data, p1.kunci_privat)
	buat v_sig_en = crypto.verify(data, sig_id, p1.public_key)

	buat sig_en = crypto.sign(data, p2.private_key)
	buat v_sig_id = kripto.verifikasi_tanda_tangan(data, sig_en, p2.kunci_publik)

	// 5. Utility functions
	buat h_sha256_id = kripto.sha256("test")
	buat h_sha256_en = crypto.sha256("test")
	buat sha256_match = (h_sha256_id == h_sha256_en)

	buat rhex_id = kripto.acak_heks(16)
	buat rhex_en = crypto.random_hex(16)
	buat rhex_len_ok = (panjang(rhex_id) == 32) dan (panjang(rhex_en) == 32)

	buat uuid_id = kripto.uuid()
	buat uuid_en = crypto.uuid_v4()
	buat uuid_len_ok = (panjang(uuid_id) == 36) dan (panjang(uuid_en) == 36)

	buat t_cmp_id = kripto.konstanta_sama("abc", "abc")
	buat t_cmp_en = crypto.timing_safe_equal("abc", "abc")

	[
		p1_keys_ok,
		p2_keys_ok,
		dec_en == "pesan rahasia",
		dec_id == "secret message",
		v_en,
		v_id,
		v_sig_en,
		v_sig_id,
		sha256_match,
		rhex_len_ok,
		uuid_len_ok,
		t_cmp_id dan t_cmp_en
	]
	`

	res := runCryptoScript(script)
	arr, ok := res.(*evaluator.Array)
	if !ok {
		t.Fatalf("bilingual parity script failed: %v", res)
	}

	for i, elem := range arr.Elements {
		if elem != evaluator.TRUE {
			t.Errorf("bilingual parity assertion [%d] failed: got %v", i, elem)
		}
	}
}

// TestAdversarial_UniversalMode_StrictRejection tests that "use universal" and 'use universal'
// strictly reject any reference to 'kripto' with an empathetic suggestion, while allowing 'crypto'.
func TestAdversarial_UniversalMode_StrictRejection(t *testing.T) {
	testCases := []struct {
		name          string
		code          string
		expectedError string
		expectedHint  string
	}{
		{
			name: "kripto identifier access in double-quoted directive",
			code: `"use universal"
let x = kripto.buat_kunci()
`,
			expectedError: "Indonesian module 'kripto' is not allowed in Universal Mode",
			expectedHint:  "Use 'crypto' instead.",
		},
		{
			name: "kripto identifier access in single-quoted directive",
			code: `'use universal'
let x = kripto.generate_key()
`,
			expectedError: "Indonesian module 'kripto' is not allowed in Universal Mode",
			expectedHint:  "Use 'crypto' instead.",
		},
		{
			name: "kripto import in double-quoted directive",
			code: `"use universal"
import "kripto"
`,
			expectedError: "Indonesian module 'kripto' is not allowed in Universal Mode",
			expectedHint:  "Import 'crypto' instead.",
		},
		{
			name: "kripto named import in single-quoted directive",
			code: `'use universal'
import { generate_key } from "kripto"
`,
			expectedError: "Indonesian module 'kripto' is not allowed in Universal Mode",
			expectedHint:  "Import 'crypto' instead.",
		},
		{
			name: "use english directive rejecting kripto",
			code: `"use english"
let x = kripto.sha256("test")
`,
			expectedError: "Indonesian module 'kripto' is not allowed in Universal Mode",
			expectedHint:  "Use 'crypto' instead.",
		},
		{
			name: "kripto with leading comments and whitespace",
			code: `// Leading header comment
// Second comment line
"use universal"
let x = kripto.uuid()
`,
			expectedError: "Indonesian module 'kripto' is not allowed in Universal Mode",
			expectedHint:  "Use 'crypto' instead.",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := runCryptoScript(tc.code)
			errObj, ok := res.(*evaluator.Error)
			if !ok {
				t.Fatalf("expected *evaluator.Error for '%s', got: %T (%v)", tc.name, res, res)
			}
			if !strings.Contains(errObj.Message, tc.expectedError) {
				t.Errorf("error message mismatch for '%s':\nexpected substr: %s\ngot:             %s",
					tc.name, tc.expectedError, errObj.Message)
			}
			if !strings.Contains(errObj.FriendlyHint, tc.expectedHint) {
				t.Errorf("hint mismatch for '%s':\nexpected substr: %s\ngot:             %s",
					tc.name, tc.expectedHint, errObj.FriendlyHint)
			}
		})
	}
}

// TestAdversarial_UniversalMode_CryptoAccepted verifies that all crypto features
// execute flawlessly under "use universal" without tripping any language filters.
func TestAdversarial_UniversalMode_CryptoAccepted(t *testing.T) {
	script := `"use universal"
import { generate_key, encrypt, decrypt, hash_password, verify_password, derive_key, generate_keypair, sign, verify, sha256, sha512, random_hex, random_bytes, uuid_v4, timing_safe_equal } from "crypto"

let k = generate_key()
let enc = encrypt("universal payload test", k)
let dec = decrypt(enc, k)

let pwHash = hash_password("strongPassword!@#")
let pwOk = verify_password("strongPassword!@#", pwHash)
let pwBad = verify_password("wrongPassword", pwHash)

let salt = "mysalt1234567890"
let derivedHex = derive_key("strongPassword!@#", salt, 32, "hex")

let pair = generate_keypair()
let sig = sign("contract v1.0", pair.private_key)
let sigOk = verify("contract v1.0", sig, pair.public_key)
let sigBad = verify("contract v1.1", sig, pair.public_key)

let hashHex = sha256("celang-universal")
let hexToken = random_hex(16)
let id = uuid_v4()
let eqOk = timing_safe_equal("token123", "token123")

[
    dec == "universal payload test",
    pwOk == true,
    pwBad == false,
    len(derivedHex) == 64,
    sigOk == true,
    sigBad == false,
    len(hashHex) == 64,
    len(hexToken) == 32,
    len(id) == 36,
    eqOk == true
]
`
	res := runCryptoScript(script)
	arr, ok := res.(*evaluator.Array)
	if !ok {
		t.Fatalf("universal mode crypto suite execution failed: %v", res)
	}

	for i, elem := range arr.Elements {
		if elem != evaluator.TRUE {
			t.Errorf("universal mode assertion [%d] failed: got %v", i, elem)
		}
	}
}

// TestAdversarial_Ed25519_RFC8032_Test1024 verifies the official 1023-byte message vector from RFC 8032 Section 7.1.
func TestAdversarial_Ed25519_RFC8032_Test1024(t *testing.T) {
	const rfc1024MsgHex = "08b8b2b733424243760fe426a4b54908" +
		"632110a66c2f6591eabd3345e3e4eb98" +
		"fa6e264bf09efe12ee50f8f54e9f77b1" +
		"e355f6c50544e23fb1433ddf73be84d8" +
		"79de7c0046dc4996d9e773f4bc9efe57" +
		"38829adb26c81b37c93a1b270b20329d" +
		"658675fc6ea534e0810a4432826bf58c" +
		"941efb65d57a338bbd2e26640f89ffbc" +
		"1a858efcb8550ee3a5e1998bd177e93a" +
		"7363c344fe6b199ee5d02e82d522c4fe" +
		"ba15452f80288a821a579116ec6dad2b" +
		"3b310da903401aa62100ab5d1a36553e" +
		"06203b33890cc9b832f79ef80560ccb9" +
		"a39ce767967ed628c6ad573cb116dbef" +
		"efd75499da96bd68a8a97b928a8bbc10" +
		"3b6621fcde2beca1231d206be6cd9ec7" +
		"aff6f6c94fcd7204ed3455c68c83f4a4" +
		"1da4af2b74ef5c53f1d8ac70bdcb7ed1" +
		"85ce81bd84359d44254d95629e9855a9" +
		"4a7c1958d1f8ada5d0532ed8a5aa3fb2" +
		"d17ba70eb6248e594e1a2297acbbb39d" +
		"502f1a8c6eb6f1ce22b3de1a1f40cc24" +
		"554119a831a9aad6079cad88425de6bd" +
		"e1a9187ebb6092cf67bf2b13fd65f270" +
		"88d78b7e883c8759d2c4f5c65adb7553" +
		"878ad575f9fad878e80a0c9ba63bcbcc" +
		"2732e69485bbc9c90bfbd62481d9089b" +
		"eccf80cfe2df16a2cf65bd92dd597b07" +
		"07e0917af48bbb75fed413d238f5555a" +
		"7a569d80c3414a8d0859dc65a46128ba" +
		"b27af87a71314f318c782b23ebfe808b" +
		"82b0ce26401d2e22f04d83d1255dc51a" +
		"ddd3b75a2b1ae0784504df543af8969b" +
		"e3ea7082ff7fc9888c144da2af58429e" +
		"c96031dbcad3dad9af0dcbaaaf268cb8" +
		"fcffead94f3c7ca495e056a9b47acdb7" +
		"51fb73e666c6c655ade8297297d07ad1" +
		"ba5e43f1bca32301651339e22904cc8c" +
		"42f58c30c04aafdb038dda0847dd988d" +
		"cda6f3bfd15c4b4c4525004aa06eeff8" +
		"ca61783aacec57fb3d1f92b0fe2fd1a8" +
		"5f6724517b65e614ad6808d6f6ee34df" +
		"f7310fdc82aebfd904b01e1dc54b2927" +
		"094b2db68d6f903b68401adebf5a7e08" +
		"d78ff4ef5d63653a65040cf9bfd4aca7" +
		"984a74d37145986780fc0b16ac451649" +
		"de6188a7dbdf191f64b5fc5e2ab47b57" +
		"f7f7276cd419c17a3ca8e1b939ae49e4" +
		"88acba6b965610b5480109c8b17b80e1" +
		"b7b750dfc7598d5d5011fd2dcc5600a3" +
		"2ef5b52a1ecc820e308aa342721aac09" +
		"43bf6686b64b2579376504ccc493d97e" +
		"6aed3fb0f9cd71a43dd497f01f17c0e2" +
		"cb3797aa2a2f256656168e6c496afc5f" +
		"b93246f6b1116398a346f1a641f3b041" +
		"e989f7914f90cc2c7fff357876e506b5" +
		"0d334ba77c225bc307ba537152f3f161" +
		"0e4eafe595f6d9d90d11faa933a15ef1" +
		"369546868a7f3a45a96768d40fd9d034" +
		"12c091c6315cf4fde7cb68606937380d" +
		"b2eaaa707b4c4185c32eddcdd306705e" +
		"4dc1ffc872eeee475a64dfac86aba41c" +
		"0618983f8741c5ef68d3a101e8a3b8ca" +
		"c60c905c15fc910840b94c00a0b9d0"

	seedHex := "f5e5767cf153319517630f226876b86c8160cc583bc013744c6bf255f5cc0ee5"
	pubHex := "278117fc144c72340f67d0f2316e8386ceffbf2b2428c9c51fef7c597f1d426e"
	expectedSig := "0aab4c900501b3e24d7cdf4663326a3a87df5e4843b2cbdb67cbf6e460fec350aa5371b1508f9f4528ecea23c436d94b5e8fcd4f681e30a6ac00a9704a188a03"

	msgBytes, err := hex.DecodeString(rfc1024MsgHex)
	if err != nil {
		t.Fatalf("failed to decode 1024 msg hex: %v", err)
	}
	if len(msgBytes) != 1023 {
		t.Fatalf("expected 1023 bytes, got %d", len(msgBytes))
	}

	// Format elements for Celang bita literal: bita([0x08, 0xb8, ...])
	var byteLiterals []string
	for _, b := range msgBytes {
		byteLiterals = append(byteLiterals, fmt.Sprintf("%d", b))
	}
	bitaArrayStr := strings.Join(byteLiterals, ", ")

	script := fmt.Sprintf(`
		buat data = bita([%s])
		buat sig = crypto.sign(data, "%s")
		buat ok = crypto.verify(data, sig, "%s")
		[sig, ok]
	`, bitaArrayStr, seedHex, pubHex)

	res := runCryptoScript(script)
	arr, ok := res.(*evaluator.Array)
	if !ok {
		t.Fatalf("RFC 8032 Test 1024 script failed: %v", res)
	}

	sigGot := arr.Elements[0].Inspect()
	if sigGot != expectedSig {
		t.Errorf("RFC 8032 Test 1024 signature mismatch:\nexpected: %s\ngot:      %s", expectedSig, sigGot)
	}
	if arr.Elements[1] != evaluator.TRUE {
		t.Errorf("RFC 8032 Test 1024 verification failed: got %v", arr.Elements[1])
	}
}

// TestAdversarial_Ed25519_UppercaseAndMixedHexSeed ensures hex seeds with uppercase
// letters are accepted and produce deterministic signatures identical to lowercase.
func TestAdversarial_Ed25519_UppercaseAndMixedHexSeed(t *testing.T) {
	seedLower := "9d61b19deffd5a60ba844af492ec2cc44449c5697b326919703bac031cae7f60"
	seedUpper := strings.ToUpper(seedLower)
	pubLower := "d75a980182b10ab7d54bfed3c964073a0ee172f3daa62325af021a68f707511a"
	pubUpper := strings.ToUpper(pubLower)

	script := fmt.Sprintf(`
		buat sig1 = kripto.tanda_tangan("pesan", "%s")
		buat sig2 = kripto.tanda_tangan("pesan", "%s")
		buat v1 = kripto.verifikasi_tanda_tangan("pesan", sig1, "%s")
		buat v2 = kripto.verifikasi_tanda_tangan("pesan", sig2, "%s")
		[sig1 == sig2, v1, v2]
	`, seedLower, seedUpper, pubLower, pubUpper)

	res := runCryptoScript(script)
	arr, ok := res.(*evaluator.Array)
	if !ok {
		t.Fatalf("uppercase seed test script failed: %v", res)
	}
	if arr.Elements[0] != evaluator.TRUE || arr.Elements[1] != evaluator.TRUE || arr.Elements[2] != evaluator.TRUE {
		t.Errorf("uppercase seed test failed: %v", arr)
	}
}

// TestAdversarial_Ed25519_EmptyBinaryMessage verifies that empty bita([]) works cleanly.
func TestAdversarial_Ed25519_EmptyBinaryMessage(t *testing.T) {
	seed := "9d61b19deffd5a60ba844af492ec2cc44449c5697b326919703bac031cae7f60"
	pub := "d75a980182b10ab7d54bfed3c964073a0ee172f3daa62325af021a68f707511a"

	script := fmt.Sprintf(`
		buat data = bita([])
		buat sig = crypto.sign(data, "%s")
		buat ok = crypto.verify(data, sig, "%s")
		[sig, ok]
	`, seed, pub)

	res := runCryptoScript(script)
	arr, ok := res.(*evaluator.Array)
	if !ok {
		t.Fatalf("empty bita test failed: %v", res)
	}
	// For 0-byte message, RFC 8032 Test 1 expected signature:
	expSig1 := "e5564300c360ac729086e2cc806e828a84877f1eb8e5d974d873e065224901555fb8821590a33bacc61e39701cf9b46bd25bf5f0595bbe24655141438e7a100b"
	if arr.Elements[0].Inspect() != expSig1 {
		t.Errorf("empty bita signature mismatch:\nexpected: %s\ngot:      %s", expSig1, arr.Elements[0].Inspect())
	}
	if arr.Elements[1] != evaluator.TRUE {
		t.Errorf("empty bita verify returned non-true: %v", arr.Elements[1])
	}
}

// TestAdversarial_UniversalMode_ComplexExpressions ensures that attempting to pass,
// nest, or alias 'kripto' in any syntactic position under Universal Mode triggers the empathetic error.
func TestAdversarial_UniversalMode_ComplexExpressions(t *testing.T) {
	cases := []struct {
		name string
		code string
	}{
		{
			name: "kripto in array literal",
			code: `"use universal"
let arr = [1, kripto, 3]`,
		},
		{
			name: "kripto in map value",
			code: `"use universal"
let m = {"module": kripto}`,
		},
		{
			name: "kripto passed to function",
			code: `"use universal"
def inspect(m) { return m }
inspect(kripto)`,
		},
		{
			name: "chained member call on kripto",
			code: `"use universal"
let k = kripto.buat_pasangan_kunci().kunci_publik`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res := runCryptoScript(tc.code)
			errObj, ok := res.(*evaluator.Error)
			if !ok {
				t.Fatalf("expected *evaluator.Error for '%s', got: %T (%v)", tc.name, res, res)
			}
			if !strings.Contains(errObj.Message, "Indonesian module 'kripto' is not allowed in Universal Mode") {
				t.Errorf("unexpected error message: %s", errObj.Message)
			}
			if !strings.Contains(errObj.FriendlyHint, "Use 'crypto' instead.") {
				t.Errorf("unexpected friendly hint: %s", errObj.FriendlyHint)
			}
		})
	}
}

