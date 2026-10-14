package evaluator

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"celang/pkg/parser"
)

func runCelangTestCode(t *testing.T, input string) Object {
	t.Helper()
	l := parser.NewLexer(input)
	p := parser.NewParser(l)
	prog, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("parser error: %v", err)
	}
	env := NewEnvironment()
	var out bytes.Buffer
	RegisterBuiltins(env, &out)
	return Eval(prog, env)
}

func TestStdlib_BerkasAndFS(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")

	// 1. Tulis berkas (Write)
	codeWrite := `
	buat sukses = berkas.tulis("` + testFile + `", "Halo Celang!")
	sukses
	`
	res := runCelangTestCode(t, codeWrite)
	if boolVal, ok := res.(*Boolean); !ok || !boolVal.Value {
		t.Fatalf("berkas.tulis failed: %s", res.Inspect())
	}

	// 2. Cek keberadaan (Exists)
	codeAda := `
	berkas.ada("` + testFile + `")
	`
	res = runCelangTestCode(t, codeAda)
	if boolVal, ok := res.(*Boolean); !ok || !boolVal.Value {
		t.Fatalf("berkas.ada expected true, got: %s", res.Inspect())
	}

	// 3. Baca berkas (Read)
	codeBaca := `
	berkas.baca("` + testFile + `")
	`
	res = runCelangTestCode(t, codeBaca)
	if strVal, ok := res.(*String); !ok || strVal.Value != "Halo Celang!" {
		t.Fatalf("berkas.baca expected 'Halo Celang!', got: %s", res.Inspect())
	}

	// 4. Tambah isi berkas (Append) via English alias fs.append
	codeAppend := `
	fs.append("` + testFile + `", " Tambahan.")
	fs.read("` + testFile + `")
	`
	res = runCelangTestCode(t, codeAppend)
	if strVal, ok := res.(*String); !ok || strVal.Value != "Halo Celang! Tambahan." {
		t.Fatalf("fs.read after fs.append got: %s", res.Inspect())
	}

	// 5. Hapus berkas (Remove)
	codeHapus := `
	berkas.hapus("` + testFile + `")
	berkas.ada("` + testFile + `")
	`
	res = runCelangTestCode(t, codeHapus)
	if boolVal, ok := res.(*Boolean); !ok || boolVal.Value {
		t.Fatalf("berkas.ada after hapus expected false, got: %s", res.Inspect())
	}

	// 6. Empathetic Error on reading non-existent file
	codeError := `
	coba {
		berkas.baca("` + filepath.Join(tmpDir, "tidak_ada.txt") + `")
	} tangkap (err) {
		err
	}
	`
	res = runCelangTestCode(t, codeError)
	if !strings.Contains(res.Inspect(), "tidak ada") && !strings.Contains(res.Inspect(), "Gagal membaca berkas") {
		t.Fatalf("expected empathetic error message, got: %s", res.Inspect())
	}
}

