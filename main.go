package main

import (
	"bufio"
	"fmt"
	"goDB/db_operations"
	"os"
	"strconv"
	"strings"
)

const Version = "1.0.0"

func main() {

	var input int = 1

	for input != 0 {
		/*Choose operation*/
		fmt.Print("1. Create (Press 1)\n2. Read (Press 2)\n3. Update (Press 3)\n4. Delete (Press 4)\nTo Exit Press 0\nChoose Any Operation To Continue...\n")

		/*Takes input from user*/
		input = GetInputNumber()

		if input == 1 { /*Create*/

			fmt.Print("1. Create Table (Press 1)\n2. Insert Row in Existing Table (Press 2)\nChoose Operation...\n")
			/*Takes input from user*/
			secondInput := GetInputNumber()
			if secondInput == 1 {
				fmt.Print("Write 'CREATE TABLE <TABLE_NAME> <PRIMARY_KEY>'\nFor example, to create a table called 'TEST' write 'CREATE TABLE TEST TEST_NAME'\n")
				thirdInput := GetInputString("create")

				fmt.Print(db_operations.CreateTablePk(strings.Fields(thirdInput)[2], strings.Fields(thirdInput)[3]))
			} else if secondInput == 2 {
				fmt.Print(db_operations.Create())
			} else {
				fmt.Print("Invalid input. Please, give valid input!\n")
			}

		} else if input == 2 { /*Read*/

			fmt.Print("1. Read all from a table (Press 1)\n2. Read specific row data from a table (Press 2)\nTo Exit Press 0\nChoose Any Operation To Continue...\n")
			/*Takes input from user*/
			secondInput := GetInputNumber()

			if secondInput == 1 {
				fmt.Print("Write 'READ * <TABLE_NAME>'\nFor example, to read from 'test' table write 'READ * TEST'\n")
				thirdInput := GetInputString("read")

				if strings.Fields(thirdInput)[1] == "*" {
					fmt.Print(db_operations.ReadAll(strings.Fields(thirdInput)[2]))
				} else {
					fmt.Print("Invalid input. Please, give valid input!\n")
				}

			} else if secondInput == 2 {
				fmt.Print("Write 'READ <PRIMARY_KEY> <TABLE_NAME>'\nFor example, to read from 'test' table and primary key of a row is 'MIDTERM' write 'READ MIDTERM TEST'\n")
				thirdInput := GetInputString("read")

				fmt.Print(db_operations.Read(strings.Fields(thirdInput)[1], strings.Fields(thirdInput)[2]))
			} else {
				fmt.Print("Invalid input. Please, give valid input!\n")
			}
		} else if input == 3 { /*Update*/
			fmt.Print("1. Update table name (Press 1)\n2. Update row data (Press 2)\nTo Exit Press 0\nChoose Any Operation To Continue...\n")
			/*Takes input from user*/
			secondInput := GetInputNumber()

			if secondInput == 1 {
				fmt.Print("Write 'UPDATE TABLENAME <OLD_TABLE_NAME> <NEW_TABLE_NAME>'\nFor example, to update a table called 'TEST' to 'EXAM' write 'UPDATE TABLENAME TEST EXAM\n")
				thirdInput := GetInputString("update")

				if WordCount(thirdInput) == 4 && (strings.Fields(thirdInput)[1] == "tablename" || strings.Fields(thirdInput)[1] == "TABLENAME") {
					fmt.Print(db_operations.UpdateTableName(strings.Fields(thirdInput)[2], strings.Fields(thirdInput)[3]))
				} else {
					fmt.Print("Invalid input. Please, give valid input!\n")
				}
			} else if secondInput == 2 {
				var thirdInput string

				/*validates input from user*/
				validInput := false
				for validInput == false {
					fmt.Print("Write 'UPDATE <PRIMARY_KEY> <TABLE_NAME>' to update specific row of a table\nFor example, to update a row containing primary key 'MIDTERM' of table 'EXAM' write 'UPDATE MIDTERM EXAM'\n")
					thirdInput = GetInputString("update")

					if WordCount(thirdInput) == 3 {
						/*Checks if directory exists*/
						if !db_operations.CheckDirectory(strings.Fields(thirdInput)[2]) {
							fmt.Print("Invalid table name. Please, give valid table name!\n")
						} else if !db_operations.CheckFile(strings.Fields(thirdInput)[1], strings.Fields(thirdInput)[2]) { /*Checks if file exists*/
							fmt.Print("Invalid primary key name. Please, give valid primary key!\n")
						} else {
							validInput = true
						}

					} else {
						fmt.Print("Invalid input. Please, give valid input!\n")
					}
				}

				/*updates row data*/
				fmt.Print(db_operations.UpdateRow(strings.Fields(thirdInput)[1], strings.Fields(thirdInput)[2]))
			} else {
				fmt.Print("Invalid input. Please, give valid input!\n")
			}
		} else if input == 4 { /*Delete*/
			fmt.Print(input)
		} else {
			fmt.Print("Invalid input. Please, give valid input!\n")
		}
		fmt.Print("\n")
	}

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
func GetInputString(operation string) string {
	validInput := false
	var line string
	for validInput == false {
		/*Read input from the user*/
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		line = scanner.Text()

		/*Command validations*/
		if operation == "create" {
			if strings.HasPrefix(line, "create table") || strings.HasPrefix(line, "CREATE TABLE") {
				if WordCount(line) != 4 {
					fmt.Println("Invalid command. Please enter valid command!")
				} else {
					validInput = true
				}
			} else {
				fmt.Println("Invalid command. Please enter valid command!")
			}
		} else if operation == "read" {
			if strings.HasPrefix(line, "read") || strings.HasPrefix(line, "READ") {
				if WordCount(line) != 3 {
					fmt.Println("Invalid command. Please enter valid command!")
				} else {
					validInput = true
				}
			} else {
				fmt.Println("Invalid command. Please enter valid command!")
			}
		} else if operation == "update" {
			if strings.HasPrefix(line, "update") || strings.HasPrefix(line, "UPDATE") {
				if WordCount(line) >= 3 {
					validInput = true
				} else {
					fmt.Println("Invalid command. Please enter valid command!")
				}
			} else {
				fmt.Println("Invalid command. Please enter valid command!")
			}
		}
	}

	return line
}

// WordCount counts total words
func WordCount(s string) int {
	return len(strings.Fields(s))
}
