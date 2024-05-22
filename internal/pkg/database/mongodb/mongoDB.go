package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"QA-System/internal/global/config"
	"QA-System/internal/pkg/log"
)



func MongodbInit() *mongo.Collection {
	// Get MongoDB connection information from the configuration file
	user := global.Config.GetString("mongodb.user")
	pass := global.Config.GetString("mongodb.pass")
	host := global.Config.GetString("mongodb.host")
	name := global.Config.GetString("mongodb.db")
	collection := global.Config.GetString("mongodb.collection")

	// Build the MongoDB connection string
	dsn := fmt.Sprintf("mongodb://%v:%v@%v/%v", user, pass, host, name)

	// Create MongoDB client options
	clientOptions := options.Client().ApplyURI(dsn)

	// Create connection context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create MongoDB client
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Logger.Fatal("Failed to connect to MongoDB:"+ err.Error())
	}

	if err :=client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Logger.Fatal("Failed to ping MongoDB:"+ err.Error())
	}

	// Set the MongoDB database
	mdb := client.Database(name).Collection(collection)

	// Print a log message to indicate successful connection to MongoDB
	log.Logger.Info("Connected to MongoDB")
	return mdb
}
