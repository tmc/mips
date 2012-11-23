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
	Pipeline           Pipeline
}

func NewCPU() *CPU {
	cpu := &CPU{
		Code:   make([]*Instruction, 0),
		Labels: make(map[Label]int),
	}
	pipeline, err := NewPipeline(cpu,
		new(IF1),
		new(IF2),
		new(EX),
		new(WB),
	)
	if err != nil {
		panic(err)
	}
	cpu.Pipeline = pipeline
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
	for i := len(cpu.Pipeline) - 1; i >= 0; i-- {
		cpu.Pipeline[i].TransferInstruction()
	}

	for i, stage := range cpu.Pipeline {
		fmt.Println("#################### CYCLE", cpu.Cycle, "stage", i, stage)
		err := stage.Step()
		if err != nil {
			return err
		}
	}

	if cpu.Pipeline.Empty() && cpu.InstructionPointer >= len(cpu.Code) {
		return CPUFinished
	}

	cpu.Cycle += 1
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
