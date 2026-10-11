---
layout: doc
title: Live Web Playground
description: Coba dan eksekusi kode Celang secara langsung di browser Anda menggunakan WebAssembly (WASM).
---

# 🎮 WebAssembly Interactive Playground

Coba dan jalankan kode Celang secara instan langsung di peramban web (*browser*) Anda tanpa perlu menginstal compiler atau dependensi apa pun! 

Mesin eksekusi di bawah ini dijalankan sepenuhnya di sisi klien (*client-side*) menggunakan biner **WebAssembly (WASM)** yang dikompilasi dari kernel Go Celang, mendukung mode **Bytecode VM** berkecepatan tinggi maupun mode **Tree-Walking AST**.

---

<ClientOnly>
  <Playground />
</ClientOnly>

---

## 💡 Fitur yang Didukung di Browser

- **Pilihan Contoh Kode Bawaan**: Uji coba Halo Dunia, Variabel Time-Travel, Bentuk & Metode, Pipeline Operator, hingga Destrukturisasi.
- **Pilihan Mesin Ganda**: Beralih antara **Bytecode VM** (linear stack-based, performa tinggi) dan **AST Evaluator**.
- **Indikator Benchmark Real-Time**: Melihat waktu eksekusi langsung dalam milidetik (ms).
- **Keamanan Penuh**: Seluruh komputasi berjalan di dalam sandbox WebAssembly lokal di browser Anda tanpa mengirim data ke server luar.
