package evaluator_test

import (
	"celang/pkg/evaluator"
	"celang/pkg/parser"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"
)

func runCryptoScript(input string) evaluator.Object {
	l := parser.NewLexer(input)
	p := parser.NewParser(l)
	prog, err := p.ParseProgram()
	if err != nil {
		return &evaluator.Error{Message: err.Error()}
	}
	env := evaluator.NewEnvironment()
	env.IsUniversal = p.IsUniversal
	evaluator.RegisterBuiltins(env, nil)
	evaluator.RegisterStdlibGlobals(env)
	return evaluator.Eval(prog, env)
}

func TestCrypto_Hashes(t *testing.T) {
	script := `
	buat h_sha256_kosong = kripto.sha256("")
	buat h_sha256_teks = kripto.sha256("halo dunia")
	buat h_sha256_bita = kripto.sha256(bita("halo dunia"))

	buat h_sha512_kosong = crypto.sha512("")
	buat h_sha512_teks = crypto.sha512("halo dunia")

	buat h_sha1_kosong = kripto.sha1("")
	buat h_sha1_teks = kripto.sha1("celang")

	buat h_md5_kosong = crypto.md5("")
	buat h_md5_teks = crypto.md5("celang")

	[
		h_sha256_kosong,
		h_sha256_teks,
		h_sha256_bita,
		h_sha512_kosong,
		h_sha512_teks,
		h_sha1_kosong,
		h_sha1_teks,
		h_md5_kosong,
		h_md5_teks
	]
	`
	res := runCryptoScript(script)
	arr, ok := res.(*evaluator.Array)
	if !ok {
		t.Fatalf("evaluasi gagal: %v", res)
	}

	// 1. SHA-256 empty string NIST vector
	expSha256Empty := "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	if arr.Elements[0].Inspect() != expSha256Empty {
		t.Errorf("sha256 kosong salah: diharapkan %s, dapat %s", expSha256Empty, arr.Elements[0].Inspect())
	}

	// 2. SHA-256 string vs bita must match
	if arr.Elements[1].Inspect() != arr.Elements[2].Inspect() {
		t.Errorf("sha256 teks dan bita tidak cocok: %s vs %s", arr.Elements[1].Inspect(), arr.Elements[2].Inspect())
	}

	// 3. SHA-512 empty string NIST vector
	expSha512Empty := "cf83e1357eefb8bdf1542850d66d8007d620e4050b5715dc83f4a921d36ce9ce47d0d13c5d85f2b0ff8318d2877eec2f63b931bd47417a81a538327af927da3e"
	if arr.Elements[3].Inspect() != expSha512Empty {
		t.Errorf("sha512 kosong salah: diharapkan %s, dapat %s", expSha512Empty, arr.Elements[3].Inspect())
	}
	if len(arr.Elements[4].Inspect()) != 128 {
		t.Errorf("panjang sha512 salah: diharapkan 128 heksadesimal, dapat %d", len(arr.Elements[4].Inspect()))
	}

	// 4. SHA-1 empty string NIST vector
	expSha1Empty := "da39a3ee5e6b4b0d3255bfef95601890afd80709"
	if arr.Elements[5].Inspect() != expSha1Empty {
		t.Errorf("sha1 kosong salah: diharapkan %s, dapat %s", expSha1Empty, arr.Elements[5].Inspect())
	}

	// 5. MD5 empty string RFC vector
	expMD5Empty := "d41d8cd98f00b204e9800998ecf8427e"
	if arr.Elements[7].Inspect() != expMD5Empty {
		t.Errorf("md5 kosong salah: diharapkan %s, dapat %s", expMD5Empty, arr.Elements[7].Inspect())
	}
}

func TestCrypto_HMAC(t *testing.T) {
	script := `
	buat kunci = "rahasia-celang"
	buat pesan = "transaksi-12345"
	buat mac1 = kripto.hmac_sha256(kunci, pesan)
	buat mac2 = crypto.hmac(kunci, bita(pesan))

	[mac1, mac2, mac1 == mac2]
	`
	res := runCryptoScript(script)
	arr, ok := res.(*evaluator.Array)
	if !ok {
		t.Fatalf("evaluasi gagal: %v", res)
	}

	if arr.Elements[2] != evaluator.TRUE {
		t.Errorf("hmac_sha256 dengan teks dan bita harus sama: %v vs %v", arr.Elements[0], arr.Elements[1])
	}
	if len(arr.Elements[0].Inspect()) != 64 {
		t.Errorf("panjang hmac_sha256 salah: diharapkan 64 karakter heks, dapat %d", len(arr.Elements[0].Inspect()))
	}
}

func TestCrypto_RandomAndUUID(t *testing.T) {
	script := `
	buat hex16 = kripto.acak_heks(16)
	buat hex0 = kripto.acak_heks(0)
	buat bita10 = crypto.random_bytes(10)
	buat id1 = kripto.uuid()
	buat id2 = crypto.uuid()

	[hex16, hex0, panjang(bita10), id1, id2, id1 != id2]
	`
	res := runCryptoScript(script)
	arr, ok := res.(*evaluator.Array)
	if !ok {
		t.Fatalf("evaluasi gagal: %v", res)
	}

	hex16 := arr.Elements[0].Inspect()
	if len(hex16) != 32 {
		t.Errorf("acak_heks(16) harus menghasilkan 32 karakter heks, dapat %d (%s)", len(hex16), hex16)
	}

	hex0 := arr.Elements[1].Inspect()
	if hex0 != "" {
		t.Errorf("acak_heks(0) harus kosong, dapat %s", hex0)
	}

	intLen, ok := arr.Elements[2].(*evaluator.Integer)
	if !ok || intLen.Value != 10 {
		t.Errorf("random_bytes(10) harus menghasilkan bita sepanjang 10, dapat %v", arr.Elements[2])
	}

	id1 := arr.Elements[3].Inspect()
	id2 := arr.Elements[4].Inspect()
	uuidRegex := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
	if !uuidRegex.MatchString(id1) {
		t.Errorf("uuid() '%s' bukan format UUID v4 RFC 4122 yang valid", id1)
	}
	if !uuidRegex.MatchString(id2) {
		t.Errorf("uuid() '%s' bukan format UUID v4 RFC 4122 yang valid", id2)
	}

	if arr.Elements[5] != evaluator.TRUE {
		t.Errorf("dua panggilan uuid() berturut-turut menghasilkan nilai yang sama: %s == %s", id1, id2)
	}
}

