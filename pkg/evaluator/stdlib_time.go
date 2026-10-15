package evaluator

import (
	"fmt"
	"math"
	"strings"
	"time"
)

var (
	namaHariIndo = []string{"Minggu", "Senin", "Selasa", "Rabu", "Kamis", "Jumat", "Sabtu"}
	namaBulanIndo = []string{
		"", "Januari", "Februari", "Maret", "April", "Mei", "Juni",
		"Juli", "Agustus", "September", "Oktober", "November", "Desember",
	}
)

func initStdlibWaktu(env *Environment) {
	// -------------------------------------------------------------
	// 1. Waktu Dasar & Pengambilan Saat Ini
	// -------------------------------------------------------------

	// sekarang / now() -> milidetik
	sekarangFn := &Builtin{
		Fn: func(args ...Object) Object {
			ms := time.Now().UnixNano() / int64(time.Millisecond)
			return &Integer{Value: ms}
		},
	}
	definePublicFunction(env, "sekarang", sekarangFn)
	definePublicFunction(env, "now", sekarangFn)

	// sekarang_detik / now_seconds() -> detik
	sekarangDetikFn := &Builtin{
		Fn: func(args ...Object) Object {
			return &Integer{Value: time.Now().Unix()}
		},
	}
	definePublicFunction(env, "sekarang_detik", sekarangDetikFn)
	definePublicFunction(env, "now_seconds", sekarangDetikFn)

	// sekarang_nano / now_nanos() -> nanodetik
	sekarangNanoFn := &Builtin{
		Fn: func(args ...Object) Object {
			return &Integer{Value: time.Now().UnixNano()}
		},
	}
	definePublicFunction(env, "sekarang_nano", sekarangNanoFn)
	definePublicFunction(env, "now_nanos", sekarangNanoFn)

	// tidur / sleep(ms)
	tidurFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "waktu.tidur mengharapkan 1 argumen: milidetik", "Contoh: waktu.tidur(500)")
			}
			msVal, _, isInt, ok := CoerceToNumeric(args[0])
			if !ok || !isInt {
				return newError(nil, fmt.Sprintf("argumen milidetik waktu.tidur harus bilangan bulat, menerima %s", args[0].Type()), "")
			}
			if msVal < 0 {
				return newError(nil, "durasi waktu.tidur tidak boleh negatif", "")
			}
			time.Sleep(time.Duration(msVal) * time.Millisecond)
			return NULL
		},
	}
	definePublicFunction(env, "tidur", tidurFn)
	definePublicFunction(env, "sleep", tidurFn)

	// -------------------------------------------------------------
	// 2. Generator Durasi Praktis (detik, menit, jam, hari, minggu)
	// -------------------------------------------------------------

	makeDurationFn := func(mult int64) *Builtin {
		return &Builtin{
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return newError(nil, "mengharapkan 1 argumen angka", "")
				}
				_, fVal, _, ok := CoerceToNumeric(args[0])
				if !ok {
					return newError(nil, "argumen harus berupa angka", "")
				}
				return &Integer{Value: int64(fVal * float64(mult))}
			},
		}
	}

	detikFn := makeDurationFn(1000)
	definePublicFunction(env, "detik", detikFn)
	definePublicFunction(env, "seconds", detikFn)

	menitFn := makeDurationFn(60000)
	definePublicFunction(env, "menit", menitFn)
	definePublicFunction(env, "minutes", menitFn)

	jamFn := makeDurationFn(3600000)
	definePublicFunction(env, "jam", jamFn)
	definePublicFunction(env, "hours", jamFn)

	hariFn := makeDurationFn(86400000)
	definePublicFunction(env, "hari", hariFn)
	definePublicFunction(env, "days", hariFn)

	mingguFn := makeDurationFn(604800000)
	definePublicFunction(env, "minggu", mingguFn)
	definePublicFunction(env, "weeks", mingguFn)

	// -------------------------------------------------------------
	// 3. Format Penanggalan Cerdas & Format Indonesia
	// -------------------------------------------------------------

	// format(ms?, layout?)
	formatFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) > 2 {
				return newError(nil, "waktu.format mengharapkan maksimal 2 argumen: [ms], [pola]", "Contoh: waktu.format(waktu.sekarang(), \"YYYY-MM-DD HH:mm:ss\")")
			}
			t := time.Now()
			layout := "YYYY-MM-DD HH:mm:ss"

			if len(args) == 1 {
				if str, ok := args[0].(*String); ok {
					layout = str.Value
				} else if msVal, _, _, ok := CoerceToNumeric(args[0]); ok {
					t = time.Unix(0, msVal*int64(time.Millisecond))
				} else {
					return newError(nil, "argumen waktu.format tidak valid", "")
				}
			} else if len(args) == 2 {
				msVal, _, _, ok := CoerceToNumeric(args[0])
				if !ok {
					return newError(nil, "argumen pertama waktu.format harus stempel waktu angka", "")
				}
				t = time.Unix(0, msVal*int64(time.Millisecond))
				str, ok := args[1].(*String)
				if !ok {
					return newError(nil, "argumen kedua waktu.format harus berupa teks pola", "")
				}
				layout = str.Value
			}

			goLayout := translateDateTokensToGoLayout(layout)
			return &String{Value: t.Format(goLayout)}
		},
	}
	definePublicFunction(env, "format", formatFn)

	// format_indonesia(ms?, [gaya="lengkap"])
	formatIndoFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) > 2 {
				return newError(nil, "waktu.format_indonesia mengharapkan maksimal 2 argumen: [ms], [gaya]", "Contoh: waktu.format_indonesia(waktu.sekarang(), \"lengkap\")")
			}
			t := time.Now()
			gaya := "lengkap"

			if len(args) == 1 {
				if str, ok := args[0].(*String); ok {
					gaya = strings.ToLower(str.Value)
				} else if msVal, _, _, ok := CoerceToNumeric(args[0]); ok {
					t = time.Unix(0, msVal*int64(time.Millisecond))
				}
			} else if len(args) == 2 {
				msVal, _, _, ok := CoerceToNumeric(args[0])
				if !ok {
					return newError(nil, "argumen pertama harus stempel waktu angka", "")
				}
				t = time.Unix(0, msVal*int64(time.Millisecond))
				str, ok := args[1].(*String)
				if !ok {
					return newError(nil, "argumen kedua harus berupa teks gaya format", "")
				}
				gaya = strings.ToLower(str.Value)
			}

			zoneStr := getIndonesianZone(t)
			dayName := namaHariIndo[t.Weekday()]
			monthName := namaBulanIndo[t.Month()]

			switch gaya {
			case "tanggal":
				return &String{Value: fmt.Sprintf("%d %s %d", t.Day(), monthName, t.Year())}
			case "hari_tanggal":
				return &String{Value: fmt.Sprintf("%s, %d %s %d", dayName, t.Day(), monthName, t.Year())}
			case "waktu":
				return &String{Value: fmt.Sprintf("%02d:%02d %s", t.Hour(), t.Minute(), zoneStr)}
			case "waktu_detik":
				return &String{Value: fmt.Sprintf("%02d:%02d:%02d %s", t.Hour(), t.Minute(), t.Second(), zoneStr)}
			case "ringkas":
				shortMonth := monthName
				if len(shortMonth) > 3 {
					shortMonth = shortMonth[:3]
				}
				return &String{Value: fmt.Sprintf("%d %s %d", t.Day(), shortMonth, t.Year())}
			default: // "lengkap"
				return &String{Value: fmt.Sprintf("%s, %d %s %d, %02d:%02d %s", dayName, t.Day(), monthName, t.Year(), t.Hour(), t.Minute(), zoneStr)}
			}
		},
	}
	definePublicFunction(env, "format_indonesia", formatIndoFn)
	definePublicFunction(env, "format_id", formatIndoFn)

	// -------------------------------------------------------------
	// 4. Waktu Relatif (lalu / relatif / time_ago)
	// -------------------------------------------------------------

	laluFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 1 || len(args) > 2 {
				return newError(nil, "waktu.lalu mengharapkan 1 atau 2 argumen: target_ms, [base_ms]", "Contoh: waktu.lalu(timestamp)")
			}
			targetMs, _, _, ok1 := CoerceToNumeric(args[0])
			if !ok1 {
				return newError(nil, "argumen target waktu.lalu harus berupa angka", "")
			}
			baseMs := time.Now().UnixNano() / int64(time.Millisecond)
			if len(args) == 2 {
				bVal, _, _, ok2 := CoerceToNumeric(args[1])
				if !ok2 {
					return newError(nil, "argumen base waktu.lalu harus berupa angka", "")
				}
				baseMs = bVal
			}
			return &String{Value: formatRelativeTime(targetMs, baseMs)}
		},
	}
	definePublicFunction(env, "lalu", laluFn)
	definePublicFunction(env, "relatif", laluFn)
	definePublicFunction(env, "time_ago", laluFn)

	// -------------------------------------------------------------
	// 5. Pengukuran Eksekusi Kode (ukur / benchmark)
	// -------------------------------------------------------------

	ukurFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "waktu.ukur mengharapkan 1 argumen fungsi", "Contoh: waktu.ukur(fungsi() { ... })")
			}
			fn := args[0]
			start := time.Now()
			res := applyFunction(fn, []Object{}, nil)
			elapsed := time.Since(start)
			if isError(res) {
				return res
			}
			msElapsed := float64(elapsed.Nanoseconds()) / 1e6
			text := fmt.Sprintf("%.2f ms", msElapsed)
			if msElapsed < 1.0 {
				text = fmt.Sprintf("%.0f µs", float64(elapsed.Nanoseconds())/1e3)
			}
			return &Map{
				Pairs: map[string]Object{
					"hasil":       res,
					"durasi_ms":   &Float{Value: msElapsed},
					"durasi_teks": &String{Value: text},
				},
			}
		},
	}
	definePublicFunction(env, "ukur", ukurFn)
	definePublicFunction(env, "benchmark", ukurFn)

	// -------------------------------------------------------------
	// 6. Rincian & Komponen Waktu (info / rincian)
	// -------------------------------------------------------------

	infoFn := &Builtin{
		Fn: func(args ...Object) Object {
			t := time.Now()
			if len(args) > 0 {
				if msVal, _, _, ok := CoerceToNumeric(args[0]); ok {
					t = time.Unix(0, msVal*int64(time.Millisecond))
				}
			}
			weekday := int64(t.Weekday())
			isoDay := weekday
			if isoDay == 0 {
				isoDay = 7
			}
			zoneStr := getIndonesianZone(t)

			return &Map{
				Pairs: map[string]Object{
					"tahun":             &Integer{Value: int64(t.Year())},
					"bulan":             &Integer{Value: int64(t.Month())},
					"nama_bulan":        &String{Value: namaBulanIndo[t.Month()]},
					"hari":              &Integer{Value: int64(t.Day())},
					"nama_hari":         &String{Value: namaHariIndo[weekday]},
					"hari_dalam_minggu": &Integer{Value: weekday},
					"iso_hari":          &Integer{Value: isoDay},
					"jam":               &Integer{Value: int64(t.Hour())},
					"menit":             &Integer{Value: int64(t.Minute())},
					"detik":             &Integer{Value: int64(t.Second())},
					"milidetik":         &Integer{Value: int64(t.Nanosecond() / 1e6)},
					"zona_waktu":        &String{Value: zoneStr},
				},
			}
		},
	}
	definePublicFunction(env, "info", infoFn)
	definePublicFunction(env, "rincian", infoFn)

	// Getter cepat: tahun, bulan, hari
	tahunFn := &Builtin{
		Fn: func(args ...Object) Object {
			t := time.Now()
			if len(args) > 0 {
				if msVal, _, _, ok := CoerceToNumeric(args[0]); ok {
					t = time.Unix(0, msVal*int64(time.Millisecond))
				}
			}
			return &Integer{Value: int64(t.Year())}
		},
	}
	definePublicFunction(env, "tahun", tahunFn)
	definePublicFunction(env, "year", tahunFn)

	bulanFn := &Builtin{
		Fn: func(args ...Object) Object {
			t := time.Now()
			if len(args) > 0 {
				if msVal, _, _, ok := CoerceToNumeric(args[0]); ok {
					t = time.Unix(0, msVal*int64(time.Millisecond))
				}
			}
			return &Integer{Value: int64(t.Month())}
		},
	}
	definePublicFunction(env, "bulan", bulanFn)
	definePublicFunction(env, "month", bulanFn)

	hariKeFn := &Builtin{
		Fn: func(args ...Object) Object {
			t := time.Now()
			if len(args) > 0 {
				if msVal, _, _, ok := CoerceToNumeric(args[0]); ok {
					t = time.Unix(0, msVal*int64(time.Millisecond))
				}
			}
			return &Integer{Value: int64(t.Day())}
		},
	}
	definePublicFunction(env, "hari_ke", hariKeFn)
	definePublicFunction(env, "day_of_month", hariKeFn)

	// -------------------------------------------------------------
	// 7. Aritmatika Tanggal & Manipulasi
	// -------------------------------------------------------------

	// tambah / add(ms, durasi)
	tambahFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "waktu.tambah mengharapkan 2 argumen: timestamp dan durasi", "Contoh: waktu.tambah(ms, {\"hari\": 3})")
			}
			msVal, _, _, ok := CoerceToNumeric(args[0])
			if !ok {
				return newError(nil, "argumen pertama waktu.tambah harus berupa timestamp angka", "")
			}
			t := time.Unix(0, msVal*int64(time.Millisecond))

			if dMs, _, _, ok := CoerceToNumeric(args[1]); ok {
				return &Integer{Value: msVal + dMs}
			}

			if mp, ok := args[1].(*Map); ok {
				t = applyDurationMap(t, mp, 1)
				return &Integer{Value: t.UnixNano() / int64(time.Millisecond)}
			}
			return newError(nil, "argumen kedua harus berupa angka milidetik atau Kamus durasi", "")
		},
	}
	definePublicFunction(env, "tambah", tambahFn)
	definePublicFunction(env, "add", tambahFn)

	// kurang / sub(ms, durasi)
	kurangFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "waktu.kurang mengharapkan 2 argumen: timestamp dan durasi", "Contoh: waktu.kurang(ms, {\"hari\": 7})")
			}
			msVal, _, _, ok := CoerceToNumeric(args[0])
			if !ok {
				return newError(nil, "argumen pertama waktu.kurang harus berupa timestamp angka", "")
			}
			t := time.Unix(0, msVal*int64(time.Millisecond))

			if dMs, _, _, ok := CoerceToNumeric(args[1]); ok {
				return &Integer{Value: msVal - dMs}
			}

			if mp, ok := args[1].(*Map); ok {
				t = applyDurationMap(t, mp, -1)
				return &Integer{Value: t.UnixNano() / int64(time.Millisecond)}
			}
			return newError(nil, "argumen kedua harus berupa angka milidetik atau Kamus durasi", "")
		},
	}
	definePublicFunction(env, "kurang", kurangFn)
	definePublicFunction(env, "sub", kurangFn)

	// selisih / diff(msA, msB)
	selisihFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "waktu.selisih mengharapkan 2 argumen: waktuA dan waktuB", "Contoh: waktu.selisih(t1, t2)")
			}
			msA, _, _, ok1 := CoerceToNumeric(args[0])
			msB, _, _, ok2 := CoerceToNumeric(args[1])
			if !ok1 || !ok2 {
				return newError(nil, "argumen waktu.selisih harus berupa angka timestamp", "")
			}
			diff := msB - msA
			absDiff := diff
			if absDiff < 0 {
				absDiff = -absDiff
			}
			totalSec := absDiff / 1000
			days := totalSec / 86400
			remSec := totalSec % 86400
			hours := remSec / 3600
			remSec %= 3600
			minutes := remSec / 60
			seconds := remSec % 60
			remMs := absDiff % 1000

			return &Map{
				Pairs: map[string]Object{
					"total_ms":    &Integer{Value: diff},
					"total_detik": &Integer{Value: totalSec},
					"hari":        &Integer{Value: days},
					"jam":         &Integer{Value: hours},
					"menit":       &Integer{Value: minutes},
					"detik":       &Integer{Value: seconds},
					"milidetik":   &Integer{Value: remMs},
				},
			}
		},
	}
	definePublicFunction(env, "selisih", selisihFn)
	definePublicFunction(env, "diff", selisihFn)

	// awal_hari / start_of_day(ms?)
	awalHariFn := &Builtin{
		Fn: func(args ...Object) Object {
			t := time.Now()
			if len(args) > 0 {
				if msVal, _, _, ok := CoerceToNumeric(args[0]); ok {
					t = time.Unix(0, msVal*int64(time.Millisecond))
				}
			}
			start := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
			return &Integer{Value: start.UnixNano() / int64(time.Millisecond)}
		},
	}
	definePublicFunction(env, "awal_hari", awalHariFn)
	definePublicFunction(env, "start_of_day", awalHariFn)

	// akhir_hari / end_of_day(ms?)
	akhirHariFn := &Builtin{
		Fn: func(args ...Object) Object {
			t := time.Now()
			if len(args) > 0 {
				if msVal, _, _, ok := CoerceToNumeric(args[0]); ok {
					t = time.Unix(0, msVal*int64(time.Millisecond))
				}
			}
			end := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, int(999*time.Millisecond), t.Location())
			return &Integer{Value: end.UnixNano() / int64(time.Millisecond)}
		},
	}
	definePublicFunction(env, "akhir_hari", akhirHariFn)
	definePublicFunction(env, "end_of_day", akhirHariFn)

	// apakah_kabisat / is_leap_year(tahun)
	apakahKabisatFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "waktu.apakah_kabisat mengharapkan 1 argumen: tahun", "")
			}
			yVal, _, _, ok := CoerceToNumeric(args[0])
			if !ok {
				return newError(nil, "argumen tahun harus berupa angka", "")
			}
			isLeap := (yVal%4 == 0 && yVal%100 != 0) || (yVal%400 == 0)
			return nativeBoolToBoolean(isLeap)
		},
	}
	definePublicFunction(env, "apakah_kabisat", apakahKabisatFn)
	definePublicFunction(env, "is_leap_year", apakahKabisatFn)

	// -------------------------------------------------------------
	// 8. Pengurai Teks Waktu (urai / parse)
	// -------------------------------------------------------------

	uraiFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 1 || len(args) > 2 {
				return newError(nil, "waktu.urai mengharapkan 1 atau 2 argumen: teks, [pola]", "Contoh: waktu.urai(\"2026-10-15 17:00:00\")")
			}
			s, ok := args[0].(*String)
			if !ok {
				return newError(nil, "argumen pertama waktu.urai harus berupa teks", "")
			}
			raw := strings.TrimSpace(s.Value)
			if len(args) == 2 {
				p, ok2 := args[1].(*String)
				if !ok2 {
					return newError(nil, "argumen pola waktu.urai harus berupa teks", "")
				}
				layout := translateDateTokensToGoLayout(p.Value)
				t, err := time.ParseInLocation(layout, raw, time.Local)
				if err != nil {
					return newError(nil, fmt.Sprintf("Gagal mengurai waktu '%s' dengan pola '%s': %v", raw, p.Value, err), "")
				}
				return &Integer{Value: t.UnixNano() / int64(time.Millisecond)}
			}

			// Deteksi otomatis format umum
			commonLayouts := []string{
				time.RFC3339Nano,
				time.RFC3339,
				"2006-01-02 15:04:05",
				"2006-01-02 15:04",
				"2006-01-02",
				"02/01/2006 15:04:05",
				"02/01/2006 15:04",
				"02/01/2006",
				"02-01-2006 15:04:05",
				"02-01-2006",
			}
			for _, l := range commonLayouts {
				if t, err := time.ParseInLocation(l, raw, time.Local); err == nil {
					return &Integer{Value: t.UnixNano() / int64(time.Millisecond)}
				}
			}
			return newError(nil, fmt.Sprintf("Format waktu '%s' tidak dikenali", raw), "Gunakan pola format eksplisit sebagai argumen kedua, misal: waktu.urai(\"15/10/2026\", \"DD/MM/YYYY\")")
		},
	}
	definePublicFunction(env, "urai", uraiFn)
	definePublicFunction(env, "parse", uraiFn)

	// -------------------------------------------------------------
	// 9. Format Durasi (format_durasi / format_duration)
	// -------------------------------------------------------------

	formatDurasiFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "waktu.format_durasi mengharapkan 1 argumen: milidetik", "Contoh: waktu.format_durasi(7385000)")
			}
			msVal, _, _, ok := CoerceToNumeric(args[0])
			if !ok {
				return newError(nil, "argumen waktu.format_durasi harus berupa angka", "")
			}
			return &String{Value: formatHumanDuration(msVal)}
		},
	}
	definePublicFunction(env, "format_durasi", formatDurasiFn)
	definePublicFunction(env, "format_duration", formatDurasiFn)
}

