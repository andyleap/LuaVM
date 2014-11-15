// LuaVM project LuaVM.go
package LuaVM

import (
	"os"
	"testing"
)

func TestTimeConsuming(t *testing.T) {
	f, err := os.Open("test.luac")
	if err != nil {
		t.Error("File Open Failed: ", err)
		return
	}
	_, err = ReadLuaC(f)
	if err != nil {
		t.Error("File Read Failed: ", err)
		return
	}

}
