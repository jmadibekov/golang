package http

import (
	"context"
	"encoding/json"
	"example/hello/project/internal/models"
	"example/hello/project/internal/store"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type Server struct {
	ctx         context.Context
	idleConnsCh chan struct{}
	store       store.Store

	Address string
}

func NewServer(ctx context.Context, address string, store store.Store) *Server {
	return &Server{
		ctx:         ctx,
		idleConnsCh: make(chan struct{}),
		store:       store,

		Address: address,
	}
}

func (s *Server) basicHandler() chi.Router {
	r := chi.NewRouter()

	r.Post("/songs", func(rw http.ResponseWriter, r *http.Request) {
		song := new(models.Song)
		if err := json.NewDecoder(r.Body).Decode(song); err != nil {
			fmt.Fprintf(rw, "Unknown err: %v", err)
			return
		}

		s.store.Songs().Create(r.Context(), song)
	})

	r.Get("/songs", func(rw http.ResponseWriter, r *http.Request) {
		songs, err := s.store.Songs().All(r.Context())
		if err != nil {
			fmt.Fprintf(rw, "Unknown err: %v", err)
			return
		}

		render.JSON(rw, r, songs)
	})

	r.Get("/songs/{id}", func(rw http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			fmt.Fprintf(rw, "Unknown err: %v", err)
			return
		}

		song, err := s.store.Songs().ByID(r.Context(), id)
		if err != nil {
			fmt.Fprintf(rw, "Unknown err: %v", err)
			return
		}

		render.JSON(rw, r, song)
	})

	r.Put("/songs", func(rw http.ResponseWriter, r *http.Request) {
		song := new(models.Song)
		if err := json.NewDecoder(r.Body).Decode(song); err != nil {
			fmt.Fprintf(rw, "Unknown err: %v", err)
			return
		}

		s.store.Songs().Update(r.Context(), song)
	})

	r.Delete("/songs/{id}", func(rw http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			fmt.Fprintf(rw, "Unknown err: %v", err)
			return
		}

		s.store.Songs().Delete(r.Context(), id)
	})

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

	log.Println("Server running on", s.Address)
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
