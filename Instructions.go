package LuaVM

import (
	"math"
)

type Stackframe struct {
	Regs     []*Value
	Upvalues []*Value
}

type Closure struct {
	Stack       *Stackframe
	Function    *FunctionPrototype
	PC          int64
	ReturnCount uint64
	ReturnPos   uint64
	ReturnFunc  func(*Closure, *VM)
}

type VM struct {
	G            map[string]*Value
	ClosureStack []*Closure
	C            *Closure
}

func Move(i *Instr, c *Closure, v *VM) {
	c.Stack.Regs[i.A] = &Value{
		Type: c.Stack.Regs[i.B].Type,
		Val:  c.Stack.Regs[i.B].Val,
	}
}

func LoadNil(i *Instr, c *Closure, v *VM) {
	for l1 := int32(i.A); l1 <= i.B; l1++ {
		c.Stack.Regs[l1] = &Value{
			Type: NIL,
		}
	}
}

func LoadK(i *Instr, c *Closure, v *VM) {
	c.Stack.Regs[i.A] = &Value{
		Type: c.Function.Constants[i.B].Type,
		Val:  c.Function.Constants[i.B].Val,
	}
}

func LoadBool(i *Instr, c *Closure, v *VM) {
	c.Stack.Regs[i.A] = &Value{
		Type: BOOLEAN,
		Val:  i.B,
	}
	if i.C != 0 {
		c.PC++
	}
}

func GetGlobal(i *Instr, c *Closure, v *VM) {
	if c.Function.Constants[i.B].Type != STRING {
		panic("Constant type is not string")
	}
	c.Stack.Regs[i.A] = &Value{
		Type: v.G[c.Function.Constants[i.B].Val.(string)].Type,
		Val:  v.G[c.Function.Constants[i.B].Val.(string)].Val,
	}
}

func SetGlobal(i *Instr, c *Closure, v *VM) {
	if c.Function.Constants[i.B].Type != STRING {
		panic("Constant type is not string")
	}
	v.G[c.Function.Constants[i.B].Val.(string)] = &Value{
		Type: c.Stack.Regs[i.A].Type,
		Val:  c.Stack.Regs[i.A].Val,
	}
}

func GetUpVal(i *Instr, c *Closure, v *VM) {
	c.Stack.Regs[i.A] = &Value{
		Type: c.Stack.Upvalues[i.B].Type,
		Val:  c.Stack.Upvalues[i.B].Val,
	}
}

func SetUpVal(i *Instr, c *Closure, v *VM) {
	c.Stack.Upvalues[i.B] = &Value{
		Type: c.Stack.Regs[i.A].Type,
		Val:  c.Stack.Regs[i.A].Val,
	}
}

func GetTable(i *Instr, c *Closure, v *VM) {
	var key *Value
	if i.C&256 == 256 {
		key = &c.Function.Constants[i.C&255]
	} else {
		key = c.Stack.Regs[i.C]
	}
	if c.Stack.Regs[i.B].Type != TABLE {
		panic("Value is not table")
	}
	val := c.Stack.Regs[i.B].Val.(*Table).Get(*key)
	c.Stack.Regs[i.A] = &Value{
		Type: val.Type,
		Val:  val.Val,
	}
}

func SetTable(i *Instr, c *Closure, v *VM) {
	var key *Value
	if i.B&256 == 256 {
		key = &c.Function.Constants[i.B&255]
	} else {
		key = c.Stack.Regs[i.B]
	}
	var val *Value
	if i.C&256 == 256 {
		val = &c.Function.Constants[i.C&255]
	} else {
		val = c.Stack.Regs[i.C]
	}
	if c.Stack.Regs[i.A].Type != TABLE {
		panic("Value is not table")
	}
	c.Stack.Regs[i.A].Val.(*Table).Set(*key, &Value{
		Type: val.Type,
		Val:  val.Val,
	})
}

func Add(i *Instr, c *Closure, v *VM) {
	var bval *Value
	if i.B&256 == 256 {
		bval = &c.Function.Constants[i.B&255]
	} else {
		bval = c.Stack.Regs[i.B]
	}
	var cval *Value
	if i.C&256 == 256 {
		cval = &c.Function.Constants[i.C&255]
	} else {
		cval = c.Stack.Regs[i.C]
	}
	if bval.Type != NUMBER || cval.Type != NUMBER {
		panic("Trying to add non-numbers")
	}
	c.Stack.Regs[i.A] = &Value{
		Type: NUMBER,
		Val:  bval.Val.(Number) + cval.Val.(Number),
	}
}

