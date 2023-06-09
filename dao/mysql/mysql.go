package mysql

import (
	"fmt"
	"go_server/setting"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var db *sqlx.DB

func Close() {
	_ = db.Close()
}

func Init(config *setting.MysqlConfig) (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.User, config.Password,
		config.Host, config.Port,
		config.DbName,
	)
	// 也可以使用MustConnect连接不成功就panic
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		zap.L().Error("connect DB failed, err:%v\n", zap.Error(err))
		return
	}
	db.SetMaxOpenConns(viper.GetInt("mysql.max_open_connection"))
	db.SetMaxIdleConns(viper.GetInt("mysql.max_idle_connection"))
	return
}
