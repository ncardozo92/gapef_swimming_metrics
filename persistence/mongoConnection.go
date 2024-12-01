package persistence

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ncardozo92/gapef_swimming_metrics/logging"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const DATABASE_NAME string = "gapef_swimming_metrics"

var database *mongo.Database

func GetDatabase() *mongo.Database {

	host := os.Getenv("MONGODB_HOST")
	username := os.Getenv("MONGODB_USER")
	password := os.Getenv("MONGODB_PASS")
	port := os.Getenv("MONGODB_PORT")
	dbName := os.Getenv("MONGODB_DATABASE_NAME")

	connectionString := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?authSource=admin", username, password, host, port, dbName)
	logging.LogInfo("Database connection string: %s", connectionString)

	if database == nil {

		mongoClient, getClientErr := mongo.Connect(context.TODO(), options.Client().ApplyURI(connectionString))

		if getClientErr != nil {
			log.Fatalln(getClientErr.Error())
		}

		database = mongoClient.Database(DATABASE_NAME)
	}
	return database
}
