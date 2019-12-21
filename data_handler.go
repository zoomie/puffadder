package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

const path = "/Users/andrew/go/src/github.com/zoomie/puffadder/data.puff"

// SetValue inputs the value into the .puff file
func SetValue(key, value string) {
	var err error
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("opening file", err)
	}
	_, err = f.Write([]byte(key + " " + value + "\n"))
	if err != nil {
		fmt.Println("writing file", err)
	}
}

// GetValue retrives the value from the .puff file
func GetValue(key string) string {
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
