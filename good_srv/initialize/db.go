package initialize

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"shop_srvs/good_srv/global"
	"time"
)

func InitDB() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		global.AppConf.Mysql.User, global.AppConf.Mysql.Password, global.AppConf.Mysql.IP, global.AppConf.Mysql.Port, global.AppConf.Mysql.DbName)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,        // Don't include params in the SQL log
			Colorful:                  false,       // Disable color
		},
	)

	// Globally mode
	var err error
	global.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return err
	}
	//err = global.DB.AutoMigrate(&model.User{})
	//if err != nil {
	//	panic(err)
	//}
	return nil
}
