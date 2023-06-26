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

func (g *GoodServer) BrandList(ctx context.Context, req *proto.BrandFilterRequest) (*proto.BrandListResponse, error) {
	var brands []model.Brand
	result := global.DB.Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Find(&brands)
	if result.Error != nil {
		return nil, result.Error
	}
	var count int64
	global.DB.Model(model.Brand{}).Count(&count)
	var brandRes []*proto.BrandInfoResponse
	for _, brand := range brands {
		brandRes = append(brandRes, &proto.BrandInfoResponse{
			Id:   brand.ID,
			Name: brand.Name,
			Logo: brand.Logo,
		})
	}
	return &proto.BrandListResponse{
		Total: int32(count),
		Data:  brandRes,
	}, nil
}

func (g *GoodServer) CreateBrand(ctx context.Context, req *proto.BrandRequest) (*proto.BrandInfoResponse, error) {
	//新建品牌
	if result := global.DB.Where("name=?", req.Name).First(&model.Brand{}); result.RowsAffected == 1 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌已存在")
	}

	brand := model.Brand{
		Name: req.Name,
		Logo: req.Logo,
	}
	result := global.DB.Save(&brand)
	if result.Error != nil {
		return nil, result.Error
	}
	return &proto.BrandInfoResponse{Id: brand.ID, Name: req.Name, Logo: req.Logo}, nil
}

func (g *GoodServer) DeleteBrand(ctx context.Context, req *proto.BrandRequest) (*empty.Empty, error) {
	if result := global.DB.Delete(&model.Brand{}, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌不存在")
	}
	return &empty.Empty{}, nil
}

func (g *GoodServer) UpdateBrand(ctx context.Context, req *proto.BrandRequest) (*empty.Empty, error) {
	brand := model.Brand{}
	if result := global.DB.Delete(&model.Brand{}, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌不存在")
	}
	if req.Name != "" {
		brand.Name = req.Name
	}
	if req.Logo != "" {
		brand.Logo = req.Logo
	}
	global.DB.Save(&brand)
	return &empty.Empty{}, nil
}
