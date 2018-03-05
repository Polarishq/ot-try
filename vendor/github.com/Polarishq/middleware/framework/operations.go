package framework

import (
	"net/http"

	"github.com/Polarishq/middleware/framework/log"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
)

type apiOperation struct {
	Response interface{}
}

// WriteResponse writes the response to the given producer
func (a *apiOperation) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {
	var err error
	switch t := a.Response.(type) {
	case error:
		ServeError(rw, nil, t)
	default:
		// success sent a 2xx response
		rw.WriteHeader(int(200))
		err = producer.Produce(rw, a.Response)
	}

	if err != nil {
		log.Errorf("failed to send response to client error=%s", err.Error())
	}
}

// HandleAPIRequestWithError checks if an error occurred and if so returns a standardized error message
func HandleAPIRequestWithError(response interface{}, e error) middleware.Responder {
	op := apiOperation{}
	if e != nil {
		op.Response = e
	} else {
		op.Response = response
	}

	return &op
}

type redirectOperation struct {
	Location   string
	StatusCode int
}

//HandleRedirect sets up the redirect op for a user
func HandleRedirect(url string, statusCode int) middleware.Responder {
	ro := redirectOperation{
		Location:   url,
		StatusCode: statusCode,
	}
	return &ro
}

//WriteResponse writes the url to the location header and sends back the redirect status code
func (ro *redirectOperation) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {
	rw.Header().Set("Location", ro.Location)
	rw.WriteHeader(ro.StatusCode)
	err := producer.Produce(rw, nil)

	if err != nil {
		log.Errorf("failed to send response to client error=%s", err.Error())
	}
}
