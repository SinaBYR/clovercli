package internal

import (
	"testing"

	d "github.com/ostafen/clover/v2/document"
	"github.com/sinabyr/clovercli/util"
)

// movieDoc returns a document used by every filter test.
func movieDoc() *d.Document {
	return d.NewDocumentOf(map[string]interface{}{
		"title":   "Inception",
		"year":    2010,
		"rating":  8.8,
		"runtime": 148,
		"author":  "Christopher Nolan",
		"cast":    []string{"Leonardo DiCaprio", "Tom Hardy"},
		"meta": map[string]interface{}{
			"score": 9,
		},
	})
}

func TestParseFilter(t *testing.T) {
	doc := movieDoc()

	tests := []struct {
		name   string
		filter string
		want   bool
	}{
		// implicit equality
		{"implicit equality", `{author: "Christopher Nolan"}`, true},
		{"implicit equality mismatch", `{author: "Wes Anderson"}`, false},
		{"implicit equality on missing field", `{director: "Someone"}`, false},
		{"implicit equality across number types", `{year: 2010.0}`, true},
		{"implicit equality on array field", `{cast: "Tom Hardy"}`, false},

		// $eq / $ne
		{"$eq", `{year: {$eq: 2010}}`, true},
		{"$eq mismatch", `{year: {$eq: 1994}}`, false},
		{"$ne", `{year: {$ne: 1994}}`, true},
		{"$ne mismatch", `{year: {$ne: 2010}}`, false},
		{"$ne on missing field", `{director: {$ne: "Someone"}}`, true},

		// $gt / $gte / $lt / $lte
		{"$gt", `{year: {$gt: 2009}}`, true},
		{"$gt mismatch", `{year: {$gt: 2010}}`, false},
		{"$gte", `{year: {$gte: 2010}}`, true},
		{"$gte mismatch", `{year: {$gte: 2011}}`, false},
		{"$lt", `{runtime: {$lt: 149}}`, true},
		{"$lt mismatch", `{runtime: {$lt: 148}}`, false},
		{"$lte", `{runtime: {$lte: 148}}`, true},
		{"$lte mismatch", `{runtime: {$lte: 147}}`, false},
		{"$gt on strings", `{title: {$gt: "A"}}`, true},

		// $in / $nin
		{"$in", `{year: {$in: [1994, 2010]}}`, true},
		{"$in mismatch", `{year: {$in: [1994, 2000]}}`, false},
		{"$in across number types", `{year: {$in: [2010.0]}}`, true},
		{"$in with a different type", `{year: {$in: ["2010"]}}`, false},
		{"$nin", `{year: {$nin: [1994, 2000]}}`, true},
		{"$nin mismatch", `{year: {$nin: [2010]}}`, false},

		// $exists
		{"$exists true", `{title: {$exists: true}}`, true},
		{"$exists false", `{title: {$exists: false}}`, false},
		{"$exists false on missing field", `{director: {$exists: false}}`, true},
		{"$exists true on missing field", `{director: {$exists: true}}`, false},

		// $regex / $like
		{"$regex", `{title: {$regex: "^Incept"}}`, true},
		{"$regex mismatch", `{title: {$regex: "^Interstellar$"}}`, false},
		{"$regex is case sensitive", `{author: {$regex: "nolan"}}`, false},
		{"$regex on a number", `{year: {$regex: "^2010$"}}`, false},
		{"$like", `{title: {$like: "nception"}}`, true},
		{"$like mismatch", `{title: {$like: "nolan"}}`, false},

		// $not
		{"$not", `{year: {$not: {$gt: 2011}}}`, true},
		{"$not mismatch", `{year: {$not: {$gte: 2010}}}`, false},
		{"$not with operators", `{year: {$not: {$gt: 2000, $lt: 2020}}}`, false},

		// multiple keys and multiple operators are combined with AND
		{"multiple keys", `{author: "Christopher Nolan", year: 2010}`, true},
		{"multiple keys mismatch", `{author: "Wes Anderson", year: 2010}`, false},
		{"multiple operators", `{year: {$gt: 2009, $lt: 2011}}`, true},
		{"multiple operators mismatch", `{year: {$gt: 2010, $lt: 2011}}`, false},

		// nested fields
		{"nested field equality", `{"meta.score": 9}`, true},
		{"nested field operator", `{"meta.score": {$gt: 8}}`, true},
		{"nested field operator mismatch", `{"meta.score": {$lt: 8}}`, false},
		{"nested missing field", `{"meta.director": "Someone"}`, false},

		// $and / $or / $nor / $not
		{"$and", `{$and: [{year: 2010}, {rating: 8.8}]}`, true},
		{"$and mismatch", `{$and: [{year: 2010}, {rating: 9}]}`, false},
		{"$and single filter", `{$and: [{year: 2010}]}`, true},
		{"$or", `{$or: [{year: 1994}, {title: "Inception"}]}`, true},
		{"$or mismatch", `{$or: [{year: 1994}, {title: "Other"}]}`, false},
		{"$nor", `{$nor: [{year: 1994}, {title: "Other"}]}`, true},
		{"$nor mismatch", `{$nor: [{year: 1994}, {title: "Inception"}]}`, false},
		{"$not filter", `{$not: {year: 1994}}`, true},
		{"$not filter mismatch", `{$not: {year: 2010}}`, false},
		{"top level $not and key", `{author: "Christopher Nolan", $not: {year: 2010}}`, false},

		// logical operators can be nested
		{"nested $and in $or", `{$or: [{$and: [{year: 2010}, {rating: 8.8}]}, {author: "Wes Anderson"}]}`, true},
		{"nested $or in $and", `{$and: [{$or: [{year: 1994}, {rating: 8.8}]}]}`, true},
		{"directive nested in $or", `{$or: [{year: {$gt: 2005}}]}`, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter, err := util.ParseJson(tt.filter)
			if err != nil {
				t.Fatalf("ParseJson(%s) returned an error: %v", tt.filter, err)
			}

			criteria, err := parseFilter(filter)
			if err != nil {
				t.Fatalf("parseFilter(%s) returned an error: %v", tt.filter, err)
			}

			if criteria == nil {
				t.Fatalf("parseFilter(%s) returned no criteria", tt.filter)
			}

			if got := criteria.Satisfy(doc); got != tt.want {
				t.Errorf("parseFilter(%s).Satisfy() = %v, want %v", tt.filter, got, tt.want)
			}
		})
	}
}

