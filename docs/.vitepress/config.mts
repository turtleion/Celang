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
            match: '\\b(jika|selain jika|selain|maka|cocok|kasus|lainnya|switch|case|default|selama|untuk|dalam|kembalikan|mengembalikan|berhenti|lanjut|sambil|bersamaan|mundur ke|coba|tangkap|akhirnya|lempar|try|catch|finally|throw|gaya|mungkin|saat_selesai|kalau|bila|defer|titip|jemput|harus|if|else|while|for|return|break|continue|async|parallel|lamun|salami|pikeun|balikkeun)\\b'
          },
          {
            name: 'keyword.declaration.celang',
            match: '\\b(buat|konstanta|fungsi|bentuk|itu|rujuk|sebagai|tipe|privat|sembarang|saluran|let|const|var|function|def|jieun|tetep|struct|ref|as)\\b'
          },
          {
            name: 'constant.language.celang',
            match: '\\b(benar|salah|kosong|true|false|null|None|bener)\\b'
          },
          {
            name: 'support.function.celang',
            match: '\\b(cetak|masukan|panjang|tipe|rentang|saring|petakan|potong_spasi|kapitalisasi|salin|clone|Saluran|jemput_semua|jemput_pertama|Angka|Pecahan|Teks|Logika|print|input|len|type|range|filter|map|trim)\\b'
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
      { text: 'Playground', link: '/playground' },
      { text: 'Fitur Unggulan', link: '/fitur-khas/time-travel' },
      { text: 'Gaya Bahasa', link: '/gaya-sintaks/multi-dialek' },
      { text: 'Alat Bawaan', link: '/pustaka-standar/matematika' },
      { text: 'Cara Kerja Mesin', link: '/arsitektur/bytecode-vm' }
    ],
    sidebar: [
      {
        text: '🚀 Panduan Memulai',
        items: [
          { text: 'Filosofi & Pengenalan', link: '/panduan/pengenalan' },
          { text: 'Mulai Cepat & Langkah Pertama', link: '/panduan/mulai-cepat' },
          { text: 'Perintah Terminal & REPL', link: '/panduan/cli-repl' }
        ]
      },
      {
        text: '🧠 Dasar Bahasa',
        items: [
          { text: 'Variabel & Jenis Data', link: '/dasar/variabel-tipe' },
          { text: 'Tanda Hitung & Penghubung', link: '/dasar/operator' },
          { text: 'Keputusan, Perulangan, & Kendala', link: '/dasar/kontrol-alur' },
          { text: 'Fungsi & Tindakan', link: '/dasar/fungsi' },
          { text: 'Perintah Bawaan Siap Pakai', link: '/dasar/fungsi-bawaan' }
        ]
      },
      {
        text: '✨ Fitur Khas Celang',
        items: [
          { text: 'Variabel yang Mengingat Masa Lalu', link: '/fitur-khas/time-travel' },
          { text: 'Bentuk Data Kustom (bentuk)', link: '/fitur-khas/bentuk-struct' },
          { text: 'Wadah Fleksibel (Generik)', link: '/fitur-khas/generik' },
          { text: 'Titip & Jemput (Tugas Latar Belakang)', link: '/fitur-khas/titip-jemput' },
          { text: 'Saluran Pengiriman Pesan', link: '/fitur-khas/saluran' },
          { text: 'Rumus Cepat Panah (=>)', link: '/fitur-khas/fungsi-ringkas' },
          { text: 'Membongkar Isi Paket (Destrukturisasi)', link: '/fitur-khas/destrukturisasi' },
          { text: 'Menyisipkan Nilai ke Kalimat', link: '/fitur-khas/interpolasi-teks' },
          { text: 'Penyaluran Berantai (|>)', link: '/fitur-khas/pipeline' },
          { text: 'Keamanan Nilai Kosong & Kenyamanan', link: '/fitur-khas/ergonomi-dan-null-safety' },
          { text: 'Pekerjaan Bersamaan (sambil & bersamaan)', link: '/fitur-khas/konkurensi' }
        ]
      },
      {
        text: '🎨 Gaya Bahasa & Penulisan',
        items: [
          { text: 'Gaya Bahasa Fleksibel (Multi-Dialek)', link: '/gaya-sintaks/multi-dialek' },
          { text: 'Mengubah Gaya di Tengah Kode', link: '/gaya-sintaks/in-code-styling' },
          { text: 'Membuat Gaya Bahasa Sendiri (.style)', link: '/gaya-sintaks/berkas-style' }
        ]
      },
      {
        text: '📦 Alat Bantu Bawaan',
        items: [
          { text: 'Alat Perhitungan (matematika)', link: '/pustaka-standar/matematika' },
          { text: 'Pengolah Teks (teks)', link: '/pustaka-standar/teks' },
          { text: 'Penunjuk Waktu (waktu)', link: '/pustaka-standar/waktu' },
          { text: 'Pengelola Berkas (berkas)', link: '/pustaka-standar/berkas' },
          { text: 'Pertukaran Data (json)', link: '/pustaka-standar/json' }
        ]
      },
      {
        text: '⚙️ Cara Kerja di Balik Layar',
        items: [
          { text: 'Pohon Logika Murni (U-AST)', link: '/arsitektur/uast' },
          { text: 'Mesin Eksekusi Cepat (Bytecode VM)', link: '/arsitektur/bytecode-vm' }
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
