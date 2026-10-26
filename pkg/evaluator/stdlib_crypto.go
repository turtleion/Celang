package evaluator

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ed25519"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"hash"
	"hash/crc32"
	"io"
	"os"
	"strconv"
	"strings"
)

func getCryptoDataBytes(arg Object, funcName string, argPos string) ([]byte, Object) {
	switch v := arg.(type) {
	case *String:
		return []byte(v.Value), nil
	case *Bytes:
		return v.Elements, nil
	default:
		return nil, newError(nil, fmt.Sprintf("%s: argumen %sharus berupa teks atau Bita, menerima %s", funcName, argPos, arg.Type()), "Contoh: "+funcName+"(\"data\") atau "+funcName+"(bita([65, 66]))")
	}
}

// parseSymmetricKey mengurai argumen kunci menjadi 32 bita murni ([]byte).
func parseSymmetricKey(keyArg Object, funcName string) ([]byte, Object) {
	switch v := keyArg.(type) {
	case *Bytes:
		if len(v.Elements) != 32 {
			return nil, newError(nil,
				fmt.Sprintf("%s: panjang kunci Bita tidak valid (ditemukan %d bita, diharapkan 32 bita / 256-bit)", funcName, len(v.Elements)),
				"Gunakan kripto.buat_kunci() untuk menghasilkan kunci 32-bita yang aman.")
		}
		return v.Elements, nil

	case *String:
		str := strings.TrimSpace(v.Value)

		// 1. Coba decode Base64 standar / URL (panjang 44 atau 43 karakter)
		if len(str) == 44 || len(str) == 43 {
			if b64, err := base64.StdEncoding.DecodeString(str); err == nil && len(b64) == 32 {
				return b64, nil
			}
			if b64Raw, err := base64.RawStdEncoding.DecodeString(str); err == nil && len(b64Raw) == 32 {
				return b64Raw, nil
			}
		}

		// 2. Coba decode Heksadesimal (64 heks)
		if len(str) == 64 {
			if hx, err := hex.DecodeString(str); err == nil && len(hx) == 32 {
				return hx, nil
			}
		}

		// 3. String teks mentah sepanjang 32 bita
		if len([]byte(str)) == 32 {
			return []byte(str), nil
		}

		// Upaya decode Base64 umum
		if b64, err := base64.StdEncoding.DecodeString(str); err == nil && len(b64) == 32 {
			return b64, nil
		}

		return nil, newError(nil,
			fmt.Sprintf("%s: format atau panjang kunci tidak valid (diharapkan Base64 44-karakter, Hex 64-karakter, atau 32 bita)", funcName),
			"Gunakan kripto.buat_kunci() untuk membuat kunci 32-bita (256-bit) yang valid.")

	default:
		return nil, newError(nil,
			fmt.Sprintf("%s: argumen kunci harus berupa teks atau Bita, menerima %s", funcName, keyArg.Type()),
			"Contoh: "+funcName+"(data, kunci)")
	}
}

// pbkdf2SHA256 computes PBKDF2 using HMAC-SHA256 according to RFC 2898 / RFC 8018.
// It is pure Go and relies exclusively on crypto/hmac and crypto/sha256.
func pbkdf2SHA256(password, salt []byte, iter, keyLen int) []byte {
	if keyLen <= 0 || iter <= 0 {
		return nil
	}

	prf := hmac.New(sha256.New, password)
	hashLen := prf.Size() // 32
	numBlocks := (keyLen + hashLen - 1) / hashLen

	var buf [4]byte
	dk := make([]byte, numBlocks*hashLen)
	u := make([]byte, hashLen)

	for block := 1; block <= numBlocks; block++ {
		// U_1 = PRF(password, salt || INT_32_BE(block))
		prf.Reset()
		prf.Write(salt)
		buf[0] = byte(block >> 24)
		buf[1] = byte(block >> 16)
		buf[2] = byte(block >> 8)
		buf[3] = byte(block)
		prf.Write(buf[:4])
		u = prf.Sum(u[:0])

		blockStart := (block - 1) * hashLen
		copy(dk[blockStart:blockStart+hashLen], u)

		// U_2 through U_iter
		for i := 2; i <= iter; i++ {
			prf.Reset()
			prf.Write(u)
			u = prf.Sum(u[:0])
			for j := 0; j < hashLen; j++ {
				dk[blockStart+j] ^= u[j]
			}
		}
	}

	return dk[:keyLen]
}

// PBKDF2SHA256 mengekspos pbkdf2SHA256 untuk pengujian langsung vektor RFC.
func PBKDF2SHA256(password, salt []byte, iter, keyLen int) []byte {
	return pbkdf2SHA256(password, salt, iter, keyLen)
}

// decodeBase64Flex tries unpadded RawStdEncoding first, then standard StdEncoding.
func decodeBase64Flex(s string) ([]byte, error) {
	if b, err := base64.RawStdEncoding.DecodeString(s); err == nil {
		return b, nil
	}
	return base64.StdEncoding.DecodeString(s)
}

