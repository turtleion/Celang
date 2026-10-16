package evaluator

import (
	"fmt"
	"math"
	"strings"
)

var satuanKata = []string{
	"nol", "satu", "dua", "tiga", "empat", "lima",
	"enam", "tujuh", "delapan", "sembilan", "sepuluh", "sebelas",
}

// TerbilangAngka mengonversi bilangan bulat (int64) ke dalam kata-kata bahasa Indonesia.
// Contoh: 15 -> "lima belas", 1000 -> "seribu", 50000 -> "lima puluh ribu"
func TerbilangAngka(n int64) string {
	if n == math.MinInt64 {
		return "minus sembilan kuantiliun dua ratus dua puluh tiga triliun tiga ratus tujuh puluh dua miliar tiga puluh enam juta delapan ratus lima puluh empat ribu tujuh ratus tujuh puluh delapan"
	}
	if n < 0 {
		return "minus " + TerbilangAngka(-n)
	}
	if n <= 11 {
		return satuanKata[n]
	}
	if n < 20 {
		return satuanKata[n-10] + " belas"
	}
	if n < 100 {
		hasil := satuanKata[n/10] + " puluh"
		sisa := n % 10
		if sisa != 0 {
			hasil += " " + satuanKata[sisa]
		}
		return hasil
	}
	if n < 200 {
		hasil := "seratus"
		sisa := n - 100
		if sisa != 0 {
			hasil += " " + TerbilangAngka(sisa)
		}
		return hasil
	}
	if n < 1000 {
		hasil := satuanKata[n/100] + " ratus"
		sisa := n % 100
		if sisa != 0 {
			hasil += " " + TerbilangAngka(sisa)
		}
		return hasil
	}
	if n < 2000 {
		hasil := "seribu"
		sisa := n - 1000
		if sisa != 0 {
			hasil += " " + TerbilangAngka(sisa)
		}
		return hasil
	}
	if n < 1000000 {
		hasil := TerbilangAngka(n/1000) + " ribu"
		sisa := n % 1000
		if sisa != 0 {
			hasil += " " + TerbilangAngka(sisa)
		}
		return hasil
	}
	if n < 1000000000 {
		hasil := TerbilangAngka(n/1000000) + " juta"
		sisa := n % 1000000
		if sisa != 0 {
			hasil += " " + TerbilangAngka(sisa)
		}
		return hasil
	}
	if n < 1000000000000 {
		hasil := TerbilangAngka(n/1000000000) + " miliar"
		sisa := n % 1000000000
		if sisa != 0 {
			hasil += " " + TerbilangAngka(sisa)
		}
		return hasil
	}
	hasil := TerbilangAngka(n/1000000000000) + " triliun"
	sisa := n % 1000000000000
	if sisa != 0 {
		hasil += " " + TerbilangAngka(sisa)
	}
	return hasil
}

// Terbilang mengonversi Object numerik (Integer atau Float) ke ejaan kata lengkap (dengan satuan "rupiah").
func Terbilang(obj Object) string {
	intVal, floatVal, isInt, ok := CoerceToNumeric(obj)
	if !ok {
		return ""
	}
	if isInt {
		return TerbilangAngka(intVal) + " rupiah"
	}

	// Bilangan berkoma (Float)
	prefix := ""
	if floatVal < 0 {
		prefix = "minus "
		floatVal = -floatVal
	}
	intPart := int64(math.Trunc(floatVal))
	decPart := math.Abs(floatVal - float64(intPart))
	str := prefix + TerbilangAngka(intPart)

	if decPart > 0.0001 {
		decStr := fmt.Sprintf("%.2f", decPart)
		decStr = strings.TrimPrefix(decStr, "0.")
		str += " koma"
		for _, ch := range decStr {
			digit := int(ch - '0')
			if digit >= 0 && digit <= 9 {
				str += " " + satuanKata[digit]
			}
		}
	}
	return str + " rupiah"
}

