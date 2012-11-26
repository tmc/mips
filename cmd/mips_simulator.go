package main

import (
	"bufio"
	"fmt"
	"mips"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

var header = `MIPS Pipeline Simulator
EECS645 - cline@ku.edu
`

const (
	modeNoForwarding = iota
	modePredictTaken
	modePredictNotTaken
)

var modeNames = map[int]string{
	modeNoForwarding:    "No Forwarding/Bypassing",
	modePredictTaken:    "Predict Branches Taken",
	modePredictNotTaken: "Predict Branches Not Taken",
}

type simulator struct {
	i io.Reader
	o io.Writer
	//e io.Writer
	inputFile, registerFile, timingFile string
	//input, register, timing io.WriteCloser
	simulatorMode int
}

func NewSimulator(i io.Reader, o io.Writer) (*simulator, error) {
	s := new(simulator)
	s.i, s.o = i, o
	//s.e = e
	return s, s.collectInput()
}

type niceReader struct {
	*bufio.Reader
}

func (nr niceReader) ReadLine() (string, error) {
	result, err := nr.ReadString('\n')
	return strings.TrimRight(result, "\n"), err
}

func (s *simulator) collectInput() (err error) {
	buf := niceReader{bufio.NewReader(s.i)}

	// input
	fmt.Fprintf(s.o, "Enter input filename: ")
	if s.inputFile, err = buf.ReadLine(); err != nil {
		return err
	}

	// mode
	for _, mode := range []int{
		modeNoForwarding,
		modePredictTaken,
		modePredictNotTaken,
	} {
		fmt.Fprintf(s.o, "%d: %s\n", mode, modeNames[mode])
	}
	fmt.Fprintf(s.o, "\nSelect Mode: ")
	if strMode, err := buf.ReadLine(); err != nil {
		return err
	} else if s.simulatorMode, err = strconv.Atoi(strMode); err != nil {
		return err
	}

	// register
	fmt.Fprintf(s.o, "\nEnter register filename: ")
	if s.registerFile, err = buf.ReadLine(); err != nil {
		return err
	}

	// input
	fmt.Fprintf(s.o, "\nEnter timing filename: ")
	if s.timingFile, err = buf.ReadLine(); err != nil {
		return err
	}

	return nil
}

func (s *simulator) xcollectFiles() error {

	for _, info := range []struct {
		name string
		dest *string
	}{
		{"input", &s.inputFile},
		{"timing", &s.timingFile},
		{"register", &s.registerFile},
	} {

		fmt.Fprintf(s.o, "Enter %s filename: ", info.name)
		b := bufio.NewReader(s.i)
		str, err := b.ReadString('\n')
		*(info.dest) = str
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *simulator) run() error {
	input, err := ioutil.ReadFile(s.inputFile)
	if err != nil {
		return err
	}

	cpu, err := mips.ParseCPUString(string(input))
	if err != nil {
		return err
	}

	switch s.simulatorMode {
	case modeNoForwarding:
		cpu.ForwardingEnabled = false
		cpu.BranchMode = mips.BranchPolicyFlush
	case modePredictTaken:
		cpu.ForwardingEnabled = true
		cpu.BranchMode = mips.BranchPolicyPredictTaken
	case modePredictNotTaken:
		cpu.ForwardingEnabled = true
		cpu.BranchMode = mips.BranchPolicyPredictTaken
	}

	fmt.Fprintln(s.o, "\nRunning Simulation.\n")
	if err := cpu.Run(10000); err != nil {
		fmt.Fprintln(os.Stderr, "Error running simulation:", err)
		return err
	}

	if err := ioutil.WriteFile(s.timingFile, []byte(cpu.RenderTiming()), 0666); err != nil {
		return err
	}
	fmt.Fprintln(s.o, "Wrote timing to", s.timingFile)
	if err := ioutil.WriteFile(s.registerFile, []byte(cpu.Registers.String()), 0666); err != nil {
		return err
	}
	fmt.Fprintln(s.o, "Wrote register contents to", s.registerFile)
	return nil
}

func main() {
	// @todo prompt user for input and output files
	//m := mips.NewMachine()

	fmt.Println(header)
	s, err := NewSimulator(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	if err := s.run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("done.")
}