func TestCrypto_ConstantTimeCompare(t *testing.T) {
	script := `
	buat token1 = "abcdef123456"
	buat token2 = "abcdef123456"
	buat token3 = "abcdef123457"

	buat sama = kripto.konstanta_sama(token1, token2)
	buat beda = crypto.timing_safe_equal(token1, token3)
	buat bitaSama = kripto.banding_aman(bita(token1), bita(token2))

	[sama, beda, bitaSama]
	`
	res := runCryptoScript(script)
	arr, ok := res.(*evaluator.Array)
	if !ok {
		t.Fatalf("evaluasi gagal: %v", res)
	}

	if arr.Elements[0] != evaluator.TRUE {
		t.Errorf("konstanta_sama harus bernilai benar untuk token yang sama")
	}
	if arr.Elements[1] != evaluator.FALSE {
		t.Errorf("timing_safe_equal harus bernilai salah untuk token yang berbeda")
	}
	if arr.Elements[2] != evaluator.TRUE {
		t.Errorf("banding_aman harus bernilai benar untuk bita yang sama")
	}
}

func TestCrypto_ModuleImportAndUniversal(t *testing.T) {
	// 1. Indonesian import
	scriptIndo := `
	impor { sha256, uuid } dari "kripto"
	buat h = sha256("tes")
	buat u = uuid()
	[panjang(h), panjang(u)]
	`
	resIndo := runCryptoScript(scriptIndo)
	arrIndo, ok := resIndo.(*evaluator.Array)
	if !ok {
		t.Fatalf("import 'kripto' gagal: %v", resIndo)
	}
	if arrIndo.Elements[0].(*evaluator.Integer).Value != 64 {
		t.Errorf("sha256 dari import salah: %v", arrIndo.Elements[0])
	}
	if arrIndo.Elements[1].(*evaluator.Integer).Value != 36 {
		t.Errorf("uuid dari import salah: %v", arrIndo.Elements[1])
	}

	// 2. Universal Mode English import
	scriptUniversal := `"use universal"
	import { sha256, random_hex } from "crypto"
	let h = sha256("tes")
	let r = random_hex(4)
	[len(h), len(r)]
	`
	resUniv := runCryptoScript(scriptUniversal)
	arrUniv, ok := resUniv.(*evaluator.Array)
	if !ok {
		t.Fatalf("universal import 'crypto' gagal: %v", resUniv)
	}
	if arrUniv.Elements[0].(*evaluator.Integer).Value != 64 {
		t.Errorf("universal sha256 salah: %v", arrUniv.Elements[0])
	}
	if arrUniv.Elements[1].(*evaluator.Integer).Value != 8 {
		t.Errorf("universal random_hex salah: %v", arrUniv.Elements[1])
	}

	// 3. Universal Mode rejecting 'kripto' in import
	scriptRejectImport := `"use universal"
	import "kripto"
	`
	resReject1 := runCryptoScript(scriptRejectImport)
	errObj1, ok := resReject1.(*evaluator.Error)
	if !ok || !strings.Contains(errObj1.Message, "Indonesian module 'kripto' is not allowed in Universal Mode") {
		t.Errorf("harus menolak import 'kripto' pada Universal Mode, dapat: %v", resReject1)
	}

	// 4. Universal Mode rejecting 'kripto' global identifier
	scriptRejectIdent := `"use universal"
	let x = kripto.sha256("tes")
	`
	resReject2 := runCryptoScript(scriptRejectIdent)
	errObj2, ok := resReject2.(*evaluator.Error)
	if !ok || !strings.Contains(errObj2.Message, "Indonesian module 'kripto' is not allowed in Universal Mode") {
		t.Errorf("harus menolak identifier 'kripto' pada Universal Mode, dapat: %v", resReject2)
	}
}

