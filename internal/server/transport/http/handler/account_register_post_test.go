package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/golovanevvs/confidant/internal/customerrors"
	"github.com/golovanevvs/confidant/internal/server/repository"
	"github.com/golovanevvs/confidant/internal/server/repository/mocks"
	"github.com/golovanevvs/confidant/internal/server/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestAccountRegisterPost(t *testing.T) {
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
		contentType   string // content-type
		body          []byte // request body
	}

	// expected response
	type expectedResponse struct {
		httpStatus int    // http status
		bodyErr    string // body if error
	}

	// test parameters
	tests := []struct {
		name             string           // test name
		request          request          // request
		setupMock        func()           // setup mock
		expectedResponse expectedResponse // expected response
	}{
		{
			name: "successfull registration",
			request: request{
				targetRequest: "/api/register",
				method:        "POST",
				contentType:   "application/json",
				body:          []byte(`{"email":"test@test.ru","password":"Gt56.@rf"}`),
			},
			setupMock: func() {
				accountRpMock.EXPECT().
					SaveAccount(
						gomock.Any(),
						gomock.Any(),
					).Return(1, nil)
			},
			expectedResponse: expectedResponse{
				httpStatus: http.StatusOK,
				bodyErr:    "",
			},
		},
		{
			name: "invalid content-type",
			request: request{
				targetRequest: "/api/register",
				method:        "POST",
				contentType:   "text/plain",
				body:          []byte(`{"email":"test@test.ru","password":"Gt56.@rf"}`),
			},
			expectedResponse: expectedResponse{
				httpStatus: http.StatusBadRequest,
				bodyErr:    customerrors.ErrContentType400.Error(),
			},
			setupMock: func() {
				accountRpMock.EXPECT().
					SaveAccount(
						gomock.Any(),
						gomock.Any(),
					).Times(0)
			},
		},
	}

	//! runnig the tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// setuping mock
			test.setupMock()

			// creating a request
			request, err := http.NewRequest(test.request.method, ts.URL+test.request.targetRequest, bytes.NewBuffer(test.request.body))
			require.NoError(t, err)
			request.Header.Set("Content-Type", test.request.contentType)

			// sending a request to the test server and receiving a response
			response, err := ts.Client().Do(request)
			require.NoError(t, err)
			defer response.Body.Close()

			// creating the test server response body
			responseBody, err := io.ReadAll(response.Body)
			require.NoError(t, err)

			require.Equal(t, test.expectedResponse.httpStatus, response.StatusCode)

			switch {
			case strings.Contains(request.Header.Get("Content-Type"), "application/json"):
				var m map[string]any
				err = json.Unmarshal(responseBody, &m)
				require.NoError(t, err)

				assert.Equal(t, "test@test.ru", m["email"])
				assert.Equal(t, "1", m["accountID"])
			case strings.Contains(request.Header.Get("Content-Type"), "text/plain"):
				require.Contains(t, string(responseBody), test.expectedResponse.bodyErr)
			default:
				fmt.Println("доделать")
			}

		})
	}

}
