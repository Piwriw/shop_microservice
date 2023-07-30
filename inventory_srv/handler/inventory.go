package handler

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"shop_srvs/inventory_srv/global"
	"shop_srvs/inventory_srv/model"
	"shop_srvs/inventory_srv/proto"
)

type InventoryServer struct {
}

func (i *InventoryServer) SetInv(ctx context.Context, req *proto.GoodsInvInfo) (*empty.Empty, error) {
	// 设置库存，需要更新库存
	var inv model.Inventory
	global.DB.Where(&model.Inventory{Goods: req.GoodsID}).First(&inv)
	if inv.Goods == 0 {
		inv.Goods = req.GoodsID
	}
	inv.Stocks = req.Num

	global.DB.Save(&inv)
	return &empty.Empty{}, nil
}

func (i *InventoryServer) InvDetail(ctx context.Context, req *proto.GoodsInvInfo) (*proto.GoodsInvInfo, error) {
	var inv model.Inventory
	if result := global.DB.Where(&model.Inventory{Goods: req.GoodsID}).First(&inv); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "没有库存信息")
	}
	return &proto.GoodsInvInfo{
		GoodsID: req.GoodsID,
		Num:     inv.Stocks,
	}, nil
}

func (i *InventoryServer) Sell(ctx context.Context, req *proto.SellInfo) (*empty.Empty, error) {
	// 扣减库存 本地事务
	tx := global.DB.Begin()
	for _, goodInfo := range req.GoodsInfo {
		var inv model.Inventory
		mutex := global.RedisLock.NewMutex(fmt.Sprintf("goods_%d", goodInfo.GoodsID))
		if err := mutex.Lock(); err != nil {
			return nil, status.Errorf(codes.Internal, "获取Redis分布式锁异常")
		}
		if result := global.DB.Where(&model.Inventory{Goods: goodInfo.GoodsID}).First(&inv); result.RowsAffected == 0 {
			tx.Rollback() // 回滚
			return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
		}
		if inv.Stocks < goodInfo.Num {
			tx.Rollback()
			return nil, status.Errorf(codes.ResourceExhausted, "库存不足")
		}
		// 扣减
		inv.Stocks -= goodInfo.Num
		global.DB.Save(&inv)
		if ok, err := mutex.Unlock(); !ok || err != nil {
			return nil, status.Errorf(codes.Internal, "释放Redis分布式锁异常")
		}
	}
	tx.Commit() // 手动提交
	return &empty.Empty{}, nil
}

func (i *InventoryServer) Reback(ctx context.Context, req *proto.SellInfo) (*empty.Empty, error) {
	// 1. 订单超时归还 2. 订单创建失败 ，归还扣减库存 3.手动归还
	tx := global.DB.Begin()
	for _, goodInfo := range req.GoodsInfo {
		var inv model.Inventory
		if result := global.DB.Where(&model.Inventory{Goods: goodInfo.GoodsID}); result.RowsAffected == 0 {
			tx.Rollback()
			return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
		}
		inv.Stocks += goodInfo.Num
		global.DB.Save(&inv)
	}
	tx.Commit() // 手动提交
	return &empty.Empty{}, nil
}
