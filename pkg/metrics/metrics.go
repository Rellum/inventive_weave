package metrics

import (
	"context"
	"errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/xerrors"
	"log/slog"
	"net/http"
)

func Serve(ctx context.Context, address string, reg *prometheus.Registry) func(ctx context.Context) error {
	httpServer := &http.Server{Addr: address}

	reg.MustRegister(
		collectors.NewBuildInfoCollector(),
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
		httpServer.Handler = mux
		slog.Info("serving metrics", slog.String("addr", httpServer.Addr), slog.String("path", "/metrics"))
		err := httpServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.ErrorContext(ctx, "http.Server.ListenAndServe error", slog.Any("error", xerrors.Errorf("%w", err)))
		}
	}()

	return func(ctx context.Context) error {
		err := httpServer.Shutdown(ctx)
		if err != nil {
			return xerrors.Errorf("%w", err)
		}
		return nil
	}
}
