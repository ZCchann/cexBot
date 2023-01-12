package initialize

import (
	"dexBot/initialize/conf"
	"dexBot/pkg/log"
	"dexBot/pkg/logger"
	"os"
)

func InitLogger() {
	var mode = "prod"
	if conf.Conf().Debug || os.Getenv("DEBUG") == "true" {
		mode = "debug"
	}
	err := logger.NewLogger(&logger.Options{
		Mode: mode,
	})
	if err != nil {
		log.Fatalln(err)
	}
}
