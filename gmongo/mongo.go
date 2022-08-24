package gmongo

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	_ "go.mongodb.org/mongo-driver/x/mongo/driver/operation"
)

type MongoDBConf struct {
	Enable   bool   `mapstructure:"MONGODB_ENABLE"`
	URI      string `mapstructure:"MONGODB_SERVER_URI"`
	Database string `mapstructure:"MONGODB_DATABASE"`
}

var (
	ctx context.Context
	// mongoClient   *mongo.Client
	mongoDatabase *mongo.Database
	logger        *logrus.Entry
)

func Init(mongoConf *MongoDBConf, log *logrus.Entry) {

	logger = log

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	// mongoconn := options.Client().ApplyURI("mongodb://root:password123@localhost:6000")
	mongoConn := options.Client().ApplyURI(mongoConf.URI)
	mongoClient, err := mongo.Connect(ctx, mongoConn)

	if err != nil {
		panic(err)
	}

	if err := mongoClient.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	mongoDatabase = mongoClient.Database(mongoConf.Database)

	logger.Info(fmt.Sprintf("MongoDB successfully connected: Server=%s, Database=%s", mongoConn.GetURI(), mongoConf.Database))

	defer mongoClient.Disconnect(ctx)

}

func InsertOne(collectionName string, document map[string]interface{}) primitive.ObjectID {

	collection := mongoDatabase.Collection(collectionName)
	res, err := collection.InsertOne(ctx, document)

	if err != nil {
		panic(err)
	}

	// return res.InsertedID.(primitive.ObjectID).Hex()
	return res.InsertedID.(primitive.ObjectID)

}

// f := &foo{}
// f, ok := baz.(*foo)
func FindOne(collectionName string, objectid string, result map[string]interface{}) {

	err := mongoDatabase.Collection(collectionName).FindOne(ctx, bson.M{"_id": objectid}).Decode(&result)

	if err != nil {
		panic(err)
	}

}

func QueryOne(collectionName string, filter map[string]interface{}, result map[string]interface{}) {

	query, err := bson.Marshal(filter)

	if err != nil {
		panic(err)
	}

	err = mongoDatabase.Collection(collectionName).FindOne(ctx, query).Decode(&result)

	if err != nil {
		panic(err)
	}

}

func Find(collectionName string, filter map[string]interface{}) []interface{} {

	query, err := bson.Marshal(filter)

	if err != nil {
		panic(err)
	}

	cursor, err := mongoDatabase.Collection(collectionName).Find(ctx, query)

	if err != nil {
		panic(err)
	}

	defer cursor.Close(ctx)

	var interfaceSlice []interface{} = make([]interface{}, 0)
	var i int = 0

	for cursor.Next(ctx) {

		var result bson.D
		err := cursor.Decode(&result)

		if err != nil {
			panic(err)
		}

		interfaceSlice[i] = result

		i++

	}

	if err := cursor.Err(); err != nil {
		panic(err)
	}

	return interfaceSlice

}
