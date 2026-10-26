package evaluator

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// sanitizeExtractPath memvalidasi bahwa jalur ekstraksi berkas berada di dalam direktori tujuan
// dan tidak mengalami path traversal (pencegahan kerentanan Zip Slip / Tar Slip).
func sanitizeExtractPath(targetDir, entryName string) (string, error) {
	cleanTarget, err := filepath.Abs(filepath.Clean(targetDir))
	if err != nil {
		return "", fmt.Errorf("direktori tujuan tidak valid: %w", err)
	}

	// Normalisasi pemisah direktori
	cleanEntry := filepath.Clean(filepath.FromSlash(entryName))
	if filepath.IsAbs(cleanEntry) || strings.HasPrefix(cleanEntry, ".."+string(filepath.Separator)) || cleanEntry == ".." {
		return "", fmt.Errorf("terdeteksi ancaman path traversal (Zip/Tar Slip) pada entri: '%s'", entryName)
	}

	fullPath := filepath.Join(cleanTarget, cleanEntry)
	cleanFull := filepath.Clean(fullPath)

	if !strings.HasPrefix(cleanFull, cleanTarget+string(filepath.Separator)) && cleanFull != cleanTarget {
		return "", fmt.Errorf("terdeteksi ancaman path traversal (Zip/Tar Slip) pada entri: '%s'", entryName)
	}

	return cleanFull, nil
}

// getArchiveDataBytes mengekstrak byte slice dari argumen String atau Bytes.
func getArchiveDataBytes(arg Object, funcName string) ([]byte, Object) {
	switch v := arg.(type) {
	case *String:
		return []byte(v.Value), nil
	case *Bytes:
		return v.Elements, nil
	default:
		return nil, newError(nil,
			fmt.Sprintf("%s: argumen harus berupa teks (String) atau Bita (Bytes), menerima %s", funcName, arg.Type()),
			"Contoh: "+funcName+"(\"halo dunia\") atau "+funcName+"(bita([65, 66]))")
	}
}

// resolveSingleDest menentukan lokasi akhir berkas tunggal yang diekstrak dengan aman.
func resolveSingleDest(destPath, entryName string) (string, error) {
	cleanDest, err := filepath.Abs(destPath)
	if err != nil {
		return "", fmt.Errorf("jalur tujuan tidak valid: %w", err)
	}

	fi, statErr := os.Stat(cleanDest)
	isDir := (statErr == nil && fi.IsDir()) || strings.HasSuffix(destPath, "/") || strings.HasSuffix(destPath, "\\")

	if isDir {
		base := filepath.Base(filepath.Clean(filepath.FromSlash(entryName)))
		if base == "." || base == "/" || base == ".." || base == "\\" {
			return "", fmt.Errorf("nama entri tidak valid untuk ekstraksi tunggal: '%s'", entryName)
		}
		finalDest := filepath.Join(cleanDest, base)
		cleanFinal := filepath.Clean(finalDest)
		if !strings.HasPrefix(cleanFinal, cleanDest+string(filepath.Separator)) && cleanFinal != cleanDest {
			return "", fmt.Errorf("terdeteksi ancaman path traversal: '%s'", entryName)
		}
		return cleanFinal, nil
	}

	return cleanDest, nil
}

// parseExtractOptions mengurai opsi timpa (overwrite) dan fungsi saring (filter) dari argumen Map opsional.
func parseExtractOptions(opts Object) (bool, Object) {
	overwrite := true
	var filterFn Object

	if m, ok := opts.(*Map); ok {
		if val, exists := m.Pairs["timpa"]; exists {
			overwrite = CoerceToBoolean(val)
		} else if val, exists := m.Pairs["overwrite"]; exists {
			overwrite = CoerceToBoolean(val)
		}

		if val, exists := m.Pairs["saring"]; exists {
			filterFn = val
		} else if val, exists := m.Pairs["filter"]; exists {
			filterFn = val
		} else if val, exists := m.Pairs["<builtin function: saring>"]; exists {
			filterFn = val
		} else if val, exists := m.Pairs["<builtin function: filter>"]; exists {
			filterFn = val
		}
	}
	return overwrite, filterFn
}


// ==============================================================================
// 1. ENGINE ZIP (archive/zip)
// ==============================================================================

// createZipArchive mengemas berkas atau folder menjadi file ZIP.
func createZipArchive(sourceArg Object, destZipPath string) error {
	var srcStringAbs string
	var srcStringIsDir bool
	var fileListAbs []string

	switch src := sourceArg.(type) {
	case *String:
		var err error
		srcStringAbs, err = filepath.Abs(src.Value)
		if err != nil {
			return fmt.Errorf("jalur sumber tidak valid: %w", err)
		}
		info, err := os.Stat(srcStringAbs)
		if err != nil {
			return fmt.Errorf("sumber '%s' tidak ditemukan: %w", src.Value, err)
		}
		srcStringIsDir = info.IsDir()

	case *Array:
		for _, el := range src.Elements {
			if _, ok := el.(*String); !ok {
				return fmt.Errorf("daftar berkas harus berupa teks (String), menerima %s", el.Type())
			}
		}
		for _, el := range src.Elements {
			strEl := el.(*String)
			filePath, err := filepath.Abs(strEl.Value)
			if err != nil {
				return fmt.Errorf("jalur berkas '%s' tidak valid: %w", strEl.Value, err)
			}
			info, err := os.Stat(filePath)
			if err != nil {
				return fmt.Errorf("berkas '%s' tidak ditemukan: %w", strEl.Value, err)
			}
			if info.IsDir() {
				return fmt.Errorf("entri '%s' berupa folder; untuk folder gunakan path folder langsung", strEl.Value)
			}
			fileListAbs = append(fileListAbs, filePath)
		}

	default:
		return fmt.Errorf("sumber harus berupa jalur string atau array daftar berkas, menerima %s", sourceArg.Type())
	}

	destAbs, err := filepath.Abs(destZipPath)
	if err != nil {
		return fmt.Errorf("jalur berkas ZIP tujuan tidak valid: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(destAbs), 0755); err != nil {
		return fmt.Errorf("gagal membuat direktori induk tujuan: %w", err)
	}

	zipFile, err := os.Create(destAbs)
	if err != nil {
		return fmt.Errorf("gagal membuat berkas ZIP '%s': %w", destAbs, err)
	}
	defer zipFile.Close()

	zw := zip.NewWriter(zipFile)
	defer zw.Close()

	if len(fileListAbs) > 0 {
		for _, filePath := range fileListAbs {
			if err := addFileToZip(zw, filePath, filepath.Base(filePath)); err != nil {
				os.Remove(destAbs)
				return err
			}
		}
		return nil
	}

	if !srcStringIsDir {
		return addFileToZip(zw, srcStringAbs, filepath.Base(srcStringAbs))
	}

	// Folder rekursif
	err = filepath.Walk(srcStringAbs, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		absItem, _ := filepath.Abs(path)
		if absItem == destAbs {
			return nil
		}
		relPath, err := filepath.Rel(srcStringAbs, path)
		if err != nil {
			return err
		}
		if relPath == "." {
			return nil
		}

		headerName := filepath.ToSlash(relPath)
		if fi.IsDir() {
			headerName += "/"
			header, err := zip.FileInfoHeader(fi)
			if err != nil {
				return err
			}
			header.Name = headerName
			_, err = zw.CreateHeader(header)
			return err
		}

		return addFileToZip(zw, path, headerName)
	})
	if err != nil {
		os.Remove(destAbs)
		return fmt.Errorf("gagal mengemas folder ke ZIP: %w", err)
	}

	return nil
}

