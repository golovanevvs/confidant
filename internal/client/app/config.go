package app

import (
	"flag"
	"os"
)

type config struct {
	addr string
}

func NewConfig() *config {
	var flagRunAddr string

	flag.StringVar(&flagRunAddr, "a", "localhost:7541", "address and port of server")
	flag.Parse()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		flagRunAddr = envRunAddr
	}

	return &config{
		addr: flagRunAddr,
	}
}
