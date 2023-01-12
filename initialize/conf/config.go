package conf

import (
	"sync"
)

type (
	conf struct {
		Debug    bool           `json:"debug"            desc:"是否开启Debug模式"`
		Binance  binanceConfig  `json:"binance_config"   desc:"binance配置"`
		OKX      okxConfig      `json:"okx_config"       desc:"okx配置"`
		Bybit    bybitConfig    `json:"bybit_config"     desc:"bybit配置"`
		Mongodb  mgoConfig      `json:"mongodb"          desc:"mongodb"`
		Telegram telegramConfig `json:"telegram"         desc:"telegram机器人信息"`
	}

	binanceConfig struct {
		ApiKey    string `json:"api_key"`
		SecretKey string `json:"secret_key"`
	}

	okxConfig struct {
		ApiKey     string `json:"api_key"`
		SecretKey  string `json:"secret_key"`
		Passphrase string `json:"passphrase"`
	}

	bybitConfig struct {
		ApiKey    string `json:"api_key"`
		SecretKey string `json:"secret_key"`
	}

	mgoConfig struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Host     string `json:"host"`
		Port     string `json:"port"`
		Database string `json:"database"`
	}

	telegramConfig struct {
		BotId          string `json:"bot_id"           desc:"telegram bot id"`
		ChannelId      string `json:"channel_id"       desc:"telegram channel id"`
		ErrorChannelId string `json:"error_channel_id" desc:"报错信息推送频道"`
	}
)

var (
	c    = new(conf)
	lock = new(sync.RWMutex)
)

func Conf() *conf {
	lock.RLock()
	defer lock.RUnlock()
	return c
}
