package main

import (
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"shop_srvs/user_srv/global"
	handler "shop_srvs/user_srv/handler/user"
	"shop_srvs/user_srv/initialize"
	"shop_srvs/user_srv/proto"
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
	zap.S().Infof("ip:%s port:%d", global.AppConf.Grpc.IP, global.AppConf.Grpc.Port)
	server := grpc.NewServer()
	proto.RegisterUserServer(server, &handler.UserServer{})
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", global.AppConf.Grpc.IP, global.AppConf.Grpc.Port))
	if err != nil {
		panic("failed to listen:" + err.Error())
	}
	//注册服务健康检查
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	go func() {
		err = server.Serve(lis)
		if err != nil {
			panic("failed to start grpc:" + err.Error())
		}
	}()

	//接收终止信号
	//quit := make(chan os.Signal)
	//signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	//<-quit
	//if err = global.Client.Agent().ServiceDeregister("servcieid"); err != nil {
	//	zap.S().Info("注销失败")
	//}
	//zap.S().Info("注销成功")
}