// KataAngka mengonversi Object numerik ke kata tanpa imbuhan "rupiah".
func KataAngka(obj Object) string {
	intVal, floatVal, isInt, ok := CoerceToNumeric(obj)
	if !ok {
		return ""
	}
	if isInt {
		return TerbilangAngka(intVal)
	}

	intPart := int64(math.Trunc(floatVal))
	decPart := math.Abs(floatVal - float64(intPart))
	prefix := ""
	if floatVal < 0 {
		prefix = "minus "
		floatVal = -floatVal
		intPart = int64(math.Trunc(floatVal))
		decPart = math.Abs(floatVal - float64(intPart))
	}
	str := prefix + TerbilangAngka(intPart)

	if decPart > 0.0001 {
		decStr := fmt.Sprintf("%.2f", decPart)
		decStr = strings.TrimPrefix(decStr, "0.")
		str += " koma"
		for _, ch := range decStr {
			digit := int(ch - '0')
			if digit >= 0 && digit <= 9 {
				str += " " + satuanKata[digit]
			}
		}
	}
	return str
}

// FormatRupiah memformat angka ke format standar mata uang Indonesia: "Rp 50.000" atau "Rp 12.500,50".
func FormatRupiah(obj Object) string {
	intVal, floatVal, isInt, ok := CoerceToNumeric(obj)
	if !ok {
		return ""
	}

	if isInt {
		if intVal < 0 {
			if intVal == math.MinInt64 {
				return "-Rp 9.223.372.036.854.775.808"
			}
			return "-Rp " + formatRibuan(-intVal)
		}
		return "Rp " + formatRibuan(intVal)
	}

	sign := ""
	if floatVal < 0 {
		sign = "-"
		floatVal = -floatVal
	}

	intPart := int64(math.Trunc(floatVal))
	decPart := math.Abs(floatVal - float64(intPart))
	if decPart < 0.0001 {
		return sign + "Rp " + formatRibuan(intPart)
	}
	decStr := fmt.Sprintf("%.2f", decPart)
	decStr = strings.TrimPrefix(decStr, "0.")
	return sign + "Rp " + formatRibuan(intPart) + "," + decStr
}

func formatRibuan(n int64) string {
	if n == math.MinInt64 {
		return "-9.223.372.036.854.775.808"
	}
	sign := ""
	if n < 0 {
		sign = "-"
		n = -n
	}
	s := fmt.Sprintf("%d", n)
	if len(s) <= 3 {
		return sign + s
	}

	var parts []string
	for len(s) > 3 {
		parts = append([]string{s[len(s)-3:]}, parts...)
		s = s[:len(s)-3]
	}
	if len(s) > 0 {
		parts = append([]string{s}, parts...)
	}
	return sign + strings.Join(parts, ".")
}

var (
	onesEnglish = []string{
		"zero", "one", "two", "three", "four", "five",
		"six", "seven", "eight", "nine", "ten", "eleven",
		"twelve", "thirteen", "fourteen", "fifteen", "sixteen",
		"seventeen", "eighteen", "nineteen",
	}
	tensEnglish = []string{
		"", "", "twenty", "thirty", "forty", "fifty",
		"sixty", "seventy", "eighty", "ninety",
	}
)

