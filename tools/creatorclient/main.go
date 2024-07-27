package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/Rellum/inventive_weave/pkg/json"
	pb "github.com/Rellum/inventive_weave/svc/creators/creatorspb"
	"github.com/Rellum/inventive_weave/svc/creators/types"
	"golang.org/x/xerrors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"os"
	"os/signal"
	"time"
)

func run(ctx context.Context, r io.Reader, w io.Writer, args []string) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	// Parse flags
	fs := flag.NewFlagSet("", flag.PanicOnError)
	creatorSvc := fs.String("creators_svc", "localhost:9070", "The address of the creator service")
	err := fs.Parse(args[1:])
	if err != nil {
		return xerrors.Errorf("FlagSet.Parse: %w", err)
	}

	// Setup creator service client
	conn, err := grpc.NewClient(*creatorSvc, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return xerrors.Errorf("grpc.NewClient: %w", err)
	}
	defer conn.Close()
	client := pb.NewCreatorsClient(conn)

	// Parse input
	data, err := json.Decode[types.Data](r)
	if err != nil {
		return err
	}

	res, err := client.MostActiveCreators(context.Background(), pb.ToProto(data))
	if err != nil {
		return xerrors.Errorf("client.MostActiveCreators: %w", err)
	}

	// Print summary
	fmt.Fprintln(w, "The top creators are:")
	for i := 0; i < 3; i++ {
		if len(res.CreatorStats) <= i {
			break
		}
		fmt.Fprintf(w, "%d: %s (products: %d, most recent creation: %v)\n", i+1, res.CreatorStats[i].Creator.Email, res.CreatorStats[i].ProductCount, res.CreatorStats[i].MostRecentCreateTime.AsTime().Local().Format(time.RFC3339Nano))
	}
	fmt.Fprintln(w, "---")

	return nil
}

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Stdin, os.Stdout, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}
