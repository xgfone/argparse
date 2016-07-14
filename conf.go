package goargparse

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

const (
	TAG_NAME    = "name"
	TAG_DEFAULT = "default"
	TAG_HELP    = "help"
)

var (
	// Error
	NotPointerError = errors.New("Not a pointer to a struct")
	InvalidError    = errors.New("The value is invalid")
	ExistError      = errors.New("This group has been registered")
	CanNotSetError  = errors.New("Can not been set")
	TypeError       = errors.New("Don't support this field type")
)

type Parser struct {
	default_group string
	cache         map[string]interface{}
	flagSet       *flag.FlagSet
}

func NewParser() *Parser {
	return &Parser{
		default_group: "Default",
		cache:         make(map[string]interface{}),
		flagSet:       flag.NewFlagSet(os.Args[0], flag.ExitOnError),
	}
}

// Set the default group name. Must be before registering the options.
func (p *Parser) SetDefaultGroup(name string) *Parser {
	p.default_group = name
	return p
}

func (p *Parser) Parse(args []string) error {
	if args == nil {
		args = os.Args[1:]
	}
	return p.flagSet.Parse(args)
}

func (p *Parser) Register(group interface{}) error {
	vg := reflect.ValueOf(group)
	if vg.Kind() != reflect.Ptr || vg.IsNil() {
		return NotPointerError
	}

	if !vg.IsValid() {
		return InvalidError
	}

	vg = vg.Elem()

	name := vg.Type().Name()
	if _, ok := p.cache[name]; ok {
		return ExistError
	}

	p.register_flag(name, vg)
	p.cache[name] = group
	return nil
}

func (p *Parser) register_flag(group string, vg reflect.Value) {
	tg := vg.Type()

	num := vg.NumField()
	for i := 0; i < num; i++ {
		field := tg.Field(i)
		name := getTagVar(field.Tag, TAG_NAME, strings.ToLower(field.Name))
		if group != p.default_group {
			name = group + name
		}
		_default := getTagVar(field.Tag, TAG_DEFAULT, "")
		usage := getTagVar(field.Tag, TAG_HELP, "")

		fv := vg.Field(i)
		kv := fv.Kind()
		fmt.Println(field.Name, field.Type)
		switch kv {
		case reflect.Bool:
			v := fv.Interface().(bool)
			value := false
			if _v, err := strconv.ParseBool(_default); err != nil {
				value = _v
			}
			p.flagSet.BoolVar(&v, name, value, usage)
		case reflect.String:
			v := fv.Interface().(string)
			fmt.Println("====", &v)
			p.flagSet.StringVar(&v, name, _default, usage)
		case reflect.Float64:
			v := fv.Interface().(float64)
			value := 0.0
			if d, err := strconv.ParseFloat(_default, 64); err != nil {
				value = d
			}
			p.flagSet.Float64Var(&v, name, value, usage)
		// case reflect.Int:
		// 	v := fv.Interface().(int)
		// 	value := 0
		// 	if d, err := strconv.ParseInt(_default, 32); err != nil {
		// 		value = int(d)
		// 	}
		// 	p.flagSet.Float64Var(&v, name, value, usage)
		case reflect.Int64:
			v := fv.Interface().(int64)
			var value int64 = 0
			if d, err := strconv.ParseInt(_default, 10, 64); err != nil {
				value = d
			}
			p.flagSet.Int64Var(&v, name, value, usage)
		// case reflect.Uint:
		case reflect.Uint64:
			v := fv.Interface().(uint64)
			var value uint64 = 0
			if d, err := strconv.ParseUint(_default, 10, 64); err != nil {
				value = d
			}
			p.flagSet.Uint64Var(&v, name, value, usage)
		default:
			fmt.Println(field.Name, field.Type)
		}
	}

}
