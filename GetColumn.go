package sqlutil

import (
	"database/sql"
	"fmt"
	"sync"
)

func GetColumns(db *sql.DB, schema, table string) ([]Column, error) {

	// get constraint names
	results, err := Query(db, fmt.Sprintf(
		sqlGetColumnNames,
		Text(schema),
		Text(table),
	))
	if err != nil {
		return nil,nil
	}
	if len(results) < 1 {
		return nil,nil
	}

	cols := make([]string, len(results))

	// string array of columns
	for i := 0; i < len(cols); i++ {
		if v, ok := results[i]["name"]; ok {
			if n, ok := v.(string); ok {
				cols[i] = n
			}
		}
	}

	// async get columns
	c := make([]Column,len(cols))
	columns := make(chan Column, len(cols))
	errs := make(chan error, len(cols))


	// fetchers
	var fetchers sync.WaitGroup
	for _, col := range cols {
		fetchers.Add(1)
		go func(name string) {
			defer fetchers.Done()
			if column, err := GetColumn(db, schema, table, name); err == nil {
				columns <- column
			} else {
				errs <- err
			}
		}(col)
	}


	fetchers.Wait()
	close(errs)
	close(columns)

	for e := range errs {
		if e != nil {
			return c, e
		}
	}

	for col := range columns {
		fmt.Println(col)
		c[(col.Position - 1)] = col
	}

	return c, nil
}

// GetColumn queries a postgres database for column information, and return false if the column does not exits
func GetColumn(db *sql.DB, schema, table, name string) (c Column, err error) {
	row := db.QueryRow(fmt.Sprintf(
		sqlGetColumnInfo,
		Text(schema),
		Text(table),
		Text(name),
	))
	err = row.Scan(
		&c.Database,
		&c.Schema,
		&c.Table,
		&c.Name,
		&c.Position,
		&c.DataType,
		&c.Default,
		&c.Nullable,
		&c.Length,
		&c.Precision,
		&c.Scale,
		&c.IsKey,
	)
	return
}

type Column struct {
	Database  string
	Schema    string
	Table     string
	Name      string
	Position  uint
	DataType  string
	Default   string
	Nullable  bool
	Length    int64
	Precision int64
	Scale     int64
	IsKey     bool
}

var sqlGetColumnNames = `
SELECT column_name AS name
FROM information_schema.columns
WHERE
    table_schema = %s
    AND table_name = %s;
`

var sqlGetColumnInfo = `
SELECT
 	c.table_catalog as "Database",
    c.table_schema as "Schema",
    c.table_name as "Table",
    c.column_name AS "Name",
    c.ordinal_position AS "Position",
    c.data_type AS "DataType",
    COALESCE(column_default,'') AS "Default",
    CASE c.is_nullable WHEN 'YES' THEN true ELSE false END AS "Nullable",
    COALESCE(c.ordinal_position,0) AS "Position",
    COALESCE(c.numeric_precision,c.datetime_precision,c.interval_precision,0) AS "Precision",
    COALESCE(c.numeric_scale,0) AS "Scale",
    case when k.constraint_name is not null then true else false end as "IsKey"
FROM information_schema.columns AS c
left join (
    select k.*, t.constraint_type
    from information_schema.key_column_usage as k
    inner join information_schema.table_constraints as t
        on  t.table_catalog = k.table_catalog
        and t.table_schema = k.table_schema
        and t.table_name = k.table_name
        and k.constraint_catalog = t.constraint_catalog
        and k.constraint_schema = t.constraint_schema
        and k.constraint_name = t.constraint_name
    where
        t.constraint_type = 'PRIMARY KEY'
) as k
    on  k.table_catalog = c.table_catalog
    and k.table_schema = c.table_schema
    and k.table_name = c.table_name
    and k.column_name = c.column_name
WHERE
    c.table_schema = %s
    AND c.table_name = %s
    AND c.column_name = %s

ORDER BY c.ordinal_position;
`
