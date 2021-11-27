package httpserver

import (
	"encoding/json"
	"example/hello/project/internal/models"
	"example/hello/project/internal/store"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"net/http"
	"strconv"
)

type ArtistResource struct {
	store store.Store
}

func NewArtistResource(store store.Store) *ArtistResource {
	return &ArtistResource{
		store: store,
	}
}

func (ar *ArtistResource) Routes() chi.Router {
	r := chi.NewRouter()

	// RESTy routes for "artists" resource
	r.Post("/", ar.CreateArtist)
	r.Get("/", ar.AllArtists)
	r.Get("/{id}", ar.ByID)
	r.Put("/", ar.UpdateArtist)
	r.Delete("/{id}", ar.DeleteArtist)

	return r
}

func (ar *ArtistResource) CreateArtist(rw http.ResponseWriter, r *http.Request) {
	artist := new(models.Artist)
	if err := json.NewDecoder(r.Body).Decode(artist); err != nil {
		rw.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = fmt.Fprintf(rw, "Unknown err: %v", err)
		return
	}

	if err := ar.store.Artists().Create(r.Context(), artist); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(rw, "DB err: %v", err)
		return
	}

	rw.WriteHeader(http.StatusCreated)
}

func (ar *ArtistResource) AllArtists(rw http.ResponseWriter, r *http.Request) {
	artists, err := ar.store.Artists().All(r.Context())
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(rw, "DB err: %v", err)
		return
	}

	render.JSON(rw, r, artists)
}

func (ar *ArtistResource) ByID(rw http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(rw, "Unknown err: %v", err)
		return
	}

	artist, err := ar.store.Artists().ByID(r.Context(), id)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(rw, "DB err: %v", err)
		return
	}

	render.JSON(rw, r, artist)
}

func (ar *ArtistResource) UpdateArtist(rw http.ResponseWriter, r *http.Request) {
	artist := new(models.Artist)
	if err := json.NewDecoder(r.Body).Decode(artist); err != nil {
		rw.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = fmt.Fprintf(rw, "Unknown err: %v", err)
		return
	}

	err := validation.ValidateStruct(
		artist,
		validation.Field(&artist.ID, validation.Required),
		validation.Field(&artist.FullName, validation.Required),
	)
	if err != nil {
		rw.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = fmt.Fprintf(rw, "Unknown err: %v", err)
		return
	}

	if err := ar.store.Artists().Update(r.Context(), artist); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(rw, "DB err: %v", err)
		return
	}
}

func (ar *ArtistResource) DeleteArtist(rw http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(rw, "Unknown err: %v", err)
		return
	}

	if err := ar.store.Artists().Delete(r.Context(), id); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(rw, "DB err: %v", err)
		return
	}
}
