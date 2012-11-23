package mips

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
)

const DEBUG = 2
const WARN = 1
const ERROR = 0

const logger = logging(DEBUG) // set to the level you want

type logging int

func getCaller() string {
	if _, file, line, ok := runtime.Caller(2); ok {
		return fmt.Sprintf("%s:%d", filepath.Base(file), line)
	}
	return "unknown"
}

func print(args ...interface{}) {
	args = append([]interface{}{getCaller()}, args...)
	log.Println(args...)
}

func printf(format string, args ...interface{}) {
	format = fmt.Sprintf("%s %s", getCaller(), format)
	log.Printf(format, args...)
}

func (l logging) Debug(args ...interface{}) {
	if l >= DEBUG {
		print(args...)
	}
}

func (l logging) Warn(args ...interface{}) {
	if l >= WARN {
		print(args...)
	}
}

func (l logging) Error(args ...interface{}) {
	if l >= ERROR {
		print(args...)
	}
}

func (l logging) Debugf(format string, args ...interface{}) {
	if l >= DEBUG {
		printf(format, args...)
	}
}

func (l logging) Warnf(format string, args ...interface{}) {
	if l >= WARN {
		printf(format, args...)
	}
}

func (l logging) Errorf(format string, args ...interface{}) {
	if l >= ERROR {
		printf(format, args...)
	}
}
