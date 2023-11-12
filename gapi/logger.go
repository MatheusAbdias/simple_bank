package gapi

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GrpcLogger(
	ctx context.Context,
	request any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (response any, err error) {
	startTime := time.Now()
	slog.Info("gRPC method called")
	slog.Info(fmt.Sprintf("Method: %s", info.FullMethod))

	result, err := handler(ctx, request)

	duration := time.Since(startTime)
	slog.Duration("Duration", duration)

	statusCode := codes.Unknown
	if st, ok := status.FromError(err); ok {
		statusCode = st.Code()
	}
	slog.Info(fmt.Sprintf("Status Code: %s", statusCode.String()))

	return result, err
}
