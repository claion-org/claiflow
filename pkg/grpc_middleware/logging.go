package grpc_middleware

import (
	"context"
	"fmt"

	"github.com/claion-org/claiflow/pkg/logger/verbose"
	"github.com/go-logr/logr"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func UnaryInterceptor_Logging(logger logr.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		ctx, reqID := XRequestID(DefaultCounter, ctx)

		logger := logger.WithName("gRPC").WithName("unary").
			WithValues(
				"grpc-method", info.FullMethod,
				"in-proto", fmt.Sprintf("%T", req),
				echo.HeaderXRequestID, reqID)

		logProtoMessageAsJson(req, func(v interface{}) {
			logger.V(verbose.TRACE).Info("rpc begin",
				"in-payload", v,
			)
		})

		resp, err = handler(ctx, req)
		logger = logger.WithValues(
			"out-proto", fmt.Sprintf("%T", resp),
		)

		if err != nil {
			logger.Error(err, "rpc ERROR")
		} else {
			logger.Info("rpc OK")
		}

		return resp, err
	}
}

func StreamInterceptor_Logging(logger logr.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ss, reqID := XRequestID_Stream(DefaultCounter, ss)

		logger := logger.WithName("gRPC").WithName("stream").
			WithValues(
				"grpc-method", info.FullMethod,
				echo.HeaderXRequestID, reqID)

		newServerStream := &wrappedStream{
			ServerStream: ss,
			WrappedRecvMsg: func(m interface{}) error {
				logger := logger.WithValues(
					"in-proto", fmt.Sprintf("%T", m))

				logger.V(verbose.DEBUG).Info("rpc-stream in")

				logProtoMessageAsJson(m, func(v interface{}) {
					logger.V(verbose.TRACE).Info("rpc-stream in",
						"in-payload", v,
					)
				})

				return ss.RecvMsg(m)
			},
			WrappedSendMsg: func(m interface{}) error {
				logger := logger.WithValues(
					"out-proto", fmt.Sprintf("%T", m))

				logger.V(verbose.DEBUG).Info("rpc-stream out")

				logProtoMessageAsJson(m, func(v interface{}) {
					logger.V(verbose.TRACE).Info("rpc-stream out",
						"out-payload", v,
					)
				})

				return ss.SendMsg(m)
			},
		}

		logger.V(verbose.TRACE).Info("rpc-stream begin")

		err := handler(srv, newServerStream)
		if err != nil {
			logger.Error(err, "rpc-stream ERROR")
		} else {
			logger.Info("rpc-stream OK")
		}

		return err
	}
}

func logProtoMessageAsJson(pbMsg interface{}, fn func(v interface{})) {
	if p, ok := pbMsg.(proto.Message); ok {
		fn(&jsonpbObjectMarshaler{pb: p})
	}
}

type jsonpbObjectMarshaler struct {
	pb proto.Message
}

func (j *jsonpbObjectMarshaler) MarshalLogObject(e zapcore.ObjectEncoder) error {
	// ZAP jsonEncoder deals with AddReflect by using json.MarshalObject. The same thing applies for consoleEncoder.
	return e.AddReflected("msg", j)
}

func (j *jsonpbObjectMarshaler) MarshalJSON() ([]byte, error) {
	b, err := JsonPbMarshaller.Marshal(j.pb)
	if err != nil {
		return nil, fmt.Errorf("jsonpb serializer failed: %v", err)
	}
	return b, nil
}

var JsonPbMarshaller = &protojson.MarshalOptions{}
