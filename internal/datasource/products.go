package datasource

import (
	"time"

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

func NewProductModel(name string, price uint64) *Product {
	return &Product{
		name:             name,
		price:            price,
		lastPriceUpdate:  time.Now(),
		priceUpdateCount: 1,
	}
}

type Products struct {
	client *mongo.Client
}

func NewProducts(client *mongo.Client) *Products {
	return &Products{
		client: client,
	}
}

func (p *Products) Update(model *Product) {
	grpclog.Infof("update model %+v", model)
}
