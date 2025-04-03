package app

import (
	"log/slog"
	"net/http"
	"time"
	"fmt"
	"log"

	"github.com/behummble/1c_exporter/internal/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
}

func New(metricService Metric, config *config.Config, log *slog.Logger) *ExporterService {
	server := &http.Server{
		Addr: fmt.Sprintf("%s:%s", "localhost", "8152"),
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
	go svc.Run("1C_Programm_Lic_exporter", s)
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
	log.Fatal(s.server.ListenAndServe())
}

func(s *ExporterService) Execute(args []string, r <-chan svc.ChangeRequest, status chan<- svc.Status) (bool, uint32) {

    const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown | svc.AcceptPauseAndContinue

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
					s.stop()
					log.Print("Shutting service...!")
					break loop
				case svc.Pause:
					status <- svc.Status{State: svc.Paused, Accepts: cmdsAccepted}
				case svc.Continue:
					status <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
				default:
					log.Printf("Unexpected service control request #%d", c)
				}
			}
		}

    status <- svc.Status{State: svc.StopPending}
    return false, 1
}


func(s *ExporterService) stop() {

}
