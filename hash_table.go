package main

import (
	"hash/fnv"
)

const startingLength = 8

// This means the fraction of fill slots the dictionary will
// create a new array and copy the contents over.
const capacityToExpand = 0.7

// The rate that the new dictionary will grow by.
const growthMultiplication = 2

type htuple struct {
	key   string
	value int
}

type valuesArray [][]htuple

type hashTable struct {
	length        int
	currentNumber int
	array         valuesArray
}

func hashString(s string, lenght int) int {
	h := fnv.New64a()
	h.Write([]byte(s))
	largeValue := int(h.Sum64())
	index := largeValue % lenght
	if index < 0 {
		return -index
	}
	return index
}

func addToHashMap(array valuesArray, length int, tup htuple) {
	// can pass array in a a direct values becuase it is a slice and will be
	// updated outside of this function scope.
	index := hashString(tup.key, length)
	// This loop is used if the key already exists and there is
	// a collision at the index.
	for i := range array[index] {
		if array[index][i].key == tup.key {
			array[index][i].value = tup.value
			return
		}
	}
	array[index] = append(array[index], tup)
}

func rebuildHashTable(array valuesArray, newLength int) valuesArray {
	newArray := make(valuesArray, newLength)
	for _, collisionArray := range array {
		if collisionArray == nil {
			continue
		}
		for _, tuple := range collisionArray {
			addToHashMap(newArray, newLength, tuple)
		}
	}
	return newArray
}

func (h *hashTable) add(key string, value int) {
	if h.array == nil {
		// need to make the valuesArray as the dafault values is nil
		h.array = make(valuesArray, startingLength)
		h.length = startingLength
	}
	// Increase size of the array once a certian capacity has been reached.
	fractionFull := int(float64(h.length) * capacityToExpand)
	if h.currentNumber > fractionFull {
		h.length = h.length * growthMultiplication
		h.array = rebuildHashTable(h.array, h.length)
	}
	h.currentNumber++
	newTuple := htuple{key: key, value: value}
	addToHashMap(h.array, h.length, newTuple)
}

func (h *hashTable) get(key string) (int, bool) {
	if h.array == nil {
		return 0, false
	}
	index := hashString(key, h.length)
	for _, tup := range h.array[index] {
		if tup.key == key {
			return tup.value, true
		}
	}
	return 0, false
}
