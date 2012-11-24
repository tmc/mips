// Simple simulator of a subset of the MIPS instruction set to show pipelining
package mips

import (
	"errors"
	"fmt"
)

// Branch prediction modes
type CPUBranchMode int

const (
	branchModeFlush = iota // no prediction, flush pipeline
	branchModePredictTaken
	branchModePredictNotTaken
)

type InstructionCache []Instruction

var CPUFinished = errors.New("CPU Finished.")

type CPU struct {
	Registers          *Registers
	BranchMode         CPUBranchMode
	ForwardingEnabled  bool
	Cycle              int
	Ram                Memory
	InstructionCache   InstructionCache
	InstructionPointer int
	Instructions       []*ExecutedInstruction
	Labels             map[Label]int // label to Code index mapping
	Pipeline           Pipeline
}

func NewCPU() *CPU {
	cpu := &CPU{
		InstructionCache: make([]Instruction, 0),
		Labels:           make(map[Label]int),
		Registers:        NewRegisters(),
	}
	pipeline, err := NewPipeline(cpu,
		new(IF1),
		new(IF2),
		new(IF3),
		new(ID),
		new(EX),
		new(MEM1),
		new(MEM2),
		new(MEM3),
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
	fmt.Println("#################### RUN ##########################")
	for err == nil {
		err = cpu.Step()

		fmt.Print(cpu.RenderState())
	}
	if err == CPUFinished {
		return nil
	}
	return err
}

func (cpu *CPU) Step() error {

	cpu.Cycle += 1

	//fmt.Println("#################### CYCLE", cpu.Cycle, "####################")

	// Move instructions to next stage of pipeline (unless a stage is stalled)
	if err := cpu.Pipeline.TransferInstructions(); err != nil {
		return err
	}

	if err := cpu.Pipeline.Execute(); err != nil {
		return err
	}

	if cpu.Pipeline.Empty() && cpu.InstructionCacheEmpty() {
		//fmt.Println("pipeline empty", len(cpu.Pipeline.ActiveInstructions()))
		return CPUFinished
	}

	return nil
}

func (cpu *CPU) RenderState() string {
	result := ""
	// spacing helper
	print := spacingHelper(12)

	result += print("c#%d", cpu.Cycle)
	for _, inst := range cpu.Instructions {
		inPipeline := false
		for _, iip := range cpu.Pipeline.ActiveInstructions() {
			if iip == inst {
				inPipeline = true
				stalled := iip.Stage.Stalled()
				if stalled {
					//print("(s) %s", iip.OpCode())
					result += print("(s)")
				} else {
					//print("%s %s", iip.Stage, iip.OpCode())
					if iip.Stage.String() == "IF1" {
						result += print("%s:%s", iip.Stage, iip.Instruction.OpCode())
					} else {
						result += print("%s", iip.Stage)
					}

				}
			}
		}
		if inPipeline == false {
			result += print("")
		}
	}
	result += "\n"
	return result
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
