package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"

	oapimiddleware "github.com/oapi-codegen/nethttp-middleware"
)

const (
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

func Start(ctx context.Context, port int) error {
	router, err := newRouter()
	if err != nil {
		return err
	}

	httpServer := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           router,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	serverErr := make(chan error, 1)
	go func() {
		serverErr <- httpServer.ListenAndServe()
	}()

	slog.Info("Server listening", "address", httpServer.Addr)

	select {
	case err := <-serverErr:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("failed to shut down server: %w", err)
		}
		return nil
	}
}

func newRouter() (http.Handler, error) {
	openapiValidationMiddleware, err := createOapiValidationMiddleware()
	if err != nil {
		return nil, fmt.Errorf("failed to create openapi validation middleware: %w", err)
	}

	apiRouter := chi.NewRouter()
	apiRouter.Use(openapiValidationMiddleware)
	apiHandler := HandlerFromMux(NewServer(), apiRouter)

	router := chi.NewRouter()
	router.Use(createLoggerMiddleware())
	router.Mount("/api", apiHandler)

	return router, nil
}

func createLoggerMiddleware() func(http.Handler) http.Handler {
	// This middleware includes the recover middleware
	return httplog.RequestLogger(
		httplog.NewLogger(
			"goilerplate",
			httplog.Options{
				LogLevel:       slog.LevelDebug,
				RequestHeaders: true,
				Concise:        true,
			},
		),
	)
}

func createOapiValidationMiddleware() (func(http.Handler) http.Handler, error) {
	spec, err := GetSpec()
	if err != nil {
		return nil, fmt.Errorf("failed to load openapi spec: %w", err)
	}

	// Set server endpoint to the route so request validation matches the mount path.
	// See: https://github.com/oapi-codegen/oapi-codegen/issues/1123
	spec.Servers = openapi3.Servers{{URL: "/api"}}

	return oapimiddleware.OapiRequestValidatorWithOptions(
		spec,
		&oapimiddleware.Options{
			SilenceServersWarning: true,
			Options: openapi3filter.Options{
				AuthenticationFunc: func(_ context.Context, _ *openapi3filter.AuthenticationInput) error {
					return nil
				},
			},
		},
	), nil
}
