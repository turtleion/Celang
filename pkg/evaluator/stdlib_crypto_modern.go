package evaluator

import (
	"bytes"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/hkdf"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha3"
	"crypto/subtle"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/zeebo/blake3"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/blake2s"
	"golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/crypto/poly1305"
)

// registerModernCrypto mendaftarkan primitif kriptografi modern tingkat lanjut (Pilar 1 - 4).
func registerModernCrypto(env *Environment) {
	// =========================================================================
	// PILAR 1: HASH FUNCTIONS & MESSAGE AUTHENTICATION CODES (MAC)
	// =========================================================================

	// sha3_256(data)
	sha3_256Fn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "kripto.sha3_256 membutuhkan 1 argumen: (data teks atau Bita)", "Contoh: kripto.sha3_256(\"halo dunia\")")
			}
			b, errObj := getCryptoDataBytes(args[0], "kripto.sha3_256", "")
			if errObj != nil {
				return errObj
			}
			h := sha3.Sum256(b)
			return &String{Value: hex.EncodeToString(h[:])}
		},
	}
	definePublicFunction(env, "sha3_256", sha3_256Fn)

	// sha3_512(data)
	sha3_512Fn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "kripto.sha3_512 membutuhkan 1 argumen: (data teks atau Bita)", "Contoh: kripto.sha3_512(\"halo dunia\")")
			}
			b, errObj := getCryptoDataBytes(args[0], "kripto.sha3_512", "")
			if errObj != nil {
				return errObj
			}
			h := sha3.Sum512(b)
			return &String{Value: hex.EncodeToString(h[:])}
		},
	}
	definePublicFunction(env, "sha3_512", sha3_512Fn)

	// shake128(data, panjang_output_bita)
	shake128Fn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "kripto.shake128 membutuhkan 2 argumen: (data, panjang_output_bita)", "Contoh: kripto.shake128(\"pesan\", 32)")
			}
			b, errObj := getCryptoDataBytes(args[0], "kripto.shake128", "data ")
			if errObj != nil {
				return errObj
			}
			lenInt, ok := args[1].(*Integer)
			if !ok || lenInt.Value <= 0 || lenInt.Value > 1024*1024 {
				return newError(nil, "kripto.shake128: panjang output harus berupa bilangan bulat positif (maks 1MB)", "")
			}
			out := make([]byte, int(lenInt.Value))
			s := sha3.NewSHAKE128()
			s.Write(b)
			s.Read(out)
			return &String{Value: hex.EncodeToString(out)}
		},
	}
	definePublicFunction(env, "shake128", shake128Fn)

	// shake256(data, panjang_output_bita)
	shake256Fn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "kripto.shake256 membutuhkan 2 argumen: (data, panjang_output_bita)", "Contoh: kripto.shake256(\"pesan\", 64)")
			}
			b, errObj := getCryptoDataBytes(args[0], "kripto.shake256", "data ")
			if errObj != nil {
				return errObj
			}
			lenInt, ok := args[1].(*Integer)
			if !ok || lenInt.Value <= 0 || lenInt.Value > 1024*1024 {
				return newError(nil, "kripto.shake256: panjang output harus berupa bilangan bulat positif (maks 1MB)", "")
			}
			out := make([]byte, int(lenInt.Value))
			s := sha3.NewSHAKE256()
			s.Write(b)
			s.Read(out)
			return &String{Value: hex.EncodeToString(out)}
		},
	}
	definePublicFunction(env, "shake256", shake256Fn)

	// blake2b(data, [kunci])
	blake2bFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 1 || len(args) > 2 {
				return newError(nil, "kripto.blake2b membutuhkan 1 atau 2 argumen: (data, [kunci])", "Contoh: kripto.blake2b(\"pesan\") atau kripto.blake2b(\"pesan\", kunci)")
			}
			b, errObj := getCryptoDataBytes(args[0], "kripto.blake2b", "data ")
			if errObj != nil {
				return errObj
			}
			if len(args) == 2 && args[1] != NULL {
				key, errKey := getCryptoDataBytes(args[1], "kripto.blake2b", "kunci ")
				if errKey != nil {
					return errKey
				}
				h, err := blake2b.New512(key)
				if err != nil {
					return newError(nil, fmt.Sprintf("kripto.blake2b gagal: %v", err), "Kunci BLAKE2b maksimal 64 bita.")
				}
				h.Write(b)
				return &String{Value: hex.EncodeToString(h.Sum(nil))}
			}
			h := blake2b.Sum512(b)
			return &String{Value: hex.EncodeToString(h[:])}
		},
	}
	definePublicFunction(env, "blake2b", blake2bFn)

	// blake2s(data, [kunci])
	blake2sFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 1 || len(args) > 2 {
				return newError(nil, "kripto.blake2s membutuhkan 1 atau 2 argumen: (data, [kunci])", "Contoh: kripto.blake2s(\"pesan\")")
			}
			b, errObj := getCryptoDataBytes(args[0], "kripto.blake2s", "data ")
			if errObj != nil {
				return errObj
			}
			if len(args) == 2 && args[1] != NULL {
				key, errKey := getCryptoDataBytes(args[1], "kripto.blake2s", "kunci ")
				if errKey != nil {
					return errKey
				}
				h, err := blake2s.New256(key)
				if err != nil {
					return newError(nil, fmt.Sprintf("kripto.blake2s gagal: %v", err), "Kunci BLAKE2s maksimal 32 bita.")
				}
				h.Write(b)
				return &String{Value: hex.EncodeToString(h.Sum(nil))}
			}
			h := blake2s.Sum256(b)
			return &String{Value: hex.EncodeToString(h[:])}
		},
	}
	definePublicFunction(env, "blake2s", blake2sFn)

	// blake3(data, [kunci])
	blake3Fn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 1 || len(args) > 2 {
				return newError(nil, "kripto.blake3 membutuhkan 1 atau 2 argumen: (data, [kunci])", "Contoh: kripto.blake3(\"pesan\")")
			}
			b, errObj := getCryptoDataBytes(args[0], "kripto.blake3", "data ")
			if errObj != nil {
				return errObj
			}
			if len(args) == 2 && args[1] != NULL {
				key, errKey := getCryptoDataBytes(args[1], "kripto.blake3", "kunci ")
				if errKey != nil {
					return errKey
				}
				if len(key) != 32 {
					return newError(nil, fmt.Sprintf("kripto.blake3: kunci harus sepanjang 32 bita, menerima %d bita", len(key)), "")
				}
				h, err := blake3.NewKeyed(key)
				if err != nil {
					return newError(nil, fmt.Sprintf("kripto.blake3: kunci tidak valid: %v", err), "")
				}
				h.Write(b)
				return &String{Value: hex.EncodeToString(h.Sum(nil))}
			}
			h := blake3.Sum256(b)
			return &String{Value: hex.EncodeToString(h[:])}
		},
	}
	definePublicFunction(env, "blake3", blake3Fn)

	// poly1305_buat_mac(pesan, kunci) / poly1305_auth
	poly1305AuthFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "kripto.poly1305_buat_mac membutuhkan 2 argumen: (pesan, kunci_32_bita)", "Contoh: kripto.poly1305_buat_mac(\"pesan\", kunci)")
			}
			msg, errObj := getCryptoDataBytes(args[0], "kripto.poly1305_buat_mac", "pesan ")
			if errObj != nil {
				return errObj
			}
			keyBytes, errKey := parseSymmetricKey(args[1], "kripto.poly1305_buat_mac")
			if errKey != nil {
				return errKey
			}
			var keyArr [32]byte
			copy(keyArr[:], keyBytes)
			var out [16]byte
			poly1305.Sum(&out, msg, &keyArr)
			return &String{Value: hex.EncodeToString(out[:])}
		},
	}
	definePublicFunction(env, "poly1305_buat_mac", poly1305AuthFn)
	definePublicFunction(env, "poly1305_auth", poly1305AuthFn)
	definePublicFunction(env, "poly1305_mac", poly1305AuthFn)

	// poly1305_verifikasi(pesan, tag, kunci) / poly1305_verify
	poly1305VerifyFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 3 {
				return newError(nil, "kripto.poly1305_verifikasi membutuhkan 3 argumen: (pesan, tag, kunci)", "Contoh: kripto.poly1305_verifikasi(pesan, tag, kunci)")
			}
			msg, errObj := getCryptoDataBytes(args[0], "kripto.poly1305_verifikasi", "pesan ")
			if errObj != nil {
				return errObj
			}
			tagBytes, errTag := parse16Bytes(args[1], "kripto.poly1305_verifikasi")
			if errTag != nil {
				return FALSE
			}
			keyBytes, errKey := parseSymmetricKey(args[2], "kripto.poly1305_verifikasi")
			if errKey != nil {
				return errKey
			}
			var tagArr [16]byte
			copy(tagArr[:], tagBytes)
			var keyArr [32]byte
			copy(keyArr[:], keyBytes)
			valid := poly1305.Verify(&tagArr, msg, &keyArr)
			return nativeBool(valid)
		},
	}
	definePublicFunction(env, "poly1305_verifikasi", poly1305VerifyFn)
	definePublicFunction(env, "poly1305_verify", poly1305VerifyFn)

	// =========================================================================
	// PILAR 2: KEY DERIVATION (KDF) & ADVANCED PASSWORD HASHING
	// =========================================================================

	// hkdf(secret, salt, info, panjang) / hkdf
	hkdfFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 4 {
				return newError(nil, "kripto.hkdf membutuhkan 4 argumen: (secret, salt, info, panjang)", "Contoh: kripto.hkdf(rahasia, salt, \"aplikasi-sesi\", 32)")
			}
			sec, err1 := getCryptoDataBytes(args[0], "kripto.hkdf", "secret ")
			if err1 != nil {
				return err1
			}
			var salt []byte
			if args[1] != NULL {
				var err2 Object
				salt, err2 = getCryptoDataBytes(args[1], "kripto.hkdf", "salt ")
				if err2 != nil {
					return err2
				}
			}
			var infoStr string
			switch v := args[2].(type) {
			case *String:
				infoStr = v.Value
			case *Bytes:
				infoStr = string(v.Elements)
			case *Null:
				infoStr = ""
			default:
				return newError(nil, fmt.Sprintf("kripto.hkdf: argumen info harus berupa teks atau Bita, menerima %s", args[2].Type()), "")
			}

			lenInt, ok := args[3].(*Integer)
			if !ok || lenInt.Value <= 0 || lenInt.Value > 255*32 {
				return newError(nil, "kripto.hkdf: panjang kunci harus antara 1 dan 8160 bita (255 * 32 bita untuk SHA-256 sesuai RFC 5869)", "")
			}
			key, err := hkdf.Key(sha256.New, sec, salt, infoStr, int(lenInt.Value))
			if err != nil {
				return newError(nil, fmt.Sprintf("kripto.hkdf gagal: %v", err), "")
			}
			return &Bytes{Elements: key}
		},
	}
	definePublicFunction(env, "hkdf", hkdfFn)

	// hkdf_ekstrak(secret, salt) / hkdf_extract
	hkdfExtractFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "kripto.hkdf_ekstrak membutuhkan 2 argumen: (secret, salt)", "")
			}
			sec, err1 := getCryptoDataBytes(args[0], "kripto.hkdf_ekstrak", "secret ")
			if err1 != nil {
				return err1
			}
			var salt []byte
			if args[1] != NULL {
				var err2 Object
				salt, err2 = getCryptoDataBytes(args[1], "kripto.hkdf_ekstrak", "salt ")
				if err2 != nil {
					return err2
				}
			}
			prk, err := hkdf.Extract(sha256.New, sec, salt)
			if err != nil {
				return newError(nil, fmt.Sprintf("kripto.hkdf_ekstrak gagal: %v", err), "")
			}
			return &Bytes{Elements: prk}
		},
	}
	definePublicFunction(env, "hkdf_ekstrak", hkdfExtractFn)
	definePublicFunction(env, "hkdf_extract", hkdfExtractFn)

	// hkdf_ekspansi(prk, info, panjang) / hkdf_expand
	hkdfExpandFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 3 {
				return newError(nil, "kripto.hkdf_ekspansi membutuhkan 3 argumen: (prk, info, panjang)", "")
			}
			prk, err1 := getCryptoDataBytes(args[0], "kripto.hkdf_ekspansi", "prk ")
			if err1 != nil {
				return err1
			}
			var infoStr string
			switch v := args[1].(type) {
			case *String:
				infoStr = v.Value
			case *Bytes:
				infoStr = string(v.Elements)
			case *Null:
				infoStr = ""
			default:
				return newError(nil, fmt.Sprintf("kripto.hkdf_ekspansi: argumen info harus berupa teks atau Bita, menerima %s", args[1].Type()), "")
			}
			lenInt, ok := args[2].(*Integer)
			if !ok || lenInt.Value <= 0 || lenInt.Value > 255*32 {
				return newError(nil, "kripto.hkdf_ekspansi: panjang kunci harus antara 1 dan 8160 bita (255 * 32 bita untuk SHA-256 sesuai RFC 5869)", "")
			}
			out, err := hkdf.Expand(sha256.New, prk, infoStr, int(lenInt.Value))
			if err != nil {
				return newError(nil, fmt.Sprintf("kripto.hkdf_ekspansi gagal: %v", err), "")
			}
			return &Bytes{Elements: out}
		},
	}
	definePublicFunction(env, "hkdf_ekspansi", hkdfExpandFn)
	definePublicFunction(env, "hkdf_expand", hkdfExpandFn)

	// hash_sandi_argon2(sandi, [opsi]) / hash_password_argon2
	hashArgon2Fn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 1 || len(args) > 2 {
				return newError(nil, "kripto.hash_sandi_argon2 membutuhkan minimal 1 argumen: (kata_sandi, [opsi])", "Contoh: kripto.hash_sandi_argon2(\"rahasia\")")
			}
			pwBytes, errObj := getCryptoDataBytes(args[0], "kripto.hash_sandi_argon2", "sandi ")
			if errObj != nil {
				return errObj
			}
			if len(pwBytes) == 0 {
				return newError(nil, "kripto.hash_sandi_argon2: kata sandi tidak boleh kosong", "")
			}

			time := uint32(3)
			memory := uint32(64 * 1024) // 64 MB
			threads := uint8(4)
			keyLen := uint32(32)

			if len(args) == 2 {
				m, ok := args[1].(*Map)
				if !ok {
					return newError(nil, fmt.Sprintf("kripto.hash_sandi_argon2: argumen kedua (opsi) harus berupa Peta (Map), menerima %s", args[1].Type()), "Contoh: kripto.hash_sandi_argon2(sandi, {\"waktu\": 3, \"memori\": 65536})")
				}
				if v, ok := m.Pairs["waktu"]; ok {
					if iv, ok := v.(*Integer); ok && iv.Value >= 1 && iv.Value <= 20 {
						time = uint32(iv.Value)
					} else {
						return newError(nil, "kripto.hash_sandi_argon2: opsi 'waktu' harus bernilai integer 1-20 iterasi", "")
					}
				}
				if v, ok := m.Pairs["memori"]; ok {
					if iv, ok := v.(*Integer); ok && iv.Value >= 1024 && iv.Value <= 1024*1024 {
						memory = uint32(iv.Value)
					} else {
						return newError(nil, "kripto.hash_sandi_argon2: opsi 'memori' harus antara 1024 KB (1MB) dan 1048576 KB (1GB)", "")
					}
				}
				if v, ok := m.Pairs["threads"]; ok {
					if iv, ok := v.(*Integer); ok && iv.Value >= 1 && iv.Value <= 64 {
						threads = uint8(iv.Value)
					} else {
						return newError(nil, "kripto.hash_sandi_argon2: opsi 'threads' harus antara 1 dan 64", "")
					}
				}
			}

			salt := make([]byte, 16)
			if _, err := io.ReadFull(rand.Reader, salt); err != nil {
				return newError(nil, fmt.Sprintf("kripto.hash_sandi_argon2: gagal membuat salt acak: %v", err), "")
			}

			derived := argon2.IDKey(pwBytes, salt, time, memory, threads, keyLen)
			saltB64 := base64.RawStdEncoding.EncodeToString(salt)
			hashB64 := base64.RawStdEncoding.EncodeToString(derived)

			phcString := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, memory, time, threads, saltB64, hashB64)
			return &String{Value: phcString}
		},
	}
	definePublicFunction(env, "hash_sandi_argon2", hashArgon2Fn)
	definePublicFunction(env, "hash_password_argon2", hashArgon2Fn)

	// verifikasi_sandi_argon2(sandi, hash) / verify_password_argon2
	verifyArgon2Fn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "kripto.verifikasi_sandi_argon2 membutuhkan 2 argumen: (sandi, hash_argon2)", "")
			}
			pwBytes, errObj := getCryptoDataBytes(args[0], "kripto.verifikasi_sandi_argon2", "sandi ")
			if errObj != nil {
				return errObj
			}
			hashStrObj, ok := args[1].(*String)
			if !ok {
				return newError(nil, "kripto.verifikasi_sandi_argon2: argumen kedua harus berupa teks hash", "")
			}
			hashStr := strings.TrimSpace(hashStrObj.Value)
			if !strings.HasPrefix(hashStr, "$argon2id$") {
				return FALSE
			}
			parts := strings.Split(hashStr, "$")
			// format: ["", "argon2id", "v=19", "m=65536,t=3,p=4", saltB64, hashB64]
			if len(parts) != 6 {
				return FALSE
			}

			var mem, time uint32
			var threads uint8
			for _, param := range strings.Split(parts[3], ",") {
				param = strings.TrimSpace(param)
				if strings.HasPrefix(param, "m=") {
					v, err := strconv.ParseUint(strings.TrimPrefix(param, "m="), 10, 32)
					if err != nil || v < 1024 || v > 256*1024 { // Batas aman verifikasi input luar: maks 256 MB
						return FALSE
					}
					mem = uint32(v)
				} else if strings.HasPrefix(param, "t=") {
					v, err := strconv.ParseUint(strings.TrimPrefix(param, "t="), 10, 32)
					if err != nil || v < 1 || v > 10 { // Batas aman verifikasi input luar: maks 10 iterasi
						return FALSE
					}
					time = uint32(v)
				} else if strings.HasPrefix(param, "p=") {
					v, err := strconv.ParseUint(strings.TrimPrefix(param, "p="), 10, 8)
					if err != nil || v < 1 || v > 64 { // Batas aman: maks 64 threads
						return FALSE
					}
					threads = uint8(v)
				}
			}
			if mem == 0 || time == 0 || threads == 0 {
				return FALSE
			}

			salt, err := decodeBase64Flex(parts[4])
			if err != nil || len(salt) < 8 || len(salt) > 64 {
				return FALSE
			}
			storedHash, err := decodeBase64Flex(parts[5])
			if err != nil || len(storedHash) < 16 || len(storedHash) > 128 {
				return FALSE
			}

			candidate := argon2.IDKey(pwBytes, salt, time, mem, threads, uint32(len(storedHash)))
			matched := subtle.ConstantTimeCompare(candidate, storedHash) == 1
			return nativeBool(matched)
		},
	}
	definePublicFunction(env, "verifikasi_sandi_argon2", verifyArgon2Fn)
	definePublicFunction(env, "verify_password_argon2", verifyArgon2Fn)

	// argon2id_turunkan_kunci(sandi, salt, time, mem, threads, panjang) / argon2id_derive_key
	argon2DeriveKeyFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 6 {
				return newError(nil, "kripto.argon2id_turunkan_kunci membutuhkan 6 argumen: (sandi, salt, waktu, memori, threads, panjang)", "")
			}
			pw, err1 := getCryptoDataBytes(args[0], "kripto.argon2id_turunkan_kunci", "sandi ")
			if err1 != nil {
				return err1
			}
			salt, err2 := getCryptoDataBytes(args[1], "kripto.argon2id_turunkan_kunci", "salt ")
			if err2 != nil {
				return err2
			}
			timeInt, ok := args[2].(*Integer)
			if !ok || timeInt.Value < 1 || timeInt.Value > 20 {
				return newError(nil, "kripto.argon2id_turunkan_kunci: waktu harus bernilai antara 1 dan 20 iterasi", "")
			}
			memInt, ok := args[3].(*Integer)
			if !ok || memInt.Value < 1024 || memInt.Value > 1024*1024 {
				return newError(nil, "kripto.argon2id_turunkan_kunci: memori harus antara 1024 KB (1MB) dan 1048576 KB (1GB)", "")
			}
			thrInt, ok := args[4].(*Integer)
			if !ok || thrInt.Value < 1 || thrInt.Value > 64 {
				return newError(nil, "kripto.argon2id_turunkan_kunci: threads harus antara 1 dan 64", "")
			}
			lenInt, ok := args[5].(*Integer)
			if !ok || lenInt.Value < 1 || lenInt.Value > 1024*1024 {
				return newError(nil, "kripto.argon2id_turunkan_kunci: panjang harus antara 1 dan 1048576 bita (1MB)", "")
			}

			dk := argon2.IDKey(pw, salt, uint32(timeInt.Value), uint32(memInt.Value), uint8(thrInt.Value), uint32(lenInt.Value))
			return &Bytes{Elements: dk}
		},
	}
	definePublicFunction(env, "argon2id_turunkan_kunci", argon2DeriveKeyFn)
	definePublicFunction(env, "argon2id_derive_key", argon2DeriveKeyFn)

	// =========================================================================
	// PILAR 3: ENKRIPSI SIMETRIS CHACHA20-POLY1305 & AES BLOCK MODES
	// =========================================================================

	// buat_kunci_chacha() / generate_chacha_key
	buatKunciChaChaFn := &Builtin{
		Fn: func(args ...Object) Object {
			key := make([]byte, 32)
			if _, err := io.ReadFull(rand.Reader, key); err != nil {
				return newError(nil, fmt.Sprintf("kripto.buat_kunci_chacha: gagal membuat keacakan: %v", err), "")
			}
			return &String{Value: base64.StdEncoding.EncodeToString(key)}
		},
	}
	definePublicFunction(env, "buat_kunci_chacha", buatKunciChaChaFn)
	definePublicFunction(env, "generate_chacha_key", buatKunciChaChaFn)

	// enkripsi_chacha(data, kunci, [aad]) / encrypt_chacha
	enkripsiChaChaFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 2 || len(args) > 3 {
				return newError(nil, "kripto.enkripsi_chacha membutuhkan 2 atau 3 argumen: (data, kunci, [data_tambahan_aad])", "Contoh: kripto.enkripsi_chacha(\"pesan\", kunci)")
			}
			data, errObj := getCryptoDataBytes(args[0], "kripto.enkripsi_chacha", "data ")
			if errObj != nil {
				return errObj
			}
			keyBytes, errKey := parseSymmetricKey(args[1], "kripto.enkripsi_chacha")
			if errKey != nil {
				return errKey
			}
			var aad []byte
			if len(args) == 3 && args[2] != NULL {
				var errAAD Object
				aad, errAAD = getCryptoDataBytes(args[2], "kripto.enkripsi_chacha", "aad ")
				if errAAD != nil {
					return errAAD
				}
			}

			// Menggunakan XChaCha20-Poly1305 (24-byte nonce) untuk ketahanan terhadap tabrakan nonce acak
			aead, err := chacha20poly1305.NewX(keyBytes)
			if err != nil {
				return newError(nil, fmt.Sprintf("kripto.enkripsi_chacha: gagal membuat cipher: %v", err), "")
			}
			nonce := make([]byte, aead.NonceSize())
			if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
				return newError(nil, fmt.Sprintf("kripto.enkripsi_chacha: gagal membuat nonce: %v", err), "")
			}

			sealed := aead.Seal(nil, nonce, data, aad)
			payload := make([]byte, len(nonce)+len(sealed))
			copy(payload, nonce)
			copy(payload[len(nonce):], sealed)

			return &String{Value: base64.StdEncoding.EncodeToString(payload)}
		},
	}
	definePublicFunction(env, "enkripsi_chacha", enkripsiChaChaFn)
	definePublicFunction(env, "encrypt_chacha", enkripsiChaChaFn)

	// dekripsi_chacha(ciphertext, kunci, [aad]) / decrypt_chacha
	dekripsiChaChaFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 2 || len(args) > 3 {
				return newError(nil, "kripto.dekripsi_chacha membutuhkan 2 atau 3 argumen: (ciphertext, kunci, [aad])", "")
			}
			var payload []byte
			switch v := args[0].(type) {
			case *String:
				decoded, err := decodeBase64Flex(v.Value)
				if err != nil {
					return newError(nil, "kripto.dekripsi_chacha: format ciphertext tidak valid (harus Base64)", "")
				}
				payload = decoded
			case *Bytes:
				payload = v.Elements
			default:
				return newError(nil, "kripto.dekripsi_chacha: ciphertext harus berupa teks atau Bita", "")
			}

			keyBytes, errKey := parseSymmetricKey(args[1], "kripto.dekripsi_chacha")
			if errKey != nil {
				return errKey
			}

			var aad []byte
			if len(args) == 3 && args[2] != NULL {
				var errAAD Object
				aad, errAAD = getCryptoDataBytes(args[2], "kripto.dekripsi_chacha", "aad ")
				if errAAD != nil {
					return errAAD
				}
			}

			if len(payload) < chacha20poly1305.NonceSize+16 {
				return newError(nil, "kripto.dekripsi_chacha: ukuran data tersandi terlalu pendek", "")
			}

			var plaintext []byte
			decrypted := false

			// 1. Coba dekripsi dengan XChaCha20-Poly1305 (24-byte nonce) jika ukuran payload mencukupi
			if len(payload) >= chacha20poly1305.NonceSizeX+16 {
				if aeadX, err := chacha20poly1305.NewX(keyBytes); err == nil {
					nonceX := payload[:chacha20poly1305.NonceSizeX]
					ctAndTagX := payload[chacha20poly1305.NonceSizeX:]
					if pt, err := aeadX.Open(nil, nonceX, ctAndTagX, aad); err == nil {
						plaintext = pt
						decrypted = true
					}
				}
			}

			// 2. Kompatibilitas mundur: coba ChaCha20-Poly1305 standar (12-byte nonce) jika belum berhasil
			if !decrypted && len(payload) >= chacha20poly1305.NonceSize+16 {
				if aeadStd, err := chacha20poly1305.New(keyBytes); err == nil {
					nonceStd := payload[:chacha20poly1305.NonceSize]
					ctAndTagStd := payload[chacha20poly1305.NonceSize:]
					if pt, err := aeadStd.Open(nil, nonceStd, ctAndTagStd, aad); err == nil {
						plaintext = pt
						decrypted = true
					}
				}
			}

			if !decrypted {
				return newError(nil, "kripto.dekripsi_chacha: dekripsi gagal (otentikasi gagal atau kunci salah)", "")
			}

			return &String{Value: string(plaintext)}
		},
	}
	definePublicFunction(env, "dekripsi_chacha", dekripsiChaChaFn)
	definePublicFunction(env, "decrypt_chacha", dekripsiChaChaFn)

	// aes_ctr(data, kunci, iv) / aes_ctr
	aesCtrFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 3 {
				return newError(nil, "kripto.aes_ctr membutuhkan 3 argumen: (data, kunci, iv_16_bita)", "Contoh: kripto.aes_ctr(data, kunci, iv)")
			}
			data, err1 := getCryptoDataBytes(args[0], "kripto.aes_ctr", "data ")
			if err1 != nil {
				return err1
			}
			keyBytes, err2 := parseAESKeyAny(args[1], "kripto.aes_ctr")
			if err2 != nil {
				return err2
			}
			ivBytes, err3 := parse16Bytes(args[2], "kripto.aes_ctr")
			if err3 != nil {
				return err3
			}

			block, err := aes.NewCipher(keyBytes)
			if err != nil {
				return newError(nil, fmt.Sprintf("kripto.aes_ctr: gagal membuat AES block: %v", err), "")
			}
			stream := cipher.NewCTR(block, ivBytes)
			out := make([]byte, len(data))
			stream.XORKeyStream(out, data)
			return &Bytes{Elements: out}
		},
	}
	definePublicFunction(env, "aes_ctr", aesCtrFn)

	// aes_cbc_enkripsi(data, kunci, iv) / aes_cbc_encrypt
	aesCbcEncFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 3 {
				return newError(nil, "kripto.aes_cbc_enkripsi membutuhkan 3 argumen: (data, kunci, iv_16_bita)", "")
			}
			data, err1 := getCryptoDataBytes(args[0], "kripto.aes_cbc_enkripsi", "data ")
			if err1 != nil {
				return err1
			}
			keyBytes, err2 := parseAESKeyAny(args[1], "kripto.aes_cbc_enkripsi")
			if err2 != nil {
				return err2
			}
			ivBytes, err3 := parse16Bytes(args[2], "kripto.aes_cbc_enkripsi")
			if err3 != nil {
				return err3
			}

			block, err := aes.NewCipher(keyBytes)
			if err != nil {
				return newError(nil, fmt.Sprintf("kripto.aes_cbc_enkripsi: gagal membuat cipher: %v", err), "")
			}

			// PKCS#7 Padding (salin ke slice baru untuk mencegah aliasing pada backing array)
			padLen := aes.BlockSize - (len(data) % aes.BlockSize)
			padded := make([]byte, len(data)+padLen)
			copy(padded, data)
			copy(padded[len(data):], bytes.Repeat([]byte{byte(padLen)}, padLen))

			enc := cipher.NewCBCEncrypter(block, ivBytes)
			out := make([]byte, len(padded))
			enc.CryptBlocks(out, padded)
			return &Bytes{Elements: out}
		},
	}
	definePublicFunction(env, "aes_cbc_enkripsi", aesCbcEncFn)
	definePublicFunction(env, "aes_cbc_encrypt", aesCbcEncFn)

	// aes_cbc_dekripsi(ciphertext, kunci, iv) / aes_cbc_decrypt
	aesCbcDecFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 3 {
				return newError(nil, "kripto.aes_cbc_dekripsi membutuhkan 3 argumen: (ciphertext, kunci, iv_16_bita)", "")
			}
			data, err1 := getCryptoDataBytes(args[0], "kripto.aes_cbc_dekripsi", "ciphertext ")
			if err1 != nil {
				return err1
			}
			if len(data) == 0 || len(data)%aes.BlockSize != 0 {
				return newError(nil, "kripto.aes_cbc_dekripsi: panjang data tersandi harus kelipatan 16 bita", "")
			}
			keyBytes, err2 := parseAESKeyAny(args[1], "kripto.aes_cbc_dekripsi")
			if err2 != nil {
				return err2
			}
			ivBytes, err3 := parse16Bytes(args[2], "kripto.aes_cbc_dekripsi")
			if err3 != nil {
				return err3
			}

			block, err := aes.NewCipher(keyBytes)
			if err != nil {
				return newError(nil, fmt.Sprintf("kripto.aes_cbc_dekripsi: gagal membuat cipher: %v", err), "")
			}

			dec := cipher.NewCBCDecrypter(block, ivBytes)
			out := make([]byte, len(data))
			dec.CryptBlocks(out, data)

			unpadded, err := pkcs7UnpadConstantTime(out, aes.BlockSize)
			if err != nil {
				return newError(nil, "kripto.aes_cbc_dekripsi: dekripsi gagal (padding tidak valid atau kunci/IV salah)", "")
			}
			return &Bytes{Elements: unpadded}
		},
	}
	definePublicFunction(env, "aes_cbc_dekripsi", aesCbcDecFn)
	definePublicFunction(env, "aes_cbc_decrypt", aesCbcDecFn)

	// =========================================================================
	// PILAR 4: ASIMETRIS (X25519 & RSA)
	// =========================================================================

	// buat_kunci_x25519() / generate_x25519_keypair
	buatKunciX25519Fn := &Builtin{
		Fn: func(args ...Object) Object {
			curve := ecdh.X25519()
			priv, err := curve.GenerateKey(rand.Reader)
			if err != nil {
				return newError(nil, fmt.Sprintf("kripto.buat_kunci_x25519 gagal: %v", err), "")
			}
			pub := priv.PublicKey()

			pubHex := hex.EncodeToString(pub.Bytes())
			privHex := hex.EncodeToString(priv.Bytes())
			pairs := map[string]Object{
				"kunci_publik": &String{Value: pubHex},
				"public_key":   &String{Value: pubHex},
				"kunci_privat": &String{Value: privHex},
				"private_key":  &String{Value: privHex},
				"kurva":        &String{Value: "X25519"},
				"curve":        &String{Value: "X25519"},
			}
			order := []string{"kunci_publik", "public_key", "kunci_privat", "private_key", "kurva", "curve"}
			return &Map{Pairs: pairs, Order: order}
		},
	}
	definePublicFunction(env, "buat_kunci_x25519", buatKunciX25519Fn)
	definePublicFunction(env, "generate_x25519_keypair", buatKunciX25519Fn)

	// x25519_hitung_rahasia(kunci_privat, kunci_publik_lawan) / x25519_shared_secret
	x25519SharedSecretFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "kripto.x25519_hitung_rahasia membutuhkan 2 argumen: (kunci_privat, kunci_publik_lawan)", "")
			}
			privBytes, err1 := parse32BytesFlex(args[0], "kripto.x25519_hitung_rahasia (kunci_privat)")
			if err1 != nil {
				return err1
			}
			pubBytes, err2 := parse32BytesFlex(args[1], "kripto.x25519_hitung_rahasia (kunci_publik_lawan)")
			if err2 != nil {
				return err2
			}

			curve := ecdh.X25519()
			privKey, err := curve.NewPrivateKey(privBytes)
			if err != nil {
				return newError(nil, fmt.Sprintf("kripto.x25519_hitung_rahasia: kunci privat tidak valid: %v", err), "")
			}
			pubKey, err := curve.NewPublicKey(pubBytes)
			if err != nil {
				return newError(nil, fmt.Sprintf("kripto.x25519_hitung_rahasia: kunci publik lawan tidak valid: %v", err), "")
			}

			sharedSecret, err := privKey.ECDH(pubKey)
			if err != nil {
				return newError(nil, fmt.Sprintf("kripto.x25519_hitung_rahasia gagal: %v", err), "")
			}
			return &Bytes{Elements: sharedSecret}
		},
	}
	definePublicFunction(env, "x25519_hitung_rahasia", x25519SharedSecretFn)
	definePublicFunction(env, "x25519_shared_secret", x25519SharedSecretFn)
	definePublicFunction(env, "x25519_compute_secret", x25519SharedSecretFn)

	// rsa_buat_kunci([bits]) / rsa_generate_keypair
	rsaBuatKunciFn := &Builtin{
		Fn: func(args ...Object) Object {
			bits := 2048
			if len(args) >= 1 {
				intObj, ok := args[0].(*Integer)
				if !ok {
					return newError(nil, fmt.Sprintf("kripto.rsa_buat_kunci: argumen ukuran bit harus berupa bilangan bulat (Integer), menerima %s", args[0].Type()), "Contoh: kripto.rsa_buat_kunci(2048)")
				}
				bits = int(intObj.Value)
			}
			if bits < 2048 || bits > 8192 {
				return newError(nil, "kripto.rsa_buat_kunci: ukuran bit harus antara 2048 dan 8192 (kunci 1024-bit sudah tidak aman)", "Gunakan ukuran minimal 2048 atau 4096 bit.")
			}
			privKey, err := rsa.GenerateKey(rand.Reader, bits)
			if err != nil {
				return newError(nil, fmt.Sprintf("kripto.rsa_buat_kunci gagal: %v", err), "")
			}

			privDER := x509.MarshalPKCS1PrivateKey(privKey)
			privPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: privDER})

			pubDER, err := x509.MarshalPKIXPublicKey(&privKey.PublicKey)
			if err != nil {
				return newError(nil, fmt.Sprintf("kripto.rsa_buat_kunci: gagal marshal publik DER: %v", err), "")
			}
			pubPEM := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubDER})

			pairs := map[string]Object{
				"kunci_publik": &String{Value: string(pubPEM)},
				"public_key":   &String{Value: string(pubPEM)},
				"kunci_privat": &String{Value: string(privPEM)},
				"private_key":  &String{Value: string(privPEM)},
				"bits":         &Integer{Value: int64(bits)},
			}
			order := []string{"kunci_publik", "public_key", "kunci_privat", "private_key", "bits"}
			return &Map{Pairs: pairs, Order: order}
		},
	}
	definePublicFunction(env, "rsa_buat_kunci", rsaBuatKunciFn)
	definePublicFunction(env, "rsa_generate_keypair", rsaBuatKunciFn)
	definePublicFunction(env, "rsa_generate_key", rsaBuatKunciFn)

	// rsa_enkripsi_oaep(data, kunci_publik) / rsa_encrypt_oaep
	rsaEnkripsiFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "kripto.rsa_enkripsi_oaep membutuhkan 2 argumen: (data, kunci_publik_pem)", "")
			}
			data, errObj := getCryptoDataBytes(args[0], "kripto.rsa_enkripsi_oaep", "data ")
			if errObj != nil {
				return errObj
			}
			pubKey, errKey := parseRSAPublicKey(args[1], "kripto.rsa_enkripsi_oaep")
			if errKey != nil {
				return errKey
			}

			ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pubKey, data, nil)
			if err != nil {
				return newError(nil, fmt.Sprintf("kripto.rsa_enkripsi_oaep gagal: %v", err), "Ukuran data mungkin terlalu besar untuk panjang kunci RSA.")
			}
			return &String{Value: base64.StdEncoding.EncodeToString(ciphertext)}
		},
	}
	definePublicFunction(env, "rsa_enkripsi_oaep", rsaEnkripsiFn)
	definePublicFunction(env, "rsa_encrypt_oaep", rsaEnkripsiFn)

	// rsa_dekripsi_oaep(ciphertext, kunci_privat) / rsa_decrypt_oaep
	rsaDekripsiFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "kripto.rsa_dekripsi_oaep membutuhkan 2 argumen: (ciphertext, kunci_privat_pem)", "")
			}
			var cipherBytes []byte
			switch v := args[0].(type) {
			case *String:
				decoded, err := decodeBase64Flex(v.Value)
				if err != nil {
					return newError(nil, "kripto.rsa_dekripsi_oaep: format ciphertext harus Base64", "")
				}
				cipherBytes = decoded
			case *Bytes:
				cipherBytes = v.Elements
			default:
				return newError(nil, "kripto.rsa_dekripsi_oaep: ciphertext harus berupa teks atau Bita", "")
			}

			privKey, errKey := parseRSAPrivateKey(args[1], "kripto.rsa_dekripsi_oaep")
			if errKey != nil {
				return errKey
			}

			plaintext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privKey, cipherBytes, nil)
			if err != nil {
				return newError(nil, fmt.Sprintf("kripto.rsa_dekripsi_oaep gagal: %v", err), "Kunci privat tidak cocok atau ciphertext rusak.")
			}
			return &String{Value: string(plaintext)}
		},
	}
	definePublicFunction(env, "rsa_dekripsi_oaep", rsaDekripsiFn)
	definePublicFunction(env, "rsa_decrypt_oaep", rsaDekripsiFn)

	// rsa_tanda_tangan_pss(data, kunci_privat) / rsa_sign_pss
	rsaSignFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "kripto.rsa_tanda_tangan_pss membutuhkan 2 argumen: (data, kunci_privat_pem)", "")
			}
			data, errObj := getCryptoDataBytes(args[0], "kripto.rsa_tanda_tangan_pss", "data ")
			if errObj != nil {
				return errObj
			}
			privKey, errKey := parseRSAPrivateKey(args[1], "kripto.rsa_tanda_tangan_pss")
			if errKey != nil {
				return errKey
			}

			hashed := sha256.Sum256(data)
			pssOpts := &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthEqualsHash, Hash: crypto.SHA256}
			sig, err := rsa.SignPSS(rand.Reader, privKey, crypto.SHA256, hashed[:], pssOpts)
			if err != nil {
				return newError(nil, fmt.Sprintf("kripto.rsa_tanda_tangan_pss gagal: %v", err), "")
			}
			return &String{Value: hex.EncodeToString(sig)}
		},
	}
	definePublicFunction(env, "rsa_tanda_tangan_pss", rsaSignFn)
	definePublicFunction(env, "rsa_sign_pss", rsaSignFn)

	// rsa_verifikasi_pss(data, signature, kunci_publik) / rsa_verify_pss
	rsaVerifyFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 3 {
				return newError(nil, "kripto.rsa_verifikasi_pss membutuhkan 3 argumen: (data, signature, kunci_publik_pem)", "")
			}
			data, errObj := getCryptoDataBytes(args[0], "kripto.rsa_verifikasi_pss", "data ")
			if errObj != nil {
				return errObj
			}
			pubKey, errKey := parseRSAPublicKey(args[2], "kripto.rsa_verifikasi_pss")
			if errKey != nil {
				return errKey
			}
			expectedSigLen := pubKey.Size()

			var sigBytes []byte
			switch v := args[1].(type) {
			case *String:
				s := strings.TrimSpace(v.Value)
				if len(s) == expectedSigLen*2 {
					if decoded, err := hex.DecodeString(s); err == nil {
						sigBytes = decoded
					}
				}
				if sigBytes == nil {
					if decoded, err := decodeBase64Flex(s); err == nil && len(decoded) == expectedSigLen {
						sigBytes = decoded
					}
				}
				if sigBytes == nil {
					if decoded, err := hex.DecodeString(s); err == nil {
						sigBytes = decoded
					} else if decoded, err := decodeBase64Flex(s); err == nil {
						sigBytes = decoded
					} else {
						return FALSE
					}
				}
			case *Bytes:
				sigBytes = v.Elements
			default:
				return FALSE
			}

			hashed := sha256.Sum256(data)
			// Interoperabilitas: gunakan SaltLengthAuto untuk menerima tanda tangan dari library lain
			pssOpts := &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthAuto, Hash: crypto.SHA256}
			err := rsa.VerifyPSS(pubKey, crypto.SHA256, hashed[:], sigBytes, pssOpts)
			return nativeBool(err == nil)
		},
	}
	definePublicFunction(env, "rsa_verifikasi_pss", rsaVerifyFn)
	definePublicFunction(env, "rsa_verify_pss", rsaVerifyFn)
}

