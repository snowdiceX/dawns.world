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

// Debugf output log with debug level
func Debugf(format string, msgs ...interface{}) {
	fmt.Println(fmt.Sprintf(format, msgs...))
}

// Warnf output log with warn level
func Warnf(format string, msgs ...interface{}) {
	fmt.Println(fmt.Sprintf(format, msgs...))
}

// Error output log with error level
func Error(msgs ...interface{}) {
	fmt.Println(msgs...)
}

// Errorf output log with error level
func Errorf(format string, msgs ...interface{}) {
	fmt.Println(fmt.Sprintf(format, msgs...))
}
