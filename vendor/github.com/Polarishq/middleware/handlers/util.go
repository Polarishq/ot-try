package handlers

import (
	"net/http"
	"os"

	"github.com/Polarishq/middleware/framework/log"
	"github.com/splunknova/log4nova"
)

// NewStandardNovaGlobalMiddleware constructs a default set of middleware handlers for global use.
func NewStandardNovaGlobalMiddleware(staticDir string, handler http.Handler) http.Handler {
	return NewStandardNovaGlobalMiddlewareWithLoggingOptions(staticDir, handler,
		true, true, false, nil)
}

// NewStandardNovaGlobalMiddlewareWithLoggingOptions constructs a default set of middleware handlers for global use with options to log
// request and response body.
func NewStandardNovaGlobalMiddlewareWithLoggingOptions(
	staticDir string, handler http.Handler,
	logRequestBody bool, logResponseBody bool, novaLogger bool, customRateLimitMapper RateLimitMapperFunc) (stack http.Handler) {

	//Construct the default stack
	stack = NewStatsdHandler(NewPanicHandler(NewSwaggerUIHandler(
		staticDir, NewRateLimitHandler(NewLoggingHandlerWithBody(handler, logRequestBody, logResponseBody),
			customRateLimitMapper))))

	if novaLogger {
		// Configure the nova logger
		clientID := os.Getenv("NOVA_CLIENT_ID")
		clientSecret := os.Getenv("NOVA_CLIENT_SECRET")
		novaHost := os.Getenv("NOVA_HOST")

		novaLogger, err := log4nova.NewNovaLoggerWithHost(clientID, clientSecret, novaHost)
		if err != nil {
			log.WithField("msg", err).Warn("Nova logger failed initialization")
		}

		if novaLogger != nil {
			log.WithFields(log.Fields{
				"clientID": clientID,
				"novaHost": novaHost,
			}).Info("log4nova configured successfully")
			stack = log4nova.NewNovaHandler(novaLogger, stack)
		}
	}

	return
}
