package mips

import (
	"fmt"
	"testing"
)

var CPU_TESTS = map[string]string{
	"raw_hazard": `REGISTERS
R1 1
MEMORY
0 7
CODE
      LD    R2,    0(R0) 
      DADDI R3,    R2,    #3
      SD    0(R1), R3
`,
	"branching_simple": `REGISTERS
R1 3
MEMORY
0 7
CODE
Start: DADDI R1, R1, #-1
       BNEZ R1, Start
`,
	"basic": `REGISTERS
R1 2
R3 22
MEMORY
0 7
1 6
2 20
CODE
Loop: LD    R2,    0(R1) 
      DADD  R4,    R2,    R3
      SD    0(R1), R4
`}

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

func TestRAWHazard(t *testing.T) {
	cpu, err := ParseCPUString(CPU_TESTS["raw_hazard"])
	if cpu == nil {
		t.Error("cpu == nil")
	}
	if err != nil {
		t.Error(err)
	}
	cpu.ForwardingEnabled = true
	err = cpu.Run()
	if err != nil {
		t.Error(err)
	}

	if cpu.Ram[1] != 10 {
		t.Fatal(cpu.Ram[1], "!=", 10)
	}
}

func TestBranching(t *testing.T) {
	cpu, err := ParseCPUString(CPU_TESTS["branching_simple"])
	if cpu == nil {
		t.Error("cpu == nil")
	}
	if err != nil {
		t.Error(err)
	}
	err = cpu.Run()
	if err != nil {
		t.Error(err)
	}

	if cpu.Registers.Get(R1) != 0 {
		t.Fatal(cpu.Registers.Get(R1), "!=", 0)
	}
}

func TestRunningBasicProgram(t *testing.T) {
	cpu, err := ParseCPUString(CPU_TESTS["basic"])
	if cpu == nil {
		t.Error("cpu == nil")
	}
	if err != nil {
		t.Error(err)
	}
	err = cpu.Run()
	if err != nil {
		t.Error(err)
	}
	expected := `REGISTERS: 
R1 = 2
R2 = 20
R3 = 22
R4 = 42
MEMORY: 
0x0 = 7
0x1 = 6
0x2 = 42`
	if cpu.String() != expected {
		t.Error(fmt.Sprintf("'%s' != '%s'", cpu.String(), expected))
	}
}
