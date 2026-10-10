package evaluator

import (
	"fmt"
	"os"
	"path/filepath"
)

func initStdlibBerkas(env *Environment) {
	// baca / read
	bacaFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "berkas.baca mengharapkan 1 argumen: path berkas", "Contoh: berkas.baca(\"./data.txt\")")
			}
			pathStr, ok := args[0].(*String)
			if !ok {
				return newError(nil, fmt.Sprintf("argumen path berkas.baca harus teks, menerima %s", args[0].Type()), "")
			}

			data, err := os.ReadFile(pathStr.Value)
			if err != nil {
				return newError(nil, fmt.Sprintf("Gagal membaca berkas '%s': %v", pathStr.Value, err), "Pastikan path berkas sudah benar dan berkas tersebut ada di disk.")
			}
			return &String{Value: string(data)}
		},
	}
	definePublicFunction(env, "baca", bacaFn)
	definePublicFunction(env, "read", bacaFn)

	// tulis / write
	tulisFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "berkas.tulis mengharapkan 2 argumen: path berkas dan data teks", "Contoh: berkas.tulis(\"./out.txt\", \"isi teks\")")
			}
			pathStr, ok := args[0].(*String)
			if !ok {
				return newError(nil, fmt.Sprintf("argumen path berkas.tulis harus teks, menerima %s", args[0].Type()), "")
			}
			data := args[1].Inspect()
			if s, ok := args[1].(*String); ok {
				data = s.Value
			}

			dir := filepath.Dir(pathStr.Value)
			if dir != "" && dir != "." {
				_ = os.MkdirAll(dir, 0755)
			}

			if err := os.WriteFile(pathStr.Value, []byte(data), 0644); err != nil {
				return newError(nil, fmt.Sprintf("Gagal menulis berkas '%s': %v", pathStr.Value, err), "Pastikan Anda memiliki izin tulis di direktori target.")
			}
			return TRUE
		},
	}
	definePublicFunction(env, "tulis", tulisFn)
	definePublicFunction(env, "write", tulisFn)

	// tambah / append
	tambahFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "berkas.tambah mengharapkan 2 argumen: path berkas dan data teks", "Contoh: berkas.tambah(\"./log.txt\", \"baris baru\\n\")")
			}
			pathStr, ok := args[0].(*String)
			if !ok {
				return newError(nil, fmt.Sprintf("argumen path berkas.tambah harus teks, menerima %s", args[0].Type()), "")
			}
			data := args[1].Inspect()
			if s, ok := args[1].(*String); ok {
				data = s.Value
			}

			dir := filepath.Dir(pathStr.Value)
			if dir != "" && dir != "." {
				_ = os.MkdirAll(dir, 0755)
			}

			f, err := os.OpenFile(pathStr.Value, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return newError(nil, fmt.Sprintf("Gagal membuka berkas '%s' untuk penambahan: %v", pathStr.Value, err), "")
			}
			defer f.Close()

			if _, err := f.WriteString(data); err != nil {
				return newError(nil, fmt.Sprintf("Gagal menambahkan teks ke berkas '%s': %v", pathStr.Value, err), "")
			}
			return TRUE
		},
	}
	definePublicFunction(env, "tambah", tambahFn)
	definePublicFunction(env, "append", tambahFn)

	// ada / exists
	adaFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "berkas.ada mengharapkan 1 argumen: path berkas", "")
			}
			pathStr, ok := args[0].(*String)
			if !ok {
				return newError(nil, fmt.Sprintf("argumen path berkas.ada harus teks, menerima %s", args[0].Type()), "")
			}

			_, err := os.Stat(pathStr.Value)
			if err == nil {
				return TRUE
			}
			if os.IsNotExist(err) {
				return FALSE
			}
			return FALSE
		},
	}
	definePublicFunction(env, "ada", adaFn)
	definePublicFunction(env, "exists", adaFn)

	// hapus / remove
	hapusFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "berkas.hapus mengharapkan 1 argumen: path berkas", "")
			}
			pathStr, ok := args[0].(*String)
			if !ok {
				return newError(nil, fmt.Sprintf("argumen path berkas.hapus harus teks, menerima %s", args[0].Type()), "")
			}

			if err := os.Remove(pathStr.Value); err != nil {
				return newError(nil, fmt.Sprintf("Gagal menghapus berkas '%s': %v", pathStr.Value, err), "Pastikan berkas ada dan Anda memiliki izin untuk menghapusnya.")
			}
			return TRUE
		},
	}
	definePublicFunction(env, "hapus", hapusFn)
	definePublicFunction(env, "remove", hapusFn)
}
