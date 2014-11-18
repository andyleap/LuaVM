package LuaVM

import "math"

type OPCODE int

const (
	OP_MOVE OPCODE = iota
	OP_LOADK
	OP_LOADBOOL
	OP_LOADNIL
	OP_GETUPVAL
	OP_GETGLOBAL
	OP_GETTABLE
	OP_SETGLOBAL
	OP_SETUPVAL
	OP_SETTABLE
	OP_NEWTABLE
	OP_SELF
	OP_ADD
	OP_SUB
	OP_MUL
	OP_DIV
	OP_MOD
	OP_POW
	OP_UNM
	OP_NOT
	OP_LEN
	OP_CONCAT
	OP_JMP
	OP_EQ
	OP_LT
	OP_LE
	OP_TEST
	OP_TESTSET
	OP_CALL
	OP_TAILCALL
	OP_RETURN
	OP_FORLOOP
	OP_FORPREP
	OP_TFORLOOP
	OP_SETLIST
	OP_CLOSE
	OP_CLOSURE
	OP_VARARG
)

type Stackframe struct {
	Regs        []*Value
	Params      []*Value
	Closure     *Closure
	PC          int64
	ReturnCount uint64
	ReturnPos   uint64
	ReturnFunc  func(*Stackframe, *VM)
}

type Closure struct {
	Upvalues []*Value
	Function *FunctionPrototype
}

func Op_Move(i *Instr, s *Stackframe, v *VM) {
	s.Regs[i.A] = s.Regs[i.B].Copy()
}

func Op_LoadNil(i *Instr, s *Stackframe, v *VM) {
	for l1 := int32(i.A); l1 <= i.B; l1++ {
		s.Regs[l1] = &Value{
			Type: NIL,
		}
	}
}

func Op_LoadK(i *Instr, s *Stackframe, v *VM) {
	s.Regs[i.A] = s.Closure.Function.Constants[i.B].Copy()
}

func Op_LoadBool(i *Instr, s *Stackframe, v *VM) {
	s.Regs[i.A] = &Value{
		Type: BOOLEAN,
		Val:  Integer(i.B),
	}
	if i.C != 0 {
		s.PC++
	}
}

func Op_GetGlobal(i *Instr, s *Stackframe, v *VM) {
	if s.Closure.Function.Constants[i.B].Type != STRING {
		panic("Constant type is not string")
	}
	s.Regs[i.A] = v.G.Get(s.Closure.Function.Constants[i.B]).Copy()
}

func Op_SetGlobal(i *Instr, s *Stackframe, v *VM) {
	if s.Closure.Function.Constants[i.B].Type != STRING {
		panic("Constant type is not string")
	}
	v.G.Set(s.Closure.Function.Constants[i.B], s.Regs[i.A].Copy())
}

func Op_GetUpVal(i *Instr, s *Stackframe, v *VM) {
	s.Regs[i.A] = s.Closure.Upvalues[i.B].Copy()
}

func Op_SetUpVal(i *Instr, s *Stackframe, v *VM) {
	s.Closure.Upvalues[i.B] = s.Regs[i.A].Copy()
}

func Op_GetTable(i *Instr, s *Stackframe, v *VM) {
	var key *Value
	if i.C&256 == 256 {
		key = &s.Closure.Function.Constants[i.C&255]
	} else {
		key = s.Regs[i.C]
	}
	if s.Regs[i.B].Type != TABLE {
		panic("Value is not table")
	}
	val := s.Regs[i.B].Val.(*Table).Get(*key)
	s.Regs[i.A] = val.Copy()
}

func Op_SetTable(i *Instr, s *Stackframe, v *VM) {
	var key *Value
	if i.B&256 == 256 {
		key = &s.Closure.Function.Constants[i.B&255]
	} else {
		key = s.Regs[i.B]
	}
	var val *Value
	if i.C&256 == 256 {
		val = &s.Closure.Function.Constants[i.C&255]
	} else {
		val = s.Regs[i.C]
	}
	if s.Regs[i.A].Type != TABLE {
		panic("Value is not table")
	}
	s.Regs[i.A].Val.(*Table).Set(*key, val.Copy())
}

