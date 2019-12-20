package sqlutil

import (
	"sort"
	"strings"
)

// ToInsert converts data into a sql insert statement.
// does not do any escaping or type checking.
func Insert(table string, data []map[string]string) string {

	if len(data) < 1 && table == "" {
		return ""
	}

	// get column count
	var columnCount int
	for _,m := range data {
		if len(m) > columnCount {
			columnCount = len(m)
		}
	}

	// get a distinct list of columns names
	columns := make([]string,0, columnCount )
	var isNewColumn bool
	for _,m := range data {
		for newColumn := range m {
			isNewColumn = false
			for _,column := range columns {
				if newColumn != column {
					isNewColumn = true
					break
				}
			}
			if isNewColumn {
				columns = append(columns, newColumn)
			}
		}
	}

	// make column order predictable
	sort.Strings(columns)

	// create row inserts, using the order of the columns
	rows := make([]string,len(data))
	row := make([]string,0,len(data))
	for i,m := range data {
		for j, columnName := range columns {
			if value, ok :=  m[columnName]; ok {
				row[j] = value
			} else {
				row[j] = `DEFAULT`
			}
		}
		rows[i] = `(` + strings.Join(row,`,`) + `)`
	}

	// build statement
	return `INSERT INTO ` + table + ` (` + strings.Join(columns, ", ") + `) VALUES (` + strings.Join(rows, "), (") + `);`
}
