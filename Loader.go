// LuaVM project LuaVM.go
package LuaVM

import (
	"encoding/binary"
	"errors"
	"io"
)

type ValueType uint8

const (
	NIL ValueType = iota
	BOOLEAN
	_
	NUMBER
	STRING
	TABLE
	FUNCTION
	CLOSURE
	GOFUNCTION
)

type GOFUNC func(params []*Value, v *VM) []*Value

var (
	BADSIGNATURE error = errors.New("Bad signature")
	BADVERSION   error = errors.New("Bad version")
	BADENCODING  error = errors.New("Bad encoding")
)

type Size_T uint64
type Integer uint32
type Instruction uint32
type Number float64

type Value struct {
	Type ValueType
	Val  interface{}
}

type Instr struct {
	Raw    Instruction
	Opcode OPCODE
	A      uint8
	B      int32
	C      uint16
}

type FunctionPrototype struct {
	Instructions []Instr
	Constants    []Value
	Functions    []*FunctionPrototype
	Upvalues     uint8
	Parameters   uint8
	IsVararg     uint8
	MaxStackSize uint8
}

type header struct {
	Signature        uint32
	Version          uint8
	Format           uint8
	Endianness       uint8
	Size_int         uint8
	Size_size_t      uint8
	Size_instruction uint8
	Size_number      uint8
	Integral         uint8
}

type luaFile struct {
	Size_size_t uint8
	Data        io.Reader
}

func ReadLuaC(data io.Reader) (*Closure, error) {
	l := &luaFile{Data: data}
	err := l.checkHeader()
	if err != nil {
		return nil, err
	}
	p, _ := l.readFunctionBlock()
	c := &Closure{Function: p}
	return c, nil
}

func (l *luaFile) checkHeader() error {
	var header header
	binary.Read(l.Data, binary.LittleEndian, &header)
	if header.Signature != 0x61754C1B {
		return BADSIGNATURE
	}
	if header.Version != 0x51 {
		return BADVERSION
	}
	if header.Format != 0 ||
		header.Endianness != 1 ||
		header.Size_int != 4 ||
		header.Size_instruction != 4 ||
		header.Size_number != 8 ||
		header.Integral != 0 {
		return BADENCODING
	}
	l.Size_size_t = header.Size_size_t
	return nil
}

type VarargFlag uint8

const (
	VARARG_HASARG   VarargFlag = 1
	VARARG_ISVARARG VarargFlag = 2
	VARARG_NEEDSARG VarargFlag = 4
)

type functionBlock struct {
	LineDefined     uint32
	LastLineDefined uint32
	Upvalues        uint8
	Parameters      uint8
	IsVararg        uint8
	MaxStackSize    uint8
}

func (l *luaFile) readFunctionBlock() (*FunctionPrototype, error) {
	l.readString()

	var block functionBlock
	binary.Read(l.Data, binary.LittleEndian, &block)

	Prototype := &FunctionPrototype{
		Upvalues:     block.Upvalues,
		Parameters:   block.Parameters,
		IsVararg:     block.IsVararg,
		MaxStackSize: block.MaxStackSize,
	}
	Prototype.Instructions = l.readInstructionList()
	Prototype.Constants = l.readConstantList()

	Prototype.Functions = l.readFunctionList()

	l.readSourceLinePositionList()
	l.readLocalList()
	l.readUpvalueList()

	return Prototype, nil
}

func (l *luaFile) readUpvalueList() {
	var size Integer
	binary.Read(l.Data, binary.LittleEndian, &size)
	for l1 := Integer(0); l1 < size; l1++ {
		l.readString()
	}
}

func (l *luaFile) readLocalList() {
	var size Integer
	var local Integer
	binary.Read(l.Data, binary.LittleEndian, &size)
	for l1 := Integer(0); l1 < size; l1++ {
		l.readString()
		binary.Read(l.Data, binary.LittleEndian, &local)
		binary.Read(l.Data, binary.LittleEndian, &local)
	}
}

func (l *luaFile) readSourceLinePositionList() {
	var size Integer
	var line Integer
	binary.Read(l.Data, binary.LittleEndian, &size)
	for l1 := Integer(0); l1 < size; l1++ {
		binary.Read(l.Data, binary.LittleEndian, &line)
	}
}

func (l *luaFile) readFunctionList() []*FunctionPrototype {
	var size Integer
	binary.Read(l.Data, binary.LittleEndian, &size)
	functions := make([]*FunctionPrototype, size)
	for l1 := Integer(0); l1 < size; l1++ {
		function, _ := l.readFunctionBlock()
		functions[l1] = function
	}
	return functions
}

func (l *luaFile) readInstruction() Instr {
	var instruction Instruction
	binary.Read(l.Data, binary.LittleEndian, &instruction)

	ret := Instr{}
	ret.Raw = instruction

	opcode := uint8(0x0000003F & instruction)

	ret.Opcode = OPCODE(opcode)

	switch opcode {
	case 0, 3, 2, 4, 8, 6, 9, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 28, 30, 29, 37, 11, 23, 24, 25, 26, 27, 33, 10, 34, 35: //iABC
		ret.A = uint8((instruction & 0x00003FC0) >> 6)
		ret.C = uint16((instruction & 0x007FC000) >> 14)
		ret.B = int32((instruction & 0xFF800000) >> 23)
	case 1, 5, 7, 36: //iABx
		ret.A = uint8((instruction & 0x00003FC0) >> 6)
		ret.B = int32((instruction & 0xFFFFC000) >> 14)
	case 22, 31, 32: //iAsBx
		ret.A = uint8((instruction & 0x00003FC0) >> 6)
		ret.B = int32(((instruction & 0xFFFFC000) >> 14)) - 131071
	}

	return ret
}

func (l *luaFile) readInstructionList() []Instr {
	var size Integer
	binary.Read(l.Data, binary.LittleEndian, &size)
	instructions := make([]Instr, size)
	for l1 := Integer(0); l1 < size; l1++ {
		instructions[l1] = l.readInstruction()
	}
	return instructions
}

func (l *luaFile) readConstantList() []Value {
	var size Integer
	var valuetype ValueType
	var integer Integer
	var number Number
	binary.Read(l.Data, binary.LittleEndian, &size)
	constants := make([]Value, size)
	for l1 := Integer(0); l1 < size; l1++ {
		binary.Read(l.Data, binary.LittleEndian, &valuetype)
		constants[l1].Type = valuetype
		switch {
		case valuetype == NIL:
		case valuetype == BOOLEAN:
			binary.Read(l.Data, binary.LittleEndian, &integer)
			constants[l1].Val = integer
		case valuetype == NUMBER:
			binary.Read(l.Data, binary.LittleEndian, &number)
			constants[l1].Val = number
		case valuetype == STRING:
			constants[l1].Val = l.readString()
		}
	}
	return constants
}

func (l *luaFile) readString() string {
	size := l.readSize_T()
	str := make([]byte, size)
	l.Data.Read(str)
	if size > 0 {
		str = str[:size-1]
	}
	return string(str)
}

func (l *luaFile) readSize_T() Size_T {
	if l.Size_size_t == 4 {
		var size uint32
		binary.Read(l.Data, binary.LittleEndian, &size)
		return Size_T(size)
	}
	if l.Size_size_t == 8 {
		var size uint64
		binary.Read(l.Data, binary.LittleEndian, &size)
		return Size_T(size)
	}
	return 0
}
