package b_tree

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
