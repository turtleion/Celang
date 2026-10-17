package evaluator

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	rndMu sync.Mutex
	rnd   = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func initStdlibMatematika(env *Environment) {
	// -------------------------------------------------------------
	// 1. Operasi Aritmatika Dasar & Presisi
	// -------------------------------------------------------------

	// akar / sqrt(x)
	akarFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "matematika.akar mengharapkan 1 argumen angka", "Contoh: matematika.akar(16)")
			}
			_, fVal, _, ok := CoerceToNumeric(args[0])
			if !ok {
				return newError(nil, fmt.Sprintf("argumen matematika.akar harus berupa angka, menerima %s", args[0].Type()), "")
			}
			if fVal < 0 {
				return newError(nil, "nilai untuk matematika.akar tidak boleh negatif", "Gunakan bilangan bernilai >= 0.")
			}
			res := math.Sqrt(fVal)
			if res == math.Trunc(res) && res <= float64(math.MaxInt64) && res >= float64(math.MinInt64) {
				return &Integer{Value: int64(res)}
			}
			return &Float{Value: res}
		},
	}
	definePublicFunction(env, "akar", akarFn)
	definePublicFunction(env, "sqrt", akarFn)

	// pangkat / pow(basis, eksponen)
	pangkatFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "matematika.pangkat mengharapkan 2 argumen: basis dan eksponen", "Contoh: matematika.pangkat(2, 3)")
			}
			_, aVal, aIsInt, ok1 := CoerceToNumeric(args[0])
			_, bVal, bIsInt, ok2 := CoerceToNumeric(args[1])
			if !ok1 || !ok2 {
				return newError(nil, "argumen matematika.pangkat harus berupa angka", "")
			}
			res := math.Pow(aVal, bVal)
			if aIsInt && bIsInt && bVal >= 0 && res == math.Trunc(res) && res <= float64(math.MaxInt64) && res >= float64(math.MinInt64) {
				return &Integer{Value: int64(res)}
			}
			return &Float{Value: res}
		},
	}
	definePublicFunction(env, "pangkat", pangkatFn)
	definePublicFunction(env, "pow", pangkatFn)

	// mutlak / abs(x)
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

	// bulat / round(x, [desimal=0])
	bulatFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 1 || len(args) > 2 {
				return newError(nil, "matematika.bulat mengharapkan 1 atau 2 argumen: angka, [jumlah_desimal]", "Contoh: matematika.bulat(3.14159, 2)")
			}
			_, fVal, _, ok := CoerceToNumeric(args[0])
			if !ok {
				return newError(nil, "argumen matematika.bulat harus angka", "")
			}
			decPlaces := 0
			if len(args) == 2 {
				dVal, _, _, ok2 := CoerceToNumeric(args[1])
				if !ok2 {
					return newError(nil, "argumen desimal pada matematika.bulat harus angka bulat", "")
				}
				decPlaces = int(dVal)
			}
			if decPlaces <= 0 {
				return &Integer{Value: int64(math.Round(fVal))}
			}
			factor := math.Pow(10, float64(decPlaces))
			rounded := math.Round(fVal*factor) / factor
			return &Float{Value: rounded}
		},
	}
	definePublicFunction(env, "bulat", bulatFn)
	definePublicFunction(env, "round", bulatFn)

	// potong_desimal / trunc(x, [desimal=0])
	potongDesimalFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 1 || len(args) > 2 {
				return newError(nil, "matematika.potong_desimal mengharapkan 1 atau 2 argumen: angka, [jumlah_desimal]", "")
			}
			_, fVal, _, ok := CoerceToNumeric(args[0])
			if !ok {
				return newError(nil, "argumen matematika.potong_desimal harus angka", "")
			}
			decPlaces := 0
			if len(args) == 2 {
				dVal, _, _, ok2 := CoerceToNumeric(args[1])
				if !ok2 {
					return newError(nil, "argumen desimal harus angka bulat", "")
				}
				decPlaces = int(dVal)
			}
			if decPlaces <= 0 {
				return &Integer{Value: int64(math.Trunc(fVal))}
			}
			factor := math.Pow(10, float64(decPlaces))
			truncated := math.Trunc(fVal*factor) / factor
			return &Float{Value: truncated}
		},
	}
	definePublicFunction(env, "potong_desimal", potongDesimalFn)
	definePublicFunction(env, "trunc", potongDesimalFn)

	// lantai / floor(x)
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

	// atap / ceil(x)
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

	// jepit / clamp(nilai, min, max)
	jepitFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 3 {
				return newError(nil, "matematika.jepit mengharapkan 3 argumen: nilai, batas_bawah, batas_atas", "Contoh: matematika.jepit(120, 0, 100) -> 100")
			}
			iVal, fVal, isInt1, ok1 := CoerceToNumeric(args[0])
			iMin, fMin, isInt2, ok2 := CoerceToNumeric(args[1])
			iMax, fMax, isInt3, ok3 := CoerceToNumeric(args[2])
			if !ok1 || !ok2 || !ok3 {
				return newError(nil, "seluruh argumen matematika.jepit harus berupa angka", "")
			}
			if fMin > fMax {
				fMin, fMax = fMax, fMin
				iMin, iMax = iMax, iMin
			}
			if isInt1 && isInt2 && isInt3 {
				if iVal < iMin {
					return &Integer{Value: iMin}
				}
				if iVal > iMax {
					return &Integer{Value: iMax}
				}
				return &Integer{Value: iVal}
			}
			if fVal < fMin {
				return &Float{Value: fMin}
			}
			if fVal > fMax {
				return &Float{Value: fMax}
			}
			return &Float{Value: fVal}
		},
	}
	definePublicFunction(env, "jepit", jepitFn)
	definePublicFunction(env, "clamp", jepitFn)

	// tanda / sign(x) -> -1, 0, 1
	tandaFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "matematika.tanda mengharapkan 1 argumen angka", "")
			}
			_, fVal, _, ok := CoerceToNumeric(args[0])
			if !ok {
				return newError(nil, "argumen matematika.tanda harus berupa angka", "")
			}
			if fVal > 0 {
				return &Integer{Value: 1}
			}
			if fVal < 0 {
				return &Integer{Value: -1}
			}
			return &Integer{Value: 0}
		},
	}
	definePublicFunction(env, "tanda", tandaFn)
	definePublicFunction(env, "sign", tandaFn)

	// faktorial / factorial(n)
	faktorialFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "matematika.faktorial mengharapkan 1 argumen bilangan bulat", "")
			}
			iVal, _, _, ok := CoerceToNumeric(args[0])
			if !ok {
				return newError(nil, "argumen matematika.faktorial harus berupa bilangan bulat", "")
			}
			if iVal < 0 {
				return newError(nil, "faktorial tidak terdefinisi untuk bilangan negatif", "")
			}
			if iVal <= 20 {
				var res int64 = 1
				for i := int64(2); i <= iVal; i++ {
					res *= i
				}
				return &Integer{Value: res}
			}
			var res float64 = 1.0
			for i := int64(2); i <= iVal; i++ {
				res *= float64(i)
			}
			return &Float{Value: res}
		},
	}
	definePublicFunction(env, "faktorial", faktorialFn)
	definePublicFunction(env, "factorial", faktorialFn)

	// fpb / gcd(a, b) - Faktor Persekutuan Terbesar
	fpbFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "matematika.fpb mengharapkan 2 argumen bilangan bulat", "Contoh: matematika.fpb(24, 36) -> 12")
			}
			a, _, _, ok1 := CoerceToNumeric(args[0])
			b, _, _, ok2 := CoerceToNumeric(args[1])
			if !ok1 || !ok2 {
				return newError(nil, "argumen matematika.fpb harus berupa bilangan bulat", "")
			}
			return &Integer{Value: gcd(a, b)}
		},
	}
	definePublicFunction(env, "fpb", fpbFn)
	definePublicFunction(env, "gcd", fpbFn)

	// kpk / lcm(a, b) - Kelipatan Persekutuan Terkecil
	kpkFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "matematika.kpk mengharapkan 2 argumen bilangan bulat", "Contoh: matematika.kpk(4, 6) -> 12")
			}
			a, _, _, ok1 := CoerceToNumeric(args[0])
			b, _, _, ok2 := CoerceToNumeric(args[1])
			if !ok1 || !ok2 {
				return newError(nil, "argumen matematika.kpk harus berupa bilangan bulat", "")
			}
			return &Integer{Value: lcm(a, b)}
		},
	}
	definePublicFunction(env, "kpk", kpkFn)
	definePublicFunction(env, "lcm", kpkFn)

	// acak / random([min, max])
	acakFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) == 0 {
				rndMu.Lock()
				val := rnd.Float64()
				rndMu.Unlock()
				return &Float{Value: val}
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
				rndMu.Lock()
				res := minVal + rnd.Int63n(diff)
				rndMu.Unlock()
				return &Integer{Value: res}
			}
			return newError(nil, "matematika.acak mengharapkan 0 argumen atau 2 argumen: acak(min, max)", "")
		},
	}
	definePublicFunction(env, "acak", acakFn)
	definePublicFunction(env, "random", acakFn)

	// maks / max(...angka atau [array])
	maksFn := &Builtin{
		Fn: func(args ...Object) Object {
			nums, intNums, allInts, errObj := extractNumericList(args, "matematika.maks")
			if errObj != nil {
				return errObj
			}
			if allInts {
				maxVal := intNums[0]
				for _, v := range intNums[1:] {
					if v > maxVal {
						maxVal = v
					}
				}
				return &Integer{Value: maxVal}
			}
			maxVal := nums[0]
			for _, v := range nums[1:] {
				if v > maxVal {
					maxVal = v
				}
			}
			return &Float{Value: maxVal}
		},
	}
	definePublicFunction(env, "maks", maksFn)
	definePublicFunction(env, "max", maksFn)

	// min / min(...angka atau [array])
	minFn := &Builtin{
		Fn: func(args ...Object) Object {
			nums, intNums, allInts, errObj := extractNumericList(args, "matematika.min")
			if errObj != nil {
				return errObj
			}
			if allInts {
				minVal := intNums[0]
				for _, v := range intNums[1:] {
					if v < minVal {
						minVal = v
					}
				}
				return &Integer{Value: minVal}
			}
			minVal := nums[0]
			for _, v := range nums[1:] {
				if v < minVal {
					minVal = v
				}
			}
			return &Float{Value: minVal}
		},
	}
	definePublicFunction(env, "min", minFn)

	// -------------------------------------------------------------
	// 2. Trigonometri & Sudut
	// -------------------------------------------------------------

	// sin(rad)
	sinFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "matematika.sin mengharapkan 1 argumen sudut (radian)", "")
			}
			_, fVal, _, ok := CoerceToNumeric(args[0])
			if !ok {
				return newError(nil, "argumen matematika.sin harus angka", "")
			}
			return &Float{Value: math.Sin(fVal)}
		},
	}
	definePublicFunction(env, "sin", sinFn)

	// cos(rad)
	cosFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "matematika.cos mengharapkan 1 argumen sudut (radian)", "")
			}
			_, fVal, _, ok := CoerceToNumeric(args[0])
			if !ok {
				return newError(nil, "argumen matematika.cos harus angka", "")
			}
			return &Float{Value: math.Cos(fVal)}
		},
	}
	definePublicFunction(env, "cos", cosFn)

	// tan(rad)
	tanFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "matematika.tan mengharapkan 1 argumen sudut (radian)", "")
			}
			_, fVal, _, ok := CoerceToNumeric(args[0])
			if !ok {
				return newError(nil, "argumen matematika.tan harus angka", "")
			}
			return &Float{Value: math.Tan(fVal)}
		},
	}
	definePublicFunction(env, "tan", tanFn)

	// asin(x)
	asinFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "matematika.asin mengharapkan 1 argumen angka (-1 sampai 1)", "")
			}
			_, fVal, _, ok := CoerceToNumeric(args[0])
			if !ok {
				return newError(nil, "argumen matematika.asin harus angka", "")
			}
			if fVal < -1 || fVal > 1 {
				return newError(nil, "domain matematika.asin berada di rentang [-1, 1]", "")
			}
			return &Float{Value: math.Asin(fVal)}
		},
	}
	definePublicFunction(env, "asin", asinFn)

	// acos(x)
	acosFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "matematika.acos mengharapkan 1 argumen angka (-1 sampai 1)", "")
			}
			_, fVal, _, ok := CoerceToNumeric(args[0])
			if !ok {
				return newError(nil, "argumen matematika.acos harus angka", "")
			}
			if fVal < -1 || fVal > 1 {
				return newError(nil, "domain matematika.acos berada di rentang [-1, 1]", "")
			}
			return &Float{Value: math.Acos(fVal)}
		},
	}
	definePublicFunction(env, "acos", acosFn)

	// atan(x)
	atanFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "matematika.atan mengharapkan 1 argumen angka", "")
			}
			_, fVal, _, ok := CoerceToNumeric(args[0])
			if !ok {
				return newError(nil, "argumen matematika.atan harus angka", "")
			}
			return &Float{Value: math.Atan(fVal)}
		},
	}
	definePublicFunction(env, "atan", atanFn)

	// atan2(y, x)
	atan2Fn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "matematika.atan2 mengharapkan 2 argumen: y dan x", "")
			}
			_, yVal, _, ok1 := CoerceToNumeric(args[0])
			_, xVal, _, ok2 := CoerceToNumeric(args[1])
			if !ok1 || !ok2 {
				return newError(nil, "argumen matematika.atan2 harus angka", "")
			}
			return &Float{Value: math.Atan2(yVal, xVal)}
		},
	}
	definePublicFunction(env, "atan2", atan2Fn)

	// ke_derajat / to_degrees(radian)
	keDerajatFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "matematika.ke_derajat mengharapkan 1 argumen radian", "")
			}
			_, fVal, _, ok := CoerceToNumeric(args[0])
			if !ok {
				return newError(nil, "argumen harus angka", "")
			}
			return &Float{Value: fVal * (180.0 / math.Pi)}
		},
	}
	definePublicFunction(env, "ke_derajat", keDerajatFn)
	definePublicFunction(env, "to_degrees", keDerajatFn)

	// ke_radian / to_radians(derajat)
	keRadianFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "matematika.ke_radian mengharapkan 1 argumen derajat", "")
			}
			_, fVal, _, ok := CoerceToNumeric(args[0])
			if !ok {
				return newError(nil, "argumen harus angka", "")
			}
			return &Float{Value: fVal * (math.Pi / 180.0)}
		},
	}
	definePublicFunction(env, "ke_radian", keRadianFn)
	definePublicFunction(env, "to_radians", keRadianFn)

	// -------------------------------------------------------------
	// 3. Logaritma & Eksponensial
	// -------------------------------------------------------------

	// log(x) - Logaritma Natural (ln)
	logFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "matematika.log mengharapkan 1 argumen angka positif", "")
			}
			_, fVal, _, ok := CoerceToNumeric(args[0])
			if !ok {
				return newError(nil, "argumen matematika.log harus angka", "")
			}
			if fVal <= 0 {
				return newError(nil, "matematika.log mengharuskan nilai > 0", "")
			}
			return &Float{Value: math.Log(fVal)}
		},
	}
	definePublicFunction(env, "log", logFn)

	// log10(x)
	log10Fn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "matematika.log10 mengharapkan 1 argumen angka positif", "")
			}
			_, fVal, _, ok := CoerceToNumeric(args[0])
			if !ok {
				return newError(nil, "argumen matematika.log10 harus angka", "")
			}
			if fVal <= 0 {
				return newError(nil, "matematika.log10 mengharuskan nilai > 0", "")
			}
			return &Float{Value: math.Log10(fVal)}
		},
	}
	definePublicFunction(env, "log10", log10Fn)

	// log2(x)
	log2Fn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "matematika.log2 mengharapkan 1 argumen angka positif", "")
			}
			_, fVal, _, ok := CoerceToNumeric(args[0])
			if !ok {
				return newError(nil, "argumen matematika.log2 harus angka", "")
			}
			if fVal <= 0 {
				return newError(nil, "matematika.log2 mengharuskan nilai > 0", "")
			}
			return &Float{Value: math.Log2(fVal)}
		},
	}
	definePublicFunction(env, "log2", log2Fn)

	// eksponen / exp(x) -> e^x
	expFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "matematika.eksponen mengharapkan 1 argumen angka", "")
			}
			_, fVal, _, ok := CoerceToNumeric(args[0])
			if !ok {
				return newError(nil, "argumen matematika.eksponen harus angka", "")
			}
			return &Float{Value: math.Exp(fVal)}
		},
	}
	definePublicFunction(env, "eksponen", expFn)
	definePublicFunction(env, "exp", expFn)

	// -------------------------------------------------------------
	// 4. Statistik & Analisis Kumpulan Angka (Array & Varargs)
	// -------------------------------------------------------------

	// jumlah / sum(daftar atau ...angka)
	jumlahFn := &Builtin{
		Fn: func(args ...Object) Object {
			nums, intNums, allInts, errObj := extractNumericList(args, "matematika.jumlah")
			if errObj != nil {
				return errObj
			}
			if allInts {
				var total int64
				for _, v := range intNums {
					total += v
				}
				return &Integer{Value: total}
			}
			var total float64
			for _, v := range nums {
				total += v
			}
			return &Float{Value: total}
		},
	}
	definePublicFunction(env, "jumlah", jumlahFn)
	definePublicFunction(env, "sum", jumlahFn)

	// rata_rata / average / mean(daftar atau ...angka)
	rataRataFn := &Builtin{
		Fn: func(args ...Object) Object {
			nums, _, _, errObj := extractNumericList(args, "matematika.rata_rata")
			if errObj != nil {
				return errObj
			}
			var total float64
			for _, v := range nums {
				total += v
			}
			avg := total / float64(len(nums))
			if avg == math.Trunc(avg) {
				return &Integer{Value: int64(avg)}
			}
			return &Float{Value: avg}
		},
	}
	definePublicFunction(env, "rata_rata", rataRataFn)
	definePublicFunction(env, "average", rataRataFn)
	definePublicFunction(env, "mean", rataRataFn)

	// median(daftar atau ...angka)
	medianFn := &Builtin{
		Fn: func(args ...Object) Object {
			nums, _, _, errObj := extractNumericList(args, "matematika.median")
			if errObj != nil {
				return errObj
			}
			sort.Float64s(nums)
			n := len(nums)
			if n%2 == 1 {
				mid := nums[n/2]
				if mid == math.Trunc(mid) {
					return &Integer{Value: int64(mid)}
				}
				return &Float{Value: mid}
			}
			mid := (nums[(n/2)-1] + nums[n/2]) / 2.0
			if mid == math.Trunc(mid) {
				return &Integer{Value: int64(mid)}
			}
			return &Float{Value: mid}
		},
	}
	definePublicFunction(env, "median", medianFn)

	// modus / mode(daftar atau ...angka)
	modusFn := &Builtin{
		Fn: func(args ...Object) Object {
			nums, _, _, errObj := extractNumericList(args, "matematika.modus")
			if errObj != nil {
				return errObj
			}
			counts := map[float64]int{}
			maxCount := 0
			for _, v := range nums {
				counts[v]++
				if counts[v] > maxCount {
					maxCount = counts[v]
				}
			}
			var modes []float64
			for v, c := range counts {
				if c == maxCount {
					modes = append(modes, v)
				}
			}
			sort.Float64s(modes)
			if len(modes) == 1 {
				m := modes[0]
				if m == math.Trunc(m) {
					return &Integer{Value: int64(m)}
				}
				return &Float{Value: m}
			}
			elems := make([]Object, len(modes))
			for i, m := range modes {
				if m == math.Trunc(m) {
					elems[i] = &Integer{Value: int64(m)}
				} else {
					elems[i] = &Float{Value: m}
				}
			}
			return &Array{Elements: elems}
		},
	}
	definePublicFunction(env, "modus", modusFn)
	definePublicFunction(env, "mode", modusFn)

	// standar_deviasi / std_dev(daftar, [isSample=false])
	standarDeviasiFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) == 0 {
				return newError(nil, "matematika.standar_deviasi membutuhkan argumen daftar angka", "")
			}
			isSample := false
			var listArgs []Object
			if arr, ok := args[0].(*Array); ok {
				listArgs = arr.Elements
				if len(args) >= 2 {
					if b, ok := args[1].(*Boolean); ok {
						isSample = b.Value
					}
				}
			} else {
				listArgs = args
			}
			nums, _, _, errObj := extractNumericList(listArgs, "matematika.standar_deviasi")
			if errObj != nil {
				return errObj
			}
			n := float64(len(nums))
			if isSample && n < 2 {
				return newError(nil, "standar deviasi sampel membutuhkan minimal 2 angka", "")
			}
			var total float64
			for _, v := range nums {
				total += v
			}
			mean := total / n
			var sumSq float64
			for _, v := range nums {
				diff := v - mean
				sumSq += diff * diff
			}
			divisor := n
			if isSample {
				divisor = n - 1
			}
			variance := sumSq / divisor
			return &Float{Value: math.Sqrt(variance)}
		},
	}
	definePublicFunction(env, "standar_deviasi", standarDeviasiFn)
	definePublicFunction(env, "std_dev", standarDeviasiFn)

	// -------------------------------------------------------------
	// 5. Fitur Khas & Human-Centric Celang
	// -------------------------------------------------------------

	// terbilang / spell_number(angka)
	terbilangFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "matematika.terbilang mengharapkan 1 argumen angka", "Contoh: matematika.terbilang(1250000)")
			}
			_, fVal, _, ok := CoerceToNumeric(args[0])
			if !ok {
				return newError(nil, "argumen matematika.terbilang harus berupa angka", "")
			}
			return &String{Value: spellNumber(fVal)}
		},
	}
	definePublicFunction(env, "terbilang", terbilangFn)
	definePublicFunction(env, "spell_number", terbilangFn)

	// rupiah / format_rupiah(angka, [dengan_desimal=false])
	rupiahFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 1 || len(args) > 2 {
				return newError(nil, "matematika.rupiah mengharapkan 1 atau 2 argumen: angka, [dengan_desimal]", "Contoh: matematika.rupiah(1500000)")
			}
			_, fVal, _, ok := CoerceToNumeric(args[0])
			if !ok {
				return newError(nil, "argumen pertama matematika.rupiah harus berupa angka", "")
			}
			withDecimals := false
			if len(args) == 2 {
				if b, ok := args[1].(*Boolean); ok {
					withDecimals = b.Value
				}
			}
			return &String{Value: formatRupiah(fVal, withDecimals)}
		},
	}
	definePublicFunction(env, "rupiah", rupiahFn)
	definePublicFunction(env, "format_rupiah", rupiahFn)

	// persen / percent(bagian, total) -> menghitung persentase
	persenFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "matematika.persen mengharapkan 2 argumen: bagian dan total", "Contoh: matematika.persen(25, 200) -> 12.5")
			}
			_, fPart, _, ok1 := CoerceToNumeric(args[0])
			_, fTotal, _, ok2 := CoerceToNumeric(args[1])
			if !ok1 || !ok2 {
				return newError(nil, "argumen matematika.persen harus berupa angka", "")
			}
			if fTotal == 0 {
				return newError(nil, "pembagian dengan nol pada matematika.persen", "")
			}
			res := (fPart / fTotal) * 100.0
			if res == math.Trunc(res) && res <= float64(math.MaxInt64) && res >= float64(math.MinInt64) {
				return &Integer{Value: int64(res)}
			}
			return &Float{Value: res}
		},
	}
	definePublicFunction(env, "persen", persenFn)
	definePublicFunction(env, "percent", persenFn)

	// dari_persen / from_percent(persen, total) -> menghitung nilai potongan / porsi
	dariPersenFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "matematika.dari_persen mengharapkan 2 argumen: persen dan total", "Contoh: matematika.dari_persen(15, 200000) -> 30000")
			}
			_, fPercent, _, ok1 := CoerceToNumeric(args[0])
			_, fTotal, _, ok2 := CoerceToNumeric(args[1])
			if !ok1 || !ok2 {
				return newError(nil, "argumen matematika.dari_persen harus berupa angka", "")
			}
			res := (fPercent / 100.0) * fTotal
			if res == math.Trunc(res) && res <= float64(math.MaxInt64) && res >= float64(math.MinInt64) {
				return &Integer{Value: int64(res)}
			}
			return &Float{Value: res}
		},
	}
	definePublicFunction(env, "dari_persen", dariPersenFn)
	definePublicFunction(env, "from_percent", dariPersenFn)

	// lerp(a, b, t) - Linear Interpolation
	lerpFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 3 {
				return newError(nil, "matematika.lerp mengharapkan 3 argumen: awal, akhir, faktor_t (0.0 .. 1.0)", "Contoh: matematika.lerp(0, 100, 0.5) -> 50")
			}
			_, aVal, isIntA, ok1 := CoerceToNumeric(args[0])
			_, bVal, isIntB, ok2 := CoerceToNumeric(args[1])
			_, tVal, isIntT, ok3 := CoerceToNumeric(args[2])
			if !ok1 || !ok2 || !ok3 {
				return newError(nil, "seluruh argumen matematika.lerp harus berupa angka", "")
			}
			res := aVal + (bVal-aVal)*tVal
			if isIntA && isIntB && isIntT && res == math.Trunc(res) && res <= float64(math.MaxInt64) && res >= float64(math.MinInt64) {
				return &Integer{Value: int64(res)}
			}
			return &Float{Value: res}
		},
	}
	definePublicFunction(env, "lerp", lerpFn)

	// petakan_rentang / map_range(nilai, inMin, inMax, outMin, outMax)
	petakanRentangFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 5 {
				return newError(nil, "matematika.petakan_rentang mengharapkan 5 argumen: nilai, in_min, in_max, out_min, out_max", "Contoh: matematika.petakan_rentang(50, 0, 100, 0, 1) -> 0.5")
			}
			_, val, _, ok1 := CoerceToNumeric(args[0])
			_, inMin, _, ok2 := CoerceToNumeric(args[1])
			_, inMax, _, ok3 := CoerceToNumeric(args[2])
			_, outMin, _, ok4 := CoerceToNumeric(args[3])
			_, outMax, _, ok5 := CoerceToNumeric(args[4])
			if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 {
				return newError(nil, "seluruh argumen matematika.petakan_rentang harus berupa angka", "")
			}
			if inMax == inMin {
				return newError(nil, "rentang masukan tidak boleh sama (pembagian dengan nol)", "")
			}
			res := outMin + (val-inMin)*(outMax-outMin)/(inMax-inMin)
			if res == math.Trunc(res) && res <= float64(math.MaxInt64) && res >= float64(math.MinInt64) {
				return &Integer{Value: int64(res)}
			}
			return &Float{Value: res}
		},
	}
	definePublicFunction(env, "petakan_rentang", petakanRentangFn)
	definePublicFunction(env, "map_range", petakanRentangFn)

	// jarak / distance / hypot(x1, y1, [z1], x2, y2, [z2])
	jarakFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) == 4 {
				_, x1, _, ok1 := CoerceToNumeric(args[0])
				_, y1, _, ok2 := CoerceToNumeric(args[1])
				_, x2, _, ok3 := CoerceToNumeric(args[2])
				_, y2, _, ok4 := CoerceToNumeric(args[3])
				if !ok1 || !ok2 || !ok3 || !ok4 {
					return newError(nil, "koordinat matematika.jarak harus berupa angka", "")
				}
				dist := math.Hypot(x2-x1, y2-y1)
				if dist == math.Trunc(dist) {
					return &Integer{Value: int64(dist)}
				}
				return &Float{Value: dist}
			}
			if len(args) == 6 {
				_, x1, _, ok1 := CoerceToNumeric(args[0])
				_, y1, _, ok2 := CoerceToNumeric(args[1])
				_, z1, _, ok3 := CoerceToNumeric(args[2])
				_, x2, _, ok4 := CoerceToNumeric(args[3])
				_, y2, _, ok5 := CoerceToNumeric(args[4])
				_, z2, _, ok6 := CoerceToNumeric(args[5])
				if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 {
					return newError(nil, "koordinat matematika.jarak 3D harus berupa angka", "")
				}
				dx := x2 - x1
				dy := y2 - y1
				dz := z2 - z1
				dist := math.Sqrt(dx*dx + dy*dy + dz*dz)
				if dist == math.Trunc(dist) {
					return &Integer{Value: int64(dist)}
				}
				return &Float{Value: dist}
			}
			return newError(nil, "matematika.jarak membutuhkan 4 argumen (2D: x1, y1, x2, y2) atau 6 argumen (3D: x1, y1, z1, x2, y2, z2)", "Contoh: matematika.jarak(0, 0, 3, 4) -> 5")
		},
	}
	definePublicFunction(env, "jarak", jarakFn)
	definePublicFunction(env, "distance", jarakFn)
	definePublicFunction(env, "hypot", jarakFn)

	// -------------------------------------------------------------
	// 6. Konstanta Matematika
	// -------------------------------------------------------------
	definePublicConstant(env, "PI", &Float{Value: math.Pi})
	definePublicConstant(env, "E", &Float{Value: math.E})
	definePublicConstant(env, "PHI", &Float{Value: 1.618033988749895})
	definePublicConstant(env, "RASIO_EMAS", &Float{Value: 1.618033988749895})
	definePublicConstant(env, "INF", &Float{Value: math.Inf(1)})
	definePublicConstant(env, "TAK_HINGGA", &Float{Value: math.Inf(1)})
}

