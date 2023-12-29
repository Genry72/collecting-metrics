package agent

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"github.com/Genry72/collecting-metrics/internal/models"
	"github.com/Genry72/collecting-metrics/internal/usecases/cryptor"
	pb "github.com/Genry72/collecting-metrics/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// SetHashToHeader Добавление методанных с хешем расчитаннго тела запроса
func SetHashToHeader(log *zap.Logger, password string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req interface{},
		reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption) error {

		// Для всех сообщений считаем хеш и добавляем в метаданные
		bodyByte, err := json.Marshal(req)
		if err != nil {
			log.Error("json.Marshal", zap.Error(err))
			return fmt.Errorf("json.Marshal: %w", err)
		}

		hash, err := cryptor.Encrypt(bodyByte, password)
		if err != nil {
			log.Error("cryptor.Encrypt", zap.Error(err))
			return fmt.Errorf("cryptor.Encrypt: %w", err)
		}

		ctx = metadata.AppendToOutgoingContext(ctx, models.HeaderHash, hash)

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// EncryptBodyWithPublicKey кодирование тела исходящего запроса при помощи открытого ключа
func EncryptBodyWithPublicKey(log *zap.Logger, publicKey *rsa.PublicKey) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req interface{},
		reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption) error {
		// Считываем тело запроса
		decryptedBody, ok := req.(*pb.EncryptedMessage)
		if !ok {
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		// Зашифровываем
		encrypteddData, err := cryptor.EncryptBodyWithPublicKey(decryptedBody.Data, publicKey)
		if err != nil {
			log.Error("cryptor.EncryptBodyWithPublicKey", zap.Error(err))
			return fmt.Errorf("cryptor.EncryptBodyWithPublicKey: %w", err)
		}

		decryptedBody.Data = encrypteddData

		req = decryptedBody

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
