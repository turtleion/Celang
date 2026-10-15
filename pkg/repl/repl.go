package repl

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"celang/pkg/bytecode"
	"celang/pkg/compiler"
	"celang/pkg/evaluator"
	"celang/pkg/parser"
	"celang/pkg/syntax"
	"celang/pkg/uast"
	"celang/pkg/vm"
)

const (
	PromptPrimary   = "celang> "
	PromptSecondary = "......> "
	Version         = "1.0.0"
)

// Options configures the REPL session.
type Options struct {
	In       io.Reader
	Out      io.Writer
	Style    string
	Env      *evaluator.Environment
	NoBanner bool
}

// Start launches the interactive REPL loop.
func Start(opts Options) error {
	in := opts.In
	if in == nil {
		in = os.Stdin
	}
	out := opts.Out
	if out == nil {
		out = os.Stdout
	}

	env := opts.Env
	if env == nil {
		env = evaluator.NewEnvironment()
		evaluator.RegisterBuiltinsWithIO(env, in, out)
	}

	styleName := opts.Style
	if styleName == "" {
		styleName = "indonesian"
	}
	activeProfile := syntax.GetProfile(styleName)

	if !opts.NoBanner {
		printBanner(out, activeProfile.DisplayName)
	}

	scanner := bufio.NewScanner(in)
	var buffer strings.Builder
	braceDepth := 0

	for {
		if braceDepth == 0 {
			fmt.Fprint(out, PromptPrimary)
		} else {
			fmt.Fprint(out, PromptSecondary)
		}

		if !scanner.Scan() {
			// EOF reached (Ctrl+D)
			fmt.Fprintln(out, "\nSampai jumpa lagi! 👋")
			break
		}

		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		// Handle REPL meta-commands on empty buffer
		if buffer.Len() == 0 && strings.HasPrefix(trimmed, ":") {
			shouldExit := handleMetaCommand(trimmed, env, &activeProfile, out)
			if shouldExit {
				break
			}
			continue
		}

		// Skip comments and empty lines if not in multiline mode
		if buffer.Len() == 0 && (trimmed == "" || strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "#")) {
			continue
		}

		if buffer.Len() > 0 {
			buffer.WriteString("\n")
		}
		buffer.WriteString(line)

		braceDepth += countBraceDelta(line)
		if braceDepth < 0 {
			braceDepth = 0
		}

		// Wait for more lines if braces are still open
		if braceDepth > 0 {
			continue
		}

		// Process complete input buffer
		fullInput := buffer.String()
		buffer.Reset()

		evaluateInput(fullInput, env, activeProfile, out)
	}

	return scanner.Err()
}

func printBanner(out io.Writer, styleName string) {
	fmt.Fprintln(out, "🌟 Celang Interactive REPL v"+Version+" (Bytecode VM Engine)")
	fmt.Fprintf(out, "Gaya aktif: %s | Ketik ':bantuan' untuk petunjuk, ':keluar' untuk selesai.\n\n", styleName)
}

func countBraceDelta(line string) int {
	delta := 0
	inQuote := byte(0)
	for i := 0; i < len(line); i++ {
		ch := line[i]
		if inQuote != 0 {
			if ch == '\\' {
				i++ // skip escape char
			} else if ch == inQuote {
				inQuote = 0
			}
			continue
		}
		if ch == '"' || ch == '\'' || ch == '`' {
			inQuote = ch
			continue
		}
		// Ignore comments
		if ch == '/' && i+1 < len(line) && line[i+1] == '/' {
			break
		}
		if ch == '#' {
			break
		}
		switch ch {
		case '{', '(', '[':
			delta++
		case '}', ')', ']':
			delta--
		}
	}
	return delta
}