// =========================================================================
// HELPER PARSING UNTUK KRIPTOGRAFI MODERN
// =========================================================================

func parse16Bytes(arg Object, funcName string) ([]byte, Object) {
	switch v := arg.(type) {
	case *Bytes:
		if len(v.Elements) != 16 {
			return nil, newError(nil, fmt.Sprintf("%s: panjang bita harus 16 bita, menerima %d", funcName, len(v.Elements)), "")
		}
		return v.Elements, nil
	case *String:
		s := strings.TrimSpace(v.Value)
		if len(s) == 32 {
			if hx, err := hex.DecodeString(s); err == nil && len(hx) == 16 {
				return hx, nil
			}
		}
		if b64, err := decodeBase64Flex(s); err == nil && len(b64) == 16 {
			return b64, nil
		}
		if len([]byte(s)) == 16 {
			return []byte(s), nil
		}
		return nil, newError(nil, fmt.Sprintf("%s: parameter 16-bita tidak valid (diharapkan 16 karakter, hex 32 karakter, atau Base64)", funcName), "")
	default:
		return nil, newError(nil, fmt.Sprintf("%s: argumen harus berupa teks atau Bita 16-bita", funcName), "")
	}
}

func parseAESKeyAny(arg Object, funcName string) ([]byte, Object) {
	switch v := arg.(type) {
	case *Bytes:
		l := len(v.Elements)
		if l != 16 && l != 24 && l != 32 {
			return nil, newError(nil, fmt.Sprintf("%s: panjang kunci AES harus 16, 24, atau 32 bita (menerima %d)", funcName, l), "")
		}
		return v.Elements, nil
	case *String:
		s := strings.TrimSpace(v.Value)
		if len(s) == 32 || len(s) == 48 || len(s) == 64 {
			if hx, err := hex.DecodeString(s); err == nil && (len(hx) == 16 || len(hx) == 24 || len(hx) == 32) {
				return hx, nil
			}
		}
		if b64, err := decodeBase64Flex(s); err == nil && (len(b64) == 16 || len(b64) == 24 || len(b64) == 32) {
			return b64, nil
		}
		l := len([]byte(s))
		if l == 16 || l == 24 || l == 32 {
			return []byte(s), nil
		}
		return nil, newError(nil, fmt.Sprintf("%s: panjang kunci AES harus 128, 192, atau 256-bit (16/24/32 bita)", funcName), "")
	default:
		return nil, newError(nil, fmt.Sprintf("%s: argumen kunci harus berupa teks atau Bita", funcName), "")
	}
}

