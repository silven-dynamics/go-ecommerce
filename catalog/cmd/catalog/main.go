package main

import (
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/stiffinWanjohi/go-ecommerce/catalog"
	"github.com/tinrab/retry"
)

type Config struct {
	DatabaseURL string `envconfig:"DATABASE_URL"`
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Elasticsearch URL:", cfg.DatabaseURL)

	var r catalog.CatalogRepository
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		r, err = catalog.NewElasticRepository(cfg.DatabaseURL)
		if err != nil {
			log.Println(err)
		}
		return
	})
	defer r.Close()

	log.Println("Listening on port 8001...")
	s := catalog.NewCatalogService(r)
	log.Fatal(catalog.ListenGRPC(s, 8001))
}
