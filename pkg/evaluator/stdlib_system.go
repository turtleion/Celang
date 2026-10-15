package evaluator

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"
)

var (
	osExitFunc        = os.Exit
	processStartTime  = time.Now()
	scriptArgsMu      sync.RWMutex
	currentScriptArgs []string
)

// SetScriptArgs stores the CLI arguments passed to Celang scripts.
func SetScriptArgs(args []string) {
	scriptArgsMu.Lock()
	defer scriptArgsMu.Unlock()
	currentScriptArgs = append([]string{}, args...)
}

// GetScriptArgs retrieves current CLI arguments.
func GetScriptArgs() []string {
	scriptArgsMu.RLock()
	defer scriptArgsMu.RUnlock()
	if currentScriptArgs != nil {
		return append([]string{}, currentScriptArgs...)
	}
	if len(os.Args) > 1 {
		return append([]string{}, os.Args[1:]...)
	}
	return []string{}
}

func initStdlibSistem(env *Environment) {
	// =============================================================
	// 1. Informasi Host & Platform
	// =============================================================

	// os / platform() -> string
	osFn := &Builtin{
		Fn: func(args ...Object) Object {
			return &String{Value: runtime.GOOS}
		},
	}
	definePublicFunction(env, "os", osFn)
	definePublicFunction(env, "platform", osFn)

	// arsitektur / arch() -> string
	arsitekturFn := &Builtin{
		Fn: func(args ...Object) Object {
			return &String{Value: runtime.GOARCH}
		},
	}
	definePublicFunction(env, "arsitektur", arsitekturFn)
	definePublicFunction(env, "arch", arsitekturFn)

	// jumlah_cpu / cpu_count() -> int
	jumlahCpuFn := &Builtin{
		Fn: func(args ...Object) Object {
			return &Integer{Value: int64(runtime.NumCPU())}
		},
	}
	definePublicFunction(env, "jumlah_cpu", jumlahCpuFn)
	definePublicFunction(env, "cpu_count", jumlahCpuFn)
	definePublicFunction(env, "cpus", jumlahCpuFn)

	// nama_host / hostname() -> string
	namaHostFn := &Builtin{
		Fn: func(args ...Object) Object {
			h, err := os.Hostname()
			if err != nil {
				return &String{Value: ""}
			}
			return &String{Value: h}
		},
	}
	definePublicFunction(env, "nama_host", namaHostFn)
	definePublicFunction(env, "hostname", namaHostFn)

	// folder_home / homedir() -> string
	folderHomeFn := &Builtin{
		Fn: func(args ...Object) Object {
			h, err := os.UserHomeDir()
			if err != nil {
				return &String{Value: ""}
			}
			return &String{Value: h}
		},
	}
	definePublicFunction(env, "folder_home", folderHomeFn)
	definePublicFunction(env, "homedir", folderHomeFn)

	// folder_sementara / tmpdir() -> string
	folderSementaraFn := &Builtin{
		Fn: func(args ...Object) Object {
			return &String{Value: os.TempDir()}
		},
	}
	definePublicFunction(env, "folder_sementara", folderSementaraFn)
	definePublicFunction(env, "tmpdir", folderSementaraFn)

	// pid() -> int
	pidFn := &Builtin{
		Fn: func(args ...Object) Object {
			return &Integer{Value: int64(os.Getpid())}
		},
	}
	definePublicFunction(env, "pid", pidFn)

	// waktu_aktif / uptime([satuan]) -> float (detik) / int (milidetik)
	waktuAktifFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) == 1 {
				unit, ok := args[0].(*String)
				if ok && (strings.ToLower(unit.Value) == "ms" || strings.ToLower(unit.Value) == "milidetik") {
					return &Integer{Value: time.Since(processStartTime).Milliseconds()}
				}
			}
			return &Float{Value: time.Since(processStartTime).Seconds()}
		},
	}
	definePublicFunction(env, "waktu_aktif", waktuAktifFn)
	definePublicFunction(env, "uptime", waktuAktifFn)

	// apakah_admin / is_admin() -> bool
	apakahAdminFn := &Builtin{
		Fn: func(args ...Object) Object {
			return nativeBool(isRootOrAdmin())
		},
	}
	definePublicFunction(env, "apakah_admin", apakahAdminFn)
	definePublicFunction(env, "is_admin", apakahAdminFn)
	definePublicFunction(env, "is_root", apakahAdminFn)

	// =============================================================
	// 2. Resource & Memori
	// =============================================================

	// memori / memory() -> Map
	memoriFn := &Builtin{
		Fn: func(args ...Object) Object {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)

			pairs := make(map[string]Object)
			order := []string{
				"alokasi", "alloc",
				"total_alokasi", "total_alloc",
				"sistem", "sys",
				"jumlah_gc", "num_gc",
			}
			allocObj := &Integer{Value: int64(m.Alloc)}
			totObj := &Integer{Value: int64(m.TotalAlloc)}
			sysObj := &Integer{Value: int64(m.Sys)}
			gcObj := &Integer{Value: int64(m.NumGC)}

			pairs["alokasi"] = allocObj
			pairs["alloc"] = allocObj
			pairs["total_alokasi"] = totObj
			pairs["total_alloc"] = totObj
			pairs["sistem"] = sysObj
			pairs["sys"] = sysObj
			pairs["jumlah_gc"] = gcObj
			pairs["num_gc"] = gcObj

			return &Map{Pairs: pairs, Order: order}
		},
	}
	definePublicFunction(env, "memori", memoriFn)
	definePublicFunction(env, "memory", memoriFn)

	// gc() -> bool
	gcFn := &Builtin{
		Fn: func(args ...Object) Object {
			runtime.GC()
			return TRUE
		},
	}
	definePublicFunction(env, "gc", gcFn)

	// =============================================================
	// 3. Direktori Kerja
	// =============================================================

	// direktori_kerja / cwd() -> string
	direktoriKerjaFn := &Builtin{
		Fn: func(args ...Object) Object {
			dir, err := os.Getwd()
			if err != nil {
				return newError(nil, fmt.Sprintf("Gagal mendapatkan direktori kerja: %v", err), "")
			}
			return &String{Value: dir}
		},
	}
	definePublicFunction(env, "direktori_kerja", direktoriKerjaFn)
	definePublicFunction(env, "cwd", direktoriKerjaFn)

	// ubah_direktori / chdir(path) -> bool
	ubahDirektoriFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "sistem.ubah_direktori mengharapkan 1 argumen jalur direktori", "Contoh: sistem.ubah_direktori(\"/tmp\")")
			}
			target, ok := args[0].(*String)
			if !ok {
				return newError(nil, fmt.Sprintf("argumen sistem.ubah_direktori harus berupa teks, menerima %s", args[0].Type()), "")
			}
			err := os.Chdir(target.Value)
			if err != nil {
				return newError(nil, fmt.Sprintf("Gagal mengubah direktori ke '%s': %v", target.Value, err), "Pastikan direktori tersebut ada dan memiliki hak akses baca/eksekusi.")
			}
			return TRUE
		},
	}
	definePublicFunction(env, "ubah_direktori", ubahDirektoriFn)
	definePublicFunction(env, "chdir", ubahDirektoriFn)

	// =============================================================
	// 4. Variabel Lingkungan (Environment Variables)
	// =============================================================

	// dapatkan_env / getenv(kunci, [bawaan]) -> string
	dapatkanEnvFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 1 || len(args) > 2 {
				return newError(nil, "sistem.dapatkan_env mengharapkan 1 atau 2 argumen: dapatkan_env(kunci, [bawaan])", "Contoh: sistem.dapatkan_env(\"PORT\", \"8080\")")
			}
			kunci, ok := args[0].(*String)
			if !ok {
				return newError(nil, fmt.Sprintf("kunci variabel lingkungan harus berupa teks, menerima %s", args[0].Type()), "")
			}
			val, found := os.LookupEnv(kunci.Value)
			if !found {
				if len(args) == 2 {
					return args[1]
				}
				return &String{Value: ""}
			}
			return &String{Value: val}
		},
	}
	definePublicFunction(env, "dapatkan_env", dapatkanEnvFn)
	definePublicFunction(env, "getenv", dapatkanEnvFn)
	definePublicFunction(env, "baca_env", dapatkanEnvFn)

	// setel_env / setenv(kunci, nilai) -> bool
	setelEnvFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "sistem.setel_env mengharapkan 2 argumen: setel_env(kunci, nilai)", "Contoh: sistem.setel_env(\"APP_MODE\", \"produksi\")")
			}
			kunci, ok := args[0].(*String)
			if !ok {
				return newError(nil, fmt.Sprintf("kunci variabel lingkungan harus berupa teks, menerima %s", args[0].Type()), "")
			}
			valStr := args[1].Inspect()
			if s, ok := args[1].(*String); ok {
				valStr = s.Value
			}
			err := os.Setenv(kunci.Value, valStr)
			if err != nil {
				return newError(nil, fmt.Sprintf("Gagal menyetel variabel lingkungan '%s': %v", kunci.Value, err), "")
			}
			return TRUE
		},
	}
	definePublicFunction(env, "setel_env", setelEnvFn)
	definePublicFunction(env, "setenv", setelEnvFn)
	definePublicFunction(env, "tulis_env", setelEnvFn)

	// hapus_env / unsetenv(kunci) -> bool
	hapusEnvFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "sistem.hapus_env mengharapkan 1 argumen kunci variabel", "Contoh: sistem.hapus_env(\"APP_MODE\")")
			}
			kunci, ok := args[0].(*String)
			if !ok {
				return newError(nil, fmt.Sprintf("kunci variabel lingkungan harus berupa teks, menerima %s", args[0].Type()), "")
			}
			err := os.Unsetenv(kunci.Value)
			if err != nil {
				return newError(nil, fmt.Sprintf("Gagal menghapus variabel lingkungan '%s': %v", kunci.Value, err), "")
			}
			return TRUE
		},
	}
	definePublicFunction(env, "hapus_env", hapusEnvFn)
	definePublicFunction(env, "unsetenv", hapusEnvFn)

	// ada_env / has_env(kunci) -> bool
	adaEnvFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "sistem.ada_env mengharapkan 1 argumen kunci variabel", "Contoh: sistem.ada_env(\"PATH\")")
			}
			kunci, ok := args[0].(*String)
			if !ok {
				return newError(nil, fmt.Sprintf("kunci variabel lingkungan harus berupa teks, menerima %s", args[0].Type()), "")
			}
			_, found := os.LookupEnv(kunci.Value)
			return nativeBool(found)
		},
	}
	definePublicFunction(env, "ada_env", adaEnvFn)
	definePublicFunction(env, "has_env", adaEnvFn)

	// semua_env / environ() -> Map
	semuaEnvFn := &Builtin{
		Fn: func(args ...Object) Object {
			pairs := make(map[string]Object)
			order := []string{}
			for _, entry := range os.Environ() {
				parts := strings.SplitN(entry, "=", 2)
				k := parts[0]
				v := ""
				if len(parts) > 1 {
					v = parts[1]
				}
				pairs[k] = &String{Value: v}
				order = append(order, k)
			}
			return &Map{Pairs: pairs, Order: order}
		},
	}
	definePublicFunction(env, "semua_env", semuaEnvFn)
	definePublicFunction(env, "environ", semuaEnvFn)
	definePublicFunction(env, "env", semuaEnvFn)

	// =============================================================
	// 5. Argumen CLI & Perjalanan Skrip
	// =============================================================

	// argumen / args([indeks, [bawaan]])
	argumenFn := &Builtin{
		Fn: func(args ...Object) Object {
			scriptArgs := GetScriptArgs()

			if len(args) == 0 {
				elements := make([]Object, len(scriptArgs))
				for i, arg := range scriptArgs {
					elements[i] = &String{Value: arg}
				}
				return &Array{Elements: elements}
			}

			idxObj, ok := args[0].(*Integer)
			if !ok {
				return newError(nil, "indeks sistem.argumen harus berupa angka bulat", "Contoh: sistem.argumen(0)")
			}
			idx := int(idxObj.Value)

			if idx >= 0 && idx < len(scriptArgs) {
				return &String{Value: scriptArgs[idx]}
			}

			if len(args) >= 2 {
				return args[1]
			}
			return NULL
		},
	}
	definePublicFunction(env, "argumen", argumenFn)
	definePublicFunction(env, "args", argumenFn)

	// argumen_mentah / raw_args() -> Array (termasuk biner celang)
	argumenMentahFn := &Builtin{
		Fn: func(args ...Object) Object {
			elements := make([]Object, len(os.Args))
			for i, arg := range os.Args {
				elements[i] = &String{Value: arg}
			}
			return &Array{Elements: elements}
		},
	}
	definePublicFunction(env, "argumen_mentah", argumenMentahFn)
	definePublicFunction(env, "raw_args", argumenMentahFn)

	// =============================================================
	// 6. Eksekusi Perintah & Proses Terminal (Killer Feature)
	// =============================================================

	// cari_program / which(nama) -> string | takada
	cariProgramFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "sistem.cari_program mengharapkan 1 argumen nama program", "Contoh: sistem.cari_program(\"git\")")
			}
			nameObj, ok := args[0].(*String)
			if !ok {
				return newError(nil, fmt.Sprintf("nama program harus berupa teks, menerima %s", args[0].Type()), "")
			}
			path, err := exec.LookPath(nameObj.Value)
			if err != nil {
				return NULL
			}
			return &String{Value: path}
		},
	}
	definePublicFunction(env, "cari_program", cariProgramFn)
	definePublicFunction(env, "which", cariProgramFn)
	definePublicFunction(env, "look_path", cariProgramFn)

	// jalankan / exec(perintah, [opsi]) -> Map
	jalankanFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 1 || len(args) > 2 {
				return newError(nil, "sistem.jalankan mengharapkan 1 atau 2 argumen: jalankan(perintah, [opsi])", "Contoh: sistem.jalankan(\"echo halo\") atau sistem.jalankan([\"ls\", \"-la\"])")
			}

			var opsiMap *Map
			if len(args) == 2 {
				m, ok := args[1].(*Map)
				if !ok {
					return newError(nil, fmt.Sprintf("opsi sistem.jalankan harus berupa Map, menerima %s", args[1].Type()), "")
				}
				opsiMap = m
			}

			ctx := context.Background()
			var cancel context.CancelFunc = func() {}
			if opsiMap != nil {
				timeoutMs := getSysIntOption(opsiMap, "timeout", "batas_waktu")
				if timeoutMs > 0 {
					ctx, cancel = context.WithTimeout(ctx, time.Duration(timeoutMs)*time.Millisecond)
				}
			}
			defer cancel()

			var cmd *exec.Cmd
			switch p := args[0].(type) {
			case *String:
				if runtime.GOOS == "windows" {
					cmd = exec.CommandContext(ctx, "cmd", "/C", p.Value)
				} else {
					cmd = exec.CommandContext(ctx, "sh", "-c", p.Value)
				}
			case *Array:
				if len(p.Elements) == 0 {
					return newError(nil, "daftar perintah tidak boleh kosong", "")
				}
				parts := make([]string, len(p.Elements))
				for i, el := range p.Elements {
					if s, ok := el.(*String); ok {
						parts[i] = s.Value
					} else {
						parts[i] = el.Inspect()
					}
				}
				cmd = exec.CommandContext(ctx, parts[0], parts[1:]...)
			default:
				return newError(nil, fmt.Sprintf("perintah harus berupa teks atau daftar teks, menerima %s", args[0].Type()), "")
			}

			if opsiMap != nil {
				if cwd, ok := getSysStringOption(opsiMap, "dir", "direktori", "cwd"); ok && cwd != "" {
					cmd.Dir = cwd
				}

				if envMap, ok := getSysMapOption(opsiMap, "env", "lingkungan"); ok {
					envList := os.Environ()
					for k, v := range envMap.Pairs {
						strVal := v.Inspect()
						if s, ok := v.(*String); ok {
							strVal = s.Value
						}
						envList = append(envList, fmt.Sprintf("%s=%s", k, strVal))
					}
					cmd.Env = envList
				}

				if stdinStr, ok := getSysStringOption(opsiMap, "stdin", "masukan"); ok {
					cmd.Stdin = strings.NewReader(stdinStr)
				}
			}

			var stdoutBuf, stderrBuf bytes.Buffer
			cmd.Stdout = &stdoutBuf

			gabungKeluar := false
			if opsiMap != nil {
				gabungKeluar = getSysBoolOption(opsiMap, "gabung_keluar", "gabung_output", "merge_output")
			}

			if gabungKeluar {
				cmd.Stderr = &stdoutBuf
			} else {
				cmd.Stderr = &stderrBuf
			}

			start := time.Now()
			runErr := cmd.Run()
			dur := time.Since(start).Milliseconds()

			exitCode := 0
			if runErr != nil {
				if exitErr, ok := runErr.(*exec.ExitError); ok {
					exitCode = exitErr.ExitCode()
				} else {
					exitCode = -1
				}
			}

			outStr := stdoutBuf.String()
			errStr := stderrBuf.String()

			if ctx.Err() == context.DeadlineExceeded {
				if errStr != "" {
					errStr += "\n"
				}
				errStr += "perintah dibatalkan karena batas waktu (timeout)"
				if exitCode == 0 {
					exitCode = -1
				}
			} else if runErr != nil && errStr == "" && exitCode == -1 {
				errStr = runErr.Error()
			}

			pairs := make(map[string]Object)
			order := []string{
				"keluaran", "stdout",
				"galat", "stderr",
				"kode_keluar", "exit_code",
				"sukses", "success",
				"durasi", "duration",
			}
			outObj := &String{Value: outStr}
			errObj := &String{Value: errStr}
			codeObj := &Integer{Value: int64(exitCode)}
			succObj := nativeBool(exitCode == 0)
			durObj := &Integer{Value: dur}

			pairs["keluaran"] = outObj
			pairs["stdout"] = outObj
			pairs["galat"] = errObj
			pairs["stderr"] = errObj
			pairs["kode_keluar"] = codeObj
			pairs["exit_code"] = codeObj
			pairs["sukses"] = succObj
			pairs["success"] = succObj
			pairs["durasi"] = durObj
			pairs["duration"] = durObj

			return &Map{Pairs: pairs, Order: order}
		},
	}
	definePublicFunction(env, "jalankan", jalankanFn)
	definePublicFunction(env, "eksekusi", jalankanFn)
	definePublicFunction(env, "exec", jalankanFn)
	definePublicFunction(env, "run", jalankanFn)

	// =============================================================
	// 7. Kontrol Siklus Hidup Proses
	// =============================================================

	// keluar / exit([kode=0])
	keluarFn := &Builtin{
		Fn: func(args ...Object) Object {
			code := 0
			if len(args) == 1 {
				if intObj, ok := args[0].(*Integer); ok {
					code = int(intObj.Value)
				}
			}
			osExitFunc(code)
			return NULL
		},
	}
	definePublicFunction(env, "keluar", keluarFn)
	definePublicFunction(env, "exit", keluarFn)

	// =============================================================
	// 8. Konstanta Sistem
	// =============================================================

	definePublicConstant(env, "PEMISAH_JALUR", &String{Value: string(os.PathListSeparator)})
	definePublicConstant(env, "PATH_SEPARATOR", &String{Value: string(os.PathListSeparator)})

	definePublicConstant(env, "PEMISAH_BERKAS", &String{Value: string(os.PathSeparator)})
	definePublicConstant(env, "DIR_SEPARATOR", &String{Value: string(os.PathSeparator)})

	definePublicConstant(env, "DEV_NULL", &String{Value: os.DevNull})

	eol := "\n"
	if runtime.GOOS == "windows" {
		eol = "\r\n"
	}
	definePublicConstant(env, "GARIS_BARU", &String{Value: eol})
	definePublicConstant(env, "EOL", &String{Value: eol})
}

