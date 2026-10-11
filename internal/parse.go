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

func Parse(line string, db *clover.DB, dbObj *Object) {
	if strings.HasPrefix(line, "show ") {
		parseShow(line, db)
		return
	}

	if strings.HasPrefix(line, "db.") {
		parseDB(line, db, dbObj)
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
func parseDB(line string, db *clover.DB, dbObj *Object) {
	rest := strings.TrimPrefix(strings.TrimSpace(line), "db.")

	open, closeIdx := util.FindArgsBounds(rest)
	if open == -1 {
		if closeIdx != -1 || strings.Contains(rest, "(") {
			fmt.Println("Error: Unbalanced parentheses in command!")
			return
		}

		parseDBRef(rest, dbObj)
		return
	}

	// The argument is taken verbatim from the parentheses instead of splitting
	// the line on dots, since a filter may contain any of them: dots in strings
	// or field names, parentheses inside a regex, etc.
	head := rest[:open]
	arg := strings.TrimSpace(rest[open+1 : closeIdx])

	// db.movies.find({year: 2010})
	if dot := strings.LastIndex(head, "."); dot != -1 {
		dispatchCollectionMethod(db, dbObj, head[:dot], strings.TrimSpace(head[dot+1:]), arg)
		return
	}

	// db.CreateCollection("movies")
	dispatchTopLevelMethod(db, dbObj, head, rest[:closeIdx+1])
}

// parseDBRef handles commands with no parentheses, like "db.movies" or
// "db.CreateCollection", which are echoed back to the user.
func parseDBRef(rest string, dbObj *Object) {
	head, _, _ := strings.Cut(rest, ".")

	for _, coll := range dbObj.Collections {
		if coll.Name == head {
			fmt.Println("db." + rest)
			return
		}
	}

	for _, m := range dbObj.Methods {
		if m.Name == head {
			fmt.Println("db." + rest)
			return
		}
	}

	fmt.Printf("Error: \"%s\" not found!\n", head)
}

// dispatchCollectionMethod looks up a method on a collection and calls it.
func dispatchCollectionMethod(db *clover.DB, dbObj *Object, collectionName, methodName, arg string) {
	for _, coll := range dbObj.Collections {
		if coll.Name != collectionName {
			continue
		}

		for _, m := range coll.Methods {
			if m.Name == methodName {
				m.Func(db, collectionName, methodName, arg)
				return
			}
		}

		fmt.Printf("Error: Method \"%s\" not found on collection \"%s\"!\n", methodName, collectionName)
		return
	}

	fmt.Printf("Error: \"%s\" not found!\n", collectionName)
}

// dispatchTopLevelMethod looks up a method of the db object and calls it.
func dispatchTopLevelMethod(db *clover.DB, dbObj *Object, methodName, function string) {
	for _, m := range dbObj.Methods {
		if m.Name == methodName {
			m.Func(function, db)
			return
		}
	}

	fmt.Printf("Error: Method \"%s\" not found!\n", methodName)
}
