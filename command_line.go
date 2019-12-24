package main

import (
	"fmt"
	"os"
)

func main() {
	operation := os.Args[1]
	name := os.Args[2]
	if operation == "get" {
		amount, err := getAmount(name)
		if err != nil {
			panic((err))
		}
		fmt.Println(amount)
	} else if operation == "set" {
		amount := os.Args[3]
		err := setAmount(name, amount)
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Println("Incorrect operation")
	}
}
