package evaluator_test

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"celang/pkg/evaluator"
)

// createMaliciousZipFile membuat berkas zip dengan entri Zip Slip (path traversal "../...").
func createMaliciousZipFile(zipPath, maliciousEntryName string) error {
	f, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer f.Close()

	zw := zip.NewWriter(f)
	defer zw.Close()

	w, err := zw.Create(maliciousEntryName)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte("PAYLOAD BAHAYA"))
	return err
}

// createMaliciousTarGzFile membuat berkas tar.gz dengan entri Tar Slip.
func createMaliciousTarGzFile(tarGzPath, maliciousEntryName string) error {
	f, err := os.Create(tarGzPath)
	if err != nil {
		return err
	}
	defer f.Close()

	gw := gzip.NewWriter(f)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

	hdr := &tar.Header{
		Name: maliciousEntryName,
		Mode: 0644,
		Size: int64(len("PAYLOAD BAHAYA")),
	}
	if err := tw.WriteHeader(hdr); err != nil {
		return err
	}
	_, err = tw.Write([]byte("PAYLOAD BAHAYA"))
	return err
}

func TestArchiveAdversarial_ZipSlipDetection(t *testing.T) {
	tempDir := t.TempDir()
	maliciousZip := filepath.Join(tempDir, "evil.zip")
	destDir := filepath.Join(tempDir, "safe_target")
	outsideFile := filepath.Join(tempDir, "pwned.txt")

	// Entri mengarah ke luar folder target
	if err := createMaliciousZipFile(maliciousZip, "../pwned.txt"); err != nil {
		t.Fatal(err)
	}

	script := `
	buat kenaGalat = salah
	coba {
		arsip.ekstrak_zip("` + filepath.ToSlash(maliciousZip) + `", "` + filepath.ToSlash(destDir) + `")
	} tangkap (e) {
		kenaGalat = benar
	}
	kenaGalat
	`

	res, _ := runArchiveScript(t, script, filepath.Join(tempDir, "main.ce"))
	if res != evaluator.TRUE {
		t.Errorf("Expected Zip Slip to be caught and return true, got %v", res)
	}

	// Pastikan berkas di luar destDir TIDAK dibuat
	if _, err := os.Stat(outsideFile); err == nil {
		t.Fatalf("SECURITY VULNERABILITY: Zip Slip berhasil menulis berkas di luar target: %s", outsideFile)
	}
}

func TestArchiveAdversarial_TarSlipDetection(t *testing.T) {
	tempDir := t.TempDir()
	maliciousTarGz := filepath.Join(tempDir, "evil.tar.gz")
	destDir := filepath.Join(tempDir, "safe_target_tar")
	outsideFile := filepath.Join(tempDir, "pwned_tar.txt")

	if err := createMaliciousTarGzFile(maliciousTarGz, "../../pwned_tar.txt"); err != nil {
		t.Fatal(err)
	}

	script := `
	buat kenaGalat = salah
	coba {
		arsip.ekstrak_tar_gz("` + filepath.ToSlash(maliciousTarGz) + `", "` + filepath.ToSlash(destDir) + `")
	} tangkap (e) {
		kenaGalat = benar
	}
	kenaGalat
	`

	res, _ := runArchiveScript(t, script, filepath.Join(tempDir, "main.ce"))
	if res != evaluator.TRUE {
		t.Errorf("Expected Tar Slip to be caught and return true, got %v", res)
	}

	if _, err := os.Stat(outsideFile); err == nil {
		t.Fatalf("SECURITY VULNERABILITY: Tar Slip berhasil menulis berkas di luar target: %s", outsideFile)
	}
}

