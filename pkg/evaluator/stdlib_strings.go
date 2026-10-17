package evaluator

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

func initStdlibTeks(env *Environment) {
	// -------------------------------------------------------------
	// 1. Operasi Dasar & Pemotongan
	// -------------------------------------------------------------

	// pisah / split(str, sep)
	pisahFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "teks.pisah mengharapkan 2 argumen: string dan pemisah", "Contoh: teks.pisah(\"a,b,c\", \",\")")
			}
			str, ok1 := args[0].(*String)
			sep, ok2 := args[1].(*String)
			if !ok1 || !ok2 {
				return newError(nil, "argumen untuk teks.pisah harus berupa teks", "")
			}
			parts := strings.Split(str.Value, sep.Value)
			elements := make([]Object, len(parts))
			for i, p := range parts {
				elements[i] = &String{Value: p}
			}
			return &Array{Elements: elements}
		},
	}
	definePublicFunction(env, "pisah", pisahFn)
	definePublicFunction(env, "split", pisahFn)

	// gabung / join(arr, sep)
	gabungFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "teks.gabung mengharapkan 2 argumen: array dan pemisah", "Contoh: teks.gabung([\"a\", \"b\"], \", \")")
			}
			arr, ok1 := args[0].(*Array)
			sep, ok2 := args[1].(*String)
			if !ok1 || !ok2 {
				return newError(nil, "argumen teks.gabung harus berupa (array, pemisah teks)", "")
			}
			strParts := make([]string, len(arr.Elements))
			for i, el := range arr.Elements {
				if s, ok := el.(*String); ok {
					strParts[i] = s.Value
				} else {
					strParts[i] = el.Inspect()
				}
			}
			return &String{Value: strings.Join(strParts, sep.Value)}
		},
	}
	definePublicFunction(env, "gabung", gabungFn)
	definePublicFunction(env, "join", gabungFn)

	// kapital / upper(str)
	kapitalFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "teks.kapital mengharapkan 1 argumen teks", "")
			}
			str, ok := args[0].(*String)
			if !ok {
				return newError(nil, fmt.Sprintf("argumen teks.kapital harus berupa teks, menerima %s", args[0].Type()), "")
			}
			return &String{Value: strings.ToUpper(str.Value)}
		},
	}
	definePublicFunction(env, "kapital", kapitalFn)
	definePublicFunction(env, "upper", kapitalFn)

	// kecil / lower(str)
	kecilFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "teks.kecil mengharapkan 1 argumen teks", "")
			}
			str, ok := args[0].(*String)
			if !ok {
				return newError(nil, fmt.Sprintf("argumen teks.kecil harus berupa teks, menerima %s", args[0].Type()), "")
			}
			return &String{Value: strings.ToLower(str.Value)}
		},
	}
	definePublicFunction(env, "kecil", kecilFn)
	definePublicFunction(env, "lower", kecilFn)

	// potong / trim(str, [cutset])
	potongFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 1 || len(args) > 2 {
				return newError(nil, "teks.potong mengharapkan 1 atau 2 argumen: teks, [karakter_potong]", "")
			}
			str, ok := args[0].(*String)
			if !ok {
				return newError(nil, "argumen pertama teks.potong harus berupa teks", "")
			}
			if len(args) == 2 {
				if cut, ok := args[1].(*String); ok {
					return &String{Value: strings.Trim(str.Value, cut.Value)}
				}
			}
			return &String{Value: strings.TrimSpace(str.Value)}
		},
	}
	definePublicFunction(env, "potong", potongFn)
	definePublicFunction(env, "trim", potongFn)

	// potong_kiri / trim_left
	potongKiriFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 1 || len(args) > 2 {
				return newError(nil, "teks.potong_kiri mengharapkan 1 atau 2 argumen: teks, [karakter_potong]", "")
			}
			str, ok := args[0].(*String)
			if !ok {
				return newError(nil, "argumen teks.potong_kiri harus berupa teks", "")
			}
			if len(args) == 2 {
				if cut, ok := args[1].(*String); ok {
					return &String{Value: strings.TrimLeft(str.Value, cut.Value)}
				}
			}
			return &String{Value: strings.TrimLeftFunc(str.Value, unicode.IsSpace)}
		},
	}
	definePublicFunction(env, "potong_kiri", potongKiriFn)
	definePublicFunction(env, "trim_left", potongKiriFn)

	// potong_kanan / trim_right
	potongKananFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 1 || len(args) > 2 {
				return newError(nil, "teks.potong_kanan mengharapkan 1 atau 2 argumen: teks, [karakter_potong]", "")
			}
			str, ok := args[0].(*String)
			if !ok {
				return newError(nil, "argumen teks.potong_kanan harus berupa teks", "")
			}
			if len(args) == 2 {
				if cut, ok := args[1].(*String); ok {
					return &String{Value: strings.TrimRight(str.Value, cut.Value)}
				}
			}
			return &String{Value: strings.TrimRightFunc(str.Value, unicode.IsSpace)}
		},
	}
	definePublicFunction(env, "potong_kanan", potongKananFn)
	definePublicFunction(env, "trim_right", potongKananFn)

	// ganti / replace(str, old, new)
	gantiFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 3 {
				return newError(nil, "teks.ganti mengharapkan 3 argumen: teks, teks_lama, teks_baru", "Contoh: teks.ganti(\"halo dunia\", \"dunia\", \"celang\")")
			}
			str, ok1 := args[0].(*String)
			oldS, ok2 := args[1].(*String)
			newS, ok3 := args[2].(*String)
			if !ok1 || !ok2 || !ok3 {
				return newError(nil, "seluruh argumen teks.ganti harus berupa teks", "")
			}
			return &String{Value: strings.ReplaceAll(str.Value, oldS.Value, newS.Value)}
		},
	}
	definePublicFunction(env, "ganti", gantiFn)
	definePublicFunction(env, "replace", gantiFn)

	// -------------------------------------------------------------
	// 2. Pemeriksaan & Pencarian
	// -------------------------------------------------------------

	// berisi / contains(str, substr)
	berisiFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "teks.berisi mengharapkan 2 argumen: teks dan teks_cari", "")
			}
			str, ok1 := args[0].(*String)
			sub, ok2 := args[1].(*String)
			if !ok1 || !ok2 {
				return newError(nil, "seluruh argumen teks.berisi harus berupa teks", "")
			}
			return nativeBoolToBoolean(strings.Contains(str.Value, sub.Value))
		},
	}
	definePublicFunction(env, "berisi", berisiFn)
	definePublicFunction(env, "contains", berisiFn)

	// dimulai_dengan / starts_with(str, prefix)
	dimulaiFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "teks.dimulai_dengan mengharapkan 2 argumen: teks dan awalan", "")
			}
			str, ok1 := args[0].(*String)
			pref, ok2 := args[1].(*String)
			if !ok1 || !ok2 {
				return newError(nil, "argumen teks.dimulai_dengan harus berupa teks", "")
			}
			return nativeBoolToBoolean(strings.HasPrefix(str.Value, pref.Value))
		},
	}
	definePublicFunction(env, "dimulai_dengan", dimulaiFn)
	definePublicFunction(env, "starts_with", dimulaiFn)

	// diakhiri_dengan / ends_with(str, suffix)
	diakhiriFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "teks.diakhiri_dengan mengharapkan 2 argumen: teks dan akhiran", "")
			}
			str, ok1 := args[0].(*String)
			suff, ok2 := args[1].(*String)
			if !ok1 || !ok2 {
				return newError(nil, "argumen teks.diakhiri_dengan harus berupa teks", "")
			}
			return nativeBoolToBoolean(strings.HasSuffix(str.Value, suff.Value))
		},
	}
	definePublicFunction(env, "diakhiri_dengan", diakhiriFn)
	definePublicFunction(env, "ends_with", diakhiriFn)

	// indeks_dari / index_of(str, sub, [startPos])
	indeksDariFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 2 || len(args) > 3 {
				return newError(nil, "teks.indeks_dari mengharapkan 2 atau 3 argumen: teks, teks_cari, [posisi_awal]", "")
			}
			str, ok1 := args[0].(*String)
			sub, ok2 := args[1].(*String)
			if !ok1 || !ok2 {
				return newError(nil, "argumen teks dan teks_cari harus berupa teks", "")
			}
			startPos := 0
			if len(args) == 3 {
				if i, ok := args[2].(*Integer); ok {
					startPos = int(i.Value)
				}
			}
			runes := []rune(str.Value)
			if startPos < 0 || startPos >= len(runes) {
				if startPos != 0 || len(runes) != 0 {
					return &Integer{Value: -1}
				}
			}
			targetSlice := string(runes[startPos:])
			idx := strings.Index(targetSlice, sub.Value)
			if idx == -1 {
				return &Integer{Value: -1}
			}
			// Convert byte offset in targetSlice to rune index
			runeOffset := len([]rune(targetSlice[:idx]))
			return &Integer{Value: int64(startPos + runeOffset)}
		},
	}
	definePublicFunction(env, "indeks_dari", indeksDariFn)
	definePublicFunction(env, "index_of", indeksDariFn)

	// indeks_terakhir / last_index_of(str, sub)
	indeksTerakhirFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "teks.indeks_terakhir mengharapkan 2 argumen: teks dan teks_cari", "")
			}
			str, ok1 := args[0].(*String)
			sub, ok2 := args[1].(*String)
			if !ok1 || !ok2 {
				return newError(nil, "argumen harus berupa teks", "")
			}
			idx := strings.LastIndex(str.Value, sub.Value)
			if idx == -1 {
				return &Integer{Value: -1}
			}
			runeIdx := len([]rune(str.Value[:idx]))
			return &Integer{Value: int64(runeIdx)}
		},
	}
	definePublicFunction(env, "indeks_terakhir", indeksTerakhirFn)
	definePublicFunction(env, "last_index_of", indeksTerakhirFn)

	// hitung / count(str, sub)
	hitungFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "teks.hitung mengharapkan 2 argumen: teks dan teks_cari", "")
			}
			str, ok1 := args[0].(*String)
			sub, ok2 := args[1].(*String)
			if !ok1 || !ok2 {
				return newError(nil, "argumen harus berupa teks", "")
			}
			return &Integer{Value: int64(strings.Count(str.Value, sub.Value))}
		},
	}
	definePublicFunction(env, "hitung", hitungFn)
	definePublicFunction(env, "count", hitungFn)

	// ulangi / repeat(str, n)
	ulangiFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "teks.ulangi mengharapkan 2 argumen: teks dan jumlah_pengulangan", "")
			}
			str, ok1 := args[0].(*String)
			cnt, ok2 := args[1].(*Integer)
			if !ok1 || !ok2 {
				return newError(nil, "argumen teks.ulangi harus berupa (teks, angka)", "")
			}
			if cnt.Value <= 0 {
				return &String{Value: ""}
			}
			if cnt.Value > 1_000_000 || int64(len(str.Value))*cnt.Value > 50_000_000 {
				return newError(nil, "hasil pengulangan teks melebihi batas aman (maksimal 50MB)", "")
			}
			return &String{Value: strings.Repeat(str.Value, int(cnt.Value))}
		},
	}
	definePublicFunction(env, "ulangi", ulangiFn)
	definePublicFunction(env, "repeat", ulangiFn)

	// -------------------------------------------------------------
	// 3. Transformasi Gaya Huruf (Case Styles)
	// -------------------------------------------------------------

	// kapital_awal / capitalize("celang") -> "Celang"
	kapitalAwalFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "teks.kapital_awal mengharapkan 1 argumen teks", "")
			}
			str, ok := args[0].(*String)
			if !ok {
				return newError(nil, "argumen harus berupa teks", "")
			}
			return &String{Value: capitalizeWord(str.Value)}
		},
	}
	definePublicFunction(env, "kapital_awal", kapitalAwalFn)
	definePublicFunction(env, "capitalize", kapitalAwalFn)

	// judul / title("halo dunia") -> "Halo Dunia"
	judulFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "teks.judul mengharapkan 1 argumen teks", "")
			}
			str, ok := args[0].(*String)
			if !ok {
				return newError(nil, "argumen harus berupa teks", "")
			}
			words := strings.Fields(str.Value)
			for i, w := range words {
				words[i] = capitalizeWord(w)
			}
			return &String{Value: strings.Join(words, " ")}
		},
	}
	definePublicFunction(env, "judul", judulFn)
	definePublicFunction(env, "title", judulFn)

	// unta / camel("halo_dunia") -> "haloDunia"
	untaFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "teks.unta mengharapkan 1 argumen teks", "")
			}
			str, ok := args[0].(*String)
			if !ok {
				return newError(nil, "argumen harus berupa teks", "")
			}
			words := splitIntoWords(str.Value)
			if len(words) == 0 {
				return &String{Value: ""}
			}
			var b strings.Builder
			b.WriteString(strings.ToLower(words[0]))
			for _, w := range words[1:] {
				b.WriteString(capitalizeWord(w))
			}
			return &String{Value: b.String()}
		},
	}
	definePublicFunction(env, "unta", untaFn)
	definePublicFunction(env, "camel", untaFn)

	// ular / snake("haloDunia") -> "halo_dunia"
	ularFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "teks.ular mengharapkan 1 argumen teks", "")
			}
			str, ok := args[0].(*String)
			if !ok {
				return newError(nil, "argumen harus berupa teks", "")
			}
			words := splitIntoWords(str.Value)
			for i, w := range words {
				words[i] = strings.ToLower(w)
			}
			return &String{Value: strings.Join(words, "_")}
		},
	}
	definePublicFunction(env, "ular", ularFn)
	definePublicFunction(env, "snake", ularFn)

	// kebab / kebab("Halo Dunia") -> "halo-dunia"
	kebabFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "teks.kebab mengharapkan 1 argumen teks", "")
			}
			str, ok := args[0].(*String)
			if !ok {
				return newError(nil, "argumen harus berupa teks", "")
			}
			words := splitIntoWords(str.Value)
			for i, w := range words {
				words[i] = strings.ToLower(w)
			}
			return &String{Value: strings.Join(words, "-")}
		},
	}
	definePublicFunction(env, "kebab", kebabFn)

	// -------------------------------------------------------------
	// 4. Perataan & Padding (Alignment)
	// -------------------------------------------------------------

	// isi_kiri / pad_start(str, targetLen, [padStr=" "])
	isiKiriFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 2 || len(args) > 3 {
				return newError(nil, "teks.isi_kiri mengharapkan 2 atau 3 argumen: teks, panjang_target, [karakter_pengisi]", "Contoh: teks.isi_kiri(\"42\", 5, \"0\")")
			}
			str, ok1 := args[0].(*String)
			targetLen, ok2 := args[1].(*Integer)
			if !ok1 || !ok2 {
				return newError(nil, "argumen teks.isi_kiri harus berupa (teks, angka, [teks])", "")
			}
			padStr := " "
			if len(args) == 3 {
				if p, ok := args[2].(*String); ok && p.Value != "" {
					padStr = p.Value
				}
			}
			runes := []rune(str.Value)
			diff := int(targetLen.Value) - len(runes)
			if diff <= 0 {
				return str
			}
			padRunes := []rune(padStr)
			var b strings.Builder
			for i := 0; i < diff; i++ {
				b.WriteRune(padRunes[i%len(padRunes)])
			}
			b.WriteString(str.Value)
			return &String{Value: b.String()}
		},
	}
	definePublicFunction(env, "isi_kiri", isiKiriFn)
	definePublicFunction(env, "pad_start", isiKiriFn)
	definePublicFunction(env, "pad_left", isiKiriFn)

	// isi_kanan / pad_end(str, targetLen, [padStr=" "])
	isiKananFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 2 || len(args) > 3 {
				return newError(nil, "teks.isi_kanan mengharapkan 2 atau 3 argumen: teks, panjang_target, [karakter_pengisi]", "Contoh: teks.isi_kanan(\"item\", 8, \" \")")
			}
			str, ok1 := args[0].(*String)
			targetLen, ok2 := args[1].(*Integer)
			if !ok1 || !ok2 {
				return newError(nil, "argumen teks.isi_kanan harus berupa (teks, angka, [teks])", "")
			}
			padStr := " "
			if len(args) == 3 {
				if p, ok := args[2].(*String); ok && p.Value != "" {
					padStr = p.Value
				}
			}
			runes := []rune(str.Value)
			diff := int(targetLen.Value) - len(runes)
			if diff <= 0 {
				return str
			}
			padRunes := []rune(padStr)
			var b strings.Builder
			b.WriteString(str.Value)
			for i := 0; i < diff; i++ {
				b.WriteRune(padRunes[i%len(padRunes)])
			}
			return &String{Value: b.String()}
		},
	}
	definePublicFunction(env, "isi_kanan", isiKananFn)
	definePublicFunction(env, "pad_end", isiKananFn)
	definePublicFunction(env, "pad_right", isiKananFn)

	// tengah / center(str, targetLen, [padStr=" "])
	tengahFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 2 || len(args) > 3 {
				return newError(nil, "teks.tengah mengharapkan 2 atau 3 argumen: teks, panjang_target, [karakter_pengisi]", "Contoh: teks.tengah(\"Judul\", 11, \"=\")")
			}
			str, ok1 := args[0].(*String)
			targetLen, ok2 := args[1].(*Integer)
			if !ok1 || !ok2 {
				return newError(nil, "argumen teks.tengah harus berupa (teks, angka, [teks])", "")
			}
			padStr := " "
			if len(args) == 3 {
				if p, ok := args[2].(*String); ok && p.Value != "" {
					padStr = p.Value
				}
			}
			runes := []rune(str.Value)
			diff := int(targetLen.Value) - len(runes)
			if diff <= 0 {
				return str
			}
			leftPadLen := diff / 2
			rightPadLen := diff - leftPadLen
			padRunes := []rune(padStr)

			var b strings.Builder
			for i := 0; i < leftPadLen; i++ {
				b.WriteRune(padRunes[i%len(padRunes)])
			}
			b.WriteString(str.Value)
			for i := 0; i < rightPadLen; i++ {
				b.WriteRune(padRunes[i%len(padRunes)])
			}
			return &String{Value: b.String()}
		},
	}
	definePublicFunction(env, "tengah", tengahFn)
	definePublicFunction(env, "center", tengahFn)

	// -------------------------------------------------------------
	// 5. Pola & Pencocokan Reguler (Regex)
	// -------------------------------------------------------------

	// cocok / match(str, pattern) -> boolean
	cocokFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "teks.cocok mengharapkan 2 argumen: teks dan pola_regex", "Contoh: teks.cocok(\"user@mail.com\", \"^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$\")")
			}
			str, ok1 := args[0].(*String)
			pat, ok2 := args[1].(*String)
			if !ok1 || !ok2 {
				return newError(nil, "argumen teks.cocok harus berupa teks", "")
			}
			matched, err := regexp.MatchString(pat.Value, str.Value)
			if err != nil {
				return newError(nil, fmt.Sprintf("Pola regex tidak valid '%s': %v", pat.Value, err), "")
			}
			return nativeBoolToBoolean(matched)
		},
	}
	definePublicFunction(env, "sesuai", cocokFn)
	definePublicFunction(env, "cocok", cocokFn)
	definePublicFunction(env, "apakah_cocok", cocokFn)
	definePublicFunction(env, "match", cocokFn)

	// temukan_semua / find_all(str, pattern) -> Array
	temukanSemuaFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "teks.temukan_semua mengharapkan 2 argumen: teks dan pola_regex", "Contoh: teks.temukan_semua(\"beli 2 apel dan 5 jeruk\", \"[0-9]+\")")
			}
			str, ok1 := args[0].(*String)
			pat, ok2 := args[1].(*String)
			if !ok1 || !ok2 {
				return newError(nil, "argumen teks.temukan_semua harus berupa teks", "")
			}
			re, err := regexp.Compile(pat.Value)
			if err != nil {
				return newError(nil, fmt.Sprintf("Pola regex tidak valid '%s': %v", pat.Value, err), "")
			}
			matches := re.FindAllString(str.Value, -1)
			elements := make([]Object, len(matches))
			for i, m := range matches {
				elements[i] = &String{Value: m}
			}
			return &Array{Elements: elements}
		},
	}
	definePublicFunction(env, "temukan_semua", temukanSemuaFn)
	definePublicFunction(env, "find_all", temukanSemuaFn)

	// ganti_pola / replace_regex(str, pattern, replacement)
	gantiPolaFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 3 {
				return newError(nil, "teks.ganti_pola mengharapkan 3 argumen: teks, pola_regex, pengganti", "Contoh: teks.ganti_pola(\"halo 123\", \"[0-9]+\", \"xxx\")")
			}
			str, ok1 := args[0].(*String)
			pat, ok2 := args[1].(*String)
			rep, ok3 := args[2].(*String)
			if !ok1 || !ok2 || !ok3 {
				return newError(nil, "seluruh argumen teks.ganti_pola harus berupa teks", "")
			}
			re, err := regexp.Compile(pat.Value)
			if err != nil {
				return newError(nil, fmt.Sprintf("Pola regex tidak valid '%s': %v", pat.Value, err), "")
			}
			return &String{Value: re.ReplaceAllString(str.Value, rep.Value)}
		},
	}
	definePublicFunction(env, "ganti_pola", gantiPolaFn)
	definePublicFunction(env, "replace_regex", gantiPolaFn)

	// -------------------------------------------------------------
	// 6. Fitur Cerdas & Unik Celang (Slug, Mask, Similarity, Format)
	// -------------------------------------------------------------

	// slug(judul) -> URL-friendly slug
	slugFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "teks.slug mengharapkan 1 argumen teks", "Contoh: teks.slug(\"Belajar Celang 2.0!\")")
			}
			str, ok := args[0].(*String)
			if !ok {
				return newError(nil, "argumen teks.slug harus berupa teks", "")
			}
			return &String{Value: generateSlug(str.Value)}
		},
	}
	definePublicFunction(env, "slug", slugFn)

	// samarkan / mask(str, [maskChar="*"])
	samarkanFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 1 || len(args) > 2 {
				return newError(nil, "teks.samarkan mengharapkan 1 atau 2 argumen: teks, [karakter_sensor]", "Contoh: teks.samarkan(\"budi@gmail.com\")")
			}
			str, ok := args[0].(*String)
			if !ok {
				return newError(nil, "argumen pertama teks.samarkan harus berupa teks", "")
			}
			maskChar := "*"
			if len(args) == 2 {
				if m, ok := args[1].(*String); ok && m.Value != "" {
					maskChar = m.Value
				}
			}
			return &String{Value: maskSensitiveText(str.Value, maskChar)}
		},
	}
	definePublicFunction(env, "samarkan", samarkanFn)
	definePublicFunction(env, "mask", samarkanFn)

	// kemiripan / similarity(strA, strB) -> Float (0.0 .. 1.0)
	kemiripanFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "teks.kemiripan mengharapkan 2 argumen teks untuk dibandingkan", "Contoh: teks.kemiripan(\"jakarta\", \"jakartaa\")")
			}
			strA, ok1 := args[0].(*String)
			strB, ok2 := args[1].(*String)
			if !ok1 || !ok2 {
				return newError(nil, "kedua argumen teks.kemiripan harus berupa teks", "")
			}
			score := calculateSimilarity(strA.Value, strB.Value)
			return &Float{Value: score}
		},
	}
	definePublicFunction(env, "kemiripan", kemiripanFn)
	definePublicFunction(env, "similarity", kemiripanFn)

	// terdekat / closest(target, choicesArray) -> String or NULL
	terdekatFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "teks.terdekat mengharapkan 2 argumen: target teks dan daftar pilihan", "Contoh: teks.terdekat(\"aple\", [\"apel\", \"jeruk\"])")
			}
			target, ok1 := args[0].(*String)
			choices, ok2 := args[1].(*Array)
			if !ok1 || !ok2 {
				return newError(nil, "argumen harus berupa (teks, daftar_pilihan)", "")
			}
			if len(choices.Elements) == 0 {
				return NULL
			}
			bestMatch := ""
			bestScore := -1.0
			for _, el := range choices.Elements {
				choiceStr := el.Inspect()
				if s, ok := el.(*String); ok {
					choiceStr = s.Value
				}
				score := calculateSimilarity(target.Value, choiceStr)
				if score > bestScore {
					bestScore = score
					bestMatch = choiceStr
				}
			}
			return &String{Value: bestMatch}
		},
	}
	definePublicFunction(env, "terdekat", terdekatFn)
	definePublicFunction(env, "closest", terdekatFn)

	// format(template, dataMapOrArgs...)
	formatFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 2 {
				return newError(nil, "teks.format mengharapkan template teks dan data", "Contoh: teks.format(\"Halo {nama}!\", {\"nama\": \"Andi\"})")
			}
			tmplStr, ok := args[0].(*String)
			if !ok {
				return newError(nil, "argumen pertama teks.format harus berupa teks template", "")
			}

			// Kasus A: Argumen kedua adalah Kamus (Map)
			if mp, ok := args[1].(*Map); ok {
				result := tmplStr.Value
				for k, v := range mp.Pairs {
					result = strings.ReplaceAll(result, "{"+k+"}", v.Inspect())
				}
				return &String{Value: result}
			}

			// Kasus B: Argumen kedua adalah Daftar (Array)
			if arr, ok := args[1].(*Array); ok {
				result := tmplStr.Value
				for i, el := range arr.Elements {
					result = strings.ReplaceAll(result, "{"+strconv.Itoa(i)+"}", el.Inspect())
				}
				return &String{Value: result}
			}

			// Kasus C: Argumen posisi ganda: teks.format("{0} {1}", "Halo", "Dunia")
			result := tmplStr.Value
			for i := 1; i < len(args); i++ {
				idx := i - 1
				result = strings.ReplaceAll(result, "{"+strconv.Itoa(idx)+"}", args[i].Inspect())
			}
			return &String{Value: result}
		},
	}
	definePublicFunction(env, "format", formatFn)

	// -------------------------------------------------------------
	// 7. Pemotongan & Pemecahan Canggih (Cut, Fields, Split variants)
	// -------------------------------------------------------------

	// potong_teks / cut(str, sep) -> Map {sebelum, sesudah, ditemukan, before, after, found}
	potongTeksFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "teks.potong_teks mengharapkan 2 argumen: teks dan pemisah", "Contoh: teks.potong_teks(\"kunci:nilai\", \":\")")
			}
			str, ok1 := args[0].(*String)
			sep, ok2 := args[1].(*String)
			if !ok1 || !ok2 {
				return newError(nil, "argumen untuk teks.potong_teks harus berupa teks", "")
			}
			before, after, found := strings.Cut(str.Value, sep.Value)
			bObj := &String{Value: before}
			aObj := &String{Value: after}
			fObj := nativeBool(found)
			return &Map{
				Pairs: map[string]Object{
					"sebelum":   bObj,
					"sesudah":   aObj,
					"ditemukan": fObj,
					"before":    bObj,
					"after":     aObj,
					"found":     fObj,
				},
				Order: []string{"sebelum", "sesudah", "ditemukan"},
			}
		},
	}
	definePublicFunction(env, "potong_teks", potongTeksFn)
	definePublicFunction(env, "cut", potongTeksFn)

	// potong_teks_terakhir / cut_last(str, sep) -> Map
	potongTeksTerakhirFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "teks.potong_teks_terakhir mengharapkan 2 argumen: teks dan pemisah", "Contoh: teks.potong_teks_terakhir(\"a.b.c\", \".\")")
			}
			str, ok1 := args[0].(*String)
			sep, ok2 := args[1].(*String)
			if !ok1 || !ok2 {
				return newError(nil, "argumen untuk teks.potong_teks_terakhir harus berupa teks", "")
			}
			idx := strings.LastIndex(str.Value, sep.Value)
			var before, after string
			found := idx != -1
			if found {
				before = str.Value[:idx]
				after = str.Value[idx+len(sep.Value):]
			} else {
				before = str.Value
				after = ""
			}
			bObj := &String{Value: before}
			aObj := &String{Value: after}
			fObj := nativeBool(found)
			return &Map{
				Pairs: map[string]Object{
					"sebelum":   bObj,
					"sesudah":   aObj,
					"ditemukan": fObj,
					"before":    bObj,
					"after":     aObj,
					"found":     fObj,
				},
				Order: []string{"sebelum", "sesudah", "ditemukan"},
			}
		},
	}
	definePublicFunction(env, "potong_teks_terakhir", potongTeksTerakhirFn)
	definePublicFunction(env, "cut_last", potongTeksTerakhirFn)

	// pangkas_awalan / trim_prefix(str, prefix)
	pangkasAwalanFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "teks.pangkas_awalan mengharapkan 2 argumen: teks dan awalan", "Contoh: teks.pangkas_awalan(\"v1.0.0\", \"v\")")
			}
			str, ok1 := args[0].(*String)
			pref, ok2 := args[1].(*String)
			if !ok1 || !ok2 {
				return newError(nil, "argumen teks.pangkas_awalan harus berupa teks", "")
			}
			return &String{Value: strings.TrimPrefix(str.Value, pref.Value)}
		},
	}
	definePublicFunction(env, "pangkas_awalan", pangkasAwalanFn)
	definePublicFunction(env, "trim_prefix", pangkasAwalanFn)

	// pangkas_akhiran / trim_suffix(str, suffix)
	pangkasAkhiranFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "teks.pangkas_akhiran mengharapkan 2 argumen: teks dan akhiran", "Contoh: teks.pangkas_akhiran(\"file.ce\", \".ce\")")
			}
			str, ok1 := args[0].(*String)
			suff, ok2 := args[1].(*String)
			if !ok1 || !ok2 {
				return newError(nil, "argumen teks.pangkas_akhiran harus berupa teks", "")
			}
			return &String{Value: strings.TrimSuffix(str.Value, suff.Value)}
		},
	}
	definePublicFunction(env, "pangkas_akhiran", pangkasAkhiranFn)
	definePublicFunction(env, "trim_suffix", pangkasAkhiranFn)

	// kata / fields(str)
	kataFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "teks.kata mengharapkan 1 argumen teks", "Contoh: teks.kata(\"satu dua tiga\")")
			}
			str, ok := args[0].(*String)
			if !ok {
				return newError(nil, "argumen teks.kata harus berupa teks", "")
			}
			words := strings.Fields(str.Value)
			elems := make([]Object, len(words))
			for i, w := range words {
				elems[i] = &String{Value: w}
			}
			return &Array{Elements: elems}
		},
	}
	definePublicFunction(env, "kata", kataFn)
	definePublicFunction(env, "fields", kataFn)

	// kata_fungsi / fields_func(str, fn)
	kataFungsiFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "teks.kata_fungsi mengharapkan 2 argumen: teks dan fungsi", "")
			}
			str, ok := args[0].(*String)
			if !ok {
				return newError(nil, "argumen pertama teks.kata_fungsi harus teks", "")
			}
			fn := args[1]
			words := strings.FieldsFunc(str.Value, func(r rune) bool {
				res := applyFunction(fn, []Object{&String{Value: string(r)}}, nil)
				return CoerceToBoolean(res)
			})
			elems := make([]Object, len(words))
			for i, w := range words {
				elems[i] = &String{Value: w}
			}
			return &Array{Elements: elems}
		},
	}
	definePublicFunction(env, "kata_fungsi", kataFungsiFn)
	definePublicFunction(env, "fields_func", kataFungsiFn)

	// sama_abaikan_kasus / equal_fold(s1, s2)
	samaAbaikanKasusFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "teks.sama_abaikan_kasus mengharapkan 2 argumen teks", "Contoh: teks.sama_abaikan_kasus(\"abc\", \"ABC\")")
			}
			s1, ok1 := args[0].(*String)
			s2, ok2 := args[1].(*String)
			if !ok1 || !ok2 {
				return newError(nil, "kedua argumen teks.sama_abaikan_kasus harus berupa teks", "")
			}
			return nativeBool(strings.EqualFold(s1.Value, s2.Value))
		},
	}
	definePublicFunction(env, "sama_abaikan_kasus", samaAbaikanKasusFn)
	definePublicFunction(env, "equal_fold", samaAbaikanKasusFn)

	// berisi_karakter / contains_any(str, chars)
	berisiKarakterFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "teks.berisi_karakter mengharapkan 2 argumen: teks dan kumpulan_karakter", "")
			}
			str, ok1 := args[0].(*String)
			chars, ok2 := args[1].(*String)
			if !ok1 || !ok2 {
				return newError(nil, "argumen teks.berisi_karakter harus berupa teks", "")
			}
			return nativeBool(strings.ContainsAny(str.Value, chars.Value))
		},
	}
	definePublicFunction(env, "berisi_karakter", berisiKarakterFn)
	definePublicFunction(env, "contains_any", berisiKarakterFn)

	// indeks_karakter / index_any(str, chars)
	indeksKarakterFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "teks.indeks_karakter mengharapkan 2 argumen: teks dan kumpulan_karakter", "")
			}
			str, ok1 := args[0].(*String)
			chars, ok2 := args[1].(*String)
			if !ok1 || !ok2 {
				return newError(nil, "argumen teks.indeks_karakter harus berupa teks", "")
			}
			idx := strings.IndexAny(str.Value, chars.Value)
			if idx == -1 {
				return &Integer{Value: -1}
			}
			runeIdx := len([]rune(str.Value[:idx]))
			return &Integer{Value: int64(runeIdx)}
		},
	}
	definePublicFunction(env, "indeks_karakter", indeksKarakterFn)
	definePublicFunction(env, "index_any", indeksKarakterFn)

	// indeks_karakter_terakhir / last_index_any(str, chars)
	indeksKarakterTerakhirFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "teks.indeks_karakter_terakhir mengharapkan 2 argumen: teks dan kumpulan_karakter", "")
			}
			str, ok1 := args[0].(*String)
			chars, ok2 := args[1].(*String)
			if !ok1 || !ok2 {
				return newError(nil, "argumen teks.indeks_karakter_terakhir harus berupa teks", "")
			}
			idx := strings.LastIndexAny(str.Value, chars.Value)
			if idx == -1 {
				return &Integer{Value: -1}
			}
			runeIdx := len([]rune(str.Value[:idx]))
			return &Integer{Value: int64(runeIdx)}
		},
	}
	definePublicFunction(env, "indeks_karakter_terakhir", indeksKarakterTerakhirFn)
	definePublicFunction(env, "last_index_any", indeksKarakterTerakhirFn)

	// berisi_fungsi / contains_func(str, fn)
	berisiFungsiFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "teks.berisi_fungsi mengharapkan 2 argumen: teks dan fungsi", "")
			}
			str, ok := args[0].(*String)
			if !ok {
				return newError(nil, "argumen pertama teks.berisi_fungsi harus berupa teks", "")
			}
			fn := args[1]
			for _, r := range str.Value {
				res := applyFunction(fn, []Object{&String{Value: string(r)}}, nil)
				if CoerceToBoolean(res) {
					return TRUE
				}
			}
			return FALSE
		},
	}
	definePublicFunction(env, "berisi_fungsi", berisiFungsiFn)
	definePublicFunction(env, "contains_func", berisiFungsiFn)

	// indeks_fungsi / index_func(str, fn)
	indeksFungsiFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "teks.indeks_fungsi mengharapkan 2 argumen: teks dan fungsi", "")
			}
			str, ok := args[0].(*String)
			if !ok {
				return newError(nil, "argumen pertama teks.indeks_fungsi harus berupa teks", "")
			}
			fn := args[1]
			runes := []rune(str.Value)
			for i, r := range runes {
				res := applyFunction(fn, []Object{&String{Value: string(r)}}, nil)
				if CoerceToBoolean(res) {
					return &Integer{Value: int64(i)}
				}
			}
			return &Integer{Value: -1}
		},
	}
	definePublicFunction(env, "indeks_fungsi", indeksFungsiFn)
	definePublicFunction(env, "index_func", indeksFungsiFn)

	// indeks_fungsi_terakhir / last_index_func(str, fn)
	indeksFungsiTerakhirFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "teks.indeks_fungsi_terakhir mengharapkan 2 argumen: teks dan fungsi", "")
			}
			str, ok := args[0].(*String)
			if !ok {
				return newError(nil, "argumen pertama teks.indeks_fungsi_terakhir harus berupa teks", "")
			}
			fn := args[1]
			runes := []rune(str.Value)
			for i := len(runes) - 1; i >= 0; i-- {
				res := applyFunction(fn, []Object{&String{Value: string(runes[i])}}, nil)
				if CoerceToBoolean(res) {
					return &Integer{Value: int64(i)}
				}
			}
			return &Integer{Value: -1}
		},
	}
	definePublicFunction(env, "indeks_fungsi_terakhir", indeksFungsiTerakhirFn)
	definePublicFunction(env, "last_index_func", indeksFungsiTerakhirFn)

	// ganti_banyak / replace_many(str, kamusAtauArray)
	gantiBanyakFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "teks.ganti_banyak mengharapkan 2 argumen: teks dan kamus_pengganti (atau array pasangan)", "Contoh: teks.ganti_banyak(\"a b\", {\"a\": \"1\", \"b\": \"2\"})")
			}
			str, ok := args[0].(*String)
			if !ok {
				return newError(nil, "argumen pertama teks.ganti_banyak harus berupa teks", "")
			}
			var pairs []string
			if mp, ok := args[1].(*Map); ok {
				for _, k := range mp.GetKeys() {
					pairs = append(pairs, k, mp.Pairs[k].Inspect())
					if s, ok := mp.Pairs[k].(*String); ok {
						pairs[len(pairs)-1] = s.Value
					}
				}
			} else if arr, ok := args[1].(*Array); ok {
				for _, el := range arr.Elements {
					if subArr, ok := el.(*Array); ok && len(subArr.Elements) >= 2 {
						kStr := subArr.Elements[0].Inspect()
						if s, ok := subArr.Elements[0].(*String); ok {
							kStr = s.Value
						}
						vStr := subArr.Elements[1].Inspect()
						if s, ok := subArr.Elements[1].(*String); ok {
							vStr = s.Value
						}
						pairs = append(pairs, kStr, vStr)
					}
				}
			} else {
				return newError(nil, "argumen kedua teks.ganti_banyak harus berupa Kamus atau Array pasangan", "")
			}
			if len(pairs) == 0 {
				return str
			}
			replacer := strings.NewReplacer(pairs...)
			return &String{Value: replacer.Replace(str.Value)}
		},
	}
	definePublicFunction(env, "ganti_banyak", gantiBanyakFn)
	definePublicFunction(env, "replace_many", gantiBanyakFn)

	// pisah_n / split_n(str, sep, n)
	pisahNFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 3 {
				return newError(nil, "teks.pisah_n mengharapkan 3 argumen: teks, pemisah, batas_jumlah", "Contoh: teks.pisah_n(\"a:b:c\", \":\", 2)")
			}
			str, ok1 := args[0].(*String)
			sep, ok2 := args[1].(*String)
			nVal, ok3 := args[2].(*Integer)
			if !ok1 || !ok2 || !ok3 {
				return newError(nil, "argumen teks.pisah_n harus berupa (teks, teks, angka)", "")
			}
			parts := strings.SplitN(str.Value, sep.Value, int(nVal.Value))
			elems := make([]Object, len(parts))
			for i, p := range parts {
				elems[i] = &String{Value: p}
			}
			return &Array{Elements: elems}
		},
	}
	definePublicFunction(env, "pisah_n", pisahNFn)
	definePublicFunction(env, "split_n", pisahNFn)

	// pisah_setelah / split_after(str, sep)
	pisahSetelahFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "teks.pisah_setelah mengharapkan 2 argumen: teks dan pemisah", "")
			}
			str, ok1 := args[0].(*String)
			sep, ok2 := args[1].(*String)
			if !ok1 || !ok2 {
				return newError(nil, "argumen teks.pisah_setelah harus berupa teks", "")
			}
			parts := strings.SplitAfter(str.Value, sep.Value)
			elems := make([]Object, len(parts))
			for i, p := range parts {
				elems[i] = &String{Value: p}
			}
			return &Array{Elements: elems}
		},
	}
	definePublicFunction(env, "pisah_setelah", pisahSetelahFn)
	definePublicFunction(env, "split_after", pisahSetelahFn)

	// pisah_setelah_n / split_after_n(str, sep, n)
	pisahSetelahNFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 3 {
				return newError(nil, "teks.pisah_setelah_n mengharapkan 3 argumen: teks, pemisah, batas", "")
			}
			str, ok1 := args[0].(*String)
			sep, ok2 := args[1].(*String)
			nVal, ok3 := args[2].(*Integer)
			if !ok1 || !ok2 || !ok3 {
				return newError(nil, "argumen harus berupa (teks, teks, angka)", "")
			}
			parts := strings.SplitAfterN(str.Value, sep.Value, int(nVal.Value))
			elems := make([]Object, len(parts))
			for i, p := range parts {
				elems[i] = &String{Value: p}
			}
			return &Array{Elements: elems}
		},
	}
	definePublicFunction(env, "pisah_setelah_n", pisahSetelahNFn)
	definePublicFunction(env, "split_after_n", pisahSetelahNFn)

	// ke_valid_utf8 / to_valid_utf8(str, replacement)
	keValidUtf8Fn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "teks.ke_valid_utf8 mengharapkan 2 argumen: teks dan teks_pengganti", "")
			}
			str, ok1 := args[0].(*String)
			rep, ok2 := args[1].(*String)
			if !ok1 || !ok2 {
				return newError(nil, "argumen harus berupa teks", "")
			}
			return &String{Value: strings.ToValidUTF8(str.Value, rep.Value)}
		},
	}
	definePublicFunction(env, "ke_valid_utf8", keValidUtf8Fn)
	definePublicFunction(env, "to_valid_utf8", keValidUtf8Fn)
}

