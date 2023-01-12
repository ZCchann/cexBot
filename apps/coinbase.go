package apps

import (
	"dexBot/model"
	"dexBot/pkg/logger"
	"dexBot/pkg/telegram"
	"fmt"
)

type CoinBaseData struct {
	Id string `json:"id"`
}

func CoinBaseDiffer() {
	url := "https://api.exchange.coinbase.com/currencies"
	resp, err := Get(url)
	if err != nil {
		logger.Error("coinbase API获取异常 请检查: ", err)
		telegram.SendError(fmt.Sprintf("coinbase API获取异常 请检查: %s", err))
		return
	}
	var ret []CoinBaseData
	err = resp.JSON(&ret)
	if err != nil {
		logger.Error("绑定到结构体失败: ", err)
		return
	}

	var apiData []string
	for _, i := range ret {
		apiData = append(apiData, i.Id)
	}

	add, data, err := Check(apiData, model.CoinBaseDBTableName)
	if err != nil {
		logger.Error(err)
		return
	}
	if !add {
		// 无新数据
		return
	}
	for _, i := range data {
		t := &model.CoinBaseTable{
			Coin: i,
		}
		if err := t.Save(); err != nil {
			logger.Error("新Token入库失败：", err)
			telegram.SendError(fmt.Sprintf("新Token入库失败：%s", err))
			return
		}
		text := fmt.Sprintf("coinbase交易所钱包增加了新的币种 #%s,请注意coinbase公告", i)
		logger.Debug(text)
		telegram.SendMessage(text)
	}
}
