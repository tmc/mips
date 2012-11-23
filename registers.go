// Simple simulator of a subset of the MIPS instruction set to show pipelining
package mips

import (
	"fmt"
)

type Register int

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
	numRegisters
)

type Registers [numRegisters]Word


func (r Registers) String() string {
	result := ""
	for i := 0; i < numRegisters; i++ {
		if r[i] != 0 {
			result += fmt.Sprintf("\n%s = %d", i, r[i])
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

func (r Register) String() string {
	return fmt.Sprintf("R%d", uint64(r))
}