func handleMetaCommand(cmd string, env *evaluator.Environment, profile **syntax.StyleProfile, out io.Writer) bool {
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return false
	}

	switch parts[0] {
	case ":keluar", ":exit", ":quit", ":q":
		fmt.Fprintln(out, "Sampai jumpa lagi! 👋")
		return true

	case ":bersih", ":clear", ":cls":
		fmt.Fprint(out, "\033[H\033[2J")

	case ":lingkungan", ":env", ":vars":
		entries := env.AllEntries()
		if len(entries) == 0 {
			fmt.Fprintln(out, "(Lingkungan kosong - belum ada variabel yang dideklarasikan)")
		} else {
			fmt.Fprintln(out, "📦 Variabel & Simbol di Lingkungan Sesi:")
			for name, entry := range entries {
				val := entry.Current()
				valType := "Sembarang"
				if val != nil {
					valType = string(val.Type())
				}
				valStr := "null"
				if val != nil {
					valStr = val.Inspect()
				}
				fmt.Fprintf(out, "  • %s (%s): %s\n", name, valType, valStr)
			}
		}

	case ":gaya", ":style":
		if len(parts) < 2 {
			fmt.Fprintf(out, "Gaya aktif saat ini: %s\n", (*profile).DisplayName)
			fmt.Fprintln(out, "Gunakan: :gaya <indonesian | javascript | pythonic | sunda | path/ke/file.style>")
		} else {
			target := parts[1]
			newProf := syntax.GetProfile(target)
			*profile = newProf
			fmt.Fprintf(out, "=> Gaya berhasil diganti ke: %s\n", newProf.DisplayName)
		}

	case ":ast", ":uast":
		if len(parts) < 2 {
			fmt.Fprintln(out, "Penggunaan: :ast <kode_ekspresi>")
		} else {
			codeToParse := strings.TrimPrefix(cmd, parts[0]+" ")
			l := parser.NewLexer(codeToParse)
			p := parser.NewParser(l)
			prog, err := p.ParseProgram()
			if err != nil {
				fmt.Fprintf(out, "❌ Error parsing: %v\n", err)
			} else {
				jsonBytes, _ := json.MarshalIndent(prog, "", "  ")
				fmt.Fprintln(out, string(jsonBytes))
			}
		}

	case ":bc", ":disasm", ":bytecode":
		if len(parts) < 2 {
			fmt.Fprintln(out, "Penggunaan: :bc <kode_ekspresi>")
		} else {
			codeToParse := strings.TrimPrefix(cmd, parts[0]+" ")
			l := parser.NewLexer(codeToParse)
			p := parser.NewParser(l)
			prog, err := p.ParseProgram()
			if err != nil {
				fmt.Fprintf(out, "❌ Error parsing: %v\n", err)
			} else {
				comp := compiler.NewCompiler()
				chunk, err := comp.Compile(prog)
				if err != nil {
					fmt.Fprintf(out, "❌ Error kompilasi: %v\n", err)
				} else {
					fmt.Fprint(out, chunk.Disassemble("repl_snippet"))
				}
			}
		}

	case ":reset", ":ulang":
		*env = *evaluator.NewEnvironment()
		evaluator.RegisterBuiltins(env, out)
		fmt.Fprintln(out, "=> Lingkungan sesi berhasil direset kembali ke kondisi awal.")

	case ":bantuan", ":help", ":?":
		printHelp(out)

	default:
		fmt.Fprintf(out, "Perintah khusus tidak dikenal: '%s'. Ketik ':bantuan' untuk daftar perintah.\n", parts[0])
	}

	return false
}

func printHelp(out io.Writer) {
	fmt.Fprintln(out, "📖 Daftar Perintah Khusus Celang REPL:")
	fmt.Fprintln(out, "  :bantuan, :help, :?         Menampilkan panduan ini")
	fmt.Fprintln(out, "  :keluar, :exit, :quit, :q   Keluar dari sesi REPL")
	fmt.Fprintln(out, "  :bersih, :clear             Membersihkan layar terminal")
	fmt.Fprintln(out, "  :lingkungan, :env, :vars    Melihat variabel & fungsi yang aktif di memori")
	fmt.Fprintln(out, "  :gaya <nama|path>           Mengganti profil gaya sintaks aktif")
	fmt.Fprintln(out, "  :bc <ekspresi>              Melihat instruksi assembly Bytecode")
	fmt.Fprintln(out, "  :ast <ekspresi>             Melihat pohon U-AST JSON dari baris kode")
	fmt.Fprintln(out, "  :reset, :ulang              Mengosongkan memori variabel kembali ke awal")
}

func evaluateInput(input string, env *evaluator.Environment, profile *syntax.StyleProfile, out io.Writer) {
	l := parser.NewLexer(input)
	p := parser.NewParser(l)
	prog, err := p.ParseProgram()
	if err != nil {
		fmt.Fprintf(out, "❌ %v\n", err)
		return
	}

	for _, stmt := range prog.Body {
		comp := compiler.NewCompilerWithEnv(env)
		var chunk *bytecode.Chunk
		var compErr error
		isExprStmt := false

		if exprStmt, ok := stmt.(*uast.ExpressionStatement); ok {
			isExprStmt = true
			chunk, compErr = comp.CompileExpression(exprStmt.Expression)
		} else {
			progSingle := &uast.Program{Body: []uast.Statement{stmt}}
			chunk, compErr = comp.Compile(progSingle)
		}

		if compErr != nil {
			fmt.Fprintf(out, "❌ Kesalahan kompilasi: %v\n", compErr)
			return
		}

		virtualMachine := vm.NewWithEnv(chunk, env)
		result, err := virtualMachine.Run()
		if err != nil {
			fmt.Fprintf(out, "❌ %v\n", err)
			return
		}

		if isExprStmt && result != nil {
			fmt.Fprintf(out, "=> %s\n", result.Inspect())
		}
	}
}
