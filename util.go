package argparse

import (
	"fmt"
	"reflect"
	"strings"
)

// Strategy Sets
const (
	// If there is this strategy in a certain option, don't register it.
	STRATEGY_SKIP = "skip"
)

func Debugf(format string, a ...interface{}) (int, error) {
	if Debug {
		return 0, nil
	}
	f := fmt.Sprintf("[Debug] %v\n", format)
	return fmt.Printf(f, a...)
}

func Infof(format string, a ...interface{}) (int, error) {
	f := fmt.Sprintf("[Info] %v\n", format)
	return fmt.Printf(f, a...)
}

func Errorf(format string, a ...interface{}) (int, error) {
	f := fmt.Sprintf("[Error] %v\n", format)
	return fmt.Printf(f, a...)
}

// Get the value of key from tag.
//
// This is the proxy of reflect.Tag. You can convert tag to reflect.Tag by yourself.
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

func checkStrategy(node, sets string) bool {
	_sets := strings.Split(sets, ",")
	for _, s := range _sets {
		if strings.TrimSpace(s) == node {
			return true
		}
	}
	return false
}

func validStrategy(tag reflect.StructTag) bool {
	stag := strings.TrimSpace(tag.Get(TAG_STRATEGY))

	if checkStrategy(STRATEGY_SKIP, stag) {
		return false
	}

	return true
}