func Sub(i *Instr, c *Closure, v *VM) {
	var bval *Value
	if i.B&256 == 256 {
		bval = &c.Function.Constants[i.B&255]
	} else {
		bval = c.Stack.Regs[i.B]
	}
	var cval *Value
	if i.C&256 == 256 {
		cval = &c.Function.Constants[i.C&255]
	} else {
		cval = c.Stack.Regs[i.C]
	}
	if bval.Type != NUMBER || cval.Type != NUMBER {
		panic("Trying to sub non-numbers")
	}
	c.Stack.Regs[i.A] = &Value{
		Type: NUMBER,
		Val:  bval.Val.(Number) - cval.Val.(Number),
	}
}

func Mul(i *Instr, c *Closure, v *VM) {
	var bval *Value
	if i.B&256 == 256 {
		bval = &c.Function.Constants[i.B&255]
	} else {
		bval = c.Stack.Regs[i.B]
	}
	var cval *Value
	if i.C&256 == 256 {
		cval = &c.Function.Constants[i.C&255]
	} else {
		cval = c.Stack.Regs[i.C]
	}
	if bval.Type != NUMBER || cval.Type != NUMBER {
		panic("Trying to mul non-numbers")
	}
	c.Stack.Regs[i.A] = &Value{
		Type: NUMBER,
		Val:  bval.Val.(Number) * cval.Val.(Number),
	}
}

func Div(i *Instr, c *Closure, v *VM) {
	var bval *Value
	if i.B&256 == 256 {
		bval = &c.Function.Constants[i.B&255]
	} else {
		bval = c.Stack.Regs[i.B]
	}
	var cval *Value
	if i.C&256 == 256 {
		cval = &c.Function.Constants[i.C&255]
	} else {
		cval = c.Stack.Regs[i.C]
	}
	if bval.Type != NUMBER || cval.Type != NUMBER {
		panic("Trying to div non-numbers")
	}
	c.Stack.Regs[i.A] = &Value{
		Type: NUMBER,
		Val:  bval.Val.(Number) / cval.Val.(Number),
	}
}

func Mod(i *Instr, c *Closure, v *VM) {
	var bval *Value
	if i.B&256 == 256 {
		bval = &c.Function.Constants[i.B&255]
	} else {
		bval = c.Stack.Regs[i.B]
	}
	var cval *Value
	if i.C&256 == 256 {
		cval = &c.Function.Constants[i.C&255]
	} else {
		cval = c.Stack.Regs[i.C]
	}
	if bval.Type != NUMBER || cval.Type != NUMBER {
		panic("Trying to add non-numbers")
	}
	c.Stack.Regs[i.A] = &Value{
		Type: NUMBER,
		Val:  math.Mod(bval.Val.(float64), cval.Val.(float64)),
	}
}

func Pow(i *Instr, c *Closure, v *VM) {
	var bval *Value
	if i.B&256 == 256 {
		bval = &c.Function.Constants[i.B&255]
	} else {
		bval = c.Stack.Regs[i.B]
	}
	var cval *Value
	if i.C&256 == 256 {
		cval = &c.Function.Constants[i.C&255]
	} else {
		cval = c.Stack.Regs[i.C]
	}
	if bval.Type != NUMBER || cval.Type != NUMBER {
		panic("Trying to add non-numbers")
	}
	c.Stack.Regs[i.A] = &Value{
		Type: NUMBER,
		Val:  math.Pow(bval.Val.(float64), cval.Val.(float64)),
	}
}

func Unm(i *Instr, c *Closure, v *VM) {
	if c.Stack.Regs[i.B].Type != NUMBER {
		panic("Trying to unm non-number")
	}
	c.Stack.Regs[i.A] = &Value{
		Type: NUMBER,
		Val:  -(c.Stack.Regs[i.B].Val.(Number)),
	}
}

func Not(i *Instr, c *Closure, v *VM) {
	bval := c.Stack.Regs[i.B]
	if bval.Type == NIL {
		c.Stack.Regs[i.A] = &Value{
			Type: BOOLEAN,
			Val:  1,
		}
	}
	if bval.Type == BOOLEAN {
		val := 1
		if bval.Val != 0 {
			val = 0
		}
		c.Stack.Regs[i.A] = &Value{
			Type: BOOLEAN,
			Val:  val,
		}
	}

}

