package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type ProgrammLic struct {
	allLic *prometheus.CounterVec
	licSum prometheus.Counter
	pathToLic string
}

func NewProgrammLic(pathToLic string) ProgrammLic {
	return ProgrammLic{}
}

func(m ProgrammLic) Register(r *prometheus.Registry) ProgrammLic {
	allLic := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "programm_1c_lic_quantity",
			Help: "How ",
		},
		[]string{"name", "id"},
	)
	licSum := prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "programm_1c_lic_quantity",
			Help: "How much programm 1c lic we have total",
		},
	)
	r.MustRegister(allLic)
	r.MustRegister(licSum)
	total := 0
	for _, v := range lic {
		allLic.WithLabelValues(v.name, v.id).Add()
		total++
	}
	licSum.Add(float64(total))

	return &metrics{
		allLic: allLic,
		licSum: licSum,
	}
}

func(m ProgrammLic) Unregister() {

}

type lic struct {
	name string
	id string
	quantity int
}