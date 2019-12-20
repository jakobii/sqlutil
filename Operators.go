package sqlutil



type BinaryOperator func(x string, y string) string




func Like(x,y string) string {
	return x + ` LIKE ` + y + ` ESCAPE '\'`
}








