package webhook

import (
	"database/sql"
	"net/http"
	"strings"
	"time"

	"github.com/claion-org/claiflow/pkg/echov4"
	"github.com/claion-org/claiflow/pkg/server/control"
	"github.com/claion-org/claiflow/pkg/server/database/vanilla/stmt"
	"github.com/claion-org/claiflow/pkg/server/macro"
	"github.com/claion-org/claiflow/pkg/server/macro/generic"
	"github.com/claion-org/claiflow/pkg/server/macro/logs"
	"github.com/claion-org/claiflow/pkg/server/model"
	"github.com/claion-org/claiflow/pkg/webhook"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type Header = map[string][]string

type Create struct {
	Uuid               *string `json:"uuid,omitempty"` // (optional)
	Name               string  `json:"name"`
	Summary            *string `json:"summary,omitempty"` // (optional)
	URL                string  `json:"url"`
	Method             string  `json:"method"`
	Headers            Header  `json:"headers"`                      // (optional)
	Timeout            *string `json:"timeout,omitempty"`            // (optional)
	ConditionValidator string  `json:"conditionValidator,omitempty"` // (optional)
	ConditionFilter    *string `json:"conditionFilter,omitempty"`    // (optional)
}

type Update struct {
	Name               *string `json:"name,omitempty"`               // (optional)
	Summary            *string `json:"summary,omitempty"`            // (optional)
	URL                *string `json:"url,omitempty"`                // (optional)
	Method             *string `json:"method,omitempty"`             // (optional)
	Headers            Header  `json:"headers,omitempty"`            // (optional)
	Timeout            *string `json:"timeout,omitempty"`            // (optional)
	ConditionValidator *string `json:"conditionValidator,omitempty"` // (optional)
	ConditionFilter    *string `json:"conditionFilter,omitempty"`    // (optional)
}

type Webhook struct {
	Uuid               string    `json:"uuid"`
	Name               string    `json:"name"`
	Summary            *string   `json:"summary,omitempty"`
	URL                string    `json:"url"`
	Method             string    `json:"method"`
	Headers            Header    `json:"headers,omitempty" swaggertype:"object"`
	Timeout            string    `json:"timeout,omitempty"`
	ConditionValidator string    `json:"conditionValidator,omitempty"`
	ConditionFilter    *string   `json:"conditionFilter,omitempty"`
	Created            time.Time `json:"created"`
	Updated            time.Time `json:"updated"`
}

func NewWebhook(m model.Webhook) Webhook {
	return Webhook{
		Uuid:               m.UUID,
		Name:               m.Name,
		Summary:            macro.FromNullString(m.Summary),
		URL:                m.URL,
		Method:             m.Method,
		Headers:            m.Headers,
		Timeout:            (time.Duration(m.Timeout.Int32) * time.Second).String(),
		ConditionValidator: model.WebhookConditionValidator(m.ConditionValidator.Int32).String(),
		ConditionFilter:    macro.FromNullString(m.ConditionFilter),
		Created:            m.Created,
		Updated:            m.Updated,
	}
}

// @Description Create a webhook
// @Accept      json
// @Produce     json
// @Tags        /api/v1/webhook
// @Router      /api/v1/webhook [post]
// @Param       webhook body     Create true "Create"
// @Success     200     {object} Webhook
func CreateWebhook(ctx echo.Context) error {
	var body Create
	if err := ctx.Bind(&body); err != nil {
		err := errors.Wrapf(err, "bind request (body=%T)", body)

		return echov4.HttpError(err, http.StatusBadRequest)
	}

	if len(body.Name) == 0 {
		err := errors.Wrapf(echov4.ErrorInvalidRequestParameter, "param=%q", ".name")

		return echov4.HttpError(err, http.StatusBadRequest)
	}

	if len(body.URL) == 0 {
		err := errors.Wrapf(echov4.ErrorInvalidRequestParameter, "param=%q", ".url")

		return echov4.HttpError(err, http.StatusBadRequest)
	}

	if len(body.Method) == 0 {
		err := errors.Wrapf(echov4.ErrorInvalidRequestParameter, "param=%q", ".method")

		return echov4.HttpError(err, http.StatusBadRequest)
	}

	timeout, err := ParseTimeout(body.Timeout)
	if err != nil {
		err := errors.Wrapf(echov4.ErrorInvalidRequestParameter, "param=%q", ".timeout")

		return echov4.HttpError(err, http.StatusBadRequest)
	}

	timeNow := time.Now()

	UUID := func() string {
		if body.Uuid == nil {
			return macro.NewUuidString()
		}

		return *body.Uuid
	}

	newWebhook := model.Webhook{}
	newWebhook.UUID = UUID()
	newWebhook.Name = body.Name
	newWebhook.Summary = macro.NewNullString(body.Summary)
	newWebhook.URL = body.URL
	newWebhook.Method = strings.ToUpper(body.Method)
	newWebhook.Headers = body.Headers
	newWebhook.Timeout = timeout
	newWebhook.ConditionValidator = sql.NullInt32{Int32: int32(generic.Left(model.ParseWebhookConditionValidator(body.ConditionValidator))), Valid: true}
	newWebhook.ConditionFilter = macro.NewNullString(body.ConditionFilter)
	newWebhook.Created = timeNow
	newWebhook.Updated = timeNow

	// check condition validator and condition filter
	err = webhook.CheckConditionFilter(
		model.WebhookConditionValidator(newWebhook.ConditionValidator.Int32),
		newWebhook.ConditionFilter.String)
	if err != nil {
		err := errors.Wrapf(err, "invalid condition validator")

		return echov4.HttpError(err, http.StatusBadRequest)
	}

	if err := control.CreateWebhook(ctx.Request().Context(), &newWebhook); err != nil {
		err := errors.Wrapf(err, "save new webhook")

		return echov4.HttpError(err, http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, NewWebhook(newWebhook))
}

// @Description Find webhooks
// @Accept      x-www-form-urlencoded
// @Produce     json
// @Tags        /api/v1/webhook
// @Router      /api/v1/webhook [get]
// @Param       q   query   string false "query  github.com/claion-org/claiflow/pkg/server/database/vanilla/stmt/README.md"
// @Param       o   query   string false "order  github.com/claion-org/claiflow/pkg/server/database/vanilla/stmt/README.md"
// @Param       p   query   string false "paging github.com/claion-org/claiflow/pkg/server/database/vanilla/stmt/README.md"
// @Success     200 {array} Webhook
func FindWebhook(ctx echo.Context) error {
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

	records, err := control.FindWebhook(ctx.Request().Context(), q, o, p)
	if err != nil {
		err := errors.Wrapf(err, "query webhook")

		return echov4.HttpError(err, http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, []Webhook(generic.Map(records, NewWebhook)))
}

// @Description Get a webhook
// @Accept      json
// @Produce     json
// @Tags        /api/v1/webhook
// @Router      /api/v1/webhook/{uuid} [get]
// @Param       uuid path     string true "webhook UUID"
// @Success     200  {object} Webhook
func GetWebhook(ctx echo.Context) error {
	const (
		UUID = "uuid"
	)

	if len(echov4.PathParam(ctx)[UUID]) == 0 {
		err := errors.Wrapf(echov4.ErrorInvalidRequestPath, "path=%q", UUID)

		return echov4.HttpError(err, http.StatusBadRequest)
	}

	uuid := echov4.PathParam(ctx)[UUID]

	record, err := control.GetWebhook(ctx.Request().Context(), uuid)
	if err != nil && err == sql.ErrNoRows {
		err := errors.Wrapf(err, "get webhook (uuid=%v)", uuid)

		return echov4.HttpError(err, http.StatusNotFound)
	}
	if err != nil && err != sql.ErrNoRows {
		err := errors.Wrapf(err, "get webhook")

		return echov4.HttpError(err, http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, Webhook(NewWebhook(*record)))
}

// @Description Update a webhook
// @Accept      json
// @Produce     json
// @Tags        /api/v1/webhook
// @Router      /api/v1/webhook/{uuid} [put]
// @Param       uuid    path     string true "Webhook UUID"
// @Param       webhook body     Update true "Update"
// @Success     200     {object} Webhook
func UpdateWebhook(ctx echo.Context) error {
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

	timeout, err := ParseTimeout(body.Timeout)
	if err != nil {
		err := errors.Wrapf(echov4.ErrorInvalidRequestParameter, "param=%q", ".timeout")

		return echov4.HttpError(err, http.StatusBadRequest)
	}

	// get a webhook
	uuid := echov4.PathParam(ctx)[UUID]

	webhook_, err := control.GetWebhook(ctx.Request().Context(), uuid)
	if err != nil && err == sql.ErrNoRows {
		err := errors.Wrapf(err, "get webhook (uuid=%v)", uuid)

		return echov4.HttpError(err, http.StatusNotFound)
	}
	if err != nil && err != sql.ErrNoRows {
		err := errors.Wrapf(err, "get webhook")

		return echov4.HttpError(err, http.StatusInternalServerError)
	}

	// update a webhook
	updateColumns := make([]string, 0, 9)

	if body.Name != nil {
		webhook_.Name = *body.Name
		updateColumns = append(updateColumns, model.WebhookFieldsName.String())
	}

	if body.Summary != nil {
		webhook_.Summary = macro.NewNullString(body.Summary)
		updateColumns = append(updateColumns, model.WebhookFieldsSummary.String())
	}

	if body.URL != nil {
		webhook_.URL = *body.URL
		updateColumns = append(updateColumns, model.WebhookFieldsUrl.String())
	}

	if body.Method != nil {
		webhook_.Method = *body.Method
		updateColumns = append(updateColumns, model.WebhookFieldsMethod.String())
	}

	if body.Headers != nil {
		webhook_.Headers = body.Headers
		updateColumns = append(updateColumns, model.WebhookFieldsHeaders.String())
	}

	if body.Timeout != nil {
		webhook_.Timeout = timeout
		updateColumns = append(updateColumns, model.WebhookFieldsTimeout.String())
	}

	if body.ConditionValidator != nil {
		webhook_.ConditionValidator = sql.NullInt32{
			Int32: int32(generic.Left(model.ParseWebhookConditionValidator(*body.ConditionValidator))),
			Valid: true}
		updateColumns = append(updateColumns, model.WebhookFieldsConditionValidator.String())
	}

	if body.ConditionFilter != nil {
		webhook_.ConditionFilter = macro.NewNullString(body.ConditionFilter)
		updateColumns = append(updateColumns, model.WebhookFieldsConditionFilter.String())
	}

	if 0 < len(updateColumns) {
		webhook_.Updated = time.Now()
		updateColumns = append(updateColumns, model.WebhookFieldsUpdated.String())
	}

	// something changed?
	if len(updateColumns) == 0 {
		err := errors.Wrapf(echov4.ErrorInvalidRequestParameter, "nothing to change")

		return echov4.HttpError(err, http.StatusBadRequest)
	}

	// check condition validator and condition filter
	err = webhook.CheckConditionFilter(
		model.WebhookConditionValidator(webhook_.ConditionValidator.Int32),
		webhook_.ConditionFilter.String)
	if err != nil {
		err := errors.Wrapf(err, "invalid condition validator")

		return echov4.HttpError(err, http.StatusBadRequest)
	}

	// save
	if err := control.UpsertWebhook(ctx.Request().Context(), webhook_, updateColumns); err != nil {
		err := errors.Wrapf(err, "save the changes to the webhook")

		return echov4.HttpError(err, http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, Webhook(NewWebhook(*webhook_)))
}

// @Description Delete a webhook
// @Accept      json
// @Produce     json
// @Tags        /api/v1/webhook
// @Router      /api/v1/webhook/{uuid} [delete]
// @Param       uuid path string true "Webhook UUID"
// @Success     200
func DeleteWebhook(ctx echo.Context) error {
	const (
		UUID = "uuid"
	)

	if len(echov4.PathParam(ctx)[UUID]) == 0 {
		err := errors.Wrapf(echov4.ErrorInvalidRequestPath, "path=%q", UUID)

		return echov4.HttpError(err, http.StatusBadRequest)
	}

	uuid := echov4.PathParam(ctx)[UUID]

	if err := control.DeleteWebhook(ctx.Request().Context(), uuid); err != nil {
		err := errors.Wrapf(err, "failed to delete a webhook")

		return echov4.HttpError(err, http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, echov4.OK())
}

// @Description Publish a webhook message
// @Accept      json
// @Produce     json
// @Tags        /api/v1/webhook
// @Router      /api/v1/webhook/{uuid}/publish [post]
// @Param       uuid    path string true "webhook UUID"
// @Param       message body []byte true "Publish message"
// @Success     200
func Publish(ctx echo.Context) error {
	const (
		UUID = "uuid"
	)

	var body interface{}
	if err := ctx.Bind(&body); err != nil {
		err := errors.Wrapf(err, "bind request (body=%T)", body)

		return echov4.HttpError(err, http.StatusBadRequest)
	}

	if len(echov4.PathParam(ctx)[UUID]) == 0 {
		err := errors.Wrapf(echov4.ErrorInvalidRequestPath, "path=%q", UUID)

		return echov4.HttpError(err, http.StatusBadRequest)
	}

	uuid := echov4.PathParam(ctx)[UUID]

	record, err := control.GetWebhook(ctx.Request().Context(), uuid)
	if err != nil && err == sql.ErrNoRows {
		err := errors.Wrapf(err, "get webhook (uuid=%v)", uuid)

		return echov4.HttpError(err, http.StatusNotFound)
	}
	if err != nil && err != sql.ErrNoRows {
		err := errors.Wrapf(err, "get webhook")

		return echov4.HttpError(err, http.StatusInternalServerError)
	}

	config := webhook.Config{
		URL:                record.URL,
		Method:             record.Method,
		Headers:            http.Header(record.Headers),
		ConditionValidator: model.WebhookConditionValidator(record.ConditionValidator.Int32),
		ConditionFilter:    record.ConditionFilter.String,
		Timeout:            time.Duration(record.Timeout.Int32) * time.Second,
	}

	if err := config.Publish(ctx.Request().Context(), body); err != nil {
		err := errors.Wrapf(err, "failed to publish webhook")

		return echov4.HttpError(err, http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, echov4.OK())
}

func ParseTimeout(p *string) (sql.NullInt32, error) {
	if p == nil {
		return sql.NullInt32{}, nil
	}

	d, err := time.ParseDuration(*p)
	if err != nil {
		return sql.NullInt32{}, err
	}

	return sql.NullInt32{Int32: int32(d / time.Second), Valid: true}, nil
}
