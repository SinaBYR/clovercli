package internal

import "fmt"

func GenerateHelp() string {
	return `CloverCLI - A simple interface to interact with CloverDB

Usage: clovercli [OPTIONS] <path>

Arguments:
  <path>		Path to a cloverdb folder.

Options:
  -h, --help		Show this page
  -v, --version		Print version
`
}

func GenerateBadOption(opt string) string {
	return fmt.Sprintf(`CloverCLI - A simple interface to interact with CloverDB

Usage: clovercli [OPTIONS] <path>

Try 'clovercli --help' for help.

Error: Bad option: %s
`, opt)
}

func GenerateVersion() string {
	return fmt.Sprintf(`CloverCLI - %s`, "v0.1.0")
}
