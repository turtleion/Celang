# 🌟 Celang - Bahasa Pemrograman yang Ramah dan Mudah Dipahami

Selamat datang di **Celang**! Celang adalah bahasa pemrograman yang dirancang agar terasa alami, ramah, dan mudah dipahami oleh manusia. 

Di Celang, Anda dapat menulis instruksi untuk komputer menggunakan bahasa sehari-hari—seperti Bahasa Indonesia dan Basa Jawa—atau menggunakan gaya penulisan populer lainnya seperti JavaScript dan Python. Yang terpenting, Anda tidak perlu khawatir tentang perdebatan gaya penulisan, karena satu program dapat diubah tampilannya ke gaya apa pun secara otomatis tanpa mengubah cara kerjanya.

---

## Daftar Isi

1. [A. Gambaran Umum](#a-gambaran-umum)
2. [B. Panduan Penggunaan Langkah demi Langkah](#b-panduan-penggunaan-langkah-demi-langkah)
3. [C. Penjelasan Fitur-Fitur Utama](#c-penjelasan-fitur-fitur-utama)
4. [D. Cara Kerja Aplikasi di Balik Layar](#d-cara-kerja-aplikasi-di-balik-layar)
5. [E. Informasi dan Penyimpanan Data](#e-informasi-dan-penyimpanan-data)
6. [F. Pengaturan dan Hak Akses](#f-pengaturan-dan-hak-akses)
7. [G. Masalah Umum dan Cara Mengatasinya](#g-masalah-umum-dan-cara-mengatasinya)
8. [H. Informasi untuk Pengembang Proyek](#h-informasi-untuk-pengembang-proyek)

---

## A. Gambaran Umum

### Apa itu Celang?
Celang adalah alat yang menerjemahkan tulisan manusia menjadi instruksi yang dapat dijalankan oleh komputer. Berbeda dari kebanyakan bahasa pemrograman yang kaku dan mewajibkan istilah bahasa Inggris, Celang memisahkan antara **makna logika program** dan **gaya tampilan tulisan**.

### Masalah Apa yang Diselesaikan?
1. **Kendala Bahasa bagi Pemula**: Banyak orang merasa kesulitan saat pertama kali belajar logika komputer karena harus menghafal istilah teknis dalam bahasa asing. Celang menyediakan kata kunci alami yang langsung dipahami maknanya, seperti `buat`, `jika`, `selama`, `fungsi`, dan `kembalikan`.
2. **Perdebatan Gaya Penulisan**: Dalam satu tim, sering kali ada orang yang lebih suka gaya penulisan tertentu (misalnya kurung kurawal gaya JavaScript vs indentasi gaya Python). Di Celang, kode disimpan sebagai data murni, sehingga setiap orang dapat melihat kode dalam gaya favorit mereka tanpa merusak berkas rekan satu timnya.
3. **Pesan Kesalahan yang Membingungkan**: Kebanyakan program komputer mengeluarkan pesan kesalahan yang dingin dan penuh istilah teknis yang sulit dipahami. Celang menyajikan pesan ramah dengan kotak petunjuk yang memberi tahu baris mana yang salah dan saran perbaikan praktisnya.

### Siapa yang Menggunakannya?
- **Pelajar dan Pemula**: Siapa pun yang ingin memahami logika pembuatan aplikasi tanpa terbebani kerumitan sintaks bahasa asing.
- **Pengembang dan Praktisi**: Orang yang ingin kemudahan mencatat riwayat perubahan nilai (*time-travel*), mengelola tugas latar belakang dengan tenang, atau menikmati fleksibilitas tampilan kode.

### Manfaat Utama
- **Alami dan Manusiawi**: Kata-kata yang digunakan mudah dimengerti sejak hari pertama.
- **Bisa Berganti Gaya**: Ubah tampilan kode antar bahasa hanya dengan satu perintah.
- **Variabel yang Mengingat Masa Lalu**: Variabel secara otomatis mencatat nilainya di masa lalu, nilainya saat pertama kali dibuat, dan linimasa perubahannya.
- **Cepat dan Mandiri**: Program langsung dijalankan oleh mesin virtual bawaan berkecepatan tinggi tanpa memerlukan aplikasi pihak ketiga yang rumit.

---

## B. Panduan Penggunaan Langkah demi Langkah

### 1. Memasang Celang di Komputer Anda

Jika Anda memiliki kode sumber Celang dan telah memasang perangkat Go di komputer:

1. Buka aplikasi terminal (Command Prompt di Windows, atau Terminal di macOS/Linux).
2. Masuk ke folder proyek Celang:
   ```bash
   cd Celang
   ```
3. Bangun aplikasi dengan perintah:
   ```bash
   go build -o celang ./cmd/celang
   ```
4. Periksa apakah Celang sudah siap digunakan:
   ```bash
   ./celang versi
   ```
   Aplikasi akan menampilkan nomor versi dan informasi sistem Anda.

---

### 2. Bermain Langsung di Jendela Interaktif (REPL)

Cara paling mudah untuk mencoba Celang adalah melalui jendela interaktif. Di sini, Anda bisa mengetik sebaris perintah dan komputer akan langsung menjawab seketika:

1. Buka jendela interaktif dengan mengetik:
   ```bash
   ./celang
   ```
2. Ketik perintah sederhana, lalu tekan **Enter**:
   ```celang
   celang> cetak("Halo, selamat datang di Celang!")
   Halo, selamat datang di Celang!
   ```
3. Coba buat simpanan nilai dan hitung:
   ```celang
   celang> buat saldo = 50000
   celang> saldo = saldo + 25000
   celang> cetak(saldo)
   75000
   ```
4. Untuk keluar dari jendela interaktif, ketik:
   ```celang
   celang> :keluar
   ```

---

### 3. Membuat dan Menjalankan Program Pertama

Anda juga bisa menulis kode di dalam berkas teks lalu menyuruh Celang menjalankannya:

1. Buat berkas baru bernama `salam.ce` menggunakan editor teks apa pun (seperti VS Code atau Notepad).
2. Tulis isi program berikut:
   ```celang
   // salam.ce
   buat nama = masukan("Siapa namamu? ")
   cetak("Senang berkenalan denganmu,", nama)

   buat hewan = ["Kucing", "Kelinci", "Burung"]
   untuk piaraan dalam hewan {
       cetak("Hewan kesukaan:", piaraan)
   }
   ```
3. Simpan berkas tersebut.
4. Buka terminal dan jalankan program Anda:
   ```bash
   ./celang run salam.ce
   ```
5. Komputer akan menanyakan nama Anda dan menampilkan salam serta daftar hewan secara berurutan.

---

### 4. Mencoba di Peramban Web Tanpa Instalasi

Jika Anda belum ingin memasang apa pun di komputer, Celang menyediakan halaman **Playground** berbasis peramban web:

1. Buka situs dokumentasi Celang dan klik menu **Playground**.
2. Ketik kode Anda pada kotak editor yang tersedia.
3. Klik tombol **Jalankan**.
4. Kode akan langsung diproses di dalam browser Anda tanpa mengirim data ke internet luar.

---

### 5. Mengubah Tampilan Kode Antar Gaya (Format)

Jika Anda memiliki berkas kode yang ditulis dalam Bahasa Indonesia dan ingin melihat bagaimana tampilannya dalam gaya lain:

```bash
# Menampilkan kode dalam gaya Python:
./celang format salam.ce --style=pythonic

# Menampilkan kode dalam gaya JavaScript:
./celang format salam.ce --style=javascript

# Langsung mengubah isi berkas secara permanen (gunakan opsi -w):
./celang format -w salam.ce --style=javascript
```

---

## C. Penjelasan Fitur-Fitur Utama

Setiap fitur di Celang dibuat untuk menyelesaikan kebutuhan nyata yang dihadapi orang saat membuat program.

### 1. Variabel Penjelajah Waktu (Time-Travel)
* **Kegunaannya**: Memungkinkan sebuah variabel mengingat nilai-nilai yang pernah dimilikinya di masa lalu.
* **Kapan Dibutuhkan**: Saat Anda ingin membandingkan nilai sekarang dengan nilai sebelumnya (misalnya kenaikan harga barang, skor permainan, atau status formulir) tanpa harus membuat banyak variabel tambahan.
* **Cara Menggunakannya**:
  ```celang
  buat harga = 10000
  harga = 15000
  harga = 20000

  cetak("Harga sekarang   :", harga)            // 20000
  cetak("Harga sebelumnya :", harga.sebelumnya) // 15000
  cetak("Harga mula-mula  :", harga.awal)       // 10000
  cetak("Histori kenaikan :", harga.histori)    // [10000, 15000, 20000]
  ```
* **Batasan**: Catatan riwayat (`.histori`) mengingat maksimal 50 perubahan terakhir agar memori komputer tetap ringan. Namun, nilai mula-mula (`.awal`) dan nilai sebelum perubahan terakhir (`.sebelumnya`) selalu diingat.

---

### 2. Bentuk Data Kustom (`bentuk`)
* **Kegunaannya**: Mengelompokkan beberapa informasi yang saling berhubungan ke dalam satu wadah rapi, serta melengkapinya dengan tindakan khusus.
* **Kapan Dibutuhkan**: Saat Anda memodelkan benda nyata, seperti rekening bank, data karyawan, atau keranjang belanja.
* **Cara Menggunakannya**:
  ```celang
  bentuk Rekening {
      pemilik,
      privat saldo = 0

      fungsi setor(jumlah) {
          @saldo += jumlah
          cetak("Setor berhasil. Saldo saat ini:", @saldo)
      }
  }

  buat r = Rekening{ pemilik: "Budi", saldo: 100000 }
  r.setor(50000)
  ```
* **Kondisi Khusus**: Kata kunci `privat` menandai bahwa suatu data hanya boleh diubah dari dalam bentuk itu sendiri. Tanda `@` adalah cara singkat untuk menyebut data milik dirinya sendiri.

---

### 3. Menitipkan Pekerjaan di Latar Belakang (`titip` & `jemput`)
* **Kegunaannya**: Menjalankan pekerjaan yang memakan waktu lama di latar belakang tanpa membuat program utama berhenti menunggu.
* **Analogi**: Seperti menaruh cucian di tempat laundry. Anda menitipkan cucian (`titip`), kasir memberi Anda tiket bukti titip, dan Anda bebas melakukan aktivitas lain. Ketika cucian sudah selesai, Anda menukar tiket tersebut untuk mengambil hasilnya (`jemput`).
* **Cara Menggunakannya**:
  ```celang
  fungsi hitungLaporan() {
      waktu.tidur(1000) // pura-pura proses butuh 1 detik
      kembalikan "Laporan Lengkap 2026"
  }

  // 1. Titipkan pekerjaan
  buat tiket = titip hitungLaporan()

  // 2. Lakukan hal lain sambil menunggu
  cetak("Melayani pelanggan lain...")

  // 3. Ambil hasilnya saat sudah siap
  buat hasil = jemput tiket
  cetak("Hasil diterima:", hasil)
  ```
* **Batas Waktu**: Anda dapat memberi batas waktu tunggu agar program tidak menunggu selamanya, misalnya: `jemput tiket dalam 3.detik`.

---

### 4. Saluran Pengiriman Pesan (`saluran`)
* **Kegunaannya**: Jalur komunikasi yang menghubungkan dua bagian program agar bisa saling bertukar pesan dengan aman dan teratur.
* **Analogi**: Seperti pipa pipa pengantar surat atau ban berjalan di pabrik. Satu petugas memasukkan barang ke satu ujung pipa, dan petugas lain menerimanya di ujung yang lain.
* **Cara Menggunakannya**:
  ```celang
  buat pipa = saluran(5) // pipa berkapasitas 5 pesan

  // Petugas 1 memasukkan pesan
  pipa.kirim("Pesanan #101")
  pipa.kirim("Pesanan #102")
  pipa.tutup()

  // Petugas 2 membaca semua pesan yang masuk
  untuk pesanan dalam pipa {
      cetak("Memproses:", pesanan)
  }
  ```
* **Kondisi Khusus**: Jangan mengirim data ke saluran yang sudah ditutup. Celang akan memberi tahu jika ada usaha pengiriman ke saluran tertutup.

---

### 5. Penyaluran Berantai (`|>`)
* **Kegunaannya**: Mengalirkan data dari satu proses ke proses berikutnya secara berurutan dari kiri ke kanan.
* **Analogi**: Seperti ban berjalan di restoran sushi, di mana setiap hidangan disiapkan tahap demi tahap.
* **Cara Menggunakannya**:
  ```celang
  buat angka = [1, 2, 3, 4, 5, 6]

  buat hasil = angka
      |> saring(fungsi(x) => x % 2 == 0)
      |> petakan(fungsi(x) => x * 10)

  cetak(hasil) // [20, 40, 60]
  ```

---

### 6. Jaring Pengaman Nilai Kosong (`mungkin` & `kalau kosong`)
* **Kegunaannya**: Mencegah program mendadak macet ketika data yang dicari ternyata tidak ada atau bernilai kosong.
* **Cara Menggunakannya**:
  ```celang
  buat profil = { nama: "Siti" }

  // 'mungkin' menjaga agar program tidak error jika properti alamat tidak ada
  buat kota = mungkin profil.alamat.kota

  // 'kalau kosong' memberikan nilai cadangan jika hasilnya kosong
  buat kotaPasti = kota kalau kosong "Kota Belum Diisi"

  cetak("Kota:", kotaPasti) // Kota Belum Diisi
  ```

---

### 7. Pembersihan Otomatis Saat Selesai (`saat_selesai`)
* **Kegunaannya**: Menjamin sebuah tindakan perapian atau penutupan selalu dikerjakan sebelum suatu fungsi selesai, apa pun yang terjadi (bahkan jika ada kesalahan).
* **Analogi**: Seperti menyalakan lampu peringatan saat masuk ruangan, dan memesan bahwa sakelar lampu harus dimatikan saat melangkah keluar pintu.
* **Cara Menggunakannya**:
  ```celang
  fungsi bacaCatatan() {
      cetak("1. Membuka koneksi...")
      saat_selesai {
          cetak("3. Menutup koneksi (selalu dijalankan!)")
      }
      cetak("2. Membaca data...")
  }

  bacaCatatan()
  ```

---

### 8. Pustaka Alat Bawaan
Celang dilengkapi lima kelompok alat bantu siap pakai yang dapat dipanggil dengan nama Bahasa Indonesia maupun padanan internasionalnya:

| Kelompok Alat | Contoh Perintah | Kegunaan |
| :--- | :--- | :--- |
| **Matematika (`matematika` / `math`)** | `matematika.akar(144)`, `matematika.pangkat(2, 8)`, `matematika.acak(1, 100)` | Menghitung akar, pangkat, nilai acak, nilai tertinggi, dan nilai mutlak. |
| **Pengolah Kata (`teks` / `strings`)** | `teks.kapital("halo")`, `teks.pisah("a,b,c", ",")`, `teks.gabung(["a", "b"], "-")` | Mengubah huruf kapital, memotong spasi, memecah kalimat, dan menyambung kata. |
| **Penunjuk Waktu (`waktu` / `time`)** | `waktu.sekarang()`, `waktu.format(waktu.sekarang())`, `waktu.tidur(500)` | Mengetahui jam dan tanggal saat ini, menunda jeda, dan mengukur durasi. |
| **Pengelola Berkas (`berkas` / `fs`)** | `berkas.baca("data.txt")`, `berkas.tulis("catatan.txt", "isi")`, `berkas.ada("data.txt")` | Membaca dokumen teks di komputer, menulis catatan baru, atau memeriksa keberadaan berkas. |
| **Pertukaran Data (`json`)** | `json.uraikan('{"id": 1}')`, `json.format(data, 2)` | Membaca data teks berformat JSON menjadi objek Celang atau menyusun data Celang menjadi teks JSON rapi. |

---

## D. Cara Kerja Aplikasi di Balik Layar

Ketika Anda menekan tombol Enter untuk menjalankan perintah `celang run berkas.ce`, berikut alur yang dilalui program Anda:

```
┌─────────────────┐       ┌─────────────────┐       ┌─────────────────┐       ┌─────────────────┐
│ 1. Membaca Teks │  ──►  │ 2. Memahami     │  ──►  │ 3. Menyusun     │  ──►  │ 4. Menjalankan  │
│    Kode Sumber  │       │    Kata-kata    │       │    Pohon Logika │       │    Instruksi    │
│    (File .ce)   │       │    (Pemecah Kata│       │    (Universal   │       │    (Mesin Cepat │
│                 │       │     / Lexer)    │       │     AST / U-AST)│       │     Bytecode VM)│
└─────────────────┘       └─────────────────┘       └─────────────────┘       └─────────────────┘
```

1. **Membaca Berkas**: Celang membuka berkas `.ce` yang Anda tunjuk dan membaca seluruh kalimat di dalamnya.
2. **Memahami Kata-kata**: Celang mengenali kata mana yang merupakan perintah (`buat`, `jika`), nama yang Anda buat sendiri (`saldo`, `nama`), atau angka dan teks biasa.
3. **Menyusun Pohon Logika (Universal AST)**: Celang menghubungkan kata-kata tersebut menjadi sebuah bagan pemahaman murni. Bagan ini bebas dari bahasa tertentu—ia hanya menyimpan inti ide dari program Anda.
4. **Menjalankan Instruksi (Mesin Virtual Cepat)**: Bagan pemahaman tersebut diubah menjadi daftar instruksi ringkas dan dijalankan secara berurutan oleh mesin virtual. Hasil keluaran langsung dicetak ke layar Anda.
5. **Menjaga Bila Terjadi Masalah (Diagnostik Empati)**: Jika Anda salah mengetik perintah atau melakukan kesalahan perhitungan (misalnya membagi angka dengan nol), Celang tidak menampilkan kode rahasia yang menyeramkan. Celang membungkus kesalahan tersebut ke dalam kotak petunjuk yang jelas, menunjukkan baris yang perlu diperbaiki beserta usulan solusinya.

---

## E. Informasi dan Penyimpanan Data

Agar Anda merasa tenang, berikut rincian tentang apa saja yang disimpan dan di mana lokasinya:

### Informasi yang Disimpan di Memori Kerja Sementara (RAM)
* **Variabel dan Nilai Program**: Seluruh data yang dibuat oleh program Anda (angka, teks, daftar, bentuk data) hanya hidup di dalam memori kerja (RAM) komputer selama program berjalan.
* **Riwayat Perubahan Nilai**: Setiap variabel menyimpan riwayat perubahan nilainya di RAM sebanyak maksimal 50 catatan.
* **Penghapusan Otomatis**: Begitu program selesai berjalan atau terminal ditutup, seluruh memori kerja tersebut langsung dibersihkan oleh sistem operasi.

### Informasi yang Disimpan di Penyimpanan Komputer (Disk)
* Celang tidak membuat basis data rahasia atau menyimpan log tersembunyi.
* Celang hanya membaca atau menulis berkas di komputer jika kode program Anda secara sengaja memanggil perintah `berkas.tulis` atau `berkas.baca`.
* Berkas proyek yang dihasilkan dari perintah `celang init` hanyalah berkas teks biasa:
  - `celang.mod`: Berisi catatan nama proyek dan nomor versinya.
  - `utama.ce`: Berkas tempat Anda mulai menulis kode.
  - `README.md`: Berkas petunjuk ringkas proyek.

### Keamanan dan Privasi
* Celang berjalan sepenuhnya secara luring (*offline*) di komputer Anda.
* Celang tidak mengirim kode, data pribadi, atau statistik apa pun ke internet atau server luar.
* Pada mode peramban web (Playground), seluruh perhitungan berjalan di dalam lingkungan tertutup (*sandbox*) browser Anda sendiri.

---

## F. Pengaturan dan Hak Akses

### 1. Pilihan Perintah Terminal

Anda dapat menjalankan perintah berikut di terminal:

| Perintah | Contoh | Kegunaan |
| :--- | :--- | :--- |
| `celang` | `./celang` | Membuka sesi terminal interaktif secara langsung. |
| `celang run` | `./celang run program.ce` | Menjalankan program menggunakan mesin virtual cepat. |
| `celang run ... --ast` | `./celang run program.ce --ast` | Menjalankan program menggunakan penafsir pohon AST langsung (opsi pembanding). |
| `celang format` | `./celang format program.ce --style=jawa` | Menampilkan atau menata ulang format kode sesuai gaya bahasa yang dipilih. |
| `celang compile` | `./celang compile program.ce -o program.ast` | Mengubah kode teks menjadi berkas data logika U-AST berformat JSON. |
| `celang init` | `./celang init toko_buku` | Membuat folder proyek baru lengkap dengan berkas pembukanya. |
| `celang versi` | `./celang versi` | Memeriksa nomor versi Celang yang terpasang. |
| `celang bantuan` | `./celang bantuan` | Menampilkan ringkasan bantuan semua perintah. |

---

### 2. Hak Akses Antar-Berkas Kode (Visibilitas)

Ketika proyek Anda bertambah besar dan dipecah ke dalam beberapa berkas, Celang melindungi data Anda dengan aturan berikut:

* **Secara bawaan bersifat rahasia (`privat!`)**: Semua variabel dan fungsi di dalam satu berkas hanya bisa dilihat oleh berkas itu sendiri. Berkas lain tidak bisa sembarangan mengintip atau mengubahnya.
* **Terbuka untuk umum (`publik!` atau `export`)**: Tambahkan kata `publik!` jika Anda sengaja ingin fungsi atau data tersebut bisa digunakan oleh berkas lain:
  ```celang
  publik! fungsi hitungPajak(nominal) {
      kembalikan nominal * 0.11
  }
  ```
* **Terbuka hanya untuk satu modul (`privat-modul!`)**: Data hanya boleh digunakan oleh berkas-berkas yang berada di dalam folder modul yang sama.

---

## G. Masalah Umum dan Cara Mengatasinya

Berikut beberapa kondisi yang mungkin Anda temui dan langkah penyelesaiannya:

### 1. Pesan: `gagal membaca file '...' no such file or directory`
* **Penyebab Terbukti**: Nama berkas yang Anda ketik salah, atau berkas tersebut berada di folder lain.
* **Langkah Penyelesaian**: Periksa kembali ejaan nama berkas. Pastikan terminal Anda sedang berada di folder yang sama dengan berkas tersebut (gunakan perintah `ls` di Linux/Mac atau `dir` di Windows untuk melihat daftar berkas).

### 2. Pesan: `Kesalahan: Tidak bisa membagi dengan angka nol!`
* **Penyebab Terbukti**: Ada perhitungan matematika di dalam kode yang membagi suatu angka dengan angka 0 (misalnya `100 / 0`). Komputer tidak dapat menghitung pembagian nol.
* **Langkah Penyelesaian**: Pastikan angka pembagi bernilai lebih besar dari nol, atau pasang pengaman menggunakan blok `coba ... tangkap`:
  ```celang
  coba {
      buat hasil = total / jumlah_orang
  } tangkap (kendala) {
      cetak("Terjadi kendala perhitungan:", kendala)
  }
  ```

### 3. Pesan: `Tidak dapat mengirim data ke saluran yang sudah ditutup`
* **Penyebab Terbukti**: Program Anda memanggil `saluran.kirim(...)` setelah sebelumnya memanggil `saluran.tutup()`.
* **Langkah Penyelesaian**: Pastikan seluruh pengiriman data sudah selesai sebelum Anda menutup saluran.

### 4. Perintah tidak dikenal saat mengetik `celang` di terminal
* **Penyebab Terbukti**: Biner aplikasi belum dipindahkan ke folder PATH sistem operasi Anda.
* **Langkah Penyelesaian**: Jalankan perintah dengan menyertakan tanda titik garis miring di depannya, misalnya `./celang run berkas.ce`, atau pindahkan biner `celang` ke folder sistem (seperti `/usr/local/bin` di Linux/Mac).

---

## H. Informasi untuk Pengembang Proyek

Bagian ini ditujukan bagi Anda yang ingin ikut memelihara atau mengembangkan kode internal Celang.

### Prasyarat Pengembangan
* Bahasa **Go** versi 1.20 atau yang lebih baru.
* Sistem operasi Linux, macOS, atau Windows.

### Menjalankan Seluruh Pengujian Otomatis
Sebelum mengajukan perubahan kode, pastikan semua unit pengujian berjalan sukses:
```bash
go test -v ./...
```

### Struktur Folder Kode Sumber
* [`cmd/celang`](file:///home/yosatsu/Coding/Celang/cmd/celang): Pintu masuk utama aplikasi terminal (`main.go`).
* [`cmd/celang-wasm`](file:///home/yosatsu/Coding/Celang/cmd/celang-wasm): Titik masuk kompilasi untuk peramban web (WebAssembly).
* [`pkg/parser`](file:///home/yosatsu/Coding/Celang/pkg/parser): Menguraikan teks kode sumber menjadi unit-unit kata dan pohon data.
* [`pkg/uast`](file:///home/yosatsu/Coding/Celang/pkg/uast): Struktur data Universal AST yang menyimpan inti logika program.
* [`pkg/syntax`](file:///home/yosatsu/Coding/Celang/pkg/syntax): Profil gaya bahasa dan pencetak format ulang kode.
* [`pkg/evaluator`](file:///home/yosatsu/Coding/Celang/pkg/evaluator): Eksekutor program, tipe objek, penanganan waktu, dan pustaka standar.
* [`pkg/compiler`](file:///home/yosatsu/Coding/Celang/pkg/compiler) & [`pkg/vm`](file:///home/yosatsu/Coding/Celang/pkg/vm): Kompiler instruksi mesin dan mesin virtual tumpukan berkecepatan tinggi.
* [`pkg/empathy`](file:///home/yosatsu/Coding/Celang/pkg/empathy): Penyusun pesan bantuan dan petunjuk kesalahan ramah.
* [`pkg/repl`](file:///home/yosatsu/Coding/Celang/pkg/repl): Lingkungan terminal interaktif.
* [`editors/vscode`](file:///home/yosatsu/Coding/Celang/editors/vscode): Berkas ekstensi resmi untuk editor Visual Studio Code.
* [`docs`](file:///home/yosatsu/Coding/Celang/docs): Berkas dokumentasi situs web yang ditenagai oleh VitePress.

---

## Lisensi

Proyek Celang dirilis di bawah lisensi terbuka **MIT License**. Anda bebas menggunakan, mempelajari, memodifikasi, dan membagikan aplikasi ini kepada siapa saja.
