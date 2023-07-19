package model

import (
	"database/sql"
	"time"

	"github.com/claion-org/claiflow/pkg/cryptography"
)

//go:generate go run -mod=mod github.com/abice/go-enum --file=service_result.go --names --nocase=true

/*
ENUM(
pdate
cluster_uuid
uuid
result_type
result
created
)
*/
type ClusterServiceResultFields int

type ClusterServiceResult struct {
	PartitionDate  time.Time                 `column:"pdate"`        // pk date
	ClusterUUID    string                    `column:"cluster_uuid"` // pk char(32) cluster.uuid
	UUID           string                    `column:"uuid"`         // pk char(32) service.uuid
	ResultSaveType ResultSaveType            `column:"result_type"`
	Result         cryptography.CipherString `column:"result"`
	Created        time.Time                 `column:"created"`
}

var _ ModelContext = (*ClusterServiceResult)(nil)

func (ClusterServiceResult) TableName() string { return "service_result" }

func (ClusterServiceResult) ColumnNames() []string {
	return ClusterServiceResultFieldsNames()
}

func (r *ClusterServiceResult) Scan(s Scanner) error {
	var (
		created sql.NullTime
	)
	defer func() {
		r.Created = created.Time
	}()
	return s.Scan(
		&r.PartitionDate,
		&r.ClusterUUID,
		&r.UUID,
		&r.ResultSaveType,
		&r.Result,
		&created,
	)
}

func (r ClusterServiceResult) ColumnValues() []any {
	return []any{
		r.PartitionDate,
		r.ClusterUUID,
		r.UUID,
		r.ResultSaveType,
		r.Result,
		r.Created,
	}
}
