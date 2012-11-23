// Simple simulator of a subset of the MIPS instruction set to show pipelining
package mips

import (
	"errors"
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
	numRegisters
)

type Registers [numRegisters]Word

type Memory [memorySize]Word

type Code []*Instruction

type Operation struct {
	op string
}

type OperandType int

const (
	operandTypeInvalid OperandType = iota
	operandTypeImmediate
	operandTypeOffset
	operandTypeNormal
	operandTypeLabel
)

type Operand struct {
	text     string
	Register Register
	Offset   int
	Type     OperandType
}

type Label string

type Instruction struct {
	Label       Label
	text        string
	Operation   Operation
	Destination Operand
	OperandA    Operand
	OperandB    Operand
}

type InstructionInPipeline struct {
	*Instruction
	Stage PipelineStage
}

type PipelineStage interface {
	Initialize(cpu *CPU)
	String() string
	Step() error
	GetInstruction() *InstructionInPipeline
	TransferInstruction()
	Next() PipelineStage
	Prev() PipelineStage
	SetNext(PipelineStage) 
	SetPrev(PipelineStage) 
}

type baseStage struct {
	Instruction *InstructionInPipeline
	cpu *CPU
	next PipelineStage
	prev PipelineStage
}

func (s *baseStage) Initialize(cpu *CPU) {
	s.cpu = cpu
}

func (s baseStage) String() string {
	return "unknown"
}

func (s *baseStage) Step() error {
	fmt.Println("baseStage executed, halting.")
	return CPUFinished
}

func (s *baseStage) Prev() PipelineStage {
	return s.prev

}
func (s *baseStage) Next() PipelineStage {
	return s.next
}

func (s *baseStage) SetPrev(p PipelineStage) {
	s.prev = p

}
func (s *baseStage) SetNext(p PipelineStage) {
	s.next = p
}

func (s *baseStage) GetInstruction() *InstructionInPipeline {
	return s.Instruction
}

func (s *baseStage) TransferInstruction() {
	prev := s.Prev()
	if prev == nil {
		return
	}
	fmt.Println("Transferrring instruction from", prev, ":", s.Prev().GetInstruction())
	s.Instruction = s.Prev().GetInstruction()
	if s.Instruction != nil {
		s.Instruction.Stage = s
	}
}

/////////////////////////////////////////////////////////////////////////////
// IF1
/////////////////////////////////////////////////////////////////////////////

type IF1 struct {
	baseStage
}

func (s *IF1) Step() error {
	fmt.Println("IF1 executed.")
	
	fmt.Println(s.cpu)
	if s.cpu.InstructionPointer == len(s.cpu.Code) {
		fmt.Println("No more instructions")
		s.Instruction = nil
	} else {
		s.Instruction = &InstructionInPipeline{
			s.cpu.Code[s.cpu.InstructionPointer],
			s,
		}
		fmt.Println("Loaded new instruction:", s.Instruction)
		s.cpu.InstructionPointer += 1
	}
	return nil
}

func (s IF1) String() string {
	return "IF1"
}


/////////////////////////////////////////////////////////////////////////////
// 
/////////////////////////////////////////////////////////////////////////////

type IF2 struct {
	baseStage
}

func (s *IF2) Step() error {
	fmt.Println("IF2 executed.")
	s.TransferInstruction()
	return nil
}

func (s IF2) String() string {
	return "IF2"
}

/////////////////////////////////////////////////////////////////////////////
// 
/////////////////////////////////////////////////////////////////////////////

type IF3 struct {
	baseStage
}

/////////////////////////////////////////////////////////////////////////////
// 
/////////////////////////////////////////////////////////////////////////////

type ID struct {
	baseStage
}

/////////////////////////////////////////////////////////////////////////////
// 
/////////////////////////////////////////////////////////////////////////////

type EX struct {
	baseStage
}

func (s EX) String() string {
	return "EX"
}

func (s *EX) Step() error {
	fmt.Println("EX executed.")
	s.TransferInstruction()
	return nil
}

/////////////////////////////////////////////////////////////////////////////
// 
/////////////////////////////////////////////////////////////////////////////

type MEM1 struct {
	baseStage
}

