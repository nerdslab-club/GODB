package db_operations

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func (d *Driver) Update(oldName, newName string) error {

	path := filepath.Join(root, oldName)
	newPath := filepath.Join(root, newName)
	mutex := d.getOrCreateMutex(path)
	mutex.Lock()
	defer mutex.Unlock()

	//fmt.Print(dir)
	switch fi, err := stat(path); {
	case fi == nil, err != nil:
		return fmt.Errorf("unable to find file or directory named %v\n", path)

	case fi.Mode().IsDir():
		err := os.Rename(path, newPath)
		if err != nil {
			return err
		}
	}
	return nil
}

func UpdateTableName(oldName string, newName string) string {

	/*Checks if directory exists*/
	if !CheckDirectory(oldName) {
		return ""
	}

	/*Create new directory at root*/
	db, err := New(filepath.Join(root, oldName), nil)
	if err != nil {
		fmt.Println("Error", err)
	}

	err = db.Update(oldName, newName)
	if err != nil {
		fmt.Println("Error : ", err)
		return ""
	}

	//var pk string
	filePath := filepath.Join(root, "table_pk.json")

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	// Unmarshal the JSON data into a []map[string]interface{}
	var objmaps []map[string]interface{}
	err = json.Unmarshal(data, &objmaps)
	if err != nil {
		log.Fatal(err)
	}

	// Loop through each object in the array and check if the key "Name" exists and if name matches table name then updates table name
	for _, objmap := range objmaps {
		if name, ok := objmap["Name"]; ok {
			if name == oldName {
				objmap["Name"] = newName
			}
		}
	}

	/*converts string to JSON*/
	b, err := json.MarshalIndent(objmaps, "", "\t")
	if err != nil {
		fmt.Println("Error: ", err)
		return ""
	}

	b = append(b, byte('\n'))

	if err := ioutil.WriteFile(filePath, b, 0755); err != nil {
		fmt.Println("Error: ", err)
		return ""
	}

	return "Successfully table name updated!"
}

func UpdateRow(pk, tableName string) string {
	/*Checks if directory exists*/
	if !CheckDirectory(tableName) {
		return ""
	}

	/*Checks if file exists*/
	if !CheckFile(pk, tableName) {
		return ""
	}

	fmt.Println("Please enter valid json file containing table's primary key to insert new row.")
	/*Taking input from user*/
	input := GetInput()

	/*Create new directory at root*/
	db, err := New(root, nil)
	if err != nil {
		fmt.Println("Error", err)
	}

	/*For random json structures*/
	var data map[string]interface{}

	/*Decode a JSON string to GO value*/
	err = json.Unmarshal([]byte(input), &data)
	if err != nil {
		//return errors.New("invalid json body provided for the request")
		fmt.Println("Error: Invalid json body provided for the request!")
		return "Failed!"
	}

	/*Get chosen table's primary key*/
	tablePk := GetPrimaryKey(tableName)

	pkValue := CheckPKValue(data, tablePk)

	/*Delete user by name*/
	if err := db.Delete(tableName, pk); err != nil {
		fmt.Println("Error", err)
	}

	if pkValue == "" || data == nil {
		return "Error occurred!"
	}

	db.Write(tableName, pkValue, data)
	return "Successfully updated row!"
}
