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

	// Move instructions to next stage of pipeline unless we encounter a stall
	for i := len(cpu.Pipeline) - 1; i >= 0; i-- {
		stage := cpu.Pipeline[i]
		err := cpu.Pipeline.TransferInstruction(stage)
		if err == PipelineStall {
			//fmt.Println("Encountered stall, stopping instruction transfer.")
			//break
		} else if err != nil {
			return err
		}
	}

	for i := len(cpu.Pipeline) - 1; i >= 0; i-- {
		stage := cpu.Pipeline[i]
		//fmt.Println("################ stage", stage, stage.GetInstruction)
		stage.Unstall()
		switch err := stage.Step(); {
		case err == RAWException:
			//fmt.Println("STALL in", stage, stage.GetInstruction())
			stage.Stall()
		case err != nil:
			return errors.New(fmt.Sprintf("Error while executing %s of %s: %s", stage, stage.GetInstruction(), err))
		}
	}

	if cpu.Pipeline.Empty() && cpu.InstructionCacheEmpty() {
		//fmt.Println("pipeline empty", len(cpu.Pipeline.ActiveInstructions()))
		return CPUFinished
	}

	return nil
}

func (cpu *CPU) PrintTiming(printHeader bool) {
	print := func(format string, args ...interface{}) {
		fmt.Printf("%-12s", fmt.Sprintf(format, args...))
	}
	if printHeader {
		print("")
		for i, _ := range cpu.InstructionCache {
			print("I#%d", i+1)
		}
		fmt.Println("")
	}

	print("c#%d", cpu.Cycle)
	for _, inst := range cpu.InstructionCache {
		inPipeline := false
		for _, iip := range cpu.Pipeline.ActiveInstructions() {
			if iip.Instruction == inst {
				inPipeline = true
				stalled := iip.Stage.Stalled()
				if stalled {
					//print("(s) %s", iip.OpCode())
					print("(s)")
				} else {
					//print("%s %s", iip.Stage, iip.OpCode())
					print("%s", iip.Stage)

				}
			}
		}
		if inPipeline == false {
			fmt.Printf("%-12s", "")
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
