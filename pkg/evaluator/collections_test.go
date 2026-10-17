package evaluator

import (
	"strings"
	"testing"
)

func TestArray_BasicMethods(t *testing.T) {
	code := `
	buat arr = [1, 2, 3]
	arr.tambah(4, 5)
	buat akhir = arr.hapus_terakhir()
	buat awal = arr.hapus_awal()
	arr.sisipkan_depan(0)
	[arr, akhir, awal]
	`
	res := runCelangTestCode(t, code)
	arr, ok := res.(*Array)
	if !ok {
		t.Fatalf("expected Array result, got %s", res.Inspect())
	}
	// arr.Elements[0] should be [0, 2, 3, 4]
	modArr := arr.Elements[0].(*Array)
	if len(modArr.Elements) != 4 {
		t.Fatalf("expected 4 elements, got %d", len(modArr.Elements))
	}
	if modArr.Inspect() != "[0, 2, 3, 4]" {
		t.Errorf("expected [0, 2, 3, 4], got %s", modArr.Inspect())
	}
	// akhir should be 5
	if arr.Elements[1].Inspect() != "5" {
		t.Errorf("expected 5, got %s", arr.Elements[1].Inspect())
	}
	// awal should be 1
	if arr.Elements[2].Inspect() != "1" {
		t.Errorf("expected 1, got %s", arr.Elements[2].Inspect())
	}
}

func TestArray_InsertAndRemove(t *testing.T) {
	code := `
	buat arr = ["a", "c"]
	arr.sisipkan(1, "b")
	buat removed = arr.hapus_indeks(0)
	[arr, removed]
	`
	res := runCelangTestCode(t, code)
	arr, ok := res.(*Array)
	if !ok {
		t.Fatalf("expected Array result, got %s", res.Inspect())
	}
	modArr := arr.Elements[0].(*Array)
	if modArr.Inspect() != "[b, c]" {
		t.Errorf("expected [b, c], got %s", modArr.Inspect())
	}
	if arr.Elements[1].Inspect() != "a" {
		t.Errorf("expected 'a', got %s", arr.Elements[1].Inspect())
	}
}

func TestArray_ReverseAndSort(t *testing.T) {
	code := `
	buat arr1 = [3, 1, 4, 1, 5, 9, 2]
	arr1.urut()
	buat arr2 = ["jeruk", "apel", "mangga"]
	arr2.urut()
	arr2.balik()
	[arr1, arr2]
	`
	res := runCelangTestCode(t, code)
	arr, ok := res.(*Array)
	if !ok {
		t.Fatalf("expected Array result, got %s", res.Inspect())
	}
	arr1 := arr.Elements[0].(*Array)
	if arr1.Inspect() != "[1, 1, 2, 3, 4, 5, 9]" {
		t.Errorf("expected sorted numbers, got %s", arr1.Inspect())
	}
	arr2 := arr.Elements[1].(*Array)
	if arr2.Inspect() != "[mangga, jeruk, apel]" {
		t.Errorf("expected sorted & reversed strings, got %s", arr2.Inspect())
	}
}

func TestArray_CustomSort(t *testing.T) {
	code := `
	buat arr = [10, 5, 30, 20]
	// descending sort
	arr.sort((a, b) => b - a)
	arr
	`
	res := runCelangTestCode(t, code)
	arr, ok := res.(*Array)
	if !ok || arr.Inspect() != "[30, 20, 10, 5]" {
		t.Fatalf("expected [30, 20, 10, 5], got %s", res.Inspect())
	}
}

func TestArray_SearchAndJoin(t *testing.T) {
	code := `
	buat arr = ["kucing", "anjing", "kelinci"]
	buat adaKucing = arr.berisi("kucing")
	buat adaSinga = arr.berisi("singa")
	buat idxAnjing = arr.indeks_dari("anjing")
	buat idxSinga = arr.indeks_dari("singa")
	buat gabungan = arr.gabung(" - ")
	[adaKucing, adaSinga, idxAnjing, idxSinga, gabungan]
	`
	res := runCelangTestCode(t, code)
	arr, ok := res.(*Array)
	if !ok {
		t.Fatalf("expected Array result, got %s", res.Inspect())
	}
	if arr.Elements[0].Inspect() != "true" || arr.Elements[1].Inspect() != "false" {
		t.Errorf("contains failed")
	}
	if arr.Elements[2].Inspect() != "1" || arr.Elements[3].Inspect() != "-1" {
		t.Errorf("index_of failed: got %s, %s", arr.Elements[2].Inspect(), arr.Elements[3].Inspect())
	}
	if arr.Elements[4].Inspect() != "kucing - anjing - kelinci" {
		t.Errorf("join failed: got %s", arr.Elements[4].Inspect())
	}
}

func TestArray_SliceConcatClear(t *testing.T) {
	code := `
	buat arr = [1, 2, 3, 4, 5]
	buat sl = arr.slice(1, 4)
	buat c = sl.concat([6, 7], 8)
	arr.clear()
	[sl, c, arr, arr.length, c.size()]
	`
	res := runCelangTestCode(t, code)
	arr, ok := res.(*Array)
	if !ok {
		t.Fatalf("expected Array result, got %s", res.Inspect())
	}
	if arr.Elements[0].Inspect() != "[2, 3, 4]" {
		t.Errorf("slice failed: %s", arr.Elements[0].Inspect())
	}
	if arr.Elements[1].Inspect() != "[2, 3, 4, 6, 7, 8]" {
		t.Errorf("concat failed: %s", arr.Elements[1].Inspect())
	}
	if arr.Elements[2].Inspect() != "[]" {
		t.Errorf("clear failed: %s", arr.Elements[2].Inspect())
	}
	if arr.Elements[3].Inspect() != "0" {
		t.Errorf("length failed: %s", arr.Elements[3].Inspect())
	}
	if arr.Elements[4].Inspect() != "6" {
		t.Errorf("size() failed: %s", arr.Elements[4].Inspect())
	}
}

