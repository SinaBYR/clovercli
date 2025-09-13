package internal

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/ostafen/clover/v2"
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
	parts := strings.Split(line, ".")
	open := strings.Index(parts[1], "(")
	close := strings.LastIndex(parts[1], ")")
	if open == -1 && close == -1 {
		// Check collections
		currCollectionName := parts[1]

		for _, coll := range dbObj.Collections {
			if coll.Name == currCollectionName {
				if len(parts) == 2 {
					fmt.Println(line)
					return
				}

				open := strings.Index(parts[2], "(")
				close := strings.LastIndex(parts[2], ")")

				if open != -1 && close != -1 {
					currMethodName := strings.TrimSpace(parts[2][:open])

					for _, m := range coll.Methods {
						if m.Name == currMethodName {
							fmt.Println("func found!")
							m.Func(currMethodName, currCollectionName, db)
							return
						}
					}
				}

				fmt.Println(line)
				return
			}
		}

		fmt.Printf("Error: %s not found!\n", parts[1])
		return
	}

	if open != -1 && close != -1 {
		// Check methods
		currMethodName := strings.TrimSpace(parts[1][:open])

		for _, m := range dbObj.Methods {
			if m.Name == currMethodName {
				m.Func(parts[1], db)
				return
			}
		}

		fmt.Printf("Error: Method \"%s\" not found!\n", currMethodName)
		return
	}
}
