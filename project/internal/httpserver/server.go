package httpserver

import (
	"context"
	"example/hello/project/internal/store"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

type Server struct {
	ctx         context.Context
	idleConnsCh chan struct{}
	store       store.Store

	Address string
}

func NewServer(ctx context.Context, opts ...ServerOption) *Server {
	srv := &Server{
		ctx:         ctx,
		idleConnsCh: make(chan struct{}),
	}

	for _, opt := range opts {
		opt(srv)
	}

	return srv
}

func (s *Server) basicHandler() chi.Router {
	r := chi.NewRouter()

	r.Get("/", func(rw http.ResponseWriter, r *http.Request) {
		_, err := rw.Write([]byte("Hello world, I am Lostify!"))
		if err != nil {
			return
		}
	})
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("pong"))
		if err != nil {
			return
		}
	})
	r.Get("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("test")
	})

	// mounting routes of /songs resource
	songsResource := NewSongResource(s.store)
	r.Mount("/songs", songsResource.Routes())

	// mounting routes of /artists resource
	artistsResource := NewArtistResource(s.store)
	r.Mount("/artists", artistsResource.Routes())

	return r
}

func (s *Server) Run() error {
	server := &http.Server{
		Addr:         s.Address,
		Handler:      s.basicHandler(),
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 30,
	}
	go s.ListenCtxForGT(server)

	log.Printf("http server running on %v", s.Address)
	return server.ListenAndServe()
}

func (s *Server) ListenCtxForGT(server *http.Server) {
	<-s.ctx.Done() // blocks while context isn't finished

	if err := server.Shutdown(context.Background()); err != nil {
		log.Println("Got error while shutting down", err)
	}

	log.Println("Processed all idle connections")
	close(s.idleConnsCh)
}

func (s *Server) WaitForGracefulTermination() {
	<-s.idleConnsCh
}
