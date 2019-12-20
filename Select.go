package sqlutil

import (
	"bytes"
	"strconv"
	"strings"
)

// ToSelect generates a sql select statement.
// does not do any escaping or type checking.
func Select(top uint64, distinct bool, columns []string, table string, where string, groupby []string, orderby []string, terminate bool) string {

	if table == "" {
		return ""
	}

	var stmt bytes.Buffer

	stmt.WriteString(`SELECT`)

	// top
	if top > 0 {
		stmt.WriteString(` TOP `)
		stmt.WriteString(strconv.FormatUint(top, 10))
	}

	if distinct {
		stmt.WriteString(` DISTINCT`)
	}

	// columns
	if len(columns) > 0 {
		stmt.WriteString(" " + strings.Join(columns, `, `))
	} else {
		stmt.WriteString(" *")
	}

	// table
	stmt.WriteString(` FROM ` + table)

	// where
	if where != "" {
		stmt.WriteString(` WHERE ` + where)
	}

	// Group by
	if len(groupby) > 0 {
		stmt.WriteString(" GROUP BY " +  strings.Join(groupby, ", ") )
	}

	// order by
	if len(orderby) > 0 {
		stmt.WriteString(" ORDER BY " +  strings.Join(orderby, ", ") )
	}

	// terminate statement
	if terminate {
		stmt.WriteString(";")
	}

	return stmt.String()
}