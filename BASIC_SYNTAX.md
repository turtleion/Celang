# Spesifikasi Sintaks Dasar Celang (Draft Awal)

> **Catatan untuk Kamu:**  
> Dokumen ini adalah rancangan awal seluruh aturan dan bentuk sintaks Celang.  
> **Silakan edit, ubah keyword, hapus, atau tambahkan aturan sesuai seleramu.**  
> Setelah selesai mengedit file ini, kabari aku dan aku akan membaca file ini untuk membangun Lexer, Parser, dan Style Profile berdasarkan versi editanmu!

---

## 1. Konsep Gaya Sintaks (Syntax Styling)

Karena Celang menyimpan kode dalam bentuk Universal AST (U-AST), kode yang sama dapat ditampilkan dan ditulis dalam beberapa gaya visual (_Style Profiles_).

Sebagai referensi, dokumen ini mendefinisikan:

- **Gaya Default (Bahasa Indonesia / Human-Centric)**
- **Gaya Alternatif (Pythonic & C-Style/JS)** sebagai pembanding.

Kamu bebas mengubah keyword mana pun di bawah ini.

---

## 2. Komentar (Comments)

### Usulan:

- **Satu baris:** Menggunakan `//` atau `#`
    ```celang
    // Ini komentar satu baris
    # Ini juga komentar satu baris
    ```
- **Banyak baris:** Menggunakan `/* ... */`
    ```celang
    /*
       Ini komentar
       beberapa baris
    */
    ```

---

## 3. Variabel & Tipe Data

### 3.1 Deklarasi Variabel

- **Variabel yang bisa diubah (Mutable):**

    ```celang
    // Gaya Indonesia:
    buat nama = "Budi"
    buat umur itu Angka
    umur = 25

    // Atau gaya klasik:
    let skor = 100
    var level = 1
    ```

- **Variabel tetap / Konstanta (Immutable):**

    ```celang
    // Gaya Indonesia:
    konstanta MAX_LOGIN = 3

    // Atau gaya klasik:
    const API_KEY = "xyz123"
    ```

### 3.2 Penugasan Ulang (Assignment)

```celang
umur = 26
skor = skor + 10
skor += 5
```

### 3.3 Tipe Data Dasar

| Tipe                     | Contoh Literal                                 | Catatan                                                                  |
| ------------------------ | ---------------------------------------------- | ------------------------------------------------------------------------ |
| **Angka (Integer)**      | `42`, `-10`, `1_000_000`                       | Mendukung pemisah underscore                                             |
| **Pecahan (Float)**      | `3.14`, `-0.05`                                |                                                                          |
| **Teks (String)**        | `"Halo Dunia"`, `'Halo'`                       | Mendukung interpolasi `"Halo ${nama}"`                                   |
| **Kebenaran (Boolean)**  | `benar` / `true`, `salah` / `false`            | Dua-duanya valid                                                         |
| **Kosong (Null)**        | `kosong`, `nihil`, `null`                      | Menandakan ketiadaan nilai                                               |
| **Array (Array / List)** | `[1, 2, 3, "empat"]`                           | Index berbasis 0: `data[0]`                                              |
| **Map (Dictionary)**     | `{"poin1": "blbalba", "poin2" itu "vlalabla"}` | Bisa pake : dan bisa pake "itu"                                          |
| **Sembarang (Any)**      | Sembarang                                      | Tipe data ini bisa apa aja, bukan tipe data default tapi bisa digunakan. |

### 3.4 Tipe Data Spesial

| Tipe                                                  | Contoh Literal                     | Catatan                                                                            |
| ----------------------------------------------------- | ---------------------------------- | ---------------------------------------------------------------------------------- |
| **Fungsi**                                            | `callback(fungsiA)`                | digunakan untuk memasukkan sebuah fungsi ke dalam parameter fungsi lain (callback) |
| **Fungsi-TIPEDATADASAR** (Explicit Function Typedata) | `callbackInteger(fungsiReturnInt)` | `Bisa juga explisit kayak gitu. Fungsi-Angka Fungsi-Kebenaran, dll.`               |

