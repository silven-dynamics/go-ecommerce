package main

import (
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/silven-dynamics/go-ecommerce/account"
	"github.com/tinrab/retry"
)

type Config struct {
	DatabaseURL string `envconfig:"DATABASE_URL"`
}

func main() {
	err := godotenv.Load("./account/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var cfg Config
	err = envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Database URL:", cfg.DatabaseURL)

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