// -------------------------------------------------------------
// Helper Functions (Internal Math, Formatting, Speller)
// -------------------------------------------------------------

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func lcm(a, b int64) int64 {
	if a == 0 || b == 0 {
		return 0
	}
	g := gcd(a, b)
	res := (a / g) * b
	if res < 0 {
		return -res
	}
	return res
}

func extractNumericList(args []Object, funcName string) ([]float64, []int64, bool, Object) {
	var items []Object
	if len(args) == 1 {
		if arr, ok := args[0].(*Array); ok {
			items = arr.Elements
		} else {
			items = args
		}
	} else {
		items = args
	}

	if len(items) == 0 {
		return nil, nil, false, newError(nil, fmt.Sprintf("%s membutuhkan minimal 1 elemen angka", funcName), "")
	}

	var nums []float64
	var intNums []int64
	allInts := true

	for _, it := range items {
		iVal, fVal, isInt, ok := CoerceToNumeric(it)
		if !ok {
			return nil, nil, false, newError(nil, fmt.Sprintf("elemen dalam %s harus berupa angka, menerima %s", funcName, it.Type()), "")
		}
		nums = append(nums, fVal)
		intNums = append(intNums, iVal)
		if !isInt {
			allInts = false
		}
	}
	return nums, intNums, allInts, nil
}

