// Simple simulator of a subset of the MIPS instruction set to show pipelining
package mips

import (
	"errors"
	"fmt"
)

type InstructionInPipeline struct {
	*Instruction
	Stage PipelineStage
}

type Pipeline []PipelineStage

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

func NewPipeline(cpu *CPU, stages ...PipelineStage) (Pipeline, error) {
	pipeline := make([]PipelineStage, 0)

	for i, stage := range stages {
		stage.Initialize(cpu)
		pipeline = append(pipeline, stage)
		if i > 0 {
			stage.SetPrev(pipeline[i-1])
			pipeline[i-1].SetNext(stage)
		}

	}
	if len(pipeline) == 0 {
		return nil, errors.New("Must have at least one stage")
	}
	return pipeline, nil
}

func (p Pipeline) GetNextStage(s PipelineStage) PipelineStage {
	index := -1
	for i, stage := range p {
		if index != -1 {
			return stage
		}
		if stage == s {
			index = i
		}
	}
	return p[0]
}

func (p Pipeline) GetPreviousStage(s PipelineStage) PipelineStage {
	index := -1
	for i, stage := range p {
		fmt.Println("eq?", stage, s)
		if stage == s {
			index = i
		}
	}
	if index > 0 {
		return p[index]
	}
	return nil
}

func (p Pipeline) Empty() bool {
	allEmpty := true
	for _, stage := range p {
		if stage.GetInstruction() != nil {
			allEmpty = false
		}
	}
	return allEmpty
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
	fmt.Println("IF1 executing")

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

		s.Instruction.IF1()

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
	if s.Instruction != nil {
		s.Instruction.IF2()
	}
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
	if s.Instruction != nil {
		s.Instruction.EX()
	}
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
	if s.Instruction != nil {
		s.Instruction.EX()
	}
	return nil
}
