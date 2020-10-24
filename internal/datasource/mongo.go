package datasource

import (
	"context"

	"github.com/nkolosov/tendigma-test/internal/config"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func MustMongoDB(cfg *config.MongoConfig) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.ConnectTimeout)
	defer cancel()

	opts := options.Client()
	opts.ApplyURI(cfg.Connection)
	opts.SetConnectTimeout(cfg.ConnectTimeout)
	opts.SetSocketTimeout(cfg.SocketTimeout)
	opts.SetMinPoolSize(cfg.MinPoolSize)
	opts.SetMaxPoolSize(cfg.MaxPoolSize)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, errors.Wrapf(err, "error on mongo connection")
	}

	ctx, cancel = context.WithTimeout(context.Background(), cfg.ConnectTimeout)
	defer cancel()

	err = client.Ping(ctx, readpref.Secondary())
	if err != nil {
		return nil, errors.Wrapf(err, "error on ping database")
	}

	return client, nil
}
