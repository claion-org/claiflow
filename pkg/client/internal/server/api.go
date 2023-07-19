package server

import (
	"compress/gzip"
	"context"
	"fmt"
	"sync/atomic"

	"google.golang.org/grpc"
	ggzip "google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/metadata"

	"github.com/claion-org/claiflow/pkg/client/internal/service"
	"github.com/claion-org/claiflow/pkg/client/log"
	apiclient "github.com/claion-org/claiflow/pkg/server/api/client"
)

func init() {
	ggzip.SetLevel(gzip.BestCompression)
}

type ServerAPIInterface interface {
	GetToken() string
	Auth(ctx context.Context, auth *apiclient.AuthRequestV1) error
	GetServices(ctx context.Context, limit int) ([]*apiclient.ServicePollingResponseV1_Data, error)
	UpdateServices(ctx context.Context, data *service.UpdateService) error
}

var _ ServerAPIInterface = &ServerAPI{}

type ServerAPI struct {
	grpcClient apiclient.ClientServiceClient
	authToken  atomic.Value
}

func NewServerAPI(clientService apiclient.ClientServiceClient) *ServerAPI {
	return &ServerAPI{grpcClient: clientService}
}

func (s *ServerAPI) GetToken() string {
	x := s.authToken.Load()

	return x.(string)
}

func (s *ServerAPI) Auth(ctx context.Context, auth *apiclient.AuthRequestV1) error {
	if auth == nil {
		return fmt.Errorf("auth is nil")
	}

	out, err := s.grpcClient.AuthV1(ctx, auth)
	if err != nil {
		return err
	}

	// get session token
	token := out.GetToken()
	if token == "" {
		return fmt.Errorf("auth's result(token) is empty")
	}
	s.authToken.Store(token)

	return nil
}

func (s *ServerAPI) GetServices(ctx context.Context, limit int) ([]*apiclient.ServicePollingResponseV1_Data, error) {
	token := s.GetToken()
	if token == "" {
		return nil, fmt.Errorf("session token is empty")
	}
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token)

	log.Debugf("request get_service: limit=%d\n", limit)

	out, err := s.grpcClient.ServicePollingV1(ctx, &apiclient.ServicePollingRequestV1{Limit: int32(limit)})
	if err != nil {
		return nil, err
	}

	services := out.GetDatas()

	return services, nil
}

func (s *ServerAPI) UpdateServices(ctx context.Context, data *service.UpdateService) error {
	if data == nil {
		return fmt.Errorf("service is nil")
	}

	resultLen := 0
	if data.GetResult().Err == nil {
		resultLen = len(data.GetResult().Body)
	} else {
		resultLen = len(data.GetResult().Err.Error())
	}
	log.Debugf("request update_service: service.uuid=%s, service.step_count=%d, service.sequence=%d, service.status=%s, service.result_len=%d\n", data.GetId(), data.GetStepCount(), data.GetSequence(), data.GetStatus().String(), resultLen)

	// convert to protobuf message
	in := service.ConvertServiceStepUpdateClientToServer(data)

	token := s.GetToken()
	if token == "" {
		return fmt.Errorf("session token is empty")
	}
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token)

	if _, err := s.grpcClient.UpdateServiceStatusV1(ctx, in,
		grpc.UseCompressor(ggzip.Name),
	); err != nil {
		return err
	}

	return nil
}