func addFileToZip(zw *zip.Writer, absPath, headerName string) error {
	fi, err := os.Stat(absPath)
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(fi)
	if err != nil {
		return err
	}
	header.Name = filepath.ToSlash(headerName)
	header.Method = zip.Deflate

	writer, err := zw.CreateHeader(header)
	if err != nil {
		return err
	}

	file, err := os.Open(absPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(writer, file)
	return err
}

// isZipArchive memeriksa apakah berkas merupakan arsip ZIP yang valid.
func isZipArchive(zipPath string) bool {
	cleanZip, err := filepath.Abs(zipPath)
	if err != nil {
		return false
	}
	info, err := os.Stat(cleanZip)
	if err != nil || info.IsDir() {
		return false
	}
	zr, err := zip.OpenReader(cleanZip)
	if err != nil {
		return false
	}
	zr.Close()
	return true
}

// createZipEntryInfo membangun objek Map representasi metadata entri ZIP.
func createZipEntryInfo(f *zip.File) *Map {
	isDir := f.FileInfo().IsDir()
	isSymlink := f.FileInfo().Mode()&os.ModeSymlink != 0
	typeStrID := "berkas"
	typeStrEN := "file"
	if isDir {
		typeStrID = "direktori"
		typeStrEN = "directory"
	} else if isSymlink {
		typeStrID = "tautan"
		typeStrEN = "symlink"
	}

	pairs := map[string]Object{
		"nama":               &String{Value: f.Name},
		"name":               &String{Value: f.Name},
		"ukuran":             &Integer{Value: int64(f.UncompressedSize64)},
		"size":               &Integer{Value: int64(f.UncompressedSize64)},
		"ukuran_terkompresi": &Integer{Value: int64(f.CompressedSize64)},
		"compressed_size":    &Integer{Value: int64(f.CompressedSize64)},
		"apakah_direktori":   nativeBool(isDir),
		"is_dir":             nativeBool(isDir),
		"apakah_tautan":      nativeBool(isSymlink),
		"is_symlink":         nativeBool(isSymlink),
		"tipe":               &String{Value: typeStrID},
		"type":               &String{Value: typeStrEN},
		"mode":               &Integer{Value: int64(f.Mode())},
		"waktu_modifikasi":   &String{Value: f.Modified.Format(time.RFC3339)},
		"modified":           &String{Value: f.Modified.Format(time.RFC3339)},
	}
	order := []string{"nama", "name", "ukuran", "size", "ukuran_terkompresi", "compressed_size", "apakah_direktori", "is_dir", "apakah_tautan", "is_symlink", "tipe", "type", "mode", "waktu_modifikasi", "modified"}
	return &Map{Pairs: pairs, Order: order}
}

// extractZipArchive mengekstrak seluruh arsip ZIP ke targetDir dengan proteksi Zip Slip dan opsi penyaring.
func extractZipArchive(zipPath, targetDir string, overwrite bool, filterFn Object) (int, error) {
	cleanZip, err := filepath.Abs(zipPath)
	if err != nil {
		return 0, fmt.Errorf("jalur berkas ZIP tidak valid: %w", err)
	}

	cleanTarget, err := filepath.Abs(targetDir)
	if err != nil {
		return 0, fmt.Errorf("jalur target tidak valid: %w", err)
	}

	zr, err := zip.OpenReader(cleanZip)
	if err != nil {
		return 0, fmt.Errorf("gagal membuka berkas ZIP '%s': %w", filepath.Base(cleanZip), err)
	}
	defer zr.Close()

	count := 0
	for _, f := range zr.File {
		if filterFn != nil && filterFn != NULL {
			infoMap := createZipEntryInfo(f)
			res := applyFunction(filterFn, []Object{infoMap}, nil)
			if isError(res) {
				return count, fmt.Errorf("kesalahan pada fungsi saring: %s", res.(*Error).Message)
			}
			if !CoerceToBoolean(res) {
				continue
			}
		}

		destPath, err := sanitizeExtractPath(cleanTarget, f.Name)
		if err != nil {
			return count, err
		}

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(destPath, f.Mode()); err != nil {
				return count, fmt.Errorf("gagal membuat direktori '%s': %w", destPath, err)
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return count, fmt.Errorf("gagal membuat direktori induk '%s': %w", filepath.Dir(destPath), err)
		}

		if !overwrite {
			if _, err := os.Stat(destPath); err == nil {
				// Berkas sudah ada dan opsi overwrite dimatikan
				continue
			}
		}

		rc, err := f.Open()
		if err != nil {
			return count, fmt.Errorf("gagal membuka entri '%s' di dalam ZIP: %w", f.Name, err)
		}

		outFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			rc.Close()
			return count, fmt.Errorf("gagal menulis berkas target '%s': %w", destPath, err)
		}

		_, copyErr := io.Copy(outFile, rc)
		rc.Close()
		outFile.Close()

		if copyErr != nil {
			return count, fmt.Errorf("gagal mengekstrak isi entri '%s': %w", f.Name, copyErr)
		}
		count++
	}

	return count, nil
}

// extractSingleFileFromZip mengekstrak satu berkas spesifik dari ZIP ke lokasi tujuan.
func extractSingleFileFromZip(zipPath, entryName, destPath string, overwrite bool) error {
	cleanZip, err := filepath.Abs(zipPath)
	if err != nil {
		return fmt.Errorf("jalur berkas ZIP tidak valid: %w", err)
	}

	zr, err := zip.OpenReader(cleanZip)
	if err != nil {
		return fmt.Errorf("gagal membuka berkas ZIP '%s': %w", filepath.Base(cleanZip), err)
	}
	defer zr.Close()

	normTarget := filepath.ToSlash(filepath.Clean(entryName))
	var targetFile *zip.File
	for _, f := range zr.File {
		if filepath.ToSlash(filepath.Clean(f.Name)) == normTarget {
			targetFile = f
			break
		}
	}

	if targetFile == nil {
		return fmt.Errorf("berkas '%s' tidak ditemukan di dalam arsip ZIP", entryName)
	}

	if targetFile.FileInfo().IsDir() {
		return fmt.Errorf("entri '%s' berupa direktori, bukan berkas", entryName)
	}

	finalDest, err := resolveSingleDest(destPath, targetFile.Name)
	if err != nil {
		return err
	}

	if !overwrite {
		if _, err := os.Stat(finalDest); err == nil {
			return fmt.Errorf("berkas tujuan '%s' sudah ada dan opsi timpa bernilai salah", finalDest)
		}
	}

	if err := os.MkdirAll(filepath.Dir(finalDest), 0755); err != nil {
		return fmt.Errorf("gagal membuat direktori induk '%s': %w", filepath.Dir(finalDest), err)
	}

	rc, err := targetFile.Open()
	if err != nil {
		return fmt.Errorf("gagal membuka entri '%s': %w", targetFile.Name, err)
	}
	defer rc.Close()

	outFile, err := os.OpenFile(finalDest, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, targetFile.Mode())
	if err != nil {
		return fmt.Errorf("gagal membuat berkas tujuan '%s': %w", finalDest, err)
	}
	defer outFile.Close()

	if _, err := io.Copy(outFile, rc); err != nil {
		return fmt.Errorf("gagal menulis isi berkas ke '%s': %w", finalDest, err)
	}

	return nil
}

