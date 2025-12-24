package main

import (
	"log"

	"github.com/stefanobassani-dev/money-tracker/internal/api"
	"github.com/stefanobassani-dev/money-tracker/internal/config"
)

func main() {
	cfg := config.Load()
	srv := api.NewServer(cfg)

	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}
