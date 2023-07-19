package model

import (
	"database/sql"
	"time"

	"github.com/claion-org/claiflow/pkg/cryptography"
)

//go:generate go run -mod=mod github.com/abice/go-enum --file=webhook.go --names --nocase=true

/*
ENUM(
none
jq
)
*/
type WebhookConditionValidator int32

/*
ENUM(
id
uuid
name
summary
url
method
headers
timeout
condition_validator
condition_filter
created
updated
)
*/
type WebhookFields int

type Webhook struct {
	ID                 int64                     `column:"id"`   // pk
	UUID               string                    `column:"uuid"` // uuid
	Name               string                    `column:"name"`
	Summary            sql.NullString            `column:"summary"`
	URL                string                    `column:"url"`
	Method             string                    `column:"method"`
	Headers            cryptography.CipherHeader `column:"headers"`
	Timeout            sql.NullInt32             `column:"timeout"`
	ConditionValidator sql.NullInt32             `column:"condition_validator"`
	ConditionFilter    sql.NullString            `column:"condition_filter"`
	Created            time.Time                 `column:"created"`
	Updated            time.Time                 `column:"updated"`
}

var _ ModelContext = (*Webhook)(nil)

func (Webhook) TableName() string { return "webhook" }

func (Webhook) ColumnNames() []string {
	return WebhookFieldsNames()
}

func (r *Webhook) Scan(s Scanner) error {
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
		&r.URL,
		&r.Method,
		&r.Headers,
		&r.Timeout,
		&r.ConditionValidator,
		&r.ConditionFilter,
		&created,
		&updated,
	)
}

func (r Webhook) ColumnValues() []any {
	return []any{
		r.ID,
		r.UUID,
		r.Name,
		r.Summary,
		r.URL,
		r.Method,
		r.Headers,
		r.Timeout,
		r.ConditionValidator,
		r.ConditionFilter,
		r.Created,
		r.Updated,
	}
}
