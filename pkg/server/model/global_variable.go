package model

import (
	"database/sql"
	"time"
)

type GlobalVariable struct {
	ID      int64          `column:"id"`   // pk
	UUID    string         `column:"uuid"` // uuid
	Name    string         `column:"name"`
	Summary sql.NullString `column:"summary"`
	Value   string         `column:"value"`
	Created time.Time      `column:"created"`
	Updated time.Time      `column:"updated"`
}

var _ ModelContext = (*GlobalVariable)(nil)

func (GlobalVariable) TableName() string { return "global_variable" }

func (GlobalVariable) ColumnNames() []string {
	return GlobalVariableFieldsNames()
}

func (r *GlobalVariable) Scan(s Scanner) error {
	var (
		created sql.NullTime
		updated sql.NullTime
	)
	defer func() {
		r.Created = created.Time
		r.Updated = updated.Time
	}()
	return s.Scan(
		&r.ID,
		&r.UUID,
		&r.Name,
		&r.Summary,
		&r.Value,
		&created,
		&updated,
	)
}

func (r GlobalVariable) ColumnValues() []any {
	return []any{
		r.ID,
		r.UUID,
		r.Name,
		r.Summary,
		r.Value,
		r.Created,
		r.Updated,
	}
}

//go:generate go run -mod=mod github.com/abice/go-enum --file=global_variable.go --names --nocase=true

/*
ENUM(
id
uuid
name
summary
value
created
updated
)
*/
type GlobalVariableFields int

// func (globvar GlobalVariable) Valid() error {
// 	if !globvar.Summary.Valid {
// 		return nil
// 	}

// 	query, err := url.ParseQuery(globvar.Summary.String)
// 	if err != nil {
// 		return err
// 	}

// 	if !query.Has("gotype") {
// 		return nil
// 	}

// 	for _, it := range query["gotype"] {
// 		switch it {
// 		case "time.Duration":
// 			_, err := time.ParseDuration(globvar.Value)
// 			return err
// 		case "string":
// 			return nil
// 		}
// 	}

// 	return nil
// }
