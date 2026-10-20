---
outline: deep
---

# Antarmuka Fungsi Asing C (`ffi`)

Pernahkah Anda ingin membuat game grafis 3D dengan **Raylib**, mengakses database **SQLite**, atau memunculkan kotak pesan asli Windows (**Win32 MessageBox**), tetapi pustaka-pustaka tersebut ditulis dalam bahasa **C atau C++**?

Modul **`ffi`** (*Foreign Function Interface*) adalah jembatan emas di Celang. Dengan `ffi`, program Celang Anda dapat langsung memanggil fungsi apa pun yang ada di dalam berkas biner pustaka C (`.so` di Linux, `.dll` di Windows, dan `.dylib` di macOS) pada saat runtime!

---

## 💡 Mengapa FFI di Celang Sangat Spesial?

Bagi para programmer yang pernah menggunakan FFI di bahasa lain (seperti *CGO* di Go atau *ctypes* di Python), prosesnya sering kali membuat frustrasi:
- Anda diwajibkan memasang kompiler C berukuran gigabita (GCC, Clang, atau Visual Studio C++).
- Jika ingin membagikan aplikasi ke komputer teman (misalnya Anda memakai Linux dan teman memakai Windows), kompilasi silang (*cross-compile*) hampir selalu gagal karena ketergantungan library C lokal.

**Celang memecahkan masalah ini dengan pendekatan Dynamic ABI Invocation:**
1. **Nol Ketergantungan Kompiler C**: Celang berbicara langsung ke prosesor melalui konvensi biner mesin (ABI register x86-64 & ARM64). Anda **tidak perlu menginstal GCC, Clang, atau MinGW** sama sekali!
2. **Kompilasi Mandiri Mulus (`ce build`)**: Anda bisa menjalankan `celang build game.ce --target windows` di laptop Linux, dan biner `.exe` yang dihasilkan langsung dapat dijalankan di Windows tanpa komplain CGO!
3. **Cepat & Ringan**: Tidak ada proses kompilasi ulang kode C yang memakan waktu lama. Cukup letakkan berkas `.dll` atau `.so`, buka, dan panggil!

---

## 🗺️ Panduan 4 Langkah Menggunakan FFI

Gunakan alur 4 langkah sederhana ini setiap kali Anda ingin memanfaatkan pustaka C:

<div class="ffi-flow-container">
  <div class="ffi-flow-step">
    <div class="ffi-step-badge">1</div>
    <div class="ffi-step-icon">📂</div>
    <div class="ffi-step-title">Muat Pustaka</div>
    <div class="ffi-step-code">ffi.open("pustaka")</div>
    <div class="ffi-step-desc">Membuka berkas biner pustaka C <code>.so</code> / <code>.dll</code> / <code>.dylib</code></div>
  </div>
  <div class="ffi-flow-arrow">➔</div>
  <div class="ffi-flow-step">
    <div class="ffi-step-badge">2</div>
    <div class="ffi-step-icon">✍️</div>
    <div class="ffi-step-title">Daftarkan Fungsi</div>
    <div class="ffi-step-code">lib.func("nama", [in], out)</div>
    <div class="ffi-step-desc">Menentukan nama dan tipe data parameter serta kembalian</div>
  </div>
  <div class="ffi-flow-arrow">➔</div>
  <div class="ffi-flow-step">
    <div class="ffi-step-badge">3</div>
    <div class="ffi-step-icon">⚡</div>
    <div class="ffi-step-title">Panggil Fungsi</div>
    <div class="ffi-step-code">hasil = fungsiC(...)</div>
    <div class="ffi-step-desc">Eksekusi langsung lewat konvensi register ABI CPU native</div>
  </div>
  <div class="ffi-flow-arrow">➔</div>
  <div class="ffi-flow-step">
    <div class="ffi-step-badge">4</div>
    <div class="ffi-step-icon">🧹</div>
    <div class="ffi-step-title">Bebaskan Memori</div>
    <div class="ffi-step-code">ffi.free(ptr)</div>
    <div class="ffi-step-desc">Membersihkan alokasi memori buffer mentah jika diperlukan</div>
  </div>
