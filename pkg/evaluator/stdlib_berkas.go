package evaluator

import (
	"bufio"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func initStdlibBerkas(env *Environment) {
	// -------------------------------------------------------------
	// 1. Membaca & Menulis (I/O)
	// -------------------------------------------------------------

	// baca / read
	bacaFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "berkas.baca mengharapkan 1 argumen: path berkas", "Contoh: berkas.baca(\"./data.txt\")")
			}
			pathStr, ok := args[0].(*String)
			if !ok {
				return newError(nil, fmt.Sprintf("argumen path berkas.baca harus teks, menerima %s", args[0].Type()), "")
			}

			data, err := os.ReadFile(pathStr.Value)
			if err != nil {
				suggestion := suggestSimilarFile(pathStr.Value)
				hint := "Pastikan path berkas sudah benar dan berkas tersebut ada di disk."
				msg := fmt.Sprintf("Gagal membaca berkas '%s': %v", pathStr.Value, err)
				if suggestion != "" {
					hint = fmt.Sprintf("Maksud Anda mungkin '%s'? Berkas serupa ditemukan di direktori target.", suggestion)
					msg = fmt.Sprintf("Berkas '%s' tidak ditemukan. Maksud Anda mungkin: '%s'?", pathStr.Value, suggestion)
				}
				return newError(nil, msg, hint)
			}
			return &String{Value: string(data)}
		},
	}
	definePublicFunction(env, "baca", bacaFn)
	definePublicFunction(env, "read", bacaFn)

	// tulis / write (Atomic Write)
	tulisFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "berkas.tulis mengharapkan 2 argumen: path berkas dan data teks", "Contoh: berkas.tulis(\"./out.txt\", \"isi teks\")")
			}
			pathStr, ok := args[0].(*String)
			if !ok {
				return newError(nil, fmt.Sprintf("argumen path berkas.tulis harus teks, menerima %s", args[0].Type()), "")
			}
			data := args[1].Inspect()
			if s, ok := args[1].(*String); ok {
				data = s.Value
			}

			if err := atomicWriteFile(pathStr.Value, []byte(data), 0644); err != nil {
				return newError(nil, fmt.Sprintf("Gagal menulis berkas '%s': %v", pathStr.Value, err), "Pastikan Anda memiliki izin tulis di direktori target.")
			}
			return TRUE
		},
	}
	definePublicFunction(env, "tulis", tulisFn)
	definePublicFunction(env, "write", tulisFn)

	// baca_bita / read_bytes
	bacaBitaFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "berkas.baca_bita mengharapkan 1 argumen: path berkas", "Contoh: berkas.baca_bita(\"./gambar.png\")")
			}
			pathStr, ok := args[0].(*String)
			if !ok {
				return newError(nil, fmt.Sprintf("argumen path berkas.baca_bita harus teks, menerima %s", args[0].Type()), "")
			}
			data, err := os.ReadFile(pathStr.Value)
			if err != nil {
				return newError(nil, fmt.Sprintf("Gagal membaca berkas '%s': %v", pathStr.Value, err), "")
			}
			return &Bytes{Elements: data}
		},
	}
	definePublicFunction(env, "baca_bita", bacaBitaFn)
	definePublicFunction(env, "read_bytes", bacaBitaFn)

	// tulis_bita / write_bytes
	tulisBitaFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "berkas.tulis_bita mengharapkan 2 argumen: path berkas dan data bita", "Contoh: berkas.tulis_bita(\"./out.bin\", data_bita)")
			}
			pathStr, ok := args[0].(*String)
			if !ok {
				return newError(nil, fmt.Sprintf("argumen path berkas.tulis_bita harus teks, menerima %s", args[0].Type()), "")
			}
			var byteData []byte
			if b, ok := args[1].(*Bytes); ok {
				byteData = b.Elements
			} else if s, ok := args[1].(*String); ok {
				byteData = []byte(s.Value)
			} else if arr, ok := args[1].(*Array); ok {
				byteData = make([]byte, len(arr.Elements))
				for i, el := range arr.Elements {
					if iVal, _, isInt, ok := CoerceToNumeric(el); ok && isInt {
						byteData[i] = byte(iVal)
					}
				}
			} else {
				return newError(nil, fmt.Sprintf("argumen kedua berkas.tulis_bita harus Bita, menerima %s", args[1].Type()), "")
			}

			if err := atomicWriteFile(pathStr.Value, byteData, 0644); err != nil {
				return newError(nil, fmt.Sprintf("Gagal menulis berkas '%s': %v", pathStr.Value, err), "")
			}
			return TRUE
		},
	}
	definePublicFunction(env, "tulis_bita", tulisBitaFn)
	definePublicFunction(env, "write_bytes", tulisBitaFn)

	// tambah / append
	tambahFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "berkas.tambah mengharapkan 2 argumen: path berkas dan data teks", "Contoh: berkas.tambah(\"./log.txt\", \"baris baru\\n\")")
			}
			pathStr, ok := args[0].(*String)
			if !ok {
				return newError(nil, fmt.Sprintf("argumen path berkas.tambah harus teks, menerima %s", args[0].Type()), "")
			}
			data := args[1].Inspect()
			if s, ok := args[1].(*String); ok {
				data = s.Value
			}

			dir := filepath.Dir(pathStr.Value)
			if dir != "" && dir != "." {
				_ = os.MkdirAll(dir, 0755)
			}

			f, err := os.OpenFile(pathStr.Value, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return newError(nil, fmt.Sprintf("Gagal membuka berkas '%s' untuk penambahan: %v", pathStr.Value, err), "")
			}
			defer f.Close()

			if _, err := f.WriteString(data); err != nil {
				return newError(nil, fmt.Sprintf("Gagal menambahkan teks ke berkas '%s': %v", pathStr.Value, err), "")
			}
			return TRUE
		},
	}
	definePublicFunction(env, "tambah", tambahFn)
	definePublicFunction(env, "append", tambahFn)

	// ganti_teks / replace_text (In-Place Find & Replace)
	gantiTeksFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 3 {
				return newError(nil, "berkas.ganti_teks mengharapkan 3 argumen: path, teksLama, teksBaru", "Contoh: berkas.ganti_teks(\"./file.txt\", \"lama\", \"baru\")")
			}
			pathStr, ok1 := args[0].(*String)
			oldStr, ok2 := args[1].(*String)
			newStr, ok3 := args[2].(*String)
			if !ok1 || !ok2 || !ok3 {
				return newError(nil, "argumen berkas.ganti_teks harus berupa teks", "")
			}

			dataBytes, err := os.ReadFile(pathStr.Value)
			if err != nil {
				suggestion := suggestSimilarFile(pathStr.Value)
				hint := "Pastikan berkas ada."
				if suggestion != "" {
					hint = fmt.Sprintf("Maksud Anda mungkin '%s'?", suggestion)
				}
				return newError(nil, fmt.Sprintf("Gagal membaca berkas '%s': %v", pathStr.Value, err), hint)
			}

			content := strings.ReplaceAll(string(dataBytes), oldStr.Value, newStr.Value)
			if err := atomicWriteFile(pathStr.Value, []byte(content), 0644); err != nil {
				return newError(nil, fmt.Sprintf("Gagal menulis hasil ganti teks: %v", err), "")
			}
			return TRUE
		},
	}
	definePublicFunction(env, "ganti_teks", gantiTeksFn)
	definePublicFunction(env, "replace_text", gantiTeksFn)

	// kepala / head (Baca n baris pertama)
	kepalaFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 1 {
				return newError(nil, "berkas.kepala mengharapkan setidaknya 1 argumen: path berkas, [jumlah_baris=10]", "Contoh: berkas.kepala(\"./app.log\", 5)")
			}
			pathStr, ok := args[0].(*String)
			if !ok {
				return newError(nil, "argumen path harus berupa teks", "")
			}
			nLines := int64(10)
			if len(args) > 1 {
				if i, ok := args[1].(*Integer); ok {
					nLines = i.Value
				}
			}
			if nLines <= 0 {
				return &Array{Elements: []Object{}}
			}

			f, err := os.Open(pathStr.Value)
			if err != nil {
				suggestion := suggestSimilarFile(pathStr.Value)
				hint := "Pastikan berkas ada."
				if suggestion != "" {
					hint = fmt.Sprintf("Maksud Anda mungkin '%s'?", suggestion)
				}
				return newError(nil, fmt.Sprintf("Gagal membuka berkas '%s': %v", pathStr.Value, err), hint)
			}
			defer f.Close()

			scanner := bufio.NewScanner(f)
			var lines []Object
			for scanner.Scan() && int64(len(lines)) < nLines {
				lines = append(lines, &String{Value: scanner.Text()})
			}
			return &Array{Elements: lines}
		},
	}
	definePublicFunction(env, "kepala", kepalaFn)
	definePublicFunction(env, "head", kepalaFn)

	// ekor / tail (Baca n baris terakhir)
	ekorFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 1 {
				return newError(nil, "berkas.ekor mengharapkan setidaknya 1 argumen: path berkas, [jumlah_baris=10]", "Contoh: berkas.ekor(\"./app.log\", 10)")
			}
			pathStr, ok := args[0].(*String)
			if !ok {
				return newError(nil, "argumen path harus berupa teks", "")
			}
			nLines := int64(10)
			if len(args) > 1 {
				if i, ok := args[1].(*Integer); ok {
					nLines = i.Value
				}
			}
			if nLines <= 0 {
				return &Array{Elements: []Object{}}
			}

			f, err := os.Open(pathStr.Value)
			if err != nil {
				suggestion := suggestSimilarFile(pathStr.Value)
				hint := "Pastikan berkas ada."
				if suggestion != "" {
					hint = fmt.Sprintf("Maksud Anda mungkin '%s'?", suggestion)
				}
				return newError(nil, fmt.Sprintf("Gagal membuka berkas '%s': %v", pathStr.Value, err), hint)
			}
			defer f.Close()

			scanner := bufio.NewScanner(f)
			var ring []string
			for scanner.Scan() {
				ring = append(ring, scanner.Text())
				if int64(len(ring)) > nLines {
					ring = ring[1:]
				}
			}
			var lines []Object
			for _, l := range ring {
				lines = append(lines, &String{Value: l})
			}
			return &Array{Elements: lines}
		},
	}
	definePublicFunction(env, "ekor", ekorFn)
	definePublicFunction(env, "tail", ekorFn)

	// hash / checksum (SHA-256 / MD5 / SHA-1)
	hashFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 1 {
				return newError(nil, "berkas.hash mengharapkan 1 atau 2 argumen: path, [algoritma='sha256']", "Contoh: berkas.hash(\"./file.pdf\")")
			}
			pathStr, ok := args[0].(*String)
			if !ok {
				return newError(nil, "argumen path harus berupa teks", "")
			}
			algo := "sha256"
			if len(args) > 1 {
				if a, ok := args[1].(*String); ok {
					algo = strings.ToLower(a.Value)
				}
			}

			var hasher hash.Hash
			switch algo {
			case "md5":
				hasher = md5.New()
			case "sha1":
				hasher = sha1.New()
			case "sha512":
				hasher = sha512.New()
			default:
				hasher = sha256.New()
			}

			f, err := os.Open(pathStr.Value)
			if err != nil {
				suggestion := suggestSimilarFile(pathStr.Value)
				hint := "Pastikan berkas ada."
				if suggestion != "" {
					hint = fmt.Sprintf("Maksud Anda mungkin '%s'?", suggestion)
				}
				return newError(nil, fmt.Sprintf("Gagal membuka berkas '%s': %v", pathStr.Value, err), hint)
			}
			defer f.Close()

			if _, err := io.Copy(hasher, f); err != nil {
				return newError(nil, fmt.Sprintf("Gagal menghitung hash berkas '%s': %v", pathStr.Value, err), "")
			}
			return &String{Value: hex.EncodeToString(hasher.Sum(nil))}
		},
	}
	definePublicFunction(env, "hash", hashFn)
	definePublicFunction(env, "checksum", hashFn)

	// -------------------------------------------------------------
	// 2. Fitur Cerdas: muat & simpan (Auto JSON / Text)
	// -------------------------------------------------------------

	// muat / load
	muatFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "berkas.muat mengharapkan 1 argumen: path berkas", "Contoh: berkas.muat(\"./data.json\")")
			}
			pathStr, ok := args[0].(*String)
			if !ok {
				return newError(nil, fmt.Sprintf("argumen path berkas.muat harus teks, menerima %s", args[0].Type()), "")
			}

			dataBytes, err := os.ReadFile(pathStr.Value)
			if err != nil {
				suggestion := suggestSimilarFile(pathStr.Value)
				hint := "Pastikan path berkas sudah benar dan berkas tersebut ada."
				msg := fmt.Sprintf("Gagal memuat berkas '%s': %v", pathStr.Value, err)
				if suggestion != "" {
					hint = fmt.Sprintf("Maksud Anda mungkin '%s'? Berkas serupa ditemukan di direktori target.", suggestion)
					msg = fmt.Sprintf("Berkas '%s' tidak ditemukan. Maksud Anda mungkin: '%s'?", pathStr.Value, suggestion)
				}
				return newError(nil, msg, hint)
			}

			ext := strings.ToLower(filepath.Ext(pathStr.Value))
			if ext == ".json" {
				decoder := json.NewDecoder(strings.NewReader(string(dataBytes)))
				decoder.UseNumber()
				var raw any
				if err := decoder.Decode(&raw); err != nil {
					return newError(nil, fmt.Sprintf("Gagal mengurai JSON dari berkas '%s': %v", pathStr.Value, err), "Periksa kevalidan format JSON di dalam berkas.")
				}
				return convertGoJSONToCelang(raw)
			}

			return &String{Value: string(dataBytes)}
		},
	}
	definePublicFunction(env, "muat", muatFn)
	definePublicFunction(env, "load", muatFn)

	// simpan / save
	simpanFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "berkas.simpan mengharapkan 2 argumen: path berkas dan data", "Contoh: berkas.simpan(\"./data.json\", objekData)")
			}
			pathStr, ok := args[0].(*String)
			if !ok {
				return newError(nil, fmt.Sprintf("argumen path berkas.simpan harus teks, menerima %s", args[0].Type()), "")
			}

			ext := strings.ToLower(filepath.Ext(pathStr.Value))
			var contentBytes []byte

			if ext == ".json" || args[1].Type() == MAP_OBJ || args[1].Type() == ARRAY_OBJ {
				nativeVal, err := convertCelangToJSONNative(args[1])
				if err != nil {
					return newError(nil, fmt.Sprintf("Gagal mengubah data ke JSON: %v", err), "")
				}
				jsonBytes, err := json.MarshalIndent(nativeVal, "", "  ")
				if err != nil {
					return newError(nil, fmt.Sprintf("Gagal memformat JSON: %v", err), "")
				}
				contentBytes = jsonBytes
			} else if s, ok := args[1].(*String); ok {
				contentBytes = []byte(s.Value)
			} else {
				contentBytes = []byte(args[1].Inspect())
			}

			if err := atomicWriteFile(pathStr.Value, contentBytes, 0644); err != nil {
				return newError(nil, fmt.Sprintf("Gagal menyimpan berkas '%s': %v", pathStr.Value, err), "")
			}
			return TRUE
		},
	}
	definePublicFunction(env, "simpan", simpanFn)
	definePublicFunction(env, "save", simpanFn)

	// -------------------------------------------------------------
	// 3. Streaming Jalur Data: alir / stream (CSP Channel)
	// -------------------------------------------------------------

	alirFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 1 {
				return newError(nil, "berkas.alir mengharapkan 1 argumen: path berkas", "Contoh: untuk baris dalam berkas.alir(\"./log.txt\") { ... }")
			}
			pathStr, ok := args[0].(*String)
			if !ok {
				return newError(nil, fmt.Sprintf("argumen path berkas.alir harus teks, menerima %s", args[0].Type()), "")
			}

			f, err := os.Open(pathStr.Value)
			if err != nil {
				suggestion := suggestSimilarFile(pathStr.Value)
				hint := "Pastikan berkas ada dan dapat dibaca."
				if suggestion != "" {
					hint = fmt.Sprintf("Maksud Anda mungkin '%s'? Berkas tersebut ditemukan di direktori target.", suggestion)
				}
				return newError(nil, fmt.Sprintf("Gagal membuka aliran berkas '%s': %v", pathStr.Value, err), hint)
			}

			ch := NewChannel(64, "Teks")

			go func() {
				defer f.Close()
				defer ch.Close()

				scanner := bufio.NewScanner(f)
				buf := make([]byte, 64*1024)
				scanner.Buffer(buf, 1024*1024)

				for scanner.Scan() {
					ch.mu.Lock()
					closed := ch.Closed
					ch.mu.Unlock()
					if closed {
						return
					}

					lineText := scanner.Text()
					ch.Ch <- &String{Value: lineText}
				}
			}()

			return ch
		},
	}
	definePublicFunction(env, "alir", alirFn)
	definePublicFunction(env, "stream", alirFn)

	// -------------------------------------------------------------
	// 4. Manajemen Direktori & Penelusuran (Glob, Watch, Walk)
	// -------------------------------------------------------------

	// cari / glob (Recursive globbing: **/*.ce)
	cariFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "berkas.cari mengharapkan 1 argumen pola path", "Contoh: berkas.cari(\"./src/**/*.ce\")")
			}
			patternStr, ok := args[0].(*String)
			if !ok {
				return newError(nil, "argumen pola harus berupa teks", "")
			}

			matches, err := globRecursive(patternStr.Value)
			if err != nil {
				return newError(nil, fmt.Sprintf("Pola pencarian tidak valid '%s': %v", patternStr.Value, err), "")
			}

			elements := make([]Object, 0, len(matches))
			for _, m := range matches {
				elements = append(elements, &String{Value: m})
			}
			return &Array{Elements: elements}
		},
	}
	definePublicFunction(env, "cari", cariFn)
	definePublicFunction(env, "glob", cariFn)

	// telusuri / walk (Recursive directory walker with callback)
	telusuriFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "berkas.telusuri mengharapkan 2 argumen: path dan fungsi callback", "Contoh: berkas.telusuri(\"./src\", (jalur, info) => { ... })")
			}
			pathStr, ok := args[0].(*String)
			if !ok {
				return newError(nil, "argumen path harus berupa teks", "")
			}
			callback := args[1]

			err := filepath.WalkDir(pathStr.Value, func(p string, d os.DirEntry, err error) error {
				if err != nil {
					return nil
				}
				size := int64(0)
				if fi, sErr := d.Info(); sErr == nil {
					size = fi.Size()
				}
				infoMap := &Map{
					Pairs: map[string]Object{
						"nama":          &String{Value: d.Name()},
						"apakah_folder": nativeBoolToBoolean(d.IsDir()),
						"apakah_berkas": nativeBoolToBoolean(!d.IsDir()),
						"ekstensi":      &String{Value: filepath.Ext(d.Name())},
						"ukuran":        &Integer{Value: size},
					},
					Order: []string{"nama", "apakah_folder", "apakah_berkas", "ekstensi", "ukuran"},
				}

				res := ApplyFunction(callback, []Object{&String{Value: p}, infoMap}, nil)
				if b, ok := res.(*Boolean); ok && !b.Value {
					if d.IsDir() {
						return filepath.SkipDir
					}
				}
				return nil
			})

			if err != nil {
				return newError(nil, fmt.Sprintf("Gagal menelusuri '%s': %v", pathStr.Value, err), "")
			}
			return TRUE
		},
	}
	definePublicFunction(env, "telusuri", telusuriFn)
	definePublicFunction(env, "walk", telusuriFn)

	// pantau / watch (File & Folder watcher via Saluran or Callback)
	pantauFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 1 {
				return newError(nil, "berkas.pantau mengharapkan setidaknya 1 argumen: path, [callback]", "Contoh: untuk kejadian dalam berkas.pantau(\"./src\") { ... }")
			}
			pathStr, ok := args[0].(*String)
			if !ok {
				return newError(nil, "argumen path harus berupa teks", "")
			}

			var callback Object
			if len(args) > 1 {
				callback = args[1]
			}

			ch := NewChannel(32, "Kamus")

			// Ambil snapshot awal
			snapshot := make(map[string]time.Time)
			_ = filepath.WalkDir(pathStr.Value, func(p string, d os.DirEntry, err error) error {
				if err != nil || d.IsDir() {
					return nil
				}
				if fi, err := d.Info(); err == nil {
					snapshot[p] = fi.ModTime()
				}
				return nil
			})

			go func() {
				ticker := time.NewTicker(300 * time.Millisecond)
				defer ticker.Stop()
				defer ch.Close()

				for range ticker.C {
					ch.mu.Lock()
					closed := ch.Closed
					ch.mu.Unlock()
					if closed {
						return
					}

					current := make(map[string]time.Time)
					_ = filepath.WalkDir(pathStr.Value, func(p string, d os.DirEntry, err error) error {
						if err != nil || d.IsDir() {
							return nil
						}
						if fi, err := d.Info(); err == nil {
							current[p] = fi.ModTime()
						}
						return nil
					})

					// 1. Berkas baru / dimodifikasi
					for p, modTime := range current {
						prevTime, exists := snapshot[p]
						if !exists {
							sendWatchEvent(ch, callback, p, "buat", "create")
						} else if modTime.After(prevTime) {
							sendWatchEvent(ch, callback, p, "ubah", "modify")
						}
					}

					// 2. Berkas dihapus
					for p := range snapshot {
						if _, exists := current[p]; !exists {
							sendWatchEvent(ch, callback, p, "hapus", "delete")
						}
					}

					snapshot = current
				}
			}()

			return ch
		},
	}
	definePublicFunction(env, "pantau", pantauFn)
	definePublicFunction(env, "watch", pantauFn)

	// ada / exists
	adaFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "berkas.ada mengharapkan 1 argumen: path berkas", "")
			}
			pathStr, ok := args[0].(*String)
			if !ok {
				return newError(nil, fmt.Sprintf("argumen path berkas.ada harus teks, menerima %s", args[0].Type()), "")
			}

			_, err := os.Stat(pathStr.Value)
			if err == nil {
				return TRUE
			}
			return FALSE
		},
	}
	definePublicFunction(env, "ada", adaFn)
	definePublicFunction(env, "exists", adaFn)

	// hapus / remove
	hapusFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "berkas.hapus mengharapkan 1 argumen: path berkas", "")
			}
			pathStr, ok := args[0].(*String)
			if !ok {
				return newError(nil, fmt.Sprintf("argumen path berkas.hapus harus teks, menerima %s", args[0].Type()), "")
			}

			if err := os.Remove(pathStr.Value); err != nil {
				suggestion := suggestSimilarFile(pathStr.Value)
				hint := "Pastikan berkas ada dan Anda memiliki izin untuk menghapusnya."
				if suggestion != "" {
					hint = fmt.Sprintf("Maksud Anda mungkin '%s'? Berkas tersebut ditemukan di direktori target.", suggestion)
				}
				return newError(nil, fmt.Sprintf("Gagal menghapus berkas '%s': %v", pathStr.Value, err), hint)
			}
			return TRUE
		},
	}
	definePublicFunction(env, "hapus", hapusFn)
	definePublicFunction(env, "remove", hapusFn)

	// hapus_semua / remove_all
	hapusSemuaFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "berkas.hapus_semua mengharapkan 1 argumen: path direktori", "Contoh: berkas.hapus_semua(\"./folder_lama\")")
			}
			pathStr, ok := args[0].(*String)
			if !ok {
				return newError(nil, fmt.Sprintf("argumen path berkas.hapus_semua harus teks, menerima %s", args[0].Type()), "")
			}

			if err := os.RemoveAll(pathStr.Value); err != nil {
				return newError(nil, fmt.Sprintf("Gagal menghapus folder '%s': %v", pathStr.Value, err), "")
			}
			return TRUE
		},
	}
	definePublicFunction(env, "hapus_semua", hapusSemuaFn)
	definePublicFunction(env, "remove_all", hapusSemuaFn)

	// buat_folder / mkdir
	buatFolderFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "berkas.buat_folder mengharapkan 1 argumen: path folder", "Contoh: berkas.buat_folder(\"./data/laporan\")")
			}
			pathStr, ok := args[0].(*String)
			if !ok {
				return newError(nil, fmt.Sprintf("argumen path berkas.buat_folder harus teks, menerima %s", args[0].Type()), "")
			}

			if err := os.MkdirAll(pathStr.Value, 0755); err != nil {
				return newError(nil, fmt.Sprintf("Gagal membuat folder '%s': %v", pathStr.Value, err), "")
			}
			return TRUE
		},
	}
	definePublicFunction(env, "buat_folder", buatFolderFn)
	definePublicFunction(env, "mkdir", buatFolderFn)

	// daftar / list / ls
	daftarFn := &Builtin{
		Fn: func(args ...Object) Object {
			targetDir := "."
			if len(args) > 0 {
				if s, ok := args[0].(*String); ok {
					targetDir = s.Value
				}
			}

			entries, err := os.ReadDir(targetDir)
			if err != nil {
				suggestion := suggestSimilarFile(targetDir)
				hint := "Pastikan direktori target ada."
				if suggestion != "" {
					hint = fmt.Sprintf("Maksud Anda mungkin folder '%s'?", suggestion)
				}
				return newError(nil, fmt.Sprintf("Gagal membaca direktori '%s': %v", targetDir, err), hint)
			}

			elements := make([]Object, 0, len(entries))
			for _, entry := range entries {
				elements = append(elements, &String{Value: entry.Name()})
			}
			return &Array{Elements: elements}
		},
	}
	definePublicFunction(env, "daftar", daftarFn)
	definePublicFunction(env, "list", daftarFn)
	definePublicFunction(env, "ls", daftarFn)

	// salin / copy
	salinFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "berkas.salin mengharapkan 2 argumen: asal dan tujuan", "Contoh: berkas.salin(\"./asal.txt\", \"./tujuan.txt\")")
			}
			srcStr, ok1 := args[0].(*String)
			dstStr, ok2 := args[1].(*String)
			if !ok1 || !ok2 {
				return newError(nil, "argumen berkas.salin harus berupa teks", "")
			}

			fi, err := os.Stat(srcStr.Value)
			if err != nil {
				suggestion := suggestSimilarFile(srcStr.Value)
				hint := "Pastikan berkas atau folder asal ada."
				if suggestion != "" {
					hint = fmt.Sprintf("Maksud Anda mungkin '%s'?", suggestion)
				}
				return newError(nil, fmt.Sprintf("Gagal menyalin '%s': %v", srcStr.Value, err), hint)
			}

			if fi.IsDir() {
				if err := copyDir(srcStr.Value, dstStr.Value); err != nil {
					return newError(nil, fmt.Sprintf("Gagal menyalin folder: %v", err), "")
				}
			} else {
				if err := copyFile(srcStr.Value, dstStr.Value); err != nil {
					return newError(nil, fmt.Sprintf("Gagal menyalin berkas: %v", err), "")
				}
			}
			return TRUE
		},
	}
	definePublicFunction(env, "salin", salinFn)
	definePublicFunction(env, "copy", salinFn)

	// pindah / move / ganti_nama / rename
	pindahFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "berkas.pindah mengharapkan 2 argumen: asal dan tujuan", "Contoh: berkas.pindah(\"./lama.txt\", \"./baru.txt\")")
			}
			srcStr, ok1 := args[0].(*String)
			dstStr, ok2 := args[1].(*String)
			if !ok1 || !ok2 {
				return newError(nil, "argumen berkas.pindah harus berupa teks", "")
			}

			dstDir := filepath.Dir(dstStr.Value)
			if dstDir != "" && dstDir != "." {
				_ = os.MkdirAll(dstDir, 0755)
			}

			err := os.Rename(srcStr.Value, dstStr.Value)
			if err != nil {
				fi, statErr := os.Stat(srcStr.Value)
				if statErr != nil {
					return newError(nil, fmt.Sprintf("Gagal memindahkan '%s': %v", srcStr.Value, err), "")
				}
				if fi.IsDir() {
					if cErr := copyDir(srcStr.Value, dstStr.Value); cErr != nil {
						return newError(nil, fmt.Sprintf("Gagal memindahkan folder: %v", cErr), "")
					}
					_ = os.RemoveAll(srcStr.Value)
				} else {
					if cErr := copyFile(srcStr.Value, dstStr.Value); cErr != nil {
						return newError(nil, fmt.Sprintf("Gagal memindahkan berkas: %v", cErr), "")
					}
					_ = os.Remove(srcStr.Value)
				}
			}
			return TRUE
		},
	}
	definePublicFunction(env, "pindah", pindahFn)
	definePublicFunction(env, "move", pindahFn)
	definePublicFunction(env, "ganti_nama", pindahFn)
	definePublicFunction(env, "rename", pindahFn)

	// -------------------------------------------------------------
	// 5. Pemeriksaan Tipe & Metadata
	// -------------------------------------------------------------

	// apakah_folder / is_dir
	apakahFolderFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "berkas.apakah_folder mengharapkan 1 argumen: path", "")
			}
			pathStr, ok := args[0].(*String)
			if !ok {
				return FALSE
			}
			fi, err := os.Stat(pathStr.Value)
			if err != nil {
				return FALSE
			}
			if fi.IsDir() {
				return TRUE
			}
			return FALSE
		},
	}
	definePublicFunction(env, "apakah_folder", apakahFolderFn)
	definePublicFunction(env, "is_dir", apakahFolderFn)

	// apakah_berkas / is_file
	apakahBerkasFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "berkas.apakah_berkas mengharapkan 1 argumen: path", "")
			}
			pathStr, ok := args[0].(*String)
			if !ok {
				return FALSE
			}
			fi, err := os.Stat(pathStr.Value)
			if err != nil {
				return FALSE
			}
			if !fi.IsDir() {
				return TRUE
			}
			return FALSE
		},
	}
	definePublicFunction(env, "apakah_berkas", apakahBerkasFn)
	definePublicFunction(env, "is_file", apakahBerkasFn)

	// ukuran / size
	ukuranFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "berkas.ukuran mengharapkan 1 argumen: path berkas", "")
			}
			pathStr, ok := args[0].(*String)
			if !ok {
				return newError(nil, "argumen berkas.ukuran harus berupa teks", "")
			}
			fi, err := os.Stat(pathStr.Value)
			if err != nil {
				suggestion := suggestSimilarFile(pathStr.Value)
				hint := "Pastikan berkas ada."
				if suggestion != "" {
					hint = fmt.Sprintf("Maksud Anda mungkin '%s'?", suggestion)
				}
				return newError(nil, fmt.Sprintf("Gagal mendapatkan ukuran berkas '%s': %v", pathStr.Value, err), hint)
			}
			return &Integer{Value: fi.Size()}
		},
	}
	definePublicFunction(env, "ukuran", ukuranFn)
	definePublicFunction(env, "size", ukuranFn)

	// info / stat
	infoFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "berkas.info mengharapkan 1 argumen: path", "")
			}
			pathStr, ok := args[0].(*String)
			if !ok {
				return newError(nil, "argumen berkas.info harus berupa teks", "")
			}
			fi, err := os.Stat(pathStr.Value)
			if err != nil {
				suggestion := suggestSimilarFile(pathStr.Value)
				hint := "Pastikan path berkas atau folder ada."
				if suggestion != "" {
					hint = fmt.Sprintf("Maksud Anda mungkin '%s'?", suggestion)
				}
				return newError(nil, fmt.Sprintf("Gagal memeriksa info '%s': %v", pathStr.Value, err), hint)
			}

			pairs := map[string]Object{
				"nama":          &String{Value: fi.Name()},
				"ukuran":        &Integer{Value: fi.Size()},
				"apakah_folder": nativeBoolToBoolean(fi.IsDir()),
				"apakah_berkas": nativeBoolToBoolean(!fi.IsDir()),
				"diubah_pada":   &String{Value: fi.ModTime().Format("2006-01-02 15:04:05")},
				"izin":          &String{Value: fi.Mode().String()},
			}
			order := []string{"nama", "ukuran", "apakah_folder", "apakah_berkas", "diubah_pada", "izin"}
			return &Map{Pairs: pairs, Order: order}
		},
	}
	definePublicFunction(env, "info", infoFn)
	definePublicFunction(env, "stat", infoFn)

	// -------------------------------------------------------------
	// 6. Utilitas Jalur (Path Helpers)
	// -------------------------------------------------------------

	// gabung / join
	gabungFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) == 0 {
				return &String{Value: ""}
			}
			parts := make([]string, 0, len(args))
			for _, arg := range args {
				if s, ok := arg.(*String); ok {
					parts = append(parts, s.Value)
				} else {
					parts = append(parts, arg.Inspect())
				}
			}
			return &String{Value: filepath.Join(parts...)}
		},
	}
	definePublicFunction(env, "gabung", gabungFn)
	definePublicFunction(env, "join", gabungFn)

	// nama / basename
	namaFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "berkas.nama mengharapkan 1 argumen path", "")
			}
			s, ok := args[0].(*String)
			if !ok {
				return &String{Value: ""}
			}
			return &String{Value: filepath.Base(s.Value)}
		},
	}
	definePublicFunction(env, "nama", namaFn)
	definePublicFunction(env, "basename", namaFn)

	// induk / dirname
	indukFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "berkas.induk mengharapkan 1 argumen path", "")
			}
			s, ok := args[0].(*String)
			if !ok {
				return &String{Value: ""}
			}
			return &String{Value: filepath.Dir(s.Value)}
		},
	}
	definePublicFunction(env, "induk", indukFn)
	definePublicFunction(env, "dirname", indukFn)

	// ekstensi / extname
	ekstensiFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "berkas.ekstensi mengharapkan 1 argumen path", "")
			}
			s, ok := args[0].(*String)
			if !ok {
				return &String{Value: ""}
			}
			return &String{Value: filepath.Ext(s.Value)}
		},
	}
	definePublicFunction(env, "ekstensi", ekstensiFn)
	definePublicFunction(env, "extname", ekstensiFn)

	// bersihkan_jalur / clean
	bersihkanFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "berkas.bersihkan_jalur mengharapkan 1 argumen path", "")
			}
			s, ok := args[0].(*String)
			if !ok {
				return &String{Value: ""}
			}
			return &String{Value: filepath.Clean(s.Value)}
		},
	}
	definePublicFunction(env, "bersihkan_jalur", bersihkanFn)
	definePublicFunction(env, "clean", bersihkanFn)

	// jalur_mutlak / absolute
	mutlakFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "berkas.jalur_mutlak mengharapkan 1 argumen path", "")
			}
			s, ok := args[0].(*String)
			if !ok {
				return &String{Value: ""}
			}
			abs, err := filepath.Abs(s.Value)
			if err != nil {
				return &String{Value: s.Value}
			}
			return &String{Value: abs}
		},
	}
	definePublicFunction(env, "jalur_mutlak", mutlakFn)
	definePublicFunction(env, "absolute", mutlakFn)
}