func TestStdlib_Berkas_AdvancedOperations(t *testing.T) {
	tmpDir := t.TempDir()
	folder1 := filepath.Join(tmpDir, "folder1")
	subFolder := filepath.Join(folder1, "sub")
	fileA := filepath.Join(subFolder, "a.txt")
	fileB := filepath.Join(subFolder, "b.txt")

	// 1. Buat folder rekursif
	codeMkdir := `
	berkas.buat_folder("` + subFolder + `")
	berkas.apakah_folder("` + subFolder + `")
	`
	res := runCelangTestCode(t, codeMkdir)
	if boolVal, ok := res.(*Boolean); !ok || !boolVal.Value {
		t.Fatalf("buat_folder failed: %s", res.Inspect())
	}

	// 2. Tulis file dan cek ukuran & apakah_berkas
	codeWrite := `
	berkas.tulis("` + fileA + `", "12345")
	buat uk = berkas.ukuran("` + fileA + `")
	buat isF = berkas.apakah_berkas("` + fileA + `")
	buat isDir = berkas.apakah_folder("` + fileA + `")
	[uk, isF, isDir]
	`
	res = runCelangTestCode(t, codeWrite)
	if arr, ok := res.(*Array); !ok || arr.Elements[0].Inspect() != "5" || arr.Elements[1].Inspect() != "true" || arr.Elements[2].Inspect() != "false" {
		t.Fatalf("ukuran / check failed: %s", res.Inspect())
	}

	// 3. Salin file
	codeCopy := `
	berkas.salin("` + fileA + `", "` + fileB + `")
	berkas.baca("` + fileB + `")
	`
	res = runCelangTestCode(t, codeCopy)
	if s, ok := res.(*String); !ok || s.Value != "12345" {
		t.Fatalf("salin failed: %s", res.Inspect())
	}

	// 4. Daftar isi folder
	codeList := `
	buat d = berkas.daftar("` + subFolder + `")
	d
	`
	res = runCelangTestCode(t, codeList)
	arr, ok := res.(*Array)
	if !ok || len(arr.Elements) != 2 {
		t.Fatalf("daftar expected 2 files, got: %s", res.Inspect())
	}

	// 5. Pindah file
	fileC := filepath.Join(subFolder, "c.txt")
	codeMove := `
	berkas.pindah("` + fileB + `", "` + fileC + `")
	[berkas.ada("` + fileB + `"), berkas.ada("` + fileC + `")]
	`
	res = runCelangTestCode(t, codeMove)
	arrMove, ok := res.(*Array)
	if !ok || arrMove.Elements[0].Inspect() != "false" || arrMove.Elements[1].Inspect() != "true" {
		t.Fatalf("pindah failed: %s", res.Inspect())
	}

	// 6. Info / stat
	codeInfo := `
	buat inf = berkas.info("` + fileC + `")
	[inf["nama"], inf["ukuran"], inf["apakah_berkas"]]
	`
	res = runCelangTestCode(t, codeInfo)
	arrInfo, ok := res.(*Array)
	if !ok || arrInfo.Elements[0].Inspect() != "c.txt" || arrInfo.Elements[1].Inspect() != "5" || arrInfo.Elements[2].Inspect() != "true" {
		t.Fatalf("info failed: %s", res.Inspect())
	}

	// 7. Hapus semua rekursif
	codeRmall := `
	berkas.hapus_semua("` + folder1 + `")
	berkas.ada("` + folder1 + `")
	`
	res = runCelangTestCode(t, codeRmall)
	if boolVal, ok := res.(*Boolean); !ok || boolVal.Value {
		t.Fatalf("hapus_semua failed, folder still exists: %s", res.Inspect())
	}
}

func TestStdlib_Berkas_PathHelpers(t *testing.T) {
	code := `
	buat g = berkas.gabung("folder", "sub", "file.txt")
	buat n = berkas.nama(g)
	buat ext = berkas.ekstensi(g)
	buat ind = berkas.induk(g)
	[n, ext, ind]
	`
	res := runCelangTestCode(t, code)
	arr, ok := res.(*Array)
	if !ok {
		t.Fatalf("path helpers failed: %s", res.Inspect())
	}
	if arr.Elements[0].Inspect() != "file.txt" {
		t.Errorf("expected file.txt, got %s", arr.Elements[0].Inspect())
	}
	if arr.Elements[1].Inspect() != ".txt" {
		t.Errorf("expected .txt, got %s", arr.Elements[1].Inspect())
	}
	if arr.Elements[2].Inspect() != "folder/sub" && arr.Elements[2].Inspect() != "folder\\sub" {
		t.Errorf("expected folder/sub, got %s", arr.Elements[2].Inspect())
	}
}

func TestStdlib_Berkas_SmartLoadSaveJSON(t *testing.T) {
	tmpDir := t.TempDir()
	jsonFile := filepath.Join(tmpDir, "config.json")

	// 1. Simpan objek Map langsung ke .json
	codeSave := `
	buat cfg = {"nama": "Celang", "versi": 2, "fitur": ["fs", "concurrency"]}
	berkas.simpan("` + jsonFile + `", cfg)
	`
	res := runCelangTestCode(t, codeSave)
	if boolVal, ok := res.(*Boolean); !ok || !boolVal.Value {
		t.Fatalf("simpan failed: %s", res.Inspect())
	}

	// 2. Muat otomatis dari .json langsung menjadi Objek Celang tanpa parse manual!
	codeLoad := `
	buat dimuat = berkas.muat("` + jsonFile + `")
	[dimuat["nama"], dimuat["versi"], dimuat["fitur"][0]]
	`
	res = runCelangTestCode(t, codeLoad)
	arr, ok := res.(*Array)
	if !ok || arr.Elements[0].Inspect() != "Celang" || arr.Elements[1].Inspect() != "2" || arr.Elements[2].Inspect() != "fs" {
		t.Fatalf("muat json failed: %s", res.Inspect())
	}
}

