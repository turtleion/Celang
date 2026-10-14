---
layout: home

hero:
  name: "Celang"
  text: "Bahasa Pemrograman Ramah & Fleksibel"
  tagline: "Tulis kode dalam bahasa yang Anda pahami, nikmati eksekusi cepat, dan telusuri riwayat variabel dengan mudah."
  image:
    src: /logo.png
    alt: Celang Logo
  actions:
    - theme: brand
      text: 🚀 Mulai Sekarang
      link: /panduan/mulai-cepat
    - theme: brand
      text: 🎮 Coba di Browser (Playground)
      link: /playground
    - theme: alt
      text: ⚡ Jelajahi Riwayat Variabel
      link: /fitur-khas/time-travel

features:
  - icon: 🎨
    title: Gaya Bahasa Fleksibel
    details: Satu program dapat dibaca dan diubah tampilannya dengan bebas ke gaya Bahasa Indonesia, Basa Sunda, JavaScript, atau Python tanpa merusak kode.
  - icon: ⏳
    title: Variabel yang Mengingat Masa Lalu
    details: Mengetahui nilai sebelumnya, nilai awal, dan riwayat perubahan variabel secara instan, tetap ringan dan aman untuk komputer Anda.
  - icon: ⚡
    title: Eksekusi Cepat Tanpa Repot
    details: Kode langsung dijalankan oleh mesin virtual bawaan berkecepatan tinggi tanpa perlu pengaturan tambahan yang rumit.
  - icon: 🧩
    title: Bentuk Data & Kemudahan Menulis
    details: Kelompokkan data dengan rapi, bongkar isi daftar secara instan, dan sisipkan nilai ke dalam kalimat dengan sangat alami.
  - icon: 🌐
    title: Pustaka Alat Dwibasa
    details: Alat bantu matematika, pengolah teks, penunjuk waktu, pengelola berkas, dan pertukaran data JSON siap pakai dalam bahasa lokal maupun internasional.
  - icon: 💖
    title: Petunjuk Kesalahan yang Ramah
    details: Saat terjadi kesalahan, aplikasi menyajikan pesan yang menuntun dengan saran solusi nyata, bukan kode acak yang membingungkan.
---

## Sekilas Kode Celang

::: code-group

```celang [🇮🇩 Gaya Indonesia]
// Bentuk data dan Riwayat Nilai
buat gaji = 5000000
gaji = 7500000
gaji = 9500000

cetak("Gaji saat ini   :", gaji)
cetak("Gaji sebelumnya :", gaji.sebelumnya)
cetak("Gaji mula-mula  :", gaji.awal)
cetak("Histori gaji    :", gaji.histori)

// Penyaluran Proses Bertahap
buat angka = [1, 2, 3, 4, 5, 6]
buat hasil = angka
    |> saring(fungsi(x) { kembalikan x % 2 == 0 })
    |> petakan(fungsi(x) { kembalikan x * 10 })

cetak("Hasil akhir     :", hasil) // [20, 40, 60]
```

```celang [🌾 Gaya Sunda]
// Wangun data sareng Riwayat Nilai
jieun gaji = 5000000
gaji = 7500000
gaji = 9500000

cetak("Gaji ayeuna     :", gaji)
cetak("Gaji sateuacanna:", gaji.sebelumnya)
cetak("Gaji mimitina   :", gaji.awal)

// Transformasi Koléksi
jieun angka = [1, 2, 3, 4, 5, 6]
jieun hasil = angka
    |> saring(fungsi(x) { balikkeun x % 2 == 0 })
    |> petakan(fungsi(x) { balikkeun x * 10 })

cetak("Hasilna         :", hasil)
```

```celang [⚡ Gaya JavaScript]
// Object-like & Destructuring
let gaji = 5000000
gaji = 7500000
gaji = 9500000

cetak("Current salary  :", gaji)
cetak("Previous salary :", gaji.sebelumnya)

let [pertama, kedua, ...sisa] = [10, 20, 30, 40]
cetak(`Pertama: ${pertama}, Kedua: ${kedua}`)
```

```celang [🐍 Gaya Pythonic]
# Pythonic styling with def & match
gaji = 5000000
gaji = 7500000
gaji = 9500000

print("Current salary  :", gaji)
print("Previous salary :", gaji.sebelumnya)

def hitung_bonus(nilai):
    if nilai > 8000000:
        return nilai * 0.15
    else:
        return nilai * 0.05

print("Bonus:", hitung_bonus(gaji))
```

:::
