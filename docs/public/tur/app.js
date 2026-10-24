const UI_TEXTS = {
  "id": {
    "subtitle": "A Tour of Celang",
    "tocBtn": "Daftar Isi",
    "footerPrev": "◀ Sebelumnya",
    "footerNext": "Selanjutnya ▶",
    "runBtn": "▶ Jalankan Kode",
    "running": "⏳ Menjalankan...",
    "resetBtn": "↺ Reset",
    "editorLabel": "KODE CELANG",
    "termTitle": "TERMINAL OUTPUT",
    "clearTerm": "Bersihkan",
    "expectedToggle": "📋 Ekspektasi Output Terminal",
    "challengeTitle": "Tantangan Praktik:",
    "docsLink": "📚 Dokumentasi",
    "tocTitle": "Daftar Isi Tur Celang",
    "tocSearchPlaceholder": "Cari materi tur (misal: variabel, transaksi, reaktif, fungsi)...",
    "tocNoResults": "🔎 Tidak ada materi tur yang cocok dengan pencarian.",
    "tocActiveBadge": "Aktif",
    "secPengantar": "🚀 Pengantar & Pilihan Bahasa",
    "secDasar": "📘 Bagian 1: Dasar Bahasa",
    "secSpesial": "✨ Bagian 2: Fitur Khas & Eksklusif",
    "badgePengantar": "Pengantar",
    "badgeDasar": "Dasar Bahasa",
    "badgeSpesial": "Fitur Spesial",
    "termPlaceholder": "// Tekan \"▶ Jalankan Kode\" atau Ctrl+Enter untuk melihat output...",
    "emptyCode": "// Kode sumber kosong!",
    "loadingWasm": "⏳ Mesin Celang WebAssembly sedang dimuat, silakan tunggu sebentar...",
    "wasmReady": "WASM Siap ✨",
    "wasmError": "❌ Gagal memuat WASM",
    "termDoneNoOutput": "(Program selesai tanpa luaran output)",
    "readDocs": "Baca dokumentasi lengkap materi ini →"
  },
  "en": {
    "subtitle": "A Tour of Celang",
    "tocBtn": "Table of Contents",
    "footerPrev": "◀ Previous",
    "footerNext": "Next ▶",
    "runBtn": "▶ Run Code",
    "running": "⏳ Running...",
    "resetBtn": "↺ Reset",
    "editorLabel": "CELANG CODE",
    "termTitle": "TERMINAL OUTPUT",
    "clearTerm": "Clear",
    "expectedToggle": "📋 Expected Console Output",
    "challengeTitle": "Practical Challenge:",
    "docsLink": "📚 Documentation",
    "tocTitle": "Celang Tour Table of Contents",
    "tocSearchPlaceholder": "Search tour lessons (e.g. variables, transactions, reactivity, functions)...",
    "tocNoResults": "🔎 No tour lessons found matching your search.",
    "tocActiveBadge": "Active",
    "secPengantar": "🚀 Introduction & Language",
    "secDasar": "📘 Part 1: Language Basics",
    "secSpesial": "✨ Part 2: Special & Exclusive Features",
    "badgePengantar": "Introduction",
    "badgeDasar": "Language Basics",
    "badgeSpesial": "Special Features",
    "termPlaceholder": "// Press \"▶ Run Code\" or Ctrl+Enter to see output...",
    "emptyCode": "// Source code is empty!",
    "loadingWasm": "⏳ Loading Celang WebAssembly engine, please wait...",
    "wasmReady": "WASM Ready ✨",
    "wasmError": "❌ Failed to load WASM",
    "termDoneNoOutput": "(Program finished with no output)",
    "readDocs": "Read full documentation on this topic →"
  }
};

