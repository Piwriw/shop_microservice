package global

import (
	"github.com/hashicorp/consul/api"
	"gorm.io/gorm"
)

var (
	AppConf   AppConfig
	NacosConf Nacos

	DB     *gorm.DB
	Client *api.Client
)