</div>

<style>
.ffi-flow-container {
  display: flex;
  align-items: stretch;
  justify-content: space-between;
  gap: 12px;
  margin: 24px 0;
  padding: 18px;
  background: var(--vp-c-bg-soft);
  border: 1px solid var(--vp-c-divider);
  border-radius: 12px;
  overflow-x: auto;
}
.ffi-flow-step {
  flex: 1;
  min-width: 170px;
  background: var(--vp-c-bg);
  border: 1px solid var(--vp-c-divider);
  border-radius: 10px;
  padding: 16px 14px;
  text-align: center;
  position: relative;
  transition: all 0.2s ease;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}
.ffi-flow-step:hover {
  transform: translateY(-2px);
  border-color: var(--vp-c-brand-1);
  box-shadow: 0 4px 12px rgba(62, 175, 124, 0.15);
}
.ffi-step-badge {
  position: absolute;
  top: -10px;
  left: 12px;
  background: var(--vp-c-brand-1);
  color: white;
  width: 22px;
  height: 22px;
  border-radius: 50%;
  font-size: 12px;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
}
.ffi-step-icon {
  font-size: 26px;
  margin-bottom: 6px;
}
.ffi-step-title {
  font-weight: 700;
  font-size: 15px;
  color: var(--vp-c-text-1);
  margin-bottom: 6px;
}
.ffi-step-code {
  display: inline-block;
  font-size: 11px;
  padding: 3px 6px;
  border-radius: 6px;
  background: var(--vp-c-default-soft);
  color: var(--vp-c-brand-1);
  font-family: var(--vp-font-family-mono);
  margin-bottom: 8px;
  word-break: break-all;
}
.ffi-step-desc {
  font-size: 12px;
  color: var(--vp-c-text-2);
  line-height: 1.4;
}
.ffi-flow-arrow {
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  color: var(--vp-c-brand-1);
  font-weight: 700;
  opacity: 0.7;
  padding: 0 4px;
}
@media (max-width: 768px) {
  .ffi-flow-container {
    flex-direction: column;
    align-items: center;
  }
  .ffi-flow-step {
    width: 100%;
    min-width: 0;
  }
  .ffi-flow-arrow {
    transform: rotate(90deg);
    padding: 6px 0;
  }
}
</style>

---

## 1. Memuat Pustaka Dinamis (`open` / `buka`)

Langkah pertama adalah membuka berkas pustaka bersama (*shared library*):

```celang
impor ffi

// Di Linux:
let libc = ffi.open("libc.so.6")

// Di Windows:
let user32 = ffi.open("user32.dll")

// Di macOS:
let libm = ffi.open("libSystem.B.dylib")
```

> [!TIP]
> Jika berkas pustaka berada di satu folder dengan skrip Anda (misalnya `raylib.dll` atau `libsqlite3.so`), Anda cukup menuliskan `./raylib.dll` atau `./libsqlite3.so`.

---

## 2. Kamus Tipe Data FFI

Saat memberitahu Celang tentang bentuk fungsi C, Anda perlu menentukan jenis data parameter dan nilai kembaliannya. Celang menyediakan konstanta tipe yang jelas:

