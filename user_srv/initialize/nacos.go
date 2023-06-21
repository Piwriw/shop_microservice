package initialize

import (
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"gopkg.in/yaml.v2"
	"shop_srvs/user_srv/global"
)

func InitNacos() error {
	sc := []constant.ServerConfig{
		{
			IpAddr: global.NacosConf.Host,
			Port:   global.NacosConf.Port,
		},
	}

	cc := constant.ClientConfig{
		NamespaceId: global.NacosConf.NameSpace, // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId
		//NamespaceId:         "93e9005f-d1ec-4988-8586-332927badab0", // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "./tmp/nacos/log",
		CacheDir:            "./tmp/nacos/cache",
		//RotateTime:          "1h",
		//MaxAge:              3,
		LogLevel: "debug",
	}

	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})
	if err != nil {
		return err
	}
	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: global.NacosConf.DataID,
		Group:  global.NacosConf.Group,
	})

	if err != nil {
		return err
	}
	//fmt.Println(content) //字符串 - yaml
	//serverConfig := global.AppConfig{}
	//想要将一个json字符串转换成struct，需要去设置这个struct的tag
	//err = json.Unmarshal([]byte(content), &global.AppConf)
	err = yaml.Unmarshal([]byte(content), &global.AppConf)
	if err != nil {
		return err
	}
	//err = configClient.ListenConfig(vo.ConfigParam{
	//	DataId: "good-web.json",
	//	Group:  "dev",
	//	OnChange: func(namespace, group, dataId, data string) {
	//		fmt.Println("配置文件变化")
	//		fmt.Println("group:" + group + ", dataId:" + dataId + ", data:" + data)
	//	},
	//})
	//time.Sleep(3000 * time.Second)
	return nil
}
