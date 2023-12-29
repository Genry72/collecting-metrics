package handlers

import (
	"bytes"
	"context"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"github.com/Genry72/collecting-metrics/internal/models"
	cryptor2 "github.com/Genry72/collecting-metrics/internal/usecases/cryptor"
	interceptor "github.com/Genry72/collecting-metrics/internal/usecases/interceptor/server"
	"os"

	"github.com/Genry72/collecting-metrics/internal/usecases/server"
	pb "github.com/Genry72/collecting-metrics/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"log"
	"net"
)

type Server struct {
	// нужно встраивать тип pb.Unimplemented<TypeName>
	// для совместимости с будущими версиями
	pb.UnimplementedServerServer
	useCases      *server.Server
	log           *zap.Logger
	keyHash       *string
	privateKey    *rsa.PrivateKey
	trustedSubnet *string
	interceptors  []grpc.UnaryServerInterceptor
}

func NewGrpsServer(useCases *server.Server, log *zap.Logger, keyHash *string, privateKeyPath, trustedSubnet *string) *Server {
	var (
		privKey *rsa.PrivateKey
		err     error
	)
	// При передаче ключа, подключаем обработчик по расшифровке тела запроса приватным ключем
	if privateKeyPath != nil && *privateKeyPath != "" {
		// Проверка наличия ключа по указанному пути, если его нет, то генерируем новый набор
		privKey, err = cryptor2.GetPrivateKeyFromFile(*privateKeyPath)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				log.Fatal("The path to the private key provided by the file was not found")
			} else {
				log.Fatal("cert.GetPrivateKeyFromFile", zap.Error(err))
			}
		}
	}

	return &Server{
		useCases:      useCases,
		keyHash:       keyHash,
		log:           log,
		privateKey:    privKey,
		trustedSubnet: trustedSubnet,
	}
}

func (h *Server) RunServer(hostPort string) {
	listen, err := net.Listen("tcp", hostPort)
	if err != nil {
		log.Fatal(err)
	}
	// Логирование
	h.use(interceptor.Logging(h.log))

	// Проверка ip адреса отправителя в метаданных
	if h.trustedSubnet != nil && *h.trustedSubnet != "" {
		h.use(interceptor.CheckIpFromHeader(h.log, *h.trustedSubnet))
	}

	// Проверка хеша тела запроса
	if h.keyHash != nil && *h.keyHash != "" {
		h.use(interceptor.CheckHashFromHeader(h.log, *h.keyHash))
	}

	// Расшифровываем тело запроса
	if h.privateKey != nil {
		h.use(interceptor.DecryptBodyWithPrivateKey(h.log, h.privateKey))
	}

	// создаём gRPC-сервер, подключаем обработчики
	s := grpc.NewServer(grpc.ChainUnaryInterceptor(h.interceptors...))

	// Регистрация серверного отражения
	reflection.Register(s)

	// регистрируем сервис
	pb.RegisterServerServer(s, h)

	h.log.Info("Сервер gRPC начал работу")
	// получаем запрос gRPC
	if err := s.Serve(listen); err != nil {
		h.log.Error("s.Serve", zap.Error(err))
	}
}

func (h *Server) GetAllMetrics(ctx context.Context, in *pb.EmptyMessage) (*pb.String, error) {
	var response pb.String

	val, statusCode, err := h.useCases.GetAllMetrics(ctx)
	if err != nil {
		h.log.Error("h.useCases.GetAllMetrics", zap.Error(err))
		return nil, checkStatus(statusCode, err.Error())
	}

	b := bytes.Buffer{}

	err = tmpl.Execute(&b, val)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "status not ok: %d body: %s", statusCode, err.Error())
	}
	//
	response.Result = b.String()

	return &response, nil
}

func (h *Server) SetMetrics(ctx context.Context, in *pb.Metrics) (*pb.EmptyMessage, error) {
	metricParams := pb.MetricsDpToMetrics(in)

	statusCode, err := h.useCases.SetMetric(ctx, metricParams...)
	if err != nil {
		return &pb.EmptyMessage{}, checkStatus(statusCode, err.Error())
	}

	return &pb.EmptyMessage{}, nil
}

func (h *Server) SetMetricsEncrypted(ctx context.Context, in *pb.EncryptedMessage) (*pb.EmptyMessage, error) {
	var metricParams []*models.Metric

	if err := json.Unmarshal([]byte(in.Data), &metricParams); err != nil {
		h.log.Error("json.Unmarshal", zap.Error(err))
		return nil, status.Errorf(codes.Unimplemented, "json.Unmarshal: %v", err)
	}

	statusCode, err := h.useCases.SetMetric(ctx, metricParams...)
	if err != nil {
		return nil, checkStatus(statusCode, err.Error())
	}

	return &pb.EmptyMessage{}, nil
}

func checkStatus(statusCode int, body string) error {
	switch {
	case statusCode >= 200 && statusCode < 400:
		return nil
	case statusCode >= 400 && statusCode < 500:
		// повтор не нужен
		return status.Errorf(codes.Unimplemented, "status not ok: %d body: %s", statusCode, body)
	case statusCode >= 500:
		// нужен повтор
		err := status.Errorf(codes.Internal, "status not ok: %d body: %s", statusCode, body)
		return models.NewRetryError(err)
	default:
		return nil
	}
}

func (h *Server) use(g grpc.UnaryServerInterceptor) {
	h.interceptors = append(h.interceptors, g)
}
