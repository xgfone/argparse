package argparse

import (
	"errors"
	"strings"
)

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
