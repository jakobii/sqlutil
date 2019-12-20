package sqlutil

import "strings"


// https://www.postgresql.org/docs/current/functions-logical.html

func And(values ...string) string {
	return strings.Join(values, ` AND `)
}

func Or(values ...string) string {
	return strings.Join(values, ` OR `)
}

func Not(values string) string {
	return `NOT ` + values
}