func TestStdlibCrypto_Symmetric_RoundTrip(t *testing.T) {
	script := `
	buat kunci = kripto.buat_kunci()
	buat teksAsli = "Pesan rahasia tingkat tinggi dalam Celang!"
	buat terenkripsi = kripto.enkripsi(teksAsli, kunci)
	buat pulih = kripto.dekripsi(terenkripsi, kunci)

	buat terenkripsi2 = kripto.enkripsi(teksAsli, kunci)
	buat bedaCiphertext = (terenkripsi != terenkripsi2) // nonce acak berbeda setiap pemanggilan
	buat pulih2 = kripto.dekripsi(terenkripsi2, kunci)

	// Uji masukan Bita
	buat dataBita = bita([72, 97, 108, 111]) // "Halo"
	buat encBita = kripto.enkripsi(dataBita, kunci)
	buat decBita = kripto.dekripsi(encBita, kunci)

	// Uji string kosong
	buat encKosong = kripto.enkripsi("", kunci)
	buat decKosong = kripto.dekripsi(encKosong, kunci)

	[pulih, bedaCiphertext, pulih2, decBita, decKosong, panjang(kunci)]
	`
	res := runCryptoScript(script)
	arr, ok := res.(*evaluator.Array)
	if !ok {
		t.Fatalf("evaluasi round-trip simetris gagal: %v", res)
	}

	if arr.Elements[0].Inspect() != "Pesan rahasia tingkat tinggi dalam Celang!" {
		t.Errorf("dekripsi teks biasa gagal: %s", arr.Elements[0].Inspect())
	}
	if arr.Elements[1] != evaluator.TRUE {
		t.Errorf("dua pemanggilan enkripsi harus menghasilkan ciphertext berbeda karena nonce acak")
	}
	if arr.Elements[2].Inspect() != "Pesan rahasia tingkat tinggi dalam Celang!" {
		t.Errorf("dekripsi kedua gagal: %s", arr.Elements[2].Inspect())
	}
	if arr.Elements[3].Inspect() != "Halo" {
		t.Errorf("dekripsi data Bita gagal: %s", arr.Elements[3].Inspect())
	}
	if arr.Elements[4].Inspect() != "" {
		t.Errorf("dekripsi data kosong gagal: %s", arr.Elements[4].Inspect())
	}
	keyLen := arr.Elements[5].(*evaluator.Integer).Value
	if keyLen != 44 {
		t.Errorf("panjang kunci Base64 dari buat_kunci() harus 44 karakter, dapat: %d", keyLen)
	}
}

func TestStdlibCrypto_Symmetric_TamperDetection(t *testing.T) {
	initScript := `
	buat kunci = kripto.buat_kunci()
	buat kunciLain = kripto.buat_kunci()
	buat pesan = "Integritas data Celang"
	buat terenkripsi = kripto.enkripsi(pesan, kunci)
	[kunci, kunciLain, terenkripsi]
	`
	resInit := runCryptoScript(initScript)
	arrInit, ok := resInit.(*evaluator.Array)
	if !ok {
		t.Fatalf("evaluasi init tamper test gagal: %v", resInit)
	}
	kunci := arrInit.Elements[0].Inspect()
	kunciLain := arrInit.Elements[1].Inspect()
	ciphertext := arrInit.Elements[2].Inspect()

	// 1. Bit flip test: decode base64, flip 1 bit in tag, re-encode
	payload, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		t.Fatalf("decode base64 ciphertext gagal: %v", err)
	}
	payload[len(payload)-1] ^= 0x01 // flip 1 bit
	tamperedCiphertext := base64.StdEncoding.EncodeToString(payload)

	script := fmt.Sprintf(`
	buat kunci = "%s"
	buat kunciLain = "%s"
	buat terenkripsi = "%s"
	buat terenkripsiRusak = "%s"

	// 1. Kunci salah
	buat errKunciSalah = ""
	coba {
		kripto.dekripsi(terenkripsi, kunciLain)
	} tangkap (e) {
		errKunciSalah = e
	}

	// 2. Modifikasi 1-bit ciphertext
	buat errTamper = ""
	coba {
		kripto.dekripsi(terenkripsiRusak, kunci)
	} tangkap (e) {
		errTamper = e
	}

	// 3. Ciphertext terlalu pendek (< 28 bita payload)
	buat errPendek = ""
	coba {
		kripto.dekripsi("dGVz", kunci) // "tes" Base64 (3 bita)
	} tangkap (e) {
		errPendek = e
	}

	// 4. Ciphertext bukan Base64
	buat errBukanB64 = ""
	coba {
		kripto.dekripsi("!@#$%%^&*()_+", kunci)
	} tangkap (e) {
		errBukanB64 = e
	}

	[errKunciSalah, errTamper, errPendek, errBukanB64]
	`, kunci, kunciLain, ciphertext, tamperedCiphertext)

	res := runCryptoScript(script)
	arr, ok := res.(*evaluator.Array)
	if !ok {
		t.Fatalf("evaluasi tamper detection script gagal: %v", res)
	}

	errKunciSalah := arr.Elements[0].Inspect()
	if !strings.Contains(errKunciSalah, "dekripsi gagal") {
		t.Errorf("pesan galat kunci salah tidak sesuai: %s", errKunciSalah)
	}

	errTamper := arr.Elements[1].Inspect()
	if !strings.Contains(errTamper, "dekripsi gagal") {
		t.Errorf("pesan galat tamper 1-bit tidak sesuai: %s", errTamper)
	}

	errPendek := arr.Elements[2].Inspect()
	if !strings.Contains(errPendek, "format data tersandi tidak valid") {
		t.Errorf("pesan galat data pendek tidak sesuai: %s", errPendek)
	}

	errBukanB64 := arr.Elements[3].Inspect()
	if !strings.Contains(errBukanB64, "format ciphertext tidak valid") {
		t.Errorf("pesan galat bukan Base64 tidak sesuai: %s", errBukanB64)
	}
}

