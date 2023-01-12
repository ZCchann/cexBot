package apps

import (
	"dexBot/model"
	"dexBot/pkg/logger"
	"dexBot/pkg/telegram"
	"fmt"
)

type MexcV2Ret struct {
	Code int
	Data []MexcV2Data
}

type MexcV2Data struct {
	Currency string `json:"currency"`
}

func MexcDiffer() {
	url := "https://www.mexc.com/open/api/v2/market/coin/list"

	resp, err := Get(url)
	if err != nil {
		logger.Error("Mexc API获取异常 请检查: ", err)
		telegram.SendError(fmt.Sprintf("Mexc API获取异常 请检查: %s", err))
		return
	}
	if err != nil {
		return
	}
	var ret MexcV2Ret
	err = resp.JSON(&ret)
	if err != nil {
		logger.Error("绑定到结构体失败: ", err)
		return
	}
	if ret.Code != 200 {
		return
	}

	var apiData []string
	for _, i := range ret.Data {
		apiData = append(apiData, i.Currency)
	}

	add, data, err := Check(apiData, model.MexcDBTableName)
	if err != nil {
		logger.Error(err)
		return
	}
	if !add {
		// 无新数据
		return
	}
	for _, i := range data {
		t := &model.MexcTable{
			Coin: i,
		}
		if err := t.Save(); err != nil {
			logger.Error("新Token入库失败：", err)
			telegram.SendError(fmt.Sprintf("新Token入库失败：%s", err))
			return
		}
		text := fmt.Sprintf("MEXC交易所钱包增加了新的币种 #%s,请注意MEXC公告", i)
		logger.Debug(text)
		telegram.SendMessage(text)
	}
}
