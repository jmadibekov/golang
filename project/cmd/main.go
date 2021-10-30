package main

import (
	"context"
	"example/hello/project/internal/grpcserver"
	"example/hello/project/internal/httpserver"
	"example/hello/project/internal/store/inmemory"
	"log"
	"os"
)

func main() {
	serverType := os.Args[1]
	store := inmemory.NewDB()

	if serverType == "http" {
		server := httpserver.NewServer(context.Background(), ":8080", store)
		if err := server.Run(); err != nil {
			log.Println(err)
		}

		server.WaitForGracefulTermination()
	} else if serverType == "grpc" {
		grpcserver.NewServer(":8000", store)
	} else {
		panic("test")
	}
}
