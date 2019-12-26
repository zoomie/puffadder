package main

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"os"
	"path"
	"strconv"
	"strings"
)

const splitToken = ":"
const namePadding = "-"
const lineOffSet = 22
const nameMaxLength = 10
const algorithmType = "binaryTree"

type indexOffset interface {
	get(key string) (int64, bool)
	add(key string, value int64)
}

var dataPath string
var currentIndex indexOffset

func initCheckDataFileExists() {
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

func chooseIndex() {
	if algorithmType == "binaryTree" {
		currentIndex = &btree{}
	} else {
		currentIndex = hashTable{}
	}
}

func initLoadInMemoryMapping() {
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
		currentIndex.add(name, currentOffset)
		currentOffset += lineOffSet
	}
}

func init() {
	initCheckDataFileExists()
	chooseIndex()
	initLoadInMemoryMapping()
}

func encodeLine(name string, number string) (string, error) {
	if len(name) > 10 {
		return "", errors.New("name too long")
	}
	namePadded := name + strings.Repeat(namePadding, 10-len(name))
	numberPadded := strings.Repeat(namePadding, 10-len(number)) + number
	return namePadded + splitToken + numberPadded, nil
}

func decodeLine(raw string) (string, int, error) {
	splitRaw := strings.Split(raw, splitToken)
	name := strings.Trim(splitRaw[0], namePadding)
	amountString := strings.Trim(splitRaw[1], namePadding)
	amount, err := strconv.Atoi(amountString)
	if err != nil {
		return "", 0, fmt.Errorf("could not decode line: %w", err)
	}
	return name, amount, nil
}
func readAtByteOffset(offset int64) string {
	file, _ := os.Open(dataPath)
	// should you check err before calling defer.Close()
	// what if the file open fails, would defer.Close fail?
	defer file.Close()
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

func getAmount(name string) (int, error) {
	offset, ok := currentIndex.get(name)
	if !ok {
		// should I push this error up of handle it here?
		return 0, errors.New("name not in db")
	}
	line := readAtByteOffset(offset)
	_, amount, err := decodeLine(line)
	if err != nil {
		return 0, err
	}
	return amount, nil
}

func appendAmountToData(raw string) (int64, error) {
	var err error
	file, err := os.OpenFile(dataPath, os.O_APPEND|os.O_WRONLY, 0644)
	defer file.Close()
	// need to read what is the best practise for handling errors
	if err != nil {
		return 0, fmt.Errorf("unable to open file: %w", err)
	}
	_, err = file.WriteString(raw + "\n")
	if err != nil {
		return 0, fmt.Errorf("unable to write string: %w", err)
	}
	fileInto, err := file.Stat()
	if err != nil {
		return 0, fmt.Errorf("unable to get byte offset %w", err)
	}
	// could also keep track of file size with a var
	size := fileInto.Size()
	return size, nil
}

func setAmount(name, amount string) error {
	amountInt, err := strconv.Atoi(amount)
	if err != nil {
		return fmt.Errorf("unable to convert amount to int: %w", err)
	}
	if amountInt > math.MaxInt32 {
		return errors.New("the amout is too large")
	}
	raw, err := encodeLine(name, amount)
	if err != nil {
		return err
	}
	endOfFileOffset, err := appendAmountToData(raw)
	if err != nil {
		return fmt.Errorf("unable to get file offset: %w", err)
	}
	currentIndex.add(name, endOfFileOffset)
	return nil
}
