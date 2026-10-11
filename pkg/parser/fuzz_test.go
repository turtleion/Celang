package parser_test

import (
	"testing"

	"celang/pkg/compiler"
	"celang/pkg/parser"
	"celang/pkg/vm"
)

// FuzzParser ensures that random/malformed inputs never cause unhandled panics
// in the lexer, parser, or compiler pipelines.
func FuzzParser(f *testing.F) {
	// Seed corpus with valid and edge-case snippets
	seeds := []string{
		`buat x = 10`,
		`jika x > 5 { cetak("ok") } selain { cetak("no") }`,
		`fungsi kali(a, b) { kembalikan a * b }`,
		`bentuk Titik { x: Angka, y: Angka }`,
		`[1, 2, 3] |> panjang`,
		`sambil { waktu.tidur(10) }`,
		`bersamaan { buat a = 1; buat b = 2 }`,
		`""`,
		`{}`,
		`[]`,
		`// komentar`,
		`/* blok komentar */`,
		`123 + 456 / (789 - 10)`,
		`mundur ke nilai 5`,
		`"teks dengan interpolasi {10 + 20}"`,
		`[[[[[[[[[[[]]]]]]]]]]]`,
		`(((((((((())))))))))`,
		`{ "a": { "b": { "c": 123 } } }`,
	}

	for _, seed := range seeds {
		f.Add([]byte(seed))
	}

	f.Fuzz(func(t *testing.T, data []byte) {
		input := string(data)

		// 1. Lexer & Parser Fuzzing
		l := parser.NewLexer(input)
		p := parser.NewParser(l)
		prog, err := p.ParseProgram()
		if err != nil {
			// Parsing errors are expected on random bytes, should NOT panic
			return
		}

		if prog == nil {
			return
		}

		// 2. Compiler Fuzzing on successfully parsed AST
		comp := compiler.NewCompiler()
		chunk, err := comp.Compile(prog)
		if err != nil {
			// Compilation errors are fine, should NOT panic
			return
		}

		// 3. Optional VM dry run with safe instruction limit or instantiation
		if chunk != nil {
			vmInst := vm.New(chunk)
			// VM run may fail with runtime error, should not crash Go runtime
			_, _ = vmInst.Run()
		}
	})
}
