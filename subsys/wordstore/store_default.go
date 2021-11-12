package wordstore

import (
	"context"
	"fmt"

	"github.com/avrebarra/valeed"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	MongoCollecName = "users"
)

type ConfigMongo struct {
	DB               *mongo.Database `validate:"required"`
	CollecNamePrefix string          `validate:"-"`
}

type Mongo struct {
	config ConfigMongo
}

func NewMongo(cfg ConfigMongo) (Store, error) {
	if err := valeed.Validate(cfg); err != nil {
		return nil, err
	}

	e := &Mongo{config: cfg}

	return e, nil
}

func (e *Mongo) Find(ctx context.Context, in FindInput) (out []Word, total int64, err error) {
	type selector struct {
		ID string `bson:"id,omitempty"`
	}

	// prep and validate
	collecname := e.config.CollecNamePrefix + MongoCollecName
	if err = valeed.Validate(in); err != nil {
		err = fmt.Errorf("bad input: %w", err)
		return
	}

	// build data
	opts := options.Find().SetLimit(int64(in.PagingLimit)).SetSkip(int64(in.PagingOffset))
	sel := selector{ID: in.ID}

	// perform count
	count, err := e.config.DB.Collection(collecname).CountDocuments(ctx, sel)
	if err != nil {
		err = fmt.Errorf("count failure: %w", err)
		return
	}
	total = int64(count)
	if total == 0 {
		out = []Word{}
		return
	}

	// perform query
	results := []MongoUser{}
	cursor, err := e.config.DB.Collection(collecname).Find(ctx, sel, opts)
	if err != nil {
		err = fmt.Errorf("query failure: %w", err)
		return
	}
	for cursor.Next(ctx) {
		r := MongoUser{}
		if cursor.Decode(&r); err != nil {
			err = fmt.Errorf("failure decoding mongo cursor: %w", err)
			return
		}
		results = append(results, r)
	}

	// build results
	out = []Word{}
	for _, v := range results {
		out = append(out, Word(v))
	}

	return
}

func (e *Mongo) Persist(ctx context.Context, in *Word) (err error) {
	type selector struct {
		ID string `bson:"id"`
	}

	// prep and validate
	collecname := e.config.CollecNamePrefix + MongoCollecName
	if err = valeed.Validate(in); err != nil {
		err = fmt.Errorf("bad input: %w", err)
		return
	}

	// prepare data
	sel := selector{ID: in.ID}
	query := bson.D{{Key: "$set", Value: MongoUser(*in)}}
	opts := options.Update().SetUpsert(true)

	// perform persist
	_, err = e.config.DB.Collection(collecname).UpdateOne(ctx, sel, query, opts)
	if err != nil {
		err = fmt.Errorf("persist failed: %w", err)
		return
	}

	return
}

// ***

type MongoUser struct {
	ID   string `validate:"required"`
	Word string `validate:"-"`
}
