package mips

import (
	"fmt"
	"testing"
)

func TestRunningEmptyCPU(t *testing.T) {
	m := NewCPU()
	if m == nil {
		t.Error("cpu == nil")
	}
	err := m.Run()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(m.Cycle)
}

func TestRunningSimpleCPU(t *testing.T) {
	m, err := ParseCPUString(`REGISTERS
MEMORY
CODE
Loop: LD    R2,    0(R1) 
      DADD  R4,    R2,    R3
      SD    0(R1), R4 
`)
	if m == nil {
		t.Error("cpu == nil")
	}
	if err != nil {
		t.Error(err)
	}
	err = m.Run()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(m.Cycle)
}
