package restapi

import (
	"crypto/tls"
	"net/http"

	"github.com/Sef1995/zipkin-swagger-poc/service/restapi/operations"
	interpose "github.com/carbocation/interpose/middleware"
	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/middleware/http"
	"github.com/openzipkin/zipkin-go/reporter"
	httpreporter "github.com/openzipkin/zipkin-go/reporter/http"
	"github.com/sirupsen/logrus"
	graceful "github.com/tylerb/graceful"
)

var tr *zipkin.Tracer = nil

const (
	zipkinURL   = "http://localhost:9411/api/v2/spans"
	serviceName = "service1"
	host        = "127.0.0.1"
	port        = ":8001"
)

// This file is safe to edit. Once it exists it will not be overwritten

//go:generate swagger generate server --target .. --name service1 --spec ../../../../../../../../../../tmp/swagger.json

func configureFlags(api *operations.Service1API) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.Service1API) http.Handler {

	api.ServeError = errors.ServeError
	api.Logger = logrus.Printf

	api.JSONConsumer = runtime.JSONConsumer()
	api.JSONProducer = runtime.JSONProducer()

	rep := httpreporter.NewReporter(zipkinURL)
	defer rep.Close()
	initializeTracer(rep)

	api.SomeFunctionHandler = operations.SomeFunctionHandlerFunc(func(params operations.SomeFunctionParams) middleware.Responder {
		return middleware.NotImplemented("someFunction is called.")
	})

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

func initializeTracer(r reporter.Reporter) {
	endpoint, err := zipkin.NewEndpoint(serviceName, host+port)
	if err != nil {
		logrus.Fatalf("unable to create local endpoint: %+v\n", err)
	}

	tracer, err := zipkin.NewTracer(r, zipkin.WithLocalEndpoint(endpoint))
	if err != nil {
		logrus.Fatalf("unable to create our tracer: %+v\n", err)
	}

	tr = tracer
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
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	logger := interpose.NegroniLogrus()
	serverMiddleware := zipkinhttp.NewServerMiddleware(
		tr, zipkinhttp.TagResponseSize(true),
	)

	handler = logger(handler)
	handler = serverMiddleware(handler)

	return handler
}
