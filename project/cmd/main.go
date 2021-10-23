package main

import (
	"context"
	"example/hello/project/internal/http"
	"example/hello/project/internal/store/inmemory"
	"log"
)

func main() {
	store := inmemory.NewDB()

	server := http.NewServer(context.Background(), ":8080", store)
	if err := server.Run(); err != nil {
		log.Println(err)
	}

	server.WaitForGracefulTermination()
}
