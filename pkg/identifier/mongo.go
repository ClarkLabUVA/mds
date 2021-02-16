//© 2020 By The Rector And Visitors Of The University Of Virginia

//Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
//The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

//THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
package identifier

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"
	bson "go.mongodb.org/mongo-driver/bson"
	mongo "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"github.com/rs/zerolog"
)

var (
	mongoLogger = zerolog.New(os.Stderr).With().Timestamp().Str("backend", "mongo").Logger()
)

type MongoServer struct {
	URI        string
	Database   string
	Collection string
    Client *mongo.Client
}

func (ms MongoServer) Connect(clientContext context.Context) (client *mongo.Client, err error) {

	// establish connection with mongo backend

    opt := options.Client().ApplyURI(ms.URI).SetMaxConnIdleTime(0).SetMaxPoolSize(1000)

	client, err = mongo.NewClient(opt)

    if err != nil {
        // log error for failing to connect
        mongoLogger.Error().
            Err(err).
            Str("operation", "ConnectMongo").
            Msg("Failed client creation")

        return
    }

	// connect to the client
	err = client.Connect(clientContext)

    if err != nil {
        // log error for failing to connect
        mongoLogger.Error().
            Err(err).
            Str("operation", "ConnectMongo").
            Msg("Client failed to connect to mongo")
    }

	return
}

func (ms MongoServer) InsertOne(record interface{}) (err error) {

    // create a new context for the operation
    mongoCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

	col := ms.Client.Database(ms.Database).Collection(ms.Collection)
	_, err = col.InsertOne(mongoCtx, record)

	if err != nil {
		mongoLogger.Error().
			Err(err).
			Str("operation", "InsertOne").
			Interface("record", record).
			Msg("failed insert one operation to mongo")

		return
	}

	mongoLogger.Info().
		Str("operation", "InsertOne").
		Interface("record", record).
		Msg("created record in mongo")

	return
}

func (ms MongoServer) FindOne(query bson.D) (record []byte, err error) {

    // create a new context for the operation
    mongoCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

	recordMap := make(map[string]interface{})
	col := ms.Client.Database(ms.Database).Collection(ms.Collection)
	err = col.FindOne(mongoCtx, query).Decode(&recordMap)

	if err != nil {
		mongoLogger.Error().
			Err(err).
			Str("operation", "FindOne").
			Interface("query", query).
			Msg("error finding document in mongo")

		return
	}

	record, err = json.Marshal(recordMap)

	mongoLogger.Info().
		Str("operation", "FindOne").
		Interface("query", query).
		Str("record", string(record)).
		Msg("found one success")

	return
}

func (ms MongoServer) FindMany(query bson.D, results interface{}) (err error) {

    // create a new context for the operation
    mongoCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

	col := ms.Client.Database(ms.Database).Collection(ms.Collection)
	cur, err := col.Find(mongoCtx, query)

	if err != nil {

		mongoLogger.Error().
			Err(err).
			Str("operation", "FindMany").
			Interface("query", query).
			Msg("failed to execute query")

		return
	}

	err = cur.All(mongoCtx, &results)

	if err != nil {

		mongoLogger.Error().
			Err(err).
			Str("operation", "FindMany").
			Interface("query", query).
			Msg("failed to unmarshal results from query cursor")

		return
	}

	mongoLogger.Info().
		Str("operation", "FindMany").
		Interface("query", query).
		Msg("success")

	return

}

func (ms MongoServer) DeleteOne(query bson.D) (record map[string]interface{}, err error) {

    // create a new context for the operation
    mongoCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

	col := ms.Client.Database(ms.Database).Collection(ms.Collection)
	err = col.FindOneAndDelete(mongoCtx, query).Decode(&record)

	if err != nil {

		mongoLogger.Error().
			Err(err).
			Str("operation", "DeleteOne").
			Interface("query", query).
			Msg("error deleting document from mongo")

		return
	}

	mongoLogger.Info().
		Str("operation", "DeleteOne").
		Interface("query", query).
		Interface("record", record).
		Msg("successfully deleted document from mongo")

	return
}

func (ms MongoServer) UpdateOne(query bson.D, update []byte) (record []byte, err error) {

    // create a new context for the operation
    mongoCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

	// nestedUpdate converts to bson
	nested, err := nestedUpdate(update)
	if err != nil {
		err = fmt.Errorf(`{"message": "Failed to Convert to Nested Dot Format", "error": "%s"}`, err.Error())

		mongoLogger.Error().
			Err(err).
			Str("operation", "UpdateOne").
			Interface("query", query).
			Bytes("update", update).
			Msg("failed to convert update to dot notation")

		return
	}

	col := ms.Client.Database(ms.Database).Collection(ms.Collection)
	opt := options.FindOneAndUpdate().SetReturnDocument(options.After)

	rec := make(map[string]interface{})
	err = col.FindOneAndUpdate(mongoCtx, query, nested, opt).Decode(&rec)

	if err != nil {
		err = fmt.Errorf(`{"message": "Mongo Update Operation Failed", "error": "%s"}`, err.Error())

		mongoLogger.Error().
			Err(err).
			Str("operation", "UpdateOne").
			Interface("query", query).
			Bytes("update", update).
			Msg("failed UpdateOne operation in mongo")

		return
	}

	record, err = json.Marshal(rec)

	mongoLogger.Info().
		Str("operation", "UpdateOne").
		Interface("query", query).
		Bytes("update", update).
		Bytes("nested", nested).
		Interface("record", record).
		Msg("succeeded UpdateOne operation in mongo")

	return
}

type tuple struct {
	Key   string
	Value interface{}
}

func nestedUpdate(update []byte) (bsonUpdate []byte, err error) {
	processedMap := make(map[string]interface{})
	updateMap := make(map[string]interface{})
	err = json.Unmarshal(update, &updateMap)
	if err != nil {
		return
	}

	resChan := make(chan tuple, 50)
	dotConvert("", updateMap, resChan)

	for {
		select {
		case elem := <-resChan:
			processedMap[elem.Key] = elem.Value
		default:
			close(resChan)
			bsonUpdate, err = bson.Marshal(
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
