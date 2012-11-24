// Parses input files for simple mips pipeline simulator
package mips

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"strings"
	"testing"
)

func testFile(fn string) io.Reader {
	_, filename, _, _ := runtime.Caller(1)
	fh, err := os.Open(path.Join(path.Dir(filename), "test_data", fn))
	if err != nil {
		panic(err)
	}
	return fh
}

func testLiteral(test string) io.Reader {
	return strings.NewReader(test)
}

var PARSER_CPU_TESTS = []string{
	// Empty case
	`REGISTERS
MEMORY
CODE`,
	// Test setting of a register
	`REGISTERS
R1 4
MEMORY
CODE`,
	// Test setting registers and memory
	`REGISTERS
R1 42
R0 9
R7 13
MEMORY
42 31337
CODE`,
}

func testCPUsEqual(m1, m2 *CPU) bool {
	if m1 == nil && m2 == nil {
		return true
	}
	stateAndRamEqual := m2.Registers == m2.Registers && m1.Ram == m2.Ram
	return stateAndRamEqual && m1.InstructionCache.Equals(m2.InstructionCache)
}

func TestParsingSmallest(t *testing.T) {
	cpu, err := ParseCPU(testLiteral(PARSER_CPU_TESTS[0]))
	if err != nil {
		t.Error(err)
	}
	if testCPUsEqual(cpu, NewCPU()) == false {
		t.Error("Parsing of smallest doesn't match empty cpu")
	}
}

func TestParsingSimpleRegisterSet(t *testing.T) {
	cpu, err := ParseCPU(testLiteral(PARSER_CPU_TESTS[1]))
	if err != nil {
		t.Error(err)
	}
	if cpu == nil {
		t.Error("cpu == nil")
	}
	if cpu.Registers.Get(R1) != 4 {
		t.Fail()
	}
}

func TestParsingRegistersAndMemorySet(t *testing.T) {
	cpu, err := ParseCPU(testLiteral(PARSER_CPU_TESTS[2]))
	if err != nil {
		t.Error(err)
	}
	if cpu == nil {
		t.Error("cpu == nil")
	}
	if cpu.Registers.Get(R1) != 42 {
		t.Error("R1 != 42")
	}
	if cpu.Registers.Get(R7) != 13 {
		t.Error("R7 != 13")
	}
	if cpu.Registers.Get(R0) != 0 {
		t.Error("R0 != 0")
	}
	if cpu.Ram[42] != 31337 {
		t.Error("Mem[42] != 31337")
	}
}

func TestParsingProvided0(t *testing.T) {
	cpu, err := ParseCPU(testFile("input-0.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if testCPUsEqual(cpu, NewCPU()) == true {
		t.Error("Parsing provided example 0 produced empty cpu")
	}
}

func TestParsingProvided1(t *testing.T) {
	cpu, err := ParseCPU(testFile("input-1.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if testCPUsEqual(cpu, NewCPU()) == true {
		t.Error("Parsing provided example 1 produced empty cpu")
	}
}

func TestParsingProvided2(t *testing.T) {
	cpu, err := ParseCPU(testFile("input-2.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if testCPUsEqual(cpu, NewCPU()) == true {
		t.Error("Parsing provided example 2 produced empty cpu")
	}
}

var INSTRUCTION_TESTS = `Loop: LD    R2,    0(R1) 
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
`

func TestInstructionParsing(t *testing.T) {
	actual := ""
	expected := `LD R2 0(R1) (label: Loop)
DADD R4 R2 R3
SD 0(R1) R4
DADDI R1 R1 #-8
BNEZ R1 Loop
DADDI R1 R1 #-8
BNEZ R1 Next
DADD R3 R4 R5
LD R6 0(R5) (label: Next)
DADD R4 R2 R3
SD 0(R5) R4
DADDI R1 R1 #-8
`
	for _, line := range strings.Split(strings.TrimSpace(INSTRUCTION_TESTS), "\n") {
		i, err := ParseInstruction(strings.NewReader(line))
		if err != nil {
			t.Error(line, err)
		}
		actual += fmt.Sprintln(i)
	}
	if expected != actual {
		t.Errorf("Expected != Actual: '%s' != '%s'", expected, actual)
	}
}
