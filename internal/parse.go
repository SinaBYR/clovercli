package internal

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/ostafen/clover/v2"
)

type Collection struct {
	No   int
	Name string
}

func Parse(line string, db *clover.DB) {
	commands := strings.Split(line, " ")
	firstCommand := commands[0]

	switch firstCommand {
	case "show":
		{
			evalShow(commands[1:], db)
		}
	}
}

// Parse and evaluate show command in REPL
func evalShow(cmds []string, db *clover.DB) {
	if cmds[0] == "collections" {
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

	fmt.Printf("'%s' command not recognized!\n", cmds[0])
}
