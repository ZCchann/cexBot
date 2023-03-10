package model

import (
	"dexBot/initialize/db"
	"dexBot/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BybitTable struct {
	ID   primitive.ObjectID `json:"id"           form:"_id"         bson:"_id"          desc:"ID"`
	Coin string             `json:"coin"         form:"coin"        bson:"coin"         desc:"coin"`
}

const BybitDBTableName = "bybit"

func (*BybitTable) TableName() string {
	return BybitDBTableName
}

func (m *BybitTable) Save() error {
	lock.Lock()
	defer lock.Unlock()
	var err error
	if m.ID.IsZero() {
		m.ID = primitive.NewObjectID()
		err = db.Mgo().Table(m.TableName()).InsertOne(m)
	} else {
		err = db.Mgo().Table(m.TableName()).UpdateMany(bson.M{"_id": m.ID}, m)
	}
	if err != nil {
		logger.Error("保存数据失败：", err)
	}
	return err
}

func (m *BybitTable) Delete() error {
	return db.Mgo().Table(m.TableName()).DeleteMany(bson.M{"coin": m.Coin})
}
