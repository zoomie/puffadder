package main

type tuple struct {
	key   string
	value int64
}

type orderedArray []tuple

func sort(o *orderedArray) {
	// implement quick sort here
}

func (o *orderedArray) add(key string, value int64) {
	tup := tuple{key: key, value: value}
	if o == nil {
		*o = []tuple{tup}
	} else {
		*o = append(*o, tup)
	}
	sort(o)
}

func (o *orderedArray) get(key string) (int64, bool) {
	for _, tup := range *o {
		if tup.key == key {
			return tup.value, true
		}
	}
	return 0, false
}

// func main() {
// 	var o orderedArray
// 	o.add("zoomie", 432)
// 	t, _ := o.get("zoomie")
// 	fmt.Println(t)
// }
