package LuaVM

import (
	"math"
)

type Table struct {
	Array     []*Value
	Hash      map[Value]*Value
	ArraySize uint64
	MaxN      uint64
}

//TODO: add metamethod support

func (t *Table) Set(key Value, val *Value) {
	if key.Type == NUMBER {
		if math.Floor(key.Val.(float64)) == key.Val.(float64) && uint64(key.Val.(float64)) < t.ArraySize {
			t.Array[uint64(key.Val.(float64))] = val
			t.CalcMaxN()
			return
		}
	}
	t.Hash[key] = val
}

func (t *Table) Get(key Value) *Value {
	if key.Type == NUMBER {
		if math.Floor(key.Val.(float64)) == key.Val.(float64) && uint64(key.Val.(float64)) < t.ArraySize {
			return t.Array[uint64(key.Val.(float64))]
		}
	}
	return t.Hash[key]
}

func (t *Table) CalcMaxN() {
	for k, v := range t.Array {
		if v.Type == NIL {
			t.MaxN = uint64(k - 1)
		}
	}
}

func (t *Table) Len() *Value {
	return &Value{Type: NUMBER, Val: float64(t.MaxN)}
}
