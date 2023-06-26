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

func (g *GoodServer) CategoryBrandList(ctx context.Context, req *proto.CategoryBrandFilterRequest) (*proto.CategoryBrandListResponse, error) {
	var total int64
	global.DB.Model(&model.GoodsCategoryBrand{}).Count(&total)

	var categoryBrand []model.GoodsCategoryBrand
	global.DB.Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Find(&categoryBrand)

	var categoryBrandRes []*proto.CategoryBrandResponse
	for _, categoryBrandItem := range categoryBrand {
		categoryBrandRes = append(categoryBrandRes, &proto.CategoryBrandResponse{
			Id: categoryBrandItem.ID,
			Brand: &proto.BrandInfoResponse{
				Id:   categoryBrandItem.Brand.ID,
				Name: categoryBrandItem.Brand.Name,
				Logo: categoryBrandItem.Brand.Logo,
			},
			Category: &proto.CategoryInfoResponse{
				Id:             categoryBrandItem.Category.ID,
				Name:           categoryBrandItem.Category.Name,
				ParentCategory: categoryBrandItem.Category.ParentCategoryID,
				Level:          categoryBrandItem.Category.Level,
				IsTab:          categoryBrandItem.Category.IsTab,
			},
		})
	}
	return &proto.CategoryBrandListResponse{
		Total: int32(total),
		Data:  categoryBrandRes,
	}, nil
}

func (g *GoodServer) GetCategoryBrandList(ctx context.Context, req *proto.CategoryInfoRequest) (*proto.BrandListResponse, error) {
	var category model.Category
	if result := global.DB.Find(&category, req.Id).First(&category); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "商品分类不存在")
	}
	var categoryBrand []model.GoodsCategoryBrand
	total := int32(1)
	if result := global.DB.Where(&model.GoodsCategoryBrand{CategoryID: req.Id}).Find(&categoryBrand); result.RowsAffected > 0 {
		total = int32(result.RowsAffected)
	}

	var brandInfoRes []*proto.BrandInfoResponse
	for _, item := range categoryBrand {
		brandInfoRes = append(brandInfoRes, &proto.BrandInfoResponse{
			Id:   item.Brand.ID,
			Name: item.Brand.Name,
			Logo: item.Brand.Logo,
		})
	}
	return &proto.BrandListResponse{
		Total: total,
		Data:  brandInfoRes,
	}, nil
}

func (g *GoodServer) CreateCategoryBrand(ctx context.Context, req *proto.CategoryBrandRequest) (*proto.CategoryBrandResponse, error) {
	var category model.Category
	if result := global.DB.First(&category, req.CategoryId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}
	var brand model.Brand
	if result := global.DB.First(&brand, req.BrandId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "品牌不存在")
	}
	categoryBrand := model.GoodsCategoryBrand{
		CategoryID: req.CategoryId,
		Category:   category,
		BrandID:    req.BrandId,
		Brand:      brand,
	}
	global.DB.Save(categoryBrand)
	return &proto.CategoryBrandResponse{
		Id: req.Id,
		//Brand:    &proto.BrandInfoResponse{},
		//Category: &proto.CategoryInfoResponse{},
	}, nil
}

func (g *GoodServer) DeleteCategoryBrand(ctx context.Context, req *proto.CategoryBrandRequest) (*empty.Empty, error) {
	if result := global.DB.Delete(model.GoodsCategoryBrand{}, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "品牌分类不存在")
	}
	return &empty.Empty{}, nil
}

func (g *GoodServer) UpdateCategoryBrand(ctx context.Context, req *proto.CategoryBrandRequest) (*empty.Empty, error) {
	var categoryBrand model.GoodsCategoryBrand
	if result := global.DB.First(&categoryBrand); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "品牌分类不存在")
	}

	var category model.Category
	if result := global.DB.First(&category, req.CategoryId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}
	var brand model.Brand
	if result := global.DB.First(&brand, req.BrandId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "品牌不存在")
	}
	categoryBrand = model.GoodsCategoryBrand{
		CategoryID: req.CategoryId,
		BrandID:    req.BrandId,
	}
	global.DB.Save(categoryBrand)
	return &empty.Empty{}, nil
}
