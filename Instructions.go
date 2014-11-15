package LuaVM

type Stackframe struct {
	Regs     []*Value
	Upvalues []*Value
}

type Closure struct {
	Stack    *Stackframe
	Function *FunctionPrototype
	PC       uint64
}

type VM struct {
	G map[string]*Value
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
	if c.Function.Constants[i.B].Type != string {
		panic("Constant type is not string")
	}
	c.Stack.Regs[i.A] = &Value{
		Type: v.G[c.Function.Constants[i.B]].Type,
		Val:  v.G[c.Function.Constants[i.B]].Val,
	}
}

func SetGlobal(i *Instr, c *Closure, v *VM) {
	if c.Function.Constants[i.B].Type != string {
		panic("Constant type is not string")
	}
	v.G[c.Function.Constants[i.B]] = &Value{
		Type: c.Stack.Regs[i.A].Type,
		Val:  vc.Stack.Regs[i.A].Val,
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
