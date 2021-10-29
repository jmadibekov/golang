package httpserver

import (
	"context"
	"example/hello/project/internal/store"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"

	"encoding/json"
	"example/hello/project/internal/models"
	"fmt"
	"strconv"

	"github.com/go-chi/render"
)

type Server struct {
	ctx         context.Context
	idleConnsCh chan struct{}
	Store       store.Store

	Address string
}

func NewServer(ctx context.Context, address string, store store.Store) *Server {
	return &Server{
		ctx:         ctx,
		idleConnsCh: make(chan struct{}),
		Store:       store,

		Address: address,
	}
}

func (s *Server) basicHandler() chi.Router {
	r := chi.NewRouter()

	r.Get("/", func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("Hello world, I am Lostify!"))
	})
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})
	r.Get("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("test")
	})

	// RESTy routes for "songs" resource
	r.Route("/songs", func(r chi.Router) {
		r.Post("/", func(rw http.ResponseWriter, r *http.Request) {
			song := new(models.Song)
			if err := json.NewDecoder(r.Body).Decode(song); err != nil {
				fmt.Fprintf(rw, "Unknown err: %v", err)
				return
			}

			s.Store.Songs().Create(r.Context(), song)
		})

		r.Get("/", func(rw http.ResponseWriter, r *http.Request) {
			songs, err := s.Store.Songs().All(r.Context())
			if err != nil {
				fmt.Fprintf(rw, "Unknown err: %v", err)
				return
			}

			render.JSON(rw, r, songs)
		})

		r.Get("/{id}", func(rw http.ResponseWriter, r *http.Request) {
			idStr := chi.URLParam(r, "id")
			id, err := strconv.Atoi(idStr)
			if err != nil {
				fmt.Fprintf(rw, "Unknown err: %v", err)
				return
			}

			song, err := s.Store.Songs().ByID(r.Context(), id)
			if err != nil {
				fmt.Fprintf(rw, "Unknown err: %v", err)
				return
			}

			render.JSON(rw, r, song)
		})

		r.Put("/", func(rw http.ResponseWriter, r *http.Request) {
			song := new(models.Song)
			if err := json.NewDecoder(r.Body).Decode(song); err != nil {
				fmt.Fprintf(rw, "Unknown err: %v", err)
				return
			}

			s.Store.Songs().Update(r.Context(), song)
		})

		r.Delete("/{id}", func(rw http.ResponseWriter, r *http.Request) {
			idStr := chi.URLParam(r, "id")
			id, err := strconv.Atoi(idStr)
			if err != nil {
				fmt.Fprintf(rw, "Unknown err: %v", err)
				return
			}

			s.Store.Songs().Delete(r.Context(), id)
		})
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
