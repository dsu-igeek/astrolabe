// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"github.com/vmware-tanzu/astrolabe/gen/restapi/operations"
)

//go:generate swagger generate server --target ../../gen --name Astrolabe --spec ../../openapi/astrolabe_api.yaml --exclude-main

func configureFlags(api *operations.AstrolabeAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.AstrolabeAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	if api.GetAstrolabeTasksNexusTaskNexusIDHandler == nil {
		api.GetAstrolabeTasksNexusTaskNexusIDHandler = operations.GetAstrolabeTasksNexusTaskNexusIDHandlerFunc(func(params operations.GetAstrolabeTasksNexusTaskNexusIDParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.GetAstrolabeTasksNexusTaskNexusID has not yet been implemented")
		})
	}
	if api.PostAstrolabeTasksNexusHandler == nil {
		api.PostAstrolabeTasksNexusHandler = operations.PostAstrolabeTasksNexusHandlerFunc(func(params operations.PostAstrolabeTasksNexusParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.PostAstrolabeTasksNexus has not yet been implemented")
		})
	}
	if api.CopyProtectedEntityHandler == nil {
		api.CopyProtectedEntityHandler = operations.CopyProtectedEntityHandlerFunc(func(params operations.CopyProtectedEntityParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.CopyProtectedEntity has not yet been implemented")
		})
	}
	if api.CreateSnapshotHandler == nil {
		api.CreateSnapshotHandler = operations.CreateSnapshotHandlerFunc(func(params operations.CreateSnapshotParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.CreateSnapshot has not yet been implemented")
		})
	}
	if api.DeleteProtectedEntityHandler == nil {
		api.DeleteProtectedEntityHandler = operations.DeleteProtectedEntityHandlerFunc(func(params operations.DeleteProtectedEntityParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.DeleteProtectedEntity has not yet been implemented")
		})
	}
	if api.GetProtectedEntityInfoHandler == nil {
		api.GetProtectedEntityInfoHandler = operations.GetProtectedEntityInfoHandlerFunc(func(params operations.GetProtectedEntityInfoParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.GetProtectedEntityInfo has not yet been implemented")
		})
	}
	if api.GetTaskInfoHandler == nil {
		api.GetTaskInfoHandler = operations.GetTaskInfoHandlerFunc(func(params operations.GetTaskInfoParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.GetTaskInfo has not yet been implemented")
		})
	}
	if api.ListProtectedEntitiesHandler == nil {
		api.ListProtectedEntitiesHandler = operations.ListProtectedEntitiesHandlerFunc(func(params operations.ListProtectedEntitiesParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.ListProtectedEntities has not yet been implemented")
		})
	}
	if api.ListServicesHandler == nil {
		api.ListServicesHandler = operations.ListServicesHandlerFunc(func(params operations.ListServicesParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.ListServices has not yet been implemented")
		})
	}
	if api.ListSnapshotsHandler == nil {
		api.ListSnapshotsHandler = operations.ListSnapshotsHandlerFunc(func(params operations.ListSnapshotsParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.ListSnapshots has not yet been implemented")
		})
	}
	if api.ListTaskNexusHandler == nil {
		api.ListTaskNexusHandler = operations.ListTaskNexusHandlerFunc(func(params operations.ListTaskNexusParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.ListTaskNexus has not yet been implemented")
		})
	}
	if api.ListTasksHandler == nil {
		api.ListTasksHandler = operations.ListTasksHandlerFunc(func(params operations.ListTasksParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.ListTasks has not yet been implemented")
		})
	}

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

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
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