// -----------------------------------------------------------------
// Internal Helpers for Option Maps
// -----------------------------------------------------------------

func getSysStringOption(m *Map, keys ...string) (string, bool) {
	if m == nil {
		return "", false
	}
	for _, k := range keys {
		if val, exists := m.Pairs[k]; exists {
			if s, ok := val.(*String); ok {
				return s.Value, true
			}
			return val.Inspect(), true
		}
	}
	return "", false
}

func getSysIntOption(m *Map, keys ...string) int64 {
	if m == nil {
		return 0
	}
	for _, k := range keys {
		if val, exists := m.Pairs[k]; exists {
			if num, ok := val.(*Integer); ok {
				return num.Value
			}
		}
	}
	return 0
}

func getSysBoolOption(m *Map, keys ...string) bool {
	if m == nil {
		return false
	}
	for _, k := range keys {
		if val, exists := m.Pairs[k]; exists {
			if b, ok := val.(*Boolean); ok {
				return b.Value
			}
		}
	}
	return false
}

func getSysMapOption(m *Map, keys ...string) (*Map, bool) {
	if m == nil {
		return nil, false
	}
	for _, k := range keys {
		if val, exists := m.Pairs[k]; exists {
			if childMap, ok := val.(*Map); ok {
				return childMap, true
			}
		}
	}
	return nil, false
}
