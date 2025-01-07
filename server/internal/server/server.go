package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/rs/zerolog"
)

type Config struct {
	Port            int           `mapstructure:"port"`
	IdleTimeout     time.Duration `mapstructure:"idle-timeout"`
	ReadTimeout     time.Duration `mapstructure:"read-timeout"`
	WriteTimeout    time.Duration `mapstructure:"write-timeout"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown-timeout"`
}

type Server struct {
	Config    *Config
	Logger    *zerolog.Logger
	WaitGroup sync.WaitGroup
}

func New(cfg *Config, logger *zerolog.Logger) *Server {
	return &Server{Config: cfg, Logger: logger}
}

func (s *Server) Start(handler http.Handler) error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", s.Config.Port),
		Handler:      handler,
		ErrorLog:     log.New(s.Logger, "", 0),
		IdleTimeout:  s.Config.IdleTimeout,
		ReadTimeout:  s.Config.ReadTimeout,
		WriteTimeout: s.Config.WriteTimeout,
	}

	shutdownError := make(chan error)
	go s.finishBackgroudTasks(shutdownError, srv)

	s.Logger.Info().Str("addr", srv.Addr).Msg("starting server")

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	s.Logger.Info().Str("addr", srv.Addr).Msg("stopped server")

	return nil
}

func (s *Server) finishBackgroudTasks(shutdownError chan error, srv *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit

	s.Logger.Info().Str("signal", sig.String()).Msg("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), s.Config.ShutdownTimeout)
	defer cancel()

	err := srv.Shutdown(ctx)
	if err != nil {
		shutdownError <- err
	}

	s.Logger.Info().Str("addr", srv.Addr).Msg("completing background tasks")

	s.WaitGroup.Wait()
	shutdownError <- nil
}
