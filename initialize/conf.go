package initialize

import (
	"dexBot/initialize/conf"
	"dexBot/pkg/config"
	"dexBot/pkg/log"
)

func ParseConfig(file string) {
	err := config.BindJSON(file, conf.Conf())
	if err != nil {
		log.Fatalln("解析配置文件失败：", err)
	}
}
