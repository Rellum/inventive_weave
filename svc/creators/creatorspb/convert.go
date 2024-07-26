package creatorspb

import (
	"github.com/Rellum/inventive_weave/svc/creators/types"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToProto(in types.Data) *MostActiveCreatorsReq {
	var req MostActiveCreatorsReq
	for i := range in.Creators {
		req.Creators = append(req.Creators, &Creator{Id: in.Creators[i].Id, Email: in.Creators[i].Email})
	}
	for i := range in.Products {
		req.Products = append(req.Products, &Product{Id: in.Products[i].Id, CreatorId: in.Products[i].CreatorId, CreateTime: timestamppb.New(in.Products[i].CreateTime)})
	}
	return &req
}
