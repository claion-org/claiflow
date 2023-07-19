package cluster

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/claion-org/claiflow/pkg/echov4"
	"github.com/claion-org/claiflow/pkg/server/control"
	"github.com/claion-org/claiflow/pkg/server/database/vanilla/stmt"
	"github.com/claion-org/claiflow/pkg/server/macro"
	"github.com/claion-org/claiflow/pkg/server/macro/generic"
	"github.com/claion-org/claiflow/pkg/server/macro/logs"
	"github.com/claion-org/claiflow/pkg/server/model"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type Create struct {
	Uuid    *string `json:"uuid,omitempty"` // (optional)
	Name    string  `json:"name"`
	Summary *string `json:"summary,omitempty"` // (optional)
}

type Update struct {
	Name    *string `json:"name,omitempty"`    // (optional)
	Summary *string `json:"summary,omitempty"` // (optional)
}

type ClusterClientSessionStatus struct {
	Alive            bool   `json:"alive"`
	ClientVersion    string `json:"clientVersion,omitempty"`
	ClientLibVersion string `json:"clientLibVersion,omitempty"`
}

type Cluster struct {
	Uuid    string    `json:"uuid"`
	Name    string    `json:"name"`
	Summary *string   `json:"summary,omitempty"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

func NewCluster(m model.Cluster) Cluster {
	return Cluster{
		Uuid:    m.UUID,
		Name:    m.Name,
		Summary: macro.FromNullString(m.Summary),
		Created: m.Created,
		Updated: m.Updated,
	}
}

// @Description Create a cluster
// @Accept      json
// @Produce     json
// @Tags        /api/v1/cluster
// @Router      /api/v1/cluster [post]
// @Param       cluster body     Create true "Create"
// @Success     200     {object} Cluster
func CreateCluster(ctx echo.Context) error {
	var body Create
	if err := ctx.Bind(&body); err != nil {
		err := errors.Wrapf(err, "bind request (body=%T)", body)

		return echov4.HttpError(err, http.StatusBadRequest)
	}

	if len(body.Name) == 0 {
		err := errors.Wrapf(echov4.ErrorInvalidRequestParameter, "param=%q", ".name")

		return echov4.HttpError(err, http.StatusBadRequest)
	}

	timeNow := time.Now()

	UUID := func() string {
		if body.Uuid == nil {
			return macro.NewUuidString()
		}

		return *body.Uuid
	}

	newCluster := model.Cluster{}
	newCluster.UUID = UUID()
	newCluster.Name = body.Name
	newCluster.Summary = macro.NewNullString(body.Summary)
	newCluster.Created = timeNow
	newCluster.Updated = timeNow

	if err := control.CreateCluster(ctx.Request().Context(), &newCluster); err != nil {
		err := errors.Wrapf(err, "save new cluster")

		return echov4.HttpError(err, http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, NewCluster(newCluster))
}

// @Description Find clusters
// @Accept      x-www-form-urlencoded
// @Produce     json
// @Tags        /api/v1/cluster
// @Router      /api/v1/cluster [get]
// @Param       q   query   string false "query  github.com/claion-org/claiflow/pkg/server/database/vanilla/stmt/README.md"
// @Param       o   query   string false "order  github.com/claion-org/claiflow/pkg/server/database/vanilla/stmt/README.md"
// @Param       p   query   string false "paging github.com/claion-org/claiflow/pkg/server/database/vanilla/stmt/README.md"
// @Success     200 {array} Cluster
func FindCluster(ctx echo.Context) error {
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

	records, err := control.FindCluster(ctx.Request().Context(), q, o, p)
	if err != nil {
		err := errors.Wrapf(err, "query cluster")

		return echov4.HttpError(err, http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, []Cluster(generic.Map(records, NewCluster)))
}

// @Description Get a cluster
// @Accept      json
// @Produce     json
// @Tags        /api/v1/cluster
// @Router      /api/v1/cluster/{uuid} [get]
// @Param       uuid path     string true "Cluster UUID"
// @Success     200  {object} Cluster
func GetCluster(ctx echo.Context) error {
	const (
		UUID = "uuid"
	)

	if len(echov4.PathParam(ctx)[UUID]) == 0 {
		err := errors.Wrapf(echov4.ErrorInvalidRequestPath, "path=%q", UUID)

		return echov4.HttpError(err, http.StatusBadRequest)
	}

	uuid := echov4.PathParam(ctx)[UUID]

	record, err := control.GetCluster(ctx.Request().Context(), uuid)
	if err != nil && err == sql.ErrNoRows {
		err := errors.Wrapf(err, "get cluster (uuid=%v)", uuid)

		return echov4.HttpError(err, http.StatusNotFound)
	}
	if err != nil && err != sql.ErrNoRows {
		err := errors.Wrapf(err, "get cluster")

		return echov4.HttpError(err, http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, Cluster(NewCluster(*record)))
}

// @Description Check Alive a Cluster Client Session
// @Accept      json
// @Produce     json
// @Tags        /api/v1/cluster
// @Router      /api/v1/cluster/{cluster_uuid}/session/alive [get]
// @Param       cluster_uuid path     string true "Cluster UUID"
// @Success     200          {object} ClusterClientSessionStatus
func GetClusterClientSessionAlive(ctx echo.Context) error {
	const __CLUSTER_UUID__ = "cluster_uuid"

	if len(echov4.PathParam(ctx)[__CLUSTER_UUID__]) == 0 {
		err := errors.Wrapf(echov4.ErrorInvalidRequestPath, "path=%q", __CLUSTER_UUID__)

		return echov4.HttpError(err, http.StatusBadRequest)
	}

	clusterUUID := echov4.PathParam(ctx)[__CLUSTER_UUID__]

	aliveCond := stmt.And(
		stmt.Equal(model.ClusterClientSessionFieldsClusterUuid.String(), clusterUUID),
	)
	aliveOrder := stmt.Desc(model.ClusterClientSessionFieldsExpirationTime.String())
	alivePage := stmt.Limit(1)

	clusterClientSessions, err := control.FindClusterClientSession(ctx.Request().Context(), aliveCond, aliveOrder, alivePage)
	if err != nil {
		err := errors.Wrapf(err, "query cluster client session")

		return echov4.HttpError(err, http.StatusInternalServerError)
	}

	if len(clusterClientSessions) == 0 {
		return ctx.JSON(http.StatusOK, ClusterClientSessionStatus{
			Alive: false,
		})
	}

	var clusterClientSession model.ClusterClientSession
	for i := range clusterClientSessions {
		clusterClientSession = clusterClientSessions[i]
	}

	var alive bool = false
	if !clusterClientSession.ExpirationTime.IsZero() {
		alive = time.Now().Before(clusterClientSession.ExpirationTime)
	}

	var clientClaims model.ClusterClientSessionClaim
	token, _, err := jwt.NewParser().ParseUnverified(clusterClientSession.Token, &clientClaims)
	if err != nil {
		err := errors.New("parse cluster session token")

		return echov4.HttpError(err, http.StatusBadRequest)
	}

	if _, ok := token.Claims.(*model.ClusterClientSessionClaim); !ok {
		err := errors.Errorf("invalid claims (type=%T)", token.Claims)

		return echov4.HttpError(err, http.StatusBadRequest)
	}

	return ctx.JSON(http.StatusOK, ClusterClientSessionStatus{
		Alive:            alive,
		ClientVersion:    clientClaims.ClientVersion,
		ClientLibVersion: clientClaims.ClientLibVersion,
	})
}

// @Description Update a cluster
// @Accept      json
// @Produce     json
// @Tags        /api/v1/cluster
// @Router      /api/v1/cluster/{uuid} [put]
// @Param       uuid    path     string true "Cluster UUID"
// @Param       cluster body     Update true "Update"
// @Success     200     {object} Cluster
func UpdateCluster(ctx echo.Context) error {
	const (
		UUID = "uuid"
	)

	var body Update
	if err := ctx.Bind(&body); err != nil {
		err := errors.Wrapf(err, "bind request (body=%T)", body)

		return echov4.HttpError(err, http.StatusBadRequest)
	}

	if len(echov4.PathParam(ctx)[UUID]) == 0 {
		err := errors.Wrapf(echov4.ErrorInvalidRequestPath, "path=%q", UUID)

		return echov4.HttpError(err, http.StatusBadRequest)
	}

	// get a cluster
	uuid := echov4.PathParam(ctx)[UUID]

	cluster, err := control.GetCluster(ctx.Request().Context(), uuid)
	if err != nil && err == sql.ErrNoRows {
		err := errors.Wrapf(err, "get cluster (uuid=%v)", uuid)

		return echov4.HttpError(err, http.StatusNotFound)
	}
	if err != nil && err != sql.ErrNoRows {
		err := errors.Wrapf(err, "get cluster")

		return echov4.HttpError(err, http.StatusInternalServerError)
	}

	// update a cluster
	updateColumns := make([]string, 0, 3)

	if body.Name != nil {
		cluster.Name = *body.Name
		updateColumns = append(updateColumns, model.ClusterFieldsName.String())
	}

	if body.Summary != nil {
		cluster.Summary = macro.NewNullString(body.Summary)
		updateColumns = append(updateColumns, model.ClusterFieldsSummary.String())
	}

	if 0 < len(updateColumns) {
		cluster.Updated = time.Now()
		updateColumns = append(updateColumns, model.ClusterFieldsUpdated.String())
	}

	// something changed?
	if len(updateColumns) == 0 {
		err := errors.Wrapf(echov4.ErrorInvalidRequestParameter, "nothing to change")

		return echov4.HttpError(err, http.StatusBadRequest)
	}

	// save
	if err := control.UpsertCluster(ctx.Request().Context(), cluster, updateColumns); err != nil {
		err := errors.Wrapf(err, "save the changes to the cluster")

		return echov4.HttpError(err, http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, Cluster(NewCluster(*cluster)))
}

// @Description Delete a cluster
// @Accept      json
// @Produce     json
// @Tags        /api/v1/cluster
// @Router      /api/v1/cluster/{uuid} [delete]
// @Param       uuid path string true "Cluster UUID"
// @Success     200
func DeleteCluster(ctx echo.Context) error {
	const (
		UUID = "uuid"
	)

	if len(echov4.PathParam(ctx)[UUID]) == 0 {
		err := errors.Wrapf(echov4.ErrorInvalidRequestPath, "path=%q", UUID)

		return echov4.HttpError(err, http.StatusBadRequest)
	}

	uuid := echov4.PathParam(ctx)[UUID]

	if err := control.DeleteCluster(ctx.Request().Context(), uuid); err != nil {
		err := errors.Wrapf(err, "failed to delete a cluster")

		return echov4.HttpError(err, http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, echov4.OK())
}
