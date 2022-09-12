package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func catchError(err error, msg string) {
	if err != nil {
		if msg == "" {
			log.Fatal(err)
		} else {
			log.Fatal(msg)
		}
	}
}

func InitDB() *mongo.Client {
	err := godotenv.Load(".env")
	catchError(err, "Error Loading .env file for MONGO URL addr")
	MONGO_URL := os.Getenv("MONGODB_URL")
	client, err := mongo.NewClient(options.Client().ApplyURI(MONGO_URL))
	catchError(err, "")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	catchError(err, "")
	fmt.Println("Connected to MongoDb")
	return client
}

var Client *mongo.Client = InitDB()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database(os.Getenv("DATABASE_NAME")).Collection(collectionName)
	return collection
}
