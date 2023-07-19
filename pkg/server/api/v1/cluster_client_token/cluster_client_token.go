package cluster_client_token

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
	"github.com/claion-org/claiflow/pkg/server/status/globvar"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type Create struct {
	Uuid        *string `json:"uuid,omitempty"` // (optional)
	Name        string  `json:"name"`
	Summary     *string `json:"summary,omitempty"` // (optional)
	ClusterUuid string  `json:"cluster_uuid"`
	Token       *string `json:"token,omitempty"` // (optional)
}

type Update struct {
	Name           *string    `json:"name,omitempty"`            // (optional)
	Summary        *string    `json:"summary,omitempty"`         // (optional)
	Token          *string    `json:"token,omitempty"`           // (optional)
	IssuedAtTime   *time.Time `json:"issued_at_time,omitempty"`  // (optional)
	ExpirationTime *time.Time `json:"expiration_time,omitempty"` // (optional)
}

type UpdateRefresh struct {
	ExpirationTime *time.Time `json:"expiration_time,omitempty"` // (optional) (force)
}

type ClusterClientToken struct {
	Uuid           string    `json:"uuid"`
	Name           string    `json:"name"`
	Summary        *string   `json:"summary,omitempty"`
	ClusterUuid    string    `json:"cluster_uuid"`
	Token          string    `json:"token"`
	IssuedAtTime   time.Time `json:"issued_at_time"`
	ExpirationTime time.Time `json:"expiration_time"`
	Created        time.Time `json:"created"`
	Updated        time.Time `json:"updated"`
}

func NewClusterClientToken(m model.ClusterClientToken) ClusterClientToken {
	return ClusterClientToken{
		Uuid:           m.UUID,
		Name:           m.Name,
		Summary:        macro.FromNullString(m.Summary),
		ClusterUuid:    m.ClusterUUID,
		Token:          m.Token,
		IssuedAtTime:   m.IssuedAtTime,
		ExpirationTime: m.ExpirationTime,
		Created:        m.Created,
		Updated:        m.Updated,
	}
}