func TestStdlib_Berkas_TypoSuggestion(t *testing.T) {
	tmpDir := t.TempDir()
	actualFile := filepath.Join(tmpDir, "pengaturan_sistem.json")
	typoPath := filepath.Join(tmpDir, "pengaturan_sistemm.json")

	_ = os.WriteFile(actualFile, []byte("{}"), 0644)

	code := `
	coba {
		berkas.baca("` + typoPath + `")
	} tangkap (err) {
		err
	}
	`
	res := runCelangTestCode(t, code)
	errMsg := res.Inspect()
	if !strings.Contains(errMsg, "pengaturan_sistem.json") {
		t.Errorf("expected typo suggestion for pengaturan_sistem.json, got: %s", errMsg)
	}
}

func TestStdlib_Berkas_StreamCSP(t *testing.T) {
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "stream.log")
	content := "baris 1: mulai\nbaris 2: proses\nbaris 3: selesai\n"
	_ = os.WriteFile(logFile, []byte(content), 0644)

	code := `
	buat hasil = []
	untuk baris dalam berkas.alir("` + logFile + `") {
		hasil.tambah(baris)
	}
	hasil
	`
	res := runCelangTestCode(t, code)
	arr, ok := res.(*Array)
	if !ok || len(arr.Elements) != 3 {
		t.Fatalf("stream CSP expected 3 lines, got: %s", res.Inspect())
	}
	if arr.Elements[0].Inspect() != "baris 1: mulai" || arr.Elements[2].Inspect() != "baris 3: selesai" {
		t.Errorf("unexpected stream content: %s", res.Inspect())
	}
}

func TestStdlib_Berkas_GantiTeks(t *testing.T) {
	tmpDir := t.TempDir()
	f := filepath.Join(tmpDir, "config.env")
	_ = os.WriteFile(f, []byte("PORT=3000\nENV=dev\n"), 0644)

	code := `
	berkas.ganti_teks("` + f + `", "ENV=dev", "ENV=prod")
	berkas.baca("` + f + `")
	`
	res := runCelangTestCode(t, code)
	if s, ok := res.(*String); !ok || !strings.Contains(s.Value, "ENV=prod") {
		t.Fatalf("ganti_teks failed: %s", res.Inspect())
	}
}

func TestStdlib_Berkas_KepalaDanEkor(t *testing.T) {
	tmpDir := t.TempDir()
	f := filepath.Join(tmpDir, "log.txt")
	content := "satu\ndua\ntiga\nempat\nlima\n"
	_ = os.WriteFile(f, []byte(content), 0644)

	code := `
	buat awal = berkas.kepala("` + f + `", 2)
	buat akhir = berkas.ekor("` + f + `", 2)
	[awal[0], awal[1], akhir[0], akhir[1]]
	`
	res := runCelangTestCode(t, code)
	arr, ok := res.(*Array)
	if !ok || len(arr.Elements) != 4 {
		t.Fatalf("kepala dan ekor failed: %s", res.Inspect())
	}
	if arr.Elements[0].Inspect() != "satu" || arr.Elements[1].Inspect() != "dua" {
		t.Errorf("kepala mismatch: %s, %s", arr.Elements[0].Inspect(), arr.Elements[1].Inspect())
	}
	if arr.Elements[2].Inspect() != "empat" || arr.Elements[3].Inspect() != "lima" {
		t.Errorf("ekor mismatch: %s, %s", arr.Elements[2].Inspect(), arr.Elements[3].Inspect())
	}
}

