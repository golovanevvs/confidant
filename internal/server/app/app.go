package app

import (
	"context"
	"encoding/json"
	"fmt"
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
	_, err := NewConfig()
	if err != nil {
		lg.Fatalf("application configuration initialization error: %s", err.Error())
	}

	//initializing the repository
	var rp *repository.Repository
	for i := 5430; i <= 5440; i++ {
		databaseURI := fmt.Sprintf("host=localhost port=%d user=postgres password=password dbname=confidant sslmode=disable", i)
		lg.Debugf("Connecting to DB: port %d...", i)
		rp, err = repository.New(databaseURI)
		if err != nil {
			if i == 5440 {
				lg.Fatalf("postgres DB initialization error: %s", err.Error())
			}
			lg.Debugf("Connect to DB: error: %s", err.Error())
			lg.Debugf("Repeating...")
		} else {
			lg.Infof("Connecting to DB: success")
			break
		}
	}

	// initializing the service
	// accountSv := accountservice.New(rp.IAccountRepository)
	// sv := transport.NewService(accountSv)
	sv := service.New(rp)

	//initializing the handler
	hd := handler.New(sv, lg)

	// initializing the server
	srv := newServer()

	// starting the server
	go func() {
		lg.Infof("The 'confidant' server is running")
		if err := srv.RunServer(":8080", hd.InitRoutes()); err != nil {
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
	time.Sleep(time.Second * 2)
}
