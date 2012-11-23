// Simple simulator of a subset of the MIPS instruction set to show pipelining
package mips

import (
	"errors"
	"fmt"
)

type Code []*Instruction

var CPUFinished = errors.New("CPU Finished.")

type CPU struct {
	Registers          Registers
	Cycle              int
	Ram                Memory
	Code               Code
	InstructionPointer int
	Labels             map[Label]int // label to Code index mapping
	//Pipeline *InstructionPipeline
	Stages []PipelineStage
}

func NewCPU() *CPU {
	cpu := &CPU{
		Code:   make([]*Instruction, 0),
		Labels: make(map[Label]int),
		Stages: []PipelineStage{
			new(IF1),
			new(IF2),
			new(EX),
			new(WB),
		},
	}
	for i, stage := range cpu.Stages {
		stage.Initialize(cpu)
		if i > 0 {
			stage.SetPrev(cpu.Stages[i-1])
			cpu.Stages[i-1].SetNext(stage)
		}
	}
	return cpu
}

func (cpu *CPU) Run() (err error) {
	for err == nil {
		err = cpu.Step()
	}
	if err == CPUFinished {
		return nil
	}
	return err
}

func (cpu *CPU) Step() error {
	fmt.Println("#################### CYCLE", cpu.Cycle, "####################")

	// Move instructions to next stage of pipeline
	for i := len(cpu.Stages) - 1; i >= 0; i-- {
		cpu.Stages[i].TransferInstruction()
	}

	for i, stage := range cpu.Stages {
		fmt.Println("#################### CYCLE", cpu.Cycle, "stage", i, stage)
		err := stage.Step()
		if err != nil {
			return err
		}
	}

	cpu.Cycle += 1
	return nil
}

func (cpu *CPU) GetNextStage(s1 PipelineStage) PipelineStage {
	s1Index := -1
	for i, stage := range cpu.Stages {
		if s1Index != -1 {
			return stage
		}
		if stage == s1 {
			s1Index = i
		}
	}
	return cpu.Stages[0]
}

func (cpu *CPU) GetPreviousStage(s1 PipelineStage) PipelineStage {
	s1Index := -1
	for i, stage := range cpu.Stages {
		fmt.Println("eq?", stage, s1)
		if stage == s1 {
			s1Index = i
		}
	}
	fmt.Println("GPS", s1Index)
	if s1Index > 0 {
		return cpu.Stages[s1Index]
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

func (cpu *CPU) String() string {
	return fmt.Sprintf("REGISTERS: %s\nMEMORY: %s", cpu.Registers, cpu.Ram)
}
