package argparse

import (
	"errors"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
)

const (
	digits = "0123456789"
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

// Validate whether the length of the value is in the range, [min, max].
//
// min and max is from the tag, that's, 'reflect.StructTag', which are
// the key-value pairs in the tag of the corresponding field.
//
// The value must be string. If not, return an error.
//
// This validation has been registered as "validate_str_len". so you can use
// it through the tag of `validate:"validate_str_len"`. min and max are given
// by `min:"MIN_VALUE" max:"MAX_VALUE"`, which are converted to the integers
// based on the base 10. If failed to convert, return an error.
// min or max or both maybe been omitted. If either is been omitted,
// it is considered to pass the validation.
//
// Notice: the leading and tail whitespaces of the value will be trimed down,
// then calculate.
func ValidateStrLen(tag string, value interface{}) error {
	if v, ok := value.(string); !ok {
		return errors.New("The type of the value is not string")
	} else {
		_len := int64(len(v))

		if min := strings.TrimSpace(TagGet(tag, "min")); min != "" {
			if vmin, err := strconv.ParseInt(min, 10, 0); err != nil {
				return errors.New(fmt.Sprintf("[min] %v", err))
			} else if vmin > _len {
				return errors.New("The length of the value is less than min")
			}
		}

		if max := strings.TrimSpace(TagGet(tag, "max")); max != "" {
			if vmax, err := strconv.ParseInt(max, 10, 0); err != nil {
				return errors.New(fmt.Sprintf("[max] %v", err))
			} else if _len > vmax {
				return errors.New("The length of the value is greater than max")
			}
		}

		return nil
	}
}

// Validate whether the value matches the pattern that is from the tag of "pattern".
//
// The value must be a string. If not, return an error. If the pattern is empty
// or doesn't exist, return nil. If matching successfully, return nil. Or an error.
//
// This validation has been registered as "validate_str_regexp". so you can use
// it through the tag of `validate:"validate_str_regexp"`. The pattern is acquired
// by the tag `pattern:"PATTERN"`. The validation way is regexp.MatchString().
//
// Notice: the leading and tail whitespaces of the value will be trimed down,
// then calculate.
func ValidateStrRegexp(tag string, value interface{}) error {
	if s, ok := value.(string); !ok {
		return errors.New("The type of the value is not string")
	} else if pattern := strings.TrimSpace(TagGet(tag, "pattern")); pattern == "" {
		return nil
	} else {
		ok, err := regexp.MatchString(pattern, s)
		if ok {
			return nil
		} else if err != nil {
			return err
		} else {
			return errors.New("The value doesn't match the pattern")
		}
	}
}

// Validate whether the value only contains the digits, 0-9.
//
// If the value is not a string, it will be convert to a string by
// fmt.Sprintf("%v", value).
//
// This validation has been registered as "validate_digit". so you can use
// it through the tag of `validate:"validate_digit"`.
//
// Notice: the leading and tail whitespaces of the value will be trimed down,
// then calculate.
func ValidateDigit(tag string, value interface{}) error {
	v := strings.TrimSpace(fmt.Sprintf("%v", value))
	for _, r := range v {
		if !strings.ContainsRune(digits, r) {
			return errors.New("The validation fails")
		}
	}
	return nil
}

// Validate whether the value is the valid ip by net.ParseIP, which is registered
// as "validate_ip".
func ValidateIP(tag string, value interface{}) error {
	ip := net.ParseIP(value.(string))
	if ip == nil {
		return errors.New("The format of the ip is invalid")
	}
	return nil
}

// Validate whether the value is the valid ipv4 by net.ParseIP, which is registered
// as "validate_ip4".
func ValidateIP4(tag string, value interface{}) error {
	ip := net.ParseIP(value.(string))
	if ip == nil || ip.To4() == nil {
		return errors.New("The format of the ip is invalid")
	}

	return nil
}

// Validate whether the value is the valid ipv6 by net.ParseIP, which is registered
// as "validate_ip6".
func ValidateIP6(tag string, value interface{}) error {
	ip := net.ParseIP(value.(string))
	if ip == nil || ip.To16() == nil {
		return errors.New("The format of the ip is invalid")
	}
	return nil
}

// Validate whether the value is in the string array came from the tag of array.
//
// when using this validtor, you should give the tag, array, which is separated
// by the comma.
//
// It's registered as "validate_str_array".
func ValidateStrArray(tag string, value interface{}) error {
	v, ok := value.(string)
	if !ok {
		return errors.New("The type of the value must be string")
	}
	array := strings.Split(TagGet(tag, "array"), ",")
	for _, s := range array {
		if s == v {
			return nil
		}
	}

	return errors.New(fmt.Sprintf("The value[%v] is not in %v", v, array))
}

func init() {
	RegisterValidator("validate_str_not_empty", ValidateStrNotEmpty)
	RegisterValidator("validate_str_len", ValidateStrLen)
	RegisterValidator("validate_str_regexp", ValidateStrRegexp)
	RegisterValidator("validate_digit", ValidateDigit)
	RegisterValidator("validate_str_array", ValidateStrArray)

	RegisterValidator("validate_ip", ValidateIP)
	RegisterValidator("validate_ip4", ValidateIP4)
	RegisterValidator("validate_ip6", ValidateIP6)
}
