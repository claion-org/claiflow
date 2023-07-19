package globvar

import (
	"database/sql"
	"time"

	"github.com/claion-org/claiflow/pkg/server/model"
)

var GlobVars = []GlobVar{
	&ClusterTokenExpirationTime,
	&ClientSessionSignatureSecret,
	&ClientSessionExpirationTime,
	&ClientConfigServiceValidityPeriod,
}

type GlobVar interface {
	UUID() string
	Name() string
	Summary() string
	GetValue() string
	SetValue(s string) error
	Clone() GlobVar
}

func ToModel(gv GlobVar, t time.Time) *model.GlobalVariable {
	globvar := model.GlobalVariable{}
	globvar.UUID = gv.UUID()
	globvar.Name = gv.Name()
	globvar.Summary = sql.NullString{String: gv.Summary(), Valid: true}
	globvar.Value = gv.GetValue()
	globvar.Created = time.Now()
	globvar.Updated = time.Now()

	return &globvar
}
