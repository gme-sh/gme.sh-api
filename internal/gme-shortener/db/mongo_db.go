package db

import (
	"context"
	"time"

	"github.com/full-stack-gods/GMEshortener/pkg/gme-shortener/short"
	"github.com/patrickmn/go-cache"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	// DatabaseName -> Name of the Database
	DatabaseName = "gme-shorts"
	// ShortenedCollectionName -> Name of the collection
	ShortenedCollectionName = "stonks"
)

// NewMongoDatabase -> Create a new MongoDB
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
		cache:    cache.New(10*time.Minute, 15*time.Minute),
	}, nil
}

// implements Database
type mongoDatabase struct {
	client   *mongo.Client
	context  context.Context
	database string
	cache    *cache.Cache
}

func (mdb *mongoDatabase) shortsCollection() *mongo.Collection {
	return mdb.client.Database(mdb.database).Collection(ShortenedCollectionName)
}

func (mdb *mongoDatabase) FindShortenedURL(id string) (res *short.ShortURL, err error) {
	// find in cache
	if s, found := mdb.cache.Get(id); found {
		return s.(*short.ShortURL), nil
	}

	filter := bson.M{
		"id": id,
	}

	cursor := mdb.shortsCollection().FindOne(mdb.context, filter)
	if err = cursor.Err(); err != nil {
		return
	}

	err = cursor.Decode(&res)

	// save to cache
	if err == nil {
		mdb.cache.Set(id, res, cache.DefaultExpiration)
	}

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

	// save to cache
	mdb.cache.Set(short.ID, &short, cache.DefaultExpiration)

	return nil
}

func (mdb *mongoDatabase) BreakCache(id string) (found bool) {
	_, found = mdb.cache.Get(id)
	mdb.cache.Delete(id)
	return
}

func (mdb *mongoDatabase) ShortURLAvailable(id string) bool {
	return shortURLAvailable(mdb, id)
}
