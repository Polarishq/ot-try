package handlers

import (
	"net/http"
	"runtime/debug"

	"github.com/Polarishq/middleware/framework/log"
	"github.com/go-openapi/errors"
)

// PanicHandler provides an HTTP middleware handler that recovers from panics and returns a standard error response
type PanicHandler struct {
	handler http.Handler
}

// NewPanicHandler creates an HTTP middleware handler that recovers from panics and returns a standard error response
func NewPanicHandler(handler http.Handler) *PanicHandler {
	return &PanicHandler{handler: handler}
}

func (p *PanicHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("action=panic url=%v r=%v, error=%v", r.RequestURI, r, err)
			debug.PrintStack()
			errors.ServeError(rw, nil, err.(error))
		}
	}()

	p.handler.ServeHTTP(rw, r)
}
