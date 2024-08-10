package server

import (
	"errors"
	"os"
	"slices"
	"strings"
	"testing"
)

type ParsedInput struct {
	original string
	connId   string
	dbName   string
	tables   string
	error    error
}

func TestParseInput(t *testing.T) {
	tests := []ParsedInput{
		{
			original: "main:sakila",
			connId:   "main",
			dbName:   "sakila",
			tables:   "",
			error:    nil,
		},
		{
			original: "main:sakila:",
			connId:   "main",
			dbName:   "sakila",
			tables:   "",
			error:    nil,
		},
		{
			original: ":sakila:",
			connId:   "",
			dbName:   "sakila",
			tables:   "",
			error:    nil,
		},
		{
			original: ":",
			connId:   "",
			dbName:   "",
			tables:   "",
			error:    nil,
		},
		{
			original: "main:sakila:table_1,,,,              ,,,table_2, table_3",
			connId:   "main",
			dbName:   "sakila",
			tables:   "table_1,table_2,table_3",
			error:    nil,
		},
		{
			original: "main:sakila:table,,,,,",
			connId:   "main",
			dbName:   "sakila",
			tables:   "table",
			error:    nil,
		},
	}

	for _, p := range tests {
		connId, dbName, tables, err := parseInput(p.original)
		tablesString := strings.Join(tables, ",")
		if connId != p.connId || dbName != p.dbName || tablesString != p.tables || err != nil {
			t.Errorf("test failed for input %q %q %q", p.original, tablesString, p.tables)
		}
	}

	failing := []ParsedInput{
		{
			original: "main:sakila::",
			error:    errors.New("unexpected input data format: main:sakila::"),
		},
	}

	for _, p := range failing {
		_, _, _, err := parseInput(p.original)
		if err.Error() != p.error.Error() {
			t.Errorf("wanted error %q, got %q", p.error, err)
		}
	}
}

func TestGetKeys(t *testing.T) {
	f, _ := os.CreateTemp("", "temp_authorized_keys")
	defer os.Remove(f.Name())

	keys := getKeys(f.Name())

	if !slices.Equal([]string{}, keys) {
		t.Errorf("expected %q, got %q", []string{}, keys)
	}

	f.WriteString("\n\nssh-abc 0J+f6JVJX4BE2SfEkv comment\nssh-def 123 comment\n       ssh-ghi aaa comment\n\nfoo bar baz")

	keys = getKeys(f.Name())

	expected := []string{"ssh-abc 0J+f6JVJX4BE2SfEkv comment", "ssh-def 123 comment", "ssh-ghi aaa comment", "foo bar baz"}
	if !slices.Equal(expected, keys) {
		t.Errorf("expected %q, got %q", expected, keys)
	}

}
