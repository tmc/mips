// Simple simulator of a subset of the MIPS instruction set to show pipelining
package mips

import (
	"fmt"
)

type Word uint64
type Register int
type Address int

const memorySize = 992 // Size of memory in words

const (
	None = -1 + iota
	R0
	R1
	R2
	R3
	R4
	R5
	R6
	R7
	R8
	R9
	R10
	R11
	R12
	R13
	R14
	R15
	R16
	R17
	R18
	R19
	R20
	R21
	R22
	R23
	R24
	R25
	R26
	R27
	R28
	R29
	R30
	R31
	PC
	numRegisters
)

type Registers [numRegisters]Word

type Memory [memorySize]Word

type Code []Instruction

type Operation struct {
	op string
}

type OperandType int

const (
	operandTypeImmediate OperandType = iota
	operandTypeOffset
	operandTypeLabel
)

type Operand struct {
	text     string
	Register Register
	Offset   int
	Type     OperandType
}

type Instruction struct {
	label       string
	text        string
	Operation   Operation
	Destination Operand
	OperandA    Operand
	OperandB    Operand
}

type Label string

type Machine struct {
	State  MachineState
	Ram    Memory
	Code   Code
	Labels map[Label]Word
}

type MachineState struct {
	Registers
}

func NewMachine() *Machine {
	return &Machine{}
}

func (w Word) String() string {
	return fmt.Sprintf("%#x", uint64(w))
}

func (r Registers) String() string {
	result := ""
	for i := 0; i < numRegisters; i++ {
		if r[i] != 0 {
			result += fmt.Sprintf("\nR%d = %d", i, r[i])
		}
	}
	return result
}

func (r *Registers) Set(register Register, value Word) Word {
	if register == 0 {
		return 0
	} else {
		r[register] = value
	}
	return value
}

func (r *Registers) Get(register Register) Word {
	if register == 0 {
		return 0
	}
	return r[register]
}

func (lhs Code) Equals(rhs Code) bool {
	if len(lhs) != len(rhs) {
		return false
	}
	for i, l := range lhs {
		if l != rhs[i] {
			return false
		}
	}
	return true
}

func (r Memory) String() string {
	result := ""
	for i := 0; i < memorySize; i++ {
		if r[i] != 0 {
			result += fmt.Sprintf("\n%#x = %d", i, r[i])
		}
	}
	return result
}

func (m *Machine) String() string {
	return fmt.Sprintf("REGISTERS: %s\nMEMORY: %s", m.State.Registers, m.Ram)
}
