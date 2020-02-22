package main

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

const splitToken = ":"
const namePadding = "-"
const nameMaxLength = 10

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

func getAmount(index indexOffset, name string) (int, error) {
	offset, ok := index.get(name)
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

func setAmount(index indexOffset, name string, amount int) error {
	if amount > math.MaxInt32 {
		return errors.New("the value is too large")
	}
	amountString := strconv.Itoa(amount)
	raw, err := encodeLine(name, amountString)
	if err != nil {
		return err
	}
	endOfFileOffset, err := appendAmountToData(raw)
	if err != nil {
		return fmt.Errorf("unable to get file offset: %w", err)
	}
	index.add(name, endOfFileOffset)
	return nil
}

// func transaction(name1, name2, amount1, amount2 string) error {
// 	// if both operations are valid then process the transactions.
// }
