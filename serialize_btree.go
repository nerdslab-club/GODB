package main

import (
	"encoding/json"
	"github.com/petar/GoLLRB/llrb"
	"io/ioutil"
)

func SaveSerializedData(data []byte, filename string) error {
	// Write the serialized data to a file.
	err := ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
func lessInt(a, b interface{}) bool { return a.(int) < b.(int) }
func main() {
	// Assume that "tree" is an instance of the GoLLRB/llrb B-tree.
	var data []byte
	iter := tree.IterAscend()
	for {
		item := iter.Next()
		if item == nil {
			break
		}
		key := item.Key.(string)
		value := item.Value.(string)
		nodeData := map[string]string{"key": key, "value": value}
		serializedNode, _ := json.Marshal(nodeData)
		data = append(data, serializedNode...)
	}

	// Now "data" contains the serialized B-tree data.

	// Assume that "serializedData" is a byte array containing the serialized B-tree data.
	tree := llrb.New()
	var nodeData []map[string]string
	json.Unmarshal(data, &nodeData)
	for _, data := range nodeData {
		key := data["key"]
		value := data["value"]
		tree.ReplaceOrInsert(llrb.Int(key), value)
	}

	// Now "tree" contains the deserialized B-tree.

}
