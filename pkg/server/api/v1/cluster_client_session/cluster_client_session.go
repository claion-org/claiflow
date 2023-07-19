package cluster_client_session

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/claion-org/claiflow/pkg/echov4"
	"github.com/claion-org/claiflow/pkg/server/api/v1/cluster"
	"github.com/claion-org/claiflow/pkg/server/control"
	"github.com/claion-org/claiflow/pkg/server/database/vanilla/stmt"
	"github.com/claion-org/claiflow/pkg/server/macro/generic"
	"github.com/claion-org/claiflow/pkg/server/macro/logs"
	"github.com/claion-org/claiflow/pkg/server/model"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type ClusterClientSession struct {
	ID             int64     `json:"id"`
	Uuid           string    `json:"uuid"`
	ClusterUuid    string    `json:"cluster_uuid"`
	Token          string    `json:"token"`
	IssuedAtTime   time.Time `json:"issued_at_time"`
	ExpirationTime time.Time `json:"expiration_time"`
	Created        time.Time `json:"created"`
	Updated        time.Time `json:"updated"`
	Deleted        time.Time `json:"deleted"`
}

func NewClusterClientSession(m model.ClusterClientSession) ClusterClientSession {
	return ClusterClientSession{
		Uuid:           m.UUID,
		ClusterUuid:    m.ClusterUUID,
		Token:          m.Token,
		IssuedAtTime:   m.IssuedAtTime,
		ExpirationTime: m.ExpirationTime,
		Created:        m.Created,
		Updated:        m.Updated,
	}
}

// @Description Find Cluster Client Sessions
// @Accept      x-www-form-urlencoded
// @Produce     json
// @Tags        /api/v1/session
// @Router      /api/v1/session [get]
// @Param       q   query   string false "query  github.com/claion-org/claiflow/pkg/database/vanilla/stmt/README.md"
// @Param       o   query   string false "order  github.com/claion-org/claiflow/pkg/database/vanilla/stmt/README.md"
// @Param       p   query   string false "paging github.com/claion-org/claiflow/pkg/database/vanilla/stmt/README.md"
// @Success     200 {array} ClusterClientSession
func FindClusterClientSession(ctx echo.Context) error {
	q, err := stmt.ConditionLexer.Parse(echov4.QueryParam(ctx)["q"])
	if err != nil && !logs.DeepCompare(err, stmt.ErrorInvalidArgumentEmptyString) {
		return echov4.HttpError(err, http.StatusBadRequest)
	}
	o, err := stmt.OrderLexer.Parse(echov4.QueryParam(ctx)["o"])
	if err != nil && !logs.DeepCompare(err, stmt.ErrorInvalidArgumentEmptyString) {
		return echov4.HttpError(err, http.StatusBadRequest)
	}
	p, err := stmt.PaginationLexer.Parse(echov4.QueryParam(ctx)["p"])
	if err != nil && !logs.DeepCompare(err, stmt.ErrorInvalidArgumentEmptyString) {
		return echov4.HttpError(err, http.StatusBadRequest)
	}

	records, err := control.FindClusterClientSession(ctx.Request().Context(), q, o, p)
	if err != nil {
		err := errors.Wrapf(err, "query cluster client session")

		return echov4.HttpError(err, http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, []ClusterClientSession(generic.Map(records, NewClusterClientSession)))
}

// @Description Get a Cluster Client Session
// @Accept      json
// @Produce     json
// @Tags        /api/v1/session
// @Router      /api/v1/session/{uuid} [get]
// @Param       uuid path     string true "ClusterClientSession UUID"
// @Success     200  {object} ClusterClientSession
func GetClusterClientSession(ctx echo.Context) error {
	const (
		UUID = "uuid"
	)

	if len(echov4.PathParam(ctx)[UUID]) == 0 {
		err := errors.Wrapf(echov4.ErrorInvalidRequestPath, "path=%q", UUID)

		return echov4.HttpError(err, http.StatusBadRequest)
	}

	uuid := echov4.PathParam(ctx)[UUID]

	record, err := control.GetClusterClientSession(ctx.Request().Context(), uuid)
	if err != nil && err == sql.ErrNoRows {
		err := errors.Wrapf(err, "get cluster client session (uuid=%v)", uuid)

		return echov4.HttpError(err, http.StatusNotFound)
	}
	if err != nil && err != sql.ErrNoRows {
		err := errors.Wrapf(err, "get cluster client session")

		return echov4.HttpError(err, http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, ClusterClientSession(NewClusterClientSession(*record)))
}

// @Description Delete a Cluster Client Session
// @Accept      json
// @Produce     json
// @Tags        /api/v1/session
// @Router      /api/v1/session/{uuid} [delete]
// @Param       uuid path string true "ClusterClientSession UUID"
// @Success     200
func DeleteClusterClientSession(ctx echo.Context) error {
	const (
		UUID = "uuid"
	)

	if len(echov4.PathParam(ctx)[UUID]) == 0 {
		err := errors.Wrapf(echov4.ErrorInvalidRequestPath, "path=%q", UUID)

		return echov4.HttpError(err, http.StatusBadRequest)
	}

	uuid := echov4.PathParam(ctx)[UUID]

	if err := control.DeleteClusterClientSession(ctx.Request().Context(), uuid); err != nil {
		err := errors.Wrapf(err, "failed to delete a cluster client session")

		return echov4.HttpError(err, http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, echov4.OK())
}

// @Deprecated
// @Description Check Alive a Cluster Client Session
// @Accept      json
// @Produce     json
// @Tags        /api/v1/session
// @Router      /api/v1/session/alive [get]
// @Param       cluster_uuid path     string true "Cluster UUID"
// @Success     200          {object} cluster.ClusterClientSessionStatus
func GetClusterClientSessionAlive(ctx echo.Context) error {
	return cluster.GetClusterClientSessionAlive(ctx)
}
