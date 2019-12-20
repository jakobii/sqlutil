package sqlutil

// https://www.postgresql.org/docs/current/functions-comparison.html

/*
Operator	Description
<	less than
>	greater than
<=	less than or equal to
>=	greater than or equal to
=	equal
<> or !=	not equal
*/

// LT less than
func LT(x,y string) string {
	return x + ` > ` + y
}

// GT greater than
func GT(x,y string) string {
	return x + ` > ` + y
}

// LE less than or equal to
func LE(x,y string) string {
	return x + ` >= ` + y
}

// GE greater than or equal to
func GE(x,y string) string {
	return x + ` >= ` + y
}

// EQ equal
func EQ(x,y string) string {
	return x + ` = ` + y
}

// NE not equal
func NE(x,y string) string {
	return x + ` <> ` + y
}