// listZipArchive membaca daftar entri dalam file ZIP tanpa mengekstraknya.
func listZipArchive(zipPath string) (Object, error) {
	cleanZip, err := filepath.Abs(zipPath)
	if err != nil {
		return nil, fmt.Errorf("jalur ZIP tidak valid: %w", err)
	}

	zr, err := zip.OpenReader(cleanZip)
	if err != nil {
		return nil, fmt.Errorf("gagal membuka ZIP '%s': %w", filepath.Base(cleanZip), err)
	}
	defer zr.Close()

	var entries []Object
	for _, f := range zr.File {
		entries = append(entries, createZipEntryInfo(f))
	}

	return &Array{Elements: entries}, nil
}

// readFileFromZip membaca satu berkas dari dalam ZIP dan mengembalikan *Bytes.
func readFileFromZip(zipPath, targetFile string) (Object, error) {
	cleanZip, err := filepath.Abs(zipPath)
	if err != nil {
		return nil, fmt.Errorf("jalur ZIP tidak valid: %w", err)
	}

	zr, err := zip.OpenReader(cleanZip)
	if err != nil {
		return nil, fmt.Errorf("gagal membuka ZIP '%s': %w", filepath.Base(cleanZip), err)
	}
	defer zr.Close()

	normTarget := filepath.ToSlash(filepath.Clean(targetFile))
	for _, f := range zr.File {
		if filepath.ToSlash(filepath.Clean(f.Name)) == normTarget {
			rc, err := f.Open()
			if err != nil {
				return nil, fmt.Errorf("gagal membuka entri '%s': %w", f.Name, err)
			}
			defer rc.Close()

			data, err := io.ReadAll(rc)
			if err != nil {
				return nil, fmt.Errorf("gagal membaca isi entri '%s': %w", f.Name, err)
			}
			return &Bytes{Elements: data}, nil
		}
	}

	return nil, fmt.Errorf("berkas '%s' tidak ditemukan di dalam arsip ZIP", targetFile)
}

// ==============================================================================
// 2. ENGINE TAR & TAR.GZ (archive/tar & compress/gzip)
// ==============================================================================

// createTarArchive mengemas berkas atau folder menjadi file .tar.gz (atau .tar jika isGzip false).
func createTarArchive(sourceArg Object, destPath string, isGzip bool) error {
	var srcStringAbs string
	var srcStringIsDir bool
	var fileListAbs []string

	switch src := sourceArg.(type) {
	case *String:
		var err error
		srcStringAbs, err = filepath.Abs(src.Value)
		if err != nil {
			return fmt.Errorf("jalur sumber tidak valid: %w", err)
		}
		info, err := os.Stat(srcStringAbs)
		if err != nil {
			return fmt.Errorf("sumber '%s' tidak ditemukan: %w", src.Value, err)
		}
		srcStringIsDir = info.IsDir()

	case *Array:
		for _, el := range src.Elements {
			if _, ok := el.(*String); !ok {
				return fmt.Errorf("daftar berkas harus berupa teks (String), menerima %s", el.Type())
			}
		}
		for _, el := range src.Elements {
			strEl := el.(*String)
			filePath, err := filepath.Abs(strEl.Value)
			if err != nil {
				return fmt.Errorf("jalur berkas '%s' tidak valid: %w", strEl.Value, err)
			}
			info, err := os.Stat(filePath)
			if err != nil {
				return fmt.Errorf("berkas '%s' tidak ditemukan: %w", strEl.Value, err)
			}
			if info.IsDir() {
				return fmt.Errorf("entri '%s' berupa folder; untuk folder gunakan path folder langsung", strEl.Value)
			}
			fileListAbs = append(fileListAbs, filePath)
		}

	default:
		return fmt.Errorf("sumber harus berupa jalur string atau array daftar berkas, menerima %s", sourceArg.Type())
	}

	destAbs, err := filepath.Abs(destPath)
	if err != nil {
		return fmt.Errorf("jalur tujuan tidak valid: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(destAbs), 0755); err != nil {
		return fmt.Errorf("gagal membuat direktori induk tujuan: %w", err)
	}

	tarFile, err := os.Create(destAbs)
	if err != nil {
		return fmt.Errorf("gagal membuat berkas arsip '%s': %w", destAbs, err)
	}
	defer tarFile.Close()

	var tw *tar.Writer
	var gw *gzip.Writer

	if isGzip {
		gw = gzip.NewWriter(tarFile)
		defer gw.Close()
		tw = tar.NewWriter(gw)
	} else {
		tw = tar.NewWriter(tarFile)
	}
	defer tw.Close()

	if len(fileListAbs) > 0 {
		for _, filePath := range fileListAbs {
			if err := addFileToTar(tw, filePath, filepath.Base(filePath)); err != nil {
				os.Remove(destAbs)
				return err
			}
		}
		return nil
	}

	if !srcStringIsDir {
		return addFileToTar(tw, srcStringAbs, filepath.Base(srcStringAbs))
	}

	err = filepath.Walk(srcStringAbs, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		absItem, _ := filepath.Abs(path)
		if absItem == destAbs {
			return nil
		}
		relPath, err := filepath.Rel(srcStringAbs, path)
		if err != nil {
			return err
		}
		if relPath == "." {
			return nil
		}

		headerName := filepath.ToSlash(relPath)
		if fi.IsDir() {
			headerName += "/"
			header, err := tar.FileInfoHeader(fi, "")
			if err != nil {
				return err
			}
			header.Name = headerName
			return tw.WriteHeader(header)
		}

		return addFileToTar(tw, path, headerName)
	})
	if err != nil {
		os.Remove(destAbs)
		return fmt.Errorf("gagal mengemas folder ke TAR: %w", err)
	}

	return nil
}

