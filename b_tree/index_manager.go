package b_tree

import "strconv"

var indexManager map[string]*BTree

func InitIndexManager() {
	indexManager = make(map[string]*BTree)
}

func GetIndex(name string) *BTree {
	if indexManager == nil {
		return nil
	}
	return indexManager[name]
}

func UpdateIndex(name string, tree *BTree) {
	indexManager[name] = tree
}

func CreateBtreeName(tableName, columnName string) string {
	return tableName + "_" + columnName
}

func ConvertStringToInt(val interface{}) int {
	var res int
	var num int

	str := toString(val)

	for i := 0; i < len(str); i++ {
		c := str[i]
		if c >= '0' && c <= '9' {
			num = num*10 + int(c-'0')
		} else {
			res = res*256 + int(c)
		}
	}
	return res*1000000 + num
}

func toString(val interface{}) string {
	switch v := val.(type) {
	case bool:
		return strconv.FormatBool(v)
	case int:
		return strconv.Itoa(v)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case string:
		return v
	default:
		return ""
	}
}







