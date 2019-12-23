package main

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
)

const splitToken = ":"
const namePadding = "-"
const lineOffSet = 22
const nameMaxLength = 10
// add maxInt size

var dataPath string
var nameToOffset map[string]int64

func initCheckDataFileExists() {
	workingDir, _ := os.Getwd()
	dataPath = path.Join(workingDir, "data.puff")
	if _, err := os.Stat(dataPath); err == nil {
		fmt.Println("Data path exists at: ", dataPath)
	} else {
		_, err := os.Create(dataPath)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Data path created at:", dataPath)
		}
	}
}

func initLoadInMemoryMapping() {
	file, _ := os.Open(dataPath)
	scanner := bufio.NewScanner(file)
	nameToOffset = make(map[string]int64)
	var currentOffset int64 = 0
	for scanner.Scan() {
		line := scanner.Text()
		name, _ := decodeLine(line)
		nameToOffset[name] = currentOffset
		currentOffset += lineOffSet
	}
}

func init() {
	initCheckDataFileExists()
	initLoadInMemoryMapping()
}

func encodeLine(name string, number string) string {
	namePadded := name + strings.Repeat(namePadding, 10-len(name))
	numberPadded := strings.Repeat(namePadding, 10-len(number)) + number
	return namePadded + splitToken + numberPadded
}

func decodeLine(raw string) (string, int) {
	splitRaw := strings.Split(raw, splitToken)
	name := strings.Trim(splitRaw[0], namePadding)
	amountString := strings.Trim(splitRaw[1], namePadding)
	amount, _ := strconv.Atoi(amountString)
	return name, amount
}
func readAtByteOffset(offset int64) string {
	file, _ := os.Open(dataPath)
	var whence int = 0 // read from start of file
	_, err := file.Seek(offset, whence)
	if err != nil {
		fmt.Println(err)
	}
	// the scanner defaults to using '\n' to tokenize
	scanner := bufio.NewScanner(file)
	// Only scan one line, until the next token
	scanner.Scan()
	line := scanner.Text()
	return line
}

func getAmount(key string) int {
	offset, ok := nameToOffset[key]
	if !ok {
		fmt.Println("Perons does not exist in db")
		return 0
	}
	line := readAtByteOffset(offset)
	_, amount := decodeLine(line)
	return amount
}

func writeAtByteOffset(raw string, offset int64) {
	// Not currently working
	// the writer is not actually writing to file at the offset.

	// file, _ := os.Create(dataPath) //, os.O_WRONLY, os.ModeExclusive)
	// whence := 0
	// _, err := file.Seek(offset, whence) // add error handling
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// writer := bufio.NewWriter(file)
	// // rawBytes := []byte(raw)
	// nn, err := writer.WriteString(raw)
	// if err != nil {
	// 	fmt.Println(err, nn)
	// } else {
	// 	fmt.Println("saved data to file")
	// }
}
func appendAmountToData(raw string){
	var err error
	file, err := os.OpenFile(dataPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("opening file", err)
	}
	_, err = file.WriteString(raw + "\n")
	if err != nil {
		fmt.Println("writing file", err)
	}
}

func setAmount(name, amount string) {
	offset, ok := nameToOffset[name]
	raw := encodeLine(name, amount)
	if ok {
		// update the key in place
		// writeAtByteOffset(raw, offset)
		fmt.Println("TOTO:", offset)
	} else {
		// append the key to file
		appendAmountToData(raw)
	}
}