func addFileToTar(tw *tar.Writer, absPath, headerName string) error {
	fi, err := os.Stat(absPath)
	if err != nil {
		return err
	}

	header, err := tar.FileInfoHeader(fi, "")
	if err != nil {
		return err
	}
	header.Name = filepath.ToSlash(headerName)

	if err := tw.WriteHeader(header); err != nil {
		return err
	}

	file, err := os.Open(absPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(tw, file)
	return err
}

// isTarArchive memeriksa apakah berkas merupakan arsip TAR polos yang valid.
func isTarArchive(archivePath string) bool {
	clean, err := filepath.Abs(archivePath)
	if err != nil {
		return false
	}
	info, err := os.Stat(clean)
	if err != nil || info.IsDir() {
		return false
	}

	f, err := os.Open(clean)
	if err != nil {
		return false
	}
	defer f.Close()

	// Pastikan bukan berkas gzip (gzip diawali magic byte 0x1f, 0x8b)
	magic := make([]byte, 2)
	n, _ := f.Read(magic)
	if n >= 2 && magic[0] == 0x1f && magic[1] == 0x8b {
		return false
	}

	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return false
	}

	tr := tar.NewReader(f)
	hdr, err := tr.Next()
	if err == nil && hdr != nil {
		return true
	}
	if err == io.EOF && info.Size() >= 1024 && info.Size()%512 == 0 {
		return true
	}
	return false
}

// isTarGzArchive memeriksa apakah berkas merupakan arsip TAR.GZ yang valid.
func isTarGzArchive(archivePath string) bool {
	clean, err := filepath.Abs(archivePath)
	if err != nil {
		return false
	}
	info, err := os.Stat(clean)
	if err != nil || info.IsDir() {
		return false
	}

	f, err := os.Open(clean)
	if err != nil {
		return false
	}
	defer f.Close()

	magic := make([]byte, 2)
	n, _ := f.Read(magic)
	if n < 2 || magic[0] != 0x1f || magic[1] != 0x8b {
		return false
	}

	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return false
	}

	gr, err := gzip.NewReader(f)
	if err != nil {
		return false
	}
	defer gr.Close()

	tr := tar.NewReader(gr)
	hdr, err := tr.Next()
	if err == nil && hdr != nil {
		return true
	}
	if err == io.EOF {
		return true
	}
	return false
}

// createTarEntryInfo membangun objek Map representasi metadata entri TAR.
func createTarEntryInfo(header *tar.Header) *Map {
	isDir := header.Typeflag == tar.TypeDir || strings.HasSuffix(header.Name, "/")
	isSymlink := header.Typeflag == tar.TypeSymlink || header.Typeflag == tar.TypeLink
	typeStrID := "berkas"
	typeStrEN := "file"
	if isDir {
		typeStrID = "direktori"
		typeStrEN = "directory"
	} else if isSymlink {
		typeStrID = "tautan"
		typeStrEN = "symlink"
	}

	pairs := map[string]Object{
		"nama":             &String{Value: header.Name},
		"name":             &String{Value: header.Name},
		"ukuran":           &Integer{Value: header.Size},
		"size":             &Integer{Value: header.Size},
		"apakah_direktori": nativeBool(isDir),
		"is_dir":           nativeBool(isDir),
		"apakah_tautan":    nativeBool(isSymlink),
		"is_symlink":       nativeBool(isSymlink),
		"tipe":             &String{Value: typeStrID},
		"type":             &String{Value: typeStrEN},
		"mode":             &Integer{Value: header.Mode},
		"waktu_modifikasi": &String{Value: header.ModTime.Format(time.RFC3339)},
		"modified":         &String{Value: header.ModTime.Format(time.RFC3339)},
	}
	order := []string{"nama", "name", "ukuran", "size", "apakah_direktori", "is_dir", "apakah_tautan", "is_symlink", "tipe", "type", "mode", "waktu_modifikasi", "modified"}
	return &Map{Pairs: pairs, Order: order}
}

// extractTarArchive mengekstrak arsip TAR atau TAR.GZ ke targetDir dengan proteksi Tar Slip dan opsi penyaring.
func extractTarArchive(archivePath, targetDir string, isGzip bool, overwrite bool, filterFn Object) (int, error) {
	cleanArchive, err := filepath.Abs(archivePath)
	if err != nil {
		return 0, fmt.Errorf("jalur berkas arsip tidak valid: %w", err)
	}

	cleanTarget, err := filepath.Abs(targetDir)
	if err != nil {
		return 0, fmt.Errorf("jalur direktori target tidak valid: %w", err)
	}

	file, err := os.Open(cleanArchive)
	if err != nil {
		return 0, fmt.Errorf("gagal membuka berkas arsip '%s': %w", filepath.Base(cleanArchive), err)
	}
	defer file.Close()

	var tr *tar.Reader
	if isGzip {
		gr, err := gzip.NewReader(file)
		if err != nil {
			return 0, fmt.Errorf("gagal mendekonstruksi stream gzip '%s': %w", filepath.Base(cleanArchive), err)
		}
		defer gr.Close()
		tr = tar.NewReader(gr)
	} else {
		tr = tar.NewReader(file)
	}

	count := 0
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return count, fmt.Errorf("kesalahan saat membaca entri TAR: %w", err)
		}

		if filterFn != nil && filterFn != NULL {
			infoMap := createTarEntryInfo(header)
			res := applyFunction(filterFn, []Object{infoMap}, nil)
			if isError(res) {
				return count, fmt.Errorf("kesalahan pada fungsi saring: %s", res.(*Error).Message)
			}
			if !CoerceToBoolean(res) {
				continue
			}
		}

		destPath, err := sanitizeExtractPath(cleanTarget, header.Name)
		if err != nil {
			return count, err
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(destPath, os.FileMode(header.Mode)); err != nil {
				return count, fmt.Errorf("gagal membuat direktori '%s': %w", destPath, err)
			}
		case tar.TypeReg, tar.TypeRegA:
			if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
				return count, fmt.Errorf("gagal membuat direktori induk '%s': %w", filepath.Dir(destPath), err)
			}

			if !overwrite {
				if _, err := os.Stat(destPath); err == nil {
					continue
				}
			}

			outFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				return count, fmt.Errorf("gagal membuat berkas target '%s': %w", destPath, err)
			}

			if _, err := io.Copy(outFile, tr); err != nil {
				outFile.Close()
				return count, fmt.Errorf("gagal mengekstrak isi entri '%s': %w", header.Name, err)
			}
			outFile.Close()
			count++
		}
	}

	return count, nil
}

