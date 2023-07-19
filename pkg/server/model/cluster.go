package model

import (
	"database/sql"
	"time"
)

//go:generate go run -mod=mod github.com/abice/go-enum --file=cluster.go --names --nocase=true

/*
ENUM(
id
uuid
name
summary
created
updated
)
*/
type ClusterFields int

type Cluster struct {
	ID      int64          `column:"id"`   // pk
	UUID    string         `column:"uuid"` // uuid
	Name    string         `column:"name"`
	Summary sql.NullString `column:"summary"`
	Created time.Time      `column:"created"`
	Updated time.Time      `column:"updated"`
}

var _ ModelContext = (*Cluster)(nil)

func (Cluster) TableName() string { return "cluster" }

func (Cluster) ColumnNames() []string {
	return ClusterFieldsNames()
}

func (r *Cluster) Scan(s Scanner) error {
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
		&created,
		&updated,
	)
}

func (r Cluster) ColumnValues() []any {
	return []any{
		r.ID,
		r.UUID,
		r.Name,
		r.Summary,
		r.Created,
		r.Updated,
	}
}
