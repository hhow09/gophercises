package hash

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
)

const epsilon = 1e-6

type (
	PlaceholderType = struct{}
	AnyType         = any
)

// Repr returns the string representation of v.
func Repr(v any) string {
	if v == nil {
		return ""
	}

	// if func (v *Type) String() string, we can't use Elem()
	switch vt := v.(type) {
	case fmt.Stringer:
		return vt.String()
	}

	val := reflect.ValueOf(v)
	for val.Kind() == reflect.Ptr && !val.IsNil() {
		val = val.Elem()
	}

	return reprOfValue(val)
}

func reprOfValue(val reflect.Value) string {
	switch vt := val.Interface().(type) {
	case bool:
		return strconv.FormatBool(vt)
	case error:
		return vt.Error()
	case float32:
		return strconv.FormatFloat(float64(vt), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(vt, 'f', -1, 64)
	case fmt.Stringer:
		return vt.String()
	case int:
		return strconv.Itoa(vt)
	case int8:
		return strconv.Itoa(int(vt))
	case int16:
		return strconv.Itoa(int(vt))
	case int32:
		return strconv.Itoa(int(vt))
	case int64:
		return strconv.FormatInt(vt, 10)
	case string:
		return vt
	case uint:
		return strconv.FormatUint(uint64(vt), 10)
	case uint8:
		return strconv.FormatUint(uint64(vt), 10)
	case uint16:
		return strconv.FormatUint(uint64(vt), 10)
	case uint32:
		return strconv.FormatUint(uint64(vt), 10)
	case uint64:
		return strconv.FormatUint(vt, 10)
	case []byte:
		return string(vt)
	default:
		return fmt.Sprint(val.Interface())
	}
}

func CalcEntropy(m map[any]int) float64 {
	if len(m) == 0 || len(m) == 1 {
		return 1
	}

	var entropy float64
	var total int
	for _, v := range m {
		total += v
	}

	for _, v := range m {
		proba := float64(v) / float64(total)
		if proba < epsilon {
			proba = epsilon
		}
		entropy -= proba * math.Log2(proba)
	}

	return entropy / math.Log2(float64(len(m)))
}