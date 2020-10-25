package datasource

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/grpclog"
)

type Product struct {
	id               string    `bson:"_id"`
	name             string    `bson:"name"`
	price            uint64    `bson:"price"`
	lastPriceUpdate  time.Time `bson:"lastPriceUpdate"`
	priceUpdateCount uint64    `bson:"priceUpdateCount"`
}

func CreateProductFromCSV(columns []string) (*Product, error) {
	if len(columns) != 2 {
		return nil, fmt.Errorf("invalid CSV row %#v", columns)
	}

	name := columns[0]
	price, _ := strconv.Atoi(columns[1])

	return &Product{
		name:             name,
		price:            uint64(price),
		lastPriceUpdate:  time.Now(),
		priceUpdateCount: 1,
	}, nil
}

const (
	database   = "tendigma"
	collection = "products"

	cursorTimeout = 10 * time.Second
)

type Products struct {
	client *mongo.Client
}

func NewProducts(client *mongo.Client) (*Products, error) {
	products := &Products{
		client: client,
	}

	err := products.init()
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (p *Products) Update(model *Product) error {
	grpclog.Infof("update product model %+v\n", model)

	opts := options.Update().SetUpsert(true)

	filter := bson.D{{"_id", model.id}}
	update := bson.D{{"$set", model}}

	ctx, cancel := context.WithTimeout(context.Background(), cursorTimeout)
	defer cancel()

	result, err := p.client.Database(database).Collection(collection).UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return errors.Wrapf(err, "can't update products")
	}

	fmt.Printf("result %+v", result)
	return nil
}

func (p *Products) init() error {
	indexes, err := p.client.Database(database).Collection(collection).Indexes().List(context.Background())
	if err != nil {
		return errors.Wrapf(err, "can't get indexes for collection %s", collection)
	}

	var index bson.D

	isIndexCreated := false

	for indexes.Next(context.Background()) {
		err = indexes.Decode(&index)
		if err != nil {
			return errors.Wrapf(err, "can't decode index object")
		}

		if index.Map()["name"] == "product_name_1" {
			isIndexCreated = true
			break
		}
	}

	if isIndexCreated {
		return nil
	}

	return p.createIndexes()
}

func (p *Products) createIndexes() error {
	c := p.client.Database(database).Collection(collection)
	opts := options.CreateIndexes()

	indexes := []mongo.IndexModel{
		{
			Keys: bson.M{"name": 1},
		},
		{
			Keys: bson.M{"price": 1},
		},
		{
			Keys: bson.M{"lastPriceUpdate": 1},
		},
	}

	_, err := c.Indexes().CreateMany(context.Background(), indexes, opts)
	if err != nil {
		return errors.Wrapf(err, "can't create indexes")
	}

	return nil
}
