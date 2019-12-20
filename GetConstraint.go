package sqlutil

import (
	"database/sql"
	"fmt"
	"sync"
)

func GetConstraints(db *sql.DB, schema, table string) (c []Constraint, err error) {

	ctrs, err := Query(db, fmt.Sprintf(
		sqlGetConstraintNames,
		Text(schema),
		Text(table),
	))
	if err != nil {
		return
	}
	if len(ctrs) < 1 {
		return nil, nil
	}

	// async GetConstraint
	var wg sync.WaitGroup
	constraints := make(chan Constraint, len(ctrs))
	errs := make(chan error, len(ctrs))
	for _, ctr := range ctrs {
		if v, ok := ctr["name"]; ok {
			if n, ok := v.(string); ok {
				wg.Add(1)
				go func() {
					defer wg.Done()
					if constraint, err := GetConstraint(db, schema, table, n); err == nil {
						constraints <- constraint
					} else {
						errs <- err
					}
				}()
			}
		}
	}
	wg.Wait()
	close(constraints)
	close(errs)
	for e := range errs {
		if e != nil {
			return nil, e
		}
	}

	for constraint := range constraints {
		c = append(c, constraint)
	}
	return c, nil
}

func GetConstraint(db *sql.DB, schema, table, name string) (c Constraint, err error) {

	row := db.QueryRow(fmt.Sprintf(
		sqlGetConstraint,
		Text(schema),
		Text(table),
		Text(name),
	))
	err = row.Scan(
		&c.Database,
		&c.Schema,
		&c.Table,
		&c.Name,
		&c.Type,
	)
	if err != nil {
		return
	}

	cols, err := Query(db, fmt.Sprintf(
		sqlGetConstraintColumnNames,
		Text(schema),
		Text(table),
		Text(name),
	))
	if err != nil {
		return
	}
	if len(cols) < 1 {
		return
	}

	for _, col := range cols {
		if v, ok := col["name"]; ok {
			if n, ok := v.(string); ok {
				if column, err := GetColumn(db, schema, table, n); err == nil {
					c.Columns = append(c.Columns, column)
				} else {
					return c, err
				}
			}
		}
	}
	return
}

type Constraint struct {
	Database string
	Schema   string
	Table    string
	Name     string
	Type     string
	Columns  []Column
}

var sqlGetConstraintNames = `
select constraint_name as name
from information_schema.table_constraints
where
    table_schema = %s
    and table_name = %s;
`

var sqlGetConstraint = `
select
    t.table_catalog as Database,
    t.table_schema as Schema,
    t.table_name as Table,
    t.constraint_name as Name,
    t.constraint_type as Type
from information_schema.table_constraints as t
where
    table_schema = %s
    and table_name = %s
	and constraint_name = %s;
`

var sqlGetConstraintColumnNames = `
select column_name as name
from information_schema.constraint_column_usage as c
where
    table_schema = %s
    and table_name = %s
	and constraint_name = %s;
`
