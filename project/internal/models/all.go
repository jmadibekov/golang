package models

type Artist struct {
	ID       int    `json:"id"`
	FullName string `json:"full_name"`
}

var SampleArtists = map[int]*Artist{
	1: {ID: 1, FullName: "Ed Sheeran"},
	2: {ID: 2, FullName: "Coldplay"},
}

type Song struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	ArtistID int    `json:"artist_id"`
}

var SampleSongs = map[int]*Song{
	100: {ID: 100, Name: "English Rose", ArtistID: 1},
	200: {ID: 200, Name: "Up&Up", ArtistID: 2},
}
