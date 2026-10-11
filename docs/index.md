---
layout: home

hero:
  name: "Celang"
  text: "Bahasa Pemrograman Human-Centric & Syntax-Styling"
  tagline: "Tulis kode dalam gayamu, nikmati kompilasi Bytecode VM secepat kilat, dan jelajahi riwayat variabel dengan Time-Travel."
  image:
    src: /logo.png
    alt: Celang Logo
  actions:
    - theme: brand
      text: 🚀 Mulai Sekarang
      link: /panduan/mulai-cepat
    - theme: brand
      text: 🎮 Buka Playground Web
      link: /playground
    - theme: alt
      text: ⚡ Coba Time-Travel
      link: /fitur-khas/time-travel

features:
  - icon: 🎨
    title: Syntax Styling Multi-Dialek
    details: Satu logika U-AST dapat dibaca dan diformat bebas dalam gaya Indonesia, Sunda, JavaScript, atau Pythonic tanpa merusak kode.
  - icon: ⏳
    title: Variabel Time-Travel (O(1) Memori)
    details: Lacak nilai variabel masa lalu secara instan via .sebelumnya, .awal, dan .histori dengan bounded ring-buffer tanpa memory leak.
  - icon: ⚡
    title: Bytecode VM Terintegrasi
    details: Berjalan secara default di Virtual Machine stack-based berkinerja tinggi, 2x lebih cepat dan 47% lebih hemat memori dibanding tree-walk.
  - icon: 🧩
    title: Bentuk (Struct) & Destrukturisasi
    details: Pemodelan data modern dengan method berorientasi objek, destrukturisasi array & map, serta interpolasi teks yang intuitif.
  - icon: 🌐
    title: Pustaka Standar Dwibasa
    details: Modul bawaan matematika, teks, waktu, berkas, dan JSON dapat dipanggil fleksibel menggunakan istilah lokal maupun bahasa internasional.
  - icon: 💖
    title: Diagnostik Empati Manusiawi
    details: Kesalahan kompilasi dan runtime disajikan dengan saran solusi yang ramah dan menuntun, bukan sekadar stack-trace membingungkan.
---

## Sekilas Kode Celang

::: code-group

```celang [🇮🇩 Gaya Indonesia]
// Bentuk data dan Time-Travel
buat gaji = 5000000
gaji = 7500000
gaji = 9500000

cetak("Gaji saat ini   :", gaji)
cetak("Gaji sebelumnya :", gaji.sebelumnya)
cetak("Gaji mula-mula  :", gaji.awal)
cetak("Histori gaji    :", gaji.histori)

// Pipeline Pemrosesan
buat angka = [1, 2, 3, 4, 5, 6]
buat hasil = angka
    |> saring(fungsi(x) { kembalikan x % 2 == 0 })
    |> petakan(fungsi(x) { kembalikan x * 10 })

cetak("Hasil pipeline  :", hasil) // [20, 40, 60]
```

```celang [🌾 Gaya Sunda]
// Wangun data sareng Time-Travel
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
