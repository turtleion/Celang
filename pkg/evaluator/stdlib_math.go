package evaluator

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

func initStdlibMatematika(env *Environment) {
	// akar / sqrt
	akarFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "matematika.akar mengharapkan 1 argumen angka", "")
			}
			_, fVal, _, ok := CoerceToNumeric(args[0])
			if !ok {
				return newError(nil, fmt.Sprintf("argumen matematika.akar harus berupa angka, menerima %s", args[0].Type()), "")
			}
			if fVal < 0 {
				return newError(nil, "nilai untuk matematika.akar tidak boleh negatif", "Gunakan bilangan bernilai >= 0.")
			}
			res := math.Sqrt(fVal)
			if res == math.Trunc(res) {
				return &Integer{Value: int64(res)}
			}
			return &Float{Value: res}
		},
	}
	definePublicFunction(env, "akar", akarFn)
	definePublicFunction(env, "sqrt", akarFn)

	// pangkat / pow
	pangkatFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "matematika.pangkat mengharapkan 2 argumen: basis dan eksponen", "")
			}
			_, aVal, aIsInt, ok1 := CoerceToNumeric(args[0])
			_, bVal, bIsInt, ok2 := CoerceToNumeric(args[1])
			if !ok1 || !ok2 {
				return newError(nil, "argumen matematika.pangkat harus berupa angka", "")
			}
			res := math.Pow(aVal, bVal)
			if aIsInt && bIsInt && bVal >= 0 && res == math.Trunc(res) {
				return &Integer{Value: int64(res)}
			}
			return &Float{Value: res}
		},
	}
	definePublicFunction(env, "pangkat", pangkatFn)
	definePublicFunction(env, "pow", pangkatFn)

	// mutlak / abs
	mutlakFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "matematika.mutlak mengharapkan 1 argumen angka", "")
			}
			iVal, fVal, isInt, ok := CoerceToNumeric(args[0])
			if !ok {
				return newError(nil, fmt.Sprintf("argumen matematika.mutlak harus berupa angka, menerima %s", args[0].Type()), "")
			}
			if isInt {
				if iVal < 0 {
					return &Integer{Value: -iVal}
				}
				return &Integer{Value: iVal}
			}
			return &Float{Value: math.Abs(fVal)}
		},
	}
	definePublicFunction(env, "mutlak", mutlakFn)
	definePublicFunction(env, "abs", mutlakFn)

	// bulat / round
	bulatFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "matematika.bulat mengharapkan 1 argumen angka", "")
			}
			_, fVal, _, ok := CoerceToNumeric(args[0])
			if !ok {
				return newError(nil, "argumen matematika.bulat harus angka", "")
			}
			return &Integer{Value: int64(math.Round(fVal))}
		},
	}
	definePublicFunction(env, "bulat", bulatFn)
	definePublicFunction(env, "round", bulatFn)

	// lantai / floor
	lantaiFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "matematika.lantai mengharapkan 1 argumen angka", "")
			}
			_, fVal, _, ok := CoerceToNumeric(args[0])
			if !ok {
				return newError(nil, "argumen matematika.lantai harus angka", "")
			}
			return &Integer{Value: int64(math.Floor(fVal))}
		},
	}
	definePublicFunction(env, "lantai", lantaiFn)
	definePublicFunction(env, "floor", lantaiFn)

	// atap / ceil
	atapFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "matematika.atap mengharapkan 1 argumen angka", "")
			}
			_, fVal, _, ok := CoerceToNumeric(args[0])
			if !ok {
				return newError(nil, "argumen matematika.atap harus angka", "")
			}
			return &Integer{Value: int64(math.Ceil(fVal))}
		},
	}
	definePublicFunction(env, "atap", atapFn)
	definePublicFunction(env, "ceil", atapFn)

	// acak / random
	acakFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) == 0 {
				return &Float{Value: rnd.Float64()}
			}
			if len(args) == 2 {
				minVal, _, _, ok1 := CoerceToNumeric(args[0])
				maxVal, _, _, ok2 := CoerceToNumeric(args[1])
				if !ok1 || !ok2 {
					return newError(nil, "argumen acak harus berupa bilangan bulat", "")
				}
				if maxVal < minVal {
					minVal, maxVal = maxVal, minVal
				}
				diff := maxVal - minVal + 1
				res := minVal + rnd.Int63n(diff)
				return &Integer{Value: res}
			}
			return newError(nil, "matematika.acak mengharapkan 0 argumen atau 2 argumen: acak(min, max)", "")
		},
	}
	definePublicFunction(env, "acak", acakFn)
	definePublicFunction(env, "random", acakFn)

	// maks / max
	maksFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) == 0 {
				return newError(nil, "matematika.maks membutuhkan minimal 1 argumen", "")
			}
			maxInt, maxFloat, isAllInt, ok := CoerceToNumeric(args[0])
			if !ok {
				return newError(nil, "argumen matematika.maks harus angka", "")
			}
			for i := 1; i < len(args); i++ {
				iVal, fVal, isInt, ok := CoerceToNumeric(args[i])
				if !ok {
					return newError(nil, "argumen matematika.maks harus angka", "")
				}
				if !isInt {
					isAllInt = false
				}
				if fVal > maxFloat {
					maxFloat = fVal
					maxInt = iVal
				}
			}
			if isAllInt {
				return &Integer{Value: maxInt}
			}
			return &Float{Value: maxFloat}
		},
	}
	definePublicFunction(env, "maks", maksFn)
	definePublicFunction(env, "max", maksFn)

	// min / min
	minFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) == 0 {
				return newError(nil, "matematika.min membutuhkan minimal 1 argumen", "")
			}
			minInt, minFloat, isAllInt, ok := CoerceToNumeric(args[0])
			if !ok {
				return newError(nil, "argumen matematika.min harus angka", "")
			}
			for i := 1; i < len(args); i++ {
				iVal, fVal, isInt, ok := CoerceToNumeric(args[i])
				if !ok {
					return newError(nil, "argumen matematika.min harus angka", "")
				}
				if !isInt {
					isAllInt = false
				}
				if fVal < minFloat {
					minFloat = fVal
					minInt = iVal
				}
			}
			if isAllInt {
				return &Integer{Value: minInt}
			}
			return &Float{Value: minFloat}
		},
	}
	definePublicFunction(env, "min", minFn)

	// Constants
	definePublicConstant(env, "PI", &Float{Value: math.Pi})
	definePublicConstant(env, "E", &Float{Value: math.E})
}