func Op_Add(i *Instr, s *Stackframe, v *VM) {
	var bval *Value
	if i.B&256 == 256 {
		bval = &s.Closure.Function.Constants[i.B&255]
	} else {
		bval = s.Regs[i.B]
	}
	var cval *Value
	if i.C&256 == 256 {
		cval = &s.Closure.Function.Constants[i.C&255]
	} else {
		cval = s.Regs[i.C]
	}
	if bval.Type != NUMBER || cval.Type != NUMBER {
		panic("Trying to add non-numbers")
	}
	s.Regs[i.A] = &Value{
		Type: NUMBER,
		Val:  bval.Val.(Number) + cval.Val.(Number),
	}
}

func Op_Sub(i *Instr, s *Stackframe, v *VM) {
	var bval *Value
	if i.B&256 == 256 {
		bval = &s.Closure.Function.Constants[i.B&255]
	} else {
		bval = s.Regs[i.B]
	}
	var cval *Value
	if i.C&256 == 256 {
		cval = &s.Closure.Function.Constants[i.C&255]
	} else {
		cval = s.Regs[i.C]
	}
	if bval.Type != NUMBER || cval.Type != NUMBER {
		panic("Trying to sub non-numbers")
	}
	s.Regs[i.A] = &Value{
		Type: NUMBER,
		Val:  bval.Val.(Number) - cval.Val.(Number),
	}
}

func Op_Mul(i *Instr, s *Stackframe, v *VM) {
	var bval *Value
	if i.B&256 == 256 {
		bval = &s.Closure.Function.Constants[i.B&255]
	} else {
		bval = s.Regs[i.B]
	}
	var cval *Value
	if i.C&256 == 256 {
		cval = &s.Closure.Function.Constants[i.C&255]
	} else {
		cval = s.Regs[i.C]
	}
	if bval.Type != NUMBER || cval.Type != NUMBER {
		panic("Trying to mul non-numbers")
	}
	s.Regs[i.A] = &Value{
		Type: NUMBER,
		Val:  bval.Val.(Number) * cval.Val.(Number),
	}
}

func Op_Div(i *Instr, s *Stackframe, v *VM) {
	var bval *Value
	if i.B&256 == 256 {
		bval = &s.Closure.Function.Constants[i.B&255]
	} else {
		bval = s.Regs[i.B]
	}
	var cval *Value
	if i.C&256 == 256 {
		cval = &s.Closure.Function.Constants[i.C&255]
	} else {
		cval = s.Regs[i.C]
	}
	if bval.Type != NUMBER || cval.Type != NUMBER {
		panic("Trying to div non-numbers")
	}
	s.Regs[i.A] = &Value{
		Type: NUMBER,
		Val:  bval.Val.(Number) / cval.Val.(Number),
	}
}

func Op_Mod(i *Instr, s *Stackframe, v *VM) {
	var bval *Value
	if i.B&256 == 256 {
		bval = &s.Closure.Function.Constants[i.B&255]
	} else {
		bval = s.Regs[i.B]
	}
	var cval *Value
	if i.C&256 == 256 {
		cval = &s.Closure.Function.Constants[i.C&255]
	} else {
		cval = s.Regs[i.C]
	}
	if bval.Type != NUMBER || cval.Type != NUMBER {
		panic("Trying to add non-numbers")
	}
	s.Regs[i.A] = &Value{
		Type: NUMBER,
		Val:  Number(math.Mod(bval.Val.(float64), cval.Val.(float64))),
	}
}

func Op_Pow(i *Instr, s *Stackframe, v *VM) {
	var bval *Value
	if i.B&256 == 256 {
		bval = &s.Closure.Function.Constants[i.B&255]
	} else {
		bval = s.Regs[i.B]
	}
	var cval *Value
	if i.C&256 == 256 {
		cval = &s.Closure.Function.Constants[i.C&255]
	} else {
		cval = s.Regs[i.C]
	}
	if bval.Type != NUMBER || cval.Type != NUMBER {
		panic("Trying to add non-numbers")
	}
	s.Regs[i.A] = &Value{
		Type: NUMBER,
		Val:  Number(math.Pow(bval.Val.(float64), cval.Val.(float64))),
	}
}

