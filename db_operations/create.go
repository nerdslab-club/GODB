package db_operations

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/jcelliott/lumber"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Address struct {
	City    string
	State   string
	Country string
	Pincode json.Number
}

type User struct {
	Name    string
	Age     json.Number
	Contact string
	Company string
	Address Address
}

type TablePk struct {
	Name       string
	PrimaryKey string
}

var root string = "database"

type (
	Logger interface {
		Fatal(string, ...interface{})
		Error(string, ...interface{})
		Warn(string, ...interface{})
		Info(string, ...interface{})
		Debug(string, ...interface{})
		Trace(string, ...interface{})
	}

	Driver struct {
		mutex   sync.Mutex
		mutexes map[string]*sync.Mutex
		dir     string
		log     Logger
	}
)

type Options struct {
	Logger
}

func New(dir string, options *Options) (*Driver, error) {
	dir = filepath.Clean(dir)

	opts := Options{}

	if options != nil {
		opts = *options
	}

	if opts.Logger == nil {
		opts.Logger = lumber.NewConsoleLogger((lumber.INFO))
	}

	driver := Driver{
		dir:     dir,
		mutexes: make(map[string]*sync.Mutex),
		log:     opts.Logger,
	}

	if _, err := os.Stat(dir); err == nil {
		opts.Logger.Debug("Using '%s' (database already exists)\n", dir)
		return &driver, nil
	}

	opts.Logger.Debug("Creating the database at '%s'...\n", dir)
	return &driver, os.MkdirAll(dir, 0755)
}

func (d *Driver) Write(collection, resource string, v interface{}) error {

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
		file, err := ioutil.ReadFile("database/table_pk.json")
		if err != nil {
			fmt.Println("Error 223:", err)
			return err
		}

		var tableArr []TablePk
		err = json.Unmarshal(file, &tableArr)
		if err != nil {
			fmt.Println("Error:", err)
			return err
		}

		v = append(tableArr, v.([]TablePk)...)
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

		if err := ioutil.WriteFile(tmpPath, b, 0777); err != nil {
			return err
		}

		return os.Rename(tmpPath, fnlPath)
	} else {
		return mkdir
	}
}

func (d *Driver) MakeDirectory(dir string) error {

	mutex := d.getOrCreateMutex(dir)
	mutex.Lock()
	defer mutex.Unlock()

	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return nil
}

func (d *Driver) getOrCreateMutex(collection string) *sync.Mutex {

	d.mutex.Lock()
	defer d.mutex.Unlock()
	m, ok := d.mutexes[collection]

	if !ok {
		m = &sync.Mutex{}
		d.mutexes[collection] = m
	}

	return m
}

