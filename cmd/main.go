package main

import (
	"log/slog"
	"os"

	"github.com/behummble/1c_exporter/internal/app"
	"github.com/behummble/1c_exporter/internal/config"
	"github.com/behummble/1c_exporter/internal/metrics"
)

func main() {
	config := config.MustLoad()
	log := initLog()
	metrics := metrics.New(config, log)
	app := app.New(metrics, config, log)
	metrics.Register()
	app.Run()
}

func initLog() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, nil))
}
