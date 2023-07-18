package global

import (
	"github.com/go-redsync/redsync/v4"
	"gorm.io/gorm"
)

var (
	AppConf   AppConfig
	NacosConf Nacos

	DB        *gorm.DB
	RedisLock *redsync.Redsync
	//Client *api.Client
)
