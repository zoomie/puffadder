package main

import (
	"testing"
)

func testHashTable(t *testing.T) {
	key1 := "key1"
	var value1 int64 = 1
	key2 := "key2"
	var value2 int64 = 2
	h := hashTable{}
	h.add(key1, value1)
	h.add(key2, value2)
	output1, _ := h.get(key1)
	output2, _ := h.get(key2)
	if output1 != value1 {
		t.Errorf("hash failed")
	}
	if output2 != value2 {
		t.Errorf("hash failed")
	}
}
