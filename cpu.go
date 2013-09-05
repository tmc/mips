// Package mips is a simple simulator of a subset of the MIPS instruction set to show pipelining
package mips

import (
	"bytes"
	"errors"
	"fmt"
)

// Branch prediction modes
type BranchPolicy int

const (
	BranchPolicyFlush = iota // no prediction, flush pipeline
	BranchPolicyPredictTaken
	BranchPolicyPredictNotTaken
)

type InstructionCache []Instruction

var (
	CPUFinished          = errors.New("CPU Finished.")
	MaximumCyclesReached = errors.New("Maximum Cycles Reached.")
)

type CPU struct {
	Registers          *Registers
	BranchMode         BranchPolicy
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

func (cpu *CPU) Run(maximumCycles int) (err error) {
	//fmt.Println("#################### RUN ##########################", len(cpu.InstructionCache), "instructions")
	for err == nil {
		err = cpu.Step()

		//Render state
		//fmt.Print(cpu.RenderState())

		if maximumCycles > 0 && cpu.Cycle == maximumCycles {
			return MaximumCyclesReached
		}
	}
	if err == CPUFinished {
		return nil
	}
	return err
}

func (cpu *CPU) Step() error {

	// First Move instructions to next stage of pipeline
	if err := cpu.Pipeline.TransferInstructions(); err != nil {
		return err
	}

	// Then check if execution is complete
	if cpu.Pipeline.Empty() && cpu.InstructionCacheEmpty() {
		//fmt.Println("pipeline empty", cpu.InstructionCacheEmpty(), cpu.InstructionPointer, cpu.Cycle)
		return CPUFinished
	}

	// If not, increase cycle count and execute pipeline
	//fmt.Println("#################### CYCLE", cpu.Cycle, "####################")
	cpu.Cycle += 1

	return cpu.Pipeline.Execute()
}

func (cpu *CPU) RenderState() string {
	result := new(bytes.Buffer)
	// spacing helper
	print := spacingHelper(10, result)

	print("c#%d", cpu.Cycle)
	for _, inst := range cpu.Instructions {
		inPipeline := false
		for _, iip := range cpu.Pipeline.ActiveInstructions() {
			if iip == inst {
				inPipeline = true
				stalled := iip.Stage.Stalled()
				if stalled {
					//print("(s) %s", iip.OpCode())
					print("(s)")
				} else {
					//print("%s %s", iip.Stage, iip.OpCode())
					if iip.CycleStart == cpu.Cycle {
						print("%s:%s", iip.Stage, iip.Instruction.OpCode())
					} else {
						print("%s", iip.Stage)
					}

				}
			}
		}
		if inPipeline == false {
			print("")
		}
	}
	result.WriteString("\n")
	//fmt.Println("\nstate: ", result)
	return string(result.Bytes())
}

func (cpu *CPU) RenderTiming() string {
	result := ""

	// render header
	format := "%-6s"
	result += fmt.Sprintf(format, "")
	for i, _ := range cpu.Instructions {
		result += fmt.Sprintf(format, fmt.Sprintf("I#%d", i+1))
	}
	result += "\n"

	// render for each cycle
	for i := 1; i <= cpu.Cycle; i++ {
		result += cpu.RenderTimingForCycle(i)
		result += "\n"
	}
	return result
}

func (cpu *CPU) RenderTimingForCycle(cycle int) string {
	result := new(bytes.Buffer)
	print := spacingHelper(6, result)
	print("c#%d", cycle)

	for _, inst := range cpu.Instructions {
		switch {
		case cycle < inst.CycleStart:
			print(".")
		case cycle > inst.CycleFinish:
			print("")
		case cycle == inst.CycleStart:
			print("IF1")
		// consider flushed cycles
		case cycle == inst.CycleFlush:
			print("(fl)")
		case inst.CycleFlush != -1 && cycle < inst.CycleFlush:
			print("(fl)")
		case cycle < inst.Stages["ID"]-2:
			print("(s)")
		case cycle == inst.Stages["ID"]-2:
			print("IF2")
		case cycle == inst.Stages["ID"]-1:
			print("IF3")
		default:
			if stage, ok := inst.Cycles[cycle]; ok {
				print(stage)
			} else {
				print("(s)")
			}
		}

	}

	return string(result.Bytes())
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
	return fmt.Sprintf("REGISTERS:\n%sMEMORY:\n%s", cpu.Registers, cpu.Ram)
}
