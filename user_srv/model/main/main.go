package main

import (
	"fmt"
	"shop_srvs/user_srv/global"
	"shop_srvs/user_srv/model"
	"shop_srvs/user_srv/utils"
	"time"
)

func main() {
	for i := 0; i < 10; i++ {
		now := time.Now()
		user := model.User{
			NickName: fmt.Sprintf("piwriw%d", i),
			Mobile:   fmt.Sprintf("111111%d", i),
			Password: utils.Encode("123456"),
			Birthday: &now,
		}
		global.DB.Save(&user)
	}
}
