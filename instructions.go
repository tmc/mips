// Simple simulator of a subset of the MIPS instruction set to show pipelining
package mips

import (
	"fmt"
)

type Operation struct {
	op string
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


func (op Operation) String() string {
	return op.op
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
