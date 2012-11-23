// Simple simulator of a subset of the MIPS instruction set to show pipelining
package mips

import (
	"errors"
	"fmt"
)

type Operation interface {
	OpCode() string
	String() string
	SetOpCode(opcode string)
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

type operation struct {
	opcode string
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

type Instruction struct {
	Label       Label
	text        string
	Operation   Operation
	Destination Operand
	OperandA    Operand
	OperandB    Operand
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

func (i Instruction) String() string {
	result := fmt.Sprintf("%s %s %s", i.Operation, i.Destination, i.OperandA)
	if i.OperandB.Type != operandTypeInvalid {
		result += fmt.Sprintf(" %s", i.OperandB)
	}
	//result := fmt.Sprintf("%s %s", i.Operation, i.Destination)
	if i.Label != "" {
		result += fmt.Sprintf(" (label: %s)", i.Label)
	}
	return result
}

func (i *Instruction) IF1() error  { return i.Operation.IF1() }
func (i *Instruction) IF2() error  { return i.Operation.IF2() }
func (i *Instruction) IF3() error  { return i.Operation.IF3() }
func (i *Instruction) ID() error   { return i.Operation.ID() }
func (i *Instruction) EX() error   { return i.Operation.EX() }
func (i *Instruction) MEM1() error { return i.Operation.MEM1() }
func (i *Instruction) MEM2() error { return i.Operation.MEM2() }
func (i *Instruction) MEM3() error { return i.Operation.MEM3() }
func (i *Instruction) WB() error   { return i.Operation.WB() }

////////////////////////////////////////////////////////////////
// Operations
////////////////////////////////////////////////////////////////

func NewOperation(opcode string) (op Operation, err error) {

	switch opcode {
	case "LD":
		op = new(LD)
	case "SD":
		op = new(SD)
	case "DADD":
		op = new(DADD)
	case "DADDI":
		op = new(DADDI)
	case "BNEZ":
		op = new(BNEZ)
	default:
		return nil, errors.New(fmt.Sprintf("Invalid opcode. %s", opcode))
	}
	op.SetOpCode(opcode)
	return
}

func (o operation) OpCode() string {
	return o.opcode
}

func (o operation) String() string {
	return o.opcode
}

func (o *operation) SetOpCode(opcode string) {
	o.opcode = opcode
}

func (o operation) IF1() error {
	logger.Debug("IF1", o)
	return nil
}
func (o operation) IF2() error {
	logger.Debug("IF2", o)
	return nil
}

func (o operation) IF3() error {
	logger.Debug("IF3", o)
	return nil
}
func (o operation) ID() error {
	logger.Debug("ID", o)
	return nil
}
func (o operation) EX() error {
	logger.Debug("EX", o)
	logger.Debugf("EX %s", o)
	return nil
}
func (o operation) MEM1() error {
	logger.Debug("MEM1", o)
	return nil
}
func (o operation) MEM2() error {
	logger.Debug("MEM2", o)
	return nil
}
func (o operation) MEM3() error {
	logger.Debug("MEM3", o)
	return nil
}
func (o operation) WB() error {
	logger.Debug("WB", o)
	return nil
}

type LD struct {
	operation
}

type SD struct {
	operation
}

type DADD struct {
	operation
}

type DADDI struct {
	operation
}

type BNEZ struct {
	operation
}
