package main

import (
	"context"
	"example/hello/project/internal/grpcserver"
	"example/hello/project/internal/httpserver"
	"example/hello/project/internal/store/mongodb"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("you didn't specify the type of the server: http or grpc")
	}

	serverType := os.Args[1]

	store := mongodb.NewDB()
	mongodbURI := "mongodb+srv://dbUser:dbUserPassword@cluster0.99ocj.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"

	if err := store.Connect(mongodbURI); err != nil {
		panic(err)
	}
	defer store.Close()

	if serverType == "http" {
		server := httpserver.NewServer(context.Background(), ":8080", store)
		if err := server.Run(); err != nil {
			log.Println(err)
		}

		server.WaitForGracefulTermination()
	} else if serverType == "grpc" {
		grpcserver.NewServer(":8000", store)
	} else {
		panic("such server type doesn't exist")
	}
}
