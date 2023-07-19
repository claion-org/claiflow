package service

import (
	"context"
	"database/sql"
	"encoding/json"
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
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/qri-io/jsonschema"
)

type Object = map[string]any

type CreateMultiClusters struct {
	ClusterUUIDs      []string `json:"cluster_uuids"`
	UUID              *string  `json:"uuid,omitempty"` // (optional)
	Name              string   `json:"name"`
	Summary           *string  `json:"summary,omitempty"`
	TemplateUUID      string   `json:"template_uuid"`
	Inputs            Object   `json:"inputs"`
	SubscribedChannel *string  `json:"subscribed_channel,omitempty"`
}

type Create struct {
	UUID              *string `json:"uuid,omitempty"` // (optional)
	Name              string  `json:"name"`
	Summary           *string `json:"summary,omitempty"`
	TemplateUUID      string  `json:"template_uuid"`
	Inputs            Object  `json:"inputs"`
	SubscribedChannel *string `json:"subscribed_channel,omitempty"`
}

type Update struct {
	Uuid     string           `json:"uuid"`
	Sequence int              `json:"sequence"`
	Status   model.StepStatus `json:"status"`
	Result   string           `json:"result"`
	Started  time.Time        `json:"started"`
	Ended    time.Time        `json:"ended"`
}

type StatusResponse struct {
	StepSeq int              `json:"step_seq"`
	Status  model.StepStatus `json:"status"`
	Started *time.Time       `json:"started,omitempty"`
	Ended   *time.Time       `json:"ended,omitempty"`
	Message *string          `json:"message,omitempty"`
	Created time.Time        `json:"created"`
}

type ResultResponse struct {
	ResultSaveType model.ResultSaveType `json:"save_type"`
	Result         string               `json:"result"`
	Created        time.Time            `json:"created"`
}

type ServiceResponse struct {
	ClusterUUID       string           `json:"cluster_uuid"`
	UUID              string           `json:"uuid"`
	Name              string           `json:"name"`
	Summary           *string          `json:"summary,omitempty"`
	TemplateUUID      string           `json:"template_uuid"`
	Flow              string           `json:"flow"`
	Inputs            Object           `json:"inputs" swaggertype:"object"`
	StepMax           int              `json:"step_max"`
	SubscribedChannel *string          `json:"subscribed_channel,omitempty"`
	Priority          model.Priority   `json:"priority"`
	Created           time.Time        `json:"created"`
	Statuses          []StatusResponse `json:"statuses,omitempty"`
	Result            []ResultResponse `json:"results,omitempty"`
}

func NewServiceResponse(service model.ClusterService, statuses []model.ClusterServiceStatus, results []model.ClusterServiceResult) ServiceResponse {
	return ServiceResponse{
		ClusterUUID:       service.ClusterUUID,
		UUID:              service.UUID,
		Name:              service.Name,
		Summary:           macro.FromNullString(service.Summary),
		TemplateUUID:      service.TemplateUUID,
		Flow:              service.Flow,
		Inputs:            service.Inputs,
		StepMax:           service.StepMax,
		SubscribedChannel: macro.FromNullString(service.SubscribedChannel),
		Priority:          service.Priority,
		Created:           service.Created,
		Statuses:          generic.Map(statuses, NewStatusResponse),
		Result:            generic.Map(results, NewResultResponse),
	}
}

func NewStatusResponse(status model.ClusterServiceStatus) StatusResponse {
	return StatusResponse{
		StepSeq: status.StepSeq,
		Status:  status.Status,
		Started: macro.FromNullTime(status.Started),
		Ended:   macro.FromNullTime(status.Ended),
		Message: macro.FromNullString(status.Message),
		Created: status.Created,
	}
}

func NewResultResponse(result model.ClusterServiceResult) ResultResponse {
	return ResultResponse{
		ResultSaveType: result.ResultSaveType,
		Result:         string(result.Result),
		Created:        result.Created,
	}
}

// @Description Create a Service (Multi Clusters)
// @Accept      json
// @Produce     json
// @Tags        /api/v1/service
// @Router      /api/v1/service [post]
// @Param       service body    Create true "Create"
// @Success     200     {array} ServiceResponse
func CreateServiceMultiClusters(ctx echo.Context) error {
	var body CreateMultiClusters
	if err := ctx.Bind(&body); err != nil {
		err := errors.Wrapf(err, "bind request (body=%T)", body)

		return echov4.HttpError(err, http.StatusBadRequest)
	}

	resp, err := createService(ctx, body)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, resp)
}

