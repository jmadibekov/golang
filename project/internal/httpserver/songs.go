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

type SongResource struct {
	store store.Store
}

func NewSongResource(store store.Store) *SongResource {
	return &SongResource{
		store: store,
	}
}

func (sr *SongResource) Routes() chi.Router {
	r := chi.NewRouter()

	// RESTy routes for "songs" resource
	r.Post("/", sr.CreateSong)
	r.Get("/", sr.AllSongs)
	r.Get("/{id}", sr.ByID)
	r.Put("/", sr.UpdateSong)
	r.Delete("/{id}", sr.DeleteSong)

	return r
}

func (sr *SongResource) CreateSong(rw http.ResponseWriter, r *http.Request) {
	song := new(models.Song)
	if err := json.NewDecoder(r.Body).Decode(song); err != nil {
		rw.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = fmt.Fprintf(rw, "Unknown err: %v", err)
		return
	}

	if err := sr.store.Songs().Create(r.Context(), song); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(rw, "DB err: %v", err)
		return
	}

	rw.WriteHeader(http.StatusCreated)
}

func (sr *SongResource) AllSongs(rw http.ResponseWriter, r *http.Request) {
	// TODO: add request parameter 'expand=True' to return all songs with their artist info
	songs, err := sr.store.Songs().All(r.Context())
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(rw, "DB err: %v", err)
		return
	}

	render.JSON(rw, r, songs)
}

func (sr *SongResource) ByID(rw http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(rw, "Unknown err: %v", err)
		return
	}

	song, err := sr.store.Songs().ByID(r.Context(), id)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(rw, "DB err: %v", err)
		return
	}

	render.JSON(rw, r, song)
}

func (sr *SongResource) UpdateSong(rw http.ResponseWriter, r *http.Request) {
	song := new(models.Song)
	if err := json.NewDecoder(r.Body).Decode(song); err != nil {
		rw.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = fmt.Fprintf(rw, "Unknown err: %v", err)
		return
	}

	err := validation.ValidateStruct(
		song,
		validation.Field(&song.ID, validation.Required),
		validation.Field(&song.Title, validation.Required),
	)
	if err != nil {
		rw.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = fmt.Fprintf(rw, "Unknown err: %v", err)
		return
	}

	if err := sr.store.Songs().Update(r.Context(), song); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(rw, "DB err: %v", err)
		return
	}
}

func (sr *SongResource) DeleteSong(rw http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(rw, "Unknown err: %v", err)
		return
	}

	if err := sr.store.Songs().Delete(r.Context(), id); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(rw, "DB err: %v", err)
		return
	}
}
