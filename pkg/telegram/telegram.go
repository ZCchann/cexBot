package telegram

import (
	"dexBot/initialize/conf"
	"dexBot/pkg/logger"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

// SendMessage 使用telegram api推送信息至频道
func SendMessage(message string) {
	api := "https://api.telegram.org/bot" + conf.Conf().Telegram.BotId + "/sendMessage" //拼接url
	resp, err := http.PostForm(
		api,
		url.Values{"chat_id": {conf.Conf().Telegram.ChannelId}, "text": {message}},
	)
	if err != nil {
		logger.Error(err)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
		return
	}
	m := make(map[string]interface{}) //创建一个map接收post的返回值
	_ = json.Unmarshal(body, &m)
	if m["ok"] == false { //如果出现不是ok的返回个报错 提供检查
		logger.Errorw("telegram error", "resp", m)
	}
}

// SendError 使用telegram api推送报错信息至频道
func SendError(message string) {
	api := "https://api.telegram.org/bot" + conf.Conf().Telegram.BotId + "/sendMessage" //拼接url
	resp, err := http.PostForm(
		api,
		url.Values{"chat_id": {conf.Conf().Telegram.ErrorChannelId}, "text": {message}},
	)
	if err != nil {
		logger.Error(err)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
		return
	}
	m := make(map[string]interface{}) //创建一个map接收post的返回值
	_ = json.Unmarshal(body, &m)
	if m["ok"] == false { //如果出现不是ok的返回个报错 提供检查
		logger.Errorw("telegram error", "resp", m)
	}
}
