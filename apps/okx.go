package apps

import (
	"crypto/hmac"
	"crypto/sha256"
	"dexBot/initialize/conf"
	"dexBot/model"
	"dexBot/pkg/logger"
	"dexBot/pkg/telegram"
	"encoding/base64"
	"fmt"
	"github.com/levigross/grequests"
	"time"
)

type okxRetData struct {
	Code string    `json:"code"`
	Msg  string    `json:"msg"`
	Data []okxData `json:"data"`
}
type okxData struct {
	Ccy string `json:"ccy"`
}

// getTimestamp 获取当前时间戳 并返回okx API要求的格式
func getTimestamp() string {
	return time.Now().UTC().Format("2006-01-02T15:04:05.715Z")
}

// hmacSha256Base64 执行sha256->base64加密
func hmacSha256Base64(data string, secretKey string) string {
	key := []byte(secretKey)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(data))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// OkxUnique 去重
func OkxUnique(arr []okxData) (result []string) {
	tempMap := make(map[string]bool)
	for _, item := range arr {
		_, ok := tempMap[item.Ccy]
		if !ok {
			tempMap[item.Ccy] = true
		}
	}
	for k := range tempMap {
		result = append(result, k)
	}
	return result
}

func OkxDiffer() {
	//文档地址: https://www.okx.com/docs-v5/zh/#rest-api-funding-get-currencies
	url := "https://www.okx.com/api/v5/asset/currencies"
	method := "GET"
	requestPath := "/api/v5/asset/currencies"
	body := ""
	nowtime := getTimestamp()
	sign := nowtime + method + requestPath + body
	resp, err := grequests.Get(url, &grequests.RequestOptions{
		Headers: map[string]string{
			"OK-ACCESS-KEY":        conf.Conf().OKX.ApiKey,
			"OK-ACCESS-SIGN":       hmacSha256Base64(sign, conf.Conf().OKX.SecretKey),
			"OK-ACCESS-TIMESTAMP":  nowtime,
			"OK-ACCESS-PASSPHRASE": conf.Conf().OKX.Passphrase,
			"Content-Type":         "application/json",
		},
		RequestTimeout: 20 * time.Second,
	})
	if err != nil {
		logger.Error("OKX API获取异常 请检查: ", err)
		telegram.SendError(fmt.Sprintf("OKX API获取异常 请检查: %s", err))
		return
	}
	defer resp.Close()
	if !resp.Ok {
		logger.Error("请求okex错误 请检查: ", resp.String())
		telegram.SendError("请求okex错误 请检查: " + resp.String())
		return
	}
	var ret okxRetData
	err = resp.JSON(&ret)
	if err != nil {
		logger.Error("绑定到结构体失败: ", err)
		return
	}

	if ret.Code != "0" {
		logger.Error("请求okex错误 请检查: ", ret.Msg)
		return
	}

	apiData := OkxUnique(ret.Data)
	add, data, err := Check(apiData, model.OkxDBTableName)
	if err != nil {
		logger.Error(err)
		return
	}
	if !add {
		// 无新数据
		return
	}
	for _, i := range data {
		t := &model.OkxTable{
			Coin: i,
		}
		if err := t.Save(); err != nil {
			logger.Error("新Token入库失败：", err)
			telegram.SendError(fmt.Sprintf("新Token入库失败：%s", err))
			return
		}
		text := fmt.Sprintf("OKX交易所钱包增加了新的币种 #%s,请注意OKX公告", i)
		logger.Debug(text)
		telegram.SendMessage(text)
	}

}
