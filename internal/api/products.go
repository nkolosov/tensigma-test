package api

import (
	"context"

	"github.com/nkolosov/tendigma-test/internal/csv"

	"github.com/nkolosov/tendigma-test/internal/datasource"
	"google.golang.org/grpc/grpclog"
)

type Products struct {
	ds *datasource.Products

	pipeline *csv.Pipeline
}

func NewProductsAPI(ds *datasource.Products, pipeline *csv.Pipeline) *Products {
	return &Products{
		ds:       ds,
		pipeline: pipeline,
	}
}

func (api *Products) Fetch(ctx context.Context, req *FetchRequest) (*FetchResponse, error) {
	grpclog.Infof("fetch request %+v\n", req)

	api.pipeline.Handle(req.Url)

	return &FetchResponse{}, nil
}

func (api *Products) List(ctx context.Context, req *ListRequest) (*ListResponse, error) {
	grpclog.Infof("list request %+v\n", req)

	return &ListResponse{}, nil
}
