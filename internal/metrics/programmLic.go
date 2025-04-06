package metrics

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/prometheus/client_golang/prometheus"
)

type ProgrammLic struct {
	allLic *prometheus.CounterVec
	licSum prometheus.Counter
}

func NewProgrammLic() func() Metric {
	return func()Metric {
		return ProgrammLic{
			allLic: prometheus.NewCounterVec(
				prometheus.CounterOpts{
					Name: "programm_1c_lic_quantity",
					Help: "Data on 1C software licenses",
				},
				[]string{"name", "id"},
			),
			licSum : prometheus.NewCounter(
				prometheus.CounterOpts{
					Name: "programm_1c_lic_quantity_total",
					Help: "How much programm 1c lic we have total",
				},
			),
		}
	}
}

func(m ProgrammLic) Register(r *prometheus.Registry, params map[string]string) error {
	
	err := r.Register(m.allLic)
	if err != nil {
		return err
	}
	err = r.Register(m.licSum)
	if err != nil {
		return err
	}
	
	licFile, ok := params["pathToLicFile"]
	if !ok {
		return errors.New("path to lic file not found")
	}

	lic, err := parseLic(licFile)
	if err != nil {
		return err
	}

	var total float64
	for _, v := range lic {
		m.allLic.WithLabelValues(v.Name, v.ID).Add(v.Quantity)
		total += v.Quantity
	}

	m.licSum.Add(float64(total))

	return nil

}

func(m ProgrammLic) Unregister(r *prometheus.Registry, params map[string]string) error {
	r.Unregister(m.allLic)
	r.Unregister(m.licSum)

	return nil
}

type lic struct {
	Name string
	ID string
	Quantity float64
}

func parseLic(path string) ([]lic, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	lics := make([]lic, 0, 50)
	err = json.Unmarshal(data, &lics)

	if err != nil {
		return nil, err
	}

	return lics, nil
}