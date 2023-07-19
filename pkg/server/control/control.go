package control

import (
	"database/sql"

	"github.com/claion-org/claiflow/pkg/server/database/vanilla/excute"
)

const (
	INIT_SLICE_CAP = 10
)

var Driver func() excute.SqlExcutor = nil
var Database func() *sql.DB = nil

// var ChannelClient func() *channels.Client = nil
