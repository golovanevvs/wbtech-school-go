package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golovanevvs/wbtech-school-go/L2/L2.18/cal/internal/config"
	"github.com/rs/zerolog"
)

type Server struct {
	config *config.Server
	router *gin.Engine
	logger *zerolog.Logger
}

func New(config *config.Server, router *gin.Engine, logger *zerolog.Logger) *Server {
	return &Server{
		config: config,
		router: router,
		logger: logger,
	}
}

func (s *Server) RunServerWithGracefulShutdown(DelayBeforeClosing time.Duration) {
	srv := &http.Server{
		Addr:    s.config.Addr,
		Handler: s.router,
	}

	go func() {
		s.logger.Info().Str("address", s.config.Addr).Msg("starting HTTP server")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error().Err(err).Msg("failed to start HTTP server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	s.logger.Info().Msg("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), s.config.Timeout)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		s.logger.Error().Err(err).Msg("server forced to shutdown")
	}

	s.logger.Info().Msg("server exited gracefully")

	time.Sleep(DelayBeforeClosing)
}
