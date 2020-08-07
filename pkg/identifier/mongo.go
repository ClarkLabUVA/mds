package identifier

import (
	"context"
	"encoding/json"
	"fmt"
	bson "go.mongodb.org/mongo-driver/bson"
	mongo "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
	"time"

	"github.com/rs/zerolog"
	"os"
)

var (
	mongoLogger = zerolog.New(os.Stderr).With().Timestamp().Str("backend", "mongo").Logger()
)

type MongoServer struct {
	URI        string
	Database   string
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
		// log error for failing to connect
		mongoLogger.Error().
			Err(err).
			Str("operation", "InsertOne").
			Msg("failed connection to mongo")

		return
	}

	col := client.Database(ms.Database).Collection(ms.Collection)
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

	// establish connection with mongo backend
	mongoCtx, cancel, client, err := ms.connect()
	defer cancel()

	if err != nil {
		// log error for failing to connect
		mongoLogger.Error().
			Err(err).
			Str("operation", "FindOne").
			Msg("failed connection to mongo")

		return
	}

	recordMap := make(map[string]interface{})
	col := client.Database(ms.Database).Collection(ms.Collection)
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
		Bytes("record", record).
		Msg("found one success")

	return
}

func (ms MongoServer) FindMany(query bson.D, results interface{}) (err error) {

	// establish connection with mongo backend
	mongoCtx, cancel, client, err := ms.connect()
	defer cancel()

	if err != nil {
		// log error for failing to connect
		mongoLogger.Error().
			Err(err).
			Str("operation", "FindMany").
			Msg("failed connection to mongo")

		return
	}

	col := client.Database(ms.Database).Collection(ms.Collection)
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

	// establish connection with mongo backend
	mongoCtx, cancel, client, err := ms.connect()
	defer cancel()

	if err != nil {
		// log error for failing to connect
		mongoLogger.Error().
			Err(err).
			Str("operation", "DeleteOne").
			Msg("failed connection to mongo")

		return
	}

	col := client.Database(ms.Database).Collection(ms.Collection)
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

	// establish connection with mongo backend
	mongoCtx, cancel, client, err := ms.connect()
	defer cancel()

	if err != nil {
		// log error for failing to connect
		mongoLogger.Error().
			Err(err).
			Str("operation", "UpdateOne").
			Msg("failed connection to mongo")

		return
	}

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

	col := client.Database(ms.Database).Collection(ms.Collection)

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
