package mips

import (
	"bytes"
	"fmt"
)

// takes a spacing argument and byte slice, returns a method that adds to provided slice

func spacingHelper(spaces int, buf *bytes.Buffer) func(format string, args ...interface{}) {
	formatString := fmt.Sprintf("%%-%ds", spaces)
	return func(format string, args ...interface{}) {
		buf.WriteString(fmt.Sprintf(formatString, fmt.Sprintf(format, args...)))
	}
}
