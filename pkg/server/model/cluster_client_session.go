package model

import (
	"database/sql"
	"time"
)

//go:generate go run -mod=mod github.com/abice/go-enum --file=cluster_client_session.go --names --nocase=true

/*
ENUM(
id
uuid
cluster_client_token_uuid
cluster_uuid
token
issued_at_time
expiration_time
created
updated
)
*/
type ClusterClientSessionFields int

type ClusterClientSession struct {
	ID                     int64     `column:"id"`   // pk
	UUID                   string    `column:"uuid"` // uuid
	ClusterUUID            string    `column:"cluster_uuid"`
	ClusterClientTokenUUID string    `column:"cluster_client_token_uuid"`
	Token                  string    `column:"token"`
	IssuedAtTime           time.Time `column:"issued_at_time"`
	ExpirationTime         time.Time `column:"expiration_time"`
	Created                time.Time `column:"created"`
	Updated                time.Time `column:"updated"`
}

var _ ModelContext = (*ClusterClientSession)(nil)

func (ClusterClientSession) TableName() string { return "session" }

func (ClusterClientSession) ColumnNames() []string {
	return ClusterClientSessionFieldsNames()
}

func (r *ClusterClientSession) Scan(s Scanner) error {
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
		&r.ClusterClientTokenUUID,
		&r.ClusterUUID,
		&r.Token,
		&r.IssuedAtTime,
		&r.ExpirationTime,
		&created,
		&updated,
	)
}

func (r ClusterClientSession) ColumnValues() []any {
	return []any{
		r.ID,
		r.UUID,
		r.ClusterClientTokenUUID,
		r.ClusterUUID,
		r.Token,
		r.IssuedAtTime,
		r.ExpirationTime,
		r.Created,
		r.Updated,
	}
}