// -------------------------------------------------------------
// Helper Functions (Token Translation, Formatting, Relative Time)
// -------------------------------------------------------------

func translateDateTokensToGoLayout(pattern string) string {
	if strings.Contains(pattern, "2006") || strings.Contains(pattern, "15:04") {
		return pattern
	}
	r := pattern
	r = strings.ReplaceAll(r, "YYYY", "2006")
	r = strings.ReplaceAll(r, "YY", "06")
	r = strings.ReplaceAll(r, "MM", "01")
	r = strings.ReplaceAll(r, "DD", "02")
	r = strings.ReplaceAll(r, "HH", "15")
	r = strings.ReplaceAll(r, "hh", "03")
	r = strings.ReplaceAll(r, "mm", "04")
	r = strings.ReplaceAll(r, "ss", "05")
	r = strings.ReplaceAll(r, "SSS", "000")
	r = strings.ReplaceAll(r, "Z", "-07:00")
	return r
}

func getIndonesianZone(t time.Time) string {
	_, offsetSec := t.Zone()
	offsetHours := offsetSec / 3600
	switch offsetHours {
	case 7:
		return "WIB"
	case 8:
		return "WITA"
	case 9:
		return "WIT"
	default:
		zName, _ := t.Zone()
		return zName
	}
}

