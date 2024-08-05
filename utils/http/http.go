package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/bnkamalesh/webgo/v6"
	"github.com/gin-gonic/gin"
)

const (
	sucsess                 = "sucsess"
	staticPath              = "static/images/"
	_defaultShutdownTimeout = 3 * time.Second
)

// HTTP struct holds all the dependencies required for starting HTTP server
type HTTP struct {
	server          *http.Server
	cfg             *Config
	notify          chan error
	shutdownTimeout time.Duration
}

// Start starts the HTTP server
// func (h *HTTP) Start() {
// 	webgo.LOGHANDLER.Info("HTTP server, listening on", h.cfg.Host, h.cfg.Port)
// 	h.server.ListenAndServe()
// }

// Config holds all the configuration required to start the HTTP server
type Config struct {
	Host         string
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	DialTimeout  time.Duration
}

// NewService returns an instance of HTTP with all its dependencies set
func NewService(handler *gin.Engine, cfg *Config) (*HTTP, error) {

	httpServer := &http.Server{
		Addr:              fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Handler:           handler,
		ReadTimeout:       cfg.ReadTimeout,
		ReadHeaderTimeout: cfg.ReadTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.ReadTimeout * 2,
	}
	http := &HTTP{
		server:          httpServer,
		cfg:             cfg,
		notify:          make(chan error, 1),
		shutdownTimeout: _defaultShutdownTimeout,
	}
	http.start()
	return http, nil
}
func (s *HTTP) start() {
	go func() {
		webgo.LOGHANDLER.Info("HTTP server, listening on", s.cfg.Host, s.cfg.Port)
		s.notify <- s.server.ListenAndServe()
		close(s.notify)
	}()
}

// Notify -.
func (s *HTTP) Notify() <-chan error {
	return s.notify
}

// Shutdown -.
func (s *HTTP) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.cfg.WriteTimeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}
