package db

import (
	"context"
	"dexBot/pkg/database"
	"dexBot/pkg/database/dal"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type mdb struct {
	*mongo.Database
}

var m = new(mdb)

func Mgo() *mdb {
	return m
}

func (m *mdb) Connect(username, password, addr, port, database string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := options.Client()
	opts.ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%s", username, password, addr, port))
	opts.SetMinPoolSize(uint64(8))
	opts.SetMaxPoolSize(uint64(800))
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return err
	}
	m.Database = client.Database(database)
	return nil
}

func (m *mdb) Table(name string) dal.Table {
	return database.NewMgo(m.Database).Table(name)
}
