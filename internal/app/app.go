package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/behummble/1c_exporter/internal/config"
	"golang.org/x/sys/windows/svc"
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
	//go svc.Run("1C_Programm_Lic_exporter", s)
	http.Handle("/metrics", s.metricService.Handler())
	if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
		s.log.Error(err.Error())
	}

}

func(s *ExporterService) Execute(args []string, r <-chan svc.ChangeRequest, status chan<- svc.Status) (bool, uint32) {

    const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown | svc.AcceptPauseAndContinue

	s.log.Info("start service...")

    status <- svc.Status{State: svc.StartPending}

    status <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}

	loop:
		for {
			select {
			case c := <-r:
				switch c.Cmd {
				case svc.Interrogate:
					status <- c.CurrentStatus
				case svc.Stop, svc.Shutdown:
					s.log.Info("shutdown service...")
					s.stop()
					break loop
				case svc.Pause:
					status <- svc.Status{State: svc.Paused, Accepts: cmdsAccepted}
				case svc.Continue:
					status <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
				default:
					s.log.Info(fmt.Sprintf("Unexpected service control request #%d", c))
				}
			}
		}

    status <- svc.Status{State: svc.StopPending}
    return false, 1
}


func(s *ExporterService) stop() {
	err := s.server.Shutdown(context.Background())
	if err != nil {
		s.log.Error(err.Error())
	}
	s.metricService.Unregister()
}
