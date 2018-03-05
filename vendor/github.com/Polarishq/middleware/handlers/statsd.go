package handlers

import (
	"net/http"
	"time"

	"strconv"

	"os"

	"github.com/Polarishq/middleware/framework/metric"
)

// StatsdHandler implements a middleware handler for publishing metrics
type StatsdHandler struct {
	handler http.Handler
}

// NewStatsdHandler creates a new HTTP handler which publishes statsd metrics
func NewStatsdHandler(handler http.Handler) *StatsdHandler {
	return &StatsdHandler{handler: handler}
}

// ServeHTTP implements the http handler interface
func (s *StatsdHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	beginning, recorder := s.Begin(w)
	s.handler.ServeHTTP(recorder, r)
	s.End(beginning, recorder, r)
}

// Begin starts a recorder
func (s *StatsdHandler) Begin(w http.ResponseWriter) (time.Time, RecorderResponseWriter) {
	start := time.Now()
	writer := NewRecorderResponseWriter(w, 0)
	return start, writer
}

// End closes the recorder with the recorder status
func (s *StatsdHandler) End(start time.Time, recorder RecorderResponseWriter, req *http.Request) {
	s.EndWithStatus(start, recorder.Status(), req)
}

// EndWithStatus closes the recorder with a specific status
func (s *StatsdHandler) EndWithStatus(start time.Time, status int, req *http.Request) {
	end := time.Now()
	responseTime := end.Sub(start)
	s.sendMetrics(status, req, responseTime)
}

// sendMetrics sends http status codes and response time metrics
func (s *StatsdHandler) sendMetrics(status int, req *http.Request, responseTime time.Duration) {
	component := os.Getenv("COMPONENT_NAME")
	if component == "" {
		component = "unknown-component"
	}
	dimensions := metric.Dimensions{metric.ComponentDim: component, "http_url_path": req.URL.Path}
	metric.WithDimensions(dimensions, metric.Dimensions{"http_status_code": strconv.Itoa(status)}).Inc(
		"http_status", 1)
	responseTimeValue := responseTime.Seconds() * 1000
	metric.WithDimensions(dimensions).Timing("http_response_time_ms", int64(responseTimeValue))
}
