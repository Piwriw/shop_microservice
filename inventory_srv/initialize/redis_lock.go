package initialize

import (
	"errors"
	"fmt"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	goredislib "github.com/redis/go-redis/v9"
	"shop_srvs/inventory_srv/global"
)

func InitRedisLock() error {
	client := goredislib.NewClient(&goredislib.Options{
		Addr:     fmt.Sprintf("%s:%d", global.AppConf.Redis.IP, global.AppConf.Redis.Port),
		Password: global.AppConf.Redis.Password,
	})
	if client == nil {
		return errors.New("初始化Redis失败")
	}
	pool := goredis.NewPool(client) // or, pool := redigo.NewPool(...)

	// Create an instance of redisync to be used to obtain a mutual exclusion
	// lock.
	global.RedisLock = redsync.New(pool)
	return nil
}
