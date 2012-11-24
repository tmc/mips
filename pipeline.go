// Simple simulator of a subset of the MIPS instruction set to show pipelining
package mips

import (
	"errors"
	"fmt"
)

type InstructionInPipeline struct {
	Instruction
	Stage PipelineStage
}

type Pipeline []PipelineStage

type PipelineStage interface {
	Initialize(cpu *CPU)
	String() string
	Step() error
	SetInstruction(instruction *InstructionInPipeline) 
	GetInstruction() *InstructionInPipeline
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

func (p Pipeline) ActiveInstructions() []*InstructionInPipeline {
	result := make([]*InstructionInPipeline, 0)
	
	for _, stage := range p {
		iip := stage.GetInstruction()
		if iip != nil {
			result = append(result, iip)
		}
	}
	
	return result
}

type baseStage struct {
	instruction *InstructionInPipeline
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
	return nil
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
	return s.instruction
}

func (s *baseStage) SetInstruction(instruction *InstructionInPipeline) {
	s.instruction = instruction
}

func (p *Pipeline) TransferInstruction(toStage PipelineStage) {
	fromStage := toStage.Prev()
	if fromStage == nil {
		return
	}
	inst := fromStage.GetInstruction()
	toStage.SetInstruction(inst)
	if inst != nil {
		//fmt.Println("Transferrring instruction from", fromStage, "to", toStage, ":", toStage.GetInstruction())
		inst.Stage = toStage
	}
}

/////////////////////////////////////////////////////////////////////////////
// IF1
/////////////////////////////////////////////////////////////////////////////

type IF1 struct {
	baseStage
}

func (s *IF1) Step() error {
	if s.cpu.InstructionCacheEmpty() {
		//fmt.Println("No more instructions")
		s.instruction = nil
	} else {
		s.instruction = &InstructionInPipeline{
			s.cpu.InstructionCache[s.cpu.InstructionPointer],
			s,
		}
		fmt.Println("Issue:", s.instruction)
		s.cpu.InstructionPointer += 1

		s.instruction.IF1()

	}
	return nil
}

func (s IF1) String() string {
	return "IF1"
}

/////////////////////////////////////////////////////////////////////////////
// IF2
/////////////////////////////////////////////////////////////////////////////

type IF2 struct {
	baseStage
}

func (s *IF2) Step() error {
	if s.instruction != nil {
		s.instruction.IF2()
	}
	return nil
}

func (s IF2) String() string {
	return "IF2"
}

/////////////////////////////////////////////////////////////////////////////
// IF3
/////////////////////////////////////////////////////////////////////////////

type IF3 struct {
	baseStage
}

func (s IF3) String() string {
	return "IF3"
}

/////////////////////////////////////////////////////////////////////////////
// ID
/////////////////////////////////////////////////////////////////////////////

type ID struct {
	baseStage
}

func (s ID) String() string {
	return "ID"
}

/////////////////////////////////////////////////////////////////////////////
// EX
/////////////////////////////////////////////////////////////////////////////

type EX struct {
	baseStage
}

func (s EX) String() string {
	return "EX"
}

func (s *EX) Step() error {
	if s.instruction != nil {
		s.instruction.EX()
	}
	return nil
}

/////////////////////////////////////////////////////////////////////////////
// MEM1
/////////////////////////////////////////////////////////////////////////////

type MEM1 struct {
	baseStage
}

func (s MEM1) String() string {
	return "MEM1"
}

/////////////////////////////////////////////////////////////////////////////
// MEM2
/////////////////////////////////////////////////////////////////////////////

type MEM2 struct {
	baseStage
}

func (s MEM2) String() string {
	return "MEM2"
}

/////////////////////////////////////////////////////////////////////////////
// MEM3
/////////////////////////////////////////////////////////////////////////////

type MEM3 struct {
	baseStage
}

func (s MEM3) String() string {
	return "MEM3"
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
	if s.instruction != nil {
		s.instruction.WB()
	}
	return nil
}
