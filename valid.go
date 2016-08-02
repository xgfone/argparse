package argparse

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var (
	validatorError = errors.New("The validation function doesn't exist")
	validators     = make(tValidation)
)

type Validation interface {
	// Validate the second argument, which is the value of the current field,
	// by the information in the tag of that field.
	//
	// Return nil if validating successfully, or an error, which is the reason
	// for failure.
	Validate(string, interface{}) error
}

// Register a validator to validate whether the value of the field is valid.
//
// The first argument is the name of the validator, which must be unique. The
// second is the validator, which is either a type implementing the interface
// Validation or a funciton whose type is the same as the method of Validate of
// the interface Validation.
//
// Return true if registering successfully. Return false if having been registered.
// If the validator is invalid, it will panic.
func RegisterValidator(name string, validator interface{}) bool {
	if _, ok := validators[name]; ok {
		return false
	}

	if _, ok := validator.(Validation); !ok {
		if _, ok := validator.(func(string, interface{}) error); !ok {
			panic("The validator is invalid")
		}
	}

	validators[name] = validator
	return true
}

type tValidation map[string]interface{}

func (t tValidation) call(tag reflect.StructTag, name string, value interface{}) error {
	defer func() {
		if err := recover(); err != nil {
			return errors.New("This validator panics: %v", err)
		}
	}()

	if validator, ok := t[name]; ok {
		if validation, ok := validator.(Validation); ok {
			return validation.Validate(string(tag), value)
		} else if validation, ok := validator.(func(string, interface{}) error); ok {
			return validation(string(tag), value)
		}
	}

	return validatorError
}

func (t tValidation) Validate(tag reflect.StructTag, value interface{}) error {
	validation := strings.TrimSpace(tag.Get(TAG_VALIDATE))

	for _, name := range strings.Split(validation, ",") {
		if name = strings.TrimSpace(name); name == "" {
			continue
		}
		if err := t.call(tag, name, value); err != nil {
			return errors.New(fmt.Sprintf("[%v] %v", name, err))
		}
	}

	return nil
}

// Register a funcation to validate the option corresponding to the field. (Deprecated)
//
// The first argument of the function is the tag of the field with string.
// And the second is the value of the parsed option. Refer to ValidateNumberRange.
func RegisterValidFunc(name string, f func(string, interface{}) error) bool {
	if _, ok := validators[name]; ok {
		return false
	}
	validators[name] = f
	return true
}
