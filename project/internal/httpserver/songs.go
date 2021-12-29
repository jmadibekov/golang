package httpserver

import (
	"context"
	"encoding/json"
	"example/hello/project/internal/message_broker"
	"example/hello/project/internal/models"
	"example/hello/project/internal/store"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	lru "github.com/hashicorp/golang-lru"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

type SongResource struct {
	store  store.Store
	broker message_broker.MessageBroker
	cache  *lru.TwoQueueCache
}

func NewSongResource(store store.Store, broker message_broker.MessageBroker, cache *lru.TwoQueueCache) *SongResource {
	return &SongResource{
		store:  store,
		broker: broker,
		cache:  cache,
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

func validateSong(song *models.Song) error {
	return validation.ValidateStruct(
		song,
		validation.Field(&song.ID, validation.Required),
		validation.Field(&song.Title, validation.Required),
		validation.Field(&song.ArtistID, validation.Required),
	)
}

func (sr *SongResource) setTemporaryStatus(song *models.Song) error {
	song.Lyrics = "fetching lyrics..."
	return sr.store.Songs().Update(context.TODO(), song)
}

func (sr *SongResource) setFailedStatus(song *models.Song) error {
	song.Lyrics = "failed"
	return sr.store.Songs().Update(context.TODO(), song)
}

func (sr *SongResource) setFinalStatus(song *models.Song) error {
	return sr.store.Songs().Update(context.TODO(), song)
}

const APIRootURL = "https://api.musixmatch.com/ws/1.1"
const APIKey = "e2dd130dd5117a2e12cbb07d1af40373"

func makeAPICall(apiURL string) ([]byte, error) {
	log.Println("making a GET call to", apiURL)

	response, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(response.Body)

	// checking if status is 200
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("received disturbing status code %v", response)
	}

	// reading the body from the response
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// Returns the "track_id" needed to get the lyrics of the song
func (sr *SongResource) searchForSong(song *models.Song) (string, error) {
	artist, err := sr.store.Songs().GetArtist(context.TODO(), song)
	if err != nil {
		return "", err
	}

	apiURL, err := url.Parse(APIRootURL)
	if err != nil {
		return "", err
	}
	apiURL.Path = path.Join(apiURL.Path, "track.search")
	q := apiURL.Query()
	q.Set("apikey", APIKey)
	q.Set("q_track", song.Title)
	q.Set("q_artist", artist.FullName)
	// sorted by track & artist rating, so that we get the most relevant first
	q.Set("s_track_rating", "desc")
	q.Set("s_artist_rating", "desc")
	apiURL.RawQuery = q.Encode()

	response, err := makeAPICall(apiURL.String())
	if err != nil {
		return "", err
	}
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(response), &result); err != nil {
		return "", err
	}
	result = result["message"].(map[string]interface{})
	result = result["body"].(map[string]interface{})

	tracks := result["track_list"].([]interface{})
	// for simplicity, we're only considering the first one
	// that would be the most relevant
	track := tracks[0].(map[string]interface{})
	track = track["track"].(map[string]interface{})

	song.AlbumName = track["album_name"].(string)

	trackID := strconv.FormatFloat(track["track_id"].(float64), 'f', -1, 64)
	return trackID, nil
}

// Fetches the lyrics for the song
func (sr *SongResource) getLyrics(song *models.Song, trackID string) error {
	apiURL, err := url.Parse(APIRootURL)
	if err != nil {
		return err
	}
	apiURL.Path = path.Join(apiURL.Path, "track.lyrics.get")
	q := apiURL.Query()
	q.Set("apikey", APIKey)
	q.Set("track_id", trackID)
	apiURL.RawQuery = q.Encode()

	response, err := makeAPICall(apiURL.String())
	if err != nil {
		return err
	}
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(response), &result); err != nil {
		return err
	}
	result = result["message"].(map[string]interface{})
	result = result["body"].(map[string]interface{})
	result = result["lyrics"].(map[string]interface{})

	song.Lyrics = result["lyrics_body"].(string)

	return nil
}

// Simple Saga Pattern
// Two services: , the other with external API
//  A) one works with DB
//		- setTemporaryStatus
//		- setFailedStatus
//		- setFinalStatus
//	B) the other works with external API
//		- searchForSong
//		- getLyrics
func (sr *SongResource) fetchTheLyrics(song *models.Song) error {
	log.Println("starting to get the lyrics for", song)

	if err := sr.setTemporaryStatus(song); err != nil {
		return err
	}

	trackID, err := sr.searchForSong(song)
	if err != nil {
		if err := sr.setFailedStatus(song); err != nil {
			return err
		}
		log.Println("failed on searching for song; err:", err)
		return nil
	}

	if err := sr.getLyrics(song, trackID); err != nil {
		if err := sr.setFailedStatus(song); err != nil {
			return err
		}
		log.Println("failed on getting the lyrics for the song; err:", err)
		return nil
	}

	if err := sr.setFinalStatus(song); err != nil {
		return err
	}
	return nil
}

func (sr *SongResource) CreateSong(rw http.ResponseWriter, r *http.Request) {
	song := new(models.Song)
	if err := json.NewDecoder(r.Body).Decode(song); err != nil {
		rw.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = fmt.Fprintf(rw, "Unknown err: %v", err)
		return
	}

	if err := validateSong(song); err != nil {
		rw.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = fmt.Fprintf(rw, "Validation err: %v", err)
		return
	}

	if err := sr.store.Songs().Create(r.Context(), song); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(rw, "DB err: %v", err)
		return
	}

	if err := sr.broker.Cache().Purge(); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(rw, "Received error while purging cache: %v", err)
		return
	}

	if err := sr.fetchTheLyrics(song); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Println("Received error while fetching lyrics:", err)
		return
	}

	rw.WriteHeader(http.StatusCreated)
}

