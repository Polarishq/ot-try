package cleanup

import (
	"github.com/go-openapi/runtime"
)

// Interface for client methods
type ClientInterface interface {
	Cleanup(params *CleanupParams) (*CleanupNoContent, error)

	SetTransport(transport runtime.ClientTransport)
}
