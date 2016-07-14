package goargparse

import "reflect"

func getTagVar(tag reflect.StructTag, key, default_ string) string {
	v := tag.Get(key)
	if v == "" {
		return default_
	}
	return v
}
