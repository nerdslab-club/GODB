package db_operations

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type TableIndex struct {
	Table string
	Index []string
	PK    string
}

func (d *Driver) WriteIndex(collection, resource string, tableExists bool, v interface{}) error {
	if collection == "" {
		return fmt.Errorf("Missing collection - no place to save record!")
	}

	if resource == "" {
		return fmt.Errorf("Missing resource - unable to save record (no name)!")
	}

	dir := filepath.Join(d.dir, collection)
	fnlPath := filepath.Join(dir, resource+".json")

	_, err := os.Stat(fnlPath)

	/*if the file already exists*/
	if !os.IsNotExist(err) {
		if !tableExists {
			file, err := ioutil.ReadFile("database/table_index.json")
			if err != nil {
				fmt.Println("Error: ", err)
				return err
			}

			var tableArr []TableIndex
			err = json.Unmarshal(file, &tableArr)
			if err != nil {
				fmt.Println("Error: ", err)
				return err
			}

			v = append(tableArr, v.([]TableIndex)...)
		}

		err = os.Remove("database/table_index.json")
		if err != nil {
			return err
		}
	}

	tmpPath := fnlPath + ".tmp"

	mkdir := d.MakeDirectory(dir)
	if mkdir == nil {

		/*converts string to JSON*/
		b, err := json.MarshalIndent(v, "", "\t")
		if err != nil {
			return err
		}

		b = append(b, byte('\n'))

		if err := ioutil.WriteFile(tmpPath, b, 0755); err != nil {
			return err
		}

		return os.Rename(tmpPath, fnlPath)
	} else {
		return mkdir
	}
}

// CreateIndex Creates table_index file and updates
func CreateIndex(tableName, columnName string) string {

	/*Create new directory at root*/
	db, err := New(root, nil)
	if err != nil {
		fmt.Println("Error", err)
		return "Failed!"
	}

	dir := filepath.Join(root, tableName)

	// check if the table exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		fmt.Println("Table does not exist!")
		return "Failed!"
	}

	// Read the contents of "table_index.json" into a slice of TableIndex objects
	var tableIndexes []TableIndex
	err = db.Read("/", "table_index", &tableIndexes)
	if err != nil {
		fmt.Println("Error reading table index:", err)
	}

	//var indexes []string
	// Iterates over the table_index data to see if the desired table exists
	tableExists := false
	for i, tIndex := range tableIndexes {
		if tIndex.Table == tableName {
			tableExists = true
			// Iterates over the index array and check if the index already exists
			valueExists := false
			for _, element := range tIndex.Index {
				if element == columnName {
					valueExists = true
					break
				}
			}
			if valueExists {
				fmt.Println("Column already indexed!")
				return ""
			}
			tIndex.Index = append(tIndex.Index, columnName)
			tableIndexes[i] = tIndex
			break
		}
	}

	// Print the result
	if tableExists {
		/*creates table_index.json and updates*/
		err = db.WriteIndex("/", "table_index", tableExists, tableIndexes)
	} else {
		/*creates table_index.json and updates*/
		err = db.WriteIndex("/", "table_index", tableExists, []TableIndex{
			{
				Table: tableName,
				Index: []string{columnName},
			},
		})
	}

	if err != nil {
		fmt.Println("Error ", err)
		return "Failed!"
	}

	return "Successfully Index Created!"
}

//func createBtree(tableName, columnName string) *b_tree.BTree {
//	tree := b_tree.NewBTree()
//	tree.Insert()
//}