func (sr *SongResource) AllSongs(rw http.ResponseWriter, r *http.Request) {
	dataFromCache, ok := sr.cache.Get(r.RequestURI)
	if ok {
		log.Println("found data from cache with URI =", r.RequestURI)
		render.JSON(rw, r, dataFromCache)
		return
	}

	// TODO: add request parameter 'expand=True' to return all songs with their artist info
	queryValues := r.URL.Query()
	searchQuery := queryValues.Get("searchQuery")

	filter := &models.Filter{}
	if searchQuery != "" {
		filter.Query = &searchQuery
	}

	songs, err := sr.store.Songs().All(r.Context(), filter)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(rw, "DB err: %v", err)
		return
	}

	sr.cache.Add(r.RequestURI, songs)
	render.JSON(rw, r, songs)
}

func (sr *SongResource) ByID(rw http.ResponseWriter, r *http.Request) {
	dataFromCache, ok := sr.cache.Get(r.RequestURI)
	if ok {
		log.Println("found data from cache with URI =", r.RequestURI)
		render.JSON(rw, r, dataFromCache)
		return
	}

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

	sr.cache.Add(r.RequestURI, song)
	render.JSON(rw, r, song)
}

func (sr *SongResource) UpdateSong(rw http.ResponseWriter, r *http.Request) {
	song := new(models.Song)
	if err := json.NewDecoder(r.Body).Decode(song); err != nil {
		rw.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = fmt.Fprintf(rw, "Unknown err: %v", err)
		return
	}

	if err := validateSong(song); err != nil {
		rw.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = fmt.Fprintf(rw, "Validation err: %v", err)
		return
	}

	if err := sr.store.Songs().Update(r.Context(), song); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(rw, "DB err: %v", err)
		return
	}

	if err := sr.broker.Cache().Purge(); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(rw, "Received error while purging cache: %v", err)
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

	if err := sr.broker.Cache().Purge(); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(rw, "Received error while purging cache: %v", err)
		return
	}
}
