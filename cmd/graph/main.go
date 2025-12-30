package main

import (
	"github.com/DimKa163/dalty/app/warehouse"
	"github.com/caarlos0/env"
)

func main() {
	var cfg warehouse.Config

	err := env.Parse(&cfg)
	if err != nil {
		panic(err)
	}
	server := warehouse.NewServer(&cfg)
	if err := server.AddServices(); err != nil {
		panic(err)
	}
	if err := server.AddLogging(); err != nil {
		panic(err)
	}
	server.Map()
	if err := server.Run(); err != nil {
		panic(err)
	}
}
