package vm_test

import (
	"testing"
)

func TestVM_ArrayMethods(t *testing.T) {
	code := `
	buat arr = [1, 2, 3]
	arr.push(4, 5)
	buat last = arr.pop()
	buat first = arr.shift()
	arr.unshift(0)
	arr.reverse()
	arr.join("-")
	`
	res, _, _ := runVM(t, code)
	if res.Inspect() != "4-3-2-0" {
		t.Fatalf("expected 4-3-2-0, got %s", res.Inspect())
	}
}

func TestVM_MapMethods(t *testing.T) {
	code := `
	buat m = {"a": 1, "b": 2}
	m.set("c", 3)
	buat hasB = m.has("b")
	buat valB = m.get("b", 0)
	buat valD = m.get("d", 99)
	m.delete("a")
	buat k = m.keys().join(",")
	k + ":" + valB + ":" + valD
	`
	res, _, _ := runVM(t, code)
	if res.Inspect() != "b,c:2:99" {
		t.Fatalf("expected b,c:2:99, got %s", res.Inspect())
	}
}

func TestVM_BreakAndContinue(t *testing.T) {
	code := `
	"use universal"
	let count = 0
	let i = 0
	while true {
		i = i + 1
		if i == 2 {
			continue
		}
		if i > 5 {
			break
		}
		count = count + i
	}
	count
	`
	res, _, _ := runVM(t, code)
	if res.Inspect() != "13" {
		t.Fatalf("expected 13, got %s", res.Inspect())
	}
}

func TestVM_ForEach_Locals(t *testing.T) {
	code := `
	"use universal"
	let arr = [1, 2, 3]
	let res = []
	for x in arr {
		let y = x * 10
		res.push(y)
	}
	res.join(",")
	`
	res, _, _ := runVM(t, code)
	if res.Inspect() != "10,20,30" {
		t.Fatalf("expected 10,20,30, got %s", res.Inspect())
	}
}

func TestVM_For_Continue_Update(t *testing.T) {
	code := `
	"use universal"
	let sum = 0
	for i in 0..4 {
		if i == 2 {
			continue
		}
		sum = sum + i
	}
	sum
	`
	res, _, _ := runVM(t, code)
	if res.Inspect() != "8" {
		t.Fatalf("expected 8, got %s", res.Inspect())
	}
}

func TestVM_ForEach_ContinueAndBreak(t *testing.T) {
	code := `
	"use universal"
	let arr = [1, 2, 3, 4, 5]
	let res = []
	for x in arr {
		let temp = x * 2
		if x == 2 {
			continue
		}
		if x == 4 {
			break
		}
		res.push(temp)
	}
	res.join(",")
	`
	res, _, _ := runVM(t, code)
	if res.Inspect() != "2,6" {
		t.Fatalf("expected 2,6, got %s", res.Inspect())
	}
}

func TestVM_InnerFunctionAndDeferInLoop(t *testing.T) {
	// 1. Inner function inside loop capturing outer variable
	code := `
	fungsi luar() {
		buat x = 100
		untuk i dalam [1] {
			fungsi inner() {
				kembalikan x
			}
			kembalikan inner()
		}
	}
	luar()
	`
	res, _, _ := runVM(t, code)
	if res.Inspect() != "100" {
		t.Fatalf("expected 100, got: %s", res.Inspect())
	}

	// 2. Defer inside loop in function executed on function return
	codeDefer := `
	buat log = []
	fungsi uji_defer() {
		untuk i dalam [1] {
			saat_selesai log.tambah("defer_loop")
		}
		log.tambah("selesai_fungsi")
	}
	uji_defer()
	log.gabung(",")
	`
	resDefer, _, _ := runVM(t, codeDefer)
	if resDefer.Inspect() != "selesai_fungsi,defer_loop" {
		t.Fatalf("expected 'selesai_fungsi,defer_loop', got: %s", resDefer.Inspect())
	}
}

func TestVM_FaseA_FunctionalArrayAndBita(t *testing.T) {
	// 1. Array map, filter, reduce, find, some, every in VM
	codeArray := `
	buat awal = [1, 2, 3, 4, 5, 6]
	buat genap = awal.saring(x => x % 2 == 0)
	buat kali10 = genap.petakan((x, i) => x * 10 + i)
	buat total = kali10.lipat((acc, x) => acc + x, 0)
	buat temukan_ada = awal.temukan(x => x > 4)
	buat temukan_kosong = (awal.temukan(x => x > 100)) ?? 999
	buat ada_sepuluh = awal.ada_yang(x => x == 5)
	buat semua_positif = awal.semuanya(x => x > 0)
	[total, temukan_ada, temukan_kosong, ada_sepuluh, semua_positif].gabung(":")
	`
	res, _, _ := runVM(t, codeArray)
	expected := "123:5:999:true:true"
	if res.Inspect() != expected {
		t.Fatalf("expected '%s', got '%s'", expected, res.Inspect())
	}

	// 2. Chaining in VM
	codeChain := `
	[1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
		.filter(x => x % 2 == 0)
		.map(x => x * 10)
		.reduce((acc, x) => acc + x, 0)
	`
	resChain, _, _ := runVM(t, codeChain)
	if resChain.Inspect() != "300" {
		t.Fatalf("expected 300, got '%s'", resChain.Inspect())
	}

	// 3. Bita decoding in VM
	codeBita := `
	buat b1 = bita.dari_base64("SGFsbw==")
	buat b2 = bytes.from_base64("SGFsbw==")
	buat b3 = bita.dari_heks("48656c6c6f")
	buat b4 = bytes.from_hex("0x48656c6c6f")
	buat b5 = bita("SGFsbw==", "base64")
	buat b6 = bytes("48656c6c6f", "hex")
	[b1.ke_teks(), b2.to_string(), b3.ke_teks(), b4.to_string(), b5.ke_teks(), b6.to_string()].gabung(",")
	`
	resBita, _, _ := runVM(t, codeBita)
	if resBita.Inspect() != "Halo,Halo,Hello,Hello,Halo,Hello" {
		t.Fatalf("expected 'Halo,Halo,Hello,Hello,Halo,Hello', got '%s'", resBita.Inspect())
	}
}

func TestVM_FaseA_UniversalMode(t *testing.T) {
	code := `
	"use universal"
	let arr = [10, 20, 30]
	let mapped = arr.map(x => x / 10)
	let filtered = mapped.filter(x => x > 1)
	let sum = filtered.reduce((a, b) => a + b, 0)
	let b = bytes.from_base64("T0s=")
	sum + ":" + b.to_string()
	`
	res, _, _ := runVM(t, code)
	if res.Inspect() != "5:OK" {
		t.Fatalf("expected 5:OK, got: %s", res.Inspect())
	}
}

