package restapi

import (
	"crypto/tls"
	"net/http"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"
	graceful "github.com/tylerb/graceful"

	"github.com/Polarishq/ot-try/restapi/operations"
	"github.com/Polarishq/ot-try/restapi/operations/cleanup"
	"github.com/Polarishq/ot-try/restapi/operations/health"
	"github.com/Polarishq/ot-try/backend"
	"github.com/Polarishq/middleware/framework"

	"github.com/go-openapi/swag"
	"github.com/Polarishq/middleware/framework/log"
	"github.com/Polarishq/middleware/handlers"
)

// This file is safe to edit. Once it exists it will not be overwritten
// CmdOptions is used to define command line flags
type CmdOptions struct {
	LogFile string `short:"l" long:"logfile" description:"Specify the log file" default:""`
	// always defaults to false
	Verbose bool `short:"v" long:"verbose" description:"Show verbose debug information"`
	// always defaults to false
	VeryVerbose bool `short:"V" long:"very-verbose" description:"Show verbose debug information including aws"`
	// default auto
	StaticDir string `short:"s" long:"static" description:"The path to the static dirs" default:""`
}

// CmdOptionsValues stores the CmdOptions the server has been run with
var CmdOptionsValues CmdOptions // export for testing

func configureFlags(api *operations.OtTryAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
	api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{
		{
			ShortDescription: "Server Flags",
			LongDescription:  "Server Flags",
			Options:          &CmdOptionsValues,
		},
	}
}

func configureAPI(api *operations.OtTryAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	log.SetDebug(CmdOptionsValues.Verbose || CmdOptionsValues.VeryVerbose)
	if CmdOptionsValues.LogFile != "" {
		log.SetOutput(CmdOptionsValues.LogFile)
	}

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	api.CleanupCleanupHandler = cleanup.CleanupHandlerFunc(func(params cleanup.CleanupParams) middleware.Responder {
		return framework.HandleAPIRequestWithError(backend.Cleanup())
	})
	api.HealthHealthHandler = health.HealthHandlerFunc(func(params health.HealthParams) middleware.Responder {
		return framework.HandleAPIRequestWithError(backend.Health())
	})

	api.ServerShutdown = func() {}

	backend.InitializeAll()

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *graceful.Server, scheme, addr string) {
	log.Infof("Browse swagger UI at: %s://%s/swagger-ui", scheme, addr)
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handlers.NewStandardNovaGlobalMiddleware(CmdOptionsValues.StaticDir, handler)
}
