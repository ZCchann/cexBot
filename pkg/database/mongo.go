package database

import (
	"context"
	"dexBot/pkg/database/dal"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Mgo struct {
	client  *mongo.Client
	db      *mongo.Database
	timeout time.Duration
}
type table struct {
	*Mgo
	tableName string
}
type find struct {
	*table
	filter interface{}
	opts   *options.FindOptions
}

func NewMgo(db *mongo.Database) dal.DB {
	return &Mgo{
		client:  db.Client(),
		db:      db,
		timeout: 10 * time.Second,
	}
}

func (m *Mgo) Table(name string) dal.Table {
	var t = new(table)
	t.Mgo = m
	t.tableName = name
	//t.client = m.client
	//t.db = m.db
	return t
}

func (t *table) Aggregate(pipeline, result interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), t.timeout)
	defer cancel()
	ops := options.Aggregate()
	cursor, err := t.db.Collection(t.tableName).Aggregate(ctx, pipeline, ops)
	if err != nil {
		return err
	}
	return cursor.All(ctx, result)
}
func (t *table) InsertOne(doc interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), t.timeout)
	defer cancel()
	_, err := t.db.Collection(t.tableName).InsertOne(ctx, doc)
	if err != nil {
		return err
	}
	return nil
}
func (t *table) InsertAll(docs []interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), t.timeout)
	defer cancel()
	_, err := t.db.Collection(t.tableName).InsertMany(ctx, docs)
	if err != nil {
		return err
	}
	return nil
}
func (t *table) UpdateMany(filter interface{}, docs interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), t.timeout)
	defer cancel()
	data := bson.M{"$set": docs}
	_, err := t.db.Collection(t.tableName).UpdateMany(ctx, filter, data)
	if err != nil {
		return err
	}
	return nil
}
func (t *table) DeleteMany(filter interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), t.timeout)
	defer cancel()
	_, err := t.db.Collection(t.tableName).DeleteMany(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}
func (t *table) FindAndModify(filter interface{}, update interface{}, result interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), t.timeout)
	defer cancel()
	return t.db.Collection(t.tableName).FindOneAndUpdate(ctx, filter, update).Decode(result)
}
func (t *table) Find(filter interface{}) dal.Find {
	var f = new(find)
	f.table = t
	f.filter = filter
	f.opts = options.Find()
	return f
}

func (f *find) Sort(sort map[string]int8) dal.Find {
	f.opts.SetSort(sort)
	return f
}
func (f *find) Skip(i int64) dal.Find {
	f.opts.SetSkip(i)
	return f
}
func (f *find) Limit(i int64) dal.Find {
	f.opts.SetLimit(i)
	return f
}
func (f *find) Projection(fields ...string) dal.Find {
	if len(fields) == 0 {
		return f
	}
	projection := make(map[string]bool)
	for _, field := range fields {
		if field != "" {
			projection[field] = true
		}
	}
	f.opts.SetProjection(projection)
	return f
}
func (f *find) Pagination(page, limit int64) ([]map[string]interface{}, int64, error) {
	var results = make([]map[string]interface{}, 0)

	cnt, err := f.Count()
	if err != nil || cnt == 0 {
		return results, 0, err
	}

	f.opts.SetSkip((page - 1) * limit).SetLimit(limit)
	err = f.All(&results)
	return results, cnt, err
}
func (f *find) One(result interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), f.timeout)
	defer cancel()
	return f.db.Collection(f.tableName).FindOne(ctx, f.filter).Decode(result)
}
func (f *find) All(result interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), f.timeout)
	defer cancel()
	cursor, err := f.db.Collection(f.tableName).Find(ctx, f.filter, f.opts)
	if err != nil {
		return err
	}
	return cursor.All(ctx, result)
}
func (f *find) Count() (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), f.timeout)
	defer cancel()

	var (
		count int64
		err   error
	)
	if count, err = f.db.Collection(f.tableName).CountDocuments(ctx, f.filter); err != nil {
		return 0, err
	}
	return count, nil
}