func parse32BytesFlex(arg Object, funcName string) ([]byte, Object) {
	switch v := arg.(type) {
	case *Bytes:
		if len(v.Elements) != 32 {
			return nil, newError(nil, fmt.Sprintf("%s: panjang harus 32 bita, menerima %d", funcName, len(v.Elements)), "")
		}
		return v.Elements, nil
	case *String:
		s := strings.TrimSpace(v.Value)
		if len(s) == 64 {
			if hx, err := hex.DecodeString(s); err == nil && len(hx) == 32 {
				return hx, nil
			}
		}
		if b64, err := decodeBase64Flex(s); err == nil && len(b64) == 32 {
			return b64, nil
		}
		if len([]byte(s)) == 32 {
			return []byte(s), nil
		}
		return nil, newError(nil, fmt.Sprintf("%s: parameter 32-bita tidak valid (diharapkan hex 64-karakter, Base64 44-karakter, atau 32 bita mentah)", funcName), "")
	default:
		return nil, newError(nil, fmt.Sprintf("%s: argumen harus berupa teks atau Bita", funcName), "")
	}
}

func parseRSAPrivateKey(arg Object, funcName string) (*rsa.PrivateKey, Object) {
	strObj, ok := arg.(*String)
	if !ok {
		return nil, newError(nil, fmt.Sprintf("%s: kunci privat RSA harus berupa teks PEM", funcName), "")
	}
	block, _ := pem.Decode([]byte(strObj.Value))
	if block == nil {
		return nil, newError(nil, fmt.Sprintf("%s: format PEM kunci privat RSA tidak valid", funcName), "")
	}
	if priv, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return priv, nil
	}
	if pkcs8, err := x509.ParsePKCS8PrivateKey(block.Bytes); err == nil {
		if priv, ok := pkcs8.(*rsa.PrivateKey); ok {
			return priv, nil
		}
	}
	return nil, newError(nil, fmt.Sprintf("%s: gagal mengurai kunci privat RSA (harus format PKCS#1 atau PKCS#8)", funcName), "")
}

