package internal

import (
	"fmt"

	"github.com/ostafen/clover/v2"
	q "github.com/ostafen/clover/v2/query"
	"github.com/sinabyr/clovercli/util"
)

// Parses operation function string with filtering argument, then queries the database.
//
// "args" must be: function, collectionName, db
func EvaluateFind(args ...interface{}) {
	if len(args) < 1 {
		fmt.Println("error: Too few inputs!")
		return
	}

	if len(args) > 3 {
		fmt.Println("error: Too many inputs!")
		return
	}

	function, ok := args[0].(string)
	if !ok {
		fmt.Println("error: Bad input for collection name")
		return
	}

	collectionName, ok := args[1].(string)
	if !ok {
		fmt.Println("error: Bad input for collection name")
		return
	}

	db, ok := args[2].(*clover.DB)
	if !ok {
		fmt.Println("error: Bad input for collection name")
		return
	}

	filterRawJson, hasArg := util.ExtractJsonArg(function)
	if !hasArg {
		docs, err := db.FindAll(q.NewQuery(collectionName))
		if err != nil {
			fmt.Println("error: Failed to get documents")
			return
		}

		for _, doc := range docs {
			util.PrettyPrint(doc)
		}
		return
	}

	queryFilter, err := util.ParseJson(filterRawJson)
	if err != nil {
		fmt.Println(err)
		return
	}

	query := q.NewQuery(collectionName)
	fmt.Println("queryFilter: ", queryFilter)
	for key, value := range queryFilter {
		fmt.Println("key: ", key, "value:‌", value)
		if util.IsStructLike(value) {
			innerMap, ok := value.(map[string]interface{})
			fmt.Println("innerMap:‌", innerMap)
			if ok {
				fmt.Println("Ok")
				// check for $gte, ...
				for directive, directiveValue := range innerMap {
					fmt.Println(directive, directiveValue)
				}
			} else {
				fmt.Println("Not Ok")
			}
		} else {
			fmt.Println(value)
			fmt.Println("Not a struct")
		}
		query = query.Where(q.Field(key).Eq(value))
	}

	docs, err := db.FindAll(query)
	if err != nil {
		// TODO better error handling (result not found for passed in criteria)
		fmt.Println("error: Failed to get documents")
		return
	}

	for _, doc := range docs {
		util.PrettyPrint(doc)
	}
}

// Parses operation function string, extracts collectionName from argument, then creates the collection.
//
// "args" must be: function, db
func EvaluateCreateCollection(args ...interface{}) {
	if len(args) < 1 {
		fmt.Println("error: Too few inputs!")
		return
	}

	if len(args) > 2 {
		fmt.Println("error: Too many inputs!")
		return
	}

	function, ok := args[0].(string)
	if !ok {
		fmt.Println("error: Bad input for collection name")
		return
	}

	collectionName, hasArg := util.ExtractStringArg(function)
	if !hasArg {
		fmt.Println("error: Missing collection name!")
		return
	}

	db, ok := args[1].(*clover.DB)
	if !ok {
		fmt.Println("error: Bad input for collection name")
		return
	}

	err := db.CreateCollection(collectionName)
	if err != nil {
		// TODO better error handling (result not found for passed in criteria)
		fmt.Println("error: Failed to create collection")
		return
	}

	fmt.Println("Success")
}

// Parses operation function string, extracts collectionName from argument, then drops the collection.
//
// "args" must be: function, db
func EvaluateDropCollection(args ...interface{}) {
	if len(args) < 2 {
		fmt.Println("error: Too few inputs!")
		return
	}

	if len(args) > 2 {
		fmt.Println("error: Too many inputs!")
		return
	}

	function, ok := args[0].(string)
	if !ok {
		fmt.Println("error: Bad input for collection name")
		return
	}

	collectionName, hasArg := util.ExtractStringArg(function)
	if !hasArg {
		fmt.Println("error: Missing collection name!")
		return
	}

	db, ok := args[1].(*clover.DB)
	if !ok {
		fmt.Println("error: Bad input for collection name")
		return
	}

	err := db.DropCollection(collectionName)
	if err != nil {
		if err == clover.ErrCollectionNotExist {
			fmt.Println("error: Collection not found")
			return
		}

		// TODO better error handling (result not found for passed in criteria)
		fmt.Println("error: Failed to drop collection")
		return
	}

	fmt.Println("Success")
}
