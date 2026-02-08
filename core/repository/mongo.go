package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"git.neds.sh/technology/pricekinetics/tools/codetest/core"
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
	docBytes, err := bson.Marshal(event)
	if err != nil {
		logrus.Errorf("could not marshal event %v", err)
		return err
	}

	var setDoc bson.M
	if err := bson.Unmarshal(docBytes, &setDoc); err != nil {
		logrus.Errorf("could not unmarshal event %v", err)
		return err
	}

	if event.GetStartTime() != nil {
		// add field startTimeBSONDate (native type for time) in Mongo to helps with querying on startTime
		setDoc["startTimeBSONDate"] = time.Unix(0, event.GetStartTime().GetValue())
	}

	update := bson.M{"$set": setDoc}
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

func (c *mongoRepo) SearchEvents(
	ctx context.Context,
	filter *core.SearchEventsFilter,
	pageSize uint64,
	pageToken string,
) ([]*model.Event, string, error) {
	collection := c.client.Database(databaseName).Collection(collectionName)

	query := buildSearchEventsQuery(filter, pageToken)

	if pageSize == 0 || pageSize > 5 {
		pageSize = 5
	}

	findOpts := options.Find().
		SetLimit(int64(pageSize)).
		SetSort(bson.D{{Key: "_id", Value: 1}})

	cursor, err := collection.Find(ctx, query, findOpts)
	if err != nil {
		logrus.Errorf("could not search events %v", err)
		return nil, "", err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		if err := cursor.Close(ctx); err != nil {
			logrus.Errorf("could not close cursor %v", err)
		}
	}(cursor, ctx)

	var events []*model.Event
	for cursor.Next(ctx) {
		var event model.Event
		if err := cursor.Decode(&event); err != nil {
			logrus.Errorf("could not decode event %v", err)
			return nil, "", err
		}
		events = append(events, &event)
	}
	if err := cursor.Err(); err != nil {
		logrus.Errorf("cursor error %v", err)
		return nil, "", err
	}

	var nextToken string
	if len(events) > 0 && uint64(len(events)) == pageSize {
		nextToken = events[len(events)-1].ID
	}

	return events, nextToken, nil
}

func buildSearchEventsQuery(filter *core.SearchEventsFilter, pageToken string) *bson.M {
	query := bson.M{}

	if filter.HasBettingStatus() {
		query["bettingstatus.value"] = filter.GetBettingStatus()
	}

	if filter.HasEventVisibility() {
		query["eventvisibility.value"] = filter.GetEventVisibility()
	}

	if filter.HasStartDate() || filter.HasEndDate() {
		timeRange := bson.M{}

		if filter.HasStartDate() {
			timeRange["$gte"] = filter.GetStartDate().AsTime()
		}

		if filter.HasEndDate() {
			timeRange["$lte"] = filter.GetEndDate().AsTime()
		}

		query["startTimeBSONDate"] = timeRange
	}

	if pageToken != "" {
		query["_id"] = bson.M{"$gt": pageToken}
	}

	return &query
}
