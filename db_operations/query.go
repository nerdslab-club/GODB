package db_operations

import (
	"encoding/json"
	"fmt"
	"goDB/b_tree"
)

func QueryIndex(tableName, columnName, condition, value string) []string {
	treeName := b_tree.CreateBtreeName(tableName, columnName)
	currentBTree := b_tree.GetIndex(treeName)
	if currentBTree == nil {
		fmt.Println("No tree found!")
		return nil
	}

	fileNames, err := currentBTree.Search(b_tree.ToString(value), condition)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}

	var searchResults []string
	for _, file := range fileNames {
		searchResults = append(searchResults, readFile(file, tableName))
	}

	return searchResults
}

func readFile(pk string, tableName string) string {

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

	return jsonStr
}

func StringInCondition(str string) bool {
	list := []string{">", "<", "==", ">=", "<="}
	for _, item := range list {
		if item == str {
			return true
		}
	}
	return false
}