func TestStdlibCrypto_Symmetric_KeyFormats(t *testing.T) {
	script := `
	buat pesan = "Format kunci serbaguna"

	// Kunci heksadesimal 64 karakter (32 bita)
	buat keyHex = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	buat encHex = kripto.enkripsi(pesan, keyHex)
	buat decHex = kripto.dekripsi(encHex, keyHex)

	// Kunci teks mentah 32 bita
	buat keyRaw = "12345678901234567890123456789012"
	buat encRaw = kripto.enkripsi(pesan, keyRaw)
	buat decRaw = kripto.dekripsi(encRaw, keyRaw)

	// Kunci objek Bita 32 bita
	buat keyBita = bita(keyRaw)
	buat encBita = kripto.enkripsi(pesan, keyBita)
	buat decBita = kripto.dekripsi(encBita, keyBita)

	// Dekripsi silang (enkripsi dengan keyRaw, dekripsi dengan keyBita)
	buat decSilang = kripto.dekripsi(encRaw, keyBita)

	// Kunci tidak valid (kurang dari 32 bita)
	buat errKunciPendek = ""
	coba {
		kripto.enkripsi(pesan, "kunci-terlalu-pendek")
	} tangkap (e) {
		errKunciPendek = e
	}

	[decHex == pesan, decRaw == pesan, decBita == pesan, decSilang == pesan, errKunciPendek]
	`
	res := runCryptoScript(script)
	arr, ok := res.(*evaluator.Array)
	if !ok {
		t.Fatalf("evaluasi key formats gagal: %v", res)
	}

	if arr.Elements[0] != evaluator.TRUE || arr.Elements[1] != evaluator.TRUE || arr.Elements[2] != evaluator.TRUE || arr.Elements[3] != evaluator.TRUE {
		t.Errorf("format kunci hex, raw, bita, atau silang gagal")
	}

	errPendek := arr.Elements[4].Inspect()
	if !strings.Contains(errPendek, "format atau panjang kunci tidak valid") {
		t.Errorf("pesan galat kunci pendek tidak sesuai: %s", errPendek)
	}
}

func TestStdlibCrypto_PBKDF2_RFCVectors(t *testing.T) {
	// 1. RFC 7914 Section 11 Test Vector: P="passwd", S="salt", c=1, dkLen=64
	dk1 := evaluator.PBKDF2SHA256([]byte("passwd"), []byte("salt"), 1, 64)
	exp1 := "55ac046e56e3089fec1691c22544b605f94185216dde0465e68b9d57c20dacbc49ca9cccf179b645991664b39d77ef317c71b845b1e30bd509112041d3a19783"
	if hex.EncodeToString(dk1) != exp1 {
		t.Fatalf("RFC 7914 vector 1 gagal: diharapkan %s, dapat %s", exp1, hex.EncodeToString(dk1))
	}

	// 2. Extended HMAC-SHA256 Test Vector: P="password", S="salt", c=4096, dkLen=32
	dk2 := evaluator.PBKDF2SHA256([]byte("password"), []byte("salt"), 4096, 32)
	exp2 := "c5e478d59288c841aa530db6845c4c8d962893a001ce4e11a4963873aa98134a"
	if hex.EncodeToString(dk2) != exp2 {
		t.Fatalf("PBKDF2 4096 iterasi gagal: diharapkan %s, dapat %s", exp2, hex.EncodeToString(dk2))
	}
}

func TestStdlibCrypto_PasswordHashingAndVerification(t *testing.T) {
	script := `
	buat sandi = "SandiRahasiaSuper123!"
	buat h1 = kripto.hash_sandi(sandi)
	buat h2 = crypto.hash_password(sandi)

	buat v1_ok = kripto.verifikasi_sandi(sandi, h1)
	buat v2_ok = crypto.verify_password(sandi, h2)
	buat v1_wrong = kripto.verifikasi_sandi("SandiKeliru", h1)

	// Tampered hash check (modifikasi string hash)
	buat h1_tampered = h1 + "x"
	buat v_tampered = kripto.verifikasi_sandi(sandi, h1_tampered)

	// Tampered hash digest length checks (truncated and extended)
	buat h_trunc1 = "$pbkdf2-sha256$i=100000$c2FsdHNhbHRzYWx0c2E$AQ"
	buat v_trunc1 = kripto.verifikasi_sandi(sandi, h_trunc1)
	buat h_trunc16 = "$pbkdf2-sha256$i=100000$c2FsdHNhbHRzYWx0c2E$AQIDBAUGBwgJCgsMDQ4PEA"
	buat v_trunc16 = kripto.verifikasi_sandi(sandi, h_trunc16)
	buat h_ext33 = "$pbkdf2-sha256$i=100000$c2FsdHNhbHRzYWx0c2E$AQIDBAUGBwgJCgsMDQ4PEBESExQVFhcYGRobHB0eHyEh"
	buat v_ext33 = crypto.verify_password(sandi, h_ext33)

	// Malformed hash checks
	buat v_malformed = kripto.verifikasi_sandi(sandi, "bukan-hash-valid")
	buat v_bad_iter = kripto.verifikasi_sandi(sandi, "$pbkdf2-sha256$i=abc$salt$hash")
	buat v_bad_prefix = kripto.verifikasi_sandi(sandi, "$argon2id$i=1000$salt$hash")

	// Sandi kosong harus menghasilkan error
	buat errKosong = ""
	coba {
		kripto.hash_sandi("")
	} tangkap (e) {
		errKosong = e
	}

	[h1, h2, h1 != h2, v1_ok, v2_ok, v1_wrong, v_tampered, v_trunc1, v_trunc16, v_ext33, v_malformed, v_bad_iter, v_bad_prefix, errKosong]
	`
	res := runCryptoScript(script)
	arr, ok := res.(*evaluator.Array)
	if !ok {
		t.Fatalf("evaluasi password hashing gagal: %v", res)
	}

	h1 := arr.Elements[0].Inspect()
	h2 := arr.Elements[1].Inspect()

	// Format PHC check: $pbkdf2-sha256$i=100000$<salt>$<hash>
	phcRegex := regexp.MustCompile(`^\$pbkdf2-sha256\$i=100000\$[A-Za-z0-9+/=_-]+\$[A-Za-z0-9+/=_-]+$`)
	if !phcRegex.MatchString(h1) {
		t.Errorf("format hash PHC h1 tidak valid: %s", h1)
	}
	if !phcRegex.MatchString(h2) {
		t.Errorf("format hash PHC h2 tidak valid: %s", h2)
	}

	// Uniqueness check: salt acak harus membuat h1 != h2
	if arr.Elements[2] != evaluator.TRUE {
		t.Errorf("dua pemanggilan hash_sandi berturut-turut harus menghasilkan salt unik berbeda: %s == %s", h1, h2)
	}

	// Verification check
	if arr.Elements[3] != evaluator.TRUE || arr.Elements[4] != evaluator.TRUE {
		t.Errorf("verifikasi sandi yang benar harus bernilai true: %v, %v", arr.Elements[3], arr.Elements[4])
	}
	if arr.Elements[5] != evaluator.FALSE {
		t.Errorf("verifikasi sandi keliru harus bernilai false: %v", arr.Elements[5])
	}
	for idx := 6; idx <= 12; idx++ {
		if arr.Elements[idx] != evaluator.FALSE {
			t.Errorf("verifikasi hash tampered/rusak pada index %d harus false, dapat: %v", idx, arr.Elements[idx])
		}
	}

	errKosong := arr.Elements[13].Inspect()
	if !strings.Contains(errKosong, "kata sandi tidak boleh kosong") {
		t.Errorf("pesan galat kata sandi kosong tidak sesuai: %s", errKosong)
	}
}

