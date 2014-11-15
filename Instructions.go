package LuaVM

type Stackframe struct {
	Regs     []*Value
	Upvalues []*Value
}

type Closure struct {
	Stack    *Stackframe
	Function *FunctionPrototype
}
