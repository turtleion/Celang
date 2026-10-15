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

func TestStdlib_Matematika_TrigonometryAndAngles(t *testing.T) {
	code := `
	buat rad90 = matematika.ke_radian(90)
	buat sin90 = matematika.bulat(matematika.sin(rad90), 2)
	buat cos90 = matematika.bulat(matematika.cos(rad90), 2)
	buat deg180 = math.to_degrees(math.PI)
	buat atanVal = matematika.atan2(0, 10)

	[sin90, cos90, matematika.bulat(deg180, 0), atanVal]
	`
	res := runCelangTestCode(t, code)
	arr, ok := res.(*Array)
	if !ok {
		t.Fatalf("expected array, got %s", res.Inspect())
	}
	if arr.Elements[0].Inspect() != "1" && arr.Elements[0].Inspect() != "1.0" {
		t.Errorf("sin(90 deg) expected 1, got %s", arr.Elements[0].Inspect())
	}
	if arr.Elements[1].Inspect() != "0" && arr.Elements[1].Inspect() != "0.0" {
		t.Errorf("cos(90 deg) expected 0, got %s", arr.Elements[1].Inspect())
	}
	if arr.Elements[2].Inspect() != "180" {
		t.Errorf("to_degrees(PI) expected 180, got %s", arr.Elements[2].Inspect())
	}
	if arr.Elements[3].Inspect() != "0" && arr.Elements[3].Inspect() != "0.0" {
		t.Errorf("atan2(0, 10) expected 0, got %s", arr.Elements[3].Inspect())
	}
}

func TestStdlib_Matematika_PrecisionAndAritmatika(t *testing.T) {
	code := `
	buat r_bulat = matematika.bulat(3.14159, 2)
	buat r_trunc = matematika.potong_desimal(3.89, 1)
	buat r_clamp1 = matematika.jepit(150, 0, 100)
	buat r_clamp2 = matematika.jepit(-20, 0, 100)
	buat r_sign = matematika.tanda(-42)
	buat r_fact = matematika.faktorial(5)
	buat r_gcd = matematika.fpb(24, 36)
	buat r_lcm = matematika.kpk(4, 6)
	buat r_log = math.log10(100)

	[r_bulat, r_trunc, r_clamp1, r_clamp2, r_sign, r_fact, r_gcd, r_lcm, r_log]
	`
	res := runCelangTestCode(t, code)
	arr, ok := res.(*Array)
	if !ok {
		t.Fatalf("expected array, got %s", res.Inspect())
	}
	if arr.Elements[0].Inspect() != "3.14" {
		t.Errorf("bulat(3.14159, 2) expected 3.14, got %s", arr.Elements[0].Inspect())
	}
	if arr.Elements[1].Inspect() != "3.8" {
		t.Errorf("potong_desimal(3.89, 1) expected 3.8, got %s", arr.Elements[1].Inspect())
	}
	if arr.Elements[2].Inspect() != "100" {
		t.Errorf("jepit(150, 0, 100) expected 100, got %s", arr.Elements[2].Inspect())
	}
	if arr.Elements[3].Inspect() != "0" {
		t.Errorf("jepit(-20, 0, 100) expected 0, got %s", arr.Elements[3].Inspect())
	}
	if arr.Elements[4].Inspect() != "-1" {
		t.Errorf("tanda(-42) expected -1, got %s", arr.Elements[4].Inspect())
	}
	if arr.Elements[5].Inspect() != "120" {
		t.Errorf("faktorial(5) expected 120, got %s", arr.Elements[5].Inspect())
	}
	if arr.Elements[6].Inspect() != "12" {
		t.Errorf("fpb(24, 36) expected 12, got %s", arr.Elements[6].Inspect())
	}
	if arr.Elements[7].Inspect() != "12" {
		t.Errorf("kpk(4, 6) expected 12, got %s", arr.Elements[7].Inspect())
	}
	if arr.Elements[8].Inspect() != "2" && arr.Elements[8].Inspect() != "2.0" {
		t.Errorf("log10(100) expected 2, got %s", arr.Elements[8].Inspect())
	}
}