func parseRSAPublicKey(arg Object, funcName string) (*rsa.PublicKey, Object) {
	strObj, ok := arg.(*String)
	if !ok {
		return nil, newError(nil, fmt.Sprintf("%s: kunci publik RSA harus berupa teks PEM", funcName), "")
	}
	block, _ := pem.Decode([]byte(strObj.Value))
	if block == nil {
		return nil, newError(nil, fmt.Sprintf("%s: format PEM kunci publik RSA tidak valid", funcName), "")
	}
	if pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes); err == nil {
		if pub, ok := pubInterface.(*rsa.PublicKey); ok {
			return pub, nil
		}
	}
	if pub, err := x509.ParsePKCS1PublicKey(block.Bytes); err == nil {
		return pub, nil
	}
	return nil, newError(nil, fmt.Sprintf("%s: gagal mengurai kunci publik RSA", funcName), "")
}

// pkcs7UnpadConstantTime memvalidasi dan menghapus padding PKCS#7 dalam waktu konstan
// untuk memitigasi serangan padding oracle.
func pkcs7UnpadConstantTime(src []byte, blockSize int) ([]byte, error) {
	if len(src) == 0 || len(src)%blockSize != 0 {
		return nil, fmt.Errorf("dekripsi gagal")
	}
	padLen := int(src[len(src)-1])
	valid := 1
	if padLen == 0 || padLen > blockSize || padLen > len(src) {
		valid = 0
		padLen = 1 // Clamp untuk menghindari indexing out of range
	}
	for i := 0; i < blockSize; i++ {
		idx := len(src) - 1 - i
		mask := subtle.ConstantTimeLessOrEq(i+1, padLen)
		eq := subtle.ConstantTimeByteEq(src[idx], byte(padLen))
		if (mask & (eq ^ 1)) == 1 {
			valid = 0
		}
	}
	if valid != 1 {
		return nil, fmt.Errorf("dekripsi gagal")
	}
	return src[:len(src)-padLen], nil
}