const SLIDES_ID = [
  {
    "id": 1,
    "title": "[1/36] Pengantar Tur Celang",
    "section": "Pengantar",
    "docLink": "/panduan/pengenalan",
    "description": "\n      <p>Selamat datang di <strong>Tur Celang</strong> (<em>A Tour of Celang</em>)! Di sini kamu bisa belajar bahasa pemrograman Celang langkah demi langkah secara langsung dan seru.</p>\n      <p>Layar ini dibagi menjadi dua panel:</p>\n      <ul>\n        <li><strong>Panel Kiri:</strong> Berisi penjelasan materi, panduan bahasa manusia yang mudah dipahami, serta tantangan praktik kecil.</li>\n        <li><strong>Panel Kanan:</strong> Tempat kamu mencoba menulis kode secara langsung dan melihat hasil jalannya di terminal hitam bawah.</li>\n      </ul>\n      <p>Programmu langsung dijalankan di browser menggunakan teknologi <strong>WebAssembly</strong>, jadi cepat, aman, dan tanpa perlu instal apa-apa di komputermu!</p>\n    ",
    "task": "Klik tombol <strong>Jalankan Kode ▶</strong> di panel kanan atas untuk menguji coba eksekusi program Celang pertamamu.",
    "code": "buat salam = \"Selamat Datang di Tur Celang!\"\ncetak(salam)\ncetak(\"Silakan klik tombol Selanjutnya di bawah untuk memulai belajar.\")",
    "expected": "Selamat Datang di Tur Celang!\nSilakan klik tombol Selanjutnya di bawah untuk memulai belajar."
  },
  {
    "id": 2,
    "title": "[2/36] Pemilihan Bahasa Tutorial",
    "section": "Pengantar",
    "docLink": "/panduan/pengenalan",
    "description": "\n      <p>Tur ini mendukung dua bahasa pengantar tutorial:</p>\n      <div style=\"display: grid; grid-template-columns: 1fr 1fr; gap: 12px; margin: 16px 0;\">\n        <div class=\"lang-select-card active\" id=\"card-lang-id\" style=\"padding: 14px; border: 2px solid var(--accent-primary); border-radius: 8px; cursor: pointer; text-align: center; background: rgba(59, 130, 246, 0.08);\">\n          <div style=\"font-size: 24px; margin-bottom: 4px;\">🇮🇩</div>\n          <strong>Bahasa Indonesia</strong>\n          <p style=\"font-size: 12px; margin: 4px 0 0; opacity: 0.8;\">Penjelasan dalam bahasa Indonesia sehari-hari.</p>\n        </div>\n        <div class=\"lang-select-card\" id=\"card-lang-en\" style=\"padding: 14px; border: 1px solid var(--border-color); border-radius: 8px; cursor: pointer; text-align: center;\">\n          <div style=\"font-size: 24px; margin-bottom: 4px;\">🇬🇧</div>\n          <strong>English</strong>\n          <p style=\"font-size: 12px; margin: 4px 0 0; opacity: 0.8;\">Explanations in international English.</p>\n        </div>\n      </div>\n      <p>Kamu juga bebas mengubah bahasa kapan saja melalui tombol bendera di pojok kanan atas.</p>\n    ",
    "task": "Pilih bahasa di atas, atau klik tombol <strong>Selanjutnya ▶</strong> untuk melanjutkan tur dengan Bahasa Indonesia.",
    "code": "buat bahasa = \"Bahasa Indonesia\"\ncetak(\"Bahasa pengantar terpilih:\", bahasa)\ncetak(\"Selamat menikmati Tur Celang!\")",
    "expected": "Bahasa pengantar terpilih: Bahasa Indonesia\nSelamat menikmati Tur Celang!"
  },
  {
    "id": 3,
    "title": "[3/36] Halo Dunia",
    "section": "Dasar Bahasa",
    "docLink": "/panduan/mulai-cepat",
    "description": "\n      <p>Mari mulai dari dasar. Fungsi bawaan <code>cetak(...)</code> digunakan untuk menampilkan teks atau data ke terminal output.</p>\n      <p>Teks diapit oleh tanda kutip ganda (<code>\"...\"</code>). Anda dapat memasukkan satu atau beberapa nilai sekaligus di dalam <code>cetak(...)</code> dengan pemisah koma <code>,</code>.</p>\n    ",
    "task": "Coba ubah teks <code>\"Halo, Dunia!\"</code> menjadi <code>\"Halo, Celang!\"</code>, lalu amati hasilnya di terminal.",
    "code": "cetak(\"Halo, Dunia!\")",
    "expected": "Halo, Celang!"
  },
  {
    "id": 4,
    "title": "[4/36] Variabel & Nilai",
    "section": "Dasar Bahasa",
    "docLink": "/dasar/variabel-tipe",
    "description": "\n      <p>Di Celang, variabel dideklarasikan menggunakan kata kunci <code>buat</code>. Variabel bersifat dinamis dan nilainya dapat diperbarui kapan saja tanpa menulis ulang kata kunci <code>buat</code>.</p>\n      <p>Jika Anda mencetak beberapa data yang dipisahkan dengan tanda koma <code>,</code>, Celang akan otomatis menampilkannya dengan spasi pemisah di terminal.</p>\n    ",
    "task": "Coba ubah nilai <code>umur</code> menjadi <code>25</code>, lalu amati perubahan hasil cetakannya.",
    "code": "buat nama = \"Budi\"\nbuat umur = 20\n\numur = 21\n\ncetak(nama, \"berumur\", umur, \"tahun\")",
    "expected": "Budi berumur 25 tahun"
  },
  {
    "id": 5,
    "title": "[5/36] Nilai Tetap (Konstanta)",
    "section": "Dasar Bahasa",
    "docLink": "/dasar/variabel-tipe",
    "description": "\n      <p>Gunakan kata kunci <code>konstanta</code> untuk mendeklarasikan nilai tetap yang tidak boleh diubah sepanjang eksekusi program.</p>\n      <p>Mencoba mengisi ulang (<em>reassign</em>) nilai konstanta akan menghasilkan galat (<em>error</em>). Menggunakan konstanta sangat dianjurkan untuk nilai acuan seperti konfigurasi, nilai matematika, atau identitas tetap.</p>\n    ",
    "task": "Tambahkan baris <code>konstanta VERSI = 1</code> dan <code>cetak(\"Versi:\", VERSI)</code> di akhir kode.",
    "code": "konstanta PI = 3.14159\nkonstanta NAMA_APLIKASI = \"Sistem Keuangan\"\n\ncetak(\"Aplikasi:\", NAMA_APLIKASI)\ncetak(\"Nilai PI:\", PI)",
    "expected": "Aplikasi: Sistem Keuangan\nNilai PI: 3.14159\nVersi: 1"
  },
  {
    "id": 6,
    "title": "[6/36] Tipe Data & Contoh Nilainya",
    "section": "Dasar Bahasa",
    "docLink": "/dasar/variabel-tipe",
    "description": "\n      <p>Celang memiliki beragam tipe data bawaan untuk menangani kebutuhan nyata:</p>\n      <table class=\"table-types\" style=\"width:100%; font-size:13px; border-collapse:collapse; margin-bottom:12px;\">\n        <thead>\n          <tr style=\"background:rgba(255,255,255,0.05); text-align:left;\">\n            <th style=\"padding:6px 8px; border:1px solid var(--border-color);\">Tipe Data</th>\n            <th style=\"padding:6px 8px; border:1px solid var(--border-color);\">Kegunaan</th>\n            <th style=\"padding:6px 8px; border:1px solid var(--border-color);\">Contoh Nilai</th>\n          </tr>\n        </thead>\n        <tbody>\n          <tr><td style=\"padding:6px 8px; border:1px solid var(--border-color);\"><code>Angka</code></td><td style=\"padding:6px 8px; border:1px solid var(--border-color);\">Bilangan bulat &amp; desimal</td><td style=\"padding:6px 8px; border:1px solid var(--border-color);\"><code>42</code>, <code>3.14</code></td></tr>\n          <tr><td style=\"padding:6px 8px; border:1px solid var(--border-color);\"><code>Teks</code></td><td style=\"padding:6px 8px; border:1px solid var(--border-color);\">Untaian huruf/karakter</td><td style=\"padding:6px 8px; border:1px solid var(--border-color);\"><code>\"Halo\"</code>, <code>'Budi'</code></td></tr>\n          <tr><td style=\"padding:6px 8px; border:1px solid var(--border-color);\"><code>Kebenaran</code></td><td style=\"padding:6px 8px; border:1px solid var(--border-color);\">Nilai benar/salah</td><td style=\"padding:6px 8px; border:1px solid var(--border-color);\"><code>benar</code>, <code>salah</code></td></tr>\n          <tr><td style=\"padding:6px 8px; border:1px solid var(--border-color);\"><code>Larik</code></td><td style=\"padding:6px 8px; border:1px solid var(--border-color);\">Kumpulan berurutan</td><td style=\"padding:6px 8px; border:1px solid var(--border-color);\"><code>[1, 2, 3]</code></td></tr>\n          <tr><td style=\"padding:6px 8px; border:1px solid var(--border-color);\"><code>Kamus</code></td><td style=\"padding:6px 8px; border:1px solid var(--border-color);\">Pasangan kunci-nilai</td><td style=\"padding:6px 8px; border:1px solid var(--border-color);\"><code>{\"a\": 1}</code></td></tr>\n          <tr><td style=\"padding:6px 8px; border:1px solid var(--border-color);\"><code>Kosong</code></td><td style=\"padding:6px 8px; border:1px solid var(--border-color);\">Ketiadaan nilai</td><td style=\"padding:6px 8px; border:1px solid var(--border-color);\"><code>kosong</code></td></tr>\n        </tbody>\n      </table>\n      <p>Setiap tipe data dapat dicetak langsung menggunakan perintah <code>cetak(...)</code>.</p>\n    ",
    "task": "Ubah nilai <code>sedangAktif</code> menjadi <code>salah</code>, lalu amati hasilnya.",
    "code": "buat umur = 28\nbuat tinggi = 175.5\nbuat pesan = \"Semangat Belajar!\"\nbuat sedangAktif = benar\nbuat statusLogin = kosong\n\ncetak(\"Umur (Angka):\", umur)\ncetak(\"Tinggi (Angka):\", tinggi)\ncetak(\"Pesan (Teks):\", pesan)\ncetak(\"Sedang aktif (Kebenaran):\", sedangAktif)\ncetak(\"Status login (Kosong):\", statusLogin)",
    "expected": "Umur (Angka): 28\nTinggi (Angka): 175.5\nPesan (Teks): Semangat Belajar!\nSedang aktif (Kebenaran): false\nStatus login (Kosong): null"
  },
  {
    "id": 7,
    "title": "[7/36] Interpolasi Teks (f-string)",
    "section": "Dasar Bahasa",
    "docLink": "/fitur-khas/interpolasi-teks",
    "description": "\n      <p>Untuk menggabungkan teks dengan nilai variabel atau hasil hitungan tanpa repot menyambung tanda tambah <code>+</code>, gunakan format <strong>f-string</strong> mirip Python:</p>\n      <p>Awali teks dengan huruf <code>f</code> lalu sertakan ekspresi di dalam kurung kurawal <code>${...}</code>.</p>\n      <p>Contoh: <code>f\"Halo ${nama}, umurmu ${umur + 1} tahun\"</code>.</p>\n    ",
    "task": "Ubah variabel <code>tabungan</code> menjadi <code>50000</code> sehingga baris kedua mencetak: <code>Jika ditambah 10000: Rp60000</code>.",
    "code": "buat nama = \"Andi\"\nbuat tabungan = 25000\n\ncetak(f\"Nasabah: ${nama}\")\ncetak(f\"Jika ditambah 10000: Rp${tabungan + 10000}\")",
    "expected": "Nasabah: Andi\nJika ditambah 10000: Rp60000"
  },
  {
    "id": 8,
    "title": "[8/36] Percabangan Bersyarat (jika)",
    "section": "Dasar Bahasa",
    "docLink": "/dasar/kontrol-alur",
    "description": "\n      <p>Untuk mengarahkan alur logika, gunakan <code>jika</code>, <code>selain jika</code>, dan <code>selain</code>.</p>\n      <p>Kondisi tidak wajib diapit tanda kurung, namun blok kode wajib diapit kurung kurawal <code>{ ... }</code>.</p>\n    ",
    "task": "Ubah nilai <code>nilaiUjian</code> menjadi <code>65</code> agar pesan yang tercetak adalah <code>Status: Perlu Remedial</code>.",
    "code": "buat nilaiUjian = 85\n\njika nilaiUjian >= 80 {\n    cetak(\"Status: Lulus Memuaskan\")\n} selain jika nilaiUjian >= 60 {\n    cetak(\"Status: Perlu Remedial\")\n} selain {\n    cetak(\"Status: Mengulang Ujian\")\n}",
    "expected": "Status: Perlu Remedial"
  },
  {
    "id": 9,
    "title": "[9/36] Pemilihan Nilai (pilih ... bila)",
    "section": "Dasar Bahasa",
    "docLink": "/dasar/kontrol-alur",
    "description": "\n      <p>Struktur <code>pilih</code> digunakan untuk mencocokkan suatu nilai ke beberapa kemungkinan cabang menggunakan kata <code>bila</code>.</p>\n      <p>Cabang <code>lainnya:</code> akan dieksekusi jika tidak ada satupun kondisi <code>bila</code> yang terpenuhi. Berbeda dengan bahasa turunan C, Celang tidak membutuhkan instruksi <code>break</code> manual.</p>\n    ",
    "task": "Ubah nilai <code>peran</code> menjadi <code>\"editor\"</code> sehingga peran yang tercetak adalah <code>Akses: Mengelola Konten</code>.",
    "code": "buat peran = \"tamu\"\n\npilih peran {\n    bila \"admin\":\n        cetak(\"Akses: Penuh ke Sistem\")\n    bila \"editor\":\n        cetak(\"Akses: Mengelola Konten\")\n    bila \"penulis\":\n        cetak(\"Akses: Menulis Draf\")\n    lainnya:\n        cetak(\"Akses: Hanya Membaca\")\n}",
    "expected": "Akses: Mengelola Konten"
  },
  {
    "id": 10,
    "title": "[10/36] Perulangan Bersyarat (selama)",
    "section": "Dasar Bahasa",
    "docLink": "/dasar/perulangan",
    "description": "\n      <p>Gunakan <code>selama</code> (ekuivalen dengan <code>while</code>) untuk mengeksekusi blok kode berulang kali sepanjang kondisinya bernilai <code>benar</code>.</p>\n      <p>Pastikan variabel kondisi selalu diperbarui di dalam blok agar perulangan berhenti pada saat yang tepat.</p>\n    ",
    "task": "Ubah kondisi <code>hitung &lt;= 3</code> menjadi <code>hitung &lt;= 5</code> agar angka dihitung sampai 5.",
    "code": "buat hitung = 1\n\nselama hitung <= 3 {\n    cetak(\"Hitungan ke:\", hitung)\n    hitung = hitung + 1\n}\n\ncetak(\"Selesai!\")",
    "expected": "Hitungan ke: 1\nHitungan ke: 2\nHitungan ke: 3\nHitungan ke: 4\nHitungan ke: 5\nSelesai!"
  },
  {
    "id": 11,
    "title": "[11/36] Perulangan Daftar & Rentang (untuk ... dalam)",
    "section": "Dasar Bahasa",
    "docLink": "/dasar/perulangan",
    "description": "\n      <p>Perulangan <code>untuk ... dalam</code> adalah cara paling praktis untuk melintasi isi daftar data atau deret angka secara berurutan:</p>\n      <ul>\n        <li><strong>Deret Angka (Rentang)</strong>: Gunakan titik dua kali seperti <code>1..5</code> untuk membuat urutan angka dari 1 sampai 5.</li>\n        <li><strong>Daftar Kumpulan (Array)</strong>: Gunakan <code>untuk nama in daftar</code> untuk menyapa atau mengolah setiap anggota satu per satu dari awal sampai habis.</li>\n      </ul>\n    ",
    "task": "Tambahkan buah <code>\"Pisang\"</code> ke dalam daftar buah: <code>[\"Apel\", \"Jeruk\", \"Mangga\", \"Pisang\"]</code>.",
    "code": "// 1. Melintasi rentang angka 1 sampai 3:\nuntuk angka dalam 1..3 {\n    cetak(\"Angka:\", angka)\n}\n\n// 2. Melintasi setiap isi daftar buah:\nbuat buah = [\"Apel\", \"Jeruk\", \"Mangga\"]\nuntuk item dalam buah {\n    cetak(\"Buah:\", item)\n}",
    "expected": "Angka: 1\nAngka: 2\nAngka: 3\nBuah: Apel\nBuah: Jeruk\nBuah: Mangga\nBuah: Pisang"
  },
  {
    "id": 12,
    "title": "[12/36] Lompat & Rem Perulangan (lanjut & berhenti)",
    "section": "Dasar Bahasa",
    "docLink": "/dasar/perulangan",
    "description": "\n      <p>Saat perulangan sedang berputar, terkadang ada giliran tertentu yang mau dilewati, atau ingin langsung disetop sama sekali:</p>\n      <ul>\n        <li><strong><code>lanjut</code> (Lewati Giliran / Continue)</strong>: Ibarat bilang <em>\"giliran yang ini lewati aja, langsung lanjut ke nomor berikutnya!\"</em>. Kode di bawahnya tidak akan dijalankan untuk giliran itu.</li>\n        <li><strong><code>berhenti</code> (Rem Mendadak / Break)</strong>: Ibarat bilang <em>\"sudah cukup, setop sekarang juga!\"</em>. Perulangan langsung bubar seketika dan melangkah ke baris di luar perulangan.</li>\n      </ul>\n    ",
    "task": "Hapus blok <code>jika i == 2 { lanjut }</code> agar angka 2 tidak dilewati dan tercetak.",
    "code": "untuk i dalam 1..5 {\n    jika i == 2 {\n        lanjut\n    }\n    jika i == 4 {\n        berhenti\n    }\n    cetak(\"Nilai i:\", i)\n}",
    "expected": "Nilai i: 1\nNilai i: 2\nNilai i: 3"
  },
  {
    "id": 13,
    "title": "[13/36] Kumpulan Data: Larik (Array)",
    "section": "Dasar Bahasa",
    "docLink": "/fitur-khas/koleksi-array-map",
    "description": "\n      <p>Larik (<em>Array</em>) menyimpan daftar elemen berurutan yang dibungkus oleh kurung siku <code>[...]</code>. Indeks elemen dimulai dari angka <code>0</code>.</p>\n      <p>Celang menyediakan metode praktis seperti <code>.push()</code> untuk menambahkan elemen dan fungsi <code>panjang(...)</code> untuk menghitung total elemen.</p>\n    ",
    "task": "Tambahkan baris <code>cetak(\"Elemen kedua:\", angka[1])</code> di bawah cetak elemen pertama.",
    "code": "buat angka = [10, 20, 30]\n\nangka.push(40)\ncetak(\"Elemen pertama:\", angka[0])\ncetak(\"Jumlah elemen:\", panjang(angka))\ncetak(\"Seluruh data:\", angka)",
    "expected": "Elemen pertama: 10\nElemen kedua: 20\nJumlah elemen: 4\nSeluruh data: [10, 20, 30, 40]"
  },
  {
    "id": 14,
    "title": "[14/36] Kamus Data (Map)",
    "section": "Dasar Bahasa",
    "docLink": "/fitur-khas/koleksi-array-map",
    "description": "\n      <p>Kamus (<em>Map</em>) menyimpan pasangan kunci dan nilai (<em>key-value pairs</em>) yang dibungkus dengan tanda kurung kurawal <code>{ ... }</code>.</p>\n      <p>Anda dapat membaca nilai properti kamus dengan dua cara fleksibel: melalui kurung siku <code>peta[\"kunci\"]</code> atau langsung memakai tanda titik <code>peta.kunci</code>.</p>\n    ",
    "task": "Tambahkan baris <code>cetak(\"Usia:\", pengguna.usia)</code> di akhir kode.",
    "code": "buat pengguna = {\n    \"nama\": \"Andi\",\n    \"pekerjaan\": \"Pengembang\",\n    \"usia\": 25\n}\n\ncetak(\"Nama:\", pengguna[\"nama\"])\ncetak(\"Pekerjaan:\", pengguna.pekerjaan)",
    "expected": "Nama: Andi\nPekerjaan: Pengembang\nUsia: 25"
  },
  {
    "id": 15,
    "title": "[15/36] Membongkar Data (Destrukturisasi)",
    "section": "Dasar Bahasa",
    "docLink": "/fitur-khas/destrukturisasi",
    "description": "\n      <p>Bayangkan kamu baru pulang belanja dan membawa sekantong belanjaan berisi barang. Daripada mengambil barang satu per satu berkali-kali, kamu langsung membongkar kantong dan menaruh tiap barang ke wadahnya masing-masing sekaligus.</p>\n      <p>Fitur <strong>Destrukturisasi</strong> memungkinkan kamu memecah isi daftar (<em>Array</em>) atau kamus data (<em>Map</em>) langsung ke dalam beberapa variabel terpisah dalam satu baris kode:</p>\n      <ul>\n        <li><strong>Membongkar Larik Berdasarkan Urutan</strong>: <code>buat [x, y] = koordinat</code> (elemen ke-0 masuk ke <code>x</code>, elemen ke-1 masuk ke <code>y</code>).</li>\n        <li><strong>Membongkar Kamus Berdasarkan Nama Kolom</strong>: <code>buat { nama, kota } = profil</code> (langsung mengambil kunci <code>\"nama\"</code> dan <code>\"kota\"</code>).</li>\n      </ul>\n    ",
    "task": "Ubah nilai koordinat menjadi <code>[500, 750]</code> sehingga output mencetak <code>Koordinat: 500 750</code>.",
    "code": "buat koordinat = [100, 250]\nbuat profil = {\n    \"nama\": \"Citra\",\n    \"kota\": \"Surabaya\"\n}\n\n// 1. Membongkar larik sesuai urutan:\nbuat [x, y] = koordinat\n\n// 2. Membongkar kamus sesuai nama kolom:\nbuat { nama, kota } = profil\n\ncetak(\"Koordinat:\", x, y)\ncetak(\"Profil:\", nama, \"dari\", kota)",
    "expected": "Koordinat: 500 750\nProfil: Citra dari Surabaya"
  },
  {
    "id": 16,
    "title": "[16/36] Fungsi, Pembatasan Tipe, & Nilai Balik",
    "section": "Dasar Bahasa",
    "docLink": "/dasar/fungsi",
    "description": "\n      <p>Fungsi itu ibarat resep atau mesin peracik: kamu kasih bahan masukannya (parameter), mesin bekerja, lalu memberikan hasil akhirnya kembali ke kamu.</p>\n      <p>Di Celang, kamu bisa menulis fungsi secara santai maupun dengan pembatasan tipe data yang tegas:</p>\n      <ul>\n        <li><strong>Batasan Tipe Masukan (Constraint)</strong>: Tegaskan jenis data parameter menggunakan kata <code>itu</code> (contoh: <code>panjang itu Angka</code>) atau gaya simbol <code>panjang: Angka</code>.</li>\n        <li><strong>Penegasan Nilai Balik</strong>: Gunakan kata <code>mengembalikan Angka</code> atau simbol panah <code>-&gt; Angka</code> di belakang kurung parameter untuk menegaskan tipe hasil yang dikeluarkan.</li>\n        <li><strong>Mengirimkan Hasil</strong>: Gunakan kata kunci <code>kembalikan &lt;nilai&gt;</code> di dalam fungsi.</li>\n      </ul>\n    ",
    "task": "Panggil fungsi tersebut dengan nilai panjang <code>10</code> dan lebar <code>3</code>: <code>buat hasil = hitungLuas(10, 3)</code>.",
    "code": "// Parameter dibatasi tipe Angka dan mengembalikan Angka:\nfungsi hitungLuas(panjang: Angka, lebar: Angka) -> Angka {\n    kembalikan panjang * lebar\n}\n\nbuat hasil = hitungLuas(5, 4)\ncetak(\"Luas persegi:\", hasil)",
    "expected": "Luas persegi: 30"
  },
  {
    "id": 17,
    "title": "[17/36] Rumus Ringkas (Panah =>) & Olah Data",
    "section": "Dasar Bahasa",
    "docLink": "/fitur-khas/fungsi-ringkas",
    "description": "\n      <p>Terkadang kita hanya butuh rumus perhitungan kecil tanpa mau repot mengetik kata <code>fungsi</code>, kurung kurawal, dan <code>kembalikan</code>.</p>\n      <p>Celang menyediakan tanda <strong>Panah Ringkas (<code>=&gt;</code>)</strong> (<em>arrow function</em>):</p>\n      <ul>\n        <li><strong>Satu Masukan</strong>: <code>buat kaliDua = x =&gt; x * 2</code> (langsung menghasilkan nilai di kanan panah).</li>\n        <li><strong>Dua Masukan</strong>: <code>buat jumlah = (a, b) =&gt; a + b</code>.</li>\n        <li><strong>Olah Kumpulan Data</strong>: Dipadukan dengan fungsi bawaan seperti <code>saring(daftar, rumus)</code> (<em>filter</em>) dan <code>petakan(daftar, rumus)</code> (<em>map</em>), mengolah daftar jadi sangat singkat dan bersih!</li>\n      </ul>\n    ",
    "task": "Ubah rumus <code>kaliDua</code> menjadi <code>x =&gt; x * 10</code> sehingga hasil angka genap dikali sepuluh menjadi <code>[20, 40, 60]</code>.",
    "code": "buat kaliDua = x => x * 2\nbuat jumlah = (a, b) => a + b\n\nbuat daftar = [1, 2, 3, 4, 5, 6]\n\n// Saring hanya angka genap, lalu kalikan:\nbuat genap = saring(daftar, x => x % 2 == 0)\nbuat hasil = petakan(genap, kaliDua)\n\ncetak(\"Hasil olah data:\", hasil)\ncetak(\"Total 10 + 20:\", jumlah(10, 20))",
    "expected": "Hasil olah data: [20, 40, 60]\nTotal 10 + 20: 30"
  },
  {
    "id": 18,
    "title": "[18/36] Referensi & Rujukan (rujuk)",
    "section": "Dasar Bahasa",
    "docLink": "/dasar/rujuk",
    "description": "\n      <p>Secara default, saat kamu mengirimkan angka ke variabel baru atau ke fungsi, Celang membuat salinan terpisah (<em>pass-by-value</em>). Mengubah salinan tidak akan mengubah data aslinya.</p>\n      <p>Tapi bagaimana jika kamu ingin dua nama mengendalikan kotak memori yang sama persis (seperti <em>dua remote untuk satu televisi</em>)? Di sinilah kata kunci <strong><code>rujuk</code></strong> (alias: <code>ref</code>) bekerja:</p>\n      <ul>\n        <li><strong>Alias Variabel (<code>rujuk alias = target</code>)</strong>: Memberikan nama panggilan kedua ke variabel asli. Mengubah <code>alias</code> otomatis mengubah <code>target</code>!</li>\n        <li><strong>Parameter Fungsi (<code>rujuk param</code>)</strong>: Memungkinkan fungsi langsung memodifikasi variabel pemanggil di luar tanpa sintaks pointer <code>*</code> dan <code>&</code> yang rumit.</li>\n      </ul>\n    ",
    "task": "Ubah bonus menjadi <code>50000</code> sehingga tabungan akhir menjadi 150000.",
    "code": "// 1. Deklarasi alias rujukan:\nbuat skor = 50\nrujuk skorPemain = skor\nskorPemain = 95\n\n// 2. Parameter fungsi langsung memutasi variabel luar:\nfungsi beriBonus(rujuk saldo, bonus) {\n    saldo = saldo + bonus\n}\n\nbuat tabungan = 100000\nberiBonus(tabungan, 25000)\n\ncetak(\"Skor asli:\", skor)\ncetak(\"Tabungan akhir:\", tabungan)",
    "expected": "Skor asli: 95\nTabungan akhir: 150000"
  },
  {
    "id": 19,
    "title": "[19/36] Cetak Biru Data (bentuk & struct)",
    "section": "Dasar Bahasa",
    "docLink": "/fitur-khas/bentuk-struct",
    "description": "\n      <p>Kalau kamu membuat program yang menangani data dunia nyata—seperti data rekening bank atau akun pengguna—menyimpan nama dan saldo di variabel terpisah-pisah bisa cepat bikin kode berantakan.</p>\n      <p>Di Celang, kamu bisa mengelompokkan data dan tindakannya ke dalam satu <strong>cetak biru</strong> rapi menggunakan kata <code>bentuk</code>:</p>\n      <ul>\n        <li><strong>Properti &amp; Metode</strong>: Tentukan apa saja data yang dimiliki (seperti <code>pemilik</code> dan <code>saldo</code>) serta tindakan yang bisa dilakukannya (fungsi seperti <code>setor</code>).</li>\n        <li><strong>Tanda <code>@</code> (Milik Sendiri)</strong>: Di dalam fungsi bentuk, gunakan tanda <code>@</code> untuk mengakses data atau metode milik objek ini (misal: <code>@saldo += jumlah</code>). Tanda <code>@</code> wajib diikuti nama anggota dan tidak bisa berdiri sendiri.</li>\n        <li><strong>Anti-Ambiguitas (Kapan Pakai <code>diri</code>?)</strong>: Jika ingin mengirim objek ini seutuhnya ke fungsi lain atau mengembalikannya, <strong>wajib pakai kata kunci <code>diri</code></strong>, bukan <code>@</code>! Contoh: <code>lapor(diri, @pemilik)</code>, <strong>bukan</strong> <code>lapor(@, @pemilik)</code>.</li>\n        <li><strong>Membuat Objek Nyata</strong>: Cukup panggil nama bentuknya, misalnya: <code>buat akun = Rekening(\"Budi\", 50000)</code>.</li>\n      </ul>\n    ",
    "task": "Ubah jumlah setor menjadi <code>akun.setor(50000)</code> sehingga saldo akhir menjadi 100000.",
    "code": "fungsi lapor(obj, nama) {\n    cetak(\"Laporan nasabah:\", nama)\n}\n\nbentuk Rekening {\n    pemilik: Teks,\n    saldo: Angka = 0,\n\n    fungsi setor(jumlah: Angka) {\n        @saldo += jumlah\n        kembalikan @saldo\n    }\n\n    fungsi cek() {\n        // Anti-ambiguitas: wajib pakai 'diri', bukan '@'!\n        lapor(diri, @pemilik)\n    }\n}\n\nbuat akun = Rekening(\"Budi\", 50000)\nakun.cek()\nakun.setor(25000)\n\ncetak(\"Pemilik:\", akun.pemilik)\ncetak(\"Saldo akhir:\", akun.saldo)",
    "expected": "Laporan nasabah: Budi\nPemilik: Budi\nSaldo akhir: 100000"
  },
  {
    "id": 20,
    "title": "[20/36] Wadah Fleksibel & Pembatas Tipe (Generik)",
    "section": "Dasar Bahasa",
    "docLink": "/fitur-khas/generik",
    "description": "\n      <p>Bayangkan kamu membuat sebuah kotak penyimpanan. Daripada membuat <code>KotakKhususAngka</code> dan <code>KotakKhususTeks</code> terpisah padahal bentuk kotaknya sama persis, kamu cukup membuat satu cetak biru <code>bentuk Kotak&lt;T&gt;</code>!</p>\n      <ul>\n        <li><strong>Tanda <code>&lt;T&gt;</code></strong>: Tipe data baru ditentukan saat kotak itu benar-benar dibuat (misal <code>Kotak&lt;Angka&gt;(100)</code> atau <code>Kotak&lt;Teks&gt;(\"Surat\")</code>).</li>\n        <li><strong>Pembatas <code>harus</code></strong>: Jika ingin memastikan hanya tipe tertentu yang boleh masuk, gunakan kata alami <code>harus</code> (contoh: <code>&lt;Nilai harus Angka&gt;</code>).</li>\n      </ul>\n    ",
    "task": "Ubah nilai angka di dalam <code>kotakAngka</code> menjadi <code>500</code> sehingga mencetak <code>Skor Game: 500</code>.",
    "code": "bentuk Kotak<T> {\n    label: Teks,\n    isi: T\n}\n\nbuat kotakAngka = Kotak<Angka>(\"Skor Game\", 100)\nbuat kotakTeks = Kotak<Teks>(\"Pesan\", \"Halo Dunia\")\n\ncetak(kotakAngka.label + \":\", kotakAngka.isi)\ncetak(kotakTeks.label + \":\", kotakTeks.isi)",
    "expected": "Skor Game: 500\nPesan: Halo Dunia"
  },
  {
    "id": 21,
    "title": "[21/36] Penanganan Galat (coba ... tangkap)",
    "section": "Dasar Bahasa",
    "docLink": "/dasar/penanganan-galat",
    "description": "\n      <p>Untuk menjaga program agar tidak mendadak berhenti saat terjadi kondisi tak terduga, gunakan blok <code>coba</code>, <code>tangkap</code>, dan <code>akhirnya</code>.</p>\n      <p>Anda dapat memicu galat secara manual menggunakan <code>lempar</code>. Blok <code>akhirnya</code> bersifat opsional dan dijamin akan selalu dieksekusi di akhir, baik terjadi galat maupun tidak.</p>\n    ",
    "task": "Ubah pesan galat yang dilempar menjadi <code>\"Angka pembagi tidak valid!\"</code>.",
    "code": "coba {\n    buat pembagi = 0\n    jika pembagi == 0 {\n        lempar \"Tidak bisa membagi dengan nol!\"\n    }\n} tangkap (kesalahan) {\n    cetak(\"Terjadi galat:\", kesalahan)\n} akhirnya {\n    cetak(\"Proses selesai.\")\n}",
    "expected": "Terjadi galat: Angka pembagi tidak valid!\nProses selesai."
  },
  {
    "id": 22,
    "title": "[22/36] Penjelajah Waktu Variabel (Time-Travel)",
    "section": "Fitur Spesial",
    "docLink": "/fitur-khas/time-travel",
    "description": "\n      <p>Di hampir semua bahasa pemrograman, nilai lama variabel akan langsung hilang saat ditimpa nilai baru. Namun di Celang, variabel secara alami mengingat masa lalunya!</p>\n      <p>Setiap variabel memiliki properti waktu bawaan: <code>.sebelumnya</code> untuk nilai tepat sebelum mutasi terakhir, <code>.awal</code> untuk nilai mula-mula saat pertama dibuat, dan <code>.histori</code> untuk melihat seluruh rekaman perubahannya.</p>\n    ",
    "task": "Tambahkan baris <code>skor = 100</code> di bawah <code>skor = 50</code>, lalu amati nilai pada <code>skor.sebelumnya</code>.",
    "code": "buat skor = 10\nskor = 25\nskor = 50\n\ncetak(\"Skor saat ini:\", skor)\ncetak(\"Skor sebelumnya:\", skor.sebelumnya)\ncetak(\"Skor awal:\", skor.awal)\ncetak(\"Riwayat skor:\", skor.histori)",
    "expected": "Skor saat ini: 100\nSkor sebelumnya: 50\nSkor awal: 10\nRiwayat skor: [10, 25, 50, 100]"
  },
  {
    "id": 23,
    "title": "[23/36] Transaksi & Rollback Otomatis (transaksi)",
    "section": "Fitur Spesial",
    "docLink": "/fitur-khas/transaksi",
    "description": "\n      <p>Di bahasa lain seperti Python atau JS, kalau kamu mengubah saldo beberapa orang lalu di tengah jalan program melempar error, saldo mereka terlanjur rusak separuh jalan (<em>dirty state</em>).</p>\n      <p>Celang membawa fitur transaksi database langsung ke level bahasa pemrograman:</p>\n      <ul>\n        <li><strong>Blok <code>transaksi { ... }</code></strong>: Menjalankan serangkaian perubahan variabel secara aman dan atomik.</li>\n        <li><strong><code>batalkan \"alasan\"</code></strong>: Jika ada kondisi gagal, Celang seketika mengembalikan seluruh nilai variabel ke kondisi semula (<em>rollback</em>).</li>\n        <li><strong><code>pulihkan (alasan) { ... }</code></strong>: Menangkap alasan pembatalan dan melanjutkan program dengan selamat.</li>\n      </ul>\n    ",
    "task": "Ubah saldo awal Andi menjadi <code>100.ribu</code> sehingga saldo akhir Andi tetap aman tercetak <code>Rp 100.000</code>.",
    "code": "buat saldoAndi = 50.ribu\nbuat saldoBudi = 50.ribu\n\ntransaksi {\n    saldoAndi -= 250.ribu\n    saldoBudi += 250.ribu\n\n    jika saldoAndi < 0 {\n        batalkan \"Saldo Andi tidak mencukupi untuk transfer!\"\n    }\n} pulihkan (alasan) {\n    cetak(\"Pemberitahuan:\", alasan)\n}\n\ncetak(\"Saldo akhir Andi tetap aman:\", saldoAndi.rupiah)\ncetak(\"Saldo akhir Budi tetap aman:\", saldoBudi.rupiah)",
    "expected": "Pemberitahuan: Saldo Andi tidak mencukupi untuk transfer!\nSaldo akhir Andi tetap aman: Rp 100.000\nSaldo akhir Budi tetap aman: Rp 50.000"
  },
  {
    "id": 24,
    "title": "[24/36] Variabel Reaktif & Pengamat (reaktif & pantau)",
    "section": "Fitur Spesial",
    "docLink": "/fitur-khas/reaktif",
    "description": "\n      <p>Pernah pakai spreadsheet seperti Excel? Saat kamu mengubah angka di satu sel, sel rumus total otomatis berubah seketika.</p>\n      <p>Celang memiliki fitur itu langsung di tingkat bahasa tanpa perlu pustaka luar:</p>\n      <ul>\n        <li><strong><code>reaktif total = harga * jumlah</code></strong>: Variabel <code>total</code> otomatis merekam variabel apa saja yang dipakainya dan langsung terhitung ulang saat variabel tersebut berubah.</li>\n        <li><strong><code>pantau variabel (baru, lama) { ... }</code></strong>: Pengamat otomatis yang langsung bertindak seketika setiap kali nilai suatu variabel berubah.</li>\n      </ul>\n    ",
    "task": "Ubah nilai <code>jumlah</code> menjadi <code>5</code> sehingga output notifikasi tercetak <code>Rp 125.000</code>.",
    "code": "buat harga = 25.ribu\nbuat jumlah = 2\n\n// 1. Variabel reaktif: otomatis terhitung ulang saat harga/jumlah berubah\nreaktif total = harga * jumlah\n\n// 2. Pengamat: bereaksi otomatis saat nilai total berubah\npantau total (baru, lama) {\n    cetak(\"Notifikasi: Total belanja berubah jadi\", baru.rupiah)\n}\n\ncetak(\"Total awal:\", total.rupiah)\n\n// Ubah variabel penyusunnya:\njumlah = 4",
    "expected": "Total awal: Rp 50.000\nNotifikasi: Total belanja berubah jadi Rp 125.000"
  },
  {
    "id": 25,
    "title": "[25/36] Saluran Pengiriman Data (Saluran)",
    "section": "Fitur Spesial",
    "docLink": "/fitur-khas/saluran",
    "description": "\n      <p>Bayangkan kamu dan rekanmu bekerja di dua ruangan terpisah. Untuk saling berkirim surat atau barang tanpa harus bolak-balik jalan kaki, kalian memasang sebuah <strong>pipa luncuran tabung</strong> di dinding.</p>\n      <p>Itulah analogi <strong>Saluran</strong> di Celang:</p>\n      <ul>\n        <li><strong><code>Saluran(kapasitas)</code></strong>: Menyiapkan pipa yang bisa menampung antrean barang. Misalnya <code>Saluran(2)</code> berarti pipa tersebut muat 2 barang sekaligus.</li>\n        <li><strong><code>pipa.kirim(\"Dokumen A\")</code></strong>: Memasukkan barang ke dalam pipa dari satu sisi.</li>\n        <li><strong><code>pipa.terima()</code></strong>: Mengambil barang yang meluncur keluar di sisi seberang. Barang yang masuk lebih dulu akan keluar lebih dulu secara teratur dan aman tanpa berebut!</li>\n      </ul>\n    ",
    "task": "Ubah kapasitas pipa menjadi <code>Saluran(3)</code>, kirimkan <code>pipa.kirim(\"Dokumen C\")</code>, lalu tambahkan baris <code>cetak(\"Barang diterima:\", pipa.terima())</code> di akhir kode.",
    "code": "buat pipa = Saluran(2)\n\npipa.kirim(\"Dokumen A\")\npipa.kirim(\"Dokumen B\")\n\ncetak(\"Barang diterima:\", pipa.terima())\ncetak(\"Barang diterima:\", pipa.terima())",
    "expected": "Barang diterima: Dokumen A\nBarang diterima: Dokumen B\nBarang diterima: Dokumen C"
  },
  {
    "id": 26,
    "title": "[26/36] Aliran Pipa & Pekerja Latar (sambil)",
    "section": "Fitur Spesial",
    "docLink": "/fitur-khas/konkurensi",
    "description": "\n      <p>Melanjutkan analogi pipa tadi: bagaimana kalau ada teman di ruangan lain yang terus-menerus memasukkan barang ke pipa, sementara kamu sibuk menerima barangnya satu per satu?</p>\n      <ul>\n        <li><strong><code>sambil { ... }</code></strong>: Menyuruh tugas berjalan di latar belakang (seperti mesin otomatis memasukkan barang ke pipa tanpa mengganggu pekerjaan utamamu).</li>\n        <li><strong><code>untuk barang dalam pipa { ... }</code></strong>: Kamu tinggal duduk santai menyedot setiap barang yang keluar dari pipa menggunakan perulangan biasa.</li>\n        <li><strong><code>pipa.tutup()</code></strong>: Saat temanmu selesai dan menutup pipa, perulanganmu akan otomatis selesai dengan sendirinya!</li>\n      </ul>\n    ",
    "task": "Ubah kapasitas menjadi <code>Saluran(3)</code> dan tambahkan baris <code>pipa.kirim(\"Data 3\")</code> di dalam blok <code>sambil</code> sebelum ditutup.",
    "code": "buat pipa = Saluran(2)\n\nsambil {\n    pipa.kirim(\"Data 1\")\n    pipa.kirim(\"Data 2\")\n    pipa.tutup()\n}\n\nuntuk pesan dalam pipa {\n    cetak(\"Aliran data:\", pesan)\n}",
    "expected": "Aliran data: Data 1\nAliran data: Data 2\nAliran data: Data 3"
  },
  {
    "id": 27,
    "title": "[27/36] Pemrosesan Serentak (bersamaan)",
    "section": "Fitur Spesial",
    "docLink": "/fitur-khas/konkurensi",
    "description": "\n      <p>Jika <code>sambil</code> menjalankan tugas latar belakang yang tidak perlu ditunggu, kata kunci <code>bersamaan</code> digunakan ketika Anda memiliki beberapa tugas yang dikerjakan serentak (<em>parallel</em>), dan program <strong>wajib menunggu semua tugas tersebut rampung</strong> sebelum melangkah ke baris berikutnya.</p>\n      <p>Ini seperti membagi tugas membersihkan ruangan: semua anggota bekerja serentak, dan pintu baru dikunci setelah semuanya selesai.</p>\n    ",
    "task": "Tambahkan baris <code>cetak(\"Tugas Paralel C\")</code> di dalam kurung kurawal blok <code>bersamaan</code>.",
    "code": "bersamaan {\n    cetak(\"Tugas Paralel A\")\n    cetak(\"Tugas Paralel B\")\n}\n\ncetak(\"Semua tugas selesai bersamaan!\")",
    "expected": "Tugas Paralel A\nTugas Paralel B\nTugas Paralel C\nSemua tugas selesai bersamaan!"
  },
  {
    "id": 28,
    "title": "[28/36] Tugas di Belakang Layar (Titip & Jemput)",
    "section": "Fitur Spesial",
    "docLink": "/fitur-khas/titip-jemput",
    "description": "\n      <p>Pernah pesan kopi atau cuci pakaian di laundry? Kasir memberimu selembar <strong>tiket / nota antrean</strong>. Kamu tidak perlu berdiri mematung di depan kompor menunggu kopi selesai; kamu bisa santai duduk atau ngobrol dulu. Pas kopinya sudah jadi dan kamu butuh, tinggal bawa tiket itu untuk <strong>dijemput</strong>.</p>\n      <p>Begitulah cara kerja tugas asinkron di Celang:</p>\n      <ul>\n        <li><strong><code>titip fungsiKu()</code></strong>: Menyuruh Celang menjalankan tugas berat di belakang layar, dan kamu langsung diberi selembar <strong>tiket</strong> tanpa membuat programmu macet.</li>\n        <li><strong><code>jemput tiket</code></strong>: Kapan pun kamu butuh hasilnya, panggil <code>jemput tiket</code>. Program akan mengambil hasil pekerjaan yang sudah selesai tadi.</li>\n      </ul>\n    ",
    "task": "Ubah teks yang dikembalikan fungsi menjadi: <code>\"Laporan Tahunan 2026 Disetujui\"</code>.",
    "code": "fungsi siapkanLaporan() {\n    kembalikan \"Laporan Keuangan Selesai\"\n}\n\n// 1. Titip pekerjaan ke belakang layar, langsung dapat tiket\nbuat tiket = titip siapkanLaporan()\ncetak(\"Tiket di tangan, kita bisa kerjakan yang lain...\")\n\n// 2. Jemput hasilnya saat sudah dibutuhkan\nbuat hasil = jemput tiket\ncetak(\"Hasil jemput:\", hasil)",
    "expected": "Tiket di tangan, kita bisa kerjakan yang lain...\nHasil jemput: Laporan Tahunan 2026 Disetujui"
  },
  {
    "id": 29,
    "title": "[29/36] Keamanan Nilai Kosong (Null-Safety)",
    "section": "Fitur Spesial",
    "docLink": "/fitur-khas/ergonomi-dan-null-safety",
    "description": "\n      <p>Akses properti berjenjang yang ternyata bernilai kosong (<em>null</em>) adalah momok yang sering memicu <em>crash</em>. Celang menyelesaikannya lewat dua kata kunci ergonomis: <code>mungkin</code> (<em>safe navigation</em>) dan <code>kalau kosong</code> (<em>null coalescing</em> / <code>??</code>).</p>\n      <p>Jika ada bagian data yang tidak ditemukan, <code>mungkin</code> akan membatalkan penelusuran secara tenang tanpa menimbulkan galat. Lalu <code>kalau kosong</code> akan memberikan nilai cadangan yang telah Anda siapkan.</p>\n    ",
    "task": "Ubah teks nilai cadangan <code>\"Jakarta\"</code> menjadi <code>\"Bandung\"</code>.",
    "code": "buat pengguna = {\n    \"profil\": {\n        \"nama\": \"Citra\"\n    }\n}\n\nbuat kota = mungkin pengguna.profil.alamat.kota kalau kosong \"Jakarta\"\ncetak(\"Kota pengguna:\", kota)",
    "expected": "Kota pengguna: Bandung"
  },
  {
    "id": 30,
    "title": "[30/36] Pembersihan Otomatis (Defer saat_selesai)",
    "section": "Fitur Spesial",
    "docLink": "/dasar/kontrol-alur",
    "description": "\n      <p>Kata kunci <code>saat_selesai</code> (ekuivalen dengan <code>defer</code> di Go) menjamin bahwa suatu perintah perapian—seperti menutup berkas, melepaskan kunci mutex, atau menutup koneksi jaringan—pasti dieksekusi tepat sebelum fungsi berakhir.</p>\n      <p>Instruksi <code>saat_selesai</code> didaftarkan di awal agar tidak terlupa, namun baru dieksekusi paling akhir dengan urutan terbalik (<em>LIFO</em>).</p>\n    ",
    "task": "Tambahkan baris <code>saat_selesai cetak(\"Pembersihan memori selesai.\")</code> di bawah baris <code>saat_selesai</code> pertama.",
    "code": "fungsi prosesTransaksi() {\n    saat_selesai cetak(\"Pembersihan: sesi koneksi ditutup.\")\n    cetak(\"Langkah 1: Memproses pembayaran...\")\n    cetak(\"Langkah 2: Mencetak resi...\")\n}\n\nprosesTransaksi()",
    "expected": "Langkah 1: Memproses pembayaran...\nLangkah 2: Mencetak resi...\nPembersihan memori selesai.\nPembersihan: sesi koneksi ditutup."
  },
  {
    "id": 31,
    "title": "[31/36] Pipa Penyaluran Data (Pipeline |>)",
    "section": "Fitur Spesial",
    "docLink": "/fitur-khas/pipeline",
    "description": "\n      <p>Menulis fungsi bertingkat yang bersarang (seperti <code>g(f(x))</code>) sering kali sulit dibaca karena harus ditelusuri dari bagian tengah ke luar.</p>\n      <p>Operator pipa <code>|&gt;</code> mengalirkan hasil komputasi dari sebelah kiri langsung sebagai argumen pertama ke fungsi di sebelah kanannya. Logika pemrosesan data menjadi runtut dan alami dari kiri ke kanan.</p>\n    ",
    "task": "Ubah angka masukan awal <code>5</code> menjadi <code>10</code>, lalu amati hasil perhitungannya: (10 * 10) + 10 = 110.",
    "code": "fungsi kuadrat(x) {\n    kembalikan x * x\n}\n\nfungsi tambahSepuluh(x) {\n    kembalikan x + 10\n}\n\nbuat hasil = 5 |> kuadrat() |> tambahSepuluh()\ncetak(\"Hasil pipa:\", hasil)",
    "expected": "Hasil pipa: 110"
  },
  {
    "id": 32,
    "title": "[32/36] Modul Data JSON (json.uraikan & json.rapi)",
    "section": "Fitur Spesial",
    "docLink": "/pustaka-standar/json",
    "description": "\n      <p>Hampir semua aplikasi masa kini bertukar data lewat format JSON (<em>JavaScript Object Notation</em>).</p>\n      <p>Modul bawaan <code>json</code> di Celang sudah siap pakai tanpa perlu instal apa pun:</p>\n      <ul>\n        <li><strong><code>json.rapi(objek)</code></strong>: Mengubah struktur data Celang menjadi teks JSON yang rapi dan mudah dibaca.</li>\n        <li><strong><code>json.uraikan(teksJSON)</code></strong>: Membaca teks JSON dan mengubahnya kembali menjadi objek/kamus data asli Celang.</li>\n        <li><strong><code>json.dapatkan(objek, jalur)</code></strong>: Mengambil data bersarang secara aman tanpa khawatir galat.</li>\n      </ul>\n    ",
    "task": "Ubah nilai poin menjadi <code>2500</code> sehingga data JSON mencetak <code>\"poin\": 2500</code>.",
    "code": "buat dataPengguna = {\n    \"nama\": \"Citra\",\n    \"kota\": \"Bandung\",\n    \"poin\": 1500\n}\n\n// 1. Mengubah objek ke teks JSON rapi:\nbuat teksJSON = json.rapi(dataPengguna)\ncetak(\"Teks JSON:\")\ncetak(teksJSON)\n\n// 2. Mengurai teks JSON kembali ke objek:\nbuat hasilUrai = json.uraikan(teksJSON)\ncetak(\"Nama hasil urai:\", hasilUrai.nama)",
    "expected": "Teks JSON:\n{\n  \"kota\": \"Bandung\",\n  \"nama\": \"Citra\",\n  \"poin\": 2500\n}\nNama hasil urai: Citra"
  },
  {
    "id": 33,
    "title": "[33/36] Satuan Budaya Nusantara & Terbilang Bawaan",
    "section": "Fitur Spesial",
    "docLink": "/fitur-khas/satuan-budaya",
    "description": "\n      <p>Masyarakat Indonesia jarang menyebut uang dengan mengetik angka nol yang panjang (seperti <code>2500000</code>). Kita terbiasa bilang <em>\"dua setengah juta\"</em>, <em>\"lima puluh ribu\"</em>, atau butuh kuitansi resmi bertuliskan <em>\"dua juta lima ratus ribu rupiah\"</em>.</p>\n      <p>Celang memahami kebiasaan ini langsung di tingkat sintaks tanpa perlu memasang pustaka luar:</p>\n      <ul>\n        <li><strong>Singkatan Bilangan</strong>: Tulis langsung <code>50.ribu</code>, <code>2.5.juta</code>, <code>1.miliar</code>, atau <code>10.triliun</code>. Aman dari salah hitung angka nol!</li>\n        <li><strong>Format Rupiah (<code>.rupiah</code>)</strong>: Otomatis merapikan angka dengan titik ribuan khas Indonesia (contoh: <code>Rp 2.500.000</code>).</li>\n        <li><strong>Ejaan Kuitansi Resmi (<code>.terbilang</code>)</strong>: Mengubah angka menjadi kalimat ejaan lengkap berakhiran rupiah untuk faktur/nota (contoh: <code>lima puluh ribu rupiah</code>).</li>\n        <li><strong>Ejaan Angka Murni (<code>.kata</code> / <code>.ejaan</code>)</strong>: Mengeja angka apa adanya tanpa kata rupiah (contoh: <code>105.kata</code> jadi <em>\"seratus lima\"</em>).</li>\n        <li><strong>Satuan Waktu &amp; Ukuran Data</strong>: Gunakan <code>.detik</code>, <code>.menit</code>, <code>.jam</code>, <code>.hari</code>, serta <code>.kb</code>, <code>.mb</code>, <code>.gb</code>, <code>.tb</code>. Celang otomatis mengubahnya ke hitungan milidetik dan byte!</li>\n      </ul>\n    ",
    "task": "Tambahkan baris <code>buat saldo = 750.ribu</code> dan <code>cetak(\"Saldo:\", saldo.rupiah)</code> di akhir kode.",
    "code": "buat modal = 50.ribu\nbuat investasi = 2.5.juta\nbuat waktuTunggu = 5.detik\nbuat kuota = 16.gb\n\ncetak(\"Rupiah:\", investasi.rupiah)\ncetak(\"Terbilang kuitansi:\", modal.terbilang)\ncetak(\"Ejaan angka:\", 105.kata)\ncetak(\"Waktu tunggu (ms):\", waktuTunggu)\ncetak(\"Ukuran kuota (byte):\", kuota)",
    "expected": "Rupiah: Rp 2.500.000\nTerbilang kuitansi: lima puluh ribu rupiah\nEjaan angka: seratus lima\nWaktu tunggu (ms): 5000\nUkuran kuota (byte): 17179869184\nSaldo: Rp 750.000"
  },
  {
    "id": 34,
    "title": "[34/36] Jembatan Pustaka C & Ekosistem (ffi)",
    "section": "Fitur Spesial",
    "docLink": "/pustaka-standar/ffi",
    "description": "\n      <p>Dunia pemrograman C/C++ menyimpan jutaan baris pustaka siap pakai yang sangat hebat (seperti <strong>SQLite</strong>, <strong>Raylib 3D</strong>, <strong>OpenGL</strong>, atau <strong>Win32 API</strong>). Mengapa harus menulis ulang semuanya dari nol?</p>\n      <p>Modul bawaan <strong><code>ffi</code></strong> (<em>Foreign Function Interface</em>) adalah jembatan emas Celang:</p>\n      <ul>\n        <li><strong>Nol Kompiler C</strong>: Berkat arsitektur <em>Dynamic ABI Invocation</em>, Celang memanggil biner <code>.so</code> / <code>.dll</code> langsung tanpa butuh GCC, Clang, atau CGO!</li>\n        <li><strong>Alokasi &amp; Pointer Memori</strong>: Sediakan buffer memori mentah menggunakan <code>ffi.alloc(ukuran)</code> dan <code>ffi.free(ptr)</code>.</li>\n        <li><strong>Bekerja Lintas Platform</strong>: Biner yang dibangun dengan <code>./celang build</code> bebas kompilasi silang ke Windows maupun Linux tanpa pusing Makefile C.</li>\n      </ul>\n    ",
    "task": "Ubah teks pesan yang ditulis ke memori menjadi <code>\"Halo Pustaka C!\"</code>.",
    "code": "impor ffi\n\n// 1. Pesan buffer memori mentah (64 byte):\nbuat buffer = ffi.alloc(64)\n\n// 2. Tulis teks langsung ke alamat memori:\nffi.write_string(buffer, \"Data tersimpan di memori C\")\n\n// 3. Baca kembali teks dari penunjuk memori:\ncetak(\"Pointer valid:\", !ffi.is_null(buffer))\ncetak(\"Isi memori:\", ffi.read_string(buffer))\n\n// 4. Selalu bersihkan memori setelah selesai:\nffi.free(buffer)\ncetak(\"Pembersihan memori selesai!\")",
    "expected": "Pointer valid: benar\nIsi memori: Halo Pustaka C!\nPembersihan memori selesai!"
  },
  {
    "id": 35,
    "title": "[35/36] Pengujian Unit Otomatis (celang uji)",
    "section": "Fitur Spesial",
    "docLink": "/panduan/pengujian",
    "description": "\n      <p>Kode yang dapat diandalkan adalah kode yang teruji. Celang menyertakan <strong>Test Runner</strong> dan <strong>Pustaka Asersi</strong> bawaan langsung di dalam mesin bahasa!</p>\n      <ul>\n        <li><strong>Blok Uji Mandiri</strong>: Tulis tes secara deklaratif dengan <code>uji \"nama tes\" { ... }</code>. Setiap tes terisolasi dalam cakupannya sendiri.</li>\n        <li><strong>Pustaka Asersi Lengkap</strong>: Gunakan <code>harus_sama(aktual, target)</code>, <code>harus_beda(...)</code>, <code>harus_benar(...)</code>, <code>harus_berisi(...)</code>, hingga <code>harus_galat(...)</code>.</li>\n        <li><strong>Test Runner Terminal</strong>: Di proyek Anda, cukup jalankan perintah <code>celang uji</code> untuk menemukan dan mengeksekusi semua berkas tes dalam hitungan milidetik.</li>\n      </ul>\n    ",
    "task": "Ubah ekspektasi <code>harus_sama(hasil, 30)</code> menjadi <code>harus_sama(hasil, 50)</code> agar pengujian lulus.",
    "code": "fungsi tambah(a, b) {\n  kembalikan a + b\n}\n\nuji \"penjumlahan bilangan positif\" {\n  buat hasil = tambah(20, 30)\n  harus_sama(hasil, 50, \"20 + 30 harus bernilai 50\")\n  harus_beda(hasil, 0)\n  harus_benar(hasil > 40)\n}\n\ncetak(\"✓ Pengujian unit selesai dieksekusi!\")",
    "expected": "✓ Pengujian unit selesai dieksekusi!"
  },
  {
    "id": 36,
    "title": "[36/36] Mode Universal (\"use universal\") & Gaya Sintaks",
    "section": "Fitur Spesial",
    "docLink": "/fitur-khas/mode-universal",
    "description": "\n      <p>Kamu mungkin bertanya: <em>\"Kalau kodingnya pakai Bahasa Indonesia, gimana kalau nanti kodenya mau dibagikan ke programmer luar negeri di GitHub?\"</em></p>\n      <p>Tenang! Celang punya fitur keren bernama <strong>Mode Universal</strong>:</p>\n      <ul>\n        <li><strong>Direktif <code>\"use universal\"</code></strong>: Cukup taruh baris <code>\"use universal\"</code> di baris paling atas berkasmu. Celang seketika berubah menjadi bahasa standar internasional (English).</li>\n        <li><strong>Aturan Standar Global</strong>: Di mode ini, kata kunci bahasa Indonesia (seperti <code>buat</code>, <code>jika</code>, <code>cetak</code>) dilarang dan otomatis diarahkan memakai bahasa Inggris standar (<code>let</code>, <code>if</code>, <code>print</code>) serta satuan global (<code>.k</code>, <code>.million</code>, <code>.usd</code>).</li>\n        <li><strong>Gaya di Tengah Kode (<code>gaya \"...\" { }</code>)</strong>: Di file biasa, kamu bahkan bisa meminjam gaya bahasa lain di tengah-tengah program, misalnya membungkus cuplikan kode dengan <code>gaya \"javascript\" { ... }</code> atau <code>gaya \"pythonic\" { ... }</code>.</li>\n        <li><strong>Alat Format Otomatis</strong>: Kamu juga bisa mengubah gaya satu file penuh dalam sekejap dengan perintah terminal: <code>./celang format program.ce --style=universal</code>.</li>\n      </ul>\n      <p>🎉 <strong>Langkah Selanjutnya</strong>: Pasang Celang di komputermu, dan ubah kreasimu menjadi aplikasi biner mandiri (<code>.exe</code> di Windows atau biner di Linux/Mac) yang bisa langsung diklik dua kali oleh temanmu tanpa perlu menginstal Celang atau Go dengan perintah: <code>./celang bangun program.ce -o Game.exe</code>!</p>\n    ",
    "task": "Ubah nilai <code>bonus</code> menjadi <code>20.k</code> agar total saldo menjadi $120,000.00.",
    "code": "\"use universal\"\n\nlet balance = 100.k\nconst bonus = 5.k\nlet total = balance + bonus\n\nprint(\"Total USD:\", total.usd)",
    "expected": "Total USD: $120,000.00"
  }
];

