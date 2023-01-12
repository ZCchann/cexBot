package apps

import (
	"dexBot/model"
	"dexBot/pkg/logger"
	"dexBot/pkg/telegram"
	"fmt"
)

type GateIoData struct {
	Base string
}

// GateIoUnique 去重
func GateIoUnique(arr []GateIoData) (result []string) {
	tempMap := make(map[string]bool)
	for _, item := range arr {
		_, ok := tempMap[item.Base]
		if !ok {
			tempMap[item.Base] = true
		}
	}
	for k := range tempMap {
		result = append(result, k)
	}
	return result
}

func GateIoDiffer() {
	url := "https://api.gateio.ws/api/v4/spot/currency_pairs"
	resp, err := Get(url)
	if err != nil {
		logger.Error("Gate.io API获取异常 请检查: ", err)
		telegram.SendError(fmt.Sprintf("Gate.io API获取异常 请检查: %s", err))
		return
	}
	var ret []GateIoData
	err = resp.JSON(&ret)
	if err != nil {
		logger.Error("绑定到结构体失败: ", err)
		return
	}

	apiData := GateIoUnique(ret) // 返回值去重

	add, data, err := Check(apiData, model.GateIoDBTableName)
	if err != nil {
		logger.Error(err)
		return
	}
	if !add {
		// 无新数据
		return
	}
	for _, i := range data {
		t := &model.GateIoTable{
			Coin: i,
		}
		if err := t.Save(); err != nil {
			logger.Error("新Token入库失败：", err)
			telegram.SendError(fmt.Sprintf("新Token入库失败：%s", err))
			return
		}
		text := fmt.Sprintf("Gate.io交易所钱包增加了新的币种 #%s,请注意Gate.io公告", i)
		logger.Debug(text)
		telegram.SendMessage(text)
	}
}
