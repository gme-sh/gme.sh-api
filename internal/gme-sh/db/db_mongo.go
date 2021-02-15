package db

import (
	"context"
	"github.com/full-stack-gods/gme.sh-api/internal/gme-sh/config"
	"time"

	"github.com/full-stack-gods/gme.sh-api/pkg/gme-sh/short"
	"github.com/patrickmn/go-cache"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// PersistentDatabase
type mongoDatabase struct {
	client             *mongo.Client
	context            context.Context
	cache              *cache.Cache
	database           string
	shortURLCollection string
}

var updateOptions = options.Update().SetUpsert(true)

// NewMongoDatabase -> Creates a new implementation of PersistentDatabase (mongodb),
// connects, and returns it
func NewMongoDatabase(cfg *config.MongoConfig) (db PersistentDatabase, err error) {
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
		cache:              cache.New(10*time.Minute, 15*time.Minute),
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
	if s, found := mdb.cache.Get(id.String()); found {
		return s.(*short.ShortURL), nil
	}
	result := mdb.shortURLs().FindOne(mdb.context, id.BsonFilter())
	if err = result.Err(); err != nil {
		return
	}
	// If the object was found in the MongoDB database, try to decode it to shortURL
	err = result.Decode(&shortURL)
	// Now we cache the object to spare the MongoDB database.
	// The object is removed from the cache after 10 minutes, or by the BreakCache method,
	// which is called among other things when the hint comes via Redis Pub-Sub
	if err == nil {
		mdb.cache.Set(id.String(), shortURL, cache.DefaultExpiration)
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
		mdb.cache.Set(short.ID.String(), short, cache.DefaultExpiration)
	}
	return
}

func (mdb *mongoDatabase) DeleteShortenedURL(id *short.ShortID) (err error) {
	// (Hopefully) deletes the object from the Mongo database
	_, err = mdb.shortURLs().DeleteOne(mdb.context, id.BsonFilter())
	// Also remove the object from the cache
	if err == nil {
		// remove from cache
		mdb.BreakCache(id)
	}
	return
}

/*
 * ==================================================================================================
 *                          D E F A U L T   I M P L E M E N T A T I O N S
 * ==================================================================================================
 */

func (mdb *mongoDatabase) BreakCache(id *short.ShortID) (found bool) {
	_, found = mdb.cache.Get(id.String())
	mdb.cache.Delete(id.String())
	return
}

func (mdb *mongoDatabase) ShortURLAvailable(id *short.ShortID) bool {
	if _, found := mdb.cache.Get(id.String()); found {
		return false
	}
	return shortURLAvailable(mdb, id)
}
