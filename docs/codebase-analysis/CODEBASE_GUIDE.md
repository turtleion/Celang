# Codebase Guide — Celang

Where to look, and how to get there from an entry point.

## Top-Level Layout

```
Celang/
  cmd/celang/main.go          CLI dispatcher + all subcommands
  cmd/celang-wasm/main.go     WASM/browser entry
  pkg/parser/                 lexer + Pratt parser + universal gate
  pkg/uast/                   dialect-free tree (Program/Statement/Expression/serde)
  pkg/syntax/                 StyleProfile + Formatter (not used by lexer/parser)
  pkg/compiler/               U-AST → bytecode Chunk
  pkg/bytecode/               Chunk, Opcode, CompiledFunction
  pkg/vm/                     stack VM (default executor)
  pkg/evaluator/              tree-walk executor + objects + environment + builtins + stdlib + FFI + modules
  pkg/empathy/                diagnostics rendering
  pkg/repl/                   interactive shell
  pkg/tester/                 `celang uji/test` runner
  pkg/repl/  pkg/tester/  pkg/empathy/  pkg/syntax/  pkg/uast/ etc.
  examples/*.ce               runnable samples; jawa.style is a profile example
  editors/vscode/             TextMate grammar + snippets
  docs/                       VitePress site (Indonesian)
  Makefile  go.mod  package.json  tes.ce
```

Total Go: ~36k LOC across ~40 files (`wc -l pkg/**/*.go` confirms).

## Entry Points

| Goal | Start here | Function |
|------|-----------|----------|
| Run a file via VM (default) | `cmd/celang/main.go:96` `case "run"` → `runFile()` | `runFile` reads source, `parser.NewLexer`/`NewParser`→`ParseProgram`, then `compiler.NewCompiler().Compile`→`vm.NewWithEnv(chunk, env).Run()`; if `--ast`, calls `evaluator.Eval` instead |
| Run via evaluator | same `run` path with `--ast`/`--legacy`/`--tree-walk` flag | `evaluator.Eval(program, env)` |
| REPL | `cmd/celang/main.go:38` (no args) or `case "repl"` | `repl.Start(Options{Style, In, Out})` — loop read→lex→parse→(compile→VM or eval)→print |
| Test suite | `case "uji","test"` → `tester.Run` | `pkg/tester/runner.go:Run` discovers `*.ce` and `uji` blocks, `evaluator.Eval` each, `harus_sama` etc. as builtins |
| Format | `case "format"` → `formatFile` | `syntax.LoadProfile` (or builtin names `indonesian`/`jawa`/`pythonic`/`javascript`) → `syntax.Formatter{profile}.Format(program)` |
| Compile to AST JSON | `case "compile","parse"` → `compileFile` | `uast.ToJSON(program)` to `file.ast` |
| Disassemble | `case "disasm"` → `disasmFile` | `compiler.Compile` → `chunk.Disassemble` |
| Standalone bundle | `case "bangun","build"` → `buildFile` | `resolveRunnerPath` + `gzip` payload + trailer `CELANG_BUNDLE_OK`; `tryRunEmbedded()` reverses at startup |
| WASM playground | `cmd/celang-wasm/main.go` | JS bridge, imports same `parser`/`evaluator` packages, stubs FFI |

## Finding a Feature

