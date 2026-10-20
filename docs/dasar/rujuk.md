---
outline: deep
---

# Referensi & Rujukan (`rujuk` & Pointer)

Dalam pemrograman, sering kali kita ingin bekerja dengan variabel yang sama dari tempat yang berbeda tanpa menduplikasi atau menyalin (*copy*) datanya ke memori baru.

Celang menyediakan mekanisme referensi yang sangat elegan dan ramah manusia melalui kata kunci **`rujuk`** (alias internasional: **`ref`**), serta operator penunjuk eksplisit (**`&`** dan **`*`**).

---

## Analogi Sederhana: Dua Remote, Satu Televisi

Bayangkan Anda memiliki satu unit televisi di ruang tamu:
- **Salinan Nilai (*Pass-by-value*)**: Anda membeli televisi kedua yang persis sama. Jika Anda mengganti saluran di TV kedua, TV pertama di ruang tamu sama sekali tidak terpengaruh.
- **Rujukan (*Reference / rujuk*)**: Anda membeli remote kontrol cadangan untuk televisi yang sama. Siapa pun yang menekan tombol di remote pertama maupun remote kedua, televisi yang berganti saluran adalah televisi yang sama!

---

## 1. Deklarasi Rujukan Variabel (`rujuk alias = target`)

Anda dapat memberikan nama panggilan (*alias*) baru ke variabel yang sudah ada:

```celang
buat saldoTabungan = 500000

// Buat alias rujukan ke 'saldoTabungan'
rujuk dompet = saldoTabungan

// Keduanya sekarang menunjuk ke kotak memori yang sama:
cetak("Saldo Tabungan:", saldoTabungan) // 500000
cetak("Dompet:", dompet)               // 500000

// Ubah nilai lewat 'dompet':
dompet = 750000

// Perhatikan: saldoTabungan asli otomatis ikut berubah!
cetak("Saldo Tabungan setelah dompet belanja:", saldoTabungan) // 750000
```

> [!NOTE]
> Nilai di sebelah kanan tanda `=` pada `rujuk alias = target` **wajib berupa nama variabel**, bukan angka atau hasil ekspresi langsung (seperti `rujuk x = 10` atau `rujuk x = a + b`), karena Anda merujuk pada wadah variabelnya, bukan hasil kalkulasinya.

---

## 2. Parameter Fungsi dengan `rujuk` (*Pass-by-Reference*)

Secara baku di hampir semua bahasa pemrograman modern, angka dan teks yang dikirim ke fungsi disalin nilainya. Jika fungsi mengubah nilai parameter, variabel asli di luar fungsi tidak akan berubah:

```celang
// ❌ TANPA RUJUK (Pass-by-value):
fungsi cobaUbah(angka) {
    angka = angka + 10
}

buat nilai = 50
cobaUbah(nilai)
cetak(nilai) // Tetap 50! Variabel asli tidak berubah.
```

Dengan menambahkan kata kunci **`rujuk`** pada parameter fungsi, fungsi tersebut langsung terhubung ke variabel pemanggil:

```celang
// ✅ DENGAN RUJUK (Pass-by-reference):
fungsi beriBonus(rujuk saldo, bonus) {
    saldo = saldo + bonus
}

buat saldoRekening = 100000
beriBonus(saldoRekening, 25000)

cetak("Saldo Akhir:", saldoRekening) // 125000! Variabel luar berhasil dimutasi.
```

### Contoh Nyata: Fungsi Menukar Nilai (*Swap*)

```celang
fungsi tukar(rujuk a, rujuk b) {
    buat cadangan = a
    a = b
    b = cadangan
}

buat kartuKiri = "As Hati"
buat kartuKanan = "Raja Sekop"

tukar(kartuKiri, kartuKanan)

cetak("Kartu Kiri:", kartuKiri)   // Raja Sekop
cetak("Kartu Kanan:", kartuKanan) // As Hati
```

---

## 3. Perbedaan `rujuk` vs Pointer Eksplisit (`&` dan `*`)

Celang sengaja menyediakan dua tingkatan penunjuk agar Anda bebas memilih gaya yang paling sesuai dengan kebutuhan:

| Aspek | Kata Kunci `rujuk` (Referensi) | Operator Pointer (`&` dan `*`) |
| :--- | :--- | :--- |
| **Sifat** | Alami, otomatis, dan transparan | Eksplisit, tampak alamat memorinya |
| **Sintaks Baca/Tulis** | Langsung nama variabelnya (`x = 10`) | Harus menggunakan tanda bintang (`*p = 10`) |
| **Kebutuhan Simbol** | Tanpa simbol rumit | Memakai `&` saat mengirim, `*` saat membaca |
| **Dapat Kosong (*Null*)?** | Tidak, selalu wajib merujuk variabel valid | Bisa berupa `NULL` (penunjuk kosong) |
| **Kasus Penggunaan** | Logika fungsi sehari-hari, *swap*, alias | Struktur data berantai (pohon, linked list), interaksi FFI C |

### Perbandingan Kode Berdampingan

:::: code-group
```celang [Menggunakan rujuk (Rekomendasi)]
fungsi tambahLima(rujuk angka) {
    angka = angka + 5 // Bersih dan alami tanpa simbol
}

buat skor = 90
tambahLima(skor)
cetak(skor) // 95
```

```celang [Menggunakan Pointer (& dan *)]
fungsi tambahLima(ptr) {
    *ptr = *ptr + 5 // Memerlukan dereference eksplisit
}

buat skor = 90
tambahLima(&skor) // Memerlukan operator address-of &
cetak(skor) // 95
```
::::

---

## 4. Alias Internasional (`ref`)

Dalam file yang menggunakan mode internasional atau pengembang yang terbiasa dengan C# / C++, kata kunci **`ref`** dapat digunakan sebagai sinonim persis dari `rujuk`:

```celang
let score = 100
ref playerRank = score
playerRank = 150

print("Original score:", score) // 150
```

---

## 5. Ringkasan & Praktik Terbaik

- Gunakan **`rujuk`** untuk membuat kode yang bersih, mudah dibaca anak-anak atau pemula, dan bebas dari kebingungan tanda bintang `*` dan ampersand `&`.
- Gunakan **`&` dan `*`** bila Anda membutuhkan representasi objek penunjuk eksplisit atau ingin membuat algoritma struktur data tingkat lanjut.
- Gunakan **`ffi.Pointer`** bila Anda berkomunikasi ke pustaka bahasa C eksternal.
