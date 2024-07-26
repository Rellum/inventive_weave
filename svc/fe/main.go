package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/Rellum/inventive_weave/pkg/metrics"
	pb "github.com/Rellum/inventive_weave/svc/creators/creatorspb"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/xerrors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func run(ctx context.Context, args []string) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	// Parse flags
	fs := flag.NewFlagSet("", flag.PanicOnError)
	webPort := fs.Int("web_port", 9072, "The web server port")
	metricsPort := fs.Int("metrics_port", 9073, "The metrics server port")
	creatorSvc := fs.String("creators_svc", "localhost:9070", "The address of the creator service")
	err := fs.Parse(args[1:])
	if err != nil {
		return xerrors.Errorf("FlagSet.Parse: %w", err)
	}

	// Setup metrics.
	reg := prometheus.NewRegistry()
	shutdownMetrics := metrics.Serve(ctx, fmt.Sprintf(":%d", *metricsPort), reg)

	// Setup creator service client
	slog.Info("connecting to creators service", slog.String("address", *creatorSvc))
	conn, err := grpc.NewClient(*creatorSvc, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return xerrors.Errorf("grpc.NewClient: %w", err)
	}
	creatorClient := pb.NewCreatorsClient(conn)

	// Setup web server
	shutdownWeb := serve(ctx, fmt.Sprintf(":%d", *webPort), reg, creatorClient)

	// Handle shutdown
	<-ctx.Done()
	slog.Info("attempting to shut down gracefully")

	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	err = shutdownWeb(shutdownCtx)
	if err != nil {
		return err
	}

	conn.Close()

	err = shutdownMetrics(shutdownCtx)
	if err != nil {
		return err
	}

	slog.Info("shut down gracefully")
	return nil
}

func serve(ctx context.Context, address string, reg prometheus.Registerer, creatorClient pb.CreatorsClient) func(ctx context.Context) error {
	mux := http.NewServeMux()
	var handler http.Handler = mux
	addRoutes(mux, creatorClient)

	requestsTotal := promauto.With(reg).NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Tracks the number of HTTP requests.",
		}, []string{"method", "code"},
	)
	requestDuration := promauto.With(reg).NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Tracks the latencies for HTTP requests.",
			Buckets: []float64{0.001, 0.01, 0.1, 0.3, 0.6, 1, 3, 6, 9, 20, 30},
		},
		[]string{"method", "code"},
	)
	requestSize := promauto.With(reg).NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "http_request_size_bytes",
			Help: "Tracks the size of HTTP requests.",
		},
		[]string{"method", "code"},
	)
	responseSize := promauto.With(reg).NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "http_response_size_bytes",
			Help: "Tracks the size of HTTP responses.",
		},
		[]string{"method", "code"},
	)

	handler = promhttp.InstrumentHandlerResponseSize(responseSize, handler)
	handler = promhttp.InstrumentHandlerRequestSize(requestSize, handler)
	handler = promhttp.InstrumentHandlerDuration(requestDuration, handler)
	handler = promhttp.InstrumentHandlerCounter(requestsTotal, handler)
	handler = loggingMiddleware(handler)
	httpServer := &http.Server{
		Addr:    address,
		Handler: handler,
	}
	go func() {
		slog.Info("serving frontend", slog.String("addr", httpServer.Addr))
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

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("received request", slog.String("method", r.Method), slog.String("path", r.URL.Path))
		next.ServeHTTP(w, r)
	})
}

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Args); err != nil {
		slog.ErrorContext(ctx, "run error", slog.Any("error", err))
		os.Exit(1)
	}
}
