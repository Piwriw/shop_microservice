package handler

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math/rand"
	"shop_srvs/order_srv/global"
	"shop_srvs/order_srv/model"
	"shop_srvs/order_srv/proto"
	"time"
)

type OrderServer struct {
}

func GenerateOrderSn(userId int32) string {
	/*
		订单号生成规则
		年月日分秒+用户id+2为随机数
	*/
	now := time.Now()
	rand.Seed(time.Now().UnixNano())
	orderSn := fmt.Sprintf("%d%d%d%d%d%d%d%d",
		now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Nanosecond(),
		userId, rand.Intn(90)+10)
	return orderSn
}

// CartItemList 获取用户的购物车列表
func (o *OrderServer) CartItemList(ctx context.Context, req *proto.UserInfo) (*proto.CartItemListResponse, error) {
	var shopCarts []model.ShoppingCart
	var rsp proto.CartItemListResponse
	if result := global.DB.Where(&model.ShoppingCart{User: req.Id}).Find(&shopCarts); result.Error != nil {
		zap.S().Errorf("CartItemList failed,err=%s", result.Error.Error())
		return nil, result.Error
	} else {
		rsp.Total = int32(result.RowsAffected)
	}
	for _, cart := range shopCarts {
		rsp.Data = append(rsp.Data, &proto.ShopCartInfoResponse{
			Id:      cart.ID,
			UserId:  cart.User,
			GoodsId: cart.Goods,
			Nums:    cart.Nums,
			Checked: cart.Checked,
		})
	}
	return &rsp, nil
}

// CreateCartItem 添加到购物车
func (o *OrderServer) CreateCartItem(ctx context.Context, req *proto.CartItemRequest) (*proto.ShopCartInfoResponse, error) {
	var shopCart model.ShoppingCart
	if result := global.DB.Where(&model.ShoppingCart{Goods: req.GoodsId, User: req.UserId}).First(&shopCart); result.RowsAffected == 1 {
		zap.S().Errorw("CreateCartItem failed")
		// 已存在，合并到购物车
		shopCart.Nums += req.Nums
	} else {
		shopCart.User = req.UserId
		shopCart.Goods = req.GoodsId
		shopCart.Nums = req.Nums
		shopCart.Checked = false
	}
	global.DB.Save(&shopCart)
	return &proto.ShopCartInfoResponse{
		Id: shopCart.ID,
	}, nil
}

func (o *OrderServer) UpdateCartItem(ctx context.Context, req *proto.CartItemRequest) (*empty.Empty, error) {
	var shopCart model.ShoppingCart
	if result := global.DB.Where("goods=? and user=?", req.GoodsId, req.UserId).First(&shopCart); result.RowsAffected == 0 {
		zap.S().Errorf("UpdateCartItem failed,err=%s", result.Error.Error())
		return nil, status.Errorf(codes.NotFound, "记录不存在")
	}
	shopCart.Checked = req.Checked
	if req.Nums > 0 {
		shopCart.Nums = req.Nums
	}
	res := global.DB.Save(&shopCart)
	if res.Error != nil {
		return nil, res.Error
	}
	return &empty.Empty{}, nil
}

func (o *OrderServer) DeleteCartItem(ctx context.Context, req *proto.CartItemRequest) (*empty.Empty, error) {
	if result := global.DB.Where("goods=? and user=?", req.GoodsId, req.UserId).Delete(&model.ShoppingCart{}); result.RowsAffected == 0 {
		zap.S().Errorw("DeleteCartItem failed")
		return nil, status.Errorf(codes.NotFound, "购物记录不存在")
	}
	return &empty.Empty{}, nil
}

