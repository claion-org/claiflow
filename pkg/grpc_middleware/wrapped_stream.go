package grpc_middleware

import (
	"context"

	"google.golang.org/grpc"
)

// wrappedStream wraps around the embedded grpc.ServerStream, and intercepts the RecvMsg and
// SendMsg method call.
type wrappedStream struct {
	grpc.ServerStream
	WrappedContext context.Context
	WrappedRecvMsg func(m interface{}) error
	WrappedSendMsg func(m interface{}) error
}

func (w *wrappedStream) Context() context.Context {
	if w.WrappedContext != nil {
		return w.WrappedContext
	}

	return w.ServerStream.Context()
}

func (w *wrappedStream) RecvMsg(m interface{}) error {
	if w.WrappedRecvMsg != nil {
		return w.WrappedRecvMsg(m)
	}

	return w.ServerStream.RecvMsg(m)
}

func (w *wrappedStream) SendMsg(m interface{}) error {
	if w.WrappedSendMsg != nil {
		return w.WrappedSendMsg(m)
	}

	return w.ServerStream.SendMsg(m)
}
