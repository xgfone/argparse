package argparse

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/xgfone/go-tools/compare"
	"github.com/xgfone/go-tools/parse"
)

func validateNumberRange(tag string, value interface{}) error {
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
