package sqlutil

import (
	"bytes"
	"sort"
)


// Where generates where clause expressions. the `WHERE` key word is not added
// to the generated expression. no extra space or padding is around the
// generated expression. you can combine the Where function with the AND and
// OR functions to create larger where clauses.
func Where(constraints map[string]string, operator BinaryOperator, any, enclose bool) string {
	if len(constraints) < 1 {
		return ""
	}

	// get and sort data columns
	columns := make([]string,0,len(constraints))
	for column := range constraints {
		columns = append(columns, column)
	}
	sort.Strings(columns)

	// compile where constraints.
	wheres := make([]string,0,len(constraints))
	for _,column := range columns {
		wheres = append(wheres, operator(column,constraints[column]))
	}

	var stmt bytes.Buffer

	if enclose {
		stmt.WriteString( `(` )
	}

	if any {
		stmt.WriteString(Or(wheres...))
	} else {
		stmt.WriteString(And(wheres...))
	}

	if enclose {
		stmt.WriteString( `)` )
	}

	return stmt.String()
}

// WhereRow is a convenience wrapper around the Where functions. It is intended
// to build where clauses that filter single row with a unique index. e.g. a
// rows primary keys.
func WhereRow(constraints map[string]string) string {
	return Where(constraints, EQ,false, true)
}