// CreateOrder 新建订单
func (o *OrderServer) CreateOrder(ctx context.Context, req *proto.OrderRequest) (*proto.OrderInfoResponse, error) {
	/*
		新建订单流程：
			1. 商品的价格查询 - 访问商品服务（跨微服务）
			2. 库存的扣减 - 访问库存服务（跨微服务）
			3. 从购物车中获取选中的商品
			4. 订单的基本信息表 - 订单的商品信息表
			5. 购物车中删除已购买记录
	*/
	var goodIds []int32
	var shopCarts []model.ShoppingCart
	goodsNumsMap := make(map[int32]int32)
	if result := global.DB.Where(&model.ShoppingCart{User: req.UserId, Checked: true}).Find(&shopCarts); result.RowsAffected == 0 {
		zap.S().Errorf("CreateOrder failed")
		return nil, status.Errorf(codes.InvalidArgument, "没有选中结算的商品")
	}
	for _, cart := range shopCarts {
		goodIds = append(goodIds, cart.Goods)
		goodsNumsMap[cart.Goods] = cart.Nums
	}
	// 跨服务调用 商品微服务
	goods, err := global.GoodsSrvClient.BatchGetGoods(context.Background(), &proto.BatchGoodsIdInfo{
		Id: goodIds,
	})
	if err != nil {
		zap.S().Errorf("CreateOrder failed,err=%s", err.Error())
		return nil, status.Errorf(codes.Internal, "批量查询商品信息失败 err:=", err.Error())
	}
	var orderAmount float32
	var orderGoods []*model.OrderGoods
	var goodsInvInfo []*proto.GoodsInvInfo
	for _, good := range goods.Data {
		orderAmount += good.ShopPrice * float32(goodsNumsMap[good.Id])
		orderGoods = append(orderGoods, &model.OrderGoods{
			Goods:      good.Id,
			GoodsName:  good.Name,
			GoodsImage: good.GoodsFrontImage,
			GoodsPrice: good.ShopPrice,
			Nums:       goodsNumsMap[good.Id],
		})
		goodsInvInfo = append(goodsInvInfo, &proto.GoodsInvInfo{
			GoodsID: good.Id,
			Num:     goodsNumsMap[good.Id],
		})
	}

	// 跨服务调用 库存微服务
	_, err = global.InventorySrvClient.Sell(context.Background(), &proto.SellInfo{
		GoodsInfo: goodsInvInfo,
	})
	if err != nil {
		zap.S().Errorf("CreateOrder failed,err=%s", err.Error())
		return nil, status.Errorf(codes.ResourceExhausted, "扣减库存失败 err:=", err.Error())
	}
	// 生成订单表
	tx := global.DB.Begin()
	order := model.OrderInfo{
		OrderSn:      GenerateOrderSn(req.UserId),
		OrderMount:   orderAmount,
		Address:      req.Address,
		SignerName:   req.Name,
		SingerMobile: req.Mobile,
		Post:         req.Post,
	}
	if res := tx.Save(&order); res.RowsAffected == 0 {
		tx.Rollback()
		if res.Error != nil {
			zap.S().Errorf("CreateOrder failed,err=%s", err.Error())
			return nil, status.Errorf(codes.Internal, "服务器有点异常了")
		}
		return nil, status.Errorf(codes.Internal, "更新失败")
	}
	for _, orderGood := range orderGoods {
		orderGood.Order = order.ID
	}
	// 批量插入
	if res := tx.CreateInBatches(orderGoods, 100); res.RowsAffected == 0 {
		tx.Rollback()
		if res.Error != nil {
			zap.S().Errorf("CreateOrder failed,err=%s", err.Error())
			return nil, status.Errorf(codes.Internal, "服务器有点异常了")
		}
		return nil, status.Errorf(codes.Internal, "插入失败 err:=", err)
	}
	if res := tx.Where(&model.ShoppingCart{User: req.UserId, Checked: true}).Delete(&model.ShoppingCart{}); res.RowsAffected == 0 {
		tx.Rollback()
		if res.Error != nil {
			zap.S().Errorf("CreateOrder failed,err=%s", err.Error())
			return nil, status.Errorf(codes.Internal, "服务器有点异常了")
		}
		return nil, status.Errorf(codes.Internal, "删除购物车 err:=", err)
	}
	tx.Commit()
	return &proto.OrderInfoResponse{Id: order.ID, OrderSn: order.OrderSn, Total: orderAmount}, nil
}

// OrderList 获得订单列表
func (o *OrderServer) OrderList(ctx context.Context, req *proto.OrderFilterRequest) (*proto.OrderListResponse, error) {
	var orderList []model.OrderInfo
	var rsp proto.OrderListResponse
	var total int64
	global.DB.Where(&model.OrderInfo{User: req.UserId}).Count(&total)
	rsp.Total = int32(total)

	// 分页
	global.DB.Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Find(&orderList)
	for _, order := range orderList {
		rsp.Data = append(rsp.Data, &proto.OrderInfoResponse{
			Id:      order.ID,
			UserId:  order.User,
			OrderSn: order.OrderSn,
			PayType: order.PayType,
			Status:  order.Status,
			Post:    order.Post,
			Total:   order.OrderMount,
			Address: order.Address,
			Name:    order.SignerName,
			Mobile:  order.SingerMobile,
			AddTime: order.CreateAt.Format("2006-04-02 15-04-05"),
		})
	}
	return &rsp, nil
}

func (o *OrderServer) OrderDetail(ctx context.Context, req *proto.OrderRequest) (*proto.OrderInfoDetailResponse, error) {
	var order model.OrderInfo
	var rsp proto.OrderInfoDetailResponse
	if result := global.DB.Where(&model.OrderInfo{BaseModel: model.BaseModel{ID: req.Id}, User: req.UserId}).First(&order); result.RowsAffected == 0 {
		if result.Error != nil {
			zap.S().Errorf("OrderDetail failed,err=%s", result.Error.Error())
			return nil, status.Errorf(codes.Internal, "服务器有点异常了")
		}
		return nil, status.Errorf(codes.NotFound, "没有这样的订单")
	}
	var orderGoods []model.OrderGoods
	if result := global.DB.Where(&model.OrderGoods{Order: order.ID}).Find(&orderGoods); result.Error != nil {
		zap.S().Errorf("OrderDetail failed,err=%s", result.Error.Error())
		return nil, status.Errorf(codes.Internal, "服务器有些问题了")
	}
	for _, good := range orderGoods {
		rsp.Goods = append(rsp.Goods, &proto.OrderItemResponse{
			GoodsId:    good.Goods,
			GoodsName:  good.GoodsName,
			GoodsPrice: good.GoodsPrice,
			Nums:       good.Nums,
		})
	}
	rsp.OrderInfo = &proto.OrderInfoResponse{
		Id:      order.ID,
		UserId:  order.User,
		OrderSn: order.OrderSn,
		PayType: order.PayType,
		Status:  order.Status,
		Post:    order.Post,
		Total:   order.OrderMount,
		Address: order.Address,
		Name:    order.SignerName,
		Mobile:  order.SingerMobile,
	}
	return &rsp, nil
}

func (o *OrderServer) UpdateOrderStatus(ctx context.Context, req *proto.OrderStatus) (*empty.Empty, error) {
	// 先查询，在更新 实际是俩sql
	if result := global.DB.Model(&model.OrderInfo{}).Where("order_sn=?", req.OrderSn).Update("status", req.Status); result.RowsAffected == 0 {
		if result.Error != nil {
			zap.S().Errorf("OrderDetail failed,err=%s", result.Error.Error())
			return nil, status.Errorf(codes.Internal, "服务器有点异常了")
		}
		return nil, status.Errorf(codes.NotFound, "订单不存在")
	}
	return &empty.Empty{}, nil
}
