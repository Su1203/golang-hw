package internalhttp

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Server struct {
	server *http.Server
	logger Logger
	app    Application
}

type Logger interface {
	Info(msg string)
	Error(msg string)
	Warn(msg string)
	Debug(msg string)
}

type Application interface {
}

func NewServer(logger Logger, app Application, host, port string) *Server {
	s := &Server{
		logger: logger,
		app:    app,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", s.helloHandler)
	mux.HandleFunc("/hello", s.helloHandler)

	s.server = &http.Server{
		Addr:         host + ":" + port,
		Handler:      s.loggingMiddleware(mux),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return s
}

func (s *Server) helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Hello, Calendar!")
}

func (s *Server) Start(ctx context.Context) error {
	s.logger.Info(fmt.Sprintf("Starting HTTP server on %s", s.server.Addr))

	errChan := make(chan error, 1)
	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	select {
	case err := <-errChan:
		return fmt.Errorf("server failed to start: %w", err)
	case <-ctx.Done():
		return nil
	}
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("Stopping HTTP server...")
	return s.server.Shutdown(ctx)
}