// @Description Create a Cluster Client Token
// @Accept      json
// @Produce     json
// @Tags        /api/v1/cluster_token
// @Router      /api/v1/cluster_token [post]
// @Param       object body     Create true "Create"
// @Success     200    {object} ClusterClientToken
func CreateClusterClientToken(ctx echo.Context) error {
	var body Create
	if err := ctx.Bind(&body); err != nil {
		err := errors.Wrapf(err, "bind request (body=%T)", body)

		return echov4.HttpError(err, http.StatusBadRequest)
	}

	if len(body.Name) == 0 {
		err := errors.Wrapf(echov4.ErrorInvalidRequestParameter, "param=%q", ".name")

		return echov4.HttpError(err, http.StatusBadRequest)
	}

	if len(body.ClusterUuid) == 0 {
		err := errors.Wrapf(echov4.ErrorInvalidRequestParameter, "param=%q", ".cluster_uuid")

		return echov4.HttpError(err, http.StatusBadRequest)
	}

	timeNow := time.Now()

	UUID := func() string {
		if body.Uuid == nil {
			return macro.NewUuidString()
		}

		return *body.Uuid
	}

	Token := func() string {
		if body.Token == nil {
			return macro.NewUuidString()
		}

		return *body.Token
	}

	//property
	newToken := model.ClusterClientToken{}
	newToken.UUID = UUID()
	newToken.Name = body.Name
	newToken.Summary = macro.NewNullString(body.Summary)
	newToken.ClusterUUID = body.ClusterUuid
	newToken.IssuedAtTime = IssuedAtTime(timeNow)
	newToken.ExpirationTime = ExpirationTime(timeNow)
	newToken.Token = Token()
	newToken.Created = timeNow
	newToken.Updated = timeNow

	if err := control.CreateClusterClientToken(ctx.Request().Context(), &newToken); err != nil {
		err := errors.Wrapf(err, "save new cluster client token")

		return echov4.HttpError(err, http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, ClusterClientToken(NewClusterClientToken(newToken)))
}

// @Description Find Cluster Client Tokens
// @Accept      x-www-form-urlencoded
// @Produce     json
// @Tags        /api/v1/cluster_token
// @Router      /api/v1/cluster_token [get]
// @Param       q   query   string false "query  github.com/claion-org/claiflow/pkg/server/database/vanilla/stmt/README.md"
// @Param       o   query   string false "order  github.com/claion-org/claiflow/pkg/server/database/vanilla/stmt/README.md"
// @Param       p   query   string false "paging github.com/claion-org/claiflow/pkg/server/database/vanilla/stmt/README.md"
// @Success     200 {array} ClusterClientToken
func FindClusterClientToken(ctx echo.Context) error {
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

	records, err := control.FindClusterClientToken(ctx.Request().Context(), q, o, p)
	if err != nil {
		err := errors.Wrapf(err, "query cluster client token")

		return echov4.HttpError(err, http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, []ClusterClientToken(generic.Map(records, NewClusterClientToken)))
}

// @Description Get a Cluster Client Token
// @Accept      json
// @Produce     json
// @Tags        /api/v1/cluster_token
// @Router      /api/v1/cluster_token/{uuid} [get]
// @Param       uuid path     string true "ClusterClientToken Uuid"
// @Success     200  {object} ClusterClientToken
func GetClusterClientToken(ctx echo.Context) (err error) {
	const (
		UUID = "uuid"
	)

	if len(echov4.PathParam(ctx)[UUID]) == 0 {
		err := errors.Wrapf(echov4.ErrorInvalidRequestPath, "path=%q", UUID)

		return echov4.HttpError(err, http.StatusBadRequest)
	}

	uuid := echov4.PathParam(ctx)[UUID]

	record, err := control.GetClusterClientToken(ctx.Request().Context(), uuid)
	if err != nil && err == sql.ErrNoRows {
		err := errors.Wrapf(err, "get cluster client token (uuid=%v)", uuid)

		return echov4.HttpError(err, http.StatusNotFound)
	}
	if err != nil && err != sql.ErrNoRows {
		err := errors.Wrapf(err, "get cluster client token")

		return echov4.HttpError(err, http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, ClusterClientToken(NewClusterClientToken(*record)))
}

// @Description Update a Cluster Client Token
// @Accept      json
// @Produce     json
// @Tags        /api/v1/cluster_token
// @Router      /api/v1/cluster_token/{uuid} [put]
// @Param       uuid   path     string true "ClusterClientToken UUID"
// @Param       object body     Update true "Update"
// @Success     200    {object} ClusterClientToken
func UpdateClusterClientToken(ctx echo.Context) (err error) {
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

	uuid := echov4.PathParam(ctx)[UUID]

	// get cluster client token
	clientToken, err := control.GetClusterClientToken(ctx.Request().Context(), uuid)
	if err != nil && err == sql.ErrNoRows {
		err := errors.Wrapf(err, "get cluster client token (uuid=%v)", uuid)

		return echov4.HttpError(err, http.StatusNotFound)
	}
	if err != nil && err != sql.ErrNoRows {
		err := errors.Wrapf(err, "get cluster client token")

		return echov4.HttpError(err, http.StatusInternalServerError)
	}

	// update a cluster client token
	updateColumns := make([]string, 0, 6)

	if body.Name != nil {
		clientToken.Name = *body.Name
		updateColumns = append(updateColumns, model.ClusterClientTokenFieldsName.String())
	}

	if body.Summary != nil {
		clientToken.Summary = macro.NewNullString(body.Summary)
		updateColumns = append(updateColumns, model.ClusterClientTokenFieldsSummary.String())
	}

	if body.Token != nil {
		clientToken.Token = *body.Token
		updateColumns = append(updateColumns, model.ClusterClientTokenFieldsToken.String())
	}

	if body.IssuedAtTime != nil {
		clientToken.IssuedAtTime = *body.IssuedAtTime
		updateColumns = append(updateColumns, model.ClusterClientTokenFieldsIssuedAtTime.String())
	}

	if body.ExpirationTime != nil {
		clientToken.ExpirationTime = *body.ExpirationTime
		updateColumns = append(updateColumns, model.ClusterClientTokenFieldsExpirationTime.String())
	}

	if 0 < len(updateColumns) {
		clientToken.Updated = time.Now()
		updateColumns = append(updateColumns, model.ClusterClientTokenFieldsUpdated.String())
	}

	// something changed?
	if len(updateColumns) == 0 {
		err := errors.Wrapf(echov4.ErrorInvalidRequestParameter, "nothing to change")

		return echov4.HttpError(err, http.StatusBadRequest)
	}

	// save
	if err := control.UpsertClusterClientToken(ctx.Request().Context(), clientToken, updateColumns); err != nil {
		err := errors.Wrapf(err, "save the changes to the cluster client token")

		return echov4.HttpError(err, http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, ClusterClientToken(NewClusterClientToken(*clientToken)))
}

// @Description Refresh Time of a Cluster Client Token
// @Accept      json
// @Produce     json
// @Tags        /api/v1/cluster_token
// @Router      /api/v1/cluster_token/{uuid}/refresh [put]
// @Param       uuid path     string true "ClusterClientToken UUID"
// @Success     200  {object} ClusterClientToken
func RefreshClusterClientToken(ctx echo.Context) (err error) {
	const (
		UUID = "uuid"
	)

	if len(echov4.PathParam(ctx)[UUID]) == 0 {
		err := errors.Wrapf(echov4.ErrorInvalidRequestPath, "path=%q", UUID)

		return echov4.HttpError(err, http.StatusBadRequest)
	}

	uuid := echov4.PathParam(ctx)[UUID]

	// get cluster client token
	clientToken, err := control.GetClusterClientToken(ctx.Request().Context(), uuid)
	if err != nil && err == sql.ErrNoRows {
		err := errors.Wrapf(err, "get cluster client token (uuid=%v)", uuid)

		return echov4.HttpError(err, http.StatusNotFound)
	}
	if err != nil && err != sql.ErrNoRows {
		err := errors.Wrapf(err, "get cluster client token")

		return echov4.HttpError(err, http.StatusInternalServerError)
	}

	timeNow := time.Now()

	// update a cluster client token
	updateColumns := make([]string, 0, 3)

	if false {
		clientToken.IssuedAtTime = IssuedAtTime(timeNow)
		updateColumns = append(updateColumns, model.ClusterClientTokenFieldsIssuedAtTime.String())
	}

	clientToken.ExpirationTime = ExpirationTime(timeNow)
	updateColumns = append(updateColumns, model.ClusterClientTokenFieldsExpirationTime.String())

	clientToken.Updated = time.Now()
	updateColumns = append(updateColumns, model.ClusterClientTokenFieldsUpdated.String())

	// save
	if err := control.UpsertClusterClientToken(ctx.Request().Context(), clientToken, updateColumns); err != nil {
		err := errors.Wrapf(err, "save the changes to the cluster client token")

		return echov4.HttpError(err, http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, ClusterClientToken(NewClusterClientToken(*clientToken)))
}

// @Description Expire a Cluster Client Token
// @Accept      json
// @Produce     json
// @Tags        /api/v1/cluster_token
// @Router      /api/v1/cluster_token/{uuid}/expire [put]
// @Param       uuid path     string true "ClusterClientToken UUID"
// @Success     200  {object} ClusterClientToken
func ExpireClusterClientToken(ctx echo.Context) (err error) {
	const (
		UUID = "uuid"
	)

	if len(echov4.PathParam(ctx)[UUID]) == 0 {
		err := errors.Wrapf(echov4.ErrorInvalidRequestPath, "path=%q", UUID)

		return echov4.HttpError(err, http.StatusBadRequest)
	}

	uuid := echov4.PathParam(ctx)[UUID]

	// get cluster client token
	clientToken, err := control.GetClusterClientToken(ctx.Request().Context(), uuid)
	if err != nil && err == sql.ErrNoRows {
		err := errors.Wrapf(err, "get cluster client token (uuid=%v)", uuid)

		return echov4.HttpError(err, http.StatusNotFound)
	}
	if err != nil && err != sql.ErrNoRows {
		err := errors.Wrapf(err, "get cluster client token")

		return echov4.HttpError(err, http.StatusInternalServerError)
	}

	timeNow := time.Now()

	// update a cluster client token
	updateColumns := make([]string, 0, 3)

	if false {
		clientToken.IssuedAtTime = IssuedAtTime(timeNow)
		updateColumns = append(updateColumns, model.ClusterClientTokenFieldsIssuedAtTime.String())
	}

	clientToken.ExpirationTime = SetExpiration(timeNow)
	updateColumns = append(updateColumns, model.ClusterClientTokenFieldsExpirationTime.String())

	clientToken.Updated = time.Now()
	updateColumns = append(updateColumns, model.ClusterClientTokenFieldsUpdated.String())

	// save
	if err := control.UpsertClusterClientToken(ctx.Request().Context(), clientToken, updateColumns); err != nil {
		err := errors.Wrapf(err, "save the changes to the cluster client token")

		return echov4.HttpError(err, http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, ClusterClientToken(NewClusterClientToken(*clientToken)))
}

// @Description Delete a Cluster Client Token
// @Accept      json
// @Produce     json
// @Tags        /api/v1/cluster_token
// @Router      /api/v1/cluster_token/{uuid} [delete]
// @Param       uuid path string true "ClusterClientToken UUID"
// @Success     200
func DeleteClusterClinetToken(ctx echo.Context) (err error) {
	const (
		UUID = "uuid"
	)

	if len(echov4.PathParam(ctx)[UUID]) == 0 {
		err := errors.Wrapf(echov4.ErrorInvalidRequestPath, "path=%q", UUID)

		return echov4.HttpError(err, http.StatusBadRequest)
	}

	uuid := echov4.PathParam(ctx)[UUID]

	if err := control.DeleteClusterClientToken(ctx.Request().Context(), uuid); err != nil {
		err := errors.Wrapf(err, "failed to delete a cluster client token")

		return echov4.HttpError(err, http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, echov4.OK())
}

func IssuedAtTime(t time.Time) time.Time {
	return t
}

func ExpirationTime(t time.Time) time.Time {
	return globvar.ClusterTokenExpirationTime.Add(t)
}

func SetExpiration(t time.Time) time.Time {
	return t
}
