package grpcserver

import (
	"context"
	"log"
	"net"

	"example/hello/project/internal/grpcserver/api"
	"example/hello/project/internal/store"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type songsServer struct {
	api.UnimplementedSongsServer
}

func (s *songsServer) CreateOrUpdateSong(context.Context, *api.Song) (*api.Song, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateOrUpdateSong not implemented")
}

func NewServer(port string, store store.Store) {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("cannot listen to %s: %v", port, err)
	}
	defer listener.Close()

	grpcServer := grpc.NewServer()

	api.RegisterSongsServer(grpcServer, new(songsServer))

	log.Printf("gRPC serving on %v", listener.Addr())
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve on %v: %v", listener.Addr(), err)
	}
}