func Len(i *Instr, c *Closure, v *VM) {
	bval := c.Stack.Regs[i.B]
	var val *Value
	if bval.Type == TABLE {
		val = bval.Val.(*Table).Len()
	}
	if bval.Type == STRING {
		val = &Value{Type: NUMBER, Val: Number(len(bval.Val.(string)))}
	}
	c.Stack.Regs[i.A] = val
}

func Concat(i *Instr, c *Closure, v *VM) {
	str := ""
	for l1 := i.B; l1 <= int32(i.C); l1++ {
		if c.Stack.Regs[l1].Type != STRING {
			panic("Attempting to concat non-strings")
		}
		str = str + c.Stack.Regs[l1].Val.(string)
	}
	c.Stack.Regs[i.A] = &Value{
		Type: STRING,
		Val:  str,
	}
}

func Jmp(i *Instr, c *Closure, v *VM) {
	c.PC = c.PC + int64(i.B)
}

func Call(i *Instr, c *Closure, v *VM) {
	function := c.Stack.Regs[i.A]
	if function.Type == CLOSURE {
		v.ClosureStack = append(v.ClosureStack, v.C)
		v.C = function.Val.(*Closure)
		v.C.ReturnCount = uint64(i.C)
		v.C.ReturnPos = uint64(i.A)
		v.C.Stack.Regs = make([]*Value, v.C.Function.MaxStackSize)
		if i.B == 0 {
			for l1 := int32(0); l1+int32(i.A)+1 < int32(len(c.Stack.Regs)); l1++ {
				if l1 >= int32(v.C.Function.MaxStackSize) {
					break
				}
				v.C.Stack.Regs[l1] = c.Stack.Regs[l1+int32(i.A)+1]
			}
		} else if i.B > 1 {
			for l1 := int32(0); l1 < i.B-1; l1++ {
				if l1 >= int32(v.C.Function.MaxStackSize) {
					break
				}
				v.C.Stack.Regs[l1] = c.Stack.Regs[l1+int32(i.A)+1]
			}
		}
		return
	}
	if function.Type == GOFUNCTION {
		function.Val.(GOFUNC)(c, v)
	}
}

func Return(i *Instr, c *Closure, v *VM) {
	v.C = v.ClosureStack[len(v.ClosureStack)-1]
	v.ClosureStack = v.ClosureStack[:len(v.ClosureStack)-1]

	if i.B == 0 {
		for l1 := int32(0); l1+int32(i.A) < int32(len(c.Stack.Regs)); l1++ {
			if c.ReturnCount > 0 && l1 >= int32(c.ReturnCount) {
				break
			}
			if len(v.C.Stack.Regs) <= int(l1+int32(c.ReturnPos)) {
				v.C.Stack.Regs = append(v.C.Stack.Regs, c.Stack.Regs[l1+int32(i.A)])
			} else {
				v.C.Stack.Regs[l1+int32(c.ReturnPos)] = c.Stack.Regs[l1+int32(i.A)]
			}
		}
	} else if i.B > 1 {
		for l1 := int32(0); l1 < i.B-1; l1++ {
			if c.ReturnCount > 0 && l1 >= int32(c.ReturnCount) {
				break
			}
			if len(v.C.Stack.Regs) <= int(l1+int32(c.ReturnPos)) {
				v.C.Stack.Regs = append(v.C.Stack.Regs, c.Stack.Regs[l1+int32(i.A)])
			} else {
				v.C.Stack.Regs[l1+int32(c.ReturnPos)] = c.Stack.Regs[l1+int32(i.A)]
			}
		}
	}
	if c.ReturnFunc != nil {
		c.ReturnFunc(v.C, v)
	}
}

func TailCall(i *Instr, c *Closure, v *VM) {
	function := c.Stack.Regs[i.A]
	if function.Type == CLOSURE {
		v.C = function.Val.(*Closure)
		v.C.ReturnCount = c.ReturnCount
		v.C.ReturnPos = c.ReturnPos
		v.C.Stack.Regs = make([]*Value, v.C.Function.MaxStackSize)
		if i.B == 0 {
			for l1 := int32(0); l1+int32(i.A)+1 < int32(len(c.Stack.Regs)); l1++ {
				if l1 >= int32(v.C.Function.MaxStackSize) {
					break
				}
				v.C.Stack.Regs[l1] = c.Stack.Regs[l1+int32(i.A)+1]
			}
		} else if i.B > 1 {
			for l1 := int32(0); l1 < i.B-1; l1++ {
				if l1 >= int32(v.C.Function.MaxStackSize) {
					break
				}
				v.C.Stack.Regs[l1] = c.Stack.Regs[l1+int32(i.A)+1]
			}
		}
		return
	}
	if function.Type == GOFUNCTION {
		function.Val.(GOFUNC)(c, v)
	}
}