func TestParseFilterEmptyMatchesEveryDocument(t *testing.T) {
	criteria, err := parseFilter(map[string]interface{}{})
	if err != nil {
		t.Fatalf("parseFilter() returned an error: %v", err)
	}

	if criteria != nil {
		t.Errorf("parseFilter({}) = %v, want nil criteria to match every document", criteria)
	}
}

func TestParseFilterErrors(t *testing.T) {
	tests := []struct {
		name   string
		filter string
	}{
		{"unknown field directive", `{year: {$foo: 1}}`},
		{"unknown top level directive", `{$foo: 1}`},
		{"$in must be an array", `{year: {$in: 2010}}`},
		{"$nin must be an array", `{year: {$nin: 2010}}`},
		{"$exists must be a boolean", `{year: {$exists: "yes"}}`},
		{"$regex must be a string", `{title: {$regex: 5}}`},
		{"$regex with an invalid pattern", `{title: {$regex: "([a-z"}}`},
		{"$like must be a string", `{title: {$like: 5}}`},
		{"$not must be an object", `{year: {$not: 5}}`},
		{"$not must not be empty", `{year: {$not: {}}}`},
		{"$and must be an array", `{$and: {year: 2010}}`},
		{"$and must not be empty", `{$and: []}`},
		{"$and filter must be an object", `{$and: [2010]}`},
		{"$or must be an array", `{$or: {year: 2010}}`},
		{"$or must not be empty", `{$or: []}`},
		{"$nor must be an array", `{$nor: {year: 2010}}`},
		{"$not must be an object of filters", `{$not: 5}`},
		{"top level $not must not be empty", `{$not: {}}`},
		{"empty operator list", `{year: {}}`},
		{"unknown directive inside $and", `{$and: [{year: {$foo: 1}}]}`},
		{"unknown directive inside $or", `{$or: [{year: {$foo: 1}}]}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter, err := util.ParseJson(tt.filter)
			if err != nil {
				t.Fatalf("ParseJson(%s) returned an error: %v", tt.filter, err)
			}

			if _, err := parseFilter(filter); err == nil {
				t.Errorf("parseFilter(%s) = nil error, want an error", tt.filter)
			}
		})
	}
}
