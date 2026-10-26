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

func TestArray_PetakanAndMap(t *testing.T) {
	code := `
	buat awal = [1, 2, 3]
	buat dikali2 = awal.petakan(x => x * 2)
	buat dengan_indeks = awal.petakan((x, i) => x + i)
	buat map_english = awal.map(x => x * 10)
	buat arr_kosong = [].petakan(x => x)
	[awal, dikali2, dengan_indeks, map_english, arr_kosong]
	`
	res := runCelangTestCode(t, code)
	arr, ok := res.(*Array)
	if !ok {
		t.Fatalf("expected Array, got %s", res.Inspect())
	}
	// Verify immutability of original array
	if arr.Elements[0].Inspect() != "[1, 2, 3]" {
		t.Errorf("expected original array to be [1, 2, 3], got %s", arr.Elements[0].Inspect())
	}
	if arr.Elements[1].Inspect() != "[2, 4, 6]" {
		t.Errorf("expected [2, 4, 6], got %s", arr.Elements[1].Inspect())
	}
	if arr.Elements[2].Inspect() != "[1, 3, 5]" {
		t.Errorf("expected [1, 3, 5], got %s", arr.Elements[2].Inspect())
	}
	if arr.Elements[3].Inspect() != "[10, 20, 30]" {
		t.Errorf("expected [10, 20, 30], got %s", arr.Elements[3].Inspect())
	}
	if arr.Elements[4].Inspect() != "[]" {
		t.Errorf("expected empty array, got %s", arr.Elements[4].Inspect())
	}
}

func TestArray_SaringAndFilter(t *testing.T) {
	code := `
	buat awal = [1, 2, 3, 4, 5, 6]
	buat genap = awal.saring(x => x % 2 == 0)
	buat indeks_ganjil = awal.saring((x, i) => i % 2 != 0)
	buat filter_english = awal.filter(x => x > 4)
	buat tak_cocok = awal.saring(x => x > 100)
	[awal, genap, indeks_ganjil, filter_english, tak_cocok]
	`
	res := runCelangTestCode(t, code)
	arr, ok := res.(*Array)
	if !ok {
		t.Fatalf("expected Array, got %s", res.Inspect())
	}
	// Immutability
	if arr.Elements[0].Inspect() != "[1, 2, 3, 4, 5, 6]" {
		t.Errorf("expected original array untouched, got %s", arr.Elements[0].Inspect())
	}
	if arr.Elements[1].Inspect() != "[2, 4, 6]" {
		t.Errorf("expected [2, 4, 6], got %s", arr.Elements[1].Inspect())
	}
	if arr.Elements[2].Inspect() != "[2, 4, 6]" {
		t.Errorf("expected elements at odd indices [2, 4, 6], got %s", arr.Elements[2].Inspect())
	}
	if arr.Elements[3].Inspect() != "[5, 6]" {
		t.Errorf("expected [5, 6], got %s", arr.Elements[3].Inspect())
	}
	if arr.Elements[4].Inspect() != "[]" {
		t.Errorf("expected [], got %s", arr.Elements[4].Inspect())
	}
}

func TestArray_LipatAndReduce(t *testing.T) {
	code := `
	buat angka = [1, 2, 3, 4]
	buat total_dengan_awal = angka.lipat((acc, x) => acc + x, 10)
	buat total_tanpa_awal = angka.lipat((acc, x) => acc + x)
	buat total_reduksi = angka.reduksi((acc, x) => acc + x, 5)
	buat reduce_english = angka.reduce((acc, x) => acc * x, 1)
	buat dengan_indeks = angka.lipat((acc, x, i) => acc + (x * i), 0) // 0 + 1*0 + 2*1 + 3*2 + 4*3 = 0+0+2+6+12 = 20
	buat kosong_dengan_awal = [].lipat((acc, x) => acc + x, 99)
	[total_dengan_awal, total_tanpa_awal, total_reduksi, reduce_english, dengan_indeks, kosong_dengan_awal]
	`
	res := runCelangTestCode(t, code)
	arr, ok := res.(*Array)
	if !ok {
		t.Fatalf("expected Array, got %s", res.Inspect())
	}
	if arr.Elements[0].Inspect() != "20" {
		t.Errorf("expected 20, got %s", arr.Elements[0].Inspect())
	}
	if arr.Elements[1].Inspect() != "10" {
		t.Errorf("expected 10, got %s", arr.Elements[1].Inspect())
	}
	if arr.Elements[2].Inspect() != "15" {
		t.Errorf("expected 15 from reduksi, got %s", arr.Elements[2].Inspect())
	}
	if arr.Elements[3].Inspect() != "24" {
		t.Errorf("expected 24, got %s", arr.Elements[3].Inspect())
	}
	if arr.Elements[4].Inspect() != "20" {
		t.Errorf("expected 20, got %s", arr.Elements[4].Inspect())
	}
	if arr.Elements[5].Inspect() != "99" {
		t.Errorf("expected 99, got %s", arr.Elements[5].Inspect())
	}

	// Error on empty array without initial accumulator
	codeErr := `[].lipat((acc, x) => acc + x)`
	resErr := runCelangTestCode(t, codeErr)
	if !strings.Contains(resErr.Inspect(), "array kosong tanpa nilai awal") {
		t.Errorf("expected error on empty array without initial value, got %s", resErr.Inspect())
	}
}