func Op_Unm(i *Instr, s *Stackframe, v *VM) {
	if s.Regs[i.B].Type != NUMBER {
		panic("Trying to unm non-number")
	}
	s.Regs[i.A] = &Value{
		Type: NUMBER,
		Val:  -(s.Regs[i.B].Val.(Number)),
	}
}

func Op_Not(i *Instr, s *Stackframe, v *VM) {
	bval := s.Regs[i.B]
	if bval.Type == NIL {
		s.Regs[i.A] = &Value{
			Type: BOOLEAN,
			Val:  Integer(1),
		}
	}
	if bval.Type == BOOLEAN {
		val := 1
		if bval.Val != 0 {
			val = 0
		}
		s.Regs[i.A] = &Value{
			Type: BOOLEAN,
			Val:  Integer(val),
		}
	}

}

func Op_Len(i *Instr, s *Stackframe, v *VM) {
	bval := s.Regs[i.B]
	var val *Value
	if bval.Type == TABLE {
		val = bval.Val.(*Table).Len()
	}
	if bval.Type == STRING {
		val = &Value{Type: NUMBER, Val: Number(len(bval.Val.(string)))}
	}
	s.Regs[i.A] = val
}

func Op_Concat(i *Instr, s *Stackframe, v *VM) {
	str := ""
	for l1 := int32(i.B); l1 <= int32(i.C); l1++ {
		if s.Regs[l1].Type != STRING {
			panic("Attempting to concat non-strings")
		}
		str = str + s.Regs[l1].Val.(string)
	}
	s.Regs[i.A] = &Value{
		Type: STRING,
		Val:  str,
	}
}

func Op_Jmp(i *Instr, s *Stackframe, v *VM) {
	s.PC = s.PC + int64(i.B)
}

func Op_Call(i *Instr, s *Stackframe, v *VM) {
	function := s.Regs[i.A]
	if function.Type == CLOSURE {
		v.FrameStack = append(v.FrameStack, v.S)
		v.S = &Stackframe{
			Closure: function.Val.(*Closure),
			PC:      0,
		}
		v.S.ReturnCount = uint64(i.C)
		v.S.ReturnPos = uint64(i.A)
		v.S.Regs = make([]*Value, v.S.Closure.Function.MaxStackSize)
		if i.B == 0 {
			for l1 := int32(0); l1+int32(i.A)+1 < int32(len(s.Regs)); l1++ {
				if l1 >= int32(len(v.S.Regs)) {
					v.S.Regs = append(v.S.Regs, s.Regs[l1+int32(i.A)+1])
				} else {
					v.S.Regs[l1] = s.Regs[l1+int32(i.A)+1]
				}
			}
			v.S.Params = s.Regs[i.A+1:]
		} else if i.B > 1 {
			for l1 := int32(0); l1 < i.B-1; l1++ {
				if l1 >= int32(len(v.S.Regs)) {
					v.S.Regs = append(v.S.Regs, s.Regs[l1+int32(i.A)+1])
				} else {
					v.S.Regs[l1] = s.Regs[l1+int32(i.A)+1]
				}
			}
			v.S.Params = s.Regs[i.A+1 : i.A+uint8(i.B)]
		}
		return
	}
	if function.Type == GOFUNCTION {
		var params []*Value
		if i.B == 0 {
			params = s.Regs[i.A+1:]
		} else {
			params = s.Regs[i.A+1 : i.A+uint8(i.B)]
		}
		function.Val.(GOFUNC)(params, v)

	}
}

