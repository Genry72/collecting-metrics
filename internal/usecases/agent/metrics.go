package agent

import (
	"fmt"
	"github.com/Genry72/collecting-metrics/internal/models"
	"github.com/fatih/structs"
	"math/rand"
	"runtime"
	"time"
)

type Metrics struct {
	gauge   *gaugeRunTimeMetrics
	counter *counterMetrics
	rtm     *runtime.MemStats
}

func NewMetrics() *Metrics {
	return &Metrics{
		gauge:   &gaugeRunTimeMetrics{},
		counter: &counterMetrics{},
		rtm:     &runtime.MemStats{},
	}
}

type gaugeRunTimeMetrics struct {
	Alloc         uint64
	BuckHashSys   uint64
	Frees         uint64
	GCCPUFraction float64
	GCSys         uint64
	HeapAlloc     uint64
	HeapIdle      uint64
	HeapInuse     uint64
	HeapObjects   uint64
	HeapReleased  uint64
	HeapSys       uint64
	LastGC        uint64
	Lookups       uint64
	MCacheInuse   uint64
	MCacheSys     uint64
	MSpanInuse    uint64
	MSpanSys      uint64
	Mallocs       uint64
	NextGC        uint64
	NumForcedGC   uint32
	NumGC         uint32
	OtherSys      uint64
	PauseTotalNs  uint64
	StackInuse    uint64
	StackSys      uint64
	Sys           uint64
	TotalAlloc    uint64
	RandomValue   float64
}

type counterMetrics struct {
	PollCount int64
}

// Получение улов для отправки метрик
func (m *Metrics) getUrlsMetric() []string {
	gaugeMetricData := structs.Map(m.gauge)
	counterMetricsData := structs.Map(m.counter)

	result := make([]string, 0, len(gaugeMetricData)+len(counterMetricsData))

	for metricName, value := range gaugeMetricData {
		url := fmt.Sprintf("/update/%s/%s/%v", models.MetricTypeGauge, metricName, value)
		result = append(result, url)
	}

	for metricName, value := range counterMetricsData {
		url := fmt.Sprintf("/update/%s/%s/%v", models.MetricTypeCounter, metricName, value)
		result = append(result, url)
	}

	return result
}

// Update Запуск обновления метрик с заданным интервалом
func (m *Metrics) Update(pollInterval time.Duration) {
	go func() {
		for {
			m.updateMetics()
			time.Sleep(pollInterval)
		}
	}()
}

func (m *Metrics) updateMetics() {
	runtime.ReadMemStats(m.rtm)
	m.gauge.Alloc = m.rtm.Alloc
	m.gauge.BuckHashSys = m.rtm.BuckHashSys
	m.gauge.Frees = m.rtm.Frees
	m.gauge.GCCPUFraction = m.rtm.GCCPUFraction
	m.gauge.GCSys = m.rtm.GCSys
	m.gauge.HeapAlloc = m.rtm.HeapAlloc
	m.gauge.HeapIdle = m.rtm.HeapIdle
	m.gauge.HeapInuse = m.rtm.HeapInuse
	m.gauge.HeapObjects = m.rtm.HeapObjects
	m.gauge.HeapReleased = m.rtm.HeapReleased
	m.gauge.HeapSys = m.rtm.HeapSys
	m.gauge.LastGC = m.rtm.LastGC
	m.gauge.Lookups = m.rtm.Lookups
	m.gauge.MCacheInuse = m.rtm.MCacheInuse
	m.gauge.MCacheSys = m.rtm.MCacheSys
	m.gauge.MSpanInuse = m.rtm.MSpanInuse
	m.gauge.MSpanSys = m.rtm.MSpanSys
	m.gauge.Mallocs = m.rtm.Mallocs
	m.gauge.NextGC = m.rtm.NextGC
	m.gauge.NumForcedGC = m.rtm.NumForcedGC
	m.gauge.NumGC = m.rtm.NumGC
	m.gauge.OtherSys = m.rtm.OtherSys
	m.gauge.PauseTotalNs = m.rtm.PauseTotalNs
	m.gauge.StackInuse = m.rtm.StackInuse
	m.gauge.StackSys = m.rtm.StackSys
	m.gauge.Sys = m.rtm.Sys
	m.gauge.TotalAlloc = m.rtm.TotalAlloc
	m.counter.PollCount++

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	m.gauge.RandomValue = r.Float64()
}