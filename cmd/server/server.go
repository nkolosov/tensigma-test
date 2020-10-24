package main

import (
	"context"
	"github.com/nkolosov/tendigma-test/internal/csv"
	"go.mongodb.org/mongo-driver/mongo"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"google.golang.org/grpc/grpclog"

	"github.com/nkolosov/tendigma-test/internal/api"
	"github.com/nkolosov/tendigma-test/internal/config"
	"github.com/nkolosov/tendigma-test/internal/datasource"
	"google.golang.org/grpc"
)

func main() {
	grpclog.Infoln("start products-api server")
	cfg := config.MustConfigure()

	grpclog.Infof("configs %+v", cfg)

	client, err := datasource.MustMongoDB(&cfg.Mongo)
	if err != nil {
		grpclog.Fatalf("failed to connect to MongoDB %s\n", err.Error())
	}

	ds := datasource.NewProducts(client)
	pipeline := csv.NewPipeline(cfg.Pipeline, ds)

	productsAPI := api.NewProductsAPI(ds, pipeline)

	hostPort := net.JoinHostPort(cfg.GRPC.Host, cfg.GRPC.Port)
	listener, err := net.Listen("tcp", hostPort)
	if err != nil {
		grpclog.Fatalf("failed to listen on %s with err %s\n", hostPort, err.Error())
	}

	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	defer func() {
		if r := recover(); r != nil {
			grpclog.Warningf("app crashed & recovered with: %#v\n", r)

			terminateApp(grpcServer, pipeline, client)
		}
	}()

	api.RegisterProductsAPIServer(grpcServer, productsAPI)

	if err = grpcServer.Serve(listener); err != nil {
		grpclog.Fatalf("failed to serve for listener %+v\n", listener)
	}

	mainWG := &sync.WaitGroup{}
	mainWG.Add(1)
	go func() {
		defer mainWG.Done()

		termChannel := make(chan os.Signal, 1)
		signal.Notify(termChannel, syscall.SIGTERM, syscall.SIGINT)
		<-termChannel

		terminateApp(grpcServer, pipeline, client)
	}()

	mainWG.Wait()
	grpclog.Infoln("service stopped")
}

func terminateApp(
	grpcServer *grpc.Server,
	pipeline *csv.Pipeline,
	client *mongo.Client,
) {
	grpcServer.GracefulStop()

	var err error

	if err = pipeline.Close(); err != nil {
		grpclog.Errorf("can't close pipeline with error %+v\n", err)
	}

	if err = client.Disconnect(context.Background()); err != nil {
		grpclog.Errorf("can't close datasource connection with error %+v\n", err)
	}
}
