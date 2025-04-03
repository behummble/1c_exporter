package metrics

import (
	"log/slog"

	"github.com/behummble/1c_exporter/internal/config"
	"github.com/prometheus/client_golang/prometheus"
)

type Metric interface {
	Register(r *prometheus.Registry)
	Unregister(r *prometheus.Registry)
}

type MetricService struct {
	metrics []Metric
	register *prometheus.Registry
	configPath string
}

func New(config *config.Config, log *slog.Logger) *MetricService {
	reg := prometheus.NewRegistry()
}

func(m *MetricService) Register() {

}

func(m *MetricService) Unregister() {

}

func readLic() []lic {
	data, err := os.ReadFile("./lic.json")
	if err != nil {
		panic(err)
	}
	lics := make([]lic, 0, 50)
	err = json.Unmarshal(data, lics)

	if err != nil {
		panic(err)
	}

	return lics
}