package client

import (
	context "context"
	"database/sql"
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/claion-org/claiflow/pkg/cryptography"
	"github.com/claion-org/claiflow/pkg/echov4"
	"github.com/claion-org/claiflow/pkg/logger"
	"github.com/claion-org/claiflow/pkg/server/control"
	"github.com/claion-org/claiflow/pkg/server/macro"
	"github.com/claion-org/claiflow/pkg/server/macro/generic"
	"github.com/claion-org/claiflow/pkg/server/model"
	"github.com/claion-org/claiflow/pkg/server/status/globvar"
	"github.com/claion-org/claiflow/pkg/webhook"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type Auth struct {
	ClusterUuid      string `json:"cluster_uuid"`
	Assertion        string `json:"assertion"`
	ClientVersion    string `json:"client_version"`
	ClientLibVersion string `json:"client_lib_version"`
}

type AccessTokenResponse struct {
	TokenType    string `json:"token_type"   enum:"Basic,Bearer"` // (required)
	AccessToken  string `json:"access_token"`                     // (required)
	ExpiresIn    int64  `json:"expires_in,omitempty"`             // (recommended)
	RefreshToken string `json:"refresh_token,omitempty"`          // (optional)
	Scope        string `json:"scope,omitempty"`                  // (optional)
}

type Object = map[string]any

type ServiceResponse struct {
	ClusterUUID       string         `json:"cluster_uuid"`
	UUID              string         `json:"uuid"`
	Name              string         `json:"name"`
	Summary           *string        `json:"summary,omitempty"`
	TemplateUUID      string         `json:"template_uuid"`
	Flow              string         `json:"flow"`
	Inputs            Object         `json:"inputs" swaggertype:"object"`
	StepMax           int            `json:"step_max"`
	SubscribedChannel *string        `json:"subscribed_channel,omitempty"`
	Priority          model.Priority `json:"priority"`
	Created           time.Time      `json:"created"`
}

func NewServiceResponse(service model.ClusterService) ServiceResponse {
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
	}
}

type ServiceStatus struct {
	Uuid     string           `json:"uuid"`
	Sequence int              `json:"sequence"`
	Status   model.StepStatus `json:"status"`
	Error    *string          `json:"error,omitempty"`
	Started  *time.Time       `json:"started,omitempty"`
	Ended    *time.Time       `json:"ended,omitempty"`
	Result   *string          `json:"result,omitempty"`
}

type Echo struct{}

// @Description auth client
// @Accept      json
// @Produce     json
// @Tags        /api/v1/client
// @Router      /api/v1/client/auth [post]
// @Param       body body     Auth true "Auth"
// @Success     200  {string} ok
// @Header      200  {string} x-sudory-client-token
// @Header      200  {object} AccessTokenResponse
func (Echo) Auth(ctx echo.Context) error {
	var auth Auth
	if err := ctx.Bind(&auth); err != nil {
		err := errors.Wrapf(err, "bind request (body=%T)", auth)

		return echov4.HttpError(err, http.StatusBadRequest)
	}

	token, err := control.GetClusterClientTokenByAssertion(ctx.Request().Context(), auth.ClusterUuid, auth.Assertion)
	if err != nil && err == sql.ErrNoRows {
		err := errors.Wrapf(err, "get cluster client token (cluster=%v, assertion=%v)",
			auth.ClusterUuid,
			auth.Assertion)

		return echov4.HttpError(err, http.StatusBadRequest)
	}
	if err != nil && err != sql.ErrNoRows {
		err := errors.Wrapf(err, "get cluster client token")

		return echov4.HttpError(err, http.StatusInternalServerError)
	}

	// valid expiration time
	if time.Until(token.ExpirationTime) < 0 {
		err := errors.Wrapf(ErrorInvalidToken, "token was expired")

		return echov4.HttpError(err, http.StatusBadRequest)
	}

	claimUUID := macro.NewUuidString()
	timeNow := time.Now()
	iat := timeNow
	// exp := globvar.ClientSession.ExpirationTime(timeNow)
	exp := token.ExpirationTime

	claim := model.NewClusterClientSessionClaim(
		token.ClusterUUID,
		token.UUID,
		claimUUID,
		iat,
		exp,
		auth.ClientVersion,
		auth.ClientLibVersion,
	)

	claim_, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claim).
		SignedString(globvar.ClientSessionSignatureSecret.Bytes)
	if err != nil {
		err := errors.Wrapf(err, "make JWT claim (method=%v, claim=%T)",
			jwt.SigningMethodHS256,
			claim)

		return echov4.HttpError(err, http.StatusInternalServerError)
	}

	ctx.Response().Header().Add(HTTP_HEADER_X_OLD_CLIENT_TOKEN, claim_)

	return ctx.JSON(http.StatusOK, AccessTokenResponse{
		TokenType:   "Bearer",
		AccessToken: claim_,
		ExpiresIn:   claim.ExpiresAt - claim.IssuedAt,
	})
}

