package main

const initArrLength = 5

type htuple struct {
	key   string
	value int64
}

type valuesArray [][]htuple

type hashTable struct {
	length int
	array  valuesArray
}

func hashString(s string, lenght int) int {
	// need to get sha implementation
	return 0
}

func (h *hashTable) add(key string, value int64) {
	index := hashString(key, h.length)
	if h.array == nil {
		h.array = make(valuesArray, initArrLength)
	}
	for i := range h.array[index] {
		if h.array[index][i].key == key {
			h.array[index][i].value = value
			return
		}
	}
	newTuple := htuple{key: key, value: value}
	h.array[index] = append(h.array[index], newTuple)
}

func (h *hashTable) get(key string) (int64, bool) {
	index := hashString(key, h.length)
	for _, tup := range h.array[index] {
		if tup.key == key {
			return tup.value, true
		}
	}
	return 0, false
}

// func main() {
// 	h := hashTable{}
// 	h.add("key", 543)
// }