const SLIDES_EN = [
  {
    "id": 1,
    "title": "[1/36] Introduction to Celang Tour",
    "section": "Introduction",
    "docLink": "/panduan/pengenalan",
    "description": "\n      <p>Welcome to <strong>Celang Tour</strong> (<em>A Tour of Celang</em>)! Learn the syntax, concepts, and unique superpowers of Celang interactively step by step.</p>\n      <p>The screen is divided into two panels:</p>\n      <ul>\n        <li><strong>Left Panel:</strong> Concise lessons, simple explanations, and a practical challenge.</li>\n        <li><strong>Right Panel:</strong> Interactive code editor where you can test and run Celang code instantly.</li>\n      </ul>\n      <p>Code executes securely right inside your browser powered by <strong>Celang WebAssembly</strong>.</p>\n    ",
    "task": "Click the <strong>Run Code ▶</strong> button in the top-right header to test your very first Celang execution.",
    "code": "buat salam = \"Selamat Datang di Tur Celang!\"\ncetak(salam)\ncetak(\"Silakan klik tombol Selanjutnya di bawah untuk memulai belajar.\")",
    "expected": "Selamat Datang di Tur Celang!\nSilakan klik tombol Selanjutnya di bawah untuk memulai belajar."
  },
  {
    "id": 2,
    "title": "[2/36] Choose Tutorial Language",
    "section": "Introduction",
    "docLink": "/panduan/pengenalan",
    "description": "\n      <p>This tour supports two tutorial explanation languages:</p>\n      <div style=\"display: grid; grid-template-columns: 1fr 1fr; gap: 12px; margin: 16px 0;\">\n        <div class=\"lang-select-card\" id=\"card-lang-id\" style=\"padding: 14px; border: 1px solid var(--border-color); border-radius: 8px; cursor: pointer; text-align: center;\">\n          <div style=\"font-size: 24px; margin-bottom: 4px;\">🇮🇩</div>\n          <strong>Bahasa Indonesia</strong>\n          <p style=\"font-size: 12px; margin: 4px 0 0; opacity: 0.8;\">Penjelasan dalam bahasa Indonesia sehari-hari.</p>\n        </div>\n        <div class=\"lang-select-card active\" id=\"card-lang-en\" style=\"padding: 14px; border: 2px solid var(--accent-primary); border-radius: 8px; cursor: pointer; text-align: center; background: rgba(59, 130, 246, 0.08);\">\n          <div style=\"font-size: 24px; margin-bottom: 4px;\">🇬🇧</div>\n          <strong>English</strong>\n          <p style=\"font-size: 12px; margin: 4px 0 0; opacity: 0.8;\">Explanations in international English.</p>\n        </div>\n      </div>\n      <p>You can also toggle between languages anytime using the flag button in the top navigation bar.</p>\n    ",
    "task": "Select your language above, or click <strong>Next ▶</strong> to continue the tour in English.",
    "code": "buat bahasa = \"English\"\ncetak(\"Selected guide language:\", bahasa)\ncetak(\"Enjoy your journey with Celang!\")",
    "expected": "Selected guide language: English\nEnjoy your journey with Celang!"
  },
  {
    "id": 3,
    "title": "[3/36] Hello World",
    "section": "Language Basics",
    "docLink": "/panduan/mulai-cepat",
    "description": "\n      <p>Let's start with the basics. The built-in <code>cetak(...)</code> function prints text and variables to the terminal output.</p>\n      <p>Strings are surrounded by double quotes (<code>\"...\"</code>). You can print multiple comma-separated items in a single call.</p>\n    ",
    "task": "Change <code>\"Halo, Dunia!\"</code> to <code>\"Halo, Celang!\"</code> and observe the terminal output.",
    "code": "cetak(\"Halo, Dunia!\")",
    "expected": "Halo, Celang!"
  },
  {
    "id": 4,
    "title": "[4/36] Variables & Values",
    "section": "Language Basics",
    "docLink": "/dasar/variabel-tipe",
    "description": "\n      <p>Variables in Celang are declared using the <code>buat</code> keyword. They are dynamically typed and can be updated anytime without repeating <code>buat</code>.</p>\n      <p>Multiple arguments passed to <code>cetak(...)</code> are automatically formatted with separating spaces.</p>\n    ",
    "task": "Change <code>umur</code> to <code>25</code> and run the code.",
    "code": "buat nama = \"Budi\"\nbuat umur = 20\n\numur = 21\n\ncetak(nama, \"berumur\", umur, \"tahun\")",
    "expected": "Budi berumur 25 tahun"
  },
  {
    "id": 5,
    "title": "[5/36] Constants (konstanta)",
    "section": "Language Basics",
    "docLink": "/dasar/variabel-tipe",
    "description": "\n      <p>Use the <code>konstanta</code> keyword to declare immutable values that must not change during execution.</p>\n      <p>Attempting to reassign a constant raises an immediate error. Use constants for fixed references like math constants or configuration keys.</p>\n    ",
    "task": "Add <code>konstanta VERSI = 1</code> and <code>cetak(\"Versi:\", VERSI)</code> at the end of the code.",
    "code": "konstanta PI = 3.14159\nkonstanta NAMA_APLIKASI = \"Sistem Keuangan\"\n\ncetak(\"Aplikasi:\", NAMA_APLIKASI)\ncetak(\"Nilai PI:\", PI)",
    "expected": "Aplikasi: Sistem Keuangan\nNilai PI: 3.14159\nVersi: 1"
  },
  {
    "id": 6,
    "title": "[6/36] Data Types & Example Values",
    "section": "Language Basics",
    "docLink": "/dasar/variabel-tipe",
    "description": "\n      <p>Celang includes native data types designed for real-world programming:</p>\n      <table class=\"table-types\" style=\"width:100%; font-size:13px; border-collapse:collapse; margin-bottom:12px;\">\n        <thead>\n          <tr style=\"background:rgba(255,255,255,0.05); text-align:left;\">\n            <th style=\"padding:6px 8px; border:1px solid var(--border-color);\">Type</th>\n            <th style=\"padding:6px 8px; border:1px solid var(--border-color);\">Usage</th>\n            <th style=\"padding:6px 8px; border:1px solid var(--border-color);\">Example</th>\n          </tr>\n        </thead>\n        <tbody>\n          <tr><td style=\"padding:6px 8px; border:1px solid var(--border-color);\"><code>Angka</code></td><td style=\"padding:6px 8px; border:1px solid var(--border-color);\">Integers and floats</td><td style=\"padding:6px 8px; border:1px solid var(--border-color);\"><code>42</code>, <code>3.14</code></td></tr>\n          <tr><td style=\"padding:6px 8px; border:1px solid var(--border-color);\"><code>Teks</code></td><td style=\"padding:6px 8px; border:1px solid var(--border-color);\">Strings</td><td style=\"padding:6px 8px; border:1px solid var(--border-color);\"><code>\"Halo\"</code>, <code>'Budi'</code></td></tr>\n          <tr><td style=\"padding:6px 8px; border:1px solid var(--border-color);\"><code>Kebenaran</code></td><td style=\"padding:6px 8px; border:1px solid var(--border-color);\">Booleans</td><td style=\"padding:6px 8px; border:1px solid var(--border-color);\"><code>benar</code>, <code>salah</code></td></tr>\n          <tr><td style=\"padding:6px 8px; border:1px solid var(--border-color);\"><code>Larik</code></td><td style=\"padding:6px 8px; border:1px solid var(--border-color);\">Ordered lists (Arrays)</td><td style=\"padding:6px 8px; border:1px solid var(--border-color);\"><code>[1, 2, 3]</code></td></tr>\n          <tr><td style=\"padding:6px 8px; border:1px solid var(--border-color);\"><code>Kamus</code></td><td style=\"padding:6px 8px; border:1px solid var(--border-color);\">Key-value maps</td><td style=\"padding:6px 8px; border:1px solid var(--border-color);\"><code>{\"a\": 1}</code></td></tr>\n          <tr><td style=\"padding:6px 8px; border:1px solid var(--border-color);\"><code>Kosong</code></td><td style=\"padding:6px 8px; border:1px solid var(--border-color);\">Absence of value (null)</td><td style=\"padding:6px 8px; border:1px solid var(--border-color);\"><code>kosong</code></td></tr>\n        </tbody>\n      </table>\n      <p>Values of each type can be directly printed using <code>cetak(...)</code>.</p>\n    ",
    "task": "Change <code>sedangAktif</code> to <code>salah</code> and observe the result.",
    "code": "buat umur = 28\nbuat tinggi = 175.5\nbuat pesan = \"Semangat Belajar!\"\nbuat sedangAktif = benar\nbuat statusLogin = kosong\n\ncetak(\"Umur (Angka):\", umur)\ncetak(\"Tinggi (Angka):\", tinggi)\ncetak(\"Pesan (Teks):\", pesan)\ncetak(\"Sedang aktif (Kebenaran):\", sedangAktif)\ncetak(\"Status login (Kosong):\", statusLogin)",
    "expected": "Umur (Angka): 28\nTinggi (Angka): 175.5\nPesan (Teks): Semangat Belajar!\nSedang aktif (Kebenaran): false\nStatus login (Kosong): null"
  },
  {
    "id": 7,
    "title": "[7/36] Text Interpolation (f-strings)",
    "section": "Language Basics",
    "docLink": "/fitur-khas/interpolasi-teks",
    "description": "\n      <p>Instead of manually concatenating with <code>+</code>, use Python-like <strong>f-strings</strong> for clean string formatting:</p>\n      <p>Prefix strings with <code>f</code> and embed expressions directly inside <code>${...}</code>.</p>\n      <p>Example: <code>f\"Hello ${name}, you are ${age + 1} years old\"</code>.</p>\n    ",
    "task": "Change <code>tabungan</code> to <code>50000</code> so the second line outputs: <code>Jika ditambah 10000: Rp60000</code>.",
    "code": "buat nama = \"Andi\"\nbuat tabungan = 25000\n\ncetak(f\"Nasabah: ${nama}\")\ncetak(f\"Jika ditambah 10000: Rp${tabungan + 10000}\")",
    "expected": "Nasabah: Andi\nJika ditambah 10000: Rp60000"
  },
  {
    "id": 8,
    "title": "[8/36] Conditional Branching (jika)",
    "section": "Language Basics",
    "docLink": "/dasar/kontrol-alur",
    "description": "\n      <p>Control program flow using <code>jika</code> (if), <code>selain jika</code> (else if), and <code>selain</code> (else).</p>\n      <p>Parentheses are optional around conditions, but curly braces <code>{ ... }</code> are required.</p>\n    ",
    "task": "Change <code>nilaiUjian</code> to <code>65</code> so the printed status becomes <code>Status: Perlu Remedial</code>.",
    "code": "buat nilaiUjian = 85\n\njika nilaiUjian >= 80 {\n    cetak(\"Status: Lulus Memuaskan\")\n} selain jika nilaiUjian >= 60 {\n    cetak(\"Status: Perlu Remedial\")\n} selain {\n    cetak(\"Status: Mengulang Ujian\")\n}",
    "expected": "Status: Perlu Remedial"
  },
  {
    "id": 9,
    "title": "[9/36] Value Selection (pilih ... bila)",
    "section": "Language Basics",
    "docLink": "/dasar/kontrol-alur",
    "description": "\n      <p>The <code>pilih</code> statement matches values against conditions using <code>bila</code> branches.</p>\n      <p>The <code>lainnya:</code> fallback branch executes when no cases match. Unlike C-like switch statements, Celang does not require manual <code>break</code> keywords.</p>\n    ",
    "task": "Change <code>peran</code> to <code>\"editor\"</code> so that the printed output becomes <code>Akses: Mengelola Konten</code>.",
    "code": "buat peran = \"tamu\"\n\npilih peran {\n    bila \"admin\":\n        cetak(\"Akses: Penuh ke Sistem\")\n    bila \"editor\":\n        cetak(\"Akses: Mengelola Konten\")\n    bila \"penulis\":\n        cetak(\"Akses: Menulis Draf\")\n    lainnya:\n        cetak(\"Akses: Hanya Membaca\")\n}",
    "expected": "Akses: Mengelola Konten"
  },
  {
    "id": 10,
    "title": "[10/36] Conditional Loop (selama)",
    "section": "Language Basics",
    "docLink": "/dasar/perulangan",
    "description": "\n      <p>Use <code>selama</code> (while loop) to execute a block of statements repeatedly as long as the test condition remains <code>benar</code>.</p>\n      <p>Ensure loop state updates inside the block so the loop terminates gracefully.</p>\n    ",
    "task": "Change <code>hitung &lt;= 3</code> to <code>hitung &lt;= 5</code> to count up to 5.",
    "code": "buat hitung = 1\n\nselama hitung <= 3 {\n    cetak(\"Hitungan ke:\", hitung)\n    hitung = hitung + 1\n}\n\ncetak(\"Selesai!\")",
    "expected": "Hitungan ke: 1\nHitungan ke: 2\nHitungan ke: 3\nHitungan ke: 4\nHitungan ke: 5\nSelesai!"
  },
  {
    "id": 11,
    "title": "[11/36] Loop Collections & Ranges (untuk ... dalam)",
    "section": "Language Basics",
    "docLink": "/dasar/perulangan",
    "description": "\n      <p>The <code>untuk ... dalam</code> loop is the cleanest way to iterate over ranges or collections:</p>\n      <ul>\n        <li><strong>Number Ranges</strong>: Use double-dots like <code>1..5</code> to produce an inclusive sequence from 1 to 5.</li>\n        <li><strong>Array Collections</strong>: Use <code>untuk item dalam daftar</code> to process every item sequentially from beginning to end.</li>\n      </ul>\n    ",
    "task": "Add <code>\"Pisang\"</code> to the fruits array: <code>[\"Apel\", \"Jeruk\", \"Mangga\", \"Pisang\"]</code>.",
    "code": "// 1. Iterate over range 1 to 3:\nuntuk angka dalam 1..3 {\n    cetak(\"Angka:\", angka)\n}\n\n// 2. Iterate over fruits list:\nbuat buah = [\"Apel\", \"Jeruk\", \"Mangga\"]\nuntuk item dalam buah {\n    cetak(\"Buah:\", item)\n}",
    "expected": "Angka: 1\nAngka: 2\nAngka: 3\nBuah: Apel\nBuah: Jeruk\nBuah: Mangga\nBuah: Pisang"
  },
  {
    "id": 12,
    "title": "[12/36] Loop Control (lanjut & berhenti)",
    "section": "Language Basics",
    "docLink": "/dasar/perulangan",
    "description": "\n      <p>While a loop is running, you may want to skip a specific turn or stop completely before the loop finishes:</p>\n      <ul>\n        <li><strong><code>lanjut</code> (Skip / Continue)</strong>: Like saying <em>\"skip this turn and jump right to the next number!\"</em>. Any statements below it are skipped for that round.</li>\n        <li><strong><code>berhenti</code> (Emergency Brake / Break)</strong>: Like saying <em>\"that's enough, stop right here!\"</em>. The loop immediately terminates and execution proceeds past the loop.</li>\n      </ul>\n    ",
    "task": "Remove the <code>jika i == 2 { lanjut }</code> block so number 2 is not skipped.",
    "code": "untuk i dalam 1..5 {\n    jika i == 2 {\n        lanjut\n    }\n    jika i == 4 {\n        berhenti\n    }\n    cetak(\"Nilai i:\", i)\n}",
    "expected": "Nilai i: 1\nNilai i: 2\nNilai i: 3"
  },
  {
    "id": 13,
    "title": "[13/36] Collections: Arrays",
    "section": "Language Basics",
    "docLink": "/fitur-khas/koleksi-array-map",
    "description": "\n      <p>An array is an ordered list enclosed in square brackets <code>[...]</code>. Indices are zero-based, starting at <code>0</code>.</p>\n      <p>Celang provides convenient methods such as <code>.push()</code> to append items and <code>panjang(...)</code> to retrieve the length.</p>\n    ",
    "task": "Add <code>cetak(\"Elemen kedua:\", angka[1])</code> below the first element print line.",
    "code": "buat angka = [10, 20, 30]\n\nangka.push(40)\ncetak(\"Elemen pertama:\", angka[0])\ncetak(\"Jumlah elemen:\", panjang(angka))\ncetak(\"Seluruh data:\", angka)",
    "expected": "Elemen pertama: 10\nElemen kedua: 20\nJumlah elemen: 4\nSeluruh data: [10, 20, 30, 40]"
  },
  {
    "id": 14,
    "title": "[14/36] Key-Value Maps",
    "section": "Language Basics",
    "docLink": "/fitur-khas/koleksi-array-map",
    "description": "\n      <p>Maps store key-value pairs wrapped in curly braces <code>{ ... }</code>.</p>\n      <p>You can access map properties using either bracket notation <code>peta[\"kunci\"]</code> or dot notation <code>peta.kunci</code>.</p>\n    ",
    "task": "Add a line <code>cetak(\"Usia:\", pengguna.usia)</code> at the end of the code.",
    "code": "buat pengguna = {\n    \"nama\": \"Andi\",\n    \"pekerjaan\": \"Pengembang\",\n    \"usia\": 25\n}\n\ncetak(\"Nama:\", pengguna[\"nama\"])\ncetak(\"Pekerjaan:\", pengguna.pekerjaan)",
    "expected": "Nama: Andi\nPekerjaan: Pengembang\nUsia: 25"
  },
  {
    "id": 15,
    "title": "[15/36] Unpacking Data (Destructuring)",
    "section": "Language Basics",
    "docLink": "/fitur-khas/destrukturisasi",
    "description": "\n      <p>Imagine coming home with a bag of groceries. Instead of taking items out one by one, you unpack everything into separate bowls all in one quick motion.</p>\n      <p>Celang's <strong>Destructuring</strong> feature lets you unpack Arrays or Maps into separate variables in a single line of code:</p>\n      <ul>\n        <li><strong>Unpacking Arrays by Position</strong>: <code>buat [x, y] = koordinat</code> (index 0 goes to <code>x</code>, index 1 goes to <code>y</code>).</li>\n        <li><strong>Unpacking Maps by Key Name</strong>: <code>buat { nama, kota } = profil</code> (extracts <code>\"nama\"</code> and <code>\"kota\"</code> directly).</li>\n      </ul>\n    ",
    "task": "Change coordinates to <code>[500, 750]</code> so the output prints <code>Koordinat: 500 750</code>.",
    "code": "buat koordinat = [100, 250]\nbuat profil = {\n    \"nama\": \"Citra\",\n    \"kota\": \"Surabaya\"\n}\n\n// 1. Unpack array by position:\nbuat [x, y] = koordinat\n\n// 2. Unpack map by key name:\nbuat { nama, kota } = profil\n\ncetak(\"Koordinat:\", x, y)\ncetak(\"Profil:\", nama, \"dari\", kota)",
    "expected": "Koordinat: 500 750\nProfil: Citra dari Surabaya"
  },
  {
    "id": 16,
    "title": "[16/36] Functions, Parameter Constraints, & Return Types",
    "section": "Language Basics",
    "docLink": "/dasar/fungsi",
    "description": "\n      <p>A function is like a recipe or machine: you pass in inputs (parameters), it does the work, and hands you back a result.</p>\n      <p>In Celang, you can write functions flexibly or with strict type guarantees:</p>\n      <ul>\n        <li><strong>Parameter Type Constraints</strong>: Enforce argument types using <code>itu</code> (e.g. <code>panjang itu Angka</code>) or concise symbol syntax <code>panjang: Angka</code>.</li>\n        <li><strong>Return Type Hints</strong>: Specify output types using <code>mengembalikan Angka</code> or the arrow symbol <code>-&gt; Angka</code>.</li>\n        <li><strong>Returning Results</strong>: Use the <code>kembalikan &lt;value&gt;</code> keyword inside the function.</li>\n      </ul>\n    ",
    "task": "Call the function with length <code>10</code> and width <code>3</code>: <code>buat hasil = hitungLuas(10, 3)</code>.",
    "code": "// Parameters constrained to Angka, returning Angka:\nfungsi hitungLuas(panjang: Angka, lebar: Angka) -> Angka {\n    kembalikan panjang * lebar\n}\n\nbuat hasil = hitungLuas(5, 4)\ncetak(\"Luas persegi:\", hasil)",
    "expected": "Luas persegi: 30"
  },
  {
    "id": 17,
    "title": "[17/36] Arrow Functions (=>) & Data Helpers",
    "section": "Language Basics",
    "docLink": "/fitur-khas/fungsi-ringkas",
    "description": "\n      <p>Sometimes you only need a quick one-liner formula without typing <code>fungsi</code>, braces, and <code>kembalikan</code>.</p>\n      <p>Celang provides concise <strong>Arrow Functions (<code>=&gt;</code>)</strong>:</p>\n      <ul>\n        <li><strong>Single Input</strong>: <code>buat kaliDua = x =&gt; x * 2</code> (directly returns the expression on the right).</li>\n        <li><strong>Multiple Inputs</strong>: <code>buat jumlah = (a, b) =&gt; a + b</code>.</li>\n        <li><strong>Collection Helpers</strong>: Combined with built-in functions like <code>saring(daftar, fn)</code> (filter) and <code>petakan(daftar, fn)</code> (map), data transformations stay brief and elegant!</li>\n      </ul>\n    ",
    "task": "Change <code>kaliDua</code> to <code>x =&gt; x * 10</code> so even numbers multiplied by ten become <code>[20, 40, 60]</code>.",
    "code": "buat kaliDua = x => x * 2\nbuat jumlah = (a, b) => a + b\n\nbuat daftar = [1, 2, 3, 4, 5, 6]\n\n// Filter even numbers, then multiply:\nbuat genap = saring(daftar, x => x % 2 == 0)\nbuat hasil = petakan(genap, kaliDua)\n\ncetak(\"Hasil olah data:\", hasil)\ncetak(\"Total 10 + 20:\", jumlah(10, 20))",
    "expected": "Hasil olah data: [20, 40, 60]\nTotal 10 + 20: 30"
  },
  {
    "id": 18,
    "title": "[18/36] References & Aliases (rujuk / ref)",
    "section": "Language Basics",
    "docLink": "/dasar/rujuk",
    "description": "\n      <p>By default, when passing numbers or values to a new variable or function, Celang creates an independent copy (<em>pass-by-value</em>). Changing the copy never affects the original data.</p>\n      <p>What if you want two variable names to control the exact same memory box (like <em>two remote controls for the same TV</em>)? That is what <strong><code>rujuk</code></strong> (international alias: <code>ref</code>) does:</p>\n      <ul>\n        <li><strong>Variable Aliases (<code>rujuk alias = target</code>)</strong>: Creates a second name for an existing variable. Modifying <code>alias</code> mutates <code>target</code> instantly!</li>\n        <li><strong>Function Parameters (<code>rujuk param</code>)</strong>: Allows functions to modify caller variables directly (<em>pass-by-reference</em>) without tricky pointer syntax.</li>\n      </ul>\n    ",
    "task": "Change the bonus value to <code>50000</code> so the final savings balance prints <code>150000</code>.",
    "code": "// 1. Variable reference alias:\nbuat skor = 50\nrujuk skorPemain = skor\nskorPemain = 95\n\n// 2. Reference parameter directly mutating caller:\nfungsi beriBonus(rujuk saldo, bonus) {\n    saldo = saldo + bonus\n}\n\nbuat tabungan = 100000\nberiBonus(tabungan, 25000)\n\ncetak(\"Skor asli:\", skor)\ncetak(\"Tabungan akhir:\", tabungan)",
    "expected": "Skor asli: 95\nTabungan akhir: 150000"
  },
  {
    "id": 19,
    "title": "[19/36] Custom Data Blueprints (bentuk / struct)",
    "section": "Language Basics",
    "docLink": "/fitur-khas/bentuk-struct",
    "description": "\n      <p>When managing real-world entities like bank accounts or store items, keeping separate variables for each field quickly gets messy.</p>\n      <p>Celang lets you define structured <strong>data blueprints</strong> using the <code>bentuk</code> (struct) keyword:</p>\n      <ul>\n        <li><strong>Fields &amp; Methods</strong>: Define properties (like <code>pemilik</code> and <code>saldo</code>) and member functions (like <code>setor</code>) right inside the blueprint.</li>\n        <li><strong>The <code>@</code> Prefix</strong>: Inside methods, use <code>@</code> to refer to properties/methods of the instance itself (e.g. <code>@saldo += jumlah</code>). <code>@</code> must always be followed by a property name.</li>\n        <li><strong>Anti-Ambiguity (When to use <code>diri</code>?)</strong>: If you need to pass the instance itself to another function or return it (like <code>this</code> or <code>self</code>), you <strong>must use the keyword <code>diri</code></strong>, never standalone <code>@</code>! For example: <code>lapor(diri, @pemilik)</code>, <strong>not</strong> <code>lapor(@, @pemilik)</code>.</li>\n        <li><strong>Instantiating</strong>: Create an instance by calling the blueprint name: <code>buat akun = Rekening(\"Budi\", 50000)</code>.</li>\n      </ul>\n    ",
    "task": "Change the deposit amount to <code>akun.setor(50000)</code> so that the final balance is 100000.",
    "code": "fungsi lapor(obj, nama) {\n    cetak(\"Laporan nasabah:\", nama)\n}\n\nbentuk Rekening {\n    pemilik: Teks,\n    saldo: Angka = 0,\n\n    fungsi setor(jumlah: Angka) {\n        @saldo += jumlah\n        kembalikan @saldo\n    }\n\n    fungsi cek() {\n        // Anti-ambiguity: must use 'diri', not '@'!\n        lapor(diri, @pemilik)\n    }\n}\n\nbuat akun = Rekening(\"Budi\", 50000)\nakun.cek()\nakun.setor(25000)\n\ncetak(\"Pemilik:\", akun.pemilik)\ncetak(\"Saldo akhir:\", akun.saldo)",
    "expected": "Laporan nasabah: Budi\nPemilik: Budi\nSaldo akhir: 100000"
  },
  {
    "id": 20,
    "title": "[20/36] Flexible Containers & Type Bounds (Generics)",
    "section": "Language Basics",
    "docLink": "/fitur-khas/generik",
    "description": "\n      <p>Imagine creating a storage box. Instead of making separate blueprints for numbers and strings, you can define a single versatile blueprint: <code>bentuk Kotak&lt;T&gt;</code>!</p>\n      <ul>\n        <li><strong>Generic Type <code>&lt;T&gt;</code></strong>: The actual type is specified when instantiating (e.g. <code>Kotak&lt;Angka&gt;(100)</code> or <code>Kotak&lt;Teks&gt;(\"Message\")</code>).</li>\n        <li><strong>Type Constraint <code>harus</code></strong>: Enforce allowed types naturally using the <code>harus</code> keyword (e.g. <code>&lt;Value harus Angka&gt;</code>).</li>\n      </ul>\n    ",
    "task": "Change the number inside <code>kotakAngka</code> to <code>500</code> so it prints <code>Skor Game: 500</code>.",
    "code": "bentuk Kotak<T> {\n    label: Teks,\n    isi: T\n}\n\nbuat kotakAngka = Kotak<Angka>(\"Skor Game\", 100)\nbuat kotakTeks = Kotak<Teks>(\"Pesan\", \"Halo Dunia\")\n\ncetak(kotakAngka.label + \":\", kotakAngka.isi)\ncetak(kotakTeks.label + \":\", kotakTeks.isi)",
    "expected": "Skor Game: 500\nPesan: Halo Dunia"
  },
  {
    "id": 21,
    "title": "[21/36] Error Handling (coba ... tangkap)",
    "section": "Language Basics",
    "docLink": "/dasar/penanganan-galat",
    "description": "\n      <p>To prevent programs from crashing on runtime exceptions, use <code>coba</code> (try), <code>tangkap</code> (catch), and <code>akhirnya</code> (finally) blocks.</p>\n      <p>You can raise errors manually using <code>lempar</code> (throw). The <code>akhirnya</code> block is always guaranteed to execute.</p>\n    ",
    "task": "Change the error message thrown to <code>\"Angka pembagi tidak valid!\"</code>.",
    "code": "coba {\n    buat pembagi = 0\n    jika pembagi == 0 {\n        lempar \"Tidak bisa membagi dengan nol!\"\n    }\n} tangkap (kesalahan) {\n    cetak(\"Terjadi galat:\", kesalahan)\n} akhirnya {\n    cetak(\"Proses selesai.\")\n}",
    "expected": "Terjadi galat: Angka pembagi tidak valid!\nProses selesai."
  },
  {
    "id": 22,
    "title": "[22/36] Time-Travel Variable History",
    "section": "Special Features",
    "docLink": "/fitur-khas/time-travel",
    "description": "\n      <p>In almost every programming language, overwriting a variable permanently destroys its prior value. In Celang, variables naturally remember their past!</p>\n      <p>Every variable provides built-in time-travel properties: <code>.sebelumnya</code> for the immediately previous state, <code>.awal</code> for the initial creation value, and <code>.histori</code> for the complete timeline.</p>\n    ",
    "task": "Add a line <code>skor = 100</code> below <code>skor = 50</code>, then inspect <code>skor.sebelumnya</code>.",
    "code": "buat skor = 10\nskor = 25\nskor = 50\n\ncetak(\"Skor saat ini:\", skor)\ncetak(\"Skor sebelumnya:\", skor.sebelumnya)\ncetak(\"Skor awal:\", skor.awal)\ncetak(\"Riwayat skor:\", skor.histori)",
    "expected": "Skor saat ini: 100\nSkor sebelumnya: 50\nSkor awal: 10\nRiwayat skor: [10, 25, 50, 100]"
  },
  {
    "id": 23,
    "title": "[23/36] Transactions & Automatic Rollback (transaksi)",
    "section": "Special Features",
    "docLink": "/fitur-khas/transaksi",
    "description": "\n      <p>In other languages, if a process fails halfway through mutating several variables, your data gets stuck in a broken <em>dirty state</em>.</p>\n      <p>Celang brings database-grade transaction atomicity directly to programming syntax:</p>\n      <ul>\n        <li><strong>The <code>transaksi { ... }</code> Block</strong>: Executes variable mutations atomically.</li>\n        <li><strong><code>batalkan \"reason\"</code></strong>: If a failure condition occurs, Celang immediately rolls back all variables to their original states.</li>\n        <li><strong><code>pulihkan (reason) { ... }</code></strong>: Catches the cancellation notice so your program continues smoothly.</li>\n      </ul>\n    ",
    "task": "Change Andi's initial balance to <code>100.ribu</code> so that his preserved final balance prints <code>Rp 100.000</code>.",
    "code": "buat saldoAndi = 50.ribu\nbuat saldoBudi = 50.ribu\n\ntransaksi {\n    saldoAndi -= 250.ribu\n    saldoBudi += 250.ribu\n\n    jika saldoAndi < 0 {\n        batalkan \"Saldo Andi tidak mencukupi untuk transfer!\"\n    }\n} pulihkan (alasan) {\n    cetak(\"Pemberitahuan:\", alasan)\n}\n\ncetak(\"Saldo akhir Andi tetap aman:\", saldoAndi.rupiah)\ncetak(\"Saldo akhir Budi tetap aman:\", saldoBudi.rupiah)",
    "expected": "Pemberitahuan: Saldo Andi tidak mencukupi untuk transfer!\nSaldo akhir Andi tetap aman: Rp 100.000\nSaldo akhir Budi tetap aman: Rp 50.000"
  },
  {
    "id": 24,
    "title": "[24/36] Reactive Variables & Watchers (reaktif & pantau)",
    "section": "Special Features",
    "docLink": "/fitur-khas/reaktif",
    "description": "\n      <p>Ever used a spreadsheet like Excel? When you modify a number in one cell, any computed formulas update automatically.</p>\n      <p>Celang provides first-class language-level reactivity without external libraries:</p>\n      <ul>\n        <li><strong><code>reaktif total = harga * jumlah</code></strong>: Automatically tracks dependencies and recomputes whenever source variables change.</li>\n        <li><strong><code>pantau variable (baru, lama) { ... }</code></strong>: Dedicated observers that trigger instantly whenever a variable's value mutates.</li>\n      </ul>\n    ",
    "task": "Change <code>jumlah</code> to <code>5</code> so the notification prints <code>Rp 125.000</code>.",
    "code": "buat harga = 25.ribu\nbuat jumlah = 2\n\n// 1. Reactive variable: recomputed automatically\nreaktif total = harga * jumlah\n\n// 2. Watcher: reacts whenever total updates\npantau total (baru, lama) {\n    cetak(\"Notifikasi: Total belanja berubah jadi\", baru.rupiah)\n}\n\ncetak(\"Total awal:\", total.rupiah)\n\n// Mutate source variable:\njumlah = 4",
    "expected": "Total awal: Rp 50.000\nNotifikasi: Total belanja berubah jadi Rp 125.000"
  },
  {
    "id": 25,
    "title": "[25/36] Data Channels (Saluran)",
    "section": "Special Features",
    "docLink": "/fitur-khas/saluran",
    "description": "\n      <p>Imagine you and a teammate work in separate rooms. Instead of walking back and forth to hand over letters, you install a <strong>pneumatic tube pipe</strong> in the wall.</p>\n      <p>That is the intuitive concept behind Celang's <strong>Saluran</strong> (Channels):</p>\n      <ul>\n        <li><strong><code>Saluran(capacity)</code></strong>: Creates a synchronized pipe buffer (e.g. <code>Saluran(2)</code> holds 2 items at once).</li>\n        <li><strong><code>pipa.kirim(\"Doc A\")</code></strong>: Pushes an item into the pipe from one side.</li>\n        <li><strong><code>pipa.terima()</code></strong>: Receives the item sliding out from the other side. Items emerge in first-in-first-out order safely without thread collisions!</li>\n      </ul>\n    ",
    "task": "Change capacity to <code>Saluran(3)</code>, send <code>pipa.kirim(\"Dokumen C\")</code>, and add <code>cetak(\"Barang diterima:\", pipa.terima())</code> at the end.",
    "code": "buat pipa = Saluran(2)\n\npipa.kirim(\"Dokumen A\")\npipa.kirim(\"Dokumen B\")\n\ncetak(\"Barang diterima:\", pipa.terima())\ncetak(\"Barang diterima:\", pipa.terima())",
    "expected": "Barang diterima: Dokumen A\nBarang diterima: Dokumen B\nBarang diterima: Dokumen C"
  },
  {
    "id": 26,
    "title": "[26/36] Channel Streams & Background Workers (sambil)",
    "section": "Special Features",
    "docLink": "/fitur-khas/konkurensi",
    "description": "\n      <p>Continuing the pneumatic tube pipe analogy: what if a coworker in another room continuously pushes items into the pipe while you consume them one by one?</p>\n      <ul>\n        <li><strong><code>sambil { ... }</code></strong>: Runs a task concurrently in the background without blocking the main program.</li>\n        <li><strong><code>untuk pesan dalam pipa { ... }</code></strong>: Stream items out of the channel using a simple standard loop!</li>\n        <li><strong><code>pipa.tutup()</code></strong>: When the sender closes the channel, the consumer loop finishes automatically.</li>\n      </ul>\n    ",
    "task": "Change channel buffer to <code>Saluran(3)</code> and add <code>pipa.kirim(\"Data 3\")</code> inside the <code>sambil</code> block before closing.",
    "code": "buat pipa = Saluran(2)\n\nsambil {\n    pipa.kirim(\"Data 1\")\n    pipa.kirim(\"Data 2\")\n    pipa.tutup()\n}\n\nuntuk pesan dalam pipa {\n    cetak(\"Aliran data:\", pesan)\n}",
    "expected": "Aliran data: Data 1\nAliran data: Data 2\nAliran data: Data 3"
  },
  {
    "id": 27,
    "title": "[27/36] Parallel Processing (bersamaan)",
    "section": "Special Features",
    "docLink": "/fitur-khas/konkurensi",
    "description": "\n      <p>While <code>sambil</code> runs tasks without waiting, the <code>bersamaan</code> block runs several concurrent jobs in parallel and <strong>guarantees to wait until all of them finish</strong> before moving forward.</p>\n      <p>It is like cleaning a room together: all workers clean concurrently, and the door is locked only when everyone has finished their task.</p>\n    ",
    "task": "Add <code>cetak(\"Tugas Paralel C\")</code> inside the <code>bersamaan</code> block.",
    "code": "bersamaan {\n    cetak(\"Tugas Paralel A\")\n    cetak(\"Tugas Paralel B\")\n}\n\ncetak(\"Semua tugas selesai bersamaan!\")",
    "expected": "Tugas Paralel A\nTugas Paralel B\nTugas Paralel C\nSemua tugas selesai bersamaan!"
  },
  {
    "id": 28,
    "title": "[28/36] Background Tasks (titip & jemput)",
    "section": "Special Features",
    "docLink": "/fitur-khas/titip-jemput",
    "description": "\n      <p>Ever ordered coffee or left laundry at a laundromat? You get a <strong>pickup ticket</strong>. You don't stand frozen in front of the counter; you can work on something else, then <strong>pick up</strong> the result when you need it.</p>\n      <p>That is how asynchronous operations work in Celang:</p>\n      <ul>\n        <li><strong><code>titip myFunction()</code></strong>: Delegates heavy work to the background, immediately returning a <strong>ticket</strong> without blocking execution.</li>\n        <li><strong><code>jemput ticket</code></strong>: Resolves and retrieves the completed result whenever you need it.</li>\n      </ul>\n    ",
    "task": "Change the return value to: <code>\"Laporan Tahunan 2026 Disetujui\"</code>.",
    "code": "fungsi siapkanLaporan() {\n    kembalikan \"Laporan Keuangan Selesai\"\n}\n\n// 1. Delegate work to background, receive ticket\nbuat tiket = titip siapkanLaporan()\ncetak(\"Tiket di tangan, kita bisa kerjakan yang lain...\")\n\n// 2. Pick up the result when needed\nbuat hasil = jemput tiket\ncetak(\"Hasil jemput:\", hasil)",
    "expected": "Tiket di tangan, kita bisa kerjakan yang lain...\nHasil jemput: Laporan Tahunan 2026 Disetujui"
  },
  {
    "id": 29,
    "title": "[29/36] Null Safety (mungkin & kalau kosong)",
    "section": "Special Features",
    "docLink": "/fitur-khas/ergonomi-dan-null-safety",
    "description": "\n      <p>Nested property lookups on missing or null objects often cause application crashes. Celang solves this ergonomically with <code>mungkin</code> (optional chaining) and <code>kalau kosong</code> (null coalescing / <code>??</code>).</p>\n      <p>If any segment is missing, <code>mungkin</code> short-circuits gracefully without throwing, and <code>kalau kosong</code> supplies your predefined fallback value.</p>\n    ",
    "task": "Change the default city from <code>\"Jakarta\"</code> to <code>\"Bandung\"</code>.",
    "code": "buat pengguna = {\n    \"profil\": {\n        \"nama\": \"Citra\"\n    }\n}\n\nbuat kota = mungkin pengguna.profil.alamat.kota kalau kosong \"Jakarta\"\ncetak(\"Kota pengguna:\", kota)",
    "expected": "Kota pengguna: Bandung"
  },
  {
    "id": 30,
    "title": "[30/36] Automatic Defer (saat_selesai)",
    "section": "Special Features",
    "docLink": "/dasar/kontrol-alur",
    "description": "\n      <p>The <code>saat_selesai</code> keyword (equivalent to <code>defer</code> in Go) guarantees cleanup actions—like releasing mutex locks, closing files, or tearing down connections—run immediately before the surrounding function exits.</p>\n      <p>Deferred statements are registered up front so they are never forgotten, and execute in reverse registration order (<em>LIFO</em>).</p>\n    ",
    "task": "Add <code>saat_selesai cetak(\"Pembersihan memori selesai.\")</code> below the first defer statement.",
    "code": "fungsi prosesTransaksi() {\n    saat_selesai cetak(\"Pembersihan: sesi koneksi ditutup.\")\n    cetak(\"Langkah 1: Memproses pembayaran...\")\n    cetak(\"Langkah 2: Mencetak resi...\")\n}\n\nprosesTransaksi()",
    "expected": "Langkah 1: Memproses pembayaran...\nLangkah 2: Mencetak resi...\nPembersihan memori selesai.\nPembersihan: sesi koneksi ditutup."
  },
  {
    "id": 31,
    "title": "[31/36] Pipeline Operator (|>)",
    "section": "Special Features",
    "docLink": "/fitur-khas/pipeline",
    "description": "\n      <p>Writing deeply nested function calls (like <code>g(f(x))</code>) can be hard to read. The pipeline operator <code>|&gt;</code> forwards the computed result from the left as the first argument to the function on the right.</p>\n      <p>Data transformations read naturally from left to right.</p>\n    ",
    "task": "Change initial input <code>5</code> to <code>10</code> and observe the result: (10 * 10) + 10 = 110.",
    "code": "fungsi kuadrat(x) {\n    kembalikan x * x\n}\n\nfungsi tambahSepuluh(x) {\n    kembalikan x + 10\n}\n\nbuat hasil = 5 |> kuadrat() |> tambahSepuluh()\ncetak(\"Hasil pipa:\", hasil)",
    "expected": "Hasil pipa: 110"
  },
  {
    "id": 32,
    "title": "[32/36] JSON Data Exchange (json.uraikan & json.rapi)",
    "section": "Special Features",
    "docLink": "/pustaka-standar/json",
    "description": "\n      <p>Virtually all modern web applications exchange structured data using JSON (<em>JavaScript Object Notation</em>).</p>\n      <p>Celang includes a powerful built-in <code>json</code> module ready to use:</p>\n      <ul>\n        <li><strong><code>json.rapi(object)</code></strong>: Pretty-prints Celang maps/arrays into clean, indented JSON text.</li>\n        <li><strong><code>json.uraikan(jsonString)</code></strong>: Parses a JSON string back into native Celang data structures.</li>\n        <li><strong><code>json.dapatkan(obj, path)</code></strong>: Safely traverses nested JSON paths without runtime errors.</li>\n      </ul>\n    ",
    "task": "Change points to <code>2500</code> so the formatted JSON displays <code>\"poin\": 2500</code>.",
    "code": "buat dataPengguna = {\n    \"nama\": \"Citra\",\n    \"kota\": \"Bandung\",\n    \"poin\": 1500\n}\n\n// 1. Format object to pretty JSON:\nbuat teksJSON = json.rapi(dataPengguna)\ncetak(\"Teks JSON:\")\ncetak(teksJSON)\n\n// 2. Parse JSON back into native object:\nbuat hasilUrai = json.uraikan(teksJSON)\ncetak(\"Nama hasil urai:\", hasilUrai.nama)",
    "expected": "Teks JSON:\n{\n  \"kota\": \"Bandung\",\n  \"nama\": \"Citra\",\n  \"poin\": 2500\n}\nNama hasil urai: Citra"
  },
  {
    "id": 33,
    "title": "[33/36] Indonesian Cultural Units & Built-in In-Words",
    "section": "Special Features",
    "docLink": "/fitur-khas/satuan-budaya",
    "description": "\n      <p>In Indonesia, people rarely write large amounts by typing dozens of zeros. People say <em>\"dua setengah juta\"</em> (2.5 million), <em>\"lima puluh ribu\"</em> (50 thousand), or require formal receipts spelled out in words.</p>\n      <p>Celang natively understands these cultural customs right inside the language without any external libraries:</p>\n      <ul>\n        <li><strong>Number Abbreviations</strong>: Write <code>50.ribu</code>, <code>2.5.juta</code>, <code>1.miliar</code>, or <code>10.triliun</code> without counting zeros!</li>\n        <li><strong>Rupiah Currency Format (<code>.rupiah</code>)</strong>: Standard Indonesian formatting with dots (e.g. <code>Rp 2.500.000</code>).</li>\n        <li><strong>Formal In-Words (<code>.terbilang</code>)</strong>: Converts numbers into official Indonesian in-words text ending with \"rupiah\" (e.g. <code>lima puluh ribu rupiah</code>).</li>\n        <li><strong>Pure Number Spelling (<code>.kata</code> / <code>.ejaan</code>)</strong>: Spells out numbers without \"rupiah\" (e.g. <code>105.kata</code> becomes <em>\"seratus lima\"</em>).</li>\n        <li><strong>Time &amp; Data Units</strong>: Use <code>.detik</code>, <code>.menit</code>, <code>.jam</code>, <code>.hari</code>, plus <code>.kb</code>, <code>.mb</code>, <code>.gb</code>, <code>.tb</code>. Celang converts them to milliseconds and bytes automatically!</li>\n      </ul>\n    ",
    "task": "Add <code>buat saldo = 750.ribu</code> and <code>cetak(\"Saldo:\", saldo.rupiah)</code> at the end of the code.",
    "code": "buat modal = 50.ribu\nbuat investasi = 2.5.juta\nbuat waktuTunggu = 5.detik\nbuat kuota = 16.gb\n\ncetak(\"Rupiah:\", investasi.rupiah)\ncetak(\"Terbilang kuitansi:\", modal.terbilang)\ncetak(\"Ejaan angka:\", 105.kata)\ncetak(\"Waktu tunggu (ms):\", waktuTunggu)\ncetak(\"Ukuran kuota (byte):\", kuota)",
    "expected": "Rupiah: Rp 2.500.000\nTerbilang kuitansi: lima puluh ribu rupiah\nEjaan angka: seratus lima\nWaktu tunggu (ms): 5000\nUkuran kuota (byte): 17179869184\nSaldo: Rp 750.000"
  },
  {
    "id": 34,
    "title": "[34/36] C Library Bridge & Ecosystem (ffi)",
    "section": "Special Features",
    "docLink": "/pustaka-standar/ffi",
    "description": "\n      <p>The C/C++ universe has decades of battle-tested high-performance libraries (such as <strong>SQLite</strong>, <strong>Raylib 3D</strong>, <strong>OpenGL</strong>, and <strong>Win32 API</strong>). Why reinvent the wheel?</p>\n      <p>Celang's built-in <strong><code>ffi</code></strong> (<em>Foreign Function Interface</em>) provides a zero-friction bridge:</p>\n      <ul>\n        <li><strong>Zero C Compiler Needed</strong>: Thanks to <em>Dynamic ABI Invocation</em>, Celang invokes <code>.so</code> / <code>.dll</code> binaries directly without GCC, Clang, or CGO!</li>\n        <li><strong>Raw Memory &amp; Pointers</strong>: Allocate raw memory buffers with <code>ffi.alloc(bytes)</code> and free them with <code>ffi.free(ptr)</code>.</li>\n        <li><strong>Frictionless Cross-Compiles</strong>: Build standalone executables with <code>./celang build</code> for Windows or Linux without painful C toolchain dependencies.</li>\n      </ul>\n    ",
    "task": "Change the text written into memory to <code>\"Halo Pustaka C!\"</code>.",
    "code": "impor ffi\n\n// 1. Allocate raw memory buffer (64 bytes):\nbuat buffer = ffi.alloc(64)\n\n// 2. Write text directly to memory address:\nffi.write_string(buffer, \"Data tersimpan di memori C\")\n\n// 3. Read back text from memory pointer:\ncetak(\"Pointer valid:\", !ffi.is_null(buffer))\ncetak(\"Isi memori:\", ffi.read_string(buffer))\n\n// 4. Always free memory after use:\nffi.free(buffer)\ncetak(\"Pembersihan memori selesai!\")",
    "expected": "Pointer valid: benar\nIsi memori: Halo Pustaka C!\nPembersihan memori selesai!"
  },
  {
    "id": 35,
    "title": "[35/36] Automated Unit Testing (celang test)",
    "section": "Special Features",
    "docLink": "/panduan/pengujian",
    "description": "\n      <p>Reliable software is well-tested software. Celang includes a built-in <strong>Test Runner</strong> and <strong>Assertion Library</strong> right into the language core!</p>\n      <ul>\n        <li><strong>Isolated Test Blocks</strong>: Write declarative tests with <code>test \"test name\" { ... }</code> or <code>uji \"...\" { ... }</code>. Each test runs in its own isolated scope.</li>\n        <li><strong>Empathetic Assertions</strong>: Use <code>assert_eq(actual, target)</code>, <code>assert_ne(...)</code>, <code>assert_true(...)</code>, <code>assert_contains(...)</code>, and <code>assert_throws(...)</code>.</li>\n        <li><strong>Terminal Test Runner</strong>: Run <code>celang test</code> in your project to discover and run all tests in milliseconds.</li>\n      </ul>\n    ",
    "task": "Change <code>assert_eq(result, 30)</code> to <code>assert_eq(result, 50)</code> so the test passes.",
    "code": "\"use universal\"\n\nfunction add(a, b) {\n  return a + b\n}\n\ntest \"positive numbers addition\" {\n  let result = add(20, 30)\n  assert_eq(result, 50, \"20 + 30 must equal 50\")\n  assert_ne(result, 0)\n  assert_true(result > 40)\n}\n\nprint(\"✓ Unit testing executed successfully!\")",
    "expected": "✓ Unit testing executed successfully!"
  },
  {
    "id": 36,
    "title": "[36/36] Universal Mode (\"use universal\") & Syntax Styling",
    "section": "Special Features",
    "docLink": "/fitur-khas/mode-universal",
    "description": "\n      <p>You might wonder: <em>\"If Celang is written with Indonesian keywords, what if I want to publish my code to the global open-source community?\"</em></p>\n      <p>Celang includes a powerful built-in solution called <strong>Universal Mode</strong>:</p>\n      <ul>\n        <li><strong>Directive <code>\"use universal\"</code></strong>: Place <code>\"use universal\"</code> at the top of your file. Celang immediately switches to standard international English.</li>\n        <li><strong>Strict Global Standards</strong>: Indonesian keywords (like <code>buat</code>, <code>jika</code>, <code>cetak</code>) are disallowed, guiding you toward clean English standards (<code>let</code>, <code>if</code>, <code>print</code>) and international units (<code>.k</code>, <code>.million</code>, <code>.usd</code>).</li>\n        <li><strong>In-Code Styling (<code>gaya \"...\" { }</code>)</strong>: In standard files, you can borrow other syntax styles mid-code, wrapping snippets with <code>gaya \"javascript\" { ... }</code> or <code>gaya \"pythonic\" { ... }</code>.</li>\n        <li><strong>Automated Formatter CLI</strong>: Reformat an entire project between syntax dialects in seconds: <code>./celang format program.ce --style=universal</code>.</li>\n      </ul>\n      <p>🎉 <strong>Next Step</strong>: Install Celang on your computer, and compile your creations into standalone executables (<code>.exe</code> on Windows or native binary on Linux/macOS) that your friends can double-click and run with zero dependencies using: <code>./celang build program.ce -o Game.exe</code>!</p>\n    ",
    "task": "Change the <code>bonus</code> value to <code>20.k</code> so the total balance becomes $120,000.00.",
    "code": "\"use universal\"\n\nlet balance = 100.k\nconst bonus = 5.k\nlet total = balance + bonus\n\nprint(\"Total USD:\", total.usd)",
    "expected": "Total USD: $120,000.00"
  }
];