// -------------------------------------------------------------
// Helper Functions (String Manipulations, Similarity, Slug, Mask)
// -------------------------------------------------------------

func capitalizeWord(s string) string {
	if len(s) == 0 {
		return ""
	}
	runes := []rune(s)
	return strings.ToUpper(string(runes[0])) + strings.ToLower(string(runes[1:]))
}

func splitIntoWords(s string) []string {
	var words []string
	var cur strings.Builder
	runes := []rune(s)

	for i := 0; i < len(runes); i++ {
		r := runes[i]
		if r == '_' || r == '-' || r == ' ' || r == '.' || r == '/' || r == '\\' {
			if cur.Len() > 0 {
				words = append(words, cur.String())
				cur.Reset()
			}
			continue
		}
		if i > 0 && unicode.IsUpper(r) && !unicode.IsUpper(runes[i-1]) {
			if cur.Len() > 0 {
				words = append(words, cur.String())
				cur.Reset()
			}
		}
		cur.WriteRune(r)
	}
	if cur.Len() > 0 {
		words = append(words, cur.String())
	}
	return words
}

func generateSlug(s string) string {
	s = strings.ReplaceAll(s, "&", " dan ")
	s = strings.ToLower(s)
	var b strings.Builder
	lastDash := false

	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
			lastDash = false
		} else if !lastDash && b.Len() > 0 {
			b.WriteRune('-')
			lastDash = true
		}
	}
	return strings.Trim(b.String(), "-")
}

