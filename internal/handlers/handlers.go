package handlers

import (
	"errors"
	"github.com/Genry72/collecting-metrics/internal/usecases"
	"net/http"
	"strings"
)

var (
	ErrOnlyPost  = errors.New("only POST requests are allowed")
	ErrFormatURL = errors.New("only /update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ> format are allowed")
)

type Handler struct {
	useCases *usecases.Usecase
}

func New(uc *usecases.Usecase) *Handler {
	return &Handler{
		useCases: uc,
	}
}

func (h Handler) setMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, ErrOnlyPost.Error(), http.StatusMethodNotAllowed)
		return
	}

	const countURL = 5

	urlSlice := strings.Split(r.URL.Path, "/")
	if len(urlSlice) != countURL {
		http.Error(w, ErrFormatURL.Error(),
			http.StatusNotFound)
		return
	}

	const (
		metricType = iota + 2
		metricName
		metricValue
	)

	if err := h.useCases.SetMetric(urlSlice[metricType], urlSlice[metricName], urlSlice[metricValue]); err != nil {
		var status int
		if errors.Is(err, usecases.ErrBadMetricType) || errors.Is(err, usecases.ErrParseValue) {
			status = http.StatusBadRequest
		} else {
			status = http.StatusInternalServerError
		}

		http.Error(w, err.Error(), status)

		return
	}
}