func Op_Return(i *Instr, s *Stackframe, v *VM) {
	if len(v.FrameStack) == 0 {
		v.S = nil
		return
	}
	v.S = v.FrameStack[len(v.FrameStack)-1]
	v.FrameStack = v.FrameStack[:len(v.FrameStack)-1]

	if i.B == 0 {
		for l1 := int32(0); l1+int32(i.A) < int32(len(s.Regs)); l1++ {
			if s.ReturnCount > 0 && l1 >= int32(s.ReturnCount) {
				break
			}
			if len(v.S.Regs) <= int(l1+int32(s.ReturnPos)) {
				v.S.Regs = append(v.S.Regs, s.Regs[l1+int32(i.A)])
			} else {
				v.S.Regs[l1+int32(s.ReturnPos)] = s.Regs[l1+int32(i.A)]
			}
		}
	} else if i.B > 1 {
		for l1 := int32(0); l1 < i.B-1; l1++ {
			if s.ReturnCount > 0 && l1 >= int32(s.ReturnCount) {
				break
			}
			if len(v.S.Regs) <= int(l1+int32(s.ReturnPos)) {
				v.S.Regs = append(v.S.Regs, s.Regs[l1+int32(i.A)])
			} else {
				v.S.Regs[l1+int32(s.ReturnPos)] = s.Regs[l1+int32(i.A)]
			}
		}
	}
	if s.ReturnFunc != nil {
		s.ReturnFunc(v.S, v)
	}
}

func Op_TailCall(i *Instr, s *Stackframe, v *VM) {
	function := s.Regs[i.A]
	if function.Type == CLOSURE {
		v.S = &Stackframe{
			Closure: function.Val.(*Closure),
			PC:      0,
		}
		v.S.ReturnCount = uint64(i.C)
		v.S.ReturnPos = uint64(i.A)
		v.S.Regs = make([]*Value, v.S.Closure.Function.MaxStackSize)
		if i.B == 0 {
			for l1 := int32(0); l1+int32(i.A)+1 < int32(len(s.Regs)); l1++ {
				if l1 >= int32(len(v.S.Regs)) {
					v.S.Regs = append(v.S.Regs, s.Regs[l1+int32(i.A)+1])
				} else {
					v.S.Regs[l1] = s.Regs[l1+int32(i.A)+1]
				}
			}
			v.S.Params = s.Regs[i.A+1:]
		} else if i.B > 1 {
			for l1 := int32(0); l1 < i.B-1; l1++ {
				if l1 >= int32(len(v.S.Regs)) {
					v.S.Regs = append(v.S.Regs, s.Regs[l1+int32(i.A)+1])
				} else {
					v.S.Regs[l1] = s.Regs[l1+int32(i.A)]
				}
			}
			v.S.Params = s.Regs[i.A+1 : i.A+uint8(i.B)+1]
		}
		return
	}
	if function.Type == GOFUNCTION {
		var params []*Value
		if i.B == 0 {
			params = s.Regs[i.A+1:]
		} else {
			params = s.Regs[i.A+1 : i.A+uint8(i.B)]
		}
		function.Val.(GOFUNC)(params, v)
	}
}

func Op_Self(i *Instr, s *Stackframe, v *VM) {
	s.Regs[i.A+1] = s.Regs[i.B].Copy()
	var cval *Value
	if i.C&256 == 256 {
		cval = &s.Closure.Function.Constants[i.C&255]
	} else {
		cval = s.Regs[i.C]
	}

	val := s.Regs[i.B].Val.(*Table).Get(*cval)
	s.Regs[i.A] = val.Copy()
}

func Op_Eq(i *Instr, s *Stackframe, v *VM) {
	var bval *Value
	if i.B&256 == 256 {
		bval = &s.Closure.Function.Constants[i.B&255]
	} else {
		bval = s.Regs[i.B]
	}
	var cval *Value
	if i.C&256 == 256 {
		cval = &s.Closure.Function.Constants[i.C&255]
	} else {
		cval = s.Regs[i.C]
	}

	if bval.Type != cval.Type {
		if i.A != 0 {
			s.PC = s.PC + 1
		}
		return
	}
	switch bval.Type {
	case NUMBER:
		if bval.Val.(Number) != cval.Val.(Number) {
			if i.A != 0 {
				s.PC = s.PC + 1
			}
			return
		}
	case STRING:
		if bval.Val.(string) != cval.Val.(string) {
			if i.A != 0 {
				s.PC = s.PC + 1
			}
			return
		}
	case BOOLEAN:
		if bval.Val.(Integer) != cval.Val.(Integer) {
			if i.A != 0 {
				s.PC = s.PC + 1
			}
			return
		}
	}
}

