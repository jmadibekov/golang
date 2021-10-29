package main

import (
	"context"
	"example/hello/project/internal/httpserver"
	"example/hello/project/internal/store/inmemory"
	"log"
)

func main() {
	store := inmemory.NewDB()

	server := httpserver.NewServer(context.Background(), ":8080", store)
	if err := server.Run(); err != nil {
		log.Println(err)
	}

	server.WaitForGracefulTermination()
}
