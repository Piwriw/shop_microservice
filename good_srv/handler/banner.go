package handler

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"shop_srvs/good_srv/global"
	"shop_srvs/good_srv/model"
	"shop_srvs/good_srv/proto"
)

func (g *GoodServer) BannerList(ctx context.Context, empty *empty.Empty) (*proto.BannerListResponse, error) {
	var banners []model.Banner
	result := global.DB.Find(&banners)
	var bannersRes []*proto.BannerResponse
	for _, banner := range banners {
		bannersRes = append(bannersRes, &proto.BannerResponse{
			Id:    banner.ID,
			Index: banner.Index,
			Image: banner.Image,
			Url:   banner.Url,
		})
	}
	return &proto.BannerListResponse{
		Total: int32(result.RowsAffected),
		Data:  bannersRes,
	}, nil
}

func (g *GoodServer) CreateBanner(ctx context.Context, req *proto.BannerRequest) (*proto.BannerResponse, error) {
	banner := model.Banner{
		Image: req.Image,
		Url:   req.Url,
		Index: req.Index,
	}
	global.DB.Save(&banner)
	return &proto.BannerResponse{
		Id: banner.ID,
	}, nil
}

func (g *GoodServer) DeleteBanner(ctx context.Context, req *proto.BannerRequest) (*empty.Empty, error) {
	if result := global.DB.Delete(&model.Banner{}, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "这个轮播图不存在")
	}
	return &empty.Empty{}, nil
}

func (g *GoodServer) UpdateBanner(ctx context.Context, req *proto.BannerRequest) (*empty.Empty, error) {
	banner := model.Banner{}
	if result := global.DB.Delete(&model.Banner{}, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "轮播图不存在")
	}
	if req.Url != "" {
		banner.Url = req.Url
	}
	if req.Image != "" {
		banner.Image = req.Image
	}
	if req.Index != 0 {
		banner.Index = req.Index
	}
	global.DB.Save(&banner)
	return &empty.Empty{}, nil
}
