# Menyisipkan Nilai ke dalam Kalimat (Interpolasi Teks)

Saat membuat pesan atau pengumuman, menyambung-nyambungkan kata dengan tanda tambah (`+`) secara berulang-ulang sering kali melelahkan dan mudah salah:
```celang
// Melelahkan dan mudah salah tanda spasi:
buat pesan = "Halo " + nama + ", sisa uangmu adalah Rp " + saldo + "."
```

Celang menyediakan cara yang jauh lebih alami: **menyisipkan nilai langsung ke dalam kalimat**, seperti mengisi titik-titik pada formulir template surat.

---

## 🇮🇩 Menyisipkan Nilai dengan `{ ... }`

Di dalam kalimat bertanda kutip ganda (`"..."`), cukup gunakan kurung kurawal `{nama_variabel}` untuk menyisipkan nilai:

```celang
buat pengguna = "Budi Santoso"
buat saldo = 130000

buat pesan = "Halo {pengguna}, saldo aktifmu adalah Rp {saldo}."
cetak(pesan)
// Hasil di layar: Halo Budi Santoso, saldo aktifmu adalah Rp 130000.
```

### Menghitung Langsung di dalam Kalimat
Anda juga bisa meletakkan perhitungan matematika langsung di dalam kurung kurawal:

```celang
buat pesanHitung = "Jika saldomu digandakan (x2), menjadi: Rp {saldo * 2}!"
cetak(pesanHitung)
// Hasil di layar: Jika saldomu digandakan (x2), menjadi: Rp 260000!
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
