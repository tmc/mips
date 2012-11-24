// Simple simulator of a subset of the MIPS instruction set to illustrate pipelining
package mips

import (
	"errors"
	"fmt"
)

type Instruction interface {
	SetCPU(cpu *CPU)
	OpCode() string
	SetOpCode(opcode string)
	String() string
	Label() Label
	SetLabel(label Label)
	SetText(text string)
	Destination() Operand
	SetDestination(o Operand)
	OperandA() Operand
	SetOperandA(o Operand)
	OperandB() Operand
	SetOperandB(o Operand)

	IF1() error
	IF2() error
	IF3() error
	ID() error
	EX() error
	MEM1() error
	MEM2() error
	MEM3() error
	WB() error
}

type OperandType int

const (
	operandTypeInvalid OperandType = iota
	operandTypeImmediate
	operandTypeOffset
	operandTypeNormal
	operandTypeLabel
)

type Operand struct {
	text     string
	Register Register
	Offset   int
	Type     OperandType
}

type Label string

type instruction struct {
	cpu         *CPU
	label       Label
	text        string
	opcode      string
	destination Operand
	operandA    Operand
	operandB    Operand
}

func (op Operand) String() string {
	switch op.Type {
	case operandTypeImmediate:
		return fmt.Sprintf("#%d", op.Offset)
	case operandTypeNormal:
		return fmt.Sprintf("%s", op.Register)
	case operandTypeOffset:
		return fmt.Sprintf("%d(%s)", op.Offset, op.Register)
	case operandTypeLabel:
		return op.text
	}
	return "@unknown@"
}

func (op Operand) Value(cpu *CPU) (value Word, err error) {
	switch op.Type {
	case operandTypeImmediate:
		value = Word(op.Offset)
	case operandTypeNormal:
		value = cpu.Registers.Get(op.Register)
	case operandTypeOffset:
		value = cpu.Registers.Get(op.Register) + Word(op.Offset)
	case operandTypeLabel:
		value = Word(cpu.Labels[Label(op.text)])
	default:
		err = errors.New("Invalid operand type:")
	}
	return value, err
}

////////////////////////////////////////////////////////////////
// Instruction
////////////////////////////////////////////////////////////////

func NewInstruction(opcode string) (i Instruction, err error) {

	switch opcode {
	case "LD":
		i = new(LD)
	case "SD":
		i = new(SD)
	case "DADD":
		i = new(DADD)
	case "DADDI":
		i = new(DADDI)
	case "BNEZ":
		i = new(BNEZ)
	default:
		return nil, errors.New(fmt.Sprintf("Invalid opcode. %s", opcode))
	}
	i.SetOpCode(opcode)
	return
}

func (i instruction) String() string {
	result := fmt.Sprintf("%s %s %s", i.opcode, i.destination, i.operandA)
	if i.operandB.Type != operandTypeInvalid {
		result += fmt.Sprintf(" %s", i.operandB)
	}
	if i.label != "" {
		result += fmt.Sprintf(" (label: %s)", i.label)
	}
	return result
}

func (i instruction) OpCode() string {
	return i.opcode
}

func (i *instruction) SetOpCode(opcode string) {
	i.opcode = opcode
}

func (i *instruction) Label() Label {
	return i.label
}

func (i *instruction) SetLabel(label Label) {
	i.label = label
}

func (i *instruction) SetCPU(cpu *CPU) {
	i.cpu = cpu
}

func (i *instruction) SetText(text string) {
	i.text = text
}

func (i *instruction) Destination() Operand {
	return i.destination
}

func (i *instruction) SetDestination(d Operand) {
	i.destination = d
}

func (i *instruction) OperandA() Operand {
	return i.operandA
}

func (i *instruction) SetOperandA(o Operand) {
	i.operandA = o
}

func (i *instruction) OperandB() Operand {
	return i.operandB
}

func (i *instruction) SetOperandB(o Operand) {
	i.operandB = o
}

// Default blank stage implementations

func (i *instruction) IF1() error  { return nil }
func (i *instruction) IF2() error  { return nil }
func (i *instruction) IF3() error  { return nil }
func (i *instruction) ID() error   { return nil }
func (i *instruction) EX() error   { return nil }
func (i *instruction) MEM1() error { return nil }
func (i *instruction) MEM2() error { return nil }
func (i *instruction) MEM3() error { return nil }
func (i *instruction) WB() error   { return nil }

////////////////////////////////////////////////////////////////
// Actual Instruction Implementations
////////////////////////////////////////////////////////////////

type LD struct {
	instruction
	value Word
}

func (i *LD) MEM1() error {
	fmt.Println("MEM1 LD", i)

	value, err := i.operandA.Value(i.cpu)
	if err != nil {
		return err
	}
	fmt.Println("value", value)
	fmt.Println("value", i.cpu.Ram[value])
	i.value = i.cpu.Ram[value]
	return nil
}

func (i *LD) WB() error {
	fmt.Println("WD LD", i)
	return i.cpu.Registers.Set(i.destination.Register, i.value)
}

type SD struct {
	instruction
	value Word
}

func (i *SD) WB() error {
	fmt.Println("WD SD", i)
	val, err := i.operandA.Value(i.cpu)
	if err != nil {
		return err
	}
	dest, err := i.destination.Value(i.cpu)
	if err != nil {
		return err
	}
	i.cpu.Ram[dest] = val
	// @todo memory write errors, etc
	return nil
}

type ALUInstruction struct {
	instruction
	value Word
}

func (i *ALUInstruction) WB() error {
	fmt.Println("ALU LD", i)
	return i.cpu.Registers.Set(i.destination.Register, i.value)
}

type DADD struct {
	ALUInstruction
}

func (i *DADD) EX() error {
	fmt.Println("DADD EX", i)
	a, err := i.operandA.Value(i.cpu)
	if err != nil {
		return err
	}
	b, err := i.operandB.Value(i.cpu)
	if err != nil {
		return err
	}
	i.value = a + b
	// @todo consider overflow
	return nil
}

type DADDI struct {
	ALUInstruction
}

func (i *DADDI) EX() error {
	fmt.Println("DADDI EX", i)
	a, err := i.operandA.Value(i.cpu)
	if err != nil {
		return err
	}
	b, err := i.operandB.Value(i.cpu)
	if err != nil {
		return err
	}
	i.value = a + b
	// @todo consider overflow
	return nil
}

type BNEZ struct {
	instruction
}
