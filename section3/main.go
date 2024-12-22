package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args

	if len(args) < 2 {
		fmt.Printf("Usage: %v <argument>\n", args[0])
		os.Exit(1)
	}

	fmt.Println("hello, world")
	fmt.Printf("Arguments: %v\n", args)
	fmt.Printf("Arguments without program: %v\n", args[1:])

	fmt.Printf("1st Argument: %v\n", args[1])
}
