package middleware

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
)

func LoggingUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()

	// Call handler
	resp, err := handler(ctx, req)

	// Log the request details
	duration := time.Since(start)
	log.Printf(
		"[gRPC] Method=%s | Duration=%dms | Error=%v",
		info.FullMethod,
		duration.Milliseconds(),
		err,
	)

	return resp, err
}

