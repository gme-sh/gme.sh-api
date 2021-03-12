package db

import (
	"context"
	"github.com/gme-sh/gme.sh-api/internal/gme-sh/config"
	"github.com/gme-sh/gme.sh-api/pkg/gme-sh/short"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// PersistentDatabase
type mongoDatabase struct {
	client             *mongo.Client
	context            context.Context
	cache              DBCache
	database           string
	shortURLCollection string
	metaCollection     string
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
		metaCollection:     cfg.MetaCollection,
		cache:              cache,
	}, nil
}

func (*mongoDatabase) ServiceName() string {
	return "MongoDB"
}

func (mdb *mongoDatabase) HealthCheck(ctx context.Context) (err error) {
	err = mdb.client.Ping(ctx, nil)
	return
}

////

func (mdb *mongoDatabase) shortURLs() *mongo.Collection {
	return mdb.client.Database(mdb.database).Collection(mdb.shortURLCollection)
}

func (mdb *mongoDatabase) meta() *mongo.Collection {
	return mdb.client.Database(mdb.database).Collection(mdb.metaCollection)
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

func (mdb *mongoDatabase) FindExpiredURLs() (res []*short.ShortURL, err error) {
	filter := bson.M{
		"$and": []bson.M{
			{"expiration_date": bson.M{"$exists": true}},
			{"expiration_date": bson.M{"$ne": nil}},
			{"expiration_date": bson.M{"$lt": time.Now()}},
		},
	}

	var cursor *mongo.Cursor
	cursor, err = mdb.shortURLs().Find(mdb.context, filter)
	if err != nil {
		return
	}

	for cursor.Next(mdb.context) {
		var u *short.ShortURL
		if err := cursor.Decode(&u); err != nil {
			return nil, err
		}
		res = append(res, u)
	}

	return
}

func (mdb *mongoDatabase) GetLastExpirationCheck() (m *LastExpirationCheckMeta) {
	m = &LastExpirationCheckMeta{LastCheck: time.Unix(5, 0)}
	res := mdb.meta().FindOne(mdb.context, bson.M{"key": "last_expiration"})
	if res.Err() != nil {
		m.LastCheck = time.Unix(4, 0)
	} else {
		if err := res.Decode(m); err != nil {
			m.LastCheck = time.Unix(3, 0)
		}
	}
	return
}

func (mdb *mongoDatabase) UpdateLastExpirationCheck(t time.Time) {
	m := &LastExpirationCheckMeta{LastCheck: t}
	_, _ = mdb.meta().UpdateOne(
		mdb.context,
		bson.M{"key": "last_expiration"},
		bson.M{"$set": m},
		options.Update().SetUpsert(true),
	)
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
