package main

import (
	"bufio"
	"fmt"
	"os"
	"path"
)

const algorithmType = "binaryTree"
const lineOffSet = 22

type indexOffset interface {
	get(key string) (int64, bool)
	add(key string, value int64)
}

var dataPath string
var currentIndex indexOffset

func init() {
	// Check if data file exists
	workingDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	dataPath = path.Join(workingDir, "data.puff")
	if _, err := os.Stat(dataPath); err == nil {
		fmt.Println("Data path exists at: ", dataPath)
	} else {
		_, err := os.Create(dataPath)
		if err != nil {
			panic(fmt.Errorf("Could not create data file: %w", err))
		} else {
			fmt.Println("Data path created at:", dataPath)
		}
	}
}

func chooseIndex() indexOffset {
	if algorithmType == "binaryTree" {
		return &btree{}
	} else {
		return hashTable{}
	}
}

func initLoadInMemoryMapping(index indexOffset) {
	file, _ := os.Open(dataPath)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var currentOffset int64 = 0
	for scanner.Scan() {
		line := scanner.Text()
		name, _, err := decodeLine(line)
		if err != nil {
			panic(fmt.Errorf("data file corrupt: %w", err))
		}
		index.add(name, currentOffset)
		currentOffset += lineOffSet
	}
}

func main() {
	index := chooseIndex()
	initLoadInMemoryMapping(index)
	operation := os.Args[1]
	name := os.Args[2]
	if operation == "get" {
		amount, err := getAmount(index, name)
		if err != nil {
			panic((err))
		}
		fmt.Println(amount)
	} else if operation == "set" {
		amount := os.Args[3]
		err := setAmount(index, name, amount)
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Println("Incorrect operation")
	}
}
