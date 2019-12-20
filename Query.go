package sqlutil

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
)


// Query creates a transaction every time it is called. it is intended to take the
// monotony our of running queries that will just be marshaled as an interface anyway.
func Query(db *sql.DB, q string) ([]map[string]interface{}, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()
	results, err := QueryTx(tx,q)
	return results, err
}

// QueryTx performs sql queries on an existing transaction.
func QueryTx(tx *sql.Tx, q string) ([]map[string]interface{}, error) {

	rows, err := tx.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.ColumnTypes()
	if err != nil {
		// there are no columns
		return nil, nil
	}

	vals := make([]interface{}, len(columns))
	ptrs := make([]interface{}, len(columns))
	for i := range vals {
		ptrs[i] = &vals[i]
	}

	results := make([]map[string]interface{}, 0)

	for rows.Next() {

		err = rows.Scan(ptrs...)
		if err != nil {
			return nil, err
		}

		row := make(map[string]interface{})
		for i := 0; i < len(columns); i++ {

			// for debugging
			name := columns[i].Name()
			typ := columns[i].DatabaseTypeName()
			val := vals[i]

			switch val.(type) {

			// all unknown types are stringed and sent over as []uint8
			case []uint8:

				// convert slice to string
				v := val.([]uint8)
				str := fmt.Sprintf("%s", v)

				// convert to richer type based on table type.
				switch typ {
				case "UUID":
					if id, err := uuid.Parse(str); err == nil {
						row[name] = id
					} else {
						row[name] = nil
					}
				default:
					row[name] = str
				}

			default:
				row[name] = vals[i]
			}
		}

		results = append(results, row)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return results, err
}
