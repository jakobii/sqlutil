package sqlutil

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)




// TableName constructs a full qualified sql table name
// database and schema will be omitted if left blank
func TableName(database, schema, table string) string {
	var sql bytes.Buffer
	if database != "" && schema != "" {
		sql.WriteString( database + `.`)
	}
	if schema != "" {
		sql.WriteString(schema + `.`)
	}
	sql.WriteString(table)
	return sql.String()
}

func Cast(val, typ string) string {
	return `CAST(` + As(val,typ) + `)`
}

func As(val, name string) string {
	return val + ` AS ` + name
}



// Escape escapes sql string values.
func EscapeString(x string) string {
	return strings.ReplaceAll(x, "'", "''")
}

// Identifier escapes a resource name. e.g. column, table, schema, and database names
func Identifier(x string) string {
	i := strings.ReplaceAll(x, `"`, `""`)
	return `"` + i + `"`
}

func s(x string) string {
	return strings.ReplaceAll(x, "'", "''")
}

// Expression wraps an expression in parentheses
func Expression(x string) string {
	return e(x)
}

// Expr convert string to sql expression.
func e(x string) string {
	return `(` + x + `)`
}

// Psqler is any function that can type switch an interface and return valid tsql syntax
// although there is no way to test the sql sytax in go. at the very least values should be
// escaped to prevent sql injections.
type Psqler func(interface{}) string

// String returns escaped sql string expression.
// can be used on all string data types.
func Text(i interface{}) string {
	var o string
	switch i.(type) {
	case string:
		o = e(`'` + s(i.(string)) + `'`)
	case uuid.UUID:
		o = e(`'` + s(i.(uuid.UUID).String()) + `'`)
	default:
		o = e(`'` + s(fmt.Sprint(i)) + `'`)
	}
	return o
}

// Number returns escaped sql number expression.
// can be used on all string data types.
func Number(i interface{}) string {
	var o string
	switch i.(type) {
	case bool:
		if i.(bool) {
			o = e(strconv.FormatInt(1, 10))
		} else {
			o = e(strconv.FormatInt(0, 10))
		}
	case int:
		o = e(strconv.FormatInt(int64(i.(int)), 10))
	case int8:
		o = e(strconv.FormatInt(int64(i.(int8)), 10))
	case int16:
		o = e(strconv.FormatInt(int64(i.(int16)), 10))
	case int32:
		o = e(strconv.FormatInt(int64(i.(int32)), 10))
	case int64:
		o = e(strconv.FormatInt(i.(int64), 10))
	case uint:
		o = e(strconv.FormatUint(uint64(i.(uint)), 10))
	case uint8:
		o = e(strconv.FormatUint(uint64(i.(uint8)), 10))
	case uint16:
		o = e(strconv.FormatUint(uint64(i.(uint16)), 10))
	case uint32:
		o = e(strconv.FormatUint(uint64(i.(uint32)), 10))
	case uint64:
		o = e(strconv.FormatUint(i.(uint64), 10))
	case float32:
		o = e(strconv.FormatFloat(float64(i.(float32)), 'f', -1, 32))
	case float64:
		o = e(strconv.FormatFloat(float64(i.(float64)), 'f', -1, 64))
	case time.Time:
		o = e(strconv.FormatInt(i.(time.Time).UnixNano(), 10))
	case nil:
		o = `NULL`
	default:
		o = Text(i)
	}
	return o
}

func QuoteString(s string) string {
	return "'" + s + "'"
}