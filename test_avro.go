package main

import (
	"fmt"
	"github.com/linkedin/goavro/v2"
	"io/ioutil"
)

func main() {
	schema := `{"type":"record","name":"example","fields":[{"name":"name","type":"string"},{"name":"age","type":"int"}]}`

	codec, err := goavro.NewCodec(schema)
	if err != nil {
		fmt.Println("Error creating Avro codec:", err)
		return
	}

	// Serialize the data
	data := map[string]interface{}{
		"name": "John Doe",
		"age":  30,
	}

	binaryData, err := codec.BinaryFromNative(nil, data)
	if err != nil {
		fmt.Println("Error serializing Avro data:", err)
		return
	}

	// Save serialized data to a file
	err = ioutil.WriteFile("serialized.avro", binaryData, 0644)
	if err != nil {
		fmt.Println("Error saving serialized data to file:", err)
		return
	}

	// Read serialized data from the file
	deserializedData, err := ioutil.ReadFile("serialized.avro")
	if err != nil {
		fmt.Println("Error reading serialized data from file:", err)
		return
	}

	// Deserialize the data
	nativeData, _, err := codec.NativeFromBinary(deserializedData)
	if err != nil {
		fmt.Println("Error deserializing Avro data:", err)
		return
	}

	fmt.Println("Deserialized data:", nativeData)
}
