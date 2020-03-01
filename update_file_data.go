package main

import (
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
const createEvent = "create--"
const addEvent = "add-----"
const withdrawEvent = "withdraw"

func encodeLine(name, eventType, number string) (string, error) {
	if len(name) > 10 {
		return "", errors.New("name too long")
	}
	namePadded := name + strings.Repeat(namePadding, 10-len(name))
	numberPadded := strings.Repeat(namePadding, 10-len(number)) + number
	line := namePadded + splitToken + eventType + splitToken + numberPadded
	return line, nil
}

func decodeLine(raw string) (string, string, int, error) {
	splitRaw := strings.Split(raw, splitToken)
	name := strings.Trim(splitRaw[0], namePadding)
	eventType := splitRaw[1]
	amountString := strings.Trim(splitRaw[2], namePadding)
	amount, err := strconv.Atoi(amountString)
	if err != nil {
		return "", "", 0, fmt.Errorf("could not decode line: %w", err)
	}
	return name, eventType, amount, nil
}

func appendAmountToData(raw string) error {
	var err error
	file, err := os.OpenFile(dataPath, os.O_APPEND|os.O_WRONLY, 0644)
	defer file.Close()
	// need to read what is the best practise for handling errors
	if err != nil {
		return fmt.Errorf("unable to open file: %w", err)
	}
	_, err = file.WriteString(raw + "\n")
	if err != nil {
		return fmt.Errorf("unable to write string: %w", err)
	}
	return nil
}

func subtractMoneyEvent(projection keyValueStore, name string, subtractAmount int) error {
	currentAmount, ok := projection.get(name)
	if !ok {
		return errors.New("name not in db")
	}
	updatedAmount := currentAmount - subtractAmount
	if updatedAmount < 0 {
		return errors.New("not enough money")
	}
	subtractAmountString := strconv.Itoa(subtractAmount)
	raw, err := encodeLine(name, withdrawEvent, subtractAmountString)
	if err != nil {
		return err
	}
	projection.add(name, updatedAmount)
	err = appendAmountToData(raw)
	if err != nil {
		projection.add(name, currentAmount)
		return fmt.Errorf("unable to persist data to file: %w", err)
	}
	return nil
}

func addMoneyEvent(projection keyValueStore, name string, addAmount int) error {
	currentAmount, ok := projection.get(name)
	if !ok {
		return errors.New("account does not exist")
	}
	if addAmount > math.MaxInt32 {
		return errors.New("the value is too large")
	}
	amountString := strconv.Itoa(addAmount)
	raw, err := encodeLine(name, addEvent, amountString)
	if err != nil {
		return err
	}
	updatedAmount := currentAmount + addAmount
	projection.add(name, updatedAmount)
	err = appendAmountToData(raw)
	if err != nil {
		// roll back the update
		projection.add(name, currentAmount)
		return fmt.Errorf("unable to persist data to file: %w", err)
	}
	return nil
}

func createAccountEvent(projection keyValueStore, name string) error {
	_, ok := projection.get(name)
	if ok {
		return fmt.Errorf("Account already exists")
	}
	startingValue := 0
	raw, err := encodeLine(name, createEvent, strconv.Itoa(startingValue))
	if err != nil {
		return err
	}
	projection.add(name, startingValue)
	err = appendAmountToData(raw)
	if err != nil {
		return fmt.Errorf("failed to append to file: %w", err)
	}
	return nil
}
