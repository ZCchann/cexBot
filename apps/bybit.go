package apps

import (
	"crypto/hmac"
	"crypto/sha256"
	"dexBot/initialize/conf"
	"dexBot/model"
	"dexBot/pkg/logger"
	"dexBot/pkg/telegram"
	"encoding/hex"
	"fmt"
	"github.com/levigross/grequests"
	"strconv"
	"time"
)

type bybitData struct {
	RetCode int         `json:"retCode"`
	RetMsg  string      `json:"retMsg"`
	Result  bybitResult `json:"result"`
}

type bybitResult struct {
	Rows []bybitRows `json:"rows"`
}

type bybitRows struct {
	Coin string `json:"coin"`
}

func bybitHMAC(ApiKey, SecretKey, timeStamp string) string {
	hmac256 := hmac.New(sha256.New, []byte(SecretKey))
	hmac256.Write([]byte(timeStamp + ApiKey + "5000"))
	signature := hex.EncodeToString(hmac256.Sum(nil))
	return signature
}

func BybitDiffer() {
	url := "https://api.bybit.com/asset/v3/private/coin-info/query"
	timeStamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
	resp, err := grequests.Get(url, &grequests.RequestOptions{
		Headers: map[string]string{
			"X-BAPI-SIGN":        bybitHMAC(conf.Conf().Bybit.ApiKey, conf.Conf().Bybit.SecretKey, timeStamp),
			"X-BAPI-API-KEY":     conf.Conf().Bybit.ApiKey,
			"X-BAPI-TIMESTAMP":   timeStamp,
			"X-BAPI-RECV-WINDOW": "5000",
		},
		RequestTimeout: 20 * time.Second,
	})
	if err != nil {
		logger.Error("bybit API获取异常 请检查: ", err)
		telegram.SendError(fmt.Sprintf("bybit API获取异常 请检查: %s", err))
		return
	}

	var ret bybitData
	err = resp.JSON(&ret)
	if err != nil {
		logger.Error("绑定到结构体失败: ", err)
		return
	}

	if ret.RetCode == 10002 {
		// 时间戳超时 跳过
		return
	}

	if ret.RetCode != 0 && ret.RetMsg != "OK" {
		logger.Error("bybit api请求错误 请检查: ", ret.RetMsg)
		return
	}

	// 清洗数据 把api返回结果丢到一个切片中
	apiData := make([]string, 0)
	for _, k := range ret.Result.Rows {
		apiData = append(apiData, k.Coin)
	}
	add, data, err := Check(apiData, model.BybitDBTableName)
	if err != nil {
		logger.Error(err)
		return
	}
	if !add {
		// 无新数据
		return
	}
	for _, i := range data {
		t := &model.BybitTable{
			Coin: i,
		}
		if err := t.Save(); err != nil {
			logger.Error("新Token入库失败：", err)
			telegram.SendError(fmt.Sprintf("新Token入库失败：%s", err))
			return
		}
		text := fmt.Sprintf("Bybit交易所钱包增加了新的币种 #%s,请注意Bybit公告", i)
		logger.Debug(text)
		telegram.SendMessage(text)
	}

}
