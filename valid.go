package argparse

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type mValid map[string](func(string, interface{}) error)

var (
	ValidFuncError = errors.New("The validation function doesn't exist")
	methods        = make(mValid)
)

func (m mValid) call(tag reflect.StructTag, name string, value interface{}) error {
	if mehtod, ok := m[name]; ok {
		return mehtod(string(tag), value)
	}

	return ValidFuncError
}

func RegisterValidFunc(name string, f func(string, interface{}) error) bool {
	if _, ok := methods[name]; ok {
		return false
	}
	methods[name] = f
	return true
}

func (m mValid) validate(tag reflect.StructTag, value interface{}) error {
	validation := strings.TrimSpace(tag.Get(TAG_VALIDATE))

	for _, name := range strings.Split(validation, ",") {
		if name = strings.TrimSpace(name); name == "" {
			continue
		}
		if err := m.call(tag, name, value); err != nil {
			return errors.New(fmt.Sprintf("[%v] %v", name, err))
		}
	}

	return nil
}

func init() {
	RegisterValidFunc("validate_num_range", ValidateNumberRange)
}
