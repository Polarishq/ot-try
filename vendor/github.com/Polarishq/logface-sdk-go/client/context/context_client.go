// Code generated by go-swagger; DO NOT EDIT.

package context

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"
)

// New creates a new context API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) *Client {
	return &Client{transport: transport, formats: formats}
}

/*
Client for context API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

/*
GetContextsID retrieves context

This endpoint allows for retrieving events based on an SPL query.
SPL Reference -- http://docs.splunk.com/Documentation/Splunk/latest/SearchReference/WhatsInThisManual

*/
func (a *Client) GetContextsID(params *GetContextsIDParams) (*GetContextsIDOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetContextsIDParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "GetContextsID",
		Method:             "GET",
		PathPattern:        "/contexts/{id}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &GetContextsIDReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*GetContextsIDOK), nil

}

/*
PostContexts creates a search context
*/
func (a *Client) PostContexts(params *PostContextsParams) (*PostContextsCreated, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewPostContextsParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "PostContexts",
		Method:             "POST",
		PathPattern:        "/contexts",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &PostContextsReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*PostContextsCreated), nil

}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}