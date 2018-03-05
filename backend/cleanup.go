package backend

import (
	"net/http"

	"github.com/Polarishq/middleware/deployment"
)

// Cleanup implements any cleanup tasks that are necessary execute after each test run
func Cleanup() (int, error) {
	cleanup, err := deployment.ShouldCleanup()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if cleanup {
		// Perform ALL cleanup, return an error if necessary
	}
	return http.StatusNoContent, nil
}
