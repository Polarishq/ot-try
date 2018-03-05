package cleanup

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"
	"time"

	"golang.org/x/net/context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"

	"github.com/Polarishq/middleware/framework/log"
	strfmt "github.com/go-openapi/strfmt"
)

// NewCleanupParams creates a new CleanupParams object
// with the default values initialized.
func NewCleanupParams() *CleanupParams {

	return &CleanupParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewCleanupParamsWithTimeout creates a new CleanupParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewCleanupParamsWithTimeout(timeout time.Duration) *CleanupParams {

	return &CleanupParams{

		timeout: timeout,
	}
}

// NewCleanupParamsWithContext creates a new CleanupParams object
// with the default values initialized, and the ability to set a context for a request
func NewCleanupParamsWithContext(ctx context.Context) *CleanupParams {

	return &CleanupParams{

		Context: ctx,
	}
}

// NewCleanupParamsWithHTTPClient creates a new CleanupParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewCleanupParamsWithHTTPClient(client *http.Client) *CleanupParams {

	return &CleanupParams{
		HTTPClient: client,
	}
}

/*CleanupParams contains all the parameters to send to the API endpoint
for the cleanup operation typically these are written to a http.Request
*/
type CleanupParams struct {
	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the cleanup params
func (o *CleanupParams) WithTimeout(timeout time.Duration) *CleanupParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the cleanup params
func (o *CleanupParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the cleanup params
func (o *CleanupParams) WithContext(ctx context.Context) *CleanupParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the cleanup params
func (o *CleanupParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the cleanup params
func (o *CleanupParams) WithHTTPClient(client *http.Client) *CleanupParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the cleanup params
func (o *CleanupParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WriteToRequest writes these params to a swagger request
func (o *CleanupParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}

	r.SetHeaderParam("X-POLARIS-REQ-ID", log.GetReqID())

	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
