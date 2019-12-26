package main

type hashTable map[string]int64

func (h hashTable) add(key string, value int64) {
	if h == nil {
		h = make(map[string]int64)
	}
	h[key] = value
}

func (h hashTable) get(key string) (int64, bool) {
	if h == nil {
		return 0, false
	}
	value, ok := h[key]
	if !ok {
		return 0, false
	}
	return value, true
}
