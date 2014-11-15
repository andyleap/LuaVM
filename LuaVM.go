// LuaVM project LuaVM.go
package LuaVM

import (
	"encoding/binary"
	"errors"
	"fmt"
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
	CLOSURE
)

var (
	BADSIGNATURE error = errors.New("Bad signature")
	BADVERSION   error = errors.New("Bad version")
	BADENCODING  error = errors.New("Bad encoding")
)

type Size_T uint32
type Integer uint32
type Instruction uint32
type Number float64

type Value struct {
	Type  ValueType
	Value interface{}
}

type Stackframe struct {
	Regs     []*Value
	Upvalues []*Value
}

type Instr struct {
	Opcode uint8
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

type Closure struct {
	Stack    *Stackframe
	Function *FunctionPrototype
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

func ReadLuaC(data io.Reader) (*Closure, error) {
	err := checkHeader(data)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func checkHeader(data io.Reader) error {
	var header header
	binary.Read(data, binary.LittleEndian, &header)
	if header.Signature != 0x61754C1B {
		return BADSIGNATURE
	}
	if header.Version != 0x51 {
		return BADVERSION
	}
	if header.Format != 0 ||
		header.Endianness != 1 ||
		header.Size_int != 4 ||
		header.Size_size_t != 4 ||
		header.Size_instruction != 4 ||
		header.Size_number != 8 ||
		header.Integral != 0 {
		return BADENCODING
	}
	readFunctionBlock(data)
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

func readFunctionBlock(data io.Reader) (*FunctionPrototype, error) {
	readString(data)

	var block functionBlock
	binary.Read(data, binary.LittleEndian, &block)

	Prototype := &FunctionPrototype{
		Upvalues:     block.Upvalues,
		Parameters:   block.Parameters,
		IsVararg:     block.IsVararg,
		MaxStackSize: block.MaxStackSize,
	}
	Prototype.Instructions = readInstructionList(data)
	Prototype.Constants = readConstantList(data)

	Prototype.Functions = readFunctionList(data)

	fmt.Println(Prototype.Instructions)

	readSourceLinePositionList(data)
	readLocalList(data)
	readUpvalueList(data)

	return Prototype, nil
}

func readUpvalueList(data io.Reader) {
	var size Integer
	binary.Read(data, binary.LittleEndian, &size)
	for l1 := Integer(0); l1 < size; l1++ {
		readString(data)
	}
}

func readLocalList(data io.Reader) {
	var size Integer
	var local Integer
	binary.Read(data, binary.LittleEndian, &size)
	for l1 := Integer(0); l1 < size; l1++ {
		readString(data)
		binary.Read(data, binary.LittleEndian, &local)
		binary.Read(data, binary.LittleEndian, &local)
	}
}

func readSourceLinePositionList(data io.Reader) {
	var size Integer
	var line Integer
	binary.Read(data, binary.LittleEndian, &size)
	for l1 := Integer(0); l1 < size; l1++ {
		binary.Read(data, binary.LittleEndian, &line)
	}
}

func readFunctionList(data io.Reader) []*FunctionPrototype {
	var size Integer
	binary.Read(data, binary.LittleEndian, &size)
	functions := make([]*FunctionPrototype, size)
	for l1 := Integer(0); l1 < size; l1++ {
		function, _ := readFunctionBlock(data)
		functions[l1] = function
	}
	return functions
}

func readInstruction(data io.Reader) Instr {
	var instruction Instruction
	binary.Read(data, binary.LittleEndian, &instruction)

	ret := Instr{}

	opcode := uint8(0x0000003F & instruction)

	fmt.Println(opcode)
	ret.Opcode = opcode

	switch opcode {
	case 0, 3, 2, 4, 8, 6, 9, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 28, 30, 29, 37, 11, 23, 24, 25, 26, 27, 33, 10, 34, 35: //iABC
		ret.A = uint8((instruction & 0x00003FC0) >> 6)
		ret.B = int32((instruction & 0x007FC000) >> 14)
		ret.C = uint16((instruction & 0xFF800000) >> 23)
	case 1, 5, 7, 36: //iABx
		ret.A = uint8((instruction & 0x00003FC0) >> 6)
		ret.B = int32((instruction & 0xFFFFC000) >> 14)
	case 22, 31, 32: //iAsBx
		ret.A = uint8((instruction & 0x00003FC0) >> 6)
		ret.B = int32(((instruction & 0xFFFFC000) >> 14)) - 131071
	}

	return ret
}

func readInstructionList(data io.Reader) []Instr {
	var size Integer
	binary.Read(data, binary.LittleEndian, &size)
	instructions := make([]Instr, size)
	for l1 := Integer(0); l1 < size; l1++ {
		instructions[l1] = readInstruction(data)
	}
	return instructions
}

func readConstantList(data io.Reader) []Value {
	var size Integer
	var valuetype ValueType
	var integer Integer
	var number Number
	binary.Read(data, binary.LittleEndian, &size)
	constants := make([]Value, size)
	for l1 := Integer(0); l1 < size; l1++ {
		binary.Read(data, binary.LittleEndian, &valuetype)
		constants[l1].Type = valuetype
		switch {
		case valuetype == NIL:
		case valuetype == BOOLEAN:
			binary.Read(data, binary.LittleEndian, &integer)
			constants[l1].Value = integer
		case valuetype == NUMBER:
			binary.Read(data, binary.LittleEndian, &number)
			constants[l1].Value = number
		case valuetype == STRING:
			constants[l1].Value = readString(data)
		}
	}
	return constants
}

func readString(data io.Reader) string {
	var size Size_T
	binary.Read(data, binary.LittleEndian, &size)
	str := make([]byte, size)
	data.Read(str)
	return string(str)
}
