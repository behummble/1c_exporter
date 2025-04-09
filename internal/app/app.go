package app

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/behummble/1c_exporter/internal/config"
)

type ExporterService struct {
	metricService Metric
	server *http.Server
	log *slog.Logger
}

type Metric interface {
	Register()
	Unregister()
	Handler() http.Handler
}

func New(metricService Metric, config *config.Config, log *slog.Logger) *ExporterService {
	server := &http.Server{
		Addr: fmt.Sprintf("%s:%s", config.Server.Addres, config.Server.Port),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}
	service := &ExporterService{
		metricService: metricService,
		server: server,
	}

	return service
}

func(s *ExporterService) Run() {
	http.Handle("/metrics", s.metricService.Handler())
	if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
		s.log.Error(err.Error())
	}

}
