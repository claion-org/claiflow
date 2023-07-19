package global_variables

import (
	"database/sql"
	"fmt"
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

type GlobalVariable struct {
	Uuid    string    `json:"uuid"`
	Name    string    `json:"name"`
	Summary *string   `json:"summary,omitempty"` // (optional)
	Value   string    `json:"value"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

func NewGlobalVariable(m model.GlobalVariable) GlobalVariable {
	return GlobalVariable{
		Uuid:    m.UUID,
		Name:    m.Name,
		Summary: macro.FromNullString(m.Summary),
		Value:   m.Value,
		Created: m.Created,
		Updated: m.Updated,
	}
}

type UpdateValue struct {
	Value *string `json:"value"` // (optional)
}

// @Description Find GlobalVariables
// @Accept      x-www-form-urlencoded
// @Produce     json
// @Tags        /api/v1/global_variables
// @Router      /api/v1/global_variables [get]
// @Param       q query string false "query  github.com/claion-org/claiflow/pkg/server/database/vanilla/stmt/README.md"
// @Param       o query string false "order  github.com/claion-org/claiflow/pkg/server/database/vanilla/stmt/README.md"
// @@Param      p query string false "paging github.com/claion-org/claiflow/pkg/server/database/vanilla/stmt/README.md"
// @Success     200 {array} GlobalVariable
func FindGlobalVariables(ctx echo.Context) error {
	q, err := stmt.ConditionLexer.Parse(echov4.QueryParam(ctx)["q"])
	if err != nil && !logs.DeepCompare(err, stmt.ErrorInvalidArgumentEmptyString) {
		return echov4.HttpError(err, http.StatusBadRequest)
	}
	o, err := stmt.OrderLexer.Parse(echov4.QueryParam(ctx)["o"])
	if err != nil && !logs.DeepCompare(err, stmt.ErrorInvalidArgumentEmptyString) {
		return echov4.HttpError(err, http.StatusBadRequest)
	}
	// p, err := stmt.PaginationLexer.Parse(echov4.QueryParam(ctx)["p"])
	// if err != nil && !logs.DeepCompare(err, stmt.ErrorInvalidArgumentEmptyString) {
	// 	return echov4.HttpError(err, http.StatusBadRequest)
	// }
	// // default pagination
	// if p == nil {
	// 	p = stmt.Limit(state.PAGINATION_LIMIT())
	// }

	rsp, err := control.FindGlobalVariables(ctx.Request().Context(), q, o)
	if err != nil {
		err := errors.Wrapf(err, "query global variable")

		return echov4.HttpError(err, http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, []GlobalVariable(generic.Map(rsp, NewGlobalVariable)))
}

// @Description Get a GlobalVariable
// @Accept      json
// @Produce     json
// @Tags        /api/v1/global_variables
// @Router      /api/v1/global_variables/{uuid} [get]
// @Param       uuid path     string true "GlobalVariable UUID"
// @Success     200  {object} GlobalVariable
func GetGlobalVariable(ctx echo.Context) error {
	const (
		UUID = "uuid"
	)

	if len(echov4.PathParam(ctx)[UUID]) == 0 {
		err := fmt.Errorf("%w: path=%q", echov4.ErrorInvalidRequestPath, UUID)
		return echov4.HttpError(err, http.StatusBadRequest)
	}

	uuid := echov4.PathParam(ctx)[UUID]

	globvar, err := control.GetGlobalVariable(ctx.Request().Context(), uuid)
	if err != nil && err == sql.ErrNoRows {
		err := errors.Wrapf(err, "get global variable (uuid=%v)", uuid)

		return echov4.HttpError(err, http.StatusNotFound)
	}
	if err != nil && err != sql.ErrNoRows {
		err := errors.Wrapf(err, "get global variable")

		return echov4.HttpError(err, http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, GlobalVariable(NewGlobalVariable(*globvar)))
}

// @Description Update GlobalVariable Value
// @Accept      json
// @Produce     json
// @Tags        /api/v1/global_variables
// @Router      /api/v1/global_variables/{uuid} [put]
// @Param       uuid       path     string      true "GlobalVariable UUID"
// @Param       enviroment body     UpdateValue true "UpdateValue"
// @Success     200        {object} GlobalVariable
func UpdateGlobalVariableValue(ctx echo.Context) error {
	const (
		UUID = "uuid"
	)

	var body UpdateValue
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

	_globvar, err := control.GetGlobalVariable(ctx.Request().Context(), uuid)
	if err != nil && err == sql.ErrNoRows {
		err := errors.Wrapf(err, "get global variable (uuid=%v)", uuid)

		return echov4.HttpError(err, http.StatusNotFound)
	}
	if err != nil && err != sql.ErrNoRows {
		err := errors.Wrapf(err, "get global variable")

		return echov4.HttpError(err, http.StatusInternalServerError)
	}

	// valid value
	if body.Value != nil {
		for _, gv := range globvar.GlobVars {
			if gv.UUID() == _globvar.UUID {
				if err := gv.Clone().SetValue(*body.Value); err != nil {
					err := errors.Wrapf(err, "invalid value (uuid=%v, name=%v)", gv.UUID(), gv.Name())

					return echov4.HttpError(err, http.StatusBadRequest)
				}
			}
		}
	}

	// update a cluster
	updateColumns := make([]string, 0, 2)

	if body.Value != nil {
		_globvar.Value = *body.Value
		updateColumns = append(updateColumns, model.GlobalVariableFieldsValue.String())
	}

	if 0 < len(updateColumns) {
		_globvar.Updated = time.Now()
		updateColumns = append(updateColumns, model.GlobalVariableFieldsUpdated.String())
	}

	// something changed?
	if len(updateColumns) == 0 {
		err := errors.Wrapf(echov4.ErrorInvalidRequestParameter, "nothing to change")

		return echov4.HttpError(err, http.StatusBadRequest)
	}

	// save
	err = control.UpsertGlobalVariable(ctx.Request().Context(), _globvar, updateColumns)
	if err != nil {
		err := errors.Wrapf(err, "save the changes to the global variable")

		return echov4.HttpError(err, http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, GlobalVariable(NewGlobalVariable(*_globvar)))
}
