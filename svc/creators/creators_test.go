package main_test

import (
	"context"
	"fmt"
	pb "github.com/Rellum/inventive_weave/svc/creators/creatorspb"
	"github.com/Rellum/inventive_weave/svc/creators/server"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net"
	"testing"
	"time"
)

func Test_MostActiveCreators_nilRequest(t *testing.T) {
	// Given
	client := setup(t)

	// When
	resp, err := client.MostActiveCreators(context.Background(), nil)
	assert.NoError(t, err)

	// Then
	assert.NotNil(t, resp)
}

func Test_MostActiveCreators_happyPath(t *testing.T) {
	// Given
	client := setup(t)
	t0 := timestamppb.Now()
	t1 := timestamppb.New(t0.AsTime().Add(time.Second * 3))

	// When
	resp, err := client.MostActiveCreators(context.Background(), &pb.MostActiveCreatorsReq{
		Creators: []*pb.Creator{
			{Id: "creator 1", Email: "creator 1 email"},
			{Id: "creator 2", Email: "creator 2 email"},
			{Id: "creator 3", Email: "creator 3 email"},
		},
		Products: []*pb.Product{
			{Id: "product 1", CreatorId: "creator 1", CreateTime: t0},
			{Id: "product 2", CreatorId: "creator 2", CreateTime: t1},
			{Id: "product 3", CreatorId: "creator 3", CreateTime: t0},
			{Id: "product 4", CreatorId: "creator 3", CreateTime: t0},
		},
	})
	assert.NoError(t, err)

	// Then
	want := &pb.MostActiveCreatorsRes{
		CreatorStats: []*pb.CreatorStats{
			{
				Creator: &pb.Creator{
					Id:    "creator 3",
					Email: "creator 3 email",
				},
				ProductCount:         2,
				MostRecentCreateTime: t0,
			},
			{
				Creator: &pb.Creator{
					Id:    "creator 2",
					Email: "creator 2 email",
				},
				ProductCount:         1,
				MostRecentCreateTime: t1,
			},
			{
				Creator: &pb.Creator{
					Id:    "creator 1",
					Email: "creator 1 email",
				},
				ProductCount:         1,
				MostRecentCreateTime: t0,
			},
		},
	}

	assert.EqualExportedValues(t, want, resp)
}

func setup(t *testing.T) pb.CreatorsClient {
	t.Helper()

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	t.Cleanup(func() { listener.Close() })
	s := grpc.NewServer()
	t.Cleanup(func() { s.Stop() })
	server.RegisterServer(s)
	go func() {
		err := s.Serve(listener)
		if err != nil {
			t.Errorf("failed to serve: %v", err)
		}
	}()

	port := fmt.Sprintf(":%d", listener.Addr().(*net.TCPAddr).Port)

	conn, err := grpc.NewClient(port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("failed to connect client: %v", err)
	}
	t.Cleanup(func() { conn.Close() })
	return pb.NewCreatorsClient(conn)
}