// validateEd25519ExpandedKey memverifikasi konsistensi internal kunci privat 64-byte (seed vs public key).
func validateEd25519ExpandedKey(priv ed25519.PrivateKey, funcName string) (ed25519.PrivateKey, Object) {
	expectedPub := ed25519.NewKeyFromSeed(priv[:32])[32:]
	if subtle.ConstantTimeCompare(priv[32:], expectedPub) != 1 {
		return nil, newError(nil,
			fmt.Sprintf("%s: kunci privat 64-byte tidak konsisten antara seed dan public key", funcName),
			"Gunakan kunci privat yang sah dari kripto.buat_pasangan_kunci().")
	}
	return priv, nil
}

// parseEd25519PrivateKey memvalidasi dan memuat kunci privat Ed25519.
// Menerima seed 32-byte (64 heks) atau private key lengkap 64-byte (128 heks),
// baik dalam format heksadesimal, Base64, teks mentah (seed saja), atau objek Bita.
func parseEd25519PrivateKey(arg Object, funcName string) (ed25519.PrivateKey, Object) {
	switch v := arg.(type) {
	case *Bytes:
		if len(v.Elements) == ed25519.SeedSize {
			return ed25519.NewKeyFromSeed(v.Elements), nil
		}
		if len(v.Elements) == ed25519.PrivateKeySize {
			return validateEd25519ExpandedKey(ed25519.PrivateKey(v.Elements), funcName)
		}
		return nil, newError(nil, fmt.Sprintf("%s: panjang kunci privat Bita tidak valid (%d bita, diharapkan 32 atau 64 bita)", funcName, len(v.Elements)), "Gunakan kunci privat dari kripto.buat_pasangan_kunci().")
	case *String:
		str := strings.TrimSpace(v.Value)

		// 1. Kunci privat heksadesimal 64-karakter (32-byte seed)
		if len(str) == 64 {
			decoded, err := hex.DecodeString(str)
			if err != nil || len(decoded) != ed25519.SeedSize {
				return nil, newError(nil,
					fmt.Sprintf("%s: kunci privat heksadesimal 64-karakter tidak valid", funcName),
					"Periksa kembali karakter heksadesimal kunci privat Anda.")
			}
			return ed25519.NewKeyFromSeed(decoded), nil
		}

		// 2. Kunci privat heksadesimal 128-karakter (64-byte expanded key)
		if len(str) == 128 {
			decoded, err := hex.DecodeString(str)
			if err != nil || len(decoded) != ed25519.PrivateKeySize {
				return nil, newError(nil,
					fmt.Sprintf("%s: kunci privat heksadesimal 128-karakter tidak valid", funcName),
					"Periksa kembali karakter heksadesimal kunci privat Anda.")
			}
			return validateEd25519ExpandedKey(ed25519.PrivateKey(decoded), funcName)
		}

		// 3. Kunci privat Base64 (StdEncoding / RawStdEncoding)
		if decoded, err := decodeBase64Flex(str); err == nil {
			if len(decoded) == ed25519.SeedSize {
				return ed25519.NewKeyFromSeed(decoded), nil
			}
			if len(decoded) == ed25519.PrivateKeySize {
				return validateEd25519ExpandedKey(ed25519.PrivateKey(decoded), funcName)
			}
		}

		// 4. String bita mentah (hanya seed 32-byte; 64-byte raw bita harus menggunakan objek Bita)
		if len(str) == ed25519.SeedSize {
			return ed25519.NewKeyFromSeed([]byte(str)), nil
		}

		return nil, newError(nil, fmt.Sprintf("%s: format kunci privat tidak valid (harus berupa heksadesimal 64/128 karakter, Base64, atau Bita)", funcName), "Gunakan kunci privat dari kripto.buat_pasangan_kunci().")
	default:
		return nil, newError(nil, fmt.Sprintf("%s: argumen kunci privat harus berupa teks atau Bita, menerima %s", funcName, arg.Type()), "Contoh: "+funcName+"(data, pasangan.kunci_privat)")
	}
}

// parseEd25519PublicKey mengekstrak kunci publik 32-byte. Mengembalikan nil jika tidak valid.
func parseEd25519PublicKey(arg Object) ed25519.PublicKey {
	switch v := arg.(type) {
	case *Bytes:
		if len(v.Elements) == ed25519.PublicKeySize {
			return ed25519.PublicKey(v.Elements)
		}
		return nil
	case *String:
		str := strings.TrimSpace(v.Value)
		if len(str) == 64 {
			if decoded, err := hex.DecodeString(str); err == nil && len(decoded) == ed25519.PublicKeySize {
				return ed25519.PublicKey(decoded)
			}
		}
		if decoded, err := base64.StdEncoding.DecodeString(str); err == nil && len(decoded) == ed25519.PublicKeySize {
			return ed25519.PublicKey(decoded)
		}
		if decoded, err := base64.RawStdEncoding.DecodeString(str); err == nil && len(decoded) == ed25519.PublicKeySize {
			return ed25519.PublicKey(decoded)
		}
		if len(str) == ed25519.PublicKeySize {
			return ed25519.PublicKey([]byte(str))
		}
		return nil
	default:
		return nil
	}
}

