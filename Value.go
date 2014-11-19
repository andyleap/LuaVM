// LuaVM project LuaVM.go
package LuaVM

import (
	"strconv"
)

type ValueType uint8

const (
	NIL ValueType = iota
	BOOLEAN
	_
	NUMBER
	STRING
	TABLE
	FUNCTION
	CLOSURE
	GOFUNCTION
)

type GOFUNC func(params []*Value, v *VM) []*Value

type Number float64

type Value struct {
	Type ValueType
	Val  interface{}
}

func (v *Value) String() string {
	switch v.Type {
	case NUMBER:
		return strconv.FormatFloat(float64(v.Val.(Number)), 'G', -1, 64)
	case STRING:
		return v.Val.(string)
	case BOOLEAN:
		return strconv.FormatBool(v.Val.(Integer) == 0)
	case NIL:
		return "NIL"
	case GOFUNCTION:
		return "GOFUNCTION"
	case FUNCTION:
		return "FUNCTION"
	case CLOSURE:
		return "CLOSURE"
	case TABLE:
		return "TABLE"
	}
	return ""
}

func NewNil() *Value {
	return &Value{Type: NIL}
}

func NewString(str string) *Value {
	return &Value{Type: STRING, Val: str}
}

func NewNumber(n float64) *Value {
	return &Value{Type: NUMBER, Val: Number(n)}
}
