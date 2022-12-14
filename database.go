package mongodbhelpers

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Database struct {
	client *mongo.Client
	ctx    context.Context
}

func NewDatabaseConnection(mongoURI string) (*Database, error) {
	ctx := context.Background()
	client, error := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	return &Database{
		client: client,
		ctx:    ctx,
	}, error
}

func (database *Database) IsConnected() bool {
	if err := database.client.Ping(database.ctx, readpref.Primary()); err != nil {
		return false
	}
	return true
}

func (database *Database) GetContext() context.Context {
	return database.ctx
}

func (database *Database) GetNextSequenceValue(databaseName string, collectionName string, sequenceName string) (*mongo.UpdateResult, error) {
	result, err := database.client.Database(databaseName).Collection(collectionName).UpdateOne(database.ctx, bson.M{
		"_id": sequenceName,
	}, bson.M{
		"$inc": bson.M{
			"seq": 1,
		},
	}, options.Update().SetUpsert(true))
	if result.MatchedCount == 0 {
		return nil, err
	}
	return result, nil
}

func (database *Database) InsertOne(databaseName string, collectioName string, data interface{}) (*mongo.InsertOneResult, error) {
	result, err := database.client.Database(databaseName).Collection(collectioName).InsertOne(database.ctx, data, options.InsertOne().SetBypassDocumentValidation(false))
	return result, err
}

func (database *Database) InsertMany(databaseName string, collectionName string, data []interface{}) (*mongo.InsertManyResult, error) {
	results, err := database.client.Database(databaseName).Collection(collectionName).InsertMany(database.ctx, data, options.InsertMany().SetOrdered(true))
	return results, err
}

func (database *Database) UpdateOne(databaseName string, collectioName string, filter interface{}, data interface{}) (*mongo.UpdateResult, error) {
	results, err := database.client.Database(databaseName).Collection(collectioName).UpdateOne(database.ctx, filter, bson.M{
		"$set": data,
	})
	return results, err
}

func (database *Database) UpdateMany(databaseName string, collectioName string, filter interface{}, data interface{}) (interface{}, error) {
	results, err := database.client.Database(databaseName).Collection(collectioName).UpdateMany(database.ctx, filter, bson.M{
		"$set": data,
	})
	return results.UpsertedID, err
}

func (database *Database) DocumentsCount(databaseName string, collectionName string, filter interface{}) (int64, error) {
	result, err := database.client.Database(databaseName).Collection(collectionName).CountDocuments(database.ctx, filter)
	return result, err
}

func (database *Database) FindFirst(databaseName string, collectionName string, filter interface{}) (*mongo.Cursor, error) {
	result, err := database.client.Database(databaseName).Collection(collectionName).Find(database.ctx, filter, options.Find().SetSort(
		bson.M{
			"_id": 1,
		},
	).SetLimit(1).SetMax(1))
	return result, err
}

func (database *Database) FindLast(databaseName string, collectionName string, filter interface{}) (*mongo.Cursor, error) {
	result, err := database.client.Database(databaseName).Collection(collectionName).Find(database.ctx, filter, options.Find().SetSort(
		bson.M{
			"_id": -1,
		},
	).SetLimit(1).SetMax(1))
	return result, err
}

func (database *Database) FindOne(databaseName string, collectionName string, filter interface{}, fields interface{}) (*mongo.SingleResult, error) {
	result := database.client.Database(databaseName).Collection(collectionName).FindOne(database.ctx, filter, options.FindOne().SetProjection(fields))
	if result.Err() != nil {
		return nil, result.Err()
	}
	return result, nil
}

func (database *Database) FindSorted(databaseName string, collectionName string, filter interface{}, fields interface{}, sortBy string, order int64) (*mongo.Cursor, error) {
	result, err := database.client.Database(databaseName).Collection(collectionName).Find(database.ctx, filter, options.Find().SetSort(
		bson.M{
			sortBy: order,
		},
	))

	return result, err
}

func (database *Database) Find(databaseName string, collectionName string, filter interface{}, fields interface{}) (*mongo.Cursor, error) {
	result, err := database.client.Database(databaseName).Collection(collectionName).Find(database.ctx, filter, options.Find().SetProjection(fields).SetAllowDiskUse(true))
	return result, err
}

func (database *Database) FindPaginated(databaseName string, collectionName string, filter interface{}, fields interface{}, skip int64, sort interface{}, limit int64) (*mongo.Cursor, error) {
	result, err := database.client.Database(databaseName).Collection(collectionName).Find(database.ctx, filter, options.Find().SetSkip(skip).SetSort(sort).SetLimit(limit).SetProjection(fields).SetAllowDiskUse(true))
	return result, err
}

func (database *Database) Distinct(databaseName string, collectionName string, fieldName string, filter interface{}) ([]interface{}, error) {
	result, err := database.client.Database(databaseName).Collection(collectionName).Distinct(database.ctx, fieldName, filter, options.Distinct())

	return result, err
}

func (database *Database) CustomAggregate(databaseName string, collectionName string, filter []bson.M) (*mongo.Cursor, error) {
	result, err := database.client.Database(databaseName).Collection(collectionName).Aggregate(database.ctx, filter, options.Aggregate())
	return result, err
}

func (database *Database) DeleteOne(databaseName string, collectionName string, filter interface{}) (interface{}, error) {
	result, err := database.client.Database(databaseName).Collection(collectionName).DeleteOne(database.ctx, filter)
	return result.DeletedCount, err
}

func (database *Database) DeleteMany(databaseName string, collectionName string, filter interface{}) (interface{}, error) {
	result, err := database.client.Database(databaseName).Collection(collectionName).DeleteMany(database.ctx, filter)
	return result.DeletedCount, err
}

func (database *Database) UpdateOneWithUpsert(databaseName string, collectionName string, filter interface{}, data interface{}, opts *options.UpdateOptions) (interface{}, error) {
	result, err := database.client.Database(databaseName).Collection(collectionName).UpdateOne(database.ctx, filter, bson.M{
		"$set": data,
	}, opts)
	return result, err
}
