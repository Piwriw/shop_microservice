package main

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/nacos-group/nacos-sdk-go/inner/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"os"
	"os/signal"
	"shop_srvs/order_srv/global"
	"shop_srvs/order_srv/handler"
	"shop_srvs/order_srv/initialize"
	"shop_srvs/order_srv/proto"
	"syscall"
)

func main() {
	//IP := flag.String("ip", "0.0.0.0", "ip地址")1
	//Port := flag.Int("port", 50051, "端口号")
	//flag.Parse()
	if err := initialize.InitReadNacos(); err != nil {
		fmt.Printf("init setting  failed, err:%v\n", err)
		return
	}
	if err := initialize.InitLogger(); err != nil {
		fmt.Printf("init Logger  failed, err:%v\n", err)
		return
	}
	if err := initialize.InitDB(); err != nil {
		fmt.Printf("init db  failed, err:%v\n", err)
		return
	}
	if err := initialize.InitSrvs(); err != nil {
		fmt.Printf("init Srvs  failed, err:%v\n", err)
		return
	}
	zap.S().Infof("ip:%s port:%d", global.AppConf.Grpc.IP, global.AppConf.Grpc.Port)
	server := grpc.NewServer()

	proto.RegisterOrderServer(server, &handler.OrderServer{})

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", global.AppConf.Grpc.IP, global.AppConf.Grpc.Port))
	if err != nil {
		panic("failed to listen:" + err.Error())
	}
	//注册服务健康检查
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	//服务注册
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.AppConf.Consul.IP,
		global.AppConf.Consul.Port)

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	//生成对应的检查对象
	check := &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("%s:%d", global.AppConf.Grpc.IP, global.AppConf.Grpc.Port),
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "15s",
	}

	//生成注册对象
	registration := new(api.AgentServiceRegistration)
	registration.Name = global.AppConf.Grpc.Name
	u, err := uuid.NewV4()
	if err != nil {
		panic("failed new  uuid:" + err.Error())
	}
	serviceID := fmt.Sprintf("%s", u.String())
	registration.ID = serviceID
	registration.Port = global.AppConf.Grpc.Port
	registration.Tags = global.AppConf.Grpc.Tags
	registration.Address = global.AppConf.Grpc.IP
	registration.Check = check
	//1. 如何启动两个服务
	//2. 即使我能够通过终端启动两个服务，但是注册到consul中的时候也会被覆盖
	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		panic(err)
	}

	go func() {
		err = server.Serve(lis)
		if err != nil {
			panic("failed to start grpc:" + err.Error())
		}
	}()

	//接收终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err = client.Agent().ServiceDeregister(serviceID); err != nil {
		zap.S().Info("注销失败")
	}
	zap.S().Info("注销成功")
}