## Sebenarnya aku juga mau tambahin fitur pointer gitu. Tapi nanti aja

## 4. Fitur Khusus: Time-Travel Variables

Setiap variabel otomatis mencatat histori perubahannya. Variabel memiliki properti bawaan untuk mengakses masa lalu:

```celang
buat harga = 50000 (otomatis assign tipe data, tapi tetap bukan tipe data Sembarang)
harga = 75000
harga = 90000

// Mengakses nilai saat ini:
cetak(harga)             // Output: 90000

// Mengakses nilai sebelumnya:
cetak(harga.sebelumnya)  // Output: 75000
cetak(harga.prev)        // Alias singkat

// Mengakses nilai paling awal:
cetak(harga.awal)        // Output: 50000
cetak(harga.initial)     // Alias

// Mengakses seluruh riwayat:
cetak(harga.histori)     // Output: [50000, 75000, 90000]
```

---

## 5. Fitur Khusus: Smart Fallback (Tipe Adaptif)

Operasi tipe campuran tidak melempar error kaku jika maksud manusia dapat dipahami:

```celang
// 1. Ekstraksi otomatis dari string ke angka pada operasi hitung:
buat total = "Rp 150000" - 25000     // Hasil: 125000
buat bagi = "100" / 2                // Hasil: 50

// 2. String numerik + Angka:
buat tambah = "100" + 20             // Hasil: 120

// 3. String teks + Angka (tetap jadi sambung teks):
buat label = "Total item: " + 5      // Hasil: "Total item: 5"
```

---

## 6. Operator

### 6.1 Aritmatika

`+` (tambah / gabung), `-` (kurang), `*` (kali), `/` (bagi), `%` (sisa bagi)

### 6.2 Perbandingan

`==` (sama dengan), `!=` (tidak sama dengan), `<` (kurang dari), `>` (lebih dari), `<=` (kurang dari sama dengan), `>=` (lebih dari sama dengan)

### 6.3 Logika

- Logika simbol: `&&`, `||`, `!`
- Logika kata: `dan`, `atau`, `bukan` / `tidak`

```celang
jika (umur >= 17 dan punya_ktp) {
    // boleh masuk
}
```

---

## 7. Percabangan (Conditionals)

### 7.1 Gaya Blok Kurung Kurawal `{ }`

```celang
jika (umur >= 18) {
    cetak("Dewasa")
} selain jika (umur >= 13) {
    cetak("Remaja")
} selain {
    cetak("Anak-anak")
}
```

### 7.2 Gaya Keyword Alami (`maka` ... `selesai`)

```celang
jika (umur >= 18) maka
    cetak("Dewasa")
selain jika (umur >= 14)
    cetak("Remaja")
selain
    cetak("Bukan dewasa")
selesai
```

_(dukung keduanya lewat syntax style)._

---

## 8. Perulangan (Loops)

### 8.1 Perulangan Kondisi (While)

```celang
buat hitungan = 1
selama hitungan <= 5 {
    cetak("Hitungan: " + hitungan)
    hitungan += 1
}
```

### 8.2 Perulangan Koleksi / Iterasi (For ... In)

```celang
buat daftar_angka = [10, 20, 30]
untuk item dalam daftar_angka {
    cetak(item)
}

// Rentang angka (Range):
untuk i dalam rentang(1, 5) {
    cetak(i)
}
```

---

## 9. Fungsi (Functions)

### 9.1 Deklarasi & Pemanggilan Fungsi

```celang
fungsi sapa(nama) {
    kembalikan "Halo, " + nama + "!"
}

buat pesan = sapa("Celang")
cetak(pesan)
```

### 9.2 Parameter Opsional / Nilai Default

```celang
fungsi hitung_diskon(harga, persen = 10) {
    kembalikan harga - (harga * persen / 100)
}
```

### 9.3 Determined Parameter + Result (Optional)

```celang
fungsi hitung_luas_persegi(panjang itu Angka, lebar: Angka, formatter itu Fungsi-Angka) mengembalikan Angka {kembalikan formatter(panjang*lebar)}

fungsi eksekusi(cb: Fungsi-Kebenaran) -> Kebenaran {}

Bisa -> dan juga "mengembalikan"
```

