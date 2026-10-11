<template>
  <div class="celang-playground">
    <div class="playground-header">
      <div class="header-left">
        <label for="preset-select" class="control-label">Pilih Contoh:</label>
        <select id="preset-select" v-model="selectedPreset" @change="loadPreset" class="select-control">
          <option v-for="(code, name) in presets" :key="name" :value="name">
            {{ name }}
          </option>
        </select>
      </div>

      <div class="header-right">
        <div class="engine-toggle">
          <label class="engine-label">
            <input type="radio" value="vm" v-model="engineMode" />
            <span>⚡ Bytecode VM</span>
          </label>
          <label class="engine-label">
            <input type="radio" value="ast" v-model="engineMode" />
            <span>🌳 AST Mode</span>
          </label>
        </div>

        <button class="run-button" :disabled="!wasmReady || isRunning" @click="executeCode">
          <span v-if="isRunning">⏳ Menjalankan...</span>
          <span v-else>▶ Jalankan Kode</span>
        </button>
      </div>
    </div>

    <div class="playground-body">
      <div class="editor-pane">
        <div class="pane-title">
          <span>KODE SUMBER CELANG</span>
          <span class="wasm-badge" :class="{ ready: wasmReady }">
            {{ wasmStatus }}
          </span>
        </div>
        <textarea
          v-model="sourceCode"
          class="code-textarea"
          spellcheck="false"
          placeholder="Tulis kode Celang Anda di sini..."
          @keydown.tab.prevent="insertTab"
        ></textarea>
      </div>

      <div class="terminal-pane">
        <div class="pane-title">
          <span>TERMINAL OUTPUT</span>
          <span v-if="executionTime !== null" class="exec-time">
            ⏱ {{ executionTime }} ms ({{ engineMode.toUpperCase() }})
          </span>
        </div>
        <div class="terminal-output" :class="{ error: hasError }">
          <pre v-if="outputLog">{{ outputLog }}</pre>
          <pre v-else class="placeholder-text">// Tekan 'Jalankan Kode' untuk melihat hasil eksekusi...</pre>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'

const wasmReady = ref(false)
const wasmStatus = ref('Memuat Engine WASM...')
const isRunning = ref(false)
const engineMode = ref('vm')
const executionTime = ref(null)
const outputLog = ref('')
const hasError = ref(false)

const presets = {
  '1. Halo Dunia & Variabel': `// Selamat datang di Celang Web Playground!
buat nama = "Pengembang Celang"
buat angkaFavorit = 42

cetak("Halo,", nama, "!")
cetak("Angka favorit Anda:", angkaFavorit)
cetak("Perhitungan kuadrat:", angkaFavorit * angkaFavorit)
`,
  '2. Variabel Time-Travel (State History)': `// Variabel Time-Travel mencatat riwayat perubahan nilai
buat saldo = 100000
saldo = 250000
saldo = 500000

cetak("Saldo Saat Ini : Rp", saldo)
cetak("Saldo Sebelumnya: Rp", saldo.sebelumnya)
cetak("Saldo Awal Mula : Rp", saldo.awal)
cetak("Seluruh Histori :", saldo.histori)
`,
  '3. Bentuk (Struct) & Metode': `// Deklarasi Bentuk (Struct) dengan metode 'diri'
bentuk RekeningBank {
    pemilik: Teks,
    saldo itu Angka = 0,

    fungsi setor(jumlah: Angka) {
        diri.saldo += jumlah
        kembalikan diri.saldo
    }

    fungsi tarik(jumlah: Angka) {
        diri.saldo -= jumlah
        kembalikan diri.saldo
    }
}

buat akun = RekeningBank("Budi Santoso", 50000)
akun.setor(150000)
akun.tarik(30000)

cetak("Pemilik Rekening:", akun.pemilik)
cetak("Saldo Akhir:", akun.saldo)
`,
  '4. Pipeline Operator (|>)': `// Aliran data fungsional dengan operator |>
fungsi tambahSatu(x) { kembalikan x + 1 }
fungsi kaliDua(x) { kembalikan x * 2 }
fungsi kuadrat(x) { kembalikan x * x }

buat hasil = 5 |> tambahSatu |> kaliDua |> kuadrat
cetak("Hasil pipa ekspresi 5 |> +1 |> *2 |> kuadrat =", hasil)
`,
  '5. Destrukturisasi (Array & Map)': `// Destrukturisasi praktis
buat koordinat = [120, 480]
buat [x, y] = koordinat
cetak("Koordinat X:", x, "| Y:", y)

buat profil = { "nama": "Siti", "kota": "Bandung" }
buat { nama, kota } = profil
cetak("Pengguna:", nama, "dari kota:", kota)
`,
  '6. Komputasi Cepat Fibonacci': `// Pengujian rekursi fungsi di Bytecode VM
fungsi fibonacci(n) {
    jika n <= 1 {
        kembalikan n
    }
    kembalikan fibonacci(n - 1) + fibonacci(n - 2)
}

cetak("Menghitung Fibonacci ke-15...")
buat hasil = fibonacci(15)
cetak("Hasil:", hasil)
`
}

const selectedPreset = ref('1. Halo Dunia & Variabel')
const sourceCode = ref(presets['1. Halo Dunia & Variabel'])

function loadPreset() {
  sourceCode.value = presets[selectedPreset.value] || ''
}

function insertTab(e) {
  const target = e.target
  const start = target.selectionStart
  const end = target.selectionEnd
  sourceCode.value = sourceCode.value.substring(0, start) + '    ' + sourceCode.value.substring(end)
  setTimeout(() => {
    target.selectionStart = target.selectionEnd = start + 4
  }, 0)
}

