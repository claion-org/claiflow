package model

import (
	"database/sql"
	"time"

	"github.com/claion-org/claiflow/pkg/cryptography"
)

//go:generate go run -mod=mod github.com/abice/go-enum --file=service.go --names --nocase=true

/*
ENUM(
none
database
)
*/
type ResultSaveType int

/*
ENUM(
regist 		= 0
sent 		= 1
processing	= 2
succeeded	= 4
failed		= 8
)
*/
type StepStatus int

/*
ENUM(
low
middle
high
)
*/
type Priority int

/*
ENUM(
pdate
cluster_uuid
uuid
name
summary
template_uuid
flow
inputs
step_max
subscribed_channel
priority
created
)
*/
type ClusterServiceFields int

type ClusterService struct {
	PartitionDate     time.Time                 `column:"pdate"`        // PK date
	ClusterUUID       string                    `column:"cluster_uuid"` // PK char(32) cluster.uuid
	UUID              string                    `column:"uuid"`         // PK char(32) service.uuid
	Name              string                    `column:"name"`
	Summary           sql.NullString            `column:"summary"`
	TemplateUUID      string                    `column:"template_uuid"`
	Flow              string                    `column:"flow"`
	Inputs            cryptography.CipherObject `column:"inputs"`
	StepMax           int                       `column:"step_max"`
	SubscribedChannel sql.NullString            `column:"subscribed_channel"`
	Priority          Priority                  `column:"priority"`
	Created           time.Time                 `column:"created"`
}

var _ ModelContext = (*ClusterService)(nil)

func (ClusterService) TableName() string { return "service" }

func (ClusterService) ColumnNames() []string {
	return ClusterServiceFieldsNames()
}

func (r *ClusterService) Scan(s Scanner) error {
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
		&r.Name,
		&r.Summary,
		&r.TemplateUUID,
		&r.Flow,
		&r.Inputs,
		&r.StepMax,
		&r.SubscribedChannel,
		&r.Priority,
		&created,
	)
}

func (r ClusterService) ColumnValues() []any {
	return []any{
		r.PartitionDate,
		r.ClusterUUID,
		r.UUID,
		r.Name,
		r.Summary,
		r.TemplateUUID,
		r.Flow,
		r.Inputs,
		r.StepMax,
		r.SubscribedChannel,
		r.Priority,
		r.Created,
	}
}