### 9.4 Multiple Determined Parameter + Result (Optional)

```celang
fungsi hitung_luas_segitiga(tinggi: Angka atau Pecahan atau Teks = 2, alas itu Angka|Pecahan|Teks = 4.4) mengembalikan Angka|Pecahan (bisa juga "mengembalikan Angka atau Pecahan)

Bisa dengan keyword atau dan |
```

---

## 10. Fitur Khusus: Natural Concurrency (Goroutines & Sync)

Menyederhanakan konkurensi Go ke bahasa alami:

### 10.1 Asinkron di Latar Belakang (`sambil`)

Menjalankan tugas tanpa memblokir baris berikutnya:

```celang
sambil {
    unduh_file_besar("https://contoh.com/data.zip")
    kirim_notifikasi("Unduhan selesai!")
}

cetak("Program utama lanjut berjalan...")
```

### 10.2 Paralel Terkoordinasi (`bersamaan`)

Menjalankan beberapa tugas bersamaan dan otomatis menunggu semuanya selesai (seperti `sync.WaitGroup` otomatis):

```celang
bersamaan {
    proses_gambar("foto1.png")
    proses_gambar("foto2.png")
    proses_gambar("foto3.png")
}

cetak("Semua proses gambar selesai!")
```

---

## 11. Fungsi Bawaan Standar (Built-ins)

- `cetak(...)` / `print(...)`: Mencetak teks/variabel ke terminal.
- `masukan(pesan)` / `input(pesan)`: Menerima input dari pengguna di terminal.
- `panjang(objek)` / `len(objek)`: Menghitung jumlah elemen array atau panjang string.
- `tipe(objek)` / `type(objek)`: Mengembalikan nama tipe data objek.
- `rentang(awal, akhir)` / `range(start, end)`: Menghasilkan array deret angka.

---

## 12. Contoh Program Lengkap (Perbandingan Style)

### Gaya A: Alami Indonesia (Default Human-Centric)

```celang
buat barang = "Laptop"
harga = "Rp 10000000"
harga = 8500000

fungsi cek_promo(nilai itu Angka|Pecahan) mengembalikan Teks {
    jika (nilai < 9000000) {
        kembalikan "Diskon Spesial!"
    } selain  {
        kembalikan "Harga Normal"
    }
}

cetak("Barang:", barang)
cetak("Harga saat ini:", harga)
cetak("Harga awal:", harga.awal)
cetak("Status:", cek_promo(harga))

sambil {
    cetak("Sinkronisasi data ke cloud...")
}
```

### Gaya B: Pythonic Style

```celang
barang = "Laptop"
harga = "Rp 10000000"
harga = 8500000

def cek_promo(nilai):
    if nilai < 9000000:
        return "Diskon Spesial!"
    else:
        return "Harga Normal"

print("Barang:", barang)
print("Harga saat ini:", harga)
print("Harga awal:", harga.awal)
print("Status:", cek_promo(harga))

async:
    print("Sinkronisasi data ke cloud...")
```

---

## 13. Area yang Bisa Kamu Modifikasi Sekarang

Silakan buka dan edit file ini:

1. **Keyword:** Apakah lebih suka `buat` atau `var` / `let`, atau langsung tanpa keyword seperti Python? untuk default, ubah. Tapi untuk pythonic dan classic C, buatlah seperti aslinya.
2. **Kondisi:** Apakah suka kurung kurawal `{ }` atau indentasi atau `maka ... selesai`? Default { }, bisa disetting di style nanto
3. **Konkurensi:** Apakah kata `sambil` dan `bersamaan` sudah pas, atau mau diganti kata lain? iya udah
4. **Time-Travel:** Apakah penamaan properti `.sebelumnya`, `.awal`, `.histori` sudah nyaman digunakan? Nyaman
5. **Simbol Tambahan:** Apakah ada operator atau aturan sintaks lain yang ingin kamu tambahkan? Udah aku tambahin, tinggal lihat.