| Nama Tipe (Rekomendasi) | Alias (Indonesia) | Tipe Data di Bahasa C | Digunakan Untuk |
| :--- | :--- | :--- | :--- |
| `ffi.Int` / `ffi.Integer` | `ffi.Angka` | `int`, `int32_t`, `int64_t`, `long` | Bilangan bulat, kode status return |
| `ffi.Float` / `ffi.Double` | `ffi.Pecahan` | `double`, `float` | Bilangan desimal, koordinat, trigonometri |
| `ffi.String` | `ffi.Teks` | `const char*` | Teks C berakhiran null (*null-terminated*) |
| `ffi.Pointer` | `ffi.Penunjuk` | `void*`, `uintptr_t`, handle jendela | Alamat memori mentah, objek C, buffer |
| `ffi.Bool` | `ffi.Kebenaran` | `bool`, `_Bool` | Nilai kebenaran (`true` / `false`) |
| `ffi.Void` | `ffi.Kosong` | `void` | Fungsi yang tidak mengembalikan nilai |

> [!NOTE]
> Konstanta penunjuk utama dinamai **`ffi.Pointer`** agar seragam dengan standar industri internasional, dan `ffi.Penunjuk` tetap disediakan sebagai alias bahasa Indonesia.

---

## 3. Mendaftarkan Tanda Tangan Fungsi (`lib.func` / `lib.fungsi`)

Fungsi C didaftarkan menggunakan format:
```celang
let namaFungsi = lib.func("nama_simbol_c", [daftar_tipe_parameter], tipe_kembalian)
```

### Resep 1: Memanggil Fungsi Teks C (`puts` & `strlen`)

Fungsi `puts` di C menerima teks dan menampilkannya ke layar, sedangkan `strlen` menghitung panjang teks:

```celang
impor ffi

let libc = ffi.open("libc.so.6") // Di macOS gunakan: "libSystem.B.dylib"

// int puts(const char *s)
let c_puts = libc.fungsi("puts", [ffi.String], ffi.Int)

// size_t strlen(const char *s)
let c_strlen = libc.fungsi("strlen", [ffi.Pointer], ffi.Int)

// Panggil seperti fungsi Celang biasa!
c_puts("Halo! Kalimat ini dicetak langsung oleh libc bahasa C.")

let panjang = c_strlen("Celang FFI")
cetak("Panjang kalimat:", panjang) // Output: 10
```

---

### Resep 2: Perhitungan Matematika C Presisi Tinggi (`sqrt` & `sin`)

Fungsi matematika di pustaka C memanfaatkan register floating-point perangkat keras (*XMM registers*) berkecepatan tinggi:

```celang
impor ffi

let libm = ffi.open("libm.so.6")

// double sqrt(double x)
let c_sqrt = libm.fungsi("sqrt", [ffi.Float], ffi.Float)

// double sin(double x)
let c_sin = libm.fungsi("sin", [ffi.Float], ffi.Float)

cetak("Akar kuadrat 81.0:", c_sqrt(81.0)) // 9
cetak("Sinus sudut 90 derajat (1.5707963 rad):", c_sin(1.57079632679)) // ~1.0
```

---

### Resep 3: Memunculkan Dialog GUI Asli Windows (`MessageBoxA`)

Ingin membuat aplikasi desktop interaktif untuk pengguna Windows? Anda dapat memanggil Win32 API langsung tanpa library GUI tambahan:

```celang
impor ffi

// Buka pustaka antarmuka grafis Windows
let user32 = ffi.open("user32.dll")

// int MessageBoxA(HWND hWnd, LPCSTR lpText, LPCSTR lpCaption, UINT uType)
let tampilkanPesan = user32.fungsi(
    "MessageBoxA", 
    [ffi.Pointer, ffi.String, ffi.String, ffi.Int], 
    ffi.Int
)

// Parameter:
// 1. Handle jendela induk (ffi.NULL = tidak ada)
// 2. Isi pesan
// 3. Judul jendela dialog
// 4. Tombol (0 = Tombol OK saja, 1 = OK/Cancel, 4 = Yes/No)
tampilkanPesan(ffi.NULL, "Selamat! Anda berhasil memanggil Windows API dari Celang.", "Pemberitahuan Celang", 0)
```

---

## 4. Manajemen Memori & Utilitas Pointer (`ffi.Pointer`)