func applyDurationMap(t time.Time, mp *Map, sign int) time.Time {
	for k, v := range mp.Pairs {
		vInt, _, _, ok := CoerceToNumeric(v)
		if !ok {
			continue
		}
		val := int(vInt) * sign
		switch strings.ToLower(k) {
		case "tahun", "years", "year":
			t = t.AddDate(val, 0, 0)
		case "bulan", "months", "month":
			t = t.AddDate(0, val, 0)
		case "hari", "days", "day":
			t = t.AddDate(0, 0, val)
		case "jam", "hours", "hour":
			t = t.Add(time.Duration(val) * time.Hour)
		case "menit", "minutes", "minute":
			t = t.Add(time.Duration(val) * time.Minute)
		case "detik", "seconds", "second":
			t = t.Add(time.Duration(val) * time.Second)
		case "milidetik", "ms":
			t = t.Add(time.Duration(val) * time.Millisecond)
		}
	}
	return t
}

func formatRelativeTime(targetMs, baseMs int64) string {
	diffMs := baseMs - targetMs
	if diffMs >= 0 {
		// Masa Lalu (Past)
		if diffMs < 45*1000 {
			return "baru saja"
		}
		if diffMs < 90*1000 {
			return "1 menit yang lalu"
		}
		if diffMs < 45*60*1000 {
			mins := int(math.Round(float64(diffMs) / 60000.0))
			return fmt.Sprintf("%d menit yang lalu", mins)
		}
		if diffMs < 90*60*1000 {
			return "1 jam yang lalu"
		}
		if diffMs < 22*3600*1000 {
			hours := int(math.Round(float64(diffMs) / 3600000.0))
			return fmt.Sprintf("%d jam yang lalu", hours)
		}
		// Kemarin
		tBase := time.Unix(0, baseMs*int64(time.Millisecond))
		tTarget := time.Unix(0, targetMs*int64(time.Millisecond))
		if tBase.Year() == tTarget.Year() && tBase.YearDay()-tTarget.YearDay() == 1 {
			return "kemarin"
		}
		if diffMs < 26*86400*1000 {
			days := int(math.Round(float64(diffMs) / 86400000.0))
			if days == 1 {
				return "kemarin"
			}
			return fmt.Sprintf("%d hari yang lalu", days)
		}
		if diffMs < 320*86400*1000 {
			months := int(math.Round(float64(diffMs) / (30 * 86400000.0)))
			if months <= 1 {
				return "1 bulan yang lalu"
			}
			return fmt.Sprintf("%d bulan yang lalu", months)
		}
		years := int(math.Round(float64(diffMs) / (365 * 86400000.0)))
		if years <= 1 {
			return "1 tahun yang lalu"
		}
		return fmt.Sprintf("%d tahun yang lalu", years)
	}

	// Masa Depan (Future)
	futureMs := -diffMs
	if futureMs < 45*1000 {
		return "sebentar lagi"
	}
	if futureMs < 90*1000 {
		return "1 menit lagi"
	}
	if futureMs < 45*60*1000 {
		mins := int(math.Round(float64(futureMs) / 60000.0))
		return fmt.Sprintf("%d menit lagi", mins)
	}
	if futureMs < 90*60*1000 {
		return "1 jam lagi"
	}
	if futureMs < 22*3600*1000 {
		hours := int(math.Round(float64(futureMs) / 3600000.0))
		return fmt.Sprintf("%d jam lagi", hours)
	}
	tBase := time.Unix(0, baseMs*int64(time.Millisecond))
	tTarget := time.Unix(0, targetMs*int64(time.Millisecond))
	if tTarget.Year() == tBase.Year() && tTarget.YearDay()-tBase.YearDay() == 1 {
		return "besok"
	}
	if futureMs < 26*86400*1000 {
		days := int(math.Round(float64(futureMs) / 86400000.0))
		if days == 1 {
			return "besok"
		}
		return fmt.Sprintf("%d hari lagi", days)
	}
	if futureMs < 320*86400*1000 {
		months := int(math.Round(float64(futureMs) / (30 * 86400000.0)))
		if months <= 1 {
			return "1 bulan lagi"
		}
		return fmt.Sprintf("%d bulan lagi", months)
	}
	years := int(math.Round(float64(futureMs) / (365 * 86400000.0)))
	if years <= 1 {
		return "1 tahun lagi"
	}
	return fmt.Sprintf("%d tahun lagi", years)
}

func formatHumanDuration(ms int64) string {
	if ms < 0 {
		return "-" + formatHumanDuration(-ms)
	}
	if ms < 1000 {
		return fmt.Sprintf("%d milidetik", ms)
	}
	sec := ms / 1000
	remMs := ms % 1000

	if sec < 60 {
		if remMs > 0 {
			return fmt.Sprintf("%.1f detik", float64(ms)/1000.0)
		}
		return fmt.Sprintf("%d detik", sec)
	}

	days := sec / 86400
	remSec := sec % 86400
	hours := remSec / 3600
	remSec %= 3600
	minutes := remSec / 60
	seconds := remSec % 60

	var parts []string
	if days > 0 {
		parts = append(parts, fmt.Sprintf("%d hari", days))
	}
	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%d jam", hours))
	}
	if minutes > 0 {
		parts = append(parts, fmt.Sprintf("%d menit", minutes))
	}
	if seconds > 0 {
		parts = append(parts, fmt.Sprintf("%d detik", seconds))
	}
	return strings.Join(parts, " ")
}
