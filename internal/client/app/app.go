package app

import (
	"encoding/json"

	"github.com/golovanevvs/confidant/internal/client/app/appview"
	"github.com/golovanevvs/confidant/internal/client/service"
	trhttp "github.com/golovanevvs/confidant/internal/client/transport/http"
	"go.uber.org/zap"
)

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

	// initializing the transport
	trHTTP := trhttp.New("localhost:7541")

	// initializing the service
	sv := service.New(trHTTP)

	av := appview.New(sv, lg)

	// running the app view
	if err := av.Run(); err != nil {
		lg.Fatal(err)
	}
}
