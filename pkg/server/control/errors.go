package control

import (
	"database/sql"
	"fmt"
)

var (
	ErrNoAffected = fmt.Errorf("sql: no affected")
	ErrNoRows     = sql.ErrNoRows
)