// extractSingleFileFromTar mengekstrak satu berkas spesifik dari TAR atau TAR.GZ ke lokasi tujuan.
func extractSingleFileFromTar(archivePath, entryName, destPath string, isGzip bool, overwrite bool) error {
	cleanArchive, err := filepath.Abs(archivePath)
	if err != nil {
		return fmt.Errorf("jalur arsip tidak valid: %w", err)
	}

	file, err := os.Open(cleanArchive)
	if err != nil {
		return fmt.Errorf("gagal membuka berkas '%s': %w", filepath.Base(cleanArchive), err)
	}
	defer file.Close()

	var tr *tar.Reader
	if isGzip {
		gr, err := gzip.NewReader(file)
		if err != nil {
			return fmt.Errorf("gagal membaca stream gzip '%s': %w", filepath.Base(cleanArchive), err)
		}
		defer gr.Close()
		tr = tar.NewReader(gr)
	} else {
		tr = tar.NewReader(file)
	}

	normTarget := filepath.ToSlash(filepath.Clean(entryName))
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("kesalahan saat membaca entri TAR: %w", err)
		}

		if filepath.ToSlash(filepath.Clean(header.Name)) == normTarget {
			if header.Typeflag == tar.TypeDir || strings.HasSuffix(header.Name, "/") {
				return fmt.Errorf("entri '%s' berupa direktori, bukan berkas", entryName)
			}

			finalDest, err := resolveSingleDest(destPath, header.Name)
			if err != nil {
				return err
			}

			if !overwrite {
				if _, err := os.Stat(finalDest); err == nil {
					return fmt.Errorf("berkas tujuan '%s' sudah ada dan opsi timpa bernilai salah", finalDest)
				}
			}

			if err := os.MkdirAll(filepath.Dir(finalDest), 0755); err != nil {
				return fmt.Errorf("gagal membuat direktori induk '%s': %w", filepath.Dir(finalDest), err)
			}

			outFile, err := os.OpenFile(finalDest, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				return fmt.Errorf("gagal membuat berkas tujuan '%s': %w", finalDest, err)
			}
			defer outFile.Close()

			if _, err := io.Copy(outFile, tr); err != nil {
				return fmt.Errorf("gagal menulis isi berkas ke '%s': %w", finalDest, err)
			}

			return nil
		}
	}

	return fmt.Errorf("berkas '%s' tidak ditemukan di dalam arsip", entryName)
}

// listTarArchive membaca daftar entri dalam TAR atau TAR.GZ.
func listTarArchive(archivePath string, isGzip bool) (Object, error) {
	cleanArchive, err := filepath.Abs(archivePath)
	if err != nil {
		return nil, fmt.Errorf("jalur arsip tidak valid: %w", err)
	}

	file, err := os.Open(cleanArchive)
	if err != nil {
		return nil, fmt.Errorf("gagal membuka berkas '%s': %w", filepath.Base(cleanArchive), err)
	}
	defer file.Close()

	var tr *tar.Reader
	if isGzip {
		gr, err := gzip.NewReader(file)
		if err != nil {
			return nil, fmt.Errorf("gagal membaca stream gzip '%s': %w", filepath.Base(cleanArchive), err)
		}
		defer gr.Close()
		tr = tar.NewReader(gr)
	} else {
		tr = tar.NewReader(file)
	}

	var entries []Object
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("kesalahan saat membaca entri TAR: %w", err)
		}

		entries = append(entries, createTarEntryInfo(header))
	}

	return &Array{Elements: entries}, nil
}

// readFileFromTar membaca satu berkas dari dalam TAR atau TAR.GZ dan mengembalikan *Bytes.
func readFileFromTar(archivePath, targetFile string, isGzip bool) (Object, error) {
	cleanArchive, err := filepath.Abs(archivePath)
	if err != nil {
		return nil, fmt.Errorf("jalur arsip tidak valid: %w", err)
	}

	file, err := os.Open(cleanArchive)
	if err != nil {
		return nil, fmt.Errorf("gagal membuka berkas '%s': %w", filepath.Base(cleanArchive), err)
	}
	defer file.Close()

	var tr *tar.Reader
	if isGzip {
		gr, err := gzip.NewReader(file)
		if err != nil {
			return nil, fmt.Errorf("gagal membaca stream gzip '%s': %w", filepath.Base(cleanArchive), err)
		}
		defer gr.Close()
		tr = tar.NewReader(gr)
	} else {
		tr = tar.NewReader(file)
	}

	normTarget := filepath.ToSlash(filepath.Clean(targetFile))
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("kesalahan saat membaca entri TAR: %w", err)
		}

		if filepath.ToSlash(filepath.Clean(header.Name)) == normTarget {
			data, err := io.ReadAll(tr)
			if err != nil {
				return nil, fmt.Errorf("gagal membaca isi entri '%s': %w", header.Name, err)
			}
			return &Bytes{Elements: data}, nil
		}
	}

	return nil, fmt.Errorf("berkas '%s' tidak ditemukan di dalam arsip", targetFile)
}

// ==============================================================================
// 3. KOMPRESI DALAM-MEMORI (Gzip & Zlib)
// ==============================================================================