// WordsEnglish mengonversi integer ke dalam kata-kata bahasa Inggris.
func WordsEnglish(n int64) string {
	if n == math.MinInt64 {
		return "minus nine quintillion two hundred twenty-three quadrillion three hundred seventy-two trillion thirty-six billion eight hundred fifty-four million seven hundred seventy-five thousand eight hundred eight"
	}
	if n < 0 {
		return "minus " + WordsEnglish(-n)
	}
	if n < 20 {
		return onesEnglish[n]
	}
	if n < 100 {
		res := tensEnglish[n/10]
		rem := n % 10
		if rem != 0 {
			res += "-" + onesEnglish[rem]
		}
		return res
	}
	if n < 1000 {
		res := onesEnglish[n/100] + " hundred"
		rem := n % 100
		if rem != 0 {
			res += " " + WordsEnglish(rem)
		}
		return res
	}
	if n < 1000000 {
		res := WordsEnglish(n/1000) + " thousand"
		rem := n % 1000
		if rem != 0 {
			res += " " + WordsEnglish(rem)
		}
		return res
	}
	if n < 1000000000 {
		res := WordsEnglish(n/1000000) + " million"
		rem := n % 1000000
		if rem != 0 {
			res += " " + WordsEnglish(rem)
		}
		return res
	}
	if n < 1000000000000 {
		res := WordsEnglish(n/1000000000) + " billion"
		rem := n % 1000000000
		if rem != 0 {
			res += " " + WordsEnglish(rem)
		}
		return res
	}
	res := WordsEnglish(n/1000000000000) + " trillion"
	rem := n % 1000000000000
	if rem != 0 {
		res += " " + WordsEnglish(rem)
	}
	return res
}

// WordsEnglishCoerced mengonversi bilangan numerik apapun (Integer / Float) ke bahasa Inggris.
func WordsEnglishCoerced(obj Object) string {
	intVal, floatVal, isInt, ok := CoerceToNumeric(obj)
	if !ok {
		return ""
	}
	if isInt {
		return WordsEnglish(intVal)
	}
	isNeg := floatVal < 0
	val := math.Abs(floatVal)
	intPart := int64(val)
	fracPart := val - float64(intPart)

	res := WordsEnglish(intPart)
	if fracPart > 0.000001 {
		res += " point"
		fracStr := fmt.Sprintf("%.2f", fracPart)
		if strings.HasPrefix(fracStr, "0.") {
			fracDigits := fracStr[2:]
			for _, d := range fracDigits {
				dig := int(d - '0')
				if dig >= 0 && dig <= 9 {
					res += " " + onesEnglish[dig]
				}
			}
		}
	}
	if isNeg {
		res = "minus " + res
	}
	return res
}

// DollarsInWords mengonversi nilai numerik ke kalimat mata uang US Dollar dalam bahasa Inggris.
func DollarsInWords(obj Object) string {
	intVal, floatVal, isInt, ok := CoerceToNumeric(obj)
	if !ok {
		return "zero dollars"
	}
	if isInt {
		sign := ""
		if intVal < 0 {
			sign = "minus "
			intVal = -intVal
		}
		if intVal == 1 {
			return sign + "one dollar"
		}
		return sign + WordsEnglish(intVal) + " dollars"
	}

	isNeg := floatVal < 0
	val := math.Abs(floatVal)
	dollars := int64(val)
	cents := int64(math.Round((val - float64(dollars)) * 100))
	if cents >= 100 {
		dollars++
		cents -= 100
	}

	var parts []string
	if dollars == 1 {
		parts = append(parts, "one dollar")
	} else if dollars > 0 || cents == 0 {
		parts = append(parts, WordsEnglish(dollars)+" dollars")
	}

	if cents > 0 {
		if cents == 1 {
			parts = append(parts, "one cent")
		} else {
			parts = append(parts, WordsEnglish(cents)+" cents")
		}
	}

	sign := ""
	if isNeg {
		sign = "minus "
	}
	return sign + strings.Join(parts, " and ")
}

// formatThousandsCommas memformat angka dengan pemisah koma internasional (e.g. 1,000,000).
func formatThousandsCommas(n int64) string {
	if n == math.MinInt64 {
		return "-9,223,372,036,854,775,808"
	}
	sign := ""
	if n < 0 {
		sign = "-"
		n = -n
	}
	s := fmt.Sprintf("%d", n)
	if len(s) <= 3 {
		return sign + s
	}

	var parts []string
	for len(s) > 3 {
		parts = append([]string{s[len(s)-3:]}, parts...)
		s = s[:len(s)-3]
	}
	if len(s) > 0 {
		parts = append([]string{s}, parts...)
	}
	return sign + strings.Join(parts, ",")
}

