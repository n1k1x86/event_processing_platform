package server

import (
	"context"
	"errors"
	"net/http"
	_ "net/http/pprof"

	"go.uber.org/zap"
)

type PprofServer struct {
	server *http.Server
	logger *zap.Logger
}

func NewPprofServer(addr string, logger *zap.Logger) *PprofServer {
	return &PprofServer{
		server: &http.Server{
			Addr: addr,
		},
		logger: logger,
	}
}

func (p *PprofServer) Start() error {
	p.logger.Info("pprof server starting...")
	err := p.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (p *PprofServer) Stop(ctx context.Context) error {
	err := p.server.Shutdown(ctx)
	if err != nil {
		p.logger.Error("error while shutting down pprof server", zap.Error(err))
		return err
	}
	p.logger.Info("pprof server was shutted down")
	return nil
}
