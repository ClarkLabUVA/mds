package main

import (
	bson "go.mongodb.org/mongo-driver/bson"
	mongo "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"context"
	"time"
	"encoding/json"
	"reflect"
)


type MongoServer struct {
	URI      string
	Database string
	Collection string
}


func (ms MongoServer) connect() (ctx context.Context, cancel context.CancelFunc, client *mongo.Client, err error) {

	// establish connection with mongo backend
	client, err = mongo.NewClient(options.Client().ApplyURI(ms.URI))
	if err != nil {
		return
	}

	// create a context for the connection
	ctx, cancel = context.WithTimeout(context.Background(), 20*time.Second)

	// connect to the client
	err = client.Connect(ctx)
	return
}

func (ms MongoServer) InsertOne(record interface{}) (err error) {

	// establish connection with mongo backend
	mongoCtx, cancel, client, err := ms.connect()
	defer cancel()

	if err != nil {
		return
	}

	col := client.Database(ms.Database).Collection(ms.Collection)
	_, err = col.InsertOne(mongoCtx, record)

	return
}

func (ms MongoServer) FindOne(query bson.D) (record map[string]interface{}, err error) {

	// establish connection with mongo backend
	mongoCtx, cancel, client, err := ms.connect()
	defer cancel()

	if err != nil {
		return
	}

	col := client.Database(ms.Database).Collection(ms.Collection)
	err = col.FindOne(mongoCtx, query).Decode(&record)

	return
}

func (ms MongoServer) FindMany(query bson.D, results interface{}) (err error) {

	// establish connection with mongo backend
	mongoCtx, cancel, client, err := ms.connect()
	defer cancel()

	if err != nil {
		return
	}

	col := client.Database(ms.Database).Collection(ms.Collection)
	cur, err := col.Find(mongoCtx, query)

	if err != nil {
		return
	}

	err = cur.All(mongoCtx, &results)

	return

}

func (ms MongoServer) DeleteOne(query bson.D) (record map[string]interface{}, err error) {

	// establish connection with mongo backend
	mongoCtx, cancel, client, err := ms.connect()
	defer cancel()

	if err != nil {
		return
	}

	col := client.Database(ms.Database).Collection(ms.Collection)
	err = col.FindOneAndDelete(mongoCtx, query).Decode(&record)

	return
}

func (ms MongoServer) UpdateOne(query bson.D, update []byte) (record map[string]interface{}, err error) {

	// establish connection with mongo backend
	mongoCtx, cancel, client, err := ms.connect()
	defer cancel()

	if err != nil {
		return
	}

	var updateBson bson.D
	err = bson.Unmarshal(nestedUpdate(update), &updateBson)

	col := client.Database(ms.Database).Collection(ms.Collection)

	opt := options.FindOneAndUpdate().SetReturnDocument(options.After)

	err = col.FindOneAndUpdate(mongoCtx, query, update, opt).Decode(&record)

	return
}


type tuple struct {
	Key   string
	Value interface{}
}

func nestedUpdate(update []byte) (bsonUpdate []byte) {
	processedMap := make(map[string]interface{})
	updateMap := make(map[string]interface{})
	json.Unmarshal(update, &updateMap)

	resChan := make(chan tuple, 50)
	dotConvert("", updateMap, resChan)

	for {
		select {
		case elem := <-resChan:
			processedMap[elem.Key] = elem.Value
		default:
			close(resChan)
			bsonUpdate, _ = bson.Marshal(
				map[string]interface{}{
					"$set": processedMap,
				})
			return
		}
	}
}

func dotConvert(base string, input map[string]interface{}, res chan tuple) {
	for key, val := range input {
		var newBase string
		if base == "" {
			newBase = key
		} else {
			newBase = base + "." + key
		}
		if valType := reflect.ValueOf(val); valType.Kind() == reflect.Map {
			dotConvert(newBase, val.(map[string]interface{}), res)
		} else {
			//log.Println("$set: ", newBase, " ", val)
			res <- tuple{Key: newBase, Value: val}
		}
	}

	return
}
