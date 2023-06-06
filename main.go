package testPackageLogger

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func LogInfo(message string) {
	log.Printf("Info %v", message)
}

func LogError(message string) {
	log.Printf("Error %v", message)
}

func LogWarning(message string) {
	log.Printf("Warning %v", message)
}

func ConnectToMongoDB(dbType, host, port, dbUser, dbPwd, dbName, collectionName string) (*mongo.Client, *mongo.Collection, context.Context, error) {
	// Construct MongoDB connection URI
	mongodbURI := dbType + "://" + dbUser + ":" + dbPwd + "@" + host + ":" + port

	// Configure the client connection
	clientOptions := options.Client().ApplyURI(mongodbURI)

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	// Check if the connection was successful
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	// Access the specified database and collection
	db := client.Database(dbName)
	collection := db.Collection(collectionName)

	// Create a context with a 15-second timeout
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)

	return client, collection, ctx, nil
}
