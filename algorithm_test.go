package main

import (
	"strconv"
	"testing"
)

func getInputValues(numberTests int) map[string]int {
	testCases := make(map[string]int)
	num := 0
	for num < numberTests {
		testCases[strconv.Itoa(num)] = num
		num++
	}
	return testCases
}

func TestAlgorithms(t *testing.T) {
	testCases := getInputValues(100)
	algsToTest := []keyValueStore{
		&hashTable{},
		&orderedArray{},
		&btree{},
	}
	for _, algorthim := range algsToTest {
		for key, value := range testCases {
			algorthim.add(key, value)
		}
		for key, expectedValue := range testCases {
			outputValue, _ := algorthim.get(key)
			if expectedValue != outputValue {
				t.Errorf("expectedValue=%d, outputValue=%d", expectedValue, outputValue)
			}
		}
	}
}

// todo: add test cases for failing numbers. ie to large
