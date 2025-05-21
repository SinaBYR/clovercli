package internal

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/ostafen/clover/v2"
	"github.com/sinabyr/clovercli/util"
)

func Parse(line string, db *clover.DB) {
	if strings.HasPrefix(line, "show ") {
		parseShow(line, db)
		return
	}

	if strings.HasPrefix(line, "db.") {
		parseDB(line, db)
		return
	}
}

// Parse and evaluate "show" command in REPL
func parseShow(line string, db *clover.DB) {
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

// Parse and evaluate "db" command in REPL
func parseDB(line string, db *clover.DB) {
	commands := strings.Split(line, ".")

	collectionName := commands[1]
	if len(commands) < 2 || commands[1] == "" {
		fmt.Println("error: Missing collection name!")
		return
	}

	colls, err := db.ListCollections()
	if err != nil {
		log.Fatal("Failed to list collections")
		return
	}

	collectionExists := util.ContainsString(colls, collectionName)
	if !collectionExists {
		fmt.Printf("'%s': No such collection\n", collectionName)
		return
	}

	if len(commands) < 3 || commands[2] == "" {
		fmt.Println("error: Missing operation!")
		return
	}

	operation := commands[2]
	if strings.HasPrefix(operation, "find") { // BUG finddddddd is also valid :/
		EvaluateFind(operation, collectionName, db)
		return
	}

	fmt.Printf("'%s' operation not recognized!\n", commands[2])
}
