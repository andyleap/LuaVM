package LuaVM

type VM struct {
	G          *Table
	FrameStack []*Stackframe
	S          *Stackframe
}

func NewVM() *VM {
	vm := &VM{
		G: NewTable(),
	}
	vm.G.SetFunc("getmetatable", getmetatable)
	vm.G.SetFunc("setmetatable", setmetatable)

	return vm
}

func (v *VM) RunClosure(c *Closure) {
	v.S = &Stackframe{
		Closure: c,
		Regs:    make([]*Value, c.Function.MaxStackSize),
	}
	v.DispatchLoop()
}

func (v *VM) runClosure(c *Closure, params []*Value, returnfunc func(*Stackframe, *VM, []*Value)) *Stackframe {
	s := &Stackframe{
		Closure:    c,
		Regs:       make([]*Value, c.Function.MaxStackSize),
		Params:     make([]*Value, len(params)),
		ReturnFunc: returnfunc,
	}
	for k, val := range params {
		s.Regs[k] = val.Copy()
		s.Params[k] = val.Copy()
	}
	return s
}

func (v *VM) DispatchLoop() {
	for {
		i := &v.S.Closure.Function.Instructions[v.S.PC]
		v.S.PC++
		switch i.Opcode {
		case OP_MOVE:
			Op_Move(i, v.S, v)
		case OP_LOADK:
			Op_LoadK(i, v.S, v)
		case OP_LOADBOOL:
			Op_LoadBool(i, v.S, v)
		case OP_LOADNIL:
			Op_LoadNil(i, v.S, v)
		case OP_GETUPVAL:
			Op_GetUpVal(i, v.S, v)
		case OP_GETGLOBAL:
			Op_GetGlobal(i, v.S, v)
		case OP_GETTABLE:
			Op_GetTable(i, v.S, v)
		case OP_SETGLOBAL:
			Op_SetGlobal(i, v.S, v)
		case OP_SETUPVAL:
			Op_SetUpVal(i, v.S, v)
		case OP_SETTABLE:
			Op_SetTable(i, v.S, v)
		case OP_NEWTABLE:
			Op_NewTable(i, v.S, v)
		case OP_SELF:
			Op_Self(i, v.S, v)
		case OP_ADD:
			Op_Add(i, v.S, v)
		case OP_SUB:
			Op_Sub(i, v.S, v)
		case OP_MUL:
			Op_Mul(i, v.S, v)
		case OP_DIV:
			Op_Div(i, v.S, v)
		case OP_MOD:
			Op_Mod(i, v.S, v)
		case OP_POW:
			Op_Pow(i, v.S, v)
		case OP_UNM:
			Op_Unm(i, v.S, v)
		case OP_NOT:
			Op_Not(i, v.S, v)
		case OP_LEN:
			Op_Len(i, v.S, v)
		case OP_CONCAT:
			Op_Concat(i, v.S, v)
		case OP_JMP:
			Op_Jmp(i, v.S, v)
		case OP_EQ:
			Op_Eq(i, v.S, v)
		case OP_LT:
			Op_Lt(i, v.S, v)
		case OP_LE:
			Op_Le(i, v.S, v)
		case OP_TEST:
			Op_Test(i, v.S, v)
		case OP_TESTSET:
			Op_TestSet(i, v.S, v)
		case OP_CALL:
			Op_Call(i, v.S, v)
		case OP_TAILCALL:
			Op_TailCall(i, v.S, v)
		case OP_RETURN:
			Op_Return(i, v.S, v)
		case OP_FORLOOP:
			Op_ForLoop(i, v.S, v)
		case OP_FORPREP:
			Op_ForPrep(i, v.S, v)
		case OP_TFORLOOP:
			Op_TForLoop(i, v.S, v)
		case OP_SETLIST:
			Op_SetList(i, v.S, v)
		case OP_CLOSE:
			Op_Close(i, v.S, v)
		case OP_CLOSURE:
			Op_Closure(i, v.S, v)
		case OP_VARARG:
			Op_Vararg(i, v.S, v)
		}
		if v.S == nil {
			break
		}
	}

}