func TestStdlib_Berkas_Hash(t *testing.T) {
	tmpDir := t.TempDir()
	f := filepath.Join(tmpDir, "hello.txt")
	_ = os.WriteFile(f, []byte("hello world"), 0644)

	code := `
	buat h_sha = berkas.hash("` + f + `")
	buat h_md5 = berkas.hash("` + f + `", "md5")
	[h_sha, h_md5]
	`
	res := runCelangTestCode(t, code)
	arr, ok := res.(*Array)
	if !ok || len(arr.Elements) != 2 {
		t.Fatalf("hash failed: %s", res.Inspect())
	}
	if arr.Elements[0].Inspect() != "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9" {
		t.Errorf("sha256 mismatch: %s", arr.Elements[0].Inspect())
	}
	if arr.Elements[1].Inspect() != "5eb63bbbe01eeed093cb22bb8f5acdc3" {
		t.Errorf("md5 mismatch: %s", arr.Elements[1].Inspect())
	}
}

func TestStdlib_Berkas_CariGlob(t *testing.T) {
	tmpDir := t.TempDir()
	sub1 := filepath.Join(tmpDir, "src", "controllers")
	sub2 := filepath.Join(tmpDir, "src", "models")
	_ = os.MkdirAll(sub1, 0755)
	_ = os.MkdirAll(sub2, 0755)
	_ = os.WriteFile(filepath.Join(sub1, "user.ce"), []byte("// user"), 0644)
	_ = os.WriteFile(filepath.Join(sub2, "post.ce"), []byte("// post"), 0644)
	_ = os.WriteFile(filepath.Join(tmpDir, "readme.txt"), []byte("txt"), 0644)

	code := `
	buat matches = berkas.cari("` + filepath.ToSlash(filepath.Join(tmpDir, "src", "**", "*.ce")) + `")
	panjang(matches)
	`
	res := runCelangTestCode(t, code)
	if i, ok := res.(*Integer); !ok || i.Value != 2 {
		t.Fatalf("glob expected 2 matches, got: %s", res.Inspect())
	}
}

func TestStdlib_Berkas_TelusuriWalk(t *testing.T) {
	tmpDir := t.TempDir()
	sub := filepath.Join(tmpDir, "sub")
	_ = os.MkdirAll(sub, 0755)
	_ = os.WriteFile(filepath.Join(sub, "data.json"), []byte("{}"), 0644)
	_ = os.WriteFile(filepath.Join(sub, "notes.txt"), []byte("hi"), 0644)

	code := `
	buat daftar = []
	berkas.telusuri("` + tmpDir + `", fungsi (jalur, info) {
		jika info.ekstensi == ".json" {
			daftar.tambah(info.nama)
		}
	})
	daftar
	`
	res := runCelangTestCode(t, code)
	arr, ok := res.(*Array)
	if !ok || len(arr.Elements) != 1 || arr.Elements[0].Inspect() != "data.json" {
		t.Fatalf("telusuri failed: %s", res.Inspect())
	}
}

func TestStdlib_Berkas_PantauWatch(t *testing.T) {
	tmpDir := t.TempDir()
	watchFile := filepath.Join(tmpDir, "watch_me.txt")
	_ = os.WriteFile(watchFile, []byte("awal"), 0644)

	code := `
	buat pengawas = berkas.pantau("` + tmpDir + `")
	pengawas
	`
	res := runCelangTestCode(t, code)
	ch, ok := res.(*Channel)
	if !ok {
		t.Fatalf("expected Channel from pantau, got: %s", res.Inspect())
	}
	ch.Close()
}

