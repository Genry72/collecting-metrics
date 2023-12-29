package agent

import (
	"context"
	"fmt"
	"github.com/Genry72/collecting-metrics/helpers"
	"github.com/Genry72/collecting-metrics/internal/models"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// SetIpToHeader Добавление методанных с ip адресом
func SetIpToHeader(log *zap.Logger) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req interface{},
		reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption) error {

		localIP, err := helpers.GetLocalIP()
		if err != nil {
			log.Error("helpers.GetLocalIP", zap.Error(err))
			return fmt.Errorf("helpers.GetLocalIP: %w", err)
		}

		ctx = metadata.AppendToOutgoingContext(ctx, models.HeaderTrustedSubnet, localIP.String())

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
