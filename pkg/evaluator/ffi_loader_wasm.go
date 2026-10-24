//go:build js && wasm

package evaluator

import (
	"fmt"
)

func ffiOpen(path string) (uintptr, error) {
	return 0, fmt.Errorf("FFI tidak didukung di lingkungan WebAssembly/browser")
}

func ffiSymbol(handle uintptr, name string) (uintptr, error) {
	return 0, fmt.Errorf("FFI tidak didukung di lingkungan WebAssembly/browser")
}

func ffiClose(handle uintptr) error {
	return nil
}
