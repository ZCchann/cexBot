package apps

import (
	"dexBot/initialize/conf"
	"dexBot/model"
	"dexBot/pkg/logger"
	"dexBot/pkg/telegram"
	"fmt"
	"github.com/levigross/grequests"
	"strconv"
	"time"
)

type BinanceData struct {
	Coin string `json:"coin"`
}

type BinanceError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func BinanceDiffer() {
	// 文档地址 https://binance-docs.github.io/apidocs/spot/cn/#user_data
	url := "https://api.binance.com/sapi/v1/capital/config/getall"
	t := strconv.FormatInt(time.Now().UnixMilli(), 10) // 获取当前时间戳 转换为string

	// 如果主机配置较差 币安经常返回{"code":-1021,"msg":"Timestamp for this request is outside of the recvWindow."} 建议启用下面这行
	// t := strconv.FormatInt(time.Now().UnixMilli() +100, 10) // 获取当前时间戳 转换为string

	resp, err := grequests.Get(url, &grequests.RequestOptions{
		Params: map[string]string{
			"timestamp": t,
			"signature": HmacSha256(fmt.Sprintf("timestamp=%s", t), conf.Conf().Binance.SecretKey),
		},
		Headers: map[string]string{
			"X-MBX-APIKEY": conf.Conf().Binance.ApiKey,
		},
		RequestTimeout: 20 * time.Second,
	})
	if err != nil {
		logger.Error("币安交易所获取API错误 请检查: ", err)
		telegram.SendError("币安交易所获取API错误 请检查: " + err.Error())
		return
	}
	if !resp.Ok {
		var retError BinanceError
		err = resp.JSON(&retError)
		if err != nil {
			logger.Error("绑定到结构体失败: ", err)
			return
		}

		if retError.Code == -1021 {
			// 请求超时 跳过
			return
		}

		telegram.SendError("币安交易所获取API错误回传信息 请检查: " + resp.String())
		return
	}
	var ret []BinanceData
	err = resp.JSON(&ret)
	if err != nil {
		logger.Error("绑定到结构体失败: ", err)
		return
	}

	var apiData []string
	for _, i := range ret {
		apiData = append(apiData, i.Coin)
	}

	add, data, err := Check(apiData, model.BinanceDBTableName)
	if err != nil {
		logger.Error(err)
		return
	}
	if !add {
		// 无新数据
		return
	}
	for _, i := range data {
		t := &model.BinanceTable{
			Coin: i,
		}
		if err := t.Save(); err != nil {
			logger.Error("新Token入库失败：", err)
			telegram.SendError(fmt.Sprintf("新Token入库失败：%s", err))
			return
		}
		text := fmt.Sprintf("币安交易所钱包增加了新的币种 #%s,请注意币安公告", i)
		logger.Debug(text)
		telegram.SendMessage(text)
	}
}
