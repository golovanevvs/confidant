package handler

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/golovanevvs/confidant/internal/server/repository"
	"github.com/golovanevvs/confidant/internal/server/repository/mocks"
	"github.com/golovanevvs/confidant/internal/server/service"
	"go.uber.org/zap"
)

func TestaccountRegisterPost(t *testing.T) {
	//! preparatory operations
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

	//! using a DB mock
	// creating a gomock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// creating a mock-object
	manageRpMock := mocks.NewMockIManageRepository(ctrl)
	accountRpMock := mocks.NewMockIAccountRepository(ctrl)

	// initializing the repository
	rp := repository.New(manageRpMock, accountRpMock)

	// initializing the service
	sv := service.New(rp)

	//initializing the handler
	hd := New(sv, lg)

	// initializing the test server
	ts := httptest.NewServer(hd.InitRoutes())
	defer ts.Close()

	//! setting input and output parameters
	// request
	type request struct {
		targetRequest string // endpoint
		method        string // http method
		body          []byte // request body
	}

	// expected response
	type expectedResponse struct {
		httpStatus int    // http status
		resp       []byte // server response
	}

	// test parameters
	tests := []struct {
		name             string           // test name
		request          request          // request
		expectedResponse expectedResponse // expected response
	}{
		{
			name: "positive test",
			request: request{
				targetRequest: "/api/register",
				method:        "POST",
				body:          []byte(`{"email":"test@test.ru","password":"Gt56.@rf"}`),
			},
			expectedResponse: expectedResponse{
				httpStatus: 200,
				resp:       []byte(`{"email":"ok"}`),
			},
		},
	}

}
