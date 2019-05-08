package log

import (
	"fmt"
)

// Info output log with info level
func Info(msgs ...interface{}) {
	fmt.Println(msgs...)
}

// Infof output log with info level
func Infof(format string, msgs ...interface{}) {
	fmt.Println(fmt.Sprintf(format, msgs...))
}

// Debug output log with debug level
func Debug(msgs ...interface{}) {
	fmt.Println(msgs...)
}

// Debugf output log with debug level
func Debugf(format string, msgs ...interface{}) {
	fmt.Println(fmt.Sprintf(format, msgs...))
}

// Warn output log with warn level
func Warn(msgs ...interface{}) {
	fmt.Println(msgs...)
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
