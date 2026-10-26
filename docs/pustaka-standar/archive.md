---
outline: deep
---

# Archive & Compression Library (`archive` / `compression`)

The **`archive`** module (with aliases **`compression`**, **`zip`**, **`tar`**, and Indonesian **`arsip`**) provides high-performance archiving, extraction, inspection, and in-memory compression capabilities in Celang.

Implemented in pure Go (*no external C libraries, no shell tools like tar or unzip*), this module delivers 100% cross-platform consistency across Linux, macOS, and Windows. It is **safe-by-default** with built-in protection against **Zip Slip** and **Tar Slip** vulnerabilities.

---

## 🎯 Key Capabilities

1. **ZIP Archives (`create_zip`, `extract_zip`, `list_zip`, `read_file_zip`)**:
   - Recursively archive directories or specific file lists into `.zip`.
   - Safely extract archives with path traversal verification.
   - Inspect archive file metadata (uncompressed/compressed sizes, modification time, mode) without extracting to disk.
   - Read specific files directly from archives into `Bytes`.
2. **TAR.GZ Tarballs (`create_tar_gz`, `extract_tar_gz`, `list_tar_gz`, `read_file_tar_gz`)**:
   - Industry-standard package distribution format used by the **Celang Package Manager**.
   - Preserves POSIX file modes (`0755` executable vs `0644` read-only/regular).
3. **Plain TAR (`create_tar`, `extract_tar`, `list_tar`)**:
   - Uncompressed streaming archiving for high-throughput pipeline workflows.
4. **In-Memory Compression**:
   - **Gzip**: `gzip(data)` and `gunzip(bytes)`.
   - **Zlib / Deflate**: `zlib_compress(data)` and `zlib_decompress(bytes)`.
5. **Security Hardening**:
   - **Anti Zip/Tar Slip**: Rejects malicious relative paths (`../../evil.sh`) trying to break outside the extraction root.
   - Graceful error handling on truncated or corrupt archive streams.

---

## 🚀 Usage Guide

```celang
"use universal"
import { create_zip, extract_zip, gzip, gunzip } from "archive"

// Create a zip bundle from a directory
const ok = create_zip("./project", "project.zip")

// Extract bundle to destination folder
const extractedCount = extract_zip("project.zip", "./dist")
print("Extracted files:", extractedCount)

// In-memory gzip compression
const payload = "High performance microservice payload..."
const compressed = gzip(payload)
const restored = gunzip(compressed).to_string()
```

---

## 📚 API Reference

### 1. ZIP Archives

#### `archive.create_zip(source, destZipPath)`
Packages files or directories into a `.zip` archive.
- **Parameters**:
  - `source`: `String` (path to folder or file) or `Array` of file paths.
  - `destZipPath`: `String` (output file path).
- **Returns**: `Boolean` (`true` on success).

#### `archive.extract_zip(zipPath, targetDir, [options])`
Extracts all files from a ZIP archive into `targetDir`.
- **Parameters**:
  - `zipPath`: `String`.
  - `targetDir`: `String`.
  - `options` (optional): `Map` with keys:
    - `"overwrite"`: `Boolean` (default `true`).
    - `"filter"`: Custom filter callback `fn(info) -> Boolean`.
- **Returns**: `Integer` (number of extracted files).

```celang
"use universal"
import { extract_zip } from "archive"
import { ends_with } from "strings"

const count = extract_zip("project.zip", "./dist", {
    "filter": fn(info) {
        return ends_with(info.name, ".ce")
    }
})
```

#### `archive.extract_file_zip(zipPath, entryName, destPath, [options])`
Extracts a single member file directly from the ZIP archive to disk without uncompressing the entire archive.

```celang
archive.extract_file_zip("archive.zip", "docs/spec.md", "./out/")
```

#### `archive.is_zip(path)`
Checks whether a given file is a valid ZIP archive by inspecting magic bytes and central directory headers.
- **Returns**: `Boolean`.

```celang
if (archive.is_zip("download.pkg")) {
    archive.extract_zip("download.pkg", "./target")
}
```

#### `archive.list_zip(zipPath)`
Lists archive contents without disk extraction.
- **Returns**: `Array` of `Map` objects with properties: `name`, `size`, `compressed_size`, `is_dir`, `is_symlink`, `type` (`"file"`, `"directory"`, `"symlink"`), `mode`, `modified`.

#### `archive.read_file_zip(zipPath, fileName)`
Directly reads a single file from within the ZIP archive into `Bytes`.

---

### 2. TAR.GZ Tarballs

- `archive.create_tar_gz(source, destTarGzPath)`: Compresses folder or files into `.tar.gz`.
- `archive.extract_tar_gz(tarGzPath, targetDir, [options])`: Extracts `.tar.gz` archive with Tar Slip guard and optional filter callback.
- `archive.extract_file_tar_gz(tarGzPath, entryName, destPath, [options])`: Extracts a single file from `.tar.gz` directly to disk.
- `archive.list_tar_gz(tarGzPath)`: Lists all file headers in `.tar.gz`.
- `archive.read_file_tar_gz(tarGzPath, fileName)`: Reads single file from `.tar.gz` into `Bytes`.
- `archive.is_tar_gz(path)`: Checks if file is a valid `.tar.gz` tarball.

---

### 3. Plain TAR

- `archive.create_tar(source, destTarPath)`: Creates uncompressed `.tar` archive.
- `archive.extract_tar(tarPath, targetDir, [options])`: Extracts `.tar` archive.
- `archive.extract_file_tar(tarPath, entryName, destPath, [options])`: Extracts single file from `.tar` to disk.
- `archive.read_file_tar(tarPath, fileName)`: Reads single file from uncompressed `.tar` into `Bytes`.
- `archive.list_tar(tarPath)`: Lists contents of `.tar` archive.
- `archive.is_tar(path)`: Checks if file is a valid plain `.tar` archive.

---

### 4. In-Memory Compression (Gzip & Zlib)

#### `archive.gzip(data)` & `archive.gunzip(compressedBytes)`
- `gzip(data)`: Compresses string or Bytes into RFC 1952 gzip `Bytes`.
- `gunzip(bytes)`: Decompresses gzip stream back into original `Bytes`.

#### `archive.zlib_compress(data)` & `archive.zlib_decompress(bytes)`
- `zlib_compress(data)`: Compresses data using RFC 1950 zlib format.
- `zlib_decompress(bytes)`: Decompresses zlib stream.

---

## 🛡️ Security & Slip Mitigation

> [!IMPORTANT]
> **Zip Slip & Tar Slip Defense**:
> Malicious packages could craft paths attempting to overwrite critical system files outside the target directory. Celang automatically computes canonical target boundaries and halts extraction immediately upon detecting any traversal escape, reporting an empathetic diagnostic error.
