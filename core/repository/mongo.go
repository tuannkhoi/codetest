package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"git.neds.sh/technology/pricekinetics/tools/codetest/model"
)

type mongoRepo struct {
	client *mongo.Client
}

const (
	databaseName   = "codetest"
	collectionName = "event"
)

// MongoConfig is the config for connecting to a MongoDB instance.
type MongoConfig struct {
	Host string
	Port int
}

// NewMongoRepository creates a new instance of a Repository using Mongo as the persistence layer
func NewMongoRepository(ctx context.Context, cfg MongoConfig) (Repository, error) {
	host := cfg.Host
	if host == "" {
		host = "localhost"
	}

	port := cfg.Port
	if port == 0 {
		port = 27017
	}

	uri := fmt.Sprintf("mongodb://%s:%d/", host, port)

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongodb: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping mongodb: %w", err)
	}

	return &mongoRepo{
		client: client,
	}, nil
}

func (c *mongoRepo) HealthCheck(ctx context.Context) bool {
	if err := c.client.Ping(ctx, nil); err != nil {
		logrus.Errorf("could not connect to MongoDB: %v", err)
		return false
	}

	return true
}

func (c *mongoRepo) UpdateEvent(ctx context.Context, event *model.Event) error {
	collection := c.client.Database(databaseName).Collection(collectionName)
	filter := bson.M{"_id": event.ID}
	update := bson.M{"$set": event}
	opts := options.UpdateOne().SetUpsert(true)

	if _, err := collection.UpdateOne(ctx, filter, update, opts); err != nil {
		logrus.Errorf("could not update event %v", err)
		return err
	}

	return nil
}

func (c *mongoRepo) GetEventByID(ctx context.Context, id string) (*model.Event, error) {
	collection := c.client.Database(databaseName).Collection(collectionName)

	var event model.Event
	if err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&event); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			logrus.Infof("Event not found")
			return nil, nil
		}
		logrus.Errorf("could not get event %v", err)
		return nil, err
	}

	return &event, nil
}

func (c *mongoRepo) DeleteEventByID(ctx context.Context, id string) error {
	collection := c.client.Database(databaseName).Collection(collectionName)

	if _, err := collection.DeleteOne(ctx, bson.M{"_id": id}); err != nil {
		logrus.Errorf("could not delete event %v", err)
		return err
	}

	return nil
}
