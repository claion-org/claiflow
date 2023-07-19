package server

import (
	"fmt"
	"net"

	x_grpc_middleware "github.com/claion-org/claiflow/pkg/grpc_middleware"
	clientpb "github.com/claion-org/claiflow/pkg/server/api/client"
	"github.com/claion-org/claiflow/pkg/server/api/v1/client"
	"github.com/claion-org/claiflow/pkg/server/config"
	"github.com/go-logr/logr"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	_ "google.golang.org/grpc/encoding/gzip" // Install the gzip compressor
	"google.golang.org/grpc/reflection"
)

func NewGrpcServer(_config config.Config, logger logr.Logger) (*grpc.Server, error) {
	opts := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(_config.GrpcService.MaxRecvMsgSize), // 1 GB
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_prometheus.UnaryServerInterceptor,
			// grpc_middleware_logging.UnaryInterceptor_XRequestID(),
			x_grpc_middleware.UnaryInterceptor_Logging(logger),
		)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_prometheus.StreamServerInterceptor,
			// grpc_middleware_logging.StreamInterceptor_XRequestID(),
			x_grpc_middleware.StreamInterceptor_Logging(logger),
		)),
	}

	if _config.GrpcService.Tls.Enable {
		cred, err := credentials.NewServerTLSFromFile(_config.GrpcService.Tls.CertFile, _config.GrpcService.Tls.KeyFile)
		if err != nil {
			err := fmt.Errorf("%w: new gRPC server credentials from files", err)
			return nil, err
		}

		opts = append(opts,
			grpc.Creds(cred),
		)
	}

	srv := grpc.NewServer(opts...)

	clientpb.RegisterClientServiceServer(srv, &client.GRPC{})

	reflection.Register(srv)
	grpc_prometheus.Register(srv)

	return srv, nil
}

func NewGrpcListener(_config config.Config) (net.Listener, error) {
	return net.Listen("tcp", fmt.Sprintf(":%d", _config.GrpcService.Port))
}

func HttpServeFactory(_config config.Config) func(e *echo.Echo) error {
	if _config.HttpService.Tls.Enable {
		return func(e *echo.Echo) error {
			return e.StartTLS(fmt.Sprintf(":%d", _config.HttpService.Port),
				_config.HttpService.Tls.CertFile,
				_config.HttpService.Tls.KeyFile)
		}
	} else {
		return func(e *echo.Echo) error {
			return e.Start(fmt.Sprintf(":%d", _config.HttpService.Port))
		}
	}
}
