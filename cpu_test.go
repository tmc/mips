package mips

import (
	"fmt"
	"testing"
	"strings"
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
`, "branching_simple": `REGISTERS
R1 2
MEMORY
0 7
CODE
Start: DADDI R1, R1, #-1
       BNEZ  R1, Start
       LD    R4, #0
`, "basic": `REGISTERS
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
`, "provided0": `REGISTERS
MEMORY 
CODE
      LD    R2,     0(R1)
      DADD  R4,     R2,    R3
      SD    0(R1),  R4
      BNEZ  R4,     NEXT
NEXT: DADD  R1,     R1,    R3
      DADDI R2,     R1,    #8
`, "provided1": `REGISTERS
R1  16
R3  42
MEMORY 
16  60
8   40
CODE
Loop: LD    R2,     0(R1)
      DADD  R4,     R2,     R3
      SD    0(R1),  R4
      DADDI R1,     R1,     #-8
      BNEZ  R1,     Loop
      DADD  R3,     R2,     R4	
`, "provided2": `REGISTERS
R1 16
R2 16
R3 20
R4 2
R5 8
R7 8
MEMORY 
16  8
8  12
CODE
Loop: LD    R2,    0(R1) 
      DADD  R4,    R2,    R3
      SD    0(R1), R4 
      DADDI R1,    R1,    #-8
      BNEZ  R1,   Loop
      DADDI R1,    R1,    #-8
      BNEZ  R1,    Next
      DADD  R3,    R4,    R5
Next: LD    R6,    0(R5) 
      DADD  R4,    R2,    R3
      SD    0(R5), R4
      DADDI R1,    R1,    #-8
`,
}

func xTestRunningEmptyCPU(t *testing.T) {
	cpu := NewCPU()
	if cpu == nil {
		t.Error("cpu == nil")
	}
	err := cpu.Run(100)
	if err != nil {
		t.Error(err)
	}
}

func xTestRAWHazard(t *testing.T) {
	cpu, err := ParseCPUString(CPU_TESTS["raw_hazard"])
	if cpu == nil {
		t.Error("cpu == nil")
	}
	if err != nil {
		t.Error(err)
	}
	//cpu.ForwardingEnabled = true
	err = cpu.Run(15)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(cpu.RenderTiming())
	fmt.Println(cpu)

	if cpu.Ram[1] != 10 {
		t.Fatal(cpu.Ram[1], "!=", 10)
	}
}

func xTestBranching(t *testing.T) {
	cpu, err := ParseCPUString(CPU_TESTS["branching_simple"])
	if cpu == nil {
		t.Error("cpu == nil")
	}
	if err != nil {
		t.Error(err)
	}
	cpu.ForwardingEnabled = true
	cpu.BranchMode = branchModePredictNotTaken
	cpu.BranchMode = branchModePredictTaken
	err = cpu.Run(35)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(cpu.RenderTiming())

	if cpu.Registers.Get(R1) != 0 {
		t.Fatal(cpu.Registers.Get(R1), "!=", 0)
	}
}

func xTestRunningBasicProgram(t *testing.T) {
	cpu, err := ParseCPUString(CPU_TESTS["basic"])
	if cpu == nil {
		t.Error("cpu == nil")
	}
	if err != nil {
		t.Error(err)
	}
	err = cpu.Run(100)
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

func TestSameOutputRegardlessOfFlags(t *testing.T) {
	for testName, test := range CPU_TESTS {
		if strings.HasPrefix(testName, "provided0") == false {
			continue
		}
		cpu, _ := ParseCPUString(test)
		if err := cpu.Run(100); err != nil{

		fmt.Println("######### PROGRAM ", testName)
		fmt.Println("######### CODE")
		fmt.Println(test)
		fmt.Println("######### Result")
		fmt.Println(cpu.String())
		fmt.Println("######### Timing")
		fmt.Println(cpu.RenderTiming())

			t.Fatal(testName, "nf", err)
		}

		a := cpu.String()

		cpu, _ = ParseCPUString(test)
		cpu.ForwardingEnabled = true
		cpu.BranchMode = branchModePredictTaken
		if err := cpu.Run(100); err != nil{
			t.Fatal(testName, "f, bt", err)
		}

		b := cpu.String()
		timing := cpu.RenderTiming()

		cpu, _ = ParseCPUString(test)
		cpu.ForwardingEnabled = true
		cpu.BranchMode = branchModePredictNotTaken
		if err := cpu.Run(100); err != nil{
			t.Fatal(testName, "f, bnt", err)
		}

		c := cpu.String()

		if a != b {
			t.Error("'%s' != '%s'", a, b)
		}
		if b != c {
			t.Error("'%s' != '%s'", b, c)
		}
		fmt.Println("######### PROGRAM ", testName)
		fmt.Println("######### CODE")
		fmt.Println(test)
		fmt.Println("######### Result")
		fmt.Println(a)
		fmt.Println("######### Timing")
		fmt.Println(timing)
	}
}
