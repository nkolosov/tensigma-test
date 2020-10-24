package main

import (
	"context"
	"fmt"

	"github.com/nkolosov/tendigma-test/internal/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

func main() {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	conn, err := grpc.Dial("127.0.0.1:9090", opts...)

	if err != nil {
		grpclog.Fatalf("failed to dial: %+v", err)
	}

	defer func() {
		err := conn.Close()
		if err != nil {
			fmt.Printf("error on connection closing %+v", err)
		}
	}()

	client := api.NewProductsAPIClient(conn)
	request := &api.FetchRequest{
		Url: "http://localhost:4000/data.csv",
	}

	response, err := client.Fetch(context.Background(), request)
	grpclog.Warningf("test response %+v\n", response)
	if err != nil {
		grpclog.Fatalf("fail to fetch: %+v %+v", response, err)
	}

	fmt.Printf("result %+v", response)
}
