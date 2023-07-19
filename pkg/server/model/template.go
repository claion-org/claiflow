package model

import (
	"database/sql"
	"time"
)

//go:generate go run -mod=mod github.com/abice/go-enum --file=template.go --names --nocase=true

/*
ENUM(
none
predefined
system
)
*/
type Origin int

/*
ENUM(
uuid
name
summary
flow
inputs
origin
created
updated
)
*/
type TemplateFields int

type Template struct {
	UUID    string         `column:"uuid"`
	Name    string         `column:"name"`
	Summary sql.NullString `column:"summary"`
	Flow    string         `column:"flow"`
	Inputs  string         `column:"inputs"`
	Origin  string         `column:"origin"`
	Created time.Time      `column:"created"`
	Updated time.Time      `column:"updated"`
}

var _ ModelContext = (*Template)(nil)

func (Template) TableName() string { return "template" }

func (Template) ColumnNames() []string {
	return TemplateFieldsNames()
}

func (r *Template) Scan(s Scanner) error {
	var (
		created sql.NullTime
		updated sql.NullTime
	)
	defer func() {
		r.Created = created.Time
		r.Updated = updated.Time
	}()
	return s.Scan(
		&r.UUID,
		&r.Name,
		&r.Summary,
		&r.Flow,
		&r.Inputs,
		&r.Origin,
		&created,
		&updated,
	)
}

func (r Template) ColumnValues() []any {
	return []any{
		r.UUID,
		r.Name,
		r.Summary,
		r.Flow,
		r.Inputs,
		r.Origin,
		r.Created,
		r.Updated,
	}
}
