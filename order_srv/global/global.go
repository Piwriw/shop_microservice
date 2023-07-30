package global

import (
	"gorm.io/gorm"
	"shop_srvs/order_srv/proto"
)

var (
	AppConf   AppConfig
	NacosConf Nacos

	DB *gorm.DB
	//Client *api.Client
	GoodsSrvClient     proto.GoodsClient
	InventorySrvClient proto.InventoryClient
)
