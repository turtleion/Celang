package evaluator

import (
	"fmt"
	"strings"
	"sync"

	"celang/pkg/uast"
)

var (
	stdlibMu      sync.RWMutex
	stdlibModules = map[string]*Environment{}
	stdlibMaps    = map[string]*Map{}
)

// StdlibNames maps all supported stdlib module names and aliases to canonical names.
var StdlibNames = map[string]string{
	"berkas":     "berkas",
	"fs":         "berkas",
	"matematika": "matematika",
	"math":       "matematika",
	"teks":       "teks",
	"strings":    "teks",
	"waktu":      "waktu",
	"time":       "waktu",
	"json":       "json",
	"sistem":     "sistem",
	"os":         "sistem",
	"http":       "http",
	"jaringan":   "http",
	"net":         "http",
	"ffi":         "ffi",
	"kripto":      "kripto",
	"crypto":      "kripto",
	"arsip":       "arsip",
	"archive":     "arsip",
	"kompresi":    "arsip",
	"compression": "arsip",
	"zip":         "arsip",
	"tar":         "arsip",
}

func isStdlibModule(name string) bool {
	clean := strings.ToLower(strings.TrimSpace(name))
	clean = strings.TrimPrefix(clean, "stdlib:")
	clean = strings.TrimPrefix(clean, "stdlib/")
	_, ok := StdlibNames[clean]
	return ok
}

func normalizeStdlibName(name string) string {
	clean := strings.ToLower(strings.TrimSpace(name))
	clean = strings.TrimPrefix(clean, "stdlib:")
	clean = strings.TrimPrefix(clean, "stdlib/")
	if canonical, ok := StdlibNames[clean]; ok {
		return canonical
	}
	return clean
}

func getStdlibModuleEnvironment(canonicalName string) (*Environment, error) {
	stdlibMu.RLock()
	env, ok := stdlibModules[canonicalName]
	stdlibMu.RUnlock()
	if ok {
		return env, nil
	}

	stdlibMu.Lock()
	defer stdlibMu.Unlock()

	// Double check after lock
	if env, ok := stdlibModules[canonicalName]; ok {
		return env, nil
	}

	newEnv := NewEnvironment()
	newEnv.FilePath = "stdlib:" + canonicalName

	switch canonicalName {
	case "berkas":
		initStdlibBerkas(newEnv)
	case "matematika":
		initStdlibMatematika(newEnv)
	case "teks":
		initStdlibTeks(newEnv)
	case "waktu":
		initStdlibWaktu(newEnv)
	case "json":
		initStdlibJSON(newEnv)
	case "sistem":
		initStdlibSistem(newEnv)
	case "http":
		initStdlibHTTP(newEnv)
	case "ffi":
		initStdlibFFI(newEnv)
	case "kripto":
		initStdlibKripto(newEnv)
	case "arsip":
		initStdlibArsip(newEnv)
	default:
		return nil, fmt.Errorf("modul standar '%s' belum diimplementasikan", canonicalName)
	}

	// Cache environment
	stdlibModules[canonicalName] = newEnv

	// Create and cache corresponding Map representation for global object access (e.g. berkas.baca)
	pairs := make(map[string]Object)
	order := []string{}
	for name, entry := range newEnv.AllEntries() {
		pairs[name] = entry.Current()
		order = append(order, name)
	}
	stdlibMaps[canonicalName] = &Map{Pairs: pairs, Order: order}

	return newEnv, nil
}

func getStdlibModuleMap(canonicalName string) *Map {
	if _, err := getStdlibModuleEnvironment(canonicalName); err != nil {
		return &Map{Pairs: make(map[string]Object)}
	}
	stdlibMu.RLock()
	defer stdlibMu.RUnlock()
	return stdlibMaps[canonicalName]
}

// RegisterStdlibGlobals injects standard library modules as global objects (e.g., berkas, fs, matematika, math, teks, strings, waktu, time, json, sistem, os, http, jaringan, net, arsip, archive).
func RegisterStdlibGlobals(env *Environment) {
	canonicalList := []string{"berkas", "matematika", "teks", "waktu", "json", "sistem", "http", "ffi", "kripto", "arsip"}
	for _, canon := range canonicalList {
		modMap := getStdlibModuleMap(canon)
		// Register canonical Indonesian name
		env.Define(canon, modMap, false, "Map")

		// Register English / JavaScript alias if distinct
		switch canon {
		case "berkas":
			env.Define("fs", modMap, false, "Map")
		case "matematika":
			env.Define("math", modMap, false, "Map")
		case "teks":
			env.Define("strings", modMap, false, "Map")
		case "waktu":
			env.Define("time", modMap, false, "Map")
		case "sistem":
			env.Define("os", modMap, false, "Map")
		case "http":
			env.Define("jaringan", modMap, false, "Map")
			env.Define("net", modMap, false, "Map")
		case "kripto":
			env.Define("crypto", modMap, false, "Map")
		case "arsip":
			env.Define("archive", modMap, false, "Map")
			env.Define("kompresi", modMap, false, "Map")
			env.Define("compression", modMap, false, "Map")
			env.Define("zip", modMap, false, "Map")
			env.Define("tar", modMap, false, "Map")
		}
	}
}

func definePublicFunction(env *Environment, name string, fn *Builtin) {
	env.DefineWithVisibility(name, fn, false, "function", &uast.Visibility{Kind: uast.VisibilityPublic})
}

func definePublicConstant(env *Environment, name string, val Object) {
	env.DefineWithVisibility(name, val, false, string(val.Type()), &uast.Visibility{Kind: uast.VisibilityPublic})
}
