# In-Code Styling Block

Selain menetapkan dialek untuk seluruh berkas via baris perintah atau header `@gaya`, Celang memiliki fitur unik: **In-Code Styling Block**.

Fitur ini memungkinkan pengembang berpindah dialek di tengah-tengah kode program menggunakan blok `gaya "<nama_dialek>" { ... }`.

---

## 1. Motivasi & Kegunaan

Bayangkan skenario berikut:
1. Anda sedang menulis aplikasi berlogika bisnis berbahasa Indonesia.
2. Anda menyalin pustaka kecil (*snippet*) algoritma dari komunitas JavaScript atau Python.
3. Alih-alih harus menerjemahkan keyword secara manual satu per satu, Anda cukup membungkusnya dalam blok `gaya "javascript"` atau `gaya "pythonic"`.
4. Variabel dan fungsi yang dideklarasikan di dalam blok tersebut **langsung dapat digunakan** oleh bagian kode lainnya!

---

## 2. Contoh Penggunaan

Berikut adalah contoh nyata dari berkas [examples/incode_styling.ce](file:///home/yosatsu/Coding/Celang/examples/incode_styling.ce):

```celang
// @gaya: indonesian

// 1. Kode Utama dalam Gaya Indonesia
buat nama_pengguna = "Budi"
cetak("Nama Pengguna:", nama_pengguna)

// 2. Blok dengan Gaya JavaScript (In-Code Styling Block)
gaya "javascript" {
    let status_login = true
    function cek_akses(user) {
        return user === "Budi" && status_login
    }
    let diizinkan = cek_akses(nama_pengguna)
}

// 3. Kembali ke Konteks Global (Variabel & Fungsi Tetap Tersedia)
jika (diizinkan) {
    cetak("Akses Diberikan: Selamat datang,", nama_pengguna)
} selain {
    cetak("Akses Ditolak!")
}

// 4. Blok dengan Gaya Pythonic
gaya "pythonic" {
    def hitung_proyeksi(pokok, bunga, tahun): {
        total = pokok
        for t in rentang(1, tahun): {
            total = total + (total * bunga / 100)
        }
        return total
    }
    saldo_akhir = hitung_proyeksi(1000000, 10, 3)
}

cetak("Saldo Proyeksi Akhir:", saldo_akhir)
```

---

## 3. Cara Kerja di Balik Layar

Ketika parser Celang menemukan kata kunci `gaya "<profil>" { ... }`:

```mermaid
flowchart TD
    A[Parser Membaca Kode] --> B{Ketemu 'gaya "profil"'?}
    B -- Ya --> C[Push Profil Baru ke Stack Lexer]
    C --> D[Tokenisasi Menggunakan Keyword Profil Baru]
    D --> E[Parse Statement Menjadi U-AST Standar]
    E --> F[Tutup Blok '}']
    F --> G[Pop Profil & Kembalikan ke Profil Awal]
    B -- Tidak --> H[Gunakan Profil Aktif]
```

1. **Lexer Context Switching**: Lexer beralih ke tabel kata kunci dialek yang diminta (`javascript`, `pythonic`, dll).
2. **Normalisasi ke U-AST**: Pernyataan `let`, `function`, atau `def` diubah menjadi simpul U-AST universal yang sama persis seperti `buat` dan `fungsi`.
3. **Penyatuan Lingkup (*Scope Sharing*)**: Variabel yang dibuat masuk ke dalam *Environment* aktif sehingga tidak terisolasi dan dapat saling berinteraksi secara mulus.
