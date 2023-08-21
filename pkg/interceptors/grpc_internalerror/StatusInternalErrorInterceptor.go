package grpc_internalerror

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

/*
	Every response with an "plain" error will be converted to an status-error with codes.Internal.
*/

// UnaryServerInterceptor returns a new unary server interceptor for panic recovery.
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {

	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ any, err error) {

		result, err := handler(ctx, req)
		_, ok := status.FromError(err)
		if !ok {
			err = status.Error(codes.Internal, err.Error())
		}

		return result, err
	}
}

// StreamServerInterceptor returns a new streaming server interceptor for panic recovery.
func StreamServerInterceptor() grpc.StreamServerInterceptor {

	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {

		err = handler(srv, stream)
		_, ok := status.FromError(err)
		if !ok {
			err = status.Error(codes.Internal, err.Error())
		}

		return err
	}
}
