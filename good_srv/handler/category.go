package handler

import (
	"context"
	"encoding/json"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"shop_srvs/good_srv/global"
	"shop_srvs/good_srv/model"
	"shop_srvs/good_srv/proto"
)

func (g *GoodServer) GetAllCategorysList(ctx context.Context, empty *empty.Empty) (*proto.CategoryListResponse, error) {
	var categorys []model.Category
	result := global.DB.Where(&model.Category{Level: 1}).Preload("SubCategory.SubCategory").Find(&categorys)
	json, err := json.Marshal(&categorys)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "错误的序列化JSON")
	}
	return &proto.CategoryListResponse{
		Total:    int32(result.RowsAffected),
		JsonData: string(json),
		//Data:     categorys,
	}, nil
}

// GetSubCategory 获取子分类
func (g *GoodServer) GetSubCategory(ctx context.Context, req *proto.CategoryListRequest) (*proto.SubCategoryListResponse, error) {
	var category model.Category
	if result := global.DB.First(&category, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}
	categoryInfoResponse := &proto.CategoryInfoResponse{
		Id:             category.ID,
		Name:           category.Name,
		Level:          category.Level,
		IsTab:          category.IsTab,
		ParentCategory: category.ParentCategoryID,
	}
	proloads := "SubCategory"
	if category.Level == 1 {
		proloads = "SubCategory.SubCategory"
	}
	var categorys []model.Category
	result := global.DB.Where(&model.Category{ParentCategoryID: req.Id}).Preload(proloads).Find(&categorys)

	var subCategorysRes []*proto.CategoryInfoResponse
	for _, subCategory := range categorys {
		subCategorysRes = append(subCategorysRes, &proto.CategoryInfoResponse{
			Id:             subCategory.ID,
			Name:           subCategory.Name,
			Level:          subCategory.Level,
			IsTab:          subCategory.IsTab,
			ParentCategory: subCategory.ParentCategoryID,
		})
	}
	return &proto.SubCategoryListResponse{
		Total:        int32(result.RowsAffected),
		Info:         categoryInfoResponse,
		SubCategorys: subCategorysRes,
	}, nil
}

func (g *GoodServer) CreateCategory(ctx context.Context, req *proto.CategoryInfoRequest) (*proto.CategoryInfoResponse, error) {
	category := model.Category{
		Name:  req.Name,
		Level: req.Level,
		IsTab: req.IsTab,
	}
	if req.Level != 1 {
		category.ParentCategoryID = req.ParentCategory
	}

	result := global.DB.Save(&category)
	if result.Error != nil {
		return nil, result.Error
	}
	return &proto.CategoryInfoResponse{Id: req.Id}, nil
}

func (g *GoodServer) DeleteCategory(ctx context.Context, req *proto.DeleteCategoryRequest) (*empty.Empty, error) {
	if result := global.DB.Delete(&model.Category{}, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}
	return &empty.Empty{}, nil
}

func (g *GoodServer) UpdateCategory(ctx context.Context, req *proto.CategoryInfoRequest) (*empty.Empty, error) {
	category := model.Category{}
	if result := global.DB.First(&category); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}
	if req.Name != "" {
		category.Name = req.Name
	}
	if req.ParentCategory != 0 {
		category.ParentCategoryID = req.ParentCategory
	}
	if req.Level != 0 {
		category.Level = req.Level
	}
	if req.IsTab {
		category.IsTab = req.IsTab
	}
	global.DB.Save(&category)
	return &empty.Empty{}, nil
}
