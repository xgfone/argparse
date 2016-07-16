package argparse

import "reflect"

func getFromTag(tag reflect.StructTag, key, default_ string) string {
	v := tag.Get(key)
	if v == "" {
		return default_
	}
	return v
}
