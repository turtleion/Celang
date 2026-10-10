# String Interpolation

Celang menyediakan interpolasi teks (*string interpolation*) yang ekspresif, memungkinkanmu menyisipkan variabel dan ekspresi aritmatika langsung ke dalam teks tanpa perlu menggabungkannya berulang kali dengan operator `+`.

---

## 🇮🇩 Gaya Indonesia: Tanda Kutip Ganda `{ekspresi}`

Dalam teks dengan tanda kutip ganda (`"..."`), gunakan kurung kurawal `{...}` untuk menyisipkan nilai:

```celang
buat pengguna = "Budi Santoso"
buat saldo = 130000

buat pesan = "Halo {pengguna}, saldo aktifmu adalah Rp {saldo}."
cetak(pesan)
// Output: Halo Budi Santoso, saldo aktifmu adalah Rp 130000.
```

### Ekspresi Langsung di Dalam Teks
Kamu bisa meletakkan operasi perhitungan langsung di dalam kurung kurawal:

```celang
buat pesanHitung = "Jika saldo digandakan (x2): Rp {saldo * 2}!"
cetak(pesanHitung)
// Output: Jika saldo digandakan (x2): Rp 260000!
```

---

## ⚡ Gaya JavaScript / TypeScript: Backtick `${ekspresi}`

Bagi pengembang yang terbiasa dengan sintaks template literal JavaScript, Celang mendukung tanda petik terbalik (backtick `` `...` ``) dengan sintaks `${...}`:

```celang
buat level = "VIP Platinum"
buat pesanJS = `Selamat ${pengguna}, status keanggotaanmu adalah ${level}.`

cetak(pesanJS)
// Output: Selamat Budi Santoso, status keanggotaanmu adalah VIP Platinum.
```
