package main

import (
	"github.com/DimKa163/dalty/app/rest"
	"github.com/caarlos0/env"
)

func main() {
	var cfg rest.Config

	err := env.Parse(&cfg)
	if err != nil {
		panic(err)
	}
	server, err := rest.NewServer(&cfg)
	if err != nil {
		panic(err)
	}
	server.AddServices()
	if err := server.AddLogging(); err != nil {
		panic(err)
	}
	server.Map()
	if err := server.Run(); err != nil {
		panic(err)
	}
}
