package server

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
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

func Start(ctx context.Context, port int, frontendFS fs.FS) error {
	router, err := newRouter(frontendFS)
	if err != nil {
		return fmt.Errorf("failed to create router: %w", err)
	}

	address := fmt.Sprintf(":%d", port)
	httpServer := &http.Server{
		Addr:              address,
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
		return fmt.Errorf("server stopped: %w", err)
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		err = httpServer.Shutdown(shutdownCtx)
		if err != nil {
			return fmt.Errorf("failed to shut down server: %w", err)
		}
		return nil
	}
}

func newRouter(frontendFS fs.FS) (http.Handler, error) {
	openapiValidationMiddleware, err := createOapiValidationMiddleware()
	if err != nil {
		return nil, fmt.Errorf("failed to create openapi validation middleware: %w", err)
	}

	apiRouter := chi.NewRouter()
	apiRouter.Use(openapiValidationMiddleware)
	apiServer := NewServer()
	apiHandler := HandlerFromMux(apiServer, apiRouter)

	router := chi.NewRouter()
	router.Use(createLoggerMiddleware())
	router.Mount("/api", apiHandler)
	if frontendFS != nil {
		frontendHandler := http.FileServer(http.FS(frontendFS))
		router.Handle("/*", frontendHandler)
	}

	return router, nil
}

func createLoggerMiddleware() func(http.Handler) http.Handler {
	logger := httplog.NewLogger(
		"goilerplate",
		httplog.Options{
			LogLevel:       slog.LevelDebug,
			RequestHeaders: true,
			Concise:        true,
		},
	)

	// This middleware includes the recover middleware
	return httplog.RequestLogger(logger)
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
