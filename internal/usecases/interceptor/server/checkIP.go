package server

import (
	"context"
	"github.com/Genry72/collecting-metrics/internal/models"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"net"
)

// CheckIpFromHeader Проверка методанных с ip адресом
func CheckIpFromHeader(log *zap.Logger, trustedSubnet string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {

		var ipFromHeader string

		if md, ok := metadata.FromIncomingContext(ctx); ok {
			values := md.Get(models.HeaderTrustedSubnet)
			if len(values) > 0 {
				ipFromHeader = values[0]
			}
		}

		if ipFromHeader == "" {
			log.Error("empty " + models.HeaderTrustedSubnet)
			return nil, status.Errorf(codes.PermissionDenied, "empty %v", models.HeaderTrustedSubnet)
		}

		ip := net.ParseIP(ipFromHeader)

		_, subnet, err := net.ParseCIDR(trustedSubnet)
		if err != nil {
			log.Error("net.ParseCIDR", zap.Error(err))
			return nil, status.Errorf(codes.PermissionDenied, "net.ParseCIDR %v", err)
		}

		if subnet.Contains(ip) {
			return handler(ctx, req)
		}

		log.Error("ip not trusted")

		return nil, status.Errorf(codes.PermissionDenied, "ip not trusted")
	}
}
