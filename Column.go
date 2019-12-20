package sqlutil

type PrimaryKey struct {
	Name    string
	Columns []Column
}

type ForeignKey struct {
	Name     string
	Columns  []Reference
	OnUpdate string
	OnDelete string
}

type Reference struct {
	Column Column
	Source Column
}
