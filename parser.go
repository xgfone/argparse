// Parse the CLI arguments to a struct.
//
// The inspiration is from the package of argparse in Python Std Library.
// Since Go doesn't have generic type, but, it may be a little different
// against argparse.
//
package argparse

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/xgfone/go-tools/parse"
)

const (
	TAG_NAME    = "name"
	TAG_DEFAULT = "default"
	TAG_HELP    = "help"
)

var (
	// Used to join the group name and the option name.
	Sep = "_"

	// Output the information of registering the options and parsing the argument
	Debug = false

	// Error
	NotPointerError = errors.New("Not a pointer to a struct")
	ExistError      = errors.New("This group has been registered")
)

func log(format string, a ...interface{}) (int, error) {
	f := fmt.Sprintf("%v\n", format)
	return fmt.Printf(f, a...)
}

type Parser struct {
	default_group string
	cache         map[string]interface{}
	group         map[string]interface{}
	flagSet       *flag.FlagSet
	parsed        bool
}

// New create a new parser.
func NewParser() *Parser {
	if Debug {
		log("The default group name is Default")
	}
	return &Parser{
		default_group: "Default",
		cache:         make(map[string]interface{}),
		group:         make(map[string]interface{}),
		flagSet:       flag.NewFlagSet(os.Args[0], flag.ExitOnError),
	}
}

// Set the name of the default group. Must set it before registering the options.
func (p *Parser) SetDefaultGroup(name string) {
	p.default_group = name
}

// Parse the arguments to the registered structs.
//
// If args is not nil, it's the arguments. Or use os.Args[1:].
// If it has been parsed, don't parse it again.
// For parsing it againt, you can create a new parser.
func (p *Parser) Parse(args []string) {
	if args == nil {
		args = os.Args[1:]
	}

	if p.parsed {
		return
	}
	p.flagSet.Parse(args)
	p.setSalues()
}

// Return true if parsed, or false.
func (p Parser) Parsed() bool {
	return p.parsed
}

// Register a pointer to struct.
//
// Return an error if group is not a pointer to struct or has been registered.
// Return nil if successfully.
//
// When registering a struct, only the exposed field. If the type of the field
// is not supported, skip it.
//
// When parsing the arguments, it will parse the result to the field of the struct.
//
func (p *Parser) Register(group interface{}) error {
	// group must be a pointer to a struct, and not nil.
	vg := reflect.ValueOf(group)
	if vg.Kind() != reflect.Ptr || vg.IsNil() || vg.Elem().Kind() != reflect.Struct {
		return NotPointerError
	}

	// Get the struct that the pointer points to
	vg = vg.Elem()

	// Check whether it is registered.
	name := vg.Type().Name()
	if _, ok := p.cache[name]; ok {
		return ExistError
	}

	// Register options.
	p.register_flag(name, vg)
	p.cache[name] = group
	return nil
}

func (p *Parser) setSalues() {
	for k, v := range p.cache {
		p.setGroup(k, reflect.ValueOf(v).Elem())
	}
}

func (p Parser) getName(gname, fname string) string {
	var name string
	if gname != p.default_group {
		name = gname + Sep + fname
	} else {
		name = fname
	}
	return strings.ToLower(name)
}

func (p *Parser) setGroup(gname string, group reflect.Value) {
	tg := group.Type()
	num := group.NumField()
	for i := 0; i < num; i++ {
		field := tg.Field(i)

		// If the strategies are not passed, skip it.
		if !validStrategy(field.Tag) {
			continue
		}

		fname := getFromTag(field.Tag, TAG_NAME, field.Name)
		name := p.getName(gname, fname)
		v, ok := p.group[name]
		if !ok {
			break
		}

		if Debug {
			log("Parsing [%v]:[%v] to %v.%v", name, reflect.ValueOf(v).Elem().Interface(),
				gname, fname)
		}

		vfield := group.Field(i)
		switch vfield.Kind() {
		case reflect.String:
			vfield.SetString(*v.(*string))
		case reflect.Bool:
			vfield.SetBool(*v.(*bool))
		case reflect.Float64, reflect.Float32:
			vfield.SetFloat(*v.(*float64))
		case reflect.Int64:
			vfield.SetInt(*v.(*int64))
		case reflect.Uint64:
			vfield.SetUint(*v.(*uint64))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
			vfield.SetInt(int64(*v.(*int)))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
			vfield.SetUint(uint64(*v.(*uint)))
		}
	}
}

func (p *Parser) register_flag(gname string, group reflect.Value) {
	num := group.NumField()
	for i := 0; i < num; i++ {
		// Calculate the name, the default value and help by the tag of the field.
		field := group.Type().Field(i)

		// If the strategies are not passed, skip it.
		if !validStrategy(field.Tag) {
			continue
		}

		_default := getFromTag(field.Tag, TAG_DEFAULT, "")
		usage := getFromTag(field.Tag, TAG_HELP, "")
		fname := getFromTag(field.Tag, TAG_NAME, field.Name)
		name := p.getName(gname, fname)

		if Debug {
			log("Registering the option: name[%v] default[%v] help[%v]", name, _default, usage)
		}

		switch group.Field(i).Kind() {
		case reflect.Bool:
			// For bool, the default is always false, and can't be true.
			// If true, the option is always true.
			// value := parse.ToBool(_default)
			p.group[name] = p.flagSet.Bool(name, false, usage)
		case reflect.String:
			p.group[name] = p.flagSet.String(name, _default, usage)
		case reflect.Float32:
			value := parse.ToF64(_default)
			p.group[name] = p.flagSet.Float64(name, value, usage)
		case reflect.Float64:
			value := parse.ToF64(_default)
			p.group[name] = p.flagSet.Float64(name, value, usage)
		case reflect.Int64:
			value := parse.ToI64(_default, 10)
			p.group[name] = p.flagSet.Int64(name, value, usage)
		case reflect.Uint64:
			value := parse.ToU64(_default, 10)
			p.group[name] = p.flagSet.Uint64(name, value, usage)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
			value := parse.ToInt(_default, 10)
			p.group[name] = p.flagSet.Int(name, value, usage)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
			value := parse.ToUint(_default, 10)
			p.group[name] = p.flagSet.Uint(name, value, usage)
		default:
			if Debug {
				log("Don't support the type, %v, so skip to register the option: %v.%v",
					group.Field(i).Type().String(), gname, fname)
			}
		}
	}
}
