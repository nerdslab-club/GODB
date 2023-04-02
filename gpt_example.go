package main

import (
	"encoding/binary"
	"fmt"
	"os"
)

type Object struct {
	ID    int
	Name  string
	Value float64
}

type IndexEntry struct {
	ID     int
	Offset int64
}

func main() {
	objects := []Object{
		{ID: 1, Name: "Object 1", Value: 10.0},
		{ID: 2, Name: "Object 2", Value: 20.0},
		{ID: 3, Name: "Object 3", Value: 30.0},
	}

	// create a binary file to store the objects
	file, err := os.Create("objects.bin")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// write the objects to the file and create an index
	index := make([]IndexEntry, len(objects))
	offset := int64(0)
	for i, obj := range objects {
		// write the object to the file
		err = binary.Write(file, binary.LittleEndian, &obj)
		if err != nil {
			panic(err)
		}

		// create an index entry for the object
		index[i] = IndexEntry{
			ID:     obj.ID,
			Offset: offset,
		}

		// update the offset for the next object
		offset += int64(binary.Size(obj))
	}

	// create a binary file to store the index
	indexFile, err := os.Create("index.bin")
	if err != nil {
		panic(err)
	}
	defer indexFile.Close()

	// write the index to the file
	err = binary.Write(indexFile, binary.LittleEndian, index)
	if err != nil {
		panic(err)
	}

	// read an object by ID from the file
	obj, err := getObjectByID(2, "objects.bin", "index.bin")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Object with ID %d: %+v\n", obj.ID, obj)
}

func getObjectByID(id int, dataFile string, indexFile string) (*Object, error) {
	// open the index file
	index, err := os.Open(indexFile)
	if err != nil {
		return nil, err
	}
	defer index.Close()

	// read the index entries from the file
	var indexEntries []IndexEntry
	err = binary.Read(index, binary.LittleEndian, &indexEntries)
	if err != nil {
		return nil, err
	}

	// find the index entry for the object with the given ID
	var indexEntry IndexEntry
	for _, entry := range indexEntries {
		if entry.ID == id {
			indexEntry = entry
			break
		}
	}
	if indexEntry.ID == 0 {
		return nil, fmt.Errorf("object with ID %d not found", id)
	}

	// open the data file and seek to the offset of the object
	data, err := os.Open(dataFile)
	if err != nil {
		return nil, err
	}
	defer data.Close()
	_, err = data.Seek(indexEntry.Offset, 0)
	if err != nil {
		return nil, err
	}

	// read the object from the file
	var obj Object
	err = binary.Read(data, binary.LittleEndian, &obj)
	if err != nil {
		return nil, err
	}

	return &obj, nil
}
