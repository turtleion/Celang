import { defineConfig } from 'vitepress'

export default defineConfig({
  title: "Celang",
  description: "Bahasa Pemrograman Human-Centric & Syntax-Styling",
  lang: 'id-ID',
  head: [
    ['link', { rel: 'icon', href: '/favicon.ico' }],
    ['meta', { name: 'theme-color', content: '#3eaf7c' }]
  ],
  markdown: {
    languages: [
      {
        name: 'celang',
        scopeName: 'source.celang',
        displayName: 'Celang',
        aliases: ['ce'],
        patterns: [
          {
            name: 'comment.line.double-slash.celang',
            match: '//.*$'
          },
          {
            name: 'string.quoted.double.celang',
            match: '"[^"\\\\]*(\\\\.[^"\\\\]*)*"'
          },
          {
            name: 'string.quoted.single.celang',
            match: '\'[^\'\\\\]*(\\\\.[^\'\\\\]*)*\''
          },
          {
            name: 'string.template.celang',
            match: '`[^`\\\\]*(\\\\.[^`\\\\]*)*`'
          },
          {
            name: 'keyword.control.celang',
            match: '\\b(jika|selain jika|selain|maka|cocok|kasus|lainnya|switch|case|default|selama|untuk|dalam|kembalikan|mengembalikan|berhenti|lanjut|sambil|bersamaan|mundur ke|coba|tangkap|akhirnya|lempar|try|catch|finally|throw|gaya|mungkin|saat_selesai|kalau|bila|defer|titip|jemput|harus|transaksi|batalkan|pulihkan|pantau|transaction|rollback|recover|watch|if|else|while|for|return|break|continue|async|parallel|lamun|salami|pikeun|balikkeun)\\b'
          },
          {
            name: 'keyword.declaration.celang',
            match: '\\b(buat|konstanta|fungsi|bentuk|itu|rujuk|sebagai|tipe|privat|sembarang|saluran|reaktif|reactive|let|const|var|function|def|jieun|tetep|struct|ref|as)\\b'
          },
          {
            name: 'constant.language.celang',
            match: '\\b(benar|salah|kosong|true|false|null|None|bener)\\b'
          },
          {
            name: 'support.function.celang',
            match: '\\b(cetak|masukan|panjang|tipe|rentang|saring|petakan|potong_spasi|kapitalisasi|salin|clone|Saluran|jemput_semua|jemput_pertama|terbilang|rupiah|format_rupiah|Angka|Pecahan|Teks|Logika|print|input|len|type|range|filter|map|trim|format_usd|format_currency|in_words|spelled_out|dollars_in_words|Channel|channel|await|await_all|await_race|Number|String|Boolean)\\b'
          },
          {
            name: 'constant.numeric.celang',
            match: '\\b\\d+(\\.\\d+)?\\b'
          }
        ]
      }
    ]
  },
  themeConfig: {
    logo: '🌟',
    siteTitle: 'Celang',
    nav: [
      { text: 'Panduan', link: '/panduan/mulai-cepat' },
      { text: 'Dasar Bahasa', link: '/dasar/variabel-tipe' },
      { text: 'Pustaka Standar', link: '/pustaka-standar/matematika' },
      { text: 'Playground', link: '/playground' },
      {
        text: 'Fitur Khusus',
        items: [
          { text: 'Struktur Data & Tipe', link: '/fitur-khas/bentuk-struct' },
          { text: 'Konkurensi & Saluran', link: '/fitur-khas/konkurensi' },
          { text: 'Multi-Dialek & Gaya Bahasa', link: '/gaya-sintaks/multi-dialek' },
          { text: 'Time-Travel Debugging', link: '/fitur-khas/time-travel' },
          { text: 'Arsitektur & Mesin VM', link: '/arsitektur/bytecode-vm' }
        ]
      }
    ],
    sidebar: [
      {
        text: '🚀 Panduan & Memulai',
        collapsed: false,
        items: [
          { text: 'Pengenalan & Filosofi', link: '/panduan/pengenalan' },
          { text: 'Instalasi & Mulai Cepat', link: '/panduan/mulai-cepat' },
          { text: 'CLI & REPL Interaktif', link: '/panduan/cli-repl' }
        ]
      },
      {
        text: '📘 Sintaks & Dasar Bahasa',
        collapsed: false,
        items: [
          { text: 'Variabel & Tipe Data', link: '/dasar/variabel-tipe' },
          { text: 'Operator & Ekspresi', link: '/dasar/operator' },
          { text: 'Kontrol Alur & Percabangan', link: '/dasar/kontrol-alur' },
          { text: 'Fungsi & Ruang Lingkup', link: '/dasar/fungsi' },
          { text: 'Sistem Modul & Impor', link: '/dasar/modul' },
          { text: 'Fungsi Bawaan Global', link: '/dasar/fungsi-bawaan' }
        ]
      },
      {
        text: '🏗️ Struktur Data & Tipe Lanjutan',
        collapsed: false,
        items: [
          { text: 'Bentuk Data Kustom (bentuk / struct)', link: '/fitur-khas/bentuk-struct' },
          { text: 'Tipe Generik (Generic Types)', link: '/fitur-khas/generik' },
          { text: 'Destrukturisasi Array & Map', link: '/fitur-khas/destrukturisasi' }
        ]
      },
      {
        text: '⚡ Ergonomi & Sintaks Modern',
        collapsed: false,
        items: [
          { text: 'Fungsi Panah Ringkas (=>)', link: '/fitur-khas/fungsi-ringkas' },
          { text: 'Interpolasi Teks (String Template)', link: '/fitur-khas/interpolasi-teks' },
          { text: 'Operator Penyaluran (Pipeline |>)', link: '/fitur-khas/pipeline' },
          { text: 'Keamanan Nilai Kosong (Null-Safety)', link: '/fitur-khas/ergonomi-dan-null-safety' }
        ]
      },
      {
        text: '🔄 Konkurensi & Asinkron',
        collapsed: false,
        items: [
          { text: 'Pekerjaan Bersamaan (sambil & bersamaan)', link: '/fitur-khas/konkurensi' },
          { text: 'Tugas Latar Belakang (titip & jemput)', link: '/fitur-khas/titip-jemput' },
          { text: 'Saluran Komunikasi (Saluran / Channel)', link: '/fitur-khas/saluran' }
        ]
      },
      {
        text: '🎨 Multi-Dialek & Personalisasi',
        collapsed: false,
        items: [
          { text: 'Konsep Multi-Dialek', link: '/gaya-sintaks/multi-dialek' },
          { text: 'Pengalihan Dialek di Kode (In-Code Styling)', link: '/gaya-sintaks/in-code-styling' },
          { text: 'Membuat Dialek Kustom (.style)', link: '/gaya-sintaks/berkas-style' }
        ]
      },
      {
        text: '⏳ Fitur Eksklusif',
        collapsed: false,
        items: [
          { text: 'Time-Travel Debugging (Riwayat Variabel)', link: '/fitur-khas/time-travel' },
          { text: 'Transaksi & Rollback Otomatis', link: '/fitur-khas/transaksi' },
          { text: 'Variabel Reaktif & Pengamat', link: '/fitur-khas/reaktif' },
          { text: 'Satuan Budaya & Terbilang Bawaan', link: '/fitur-khas/satuan-budaya' },
          { text: 'Mode Universal / Internasional ("use universal")', link: '/fitur-khas/mode-universal' }
        ]
      },
      {
        text: '📦 Pustaka Standar (Standard Library)',
        collapsed: false,
        items: [
          { text: 'Matematika & Numerik (matematika / math)', link: '/pustaka-standar/matematika' },
          { text: 'Pengolah Teks & String (teks / strings)', link: '/pustaka-standar/teks' },
          { text: 'Sistem Berkas & I/O (berkas / fs)', link: '/pustaka-standar/berkas' },
          { text: 'Waktu, Tanggal & Durasi (waktu / time)', link: '/pustaka-standar/waktu' },
          { text: 'Pertukaran Data JSON (json)', link: '/pustaka-standar/json' },
          { text: 'Sistem Operasi & Proses (sistem / os)', link: '/pustaka-standar/sistem' }
        ]
      },
      {
        text: '⚙️ Arsitektur & Mesin Virtual',
        collapsed: false,
        items: [
          { text: 'Pohon Sintaks Universal (U-AST)', link: '/arsitektur/uast' },
          { text: 'Kompiler Bytecode & Mesin Virtual (VM)', link: '/arsitektur/bytecode-vm' }
        ]
      }
    ],
    socialLinks: [
      { icon: 'github', link: 'https://github.com/celang/celang' }
    ],
    footer: {
      message: 'Dirilis di bawah Lisensi MIT. Dibuat dengan cinta untuk ekosistem pemrograman Indonesia.',
      copyright: 'Hak Cipta © 2026 Celang Project'
    },
    search: {
      provider: 'local'
    }
  }
})
