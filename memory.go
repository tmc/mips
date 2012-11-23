// Simple simulator of a subset of the MIPS instruction set to show pipelining
package mips

import (
	"fmt"
)

type Word uint64

const memorySize = 992 // Size of memory in words

type Memory [memorySize]Word

func (w Word) String() string {
	return fmt.Sprintf("%#x", uint64(w))
}

func (r Memory) String() string {
	result := ""
	for i := 0; i < memorySize; i++ {
		if r[i] != 0 {
			result += fmt.Sprintf("\n%#x = %d", i, r[i])
		}
	}
	return result
}
