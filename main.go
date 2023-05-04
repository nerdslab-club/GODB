package main

import (
	"bufio"
	"fmt"
	"goDB/db_operations"
	"os"
	"strconv"
	"strings"
)

func main() {

	var input int = 1

	for input != 0 {
		/*Choose operation*/
		fmt.Print(
			"1. Create (Press 1)\n" +
				"2. Read (Press 2)\n" +
				"3. Update (Press 3)\n" +
				"4. Delete (Press 4)\n" +
				"5. Query (Press 5)\n" +
				"To Exit Press 0\n" +
				"Choose Any Operation To Continue...\n")

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
			fmt.Print("1. Delete table row (Press 1)\n2. Delete one table (Press 2)\n3. Delete all table (Press 3)\nTo Exit Press 0\nChoose Any Operation To Continue...\n")

			/*Takes input from user*/
			secondInput := GetInputNumber()
			if secondInput == 1 {
				var thirdInput string
				/*validates input from user*/
				validInput := false
				for validInput == false {
					fmt.Print("Write 'DELETE <TABLE_NAME> <PRIMARY_KEY>'\nFor example, to delete row containing primary key 'MIDTERM' of table 'EXAM', write 'DELETE EXAM MIDTERM'\n")
					thirdInput = GetInputString("delete")

					if WordCount(thirdInput) == 3 {
						validInput = true
					} else {
						fmt.Print("Invalid input. Please, give valid input!\n")
					}
				}

				/*Delete one row*/
				fmt.Print(db_operations.DeleteRow(strings.Fields(thirdInput)[1], strings.Fields(thirdInput)[2]))

			} else if secondInput == 2 {
				var thirdInput string
				/*validates input from user*/
				validInput := false
				for validInput == false {
					fmt.Print("Write 'DELETE <TABLE_NAME>'\nFor example, to delete table 'EXAM', write 'DELETE EXAM'\n")
					thirdInput = GetInputString("delete")

					if WordCount(thirdInput) == 2 {
						validInput = true
					} else {
						fmt.Print("Invalid input. Please, give valid input!\n")
					}
				}

				/*Delete one row*/
				fmt.Print(db_operations.DeleteTable(strings.Fields(thirdInput)[1]))
			} else if secondInput == 3 {
				var thirdInput string
				/*validates input from user*/
				validInput := false
				for validInput == false {
					fmt.Print("Write 'DELETE *' to delete all table\nTo Exit Press 0\n")
					thirdInput = GetInputString("delete")

					if WordCount(thirdInput) == 2 && strings.Fields(thirdInput)[1] == "*" {
						if IsSure() {
							validInput = true
						} else {
							fmt.Print("Invalid input. Please, give valid input!\n")
						}
					} else {
						fmt.Print("Invalid input. Please, give valid input!\n")
					}
				}

				/*Deletes all table*/
				fmt.Print(db_operations.DeleteAll())
			} else {
				fmt.Print("Invalid input. Please, give valid input!\n")
			}
		} else if input == 5 {
			fmt.Print(
				"1. Create new index (Press 1)\n" +
					"2. Execute query (Press 2)\n")

			secondInput := GetInputNumber()
			if secondInput == 1 {
				var thirdInput string
				/*validates input from user*/
				validInput := false
				for validInput == false {
					fmt.Print("Write 'INDEX <TABLE_NAME> <COLUMN_NAME>'\nFor example, to create index for column 'MIDTERM' of table 'EXAM', write 'INDEX EXAM MIDTERM'\n")
					thirdInput = GetInputString("index")

					if WordCount(thirdInput) == 3 {
						validInput = true
					} else {
						fmt.Print("Invalid input. Please, give valid input!\n")
					}
				}

				/*Index one column*/
				fmt.Print(db_operations.CreateIndex(strings.Fields(thirdInput)[1], strings.Fields(thirdInput)[2]))

			} else if secondInput == 2 {
				var thirdInput string
				/*validates input from user*/
				validInput := false
				for validInput == false {
					fmt.Print("Write 'QUERY <TABLE_NAME> <COLUMN_NAME> <CONDITION> <VALUE>'\n" +
						"For example, to query a value 'MARK' from 'STUDENT' table, write 'QUERY STUDENT MARK G 95'\n" +
						"Where, equal='==', greater='>', lesser='<', greater_and_equal='>=', lesser_and_equal='<='")
					thirdInput = GetInputString("query")

					if !db_operations.StringInCondition(strings.Fields(thirdInput)[3]) {
						fmt.Print("Invalid condition. Need to be one of (\">\", \"<\", \"==\", \">=\", \"<=\")!\n")
					} else if WordCount(thirdInput) == 4 {
						validInput = true
					} else {
						fmt.Print("Invalid input. Please, give valid input!\n")
					}
				}

				/*Delete one row*/
				fmt.Print(db_operations.QueryIndex(strings.Fields(thirdInput)[1], strings.Fields(thirdInput)[2], strings.Fields(thirdInput)[3], strings.Fields(thirdInput)[4]))
			} else {
				fmt.Print("Invalid input. Please, give valid input!\n")
			}
		} else {
			fmt.Print("Invalid input. Please, give valid input!\n")
		}
		fmt.Print("\n")
	}
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
			number = 99 //default set for invalid number nput
		} else {
			validInput = true
		}

		/*exists if 0 pressed*/
		if number == 0 {
			os.Exit(0)
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

		/*exists if 0 pressed*/
		if line == "0" {
			os.Exit(0)
		}

		/*Command validations*/
		if operation == "create" {
			if strings.HasPrefix(line, "create table ") || strings.HasPrefix(line, "CREATE TABLE ") {
				if WordCount(line) != 4 {
					fmt.Println("Invalid command. Please enter valid command!")
				} else {
					validInput = true
				}
			} else {
				fmt.Println("Invalid command. Please enter valid command!")
			}
		} else if operation == "read" {
			if strings.HasPrefix(line, "read ") || strings.HasPrefix(line, "READ ") {
				if WordCount(line) != 3 {
					fmt.Println("Invalid command. Please enter valid command!")
				} else {
					validInput = true
				}
			} else {
				fmt.Println("Invalid command. Please enter valid command!")
			}
		} else if operation == "update" {
			if strings.HasPrefix(line, "update ") || strings.HasPrefix(line, "UPDATE ") {
				if WordCount(line) >= 3 {
					validInput = true
				} else {
					fmt.Println("Invalid command. Please enter valid command!")
				}
			} else {
				fmt.Println("Invalid command. Please enter valid command!")
			}
		} else if operation == "delete" {
			if strings.HasPrefix(line, "delete ") || strings.HasPrefix(line, "DELETE ") {
				if WordCount(line) >= 2 {
					validInput = true
				} else {
					fmt.Println("Invalid command. Please enter valid command!")
				}
			} else {
				fmt.Println("Invalid command. Please enter valid command!")
			}
		} else if operation == "index" {
			if strings.HasPrefix(line, "index ") || strings.HasPrefix(line, "INDEX ") {
				if WordCount(line) >= 2 {
					validInput = true
				} else {
					fmt.Println("Invalid command. Please enter valid command!")
				}
			} else {
				fmt.Println("Invalid command. Please enter valid command!")
			}
		} else if operation == "query" {
			if strings.HasPrefix(line, "query ") || strings.HasPrefix(line, "QUERY ") {
				if WordCount(line) >= 4 {
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

func IsSure() bool {
	fmt.Print("Are you sure? (Press Y/N)\nTo Exit Press 0\n")
	/*Read input from the user*/
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()

	if input == "Y" || input == "y" {
		return true
	} else if input == "0" { /*exists if 0 pressed*/
		os.Exit(0)
	}
	return false
}
