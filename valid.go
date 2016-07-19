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

// Register a funcation to validate the option corresponding to the field.
//
// The first argument of the function is the tag of the field with string.
// And the second is the value of the parsed option. Refer to ValidateNumberRange.
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
	RegisterValidFunc("validate_str_not_empty", ValidateStrNotEmpty)
}

// Validate whether the value is a non-empty string.
//
// Return nil if it's not empty. Or false. If the value is the type of string,
// don't validate it, that's, return nil.
//
// This function has been registered as "validate_str_not_empty", and you can
// use it with the tag of `validate:"validate_str_not_empty"`.
func ValidateStrNotEmpty(tag string, value interface{}) error {
	if v, ok := value.(string); !ok {
		return nil
	} else if strings.TrimSpace(v) == "" {
		return errors.New("The string is empty")
	} else {
		return nil
	}
}
