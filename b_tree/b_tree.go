package b_tree

import (
	"fmt"
	"github.com/google/btree"
)

type Node struct {
	Key      string
	FileName string
}

func (n Node) Less(than btree.Item) bool {
	return n.Key < than.(Node).Key
}

type BTree struct {
	tree *btree.BTree
}

func NewBTree(degree int) *BTree {
	return &BTree{
		tree: btree.New(degree),
	}
}

func (t *BTree) Insert(key string, fileName string) {
	node := Node{
		Key:      key,
		FileName: fileName,
	}
	t.tree.ReplaceOrInsert(node)
}

func (t *BTree) SearchEqual(key string) string {
	currentItem := t.tree.Get(Node{Key: key})
	if currentItem != nil {
		return currentItem.(Node).FileName
	} else {
		return ""
	}
}

func (t *BTree) SearchGreater(key string) []string {
	results := make([]string, 0)
	t.tree.AscendGreaterOrEqual(Node{Key: key}, func(i btree.Item) bool {
		if i.(Node).Key > key {
			results = append(results, i.(Node).FileName)
			return true
		}
		return true
	})
	return results
}

func (t *BTree) SearchLess(key string) []string {
	results := make([]string, 0)
	t.tree.AscendLessThan(Node{Key: key}, func(i btree.Item) bool {
		results = append(results, i.(Node).FileName)
		return true
	})
	return results
}

func (t *BTree) SearchGreaterOrEqual(key string) []string {
	results := make([]string, 0)
	t.tree.AscendGreaterOrEqual(Node{Key: key}, func(i btree.Item) bool {
		results = append(results, i.(Node).FileName)
		return true
	})
	return results
}

func (t *BTree) SearchLessOrEqual(key string) []string {
	results := t.SearchLess(key)
	equalResult := t.SearchEqual(key)
	if equalResult != "" {
		results = append(results, t.SearchEqual(key))
	}
	return results
}

func (t *BTree) Search(key string, op string) ([]string, error) {
	switch op {
	case "==":
		currentItem := t.SearchEqual(key)
		currentItemArr := make([]string, 0)
		if currentItem != "" {
			currentItemArr = append(currentItemArr, currentItem)
		}
		return currentItemArr, nil
	case ">":
		return t.SearchGreater(key), nil
	case "<":
		return t.SearchLess(key), nil
	case ">=":
		return t.SearchGreaterOrEqual(key), nil
	case "<=":
		return t.SearchLessOrEqual(key), nil
	default:
		return nil, fmt.Errorf("invalid operator: %s", op)
	}
}
