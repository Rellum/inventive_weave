package activity_test

import (
	"github.com/Rellum/inventive_weave/svc/creators/activity"
	pb "github.com/Rellum/inventive_weave/svc/creators/creatorspb"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"
)

func Test_MostActive_happyPath(t *testing.T) {
	t0 := timestamppb.Now()
	t1 := timestamppb.New(t0.AsTime().Add(time.Second * 3))

	got := activity.MostActive(
		[]*pb.Creator{
			{Id: "creator 1", Email: "creator 1 email"},
			{Id: "creator 2", Email: "creator 2 email"},
			{Id: "creator 3", Email: "creator 3 email"},
		},
		[]*pb.Product{
			{Id: "product 1", CreatorId: "creator 1", CreateTime: t0},
			{Id: "product 2", CreatorId: "creator 2", CreateTime: t1},
			{Id: "product 3", CreatorId: "creator 3", CreateTime: t0},
			{Id: "product 4", CreatorId: "creator 3", CreateTime: t0},
		},
	)

	// Then
	want := []*pb.CreatorStats{
		{Creator: &pb.Creator{Id: "creator 3", Email: "creator 3 email"}, ProductCount: 2, MostRecentCreateTime: t0},
		{Creator: &pb.Creator{Id: "creator 2", Email: "creator 2 email"}, ProductCount: 1, MostRecentCreateTime: t1},
		{Creator: &pb.Creator{Id: "creator 1", Email: "creator 1 email"}, ProductCount: 1, MostRecentCreateTime: t0},
	}

	assert.Equal(t, want, got)
}

func Test_MostActive_creatorMissingForProduct(t *testing.T) {
	t0 := timestamppb.Now()
	t1 := timestamppb.New(t0.AsTime().Add(time.Second * 3))

	got := activity.MostActive(
		[]*pb.Creator{
			{Id: "creator 1", Email: "creator 1 email"},
			{Id: "creator 2", Email: "creator 2 email"},
		},
		[]*pb.Product{
			{Id: "product 1", CreatorId: "creator 1", CreateTime: t0},
			{Id: "product 2", CreatorId: "creator 2", CreateTime: t1},
			{Id: "product 3", CreatorId: "creator 3", CreateTime: t0},
			{Id: "product 4", CreatorId: "creator 3", CreateTime: t0},
		},
	)

	// Then
	want := []*pb.CreatorStats{
		{Creator: &pb.Creator{Id: "creator 3", Email: ""}, ProductCount: 2, MostRecentCreateTime: t0},
		{Creator: &pb.Creator{Id: "creator 2", Email: "creator 2 email"}, ProductCount: 1, MostRecentCreateTime: t1},
		{Creator: &pb.Creator{Id: "creator 1", Email: "creator 1 email"}, ProductCount: 1, MostRecentCreateTime: t0},
	}

	assert.Equal(t, want, got)
}

func Test_MostActive_fewerThan3Creators(t *testing.T) {
	t0 := timestamppb.Now()
	t1 := timestamppb.New(t0.AsTime().Add(time.Second * 3))

	got := activity.MostActive(
		[]*pb.Creator{
			{Id: "creator 1", Email: "creator 1 email"},
			{Id: "creator 2", Email: "creator 2 email"},
		},
		[]*pb.Product{
			{Id: "product 1", CreatorId: "creator 1", CreateTime: t0},
			{Id: "product 2", CreatorId: "creator 2", CreateTime: t1},
		},
	)

	// Then
	want := []*pb.CreatorStats{
		{Creator: &pb.Creator{Id: "creator 2", Email: "creator 2 email"}, ProductCount: 1, MostRecentCreateTime: t1},
		{Creator: &pb.Creator{Id: "creator 1", Email: "creator 1 email"}, ProductCount: 1, MostRecentCreateTime: t0},
	}

	assert.Equal(t, want, got)
}

func Test_MostActive_nilProducts(t *testing.T) {
	got := activity.MostActive(
		[]*pb.Creator{
			{Id: "creator 1", Email: "creator 1 email"},
			{Id: "creator 2", Email: "creator 2 email"},
			{Id: "creator 3", Email: "creator 3 email"},
		},
		nil,
	)

	// Then
	want := []*pb.CreatorStats{
		{Creator: &pb.Creator{Id: "creator 1", Email: "creator 1 email"}, ProductCount: 0, MostRecentCreateTime: nil},
		{Creator: &pb.Creator{Id: "creator 2", Email: "creator 2 email"}, ProductCount: 0, MostRecentCreateTime: nil},
		{Creator: &pb.Creator{Id: "creator 3", Email: "creator 3 email"}, ProductCount: 0, MostRecentCreateTime: nil},
	}

	assert.Equal(t, want, got)
}
