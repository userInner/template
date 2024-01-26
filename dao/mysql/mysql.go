package mysql

import (
	"database/sql"
	"fmt"
	"staging/setting"
	"time"

	"go.uber.org/zap"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *sql.DB

// Init 初始化mysql
func Init(cfg *setting.MySQLConfig) (err error) {
	client, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DB), // DSN data source name
		DefaultStringSize:         256,                                                                                                                          // string 类型字段的默认长度
		DisableDatetimePrecision:  true,                                                                                                                         // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,                                                                                                                         // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,                                                                                                                         // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,                                                                                                                        // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{})

	if err != nil {
		return err
	}
	db, err = client.DB()
	db.SetMaxIdleConns(10)           // 设置空闲时连接池连接最多数量
	db.SetMaxOpenConns(50)           // 设置连接最多连接数量
	db.SetConnMaxLifetime(time.Hour) // 设置连接最大未被使用时间
	zap.L().Info("mysql connected success!")
	return
}

// Close 手动关闭连接
func Close() {
	_ = db.Close()
}