func TestStdlibCrypto_PasswordVerification_TamperedDigestLengths(t *testing.T) {
	saltB64 := base64.RawStdEncoding.EncodeToString([]byte("1234567890123456"))
	password := "RahasiaCelang2026!"

	testCases := []struct {
		name       string
		digestSize int
	}{
		{"Empty digest (0 bytes)", 0},
		{"1-byte truncated digest", 1},
		{"2-byte truncated digest", 2},
		{"16-byte truncated digest", 16},
		{"31-byte truncated digest", 31},
		{"33-byte extended digest", 33},
		{"48-byte extended digest", 48},
		{"64-byte extended digest", 64},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rawDigest := make([]byte, tc.digestSize)
			for i := range rawDigest {
				rawDigest[i] = byte(i + 1)
			}
			digestB64 := base64.RawStdEncoding.EncodeToString(rawDigest)
			phcHash := fmt.Sprintf("$pbkdf2-sha256$i=1000$%s$%s", saltB64, digestB64)

			scriptID := fmt.Sprintf(`kripto.verifikasi_sandi("%s", "%s")`, password, phcHash)
			resID := runCryptoScript(scriptID)
			if resID != evaluator.FALSE {
				t.Errorf("kripto.verifikasi_sandi dengan %s harus bernilai FALSE, dapat: %v", tc.name, resID)
			}

			scriptEN := fmt.Sprintf(`crypto.verify_password("%s", "%s")`, password, phcHash)
			resEN := runCryptoScript(scriptEN)
			if resEN != evaluator.FALSE {
				t.Errorf("crypto.verify_password dengan %s harus bernilai FALSE, dapat: %v", tc.name, resEN)
			}
		})
	}

	t.Run("CollisionResistance_TamperedValidHash", func(t *testing.T) {
		genScript := fmt.Sprintf(`kripto.hash_sandi("%s")`, password)
		resGen := runCryptoScript(genScript)
		strGen, ok := resGen.(*evaluator.String)
		if !ok || strGen.Value == "" {
			t.Fatalf("gagal membuat hash valid untuk pengujian: %v", resGen)
		}
		validPHC := strGen.Value

		tamperedVariants := []string{validPHC + "x", validPHC + "A", validPHC + "0", validPHC + "=="}
		for _, tampered := range tamperedVariants {
			verifyScript := fmt.Sprintf(`kripto.verifikasi_sandi("%s", "%s")`, password, tampered)
			res := runCryptoScript(verifyScript)
			if res != evaluator.FALSE {
				t.Errorf("tampered hash '%s' harus mengembalikan FALSE, dapat: %v", tampered, res)
			}
		}
	})
}

func TestStdlibCrypto_DeriveKey(t *testing.T) {
	script := `
	buat sandi = "kata-kunci-induk-2026"
	buat salt = kripto.acak_bita(16)
	buat key_bita = kripto.turunkan_kunci(sandi, salt, 32)
	buat key_heks = crypto.derive_key(sandi, salt, 32, "heks")
	buat key_b64 = crypto.derive_key(sandi, salt, 32, "base64")

	// Gunakan derived key langsung untuk enkripsi AES-256-GCM
	buat enc = kripto.enkripsi("pesan rahasia dari derived key", key_bita)
	buat dec = kripto.dekripsi(enc, key_bita)
	buat dec2 = kripto.dekripsi(enc, key_heks)

	// Uji kesalahan argumen
	buat errPanjangNol = ""
	coba {
		kripto.turunkan_kunci(sandi, salt, 0)
	} tangkap (e) {
		errPanjangNol = e
	}

	buat errFormatAsing = ""
	coba {
		kripto.turunkan_kunci(sandi, salt, 32, "unknown_fmt")
	} tangkap (e) {
		errFormatAsing = e
	}

	[panjang(key_bita), panjang(key_heks), dec, dec2, errPanjangNol, errFormatAsing]
	`
	res := runCryptoScript(script)
	arr, ok := res.(*evaluator.Array)
	if !ok {
		t.Fatalf("evaluasi derive_key gagal: %v", res)
	}

	if arr.Elements[0].(*evaluator.Integer).Value != 32 {
		t.Errorf("panjang key_bita harus 32 bita, dapat: %v", arr.Elements[0])
	}
	if arr.Elements[1].(*evaluator.Integer).Value != 64 {
		t.Errorf("panjang key_heks harus 64 karakter heks, dapat: %v", arr.Elements[1])
	}
	if arr.Elements[2].Inspect() != "pesan rahasia dari derived key" || arr.Elements[3].Inspect() != "pesan rahasia dari derived key" {
		t.Errorf("dekripsi menggunakan kunci hasil derive_key gagal: %s vs %s", arr.Elements[2].Inspect(), arr.Elements[3].Inspect())
	}

	errPanjang := arr.Elements[4].Inspect()
	if !strings.Contains(errPanjang, "harus lebih besar dari 0") {
		t.Errorf("pesan galat panjang nol tidak sesuai: %s", errPanjang)
	}

	errFormat := arr.Elements[5].Inspect()
	if !strings.Contains(errFormat, "tidak dikenal") {
		t.Errorf("pesan galat format tidak dikenal tidak sesuai: %s", errFormat)
	}
}

