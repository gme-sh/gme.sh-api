package db

import (
	"context"
	"github.com/gme-sh/gme.sh-api/internal/gme-sh/config"
	"github.com/gme-sh/gme.sh-api/pkg/gme-sh/short"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

// PersistentDatabase
type mongoDatabase struct {
	client             *mongo.Client
	context            context.Context
	cache              DBCache
	database           string
	shortURLCollection string
}

var updateOptions = options.Update().SetUpsert(true)

// NewMongoDatabase -> Creates a new implementation of PersistentDatabase (mongodb),
// connects, and returns it
func NewMongoDatabase(cfg *config.MongoConfig, cache DBCache) (db PersistentDatabase, err error) {
	// create client
	opts := options.Client().ApplyURI(cfg.ApplyURI)
	client, err := mongo.NewClient(opts)
	if err != nil {
		return nil, err
	}

	// connect to client
	ctx := context.TODO()
	err = client.Connect(ctx)
	if err != nil {
		return
	}

	return &mongoDatabase{
		client:             client,
		context:            ctx,
		database:           cfg.Database,
		shortURLCollection: cfg.ShortURLCollection,
		cache:              cache,
	}, nil
}

func (mdb *mongoDatabase) shortURLs() *mongo.Collection {
	return mdb.client.Database(mdb.database).Collection(mdb.shortURLCollection)
}

/*
 * ==================================================================================================
 *                            P E R M A N E N T  D A T A B A S E
 * ==================================================================================================
 */

func (mdb *mongoDatabase) FindShortenedURL(id *short.ShortID) (shortURL *short.ShortURL, err error) {
	// At first, try to load the object from the cache
	if u := mdb.cache.GetShortURL(id); u != nil {
		return u, nil
	}
	result := mdb.shortURLs().FindOne(mdb.context, id.BsonFilter())
	log.Println("-> not in cache. Result with filter", id.BsonFilter(), " (error) :", result.Err())
	if err = result.Err(); err != nil {
		return
	}
	// If the object was found in the MongoDB database, try to decode it to shortURL
	err = result.Decode(&shortURL)

	// Now we cache the object to spare the MongoDB database.
	// The object is removed from the cache after 10 minutes, or by the BreakCache method,
	// which is called among other things when the hint comes via Redis Pub-Sub
	if err == nil {
		err = mdb.cache.UpdateCache(shortURL)
	}

	return
}

func (mdb *mongoDatabase) SaveShortenedURL(short *short.ShortURL) (err error) {
	_, err = mdb.shortURLs().UpdateOne(
		mdb.context,
		short.ID.BsonFilter(),
		short.BsonUpdate(),
		updateOptions,
	)
	// Now we save/replace the object to the cache to spare the Mongo database.
	if err == nil {
		err = mdb.cache.UpdateCache(short)
	}
	return
}

func (mdb *mongoDatabase) DeleteShortenedURL(id *short.ShortID) (err error) {
	// (Hopefully) deletes the object from the Mongo database
	_, err = mdb.shortURLs().DeleteOne(mdb.context, id.BsonFilter())
	// Also remove the object from the cache
	if err == nil {
		// remove from cache
		err = mdb.cache.BreakCache(id)
	}
	return
}

/*
 * ==================================================================================================
 *                          D E F A U L T   I M P L E M E N T A T I O N S
 * ==================================================================================================
 */
func (mdb *mongoDatabase) ShortURLAvailable(id *short.ShortID) bool {
	if u := mdb.cache.GetShortURL(id); u != nil {
		return false
	}
	return shortURLAvailable(mdb, id)
}
