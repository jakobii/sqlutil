# sqlutil
A group of SQL utilities. Statement generators and datatype conversion.



This library allows you to programmatically build sql statements. statements can be easily build by combining function calls that incrementally generates a statement.
 
```go
package main
import (
	"fmt"
	"time"
    "github.com/jakobii/sqlutil"
)
func main () {
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
```

The above code produces the following output. the output has been formated for readability.
```
SELECT TOP 42 
    "id" AS "ID", 
    "fn" AS "FirstName", 
    "ln" AS "LastName" 
FROM "test"."public"."people" 
WHERE 
    (
        "fn" LIKE ('Jac%') ESCAPE '\' 
        AND "ln" LIKE ('?choa') ESCAPE '\'
    ) 
    AND ("bd" >= '1990-01-01') 
GROUP BY 
    "id", 
    "fn", 
    "ln" 
ORDER BY 
    "ln", 
    "fn";
```
