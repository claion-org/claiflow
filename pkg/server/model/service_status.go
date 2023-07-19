package model

import (
	"database/sql"
	"time"
)

type ServiceWithStatuses struct {
	ClusterService
	ClusterServiceStatuses []ClusterServiceStatus
}

//go:generate go run -mod=mod github.com/abice/go-enum --file=service_status.go --names --nocase=true

/*
ENUM(
pdate
cluster_uuid
uuid
created
step_max
step_seq
status
started
ended
message
)
*/
type ClusterServiceStatusFields int

type ClusterServiceStatus struct {
	PartitionDate time.Time      `column:"pdate"`        // PK date
	ClusterUUID   string         `column:"cluster_uuid"` // PK char(32) cluster.uuid
	Uuid          string         `column:"uuid"`         // PK char(32) service.uuid
	Created       time.Time      `column:"created"`      // PK datetime(6)
	StepMax       int            `column:"step_max"`
	StepSeq       int            `column:"step_seq"`
	Status        StepStatus     `column:"status"`
	Started       sql.NullTime   `column:"started"`
	Ended         sql.NullTime   `column:"ended"`
	Message       sql.NullString `column:"message"`
}

var _ ModelContext = (*ClusterServiceStatus)(nil)

func (ClusterServiceStatus) TableName() string { return "service_status" }

func (ClusterServiceStatus) ColumnNames() []string {
	return ClusterServiceStatusFieldsNames()
}

func (r *ClusterServiceStatus) Scan(s Scanner) error {
	var (
		created sql.NullTime
	)
	defer func() {
		r.Created = created.Time
	}()
	return s.Scan(
		&r.PartitionDate,
		&r.ClusterUUID,
		&r.Uuid,
		&created,
		&r.StepMax,
		&r.StepSeq,
		&r.Status,
		&r.Started,
		&r.Ended,
		&r.Message,
	)
}

func (r ClusterServiceStatus) ColumnValues() []any {
	return []any{
		r.PartitionDate,
		r.ClusterUUID,
		r.Uuid,
		r.Created,
		r.StepMax,
		r.StepSeq,
		r.Status,
		r.Started,
		r.Ended,
		r.Message,
	}
}
