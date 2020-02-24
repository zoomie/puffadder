package main

import "hash/fnv"

const initArrLength = 100

type htuple struct {
	key   string
	value int
}

type valuesArray [][]htuple

type hashTable struct {
	length int
	array  valuesArray
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

func (h *hashTable) add(key string, value int) {
	if h.array == nil {
		h.array = make(valuesArray, initArrLength)
		h.length = initArrLength
	}
	index := hashString(key, h.length)
	for i := range h.array[index] {
		if h.array[index][i].key == key {
			h.array[index][i].value = value
			return
		}
	}
	newTuple := htuple{key: key, value: value}
	h.array[index] = append(h.array[index], newTuple)
}

func (h *hashTable) get(key string) (int, bool) {
	index := hashString(key, h.length)
	for _, tup := range h.array[index] {
		if tup.key == key {
			return tup.value, true
		}
	}
	return 0, false
}
