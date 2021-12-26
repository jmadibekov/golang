package main

import (
	"context"
	"example/hello/project/internal/grpcserver"
	"example/hello/project/internal/httpserver"
	"example/hello/project/internal/message_broker"
	"example/hello/project/internal/message_broker/kafka"
	"example/hello/project/internal/store"
	"example/hello/project/internal/store/mongodb"
	"fmt"
	lru "github.com/hashicorp/golang-lru"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	serverType := os.Getenv("SERVER_TYPE")
	if serverType == "" {
		serverType = "http"
	}

	// for graceful termination in case of keyboard interrupt
	ctx, cancel := context.WithCancel(context.Background())
	go CatchTermination(cancel)

	// connecting to MongoDB store and deferring the closure of it
	// TODO: move hardcoded username & pw values to docker-compose.yml or .env file
	mongodbURI := "mongodb+srv://dbUser:dbUserPassword@cluster0.99ocj.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"
	mongodbStore := mongodb.NewDB()
	if err := mongodbStore.Connect(mongodbURI); err != nil {
		panic(err)
	}
	defer func(mongodbStore store.Store) {
		err := mongodbStore.Close()
		if err != nil {
			panic(err)
		}
	}(mongodbStore)

	// creating in-memory cache
	cache, err := lru.New2Q(6)
	if err != nil {
		panic(err)
	}

	// connecting to Kafka brokers
	brokers := []string{"localhost:9092"}
	broker := kafka.NewBroker(brokers, cache, "peer0")
	if err := broker.Connect(ctx); err != nil {
		panic(err)
	}
	defer func(broker message_broker.MessageBroker) {
		err := broker.Close()
		if err != nil {
			panic(err)
		}
	}(broker)

	if serverType == "http" {
		server := httpserver.NewServer(ctx,
			httpserver.WithAddress(":8080"),
			httpserver.WithStore(mongodbStore),
			httpserver.WithCache(cache),
			httpserver.WithBroker(broker),
		)
		if err := server.Run(); err != nil {
			log.Println(err)
		}

		server.WaitForGracefulTermination()
	} else if serverType == "grpc" {
		// TODO: grpc server is not fully implemented yet
		grpcserver.NewServer(":8000", mongodbStore)
	} else {
		panic(fmt.Sprintf("server type %v doesn't exist", serverType))
	}
}

func CatchTermination(cancel context.CancelFunc) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("[WARN] caught termination signal")
	cancel()
}
