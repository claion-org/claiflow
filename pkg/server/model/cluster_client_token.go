package model

import (
	"database/sql"
	"time"
)

//go:generate go run -mod=mod github.com/abice/go-enum --file=cluster_client_token.go --names --nocase=true

/*
ENUM(
id
uuid
name
summary
cluster_uuid
token
issued_at_time
expiration_time
created
updated
)
*/
type ClusterClientTokenFields int

type ClusterClientToken struct {
	ID             int64          `column:"id"`   // pk
	UUID           string         `column:"uuid"` // uuid
	Name           string         `column:"name"`
	Summary        sql.NullString `column:"summary"`
	ClusterUUID    string         `column:"cluster_uuid"`
	Token          string         `column:"token"`
	IssuedAtTime   time.Time      `column:"issued_at_time"`
	ExpirationTime time.Time      `column:"expiration_time"`
	Created        time.Time      `column:"created"`
	Updated        time.Time      `column:"updated"`
}

var _ ModelContext = (*ClusterClientToken)(nil)

func (ClusterClientToken) TableName() string { return "cluster_token" }

func (ClusterClientToken) ColumnNames() []string {
	return ClusterClientTokenFieldsNames()
}

func (r *ClusterClientToken) Scan(s Scanner) error {
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
		&r.ClusterUUID,
		&r.Token,
		&r.IssuedAtTime,
		&r.ExpirationTime,
		&created,
		&updated,
	)
}

func (r ClusterClientToken) ColumnValues() []any {
	return []any{
		r.ID,
		r.UUID,
		r.Name,
		r.Summary,
		r.ClusterUUID,
		r.Token,
		r.IssuedAtTime,
		r.ExpirationTime,
		r.Created,
		r.Updated,
	}
}
