package macro

import (
	"database/sql"
	"time"
)

func NewNullString(p *string) sql.NullString {
	if p == nil {
		return sql.NullString{}
	}

	return sql.NullString{String: *p, Valid: true}
}

func FromNullString(n sql.NullString) *string {
	if !n.Valid {
		return nil
	}

	return &n.String
}

func WithDefaultNullString(a sql.NullString, b sql.NullString) sql.NullString {
	if a.Valid {
		return a
	}

	return b
}

func NewNullInt32(p *int32) sql.NullInt32 {
	if p == nil {
		return sql.NullInt32{}
	}

	return sql.NullInt32{Int32: *p, Valid: true}
}

func FromNullInt32(n sql.NullInt32) *int32 {
	if !n.Valid {
		return nil
	}

	return &n.Int32
}

func NewNullTime(p *time.Time) sql.NullTime {
	if p == nil {
		return sql.NullTime{}
	}

	if p.IsZero() {
		return sql.NullTime{}
	}

	return sql.NullTime{Time: *p, Valid: true}
}

func FromNullTime(n sql.NullTime) *time.Time {
	if !n.Valid {
		return nil
	}

	return &n.Time
}