// @Description get []Service
// @Security    ClientAuthorization
// @Accept      json
// @Produce     json
// @Tags        /api/v1/client
// @Router      /api/v1/client/service [get]
// @param       limit query   int false "count limit of ServiceResponse""
// @Success     200   {array} ServiceResponse
func (Echo) PollService(ctx echo.Context) error {
	claimsToken, claims, err := GetClientSessionClaims(ctx.Request().Context(), ctx.Request().Header)
	if err != nil {
		err := errors.Wrapf(err, "get client session claim")

		return echov4.HttpError(err, http.StatusForbidden)
	}

	limit, _ := strconv.Atoi(ctx.QueryParam("limit"))
	if limit <= 0 {
		limit = 1
	}

	// keep alive session
	go KeepAliveClientSessionStatus(ctx.Request().Context(), claimsToken, *claims,
		globvar.ClientSessionExpirationTime.Duration,
		func(err error) {
			logger := logger.WithName("KeepAliveClientSessionStatus").WithValues(
				"cluster", claims.ClusterUUID,
				"assertion", claims.ClusterClientTokenUUID,
				echo.HeaderXRequestID, ctx.Request().Header.Get(echo.HeaderXRequestID))

			logger.Error(err, "keep client session alive")
		})

	offset, err := control.GetClusterPollingOffset(ctx.Request().Context(), claims.ClusterUUID)
	if err != nil && err != sql.ErrNoRows {
		err := errors.Wrapf(err, "get cluster service polling offset")

		return echov4.HttpError(err, http.StatusInternalServerError)
	}

	Offset := func(offset time.Time, exp time.Time) time.Time {
		if offset.IsZero() {
			return exp
		}

		return offset
	}

	var delay time.Duration
LOOP_POLL:
	for i := 0; ; i++ {
		now := time.Now()
		exp := globvar.ClientConfigServiceValidityPeriod.Add(now)
		filter := PollServiceFilter(limit, exp, now)
		offset_ := Offset(offset, exp)

		services, err := control.PollClusterService(ctx.Request().Context(), claims.ClusterUUID, offset_, filter)
		if err != nil {
			err := errors.Wrapf(err, "polling service")

			return echov4.HttpError(err, http.StatusInternalServerError)
		}

		if len(services) == 0 {
			delay = macro.IntervalSqrt(delay, time.Second, (time.Second * 3), i)
			select {
			case <-ctx.Request().Context().Done():
				err := errors.Wrapf(ctx.Request().Context().Err(), "wait polling service")

				return echov4.HttpError(err, http.StatusRequestTimeout)
			case <-time.After(delay):
			}

			continue LOOP_POLL
		}

		// sort ASC by created field
		sort.Slice(services, func(i, j int) bool { return services[i].Created.Before(services[j].Created) })

		nowTime := time.Now()

		statuses := generic.Map(services, func(service model.ClusterService) model.ClusterServiceStatus {
			// update service status
			var status = model.ClusterServiceStatus{
				PartitionDate: service.PartitionDate,
				ClusterUUID:   service.ClusterUUID,
				Uuid:          service.UUID,
				Created:       nowTime,
				StepMax:       service.StepMax,
				StepSeq:       0,
				Status:        model.StepStatusSent,
				Started:       macro.NewNullTime(&nowTime),
				Ended:         macro.NewNullTime(&nowTime),
				Message:       macro.NewNullString(nil),
			}

			return status
		})
		if err := control.CreateClusterServiceStatuses(ctx.Request().Context(), statuses); err != nil {
			err := errors.Wrapf(err, "save the service state")

			return echov4.HttpError(err, http.StatusInternalServerError)
		}

		// reset offset
		var resetOffset time.Time
		for _, it := range services {
			if resetOffset.Before(it.Created) {
				resetOffset = it.Created
			}
		}

		if err := control.UpsertClusterPollingOffset(ctx.Request().Context(), claims.ClusterUUID, resetOffset, time.Now()); err != nil {
			err := errors.Wrapf(err, "save the service polling offset")

			return echov4.HttpError(err, http.StatusInternalServerError)
		}

		return ctx.JSON(http.StatusOK, []ServiceResponse(generic.Map(services, NewServiceResponse)))
	}
}

