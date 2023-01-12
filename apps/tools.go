package apps

import (
	"crypto/hmac"
	"crypto/sha256"
	"dexBot/initialize/db"
	"dexBot/model"
	"dexBot/pkg/logger"
	"dexBot/pkg/telegram"
	"encoding/hex"
	"fmt"
	"github.com/levigross/grequests"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

// Get grequests的get请求 二次封装
func Get(url string) (ret *grequests.Response, err error) {
	resp, err := grequests.Get(url, &grequests.RequestOptions{
		RequestTimeout: 2 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	if !resp.Ok {
		logger.Error("resp !ok ", resp.StatusCode)
		logger.Error("error: ", resp.String())
		return nil, err
	}

	return resp, nil
}

// HmacSha256 返回hmac sha2256加密后的结果
func HmacSha256(data string, secretKey string) string {
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

// Minus 获取差集
func Minus(a []string, b []string) []string {
	var inter []string
	mp := make(map[string]bool)
	for _, s := range a {
		if _, ok := mp[s]; !ok {
			mp[s] = true
		}
	}
	for _, s := range b {
		if _, ok := mp[s]; ok {
			delete(mp, s)
		}
	}
	for key := range mp {
		inter = append(inter, key)
	}
	return inter
}

// mgoDelete 删除mongo中的数据
func mgoDelete(MgoTableName, data string) {
	err := db.Mgo().Table(MgoTableName).DeleteMany(bson.M{"coin": data})
	if err != nil {
		logger.Error(fmt.Sprintf("数据库%s数据%s删除失败 请检查: %s", MgoTableName, data, err))
		telegram.SendMessage(fmt.Sprintf("数据库%s数据删除失败 请检查: %s", MgoTableName, err))
		return
	}
	return

}

// Check 与数据库交互检查
func Check(ApiData []string, dbName string) (add bool, data []string, err error) {
	var objs = make([]*model.KuCoinTable, 0)
	var DBData []string
	err = db.Mgo().Table(dbName).Find(bson.M{}).All(&objs)
	if err != nil {
		return false, nil, fmt.Errorf("读取mongo表失败 请检查: ", err)
	}

	// 取出数据库数据 转换为切片
	for _, obj := range objs {
		DBData = append(DBData, obj.Coin)

	}

	// 对比长度 长度大于数据库 说明有新数据 返回差集
	if len(ApiData) > len(DBData) {
		minus := Minus(ApiData, DBData)
		return true, minus, nil
	}

	// 如果数据库数据比返回值多 返回多余数据 提示删除
	if len(DBData) > len(ApiData) {
		minus := Minus(DBData, ApiData)
		for _, i := range minus {
			mgoDelete(dbName, i)
		}
		//telegram.SendError(fmt.Sprintf("数据库%s 库内数据比API返回数据多 请检查\n 数据:%s", dbName, minus))
		//log.Println(fmt.Errorf("数据库%s 库内数据比API返回数据多 请检查 数据:%s", dbName, minus))
		return false, nil, fmt.Errorf("数据库%s 库内数据比API返回数据多 请检查 数据:%s", dbName, minus)
	}
	return false, nil, nil
}
