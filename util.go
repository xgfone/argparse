package argparse

import (
	"reflect"
	"strings"
)

const (
	TAG_STRATEGY  = "strategy"
	STRATEGY_SKIP = "skip"
)

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
