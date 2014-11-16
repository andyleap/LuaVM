}


type VM struct {
	G          map[string]*Value
	FrameStack []*Stackframe
	S          *Stackframe
}

func NewVM() *VM {
	return &VM{
		G: make(map[string]*Value),
	}
}

func (v *VM) RunClosure(c *Closure) {
	v.S = &Stackframe{
		Closure: c,
		Regs:    make([]*Value, c.Function.MaxStackSize),
	}
	v.DispatchLoop()
}
	

	

		switch i.Opcode {
		case OP_MOVE:
			Op_Move(i, v.S, v)
		case OP_LOADK:
		Op_LoadK(i, v.S, v)
		Op_LoadBool(i, v.S, v)
			Op_LoadK(i, v.S, v)
			Op_LoadBool(i, v.S, v)
		Op_LoadNil(i, v.S, v)
		Op_GetUpVal(i, v.S, v)
			Op_LoadNil(i, v.S, v)
			Op_GetUpVal(i, v.S, v)
		Op_GetGlobal(i, v.S, v)
		Op_GetTable(i, v.S, v)
			Op_GetGlobal(i, v.S, v)
			Op_GetTable(i, v.S, v)
		Op_SetGlobal(i, v.S, v)
		Op_SetUpVal(i, v.S, v)
			Op_SetGlobal(i, v.S, v)
			Op_SetUpVal(i, v.S, v)
		Op_SetTable(i, v.S, v)
		Op_NewTable(i, v.S, v)
			Op_SetTable(i, v.S, v)
			Op_NewTable(i, v.S, v)
		Op_Self(i, v.S, v)
		Op_Add(i, v.S, v)
			Op_Self(i, v.S, v)
			Op_Add(i, v.S, v)
		Op_Sub(i, v.S, v)
		Op_Mul(i, v.S, v)
			Op_Sub(i, v.S, v)
			Op_Mul(i, v.S, v)
		Op_Div(i, v.S, v)
		Op_Mod(i, v.S, v)
			Op_Div(i, v.S, v)
			Op_Mod(i, v.S, v)
		Op_Pow(i, v.S, v)
		Op_Unm(i, v.S, v)
			Op_Pow(i, v.S, v)
			Op_Unm(i, v.S, v)
		Op_Not(i, v.S, v)
		Op_Len(i, v.S, v)
			Op_Not(i, v.S, v)
			Op_Len(i, v.S, v)
		Op_Concat(i, v.S, v)
		Op_Jmp(i, v.S, v)
			Op_Concat(i, v.S, v)
			Op_Jmp(i, v.S, v)
		Op_Eq(i, v.S, v)
		Op_Lt(i, v.S, v)
			Op_Eq(i, v.S, v)
			Op_Lt(i, v.S, v)
		Op_Le(i, v.S, v)
		Op_Test(i, v.S, v)
			Op_Le(i, v.S, v)
			Op_Test(i, v.S, v)
		Op_TestSet(i, v.S, v)
		Op_Call(i, v.S, v)
			Op_TestSet(i, v.S, v)
			Op_Call(i, v.S, v)
		Op_TailCall(i, v.S, v)
			Op_TailCall(i, v.S, v)
		case OP_FORLOOP:
			Op_Return(i, v.S, v)
		case OP_FORPREP:
			Op_ForLoop(i, v.S, v)
		case OP_TFORLOOP:
			Op_ForPrep(i, v.S, v)
		case OP_SETLIST:
			Op_TForLoop(i, v.S, v)
		case OP_CLOSE:
			Op_SetList(i, v.S, v)
		case OP_CLOSURE:
			Op_Close(i, v.S, v)
		case OP_VARARG:
			Op_Closure(i, v.S, v)
		}
			Op_Vararg(i, v.S, v)
		Op_Vararg(i, v.S, v)
	