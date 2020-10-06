package store

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

const valuesPlaceholder = "@values"

type (
	Row     []interface{}
	RowFunc func(int) Row
)

func bulkImport(db *gorm.DB, query string, n int, rowfunc RowFunc) error {
	if !strings.Contains(query, valuesPlaceholder) {
		panic(fmt.Errorf("query %q does not contain @values reference", query))
	}

	// No records to process, skipping
	if n < 1 {
		return nil
	}

	var placeholders string
	var vals []interface{}

	for i := 0; i < n; i++ {
		row := rowfunc(i)
		if placeholders == "" {
			placeholders = placeholder(len(row), n)
		}
		vals = append(vals, row...)
	}

	sql := strings.Replace(query, valuesPlaceholder, placeholders, 1)

	return db.Exec(sql, vals...).Error
}

func placeholder(cols int, rows int) string {
	lines := make([]string, rows)

	for i := 0; i < rows; i++ {
		l := make([]string, cols)
		for j := 0; j < cols; j++ {
			l[j] = "?"
		}
		lines[i] = "(" + strings.Join(l, ",") + ")"
	}

	return strings.Join(lines, ",")
}
