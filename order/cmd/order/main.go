package main

import (
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/silven-dynamics/go-ecommerce/order"
	"github.com/tinrab/retry"
)

type Config struct {
	DatabaseURL string `envconfig:"DATABASE_URL"`
	AccountURL  string `envconfig:"ACCOUNT_SERVICE_URL"`
	CatalogURL  string `envconfig:"CATALOG_SERVICE_URL"`
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	var r order.OrderRepository
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		r, err = order.NewPostgresRepository(cfg.DatabaseURL)
		if err != nil {
			log.Println(err)
		}
		return
	})
	defer r.Close()

	log.Println("Listening on port 8001...")
	s := order.NewOrderService(r)
	log.Fatal(order.ListenGRPC(
		s,
		cfg.AccountURL,
		cfg.CatalogURL,
		8001),
	)
}
