package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/chzyer/readline"
	cl "github.com/ostafen/clover/v2"
	"github.com/sinabyr/clovercli/internal"
	"github.com/sinabyr/clovercli/util"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println(internal.GenerateHelp())
		return
	}

	first := args[0]

	if strings.HasPrefix(first, "--") || strings.HasPrefix(first, "-") {
		switch first {
		case "--help":
			fmt.Println(internal.GenerateHelp())
			return
		case "-h":
			fmt.Println(internal.GenerateHelp())
			return

		case "--version":
			fmt.Println(internal.GenerateVersion())
			return
		case "-v":
			fmt.Println(internal.GenerateVersion())
			return

		default:
			fmt.Println(internal.GenerateBadOption(first))
			return
		}
	}

	info, err := os.Stat(first)
	if err != nil {
		fmt.Println("error: ", err)
		return
	}
	if !info.IsDir() {
		fmt.Println("Not a directory")
		return
	}

	fileExists, err := util.PathExists(filepath.Join(first, "data.db"))
	if err != nil {
		fmt.Println("error: ", err)
		return
	}
	if !fileExists {
		fmt.Println("error: ", "data.db doesn't exist")
		return
	}

	db, err := cl.Open(first)
	defer db.Close()

	rl, err := readline.NewEx(&readline.Config{
		Prompt:          "> ",
		HistoryFile:     "/tmp/clovercli.tmp",
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer rl.Close()

	fmt.Println("Use ↑/↓ to navigate history. Type 'exit' to exit.")

	dbObj := internal.NewObject("db")
	dbObj.RegisterField(internal.Field{Name: "path", Value: first})
	dbObj.RegisterMethod(internal.Method{Name: "CreateCollection", Func: internal.EvaluateCreateCollection})
	dbObj.RegisterMethod(internal.Method{Name: "DropCollection", Func: internal.EvaluateDropCollection})

	colls, err := db.ListCollections()
	if err != nil {
		log.Fatal("Error: Something went wrong!")
		return
	}
	for _, c := range colls { // Loop over collections
		dbObj.RegisterCollection(c)
	}

	for {
		line, err := rl.Readline()
		if err != nil { // e.g. io.EOF or Ctrl+D
			break
		}
		if line == "" {
			continue
		}

		if line == "exit" || line == "quit" {
			fmt.Println("Bye!")
			return
		}

		internal.Parse(line, db, dbObj)
	}
}
