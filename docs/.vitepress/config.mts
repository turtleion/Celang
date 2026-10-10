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
            match: '\\b(jika|selain jika|selain|selama|untuk|dalam|kembalikan|berhenti|lanjut|sambil|bersamaan|mundur ke|lempar|tangkap|gaya|if|else|while|for|return|break|continue|async|parallel|lamun|salami|pikeun|balikkeun)\\b'
          },
          {
            name: 'keyword.declaration.celang',
            match: '\\b(buat|konstanta|fungsi|bentuk|itu|let|const|var|function|def|jieun|tetep|struct)\\b'
          },
          {
            name: 'constant.language.celang',
            match: '\\b(benar|salah|kosong|true|false|null|None|bener)\\b'
          },
          {
            name: 'support.function.celang',
            match: '\\b(cetak|masukan|panjang|tipe|rentang|print|input|len|type|range)\\b'
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
      { text: 'Fitur Unggulan', link: '/fitur-khas/time-travel' },
      { text: 'Syntax Styling', link: '/gaya-sintaks/multi-dialek' },
      { text: 'Stdlib API', link: '/pustaka-standar/matematika' },
      { text: 'Bytecode VM', link: '/arsitektur/bytecode-vm' }
    ],
    sidebar: [
      {
        text: '🚀 Panduan Memulai',
        items: [
          { text: 'Filosofi & Pengenalan', link: '/panduan/pengenalan' },
          { text: 'Mulai Cepat & Instalasi', link: '/panduan/mulai-cepat' },
          { text: 'CLI & Interactive REPL', link: '/panduan/cli-repl' }
        ]
      },
      {
        text: '🧠 Dasar Bahasa',
        items: [
          { text: 'Variabel & Tipe Data', link: '/dasar/variabel-tipe' },
          { text: 'Kontrol Alur (Percabangan & Loop)', link: '/dasar/kontrol-alur' },
          { text: 'Fungsi & Fungsi Anonim', link: '/dasar/fungsi' }
        ]
      },
      {
        text: '✨ Fitur Khas Celang',
        items: [
          { text: 'Variabel Time-Travel', link: '/fitur-khas/time-travel' },
          { text: 'Bentuk (Struct) & Metode', link: '/fitur-khas/bentuk-struct' },
          { text: 'Destrukturisasi (Array & Map)', link: '/fitur-khas/destrukturisasi' },
          { text: 'String Interpolation', link: '/fitur-khas/interpolasi-teks' },
          { text: 'Pipeline Operator (|>)', link: '/fitur-khas/pipeline' },
          { text: 'Konkurensi Alami (sambil & bersamaan)', link: '/fitur-khas/konkurensi' }
        ]
      },
      {
        text: '🎨 Syntax Styling',
        items: [
          { text: 'Multi-Dialek (ID, Sunda, JS, Pythonic)', link: '/gaya-sintaks/multi-dialek' },
          { text: 'In-Code Styling Block', link: '/gaya-sintaks/in-code-styling' },
          { text: 'Membuat Berkas .style Kustom', link: '/gaya-sintaks/berkas-style' }
        ]
      },
      {
        text: '📦 Pustaka Standar (Stdlib)',
        items: [
          { text: 'Modul Matematika (math)', link: '/pustaka-standar/matematika' },
          { text: 'Modul Teks (strings)', link: '/pustaka-standar/teks' },
          { text: 'Modul Waktu (time)', link: '/pustaka-standar/waktu' },
          { text: 'Modul Berkas (fs)', link: '/pustaka-standar/berkas' },
          { text: 'Modul JSON', link: '/pustaka-standar/json' }
        ]
      },
      {
        text: '⚙️ Arsitektur Internal',
        items: [
          { text: 'Universal AST (U-AST) & JSON Serde', link: '/arsitektur/uast' },
          { text: 'Bytecode Compiler & Virtual Machine', link: '/arsitektur/bytecode-vm' }
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
