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

var MACHINE_TESTS = []string{
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

func testMachinesEqual(m1, m2 *Machine) bool {
	if m1 == nil && m2 == nil {
		return true
	}
	stateAndRamEqual := m1.State == m2.State && m1.Ram == m2.Ram
	return stateAndRamEqual && m1.Code.Equals(m2.Code)
}

func TestParsingSmallest(t *testing.T) {
	m, err := ParseMachine(testLiteral(MACHINE_TESTS[0]))
	if err != nil {
		t.Error(err)
	}
	if testMachinesEqual(m, NewMachine()) == false {
		t.Error("Parsing of smallest doesn't match empty machine")
	}
}

func TestParsingSimpleRegisterSet(t *testing.T) {
	m, err := ParseMachine(testLiteral(MACHINE_TESTS[1]))
	if err != nil {
		t.Error(err)
	}
	if m == nil {
		t.Error("machine == nil")
	}
	if m.State.Registers[R1] != 4 {
		t.Fail()
	}
}

func TestParsingRegistersAndMemorySet(t *testing.T) {
	m, err := ParseMachine(testLiteral(MACHINE_TESTS[2]))
	if err != nil {
		t.Error(err)
	}
	if m == nil {
		t.Error("machine == nil")
	}
	if m.State.Registers[R1] != 42 {
		t.Error("R1 != 42")
	}
	if m.State.Registers[R7] != 13 {
		t.Error("R7 != 13")
	}
	if m.State.Registers[R0] != 0 {
		t.Error("R0 != 0")
	}
	if m.Ram[42] != 31337 {
		t.Error("Mem[42] != 31337")
	}
}

func TestParsingProvided0(t *testing.T) {
	m, err := ParseMachine(testFile("input-0.txt"))
	if err != nil {
		t.Error(err)
	}
	if testMachinesEqual(m, NewMachine()) == true {
		t.Error("Parsing provided example 0 produced empty machine")
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
	fmt.Println(None, R0, R1, R2)
	for idx, line := range strings.Split(strings.TrimSpace(INSTRUCTION_TESTS), "\n") {
		i, err := ParseInstruction(strings.NewReader(line))
		if err != nil {
			t.Error(line, err)
		}
		fmt.Println(idx, i, "\n")
	}
}
