package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Server struct {
	router *gin.Engine
	logger *zerolog.Logger
}

func New(router *gin.Engine, logger *zerolog.Logger) *Server {
	return &Server{
		router: router,
		logger: logger,
	}
}

func (s *Server) RunServerWithGracefulShutdown(addr string, timeout time.Duration) {
	srv := &http.Server{
		Addr:    addr,
		Handler: s.router,
	}

	go func() {
		s.logger.Info().Str("address", addr).Msg("starting HTTP server")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error().Err(err).Msg("failed to start HTTP server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	s.logger.Info().Msg("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		s.logger.Error().Err(err).Msg("server forced to shutdown")
	}

	s.logger.Info().Msg("server exited gracefully")
}
