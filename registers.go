package mips

import (
	"errors"
	"fmt"
)

var (
	InvalidSet        = errors.New("Invalid Register Set")
	RegisterLocked    = errors.New("Register Locked")
	RegisterNotLocked = errors.New("Register Not Locked")
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

type Registers struct {
	values []Word
	locks  []int
}

func NewRegisters() *Registers {
	return &Registers{
		values: make([]Word, numRegisters),
		locks:  make([]int, numRegisters),
	}
}

func (r Registers) String() string {
	result := ""
	for i := 0; i < numRegisters; i++ {
		if r.values[i] != 0 {
			result += fmt.Sprintf("\n%s = %d", Register(i), r.values[i])
		}
	}
	return result
}

func (r *Registers) Acquire(register Register) {
	if register == R0 {
		return
	}
	r.locks[register] += 1
}

func (r *Registers) Locked(register Register) bool {
	if register == R0 {
		return false
	}
	return r.locks[register] > 0
}

func (r *Registers) Release(register Register) {
	if register == R0 {
		return
	}
	r.locks[register] -= 1
	if r.locks[register] < 0 {
		panic("over-released")
	}
}

func (r *Registers) Set(register Register, value Word) error {
	if register == 0 {
		return InvalidSet
	} else {
		r.values[register] = value
	}
	return nil
}

func (r *Registers) Get(register Register) Word {
	if register == 0 {
		return 0
	}
	return r.values[register]
}

func (r Register) String() string {
	return fmt.Sprintf("R%d", uint64(r))
}
