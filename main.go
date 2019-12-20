package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func setValue(key, value string) {
	var err error
	path := "/Users/andrew/work/puffadder/data.puff"
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("opening file", err)
	}
	_, err = f.Write([]byte(key + " " + value + "\n"))
	if err != nil {
		fmt.Println("writing file", err)
	}
}

func getValue(key string) string {
	path := "/Users/andrew/work/puffadder/data.puff"
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	value := ""
	for scanner.Scan() {
		line := scanner.Text()
		seperate := strings.Split(line, " ")
		if key == seperate[0] {
			value = seperate[1]
		}
	}
	return value
}

func main() {
	operation := os.Args[1]
	key := os.Args[2]
	if operation == "get" {
		fmt.Println(getValue(key))
	} else if operation == "set" {
		value := os.Args[3]
		setValue(key, value)
	} else {
		fmt.Println("Incorrect operation")
	}
}
