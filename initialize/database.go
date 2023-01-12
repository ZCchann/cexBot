package initialize

import (
	"dexBot/initialize/conf"
	"dexBot/initialize/db"
	"dexBot/pkg/log"
)

func InitMongodb() {
	err := db.Mgo().Connect(
		conf.Conf().Mongodb.Username,
		conf.Conf().Mongodb.Password,
		conf.Conf().Mongodb.Host,
		conf.Conf().Mongodb.Port,
		conf.Conf().Mongodb.Database,
	)
	if err != nil {
		log.Fatalln("连接MongoDB失败:", err)
	}
}
