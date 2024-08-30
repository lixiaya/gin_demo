package config

import (
	"fmt"
	"gin_demo/global"
	"gin_demo/internal/model"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitMysql() *gorm.DB {
	dsn := viper.GetString("db.dsn")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		//fmt.Sprintln(err.Error())
		global.Logger.Error(err) //记录数据库连接错误日志
		panic("failed to connect database")
	}
	global.Logger.Info("connect database success(连接数据库成功)")
	fmt.Println("connect database success")
	//自动迁移
	/*
			检查数据库中是否已经存在一个名为users（默认表名是模型名的小写复数形式）的表。
		    如果表不存在，GORM将根据model.User结构体的定义创建一个新的表。
		    如果表已存在，GORM将比较表结构和model.User结构体的定义，并根据需要添加或修改列以匹配模型。
		    如果model.User结构体中定义了索引或唯一约束，GORM也会相应地在数据库中创建它们。
	*/
	err = db.AutoMigrate(&model.User{})
	if err != nil {
		global.Logger.Error(fmt.Sprintln("自动迁移失败TAT:", err.Error()))
		fmt.Println("自动迁移失败：", err.Error())
	} else {
		global.Logger.Info("自动迁移成功或数据库存在对应表")
		fmt.Println("自动迁移成功或数据库存在对应表")
	}
	return db
}
