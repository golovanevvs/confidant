package app

import (
	"fmt"

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
	fmt.Println()
}
