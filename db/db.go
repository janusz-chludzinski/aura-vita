package db

import (
	. "context"
	"errors"
	"fmt"
	"github.com/janusz-chludzinski/aura-vita/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

func GetConnectedClient(connectionString string, ctx Context) (*mongo.Client, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(connectionString))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("[ERROR] Getting client failed: %v", err))
	}
	err = client.Connect(ctx)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("[ERROR] Connection failed: %v", err))
	}
	return client, nil
}

func Save(entry *models.DbEntry, collection *mongo.Collection, ctx Context) (*mongo.InsertOneResult, error) {
	log.Printf("[INFO] Saving result to collection [%v] in database [%v]", collection.Name(), collection.Database().Name())
	savedEntry, err := collection.InsertOne(ctx, entry)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("[ERROR] Failed to save entry: %v", err))
	}

	log.Printf("[INFO] Created entry %v", savedEntry)
	return savedEntry, nil
}

func ConnectionString() (string, error) {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	db := os.Getenv("DB_NAME")

	if user == "" || password == "" || db == "" {
		return "", errors.New(fmt.Sprintf("[ERROR] Failed read connection string data."))
	}

	return fmt.Sprintf("mongodb://%v:%v@localhost:27017/%v", user, password, db), nil
}
