package grpcclients

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"github.com/Genry72/collecting-metrics/internal/models"
	"github.com/Genry72/collecting-metrics/internal/usecases/cryptor"
	"github.com/Genry72/collecting-metrics/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type GrpcClient struct {
	grpClient proto.ServerClient
	grpcconn  *grpc.ClientConn
	log       *zap.Logger
	keyHash   *string
	publicKey *rsa.PublicKey
}

func NewGrpcClient(grpcconn *grpc.ClientConn, log *zap.Logger, keyHash *string, publicKeyPath *string) (*GrpcClient, error) {
	var (
		publicLey *rsa.PublicKey
		err       error
	)

	if publicKeyPath != nil && *publicKeyPath != "" {
		publicLey, err = cryptor.GetPubKeyFromFile(*publicKeyPath)
		if err != nil {
			return nil, fmt.Errorf("cryptor.GetPubKeyFromFile: %w", err)
		}
	}

	client := proto.NewServerClient(grpcconn)

	return &GrpcClient{
		grpcconn:  grpcconn,
		grpClient: client,
		log:       log,
		keyHash:   keyHash,
		publicKey: publicLey,
	}, nil
}

func (c *GrpcClient) Send(ctx context.Context, metric models.Metrics) error {

	metrycByte, err := json.Marshal(metric)
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}

	_, err = c.grpClient.SetMetricsEncrypted(ctx, &proto.EncryptedMessage{Data: metrycByte})
	if err != nil {
		if ctx.Err() != nil {
			return nil
		}

		c.log.Error("resp", zap.Error(err))
		// или сеть или тело ответа
		return models.NewRetryError(err)
	}

	return nil
}

func (c *GrpcClient) Stop() error {
	return c.grpcconn.Close()
}
