package health

import (
	"github.com/go-openapi/runtime"
)

// Interface for client methods
type ClientInterface interface {
	Health(params *HealthParams) (*HealthOK, error)

	SetTransport(transport runtime.ClientTransport)
}
