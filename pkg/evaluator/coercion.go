package evaluator

import (
	"regexp"
	"strconv"
	"strings"
)

var numericRegex = regexp.MustCompile(`[-+]?[0-9]+(?:\.[0-9]+)?`)

// CoerceToNumeric attempts to extract a numeric value (int64 or float64) from an Object.
// If the object is already a number, it returns directly.
// If the object is a string, it uses Smart Extraction (RegEx/heuristic fallback)
// to extract numeric tokens (e.g., "100", "Rp 50000", "15 items").
func CoerceToNumeric(obj Object) (int64, float64, bool, bool) {
	switch o := obj.(type) {
	case *Integer:
		return o.Value, float64(o.Value), true, true
	case *Float:
		return int64(o.Value), o.Value, false, true
	case *Boolean:
		if o.Value {
			return 1, 1.0, true, true
		}
		return 0, 0.0, true, true
	case *String:
		raw := strings.TrimSpace(o.Value)
		// 1. Direct Integer parse
		if intVal, err := strconv.ParseInt(raw, 10, 64); err == nil {
			return intVal, float64(intVal), true, true
		}
		// 2. Direct Float parse
		if floatVal, err := strconv.ParseFloat(raw, 64); err == nil {
			return int64(floatVal), floatVal, false, true
		}
		// 3. Smart Fallback RegEx Extraction
		// Cleans currency/commas: e.g. "50.000" or "50,000" if no other decimal
		cleaned := strings.ReplaceAll(raw, ",", "")
		match := numericRegex.FindString(cleaned)
		if match != "" {
			if strings.Contains(match, ".") {
				if floatVal, err := strconv.ParseFloat(match, 64); err == nil {
					return int64(floatVal), floatVal, false, true
				}
			} else {
				if intVal, err := strconv.ParseInt(match, 10, 64); err == nil {
					return intVal, float64(intVal), true, true
				}
			}
		}
	}
	return 0, 0, false, false
}

// IsPurelyNumericString checks if a string consists only of valid numeric characters.
func IsPurelyNumericString(s string) bool {
	trimmed := strings.TrimSpace(s)
	if _, err := strconv.ParseInt(trimmed, 10, 64); err == nil {
		return true
	}
	if _, err := strconv.ParseFloat(trimmed, 64); err == nil {
		return true
	}
	return false
}

// CoerceToBoolean evaluates truthiness in human-centric terms.
// Falsy values: null, false, 0, 0.0, ""
func CoerceToBoolean(obj Object) bool {
	switch o := obj.(type) {
	case *Boolean:
		return o.Value
	case *Null:
		return false
	case *Integer:
		return o.Value != 0
	case *Float:
		return o.Value != 0.0
	case *String:
		return len(o.Value) > 0
	case *Array:
		return len(o.Elements) > 0
	case *Error:
		return false
	case *ReturnValue:
		return CoerceToBoolean(o.Value)
	default:
		return true
	}
}
