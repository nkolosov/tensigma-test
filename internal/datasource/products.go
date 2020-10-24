package datasource

import (
	"fmt"
	"strconv"
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

type Products struct {
	client *mongo.Client
}

func NewProducts(client *mongo.Client) *Products {
	return &Products{
		client: client,
	}
}

func (p *Products) Update(model *Product) {
	grpclog.Infof("update model %+v\n", model)
}