// Create inserts new row data
func Create() string {
	dir := "database/"

	fmt.Println("Write table name where you want to insert data or to see table list type 'table list'")
	var tableName string
	isTableName := false

	for isTableName == false {

		input := GetInput()

		if WordCount(input) == 2 && (input == "table list" || input == "TABLE LIST") { //user wants to see table list
			var count int
			// Use filepath.Walk to traverse the directory tree
			err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
				// Check if the current path is a directory
				if info.IsDir() {
					count++
					// Skip the first directory
					if count == 1 {
						return nil
					}
					fmt.Print(count-1, ". ", info.Name(), "\n")
				}
				return nil
			})
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Write table name where you want to insert data or to see table list type 'table list'")
		} else if WordCount(input) == 1 { //Table name entered by user
			//Checks if the table exists
			// Use filepath.Walk to traverse the directory tree
			err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
				// Check if the current path is a directory with the desired name
				if info.IsDir() && info.Name() == input {
					isTableName = true
					tableName = input
				}
				return nil
			})
			if err != nil {
				fmt.Println(err)
			}

			if isTableName == false {
				fmt.Print("Table not found. Please, enter a valid table name!\n")
			}
		} else {
			fmt.Print("Invalid input. Please, enter correct command!\n")
		}
	}

	/*Create new directory at root*/
	db, err := New(dir, nil)
	if err != nil {
		fmt.Println("Error", err)
	}

	/*For random json structures*/
	var data map[string]interface{}
	//{"Name": "Prat","age": 30,"isEmployed": true,"contact": {"email": "johndoe@example.com", "phone": "+1 555-555-5555"}, "hobbies": ["reading", "playing video games", "hiking"]}
	var pkValue string
	validInput := false
	for validInput == false {
		fmt.Println("Please enter valid json file containing table's primary key to insert new row.")
		/*Taking input from user*/
		input := GetInput()
		//{"Name": "John", "Age": 32, "Contact": "23344333", "Company": "Dominate", "Address": {"City": "bangalore", "State": "karnataka", "Country": "india", "Pincode": 410013}}

		/*Decode a JSON string to GO value*/
		err = json.Unmarshal([]byte(input), &data)
		if err != nil {
			//return errors.New("invalid json body provided for the request")
			fmt.Println("Error: Invalid json body provided for the request!")
			return "Failed!"
		}

		/*Get chosen table's primary key*/
		tablePk := GetPrimaryKey(tableName)

		pkValue = CheckPKValue(data, tablePk)

		/*Checks if file exists with same name*/
		filename := filepath.Join(dir, tableName, pkValue+".json")

		if _, err := os.Stat(filename); os.IsNotExist(err) {
			if pkValue == "" {
				fmt.Println("Primary key is absent in json file!")
			} else {
				validInput = true
			}
		} else {
			fmt.Println("Row with same primary key already exists. Please, give valid data!")
		}

	}

	/*test json data*/
	//{"Name": "Prattoy", "Age": 32, "Contact": "23344333", "Company": "Dominate", "Address": {"City": "bangalore", "State": "karnataka", "Country": "india", "Pincode": 410013}}

	if pkValue == "" || data == nil {
		return "Error occurred!"
	}

	db.Write(tableName, pkValue, data)
	return "Created!"
}

// CreateTablePk Creates table directory and updates table_pk file
func CreateTablePk(table string, pk string) string {

	/*Create new directory at root*/
	db, err := New(root, nil)
	if err != nil {
		fmt.Println("Error", err)
	}

	dir := filepath.Join(root, table)

	/*make table*/
	err = db.MakeDirectory(dir)
	if err != nil {
		fmt.Println("Error", err)
		return "Failed!"
	}

	/*creates table_pk.json and updates*/
	err = db.Write("/", "table_pk", []TablePk{
		{
			Name:       table,
			PrimaryKey: pk,
		},
	})
	if err != nil {
		fmt.Println("Error 123", err)
		return "Failed!"
	}

	return "Successfully table created!"
}

// GetInput gets input from user
func GetInput() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

// WordCount counts total words
func WordCount(s string) int {
	return len(strings.Fields(s))
}

// GetPrimaryKey get tables primary key
func GetPrimaryKey(tableName string) string {
	var pk string
	fmt.Print(pk)
	data, err := ioutil.ReadFile("database/table_pk.json")
	if err != nil {
		log.Fatal(err)
	}

	// Unmarshal the JSON data into a []map[string]interface{}
	var objmaps []map[string]interface{}
	err = json.Unmarshal(data, &objmaps)
	if err != nil {
		log.Fatal(err)
	}

	// Loop through each object in the array and check if the key "Name" exists
	for _, objmap := range objmaps {
		if name, ok := objmap["Name"]; ok {
			if name == tableName {
				pk = objmap["PrimaryKey"].(string)
			}
		}
	}

	return pk
}

// CheckPKValue if primary key exists then pk value is returned
func CheckPKValue(data map[string]interface{}, tablePk string) string {
	// Check if the key "tablePk" exists in the map
	if _, ok := data[tablePk]; ok {
		return data[tablePk].(string)
	}

	return ""
}
