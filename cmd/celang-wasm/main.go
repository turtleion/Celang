//go:build js && wasm

package main

import (
	"bytes"
	"fmt"
	"syscall/js"
	"time"

	"celang/pkg/compiler"
	"celang/pkg/evaluator"
	"celang/pkg/parser"
	"celang/pkg/vm"
)

func runCelangWrapper(this js.Value, args []js.Value) any {
	if len(args) < 1 {
		return map[string]any{
			"error":      "Kode sumber Celang tidak boleh kosong",
			"output":     "",
			"result":     "",
			"durationMs": 0,
		}
	}

	sourceCode := args[0].String()
	mode := "vm"
	if len(args) > 1 && args[1].Type() == js.TypeString {
		mode = args[1].String()
	}

	startTime := time.Now()

	lexer := parser.NewLexer(sourceCode)
	p := parser.NewParser(lexer)
	prog, err := p.ParseProgram()
	if err != nil {
		duration := time.Since(startTime).Milliseconds()
		return map[string]any{
			"error":      fmt.Sprintf("Kesalahan Parsing: %v", err),
			"output":     "",
			"result":     "",
			"durationMs": duration,
		}
	}

	var outBuf bytes.Buffer
	env := evaluator.NewEnvironment()
	evaluator.RegisterBuiltinsWithIO(env, nil, &outBuf)

	var finalResult evaluator.Object

	if mode == "ast" {
		res := evaluator.Eval(prog, env)
		if res != nil && res.Type() == evaluator.ERROR_OBJ {
			duration := time.Since(startTime).Milliseconds()
			return map[string]any{
				"error":      res.Inspect(),
				"output":     outBuf.String(),
				"result":     "",
				"durationMs": duration,
			}
		}
		finalResult = res
	} else {
		// Default: Bytecode VM
		comp := compiler.NewCompiler()
		chunk, err := comp.Compile(prog)
		if err != nil {
			duration := time.Since(startTime).Milliseconds()
			return map[string]any{
				"error":      fmt.Sprintf("Kesalahan Kompilasi VM: %v", err),
				"output":     outBuf.String(),
				"result":     "",
				"durationMs": duration,
			}
		}

		vmInst := vm.NewWithEnv(chunk, env)
		res, err := vmInst.Run()
		if err != nil {
			duration := time.Since(startTime).Milliseconds()
			return map[string]any{
				"error":      fmt.Sprintf("Kesalahan Runtime VM: %v", err),
				"output":     outBuf.String(),
				"result":     "",
				"durationMs": duration,
			}
		}
		finalResult = res
	}

	duration := time.Since(startTime).Milliseconds()
	resStr := ""
	if finalResult != nil {
		resStr = finalResult.Inspect()
	}

	return map[string]any{
		"error":      "",
		"output":     outBuf.String(),
		"result":     resStr,
		"durationMs": duration,
	}
}

func main() {
	// Expose runCelang global function to WebAssembly JS environment
	js.Global().Set("runCelang", js.FuncOf(runCelangWrapper))
	js.Global().Set("celangWasmReady", true)

	// Keep WebAssembly running indefinitely
	select {}
}
