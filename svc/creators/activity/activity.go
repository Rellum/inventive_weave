package activity

import (
	pb "github.com/Rellum/inventive_weave/svc/creators/creatorspb"
	"slices"
)

func MostActive(creators []*pb.Creator, products []*pb.Product) []*pb.CreatorStats {
	agg := make(map[string]*pb.CreatorStats)
	for i := range creators {
		agg[creators[i].Id] = &pb.CreatorStats{
			Creator: creators[i],
		}
	}
	for i := range products {
		creatorStat := agg[products[i].CreatorId]
		if creatorStat == nil {
			creatorStat = &pb.CreatorStats{
				Creator: &pb.Creator{
					Id:    products[i].CreatorId,
					Email: "",
				},
			}
		}
		creatorStat.ProductCount++
		if creatorStat.MostRecentCreateTime.AsTime().Before(products[i].CreateTime.AsTime()) {
			creatorStat.MostRecentCreateTime = products[i].CreateTime
		}
		agg[products[i].CreatorId] = creatorStat
	}

	var res []*pb.CreatorStats
	for i := range agg {
		res = append(res, agg[i])
	}
	slices.SortFunc(res, compare)

	return res
}

func compare(a, b *pb.CreatorStats) int {
	if a.ProductCount < b.ProductCount {
		return 1
	}
	if a.ProductCount > b.ProductCount {
		return -1
	}

	if a.MostRecentCreateTime.AsTime().Before(b.MostRecentCreateTime.AsTime()) {
		return 1
	}
	if a.MostRecentCreateTime.AsTime().After(b.MostRecentCreateTime.AsTime()) {
		return -1
	}

	return 0
}
