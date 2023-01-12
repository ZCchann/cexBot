package apps

import (
	"dexBot/model"
	"dexBot/pkg/logger"
	"dexBot/pkg/telegram"
	"fmt"
)

type KuCoinRet struct {
	Code string          `json:"code"`
	Data []KuCoinRetData `json:"data"`
}

type KuCoinRetData struct {
	BaseCurrency string `json:"baseCurrency"`
}

// 去重
func unique2(arr []KuCoinRetData) (result []string) {
	tempMap := make(map[string]bool)
	for _, item := range arr {
		_, ok := tempMap[item.BaseCurrency]
		if !ok {
			tempMap[item.BaseCurrency] = true
		}
	}
	for k := range tempMap {
		result = append(result, k)
	}
	return result
}
func KuCoinDiffer() {
	url := "https://api.kucoin.com/api/v1/symbols"
	resp, err := Get(url)
	if err != nil {
		logger.Error("Kucoin API获取异常 请检查: ", err)
		telegram.SendError(fmt.Sprintf("Kucoin API获取异常 请检查: %s", err))
		return
	}
	var ret KuCoinRet
	err = resp.JSON(&ret)
	if err != nil {
		logger.Error("绑定到结构体失败: ", err)
		return
	}

	apiData := unique2(ret.Data)
	add, data, err := Check(apiData, model.KucoinDBTableName)
	if err != nil {
		if err == fmt.Errorf("DBData delete") {

		} else {
			logger.Error(err)
			return
		}

	}
	if !add {
		// 无新数据
		return
	}
	for _, i := range data {
		t := &model.KuCoinTable{
			Coin: i,
		}
		if err := t.Save(); err != nil {
			logger.Error("新Token入库失败：", err)
			telegram.SendError(fmt.Sprintf("新Token入库失败：%s", err))
			return
		}
		text := fmt.Sprintf("kucoin交易所钱包增加了新的币种 #%s,请注意kucoin公告", i)
		logger.Debug(text)
		telegram.SendMessage(text)
	}
}