func TestArray_TemukanAndFind(t *testing.T) {
	code := `
	buat angka = [10, 20, 30, 40]
	buat hasil1 = angka.temukan(x => x > 25)
	buat hasil2 = angka.find(x => x == 20)
	buat hasil3 = angka.temukan((x, i) => i == 0)
	buat takada1 = angka.temukan(x => x > 100)
	buat cadangan = (angka.temukan(x => x > 100)) ?? 999
	[hasil1, hasil2, hasil3, takada1, cadangan]
	`
	res := runCelangTestCode(t, code)
	arr, ok := res.(*Array)
	if !ok {
		t.Fatalf("expected Array, got %s", res.Inspect())
	}
	if arr.Elements[0].Inspect() != "30" {
		t.Errorf("expected 30, got %s", arr.Elements[0].Inspect())
	}
	if arr.Elements[1].Inspect() != "20" {
		t.Errorf("expected 20, got %s", arr.Elements[1].Inspect())
	}
	if arr.Elements[2].Inspect() != "10" {
		t.Errorf("expected 10, got %s", arr.Elements[2].Inspect())
	}
	if arr.Elements[3] != NULL {
		t.Errorf("expected NULL, got %s", arr.Elements[3].Inspect())
	}
	if arr.Elements[4].Inspect() != "999" {
		t.Errorf("expected 999 from fallback, got %s", arr.Elements[4].Inspect())
	}
}

func TestArray_AdaYangAndSemuanya(t *testing.T) {
	code := `
	buat angka = [2, 4, 6, 8]
	buat ada_ganjil = angka.ada_yang(x => x % 2 != 0)
	buat ada_besar = angka.some(x => x > 5)
	buat semua_genap = angka.semuanya(x => x % 2 == 0)
	buat semua_besar = angka.every(x => x > 5)
	buat kosong_some = [].ada_yang(x => true)
	buat kosong_every = [].semuanya(x => false)
	[ada_ganjil, ada_besar, semua_genap, semua_besar, kosong_some, kosong_every]
	`
	res := runCelangTestCode(t, code)
	arr, ok := res.(*Array)
	if !ok {
		t.Fatalf("expected Array, got %s", res.Inspect())
	}
	if arr.Elements[0] != FALSE {
		t.Errorf("expected false, got %s", arr.Elements[0].Inspect())
	}
	if arr.Elements[1] != TRUE {
		t.Errorf("expected true, got %s", arr.Elements[1].Inspect())
	}
	if arr.Elements[2] != TRUE {
		t.Errorf("expected true, got %s", arr.Elements[2].Inspect())
	}
	if arr.Elements[3] != FALSE {
		t.Errorf("expected false, got %s", arr.Elements[3].Inspect())
	}
	if arr.Elements[4] != FALSE {
		t.Errorf("expected false on empty some, got %s", arr.Elements[4].Inspect())
	}
	if arr.Elements[5] != TRUE {
		t.Errorf("expected true on empty every, got %s", arr.Elements[5].Inspect())
	}
}

func TestArray_ChainingMethods(t *testing.T) {
	code := `
	buat hasil = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
		.saring(x => x % 2 == 0)
		.petakan(x => x * 10)
		.lipat((acc, x) => acc + x, 0)
	hasil
	`
	res := runCelangTestCode(t, code)
	// (2 + 4 + 6 + 8 + 10) * 10 = 30 * 10 = 300
	if res.Inspect() != "300" {
		t.Fatalf("expected 300 from method chaining, got %s", res.Inspect())
	}
}

