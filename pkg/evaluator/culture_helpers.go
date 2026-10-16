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
	intPart := int64(math.Trunc(floatVal))
	decPart := math.Abs(floatVal - float64(intPart))
	str := TerbilangAngka(intPart)

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
	str := TerbilangAngka(intPart)

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
		return "Rp " + formatRibuan(intVal)
	}

	intPart := int64(math.Trunc(floatVal))
	decPart := math.Abs(floatVal - float64(intPart))
	if decPart < 0.0001 {
		return "Rp " + formatRibuan(intPart)
	}
	decStr := fmt.Sprintf("%.2f", decPart)
	decStr = strings.TrimPrefix(decStr, "0.")
	return "Rp " + formatRibuan(intPart) + "," + decStr
}

func formatRibuan(n int64) string {
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

// EvalNumericMember mengevaluasi pengaksesan member helper pada tipe numerik (Integer & Float).
func EvalNumericMember(obj Object, prop string) (Object, bool) {
	intVal, floatVal, isInt, ok := CoerceToNumeric(obj)
	if !ok {
		return nil, false
	}

	switch prop {
	// 1. Satuan Nominal / Jumlah
	case "ribu":
		if isInt {
			return &Integer{Value: intVal * 1000}, true
		}
		res := floatVal * 1000
		if res == math.Trunc(res) {
			return &Integer{Value: int64(res)}, true
		}
		return &Float{Value: res}, true

	case "juta":
		if isInt {
			return &Integer{Value: intVal * 1000000}, true
		}
		res := floatVal * 1000000
		if res == math.Trunc(res) {
			return &Integer{Value: int64(res)}, true
		}
		return &Float{Value: res}, true

	case "miliar", "milyar":
		if isInt {
			return &Integer{Value: intVal * 1000000000}, true
		}
		res := floatVal * 1000000000
		if res == math.Trunc(res) {
			return &Integer{Value: int64(res)}, true
		}
		return &Float{Value: res}, true

	case "triliun", "trilyun":
		if isInt {
			return &Integer{Value: intVal * 1000000000000}, true
		}
		res := floatVal * 1000000000000
		if res == math.Trunc(res) {
			return &Integer{Value: int64(res)}, true
		}
		return &Float{Value: res}, true

	// 2. Format & Terbilang
	case "rupiah":
		return &String{Value: FormatRupiah(obj)}, true

	case "terbilang":
		return &String{Value: Terbilang(obj)}, true

	case "kata", "ejaan":
		return &String{Value: KataAngka(obj)}, true

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

	case "detik", "s":
		if isInt {
			return &Integer{Value: intVal * 1000}, true
		}
		return &Float{Value: floatVal * 1000}, true

	case "menit", "m":
		if isInt {
			return &Integer{Value: intVal * 60 * 1000}, true
		}
		return &Float{Value: floatVal * 60 * 1000}, true

	case "jam", "h":
		if isInt {
			return &Integer{Value: intVal * 3600 * 1000}, true
		}
		return &Float{Value: floatVal * 3600 * 1000}, true

	case "hari", "d":
		if isInt {
			return &Integer{Value: intVal * 86400 * 1000}, true
		}
		return &Float{Value: floatVal * 86400 * 1000}, true
	}

	return nil, false
}
