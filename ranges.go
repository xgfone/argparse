package argparse

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/xgfone/go-tools/compare"
	"github.com/xgfone/go-tools/parse"
)

// Validate whether the value is in the range, [min, max].
//
// min and max is from the tag, that's, 'reflect.StructTag', which are
// the key-value pairs in the tag of the corresponding field.
//
// The type of the value is one of int, int8, int16, int32, int64, uint, uint8,
// uint16, uint32, uint64, float32, float64. And min and max are converted to
// the corresponding type according to the value.
//
// This validation has been registered as "validate_num_range". so you can use
// it through the tag of `validate:"validate_num_range"`. min and max are given by
// `min:"MIN_VALUE" max:"MAX_VALUE"`. min or max or both maybe been omitted.
// If either is been omitted, it is considered to pass the validation.
func ValidateNumberRange(tag string, value interface{}) error {
	min := TagGet(tag, "min")
	max := TagGet(tag, "max")

	if min == "" && max == "" {
		return nil
	}

	var vmin, vmax interface{}
	var typ reflect.Type

	switch reflect.ValueOf(value).Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		typ = reflect.TypeOf(int64(0))
		if v := parse.ToI64(min, 10); v == 0 {
			vmin = nil
		} else {
			vmin = v
		}
		if v := parse.ToI64(max, 10); v == 0 {
			vmax = nil
		} else {
			vmax = v
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		typ = reflect.TypeOf(uint64(0))
		if v := parse.ToU64(min, 10); v == 0 {
			vmin = nil
		} else {
			vmin = v
		}
		if v := parse.ToU64(max, 10); v == 0 {
			vmax = nil
		} else {
			vmax = v
		}
	case reflect.Float32, reflect.Float64:
		typ = reflect.TypeOf(float64(0.0))
		if v := parse.ToF64(min); v == 0.0 {
			vmin = nil
		} else {
			vmin = v
		}
		if v := parse.ToF64(max); v == 0.0 {
			vmax = nil
		} else {
			vmax = v
		}
	}
	v := reflect.ValueOf(value).Convert(typ).Interface()

	if vmin != nil && !compare.LE(vmin, v) {
		return errors.New(fmt.Sprintf("the value %v is less than %v", v, vmin))
	}

	if vmax != nil && !compare.LE(v, vmax) {
		return errors.New(fmt.Sprintf("the value %v is more than %v", v, vmax))
	}

	return nil
}

func init() {
	RegisterValidator("validate_num_range", ValidateNumberRange)
}