// App State
let currentLang = 'id';
let currentIndex = 0;
let engineMode = 'vm';
let wasmReady = false;
let userModifiedCodes = {};

// DOM Elements
const elSlideIndicator = document.getElementById('slide-indicator');
const elBtnPrev = document.getElementById('btn-prev');
const elBtnNext = document.getElementById('btn-next');
const elFooterPrev = document.getElementById('footer-btn-prev');
const elFooterNext = document.getElementById('footer-btn-next');
const elLessonBadge = document.getElementById('lesson-badge');
const elLessonTitle = document.getElementById('lesson-title');
const elLessonContent = document.getElementById('lesson-content');
const elLessonTask = document.getElementById('lesson-task');
const elExpectedOutput = document.getElementById('expected-output');
const elExpectedContainer = document.getElementById('expected-container');
const elExpectedArrow = document.getElementById('expected-arrow');
const elBtnToggleExpected = document.getElementById('btn-toggle-expected');
const elCodeEditor = document.getElementById('code-editor');
const elTerminalOutput = document.getElementById('terminal-output');
const elExecutionTime = document.getElementById('execution-time');
const elWasmBadge = document.getElementById('wasm-status-badge');
const elBtnRun = document.getElementById('btn-run');
const elBtnReset = document.getElementById('btn-reset');
const elBtnClearTerm = document.getElementById('btn-clear-term');
const elBtnTheme = document.getElementById('btn-theme');
const elBtnToc = document.getElementById('btn-toc');
const elBtnCloseToc = document.getElementById('btn-close-toc');
const elTocModal = document.getElementById('toc-modal');
const elTocListPengantar = document.getElementById('toc-list-pengantar');
const elTocListDasar = document.getElementById('toc-list-dasar');
const elTocListSpesial = document.getElementById('toc-list-spesial');
const elResizer = document.getElementById('resizer');
const elLessonPanel = document.getElementById('lesson-panel');
const elBtnLang = document.getElementById('btn-lang');
const elLangFlag = document.getElementById('lang-flag');
const elLangText = document.getElementById('lang-text');
const elLabelDocsLink = document.getElementById('label-docs-link');

