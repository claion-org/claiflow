package model

import (
	"database/sql"
	"time"
)

//go:generate go run -mod=mod github.com/abice/go-enum --file=cluster_information.go --names --nocase=true

/*
ENUM(
id
cluster_uuid
polling_offset
created
updated
)
*/
type ClusterInformationFields int

type ClusterInformation struct {
	ID          uint   `column:"id"`
	ClusterUUID string `column:"cluster_uuid"`
	// PollingCount int    `column:"polling_count"`
	PollingOffset time.Time `column:"polling_offset"`
	Created       time.Time `column:"created"`
	Updated       time.Time `column:"updated"`
}

var _ ModelContext = (*ClusterInformation)(nil)

func (ClusterInformation) TableName() string { return "cluster_information" }

func (ClusterInformation) ColumnNames() []string {
	return ClusterInformationFieldsNames()
}

func (r *ClusterInformation) Scan(s Scanner) error {
	var (
		pollingOffset sql.NullTime
		created       sql.NullTime
		updated       sql.NullTime
	)
	defer func() {
		r.PollingOffset = pollingOffset.Time
		r.Created = created.Time
		r.Updated = updated.Time
	}()
	return s.Scan(
		&r.ID,
		&r.ClusterUUID,
		&pollingOffset,
		&created,
		&updated,
	)
}

func (r ClusterInformation) ColumnValues() []any {
	return []any{
		r.ID,
		r.ClusterUUID,
		r.PollingOffset,
		r.Created,
		r.Updated,
	}
}
