package mongo

import (
	"context"

	"github.com/labasubagia/realworld-backend/internal/core/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	DBName                    = "realworld"
	CollectionUser            = "users"
	CollectionUserFollow      = "user_follows"
	CollectionTag             = "tags"
	CollectionArticle         = "articles"
	CollectionComment         = "comments"
	CollectionArticleTag      = "article_tags"
	CollectionArticleFavorite = "article_favorites"
)

type DB struct {
	client *mongo.Client
}

func NewDB(config util.Config) (DB, error) {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoSource))
	if err != nil {
		return DB{}, err
	}
	db := DB{
		client: client,
	}
	if err := db.migrate(); err != nil {
		return DB{}, err
	}
	return db, nil
}

func (db *DB) Client() *mongo.Client {
	return db.client
}

func (db *DB) Collection(name string) *mongo.Collection {
	return db.Client().Database(DBName).Collection(name)
}

func (db *DB) migrate() error {
	ctx := context.Background()

	// user index
	_, err := db.Collection(CollectionUser).Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "email", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "username", Value: 1}}, Options: options.Index().SetUnique(true)},
	})
	if err != nil {
		return err
	}

	// user follow index
	_, err = db.Collection(CollectionUserFollow).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "follower_id", Value: 1}, {Key: "followee_id", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return err
	}

	// tag index
	_, err = db.Collection(CollectionTag).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "name", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return err
	}

	// article tag index
	_, err = db.Collection(CollectionArticleTag).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "article_id", Value: 1}, {Key: "tag_id", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return err
	}

	// article favorite index
	_, err = db.Collection(CollectionArticleFavorite).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "article_id", Value: 1}, {Key: "user_id", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return err
	}

	return nil
}
