package sqlutil

import (
	"fmt"
	"testing"
	"time"
)

func TestSelect(t *testing.T) {
	stmt := Select(
		// top
		42,

		// distinct
		false,

		// columns
		[]string{
			As(Identifier(`id`), Identifier(`ID`)),
			As(Identifier(`fn`), Identifier(`FirstName`)),
			As(Identifier(`ln`), Identifier(`LastName`)),
		},

		// table
		TableName(
			Identifier(`test`),   //db
			Identifier(`public`), //sch
			Identifier(`people`), //tb
		),

		// where
		And(
			Where(
				map[string]string{
					Identifier("fn"): Text("Jac%"),
					Identifier("ln"): Text("?choa"),
				},
				Like,
				false,
				true,
			),
			Where(
				map[string]string{
					Identifier("bd"): QuoteString(Date(time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC))),
				},
				GE,
				false,
				true,
			),
		),

		// group by
		[]string{
			Identifier(`id`),
			Identifier(`fn`),
			Identifier(`ln`),
		},

		// order by
		[]string{
			Identifier(`ln`),
			Identifier(`fn`),
		},

		// terminate `;`
		true,
	)
	fmt.Println(stmt)
}