function getCurrentSlides() {
  return currentLang === 'en' ? SLIDES_EN : SLIDES_ID;
}

function getTexts() {
  return UI_TEXTS[currentLang] || UI_TEXTS.id;
}

// Initialize App
function initApp() {
  loadSavedLanguage();
  buildTableOfContents();
  setupEventListeners();
  loadSavedTheme();
  loadInitialSlide();
  initWasm();
}

function loadSavedLanguage() {
  const saved = localStorage.getItem('celang_tour_lang');
  if (saved === 'en' || saved === 'id') {
    currentLang = saved;
  } else {
    currentLang = 'id';
  }
  updateLanguageUI();
}

function setLanguage(lang) {
  if (lang !== 'id' && lang !== 'en') return;
  currentLang = lang;
  localStorage.setItem('celang_tour_lang', lang);
  updateLanguageUI();
  buildTableOfContents();
  goToSlide(currentIndex);
}

function updateLanguageUI() {
  const t = getTexts();
  if (elLangFlag) elLangFlag.textContent = currentLang === 'en' ? '🇬🇧' : '🇮🇩';
  if (elLangText) elLangText.textContent = currentLang === 'en' ? 'EN' : 'ID';
  if (elLabelDocsLink) elLabelDocsLink.textContent = t.docsLink;
  const elDocLinkText = document.getElementById('doc-link-text');
  if (elDocLinkText) elDocLinkText.textContent = t.readDocs;
  if (elFooterPrev) elFooterPrev.textContent = t.footerPrev;
  if (elFooterNext) elFooterNext.textContent = t.footerNext;
  if (elBtnRun) elBtnRun.innerHTML = t.runBtn;
  if (elBtnReset) elBtnReset.textContent = t.resetBtn;
  if (elBtnClearTerm) elBtnClearTerm.textContent = t.clearTerm;

  const btnTocText = document.querySelector('.toc-btn-text');
  if (btnTocText) btnTocText.textContent = t.tocBtn;

  const tocModalTitle = document.getElementById('toc-modal-title');
  if (tocModalTitle) tocModalTitle.textContent = t.tocTitle;

  const elTocSearch = document.getElementById('toc-search-input');
  if (elTocSearch) elTocSearch.placeholder = t.tocSearchPlaceholder;

  const elTocNoRes = document.querySelector('#toc-no-results span');
  if (elTocNoRes) elTocNoRes.textContent = t.tocNoResults;
}

