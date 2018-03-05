// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// APIKey Api key
// swagger:model ApiKey
type APIKey struct {

	// Client id for the API key
	// Required: true
	ClientID *string `json:"client_id"`

	// Client secret for the API key
	// Required: true
	ClientSecret *string `json:"client_secret"`

	// If the api key is enabled
	// Required: true
	Enabled *bool `json:"enabled"`

	// Common name for an API Key
	// Required: true
	Label *string `json:"label"`
}

// Validate validates this Api key
func (m *APIKey) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateClientID(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateClientSecret(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateEnabled(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateLabel(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *APIKey) validateClientID(formats strfmt.Registry) error {

	if err := validate.Required("client_id", "body", m.ClientID); err != nil {
		return err
	}

	return nil
}

func (m *APIKey) validateClientSecret(formats strfmt.Registry) error {

	if err := validate.Required("client_secret", "body", m.ClientSecret); err != nil {
		return err
	}

	return nil
}

func (m *APIKey) validateEnabled(formats strfmt.Registry) error {

	if err := validate.Required("enabled", "body", m.Enabled); err != nil {
		return err
	}

	return nil
}

func (m *APIKey) validateLabel(formats strfmt.Registry) error {

	if err := validate.Required("label", "body", m.Label); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *APIKey) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *APIKey) UnmarshalBinary(b []byte) error {
	var res APIKey
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
