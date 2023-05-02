package db_operations

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func (d *Driver) ReadAll(collection string) ([]string, error) {

	if collection == "" {
		return nil, fmt.Errorf("Missing collection - unable to read")
	}
	dir := filepath.Join(d.dir, collection)

	if _, err := stat(dir); err != nil {
		return nil, err
	}

	files, _ := ioutil.ReadDir(dir)

	var records []string

	for _, file := range files {
		b, err := ioutil.ReadFile(filepath.Join(dir, file.Name()))
		if err != nil {
			return nil, err
		}

		records = append(records, string(b))
	}
	return records, nil
}

func (d *Driver) Read(collection, resource string, v interface{}) error {

	if collection == "" {
		return fmt.Errorf("Missing collection - unable to read!")
	}

	if resource == "" {
		return fmt.Errorf("Missing resource - unable to read record (no name)!")
	}

	record := filepath.Join(d.dir, collection, resource)

	if _, err := stat(record); err != nil {
		return err
	}

	b, err := ioutil.ReadFile(record + ".json")
	if err != nil {
		return err
	}

	return json.Unmarshal(b, &v)
}

func stat(path string) (fi os.FileInfo, err error) {
	if fi, err = os.Stat(path); os.IsNotExist(err) {
		fi, err = os.Stat(path + ".json")
	}
	return
}

func ReadAll(tableName string) (string, string) {

	/*If directory does not exist it returns*/
	if !CheckDirectory(tableName) {
		return "", ""
	}

	db, err := New(root, nil)
	if err != nil {
		fmt.Println("Error ", err)
		return "", ""
	}
	/*Read all users*/
	records, err := db.ReadAll(tableName)
	if err != nil {
		fmt.Println("Error", err)
		return "", ""
	}
	//fmt.Println(records)

	/*For any type of JSON structures*/
	var allData []map[string]interface{}

	for _, f := range records {
		var data map[string]interface{}

		/*Decodes JSON*/
		if err := json.Unmarshal([]byte(f), &data); err != nil {
			fmt.Println("Error", err)
		}
		allData = append(allData, data)
	}

	jsonBytes, err := json.MarshalIndent(allData, "", "\t")
	if err != nil {
		fmt.Println("Error", err)
		return "", ""
	}

	jsonStr := string(jsonBytes)
	fmt.Print(jsonStr)

	return "\nSuccessfully read from directory!", jsonStr
}

func Read(pk string, tableName string) string {

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

	var data map[string]interface{}

	err = db.Read(tableName, pk, &data)
	if err != nil {
		fmt.Println("Error: ", err)
		return ""
	}

	rowData, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		fmt.Println("Error", err)
		return ""
	}

	jsonStr := string(rowData)
	fmt.Print(jsonStr)

	return "\nSuccessfully read from file!"
}

// CheckDirectory checks if directory exists
func CheckDirectory(directoryName string) bool {
	if _, err := os.Stat(filepath.Join(root, directoryName)); os.IsNotExist(err) {
		fmt.Println("Sorry, no table found!")
		return false
	}
	return true
}

// CheckFile checks if file exists
func CheckFile(fileName string, directoryName string) bool {
	if _, err := os.Stat(filepath.Join(root, directoryName, fileName+".json")); os.IsNotExist(err) {
		fmt.Println("Sorry, no row found!")
		return false
	}
	return true
}
