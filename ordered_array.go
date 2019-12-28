package main

import (
	"math"
)

type tuple struct {
	key   string
	value int64
}

type orderedArray []tuple

func (array *orderedArray) add(key string, value int64) {
	newTup := tuple{key: key, value: value}
	arrayLength := len(*array)
	if arrayLength == 0 {
		*array = orderedArray{newTup}
		return
	}
	// at the start
	if key < (*array)[0].key {
		newArray := make(orderedArray, arrayLength+1)
		newArray[0] = newTup
		copy(newArray[1:], *array)
		*array = newArray
		return
	}
	for i := range *array {
		if i == 0 {
			continue
		}
		if key == (*array)[i].key {
			(*array)[i] = newTup
			return
		} else if key > (*array)[i-1].key && key < (*array)[i].key {
			start := make(orderedArray, i)
			copy(start, *array)
			end := (*array)[i:]
			*array = append(append(start, newTup), end...)
			return
		}
	}
	// add at end
	if key > (*array)[arrayLength-1].key {
		*array = append(*array, newTup)
		return
	}
}

func binarySearchGet(array orderedArray, key string) (int64, bool) {
	if len(array) == 0 {
		return 0, false
	}
	if len(array) == 1 {
		if array[0].key == key {
			return array[0].value, true
		}
		return 0, false
	}
	middleIndex := int(math.Floor(float64(len(array)) / 2))
	if key > array[middleIndex].key {
		top := array[middleIndex:] // could add middleIndex+1
		return binarySearchGet(top, key)
	} else if key < array[middleIndex].key {
		bottom := array[:middleIndex]
		return binarySearchGet(bottom, key)
	} else {
		return array[middleIndex].value, true
	}
}

func (array *orderedArray) get(key string) (int64, bool) {
	return binarySearchGet(*array, key)
}

// func main() {
// 	var o orderedArray
// 	o.add("b", 6)
// 	o.add("a", 5)
// 	o.add("d", 4)
// 	o.add("c", 1)
// 	o.add("c", 2)
// 	value, _ := o.get("a")
// 	fmt.Println(value)
// }
