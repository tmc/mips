// Simple simulator of a subset of the MIPS instruction set to show pipelining
package mips

import (
	"errors"
	"fmt"
)

type CPUMode int

const (
	ModeNoPipeline = iota
	ModeNoForwardingOrBypassing
	ModeBranchPredictionTaken
	ModeBranchPredictionNotTaken
)

type InstructionCache []Instruction

var CPUFinished = errors.New("CPU Finished.")

type CPU struct {
	Registers          Registers
	Mode               CPUMode
	Cycle              int
	Ram                Memory
	InstructionCache   InstructionCache
	InstructionPointer int
	Labels             map[Label]int // label to Code index mapping
	Pipeline           Pipeline
}

func NewCPU() *CPU {
	cpu := &CPU{
		InstructionCache: make([]Instruction, 0),
		Labels:           make(map[Label]int),
	}
	pipeline, err := NewPipeline(cpu,
		new(IF1),
		//new(IF2),
		//new(IF3),
		new(ID),
		new(EX),
		new(MEM1),
		//new(MEM2),
		//new(MEM3),
		new(WB),
	)
	if err != nil {
		panic(err)
	}
	cpu.Pipeline = pipeline
	return cpu
}

func (cpu *CPU) Run() (err error) {
	fmt.Println(len(cpu.InstructionCache), "instructions")

	for err == nil {
		err = cpu.Step()

		cpu.PrintTiming(cpu.Cycle == 1)
	}
	if err == CPUFinished {
		return nil
	}
	return err
}

func (cpu *CPU) Step() error {
	cpu.Cycle += 1

	//fmt.Println("#################### CYCLE", cpu.Cycle, "####################")

	// Move instructions to next stage of pipeline (talk pipeline stages backwards)
	for i := len(cpu.Pipeline) - 1; i >= 0; i-- {
		cpu.Pipeline.TransferInstruction(cpu.Pipeline[i])
	}

	for _, stage := range cpu.Pipeline {
		//fmt.Println("################ stage", stage)
		err := stage.Step()
		if err != nil {
			return errors.New(fmt.Sprintf("Error while executing %s of %s: %s", stage, stage.GetInstruction(), err))
		}
	}

	if cpu.Pipeline.Empty() && cpu.InstructionCacheEmpty() {
		return CPUFinished
	}

	return nil
}

func (cpu *CPU) PrintTiming(printHeader bool) {
	if printHeader {
		fmt.Printf("%8s", "")
		for i, _ := range cpu.InstructionCache {
			fmt.Printf("I#%-6d", i+1)
		}
		fmt.Println("")
	}

	fmt.Printf("c#%-6d", cpu.Cycle)
	for _, inst := range cpu.InstructionCache {
		inPipeline := false
		for _, iip := range cpu.Pipeline.ActiveInstructions() {
			if iip.Instruction == inst {
				inPipeline = true
				fmt.Printf("%-6s", iip.Stage)
			}
		}
		if inPipeline == false {
			fmt.Printf("%-6s", "")
		}
		//fmt.Print("%6s", )
	}
	fmt.Println("")
}

func (cpu *CPU) InstructionCacheEmpty() bool {
	return cpu.InstructionPointer == len(cpu.InstructionCache)
}

func (lhs InstructionCache) Equals(rhs InstructionCache) bool {
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
