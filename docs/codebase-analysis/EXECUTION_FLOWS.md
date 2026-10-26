# Execution Flows — Celang

All flows start from `cmd/celang/main.go:main()`.

## A. Program Execution — VM Path (Default: `celang run file.ce`)

```
user: celang run file.ce
  main():32  tryRunEmbedded() → false (no bundle trailer)
  main():96  case "run" → runFile(filePath, scriptArgs, useAST=false)
    runFile(): read file bytes → source string
    if suffix .ast/.json → uast.FromJSON(bytes) → program
    else
      l := parser.NewLexer(source)         // lexer.go:NewLexer — DetectUniversalDirective stored
      p := parser.NewParser(l)             // parser.go:NewParser — registers prefix/infix parsers
      program, err := p.ParseProgram()     // loops parseStatement until EOF
      if err != nil (_ ParseErrorList)    → empathy.NewSyntaxDiagnosticFromError(err, source).Render() → stderr, exit 1
    env := evaluator.NewEnvironment()      // environment.go — inherits IsUniversal from DetectUniversalDirective
    env.IsUniversal = parser.DetectUniversalDirective(source)
    evaluator.RegisterBuiltins(env, os.Stdout)
    evaluator.RegisterStdlibGlobals(env)   // stdlib.go — installs berkas/math/teks/waktu/json/os/http/ffi
    evaluator.SetScriptArgs(scriptArgs)    // stdlib_system.go — for sistem.args
    if useAST
      result := evaluator.Eval(program, env)
    else
      compiler := compiler.NewCompilerWithEnv(env)   // compiler.go:NewCompilerWithEnv
      chunk, _ := compiler.Compile(program)           // walk U-AST, emit bytecodes + constants
      vm := vm.NewWithEnv(chunk, env)                // vm.go:NewWithEnv — wraps chunk in root CompiledFunction
      result, err := vm.Run()                        // stack machine dispatch loop
      // fallback inside Run: OpEvalAST → evaluator.Eval(node, env)
    if result.Type()==ERROR_OBJ
      diag := empathy.NewRuntimeDiagnostic(result.(*Error), source)
      fmt.Fprintln(os.Stderr, diag.Render()); os.Exit(1)
    // if program is bundle (tryRunEmbedded at startup), same path but source came from gzipped trailer
```

Chunk/constants sharing: `evaluator.Object` values used as `Chunk.Constants` entries, so the compiler reuses the same type system as the evaluator; `ASTNodeWrapper{Node}` allows a chunk to embed an unevaluated node for `OpEvalAST`.

`--ast` path differs only at the executor choice; parsing, env, builtins, and diagnostics are identical.

Sequence (Mermaid):

```mermaid
sequenceDiagram
  participant U as user
  participant CLI as cmd/celang/main.go:runFile
  participant Lex as parser.Lexer
  participant Par as parser.Parser
  participant UAST as uast.Program
  participant Comp as compiler.Compiler
  participant Chunk as bytecode.Chunk
  participant VM as vm.VM
  participant Env as evaluator.Environment
  U->>CLI: celang run file.ce
  CLI->>Lex: NewLexer(source)
  Lex-->>CLI: IsUniversal flag
  CLI->>Par: NewParser(lex)
  Par->>UAST: ParseProgram()
  UAST-->>CLI: Program{Body, IsUniversal}
  CLI->>Env: NewEnvironment + RegisterBuiltins/Stdlib
  CLI->>Comp: NewCompilerWithEnv(env)
  Comp->>Chunk: Compile(program)
  Chunk-->>CLI: {Instructions, Constants, Lines}
  CLI->>VM: NewWithEnv(chunk, env)
  VM->>Env: Get/Set via activeVMEnv
  VM-->>CLI: Object result
  CLI->>U: Render / exit code
```

## B. Syntax Styling — From Profile to Text

There are two directions, and they do not mirror each other via a shared abstraction.

### Forward: source text → token → U-AST (lex/parse)

The lexer ignores profiles. Every dialect keyword is in `lexer.go:keywords` (Indo + Sunda + English). `NextToken()` maps the current word via `keywords[word]` or `IDENT`. Two-word compounds (`selain jika`→`ELSE_IF`, `kalau kosong`→`NULL_COALESCE`) are handled by a peek-ahead after emitting the first word. Numbers consume `.` for floats; strings handle `"`, `\\`, `\n`, and `${` interpolation holes by returning multiple `STRING` tokens. Result: a single token stream regardless of intended display dialect.