// -------------------------------------------------------------
// Helper Functions (Atomic Write, Copy, Typo Suggestion, Watch, Glob)
// -------------------------------------------------------------

func sendWatchEvent(ch *Channel, callback Object, path, aksi, action string) {
	eventMap := &Map{
		Pairs: map[string]Object{
			"path":   &String{Value: path},
			"aksi":   &String{Value: aksi},
			"action": &String{Value: action},
		},
		Order: []string{"path", "aksi", "action"},
	}

	if callback != nil {
		_ = ApplyFunction(callback, []Object{eventMap}, nil)
	}

	ch.mu.Lock()
	closed := ch.Closed
	ch.mu.Unlock()
	if !closed {
		select {
		case ch.Ch <- eventMap:
		default:
		}
	}
}

func globRecursive(pattern string) ([]string, error) {
	patternClean := filepath.Clean(pattern)
	patternSlash := filepath.ToSlash(patternClean)

	parts := strings.Split(patternSlash, "/")
	baseDir := "."
	wildcardIdx := -1
	for i, part := range parts {
		if strings.ContainsAny(part, "*?[]") {
			wildcardIdx = i
			break
		}
	}

	if wildcardIdx == -1 {
		if _, err := os.Stat(patternClean); err == nil {
			return []string{patternClean}, nil
		}
		return []string{}, nil
	}

	if wildcardIdx > 0 {
		baseDir = filepath.FromSlash(strings.Join(parts[:wildcardIdx], "/"))
	}

	re, err := globToRegex(patternSlash)
	if err != nil {
		return nil, err
	}

	var matches []string
	err = filepath.WalkDir(baseDir, func(p string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return nil
		}
		pSlash := filepath.ToSlash(filepath.Clean(p))
		trimmed := strings.TrimPrefix(pSlash, "./")
		pattTrimmed := strings.TrimPrefix(patternSlash, "./")

		if re.MatchString(pSlash) || re.MatchString(trimmed) || re.MatchString("./"+trimmed) {
			matches = append(matches, filepath.FromSlash(p))
		} else if strings.HasPrefix(patternSlash, "**/") {
			// Jika pola dimulai dengan **/, cocokkan nama relatif
			if re.MatchString(trimmed) {
				matches = append(matches, filepath.FromSlash(p))
			}
		} else if strings.HasPrefix(pattTrimmed, "*") && !d.IsDir() {
			if matched, _ := filepath.Match(pattTrimmed, filepath.Base(p)); matched {
				matches = append(matches, filepath.FromSlash(p))
			}
		}
		return nil
	})

	return matches, err
}

