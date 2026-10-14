# Tanda Hitung & Kata Penghubung (Operator)

Operator adalah simbol atau kata penghubung khusus yang digunakan untuk menghitung angka, membandingkan dua nilai, atau merangkai logika pengambilan keputusan di dalam program Celang.

Celang mendukung tanda hitung matematika standar maupun kata penghubung alami (seperti `dan`, `atau`, `bukan`, dan `kalau kosong`).

---

## 1. Tanda Hitung Matematika

Digunakan untuk melakukan perhitungan angka sehari-hari:

| Tanda | Kegunaan | Contoh | Hasil |
| :--- | :--- | :--- | :--- |
| `+` | Penjumlahan atau menyambung kata | `10 + 5` | `15` |
| `-` | Pengurangan | `10 - 4` | `6` |
| `*` | Perkalian | `6 * 7` | `42` |
| `/` | Pembagian | `20 / 4` | `5` |
| `%` | Sisa pembagian (Modulo) | `10 % 3` | `1` (karena 10 dibagi 3 bersisa 1) |

```celang
buat x = 15
buat y = 4

cetak(x + y) // 19
cetak(x - y) // 11
cetak(x * y) // 60
cetak(x / y) // 3.75
cetak(x % y) // 3
```

---

## 2. Tanda Pengisian Nilai Cepat

Jika Anda ingin langsung menambah atau mengurangi isi variabel yang sudah ada:

| Tanda | Contoh | Arti Sebenarnya | Keterangan |
| :--- | :--- | :--- | :--- |
| `=` | `x = 10` | Isi `x` dengan angka 10 | Pengisian nilai biasa |
| `+=` | `x += 5` | `x = x + 5` | Tambahkan 5 ke dalam isi `x` |
| `-=` | `x -= 3` | `x = x - 3` | Kurangkan 3 dari isi `x` |
| `*=` | `x *= 2` | `x = x * 2` | Kalikan isi `x` dengan 2 |
| `/=` | `x /= 4` | `x = x / 4` | Bagikan isi `x` dengan 4 |

```celang
buat tabungan = 100000
tabungan += 50000  // sekarang tabungan bernilai 150000
tabungan -= 20000  // sekarang tabungan bernilai 130000
```

---

## 3. Tanda Pembanding Nilai

Digunakan untuk membandingkan dua hal. Jawabannya selalu berupa `benar` atau `salah`:

| Tanda | Maksud | Contoh | Hasil jika `a = 10, b = 20` |
| :--- | :--- | :--- | :--- |
| `==` | Apakah nilainya sama persis? | `a == b` | `salah` |
| `!=` | Apakah nilainya tidak sama / beda? | `a != b` | `benar` |
| `>` | Apakah lebih besar? | `a > b` | `salah` |
| `<` | Apakah lebih kecil? | `a < b` | `benar` |
| `>=` | Apakah lebih besar atau sama? | `a >= 10` | `benar` |
| `<=` | Apakah lebih kecil atau sama? | `a <= 20` | `benar` |

```celang
buat batasUsia = 17
buat usiaPengguna = 19

jika usiaPengguna >= batasUsia {
    cetak("Boleh masuk: Usia sudah mencukupi.")
}
```

---

## 4. Kata Penghubung Logika (`dan`, `atau`, `bukan`)

Celang memungkinkan Anda merangkai beberapa syarat menggunakan kata penghubung bahasa Indonesia:

| Kata Kunci | Simbol Alternatif | Maksud |
| :--- | :--- | :--- |
| `dan` | `&&` | Bernilai `benar` hanya jika **kedua syarat** sama-sama terpenuhi. |
| `atau` | `\|\|` | Bernilai `benar` jika **salah satu syarat saja** sudah terpenuhi. |
| `bukan` / `tidak` | `!` | Membalikkan arti (dari `benar` menjadi `salah`, atau sebaliknya). |

```celang
buat punyaAkun = benar
buat saldoCukup = salah

// Kedua syarat harus benar:
jika punyaAkun dan saldoCukup {
    cetak("Pesanan diproses.")
} selain {
    cetak("Pesanan ditolak: saldo tidak cukup.")
}

// Membalikkan keadaan:
buat sedangMemuat = salah
jika bukan sedangMemuat {
    cetak("Layar siap dibaca.")
}
```

---

## 5. Jaring Pengaman Nilai Kosong (`mungkin` & `kalau kosong`)

Mencegah program terhenti ketika data yang dicari ternyata kosong:

- **`mungkin`** (atau `?.`): Memeriksa data bertingkat secara aman tanpa macet jika salah satu bagian kosong.
- **`kalau kosong`** (atau `??`): Menyiapkan nilai cadangan jika data yang dicari bernilai kosong.

```celang
buat kotaAsal = mungkin profil.alamat.kota kalau kosong "Jakarta"
cetak("Kota tinggal:", kotaAsal)
```

---

## 6. Penyaluran Berantai (`|>`)

Mengalirkan hasil perhitungan dari langkah sebelumnya ke langkah berikutnya secara rapi:

```celang
buat hasil = [1, 2, 3, 4]
    |> saring(x => x % 2 == 0)
    |> petakan(x => x * 10)

cetak(hasil) // [20, 40]
```
