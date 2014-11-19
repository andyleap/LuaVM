// LuaVM project LuaVM.go
package LuaVM

import (
	"fmt"
	"os"
	"testing"
)

func TestReadLuaC(t *testing.T) {
	f, err := os.Open("test.luac")
	if err != nil {
		t.Error("File Open Failed: ", err)
		return
	}
	c, err := ReadLuaC(f)
	if err != nil {
		t.Error("File Read Failed: ", err)
		return
	}
	vm := NewVM()
	vm.G.SetFunc("print", lua_print)
	vm.RunClosure(c)
}

func lua_print(params []*Value, v *VM) []*Value {
	for _, v := range params {
		fmt.Print(v.Val)
	}
	fmt.Println()
	return nil
}
