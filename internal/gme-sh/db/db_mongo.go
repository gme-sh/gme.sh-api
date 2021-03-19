package db

import (
	"context"
	"github.com/gme-sh/gme.sh-api/internal/gme-sh/config"
	"github.com/gme-sh/gme.sh-api/pkg/gme-sh/short"
	"github.com/gme-sh/gme.sh-api/pkg/gme-sh/tpl"
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
	tplCollection      string
	poolCollection     string
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
		tplCollection:      cfg.TplCollection,
		poolCollection:     cfg.PoolCollection,
		cache:              cache,
	}, nil
}

////

func (mdb *mongoDatabase) shortURLs() *mongo.Collection {
	return mdb.client.Database(mdb.database).Collection(mdb.shortURLCollection)
}

func (mdb *mongoDatabase) meta() *mongo.Collection {
	return mdb.client.Database(mdb.database).Collection(mdb.metaCollection)
}

func (mdb *mongoDatabase) tpl() *mongo.Collection {
	return mdb.client.Database(mdb.database).Collection(mdb.tplCollection)
}
func (mdb *mongoDatabase) pool() *mongo.Collection {
	return mdb.client.Database(mdb.database).Collection(mdb.poolCollection)
}

/*
 * ==================================================================================================
 *                          D E F A U L T   I M P L E M E N T A T I O N S
 * ==================================================================================================
 */

func (*mongoDatabase) ServiceName() string {
	return "MongoDB"
}

func (mdb *mongoDatabase) HealthCheck(ctx context.Context) (err error) {
	err = mdb.client.Ping(ctx, nil)
	return
}

/*
 * ==================================================================================================
 *                            P E R M A N E N T  D A T A B A S E
 * ==================================================================================================
 */

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

func (mdb *mongoDatabase) ShortURLAvailable(id *short.ShortID) bool {
	if u := mdb.cache.GetShortURL(id); u != nil {
		return false
	}
	return shortURLAvailable(mdb, id)
}

/*
 * ==================================================================================================
 *                          E X P I R A T I O N   I M P L E M E N T A T I O N S
 * ==================================================================================================
 */

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
 *                          T E M P L A T E   I M P L E M E N T A T I O N S
 * ==================================================================================================
 */

func (mdb *mongoDatabase) FindTemplates() (templates []*tpl.Template, err error) {
	filter := bson.M{}

	var cursor *mongo.Cursor
	cursor, err = mdb.tpl().Find(mdb.context, filter)
	if err != nil {
		return
	}

	for cursor.Next(mdb.context) {
		var f *tpl.Template
		err = cursor.Decode(&f)
		if err != nil {
			return
		}
		templates = append(templates, f)
	}
	return
}

func (mdb *mongoDatabase) SaveTemplate(t *tpl.Template) (err error) {
	filter := bson.M{
		"template_url": t.TemplateURL,
	}
	update := bson.M{
		"$set": t,
	}
	_, err = mdb.tpl().UpdateOne(mdb.context,
		filter,
		update,
		options.Update().SetUpsert(true))
	return
}

/*
 * ==================================================================================================
 *                             P O O L   I M P L E M E N T A T I O N S
 * ==================================================================================================
 */

func (mdb *mongoDatabase) FindPool(id *short.PoolID) (pool *short.Pool, err error) {
	filter := bson.M{
		"pool_id": id.String(),
	}
	cursor := mdb.pool().FindOne(mdb.context, filter)
	if err = cursor.Err(); err != nil {
		return
	}
	pool = new(short.Pool)
	err = cursor.Decode(pool)
	return
}

func (mdb *mongoDatabase) SavePool(pool *short.Pool) (err error) {
	filter := bson.M{
		"pool_id": pool.ID.String(),
	}
	update := bson.M{
		"$set": pool,
	}
	_, err = mdb.pool().UpdateOne(mdb.context,
		filter,
		update,
		options.Update().SetUpsert(true))
	return
}
