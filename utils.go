package mips

import (
	"fmt"
)

func spacingHelper(spaces int) func(format string, args ...interface{}) string {
	formatString := fmt.Sprintf("%%-%ds", spaces)
	return func(format string, args ...interface{}) string {
		return fmt.Sprintf(formatString, fmt.Sprintf(format, args...))
	}
}
