package app

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golovanevvs/confidant/internal/server/repository"
	"github.com/golovanevvs/confidant/internal/server/service"
	"github.com/golovanevvs/confidant/internal/server/transport/http/handler"
	"go.uber.org/zap"
)

type IHandler interface {
	InitRoutes() http.Handler
}

func RunApp() {
	// initializing the logger
	rawJSON := []byte(`{
		"level": "debug",
		"outputPaths": ["stdout"],
		"errorOutputPaths": ["stderr"],
		"encoding": "json",
		"encoderConfig": {
			"messageKey": "message",
			"levelKey": "level",
			"levelEncoder": "lowercase"
		}
	}`)
	var cfgZap zap.Config
	if err := json.Unmarshal(rawJSON, &cfgZap); err != nil {
		panic(err)
	}
	logger := zap.Must(cfgZap.Build())
	defer logger.Sync() // flushes buffer, if any
	lg := logger.Sugar()

	// initializing the config
	cfg := NewConfig()

	//initializing the repository
	rp, err := repository.New(cfg.repository.DatabaseURI)
	if err != nil {
		lg.Fatalf("postgres DB initialization error: %s", err.Error())
	}
	lg.Infof("Connecting to DB: success")

	// initializing the service
	sv := service.New(rp)

	//initializing the handler
	hd := handler.New(sv, lg)

	// initializing the server
	srv := newServer()

	// starting the server
	go func() {
		lg.Infof("The 'confidant' server is running, address: %s", cfg.server.Addr)
		if err := srv.RunServer(cfg.server.Addr, hd.InitRoutes()); err != nil {
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

	if err := rp.CloseDB(); err != nil {
		lg.Errorf("error when shutting down the DB: %v", err.Error())
	}

	lg.Infof("the server operation is completed")
	time.Sleep(time.Second * 1)
}