func formatRupiah(val float64, withDecimals bool) string {
	isNegative := val < 0
	absVal := math.Abs(val)
	intPart := int64(absVal)

	s := strconv.FormatInt(intPart, 10)
	var b strings.Builder
	n := len(s)
	for i := 0; i < n; i++ {
		if i > 0 && (n-i)%3 == 0 {
			b.WriteByte('.')
		}
		b.WriteByte(s[i])
	}
	formattedInt := b.String()

	var result string
	if withDecimals {
		dec := int64(math.Round((absVal - float64(intPart)) * 100))
		if dec >= 100 {
			dec = 0
		}
		result = fmt.Sprintf("Rp %s,%02d", formattedInt, dec)
	} else {
		result = fmt.Sprintf("Rp %s", formattedInt)
	}

	if isNegative {
		result = "-" + result
	}
	return result
}

var terbilangKata = []string{
	"", "satu", "dua", "tiga", "empat", "lima",
	"enam", "tujuh", "delapan", "sembilan", "sepuluh", "sebelas",
}

var digitWords = []string{
	"nol", "satu", "dua", "tiga", "empat", "lima",
	"enam", "tujuh", "delapan", "sembilan",
}

func spellIndonesian(n int64) string {
	if n == 0 {
		return "nol"
	}
	if n < 0 {
		return "minus " + spellIndonesian(-n)
	}
	if n < 12 {
		return terbilangKata[n]
	}
	if n < 20 {
		return spellIndonesian(n-10) + " belas"
	}
	if n < 100 {
		res := spellIndonesian(n/10) + " puluh"
		if rem := n % 10; rem > 0 {
			res += " " + spellIndonesian(rem)
		}
		return res
	}
	if n < 200 {
		res := "seratus"
		if rem := n - 100; rem > 0 {
			res += " " + spellIndonesian(rem)
		}
		return res
	}
	if n < 1000 {
		res := spellIndonesian(n/100) + " ratus"
		if rem := n % 100; rem > 0 {
			res += " " + spellIndonesian(rem)
		}
		return res
	}
	if n < 2000 {
		res := "seribu"
		if rem := n - 1000; rem > 0 {
			res += " " + spellIndonesian(rem)
		}
		return res
	}
	if n < 1000000 {
		res := spellIndonesian(n/1000) + " ribu"
		if rem := n % 1000; rem > 0 {
			res += " " + spellIndonesian(rem)
		}
		return res
	}
	if n < 1000000000 {
		res := spellIndonesian(n/1000000) + " juta"
		if rem := n % 1000000; rem > 0 {
			res += " " + spellIndonesian(rem)
		}
		return res
	}
	if n < 1000000000000 {
		res := spellIndonesian(n/1000000000) + " miliar"
		if rem := n % 1000000000; rem > 0 {
			res += " " + spellIndonesian(rem)
		}
		return res
	}
	if n < 1000000000000000 {
		res := spellIndonesian(n/1000000000000) + " triliun"
		if rem := n % 1000000000000; rem > 0 {
			res += " " + spellIndonesian(rem)
		}
		return res
	}
	return fmt.Sprintf("%d", n)
}

func spellNumber(val float64) string {
	if val == math.Trunc(val) {
		return spellIndonesian(int64(val))
	}
	isNeg := val < 0
	absVal := math.Abs(val)
	intPart := int64(math.Trunc(absVal))

	valStr := strconv.FormatFloat(absVal, 'f', -1, 64)
	parts := strings.Split(valStr, ".")

	var spelledInt string
	if intPart == 0 {
		spelledInt = "nol"
	} else {
		spelledInt = spellIndonesian(intPart)
	}

	var decWords []string
	if len(parts) > 1 {
		for _, ch := range parts[1] {
			if ch >= '0' && ch <= '9' {
				decWords = append(decWords, digitWords[ch-'0'])
			}
		}
	}

	res := spelledInt + " koma " + strings.Join(decWords, " ")
	if isNeg {
		res = "minus " + res
	}
	return res
}
