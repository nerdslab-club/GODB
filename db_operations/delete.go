package db_operations

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func (d *Driver) Delete(collection, resource string) error {

	path := filepath.Join(collection, resource)
	mutex := d.getOrCreateMutex(collection)
	mutex.Lock()
	defer mutex.Unlock()

	dir := filepath.Join(d.dir, path)

	switch fi, err := stat(dir); {
	case fi == nil, err != nil:
		return fmt.Errorf("unable to find file or directory named %v\n", path)

	case fi.Mode().IsDir():
		return os.RemoveAll(dir)

	case fi.Mode().IsRegular():
		return os.RemoveAll(dir + ".json")
	}
	return nil
}

func DeleteRow(tableName, pk string) string {
	/*Checks if directory exists*/
	if !CheckDirectory(tableName) {
		return ""
	}

	/*Checks if file exists*/
	if !CheckFile(pk, tableName) {
		return ""
	}

	db, err := New(root, nil)
	if err != nil {
		fmt.Println("Error ", err)
		return ""
	}
	/*Delete table row*/
	if err := db.Delete(tableName, pk); err != nil {
		fmt.Println("Error", err)
	}
	return "Successfully row deleted!"
}

func DeleteTable(tableName string) string {
	/*Checks if directory exists*/
	if !CheckDirectory(tableName) {
		return ""
	}

	db, err := New(root, nil)
	if err != nil {
		fmt.Println("Error ", err)
		return ""
	}

	/*To delete all users*/
	if err := db.Delete(tableName, ""); err != nil {
		fmt.Println("Error", err)
	}

	data, err := ioutil.ReadFile(root + "/table_pk.json")
	if err != nil {
		log.Fatal(err)
	}

	// Unmarshal the JSON data into a []map[string]interface{}
	var objmaps []map[string]interface{}
	err = json.Unmarshal(data, &objmaps)
	if err != nil {
		log.Fatal(err)
	}

	// Find the index of the object to delete
	var indexToDelete int
	for i, objmap := range objmaps {
		if name, ok := objmap["Name"]; ok {
			if name == tableName {
				indexToDelete = i
				break
			}
		}
	}

	// Remove the object from the array
	objmaps = append(objmaps[:indexToDelete], objmaps[indexToDelete+1:]...)

	/*converts string to JSON*/
	b, err := json.MarshalIndent(objmaps, "", "\t")
	if err != nil {
		fmt.Println("Error: ", err)
		return ""
	}

	b = append(b, byte('\n'))

	if err := ioutil.WriteFile(root+"/table_pk.json", b, 0755); err != nil {
		fmt.Println("Error: ", err)
		return ""
	}
	return "Successfully table deleted!"
}

func DeleteAll() string {
	/*creates db object*/
	db, err := New(root, nil)
	if err != nil {
		fmt.Println("Error ", err)
		return ""
	}

	/*To delete all users*/
	if err := db.Delete("/", ""); err != nil {
		fmt.Println("Error", err)
	}
	return "Successfully deleted all tables!"
}