// FormatUSD memformat nilai uang dalam standar internasional US Dollar ($1,250.00).
func FormatUSD(obj Object) string {
	intVal, floatVal, isInt, ok := CoerceToNumeric(obj)
	if !ok {
		return "$0.00"
	}
	if isInt {
		sign := ""
		if intVal < 0 {
			sign = "-"
			intVal = -intVal
		}
		return fmt.Sprintf("%s$%s.00", sign, formatThousandsCommas(intVal))
	}
	isNeg := floatVal < 0
	val := math.Abs(floatVal)
	dollars := int64(val)
	cents := int64(math.Round((val - float64(dollars)) * 100))
	if cents >= 100 {
		dollars++
		cents -= 100
	}
	sign := ""
	if isNeg {
		sign = "-"
	}
	return fmt.Sprintf("%s$%s.%02d", sign, formatThousandsCommas(dollars), cents)
}

// EvalNumericMember mengevaluasi pengaksesan member helper pada tipe numerik (Integer & Float).
func EvalNumericMember(obj Object, prop string) (Object, bool) {
	return EvalNumericMemberWithEnv(obj, prop, nil)
}

// EvalNumericMemberWithEnv mengevaluasi pengaksesan member helper numerik dengan konteks Environment.
func EvalNumericMemberWithEnv(obj Object, prop string, env *Environment) (Object, bool) {
	intVal, floatVal, isInt, ok := CoerceToNumeric(obj)
	if !ok {
		return nil, false
	}

	isUniversal := env != nil && env.IsUniversal

	// Guard against Indonesian suffixes in Universal Mode
	if isUniversal {
		switch prop {
		case "ribu":
			return &Error{Message: "Indonesian suffix 'ribu' is not allowed in Universal Mode. Use '.thousand' or '.k'."}, true
		case "juta":
			return &Error{Message: "Indonesian suffix 'juta' is not allowed in Universal Mode. Use '.million'."}, true
		case "miliar", "milyar":
			return &Error{Message: "Indonesian suffix 'miliar' is not allowed in Universal Mode. Use '.billion'."}, true
		case "triliun", "trilyun":
			return &Error{Message: "Indonesian suffix 'triliun' is not allowed in Universal Mode. Use '.trillion'."}, true
		case "rupiah":
			return &Error{Message: "Indonesian suffix 'rupiah' is not allowed in Universal Mode. Use '.usd' or '.dollars'."}, true
		case "terbilang", "kata", "ejaan":
			return &Error{Message: "Indonesian suffix 'terbilang' is not allowed in Universal Mode. Use '.words' or '.dollars_in_words'."}, true
		case "milidetik":
			return &Error{Message: "Indonesian suffix 'milidetik' is not allowed in Universal Mode. Use '.ms'."}, true
		case "detik":
			return &Error{Message: "Indonesian suffix 'detik' is not allowed in Universal Mode. Use '.sec', '.seconds', or '.s'."}, true
		case "menit":
			return &Error{Message: "Indonesian suffix 'menit' is not allowed in Universal Mode. Use '.min', '.minutes', or '.m'."}, true
		case "jam":
			return &Error{Message: "Indonesian suffix 'jam' is not allowed in Universal Mode. Use '.hour', '.hours', or '.h'."}, true
		case "hari":
			return &Error{Message: "Indonesian suffix 'hari' is not allowed in Universal Mode. Use '.day', '.days', or '.d'."}, true
		}
	}

	switch prop {
	// 1. Satuan Nominal / Jumlah (Lokal & Internasional)
	case "ribu", "thousand", "k":
		if isInt {
			return &Integer{Value: intVal * 1000}, true
		}
		res := floatVal * 1000
		if res == math.Trunc(res) {
			return &Integer{Value: int64(res)}, true
		}
		return &Float{Value: res}, true

	case "juta", "million":
		if isInt {
			return &Integer{Value: intVal * 1000000}, true
		}
		res := floatVal * 1000000
		if res == math.Trunc(res) {
			return &Integer{Value: int64(res)}, true
		}
		return &Float{Value: res}, true

	case "miliar", "milyar", "billion":
		if isInt {
			return &Integer{Value: intVal * 1000000000}, true
		}
		res := floatVal * 1000000000
		if res == math.Trunc(res) {
			return &Integer{Value: int64(res)}, true
		}
		return &Float{Value: res}, true

	case "triliun", "trilyun", "trillion", "t":
		if isInt {
			return &Integer{Value: intVal * 1000000000000}, true
		}
		res := floatVal * 1000000000000
		if res == math.Trunc(res) {
			return &Integer{Value: int64(res)}, true
		}
		return &Float{Value: res}, true

	// 2. Format & Terbilang (Lokal & Internasional)
	case "rupiah":
		return &String{Value: FormatRupiah(obj)}, true

	case "terbilang":
		return &String{Value: Terbilang(obj)}, true

	case "kata", "ejaan":
		return &String{Value: KataAngka(obj)}, true

	case "usd", "dollar", "dollars":
		return &String{Value: FormatUSD(obj)}, true

	case "words", "in_words", "spelled_out":
		return &String{Value: WordsEnglishCoerced(obj)}, true

	case "dollars_in_words":
		return &String{Value: DollarsInWords(obj)}, true

	// 3. Satuan Ukuran Data (Bytes)
	case "b", "byte", "bytes":
		if isInt {
			return &Integer{Value: intVal}, true
		}
		return &Float{Value: floatVal}, true

	case "kb", "kilobyte":
		const mul = 1024
		if isInt {
			return &Integer{Value: intVal * mul}, true
		}
		res := floatVal * mul
		if res == math.Trunc(res) {
			return &Integer{Value: int64(res)}, true
		}
		return &Float{Value: res}, true

	case "mb", "megabyte":
		const mul = 1024 * 1024
		if isInt {
			return &Integer{Value: intVal * mul}, true
		}
		res := floatVal * mul
		if res == math.Trunc(res) {
			return &Integer{Value: int64(res)}, true
		}
		return &Float{Value: res}, true

	case "gb", "gigabyte":
		const mul = 1024 * 1024 * 1024
		if isInt {
			return &Integer{Value: intVal * mul}, true
		}
		res := floatVal * mul
		if res == math.Trunc(res) {
			return &Integer{Value: int64(res)}, true
		}
		return &Float{Value: res}, true

	case "tb", "terabyte":
		const mul = 1024 * 1024 * 1024 * 1024
		if isInt {
			return &Integer{Value: intVal * mul}, true
		}
		res := floatVal * mul
		if res == math.Trunc(res) {
			return &Integer{Value: int64(res)}, true
		}
		return &Float{Value: res}, true

	// 4. Satuan Waktu (Durasi dalam milidetik)
	case "milidetik", "ms":
		if isInt {
			return &Integer{Value: intVal}, true
		}
		return &Float{Value: floatVal}, true

	case "detik", "s", "sec", "second", "seconds":
		if isInt {
			return &Integer{Value: intVal * 1000}, true
		}
		return &Float{Value: floatVal * 1000}, true

	case "menit", "m", "min", "minute", "minutes":
		if isInt {
			return &Integer{Value: intVal * 60 * 1000}, true
		}
		return &Float{Value: floatVal * 60 * 1000}, true

	case "jam", "h", "hour", "hours":
		if isInt {
			return &Integer{Value: intVal * 3600 * 1000}, true
		}
		return &Float{Value: floatVal * 3600 * 1000}, true

	case "hari", "d", "day", "days":
		if isInt {
			return &Integer{Value: intVal * 86400 * 1000}, true
		}
		return &Float{Value: floatVal * 86400 * 1000}, true
	}

	return nil, false
}
