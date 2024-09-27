package config

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SetupDatabase establishes the MongoDB connection
func SetupDatabase() (*mongo.Client, *mongo.Collection, *mongo.Collection) {
	clientOptions := options.Client().ApplyURI("mongodb+srv://kerem:H%23Ec3UCPNCDL%40ZP@cluster0.cvrfj.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0")

	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Connect to the userdb database and users and documents collections
	userCollection := client.Database("userdb").Collection("users")
	documentCollection := client.Database("userdb").Collection("documents")

	return client, userCollection, documentCollection
}