The parser is dialect-free: it switches on `TokenType` constants (`VAR_DECL`, `IF`, etc.), not on spelling. Every valid spelling parses identically.

### Reverse: U-AST → dialect text (formatter)

```
celang format file.ce --style=sunda|path/to/custom.style
  LoadProfile(nameOrPath)  // syntax/profile.go:UnmarshalJSON — flat OR categorized JSON
  l := parser.NewLexer(source) → p := parser.NewParser(l) → program := p.ParseProgram()
  fmt := syntax.NewFormatter(profile)
  out := fmt.Format(program)   // formatter.go — switches on NodeType, emits profile.VarDecl etc.
  if -w → write file else print
```

`StyleProfile` is only consulted by `Formatter`; `Lexer`/`Parser` never import `pkg/syntax`. Verified by `grep -r "syntax" pkg/parser pkg/uast --include="*.go"` → zero hits.

### In-code block `gaya "name" { ... }`

Tokenized as `STYLE_BLOCK "gaya"` + `STRING ""name""` + `LBRACE` → parsed by `parseStyleBlockStatement()` into `uast.StyleBlockStatement{StyleName, Body []Statement}`. At format time, `Formatter` recursively formats `Body` with the named profile; at eval time, `evaluator.go:StyleBlockStatement` just evaluates `Body` in the **same** `Environment` (no new scope), so variables inside are visible outside. Docs (`docs/gaya-sintaks/in-code-styling.md`) describe this correctly.

### Universal mode `"use universal"`

`DetectUniversalDirective(source)` scans first non-empty, non-comment line (`universal.go:10`). If matched, `Lexer.IsUniversal`, `Parser.IsUniversal`, `Program.IsUniversal`, `Environment.IsUniversal` are all set. Effects:
- `culture_helpers.go` blocks Indonesian suffix members (`rupiah`, `ribu`, suffix multipliers) when `env.IsUniversal`.
- `empathy/diagnostics.go` switches hint language to English.
- No change to lexing itself — Indonesian spellings still lex, but some runtime sugars are disabled.

## C. Error Handling

Three origins, one rendering pipeline:

```
Parse error
  parser.addError() → errors []string + parseErrors []ParseError{Line,Column,Message}
  ParseProgram() returns error = ParseErrorList (joined with "\n")
  runFile/repl → empathy.NewSyntaxDiagnosticFromError(err, source)
              → takes pel[0] as primary, len(pel)-1 as OtherErrorsCount
              → suggestSyntaxHint() (checks "lain jika", "maka", etc.)
  → Diagnostic.Render() → boxed "...\n Line N | source\n         ^ di sini\n 💡 Saran:..."

Runtime error (evaluator)
  evaluator.newError(meta, msg, hint) → *Error{Message, Pos(from meta), FriendlyHint}
  Each Eval branch checks isError(obj) and early-returns
  runFile → isError(result) → empathy.NewRuntimeDiagnostic(err, source) → Render()

Runtime error (VM)
  vm.Run() returns (Object, error); Go error from e.g. stack overflow
  plus evaluator-style *Error objects pushed as result and handled same as evaluator path
```

Recovery: parser attempts to continue after an error (peek expectations still advance), so one run can report multiple `[Baris N, Kolom M]` lines; diagnostics shows the first and counts the rest (`OtherErrorsCount`). Evaluator and VM do not recover inside a program — first `*Error` aborts the current block (`evalProgram`/`evalBlockStatement` checks `isError` after each statement).

Builtin assertion failures (`harus_sama` etc. in `builtins.go:registerTestingBuiltins`) return `*AssertionError{Message, Expected, Actual, Pos}` — handled like `*Error` but rendered with `Ekspektasi:/Aktual:` in `diagnostics.go`.

## D. Data and State

### Values

`Object` is the runtime value interface (`object.go:39`):

- Scalar: `Integer{int64}`, `Float{float64}`, `Boolean{bool}`, `String{string}`, `Bytes{[]byte}`, `Null`, `FFIPointer{uintptr}`.
- Composite: `Array{Elements []Object}`, `Map{Pairs map[string]Object, order []string}`, `StructDef`, `StructInstance`.
- Callable: `Function{Parameters []Parameter, Body *BlockStatement, Env *Environment}`, `Builtin{Fn func(...Object) Object}`, `CompiledFunction{Chunk, NumParams}` (in `bytecode/function.go`), all implement `Callable.Call`.
- Concurrency / control: `Ticket`, `Channel`, `Iterator`, signals (`ReturnValue`, `Error`, `AssertionError`, `BREAK_SIGNAL`, `CONTINUE_SIGNAL`, `ROLLBACK_SIGNAL`).

