package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"goDB/db_operations"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

const Version = "1.0.0"

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

func stat(path string) (fi os.FileInfo, err error) {
	if fi, err = os.Stat(path); os.IsNotExist(err) {
		fi, err = os.Stat(path + ".json")
	}
	return
}

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

func main() {
	//dir := "database/"

	var input int = 1

	for input != 0 {
		/*Choose operation*/
		fmt.Print("1. Create (Press 1)\n2. Read (Press 2)\n3. Update (Press 3)\n4. Delete (Press 4)\nTo Exit Press 0\nChoose Any Operations To Continue...\n")

		/*Takes input from user*/
		input = GetInputNumber()

		/*Create*/
		if input == 1 {
			fmt.Print("1. Create Table (Press 1)\n2. Insert Row in Existing Table (Press 2)\nChoose Operation...\n")
			/*Takes input from user*/
			secondInput := GetInputNumber()
			if secondInput == 1 {
				fmt.Print("Write 'CREATE TABLE <TABLE_NAME> <PRIMARY_KEY>'\nFor example, to create a table called 'users' write 'CREATE TABLE USERS NAME'\n")
				thirdInput := GetInputString()

				//fmt.Print(db_operations.Create())
				fmt.Print(db_operations.CreateTablePk(strings.Fields(thirdInput)[2], strings.Fields(thirdInput)[3]))
			}

		} else if input == 2 {
			fmt.Print(input)
		} else if input == 3 {
			fmt.Print(input)
		} else if input == 4 {
			fmt.Print(input)
		}
		fmt.Print("\n")
	}

	//db, err := db_operations.New(dir, nil)
	//if err != nil {
	//	fmt.Println("Error", err)
	//}
	//
	/*Read all users*/
	//records, err := db.ReadAll("users")
	//if err != nil {
	//	fmt.Println("Error", err)
	//}
	//fmt.Println(records)
	//
	//allusers := []User{}
	//
	//for _, f := range records {
	//	employeeFound := User{}
	//
	//	/*Converts string to JSON*/
	//	if err := json.Unmarshal([]byte(f), &employeeFound); err != nil {
	//		fmt.Println("Error", err)
	//	}
	//	allusers = append(allusers, employeeFound)
	//}
	//fmt.Println((allusers))

	/*Delete user by name*/
	//if err := db.Delete("users", "Albert"); err != nil {
	//	fmt.Println("Error", err)
	//}

	/*To delete all users*/
	//if err := db.Delete("users", ""); err != nil {
	//	fmt.Println("Error", err)
	//}
}

// GetInputNumber Takes number input from user
func GetInputNumber() int {
	validInput := false
	var number int
	var err error
	for validInput == false {
		/*Reads input from the user*/
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		input := scanner.Text()

		/*Converts the input to an integer*/
		number, err = strconv.Atoi(input)
		if err != nil {
			fmt.Println("Invalid input. Please enter a number!")
		} else {
			validInput = true
		}
	}
	return number
}

// GetInputString Takes string input from user
func GetInputString() string {
	validInput := false
	var line string
	for validInput == false {
		/*Read input from the user*/
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		line = scanner.Text()

		/*Command validations*/
		if strings.HasPrefix(line, "create table") || strings.HasPrefix(line, "CREATE TABLE") {
			if WordCount(line) != 4 {
				fmt.Println("Invalid command. Please enter valid command!")
			} else {
				validInput = true
			}
		} else {
			fmt.Println("Invalid command. Please enter valid command!")
		}
	}

	return line
}

// WordCount counts total words
func WordCount(s string) int {
	return len(strings.Fields(s))
}
