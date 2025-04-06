package metrics

import (
	"log/slog"
	"net/http"

	"github.com/behummble/1c_exporter/internal/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metric interface {
	Register(r *prometheus.Registry, params map[string]string) error
	Unregister(r *prometheus.Registry, params map[string]string) error
}

type MetricService struct {
	metrics []Metric
	log *slog.Logger
	register *prometheus.Registry
	params map[string]string
}

func New(config *config.Config, log *slog.Logger) *MetricService {
	reg := prometheus.NewRegistry()
	names := metricsNames(config.Metrics)
	metrics := newMetrics(names)
	return &MetricService{
		register: reg,
		metrics: metrics,
		params: newParameters(config.Metrics),
		log: log,
	}
}

func newMetrics(names []string) []Metric {
	allMetrics := allMetrics()
	metrics := make([]Metric, 0, len(allMetrics))
	for _, v := range names {
		if newMetricFunc, ok := allMetrics[v]; ok {
			metrica :=  newMetricFunc()
			metrics = append(metrics, metrica)
		}
	}

	return metrics
}

func newParameters(metrics []config.MetricConfig) map[string]string {
	res := make(map[string]string, len(metrics))
	for _, v := range metrics {
		res[v.Options.Name] = v.Options.Value
	}
	return res
}

func(m *MetricService) Register() {
	for _, metrica := range m.metrics {
		err := metrica.Register(m.register, m.params)
		if err != nil {
			m.log.Error(err.Error())
		}
	}
}

func(m *MetricService) Unregister() {
	for _, metrica := range m.metrics {
		err := metrica.Unregister(m.register, m.params)
		if err != nil {
			m.log.Error(err.Error())
		}
	}
}

func(m *MetricService) Handler() http.Handler {
	return promhttp.HandlerFor(m.register, promhttp.HandlerOpts{Registry: m.register})
}

func allMetrics() map[string]func() Metric {
	metrics := make(map[string]func()Metric, 30)
	metrics["programm_lic_1C"] = NewProgrammLic()
	return metrics
}

func metricsNames(metrics []config.MetricConfig) []string {
	names := make([]string, len(metrics))
	for k, v := range metrics {
		names[k] = v.Name
	}
	return names
}