// @Description Create a Service
// @Accept      json
// @Produce     json
// @Tags        /api/v1/service
// @Router      /api/v1/cluster/{cluster_uuid}/service [post]
// @Param       cluster_uuid path    string true "cluster UUID"
// @Param       service      body    Create true "Create"
// @Success     200          {object} ServiceResponse
func CreateService(ctx echo.Context) error {
	const (
		CLUSTER_UUID = "cluster_uuid"
	)

	var body Create
	if err := ctx.Bind(&body); err != nil {
		err := errors.Wrapf(err, "bind request (body=%T)", body)

		return echov4.HttpError(err, http.StatusBadRequest)
	}

	clusterUUID := echov4.PathParam(ctx)[CLUSTER_UUID]
	multibody := CreateMultiClusters{
		ClusterUUIDs:      []string{clusterUUID},
		UUID:              body.UUID,
		Name:              body.Name,
		Summary:           body.Summary,
		TemplateUUID:      body.TemplateUUID,
		Inputs:            body.Inputs,
		SubscribedChannel: body.SubscribedChannel,
	}

	resp, err := createService(ctx, multibody)
	if err != nil {
		return err
	}

	if len(resp) == 0 {
		return echov4.HttpError(fmt.Errorf("response was empty"), http.StatusInternalServerError)
	}

	if 1 < len(resp) {
		return echov4.HttpError(fmt.Errorf("too many responses"), http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, resp[0])
}

// @Description Find Services
// @Accept      x-www-form-urlencoded
// @Produce     json
// @Tags        /api/v1/service
// @Router      /api/v1/service [get]
// @Param       q   query   string false "query  github.com/claion-org/claiflow/pkg/server/database/vanilla/stmt/README.md"
// @Param       o   query   string false "order  github.com/claion-org/claiflow/pkg/server/database/vanilla/stmt/README.md"
// @Param       p   query   string false "paging github.com/claion-org/claiflow/pkg/server/database/vanilla/stmt/README.md"
// @Success     200 {array} ServiceResponse
func FindService(ctx echo.Context) (err error) {
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

	servicesWithStatuses, err := control.FindClusterService(ctx.Request().Context(), q, o, p)
	if err != nil {
		err := errors.Wrapf(err, "query service")

		return echov4.HttpError(err, http.StatusInternalServerError)
	}

	NewServiceResponse_ := func(a model.ServiceWithStatuses) ServiceResponse {
		return NewServiceResponse(a.ClusterService, a.ClusterServiceStatuses, nil)
	}

	return ctx.JSON(http.StatusOK, []ServiceResponse(generic.Map(servicesWithStatuses, NewServiceResponse_)))
}

// @Description Get a Service
// @Accept      json
// @Produce     json
// @Tags        /api/v1/service
// @Router      /api/v1/cluster/{cluster_uuid}/service/{uuid} [get]
// @Param       cluster_uuid path     string true "cluster UUID"
// @Param       uuid         path     string true "service UUID"
// @Success     200          {object} ServiceResponse
func GetService(ctx echo.Context) (err error) {
	const (
		CLUSTER_UUID = "cluster_uuid"
		UUID         = "uuid"
	)

	if len(echov4.PathParam(ctx)[CLUSTER_UUID]) == 0 {
		err := errors.Wrapf(echov4.ErrorInvalidRequestPath, "path=%q", CLUSTER_UUID)

		return echov4.HttpError(err, http.StatusBadRequest)
	}

	if len(echov4.PathParam(ctx)[UUID]) == 0 {
		err := errors.Wrapf(echov4.ErrorInvalidRequestPath, "path=%q", UUID)

		return echov4.HttpError(err, http.StatusBadRequest)
	}

	clusterUUID := echov4.PathParam(ctx)[CLUSTER_UUID]
	uuid := echov4.PathParam(ctx)[UUID]

	service, err := control.GetClusterService(ctx.Request().Context(), clusterUUID, uuid)
	if err != nil && err == sql.ErrNoRows {
		err := errors.Wrapf(err, "get service (cluster=%v, uuid=%v)", clusterUUID, uuid)

		return echov4.HttpError(err, http.StatusNotFound)
	}
	if err != nil && err != sql.ErrNoRows {
		err := errors.Wrapf(err, "get service")

		return echov4.HttpError(err, http.StatusInternalServerError)
	}

	statuses, err := control.GetClusterServiceStatuses(ctx.Request().Context(), clusterUUID, uuid)
	if err != nil {
		err := errors.Wrapf(err, "query service statuses")

		return echov4.HttpError(err, http.StatusInternalServerError)
	}

	results, err := control.GetClusterServiceResults(ctx.Request().Context(), clusterUUID, uuid)
	if err != nil {
		err := errors.Wrapf(err, "query service results")

		return echov4.HttpError(err, http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, ServiceResponse(NewServiceResponse(*service, statuses, results)))
}

// @Description Get a Service Result
// @Accept      json
// @Produce     json
// @Tags        /api/v1/service
// @Router      /api/v1/cluster/{cluster_uuid}/service/{uuid}/result [get]
// @Param       cluster_uuid path     string true "cluster UUID"
// @Param       uuid         path     string true "service UUID"
// @Success     200          {object} ResultResponse
func GetServiceResult(ctx echo.Context) (err error) {
	const (
		CLUSTER_UUID = "cluster_uuid"
		UUID         = "uuid"
	)

	if len(echov4.PathParam(ctx)[CLUSTER_UUID]) == 0 {
		err := errors.Wrapf(echov4.ErrorInvalidRequestPath, "path=%q", CLUSTER_UUID)

		return echov4.HttpError(err, http.StatusBadRequest)
	}

	if len(echov4.PathParam(ctx)[UUID]) == 0 {
		err := errors.Wrapf(echov4.ErrorInvalidRequestPath, "path=%q", UUID)

		return echov4.HttpError(err, http.StatusBadRequest)
	}

	clusterUUID := echov4.PathParam(ctx)[CLUSTER_UUID]
	uuid := echov4.PathParam(ctx)[UUID]

	result, err := control.GetClusterServiceResult(ctx.Request().Context(), clusterUUID, uuid)
	if err != nil && err == sql.ErrNoRows {
		err := errors.Wrapf(err, "get service (cluster=%v, uuid=%v)", clusterUUID, uuid)

		return echov4.HttpError(err, http.StatusNotFound)
	}
	if err != nil && err != sql.ErrNoRows {
		err := errors.Wrapf(err, "get service")

		return echov4.HttpError(err, http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, ResultResponse(NewResultResponse(*result)))
}

func ValidCreateService(ctx context.Context, body CreateMultiClusters, tmpl model.Template) error {
	validator := jsonschema.Schema{}

	if err := json.Unmarshal([]byte(tmpl.Inputs), &validator); err != nil {
		err = errors.Wrapf(err, "command.Args convert to json schema validator")
		return err
	}

	timeout, cancel := context.WithTimeout(ctx, 333*time.Millisecond)
	defer cancel()

	validation := validator.Validate(timeout, body.Inputs)

	validErr := *validation.Errs

	iter_verr := func() error {
		var err error
		for _, iter := range validErr {
			if err == nil {
				err = iter
				continue
			}
			err = errors.Wrap(err, iter.Error())
		}
		return err
	}
	if err := iter_verr(); err != nil {
		return err
	}

	return nil
}

func NewServiceWithStatus(body CreateMultiClusters, tmpl model.Template, nowTime time.Time) []model.ServiceWithStatuses {
	UUID := func() string {
		if body.UUID == nil {
			return macro.NewUuidString()
		}

		return *body.UUID
	}()

	StringWithDefault := func(a string, b string) string {
		if len(a) == 0 {
			return b
		}

		return a
	}

	// compute flow len
	var flow = []interface{}{}
	_ = json.Unmarshal([]byte(tmpl.Flow), &flow)
	StepMax := func() int {
		return len(flow)
	}

	Priority := func(tmpl model.Template) model.Priority {
		if tmpl.Origin == model.OriginSystem.String() {
			return model.PriorityHigh // system
		}
		return model.PriorityLow
	}

	Service := func(body CreateMultiClusters, clusterUUID string, UUID string) model.ClusterService {
		var newService model.ClusterService

		newService.PartitionDate = nowTime
		newService.Created = nowTime
		newService.ClusterUUID = clusterUUID
		newService.UUID = UUID
		newService.Name = StringWithDefault(body.Name, tmpl.Name)
		newService.Summary = macro.WithDefaultNullString(macro.NewNullString(body.Summary), tmpl.Summary)
		newService.TemplateUUID = body.TemplateUUID
		newService.Flow = tmpl.Flow
		newService.Inputs = body.Inputs
		newService.StepMax = StepMax()
		newService.SubscribedChannel = macro.NewNullString(body.SubscribedChannel)
		newService.Priority = Priority(tmpl)

		return newService
	}

	Status := func(clusterUUID string, UUID string) model.ClusterServiceStatus {
		var newStatus model.ClusterServiceStatus

		newStatus.PartitionDate = nowTime
		newStatus.Created = nowTime
		newStatus.ClusterUUID = clusterUUID
		newStatus.Uuid = UUID
		newStatus.StepMax = StepMax()

		return newStatus
	}

	var out = make([]model.ServiceWithStatuses, 0, len(body.ClusterUUIDs))
	for i := range body.ClusterUUIDs {
		clusterUUID := body.ClusterUUIDs[i]

		out = append(out, model.ServiceWithStatuses{
			ClusterService:         Service(body, clusterUUID, UUID),
			ClusterServiceStatuses: []model.ClusterServiceStatus{Status(clusterUUID, UUID)},
		})
	}

	return out
}

func createService(ctx echo.Context, body CreateMultiClusters) ([]ServiceResponse, error) {
	if len(body.TemplateUUID) == 0 {
		err := errors.Wrapf(echov4.ErrorInvalidRequestParameter, "param=%q", ".template_uuid")

		return nil, echov4.HttpError(err, http.StatusBadRequest)
	}

	if len(body.ClusterUUIDs) == 0 {
		err := errors.Wrapf(echov4.ErrorInvalidRequestParameter, "param=%q", ".cluster_uuids")

		return nil, echov4.HttpError(err, http.StatusBadRequest)
	}

	for i := range body.ClusterUUIDs {
		if len(body.ClusterUUIDs[i]) == 0 {
			err := errors.Wrapf(echov4.ErrorInvalidRequestParameter, "param=%q", ".cluster_uuids")

			return nil, echov4.HttpError(err, http.StatusBadRequest)
		}
	}

	// check cluster
	for i := range body.ClusterUUIDs {
		exists, err := control.IsExistsCluster(ctx.Request().Context(), body.ClusterUUIDs[i])
		if err != nil {
			return nil, echov4.HttpError(err, http.StatusInternalServerError)
		}

		if !exists {
			err := fmt.Errorf("%w: param=%q", echov4.ErrorInvalidRequestParameter, "cluster_uuid")
			return nil, echov4.HttpError(err, http.StatusBadRequest)
		}
	}

	// get template
	template, err := control.GetTemplate(ctx.Request().Context(), body.TemplateUUID)
	if err != nil && err == sql.ErrNoRows {
		err := errors.Wrapf(err, "get template (uuid=%v)", body.TemplateUUID)

		return nil, echov4.HttpError(err, http.StatusNotFound)
	}
	if err != nil && err != sql.ErrNoRows {
		err := errors.Wrapf(err, "get template")

		return nil, echov4.HttpError(err, http.StatusInternalServerError)
	}

	// new service && validation
	if err := ValidCreateService(ctx.Request().Context(), body, *template); err != nil {
		err := errors.Wrapf(err, "valid new service")

		return nil, echov4.HttpError(err, http.StatusInternalServerError)
	}

	nowTime := time.Now()
	newServiceWithStatus := NewServiceWithStatus(body, *template, nowTime)
	if err := control.CreateService(ctx.Request().Context(), newServiceWithStatus); err != nil {
		err := errors.Wrapf(err, "save new service")

		return nil, echov4.HttpError(err, http.StatusInternalServerError)
	}

	NewServiceResponse_ := func(a model.ServiceWithStatuses) ServiceResponse {
		return NewServiceResponse(a.ClusterService, a.ClusterServiceStatuses, nil)
	}

	return generic.Map(newServiceWithStatus, NewServiceResponse_), nil
}
