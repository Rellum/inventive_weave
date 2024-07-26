package grpc

import (
	"context"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/xerrors"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

func Serve(ctx context.Context, address string, reg prometheus.Registerer) (*grpc.Server, error) {
	srvMetrics := grpcprom.NewServerMetrics(
		grpcprom.WithServerHandlingTimeHistogram(
			grpcprom.WithHistogramBuckets([]float64{0.001, 0.01, 0.1, 0.3, 0.6, 1, 3, 6, 9, 20, 30}),
		),
	)
	reg.MustRegister(srvMetrics)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, xerrors.Errorf("net.Listen: %w", err)
	}
	s := grpc.NewServer(grpc.ChainUnaryInterceptor(srvMetrics.UnaryServerInterceptor()))

	go func() {
		slog.Info("grpc server starting", slog.String("addr", listener.Addr().String()))
		err := s.Serve(listener)
		if err != nil {
			slog.ErrorContext(ctx, "http.Server.ListenAndServe error", slog.Any("error", xerrors.Errorf("%w", err)))
		}
	}()

	return s, nil
}