// @Description update a service status
// @Security    ClientAuthorization
// @Accept      json
// @Produce     json
// @Tags        /api/v1/client
// @Router      /api/v1/client/service [put]
// @Param       body body ServiceStatus true "ServiceStatus"
// @Success     200
func (Echo) UpdateService(ctx echo.Context) error {
	var body ServiceStatus
	if err := ctx.Bind(&body); err != nil {
		err := errors.Wrapf(err, "bind request (body=%T)", body)

		return echov4.HttpError(err, http.StatusBadRequest)
	}

	_, claims, err := GetClientSessionClaims(ctx.Request().Context(), ctx.Request().Header)
	if err != nil {
		logger.Error(err, "get client session claim")

		return echov4.HttpError(err, http.StatusForbidden)
	}

	service, err := control.GetClusterService(ctx.Request().Context(), claims.ClusterUUID, body.Uuid)
	if err != nil && err == sql.ErrNoRows {
		err := errors.Wrapf(err, "get service (cluster=%v)", claims.ClusterUUID)

		return echov4.HttpError(err, http.StatusNotFound)
	}
	if err != nil && err != sql.ErrNoRows {
		err := errors.Wrapf(err, "get service")

		return echov4.HttpError(err, http.StatusInternalServerError)
	}

	nowTime := time.Now()

	Status := func(status model.StepStatus, stepSeq, stepMax int) model.StepStatus {
		if stepMax == stepSeq+1 {
			return status
		}

		if status == model.StepStatusFailed {
			return model.StepStatusFailed
		}

		return model.StepStatusProcessing
	}

	var status = model.ClusterServiceStatus{
		PartitionDate: service.PartitionDate,
		ClusterUUID:   service.ClusterUUID,
		Uuid:          service.UUID,
		Created:       nowTime,
		StepMax:       service.StepMax,
		StepSeq:       body.Sequence,
		Status:        Status(body.Status, body.Sequence, service.StepMax),
		Started:       macro.NewNullTime(body.Started),
		Ended:         macro.NewNullTime(body.Ended),
		Message:       macro.NewNullString(body.Error),
	}

	if err := control.CreateClusterServiceStatus(ctx.Request().Context(), &status); err != nil {
		err := errors.Wrapf(err, "save the service status")

		return echov4.HttpError(err, http.StatusInternalServerError)
	}

	var result = model.ClusterServiceResult{
		PartitionDate:  service.PartitionDate,
		ClusterUUID:    service.ClusterUUID,
		UUID:           service.UUID,
		ResultSaveType: model.ResultSaveTypeDatabase,
		Result:         cryptography.CipherString(*body.Result),
		Created:        nowTime,
	}

	// TODO: need service optional field
	isSaveResult := service.StepMax == int(body.Sequence)+1 &&
		model.StepStatus(body.Status) == model.StepStatusSucceeded

	if isSaveResult {
		if err := control.UpsertClusterServiceResult(ctx.Request().Context(), &result); err != nil {
			err := errors.Wrapf(err, "save the service result")

			return echov4.HttpError(err, http.StatusInternalServerError)
		}
	}

	go func(ctx context.Context, h http.Header) {
		// webhook
		if service.SubscribedChannel.Valid {
			logger := logger.
				WithName("UpdateService").WithName("Webhook").
				WithValues(
					"service", service.UUID,
					"webhook", service.SubscribedChannel,
					echo.HeaderXRequestID, h.Get(echo.HeaderXRequestID))

			webhook_, err := control.GetWebhook(ctx, service.SubscribedChannel.String)
			if err != nil {
				logger.Error(err, "get webhook")

				return
			}

			JsonRawMessage := func(s string) json.RawMessage {
				s = strings.TrimSpace(s)

				head, tail := s[0], s[len(s)-1]
				if ok := (head == '{' && tail == '}') || (head == '[' && tail == ']'); ok {
					return json.RawMessage(s)
				}

				return json.RawMessage(strconv.Quote(s))
			}

			var payload = EmbedFields(
				Map{"template_uuid": service.TemplateUUID},
				Map{"cluster_uuid": service.ClusterUUID},
				Map{"service_uuid": service.UUID},
				Map{"service_name": service.Name},
				Map{"inputs": service.Inputs},
				Map{"assigned_client_uuid": claims.UUID},
				Map{"status": int(status.Status)},
				Map{"status_description": status.Status.String()},
				Map{"step_count": service.StepMax},
				Map{"step_position": status.StepSeq},
				Map{"step_status": int(body.Status)},
				Map{"step_started": macro.NewNullTime(body.Started).Time},
				Map{"step_ended": macro.NewNullTime(body.Ended).Time},
				Map{"result": JsonRawMessage(macro.NewNullString(body.Result).String)},
				Map{"error": body.Error},
				Map{"webhook": service.SubscribedChannel.String},
			)

			var webhook = webhook.Config{
				URL:                webhook_.URL,
				Method:             webhook_.Method,
				Headers:            http.Header(webhook_.Headers),
				ConditionValidator: model.WebhookConditionValidator(webhook_.ConditionValidator.Int32),
				ConditionFilter:    webhook_.ConditionFilter.String,
				Timeout:            time.Duration(webhook_.Timeout.Int32) * time.Second,
			}

			if err := webhook.Publish(ctx, payload); err != nil {
				logger.Error(err, "publish webhook")

				return
			}
		}
	}(context.Background(), ctx.Request().Header)

	return ctx.JSON(http.StatusOK, echov4.OK())
}
