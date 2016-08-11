package argparse_test

import (
	"fmt"
	"strings"

	"github.com/xgfone/argparse"
)

func ExampleParser() {
	type Default struct {
		String string `name:"str" default:"0.0.0.0", help:"the ip to listen to" validate:"validate_str_not_empty"`
		Bool   bool   // The default is useless, and it will be ignored.

		Float32 float32 `default:"RRRR"` // The default value is ZERO
		Float64 float64 `default:"1.2"`

		Int   int   `default:"123"`
		Int8  int8  `default:"123"`
		Int16 int16 `default:"123"`
		Int32 int32 `default:"123"`
		Int64 int64 `default:"123"`

		Uint   uint   `default:"123"`
		Uint8  uint8  `default:"123"`
		Uint16 uint16 `default:"123" validate:"validate_num_range" min:"100" max:"200"`
		Uint32 uint32 `default:"123"`
		Uint64 uint64 `default:"123" strategy:"skip"`
	}

	type Group struct {
		String string `name:"str" default:"0.0.0.0" help:"the ip to listen to"`
		Bool   bool   // The default is useless, and it will be ignored.

		Float32 float32 `default:"RRRR"` // The default value is ZERO
		Float64 float64 `default:"1.2"`

		Int   int   `default:"123"`
		Int8  int8  `default:"123"`
		Int16 int16 `default:"123"`
		Int32 int32 `default:"123"`
		Int64 int64 `default:"123"`

		Uint   uint   `default:"123"`
		Uint8  uint8  `default:"123"`
		Uint16 uint16 `default:"123" validate:"validate_num_range" min:"100" max:"200"`
		Uint32 uint32 `default:"123"`
		Uint64 uint64 `default:"123" strategy:"skip"`
	}

	p := argparse.NewParser()
	default_ := Default{}
	group := Group{}

	p.Register(&default_)
	p.Register(&group)
	p.Parse(strings.Split("-str 127.0.0.1 -float32 2.5 -int32 456 -uint32 456 -bool -group_str 127.0.0.1 -group_float32 2.5 -group_int32 456 -group_uint32 456 Arg1 Arg2", " "))

	fmt.Printf("%T%+v\n", default_, default_)
	fmt.Printf("%T%+v\n", group, group)
	fmt.Printf("%v %v\n", p.NArg(), p.Args())

	// Output:
	// argparse_test.Default{String:127.0.0.1 Bool:true Float32:2.5 Float64:1.2 Int:123 Int8:123 Int16:123 Int32:456 Int64:123 Uint:123 Uint8:123 Uint16:123 Uint32:456 Uint64:0}
	// argparse_test.Group{String:127.0.0.1 Bool:false Float32:2.5 Float64:1.2 Int:123 Int8:123 Int16:123 Int32:456 Int64:123 Uint:123 Uint8:123 Uint16:123 Uint32:456 Uint64:0}
	// 2 [Arg1 Arg2]
}
