package util

import "testing"

func TestFindArgsBounds(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantArg string
	}{
		{"call", `db.movies.find({year: 2010})`, `{year: 2010}`},
		{"dot inside a string", `db.movies.find({title: "Dr. No"})`, `{title: "Dr. No"}`},
		{"parenthesis inside a string", `db.movies.find({title: "a)b"})`, `{title: "a)b"}`},
		{"single quotes", `db.movies.find({title: 'Dr. No'})`, `{title: 'Dr. No'}`},
		{"escaped quote", `db.movies.find({title: "a\")b"})`, `{title: "a\")b"}`},
		{"nested parentheses", `db.movies.find({a: (1)})`, `{a: (1)}`},
		{"only the first call", `db.movies.find({year: 2010}).find({title: "x"})`, `{year: 2010}`},
		{"no call", `db.movies`, ``},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			open, close := FindArgsBounds(tt.input)

			arg := ""
			if open != -1 && close > open {
				arg = tt.input[open+1 : close]
			}

			if arg != tt.wantArg {
				t.Errorf("FindArgsBounds(%q) = (%d, %d), argument %q, want %q", tt.input, open, close, arg, tt.wantArg)
			}
		})
	}
}

func TestFindArgsBoundsUnbalanced(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantOpen  int
		wantClose int
	}{
		{"unterminated call", `db.movies.find({`, -1, -1},
		{"closing parenthesis before the call", `db.movies.find)`, -1, 14},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			open, close := FindArgsBounds(tt.input)
			if open != tt.wantOpen || close != tt.wantClose {
				t.Errorf("FindArgsBounds(%q) = (%d, %d), want (%d, %d)", tt.input, open, close, tt.wantOpen, tt.wantClose)
			}
		})
	}
}
