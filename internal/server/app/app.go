package app

import (
	"github.com/golovanevvs/confidant/internal/server/repository"
	"github.com/golovanevvs/confidant/internal/server/service"
	"go.uber.org/zap"
)

func RunApp() {
	// initializing the logger
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	lg := logger.Sugar()

	// initializing the config
	cfg, err := newConfig()
	if err != nil {
		lg.Fatalf("application configuration initialization error: %s", err.Error())
	}

	// initializing the repository
	rp, err := repository.New(cfg.repository.DatabaseURI)
	if err != nil {
		lg.Fatalf("postgres DB initialization error: %s", err.Error())
	}

	// initializing the service
	sv := service.New(rp)
}
