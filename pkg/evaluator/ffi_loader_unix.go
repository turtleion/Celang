//go:build !windows && !wasm && !js

package evaluator

import (
	"github.com/ebitengine/purego"
)

func ffiOpen(path string) (uintptr, error) {
	// Mode 0x002 is RTLD_NOW, 0x100 is RTLD_GLOBAL
	return purego.Dlopen(path, 0x002|0x100)
}

func ffiSymbol(handle uintptr, name string) (uintptr, error) {
	return purego.Dlsym(handle, name)
}

func ffiClose(handle uintptr) error {
	return purego.Dlclose(handle)
}
