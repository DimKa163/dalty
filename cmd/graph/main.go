package main

import (
	"github.com/DimKa163/dalty/app/graph"
	"github.com/caarlos0/env"
)

func main() {
	var cfg graph.Config

	err := env.Parse(&cfg)
	if err != nil {
		panic(err)
	}
	server := graph.NewServer(&cfg)
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
