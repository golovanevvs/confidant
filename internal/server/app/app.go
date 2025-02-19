package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golovanevvs/confidant/internal/server/repository"
	"github.com/golovanevvs/confidant/internal/server/service"
	handler "github.com/golovanevvs/confidant/internal/server/transport/http"
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
	//initializing the handler
	hd := handler.New(sv, lg)
	// initializing the server
	srv := newServer()

	// starting the server
	go func() {
		lg.Infof("The \"confidant\" server is running")
		if err := srv.RunServer(":8080", hd.InitRoutes(lg)); err != nil {
			lg.Fatalf("server startup error: %s", err.Error())
		}
	}()

	// shutting down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	lg.Infof("A server shutdown signal has been received")

	if err := srv.ShutdownServer(context.Background()); err != nil {
		lg.Errorf("error when shutting down the server: %s", err.Error())
	}

	if err := rp.IManageRepository.CloseDB(); err != nil {
		lg.Errorf("error when shutting down the DB: %v", err.Error())
	}

	lg.Infof("the server operation is completed")
	time.Sleep(time.Second * 2)
}
