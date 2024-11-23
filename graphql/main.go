package main

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/kelseyhightower/envconfig"
)

type AppConfig struct {
	AccountURL string `envconfig:"ACCOUNT_SERVICE_URL" required:"true"`
	CatalogURL string `envconfig:"CATALOG_SERVICE_URL" required:"true"`
	OrderURL   string `envconfig:"ORDER_SERVICE_URL" required:"true"`
}

func main() {
	var cfg AppConfig
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	s, err := NewGraphQLServer(cfg.AccountURL, cfg.CatalogURL, cfg.OrderURL)
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/graphql", handler.New(s.ToExecutableSchema()))
	http.Handle("/playground", playground.Handler("go-ecommerce", "/graphql"))

	log.Fatal(http.ListenAndServe(":8001", nil))
	log.Println("Listening on port 8001...")
}
