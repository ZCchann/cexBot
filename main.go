package main

import (
	"context"
	"dexBot/apps"
	"dexBot/initialize"
	"dexBot/initialize/db"
	"flag"
	"log"
	"os"
	"time"
)

func init() {
	var (
		file   string
		daemon bool
	)

	flag.StringVar(&file, "c", "conf/config.json", "configuration file")
	flag.BoolVar(&daemon, "D", false, "daemon")
	flag.Parse()

	log.Println("加载配置文件:", file)
	initialize.ParseConfig(file)

}
func main() {

	pid := os.Getpid()
	log.Println("开始启动程序", pid)
	defer log.Println("程序已退出", pid)

	log.Println("初始化日志器")
	initialize.InitLogger()

	log.Println("连接mongodb")
	initialize.InitMongodb()
	defer func() {
		_ = db.Mgo().Client().Disconnect(context.Background())
	}()

	log.Println("启动进程")
	for true {
		go apps.BinanceDiffer()
		go apps.OkxDiffer()
		go apps.CoinBaseDiffer()
		go apps.GateIoDiffer()
		go apps.HuobiDiffer()
		go apps.KuCoinDiffer()
		go apps.MexcDiffer()
		go apps.BybitDiffer()
		time.Sleep(60 * time.Second)
	}
}
