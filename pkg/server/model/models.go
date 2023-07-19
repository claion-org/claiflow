package model

type Scanner interface {
	Scan(dest ...any) error
}

type ModelContext interface {
	TableName() string
	ColumnNames() []string
	ColumnValues() []any
	Scan(scanner Scanner) error
}
