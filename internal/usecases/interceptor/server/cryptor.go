package server

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"github.com/Genry72/collecting-metrics/internal/models"
	"github.com/Genry72/collecting-metrics/internal/usecases/cryptor"
	pb "github.com/Genry72/collecting-metrics/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

/*
CheckHashFromHeader Проверка входящего хедера HashSHA256. Если передан, то сверяем сумму хэша с телом запроса.
Считаем хэш тела ответа и добавляем заголовок в ответ.
*/
func CheckHashFromHeader(log *zap.Logger, password string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {
		var headerHash string

		if md, ok := metadata.FromIncomingContext(ctx); ok {
			values := md.Get(models.HeaderHash)
			if len(values) > 0 {
				headerHash = values[0]
			}
		}

		if headerHash == "" {
			return handler(ctx, req)
		}

		requestBody, err := json.Marshal(req)
		if err != nil {
			log.Error("json.Marshal", zap.Error(err))
			return nil, status.Errorf(codes.InvalidArgument, "failed to marshal request body: %v", err)
		}

		hashFromBody, err := cryptor.Encrypt(requestBody, password)
		if err != nil {
			log.Error("cryptor.Encrypt hashFromBody", zap.Error(err))
			return nil, status.Errorf(codes.InvalidArgument, "cryptor.Encrypt: %v", err)
		}

		if headerHash != hashFromBody {
			log.Error("headerHash != hashFromBody", zap.Error(models.ErrHashNotEqual))
			return nil, status.Errorf(codes.InvalidArgument, "headerHash != hashFromBody: %v", err)
		}

		// Вызов метода обработчика
		resp, err := handler(ctx, req)
		if err != nil {
			return resp, err
		}

		responseBody, err := json.Marshal(resp)
		if err != nil {
			log.Error("json.Marshal", zap.Error(err))
			return nil, status.Errorf(codes.InvalidArgument, "failed to marshal request body: %v", err)
		}

		respHash, err := cryptor.Encrypt(responseBody, password)
		if err != nil {
			log.Error("cryptor.Encrypt respHash", zap.Error(err))
			return nil, status.Errorf(codes.Unauthenticated, "failed to marshal request body: %v", err)
		}

		// Добавляем хеш в метаданные отправляемого запроса
		md := metadata.New(map[string]string{models.HeaderHash: respHash})
		if err := grpc.SetHeader(ctx, md); err != nil {
			return nil, status.Errorf(codes.Internal, "grpc.SetHeader: %v", err)
		}

		return resp, err
	}
}

// DecryptBodyWithPrivateKey декодирование тела входящего запроса при помощи закрытого ключа
func DecryptBodyWithPrivateKey(log *zap.Logger, privateKey *rsa.PrivateKey) grpc.UnaryServerInterceptor {
	return func(ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {
		// Считываем закодированное тело запроса
		encryptedBody, ok := req.(*pb.EncryptedMessage)
		if !ok {
			return handler(ctx, req)
		}

		// Расшифровываем
		decryptedData, err := cryptor.DecryptWithPrivateKey(encryptedBody.Data, privateKey)
		if err != nil {
			log.Error("grpc cryptor.DecryptWithPrivateKey", zap.Error(err))
			return nil, status.Errorf(codes.Unauthenticated, "cryptor.DecryptWithPrivateKey: %v", err)
		}

		encryptedBody.Data = decryptedData
		req = encryptedBody

		return handler(ctx, req)
	}
}
