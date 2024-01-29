package mongo

// mongo module
//
// Copyright (c) 2019-2024 - Valentin Kuznetsov <vkuznet AT gmail dot com>
//
// References :
// using mgo driver:
//   https://gist.github.com/boj/5412538
//   https://gist.github.com/border/3489566
// using mongo driver:
//   https://github.com/mongodb/mongo-go-driver/
//   https://www.loginradius.com/blog/engineering/mongodb-as-datasource-in-golang/

import (
	"context"
	"fmt"
	"html"
	"log"
	"strings"
	"time"

	bson "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	_ = iota
	ServerError
	DBError
	QueryError
	ParserError
	ValidationError
)

// ServerErrorName and others provides human based definition of the error
const (
	ServerErrorName     = "Server error"
	DBErrorName         = "MongoDB error"
	QueryErrorName      = "Server query error"
	ParserErrorName     = "Server parser error"
	ValidationErrorName = "Server validation error"
)

/*
// Record define Mongo record
type Record map[string]interface{}

// JsonStrings provides string representation of Record
func (r Record) JsonString() string {
	// create pretty JSON representation of the record
	data, _ := json.MarshalIndent(r, "", "    ")
	return string(data)
}

// ToString provides string representation of Record
func (r Record) String() string {
	var out []string
	for _, k := range MapKeys(r) {
		if k == "_id" {
			continue
		}
		switch v := r[k].(type) {
		case int, int64:
			out = append(out, fmt.Sprintf("%s:%d", k, v))
		case float64:
			d := int(v)
			if float64(d) == v {
				out = append(out, fmt.Sprintf("%s:%d", k, d))
			} else {
				out = append(out, fmt.Sprintf("%s:%f", k, v))
			}
		case []interface{}:
			var vals []string
			for i, val := range v {
				if i == len(v)-1 {
					vals = append(vals, fmt.Sprintf("%v", val))
				} else {
					vals = append(vals, fmt.Sprintf("%v,", val))
				}
			}
			out = append(out, fmt.Sprintf("%s:%s", k, vals))
		default:
			out = append(out, fmt.Sprintf("%s:%v", k, r[k]))
		}
	}
	return strings.Join(out, "\n")
}
*/

// ErrorRecord provides error record
func ErrorRecord(msg, etype string, ecode int) map[string]any {
	erec := make(map[string]any)
	erec["error"] = html.EscapeString(msg)
	erec["type"] = html.EscapeString(etype)
	erec["code"] = ecode
	return erec
}

// GetValue function to get int value from record for given key
func GetValue(rec map[string]any, key string) interface{} {
	var val map[string]any
	keys := strings.Split(key, ".")
	if len(keys) > 1 {
		value, ok := rec[keys[0]]
		if !ok {
			log.Printf("Unable to find key value in Record %v, key %v\n", rec, key)
			return ""
		}
		switch v := value.(type) {
		case map[string]any:
			val = v
		case []map[string]any:
			if len(v) > 0 {
				val = v[0]
			} else {
				return ""
			}
		case []interface{}:
			vvv := v[0]
			if vvv != nil {
				val = vvv.(map[string]any)
			} else {
				return ""
			}
		default:
			log.Printf("Unknown type %v, rec %v, key %v keys %v\n", fmt.Sprintf("%T", v), v, key, keys)
			return ""
		}
		if len(keys) == 2 {
			return GetValue(val, keys[1])
		}
		return GetValue(val, strings.Join(keys[1:], "."))
	}
	value := rec[key]
	return value
}

// helper function to return single entry (e.g. from a list) of given value
func singleEntry(data interface{}) interface{} {
	switch v := data.(type) {
	case []interface{}:
		return v[0]
	default:
		return v
	}
}

// GetStringValue function to get string value from record for given key
func GetStringValue(rec map[string]any, key string) (string, error) {
	value := GetValue(rec, key)
	val := fmt.Sprintf("%v", value)
	return val, nil
}

// GetSingleStringValue function to get string value from record for given key
func GetSingleStringValue(rec map[string]any, key string) (string, error) {
	value := singleEntry(GetValue(rec, key))
	val := fmt.Sprintf("%v", value)
	return val, nil
}

// GetIntValue function to get int value from record for given key
func GetIntValue(rec map[string]any, key string) (int, error) {
	value := GetValue(rec, key)
	val, ok := value.(int)
	if ok {
		return val, nil
	}
	return 0, fmt.Errorf("Unable to cast value for key '%s'", key)
}

// GetInt64Value function to get int value from record for given key
func GetInt64Value(rec map[string]any, key string) (int64, error) {
	value := GetValue(rec, key)
	out, ok := value.(int64)
	if ok {
		return out, nil
	}
	return 0, fmt.Errorf("Unable to cast value for key '%s'", key)
}

// Connection defines connection to MongoDB
type Connection struct {
	Client *mongo.Client
	URI    string
}

// InitMongoDB initializes MongoDB connection object
func InitMongoDB(uri string) {
	Mongo = Connection{URI: uri}
}

