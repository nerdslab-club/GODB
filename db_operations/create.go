package db_operations

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/jcelliott/lumber"
	"io/ioutil"
	"os"
	"path/filepath"
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
	fmt.Print(v)
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
		fmt.Print(tableArr)

		v = append(tableArr, v.([]TablePk)...)
	}

	tmpPath := fnlPath + ".tmp"

	mkdir := d.MakeDirectory(dir)
	if mkdir == nil {

		b, err := json.MarshalIndent(v, "", "\t")
		//b, err := json.MarshalIndent(v, "", "\t")
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

	/*Create new directory at root*/
	db, err := New(dir, nil)
	if err != nil {
		fmt.Println("Error", err)
	}

	/*Taking input from user*/
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()
	//{"Name": "John", "Age": 32, "Contact": "23344333", "Company": "Dominate", "Address": {"City": "bangalore", "State": "karnataka", "Country": "india", "Pincode": 410013}}
	employee := User{}

	/*Converts string to JSON*/
	if err := json.Unmarshal([]byte(input), &employee); err != nil {
		fmt.Println("Error", err)
	}
	//fmt.Print(employee)

	//var data map[string]interface{}
	//err = json.Unmarshal([]byte(line), &data)
	//if err != nil {
	//	return errors.New("invalid json body provided for the request")
	//}
	//{"Name": "Albert", "Age": 32, "Contact": "23344333", "Company": "Dominate", "Address": {"City": "bangalore", "State": "karnataka", "Country": "india", "Pincode": 410013}}

	employees := []User{
		//{"John", "23", "23344333", "Myrl Tech", Address{"bangalore", "karnataka", "india", "410013"}},
		//{"Paul", "25", "23344333", "Google", Address{"san francisco", "california", "USA", "410013"}},
		//{"Robert", "27", "23344333", "Microsoft", Address{"bangalore", "karnataka", "india", "410013"}},
		//{"Vince", "29", "23344333", "Facebook", Address{"bangalore", "karnataka", "india", "410013"}},
		//{"Neo", "31", "23344333", "Remote-Teams", Address{"bangalore", "karnataka", "india", "410013"}},
		//{"Albert", "32", "23344333", "Dominate", Address{"bangalore", "karnataka", "india", "410013"}},
	}

	employees = append(employees, employee)
	//fmt.Print(employees)

	/*Inserts all users in users directory*/
	for _, value := range employees {
		db.Write("users", value.Name, User{
			Name:    value.Name,
			Age:     value.Age,
			Contact: value.Contact,
			Company: value.Company,
			Address: value.Address,
		})
	}
	return "Created"
}

// CreateTablePk Creates table directory and updates table_pk file
func CreateTablePk(table string, pk string) string {
	root := "database/"

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

	tableArray := []TablePk{
		{
			Name:       table,
			PrimaryKey: pk,
		},
	}

	/*creates table_pk.json and updates*/
	err = db.Write("/", "table_pk", tableArray)
	if err != nil {
		fmt.Println("Error 123", err)
		return "Failed!"
	}

	return "Success!"
}
