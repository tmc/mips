// Simple simulator of a subset of the MIPS instruction set to show pipelining
package mips

import (
	"fmt"
)

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
	cpu         *CPU
	next        PipelineStage
	prev        PipelineStage
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
