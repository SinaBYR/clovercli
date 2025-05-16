package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/sinabyr/clovercli/internal"
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

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ") // Prompt
		if !scanner.Scan() {
			break // EOF or error
		}

		line := scanner.Text()
		line = strings.TrimSpace(line)

		if line == "exit" || line == "quit" {
			fmt.Println("Bye!")
			break
		}

		if line == "exit" || line == "quit" {
			fmt.Println("Bye!")
			break
		}

		// Evaluate the command (replace this with real logic)
		fmt.Println("You typed:", line)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading input:", err)
	}
}
