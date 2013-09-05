package mips

import (
	"errors"
	"fmt"
)

var (
	PipelineStall = errors.New("Pipeline Stalled")
)

type ExecutedInstruction struct {
	Instruction
	Stage       PipelineStage
	Stages      map[string]int // map of Stages to cycle at which that stage was entered
	Cycles      map[int]string // map of cycles to Stages
	CycleStart  int
	CycleFinish int
	CycleFlush  int
}

type Pipeline []PipelineStage

type PipelineStage interface {
	Initialize(cpu *CPU)
	CPU() *CPU
	String() string
	Step() error
	Stall()
	Unstall()
	Stalled() bool
	SetInstruction(instruction *ExecutedInstruction)
	GetInstruction() *ExecutedInstruction
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

func (p Pipeline) cpu() *CPU { return p[0].CPU() }

func (p Pipeline) Reverse() []PipelineStage {
	result := make([]PipelineStage, len(p))
	for i := 0; i < len(p); i++ {
		result[i] = p[len(p)-i-1]
	}
	return result
}

// Execute a cycle in the pipeline
func (p Pipeline) Execute() error {

	// run pipeline pipeline stages back to front to execute older instructions first
	//for _, stage := range p.Reverse() {
	for i := len(p) - 1; i >= 0; i-- {
		stage := p[i]

		stage.Unstall()
		switch err := stage.Step(); {
		case err == RAWHazard:
			//fmt.Println("RAWHazard in", stage, stage.GetInstruction(), "stalling")
			stage.Stall()
			return nil
		case err == FlushPipeline:
			// flush and finish cycle
			p.FlushBefore(stage)
			i = 0
			p.RecordTiming(stage)
		case err == BranchResolving:
			p.StallBefore(stage)
			p.RecordTiming(stage)
		case err == nil:
			// entered stage successfully, record timing if an instruction is present
			p.RecordTiming(stage)
		case err != nil:
			return errors.New(fmt.Sprintf("Error while executing %s of %s: %s", stage, stage.GetInstruction(), err))
		}

	}
	return nil
}

func (p Pipeline) TransferInstruction(fromStage PipelineStage) error {
	if fromStage.Stalled() {
		//fmt.Println("TransferInstruction fromstage stalled", fromStage)
		return nil
	}
	instruction := fromStage.GetInstruction()
	toStage := fromStage.Next()
	if toStage != nil {
		if instruction != nil {
			//fmt.Println("TransferInstruction ", fromStage, toStage, instruction)
		}
		toStage.SetInstruction(instruction)
		if instruction != nil {
			//fmt.Println("TransferInstruction ::: ", toStage.GetInstruction(), &toStage)
		}
	}
	if instruction != nil {
		instruction.Stage = toStage
	}
	fromStage.SetInstruction(nil)
	return nil
}

func (p Pipeline) RecordTiming(stage PipelineStage) {
	if i := stage.GetInstruction(); i != nil {
		i.Stages[stage.String()] = p.cpu().Cycle
		i.Cycles[p.cpu().Cycle] = stage.String()
	}
}

func (p Pipeline) TransferInstructions() error {
	for i := len(p) - 1; i >= 0; i-- {
		stage := p[i]
		err := p.TransferInstruction(stage)
		if err == PipelineStall {
			fmt.Println("Encountered stall, stopping instruction transfer.")
			break
		} else if err != nil {
			return err
		}
	}
	return nil
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

func (p Pipeline) StallBefore(stage PipelineStage) {
	stage = stage.Prev()
	for stage != nil {
		stage.Stall()
		stage = stage.Prev()
	}
}

func (p Pipeline) FlushBefore(stage PipelineStage) {
	stage = stage.Prev()
	for stage != nil {
		//fmt.Println("flushing", stage, stage.GetInstruction())
		i := stage.GetInstruction()
		if i != nil {
			i.Flush()
			i.CycleFlush = p.cpu().Cycle
			i.CycleFinish = p.cpu().Cycle
		}
		stage.SetInstruction(nil)
		//fmt.Println("flushed", stage, stage.GetInstruction())
		stage = stage.Prev()
	}
}

func (p Pipeline) xFlush(currentCycle int) {
	for _, stage := range p {
		i := stage.GetInstruction()
		if i != nil {
			i.CycleFlush = currentCycle
			i.CycleFinish = currentCycle
		}
		stage.SetInstruction(nil)
	}
}

func (p Pipeline) StageNames() (result []string) {
	if result, cached := stageStringCache[&p]; cached {
		return result
	}
	result = make([]string, 0)
	for _, stage := range p {
		result = append(result, stage.String())
	}
	stageStringCache[&p] = result
	return result
}

func (p Pipeline) ActiveInstructions() []*ExecutedInstruction {
	result := make([]*ExecutedInstruction, 0)

	for _, stage := range p {
		iip := stage.GetInstruction()
		if iip != nil {
			result = append(result, iip)
		}
	}

	return result
}

type stage struct {
	instruction *ExecutedInstruction
	stalled     bool
	cpu         *CPU
	next        PipelineStage
	prev        PipelineStage
}

func (s *stage) Initialize(cpu *CPU) {
	s.cpu = cpu
}

func (s *stage) CPU() *CPU {
	return s.cpu
}

func (s *stage) Addr() string {
	return fmt.Sprint(&s)
}

func (s stage) String() string {
	return "unknown"
}

func (s *stage) Stall() {
	s.stalled = true
}

func (s *stage) Unstall() {
	s.stalled = false
}

func (s *stage) Stalled() bool {
	return s.stalled || (s.Next() != nil && s.Next().Stalled())
}

func (s *stage) Prev() PipelineStage {
	return s.prev

}
func (s *stage) Next() PipelineStage {
	return s.next
}

func (s *stage) SetPrev(p PipelineStage) {
	s.prev = p

}
func (s *stage) SetNext(p PipelineStage) {
	s.next = p
}

func (s *stage) GetInstruction() *ExecutedInstruction {
	return s.instruction
}

func (s *stage) SetInstruction(instruction *ExecutedInstruction) {
	s.instruction = instruction
}

/////////////////////////////////////////////////////////////////////////////
// IF1
/////////////////////////////////////////////////////////////////////////////

type IF1 struct{ stage }

func (s IF1) String() string { return "IF1" }

func (s *IF1) Step() error {

	// If we already have an instruction (as a result of a stall) then no-op
	if s.instruction != nil {
		fmt.Println("not loading new instruction, stalled.", s.instruction.OpCode())
		return nil
	}

	// otherwise fetch a new instruction if possible

	if s.cpu.InstructionCacheEmpty() == false {

		s.instruction = &ExecutedInstruction{
			Instruction: s.cpu.InstructionCache[s.cpu.InstructionPointer],
			Stage:       s,
			Stages:      make(map[string]int, 0),
			Cycles:      make(map[int]string, 0),
			CycleStart:  s.cpu.Cycle, // Start
			CycleFinish: -1,
			CycleFlush:  -1,
		}

		// record instuction in cpu's list of execut(ed|ing) instructions
		s.cpu.Instructions = append(s.cpu.Instructions, s.instruction)

		//fmt.Println("Issue:", s.instruction)
		s.cpu.InstructionPointer += 1
		return s.instruction.IF1()
	}
	return nil
}

/////////////////////////////////////////////////////////////////////////////
// IF2
/////////////////////////////////////////////////////////////////////////////

type IF2 struct{ stage }

func (s IF2) String() string { return "IF2" }

func (s *IF2) Step() error {
	if s.instruction == nil {
		return nil
	}
	return s.instruction.IF2()
}

/////////////////////////////////////////////////////////////////////////////
// IF3
/////////////////////////////////////////////////////////////////////////////

type IF3 struct{ stage }

func (s IF3) String() string { return "IF3" }

func (s *IF3) Step() error {
	if s.instruction == nil {
		return nil
	}
	return s.instruction.IF3()
}

/////////////////////////////////////////////////////////////////////////////
// ID
/////////////////////////////////////////////////////////////////////////////

type ID struct{ stage }

func (s ID) String() string { return "ID" }

func (s *ID) Step() error {

	if s.instruction == nil {
		return nil
	}
	err := s.instruction.ID()
	return err
}

/////////////////////////////////////////////////////////////////////////////
// EX
/////////////////////////////////////////////////////////////////////////////

type EX struct{ stage }

func (s EX) String() string { return "EX" }

func (s *EX) Step() error {
	if s.instruction == nil {
		return nil
	}
	return s.instruction.EX()
}

/////////////////////////////////////////////////////////////////////////////
// MEM1
/////////////////////////////////////////////////////////////////////////////

type MEM1 struct{ stage }

func (s MEM1) String() string { return "MEM1" }

func (s *MEM1) Step() error {
	if s.instruction == nil {
		return nil
	}
	return s.instruction.MEM1()
}

/////////////////////////////////////////////////////////////////////////////
// MEM2
/////////////////////////////////////////////////////////////////////////////

type MEM2 struct{ stage }

func (s MEM2) String() string { return "MEM2" }

func (s *MEM2) Step() error {
	if s.instruction == nil {
		return nil
	}
	return s.instruction.MEM2()
}

/////////////////////////////////////////////////////////////////////////////
// MEM3
/////////////////////////////////////////////////////////////////////////////

type MEM3 struct{ stage }

func (s MEM3) String() string { return "MEM3" }

func (s *MEM3) Step() error {
	if s.instruction == nil {
		return nil
	}
	return s.instruction.MEM3()
}

/////////////////////////////////////////////////////////////////////////////
// 
/////////////////////////////////////////////////////////////////////////////

type WB struct{ stage }

func (s WB) String() string { return "WB" }

func (s *WB) Step() error {
	if s.instruction == nil {
		return nil
	}
	err := s.instruction.WB()
	if err == nil {
		s.instruction.CycleFinish = s.cpu.Cycle
	}
	return nil
}

// internal caching for stage list generation
var stageStringCache map[*Pipeline][]string

func init() {
	stageStringCache = make(map[*Pipeline][]string)
}