func TestBita_DecodingBase64AndHex(t *testing.T) {
	code := `
	buat b1 = bita.dari_base64("SGFsbyBDZWxhbmch")
	buat b2 = bytes.from_base64("SGFsbyBDZWxhbmch")
	buat b3 = bita.dari_heks("48656c6c6f")
	buat b4 = bytes.from_hex("0x48656c6c6f")
	buat b5 = bita("SGFsbyBDZWxhbmch", "base64")
	buat b6 = bita("48656c6c6f", "heks")
	buat b7 = bytes("48656c6c6f", "hex")
	[b1.ke_teks(), b2.to_string(), b3.ke_teks(), b4.to_string(), b5.ke_teks(), b6.ke_teks(), b7.to_string()]
	`
	res := runCelangTestCode(t, code)
	arr, ok := res.(*Array)
	if !ok {
		t.Fatalf("expected Array, got %s", res.Inspect())
	}
	if arr.Elements[0].Inspect() != "Halo Celang!" {
		t.Errorf("b1 failed: got %s", arr.Elements[0].Inspect())
	}
	if arr.Elements[1].Inspect() != "Halo Celang!" {
		t.Errorf("b2 failed: got %s", arr.Elements[1].Inspect())
	}
	if arr.Elements[2].Inspect() != "Hello" {
		t.Errorf("b3 failed: got %s", arr.Elements[2].Inspect())
	}
	if arr.Elements[3].Inspect() != "Hello" {
		t.Errorf("b4 failed: got %s", arr.Elements[3].Inspect())
	}
	if arr.Elements[4].Inspect() != "Halo Celang!" {
		t.Errorf("b5 failed: got %s", arr.Elements[4].Inspect())
	}
	if arr.Elements[5].Inspect() != "Hello" {
		t.Errorf("b6 failed: got %s", arr.Elements[5].Inspect())
	}
	if arr.Elements[6].Inspect() != "Hello" {
		t.Errorf("b7 failed: got %s", arr.Elements[6].Inspect())
	}

	// Invalid decoding tests
	codeErrBase64 := `bita.dari_base64("not!valid!base64!!!")`
	resErr1 := runCelangTestCode(t, codeErrBase64)
	if !strings.Contains(resErr1.Inspect(), "Gagal mendekode data Base64") {
		t.Errorf("expected Base64 decode error, got %s", resErr1.Inspect())
	}

	codeErrHex := `bita.dari_heks("ZZZZ")`
	resErr2 := runCelangTestCode(t, codeErrHex)
	if !strings.Contains(resErr2.Inspect(), "Gagal mendekode teks heksadesimal") {
		t.Errorf("expected Hex decode error, got %s", resErr2.Inspect())
	}
}

func TestUniversalMode_FaseA(t *testing.T) {
	codeEnglish := `
	"use universal"
	let arr = [1, 2, 3, 4]
	let doubled = arr.map(x => x * 2)
	let filtered = doubled.filter(x => x > 4)
	let sum = filtered.reduce((acc, x) => acc + x, 0)
	let found = filtered.find(x => x == 8)
	let hasEight = filtered.some(x => x == 8)
	let allBig = filtered.every(x => x >= 6)
	let b64 = bytes.from_base64("SGVsbG8=")
	let hexB = bytes.from_hex("48656c6c6f")
	[sum, found, hasEight, allBig, b64.to_string(), hexB.to_string()]
	`
	res := runCelangTestCode(t, codeEnglish)
	arr, ok := res.(*Array)
	if !ok {
		t.Fatalf("expected Array in universal mode, got %s", res.Inspect())
	}
	if arr.Elements[0].Inspect() != "14" { // 6 + 8 = 14
		t.Errorf("sum expected 14, got %s", arr.Elements[0].Inspect())
	}
	if arr.Elements[1].Inspect() != "8" {
		t.Errorf("found expected 8, got %s", arr.Elements[1].Inspect())
	}
	if arr.Elements[2] != TRUE {
		t.Errorf("hasEight expected true, got %s", arr.Elements[2].Inspect())
	}
	if arr.Elements[3] != TRUE {
		t.Errorf("allBig expected true, got %s", arr.Elements[3].Inspect())
	}
	if arr.Elements[4].Inspect() != "Hello" {
		t.Errorf("b64 expected Hello, got %s", arr.Elements[4].Inspect())
	}
	if arr.Elements[5].Inspect() != "Hello" {
		t.Errorf("hexB expected Hello, got %s", arr.Elements[5].Inspect())
	}

	// Verify Indonesian methods are prohibited in Universal Mode
	prohibitions := []struct {
		code string
		err  string
	}{
		{"\"use universal\"\n[1, 2].petakan(x => x)", "Indonesian method 'petakan' is not allowed in Universal Mode"},
		{"\"use universal\"\n[1, 2].saring(x => x > 0)", "Indonesian method 'saring' is not allowed in Universal Mode"},
		{"\"use universal\"\n[1, 2].lipat((a, x) => a + x, 0)", "Indonesian method 'lipat' is not allowed in Universal Mode"},
		{"\"use universal\"\n[1, 2].temukan(x => x == 1)", "Indonesian method 'temukan' is not allowed in Universal Mode"},
		{"\"use universal\"\n[1, 2].ada_yang(x => x == 1)", "Indonesian method 'ada_yang' is not allowed in Universal Mode"},
		{"\"use universal\"\n[1, 2].semuanya(x => x == 1)", "Indonesian method 'semuanya' is not allowed in Universal Mode"},
		{"\"use universal\"\nbita(\"Halo\")", "Indonesian type 'bita' is not allowed in Universal Mode"},
		{"\"use universal\"\nbytes.dari_base64(\"AA==\")", "Indonesian method 'dari_base64' is not allowed in Universal Mode"},
		{"\"use universal\"\nbytes.dari_heks(\"AA\")", "Indonesian method 'dari_heks' is not allowed in Universal Mode"},
	}

	for _, p := range prohibitions {
		errRes := runCelangTestCode(t, p.code)
		if !strings.Contains(errRes.Inspect(), p.err) {
			t.Errorf("code %s: expected error containing '%s', got '%s'", p.code, p.err, errRes.Inspect())
		}
	}
}

