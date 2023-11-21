package main

import (
	"context"
	"github.com/Genry72/collecting-metrics/internal/logger"
	"github.com/Genry72/collecting-metrics/internal/usecases/agent"
	"net/http"
	"net/http/httptest"
	"runtime"
	"testing"
	"time"
)

func BenchmarkHandlers(b *testing.B) {
	runtime.MemProfileRate = 0
	zapLogger := logger.NewZapLogger("info")
	keyHash := "fagfrgthagarfafserghaegferfegraferwfreagagrag"
	s := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		}),
	)
	defer s.Close()

	runtime.MemProfileRate = 1

	b.Run("new", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			runtime.MemProfileRate = 0
			metrics := agent.NewMetrics()
			agentUc := agent.NewAgent(s.URL, zapLogger, &keyHash, 100)

			runtime.MemProfileRate = 1
			ctx, c1 := context.WithTimeout(context.Background(), 10*time.Millisecond)
			metrics.Update(ctx, 1*time.Microsecond)
			c1()
			ctx2, c2 := context.WithTimeout(context.Background(), 10*time.Millisecond)
			agentUc.SendMetrics(ctx2, metrics, 1*time.Microsecond)
			c2()
		}
	})

}
