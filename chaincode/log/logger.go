package log

import (
	"fmt"
)

// Info output log with info level
func Info(msgs ...interface{}) {
	fmt.Println(msgs...)
}

// Debug output log with debug level
func Debug(msgs ...interface{}) {
	fmt.Println(msgs...)
}

// Error output log with error level
func Error(msgs ...interface{}) {
	fmt.Println(msgs...)
}
