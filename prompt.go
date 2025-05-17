package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Connected to Clover database!\n")
	for {
		fmt.Print("> ") // Prompt
		if !scanner.Scan() {
			break // EOF or error
		}

		line := scanner.Text()
		line = strings.TrimSpace(line)

		if line == "exit" || line == "quit" {
			fmt.Println("Bye!")
			return
		}

		internal.Parse(line, db)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading input:", err)
	}
}