| Feature | Where it is parsed | Where it is executed |
|---------|-------------------|---------------------|
| Variables `buat`/`konstanta`/`mungkin` + time travel `.awal` | `parser.go:parseVariableDeclaration` (908) | `evaluator.go:Eval VariableDeclaration` + `environment.go:VariableEntry/Assign` |
| `rujuk` / `*ptr` pointers | `parser.go:parseReferenceDeclaration` (1066), `token REF_DECL` | `evaluator.go:ReferenceDeclaration`, `object.go:POINTER_OBJ`, deref in `evalMemberAccess`/`evalAssignment` |
| `bentuk` struct + `privat`/`privat!` + `@` | `parseStructDeclaration` (1008), `parseStructBody` (1089) | `evaluator.go:StructDeclaration` case, `object.go:STRUCT_DEF_OBJ/STRUCT_INSTANCE_OBJ`, visibility checked via `isInternalAccess` |
| Control flow `jika`/`selama`/`untuk dalam`/`pilih` | `parseIfStatement` (1361), `parseWhileStatement` (1439), `parseForEachStatement` (1475), `parseSwitchStatement` (1273) | `evaluator.go` switch cases; `compiler.go:compileStatement` emits `OpJumpIfFalse`/`OpJump`/`OpLoop`/`OpGetIter`/`OpForIter` |
| `berhenti`/`lanjut` | `parseBreakStatement` (1665), `parseContinueStatement` (1676) | VM emits `OpJump` to loop exit/start; evaluator returns `BREAK_SIGNAL`/`CONTINUE_SIGNAL` |
| `coba`/`tangkap`/`akhirnya`/`lempar` | `parseTryCatchStatement` (642), `parseThrowStatement` (703) | `evaluator.go:TryCatchStatement` — catch binds error, finally always runs |
| `saat_selesai`/`defer` | `token DEFER`, parsed in `parseStatement` | `environment.go:PushDefer/PopAllDefers`, flushed on function return/block exit in `evaluator.go` |
| Pipeline `\|>` | `universal.go` not needed; `token PIPE_ARROW`, infix `parsePipelineExpression` | `evaluator.go:PipelineExpression` and `compiler.go:compilePipeline` |
| Interpolation `{...}` / `${...}` | `lexer.go:NextToken` string branch + `parseInterpolationOrLiteral` (1902) | `evaluator.go:TemplateLiteral` eval, `vm.go` string concat |
| `?.` `??`/`kalau kosong` | `token OPTIONAL_DOT`/`NULL_COALESCE`, Lexer handles `?.` and word-compounds | `evaluator.go:evalSafeNavigation`, `NullCoalesceExpression` |
| `..` range & slicing | `token RANGE`, `parseRangeExpression` | `evaluator.go:evalRangeExpression` + `evalArraySlice`/`evalStringSlice`/`evalBytesSlice` |
| `fungsi` / `=>` | `parseFunctionDeclaration` (1512), `parseFunctionParameters` (1561) | `object.go:FUNCTION_OBJ`, `evaluator.go:ApplyFunction`, `compiler.go` local slots; `@` sugar in struct methods |
| Generics `<T harus Angka>` | `parseTypeHint` (1618), `HITIP`/`HARUS` tokens | `evaluator.go:applyFunctionWithGenerics` validates `harus` constraints at call time |
| `titip`/`jemput`/`saluran` | `parseConcurrencyStatement` (1687), tokens `TITIP`/`JEMPUT`/`SALURAN` | `concurrency_modern.go:Ticket/Channel`, `evaluator.go:TitipExpression/JemputExpression` |
| `transaksi`/`batalkan`/`reaktif`/`pantau` | `parseTransactionStatement` (716), `parseReactiveDeclaration` (778) | `environment.go:Begin/Commit/RollbackTransaction`, `DefineReactive`/`triggerReactivity` |
| `gaya "..." { }` | `parseStyleBlockStatement` (879), token `STYLE_BLOCK` | `evaluator.go:StyleBlockStatement` — inner body evaluated in same scope; `formatter.go` switches profile per block |
| `ambil`/`impor` | `parseImportStatement` (443), `parsePythonicImportStatement` (569) | `evaluator/module.go:LoadModule`, `StdlibNames` map in `stdlib.go` |
| `uji "name" { harus_sama }` | `parseTestStatement` (852), token `TEST` | `evaluator.go:evalTestStatement` records `TestCase`; `tester/runner.go:Run` drives whole files |
| FFI `ffi.buka` etc. | `stdlib_ffi.go` registers builtins | `ffi_loader_{unix,windows,wasm}.go` dispatch |
| Time-travel `.sebelumnya`/`.awal`/`.histori` | not parsed specially — member access | `evaluator.go:evalMemberAccess` reads `VariableEntry.History` |
| Culture suffixes `rupiah`/`ribu`/`juta` etc. | not tokens — evaluated as property access | `culture_helpers.go:evalNumericMemberWithEnv` + `evaluator.go:evalMemberAccess` numeric branch, gated by `env.IsUniversal` |

## Which Files Are Central?

- Start with `pkg/parser/lexer.go` + `token.go` + `universal.go` to understand the vocabulary.
- Then `pkg/uast/node.go` + `statement.go` + `expression.go` for the shape everything else depends on.
- Then `pkg/evaluator/evaluator.go` (2671 LOC) — the reference implementation of semantics. The compiler/VM reimplement the same semantics differently, so reading evaluator first avoids duplicating mental work.
- Then `pkg/compiler/compiler.go` + `pkg/vm/vm.go` side-by-side to see the VM mapping.
- Then `pkg/evaluator/environment.go` + `object.go` for state ownership.
- For tooling, `cmd/celang/main.go` + `pkg/empathy/diagnostics.go` + `pkg/tester/runner.go`.

## Tracing Tips

- To answer “what does syntax X do?” — grep `pkg/parser/parser.go` for `parseX`, then find the same `NodeType` case in `pkg/evaluator/evaluator.go:Eval` and `pkg/compiler/compiler.go:compileStatement`.
- To answer “how does value Y behave at runtime?” — check `pkg/evaluator/object.go` for the type, and `pkg/evaluator/evaluator.go:evalBinaryOp` / `pkg/vm/vm.go:evalAdd` for operator semantics.
- To answer “why is error message Z shown?” — start from `pkg/empathy/diagnostics.go:NewRuntimeDiagnostic` / `NewSyntaxDiagnosticFromError` and search callers `diagnostics.Render()` in `cmd/celang/main.go:runFile`.
- To answer “which stdlib function exists?” — read `pkg/evaluator/stdlib.go:StdlibNames` for module registry, then the per-module file (`stdlib_math.go` etc.) for function list; `definePublicFunction` calls are the manifests.

## Tests as Navigation Aids

| Package | Test file | What it proves |
|---------|-----------|----------------|
| `pkg/parser` | `parser_test.go`, `fuzz_test.go` | Lexing + parsing of keywords, compounds, interpolation, errors |
| `pkg/uast` | `serde_test.go` | JSON round-trip (U-AST ↔ `.ast`) lossless |
| `pkg/syntax` | `profile_test.go`, `formatter_test.go` | `.style` load (flat + categorized) and Indo↔Pythonic formatting |
| `pkg/evaluator` | `evaluator_test.go`, `collections_test.go`, `module_test.go`, `stdlib_test.go`, `universal_test.go` | Semantics, collections, modules, stdlib, universal gate |
| `pkg/compiler` | `compiler_test.go` | U-AST → Chunk opcode sequences |
| `pkg/bytecode` | `chunk_test.go` | Chunk write/disassemble |
| `pkg/vm` | `vm_test.go`, `collections_test.go`, `stress_parity_test.go` | VM execution and parity with evaluator |
| `pkg/empathy` | `diagnostics_test.go` | Rendered diagnostic boxes (ID + EN) |
| `pkg/repl` | `repl_test.go` | REPL line handling and multi-line |
| `pkg/tester` | `runner_test.go` | Discovery + `harus_sama` + filter |
