package db

import (
	"context"
	"github.com/full-stack-gods/GMEshortener/pkg/gme-shortener/short"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	DatabaseName            = "gme-shorts" // Database
	ShortenedCollectionName = "stonks"     // Table
)

func NewMongoDatabase(connectionString string) (db Database, err error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(connectionString))
	if err != nil {
		return nil, err
	}

	err = client.Connect(context.TODO())
	if err != nil {
		return
	}

	return &mongoDatabase{
		client:   client,
		context:  context.TODO(),
		database: DatabaseName,
	}, nil
}

// implements Database
type mongoDatabase struct {
	client   *mongo.Client
	context  context.Context
	database string
}

func (mdb *mongoDatabase) shortsCollection() *mongo.Collection {
	return mdb.client.Database(mdb.database).Collection(ShortenedCollectionName)
}

func (mdb *mongoDatabase) FindShortenedURL(id string) (res *short.ShortURL, err error) {
	filter := bson.M{
		"id": id,
	}

	cursor := mdb.shortsCollection().FindOne(mdb.context, filter)
	if err = cursor.Err(); err != nil {
		return
	}

	err = cursor.Decode(&res)
	return
}

func (mdb *mongoDatabase) SaveShortenedURL(short short.ShortURL) (err error) {
	filter := bson.M{
		"id": short.ID,
	}
	update := bson.M{
		"$set": short,
	}
	opts := options.Update().SetUpsert(true)

	_, err = mdb.shortsCollection().UpdateOne(mdb.context, filter, update, opts)
	return nil
}