func TestStdlibCrypto_Ed25519_RFC8032(t *testing.T) {
	script := `
	// Vector 1: Pesan kosong (0 bytes)
	buat seed1 = "9d61b19deffd5a60ba844af492ec2cc44449c5697b326919703bac031cae7f60"
	buat pub1 = "d75a980182b10ab7d54bfed3c964073a0ee172f3daa62325af021a68f707511a"
	buat sig1 = kripto.tanda_tangan("", seed1)
	buat ver1 = kripto.verifikasi_tanda_tangan("", sig1, pub1)

	// Vector 2: Pesan 1 byte ('r' / 0x72)
	buat seed2 = "4ccd089b28ff96da9db6c346ec114e0f5b8a319f35aba624da8cf6ed4fb8a6fb"
	buat pub2 = "3d4017c3e843895a92b70aa74d1b7ebc9c982ccf2ec4968cc0cd55f12af4660c"
	buat sig2 = kripto.tanda_tangan("r", seed2)
	buat ver2 = kripto.verifikasi_tanda_tangan("r", sig2, pub2)

	// Vector 3: Pesan bita 2 byte (0xaf, 0x82)
	buat seed3 = "c5aa8df43f9f837bedb7442f31dcb7b166d38535076f094b85ce3a2e0b4458f7"
	buat pub3 = "fc51cd8e6218a1a38da47ed00230f0580816ed13ba3303ac5deb911548908025"
	buat data3 = bita([175, 130]) // 0xaf, 0x82
	buat sig3 = kripto.tanda_tangan(data3, seed3)
	buat ver3 = kripto.verifikasi_tanda_tangan(data3, sig3, pub3)

	[sig1, ver1, sig2, ver2, sig3, ver3]
	`
	res := runCryptoScript(script)
	arr, ok := res.(*evaluator.Array)
	if !ok {
		t.Fatalf("evaluasi RFC 8032 gagal: %v", res)
	}

	expSig1 := "e5564300c360ac729086e2cc806e828a84877f1eb8e5d974d873e065224901555fb8821590a33bacc61e39701cf9b46bd25bf5f0595bbe24655141438e7a100b"
	if arr.Elements[0].Inspect() != expSig1 {
		t.Errorf("RFC 8032 Test 1 tanda tangan salah: diharapkan %s, dapat %s", expSig1, arr.Elements[0].Inspect())
	}
	if arr.Elements[1] != evaluator.TRUE {
		t.Errorf("RFC 8032 Test 1 verifikasi gagal")
	}

	expSig2 := "92a009a9f0d4cab8720e820b5f642540a2b27b5416503f8fb3762223ebdb69da085ac1e43e15996e458f3613d0f11d8c387b2eaeb4302aeeb00d291612bb0c00"
	if arr.Elements[2].Inspect() != expSig2 {
		t.Errorf("RFC 8032 Test 2 tanda tangan salah: diharapkan %s, dapat %s", expSig2, arr.Elements[2].Inspect())
	}
	if arr.Elements[3] != evaluator.TRUE {
		t.Errorf("RFC 8032 Test 2 verifikasi gagal")
	}

	expSig3 := "6291d657deec24024827e69c3abe01a30ce548a284743a445e3680d7db5ac3ac18ff9b538d16f290ae67f760984dc6594a7c15e9716ed28dc027beceea1ec40a"
	if arr.Elements[4].Inspect() != expSig3 {
		t.Errorf("RFC 8032 Test 3 tanda tangan salah: diharapkan %s, dapat %s", expSig3, arr.Elements[4].Inspect())
	}
	if arr.Elements[5] != evaluator.TRUE {
		t.Errorf("RFC 8032 Test 3 verifikasi gagal")
	}
}

