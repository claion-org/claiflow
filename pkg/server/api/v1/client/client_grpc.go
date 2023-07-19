package client

import (
	context "context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/claion-org/claiflow/pkg/cryptography"
	"github.com/claion-org/claiflow/pkg/logger"
	"github.com/claion-org/claiflow/pkg/server/api/client"
	"github.com/claion-org/claiflow/pkg/server/api/datatype"
	"github.com/claion-org/claiflow/pkg/server/control"
	"github.com/claion-org/claiflow/pkg/server/macro"
	"github.com/claion-org/claiflow/pkg/server/macro/generic"
	"github.com/claion-org/claiflow/pkg/server/model"
	"github.com/claion-org/claiflow/pkg/server/status/globvar"
	"github.com/claion-org/claiflow/pkg/webhook"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	grpcstatus "google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var ErrorInvalidIncomingContext = fmt.Errorf("invalid incoming context")
var ErrorInvalidToken = fmt.Errorf("invalid token")

type GRPC struct {
	client.UnimplementedClientServiceServer
}

// AuthV1 implements client.ClientServiceServer
func (*GRPC) AuthV1(ctx context.Context, auth *client.AuthRequestV1) (*client.AuthResponseV1, error) {
	token, err := control.GetClusterClientTokenByAssertion(ctx, auth.ClusterUuid, auth.Assertion)
	if err != nil && err == sql.ErrNoRows {
		err := errors.Wrapf(err, "get cluster client token (cluster=%v, assertion=%v)",
			auth.ClusterUuid,
			auth.Assertion)

		return nil, grpcstatus.Errorf(codes.InvalidArgument, err.Error())
	}
	if err != nil && err != sql.ErrNoRows {
		err := errors.Wrapf(err, "get cluster client token")

		return nil, grpcstatus.Errorf(codes.Internal, err.Error())
	}

	// valid expiration time
	if time.Until(token.ExpirationTime) < 0 {
		err := errors.Wrapf(ErrorInvalidToken, "token was expired")

		return nil, grpcstatus.Errorf(codes.InvalidArgument, err.Error())
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

		return nil, grpcstatus.Errorf(codes.Internal, err.Error())
	}

	return &client.AuthResponseV1{Token: claim_}, nil
}

// ServicePollingV1 implements client.ClientServiceServer
func (*GRPC) ServicePollingV1(ctx context.Context, in *client.ServicePollingRequestV1) (*client.ServicePollingResponseV1, error) {
	inMD, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		err := errors.Wrapf(ErrorInvalidIncomingContext, "invalid incoming was missing")

		return nil, grpcstatus.Errorf(codes.Unauthenticated, err.Error())
	}

	claimsToken, claims, err := GetClientSessionClaims(ctx, Header(inMD))
	if err != nil {
		err := errors.Wrapf(err, "get client session claim")

		return nil, grpcstatus.Errorf(codes.Unauthenticated, err.Error())
	}

	limit := int(in.Limit)
	if limit <= 0 {
		limit = 1
	}

	// keep alive session
	go KeepAliveClientSessionStatus(ctx, claimsToken, *claims,
		globvar.ClientSessionExpirationTime.Duration,
		func(err error) {
			logger := logger.WithName("KeepAliveClientSessionStatus").WithValues(
				"cluster", claims.ClusterUUID,
				"assertion", claims.ClusterClientTokenUUID,
				echo.HeaderXRequestID, Header(inMD).Get(echo.HeaderXRequestID))

			logger.Error(err, "keep client session alive")
		})

	offset, err := control.GetClusterPollingOffset(ctx, claims.ClusterUUID)
	if err != nil && err != sql.ErrNoRows {
		err := errors.Wrapf(err, "get cluster service polling offset")

		return nil, grpcstatus.Errorf(codes.Internal, err.Error())
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

		services, err := control.PollClusterService(ctx, claims.ClusterUUID, offset_, filter)
		if err != nil {
			err := errors.Wrapf(err, "polling service")

			return nil, grpcstatus.Errorf(codes.Internal, err.Error())
		}

		if len(services) == 0 {
			delay = macro.IntervalSqrt(delay, time.Second, (time.Second * 3), i)
			select {
			case <-ctx.Done():
				err := errors.Wrapf(ctx.Err(), "wait polling service")

				return nil, grpcstatus.Errorf(codes.DeadlineExceeded, err.Error())
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
		if err := control.CreateClusterServiceStatuses(ctx, statuses); err != nil {
			err := errors.Wrapf(err, "save the service state")

			return nil, grpcstatus.Errorf(codes.Internal, err.Error())
		}

		// reset offset
		var resetOffset time.Time
		for _, it := range services {
			if resetOffset.Before(it.Created) {
				resetOffset = it.Created
			}
		}

		if err := control.UpsertClusterPollingOffset(ctx, claims.ClusterUUID, resetOffset, time.Now()); err != nil {
			err := errors.Wrapf(err, "save the service polling offset")

			return nil, grpcstatus.Errorf(codes.Internal, err.Error())
		}

		ToProto := func(a model.ClusterService) (*client.ServicePollingResponseV1_Data, error) {
			summary := datatype.NullString{String_: a.Summary.String, Valid: a.Summary.Valid}
			subscribedChannel := datatype.NullString{String_: a.SubscribedChannel.String, Valid: a.SubscribedChannel.Valid}
			bytes, err := json.Marshal(a.Inputs)
			if err != nil {
				return nil, err
			}

			b := client.ServicePollingResponseV1_Data{}
			b.PartitionDate = timestamppb.New(a.PartitionDate)
			b.ClusterUuid = a.ClusterUUID
			b.Uuid = a.UUID
			b.Name = a.Name
			b.Summary = &summary
			b.TemplateUuid = a.TemplateUUID
			b.Flow = a.Flow
			b.Inputs = bytes
			b.StepMax = int32(a.StepMax)
			b.SubscribedChannel = &subscribedChannel
			b.Priority = int32(a.Priority)
			b.Created = timestamppb.New(a.Created)

			return &b, nil
		}

		datas, err := generic.MapE(services, ToProto)
		if err != nil {
			err := errors.Wrapf(err, "convert the service polling response")

			return nil, grpcstatus.Errorf(codes.Internal, err.Error())
		}

		return &client.ServicePollingResponseV1{Datas: datas}, nil
	}
}

// UpdateServiceStatusV1 implements client.ClientServiceServer
func (*GRPC) UpdateServiceStatusV1(ctx context.Context, serviceStatus *client.UpdateServiceStatusRequestV1) (*client.UpdateServiceStatusResponseV1, error) {
	inMD, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		err := errors.Wrapf(ErrorInvalidIncomingContext, "invalid incoming was missing")

		return nil, grpcstatus.Errorf(codes.Unauthenticated, err.Error())
	}

	_, claims, err := GetClientSessionClaims(ctx, Header(inMD))
	if err != nil {
		err := errors.Wrapf(err, "get client session claim")

		return nil, grpcstatus.Errorf(codes.Unauthenticated, err.Error())
	}

	service, err := control.GetClusterService(ctx, claims.ClusterUUID, serviceStatus.Uuid)
	if err != nil && err == sql.ErrNoRows {
		err := errors.Wrapf(err, "get service (cluster=%v)", claims.ClusterUUID)

		return nil, grpcstatus.Errorf(codes.ResourceExhausted, err.Error())
	}
	if err != nil && err != sql.ErrNoRows {
		err := errors.Wrapf(err, "get service")

		return nil, grpcstatus.Errorf(codes.Internal, err.Error())
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
		StepSeq:       int(serviceStatus.Sequence),
		Status:        Status(model.StepStatus(serviceStatus.Status), int(serviceStatus.Sequence), int(service.StepMax)),
		Started:       NewNullTime(serviceStatus.Started),
		Ended:         NewNullTime(serviceStatus.Ended),
		Message:       NewNullString(serviceStatus.Error),
	}

	if err := control.CreateClusterServiceStatus(ctx, &status); err != nil {
		err := errors.Wrapf(err, "save the service status")

		return nil, grpcstatus.Errorf(codes.Internal, err.Error())
	}

	var result = model.ClusterServiceResult{
		PartitionDate:  service.PartitionDate,
		ClusterUUID:    service.ClusterUUID,
		UUID:           service.UUID,
		ResultSaveType: model.ResultSaveTypeDatabase,
		Result:         cryptography.CipherString(serviceStatus.Result),
		Created:        nowTime,
	}

	// TODO: need service optional field
	isSaveResult := service.StepMax == int(serviceStatus.Sequence)+1 &&
		model.StepStatus(serviceStatus.Status) == model.StepStatusSucceeded

	if isSaveResult {
		if err := control.UpsertClusterServiceResult(ctx, &result); err != nil {
			err := errors.Wrapf(err, "save the service result")

			return nil, grpcstatus.Errorf(codes.Internal, err.Error())
		}
	}

	go func(ctx context.Context) {
		// webhook
		if service.SubscribedChannel.Valid {
			logger := logger.
				WithName("UpdateServiceStatusV1").WithName("Webhook").
				WithValues(
					"service", service.UUID,
					"webhook", service.SubscribedChannel,
					echo.HeaderXRequestID, Header(inMD).Get(echo.HeaderXRequestID))

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
				Map{"step_status": int(serviceStatus.Status)},
				Map{"step_started": serviceStatus.Started.AsTime()},
				Map{"step_ended": serviceStatus.Ended.AsTime()},
				Map{"result": JsonRawMessage(serviceStatus.Result)},
				Map{"error": serviceStatus.Error},
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
	}(context.Background())

	return &client.UpdateServiceStatusResponseV1{}, nil
}

var _ client.ClientServiceServer = (*GRPC)(nil)

func NewNullTime(t *timestamppb.Timestamp) sql.NullTime {
	t_ := t.AsTime()
	if t_.IsZero() {
		return macro.NewNullTime(nil)
	}

	return macro.NewNullTime(&t_)
}

func NewNullString(s string) sql.NullString {
	return macro.NewNullString(&s)
}

func Header(md metadata.MD) http.Header {
	var h = http.Header{}
	for k, vv := range md {
		for _, v := range vv {
			if h[k] == nil {
				h.Set(k, v)
			} else {
				h.Add(k, v)
			}
		}
	}

	return h
}