func Self(i *Instr, c *Closure, v *VM) {
	c.Stack.Regs[i.A+1] = &Value{
		Type: c.Stack.Regs[i.B].Type,
		Val:  c.Stack.Regs[i.B].Val,
	}
	var cval *Value
	if i.C&256 == 256 {
		cval = &c.Function.Constants[i.C&255]
	} else {
		cval = c.Stack.Regs[i.C]
	}

	val := c.Stack.Regs[i.B].Val.(*Table).Get(*cval)
	c.Stack.Regs[i.A] = &Value{
		Type: val.Type,
		Val:  val.Val,
	}
}

func Eq(i *Instr, c *Closure, v *VM) {
	var bval *Value
	if i.B&256 == 256 {
		bval = &c.Function.Constants[i.B&255]
	} else {
		bval = c.Stack.Regs[i.B]
	}
	var cval *Value
	if i.C&256 == 256 {
		cval = &c.Function.Constants[i.C&255]
	} else {
		cval = c.Stack.Regs[i.C]
	}

	if bval.Type != cval.Type {
		if i.A != 0 {
			c.PC = c.PC + 1
		}
		return
	}
	switch bval.Type {
	case NUMBER:
		if bval.Val.(Number) != cval.Val.(Number) {
			if i.A != 0 {
				c.PC = c.PC + 1
			}
			return
		}
	case STRING:
		if bval.Val.(string) != cval.Val.(string) {
			if i.A != 0 {
				c.PC = c.PC + 1
			}
			return
		}
	case BOOLEAN:
		if bval.Val.(Integer) != cval.Val.(Integer) {
			if i.A != 0 {
				c.PC = c.PC + 1
			}
			return
		}
	}
}

func Lt(i *Instr, c *Closure, v *VM) {
	var bval *Value
	if i.B&256 == 256 {
		bval = &c.Function.Constants[i.B&255]
	} else {
		bval = c.Stack.Regs[i.B]
	}
	var cval *Value
	if i.C&256 == 256 {
		cval = &c.Function.Constants[i.C&255]
	} else {
		cval = c.Stack.Regs[i.C]
	}

	if bval.Type != cval.Type {
		panic("it just don't make sense")
		return
	}
	switch bval.Type {
	case NUMBER:
		if bval.Val.(Number) >= cval.Val.(Number) {
			if i.A != 0 {
				c.PC = c.PC + 1
			}
			return
		}
	case STRING:
		if bval.Val.(string) != cval.Val.(string) {
			panic("it just don't make sense")
			return
		}
	case BOOLEAN:
		if bval.Val.(Integer) != cval.Val.(Integer) {
			panic("it just don't make sense")
			return
		}
	}
}

func Le(i *Instr, c *Closure, v *VM) {
	var bval *Value
	if i.B&256 == 256 {
		bval = &c.Function.Constants[i.B&255]
	} else {
		bval = c.Stack.Regs[i.B]
	}
	var cval *Value
	if i.C&256 == 256 {
		cval = &c.Function.Constants[i.C&255]
	} else {
		cval = c.Stack.Regs[i.C]
	}

	if bval.Type != cval.Type {
		panic("it just don't make sense")
		return
	}
	switch bval.Type {
	case NUMBER:
		if bval.Val.(Number) > cval.Val.(Number) {
			if i.A != 0 {
				c.PC = c.PC + 1
			}
			return
		}
	case STRING:
		panic("it just don't make sense")
	case BOOLEAN:
		panic("it just don't make sense")
	}
}

func Test(i *Instr, c *Closure, v *VM) {
	val := c.Stack.Regs[i.A]
	switch val.Type {
	case NIL:
		if i.C == 0 {
			c.PC = c.PC + 1
		}
	case BOOLEAN:
		if Integer(i.C) != val.Val.(Integer) {
			c.PC = c.PC + 1
		}
	case NUMBER:
		boolval := Integer(1)
		if val.Val.(Number) == 0 {
			boolval = 0
		}
		if Integer(i.C) != boolval {
			c.PC = c.PC + 1
		}
	default:
		if i.C != 0 {
			c.PC = c.PC + 1
		}
	}
}

