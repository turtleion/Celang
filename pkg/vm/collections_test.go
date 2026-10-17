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