function executeCode() {
  if (!wasmReady.value || typeof window.runCelang !== 'function') {
    outputLog.value = 'Engine WASM belum siap!'
    hasError.value = true
    return
  }

  isRunning.value = true
  hasError.value = false

  try {
    const res = window.runCelang(sourceCode.value, engineMode.value, 'indonesian')
    executionTime.value = res.durationMs

    if (res.error) {
      hasError.value = true
      outputLog.value = (res.output ? res.output + '\n' : '') + '❌ ' + res.error
    } else {
      let finalOut = res.output
      if (res.result && res.result !== 'kosong') {
        finalOut += (finalOut ? '\n' : '') + '=> ' + res.result
      }
      outputLog.value = finalOut || '(Program selesai tanpa luaran output)'
    }
  } catch (err) {
    hasError.value = true
    outputLog.value = '❌ Kesalahan tak terduga: ' + err.message
  } finally {
    isRunning.value = false
  }
}

onMounted(() => {
  if (typeof window === 'undefined') return

  if (window.celangWasmReady && typeof window.runCelang === 'function') {
    wasmReady.value = true
    wasmStatus.value = 'Engine WASM Siap ✨'
    return
  }

  // Load wasm_exec.js dynamically
  const script = document.createElement('script')
  script.src = '/wasm_exec.js'
  script.onload = () => {
    if (typeof window.Go === 'undefined') {
      wasmStatus.value = 'Gagal memuat Go runtime!'
      return
    }

    const go = new window.Go()
    WebAssembly.instantiateStreaming(fetch('/celang.wasm'), go.importObject)
      .then((result) => {
        go.run(result.instance)
        wasmReady.value = true
        wasmStatus.value = 'Engine WASM Siap ✨'
      })
      .catch((err) => {
        console.error('WASM Init Error:', err)
        wasmStatus.value = 'Gagal inisialisasi WASM'
      })
  }
  script.onerror = () => {
    wasmStatus.value = 'Gagal memuat wasm_exec.js'
  }
  document.head.appendChild(script)
})
</script>

<style scoped>
.celang-playground {
  display: flex;
  flex-direction: column;
  border: 1px solid var(--vp-c-divider);
  border-radius: 12px;
  overflow: hidden;
  background-color: var(--vp-c-bg-soft);
  margin: 1.5rem 0;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
}

.playground-header {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: space-between;
  padding: 0.75rem 1rem;
  background-color: var(--vp-c-bg-mute);
  border-bottom: 1px solid var(--vp-c-divider);
  gap: 0.75rem;
}

.header-left, .header-right {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.control-label {
  font-size: 0.85rem;
  font-weight: 600;
  color: var(--vp-c-text-2);
}

.select-control {
  padding: 0.35rem 0.75rem;
  font-size: 0.85rem;
  border-radius: 6px;
  border: 1px solid var(--vp-c-divider);
  background: var(--vp-c-bg);
  color: var(--vp-c-text-1);
}

.engine-toggle {
  display: flex;
  background: var(--vp-c-bg);
  padding: 0.2rem;
  border-radius: 6px;
  border: 1px solid var(--vp-c-divider);
  gap: 0.5rem;
}

.engine-label {
  display: flex;
  align-items: center;
  font-size: 0.8rem;
  gap: 0.3rem;
  cursor: pointer;
  padding: 0.15rem 0.4rem;
  border-radius: 4px;
}

.run-button {
  background: var(--vp-c-brand-1);
  color: #ffffff;
  padding: 0.4rem 1.1rem;
  font-size: 0.85rem;
  font-weight: 600;
  border-radius: 6px;
  cursor: pointer;
  border: none;
  transition: opacity 0.2s;
}

.run-button:hover:not(:disabled) {
  opacity: 0.9;
}

.run-button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.playground-body {
  display: grid;
  grid-template-columns: 1fr 1fr;
  min-height: 420px;
}

@media (max-width: 768px) {
  .playground-body {
    grid-template-columns: 1fr;
  }
}

.editor-pane, .terminal-pane {
  display: flex;
  flex-direction: column;
}

.editor-pane {
  border-right: 1px solid var(--vp-c-divider);
}

.pane-title {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.4rem 0.8rem;
  font-size: 0.75rem;
  font-weight: 700;
  letter-spacing: 0.05em;
  color: var(--vp-c-text-2);
  background: var(--vp-c-bg-alt);
  border-bottom: 1px solid var(--vp-c-divider);
}

.wasm-badge {
  font-size: 0.7rem;
  padding: 0.15rem 0.5rem;
  border-radius: 12px;
  background: #f59e0b;
  color: #fff;
}

.wasm-badge.ready {
  background: #10b981;
}

.exec-time {
  font-size: 0.75rem;
  color: var(--vp-c-brand-1);
}

.code-textarea {
  flex: 1;
  width: 100%;
  padding: 1rem;
  font-family: var(--vp-font-family-mono);
  font-size: 0.88rem;
  line-height: 1.5;
  background: var(--vp-c-bg);
  color: var(--vp-c-text-1);
  border: none;
  resize: none;
  outline: none;
}

.terminal-output {
  flex: 1;
  padding: 1rem;
  background: #18181b;
  color: #f4f4f5;
  font-family: var(--vp-font-family-mono);
  font-size: 0.85rem;
  overflow-y: auto;
  line-height: 1.5;
}

.terminal-output pre {
  margin: 0;
  white-space: pre-wrap;
  word-break: break-all;
}

.terminal-output.error {
  color: #f87171;
}

.placeholder-text {
  color: #71717a;
}
</style>
