package server

import (
	"context"
	"fmt"
	"github.com/Genry72/collecting-metrics/internal/models"
	"github.com/Genry72/collecting-metrics/internal/repositories"

	"strings"
)

type Server struct {
	storage repositories.Repositories
}

func NewServerUc(repo repositories.Repositories) *Server {
	return &Server{
		storage: repo,
	}
}

func (uc *Server) SetMetric(ctx context.Context, metric *models.Metrics) (*models.Metrics, error) {

	return uc.storage.SetMetric(ctx, metric)
}

func (uc *Server) GetMetricValue(ctx context.Context, metric *models.Metrics) (*models.Metrics, error) {

	return uc.storage.GetMetricValue(ctx, metric)

}

func (uc *Server) GetAllMetrics(ctx context.Context) (string, error) {
	mapa, err := uc.storage.GetAllMetrics(ctx)
	if err != nil {
		return "", err
	}

	sb := strings.Builder{}

	for k, v := range mapa {
		sb.WriteString(fmt.Sprintf("%s : %v\n", k, v))
	}

	return sb.String(), nil
}