func TestMap_KeysValuesEntries(t *testing.T) {
	code := `
	buat user = {"nama": "Budi", "umur": 25, "kota": "Jakarta"}
	buat k = user.kunci()
	buat v = user.nilai()
	buat e = user.entri()
	[k, v, e, user.panjang, user.ukuran()]
	`
	res := runCelangTestCode(t, code)
	arr, ok := res.(*Array)
	if !ok {
		t.Fatalf("expected Array result, got %s", res.Inspect())
	}
	if arr.Elements[0].Inspect() != "[nama, umur, kota]" {
		t.Errorf("keys failed: %s", arr.Elements[0].Inspect())
	}
	if arr.Elements[1].Inspect() != "[Budi, 25, Jakarta]" {
		t.Errorf("values failed: %s", arr.Elements[1].Inspect())
	}
	if arr.Elements[2].Inspect() != "[[nama, Budi], [umur, 25], [kota, Jakarta]]" {
		t.Errorf("entries failed: %s", arr.Elements[2].Inspect())
	}
	if arr.Elements[3].Inspect() != "3" || arr.Elements[4].Inspect() != "3" {
		t.Errorf("length/size failed: %s, %s", arr.Elements[3].Inspect(), arr.Elements[4].Inspect())
	}
}

func TestMap_GetSetDeleteHas(t *testing.T) {
	code := `
	buat data = {"a": 10}
	buat hasA = data.memiliki("a")
	buat hasB = data.has("b")
	buat valA = data.ambil("a")
	buat valB = data.get("b", 99)
	data.atur("c", 30)
	buat delA = data.hapus("a")
	[hasA, hasB, valA, valB, delA, data]
	`
	res := runCelangTestCode(t, code)
	arr, ok := res.(*Array)
	if !ok {
		t.Fatalf("expected Array result, got %s", res.Inspect())
	}
	if arr.Elements[0].Inspect() != "true" || arr.Elements[1].Inspect() != "false" {
		t.Errorf("has failed")
	}
	if arr.Elements[2].Inspect() != "10" || arr.Elements[3].Inspect() != "99" {
		t.Errorf("get failed: %s, %s", arr.Elements[2].Inspect(), arr.Elements[3].Inspect())
	}
	if arr.Elements[4].Inspect() != "true" {
		t.Errorf("delete failed")
	}
	modMap := arr.Elements[5].(*Map)
	if len(modMap.Pairs) != 1 || modMap.Pairs["c"].Inspect() != "30" {
		t.Errorf("unexpected map state: %s", modMap.Inspect())
	}
}

func TestMap_MergeAndClear(t *testing.T) {
	code := `
	buat m1 = {"x": 1, "y": 2}
	buat m2 = {"y": 20, "z": 30}
	m1.gabung(m2)
	buat mergedY = m1.y
	buat mergedZ = m1.z
	m1.bersihkan()
	[mergedY, mergedZ, m1.panjang]
	`
	res := runCelangTestCode(t, code)
	arr, ok := res.(*Array)
	if !ok {
		t.Fatalf("expected Array result, got %s", res.Inspect())
	}
	if arr.Elements[0].Inspect() != "20" || arr.Elements[1].Inspect() != "30" || arr.Elements[2].Inspect() != "0" {
		t.Errorf("merge/clear failed: %s", arr.Inspect())
	}
}

func TestCollections_UniversalMode(t *testing.T) {
	// English methods in Universal mode should succeed
	codeValid := `
	"use universal"
	let arr = [3, 1, 2]
	arr.push(4)
	arr.sort()
	let last = arr.pop()
	
	let m = {"a": 1}
	m.set("b", 2)
	let k = m.keys()
	let v = m.values()
	let hasB = m.has("b")
	let removed = m.delete("a")
	[arr, last, k, v, hasB, removed]
	`
	res := runCelangTestCode(t, codeValid)
	arr, ok := res.(*Array)
	if !ok {
		t.Fatalf("Universal Mode English methods should work, got %s", res.Inspect())
	}
	if arr.Elements[0].Inspect() != "[1, 2, 3]" || arr.Elements[1].Inspect() != "4" {
		t.Errorf("Array English methods failed: %s", arr.Inspect())
	}
	if arr.Elements[4].Inspect() != "true" || arr.Elements[5].Inspect() != "true" {
		t.Errorf("Map English methods failed: %s", arr.Inspect())
	}

	// Indonesian methods in Universal mode should fail with helpful hints
	codeInvalidArray := `
	"use universal"
	let arr = [1, 2]
	arr.tambah(3)
	`
	resErrArr := runCelangTestCode(t, codeInvalidArray)
	if !strings.Contains(resErrArr.Inspect(), "Indonesian method 'tambah' is not allowed in Universal Mode") {
		t.Errorf("expected Universal Mode error for 'tambah', got %s", resErrArr.Inspect())
	}

	codeInvalidMap := `
	"use universal"
	let m = {"a": 1}
	m.kunci()
	`
	resErrMap := runCelangTestCode(t, codeInvalidMap)
	if !strings.Contains(resErrMap.Inspect(), "Indonesian method 'kunci' is not allowed in Universal Mode") {
		t.Errorf("expected Universal Mode error for 'kunci', got %s", resErrMap.Inspect())
	}
}
