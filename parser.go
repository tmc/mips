// Parses input files for simple mips pipeline simulator
package mips

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
)

type parserState byte

const (
	stateStart = iota
	stateRegisters
	stateMemory
	stateCode
	stateFinished
	stateLabel
	stateOperation
	stateDestination
	stateOperand1
	stateOperand2
)

type cpuParser struct {
	cpu         *CPU
	lines       []string
	currentLine int
	state       parserState
}

func newCPUParser(cpu *CPU, input io.Reader) (*cpuParser, error) {
	content, err := ioutil.ReadAll(input)
	mp := &cpuParser{
		cpu:   cpu,
		lines: strings.Split(string(content), "\n"),
	}
	return mp, err
}

func (mp *cpuParser) current() string {
	return strings.TrimSpace(mp.lines[mp.currentLine])
}

func (mp *cpuParser) next() (string, error) {
	mp.currentLine += 1
	if mp.currentLine >= len(mp.lines) {
		return "", io.EOF
	}
	return strings.TrimSpace(mp.lines[mp.currentLine]), nil
}

func (mp *cpuParser) peek() (string, error) {
	defer func() {
		mp.currentLine -= 1
	}()
	return mp.next()
}

func (mp *cpuParser) parseError(msg string) error {
	return errors.New(fmt.Sprintf("(line %d: %s) %s", mp.currentLine, mp.current(), msg))
}

func (mp *cpuParser) Parse() (m *CPU, err error) {
	m = mp.cpu
	for mp.state != stateFinished {
		switch mp.state {
		case stateStart:
			if mp.current() != "REGISTERS" {
				return nil, mp.parseError("REGISTERS was expected")
			}
			mp.state = stateRegisters
		case stateRegisters:
			if s, _ := mp.next(); s == "MEMORY" {
				mp.state = stateMemory
			} else {
				parts := removeEmpty(strings.Split(mp.current(), " "))

				if len(parts) != 2 {
					return nil, mp.parseError("unexpected number of parts in register statement")
				}
				reg, val := parts[0], parts[1]
				operand, err := ParseOperand(reg)
				if err != nil {
					return nil, err
				}
				intVal, err := strconv.Atoi(val)
				if err != nil {
					return nil, err
				}
				m.Registers.Set(operand.Register, Word(intVal))
			}
		case stateMemory:
			if s, _ := mp.next(); s == "CODE" {
				mp.state = stateCode
			} else {
				parts := removeEmpty(strings.Split(mp.current(), " "))

				if len(parts) != 2 {
					return nil, mp.parseError("unexpected number of pargs in memory statement")
				}
				mem, val := parts[0], parts[1]
				memPos, err := strconv.Atoi(mem)
				if err != nil {
					return nil, err
				}
				intVal, err := strconv.Atoi(val)
				if err != nil {
					return nil, err
				}
				m.Ram[memPos] = Word(intVal)
			}
		case stateCode:
			if _, e := mp.next(); e == io.EOF || mp.current() == "" {
				mp.state = stateFinished
			} else {
				//fmt.Println("process code: ", mp.current())
				instruction, err := ParseInstruction(strings.NewReader(mp.current()))
				if err != nil {
					return nil, mp.parseError(fmt.Sprintf("Instruction parse error: %s", err))
				}
				mp.cpu.InstructionCache = append(mp.cpu.InstructionCache, instruction)

				// if there's a label, store it in the label -> IC addr map
				if instruction.Label() != "" {
					mp.cpu.Labels[instruction.Label()] = len(mp.cpu.InstructionCache)
				}
			}
		}
	}
	return m, nil
}

func ParseCPUString(input string) (*CPU, error) {
	return ParseCPU(strings.NewReader(input))
}

func ParseCPU(input io.Reader) (*CPU, error) {
	p, err := newCPUParser(NewCPU(), input)
	if err != nil {
		return nil, err
	}
	return p.Parse()
}

func ParseOperand(s string) (o Operand, err error) {
	o.text = s

	//#-8
	if s[0] == '#' {
		o.Offset, err = strconv.Atoi(s[1:])
		o.Register = None
		o.Type = operandTypeImmediate
		if err != nil {
			return o, err
		}
		//R4
	} else if s[0] == 'R' {
		iVal, err := strconv.Atoi(s[1:])
		if err != nil {
			return o, err
		}
		o.Register = Register(iVal)
		o.Type = operandTypeNormal
		//16(R2)
	} else if strings.Index(s, "(") != -1 && strings.Index(s, ")") != -1 {
		parenOpen := strings.Index(s, "(")
		parenClose := strings.Index(s, ")")
		o.Offset, err = strconv.Atoi(s[:parenOpen])
		if err != nil {
			return o, err
		}
		iVal, err := strconv.Atoi(s[parenOpen+2 : parenClose])
		if err != nil {
			return o, err
		}
		o.Register = Register(iVal)
		o.Type = operandTypeOffset
		// Loop
	} else {
		o.Type = operandTypeLabel
	}
	return o, err
}

type instructionParser struct {
	line        string
	instruction *Instruction
	state       parserState
	pos         int
}

func newInstructionParser(input io.Reader) (*instructionParser, error) {
	s, err := ioutil.ReadAll(input)
	if err != nil {
		return nil, err
	}
	return &instructionParser{
		line:        string(s),
		instruction: new(Instruction),
	}, nil
}

func removeEmpty(i []string) []string {
	result := make([]string, 0)
	for _, s := range i {
		if s != "" {
			result = append(result, s)
		}
	}
	return result
}

func (ip *instructionParser) Parse() (i Instruction, err error) {
	ip.line = strings.TrimSpace(ip.line)

	line := ip.line
	parts := removeEmpty(strings.Split(line, " "))
	instruction_label := ""

	for ip.state != stateFinished {
		if len(parts) == 0 {
			return nil, errors.New("Invalid instruction input")
		}
		switch ip.state {
		case stateStart:
			if strings.Contains(parts[0], ":") {
				ip.state = stateLabel
			} else {
				ip.state = stateOperation
			}

		case stateLabel:
			instruction_label = parts[0][:strings.Index(parts[0], ":")]
			parts = parts[1:]
			ip.state = stateOperation

		case stateOperation:
			i, err = NewInstruction(parts[0])
			if err != nil {
				return nil, err
			}
			i.SetLabel(Label(instruction_label))
			i.SetText(line)

			parts = parts[1:]
			ip.state = stateDestination

		case stateDestination:
			parts[0] = strings.Trim(parts[0], ",")
			destination, err := ParseOperand(parts[0])
			if err != nil {
				return nil, err
			}
			i.SetDestination(destination)
			parts = parts[1:]
			ip.state = stateOperand1

		case stateOperand1:
			parts[0] = strings.Trim(parts[0], ",")
			operandA, err := ParseOperand(parts[0])
			if err != nil {
				return nil, err
			}
			i.SetOperandA(operandA)
			parts = parts[1:]
			if len(parts) > 0 {
				ip.state = stateOperand2
			} else {
				ip.state = stateFinished
			}

		case stateOperand2:
			parts[0] = strings.Trim(parts[0], ",")
			operandB, err := ParseOperand(parts[0])
			if err != nil {
				return nil, err
			}
			i.SetOperandB(operandB)
			parts = parts[1:]
			if len(parts) > 0 {
				return nil, errors.New(fmt.Sprintf("Extra content: %s", parts))
			} else {
				ip.state = stateFinished
			}
		}
	}
	return i, nil
}

func ParseInstruction(input io.Reader) (Instruction, error) {
	p, err := newInstructionParser(input)
	if err != nil {
		return nil, err
	}
	return p.Parse()
}
