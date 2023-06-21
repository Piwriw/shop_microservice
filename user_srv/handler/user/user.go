package handler

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"shop_srvs/user_srv/global"
	"shop_srvs/user_srv/model"
	"shop_srvs/user_srv/proto"
	"shop_srvs/user_srv/utils"
	"time"
)

type UserServer struct {
}

func Model2Response(user model.User) proto.UserInfoResponse {
	userInfoRsp := proto.UserInfoResponse{
		Id:       user.ID,
		Password: user.Password,
		//Mobile:   good.Mobile,
		Nickname: user.NickName,
		Gender:   user.Gender,
		Role:     int32(user.Role),
	}
	if user.Birthday != nil {
		userInfoRsp.BirthDay = uint64(user.Birthday.Unix())
	}
	return userInfoRsp
}
func (u *UserServer) GetUserByMobile(ctx context.Context, req *proto.MobileRequest) (*proto.UserInfoResponse, error) {
	//通过手机号码查询用户
	var user model.User
	result := global.DB.Where(&model.User{Mobile: req.Mobile}).Find(&user)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	if result.Error != nil {
		return nil, result.Error
	}
	userInfoRsp := Model2Response(user)
	return &userInfoRsp, nil
}
func (u *UserServer) GetUserByID(ctx context.Context, req *proto.IdRequest) (*proto.UserInfoResponse, error) {
	// 通过主键ID查询用户
	var user model.User
	result := global.DB.Find(&user)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	if result.Error != nil {
		return nil, result.Error
	}
	userInfoRsp := Model2Response(user)
	return &userInfoRsp, nil
}
func (u *UserServer) CreateUser(ctx context.Context, req *proto.CreateUserInfo) (*proto.UserInfoResponse, error) {
	// 新建用户
	// 查询用户是否存在
	user := new(model.User)
	result := global.DB.Where(&model.User{Mobile: req.Mobile}).Find(&user)
	if result.RowsAffected == 1 {
		return nil, status.Errorf(codes.AlreadyExists, "用户已存在")
	}
	user.Mobile = req.Mobile
	user.NickName = req.NickName
	user.Password = utils.Encode(req.PassWord)
	result = global.DB.Create(&user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	userInfoResp := Model2Response(*user)
	return &userInfoResp, nil
}
func (u *UserServer) UpdateUser(ctx context.Context, req *proto.UpdateUserInfo) (*empty.Empty, error) {
	var user model.User
	result := global.DB.First(&user, req.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	birthDay := time.Unix(int64(req.BirthDay), 0)
	user.NickName = req.NickName
	user.Birthday = &birthDay
	user.Gender = req.Gender
	result = global.DB.Save(&user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	return &empty.Empty{}, nil
}
func (u *UserServer) CheckPassWord(ctx context.Context, req *proto.CheckPasswordInfo) (*proto.CheckResponse, error) {
	// 校验密码
	newPassword := utils.Encode(req.Password)
	if newPassword != req.EncryptedPassword {
		return &proto.CheckResponse{IsSuccess: false}, nil
	}
	return &proto.CheckResponse{IsSuccess: true}, nil
}
func (u *UserServer) GetUserList(ctx context.Context, req *proto.PageInfo) (*proto.UserListResponse, error) {
	var users []model.User
	result := global.DB.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	rsp := &proto.UserListResponse{}
	rsp.Total = int32(result.RowsAffected)
	if req.Pn <= 0 {
		req.Pn = 1
	}
	switch {
	case req.PSize > 100:
		req.PSize = 100
	case req.PSize <= 0:
		req.PSize = 10
	}
	offset := (req.Pn - 1) * req.PSize
	result = global.DB.Offset(int(offset)).Limit(int(req.PSize)).Find(&users)
	//result = global.DB.Scopes(Paginate(int(req.Pn), int(req.PSize))).Find(&users)
	for _, user := range users {
		userInfoRsp := Model2Response(user)
		rsp.Data = append(rsp.Data, &userInfoRsp)
	}
	return rsp, nil
}
func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}
		offset := (page - 1) * pageSize
		return global.DB.Offset(offset).Limit(pageSize)
	}
}
