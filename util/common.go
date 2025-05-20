package util

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/ostafen/clover/v2/document"
)

func ContainsString(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func PrettyPrint(doc *document.Document) {
	m := doc.AsMap()
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(string(b))
}

// Extracts query passed into operation methods e.g. db.myCol.find({ _id: "123" })
func ExtractJSONArg(input string) (string, bool) {
	// Match function call with optional argument
	re := regexp.MustCompile(`(?i)^\s*\w+\s*\(\s*(.*?)\s*\)\s*$`)
	matches := re.FindStringSubmatch(input)
	if len(matches) < 2 {
		return "", false
	}
	arg := matches[1]
	if strings.TrimSpace(arg) == "" {
		return "", false // No argument
	}
	return arg, true // Argument found
}

// Transforms query passed to operation methods into valid golang struct e.g. db.myCol.find({ _id: "123" })
func ParseFilterQuery(input string) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(input), &result)
	return result, err

	// TODO Convert JS-style object (unquoted keys) to JSON
	// fixed := regexp.MustCompile(`([{,]\s*)([a-zA-Z_][a-zA-Z0-9_]*)(\s*:)`).ReplaceAllString(arg, `$1"$2"$3`)
}
