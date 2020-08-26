package main

import (
	"fmt"
)

func state(format string, args ...interface{}) {
	fmt.Printf("[ ] %s ...\r", fmt.Sprintf(format, args...))
}

func printErrorMsg(format string, args ...interface{}) {
	fmt.Println("[-]")
	fmt.Printf("  %s\n", fmt.Sprintf(format, args...))
}

func printError(err error) error {
	fmt.Println("[-]")
	fmt.Printf("  Error: %s\n", err)
	return err
}
