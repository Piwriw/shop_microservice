package initialize

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"shop_srvs/user_srv/global"
)

func InitReadNacos() error {

	viper.SetConfigFile("./user_srv/config/dev_nacos.yaml")

	err := viper.ReadInConfig() // 读取配置信息
	if err != nil {
		// 读取配置信息失败
		zap.S().Error("viper.ReadInConfig failed, err:%v\n", err)
		return err
	}

	// 把读取到的配置信息反序列化到 Conf 变量中
	if err := viper.Unmarshal(&global.NacosConf); err != nil {
		fmt.Printf("viper.Unmarshal failed, err:%v\n", err)
		return err
	}
	err = InitNacos()
	if err != nil {
		return err
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件修改了...")
		if err := viper.Unmarshal(&global.NacosConf); err != nil {
			fmt.Printf("viper.Unmarshal failed, err:%v\n", err)
		}
	})
	return nil
}

func InitSetting() error {

	viper.SetConfigFile("./user_srv/config/dev.yaml")

	err := viper.ReadInConfig() // 读取配置信息
	if err != nil {
		// 读取配置信息失败
		zap.S().Error("viper.ReadInConfig failed, err:%v\n", err)
		return err
	}

	// 把读取到的配置信息反序列化到 Conf 变量中
	if err := viper.Unmarshal(&global.AppConf); err != nil {
		fmt.Printf("viper.Unmarshal failed, err:%v\n", err)
		return err
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件修改了...")
		if err := viper.Unmarshal(&global.AppConf); err != nil {
			fmt.Printf("viper.Unmarshal failed, err:%v\n", err)
		}
	})
	return nil
}
