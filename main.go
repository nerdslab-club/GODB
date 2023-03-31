package main

import (
	"bufio"
	"fmt"
	"github.com/petar/GoLLRB/llrb"
	"os"
)

func getInput() string {
	fmt.Print("Enter some input: ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return input
}

type Index struct {
	Key string
	Val int
}

func (i Index) Less(other llrb.Item) bool {
	return i.Key < other.(Index).Key
}

func main() {
	//input := getInput()
	//fmt.Println("You entered:", input)
	tree := llrb.New()
	tree.ReplaceOrInsert(Index{"string90", 1})
	tree.ReplaceOrInsert(Index{"string2", 2})
	tree.ReplaceOrInsert(Index{"string65", 3})
	tree.ReplaceOrInsert(Index{"string4", 4})

	fmt.Println("Inorder traversal:")
	tree.AscendGreaterOrEqual(tree.Min(), func(i llrb.Item) bool {
		index := i.(Index)
		fmt.Printf("Key: %s, Val: %d\n", index.Key, index.Val)
		return true
	})
}