// Connect provides connection to MongoDB
func (m *Connection) Connect() *mongo.Client {
	var err error
	if m.Client != nil {
		return m.Client
	}
	client, err := mongo.NewClient(options.Client().ApplyURI(m.URI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	m.Client = client
	return client
}

// Mongo holds MongoDB connection
var Mongo Connection

// InsertAny insert records into MongoDB
func InsertAny(dbname, collname string, records []any) {
	client := Mongo.Connect()
	ctx := context.TODO()
	c := client.Database(dbname).Collection(collname)
	for _, rec := range records {
		if _, err := c.InsertOne(ctx, &rec); err != nil {
			log.Printf("Fail to insert record %v, error %v\n", rec, err)
		}
	}
}

// UpsertAny upsert records into MongoDB
func UpsertAny(dbname, collname string, records []any) error {
	client := Mongo.Connect()
	ctx := context.TODO()
	c := client.Database(dbname).Collection(collname)
	for _, rec := range records {
		spec := bson.M{}
		update := bson.D{{"$set", rec}}
		opts := options.Update().SetUpsert(true)
		if _, err := c.UpdateOne(ctx, spec, update, opts); err != nil {
			log.Printf("Fail to insert record %v, error %v\n", rec, err)
			return err
		}
	}
	return nil
}

// Insert records into MongoDB
func Insert(dbname, collname string, records []map[string]any) {
	client := Mongo.Connect()
	ctx := context.TODO()
	c := client.Database(dbname).Collection(collname)
	for _, rec := range records {
		if _, err := c.InsertOne(ctx, &rec); err != nil {
			log.Printf("Fail to insert record %v, error %v\n", rec, err)
		}
	}
}

// UpsertRecord insert record with given spec to MongoDB
func UpsertRecord(dbname, collname string, spec bson.M, rec bson.M) error {
	client := Mongo.Connect()
	ctx := context.TODO()
	c := client.Database(dbname).Collection(collname)
	opts := options.Update().SetUpsert(true)
	if _, err := c.UpdateOne(ctx, spec, rec, opts); err != nil {
		log.Printf("Fail to insert record %v, error %v\n", rec, err)
		return err
	}
	return nil
}

// Upsert records into MongoDB
func Upsert(dbname, collname, attr string, records []map[string]any) error {
	client := Mongo.Connect()
	ctx := context.TODO()
	c := client.Database(dbname).Collection(collname)
	for _, rec := range records {
		value := rec[attr].(string)
		if value == "" {
			continue
		}
		spec := bson.M{attr: value}
		update := bson.D{{"$set", rec}}
		opts := options.Update().SetUpsert(true)
		if _, err := c.UpdateOne(ctx, spec, update, opts); err != nil {
			log.Printf("Fail to insert record %v, error %v\n", rec, err)
			return err
		}
	}
	return nil
}

// Get records from MongoDB
func Get(dbname, collname string, spec bson.M, idx, limit int) []map[string]any {
	out := []map[string]any{}
	client := Mongo.Connect()
	ctx := context.TODO()
	c := client.Database(dbname).Collection(collname)
	var err error
	if limit > 0 {
		opts := options.Find().SetSkip(int64(idx)).SetLimit(int64(limit))
		cur, err := c.Find(ctx, spec, opts)
		if err != nil {
			log.Printf("ERROR: spec=%+v, error=%v", spec, err)
		}
		cur.All(ctx, &out)
	} else {
		opts := options.Find().SetSkip(int64(idx))
		cur, err := c.Find(ctx, spec, opts)
		if err != nil {
			log.Printf("ERROR: spec=%+v, error=%v", spec, err)
		}
		cur.All(ctx, &out)
	}
	if err != nil {
		log.Printf("Unable to get records, error %v\n", err)
	}
	return out
}

// GetSorted records from MongoDB sorted by given key
func GetSorted(dbname, collname string, spec bson.M, skeys []string) []map[string]any {
	out := []map[string]any{}
	client := Mongo.Connect()
	ctx := context.TODO()
	c := client.Database(dbname).Collection(collname)
	var sortSpec bson.M
	for _, s := range skeys {
		sortSpec[s] = 1
	}
	opts := options.Find().SetSort(sortSpec)
	cur, err := c.Find(ctx, spec, opts)
	cur.All(ctx, &out)
	if err != nil {
		log.Printf("Unable to sort records, error %v\n", err)
		// try to fetch all unsorted data
		cur, err := c.Find(ctx, spec)
		if err != nil {
			log.Printf("Unable to find records, error %v\n", err)
			out = append(out, ErrorRecord(fmt.Sprintf("%v", err), DBErrorName, DBError))
			return out
		}
		cur.All(ctx, &out)
	}
	return out
}

// helper function to present in bson selected fields
func sel(q ...string) (r bson.M) {
	r = make(bson.M, len(q))
	for _, s := range q {
		r[s] = 1
	}
	return
}

// Update inplace for given spec
func Update(dbname, collname string, spec, newdata bson.M) {
	client := Mongo.Connect()
	ctx := context.TODO()
	c := client.Database(dbname).Collection(collname)
	_, err := c.UpdateOne(ctx, spec, newdata)
	if err != nil {
		log.Printf("Unable to update record, spec %v, data %v, error %v\n", spec, newdata, err)
	}
}

// Count gets number records from MongoDB
func Count(dbname, collname string, spec bson.M) int {
	client := Mongo.Connect()
	ctx := context.TODO()
	c := client.Database(dbname).Collection(collname)
	nrec, err := c.CountDocuments(ctx, spec)
	if err != nil {
		log.Printf("Unable to count records, spec %v, error %v\n", spec, err)
	}
	return int(nrec)
}

// Remove records from MongoDB
func Remove(dbname, collname string, spec bson.M) {
	client := Mongo.Connect()
	ctx := context.TODO()
	c := client.Database(dbname).Collection(collname)
	_, err := c.DeleteMany(ctx, spec)
	if err != nil {
		log.Printf("Unable to remove records, spec %v, error %v\n", spec, err)
	}
}
