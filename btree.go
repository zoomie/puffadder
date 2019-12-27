package main

type btree struct {
	start *node
}

type node struct {
	key   string
	value int64
	left  *node
	right *node
}

func addNode(baseNode, newNode *node) {
	if baseNode.key < newNode.key {
		if baseNode.right == nil {
			baseNode.right = newNode
		} else {
			addNode(baseNode.right, newNode)
		}
	} else if baseNode.key > newNode.key {
		if baseNode.left == nil {
			baseNode.left = newNode
		} else {
			addNode(baseNode.left, newNode)
		}
	} else {
		baseNode.value = newNode.value
	}
}

func (b *btree) add(key string, value int64) {
	newNode := &node{key: key, value: value}
	if b.start == nil {
		b.start = newNode
	} else {
		addNode(b.start, newNode)
	}
}

func getValue(baseNode *node, key string) (int64, bool) {
	if baseNode.key < key {
		if baseNode.right == nil {
			return 0, false
		}
		return getValue(baseNode.right, key)
	} else if baseNode.key > key {
		if baseNode.left == nil {
			return 0, false
		}
		return getValue(baseNode.left, key)
	} else {
		return baseNode.value, true
	}
}

func (b *btree) get(key string) (int64, bool) {
	if b.start == nil {
		return 0, false
	}
	return getValue(b.start, key)
}
