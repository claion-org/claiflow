package grpc_middleware

import (
	"context"
	"net/http"
	"strconv"
	"sync/atomic"

	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var DefaultCounter = &atomic.Int64{}

func XRequestID(at *atomic.Int64, ctx context.Context) (context.Context, string) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.MD{}
	}

	reqID := ((http.Header)(md)).Get(echo.HeaderXRequestID)

	// reqID := First(md.Get(echo.HeaderXRequestID))
	// if 0 < len(reqID) {
	// 	return ctx, reqID
	// }

	if 0 < len(reqID) {
		return ctx, reqID
	}

	id := at.Add(1)
	reqID = strconv.FormatInt(id, 10)
	md.Set(echo.HeaderXRequestID, reqID)

	return metadata.NewIncomingContext(ctx, md), reqID
}

func XRequestID_Stream(at *atomic.Int64, ss grpc.ServerStream) (grpc.ServerStream, string) {
	newCtx, reqID := XRequestID(at, ss.Context())

	ss = &wrappedStream{
		ServerStream:   ss,
		WrappedContext: newCtx,
	}

	return ss, reqID
}

func UnaryInterceptor_XRequestID() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ctx, _ = XRequestID(DefaultCounter, ctx)
		return handler(ctx, req)
	}
}

func StreamInterceptor_XRequestID() grpc.StreamServerInterceptor {
	return func(srv interface{}, newServerStream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		newServerStream, _ = XRequestID_Stream(DefaultCounter, newServerStream)

		newServerStream = &wrappedStream{
			ServerStream:   newServerStream,
			WrappedContext: newServerStream.Context(),
		}

		return handler(srv, newServerStream)
	}
}
