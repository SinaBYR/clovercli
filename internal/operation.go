package internal

import (
	"fmt"

	"github.com/ostafen/clover/v2"
	q "github.com/ostafen/clover/v2/query"
	"github.com/sinabyr/clovercli/util"
)

// Parses operation string with any filtering criteria provided as its argument, then queries the database.
func EvaluateFind(operation string, collectionName string, db *clover.DB) {
	arg, hasArg := util.ExtractJSONArg(operation)
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

	rawQuery, err := util.ParseFilterQuery(arg)
	if err != nil {
		fmt.Println(err)
		return
	}

	query := q.NewQuery(collectionName)
	for key, value := range rawQuery {
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

	return
}
