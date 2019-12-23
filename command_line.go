package main

import (
	"fmt"
	"os"
)

func main() {
	operation := os.Args[1]
	name := os.Args[2]
	if operation == "get" {
		amount := getAmount(name)
		fmt.Println(amount)
	} else if operation == "set" {
		amount := os.Args[3]
		setAmount(name, amount)
	} else {
		fmt.Println("Incorrect operation")
	}
}
