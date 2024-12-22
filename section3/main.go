package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args
	fmt.Println("hello, world")
	fmt.Printf("Arguments: %v\n", args)
	fmt.Printf("Arguments without program: %v\n", args[1:])
}
