package handler

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"shop_srvs/good_srv/global"
	"shop_srvs/good_srv/model"
	"shop_srvs/good_srv/proto"
)

type GoodServer struct {
	proto.UnimplementedGoodsServer
}

// ModelToResponse 转化为Grpc的Res
func ModelToResponse(goods model.Good) *proto.GoodsInfoResponse {
	return &proto.GoodsInfoResponse{
		Id:              goods.ID,
		CategoryId:      goods.CategoryID,
		Name:            goods.Name,
		GoodsSn:         goods.GoodSn,
		ClickNum:        goods.ClickNum,
		SoldNum:         goods.SoldNum,
		FavNum:          goods.FavNum,
		MarketPrice:     goods.MarketPrice,
		ShopPrice:       goods.ShopPrice,
		GoodsBrief:      goods.GoodBrief,
		ShipFree:        goods.ShipFree,
		GoodsFrontImage: goods.GoodFrontImage,
		IsNew:           goods.IsNew,
		IsHot:           goods.IsHot,
		OnSale:          goods.OnSale,
		DescImages:      goods.DescImages,
		Images:          goods.Images,
		Category: &proto.CategoryBriefInfoResponse{
			Id:   goods.Category.ID,
			Name: goods.Category.Name,
		},
		Brand: &proto.BrandInfoResponse{
			Id:   goods.Brand.ID,
			Name: goods.Brand.Name,
			Logo: goods.Brand.Logo,
		},
	}
}

// GoodsList 过滤 获取商品列表
func (g *GoodServer) GoodsList(ctx context.Context, req *proto.GoodsFilterRequest) (*proto.GoodsListResponse, error) {
	var goods []model.Good
	localDB := global.DB.Model(model.Good{})
	if req.KeyWords != "" {
		// search for like
		localDB = localDB.Where("name like ?", "%s"+req.KeyWords+"%s")
	}
	if req.IsHot {
		localDB = localDB.Where(model.Good{IsHot: true})
	}
	if req.IsNew {
		localDB = localDB.Where(model.Good{IsNew: true})
	}
	if req.PriceMin > 0 {
		localDB = localDB.Where("shop_price>=?", req.PriceMin)
	}
	if req.PriceMax > 0 {
		localDB = localDB.Where("shop_price<=?", req.PriceMax)
	}
	if req.Brand > 0 {
		localDB = localDB.Where("brand_id=?", req.Brand)
	}
	// 手动拼凑sql （子查询）
	var subQuery string
	if req.TopCategory > 0 {
		var category model.Category
		if result := global.DB.First(&category, req.TopCategory); result.RowsAffected == 0 {
			return nil, status.Errorf(codes.NotFound, "商品分类不存在")
		}
		if category.Level == 1 {
			subQuery = fmt.Sprintf("select id from category where parent_category_id in (select id from category WHERE parent_category_id=%d)", req.TopCategory)
		} else if category.Level == 2 {
			subQuery = fmt.Sprintf("select id from category WHERE parent_category_id=%d", req.TopCategory)
		} else if category.Level == 3 {
			subQuery = fmt.Sprintf("select id from category WHERE id=%d", req.TopCategory)
		}
		localDB = localDB.Where(fmt.Sprintf("catgeory_id in (%s)", subQuery))
	}
	var total int64
	localDB.Count(&total)
	result := localDB.Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Find(&goods)
	if result.Error != nil {
		return nil, result.Error
	}
	var goodsInfoRes []*proto.GoodsInfoResponse
	for _, good := range goods {
		goodsInfoRes = append(goodsInfoRes, ModelToResponse(good))
	}
	return &proto.GoodsListResponse{
		Total: int32(total),
		Data:  goodsInfoRes,
	}, nil
}

// BatchGetGoods 批量获取商品列表
func (g *GoodServer) BatchGetGoods(ctx context.Context, req *proto.BatchGoodsIdInfo) (*proto.GoodsListResponse, error) {
	var goods []model.Good
	result := global.DB.Where(req.Id).Find(&goods)

	var goodListRes []*proto.GoodsInfoResponse
	for _, good := range goods {
		goodListRes = append(goodListRes, ModelToResponse(good))
	}
	return &proto.GoodsListResponse{
		Total: int32(result.RowsAffected),
		Data:  goodListRes,
	}, nil
}

// CreateGoods 创建商品
func (g *GoodServer) CreateGoods(ctx context.Context, req *proto.CreateGoodsInfo) (*proto.GoodsInfoResponse, error) {
	var category model.Category
	if result := global.DB.First(&category, req.CategoryId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "商品分类不存在")
	}

	var brand model.Brand
	if result := global.DB.First(&brand, req.BrandId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌不存在")
	}
	good := model.Good{
		CategoryID: category.ID,
		Category:   category,
		BrandID:    brand.ID,
		Brand:      brand,
		OnSale:     req.OnSale,
		ShipFree:   req.ShipFree,
		IsNew:      req.IsNew,
		// 是否热卖
		IsHot: req.IsHot,
		// 商品名字
		Name: req.Name,
		// 商品序号
		GoodSn:         req.GoodsSn,
		MarketPrice:    req.MarketPrice,
		ShopPrice:      req.ShopPrice,
		GoodBrief:      req.GoodsBrief,
		Images:         req.Images,
		DescImages:     req.DescImages,
		GoodFrontImage: req.GoodsFrontImage,
	}
	global.DB.Save(&good)
	return ModelToResponse(good), nil
}

// DeleteGoods 删除商品
func (g *GoodServer) DeleteGoods(ctx context.Context, req *proto.DeleteGoodsInfo) (*empty.Empty, error) {
	if result := global.DB.Delete(&model.Good{}, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "没有这个商品")
	}
	return &empty.Empty{}, nil
}

// UpdateGoods 更新商品
func (g *GoodServer) UpdateGoods(ctx context.Context, req *proto.CreateGoodsInfo) (*empty.Empty, error) {
	var good model.Good
	if result := global.DB.First(&good, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "没有这个商品")
	}

	var category model.Category
	if result := global.DB.First(&category, req.CategoryId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "没有这个商品分类")
	}
	var brand model.Brand
	if result := global.DB.First(&brand, req.BrandId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "没有这个品牌")
	}
	createGoodsInfo := &proto.CreateGoodsInfo{
		Name:            req.Name,
		GoodsSn:         req.GoodsSn,
		Stocks:          req.Stocks,
		MarketPrice:     req.MarketPrice,
		ShopPrice:       req.ShopPrice,
		GoodsBrief:      req.GoodsBrief,
		GoodsDesc:       req.GoodsDesc,
		ShipFree:        req.ShipFree,
		Images:          req.Images,
		DescImages:      req.DescImages,
		GoodsFrontImage: req.GoodsFrontImage,
		IsNew:           req.IsNew,
		IsHot:           req.IsHot,
		OnSale:          req.OnSale,
		CategoryId:      req.CategoryId,
		BrandId:         req.BrandId,
	}
	global.DB.Save(&createGoodsInfo)
	return &empty.Empty{}, nil
}

func (g *GoodServer) GetGoodsDetail(ctx context.Context, req *proto.GoodInfoRequest) (*proto.GoodsInfoResponse, error) {
	var good model.Good
	if result := global.DB.Preload("Category").Preload("Brand").First(&good, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}
	return ModelToResponse(good), nil
}
