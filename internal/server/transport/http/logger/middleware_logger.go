package logger

import (
	"bytes"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// structure for storing information about the response
type responseData struct {
	status      int
	size        int
	contentType string
	body        *bytes.Buffer
}

// structure with http.ResponseWriter и responseData
type loggingResponseWriter struct {
	// embedding the original http.ResponseWriter
	http.ResponseWriter
	responseData *responseData
}

// redefining the Write and WriteHeader methods of the http interface.ResponseWriter
func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	// we record the response using the original http.ResponseWriter, we get the size
	size, err := r.ResponseWriter.Write(b)
	// capturing the size
	r.responseData.size += size
	// capturing the response body
	r.responseData.body.Write(b)

	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	// we record the status code using the original http.ResponseWriter
	r.ResponseWriter.WriteHeader(statusCode)
	// capturing the status code
	r.responseData.status = statusCode
	// capturing the type of content
	r.responseData.contentType = r.Header().Get("Content-Type")
}

// WithLogging - middleware, a wrapper function that wraps http.Handler
// adds an additional code and returns a new http.Handler
func WithLogging(lg *zap.SugaredLogger) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/api/status" {
				h.ServeHTTP(w, r)
				return
			}
			// to determine the request processing time
			start := time.Now()

			// creating an instance of the responseData structure
			responseData := &responseData{
				status:      0,
				size:        0,
				contentType: "",
				body:        bytes.NewBuffer(nil),
			}

			// creating an instance of the loggingResponseWriter structure
			lw := loggingResponseWriter{
				// embedding the original http.ResponseWriter
				ResponseWriter: w,
				responseData:   responseData,
			}

			// recording the request data
			// endpoint
			reqURI := r.RequestURI
			// request method
			reqMethod := r.Method
			// type of content
			reqContentType := r.Header.Get("Content-Type")

			// обслуживание оригинального запроса c внедрённой реализацией http.ResponseWriter
			h.ServeHTTP(&lw, r)

			// recording the response data
			// status
			resStatus := responseData.status
			// type of content
			resContentType := responseData.contentType
			// size
			resSize := responseData.size
			// body
			resBody := responseData.body.String()

			// request processing duration
			duration := time.Since(start)

			// sending information to the logger
			lg.Debugf("---------------------------------------------------------------")
			lg.Debugf("Request method: %v", reqMethod)
			lg.Debugf("Request URI: %v", reqURI)
			lg.Debugf("Request Content-Type: %v", reqContentType)
			lg.Debugf("Response status: %v", resStatus)
			lg.Debugf("Response Content-Type: %v", resContentType)
			lg.Debugf("Response size: %v", resSize)
			lg.Debugf("Response body: %v", resBody)
			lg.Debugf("Duration: %v", duration)
		})
	}
}
