package apps

import (
	"dexBot/model"
	"dexBot/pkg/logger"
	"dexBot/pkg/telegram"
	"fmt"
)

type HuobiData struct {
	Status string
	Data   []string
}

func HuobiDiffer() {
	url := "https://api.huobi.pro/v1/common/currencys"
	resp, err := Get(url)
	if err != nil {
		logger.Error("火币 API获取异常 请检查: ", err)
		telegram.SendError(fmt.Sprintf("火币 API获取异常 请检查: %s", err))
		return
	}
	var ret HuobiData
	err = resp.JSON(&ret)
	if err != nil {
		logger.Error("绑定到结构体失败: ", err)
		return
	}
	if ret.Status != "ok" {
		logger.Error("获取火币api失败 请检查")
		telegram.SendError("获取火币api失败 请检查")
		return
	}

	apiData := ret.Data

	add, data, err := Check(apiData, model.HuobiDBTableName)
	if err != nil {
		logger.Error(err)
		return
	}
	if !add {
		// 无新数据
		return
	}
	for _, i := range data {
		t := &model.HuobiTable{
			Coin: i,
		}
		if err := t.Save(); err != nil {
			logger.Error("新Token入库失败：", err)
			telegram.SendError(fmt.Sprintf("新Token入库失败：%s", err))
			return
		}
		text := fmt.Sprintf("火币交易所钱包增加了新的币种 #%s,请注意火币公告", i)
		logger.Debug(text)
		telegram.SendMessage(text)
	}
}