func gzipCompressBytes(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	if _, err := zw.Write(data); err != nil {
		return nil, err
	}
	if err := zw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func gzipDecompressBytes(compressed []byte) ([]byte, error) {
	zr, err := gzip.NewReader(bytes.NewReader(compressed))
	if err != nil {
		return nil, fmt.Errorf("data bukan format gzip yang valid: %w", err)
	}
	defer zr.Close()

	data, err := io.ReadAll(zr)
	if err != nil {
		return nil, fmt.Errorf("gagal mendekompresi payload gzip: %w", err)
	}
	return data, nil
}

func zlibCompressBytes(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	zw := zlib.NewWriter(&buf)
	if _, err := zw.Write(data); err != nil {
		return nil, err
	}
	if err := zw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func zlibDecompressBytes(compressed []byte) ([]byte, error) {
	zr, err := zlib.NewReader(bytes.NewReader(compressed))
	if err != nil {
		return nil, fmt.Errorf("data bukan format zlib yang valid: %w", err)
	}
	defer zr.Close()

	data, err := io.ReadAll(zr)
	if err != nil {
		return nil, fmt.Errorf("gagal mendekompresi payload zlib: %w", err)
	}
	return data, nil
}

// ==============================================================================
// 4. REGISTRASI MODUL PUSTAKA STANDAR (initStdlibArsip)
// ==============================================================================

func initStdlibArsip(env *Environment) {
	// ==========================================
	// 1. FORMAT ZIP
	// ==========================================

	// apakah_zip(jalur) / is_zip
	apakahZipFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 1 {
				return newError(nil, "arsip.apakah_zip membutuhkan 1 argumen: (jalur_berkas)", "Contoh: arsip.apakah_zip(\"paket.zip\")")
			}
			pathStr, ok := args[0].(*String)
			if !ok {
				return newError(nil, "arsip.apakah_zip: argumen harus berupa teks (String)", "")
			}
			return nativeBool(isZipArchive(pathStr.Value))
		},
	}
	definePublicFunction(env, "apakah_zip", apakahZipFn)
	definePublicFunction(env, "is_zip", apakahZipFn)

	// buat_zip(sumber, tujuanZip) / create_zip
	buatZipFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 2 {
				return newError(nil, "arsip.buat_zip membutuhkan minimal 2 argumen: (sumber, tujuan_zip)", "Contoh: arsip.buat_zip(\"./projek\", \"projek.zip\")")
			}
			destStr, ok := args[1].(*String)
			if !ok {
				return newError(nil, "arsip.buat_zip: argumen kedua (tujuan) harus berupa teks (String)", "")
			}
			if err := createZipArchive(args[0], destStr.Value); err != nil {
				return newError(nil, fmt.Sprintf("arsip.buat_zip gagal: %v", err), "")
			}
			return TRUE
		},
	}
	definePublicFunction(env, "buat_zip", buatZipFn)
	definePublicFunction(env, "create_zip", buatZipFn)

	// ekstrak_zip(arsipZip, targetDir, [opsi]) / extract_zip
	ekstrakZipFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 2 {
				return newError(nil, "arsip.ekstrak_zip membutuhkan minimal 2 argumen: (berkas_zip, folder_tujuan)", "Contoh: arsip.ekstrak_zip(\"paket.zip\", \"./tujuan\")")
			}
			zipStr, ok := args[0].(*String)
			if !ok {
				return newError(nil, "arsip.ekstrak_zip: argumen pertama harus berupa jalur string berkas ZIP", "")
			}
			targetStr, ok := args[1].(*String)
			if !ok {
				return newError(nil, "arsip.ekstrak_zip: argumen kedua harus berupa jalur string direktori tujuan", "")
			}

			overwrite := true
			var filterFn Object
			if len(args) >= 3 {
				overwrite, filterFn = parseExtractOptions(args[2])
			}

			count, err := extractZipArchive(zipStr.Value, targetStr.Value, overwrite, filterFn)
			if err != nil {
				return newError(nil, fmt.Sprintf("arsip.ekstrak_zip gagal: %v", err), "")
			}
			return &Integer{Value: int64(count)}
		},
	}
	definePublicFunction(env, "ekstrak_zip", ekstrakZipFn)
	definePublicFunction(env, "extract_zip", ekstrakZipFn)

	// ekstrak_berkas_zip(arsipZip, namaBerkas, tujuan, [opsi]) / extract_file_zip
	ekstrakBerkasZipFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 3 {
				return newError(nil, "arsip.ekstrak_berkas_zip membutuhkan minimal 3 argumen: (berkas_zip, nama_berkas, tujuan)", "Contoh: arsip.ekstrak_berkas_zip(\"paket.zip\", \"berkas.txt\", \"./keluaran/\")")
			}
			zipStr, ok := args[0].(*String)
			if !ok {
				return newError(nil, "arsip.ekstrak_berkas_zip: argumen pertama harus berupa String", "")
			}
			fileStr, ok := args[1].(*String)
			if !ok {
				return newError(nil, "arsip.ekstrak_berkas_zip: argumen kedua harus berupa String", "")
			}
			destStr, ok := args[2].(*String)
			if !ok {
				return newError(nil, "arsip.ekstrak_berkas_zip: argumen ketiga harus berupa String", "")
			}
			overwrite := true
			if len(args) >= 4 {
				overwrite, _ = parseExtractOptions(args[3])
			}
			if err := extractSingleFileFromZip(zipStr.Value, fileStr.Value, destStr.Value, overwrite); err != nil {
				return newError(nil, fmt.Sprintf("arsip.ekstrak_berkas_zip gagal: %v", err), "")
			}
			return TRUE
		},
	}
	definePublicFunction(env, "ekstrak_berkas_zip", ekstrakBerkasZipFn)
	definePublicFunction(env, "extract_file_zip", ekstrakBerkasZipFn)

	// daftar_zip(arsipZip) / list_zip
	daftarZipFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 1 {
				return newError(nil, "arsip.daftar_zip membutuhkan 1 argumen: (berkas_zip)", "")
			}
			zipStr, ok := args[0].(*String)
			if !ok {
				return newError(nil, "arsip.daftar_zip: argumen harus berupa teks (String)", "")
			}
			list, err := listZipArchive(zipStr.Value)
			if err != nil {
				return newError(nil, fmt.Sprintf("arsip.daftar_zip gagal: %v", err), "")
			}
			return list
		},
	}
	definePublicFunction(env, "daftar_zip", daftarZipFn)
	definePublicFunction(env, "list_zip", daftarZipFn)

	// baca_berkas_zip(arsipZip, namaBerkas) / read_file_zip
	bacaBerkasZipFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 2 {
				return newError(nil, "arsip.baca_berkas_zip membutuhkan 2 argumen: (berkas_zip, nama_berkas_dalam_zip)", "")
			}
			zipStr, ok := args[0].(*String)
			if !ok {
				return newError(nil, "arsip.baca_berkas_zip: argumen pertama harus berupa String", "")
			}
			fileStr, ok := args[1].(*String)
			if !ok {
				return newError(nil, "arsip.baca_berkas_zip: argumen kedua harus berupa String", "")
			}
			data, err := readFileFromZip(zipStr.Value, fileStr.Value)
			if err != nil {
				return newError(nil, fmt.Sprintf("arsip.baca_berkas_zip gagal: %v", err), "")
			}
			return data
		},
	}
	definePublicFunction(env, "baca_berkas_zip", bacaBerkasZipFn)
	definePublicFunction(env, "read_file_zip", bacaBerkasZipFn)

	// ==========================================
	// 2. FORMAT TAR.GZ
	// ==========================================

	// apakah_tar_gz(jalur) / is_tar_gz
	apakahTarGzFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 1 {
				return newError(nil, "arsip.apakah_tar_gz membutuhkan 1 argumen: (jalur_berkas)", "Contoh: arsip.apakah_tar_gz(\"arsip.tar.gz\")")
			}
			pathStr, ok := args[0].(*String)
			if !ok {
				return newError(nil, "arsip.apakah_tar_gz: argumen harus berupa teks (String)", "")
			}
			return nativeBool(isTarGzArchive(pathStr.Value))
		},
	}
	definePublicFunction(env, "apakah_tar_gz", apakahTarGzFn)
	definePublicFunction(env, "is_tar_gz", apakahTarGzFn)

	// buat_tar_gz(sumber, tujuanTarGz) / create_tar_gz
	buatTarGzFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 2 {
				return newError(nil, "arsip.buat_tar_gz membutuhkan minimal 2 argumen: (sumber, tujuan_tar_gz)", "Contoh: arsip.buat_tar_gz(\"./projek\", \"projek.tar.gz\")")
			}
			destStr, ok := args[1].(*String)
			if !ok {
				return newError(nil, "arsip.buat_tar_gz: argumen kedua harus berupa String", "")
			}
			if err := createTarArchive(args[0], destStr.Value, true); err != nil {
				return newError(nil, fmt.Sprintf("arsip.buat_tar_gz gagal: %v", err), "")
			}
			return TRUE
		},
	}
	definePublicFunction(env, "buat_tar_gz", buatTarGzFn)
	definePublicFunction(env, "create_tar_gz", buatTarGzFn)

	// ekstrak_tar_gz(arsipTarGz, targetDir, [opsi]) / extract_tar_gz
	ekstrakTarGzFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 2 {
				return newError(nil, "arsip.ekstrak_tar_gz membutuhkan minimal 2 argumen: (arsip_tar_gz, folder_tujuan)", "")
			}
			archiveStr, ok := args[0].(*String)
			if !ok {
				return newError(nil, "arsip.ekstrak_tar_gz: argumen pertama harus berupa String", "")
			}
			targetStr, ok := args[1].(*String)
			if !ok {
				return newError(nil, "arsip.ekstrak_tar_gz: argumen kedua harus berupa String", "")
			}

			overwrite := true
			var filterFn Object
			if len(args) >= 3 {
				overwrite, filterFn = parseExtractOptions(args[2])
			}

			count, err := extractTarArchive(archiveStr.Value, targetStr.Value, true, overwrite, filterFn)
			if err != nil {
				return newError(nil, fmt.Sprintf("arsip.ekstrak_tar_gz gagal: %v", err), "")
			}
			return &Integer{Value: int64(count)}
		},
	}
	definePublicFunction(env, "ekstrak_tar_gz", ekstrakTarGzFn)
	definePublicFunction(env, "extract_tar_gz", ekstrakTarGzFn)

	// ekstrak_berkas_tar_gz(arsipTarGz, namaBerkas, tujuan, [opsi]) / extract_file_tar_gz
	ekstrakBerkasTarGzFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 3 {
				return newError(nil, "arsip.ekstrak_berkas_tar_gz membutuhkan minimal 3 argumen: (arsip_tar_gz, nama_berkas, tujuan)", "Contoh: arsip.ekstrak_berkas_tar_gz(\"arsip.tar.gz\", \"modul.ce\", \"./keluaran/\")")
			}
			tarStr, ok := args[0].(*String)
			if !ok {
				return newError(nil, "arsip.ekstrak_berkas_tar_gz: argumen pertama harus berupa String", "")
			}
			fileStr, ok := args[1].(*String)
			if !ok {
				return newError(nil, "arsip.ekstrak_berkas_tar_gz: argumen kedua harus berupa String", "")
			}
			destStr, ok := args[2].(*String)
			if !ok {
				return newError(nil, "arsip.ekstrak_berkas_tar_gz: argumen ketiga harus berupa String", "")
			}
			overwrite := true
			if len(args) >= 4 {
				overwrite, _ = parseExtractOptions(args[3])
			}
			if err := extractSingleFileFromTar(tarStr.Value, fileStr.Value, destStr.Value, true, overwrite); err != nil {
				return newError(nil, fmt.Sprintf("arsip.ekstrak_berkas_tar_gz gagal: %v", err), "")
			}
			return TRUE
		},
	}
	definePublicFunction(env, "ekstrak_berkas_tar_gz", ekstrakBerkasTarGzFn)
	definePublicFunction(env, "extract_file_tar_gz", ekstrakBerkasTarGzFn)

	// daftar_tar_gz(arsipTarGz) / list_tar_gz
	daftarTarGzFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 1 {
				return newError(nil, "arsip.daftar_tar_gz membutuhkan 1 argumen: (berkas_tar_gz)", "")
			}
			archiveStr, ok := args[0].(*String)
			if !ok {
				return newError(nil, "arsip.daftar_tar_gz: argumen harus berupa String", "")
			}
			list, err := listTarArchive(archiveStr.Value, true)
			if err != nil {
				return newError(nil, fmt.Sprintf("arsip.daftar_tar_gz gagal: %v", err), "")
			}
			return list
		},
	}
	definePublicFunction(env, "daftar_tar_gz", daftarTarGzFn)
	definePublicFunction(env, "list_tar_gz", daftarTarGzFn)

	// baca_berkas_tar_gz(arsipTarGz, namaBerkas) / read_file_tar_gz
	bacaBerkasTarGzFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 2 {
				return newError(nil, "arsip.baca_berkas_tar_gz membutuhkan 2 argumen: (berkas_tar_gz, nama_berkas_dalam_tar)", "")
			}
			archiveStr, ok := args[0].(*String)
			if !ok {
				return newError(nil, "arsip.baca_berkas_tar_gz: argumen pertama harus berupa String", "")
			}
			fileStr, ok := args[1].(*String)
			if !ok {
				return newError(nil, "arsip.baca_berkas_tar_gz: argumen kedua harus berupa String", "")
			}
			data, err := readFileFromTar(archiveStr.Value, fileStr.Value, true)
			if err != nil {
				return newError(nil, fmt.Sprintf("arsip.baca_berkas_tar_gz gagal: %v", err), "")
			}
			return data
		},
	}
	definePublicFunction(env, "baca_berkas_tar_gz", bacaBerkasTarGzFn)
	definePublicFunction(env, "read_file_tar_gz", bacaBerkasTarGzFn)

	// ==========================================
	// 3. FORMAT TAR POLOS (Uncompressed)
	// ==========================================

	// apakah_tar(jalur) / is_tar
	apakahTarFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 1 {
				return newError(nil, "arsip.apakah_tar membutuhkan 1 argumen: (jalur_berkas)", "Contoh: arsip.apakah_tar(\"arsip.tar\")")
			}
			pathStr, ok := args[0].(*String)
			if !ok {
				return newError(nil, "arsip.apakah_tar: argumen harus berupa teks (String)", "")
			}
			return nativeBool(isTarArchive(pathStr.Value))
		},
	}
	definePublicFunction(env, "apakah_tar", apakahTarFn)
	definePublicFunction(env, "is_tar", apakahTarFn)

	// buat_tar(sumber, tujuanTar) / create_tar
	buatTarFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 2 {
				return newError(nil, "arsip.buat_tar membutuhkan minimal 2 argumen: (sumber, tujuan_tar)", "")
			}
			destStr, ok := args[1].(*String)
			if !ok {
				return newError(nil, "arsip.buat_tar: argumen kedua harus berupa String", "")
			}
			if err := createTarArchive(args[0], destStr.Value, false); err != nil {
				return newError(nil, fmt.Sprintf("arsip.buat_tar gagal: %v", err), "")
			}
			return TRUE
		},
	}
	definePublicFunction(env, "buat_tar", buatTarFn)
	definePublicFunction(env, "create_tar", buatTarFn)

	// ekstrak_tar(arsipTar, targetDir, [opsi]) / extract_tar
	ekstrakTarFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 2 {
				return newError(nil, "arsip.ekstrak_tar membutuhkan minimal 2 argumen: (arsip_tar, folder_tujuan)", "")
			}
			archiveStr, ok := args[0].(*String)
			if !ok {
				return newError(nil, "arsip.ekstrak_tar: argumen pertama harus berupa String", "")
			}
			targetStr, ok := args[1].(*String)
			if !ok {
				return newError(nil, "arsip.ekstrak_tar: argumen kedua harus berupa String", "")
			}

			overwrite := true
			var filterFn Object
			if len(args) >= 3 {
				overwrite, filterFn = parseExtractOptions(args[2])
			}

			count, err := extractTarArchive(archiveStr.Value, targetStr.Value, false, overwrite, filterFn)
			if err != nil {
				return newError(nil, fmt.Sprintf("arsip.ekstrak_tar gagal: %v", err), "")
			}
			return &Integer{Value: int64(count)}
		},
	}
	definePublicFunction(env, "ekstrak_tar", ekstrakTarFn)
	definePublicFunction(env, "extract_tar", ekstrakTarFn)

	// ekstrak_berkas_tar(arsipTar, namaBerkas, tujuan, [opsi]) / extract_file_tar
	ekstrakBerkasTarFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 3 {
				return newError(nil, "arsip.ekstrak_berkas_tar membutuhkan minimal 3 argumen: (berkas_tar, nama_berkas, tujuan)", "Contoh: arsip.ekstrak_berkas_tar(\"arsip.tar\", \"berkas.txt\", \"./keluaran/\")")
			}
			tarStr, ok := args[0].(*String)
			if !ok {
				return newError(nil, "arsip.ekstrak_berkas_tar: argumen pertama harus berupa String", "")
			}
			fileStr, ok := args[1].(*String)
			if !ok {
				return newError(nil, "arsip.ekstrak_berkas_tar: argumen kedua harus berupa String", "")
			}
			destStr, ok := args[2].(*String)
			if !ok {
				return newError(nil, "arsip.ekstrak_berkas_tar: argumen ketiga harus berupa String", "")
			}
			overwrite := true
			if len(args) >= 4 {
				overwrite, _ = parseExtractOptions(args[3])
			}
			if err := extractSingleFileFromTar(tarStr.Value, fileStr.Value, destStr.Value, false, overwrite); err != nil {
				return newError(nil, fmt.Sprintf("arsip.ekstrak_berkas_tar gagal: %v", err), "")
			}
			return TRUE
		},
	}
	definePublicFunction(env, "ekstrak_berkas_tar", ekstrakBerkasTarFn)
	definePublicFunction(env, "extract_file_tar", ekstrakBerkasTarFn)

	// daftar_tar(arsipTar) / list_tar
	daftarTarFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 1 {
				return newError(nil, "arsip.daftar_tar membutuhkan 1 argumen: (berkas_tar)", "")
			}
			archiveStr, ok := args[0].(*String)
			if !ok {
				return newError(nil, "arsip.daftar_tar: argumen harus berupa String", "")
			}
			list, err := listTarArchive(archiveStr.Value, false)
			if err != nil {
				return newError(nil, fmt.Sprintf("arsip.daftar_tar gagal: %v", err), "")
			}
			return list
		},
	}
	definePublicFunction(env, "daftar_tar", daftarTarFn)
	definePublicFunction(env, "list_tar", daftarTarFn)

	// baca_berkas_tar(arsipTar, namaBerkas) / read_file_tar
	bacaBerkasTarFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 2 {
				return newError(nil, "arsip.baca_berkas_tar membutuhkan 2 argumen: (berkas_tar, nama_berkas_dalam_tar)", "")
			}
			archiveStr, ok := args[0].(*String)
			if !ok {
				return newError(nil, "arsip.baca_berkas_tar: argumen pertama harus berupa String", "")
			}
			fileStr, ok := args[1].(*String)
			if !ok {
				return newError(nil, "arsip.baca_berkas_tar: argumen kedua harus berupa String", "")
			}
			data, err := readFileFromTar(archiveStr.Value, fileStr.Value, false)
			if err != nil {
				return newError(nil, fmt.Sprintf("arsip.baca_berkas_tar gagal: %v", err), "")
			}
			return data
		},
	}
	definePublicFunction(env, "baca_berkas_tar", bacaBerkasTarFn)
	definePublicFunction(env, "read_file_tar", bacaBerkasTarFn)

	// ==========================================
	// 4. IN-MEMORY GZIP & ZLIB
	// ==========================================

	// gzip(data) / compress_gzip
	gzipFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "arsip.gzip membutuhkan 1 argumen: (data teks atau Bita)", "Contoh: arsip.gzip(\"halo dunia\")")
			}
			raw, errObj := getArchiveDataBytes(args[0], "arsip.gzip")
			if errObj != nil {
				return errObj
			}
			comp, err := gzipCompressBytes(raw)
			if err != nil {
				return newError(nil, fmt.Sprintf("arsip.gzip gagal: %v", err), "")
			}
			return &Bytes{Elements: comp}
		},
	}
	definePublicFunction(env, "gzip", gzipFn)
	definePublicFunction(env, "compress_gzip", gzipFn)

	// gunzip(bita) / decompress_gzip
	gunzipFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "arsip.gunzip membutuhkan 1 argumen: (bita terkompresi gzip)", "Contoh: arsip.gunzip(data_terkompresi)")
			}
			raw, errObj := getArchiveDataBytes(args[0], "arsip.gunzip")
			if errObj != nil {
				return errObj
			}
			decomp, err := gzipDecompressBytes(raw)
			if err != nil {
				return newError(nil, fmt.Sprintf("arsip.gunzip gagal: %v", err), "")
			}
			return &Bytes{Elements: decomp}
		},
	}
	definePublicFunction(env, "gunzip", gunzipFn)
	definePublicFunction(env, "decompress_gzip", gunzipFn)

	// zlib_kompres(data) / zlib_compress
	zlibKompresFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "arsip.zlib_kompres membutuhkan 1 argumen: (data teks atau Bita)", "")
			}
			raw, errObj := getArchiveDataBytes(args[0], "arsip.zlib_kompres")
			if errObj != nil {
				return errObj
			}
			comp, err := zlibCompressBytes(raw)
			if err != nil {
				return newError(nil, fmt.Sprintf("arsip.zlib_kompres gagal: %v", err), "")
			}
			return &Bytes{Elements: comp}
		},
	}
	definePublicFunction(env, "zlib_kompres", zlibKompresFn)
	definePublicFunction(env, "zlib_compress", zlibKompresFn)

	// zlib_dekompres(bita) / zlib_decompress
	zlibDekompresFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "arsip.zlib_dekompres membutuhkan 1 argumen: (bita terkompresi zlib)", "")
			}
			raw, errObj := getArchiveDataBytes(args[0], "arsip.zlib_dekompres")
			if errObj != nil {
				return errObj
			}
			decomp, err := zlibDecompressBytes(raw)
			if err != nil {
				return newError(nil, fmt.Sprintf("arsip.zlib_dekompres gagal: %v", err), "")
			}
			return &Bytes{Elements: decomp}
		},
	}
	definePublicFunction(env, "zlib_dekompres", zlibDekompresFn)
	definePublicFunction(env, "zlib_decompress", zlibDekompresFn)
}