func TestArchiveAdversarial_CorruptedArchives(t *testing.T) {
	tempDir := t.TempDir()
	corruptFile := filepath.Join(tempDir, "corrupt.bin")
	os.WriteFile(corruptFile, []byte("bukan arsip valid sama sekali! 1234567890"), 0644)
	destDir := filepath.Join(tempDir, "corrupt_dest")

	scripts := []string{
		`arsip.ekstrak_zip("` + filepath.ToSlash(corruptFile) + `", "` + filepath.ToSlash(destDir) + `")`,
		`arsip.daftar_zip("` + filepath.ToSlash(corruptFile) + `")`,
		`arsip.ekstrak_tar_gz("` + filepath.ToSlash(corruptFile) + `", "` + filepath.ToSlash(destDir) + `")`,
		`arsip.daftar_tar_gz("` + filepath.ToSlash(corruptFile) + `")`,
		`arsip.ekstrak_tar("` + filepath.ToSlash(corruptFile) + `", "` + filepath.ToSlash(destDir) + `")`,
	}

	for i, sc := range scripts {
		res, _ := runArchiveScript(t, sc, filepath.Join(tempDir, "main.ce"))
		if res == nil || res.Type() != evaluator.ERROR_OBJ {
			t.Errorf("[%d] Expected corrupted archive to return Error object, got %v", i, res)
		}
	}
}

func TestArchiveAdversarial_CorruptedCompressionPayload(t *testing.T) {
	tempDir := t.TempDir()
	corruptBytes := bytes.Repeat([]byte{0xFF, 0x00, 0x7F}, 10)

	script := `
	buat b = bita([255, 0, 127, 255, 0, 127])
	buat e1 = salah
	coba { arsip.gunzip(b) } tangkap (e) { e1 = benar }

	buat e2 = salah
	coba { arsip.zlib_dekompres(b) } tangkap (e) { e2 = benar }

	[e1, e2]
	`
	_ = corruptBytes

	res, _ := runArchiveScript(t, script, filepath.Join(tempDir, "main.ce"))
	arr, ok := res.(*evaluator.Array)
	if !ok || len(arr.Elements) != 2 {
		t.Fatalf("Expected array result, got %v", res)
	}
	if arr.Elements[0] != evaluator.TRUE || arr.Elements[1] != evaluator.TRUE {
		t.Errorf("Expected errors to be caught gracefully, got %v", arr)
	}
}

func TestArchiveAdversarial_NonExistentSource(t *testing.T) {
	tempDir := t.TempDir()
	script := `
	buat e = salah
	coba {
		arsip.buat_zip("./folder_yang_pasti_tidak_ada", "out.zip")
	} tangkap (err) {
		e = benar
	}
	e
	`
	res, _ := runArchiveScript(t, script, filepath.Join(tempDir, "main.ce"))
	if res != evaluator.TRUE {
		t.Errorf("Expected error for non-existent source, got %v", res)
	}
}

func TestArchiveAdversarial_TypeValidation(t *testing.T) {
	tempDir := t.TempDir()
	tests := []struct {
		code    string
		errSub  string
	}{
		{`arsip.buat_zip(123, "out.zip")`, "sumber harus berupa jalur string atau array"},
		{`arsip.buat_zip(["valid.txt", 456], "out.zip")`, "daftar berkas harus berupa teks"},
		{`arsip.ekstrak_zip(123, "out")`, "argumen pertama harus berupa jalur string"},
		{`arsip.gzip(12345)`, "argumen harus berupa teks (String) atau Bita (Bytes)"},
		{`arsip.gunzip(true)`, "argumen harus berupa teks (String) atau Bita (Bytes)"},
	}

	for i, tt := range tests {
		res, _ := runArchiveScript(t, tt.code, filepath.Join(tempDir, "main.ce"))
		if res == nil || res.Type() != evaluator.ERROR_OBJ {
			t.Errorf("[%d] Expected error, got %v", i, res)
			continue
		}
		errObj := res.(*evaluator.Error)
		if !strings.Contains(errObj.Message, tt.errSub) {
			t.Errorf("[%d] Expected error message containing %q, got %q", i, tt.errSub, errObj.Message)
		}
	}
}
