package main

import (
	"log/slog"
	"net/http"

	"github.com/behummble/1c_exporter/internal/config"
	"github.com/behummble/1c_exporter/internal/metrics"
	"github.com/behummble/1c_exporter/internal/app"

)

func main() {
	config := config.New()
	log := initLog()
	metrics := metrics.New(config, log)
	metrics.Register()
	app := app.New(metrics, config, log)
	app.Run()
}

func initLog() *slog.Logger {

}

func newMetrics(lic []lic, reg *prometheus.Registry) *metrics {
	
}
