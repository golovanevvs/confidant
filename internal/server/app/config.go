package app

import (
	"flag"
	"os"
)

type config struct {
	server     serverConfig
	repository repositoryConfig
}

type serverConfig struct {
	Addr string
}

type repositoryConfig struct {
	DatabaseURI string
}

func NewConfig() *config {
	var flagRunAddr, flagDatabaseURI string

	flag.StringVar(&flagRunAddr, "a", ":7541", "address and port of server")
	flag.StringVar(&flagDatabaseURI, "d", "", "database DSN")
	flag.Parse()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		flagRunAddr = envRunAddr
	}
	if envDatabaseURI := os.Getenv("DATABASE_DSN"); envDatabaseURI != "" {
		flagDatabaseURI = envDatabaseURI
	}

	return &config{
		serverConfig{
			Addr: flagRunAddr,
		},
		repositoryConfig{
			DatabaseURI: flagDatabaseURI,
		},
	}
}