func Op_Lt(i *Instr, s *Stackframe, v *VM) {
	var bval *Value
	if i.B&256 == 256 {
		bval = &s.Closure.Function.Constants[i.B&255]
	} else {
		bval = s.Regs[i.B]
	}
	var cval *Value
	if i.C&256 == 256 {
		cval = &s.Closure.Function.Constants[i.C&255]
	} else {
		cval = s.Regs[i.C]
	}

	if bval.Type != cval.Type {
		panic("it just don't make sense")
		return
	}
	switch bval.Type {
	case NUMBER:
		if bval.Val.(Number) >= cval.Val.(Number) {
			if i.A != 0 {
				s.PC = s.PC + 1
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

func Op_Le(i *Instr, s *Stackframe, v *VM) {
	var bval *Value
	if i.B&256 == 256 {
		bval = &s.Closure.Function.Constants[i.B&255]
	} else {
		bval = s.Regs[i.B]
	}
	var cval *Value
	if i.C&256 == 256 {
		cval = &s.Closure.Function.Constants[i.C&255]
	} else {
		cval = s.Regs[i.C]
	}

	if bval.Type != cval.Type {
		panic("it just don't make sense")
		return
	}
	switch bval.Type {
	case NUMBER:
		if bval.Val.(Number) > cval.Val.(Number) {
			if i.A != 0 {
				s.PC = s.PC + 1
			}
			return
		}
	case STRING:
		panic("it just don't make sense")
	case BOOLEAN:
		panic("it just don't make sense")
	}
}

func Op_Test(i *Instr, s *Stackframe, v *VM) {
	val := s.Regs[i.A]
	switch val.Type {
	case NIL:
		if i.C != 0 {
			s.PC = s.PC + 1
		}
	case BOOLEAN:
		if Integer(i.C) != val.Val.(Integer) {
			s.PC = s.PC + 1
		}
	case NUMBER:
		boolval := Integer(1)
		if val.Val.(Number) == 0 {
			boolval = 0
		}
		if Integer(i.C) != boolval {
			s.PC = s.PC + 1
		}
	default:
		if i.C != 0 {
			s.PC = s.PC + 1
		}
	}
}

func Op_TestSet(i *Instr, s *Stackframe, v *VM) {
	val := s.Regs[i.B]
	switch val.Type {
	case NIL:
		if i.C != 0 {
			s.Regs[i.A] = val.Copy()
		} else {
			s.PC = s.PC + 1
		}
	case BOOLEAN:
		if Integer(i.C) == val.Val.(Integer) {
			s.Regs[i.A] = val.Copy()
		} else {
			s.PC = s.PC + 1
		}
	case NUMBER:
		boolval := Integer(1)
		if val.Val.(Number) == 0 {
			boolval = 0
		}
		if Integer(i.C) == boolval {
			s.Regs[i.A] = val.Copy()
		} else {
			s.PC = s.PC + 1
		}
	default:
		if i.C == 0 {
			s.Regs[i.A] = val.Copy()
		} else {
			s.PC = s.PC + 1
		}
	}
}

func Op_ForPrep(i *Instr, s *Stackframe, v *VM) {
	s.Regs[i.A].Val = s.Regs[i.A].Val.(Number) - s.Regs[i.A+2].Val.(Number)
	s.PC += int64(i.B)
}

func Op_ForLoop(i *Instr, s *Stackframe, v *VM) {

	s.Regs[i.A].Val = s.Regs[i.A].Val.(Number) + s.Regs[i.A+2].Val.(Number)

	dirp := true
	if s.Regs[i.A+2].Val.(Number) < 0 {
		dirp = false
	}

	passed := false
	if dirp {
		if s.Regs[i.A].Val.(Number) > s.Regs[i.A+1].Val.(Number) {
			passed = true
		}
	} else {
		if s.Regs[i.A].Val.(Number) < s.Regs[i.A+1].Val.(Number) {
			passed = true
		}
	}
	if !passed {
		s.Regs[i.A+3] = s.Regs[i.A].Copy()
		s.PC += int64(i.B)
	}
}

func Op_TForLoop(i *Instr, s *Stackframe, v *VM) {
	function := s.Regs[i.A]
	if function.Type == CLOSURE {
		v.FrameStack = append(v.FrameStack, v.S)
		v.S = &Stackframe{
			Closure: function.Val.(*Closure),
			PC:      0,
		}
		v.S.ReturnCount = uint64(i.C)
		v.S.ReturnPos = uint64(i.A + 3)
		v.S.Regs = make([]*Value, v.S.Closure.Function.MaxStackSize)

		v.S.ReturnFunc = func(s *Stackframe, v *VM) {
			if s.Regs[i.A+3].Type != NIL {
				s.Regs[i.A+2] = s.Regs[i.A+3].Copy()
			} else {
				s.PC++
			}
		}

		for l1 := int32(0); l1 < 2; l1++ {
			if l1 >= int32(v.S.Closure.Function.MaxStackSize) {
				break
			}
			v.S.Regs[l1] = s.Regs[l1+int32(i.A)+1]
		}
		return
	}
	/*if function.Type == GOFUNCTION {
		function.Val.(GOFUNC)(c, v)
	}*/
}

func Op_NewTable(i *Instr, s *Stackframe, v *VM) {
	t := &Table{}
	x := i.B & 7
	e := i.B >> 3
	arraySize := uint64(x)
	if e != 0 {
		arraySize = uint64(math.Pow(float64((x+8)*2), float64(e-1)))
	}
	t.ArraySize = arraySize
	t.Array = make([]*Value, arraySize)
	t.Hash = make(map[Value]*Value)
	s.Regs[i.A] = &Value{
		Type: TABLE,
		Val:  t,
	}
}

func Op_SetList(i *Instr, s *Stackframe, v *VM) {
	t := s.Regs[i.A].Val.(*Table)
	top := int(i.B)
	block := Integer(i.C)
	if top == 0 {
		top = len(s.Regs) - int(i.A+1)
	}
	if block == 0 {
		block = Integer(s.Closure.Function.Instructions[s.PC].Raw)
		s.PC++
	}
	for l1 := Integer(1); l1 <= Integer(top); l1++ {
		t.Set(
			Value{Type: NUMBER, Val: Number(l1 + ((block - 1) * 50))},
			s.Regs[i.A+uint8(l1)].Copy())
	}
}

func Op_Closure(i *Instr, s *Stackframe, v *VM) {
	closure := &Closure{
		Function: s.Closure.Function.Functions[i.B],
	}
	destReg := i.A
	closure.Upvalues = make([]*Value, closure.Function.Upvalues)
	for l1 := uint8(0); l1 < closure.Function.Upvalues; l1++ {
		subi := s.Closure.Function.Instructions[s.PC]
		if subi.Opcode == OP_GETUPVAL {
			closure.Upvalues[l1] = s.Closure.Upvalues[subi.B]
		} else if subi.Opcode == OP_MOVE {
			closure.Upvalues[l1] = s.Regs[subi.B]
		} else {
			panic("Invalid upval reg code")
		}
		s.PC++
	}
	s.Regs[destReg] = &Value{Type: CLOSURE, Val: closure}
}

func Op_Close(i *Instr, s *Stackframe, v *VM) {
}

func Op_Vararg(i *Instr, s *Stackframe, v *VM) {
	if i.B == 0 {
		for l1 := int32(0); l1 < int32(len(s.Params)); l1++ {
			if l1+int32(i.A) >= int32(len(v.S.Regs)) {
				v.S.Regs = append(v.S.Regs, s.Params[l1])
			} else {
				v.S.Regs[l1+int32(i.A)] = s.Params[l1]
			}
		}
	} else {
		for l1 := int32(0); l1 < i.B-1; l1++ {
			if l1+int32(i.A) >= int32(len(v.S.Regs)) {
				v.S.Regs = append(v.S.Regs, s.Params[l1])
			} else {
				v.S.Regs[l1+int32(i.A)] = s.Params[l1]
			}
		}
	}

}

func (v *Value) Copy() *Value {
	return &Value{
		Type: v.Type,
		Val:  v.Val,
	}
}
