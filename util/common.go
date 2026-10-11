package util

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	json5 "github.com/yosuke-furukawa/json5/encoding/json5"

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

// Print struct doc with better indentation
func PrettyPrint(doc *document.Document) {
	m := doc.AsMap()
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(string(b))
}

// Extracts json query passed into operation methods e.g. db.myCol.find({ _id: "123" })
func ExtractJsonArg(input string) (string, bool) {
	open := strings.Index(input, "(")
	close := strings.LastIndex(input, ")")
	if open == -1 || close == -1 || close < open || close == open+1 {
		return "", false
	}
	arg := input[open+1 : close]
	return arg, true
}

// FindArgsBounds returns the index of the opening parenthesis of a function
// call along with the index of the matching closing one. Parentheses inside
// quoted strings are ignored, so that find({title: "Dr. No"}) is handled
// correctly. It returns (-1, -1) when no balanced call is present, and
// (-1, i) when a closing parenthesis comes before any opening one.
func FindArgsBounds(input string) (int, int) {
	open, depth := -1, 0
	var quote byte

	for i := 0; i < len(input); i++ {
		c := input[i]

		if quote != 0 {
			switch c {
			case '\\':
				i++ // skip the escaped character
			case quote:
				quote = 0
			}
			continue
		}

		switch c {
		case '"', '\'':
			quote = c
		case '(':
			if open == -1 {
				open = i
			}
			depth++
		case ')':
			if open == -1 {
				return -1, i
			}
			depth--
			if depth == 0 {
				return open, i
			}
		}
	}

	return -1, -1
}

// Parses json and maps key-value pairs into a map
func ParseJson(input string) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json5.Unmarshal([]byte(input), &result)
	return result, err
}

// Extracts argument passed into a function e.g. db.CreateCollection("myNewCollection")
func ExtractStringArg(input string) (string, bool) {
	// Match function call with optional argument
	re := regexp.MustCompile(`\s*\(\s*"([^"]+)"\s*\)`)
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

// Checks if an interface{} value is of type struct or not.
func IsStructLike(data interface{}) bool {
	return reflect.TypeOf(data).Kind() == reflect.Struct
}