func TestStdlibCrypto_Ed25519_RoundTripAndTamper(t *testing.T) {
	script := `
	buat pasangan = kripto.buat_pasangan_kunci()
	buat pesan = "kontrak-paket-celang-v1.0"
	buat ttd = kripto.tanda_tangan(pesan, pasangan.kunci_privat)

	buat sah = kripto.verifikasi_tanda_tangan(pesan, ttd, pasangan.kunci_publik)
	buat pesanUbah = "kontrak-paket-celang-v1.1"
	buat sahUbah = kripto.verifikasi_tanda_tangan(pesanUbah, ttd, pasangan.kunci_publik)

	// Modifikasi tanda tangan (1 karakter / bit)
	buat ttdRusak = ttd + "0"
	buat sahRusak = kripto.verifikasi_tanda_tangan(pesan, ttdRusak, pasangan.kunci_publik)

	// Kunci publik yang salah
	buat pasanganLain = kripto.buat_pasangan_kunci()
	buat sahKunciSalah = kripto.verifikasi_tanda_tangan(pesan, ttd, pasanganLain.kunci_publik)

	// Kunci privat rusak (64 karakter non-heks 'z') harus menghasilkan galat empati
	buat errKunci64z = ""
	coba {
		kripto.tanda_tangan(pesan, "zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz")
	} tangkap (e) {
		errKunci64z = e
	}

	// Pengujian seluruh alias: verifikasi_ttd, verify, verify_signature
	buat sahAlias1 = kripto.verifikasi_ttd(pesan, ttd, pasangan.kunci_publik)
	buat sahAlias2 = crypto.verify(pesan, ttd, pasangan.public_key)
	buat sahAlias3 = crypto.verify_signature(pesan, ttd, pasangan.public_key)

	[sah, sahUbah, sahRusak, sahKunciSalah, errKunci64z, sahAlias1, sahAlias2, sahAlias3, panjang(pasangan.kunci_publik), panjang(pasangan.kunci_privat), panjang(ttd)]
	`
	res := runCryptoScript(script)
	arr, ok := res.(*evaluator.Array)
	if !ok {
		t.Fatalf("evaluasi round-trip Ed25519 gagal: %v", res)
	}

	if arr.Elements[0] != evaluator.TRUE {
		t.Errorf("verifikasi tanda tangan round-trip harus bernilai benar")
	}
	if arr.Elements[1] != evaluator.FALSE {
		t.Errorf("verifikasi pesan diubah harus bernilai salah")
	}
	if arr.Elements[2] != evaluator.FALSE {
		t.Errorf("verifikasi ttd rusak harus bernilai salah")
	}
	if arr.Elements[3] != evaluator.FALSE {
		t.Errorf("verifikasi kunci publik salah harus bernilai salah")
	}
	errKunci64z := arr.Elements[4].Inspect()
	if !strings.Contains(errKunci64z, "kunci privat") {
		t.Errorf("tanda_tangan dengan kunci 64 'z' harus menghasilkan galat empati tentang kunci privat, dapat: %s", errKunci64z)
	}
	if arr.Elements[5] != evaluator.TRUE || arr.Elements[6] != evaluator.TRUE || arr.Elements[7] != evaluator.TRUE {
		t.Errorf("alias verifikasi_ttd, verify, verify_signature harus bernilai benar")
	}

	pubLen := arr.Elements[8].(*evaluator.Integer).Value
	privLen := arr.Elements[9].(*evaluator.Integer).Value
	sigLen := arr.Elements[10].(*evaluator.Integer).Value
	if pubLen != 64 {
		t.Errorf("panjang kunci publik heks harus 64, dapat: %d", pubLen)
	}
	if privLen != 64 {
		t.Errorf("panjang kunci privat heks (seed) harus 64, dapat: %d", privLen)
	}
	if sigLen != 128 {
		t.Errorf("panjang tanda tangan heks harus 128, dapat: %d", sigLen)
	}
}

func TestStdlibCrypto_Ed25519_CorruptedPrivateKey(t *testing.T) {
	badKeys := []struct {
		name string
		key  string
	}{
		{"64 'z' characters", strings.Repeat("z", 64)},
		{"64-char hex with typo at end", "9d61b19deffd5a60ba844af492ec2cc44449c5697b326919703bac031cae7f6z"},
		{"64-char hex with typo at start", "zd61b19deffd5a60ba844af492ec2cc44449c5697b326919703bac031cae7f60"},
		{"64-char arbitrary ASCII text", "this is a sixty four character long secret key for my celang app"},
		{"128 'z' characters", strings.Repeat("z", 128)},
	}

	for _, tc := range badKeys {
		t.Run(tc.name, func(t *testing.T) {
			// 1. Direct return: must return Error object, not String signature
			scriptID := fmt.Sprintf(`kripto.tanda_tangan("pesan rahasia", "%s")`, tc.key)
			resID := runCryptoScript(scriptID)
			errObjID, isErrID := resID.(*evaluator.Error)
			if !isErrID {
				t.Fatalf("kripto.tanda_tangan dengan '%s' harus menghasilkan Error, dapat: %T (%v)", tc.name, resID, resID)
			}
			if !strings.Contains(errObjID.Message, "kunci privat") {
				t.Errorf("pesan galat harus memuat 'kunci privat', dapat: %s", errObjID.Message)
			}

			// 2. English alias: crypto.sign
			scriptEN := fmt.Sprintf(`crypto.sign("secret message", "%s")`, tc.key)
			resEN := runCryptoScript(scriptEN)
			errObjEN, isErrEN := resEN.(*evaluator.Error)
			if !isErrEN {
				t.Fatalf("crypto.sign dengan '%s' harus menghasilkan Error, dapat: %T (%v)", tc.name, resEN, resEN)
			}
			if !strings.Contains(errObjEN.Message, "kunci privat") {
				t.Errorf("pesan galat crypto.sign harus memuat 'kunci privat', dapat: %s", errObjEN.Message)
			}

			// 3. Exception handling via coba-tangkap in Celang
			scriptCoba := fmt.Sprintf(`
				buat err = ""
				coba {
					kripto.tanda_tangan("pesan", "%s")
				} tangkap (e) {
					err = e
				}
				err
			`, tc.key)
			resCoba := runCryptoScript(scriptCoba)
			strCoba, ok := resCoba.(*evaluator.String)
			if !ok || strCoba.Value == "" {
				t.Errorf("coba-tangkap untuk '%s' harus menangkap galat, dapat: %v", tc.name, resCoba)
			}
		})
	}

	t.Run("InconsistentExpandedKey_128Hex", func(t *testing.T) {
		seedHex := "9d61b19deffd5a60ba844af492ec2cc44449c5697b326919703bac031cae7f60"
		fakePubHex := strings.Repeat("0", 64)
		inconsistentKey := seedHex + fakePubHex // 128 hex chars

		script := fmt.Sprintf(`kripto.tanda_tangan("pesan", "%s")`, inconsistentKey)
		res := runCryptoScript(script)
		errObj, isErr := res.(*evaluator.Error)
		if !isErr {
			t.Fatalf("kripto.tanda_tangan dengan inconsistent expanded key harus menghasilkan Error, dapat: %T (%v)", res, res)
		}
		if !strings.Contains(errObj.Message, "tidak konsisten antara seed dan public key") {
			t.Errorf("pesan galat harus memuat informasi inkonsistensi: %s", errObj.Message)
		}
	})

	t.Run("InconsistentExpandedKey_Bytes64", func(t *testing.T) {
		script := `
		buat bitaInkonsisten = bita(64)
		kripto.tanda_tangan("pesan", bitaInkonsisten)
		`
		res := runCryptoScript(script)
		errObj, isErr := res.(*evaluator.Error)
		if !isErr {
			t.Fatalf("kripto.tanda_tangan dengan 64-byte Bita inkonsisten harus Error, dapat: %T (%v)", res, res)
		}
		if !strings.Contains(errObj.Message, "tidak konsisten antara seed dan public key") {
			t.Errorf("pesan galat harus memuat informasi inkonsistensi: %s", errObj.Message)
		}
	})
}

