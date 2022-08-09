package internalhttp

import (
	"context"
	"io"
	"net"
	"net/http"

	"github.com/Galiks/OTUS_2022/hw12_13_14_15_calendar/internal/logger"
)

type Server struct { // TODO
	srv *http.Server
	mux *http.ServeMux
	app Application
}

type Logger interface { // TODO
}

type Application interface { // TODO
}

func NewServer(app Application, host string, port string) *Server {
	mux := http.NewServeMux()
	RegisterEndpoint(mux)
	server := &http.Server{
		Addr:    net.JoinHostPort(host, port),
		Handler: mux,
	}
	return &Server{
		srv: server,
		app: app,
		mux: mux,
	}
}

func (s *Server) Start() error {
	logger.Info("calendar is running...")
	if err := s.srv.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func RegisterEndpoint(mux *http.ServeMux) {
	mux.Handle("/hello", loggingMiddleware(http.HandlerFunc(testHTTP)))
}

func testHTTP(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "HELLO!")
}

func (s *Server) Stop(ctx context.Context) error {
	if err := s.srv.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}

// TODO