const SLIDE_META = [
  { icon: '🌟', tag: 'Mulai' },
  { icon: '🌐', tag: 'ID / EN' },
  { icon: '👋', tag: 'cetak()' },
  { icon: '📦', tag: 'buat' },
  { icon: '🔒', tag: 'konstanta' },
  { icon: '🏷️', tag: 'Tipe Data' },
  { icon: '💬', tag: 'f-string' },
  { icon: '🔀', tag: 'jika ... selain' },
  { icon: '🎯', tag: 'pilih ... bila' },
  { icon: '🔁', tag: 'selama' },
  { icon: '🔢', tag: 'untuk ... dalam' },
  { icon: '🛑', tag: 'lanjut & berhenti' },
  { icon: '📚', tag: 'Larik [ ... ]' },
  { icon: '🗺️', tag: 'Kamus { ... }' },
  { icon: '🧩', tag: 'Destrukturisasi' },
  { icon: '⚙️', tag: 'fungsi -> Tipe' },
  { icon: '🏹', tag: 'Panah =>' },
  { icon: '🔗', tag: 'rujuk & ref' },
  { icon: '🏗️', tag: 'bentuk & diri / @' },
  { icon: '📦', tag: 'Generik <T>' },
  { icon: '🛡️', tag: 'coba ... tangkap' },
  { icon: '⏳', tag: 'Time-Travel' },
  { icon: '💳', tag: 'transaksi & batalkan' },
  { icon: '⚡', tag: 'reaktif & pantau' },
  { icon: '📮', tag: 'Saluran' },
  { icon: '🌊', tag: 'sambil' },
  { icon: '👥', tag: 'bersamaan' },
  { icon: '🎟️', tag: 'titip & jemput' },
  { icon: '🛡️', tag: 'Null-Safety' },
  { icon: '🧹', tag: 'saat_selesai' },
  { icon: '🚰', tag: 'Pipa |>' },
  { icon: '📄', tag: 'Modul json' },
  { icon: '🇮🇩', tag: 'Satuan Budaya' },
  { icon: '🔌', tag: 'Modul ffi' },
  { icon: '🌍', tag: '"use universal"' }
];

