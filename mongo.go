package main

import (
	bson "go.mongodb.org/mongo-driver/bson"
	mongo "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	//"encoding/json"
	"context"
	"time"
	//"fmt"
)

/*
	type Server struct {
	  Mongo MongoServer
	  Validation ValidationConfig
	}


	type ValidationConfig struct {
	  Validate      bool
	  ValidationURI string
	  }

	type StardogServer struct {
	    URI      string
	    Database string
	  }
*/

type MongoServer struct {
	URI	string
	Database string
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

func (ms MongoServer) InsertOne(record interface{}, collection string) (err error) {

	// establish connection with mongo backend
	mongoCtx, cancel, client, err := ms.connect()
	defer cancel()

	if err != nil {
		return
	}

	col := client.Database(ms.Database).Collection(collection)
	_, err = col.InsertOne(mongoCtx, record)

	return
}

func (ms MongoServer) FindOne(query bson.D, collection string) (record map[string]interface{}, err error) {

	// establish connection with mongo backend
	mongoCtx, cancel, client, err := ms.connect()
	defer cancel()

	if err != nil {
		return
	}

	col := client.Database(ms.Database).Collection(collection)
	err = col.FindOne(mongoCtx, query).Decode(&record)

	return
}

func (ms MongoServer) FindMany(query bson.D, collection string, results interface{}) (err error) {

	// establish connection with mongo backend
	mongoCtx, cancel, client, err := ms.connect()
	defer cancel()

	if err != nil {
		return
	}

	col := client.Database(ms.Database).Collection(collection)
	cur, err := col.Find(mongoCtx, query)

	if err != nil {
		return
	}

	err = cur.All(mongoCtx, &results)

	return

}

func (ms MongoServer) DeleteOne(query bson.D, collection string) (record map[string]interface{}, err error) {

	// establish connection with mongo backend
	mongoCtx, cancel, client, err := ms.connect()
	defer cancel()

	if err != nil {
		return
	}

	col := client.Database(ms.Database).Collection(collection)
	err = col.FindOneAndDelete(mongoCtx, query).Decode(&record)

	return
}

func (ms MongoServer) UpdateOne(query bson.D, update bson.D, collection string) (record map[string]interface{}, err error) {

	// establish connection with mongo backend
	mongoCtx, cancel, client, err := ms.connect()
	defer cancel()

	if err != nil {
		return
	}

	col := client.Database(ms.Database).Collection(collection)

	opt := options.FindOneAndUpdate().SetReturnDocument(options.After)

	err = col.FindOneAndUpdate(mongoCtx, query, update, opt).Decode(&record)

	return
}
