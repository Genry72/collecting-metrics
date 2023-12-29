package agent

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"time"
)

func Logging(log *zap.Logger) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req interface{},
		reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption) error {
		// выполняем действия перед вызовом метода
		start := time.Now()

		// выполняем действия после вызова метода
		log.Info("Request", zap.String("method", method))

		// вызываем RPC-метод
		err := invoker(ctx, method, req, reply, cc, opts...)

		elapsed := time.Since(start).Seconds()

		// выполняем действия после вызова метода
		log.Info("Response",
			zap.String("code", status.Code(err).String()),
			zap.Float64("time", elapsed),
		)
		return err
	}
}
