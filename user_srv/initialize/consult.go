package initialize

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"shop_srvs/user_srv/global"
)

func ConsulRegister(address, name string, port int, tags []string, id string) error {
	cfg := api.DefaultConfig()
	cfg.Address = global.AppConf.Consul.IP
	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	global.Client = client
	// 生成对应的检查对象
	check := &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("http://%s:50051/health", global.AppConf.Consul.IP),
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "10",
	}
	// 生成注册对象
	registerion := &api.AgentServiceRegistration{
		Name:    name,
		ID:      id,
		Port:    port,
		Tags:    tags,
		Address: address,
		Check:   check,
	}
	err = global.Client.Agent().ServiceRegister(registerion)
	if err != nil {
		return err
	}
	return nil
}
