package main

import (
	"context"
	"example/hello/project/internal/grpcserver/api"
	"log"

	"google.golang.org/grpc"
)

const (
	port = ":8000"
)

func main() {
	ctx := context.Background()
	conn, err := grpc.Dial("localhost"+port, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("could not connect to %s: %v", port, err)
	}

	clientOne := api.NewSongsClient(conn)

	_, err = clientOne.CreateOrUpdateSong(ctx, &api.Song{
		Id:       300,
		Title:    "Make It Rain",
		ArtistId: 1,
	})
	if err != nil {
		log.Fatal(err)
	}
}
