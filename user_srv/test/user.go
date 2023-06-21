package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"shop_srvs/user_srv/proto"
)

var userClient proto.UserClient
var conn *grpc.ClientConn

func Init() {
	var err error
	conn, err = grpc.Dial("127.0.0.1:50051", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	userClient = proto.NewUserClient(conn)
}
func TestGetUserList() {
	rsp, err := userClient.GetUserList(context.Background(), &proto.PageInfo{Pn: 3, PSize: 3})
	if err != nil {
		panic(err)
	}
	for _, user := range rsp.Data {
		fmt.Println(user.Mobile, user.Nickname, user.Password)
		checkResp, err := userClient.CheckPassWord(context.Background(), &proto.CheckPasswordInfo{
			Password:          "123456",
			EncryptedPassword: user.Password,
		})
		if err != nil {
			panic(err)
		}
		fmt.Println(checkResp)
	}
}
func main() {
	Init()
	defer conn.Close()
	TestGetUserList()

}
