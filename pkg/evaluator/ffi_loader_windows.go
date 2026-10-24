//go:build windows && !wasm && !js

package evaluator

import (
	"syscall"
)

func ffiOpen(path string) (uintptr, error) {
	h, err := syscall.LoadLibrary(path)
	return uintptr(h), err
}

func ffiSymbol(handle uintptr, name string) (uintptr, error) {
	proc, err := syscall.GetProcAddress(syscall.Handle(handle), name)
	return proc, err
}

func ffiClose(handle uintptr) error {
	return syscall.FreeLibrary(syscall.Handle(handle))
}