func TestSet(i *Instr, c *Closure, v *VM) {
	val := c.Stack.Regs[i.B]
	switch val.Type {
	case NIL:
		if i.C != 0 {
			c.Stack.Regs[i.A] = &Value{
				Type: val.Type,
				Val:  val.Val,
			}
		} else {
			c.PC = c.PC + 1
		}
	case BOOLEAN:
		if Integer(i.C) == val.Val.(Integer) {
			c.Stack.Regs[i.A] = &Value{
				Type: val.Type,
				Val:  val.Val,
			}
		} else {
			c.PC = c.PC + 1
		}
	case NUMBER:
		boolval := Integer(1)
		if val.Val.(Number) == 0 {
			boolval = 0
		}
		if Integer(i.C) == boolval {
			c.Stack.Regs[i.A] = &Value{
				Type: val.Type,
				Val:  val.Val,
			}
		} else {
			c.PC = c.PC + 1
		}
	default:
		if i.C == 0 {
			c.Stack.Regs[i.A] = &Value{
				Type: val.Type,
				Val:  val.Val,
			}
		} else {
			c.PC = c.PC + 1
		}
	}
}

func ForPrep(i *Instr, c *Closure, v *VM) {
	c.Stack.Regs[i.A].Val = c.Stack.Regs[i.A].Val.(Number) - c.Stack.Regs[i.A+2].Val.(Number)
	c.PC += int64(i.B)
}

func ForLoop(i *Instr, c *Closure, v *VM) {

	c.Stack.Regs[i.A].Val = c.Stack.Regs[i.A].Val.(Number) + c.Stack.Regs[i.A+2].Val.(Number)

	dirp := true
	if c.Stack.Regs[i.A+2].Val.(Number) < 0 {
		dirp = false
	}

	passed := false
	if dirp {
		if c.Stack.Regs[i.A].Val.(Number) > c.Stack.Regs[i.A+1].Val.(Number) {
			passed = true
		}
	} else {
		if c.Stack.Regs[i.A].Val.(Number) < c.Stack.Regs[i.A+1].Val.(Number) {
			passed = true
		}
	}
	if !passed {
		c.Stack.Regs[i.A+3] = &Value{
			Type: c.Stack.Regs[i.A].Type,
			Val:  c.Stack.Regs[i.A].Val,
		}
		c.PC += int64(i.B)
	}
}

func TForLoop(i *Instr, c *Closure, v *VM) {
	function := c.Stack.Regs[i.A]
	if function.Type == CLOSURE {
		v.ClosureStack = append(v.ClosureStack, v.C)
		v.C = function.Val.(*Closure)
		v.C.ReturnCount = uint64(i.C) + 1
		v.C.ReturnPos = uint64(i.A) + 3
		v.C.ReturnFunc = func(c *Closure, v *VM) {
			if c.Stack.Regs[i.A+3].Type != NIL {
				c.Stack.Regs[i.A+2] = &Value{
					Type: c.Stack.Regs[i.A+3].Type,
					Val:  c.Stack.Regs[i.A+3].Val,
				}
			}
		}
		v.C.Stack.Regs = make([]*Value, v.C.Function.MaxStackSize)

		for l1 := int32(0); l1 < 2; l1++ {
			if l1 >= int32(v.C.Function.MaxStackSize) {
				break
			}
			v.C.Stack.Regs[l1] = c.Stack.Regs[l1+int32(i.A)+1]
		}
		return
	}
	/*if function.Type == GOFUNCTION {
		function.Val.(GOFUNC)(c, v)
	}*/
}

func NewTable(i *Instr, c *Closure, v *VM) {
	t := &Table{}
	x := i.B & 7
	e := i.B >> 3
	arraySize := x
	if e != 0 {
		arraySize = int(math.Pow(float64((x+8)*2), float64(e-1)))
	}
	t.ArraySize = arraySize
	t.Array = make([]*Value, arraySize)
	t.Hash = make(map[Value]*Value)
	c.Stack.Regs[i.A] = &Value{
		Type: TABLE,
		Val:  t,
	}
}

func SetList(i *Instr, c *Closure, v *VM) {
	t := &Table{}
	x := i.B & 7
	e := i.B >> 3
	arraySize := x
	if e != 0 {
		arraySize = int(math.Pow(float64((x+8)*2), float64(e-1)))
	}
	t.ArraySize = arraySize
	t.Array = make([]*Value, arraySize)
	t.Hash = make(map[Value]*Value)
	c.Stack.Regs[i.A] = &Value{
		Type: TABLE,
		Val:  t,
	}
}