func TestStdlib_MatematikaAndMath(t *testing.T) {
	code := `
	buat r_akar = matematika.akar(25)
	buat r_pangkat = math.pow(2, 4)
	buat r_mutlak = matematika.mutlak(-50)
	buat r_bulat = math.round(4.6)
	buat r_lantai = matematika.lantai(4.9)
	buat r_atap = math.ceil(4.1)
	buat r_acak = matematika.acak(10, 20)
	buat r_maks = matematika.maks(3, 10, 5)
	buat r_min = math.min(12, 4, 9)
	buat pi = math.PI

	[r_akar, r_pangkat, r_mutlak, r_bulat, r_lantai, r_atap, r_acak >= 10 dan r_acak <= 20, r_maks, r_min, pi > 3.14]
	`
	res := runCelangTestCode(t, code)
	arr, ok := res.(*Array)
	if !ok {
		t.Fatalf("expected array result, got %s", res.Inspect())
	}

	// Check r_akar == 5
	if arr.Elements[0].(*Integer).Value != 5 {
		t.Errorf("akar(25) expected 5, got %s", arr.Elements[0].Inspect())
	}
	// Check r_pangkat == 16
	if arr.Elements[1].(*Integer).Value != 16 {
		t.Errorf("pow(2,4) expected 16, got %s", arr.Elements[1].Inspect())
	}
	// Check r_mutlak == 50
	if arr.Elements[2].(*Integer).Value != 50 {
		t.Errorf("mutlak(-50) expected 50, got %s", arr.Elements[2].Inspect())
	}
	// Check r_bulat == 5
	if arr.Elements[3].(*Integer).Value != 5 {
		t.Errorf("round(4.6) expected 5, got %s", arr.Elements[3].Inspect())
	}
	// Check r_lantai == 4
	if arr.Elements[4].(*Integer).Value != 4 {
		t.Errorf("lantai(4.9) expected 4, got %s", arr.Elements[4].Inspect())
	}
	// Check r_atap == 5
	if arr.Elements[5].(*Integer).Value != 5 {
		t.Errorf("ceil(4.1) expected 5, got %s", arr.Elements[5].Inspect())
	}
	// Check r_acak in range
	if !arr.Elements[6].(*Boolean).Value {
		t.Errorf("acak(10, 20) expected in range [10, 20]")
	}
	// Check r_maks == 10
	if arr.Elements[7].(*Integer).Value != 10 {
		t.Errorf("maks expected 10, got %s", arr.Elements[7].Inspect())
	}
	// Check r_min == 4
	if arr.Elements[8].(*Integer).Value != 4 {
		t.Errorf("min expected 4, got %s", arr.Elements[8].Inspect())
	}
	// Check pi > 3.14 == true
	if !arr.Elements[9].(*Boolean).Value {
		t.Errorf("pi > 3.14 expected true")
	}
}

func TestStdlib_TeksAndStrings(t *testing.T) {
	code := `
	buat parts = teks.pisah("apel,jeruk,mangga", ",")
	buat gabung_str = strings.join(parts, " - ")
	buat kapital_str = teks.kapital("celang")
	buat kecil_str = strings.lower("BAHASA")
	buat potong_str = teks.potong("  rapi  ")
	buat ganti_str = strings.replace("halo kawan", "kawan", "dunia")
	buat ada_kawan = teks.berisi("halo kawan", "kawan")
	buat prefix = teks.dimulai_dengan("celang_keren", "celang")
	buat suffix = strings.ends_with("celang_keren", "keren")

	[gabung_str, kapital_str, kecil_str, potong_str, ganti_str, ada_kawan, prefix, suffix]
	`
	res := runCelangTestCode(t, code)
	arr, ok := res.(*Array)
	if !ok {
		t.Fatalf("expected array, got %s", res.Inspect())
	}

	if arr.Elements[0].(*String).Value != "apel - jeruk - mangga" {
		t.Errorf("join mismatch: %s", arr.Elements[0].Inspect())
	}
	if arr.Elements[1].(*String).Value != "CELANG" {
		t.Errorf("kapital mismatch: %s", arr.Elements[1].Inspect())
	}
	if arr.Elements[2].(*String).Value != "bahasa" {
		t.Errorf("kecil mismatch: %s", arr.Elements[2].Inspect())
	}
	if arr.Elements[3].(*String).Value != "rapi" {
		t.Errorf("potong mismatch: %s", arr.Elements[3].Inspect())
	}
	if arr.Elements[4].(*String).Value != "halo dunia" {
		t.Errorf("ganti mismatch: %s", arr.Elements[4].Inspect())
	}
	if !arr.Elements[5].(*Boolean).Value {
		t.Errorf("berisi expected true")
	}
	if !arr.Elements[6].(*Boolean).Value {
		t.Errorf("dimulai_dengan expected true")
	}
	if !arr.Elements[7].(*Boolean).Value {
		t.Errorf("diakhiri_dengan expected true")
	}
}