// parseEd25519Signature mengekstrak tanda tangan 64-byte. Mengembalikan nil jika tidak valid.
func parseEd25519Signature(arg Object) []byte {
	switch v := arg.(type) {
	case *Bytes:
		if len(v.Elements) == ed25519.SignatureSize {
			return v.Elements
		}
		return nil
	case *String:
		str := strings.TrimSpace(v.Value)
		if len(str) == 128 {
			if decoded, err := hex.DecodeString(str); err == nil && len(decoded) == ed25519.SignatureSize {
				return decoded
			}
		}
		if decoded, err := base64.StdEncoding.DecodeString(str); err == nil && len(decoded) == ed25519.SignatureSize {
			return decoded
		}
		if decoded, err := base64.RawStdEncoding.DecodeString(str); err == nil && len(decoded) == ed25519.SignatureSize {
			return decoded
		}
		if len(str) == ed25519.SignatureSize {
			return []byte(str)
		}
		return nil
	default:
		return nil
	}
}

// hashFileStream membaca berkas dari disk dalam chunk 64KB dan menghitung digest hash secara streaming.
func hashFileStream(filePath string, h hash.Hash, funcName string) (string, Object) {
	cleanPath := strings.TrimSpace(filePath)
	if cleanPath == "" {
		return "", newError(nil, fmt.Sprintf("%s: jalur berkas tidak boleh kosong", funcName), "Contoh: "+funcName+"(\"./berkas.tar.gz\")")
	}

	f, err := os.Open(cleanPath)
	if err != nil {
		return "", newError(nil, fmt.Sprintf("%s: gagal membuka berkas '%s': %v", funcName, cleanPath, err), "Pastikan berkas tersebut ada dan memiliki izin baca.")
	}
	defer f.Close()

	buf := make([]byte, 64*1024)
	if _, err := io.CopyBuffer(h, f, buf); err != nil {
		return "", newError(nil, fmt.Sprintf("%s: gagal membaca berkas '%s': %v", funcName, cleanPath, err), "")
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

func initStdlibKripto(env *Environment) {
	// =============================================================
	// 1. Fungsi Hash Kriptografis (Cryptographic Hash Functions)
	// =============================================================

	// sha256(data)
	sha256Fn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "kripto.sha256 mengharapkan 1 argumen (teks atau Bita)", "Contoh: kripto.sha256(\"halo dunia\")")
			}
			b, errObj := getCryptoDataBytes(args[0], "kripto.sha256", "")
			if errObj != nil {
				return errObj
			}
			hash := sha256.Sum256(b)
			return &String{Value: hex.EncodeToString(hash[:])}
		},
	}
	definePublicFunction(env, "sha256", sha256Fn)

	// sha512(data)
	sha512Fn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "kripto.sha512 mengharapkan 1 argumen (teks atau Bita)", "Contoh: kripto.sha512(\"halo dunia\")")
			}
			b, errObj := getCryptoDataBytes(args[0], "kripto.sha512", "")
			if errObj != nil {
				return errObj
			}
			hash := sha512.Sum512(b)
			return &String{Value: hex.EncodeToString(hash[:])}
		},
	}
	definePublicFunction(env, "sha512", sha512Fn)

	// sha1(data)
	sha1Fn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "kripto.sha1 mengharapkan 1 argumen (teks atau Bita)", "Contoh: kripto.sha1(\"halo dunia\")")
			}
			b, errObj := getCryptoDataBytes(args[0], "kripto.sha1", "")
			if errObj != nil {
				return errObj
			}
			hash := sha1.Sum(b)
			return &String{Value: hex.EncodeToString(hash[:])}
		},
	}
	definePublicFunction(env, "sha1", sha1Fn)

	// md5(data)
	md5Fn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "kripto.md5 mengharapkan 1 argumen (teks atau Bita)", "Contoh: kripto.md5(\"halo dunia\")")
			}
			b, errObj := getCryptoDataBytes(args[0], "kripto.md5", "")
			if errObj != nil {
				return errObj
			}
			hash := md5.Sum(b)
			return &String{Value: hex.EncodeToString(hash[:])}
		},
	}
	definePublicFunction(env, "md5", md5Fn)

	// sha256_berkas / sha256_file(jalur)
	sha256BerkasFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "kripto.sha256_berkas mengharapkan 1 argumen teks jalur berkas", "Contoh: kripto.sha256_berkas(\"./paket.tar.gz\")")
			}
			pathStr, ok := args[0].(*String)
			if !ok {
				return newError(nil, fmt.Sprintf("argumen jalur berkas kripto.sha256_berkas harus berupa teks, menerima %s", args[0].Type()), "")
			}
			hexVal, errObj := hashFileStream(pathStr.Value, sha256.New(), "kripto.sha256_berkas")
			if errObj != nil {
				return errObj
			}
			return &String{Value: hexVal}
		},
	}
	definePublicFunction(env, "sha256_berkas", sha256BerkasFn)
	definePublicFunction(env, "sha256_file", sha256BerkasFn)

	// sha512_berkas / sha512_file(jalur)
	sha512BerkasFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "kripto.sha512_berkas mengharapkan 1 argumen teks jalur berkas", "Contoh: kripto.sha512_berkas(\"./paket.tar.gz\")")
			}
			pathStr, ok := args[0].(*String)
			if !ok {
				return newError(nil, fmt.Sprintf("argumen jalur berkas kripto.sha512_berkas harus berupa teks, menerima %s", args[0].Type()), "")
			}
			hexVal, errObj := hashFileStream(pathStr.Value, sha512.New(), "kripto.sha512_berkas")
			if errObj != nil {
				return errObj
			}
			return &String{Value: hexVal}
		},
	}
	definePublicFunction(env, "sha512_berkas", sha512BerkasFn)
	definePublicFunction(env, "sha512_file", sha512BerkasFn)

	// md5_berkas / md5_file(jalur)
	md5BerkasFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "kripto.md5_berkas mengharapkan 1 argumen teks jalur berkas", "Contoh: kripto.md5_berkas(\"./paket.tar.gz\")")
			}
			pathStr, ok := args[0].(*String)
			if !ok {
				return newError(nil, fmt.Sprintf("argumen jalur berkas kripto.md5_berkas harus berupa teks, menerima %s", args[0].Type()), "")
			}
			hexVal, errObj := hashFileStream(pathStr.Value, md5.New(), "kripto.md5_berkas")
			if errObj != nil {
				return errObj
			}
			return &String{Value: hexVal}
		},
	}
	definePublicFunction(env, "md5_berkas", md5BerkasFn)
	definePublicFunction(env, "md5_file", md5BerkasFn)

	// crc32(data)
	crc32Fn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "kripto.crc32 mengharapkan 1 argumen (teks atau Bita)", "Contoh: kripto.crc32(\"halo dunia\")")
			}
			b, errObj := getCryptoDataBytes(args[0], "kripto.crc32", "")
			if errObj != nil {
				return errObj
			}
			val := crc32.ChecksumIEEE(b)
			return &Integer{Value: int64(val)}
		},
	}
	definePublicFunction(env, "crc32", crc32Fn)

	// crc32_berkas / crc32_file(jalur)
	crc32BerkasFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "kripto.crc32_berkas mengharapkan 1 argumen teks jalur berkas", "Contoh: kripto.crc32_berkas(\"./berkas.zip\")")
			}
			pathStr, ok := args[0].(*String)
			if !ok {
				return newError(nil, fmt.Sprintf("argumen jalur berkas kripto.crc32_berkas harus berupa teks, menerima %s", args[0].Type()), "")
			}
			cleanPath := strings.TrimSpace(pathStr.Value)
			if cleanPath == "" {
				return newError(nil, "kripto.crc32_berkas: jalur berkas tidak boleh kosong", "")
			}
			f, err := os.Open(cleanPath)
			if err != nil {
				return newError(nil, fmt.Sprintf("kripto.crc32_berkas: gagal membuka berkas '%s': %v", cleanPath, err), "Pastikan berkas tersebut ada dan memiliki izin baca.")
			}
			defer f.Close()

			h := crc32.NewIEEE()
			buf := make([]byte, 64*1024)
			if _, err := io.CopyBuffer(h, f, buf); err != nil {
				return newError(nil, fmt.Sprintf("kripto.crc32_berkas: gagal membaca berkas '%s': %v", cleanPath, err), "")
			}
			return &Integer{Value: int64(h.Sum32())}
		},
	}
	definePublicFunction(env, "crc32_berkas", crc32BerkasFn)
	definePublicFunction(env, "crc32_file", crc32BerkasFn)

	// =============================================================
	// 2. Kode Autentikasi Pesan Berbasis Hash (HMAC)
	// =============================================================

	// hmac_sha256(kunci, data)
	hmacFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "kripto.hmac_sha256 mengharapkan 2 argumen: kunci dan data", "Contoh: kripto.hmac_sha256(\"rahasia\", \"pesan\")")
			}
			key, err1 := getCryptoDataBytes(args[0], "kripto.hmac_sha256", "kunci ")
			if err1 != nil {
				return err1
			}
			data, err2 := getCryptoDataBytes(args[1], "kripto.hmac_sha256", "data ")
			if err2 != nil {
				return err2
			}
			mac := hmac.New(sha256.New, key)
			mac.Write(data)
			return &String{Value: hex.EncodeToString(mac.Sum(nil))}
		},
	}
	definePublicFunction(env, "hmac_sha256", hmacFn)
	definePublicFunction(env, "hmac", hmacFn)

	// =============================================================
	// 3. Pembangkit Keacakan Kriptografis (Cryptographic Randomness)
	// =============================================================

	// acak_heks / random_hex(panjang_bita)
	acakHeksFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "kripto.acak_heks mengharapkan 1 argumen: jumlah bita acak (angka)", "Contoh: kripto.acak_heks(16)")
			}
			intVal, ok := args[0].(*Integer)
			if !ok {
				return newError(nil, fmt.Sprintf("argumen panjang kripto.acak_heks harus berupa angka bulat, menerima %s", args[0].Type()), "")
			}
			if intVal.Value < 0 {
				return newError(nil, "panjang kripto.acak_heks tidak boleh negatif", "")
			}
			if intVal.Value > 1024*1024 {
				return newError(nil, "panjang kripto.acak_heks melebihi batas aman (1MB)", "")
			}
			buf := make([]byte, int(intVal.Value))
			if _, err := rand.Read(buf); err != nil {
				return newError(nil, fmt.Sprintf("gagal menghasilkan keacakan kriptografis: %v", err), "")
			}
			return &String{Value: hex.EncodeToString(buf)}
		},
	}
	definePublicFunction(env, "acak_heks", acakHeksFn)
	definePublicFunction(env, "random_hex", acakHeksFn)

	// acak_bita / random_bytes(panjang_bita)
	acakBitaFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "kripto.acak_bita mengharapkan 1 argumen: jumlah bita acak (angka)", "Contoh: kripto.acak_bita(16)")
			}
			intVal, ok := args[0].(*Integer)
			if !ok {
				return newError(nil, fmt.Sprintf("argumen panjang kripto.acak_bita harus berupa angka bulat, menerima %s", args[0].Type()), "")
			}
			if intVal.Value < 0 {
				return newError(nil, "panjang kripto.acak_bita tidak boleh negatif", "")
			}
			if intVal.Value > 100*1024*1024 {
				return newError(nil, "panjang kripto.acak_bita melebihi batas aman (100MB)", "")
			}
			buf := make([]byte, int(intVal.Value))
			if _, err := rand.Read(buf); err != nil {
				return newError(nil, fmt.Sprintf("gagal menghasilkan keacakan kriptografis: %v", err), "")
			}
			return &Bytes{Elements: buf}
		},
	}
	definePublicFunction(env, "acak_bita", acakBitaFn)
	definePublicFunction(env, "random_bytes", acakBitaFn)

	// uuid() / uuid_v4() -> RFC 4122 v4 UUID
	uuidFn := &Builtin{
		Fn: func(args ...Object) Object {
			var u [16]byte
			if _, err := rand.Read(u[:]); err != nil {
				return newError(nil, fmt.Sprintf("gagal menghasilkan UUID: %v", err), "")
			}
			u[6] = (u[6] & 0x0f) | 0x40 // Versi 4
			u[8] = (u[8] & 0x3f) | 0x80 // Varian 1 (RFC 4122)
			uuidStr := fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:16])
			return &String{Value: uuidStr}
		},
	}
	definePublicFunction(env, "uuid", uuidFn)
	definePublicFunction(env, "uuid_v4", uuidFn)
	definePublicFunction(env, "uuid4", uuidFn)

	// =============================================================
	// 4. Perbandingan Konstan Anti-Timing Attack
	// =============================================================

	// konstanta_sama / timing_safe_equal(a, b)
	equalFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "kripto.konstanta_sama mengharapkan 2 argumen untuk perbandingan aman", "Contoh: kripto.konstanta_sama(tokenA, tokenB)")
			}
			b1, err1 := getCryptoDataBytes(args[0], "kripto.konstanta_sama", "pertama ")
			if err1 != nil {
				return err1
			}
			b2, err2 := getCryptoDataBytes(args[1], "kripto.konstanta_sama", "kedua ")
			if err2 != nil {
				return err2
			}
			match := subtle.ConstantTimeCompare(b1, b2) == 1
			return nativeBool(match)
		},
	}
	definePublicFunction(env, "konstanta_sama", equalFn)
	definePublicFunction(env, "banding_aman", equalFn)
	definePublicFunction(env, "timing_safe_equal", equalFn)
	definePublicFunction(env, "constant_time_compare", equalFn)

	// =============================================================
	// 5. Enkripsi Simetris Terotentikasi (AES-256-GCM AEAD)
	// =============================================================

	// buat_kunci / generate_key() -> Base64 32 random bytes
	buatKunciFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 0 {
				return newError(nil, "kripto.buat_kunci tidak memerlukan argumen", "Panggil langsung: kripto.buat_kunci() atau crypto.generate_key()")
			}
			key := make([]byte, 32)
			if _, err := io.ReadFull(rand.Reader, key); err != nil {
				return newError(nil, fmt.Sprintf("kripto.buat_kunci: gagal menghasilkan keacakan: %v", err), "Periksa ketersediaan sumber entropi sistem.")
			}
			return &String{Value: base64.StdEncoding.EncodeToString(key)}
		},
	}
	definePublicFunction(env, "buat_kunci", buatKunciFn)
	definePublicFunction(env, "generate_key", buatKunciFn)

	// enkripsi / encrypt(data, kunci) -> Base64(nonce 12B || ciphertext || tag 16B)
	enkripsiFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "kripto.enkripsi mengharapkan 2 argumen: data dan kunci", "Contoh: kripto.enkripsi(\"teks rahasia\", kunci)")
			}
			dataBytes, errObj := getCryptoDataBytes(args[0], "kripto.enkripsi", "data ")
			if errObj != nil {
				return errObj
			}
			keyBytes, errKey := parseSymmetricKey(args[1], "kripto.enkripsi")
			if errKey != nil {
				return errKey
			}

			block, err := aes.NewCipher(keyBytes)
			if err != nil {
				return newError(nil, fmt.Sprintf("kripto.enkripsi: gagal membuat AES cipher: %v", err), "")
			}

			gcm, err := cipher.NewGCM(block)
			if err != nil {
				return newError(nil, fmt.Sprintf("kripto.enkripsi: gagal membuat mode GCM: %v", err), "")
			}

			nonce := make([]byte, gcm.NonceSize())
			if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
				return newError(nil, fmt.Sprintf("kripto.enkripsi: gagal membuat nonce: %v", err), "Pastikan sumber entropi sistem tersedia.")
			}

			ciphertextAndTag := gcm.Seal(nil, nonce, dataBytes, nil)
			payload := make([]byte, len(nonce)+len(ciphertextAndTag))
			copy(payload, nonce)
			copy(payload[len(nonce):], ciphertextAndTag)

			return &String{Value: base64.StdEncoding.EncodeToString(payload)}
		},
	}
	definePublicFunction(env, "enkripsi", enkripsiFn)
	definePublicFunction(env, "encrypt", enkripsiFn)

	// dekripsi / decrypt(ciphertext, kunci) -> String (plaintext)
	dekripsiFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "kripto.dekripsi mengharapkan 2 argumen: ciphertext dan kunci", "Contoh: kripto.dekripsi(terenkripsi, kunci)")
			}

			var payload []byte
			switch v := args[0].(type) {
			case *String:
				str := strings.TrimSpace(v.Value)
				decoded, err := base64.StdEncoding.DecodeString(str)
				if err != nil {
					decoded, err = base64.URLEncoding.DecodeString(str)
					if err != nil {
						decoded, err = base64.RawStdEncoding.DecodeString(str)
					}
				}
				if err != nil {
					return newError(nil, "kripto.dekripsi: format ciphertext tidak valid (harus berupa Base64)", "Pastikan ciphertext merupakan string Base64 yang dihasilkan oleh kripto.enkripsi().")
				}
				payload = decoded
			case *Bytes:
				payload = v.Elements
			default:
				return newError(nil, fmt.Sprintf("kripto.dekripsi: argumen ciphertext harus berupa teks atau Bita, menerima %s", args[0].Type()), "Contoh: kripto.dekripsi(terenkripsi, kunci)")
			}

			if len(payload) < 28 {
				return newError(nil, "kripto.dekripsi: format data tersandi tidak valid (panjang data terlalu pendek)", "Data tersandi minimal harus memuat nonce dan tag otentikasi (28 bita).")
			}

			keyBytes, errKey := parseSymmetricKey(args[1], "kripto.dekripsi")
			if errKey != nil {
				return errKey
			}

			block, err := aes.NewCipher(keyBytes)
			if err != nil {
				return newError(nil, fmt.Sprintf("kripto.dekripsi: gagal membuat AES cipher: %v", err), "")
			}

			gcm, err := cipher.NewGCM(block)
			if err != nil {
				return newError(nil, fmt.Sprintf("kripto.dekripsi: gagal membuat mode GCM: %v", err), "")
			}

			nonceSize := gcm.NonceSize()
			nonce := payload[:nonceSize]
			ciphertextAndTag := payload[nonceSize:]

			plaintext, err := gcm.Open(nil, nonce, ciphertextAndTag, nil)
			if err != nil {
				return newError(nil, "kripto.dekripsi: dekripsi gagal (otentikasi gagal, data tersandi telah dimodifikasi atau kunci salah)", "Pastikan data tersandi tidak mengalami perubahan dan kunci enkripsi yang digunakan sesuai.")
			}

			return &String{Value: string(plaintext)}
		},
	}
	definePublicFunction(env, "dekripsi", dekripsiFn)
	definePublicFunction(env, "decrypt", dekripsiFn)

	// =============================================================
	// 6. Derivasi Kunci & Password Hashing (PBKDF2-HMAC-SHA256)
	// =============================================================

	// hash_sandi / hash_password(sandi) -> PHC string
	hashPasswordFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "kripto.hash_sandi mengharapkan 1 argumen: kata sandi (teks atau Bita)", "Contoh: kripto.hash_sandi(\"kata-sandi-rahasia\")")
			}
			pwBytes, errObj := getCryptoDataBytes(args[0], "kripto.hash_sandi", "sandi ")
			if errObj != nil {
				return errObj
			}
			if len(pwBytes) == 0 {
				return newError(nil, "kripto.hash_sandi: kata sandi tidak boleh kosong", "Masukkan kata sandi yang valid untuk di-hash.")
			}

			salt := make([]byte, 16)
			if _, err := io.ReadFull(rand.Reader, salt); err != nil {
				return newError(nil, fmt.Sprintf("kripto.hash_sandi: gagal menghasilkan salt acak: %v", err), "Periksa ketersediaan sumber entropi sistem operasi.")
			}

			const iterations = 100000
			const keyLength = 32

			derived := pbkdf2SHA256(pwBytes, salt, iterations, keyLength)
			saltB64 := base64.RawStdEncoding.EncodeToString(salt)
			hashB64 := base64.RawStdEncoding.EncodeToString(derived)

			phcString := fmt.Sprintf("$pbkdf2-sha256$i=%d$%s$%s", iterations, saltB64, hashB64)
			return &String{Value: phcString}
		},
	}
	definePublicFunction(env, "hash_sandi", hashPasswordFn)
	definePublicFunction(env, "hash_password", hashPasswordFn)

	// verifikasi_sandi / verify_password(sandi, hash) -> Boolean
	verifyPasswordFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "kripto.verifikasi_sandi mengharapkan 2 argumen: kata sandi dan hash", "Contoh: kripto.verifikasi_sandi(sandi, hash)")
			}
			pwBytes, errObj := getCryptoDataBytes(args[0], "kripto.verifikasi_sandi", "sandi ")
			if errObj != nil {
				return errObj
			}
			hashStrObj, ok := args[1].(*String)
			if !ok {
				return newError(nil, fmt.Sprintf("argumen hash kripto.verifikasi_sandi harus berupa teks, menerima %s", args[1].Type()), "Contoh: kripto.verifikasi_sandi(\"rahasia\", hashTeks)")
			}
			hashStr := strings.TrimSpace(hashStrObj.Value)

			// Validate PHC structure: $pbkdf2-sha256$i=<iter>$<salt>$<hash>
			if !strings.HasPrefix(hashStr, "$pbkdf2-sha256$") {
				return FALSE
			}
			parts := strings.Split(hashStr, "$")
			// parts[0] = "", parts[1] = "pbkdf2-sha256", parts[2] = params, parts[3] = salt, parts[4] = hash
			if len(parts) != 5 || parts[1] != "pbkdf2-sha256" {
				return FALSE
			}

			// Extract iteration count parameter
			iter := 0
			for _, param := range strings.Split(parts[2], ",") {
				param = strings.TrimSpace(param)
				if strings.HasPrefix(param, "i=") {
					parsedIter, err := strconv.Atoi(strings.TrimPrefix(param, "i="))
					if err != nil {
						return FALSE
					}
					iter = parsedIter
					break
				}
			}
			if iter <= 0 || iter > 10_000_000 {
				return FALSE
			}

			// Decode salt and stored hash
			salt, err := decodeBase64Flex(parts[3])
			if err != nil || len(salt) < 16 {
				return FALSE
			}
			storedHash, err := decodeBase64Flex(parts[4])
			if err != nil || len(storedHash) != 32 {
				return FALSE
			}

			// Derive candidate hash (canonical 32 bytes for PBKDF2-HMAC-SHA256)
			candidateHash := pbkdf2SHA256(pwBytes, salt, iter, 32)

			// Constant-time comparison
			matched := subtle.ConstantTimeCompare(candidateHash, storedHash) == 1
			return nativeBool(matched)
		},
	}
	definePublicFunction(env, "verifikasi_sandi", verifyPasswordFn)
	definePublicFunction(env, "verify_password", verifyPasswordFn)

	// turunkan_kunci / derive_key(sandi, salt, panjang, [format])
	deriveKeyFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 3 || len(args) > 4 {
				return newError(nil, "kripto.turunkan_kunci mengharapkan 3 atau 4 argumen: sandi, salt, panjang, dan opsional format", "Contoh: kripto.turunkan_kunci(\"sandi\", salt, 32)")
			}
			pwBytes, err1 := getCryptoDataBytes(args[0], "kripto.turunkan_kunci", "sandi ")
			if err1 != nil {
				return err1
			}
			if len(pwBytes) == 0 {
				return newError(nil, "kripto.turunkan_kunci: kata sandi tidak boleh kosong", "Gunakan kata sandi dengan entropi yang memadai.")
			}

			saltBytes, err2 := getCryptoDataBytes(args[1], "kripto.turunkan_kunci", "salt ")
			if err2 != nil {
				return err2
			}
			if len(saltBytes) == 0 {
				return newError(nil, "kripto.turunkan_kunci: salt tidak boleh kosong", "Gunakan salt acak unik (minimal 16 bita), contoh: kripto.acak_bita(16)")
			}

			lenObj, ok := args[2].(*Integer)
			if !ok {
				return newError(nil, fmt.Sprintf("argumen panjang kripto.turunkan_kunci harus berupa angka bulat, menerima %s", args[2].Type()), "Contoh: kripto.turunkan_kunci(sandi, salt, 32)")
			}
			if lenObj.Value <= 0 {
				return newError(nil, "panjang kunci kripto.turunkan_kunci harus lebih besar dari 0", "Gunakan panjang positif, misalnya 32 untuk kunci AES-256.")
			}
			if lenObj.Value > 1024*1024 {
				return newError(nil, "panjang kunci kripto.turunkan_kunci melebihi batas aman (1MB)", "")
			}
			keyLen := int(lenObj.Value)

			format := "bita"
			if len(args) == 4 {
				fmtObj, ok := args[3].(*String)
				if !ok {
					return newError(nil, fmt.Sprintf("argumen format kripto.turunkan_kunci harus berupa teks, menerima %s", args[3].Type()), "Pilihan format: \"bita\", \"heks\"/\"hex\", atau \"base64\"")
				}
				format = strings.ToLower(strings.TrimSpace(fmtObj.Value))
			}

			const iterations = 100000
			dk := pbkdf2SHA256(pwBytes, saltBytes, iterations, keyLen)

			switch format {
			case "bita", "bytes", "raw":
				return &Bytes{Elements: dk}
			case "heks", "hex":
				return &String{Value: hex.EncodeToString(dk)}
			case "base64", "b64":
				return &String{Value: base64.StdEncoding.EncodeToString(dk)}
			default:
				return newError(nil, fmt.Sprintf("format keluaran '%s' tidak dikenal pada kripto.turunkan_kunci", format), "Pilihan format yang didukung: \"bita\", \"heks\"/\"hex\", \"base64\".")
			}
		},
	}
	definePublicFunction(env, "turunkan_kunci", deriveKeyFn)
	definePublicFunction(env, "derive_key", deriveKeyFn)

	// =============================================================
	// 7. Tanda Tangan Digital Modern Asimetris (Ed25519 - RFC 8032)
	// =============================================================

	// buat_pasangan_kunci() / generate_keypair()
	buatPasanganKunciFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 0 {
				return newError(nil, "kripto.buat_pasangan_kunci tidak memerlukan argumen", "Contoh: buat pasangan = kripto.buat_pasangan_kunci()")
			}
			pub, priv, err := ed25519.GenerateKey(rand.Reader)
			if err != nil {
				return newError(nil, fmt.Sprintf("kripto.buat_pasangan_kunci: gagal menghasilkan pasangan kunci: %v", err), "Pastikan sumber entropi sistem tersedia.")
			}
			pubHex := hex.EncodeToString(pub)
			privHex := hex.EncodeToString(priv.Seed())
			order := []string{"kunci_publik", "kunci_privat", "public_key", "private_key"}
			pairs := map[string]Object{
				"kunci_publik": &String{Value: pubHex},
				"kunci_privat": &String{Value: privHex},
				"public_key":   &String{Value: pubHex},
				"private_key":  &String{Value: privHex},
			}
			return &Map{Pairs: pairs, Order: order}
		},
	}
	definePublicFunction(env, "buat_pasangan_kunci", buatPasanganKunciFn)
	definePublicFunction(env, "generate_keypair", buatPasanganKunciFn)

	// tanda_tangan(data, kunci_privat) / sign(data, private_key)
	tandaTanganFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "kripto.tanda_tangan mengharapkan 2 argumen: data dan kunci_privat", "Contoh: kripto.tanda_tangan(\"pesan\", pasangan.kunci_privat)")
			}
			data, errObj := getCryptoDataBytes(args[0], "kripto.tanda_tangan", "pertama (data) ")
			if errObj != nil {
				return errObj
			}
			privKey, errPriv := parseEd25519PrivateKey(args[1], "kripto.tanda_tangan")
			if errPriv != nil {
				return errPriv
			}
			sig := ed25519.Sign(privKey, data)
			return &String{Value: hex.EncodeToString(sig)}
		},
	}
	definePublicFunction(env, "tanda_tangan", tandaTanganFn)
	definePublicFunction(env, "sign", tandaTanganFn)

	// verifikasi_tanda_tangan(data, signature, kunci_publik) / verify(data, signature, public_key)
	verifikasiTandaTanganFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 3 {
				return newError(nil, "kripto.verifikasi_tanda_tangan mengharapkan 3 argumen: data, signature, dan kunci_publik", "Contoh: kripto.verifikasi_tanda_tangan(data, signature, kunci_publik)")
			}
			data, errObj := getCryptoDataBytes(args[0], "kripto.verifikasi_tanda_tangan", "pertama (data) ")
			if errObj != nil {
				return errObj
			}
			sig := parseEd25519Signature(args[1])
			if sig == nil {
				return FALSE
			}
			pubKey := parseEd25519PublicKey(args[2])
			if pubKey == nil {
				return FALSE
			}
			valid := ed25519.Verify(pubKey, data, sig)
			return nativeBool(valid)
		},
	}
	definePublicFunction(env, "verifikasi_tanda_tangan", verifikasiTandaTanganFn)
	definePublicFunction(env, "verifikasi_ttd", verifikasiTandaTanganFn)
	definePublicFunction(env, "verify", verifikasiTandaTanganFn)
	definePublicFunction(env, "verify_signature", verifikasiTandaTanganFn)
}
