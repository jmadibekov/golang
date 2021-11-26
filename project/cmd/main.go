package main

import (
	"context"
	"example/hello/project/internal/grpcserver"
	"example/hello/project/internal/httpserver"
	"example/hello/project/internal/store"
	"example/hello/project/internal/store/mongodb"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("you didn't specify the type of the server: http or grpc")
	}
	serverType := os.Args[1]

	// for graceful termination in case of keyboard interrupt
	ctx, cancel := context.WithCancel(context.Background())
	go CatchTermination(cancel)

	// connecting to MongoDB store and deferring the closure of it
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

	if serverType == "http" {
		server := httpserver.NewServer(ctx,
			httpserver.WithAddress(":8080"),
			httpserver.WithStore(mongodbStore),
		)
		if err := server.Run(); err != nil {
			log.Println(err)
		}

		server.WaitForGracefulTermination()
	} else if serverType == "grpc" {
		// TODO: grpc server is not fully implemented yet
		grpcserver.NewServer(":8000", mongodbStore)
	} else {
		panic("such server type doesn't exist")
	}
}

func CatchTermination(cancel context.CancelFunc) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("[WARN] caught termination signal")
	cancel()
}