const SLIDE_META_EN = [
  { icon: '🌟', tag: 'Start' },
  { icon: '🌐', tag: 'ID / EN' },
  { icon: '👋', tag: 'print()' },
  { icon: '📦', tag: 'variables' },
  { icon: '🔒', tag: 'constants' },
  { icon: '🏷️', tag: 'Data Types' },
  { icon: '💬', tag: 'f-strings' },
  { icon: '🔀', tag: 'if ... else' },
  { icon: '🎯', tag: 'switch ... case' },
  { icon: '🔁', tag: 'while loop' },
  { icon: '🔢', tag: 'for in ranges' },
  { icon: '🛑', tag: 'break & continue' },
  { icon: '📚', tag: 'Arrays [ ... ]' },
  { icon: '🗺️', tag: 'Maps { ... }' },
  { icon: '🧩', tag: 'Destructuring' },
  { icon: '⚙️', tag: 'functions -> Type' },
  { icon: '🏹', tag: 'Arrow =>' },
  { icon: '🔗', tag: 'rujuk / ref' },
  { icon: '🏗️', tag: 'structs & self' },
  { icon: '📦', tag: 'Generics <T>' },
  { icon: '🛡️', tag: 'try ... catch' },
  { icon: '⏳', tag: 'Time-Travel' },
  { icon: '💳', tag: 'transactions' },
  { icon: '⚡', tag: 'reactivity' },
  { icon: '📮', tag: 'Channels' },
  { icon: '🌊', tag: 'async stream' },
  { icon: '👥', tag: 'parallel' },
  { icon: '🎟️', tag: 'promises' },
  { icon: '🛡️', tag: 'null safety' },
  { icon: '🧹', tag: 'defer' },
  { icon: '🚰', tag: 'pipeline |>' },
  { icon: '📄', tag: 'JSON module' },
  { icon: '🇮🇩', tag: 'Cultural units' },
  { icon: '🔌', tag: 'FFI C Bridge' },
  { icon: '🌍', tag: '"use universal"' }
];