func TestStdlib_Matematika_Statistics(t *testing.T) {
	code := `
	buat data = [2, 4, 4, 4, 5, 5, 7, 9]
	buat total = matematika.jumlah(data)
	buat rata = matematika.rata_rata(data)
	buat med = matematika.median([1, 5, 2, 8, 7])
	buat mod = matematika.modus(data)
	buat std = matematika.standar_deviasi(data)

	[total, rata, med, mod, matematika.bulat(std, 1)]
	`
	res := runCelangTestCode(t, code)
	arr, ok := res.(*Array)
	if !ok {
		t.Fatalf("expected array, got %s", res.Inspect())
	}
	if arr.Elements[0].Inspect() != "40" {
		t.Errorf("jumlah expected 40, got %s", arr.Elements[0].Inspect())
	}
	if arr.Elements[1].Inspect() != "5" {
		t.Errorf("rata_rata expected 5, got %s", arr.Elements[1].Inspect())
	}
	if arr.Elements[2].Inspect() != "5" {
		t.Errorf("median expected 5, got %s", arr.Elements[2].Inspect())
	}
	if arr.Elements[3].Inspect() != "4" {
		t.Errorf("modus expected 4, got %s", arr.Elements[3].Inspect())
	}
	if arr.Elements[4].Inspect() != "2" && arr.Elements[4].Inspect() != "2.0" {
		t.Errorf("standar_deviasi expected 2, got %s", arr.Elements[4].Inspect())
	}
}

