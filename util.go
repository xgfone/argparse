package argparse

import (
	"fmt"
	"reflect"
)

// Debugf is convenient for the validation plugins.
//
// If Debug is false, don't output the information.
func Debugf(format string, a ...interface{}) (int, error) {
	if !Debug {
		return 0, nil
	}
	f := fmt.Sprintf("[Debug] %v\n", format)
	return fmt.Printf(f, a...)
}

// Infof is convenient for the validation plugins.
func Infof(format string, a ...interface{}) (int, error) {
	f := fmt.Sprintf("[Info] %v\n", format)
	return fmt.Printf(f, a...)
}

// Errorf is convenient for the validation plugins.
func Errorf(format string, a ...interface{}) (int, error) {
	f := fmt.Sprintf("[Error] %v\n", format)
	return fmt.Printf(f, a...)
}

// Get the value of key from tag.
//
// This is the proxy of 'reflect.Tag'. You can convert tag to 'reflect.Tag' by yourself.
func TagGet(tag, key string) string {
	return reflect.StructTag(tag).Get(key)
}

func getFromTag(tag reflect.StructTag, key, default_ string) string {
	v := tag.Get(key)
	if v == "" {
		return default_
	}
	return v
}
