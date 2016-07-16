# argparse
Parse the CLI arguments to a struct.

The inspiration is from the package of argparse in Python Std Library. Since Go doesn't have generic type, but, it may be a little different against argparse.

## Installation
```shell
$ go get -u github.com/xgfone/argparse
```

## Usage
```go
import "github.com/xgfone/argparse"

parser := argparse.NewParser()
parser.Debug = true
group := Struct{}
parser.Register(&group)

parser.Parse(nil)

fmt.Println(group.Field)
```

## Example
```go
package main

import (
	"fmt"
	"strings"

	"github.com/xgfone/argparse"
)

type Default struct {
	String string `name:"str" default:"0.0.0.0", help:"the ip to listen to"`
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
	Uint16 uint16 `default:"123"`
	Uint32 uint32 `default:"123"`
	Uint64 uint64 `default:"123"`
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
	Uint16 uint16 `default:"123"`
	Uint32 uint32 `default:"123"`
	Uint64 uint64 `default:"123"`
}

func main() {

	p := argparse.NewParser()
	default_ := Default{}
	group := Group{}

	p.Register(&default_)
	p.Register(&group)
	p.Parse(strings.Split("-str 127.0.0.1 -float32 2.5 -int32 456 -uint32 456 -bool -group_str 127.0.0.1 -group_float32 2.5 -group_int32 456 -group_uint32 456", " "))

	fmt.Printf("%T%+v\n", default_, default_)
	fmt.Printf("%T%+v\n", group, group)

	// Output:
	// main.Default{String:127.0.0.1 Bool:true Float32:2.5 Float64:1.2 Int:123 Int8:123 Int16:123 Int32:456 Int64:123 Uint:123 Uint8:123 Uint16:123 Uint32:456 Uint64:123}
	// main.Group{String:127.0.0.1 Bool:false Float32:2.5 Float64:1.2 Int:123 Int8:123 Int16:123 Int32:456 Int64:123 Uint:123 Uint8:123 Uint16:123 Uint32:456 Uint64:123}
}
```
