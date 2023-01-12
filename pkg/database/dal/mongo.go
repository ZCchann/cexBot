package dal

type DB interface {
	Table(name string) Table
}

type Table interface {
	Find(filter interface{}) Find
	Aggregate(pipeline, result interface{}) error
	InsertOne(doc interface{}) error
	InsertAll(docs []interface{}) error
	UpdateMany(filter interface{}, doc interface{}) error
	DeleteMany(filter interface{}) error
}

type Find interface {
	Projection(fields ...string) Find
	Sort(sort map[string]int8) Find
	Skip(limit int64) Find
	Limit(limit int64) Find
	Pagination(page, limit int64) ([]map[string]interface{}, int64, error)
	One(result interface{}) error
	All(result interface{}) error
	Count() (int64, error)
}
