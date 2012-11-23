// Simple simulator of a subset of the MIPS instruction set to show pipelining
package mips

import (
	"errors"
	"fmt"
)

type Code []*Instruction

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


func (m *CPU) String() string {
	return fmt.Sprintf("REGISTERS: %s\nMEMORY: %s", m.State.Registers, m.Ram)
}