function updateTocProgress() {
  const slides = getCurrentSlides();
  const currentNum = currentIndex + 1;
  const pct = Math.round((currentNum / slides.length) * 100);

  const elProgressText = document.getElementById('toc-progress-text');
  const elProgressBar = document.getElementById('toc-progress-bar');
  if (elProgressText) {
    elProgressText.textContent = currentLang === 'id'
      ? `Langkah ${currentNum} dari ${slides.length} (${pct}%)`
      : `Lesson ${currentNum} of ${slides.length} (${pct}%)`;
  }
  if (elProgressBar) {
    elProgressBar.style.width = `${pct}%`;
  }
}

function buildTableOfContents() {
  const slides = getCurrentSlides();
  const t = getTexts();
  const metaList = currentLang === 'id' ? SLIDE_META : SLIDE_META_EN;

  const secTitles = document.querySelectorAll('.toc-section-title');
  if (secTitles.length >= 3) {
    secTitles[0].textContent = t.secPengantar;
    secTitles[1].textContent = t.secDasar;
    secTitles[2].textContent = t.secSpesial;
  }

  const countPengantar = document.getElementById('toc-count-pengantar');
  const countDasar = document.getElementById('toc-count-dasar');
  const countSpesial = document.getElementById('toc-count-spesial');
  const pCount = slides.filter(s => s.section === 'Pengantar' || s.section === 'Introduction').length;
  const dCount = slides.filter(s => s.section === 'Dasar Bahasa' || s.section === 'Language Basics').length;
  const sCount = slides.length - pCount - dCount;
  if (countPengantar) countPengantar.textContent = currentLang === 'id' ? `${pCount} Langkah` : `${pCount} Lessons`;
  if (countDasar) countDasar.textContent = currentLang === 'id' ? `${dCount} Langkah` : `${dCount} Lessons`;
  if (countSpesial) countSpesial.textContent = currentLang === 'id' ? `${sCount} Langkah` : `${sCount} Lessons`;

  updateTocProgress();

  if (elTocListPengantar) elTocListPengantar.innerHTML = '';
  if (elTocListDasar) elTocListDasar.innerHTML = '';
  if (elTocListSpesial) elTocListSpesial.innerHTML = '';

  slides.forEach((slide, idx) => {
    const meta = metaList[idx] || { icon: '📌', tag: 'Materi' };
    const numStr = (idx + 1).toString().padStart(2, '0');
    const cleanTitle = slide.title.replace(/^\[\d+\/\d+\]\s*/, '');

    const card = document.createElement('div');
    card.className = `toc-card${idx === currentIndex ? ' active' : ''}`;
    card.id = `toc-item-${idx}`;
    card.dataset.index = idx;
    card.dataset.title = slide.title.toLowerCase();
    card.dataset.tag = meta.tag.toLowerCase();

    card.innerHTML = `
      <div class="toc-card-left">
        <span class="toc-card-num">${numStr}</span>
        <span class="toc-card-icon">${meta.icon}</span>
        <div class="toc-card-text">
          <span class="toc-card-title">${cleanTitle}</span>
          <span class="toc-card-tag">${meta.tag}</span>
        </div>
      </div>
      <div class="toc-card-right">
        <span class="toc-badge-active">${t.tocActiveBadge || 'Aktif'}</span>
        <span class="toc-card-arrow">➔</span>
      </div>
    `;

    card.addEventListener('click', () => {
      goToSlide(idx);
      closeTocModal();
    });

    if (slide.section === 'Pengantar' || slide.section === 'Introduction') {
      if (elTocListPengantar) elTocListPengantar.appendChild(card);
    } else if (slide.section === 'Dasar Bahasa' || slide.section === 'Language Basics') {
      if (elTocListDasar) elTocListDasar.appendChild(card);
    } else {
      if (elTocListSpesial) elTocListSpesial.appendChild(card);
    }
  });
}

function loadInitialSlide() {
  const slides = getCurrentSlides();
  let initialIndex = 0;
  if (window.location.hash) {
    const match = window.location.hash.match(/\d+/);
    if (match) {
      const parsed = parseInt(match[0], 10) - 1;
      if (parsed >= 0 && parsed < slides.length) {
        initialIndex = parsed;
      }
    }
  } else {
    const saved = localStorage.getItem('celang_tour_slide');
    if (saved !== null) {
      const parsed = parseInt(saved, 10);
      if (parsed >= 0 && parsed < slides.length) {
        initialIndex = parsed;
      }
    }
  }

  goToSlide(initialIndex);
}

function goToSlide(index) {
  const slides = getCurrentSlides();
  if (index < 0 || index >= slides.length) return;

  currentIndex = index;
  const slide = slides[currentIndex];
  const t = getTexts();

  // Update Header Indicator
  elSlideIndicator.textContent = `${currentIndex + 1} / ${slides.length}`;
  elBtnPrev.disabled = currentIndex === 0;
  elFooterPrev.disabled = currentIndex === 0;
  elBtnNext.disabled = currentIndex === slides.length - 1;
  elFooterNext.disabled = currentIndex === slides.length - 1;

  // Update Lesson Content
  elLessonBadge.textContent = slide.section;
  elLessonTitle.textContent = slide.title;
  elLessonContent.innerHTML = slide.description;
  elLessonTask.innerHTML = slide.task;
  elExpectedOutput.textContent = slide.expected;

  // Update Documentation Link
  const elDocLink = document.getElementById('doc-link');
  const elDocLinkText = document.getElementById('doc-link-text');
  const elDocLinkBox = document.getElementById('doc-link-box');
  if (elDocLink && slide.docLink) {
    elDocLink.href = slide.docLink;
    if (elDocLinkText) elDocLinkText.textContent = t.readDocs;
    if (elDocLinkBox) elDocLinkBox.style.display = 'flex';
  } else if (elDocLinkBox) {
    elDocLinkBox.style.display = 'none';
  }

  const challengeHeaderStrong = document.querySelector('.challenge-header strong');
  if (challengeHeaderStrong) challengeHeaderStrong.textContent = t.challengeTitle;

  const expectedToggleSpan = document.querySelector('.expected-toggle span');
  if (expectedToggleSpan) expectedToggleSpan.textContent = t.expectedToggle;

  // Attach language selector listeners if on Slide 2
  if (currentIndex === 1) {
    const btnId = document.getElementById('card-lang-id');
    const btnEn = document.getElementById('card-lang-en');
    if (btnId) {
      btnId.addEventListener('click', () => setLanguage('id'));
    }
    if (btnEn) {
      btnEn.addEventListener('click', () => setLanguage('en'));
    }
  }

  // Maintain current collapsed/expanded state
  const isCollapsed = elExpectedContainer.classList.contains('collapsed');
  elExpectedArrow.textContent = isCollapsed ? '▶' : '▼';

  // Restore Code (from edits or original default)
  const codeKey = `${currentLang}_${currentIndex}`;
  elCodeEditor.value = userModifiedCodes[codeKey] !== undefined
    ? userModifiedCodes[codeKey]
    : slide.code;

  // Reset Terminal
  clearTerminal();

  // Update Hash & Storage
  window.location.hash = `#${currentIndex + 1}`;
  localStorage.setItem('celang_tour_slide', currentIndex);

  // Update TOC active state & progress
  document.querySelectorAll('.toc-card').forEach((c, idx) => {
    c.classList.toggle('active', idx === currentIndex);
  });
  updateTocProgress();
}

function clearTerminal() {
  const t = getTexts();
  elTerminalOutput.textContent = t.termPlaceholder;
  elTerminalOutput.className = 'terminal-pre placeholder';
  elExecutionTime.textContent = '';
}

function resetCurrentCode() {
  const slides = getCurrentSlides();
  const slide = slides[currentIndex];
  const codeKey = `${currentLang}_${currentIndex}`;
  delete userModifiedCodes[codeKey];
  elCodeEditor.value = slide.code;
  clearTerminal();
}

function executeCode() {
  const t = getTexts();
  const code = elCodeEditor.value;
  if (!code.trim()) {
    elTerminalOutput.textContent = t.emptyCode;
    elTerminalOutput.className = 'terminal-pre placeholder';
    return;
  }

  if (!wasmReady || typeof window.runCelang !== 'function') {
    elTerminalOutput.textContent = t.loadingWasm;
    elTerminalOutput.className = 'terminal-pre placeholder';
    return;
  }

  elBtnRun.disabled = true;
  elBtnRun.innerHTML = t.running;

  setTimeout(() => {
    try {
      const res = window.runCelang(code, engineMode, 'indonesian');
      elExecutionTime.textContent = `⏱ ${res.durationMs} ms (${engineMode.toUpperCase()})`;

      if (res.error) {
        elTerminalOutput.className = 'terminal-pre error';
        elTerminalOutput.textContent = (res.output ? res.output + '\n' : '') + '❌ ' + res.error;
      } else {
        elTerminalOutput.className = 'terminal-pre';
        let finalOut = res.output ? res.output.trimEnd() : '';
        if (res.result && res.result !== 'kosong' && res.result !== 'null') {
          finalOut += (finalOut ? '\n' : '') + '=> ' + res.result;
        }
        elTerminalOutput.textContent = finalOut || t.termDoneNoOutput;
      }
    } catch (err) {
      elTerminalOutput.className = 'terminal-pre error';
      elTerminalOutput.textContent = '❌ Error: ' + err.message;
    } finally {
      elBtnRun.disabled = false;
      elBtnRun.innerHTML = t.runBtn;
    }
  }, 10);
}

// WASM Initialization
function initWasm() {
  const t = getTexts();
  if (window.celangWasmReady && typeof window.runCelang === 'function') {
    setWasmReadyState();
    return;
  }

  if (typeof window.Go === 'undefined') {
    elWasmBadge.textContent = '❌ Go runtime not found';
    return;
  }

  const go = new window.Go();
  WebAssembly.instantiateStreaming(fetch('/celang.wasm'), go.importObject)
    .then((result) => {
      go.run(result.instance);
      setWasmReadyState();
    })
    .catch((err) => {
      console.error('WASM Load Error:', err);
      elWasmBadge.textContent = t.wasmError;
    });
}

function setWasmReadyState() {
  const t = getTexts();
  wasmReady = true;
  elWasmBadge.textContent = t.wasmReady;
  elWasmBadge.classList.add('ready');
}

// Event Listeners
function setupEventListeners() {
  // Next & Prev buttons
  elBtnPrev.addEventListener('click', () => goToSlide(currentIndex - 1));
  elFooterPrev.addEventListener('click', () => goToSlide(currentIndex - 1));
  elBtnNext.addEventListener('click', () => goToSlide(currentIndex + 1));
  elFooterNext.addEventListener('click', () => goToSlide(currentIndex + 1));

  // Run & Reset
  elBtnRun.addEventListener('click', executeCode);
  elBtnReset.addEventListener('click', resetCurrentCode);
  elBtnClearTerm.addEventListener('click', clearTerminal);

  // Engine Radios
  document.querySelectorAll('input[name="engine"]').forEach((radio) => {
    radio.addEventListener('change', (e) => {
      engineMode = e.target.value;
    });
  });

  // Expected output accordion
  elBtnToggleExpected.addEventListener('click', () => {
    const isCollapsed = elExpectedContainer.classList.toggle('collapsed');
    elExpectedArrow.textContent = isCollapsed ? '▶' : '▼';
  });

  // Language button in navbar
  if (elBtnLang) {
    elBtnLang.addEventListener('click', () => {
      setLanguage(currentLang === 'id' ? 'en' : 'id');
    });
  }

  // Save user code edits on input
  elCodeEditor.addEventListener('input', () => {
    const codeKey = `${currentLang}_${currentIndex}`;
    userModifiedCodes[codeKey] = elCodeEditor.value;
  });

  // Code editor tab key indentation & shortcuts
  elCodeEditor.addEventListener('keydown', (e) => {
    if (e.key === 'Tab') {
      e.preventDefault();
      const start = elCodeEditor.selectionStart;
      const end = elCodeEditor.selectionEnd;
      elCodeEditor.value = elCodeEditor.value.substring(0, start) + '    ' + elCodeEditor.value.substring(end);
      elCodeEditor.selectionStart = elCodeEditor.selectionEnd = start + 4;
      const codeKey = `${currentLang}_${currentIndex}`;
      userModifiedCodes[codeKey] = elCodeEditor.value;
    } else if ((e.ctrlKey || e.metaKey) && e.key === 'Enter') {
      e.preventDefault();
      executeCode();
    }
  });

  // Global keyboard shortcuts (Left/Right arrow when not in editor)
  window.addEventListener('keydown', (e) => {
    if (document.activeElement === elCodeEditor) return;
    const slides = getCurrentSlides();
    if (e.key === 'ArrowLeft' && currentIndex > 0) {
      goToSlide(currentIndex - 1);
    } else if (e.key === 'ArrowRight' && currentIndex < slides.length - 1) {
      goToSlide(currentIndex + 1);
    }
  });

  // Table of Contents Modal
  elBtnToc.addEventListener('click', openTocModal);
  elBtnCloseToc.addEventListener('click', closeTocModal);
  elTocModal.addEventListener('click', (e) => {
    if (e.target === elTocModal) closeTocModal();
  });

  // TOC Real-time Search
  const elTocSearchInput = document.getElementById('toc-search-input');
  if (elTocSearchInput) {
    elTocSearchInput.addEventListener('input', (e) => {
      const q = e.target.value.toLowerCase().trim();
      const cards = document.querySelectorAll('.toc-card');
      let visibleCount = 0;

      cards.forEach(card => {
        const title = card.dataset.title || '';
        const tag = card.dataset.tag || '';
        const match = !q || title.includes(q) || tag.includes(q);
        card.style.display = match ? 'flex' : 'none';
        if (match) visibleCount++;
      });

      ['pengantar', 'dasar', 'spesial'].forEach(sec => {
        const list = document.getElementById(`toc-list-${sec}`);
        const secBox = document.getElementById(`toc-section-${sec}`);
        if (list && secBox) {
          const hasVisible = Array.from(list.children).some(c => c.style.display !== 'none');
          secBox.style.display = hasVisible ? 'flex' : 'none';
        }
      });

      const elNoResults = document.getElementById('toc-no-results');
      if (elNoResults) {
        elNoResults.classList.toggle('hidden', visibleCount > 0);
      }
    });
  }

  // Theme Toggle
  elBtnTheme.addEventListener('click', toggleTheme);

  // Split resizer drag
  setupResizer();
}

function openTocModal() {
  elTocModal.classList.remove('hidden');
  updateTocProgress();

  // Reset search
  const elSearch = document.getElementById('toc-search-input');
  if (elSearch) {
    elSearch.value = '';
    elSearch.dispatchEvent(new Event('input'));
    setTimeout(() => elSearch.focus(), 80);
  }

  // Scroll active slide into view smoothly
  setTimeout(() => {
    const activeCard = document.querySelector('.toc-card.active');
    if (activeCard) {
      activeCard.scrollIntoView({ behavior: 'smooth', block: 'center' });
    }
  }, 100);
}

function closeTocModal() {
  elTocModal.classList.add('hidden');
}

function toggleTheme() {
  const isLight = document.body.classList.toggle('light-theme');
  elBtnTheme.textContent = isLight ? '☀️' : '🌙';
  localStorage.setItem('celang_tour_theme', isLight ? 'light' : 'dark');
}

function loadSavedTheme() {
  const saved = localStorage.getItem('celang_tour_theme');
  if (saved === 'light') {
    document.body.classList.add('light-theme');
    elBtnTheme.textContent = '☀️';
  } else {
    document.body.classList.remove('light-theme');
    elBtnTheme.textContent = '🌙';
  }
}

// Resizable Divider
function setupResizer() {
  let isDragging = false;

  elResizer.addEventListener('mousedown', (e) => {
    isDragging = true;
    document.body.style.cursor = 'col-resize';
    document.body.style.userSelect = 'none';
  });

  window.addEventListener('mousemove', (e) => {
    if (!isDragging) return;
    const containerWidth = window.innerWidth;
    const newWidth = Math.max(300, Math.min(e.clientX, containerWidth - 360));
    elLessonPanel.style.width = `${newWidth}px`;
  });

  window.addEventListener('mouseup', () => {
    if (isDragging) {
      isDragging = false;
      document.body.style.cursor = '';
      document.body.style.userSelect = '';
    }
  });
}

// Start app on DOM load or immediately if already ready
if (document.readyState === 'loading') {
  document.addEventListener('DOMContentLoaded', initApp);
} else {
  initApp();
}
