package main

import (
	"github.com/stefanobassani-dev/money-tracker/internal/api"
	"github.com/stefanobassani-dev/money-tracker/internal/config"
)

func main() {
	cfg := config.Load()

	server := api.NewServer(cfg)
	server.Run()
}
