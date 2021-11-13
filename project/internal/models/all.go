package models

type Artist struct {
	ID       int    `json:"id"`
	FullName string `json:"full_name"`
}

// you could use this sample data to populate db
var SampleArtists = map[int]*Artist{
	1: {ID: 1, FullName: "Ed Sheeran"},
	2: {ID: 2, FullName: "Coldplay"},
}

type Song struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	ArtistID int    `json:"artist_id"`
}

// you could use this sample data to populate db
var SampleSongs = map[int]*Song{
	100: {ID: 100, Title: "English Rose", ArtistID: 1},
	200: {ID: 200, Title: "Up&Up", ArtistID: 2},
}
