package db

import (
	"github.com/hellofresh/health-go/v4"
	"log"
	"time"
)

func NewHealthCheck(pe PersistentDatabase, st StatsDatabase, ps PubSub) (h *health.Health, err error) {
	h, err = health.New()
	if err != nil {
		log.Fatalln("Health-Service:", err)
		return
	}
	// register persistent check
	if pe != nil {
		if err := h.Register(health.Config{
			Name:      "Persistent (" + pe.ServiceName() + ")",
			Timeout:   5 * time.Second,
			SkipOnErr: true,
			Check:     pe.HealthCheck,
		}); err != nil {
			log.Println("Error registering health service:", err)
		}
	}
	if st != nil {
		if err := h.Register(health.Config{
			Name:      "Stats (" + st.ServiceName() + ")",
			Timeout:   5 * time.Second,
			SkipOnErr: true,
			Check:     st.HealthCheck,
		}); err != nil {
			log.Println("Error registering health service:", err)
		}
	}
	if ps != nil {
		if err := h.Register(health.Config{
			Name:      "PubSub (" + ps.ServiceName() + ")",
			Timeout:   5 * time.Second,
			SkipOnErr: true,
			Check:     ps.HealthCheck,
		}); err != nil {
			log.Println("Error registering health service:", err)
		}
	}
	return
}
