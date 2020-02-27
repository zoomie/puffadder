package main

import (
	"testing"
)

func getInputValues() map[string]int {
	testCases := map[string]int{
		"key1":    1,
		"key2":    2,
		"key3":    3,
		"newType": 1000,
		"aaaaaa":  2121,
		"z":       100000000000000,
	}
	return testCases
}

func TestAlgorithms(t *testing.T) {
	testCases := getInputValues()
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
