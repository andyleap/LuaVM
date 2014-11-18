package LuaVM

import "math"

type Table struct {
	Array     []*Value
	Hash      map[Value]*Value
	ArraySize uint64
	MaxN      uint64
}

//TODO: add metamethod support

func (t *Table) Set(key Value, val *Value) {
	if key.Type == NUMBER {
		if math.Floor(float64(key.Val.(Number))) == float64(key.Val.(Number)) && uint64(key.Val.(Number)) < t.ArraySize {
			t.Array[uint64(key.Val.(Number))] = val
			t.CalcMaxN()
			return
		}
	}
	t.Hash[key] = val
}

func (t *Table) Get(key Value) *Value {
	if key.Type == NUMBER {
		if math.Floor(float64(key.Val.(Number))) == float64(key.Val.(Number)) && uint64(key.Val.(Number)) < t.ArraySize {
			return t.Array[uint64(key.Val.(Number))]
		}
	}

	v, ok := t.Hash[key]
	if !ok {
		v = &Value{Type: NIL}
	}
	return v
}

func (t *Table) CalcMaxN() {
	for k, v := range t.Array {
		if v == nil || v.Type == NIL {
			t.MaxN = uint64(k - 1)
			return
		}
	}
	t.MaxN = uint64(len(t.Array))
}

func (t *Table) Len() *Value {
	return &Value{Type: NUMBER, Val: float64(t.MaxN)}
}
