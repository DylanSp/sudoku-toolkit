package utils

import "fmt"

func Assert(condition bool, failureMsg string) {
	if !condition {
		panic(failureMsg)
	}
}

func Assertf(condition bool, format string, args ...any) {
	if !condition {
		msg := fmt.Sprintf(format, args...)
		fmt.Print(msg)
		panic(msg)
	}
}