func TestStdlibCrypto_BilingualAndUniversal(t *testing.T) {
	// 1. Script dalam mode "use universal"
	scriptUniv := `"use universal"
	import { generate_key, encrypt, decrypt, hash_password, verify_password, derive_key, generate_keypair, sign, verify } from "crypto"

	let key = generate_key()
	let enc = encrypt("universal secret message", key)
	let dec = decrypt(enc, key)

	let pwHash = hash_password("mypassword")
	let pwOk = verify_password("mypassword", pwHash)

	let salt = "mysalt1234567890"
	let dKey = derive_key("mypassword", salt, 32, "hex")

	let pair = generate_keypair()
	let sig = sign("payload", pair.private_key)
	let sigOk = verify("payload", sig, pair.public_key)

	[dec == "universal secret message", pwOk, len(dKey) == 64, sigOk]
	`
	resUniv := runCryptoScript(scriptUniv)
	arrUniv, ok := resUniv.(*evaluator.Array)
	if !ok {
		t.Fatalf("evaluasi universal crypto suite gagal: %v", resUniv)
	}

	for i := 0; i < 4; i++ {
		if arrUniv.Elements[i] != evaluator.TRUE {
			t.Errorf("universal suite assertion %d gagal: %v", i, arrUniv.Elements[i])
		}
	}

	// 2. Akses global crypto dalam "use universal"
	scriptUnivGlobal := `"use universal"
	let k = crypto.generate_key()
	let p = crypto.generate_keypair()
	[len(k) == 44, len(p.public_key) == 64]
	`
	resGlobal := runCryptoScript(scriptUnivGlobal)
	arrGlobal, ok := resGlobal.(*evaluator.Array)
	if !ok {
		t.Fatalf("evaluasi universal global crypto gagal: %v", resGlobal)
	}
	if arrGlobal.Elements[0] != evaluator.TRUE || arrGlobal.Elements[1] != evaluator.TRUE {
		t.Errorf("universal global crypto assertions gagal: %v", arrGlobal)
	}
}

func TestCrypto_FileHashAndCRC32(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "celang_crypto_test_*.txt")
	if err != nil {
		t.Fatalf("gagal membuat berkas sementara: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	content := "Celang File Streaming Hash & CRC32 Verification"
	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("gagal menulis ke berkas sementara: %v", err)
	}
	tmpFile.Close()

	script := fmt.Sprintf(`
	buat h256 = kripto.sha256("%s")
	buat h256_berkas = kripto.sha256_berkas("%s")
	buat h512_berkas = crypto.sha512_file("%s")
	buat md5_berkas = crypto.md5_file("%s")

	buat c1 = kripto.crc32("%s")
	buat c2 = kripto.crc32(bita("%s"))
	buat c_berkas = kripto.crc32_berkas("%s")
	buat c_file = crypto.crc32_file("%s")

	[
		h256 == h256_berkas,
		panjang(h512_berkas) == 128,
		panjang(md5_berkas) == 32,
		c1 == c2,
		c1 == c_berkas,
		c_berkas == c_file,
		c1 > 0
	]
	`, content, tmpFile.Name(), tmpFile.Name(), tmpFile.Name(), content, content, tmpFile.Name(), tmpFile.Name())

	res := runCryptoScript(script)
	arr, ok := res.(*evaluator.Array)
	if !ok {
		t.Fatalf("evaluasi FileHashAndCRC32 gagal: %v", res)
	}

	for i := 0; i < len(arr.Elements); i++ {
		if arr.Elements[i] != evaluator.TRUE {
			t.Errorf("FileHashAndCRC32 assertion %d gagal (dapat %v)", i, arr.Elements[i])
		}
	}
}

