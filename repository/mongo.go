package mongodb

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/matryer/resync"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

var onceMongo resync.Once
var mongoConn *mongo.Client

type MongoCon struct {
	Connection *mongo.Client
}

func MongoConnect() *MongoCon {

	onceMongo.Do(func() {
		zap.L().Info("Inside mongoconnect function")
		if err := godotenv.Load(); err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}

		clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URL"))

		// Connect to MongoDB
		client, err := mongo.Connect(context.Background(), clientOptions)
		if err != nil {
			fmt.Println(err)
			zap.L().Fatal("Not able to connect to mongo", zap.Any("err", err))
		}

		// Check the connection
		err = client.Ping(context.Background(), nil)
		if err != nil {
			zap.L().Error("Failed to ping mongo", zap.Any("err", err))
		}
		zap.L().Info("Connected to MongoDB!", zap.Any("connection", client))

		mongoConn = client
	})

	return &MongoCon{Connection: mongoConn}
}

func InsertOne(connection *mongo.Client, collectionName string, document interface{}) (*mongo.InsertOneResult, error) {
	collection := connection.Database("analytics").Collection("analyticsLogs")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := collection.InsertOne(ctx, document)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (db *MongoCon) GetCollection(databaseName, collectionName string) *mongo.Collection {
	return db.Connection.Database(databaseName).Collection(collectionName)
}
