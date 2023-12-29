package server

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"time"
)

func Logging(log *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {
		// выполняем действия перед вызовом метода
		start := time.Now()

		log.Info("Request",
			zap.String("method", info.FullMethod),
		)

		// Вызов метода обработчика
		resp, err := handler(ctx, req)

		elapsed := time.Since(start).Seconds()

		log.Info("Response",
			zap.String("code", status.Code(err).String()),
			zap.Float64("time", elapsed),
		)

		return resp, err
	}
}
