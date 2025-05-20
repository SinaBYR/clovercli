package internal

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/ostafen/clover/v2"
	q "github.com/ostafen/clover/v2/query"
	"github.com/sinabyr/clovercli/util"
)

func Parse(line string, db *clover.DB) {
	if strings.HasPrefix(line, "show ") {
		evalShow(line, db)
		return
	}

	if strings.HasPrefix(line, "db.") {
		evalDB(line, db)
		return
	}
}

// Parse and evaluate show command in REPL
func evalShow(line string, db *clover.DB) {
	cmds := strings.Split(line, " ")
	if cmds[1] == "collections" {
		collectionsRes, err := db.ListCollections()
		if err != nil {
			log.Fatal("Failed to list collections")
		}

		var collections []table.Row

		for no, name := range collectionsRes {
			collections = append(collections, table.Row{no, name})
		}

		t := table.NewWriter()

		t.SetStyle(table.StyleLight)
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"#", "Name"})
		t.AppendRows(collections)
		t.Render()
		return
	}

	fmt.Printf("'%s' command not recognized!\n", cmds[1])
}

// Parse and evaluate db command in REPL
func evalDB(line string, db *clover.DB) {
	commands := strings.Split(line, ".")

	colls, err := db.ListCollections()
	if err != nil {
		log.Fatal("Failed to list collections")
		return
	}

	collectionName := commands[1]

	if collectionName == "" {
		fmt.Println("error: Missing collection name!")
		return
	}

	collectionExists := util.ContainsString(colls, collectionName)
	if !collectionExists {
		fmt.Printf("'%s': No such collection\n", collectionName)
		return
	}

	if len(commands) < 3 {
		fmt.Println("error: Missing operation!")
		return
	}

	operation := commands[2]

	if strings.HasPrefix(operation, "find") {
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

	fmt.Printf("'%s' command not recognized!\n", commands[0])
}
