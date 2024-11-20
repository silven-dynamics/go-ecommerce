package main

import (
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/silven-dynamics/go-ecommerce/account"
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

	var r account.AccountRepository
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		r, err = account.NewPostgresRepository(cfg.DatabaseURL)
		if err != nil {
			log.Println(err)
		}
		return
	})
	defer r.Close()

	log.Println("Listening on port 8001...")
	s := account.NewAccountService(r)
	log.Fatal(account.ListenGRPC(s, 8001))
}
