package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"shop_srvs/good_srv/proto"
)

var goodClient proto.GoodsClient
var conn *grpc.ClientConn

func Init() {
	var err error
	conn, err = grpc.Dial("192.168.44.190:50051", grpc.WithTransportCredentials((insecure.NewCredentials())))
	if err != nil {
		panic(err)
	}
	goodClient = proto.NewGoodsClient(conn)
}
func TestGetBrandList() {
	rsp, err := goodClient.BrandList(context.Background(), &proto.BrandFilterRequest{Pages: 2, PagePerNums: 5})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Total)
	for _, brand := range rsp.Data {
		fmt.Println(brand.Name, brand.Id, brand.Logo)

	}
}
func main() {
	Init()
	defer conn.Close()
	TestGetBrandList()

}
