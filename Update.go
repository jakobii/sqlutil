package sqlutil

import (
	"bytes"
	"sort"
	"strings"
)

// ToUpdate converts data into a sql update statement.
// values and keys should be fully escaped sql expressions.
func Update(table string, data map[string]string, where string) string {

	if len(data) < 1 && table == "" {
		return ""
	}

	// get and sort data columns
	columns := make([]string,0,len(data))
	for column := range data {
		columns = append(columns, column)
	}
	sort.Strings(columns)

	// SETs
	sets := make([]string,0,len(columns))
	for _,column := range columns {
		sets = append(sets, column + ` = ` + data[column])
	}

	// build statement
	var stmt bytes.Buffer
	stmt.WriteString(`UPDATE ` + table + ` SET ` + strings.Join(sets, `, `) )

	if where != "" {
		stmt.WriteString(`WHERE ` + where)
	}

	stmt.WriteString(`;`)

	return stmt.String()
}
