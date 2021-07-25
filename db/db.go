package db

import (
	. "context"
	. "errors"
	"fmt"
	. "github.com/janusz-chludzinski/aura-vita/models"
	"github.com/pkg/errors"
	. "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func GetConnectedClient(connectionString string, ctx Context) (*Client, error) {
	client, err := NewClient(options.Client().ApplyURI(connectionString))
	if err != nil {
		return nil, New(fmt.Sprintf("[ERROR] Getting client failed: %v", err))
	}
	err = client.Connect(ctx)
	if err != nil {
		return nil, New(fmt.Sprintf("[ERROR] Connection failed: %v", err))
	}
	return client, nil
}

func Save(entry *DbEntry, collection *Collection, ctx Context) (*InsertOneResult, error) {
	log.Printf("[INFO] Saving result to collection [%v] in database [%v]", collection.Name(), collection.Database().Name())
	savedEntry, err := collection.InsertOne(ctx, entry)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("[ERROR] Failed to save entry: %v", err))
	}

	log.Printf("[INFO] Created entry %v", savedEntry)
	return savedEntry, nil
}
