package evaluator_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"celang/pkg/evaluator"
	"celang/pkg/parser"
)

func runArchiveScript(t *testing.T, script string, filePath string) (evaluator.Object, *evaluator.Environment) {
	t.Helper()
	l := parser.NewLexer(script)
	p := parser.NewParser(l)
	prog, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("ParseProgram error: %v", err)
	}

	env := evaluator.NewEnvironment()
	env.FilePath = filePath
	env.IsUniversal = prog.IsUniversal
	evaluator.RegisterBuiltins(env, nil)

	res := evaluator.Eval(prog, env)
	return res, env
}

func TestArchive_ZipRoundTrip(t *testing.T) {
	tempDir := t.TempDir()
	srcDir := filepath.Join(tempDir, "src")
	subDir := filepath.Join(srcDir, "sub")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatal(err)
	}

	f1Path := filepath.Join(srcDir, "pesan.txt")
	f2Path := filepath.Join(subDir, "data.bin")
	if err := os.WriteFile(f1Path, []byte("Halo dari arsip ZIP!"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(f2Path, []byte{1, 2, 3, 4, 5, 255}, 0644); err != nil {
		t.Fatal(err)
	}

	zipPath := filepath.Join(tempDir, "output.zip")
	destDir := filepath.Join(tempDir, "extracted")

	script := `
	buat ok_buat = arsip.buat_zip("` + filepath.ToSlash(srcDir) + `", "` + filepath.ToSlash(zipPath) + `")
	buat daftar = arsip.daftar_zip("` + filepath.ToSlash(zipPath) + `")
	buat isi_pesan = arsip.baca_berkas_zip("` + filepath.ToSlash(zipPath) + `", "pesan.txt")
	buat total_ekstrak = arsip.ekstrak_zip("` + filepath.ToSlash(zipPath) + `", "` + filepath.ToSlash(destDir) + `")
	`

	res, env := runArchiveScript(t, script, filepath.Join(tempDir, "main.ce"))
	if res != nil && res.Type() == evaluator.ERROR_OBJ {
		t.Fatalf("Script error: %s", res.Inspect())
	}

	okVal, _ := env.Get("ok_buat")
	if okVal != evaluator.TRUE {
		t.Errorf("Expected ok_buat = true, got %v", okVal)
	}

	daftarVal, ok := env.Get("daftar")
	if !ok || daftarVal.Type() != evaluator.ARRAY_OBJ {
		t.Fatalf("Expected daftar to be Array, got %v", daftarVal)
	}
	arr := daftarVal.(*evaluator.Array)
	if len(arr.Elements) < 2 {
		t.Errorf("Expected at least 2 entries in zip list, got %d", len(arr.Elements))
	}

	isiVal, ok := env.Get("isi_pesan")
	if !ok || isiVal.Type() != evaluator.BYTES_OBJ {
		t.Fatalf("Expected isi_pesan to be Bytes, got %v", isiVal)
	}
	bytesObj := isiVal.(*evaluator.Bytes)
	if string(bytesObj.Elements) != "Halo dari arsip ZIP!" {
		t.Errorf("Expected 'Halo dari arsip ZIP!', got %q", string(bytesObj.Elements))
	}

	extF1 := filepath.Join(destDir, "pesan.txt")
	dataF1, err := os.ReadFile(extF1)
	if err != nil || string(dataF1) != "Halo dari arsip ZIP!" {
		t.Errorf("Extracted pesan.txt mismatch: err=%v, data=%q", err, string(dataF1))
	}

	extF2 := filepath.Join(destDir, "sub", "data.bin")
	dataF2, err := os.ReadFile(extF2)
	if err != nil || len(dataF2) != 6 {
		t.Errorf("Extracted data.bin mismatch: err=%v, len=%d", err, len(dataF2))
	}
}

func TestArchive_ZipFileList(t *testing.T) {
	tempDir := t.TempDir()
	f1 := filepath.Join(tempDir, "f1.txt")
	f2 := filepath.Join(tempDir, "f2.txt")
	os.WriteFile(f1, []byte("konten satu"), 0644)
	os.WriteFile(f2, []byte("konten dua"), 0644)

	zipPath := filepath.Join(tempDir, "bundle.zip")
	destDir := filepath.Join(tempDir, "out")

	script := `
	buat berkas_list = ["` + filepath.ToSlash(f1) + `", "` + filepath.ToSlash(f2) + `"]
	arsip.buat_zip(berkas_list, "` + filepath.ToSlash(zipPath) + `")
	arsip.ekstrak_zip("` + filepath.ToSlash(zipPath) + `", "` + filepath.ToSlash(destDir) + `")
	`

	res, _ := runArchiveScript(t, script, filepath.Join(tempDir, "main.ce"))
	if res != nil && res.Type() == evaluator.ERROR_OBJ {
		t.Fatalf("Script error: %s", res.Inspect())
	}

	out1, _ := os.ReadFile(filepath.Join(destDir, "f1.txt"))
	out2, _ := os.ReadFile(filepath.Join(destDir, "f2.txt"))
	if string(out1) != "konten satu" || string(out2) != "konten dua" {
		t.Errorf("Mismatch in extracted files from list: %q, %q", string(out1), string(out2))
	}
}

func TestArchive_TarGzRoundTrip(t *testing.T) {
	tempDir := t.TempDir()
	srcDir := filepath.Join(tempDir, "pkg_source")
	os.MkdirAll(srcDir, 0755)
	os.WriteFile(filepath.Join(srcDir, "modul.ce"), []byte(`publik! buat versi = "1.0.0"`), 0644)
	os.WriteFile(filepath.Join(srcDir, "baca.md"), []byte("# Dokumentasi Paket"), 0644)

	tarGzPath := filepath.Join(tempDir, "paket.tar.gz")
	destDir := filepath.Join(tempDir, "pkg_extracted")

	script := `
	buat ok_tar = arsip.buat_tar_gz("` + filepath.ToSlash(srcDir) + `", "` + filepath.ToSlash(tarGzPath) + `")
	buat daftar = arsip.daftar_tar_gz("` + filepath.ToSlash(tarGzPath) + `")
	buat isi_modul = arsip.baca_berkas_tar_gz("` + filepath.ToSlash(tarGzPath) + `", "modul.ce")
	buat jumlah = arsip.ekstrak_tar_gz("` + filepath.ToSlash(tarGzPath) + `", "` + filepath.ToSlash(destDir) + `")
	`

	res, env := runArchiveScript(t, script, filepath.Join(tempDir, "main.ce"))
	if res != nil && res.Type() == evaluator.ERROR_OBJ {
		t.Fatalf("Script error: %s", res.Inspect())
	}

	okVal, _ := env.Get("ok_tar")
	if okVal != evaluator.TRUE {
		t.Errorf("Expected ok_tar = true, got %v", okVal)
	}

	isiVal, _ := env.Get("isi_modul")
	if b, ok := isiVal.(*evaluator.Bytes); !ok || string(b.Elements) != `publik! buat versi = "1.0.0"` {
		t.Errorf("Expected modul.ce content, got %v", isiVal)
	}

	dataModul, err := os.ReadFile(filepath.Join(destDir, "modul.ce"))
	if err != nil || string(dataModul) != `publik! buat versi = "1.0.0"` {
		t.Errorf("Extracted modul.ce error: %v, data: %q", err, string(dataModul))
	}
}

func TestArchive_TarRoundTrip(t *testing.T) {
	tempDir := t.TempDir()
	srcDir := filepath.Join(tempDir, "plain_source")
	os.MkdirAll(srcDir, 0755)
	os.WriteFile(filepath.Join(srcDir, "berkas.txt"), []byte("isi tar polos"), 0644)

	tarPath := filepath.Join(tempDir, "plain.tar")
	destDir := filepath.Join(tempDir, "plain_out")

	script := `
	arsip.buat_tar("` + filepath.ToSlash(srcDir) + `", "` + filepath.ToSlash(tarPath) + `")
	buat daftar = arsip.daftar_tar("` + filepath.ToSlash(tarPath) + `")
	arsip.ekstrak_tar("` + filepath.ToSlash(tarPath) + `", "` + filepath.ToSlash(destDir) + `")
	`

	res, env := runArchiveScript(t, script, filepath.Join(tempDir, "main.ce"))
	if res != nil && res.Type() == evaluator.ERROR_OBJ {
		t.Fatalf("Script error: %s", res.Inspect())
	}

	daftarVal, _ := env.Get("daftar")
	arr := daftarVal.(*evaluator.Array)
	if len(arr.Elements) < 1 {
		t.Errorf("Expected at least 1 entry in plain tar, got %d", len(arr.Elements))
	}

	out, err := os.ReadFile(filepath.Join(destDir, "berkas.txt"))
	if err != nil || string(out) != "isi tar polos" {
		t.Errorf("Extracted plain tar error: %v, data: %q", err, string(out))
	}
}

func TestArchive_GzipAndZlib(t *testing.T) {
	tempDir := t.TempDir()
	script := `
	buat teksAsli = "Halo dunia dari algoritma kompresi Celang!"
	buat gzipped = arsip.gzip(teksAsli)
	buat gunzipped = arsip.gunzip(gzipped)
	buat pulihGzip = gunzipped.ke_teks()

	buat zlibbed = arsip.zlib_kompres(teksAsli)
	buat unzlibbed = arsip.zlib_dekompres(zlibbed)
	buat pulihZlib = unzlibbed.ke_teks()

	[pulihGzip == teksAsli, pulihZlib == teksAsli]
	`

	res, _ := runArchiveScript(t, script, filepath.Join(tempDir, "main.ce"))
	if res != nil && res.Type() == evaluator.ERROR_OBJ {
		t.Fatalf("Script error: %s", res.Inspect())
	}

	arr, ok := res.(*evaluator.Array)
	if !ok || len(arr.Elements) != 2 {
		t.Fatalf("Expected 2 bool array, got %v", res)
	}
	if arr.Elements[0] != evaluator.TRUE {
		t.Errorf("Gzip round-trip failed: %v", arr.Elements[0])
	}
	if arr.Elements[1] != evaluator.TRUE {
		t.Errorf("Zlib round-trip failed: %v", arr.Elements[1])
	}
}

func TestArchive_UniversalMode(t *testing.T) {
	tempDir := t.TempDir()
	srcDir := filepath.Join(tempDir, "univ_src")
	os.MkdirAll(srcDir, 0755)
	os.WriteFile(filepath.Join(srcDir, "app.ce"), []byte(`print("universal")`), 0644)
	zipPath := filepath.Join(tempDir, "univ.zip")
	destDir := filepath.Join(tempDir, "univ_out")

	// 1. Valid universal script using 'archive'
	scriptUniv := `"use universal"
const ok_zip = archive.create_zip("` + filepath.ToSlash(srcDir) + `", "` + filepath.ToSlash(zipPath) + `")
const count = archive.extract_zip("` + filepath.ToSlash(zipPath) + `", "` + filepath.ToSlash(destDir) + `")
const compressed = archive.gzip("universal text")
const decompressed = archive.gunzip(compressed)
[ok_zip, count > 0]
`
	resUniv, _ := runArchiveScript(t, scriptUniv, filepath.Join(tempDir, "main_univ.ce"))
	if resUniv != nil && resUniv.Type() == evaluator.ERROR_OBJ {
		t.Fatalf("Universal script error: %s", resUniv.Inspect())
	}
	arr := resUniv.(*evaluator.Array)
	if arr.Elements[0] != evaluator.TRUE || arr.Elements[1] != evaluator.TRUE {
		t.Errorf("Universal assertions failed: %v", arr)
	}

	// 2. Reject Indonesian 'arsip' in Universal Mode
	scriptReject := `"use universal"
const x = arsip.gzip("halo")
`
	resReject, _ := runArchiveScript(t, scriptReject, filepath.Join(tempDir, "reject.ce"))
	if resReject == nil || resReject.Type() != evaluator.ERROR_OBJ {
		t.Fatalf("Expected universal mode to reject 'arsip', got %v", resReject)
	}
	errObj := resReject.(*evaluator.Error)
	if !strings.Contains(errObj.Message, "Indonesian module 'arsip' is not allowed in Universal Mode") {
		t.Errorf("Expected empathy message for arsip, got: %s", errObj.Message)
	}
}

func TestArchive_Detectors(t *testing.T) {
	tempDir := t.TempDir()
	srcDir := filepath.Join(tempDir, "src")
	os.MkdirAll(srcDir, 0755)
	os.WriteFile(filepath.Join(srcDir, "test.txt"), []byte("sample"), 0644)

	zipPath := filepath.Join(tempDir, "sample.zip")
	tarPath := filepath.Join(tempDir, "sample.tar")
	tarGzPath := filepath.Join(tempDir, "sample.tar.gz")
	plainTxtPath := filepath.Join(tempDir, "sample.txt")
	os.WriteFile(plainTxtPath, []byte("bukan arsip"), 0644)

	script := `
	arsip.buat_zip("` + filepath.ToSlash(srcDir) + `", "` + filepath.ToSlash(zipPath) + `")
	arsip.buat_tar("` + filepath.ToSlash(srcDir) + `", "` + filepath.ToSlash(tarPath) + `")
	arsip.buat_tar_gz("` + filepath.ToSlash(srcDir) + `", "` + filepath.ToSlash(tarGzPath) + `")

	buat is_z = arsip.apakah_zip("` + filepath.ToSlash(zipPath) + `")
	buat is_z_false = arsip.apakah_zip("` + filepath.ToSlash(plainTxtPath) + `")

	buat is_t = arsip.apakah_tar("` + filepath.ToSlash(tarPath) + `")
	buat is_t_false = arsip.apakah_tar("` + filepath.ToSlash(tarGzPath) + `")

	buat is_tgz = arsip.apakah_tar_gz("` + filepath.ToSlash(tarGzPath) + `")
	buat is_tgz_false = arsip.apakah_tar_gz("` + filepath.ToSlash(tarPath) + `")

	[is_z, is_z_false, is_t, is_t_false, is_tgz, is_tgz_false]
	`

	res, _ := runArchiveScript(t, script, filepath.Join(tempDir, "detect.ce"))
	if res != nil && res.Type() == evaluator.ERROR_OBJ {
		t.Fatalf("Script error: %s", res.Inspect())
	}
	arr, ok := res.(*evaluator.Array)
	if !ok || len(arr.Elements) != 6 {
		t.Fatalf("Expected 6 elements array, got %v", res)
	}
	if arr.Elements[0] != evaluator.TRUE || arr.Elements[1] != evaluator.FALSE {
		t.Errorf("Zip detector failed: %v, %v", arr.Elements[0], arr.Elements[1])
	}
	if arr.Elements[2] != evaluator.TRUE || arr.Elements[3] != evaluator.FALSE {
		t.Errorf("Tar detector failed: %v, %v", arr.Elements[2], arr.Elements[3])
	}
	if arr.Elements[4] != evaluator.TRUE || arr.Elements[5] != evaluator.FALSE {
		t.Errorf("Tar.gz detector failed: %v, %v", arr.Elements[4], arr.Elements[5])
	}
}

func TestArchive_ExtractSingleFile(t *testing.T) {
	tempDir := t.TempDir()
	srcDir := filepath.Join(tempDir, "src")
	os.MkdirAll(filepath.Join(srcDir, "sub"), 0755)
	os.WriteFile(filepath.Join(srcDir, "readme.txt"), []byte("Dokumentasi Celang"), 0644)
	os.WriteFile(filepath.Join(srcDir, "sub", "rahasia.key"), []byte("KUNCI-12345"), 0644)

	zipPath := filepath.Join(tempDir, "data.zip")
	tarPath := filepath.Join(tempDir, "data.tar")
	tarGzPath := filepath.Join(tempDir, "data.tar.gz")

	outZipDir := filepath.Join(tempDir, "out_zip")
	outTarDir := filepath.Join(tempDir, "out_tar")
	outTarGzDir := filepath.Join(tempDir, "out_targz")

	script := `
	arsip.buat_zip("` + filepath.ToSlash(srcDir) + `", "` + filepath.ToSlash(zipPath) + `")
	arsip.buat_tar("` + filepath.ToSlash(srcDir) + `", "` + filepath.ToSlash(tarPath) + `")
	arsip.buat_tar_gz("` + filepath.ToSlash(srcDir) + `", "` + filepath.ToSlash(tarGzPath) + `")

	// 1. Ekstrak satu berkas ke direktori
	arsip.ekstrak_berkas_zip("` + filepath.ToSlash(zipPath) + `", "readme.txt", "` + filepath.ToSlash(outZipDir) + `/")
	arsip.ekstrak_berkas_tar("` + filepath.ToSlash(tarPath) + `", "sub/rahasia.key", "` + filepath.ToSlash(outTarDir) + `/")
	arsip.ekstrak_berkas_tar_gz("` + filepath.ToSlash(tarGzPath) + `", "readme.txt", "` + filepath.ToSlash(outTarGzDir) + `/")

	// 2. Baca berkas TAR polos
	buat isiTar = arsip.baca_berkas_tar("` + filepath.ToSlash(tarPath) + `", "sub/rahasia.key")

	isiTar.ke_teks()
	`

	res, _ := runArchiveScript(t, script, filepath.Join(tempDir, "single.ce"))
	if res != nil && res.Type() == evaluator.ERROR_OBJ {
		t.Fatalf("Script error: %s", res.Inspect())
	}
	strObj, ok := res.(*evaluator.String)
	if !ok || strObj.Value != "KUNCI-12345" {
		t.Fatalf("Expected 'KUNCI-12345', got %v", res)
	}

	// Cek berkas terekstrak
	fZip, err := os.ReadFile(filepath.Join(outZipDir, "readme.txt"))
	if err != nil || string(fZip) != "Dokumentasi Celang" {
		t.Errorf("ZIP single extract failed: err=%v, data=%q", err, string(fZip))
	}

	fTar, err := os.ReadFile(filepath.Join(outTarDir, "rahasia.key"))
	if err != nil || string(fTar) != "KUNCI-12345" {
		t.Errorf("TAR single extract failed: err=%v, data=%q", err, string(fTar))
	}

	fTarGz, err := os.ReadFile(filepath.Join(outTarGzDir, "readme.txt"))
	if err != nil || string(fTarGz) != "Dokumentasi Celang" {
		t.Errorf("TAR.GZ single extract failed: err=%v, data=%q", err, string(fTarGz))
	}
}

func TestArchive_ExtractionFilter(t *testing.T) {
	tempDir := t.TempDir()
	srcDir := filepath.Join(tempDir, "src")
	os.MkdirAll(srcDir, 0755)
	os.WriteFile(filepath.Join(srcDir, "file1.txt"), []byte("text"), 0644)
	os.WriteFile(filepath.Join(srcDir, "file2.log"), []byte("log data"), 0644)
	os.WriteFile(filepath.Join(srcDir, "file3.txt"), []byte("text again"), 0644)

	zipPath := filepath.Join(tempDir, "filter.zip")
	destDir := filepath.Join(tempDir, "out_filtered")

	script := `
	arsip.buat_zip("` + filepath.ToSlash(srcDir) + `", "` + filepath.ToSlash(zipPath) + `")

	// Saring: hanya ekstrak berkas yang berakhiran .txt
	buat total = arsip.ekstrak_zip("` + filepath.ToSlash(zipPath) + `", "` + filepath.ToSlash(destDir) + `", {
		"saring": fungsi(info) {
			kembalikan teks.diakhiri_dengan(info.nama, ".txt")
		}
	})

	total
	`

	res, _ := runArchiveScript(t, script, filepath.Join(tempDir, "filter.ce"))
	if res != nil && res.Type() == evaluator.ERROR_OBJ {
		t.Fatalf("Script error: %s", res.Inspect())
	}
	intObj, ok := res.(*evaluator.Integer)
	if !ok || intObj.Value != 2 {
		t.Fatalf("Expected 2 extracted files, got %v", res)
	}

	if _, err := os.Stat(filepath.Join(destDir, "file1.txt")); err != nil {
		t.Errorf("Expected file1.txt to exist")
	}
	if _, err := os.Stat(filepath.Join(destDir, "file3.txt")); err != nil {
		t.Errorf("Expected file3.txt to exist")
	}
	if _, err := os.Stat(filepath.Join(destDir, "file2.log")); err == nil {
		t.Errorf("Expected file2.log to be skipped by filter")
	}
}