func TestStdlib_WaktuAndTime(t *testing.T) {
	code := `
	buat ms1 = waktu.sekarang()
	waktu.tidur(10)
	buat ms2 = time.now()
	buat sec = waktu.sekarang_detik()
	buat strFormat = waktu.format(ms1, "2006")

	[ms2 >= ms1, sec > 0, strFormat != ""]
	`
	res := runCelangTestCode(t, code)
	arr, ok := res.(*Array)
	if !ok {
		t.Fatalf("expected array, got %s", res.Inspect())
	}

	if !arr.Elements[0].(*Boolean).Value {
		t.Errorf("expected ms2 >= ms1")
	}
	if !arr.Elements[1].(*Boolean).Value {
		t.Errorf("expected sec > 0")
	}
	if !arr.Elements[2].(*Boolean).Value {
		t.Errorf("expected strFormat not empty")
	}
}

func TestStdlib_JSON(t *testing.T) {
	code := `
	buat rawJSON = "{\"nama\": \"Celang\", \"versi\": 1, \"aktif\": true, \"fitur\": [\"stdlib\", \"pipeline\"]}"
	buat data = json.uraikan(rawJSON)

	buat nama = data.nama
	buat versi = data.versi
	buat aktif = data.aktif
	buat fitur1 = data.fitur[0]

	buat formatted = json.format(data)
	buat parsedKembali = json.parse(formatted)

	[nama, versi, aktif, fitur1, parsedKembali.nama]
	`
	res := runCelangTestCode(t, code)
	arr, ok := res.(*Array)
	if !ok {
		t.Fatalf("expected array, got %s", res.Inspect())
	}

	if arr.Elements[0].(*String).Value != "Celang" {
		t.Errorf("expected 'Celang', got %s", arr.Elements[0].Inspect())
	}
	if arr.Elements[1].(*Integer).Value != 1 {
		t.Errorf("expected 1, got %s", arr.Elements[1].Inspect())
	}
	if !arr.Elements[2].(*Boolean).Value {
		t.Errorf("expected true, got %s", arr.Elements[2].Inspect())
	}
	if arr.Elements[3].(*String).Value != "stdlib" {
		t.Errorf("expected 'stdlib', got %s", arr.Elements[3].Inspect())
	}
	if arr.Elements[4].(*String).Value != "Celang" {
		t.Errorf("expected 'Celang' on roundtrip, got %s", arr.Elements[4].Inspect())
	}
}

func TestStdlib_VirtualModuleImports(t *testing.T) {
	tmpDir := t.TempDir()
	callerFile := filepath.Join(tmpDir, "main.ce")

	code := `
	ambil { akar, PI } dari "matematika"
	import { split, join } from "strings"
	ambil * sebagai waktuMod dari "waktu"
	import { parse, stringify } from "json"

	buat v_akar = akar(49)
	buat v_teks = join(split("foo:bar", ":"), "-")
	buat v_ms = waktuMod.now()
	buat v_json = parse("{\"ok\": true}").ok

	[v_akar, v_teks, v_ms > 0, v_json, PI > 3.0]
	`
	err := os.WriteFile(callerFile, []byte(code), 0644)
	if err != nil {
		t.Fatalf("failed to write test caller file: %v", err)
	}

	l := parser.NewLexer(code)
	p := parser.NewParser(l)
	prog, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("parser error: %v", err)
	}
	env := NewEnvironment()
	env.FilePath = callerFile
	var out bytes.Buffer
	RegisterBuiltins(env, &out)

	res := Eval(prog, env)
	arr, ok := res.(*Array)
	if !ok {
		t.Fatalf("expected array, got %s", res.Inspect())
	}

	if arr.Elements[0].(*Integer).Value != 7 {
		t.Errorf("expected akar(49) == 7, got %s", arr.Elements[0].Inspect())
	}
	if arr.Elements[1].(*String).Value != "foo-bar" {
		t.Errorf("expected 'foo-bar', got %s", arr.Elements[1].Inspect())
	}
	if !arr.Elements[2].(*Boolean).Value {
		t.Errorf("expected v_ms > 0")
	}
	if !arr.Elements[3].(*Boolean).Value {
		t.Errorf("expected v_json true")
	}
	if !arr.Elements[4].(*Boolean).Value {
		t.Errorf("expected PI > 3.0")
	}
}
