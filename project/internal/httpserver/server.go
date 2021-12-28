package httpserver

import (
	"context"
	"example/hello/project/internal/message_broker"
	"example/hello/project/internal/store"
	"github.com/go-chi/chi/middleware"
	lru "github.com/hashicorp/golang-lru"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

type Server struct {
	ctx         context.Context
	idleConnsCh chan struct{}
	store       store.Store
	cache       *lru.TwoQueueCache
	broker      message_broker.MessageBroker

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

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

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
	songsResource := NewSongResource(s.store, s.broker, s.cache)
	r.Mount("/songs", songsResource.Routes())

	// mounting routes of /artists resource
	artistsResource := NewArtistResource(s.store, s.broker, s.cache)
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

	log.Printf("[HTTP] Server running on %v", s.Address)
	return server.ListenAndServe()
}

func (s *Server) ListenCtxForGT(server *http.Server) {
	<-s.ctx.Done() // blocks while context isn't finished

	if err := server.Shutdown(context.Background()); err != nil {
		log.Println("[HTTP] Got error while shutting down", err)
	}

	log.Println("[HTTP] Processed all idle connections")
	close(s.idleConnsCh)
}

func (s *Server) WaitForGracefulTermination() {
	<-s.idleConnsCh
}