func maskSensitiveText(s string, maskChar string) string {
	if maskChar == "" {
		maskChar = "*"
	}
	if strings.Contains(s, "@") {
		parts := strings.SplitN(s, "@", 2)
		user, domain := parts[0], parts[1]
		uRunes := []rune(user)
		if len(uRunes) <= 2 {
			return string(uRunes[0]) + strings.Repeat(maskChar, 3) + "@" + domain
		}
		masked := string(uRunes[0]) + strings.Repeat(maskChar, 3) + string(uRunes[len(uRunes)-1])
		return masked + "@" + domain
	}

	runes := []rune(s)
	l := len(runes)
	if l <= 4 {
		return strings.Repeat(maskChar, l)
	}
	if l >= 10 {
		middleCount := l - 8
		if middleCount < 4 {
			middleCount = 4
		}
		return string(runes[:4]) + strings.Repeat(maskChar, middleCount) + string(runes[l-4:])
	}
	return string(runes[:1]) + strings.Repeat(maskChar, l-2) + string(runes[l-1:])
}

func calculateSimilarity(s1, s2 string) float64 {
	if s1 == s2 {
		return 1.0
	}
	r1, r2 := []rune(s1), []rune(s2)
	maxLen := len(r1)
	if len(r2) > maxLen {
		maxLen = len(r2)
	}
	if maxLen == 0 {
		return 1.0
	}
	dist := levenshtein(s1, s2)
	sim := 1.0 - (float64(dist) / float64(maxLen))
	if sim < 0 {
		return 0.0
	}
	return sim
}
