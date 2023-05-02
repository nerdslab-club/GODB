package db_operations

func QueryIndex(tableName, columnName string, condition string, value string) string {
	/*Checks if directory exists*/
	//if !CheckDirectory(tableName) {
	//	return ""
	//}
	//
	///*Checks if file exists*/
	//if !CheckFile(pk, tableName) {
	//	return ""
	//}
	//
	//db, err := New(root, nil)
	//if err != nil {
	//	fmt.Println("Error ", err)
	//	return ""
	//}
	///*Delete table row*/
	//if err := db.Delete(tableName, pk); err != nil {
	//	fmt.Println("Error", err)
	//}
	return "Successfully Index Created!"
}

func StringInCondition(str string) bool {
	list := []string{"G", "L", "E", "GE", "LE"}
	for _, item := range list {
		if item == str {
			return true
		}
	}
	return false
}
