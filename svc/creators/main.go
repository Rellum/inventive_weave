package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/Rellum/inventive_weave/pkg/grpc"
	"github.com/Rellum/inventive_weave/pkg/metrics"
	creatorssrv "github.com/Rellum/inventive_weave/svc/creators/server"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/xerrors"
	"log/slog"
	"os"
	"os/signal"
	"time"
)

func run(ctx context.Context, args []string) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	// Parse flags.
	fs := flag.NewFlagSet("", flag.PanicOnError)
	grpcPort := fs.Int("grpc_port", 9070, "The grpc server port")
	metricsPort := fs.Int("metrics_port", 9071, "The metrics server port")
	err := fs.Parse(args[1:])
	if err != nil {
		return xerrors.Errorf("FlagSet.Parse: %w", err)
	}

	// Setup metrics.
	reg := prometheus.NewRegistry()
	shutdown := metrics.Serve(ctx, fmt.Sprintf(":%d", *metricsPort), reg)

	// Setup gRPC server
	s, err := grpc.Serve(ctx, fmt.Sprintf(":%d", *grpcPort), reg)
	if err != nil {
		return err
	}
	creatorssrv.RegisterServer(s)

	// Handle shutdown
	<-ctx.Done()
	slog.Info("attempting to shut down gracefully")

	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	s.GracefulStop()

	err = shutdown(shutdownCtx)
	if err != nil {
		return err
	}

	slog.Info("shut down gracefully")
	return nil
}

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Args); err != nil {
		slog.ErrorContext(ctx, "run error", slog.Any("error", err))
		os.Exit(1)
	}
}
