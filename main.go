package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	handler "cart/handler"
	productProto "cart/product"
	pb "cart/proto"

	"github.com/micro/go-micro/v2/logger"
	"github.com/micro/micro/v3/service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateClient(ctx context.Context, uri string, retry int32) (*mongo.Client, error) {
	conn, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err := conn.Ping(ctx, nil); err != nil {
		if retry >= 3 {
			fmt.Printf("Failed connect")
			return nil, err
		}
		retry = retry + 1
		time.Sleep(time.Second * 2)
		return CreateClient(ctx, uri, retry)
	}
	fmt.Printf("Connect with database")
	return conn, err
}

func main() {
	uri := os.Getenv("DB_HOST")

	srv := service.New(
		service.Name("cart"),
		service.Version("latest"),
	)

	srv.Init()

	client, err := CreateClient(context.Background(), uri, 0)
	if err != nil {
		log.Panic(err)
	}
	defer client.Disconnect(context.Background())

	cartCollection := client.Database("cart").Collection("cart")
	productClient := productProto.NewProductService("product", srv.Client())
	repo := &handler.MongoRepository{
		Collection: cartCollection,
	}

	h := &handler.Handler{
		Repository:    repo,
		ProductClient: productClient,
	}

	pb.RegisterCartHandler(srv.Server(), h)

	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