func TestStdlib_Matematika_KillerFeatures(t *testing.T) {
	code := `
	buat ejaan1 = matematika.terbilang(1250500)
	buat ejaan2 = matematika.terbilang(15.5)
	buat rp1 = matematika.rupiah(1500000)
	buat rp2 = matematika.rupiah(1500000, benar)
	buat rp3 = matematika.rupiah(-50000)
	buat pct = matematika.persen(25, 200)
	buat dari_pct = matematika.dari_persen(15, 200000)
	buat nilai_lerp = matematika.lerp(0, 100, 0.5)
	buat mapped = matematika.petakan_rentang(50, 0, 100, 0, 1)
	buat dist2d = matematika.jarak(0, 0, 3, 4)
	buat dist3d = matematika.jarak(0, 0, 0, 2, 3, 6)

	[ejaan1, ejaan2, rp1, rp2, rp3, pct, dari_pct, nilai_lerp, mapped, dist2d, dist3d]
	`
	res := runCelangTestCode(t, code)
	arr, ok := res.(*Array)
	if !ok {
		t.Fatalf("expected array, got %s", res.Inspect())
	}
	if arr.Elements[0].Inspect() != "satu juta dua ratus lima puluh ribu lima ratus" {
		t.Errorf("terbilang(1250500) mismatch: %s", arr.Elements[0].Inspect())
	}
	if arr.Elements[1].Inspect() != "lima belas koma lima" {
		t.Errorf("terbilang(15.5) mismatch: %s", arr.Elements[1].Inspect())
	}
	if arr.Elements[2].Inspect() != "Rp 1.500.000" {
		t.Errorf("rupiah(1500000) mismatch: %s", arr.Elements[2].Inspect())
	}
	if arr.Elements[3].Inspect() != "Rp 1.500.000,00" {
		t.Errorf("rupiah with dec mismatch: %s", arr.Elements[3].Inspect())
	}
	if arr.Elements[4].Inspect() != "-Rp 50.000" {
		t.Errorf("negative rupiah mismatch: %s", arr.Elements[4].Inspect())
	}
	if arr.Elements[5].Inspect() != "12.5" {
		t.Errorf("persen(25, 200) expected 12.5, got %s", arr.Elements[5].Inspect())
	}
	if arr.Elements[6].Inspect() != "30000" {
		t.Errorf("dari_persen(15, 200000) expected 30000, got %s", arr.Elements[6].Inspect())
	}
	if arr.Elements[7].Inspect() != "50" {
		t.Errorf("lerp(0, 100, 0.5) expected 50, got %s", arr.Elements[7].Inspect())
	}
	if arr.Elements[8].Inspect() != "0.5" {
		t.Errorf("petakan_rentang expected 0.5, got %s", arr.Elements[8].Inspect())
	}
	if arr.Elements[9].Inspect() != "5" {
		t.Errorf("jarak 2D expected 5, got %s", arr.Elements[9].Inspect())
	}
	if arr.Elements[10].Inspect() != "7" {
		t.Errorf("jarak 3D expected 7, got %s", arr.Elements[10].Inspect())
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

func TestStdlib_Teks_CaseStyles(t *testing.T) {
	code := `
	buat c_awal = teks.kapital_awal("celang")
	buat c_judul = teks.judul("halo dunia celang")
	buat c_unta = teks.unta("halo_dunia_celang")
	buat c_ular = teks.ular("haloDuniaCelang")
	buat c_kebab = teks.kebab("Halo Dunia Celang")

	[c_awal, c_judul, c_unta, c_ular, c_kebab]
	`
	res := runCelangTestCode(t, code)
	arr, ok := res.(*Array)
	if !ok {
		t.Fatalf("expected array, got %s", res.Inspect())
	}
	if arr.Elements[0].Inspect() != "Celang" {
		t.Errorf("kapital_awal expected 'Celang', got %s", arr.Elements[0].Inspect())
	}
	if arr.Elements[1].Inspect() != "Halo Dunia Celang" {
		t.Errorf("judul expected 'Halo Dunia Celang', got %s", arr.Elements[1].Inspect())
	}
	if arr.Elements[2].Inspect() != "haloDuniaCelang" {
		t.Errorf("unta expected 'haloDuniaCelang', got %s", arr.Elements[2].Inspect())
	}
	if arr.Elements[3].Inspect() != "halo_dunia_celang" {
		t.Errorf("ular expected 'halo_dunia_celang', got %s", arr.Elements[3].Inspect())
	}
	if arr.Elements[4].Inspect() != "halo-dunia-celang" {
		t.Errorf("kebab expected 'halo-dunia-celang', got %s", arr.Elements[4].Inspect())
	}
}

func TestStdlib_Teks_SearchAndPadding(t *testing.T) {
	code := `
	buat idx1 = teks.indeks_dari("halo dunia", "dunia")
	buat idx2 = teks.indeks_terakhir("halo halo", "halo")
	buat cnt = teks.hitung("banana", "an")
	buat rep = teks.ulangi("ha", 3)
	buat pad_kiri = teks.isi_kiri("42", 5, "0")
	buat pad_kanan = teks.isi_kanan("item", 8, " ")
	buat pad_tengah = teks.tengah("Judul", 11, "=")

	[idx1, idx2, cnt, rep, pad_kiri, pad_kanan, pad_tengah]
	`
	res := runCelangTestCode(t, code)
	arr, ok := res.(*Array)
	if !ok {
		t.Fatalf("expected array, got %s", res.Inspect())
	}
	if arr.Elements[0].Inspect() != "5" {
		t.Errorf("indeks_dari expected 5, got %s", arr.Elements[0].Inspect())
	}
	if arr.Elements[1].Inspect() != "5" {
		t.Errorf("indeks_terakhir expected 5, got %s", arr.Elements[1].Inspect())
	}
	if arr.Elements[2].Inspect() != "2" {
		t.Errorf("hitung expected 2, got %s", arr.Elements[2].Inspect())
	}
	if arr.Elements[3].Inspect() != "hahaha" {
		t.Errorf("ulangi expected 'hahaha', got %s", arr.Elements[3].Inspect())
	}
	if arr.Elements[4].Inspect() != "00042" {
		t.Errorf("isi_kiri expected '00042', got %s", arr.Elements[4].Inspect())
	}
	if arr.Elements[5].Inspect() != "item    " {
		t.Errorf("isi_kanan expected 'item    ', got '%s'", arr.Elements[5].Inspect())
	}
	if arr.Elements[6].Inspect() != "===Judul===" {
		t.Errorf("tengah expected '===Judul===', got %s", arr.Elements[6].Inspect())
	}
}

func TestStdlib_Teks_Regex(t *testing.T) {
	code := `
	buat is_email = teks.sesuai("user@mail.com", "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$")
	buat is_not_email = teks.sesuai("bukan-email", "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$")
	buat numbers = teks.temukan_semua("beli 2 apel dan 5 jeruk", "[0-9]+")
	buat replaced = teks.ganti_pola("halo 123 dunia", "[0-9]+", "xxx")

	[is_email, is_not_email, numbers[0], numbers[1], replaced]
	`
	res := runCelangTestCode(t, code)
	arr, ok := res.(*Array)
	if !ok {
		t.Fatalf("expected array, got %s", res.Inspect())
	}
	if arr.Elements[0].Inspect() != "true" {
		t.Errorf("sesuai email expected true")
	}
	if arr.Elements[1].Inspect() != "false" {
		t.Errorf("sesuai not_email expected false")
	}
	if arr.Elements[2].Inspect() != "2" || arr.Elements[3].Inspect() != "5" {
		t.Errorf("temukan_semua expected 2 and 5")
	}
	if arr.Elements[4].Inspect() != "halo xxx dunia" {
		t.Errorf("ganti_pola expected 'halo xxx dunia', got %s", arr.Elements[4].Inspect())
	}
}

func TestStdlib_Teks_SmartFeatures(t *testing.T) {
	code := `
	buat slug_res = teks.slug("Belajar Celang 2.0: Sangat Mudah & Seru!")
	buat mask_email = teks.samarkan("budi.santoso@gmail.com")
	buat mask_phone = teks.samarkan("081234567890")
	buat sim_score = teks.kemiripan("jakarta", "jakartaa")
	buat terdekat_buah = teks.terdekat("aple", ["apel", "jeruk", "mangga"])
	buat fmt_map = teks.format("Halo \{nama}, total \{total}.", {"nama": "Andi", "total": 50000})
	buat fmt_pos = teks.format("Halo \{0}, skor \{1}.", ["Budi", 100])

	[slug_res, mask_email, mask_phone, sim_score > 0.8, terdekat_buah, fmt_map, fmt_pos]
	`
	res := runCelangTestCode(t, code)
	arr, ok := res.(*Array)
	if !ok {
		t.Fatalf("expected array, got %s", res.Inspect())
	}
	if arr.Elements[0].Inspect() != "belajar-celang-2-0-sangat-mudah-dan-seru" {
		t.Errorf("slug expected 'belajar-celang-2-0-sangat-mudah-dan-seru', got %s", arr.Elements[0].Inspect())
	}
	if arr.Elements[1].Inspect() != "b***o@gmail.com" {
		t.Errorf("samarkan email expected 'b***o@gmail.com', got %s", arr.Elements[1].Inspect())
	}
	if arr.Elements[2].Inspect() != "0812****7890" {
		t.Errorf("samarkan phone expected '0812****7890', got %s", arr.Elements[2].Inspect())
	}
	if arr.Elements[3].Inspect() != "true" {
		t.Errorf("kemiripan score expected > 0.8")
	}
	if arr.Elements[4].Inspect() != "apel" {
		t.Errorf("terdekat expected 'apel', got %s", arr.Elements[4].Inspect())
	}
	if arr.Elements[5].Inspect() != "Halo Andi, total 50000." {
		t.Errorf("format expected 'Halo Andi, total 50000.', got %s", arr.Elements[5].Inspect())
	}
	if arr.Elements[6].Inspect() != "Halo Budi, skor 100." {
		t.Errorf("format pos expected 'Halo Budi, skor 100.', got %s", arr.Elements[6].Inspect())
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

func TestStdlib_Waktu_FormattingAndIndo(t *testing.T) {
	code := `
	buat t0 = waktu.urai("2026-10-14 20:30:00")
	buat fmt_iso = waktu.format(t0, "YYYY-MM-DD HH:mm:ss")
	buat fmt_slash = waktu.format(t0, "DD/MM/YYYY")
	buat fmt_id_tgl = waktu.format_indonesia(t0, "tanggal")
	buat fmt_id_hari = waktu.format_indonesia(t0, "hari_tanggal")

	[fmt_iso, fmt_slash, fmt_id_tgl, fmt_id_hari]
	`
	res := runCelangTestCode(t, code)
	arr, ok := res.(*Array)
	if !ok {
		t.Fatalf("expected array, got %s", res.Inspect())
	}
	if arr.Elements[0].Inspect() != "2026-10-14 20:30:00" {
		t.Errorf("format YYYY-MM-DD HH:mm:ss expected '2026-10-14 20:30:00', got %s", arr.Elements[0].Inspect())
	}
	if arr.Elements[1].Inspect() != "14/10/2026" {
		t.Errorf("format DD/MM/YYYY expected '14/10/2026', got %s", arr.Elements[1].Inspect())
	}
	if arr.Elements[2].Inspect() != "14 Oktober 2026" {
		t.Errorf("format_indonesia tanggal expected '14 Oktober 2026', got %s", arr.Elements[2].Inspect())
	}
	if arr.Elements[3].Inspect() != "Rabu, 14 Oktober 2026" {
		t.Errorf("format_indonesia hari_tanggal expected 'Rabu, 14 Oktober 2026', got %s", arr.Elements[3].Inspect())
	}
}

func TestStdlib_Waktu_RelativeTime(t *testing.T) {
	code := `
	buat base = 1700000000000
	buat t_just_now = waktu.lalu(base - 10000, base)
	buat t_mins_ago = waktu.lalu(base - 5 * 60000, base)
	buat t_hours_ago = waktu.lalu(base - 2 * 3600000, base)
	buat t_days_ago = waktu.lalu(base - 3 * 86400000, base)
	buat t_future = waktu.lalu(base + 10 * 60000, base)

	[t_just_now, t_mins_ago, t_hours_ago, t_days_ago, t_future]
	`
	res := runCelangTestCode(t, code)
	arr, ok := res.(*Array)
	if !ok {
		t.Fatalf("expected array, got %s", res.Inspect())
	}
	if arr.Elements[0].Inspect() != "baru saja" {
		t.Errorf("expected 'baru saja', got %s", arr.Elements[0].Inspect())
	}
	if arr.Elements[1].Inspect() != "5 menit yang lalu" {
		t.Errorf("expected '5 menit yang lalu', got %s", arr.Elements[1].Inspect())
	}
	if arr.Elements[2].Inspect() != "2 jam yang lalu" {
		t.Errorf("expected '2 jam yang lalu', got %s", arr.Elements[2].Inspect())
	}
	if arr.Elements[3].Inspect() != "3 hari yang lalu" {
		t.Errorf("expected '3 hari yang lalu', got %s", arr.Elements[3].Inspect())
	}
	if arr.Elements[4].Inspect() != "10 menit lagi" {
		t.Errorf("expected '10 menit lagi', got %s", arr.Elements[4].Inspect())
	}
}

func TestStdlib_Waktu_DateMathAndComponents(t *testing.T) {
	code := `
	buat t0 = waktu.urai("2026-10-14 20:30:00")
	buat info = waktu.info(t0)
	buat t_plus = waktu.tambah(t0, {"hari": 2, "jam": 1})
	buat info_plus = waktu.info(t_plus)
	buat t_minus = waktu.kurang(t0, {"hari": 1})
	buat info_minus = waktu.info(t_minus)
	buat diff = waktu.selisih(t0, t_plus)
	buat kabisat_2024 = waktu.apakah_kabisat(2024)
	buat kabisat_2023 = waktu.apakah_kabisat(2023)

	[info.tahun, info.nama_bulan, info.hari, info_plus.hari, info_plus.jam, info_minus.hari, diff.hari, diff.jam, kabisat_2024, kabisat_2023]
	`
	res := runCelangTestCode(t, code)
	arr, ok := res.(*Array)
	if !ok {
		t.Fatalf("expected array, got %s", res.Inspect())
	}
	if arr.Elements[0].Inspect() != "2026" {
		t.Errorf("info.tahun expected 2026, got %s", arr.Elements[0].Inspect())
	}
	if arr.Elements[1].Inspect() != "Oktober" {
		t.Errorf("info.nama_bulan expected 'Oktober', got %s", arr.Elements[1].Inspect())
	}
	if arr.Elements[2].Inspect() != "14" {
		t.Errorf("info.hari expected 14, got %s", arr.Elements[2].Inspect())
	}
	if arr.Elements[3].Inspect() != "16" {
		t.Errorf("info_plus.hari expected 16, got %s", arr.Elements[3].Inspect())
	}
	if arr.Elements[4].Inspect() != "21" {
		t.Errorf("info_plus.jam expected 21, got %s", arr.Elements[4].Inspect())
	}
	if arr.Elements[5].Inspect() != "13" {
		t.Errorf("info_minus.hari expected 13, got %s", arr.Elements[5].Inspect())
	}
	if arr.Elements[6].Inspect() != "2" {
		t.Errorf("diff.hari expected 2, got %s", arr.Elements[6].Inspect())
	}
	if arr.Elements[7].Inspect() != "1" {
		t.Errorf("diff.jam expected 1, got %s", arr.Elements[7].Inspect())
	}
	if arr.Elements[8].Inspect() != "true" {
		t.Errorf("kabisat 2024 expected true")
	}
	if arr.Elements[9].Inspect() != "false" {
		t.Errorf("kabisat 2023 expected false")
	}
}

func TestStdlib_Waktu_DurasiAndUkur(t *testing.T) {
	code := `
	buat dur_str = waktu.format_durasi(7385000)
	buat ms_detik = waktu.detik(5)
	buat ms_menit = waktu.menit(2)
	buat ms_jam = waktu.jam(1)
	buat ms_hari = waktu.hari(1)
	buat prof = waktu.ukur(fungsi() {
		buat x = 10 + 20
		kembalikan x * 2
	})

	[dur_str, ms_detik, ms_menit, ms_jam, ms_hari, prof.hasil, prof.durasi_ms >= 0]
	`
	res := runCelangTestCode(t, code)
	arr, ok := res.(*Array)
	if !ok {
		t.Fatalf("expected array, got %s", res.Inspect())
	}
	if arr.Elements[0].Inspect() != "2 jam 3 menit 5 detik" {
		t.Errorf("format_durasi expected '2 jam 3 menit 5 detik', got %s", arr.Elements[0].Inspect())
	}
	if arr.Elements[1].Inspect() != "5000" {
		t.Errorf("waktu.detik(5) expected 5000, got %s", arr.Elements[1].Inspect())
	}
	if arr.Elements[2].Inspect() != "120000" {
		t.Errorf("waktu.menit(2) expected 120000, got %s", arr.Elements[2].Inspect())
	}
	if arr.Elements[3].Inspect() != "3600000" {
		t.Errorf("waktu.jam(1) expected 3600000, got %s", arr.Elements[3].Inspect())
	}
	if arr.Elements[4].Inspect() != "86400000" {
		t.Errorf("waktu.hari(1) expected 86400000, got %s", arr.Elements[4].Inspect())
	}
	if arr.Elements[5].Inspect() != "60" {
		t.Errorf("prof.hasil expected 60, got %s", arr.Elements[5].Inspect())
	}
	if arr.Elements[6].Inspect() != "true" {
		t.Errorf("prof.durasi_ms >= 0 expected true")
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

func TestStdlib_JSON_PathAndMutation(t *testing.T) {
	code := `
	buat raw = "{\"pengguna\": {\"nama\": \"Andi\", \"alamat\": {\"kota\": \"Jakarta\"}}, \"hobi\": [\"coding\", \"gaming\"]}"
	buat data = json.uraikan(raw)

	buat kota = json.dapatkan(data, "pengguna.alamat.kota")
	buat hobi2 = json.dapatkan(data, "hobi[1]")
	buat fallback = json.dapatkan(data, "pengguna.kontak.telepon", "0812")

	json.setel(data, "pengguna.alamat.kota", "Bandung")
	buat kota_baru = json.dapatkan(data, "pengguna.alamat.kota")

	json.setel(data, "pengguna.profil.umur", 25)
	buat umur = json.dapatkan(data, "pengguna.profil.umur")

	json.hapus(data, "pengguna.nama")
	buat nama_dihapus = json.dapatkan(data, "pengguna.nama", "KOSONG")

	[kota, hobi2, fallback, kota_baru, umur, nama_dihapus]
	`
	res := runCelangTestCode(t, code)
	arr, ok := res.(*Array)
	if !ok {
		t.Fatalf("expected array, got %s", res.Inspect())
	}
	if arr.Elements[0].Inspect() != "Jakarta" {
		t.Errorf("expected 'Jakarta', got %s", arr.Elements[0].Inspect())
	}
	if arr.Elements[1].Inspect() != "gaming" {
		t.Errorf("expected 'gaming', got %s", arr.Elements[1].Inspect())
	}
	if arr.Elements[2].Inspect() != "0812" {
		t.Errorf("expected '0812', got %s", arr.Elements[2].Inspect())
	}
	if arr.Elements[3].Inspect() != "Bandung" {
		t.Errorf("expected 'Bandung', got %s", arr.Elements[3].Inspect())
	}
	if arr.Elements[4].Inspect() != "25" {
		t.Errorf("expected 25, got %s", arr.Elements[4].Inspect())
	}
	if arr.Elements[5].Inspect() != "KOSONG" {
		t.Errorf("expected 'KOSONG', got %s", arr.Elements[5].Inspect())
	}
}

func TestStdlib_JSON_MergeAndMask(t *testing.T) {
	code := `
	buat defConfig = {"tema": "terang", "log": {"level": "info", "simpan": salah}}
	buat userConfig = {"tema": "gelap", "log": {"simpan": benar, "jalur": "/var/log"}}
	buat merged = json.gabung(defConfig, userConfig)

	buat sensitive = {
		"nama": "Budi",
		"kata_sandi": "rahasia123",
		"kartu_kredit": "4111222233334444",
		"detail": {"api_key": "xyz987"}
	}
	buat masked = json.samarkan(sensitive)

	buat isValid1 = json.apakah_valid("{\"a\": 1}")
	buat isValid2 = json.apakah_valid("\{bukan json}")

	[merged.tema, merged.log.level, merged.log.simpan, merged.log.jalur, masked.kata_sandi, masked.kartu_kredit, masked.detail.api_key, isValid1, isValid2]
	`
	res := runCelangTestCode(t, code)
	arr, ok := res.(*Array)
	if !ok {
		t.Fatalf("expected array, got %s", res.Inspect())
	}
	if arr.Elements[0].Inspect() != "gelap" {
		t.Errorf("merged.tema expected 'gelap', got %s", arr.Elements[0].Inspect())
	}
	if arr.Elements[1].Inspect() != "info" {
		t.Errorf("merged.log.level expected 'info', got %s", arr.Elements[1].Inspect())
	}
	if arr.Elements[2].Inspect() != "true" {
		t.Errorf("merged.log.simpan expected true, got %s", arr.Elements[2].Inspect())
	}
	if arr.Elements[3].Inspect() != "/var/log" {
		t.Errorf("merged.log.jalur expected '/var/log', got %s", arr.Elements[3].Inspect())
	}
	if arr.Elements[4].Inspect() != "***" {
		t.Errorf("masked kata_sandi expected '***', got %s", arr.Elements[4].Inspect())
	}
	if arr.Elements[5].Inspect() != "***" {
		t.Errorf("masked kartu_kredit expected '***', got %s", arr.Elements[5].Inspect())
	}
	if arr.Elements[6].Inspect() != "***" {
		t.Errorf("masked api_key expected '***', got %s", arr.Elements[6].Inspect())
	}
	if arr.Elements[7].Inspect() != "true" {
		t.Errorf("isValid1 expected true")
	}
	if arr.Elements[8].Inspect() != "false" {
		t.Errorf("isValid2 expected false")
	}
}

func TestStdlib_JSON_LinesAndFlatten(t *testing.T) {
	code := `
	buat jsonl = "{\"id\": 1, \"pesan\": \"halo\"}\n{\"id\": 2, \"pesan\": \"dunia\"}\n"
	buat lines = json.uraikan_baris(jsonl)
	buat formatted_lines = json.format_baris(lines)

	buat nested = {"user": {"nama": "Andi", "usia": 20}}
	buat flat = json.ratakan(nested)
	buat unflat = json.susun(flat)

	[lines[0].pesan, lines[1].id, flat["user.nama"], unflat.user.usia]
	`
	res := runCelangTestCode(t, code)
	arr, ok := res.(*Array)
	if !ok {
		t.Fatalf("expected array, got %s", res.Inspect())
	}
	if arr.Elements[0].Inspect() != "halo" {
		t.Errorf("lines[0].pesan expected 'halo', got %s", arr.Elements[0].Inspect())
	}
	if arr.Elements[1].Inspect() != "2" {
		t.Errorf("lines[1].id expected 2, got %s", arr.Elements[1].Inspect())
	}
	if arr.Elements[2].Inspect() != "Andi" {
		t.Errorf("flat['user.nama'] expected 'Andi', got %s", arr.Elements[2].Inspect())
	}
	if arr.Elements[3].Inspect() != "20" {
		t.Errorf("unflat.user.usia expected 20, got %s", arr.Elements[3].Inspect())
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
