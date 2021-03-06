// Code generated by go-swagger; DO NOT EDIT.

package cleanup

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// CleanupHandlerFunc turns a function with the right signature into a cleanup handler
type CleanupHandlerFunc func(CleanupParams) middleware.Responder

// Handle executing the request and returning a response
func (fn CleanupHandlerFunc) Handle(params CleanupParams) middleware.Responder {
	return fn(params)
}

// CleanupHandler interface for that can handle valid cleanup params
type CleanupHandler interface {
	Handle(CleanupParams) middleware.Responder
}

// NewCleanup creates a new http.Handler for the cleanup operation
func NewCleanup(ctx *middleware.Context, handler CleanupHandler) *Cleanup {
	return &Cleanup{Context: ctx, Handler: handler}
}

/*Cleanup swagger:route POST /cleanup Cleanup cleanup

Endpoint to execute all cleanups necessary after a test run

Endpoint to be called during test teardown.  It will execute all cleanups necessary after a test run


*/
type Cleanup struct {
	Context *middleware.Context
	Handler CleanupHandler
}

func (o *Cleanup) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewCleanupParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
