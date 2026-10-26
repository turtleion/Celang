# Project Overview — Celang

## What Is Celang?

Celang is a human-centric programming language implemented in Go. Source files use the `.ce` extension. The compiler and virtual machine ship as a single binary (`celang`) that runs programs, formats code, compiles to an intermediate tree, and packages standalone executables. A WebAssembly build (`cmd/celang-wasm`) powers a browser playground.

The central idea is **Syntax Styling**: the logic of a program (its tree structure) is stored separately from how it looks on screen. One program can be displayed in Indonesian keywords (`buat`, `jika`, `fungsi`), in Sunda (`jieun`, `lamun`), or in English/JS/Python style without changing what it does. The tree is called the **Universal AST (U-AST)** and is serialised as JSON (`.ast` files).

## Primary Goals

1. **Lower the language barrier** — beginners write with Indonesian or regional keywords they already understand, not English jargon.
2. **End style arguments** — teams each view the same file in their preferred dialect; formatting is a projection, not a rewrite.
3. **Friendly errors** — diagnostics show a framed box with the exact line, an arrow, and a suggestion in Indonesian (or English in universal mode).
4. **Modern ergonomics** — null safety (`?.`, `kalau kosong`), pipelines (`|>`), transactions, reactivity, concurrency (`titip`/`jemput`, `saluran`), generics (`<T harus Angka>`), and destructuring, all on a fast bytecode VM.

Evidence: `README.md`, `ROADMAP.md`, `cmd/celang/main.go` dispatch table, `pkg/empathy/diagnostics.go`.

## Current Capabilities (Verified)

All items below pass `go test ./...` (10 packages, ~203 tests) and run via `celang run`.

- **Variables**: `buat` (mutable), `konstanta`/`tetep` (immutable), `mungkin` (nullable), `rujuk` (reference/pointer `*`), `reaktif` (reactive cell). Time-travel accessors `.sebelumnya` / `.awal` / `.histori` are implemented on `VariableEntry` with a bounded ring buffer of 50 entries (`pkg/evaluator/environment.go:12`).
- **Control flow**: `jika`/`selain jika`/`selain`, `selama`, `untuk .. dalam`, `pilih`/`bila`/`lainnya`, `berhenti`/`lanjut`, `coba`/`tangkap`/`akhirnya`/`lempar`.
- **Functions**: `fungsi` with `kembalikan`, fat-arrow shorthand `=>`, closures, recursion, `saat_selesai`/`defer`.
- **Collections**: arrays, maps, bytes (`bita`), slicing with `..` ranges, ~40 collection helpers in `pkg/evaluator/collections.go`.
- **Data modelling**: `bentuk` (struct) with fields, methods, `privat`, embed, `@` shorthand, operator overloading.
- **Null safety**: `?.` safe navigation, `??` / `kalau kosong` / `bila_kosong` coalescence.
- **Pipeline**: `|>` and string interpolation `{...}` / `${...}`.
- **Types**: `tipe` aliases, generics `<T harus Angka>` / `sembarang`/`any`, `itu` type annotation, cast checked at runtime.
- **Concurrency**: `titip` (spawn -> `Ticket`), `jemput`/`jemput_semua`/`jemput_pertama`, `saluran` (channel with `.kirim`/`.terima`/`.tutup`).
- **Transactions & reactivity**: `transaksi` / `batalkan` / `pulihkan`, `reaktif`, `pantau`/`watch`.
- **Syntax tools**: `gaya "nama" { ... }` in-code block, `celang format --style=...`, custom `.style` profiles (e.g., `examples/sunda.style`), `celang compile` to `.ast`, `celang disasm`.
- **Stdlib + FFI**: `matematika`/`math`, `teks`/`strings`, `waktu`/`time`, `berkas`/`fs`, `json`, `sistem`/`os`, `http`/`jaringan`/`net`, `ffi` (dynamic C ABI via `purego`, platform-split loaders).
- **Testing**: `celang uji` / `celang test` runner (`pkg/tester/runner.go`) plus embedded `uji "nama" { harus_sama(...) }` blocks; `celang bangun`/`build` produces a self-contained executable (runner + gzipped payload + trailer).
- **Editor**: `editors/vscode/` extension with `celang.tmLanguage.json` and snippets.

Confidence: **Verified** (tests + manual `celang run` on `examples/*.ce`).

## Technology Stack

| Layer | Choice | Evidence |
|-------|--------|----------|
| Language | Go 1.27.1 | `go.mod` |
| Binary | Single `celang` from `cmd/celang` | `Makefile` `CMD_PATH := ./cmd/celang` |
| WASM | `cmd/celang-wasm` (js/wasm build) | `ffi_loader_wasm.go`, `stdlib_ffi_wasm.go` |
| Stdlib ext | `github.com/ebitengine/purego v0.11.2` for FFI | `go.mod` |
| Docs | VitePress 1.6.4 | `package.json` |
| Tests | `go test` + custom `tester` package | `pkg/tester/runner.go` |

## Execution Model

Two interchangeable executors consume the same U-AST:

1. **Bytecode VM (default)** — `pkg/compiler` lowers U-AST to a `bytecode.Chunk` (instructions + constant pool + line table), then `pkg/vm` executes it on a stack machine (256 call frames, 2048 stack slots). Selected by `celang run file.ce`.
2. **Tree-walk evaluator** — `pkg/evaluator.Eval` walks U-AST directly against an `Environment`. Selected by `celang run --ast file.ce` (also `--legacy`/`--tree-walk`).

Both share the same `Environment`, `Object` type system, builtins, and stdlib. Behavioural parity is guarded by `pkg/vm/stress_parity_test.go`.

## Implementation Status

- Passes `go test -count=1 ./...` on Linux on every commit checked (bytecode, compiler, empathy, evaluator, parser, repl, syntax, tester, uast, vm).
- Native binary builds with `make build` (`go build -o celang ./cmd/celang`), ~13 MB on Linux.
- Examples in `examples/` run end-to-end (verified by running each via `celang run`).
- Web demo build is present but not exercised in CI here; requires `GOOS=js GOARCH=wasm`.

## What Celang Does Not Currently Do

- **No package manager** — `ROADMAP.md` lists this as a planned future item; `celang.mod` today is only a module-root marker, not a dependency manifest.
- **No ahead-of-time native binary** — `celang bangun` embeds the runner + source in a self-extracting binary; it is not a native-to-machine compiler.
- **No Language Server Protocol** — listed under "Rencana Masa Depan" in `ROADMAP.md`; `position` in U-AST nodes exists (`pkg/uast/node.go`) but no LSP server ships.
- **No static type checking at compile time** — type hints are enforced at runtime via `castValueToType`; generics constraints (`harus`) are validated at call time.
- **No separate semantic-analysis pass** — name resolution and scope are handled during evaluation/compilation; there is no standalone analyser/typer.
