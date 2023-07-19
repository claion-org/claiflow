package migrations

import (
	"bytes"
	"embed"
)

//go:embed migrations
var Migrations embed.FS

//go:embed migrations/mysql/latest
var mysqlLatest []byte

var Latests = map[string]*Latest{
	"migrations/mysql": new(Latest).SetReader(bytes.NewBuffer(mysqlLatest)),
}