/////////////////////////////////////////////////////////////////////////////
// 
/////////////////////////////////////////////////////////////////////////////

type MEM2 struct {
	baseStage
}

/////////////////////////////////////////////////////////////////////////////
// 
/////////////////////////////////////////////////////////////////////////////

type MEM3 struct {
	baseStage
}

/////////////////////////////////////////////////////////////////////////////
// 
/////////////////////////////////////////////////////////////////////////////

type WB struct {
	baseStage
}

func (s WB) String() string {
	return "WB"
}

func (s *WB) Step() error {
	fmt.Println("WB executed.")
	s.TransferInstruction()
	if s.Instruction == nil {
		return CPUFinished
	}
	return nil
}

//type InstructionPipeline struct {
//	stages	[]PipelineStage
//}

var CPUFinished = errors.New("CPU Finished.")

type CPU struct {
	State              CPUState
	Cycle              int
	Ram                Memory
	Code               Code
	InstructionPointer int
	Labels             map[Label]int // label to Code index mapping
	//Pipeline *InstructionPipeline
	Stages       []PipelineStage
}

type CPUState struct {
	Registers
}

func NewCPU() *CPU {
	m := &CPU{
		Code:   make([]*Instruction, 0),
		Labels: make(map[Label]int),
		Stages: []PipelineStage{
			new(IF1),
			new(IF2),
			new(EX),
			new(WB),
		},
	}
	for i,stage := range m.Stages {
		stage.Initialize(m)
		if i > 0 {
			stage.SetPrev(m.Stages[i-1])
			m.Stages[i-1].SetNext(stage)
		}
	}
	return m
}

func (m *CPU) Run() (err error) {
	for err == nil {
		err = m.Step()
	}
	if err == CPUFinished {
		return nil
	}
	return err
}

func (m *CPU) Step() error {
	fmt.Println("#################### CYCLE", m.Cycle, "####################")

	// Move instructions to next stage of pipeline
	for i := len(m.Stages)-1; i >= 0; i-- {
		m.Stages[i].TransferInstruction()
	}
	
	for i, stage := range m.Stages {
		fmt.Println("#################### CYCLE", m.Cycle, "stage", i, stage)
		err := stage.Step()
		if err != nil {
			return err
		}
	}

	m.Cycle += 1
	return nil
}

func (m *CPU) GetNextStage(s1 PipelineStage) PipelineStage {
	s1Index := -1
	for i,stage := range m.Stages {
		if s1Index != -1 {
			return stage
		}
		if stage == s1 {
			s1Index = i
		}
	}
	return m.Stages[0]
}

func (m *CPU) GetPreviousStage(s1 PipelineStage) PipelineStage {
	s1Index := -1
	for i,stage := range m.Stages {
		fmt.Println("eq?", stage, s1)
		if stage == s1 {
			s1Index = i
		}
	}
	fmt.Println("GPS", s1Index)
	if s1Index > 0 {
		return m.Stages[s1Index]
	}
	return nil
}


func (w Word) String() string {
	return fmt.Sprintf("%#x", uint64(w))
}

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

func (op Operation) String() string {
	return op.op
}

func (op Operand) String() string {
	switch op.Type {
	case operandTypeImmediate:
		return fmt.Sprintf("#%d", op.Offset)
	case operandTypeNormal:
		return fmt.Sprintf("%s", op.Register)
	case operandTypeOffset:
		return fmt.Sprintf("%d(%s)", op.Offset, op.Register)
	case operandTypeLabel:
		return op.text
	}
	return "@unknown@"
}

func (i Instruction) String() string {
	result := fmt.Sprintf("%s %s %s", i.Operation, i.Destination, i.OperandA)
	if i.OperandB.Type != operandTypeInvalid {
		result += fmt.Sprintf(" %s", i.OperandB)
	}
	//result := fmt.Sprintf("%s %s", i.Operation, i.Destination)
	if i.Label != "" {
		result += fmt.Sprintf(" (label: %s)", i.Label)
	}
	return result
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

func (m *CPU) String() string {
	return fmt.Sprintf("REGISTERS: %s\nMEMORY: %s", m.State.Registers, m.Ram)
}
