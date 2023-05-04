package b_tree

import (
	"fmt"
	"sort"
)

const M = 4 // Maximum degree of the B-Tree

type Key int // Type of keys stored in the B-Tree

type Node struct {
	IsLeaf    bool
	Keys      []Key
	FileNames []string
	Children  []*Node
}

func NewNode(isLeaf bool) *Node {
	return &Node{
		IsLeaf:    isLeaf,
		Keys:      make([]Key, 0, M-1),
		FileNames: make([]string, 0, M-1),
		Children:  make([]*Node, 0, M),
	}
}

type BTree struct {
	Root *Node
}

func NewBTree() *BTree {
	return &BTree{Root: NewNode(true)}
}

func (t *BTree) Insert(key Key, fileName string) {
	if len(t.Root.Keys) == M-1 {
		newRoot := NewNode(false)
		newRoot.Children = append(newRoot.Children, t.Root)
		t.Root = newRoot
		t.splitChild(0, t.Root)
	}
	t.insertNonFull(key, fileName, t.Root)
}

func (t *BTree) insertNonFull(key Key, fileName string, node *Node) {
	i := len(node.Keys) - 1
	if node.IsLeaf {
		node.Keys = append(node.Keys, 0)
		node.FileNames = append(node.FileNames, "")
		for j := i; j >= 0 && key < node.Keys[j]; j-- {
			node.Keys[j+1] = node.Keys[j]
			node.FileNames[j+1] = node.FileNames[j]
		}
		node.Keys[i+1] = key
		node.FileNames[i+1] = fileName
	} else {
		for i >= 0 && key < node.Keys[i] {
			i--
		}
		i++
		if len(node.Children[i].Keys) == M-1 {
			t.splitChild(i, node)
			if key > node.Keys[i] {
				i++
			}
		}
		t.insertNonFull(key, fileName, node.Children[i])
	}
}

func (t *BTree) splitChild(i int, node *Node) {
	y := node.Children[i]
	z := NewNode(y.IsLeaf)
	z.Keys = append(z.Keys, y.Keys[M/2:M]...)
	z.FileNames = append(z.FileNames, y.FileNames[M/2:M]...)
	y.Keys = y.Keys[0 : M/2]
	y.FileNames = y.FileNames[0 : M/2]
	if !y.IsLeaf {
		z.Children = append(z.Children, y.Children[M/2+1:]...)
		y.Children = y.Children[0 : M/2+1]
	}
	node.Children = append(node.Children, nil)
	copy(node.Children[i+2:], node.Children[i+1:])
	node.Children[i+1] = z
	node.Keys = append(node.Keys, 0)
	node.FileNames = append(node.FileNames, "")
	copy(node.Keys[i+1:], node.Keys[i:])
	copy(node.FileNames[i+1:], node.FileNames[i:])
	node.Keys[i] = y.Keys[M/2]
	node.FileNames[i] = y.FileNames[M/2]
}

func (t *BTree) Search(key Key, op string) ([]string, error) {
	switch op {
	case "==":
		return t.searchEqual(key, t.Root), nil
	case ">":
		return t.searchGreaterThan(key, t.Root), nil
	case "<":
		return t.searchLessThan(key, t.Root), nil
	case ">=":
		return t.searchGreaterThanOrEqual(key, t.Root), nil
	case "<=":
		return t.searchLessThanOrEqual(key, t.Root), nil
	default:
		return nil, fmt.Errorf("invalid operator: %s", op)
	}
}

func (t *BTree) searchEqual(key Key, node *Node) []string {
	i := sort.Search(len(node.Keys), func(i int) bool { return node.Keys[i] >= key })
	if i < len(node.Keys) && node.Keys[i] == key {
		return []string{node.FileNames[i]}
	}
	if node.IsLeaf {
		return nil
	}
	return t.searchEqual(key, node.Children[i])
}

func (t *BTree) searchGreaterThan(key Key, node *Node) []string {
	i := sort.Search(len(node.Keys), func(i int) bool { return node.Keys[i] >= key })
	if node.IsLeaf {
		return node.FileNames[i:]
	}
	if i == len(node.Keys) {
		return t.searchGreaterThan(key, node.Children[i])
	}
	return append(node.FileNames[i:], t.searchGreaterThan(key, node.Children[i])...)
}

func (t *BTree) searchLessThan(key Key, node *Node) []string {
	i := sort.Search(len(node.Keys), func(i int) bool { return node.Keys[i] >= key })
	if node.IsLeaf {
		return node.FileNames[0:i]
	}
	if i == 0 {
		return t.searchLessThan(key, node.Children[i])
	}
	return append(t.searchLessThan(key, node.Children[i-1]), node.FileNames[0:i]...)
}

func (t *BTree) searchGreaterThanOrEqual(key Key, node *Node) []string {
	i := sort.Search(len(node.Keys), func(i int) bool { return node.Keys[i] >= key })
	if node.IsLeaf {
		return node.FileNames[i:]
	}
	if i == len(node.Keys) {
		return t.searchGreaterThanOrEqual(key, node.Children[i])
	}
	return append(node.FileNames[i:], t.searchGreaterThanOrEqual(key, node.Children[i])...)
}

func (t *BTree) searchLessThanOrEqual(key Key, node *Node) []string {
	i := sort.Search(len(node.Keys), func(i int) bool { return node.Keys[i] >= key })
	if node.IsLeaf {
		return node.FileNames[0:i]
	}
	if i == 0 {
		return t.searchLessThanOrEqual(key, node.Children[i])
	}
	return append(t.searchLessThanOrEqual(key, node.Children[i-1]), node.FileNames[0:i]...)
}