func globToRegex(glob string) (*regexp.Regexp, error) {
	var b strings.Builder
	b.WriteString("^")
	i := 0
	for i < len(glob) {
		if i+1 < len(glob) && glob[i:i+2] == "**" {
			if i+2 < len(glob) && glob[i+2] == '/' {
				b.WriteString("(?:.*/)?")
				i += 3
			} else {
				b.WriteString(".*")
				i += 2
			}
		} else if glob[i] == '*' {
			b.WriteString("[^/]*")
			i++
		} else if glob[i] == '?' {
			b.WriteString("[^/]")
			i++
		} else {
			char := glob[i]
			if strings.ContainsRune(`.+()|[]{}^$\`, rune(char)) {
				b.WriteByte('\\')
			}
			b.WriteByte(char)
			i++
		}
	}
	b.WriteString("$")
	return regexp.Compile(b.String())
}

func nativeBoolToBoolean(b bool) *Boolean {
	if b {
		return TRUE
	}
	return FALSE
}

func atomicWriteFile(filePath string, data []byte, perm os.FileMode) error {
	dir := filepath.Dir(filePath)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	tmpFile, err := os.CreateTemp(dir, ".celang-tmp-*")
	if err != nil {
		return os.WriteFile(filePath, data, perm)
	}
	tmpName := tmpFile.Name()
	defer func() {
		_ = os.Remove(tmpName)
	}()

	if _, err := tmpFile.Write(data); err != nil {
		tmpFile.Close()
		return err
	}
	if err := tmpFile.Sync(); err != nil {
		tmpFile.Close()
		return err
	}
	if err := tmpFile.Close(); err != nil {
		return err
	}

	_ = os.Chmod(tmpName, perm)
	return os.Rename(tmpName, filePath)
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	dstDir := filepath.Dir(dst)
	if dstDir != "" && dstDir != "." {
		_ = os.MkdirAll(dstDir, 0755)
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}

func copyDir(src, dst string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dst, 0755); err != nil {
		return err
	}
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())
		if entry.IsDir() {
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}
	return nil
}

func levenshtein(s1, s2 string) int {
	r1, r2 := []rune(strings.ToLower(s1)), []rune(strings.ToLower(s2))
	n, m := len(r1), len(r2)
	if n == 0 {
		return m
	}
	if m == 0 {
		return n
	}
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, m+1)
		dp[i][0] = i
	}
	for j := 0; j <= m; j++ {
		dp[0][j] = j
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			cost := 0
			if r1[i-1] != r2[j-1] {
				cost = 1
			}
			dp[i][j] = minInt(dp[i-1][j]+1, minInt(dp[i][j-1]+1, dp[i-1][j-1]+cost))
		}
	}
	return dp[n][m]
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func suggestSimilarFile(targetPath string) string {
	dir := filepath.Dir(targetPath)
	targetBase := filepath.Base(targetPath)
	if dir == "" {
		dir = "."
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return ""
	}
	bestMatch := ""
	bestDist := 4
	for _, entry := range entries {
		name := entry.Name()
		dist := levenshtein(targetBase, name)
		if dist < bestDist {
			bestDist = dist
			bestMatch = name
		}
	}
	if bestMatch != "" {
		if dir == "." || dir == "" {
			return bestMatch
		}
		return filepath.Join(dir, bestMatch)
	}
	return ""
}
