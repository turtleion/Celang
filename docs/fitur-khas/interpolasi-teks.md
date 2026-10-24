# Menyisipkan Nilai ke dalam Kalimat (Interpolasi Teks)

Saat membuat pesan atau pengumuman, menyambung-nyambungkan kata dengan tanda tambah (`+`) secara berulang-ulang sering kali melelahkan dan mudah salah:
```celang
// Melelahkan dan mudah salah tanda spasi:
buat pesan = "Halo " + nama + ", sisa uangmu adalah Rp " + saldo + "."
```

Celang menyediakan cara yang jauh lebih alami: **menyisipkan nilai langsung ke dalam kalimat**, seperti *f-string* pada Python maupun *template literal* pada JavaScript.

---

## 🐍 Gaya Format Teks Awalan `f` (f-string)

Mirip seperti di Python, Anda dapat menambahkan awalan `f` di depan tanda kutip: `f"..."` atau `f'...'`. Di dalamnya, Anda dapat menggunakan `${...}` maupun `{...}` untuk menyisipkan variabel:

```celang
buat pengguna = "Budi Santoso"
buat saldo = 130000

buat pesan = f"Halo ${pengguna}, saldo aktifmu adalah Rp${saldo}."
cetak(pesan)
// Hasil di layar: Halo Budi Santoso, saldo aktifmu adalah Rp130000.
```

### Menghitung Langsung di dalam Kalimat
Anda juga bisa meletakkan perhitungan matematika langsung di dalam kurung kurawal:

```celang
buat pesanHitung = f"Jika ditambah 10000: Rp${saldo + 10000}!"
cetak(pesanHitung)
// Hasil di layar: Jika ditambah 10000: Rp140000!
```

---

## 🇮🇩 Menyisipkan Nilai dengan `{ ... }`

Di dalam teks bertanda kutip ganda standar, Anda juga dapat langsung menggunakan kurung kurawal `{nama_variabel}`:

```celang
buat pesan = "Halo {pengguna}, saldo aktifmu adalah Rp {saldo}."
cetak(pesan)
```

---

## ⚡ Gaya Template Petik Terbalik (Backtick `${...}`)

Jika Anda terbiasa dengan gaya JavaScript atau TypeScript, Celang juga mendukung penulisan menggunakan tanda petik terbalik (`` `...` ``) dengan simbol `${...}`:

```celang
buat tingkatan = "VIP Platinum"
buat pesanKhusus = `Selamat ${pengguna}, status keanggotaanmu adalah ${tingkatan}.`

cetak(pesanKhusus)
// Hasil di layar: Selamat Budi Santoso, status keanggotaanmu adalah VIP Platinum.
```