Banyak fungsi C tingkat lanjut yang meminta Anda menyediakan ruang memori (*buffer*) terlebih dahulu. Modul `ffi` menyediakan kendali memori mentah yang aman:

| Perintah | Alias | Penjelasan Ramah |
| :--- | :--- | :--- |
| `ffi.alloc(ukuran)` | `ffi.alokasi(ukuran)` | Memesan kapling memori kosong sebesar sekian byte |
| `ffi.free(ptr)` | `ffi.bebas(ptr)` | Mengembalikan kapling memori setelah selesai digunakan |
| `ffi.pointer(angka)` | `ffi.to_pointer(angka)` | Mengubah angka alamat menjadi objek `ffi.Pointer` |
| `ffi.address(ptr)` | `ffi.alamat(ptr)` | Melihat angka alamat asli penunjuk memori |
| `ffi.is_null(ptr)` | `ffi.apakah_kosong(ptr)` | Mengecek apakah penunjuk kosong (`ffi.NULL` / alamat 0) |
| `ffi.read_string(ptr)` | `ffi.baca_teks(ptr)` | Membaca rangkaian teks dari alamat memori |
| `ffi.write_string(ptr, teks)` | `ffi.tulis_teks(ptr, teks)` | Menuliskan teks ke alamat memori |
| `ffi.read_bytes(ptr, panjang)`| `ffi.baca_bita(ptr, panjang)`| Mengambil data biner mentah (`Bita`) sepanjang $N$ byte |
| `ffi.write_bytes(ptr, bita)`  | `ffi.tulis_bita(ptr, bita)`  | Menyalin data biner mentah ke alamat memori |

### Contoh Mengelola Buffer Memori Mentah:

```celang
impor ffi

// 1. Pesan buffer memori sebesar 128 byte
let buffer = ffi.alloc(128)

jika (!ffi.is_null(buffer)) {
    cetak("Buffer berhasil dibuat di alamat:", ffi.address(buffer))
    
    // 2. Tulis data ke dalam memori
    ffi.write_string(buffer, "Data rahasia tersimpan di memori C")
    
    // 3. Baca kembali kapan saja
    let hasil = ffi.read_string(buffer)
    cetak("Hasil pembacaan:", hasil)
}

// 4. Selalu bersihkan memori untuk mencegah memory leak!
ffi.free(buffer)
```

---

## 🛡️ Tips Keamanan & Pemecahan Masalah (Troubleshooting)

1. **Pustaka Tidak Ditemukan (*Library not found*)**:
   - Pastikan nama berkas pustaka dieja dengan benar.
   - Di Linux, nama pustaka biasanya memiliki ekstensi `.so` (contoh: `libc.so.6`).
   - Di Windows, pastikan ekstensi berkas adalah `.dll` (contoh: `user32.dll` atau `./raylib.dll`).
2. **Nama Fungsi Bersifat Peka Huruf Besar-Kecil (*Case-Sensitive*)**:
   - Jika fungsi C bernama `MessageBoxA`, memanggil `messageboxa` akan menghasilkan galat simbol tidak ditemukan.
3. **Selalu `free()` Memori yang di-`alloc()`**:
   - Setiap pemanggilan `ffi.alloc()` mengalokasikan memori sistem. Biasakan selalu memanggil `ffi.free(ptr)` setelah selesai, atau gunakan blok `saat_selesai` (defer) agar pembersihan terjadi otomatis saat fungsi berakhir:
   ```celang
   let ptr = ffi.alloc(256)
   saat_selesai { ffi.free(ptr) } // Otomatis dibebaskan saat blok selesai!
   ```
4. **Lingkungan WebAssembly / Browser**:
   - Karena keterbatasan keamanan peramban web (*browser sandbox*), WebAssembly tidak memiliki akses langsung ke pustaka `.so` atau `.dll` komputer lokal. Fitur FFI didesain untuk berjalan optimal pada lingkungan komputer asli (Linux, Windows, macOS).