### Variables & Scope

`Environment` (`environment.go`) is a linked chain:

```
global env ─outer→ function env ─outer→ block env ─outer→ block env
  store: map[name]*VariableEntry{Value, IsMutable, Visibility, History ring[50], SourcePos}
```

- `Define(name, value, isMutable)` checks existing entry mutability, writes `VariableEntry`, records history (push to ring, evict oldest when >50), calls `RecordSavepointForVar` if inside a `transaksi`, and `triggerReactivity` for watchers.
- `Get(name)` walks `outer` chain; `GetEntry` returns `*VariableEntry` for history accessors.
- `NewEnclosedEnvironment(outer)` copies `FilePath`, `ModuleRoot`, `Registry`, `IsUniversal`.
- Function call in `evaluator.go:applyFunction` creates a new enclosed env, binds params, then `evalBlockStatement`.
- Thread safety: `sync.RWMutex` + `atomic` for `History` increments; goroutines from `titip` share the same outer env via pointer (not copied), so cross-goroutine mutation is racy by design (see Known Limitations).

Time-travel accessors `.sebelumnya`/`.awal`/`.histori` are not fields — they are intercepted in `evalMemberAccess`: if `obj` is a variable lookup result and `property` matches one of the three names, `GetEntry(name).History` is consulted instead of struct field access.

### Module & Visibility

`module.go:findModuleRoot` climbs directories looking for `celang.mod`/`modul.ce`; `LoadModule(path, env)` resolves relative to `env.ModuleRoot`, parses the target, evaluates it in a fresh env, then copies `public` (`publik!`) bindings into the caller. Visibility stored per `VariableEntry.Visibility` enum (`private`/`module`/`public`/`allowlist`).

### Transactions & Reactivity

`Environment` holds `txStack []*Transaction{savepoints map[name]Object, createdVars []string}`. `BeginTransaction`/`Commit`/`Rollback` are called from `evaluator.go:TransactionStatement`. `DefineReactive(name, computeFn)` registers a derived cell; `triggerReactivity(changedVar)` walks `reactiveDeps` graph (with `visited` set) and recomputes dependents, firing `Watch` callbacks.

## E. Build and Testing

### Build

```
Makefile:  build  →  go build -o celang ./cmd/celang
           test   →  go test -count=1 ./...
           docs-* →  npm run docs:build/dev/preview  (VitePress)
go.mod:    go 1.27.1, dep: purego v0.11.2
Binary:    ~13 MB on linux/amd64
```

`go test ./...` covers 10 packages; `cmd/celang` has no tests (CLI is integration-tested via examples).

### Test Execution

1. **Unit tests** — `go test ./...` via Makefile target `test`. ~203 cases across `bytecode`, `compiler`, `empathy`, `evaluator`, `parser`, `repl`, `syntax`, `tester`, `uast`, `vm`.
2. **In-language tests** — `celang uji [path] [-f filter]` / `celang test ...` (`pkg/tester/runner.go:Run`). Discovery: scan args for files/dirs, collect `*.ce` files, sort. Per file: parse→resolve imports→`evaluator.Eval`. Collection: top-level `uji "desc" { ... }` nodes become `TestCase{Name: desc}`. Builtins `harus_sama(actual, expected, msg?)` / `harus_benar` / `harus_salah` etc. set pass/fail via returned `AssertionError`. Runner aggregates `FileResult{FilePath, Cases, Duration, Error}` and prints colored `✓`/`✕` + summary `SUKSES/GAGAL`; exit 1 on any failure.
3. **Examples as smoke tests** — `examples/*.ce` are exercised manually via `celang run examples/<file>.ce`; no harness auto-runs them.

Confidence labels:
- VM default path, parsing, evaluator, and VM parity: **Verified** (tests + manual runs).
- Tester `harus_*` family completeness: **Code-supported** (exists in `builtins.go:registerTestingBuiltins`, but the exact set undocumented beyond examples).
- WASM playground path: **Inferred** (code exists, not exercised here).
