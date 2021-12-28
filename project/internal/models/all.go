// couple of remarks:
//
// 1) I could get rid of ID field and use the default `_id` field that mongodb provides
// but in that case, ID field needs to be ObjectID type which is 24 characters long string
// (which is not needed, in my case)
// link: https://stackoverflow.com/questions/55921098/making-a-unique-field-in-mongo-go-driver

package models

type Artist struct {
	ID       int    `json:"id"`
	FullName string `json:"full_name"`
}

type Song struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	ArtistID int    `json:"artist_id"`
}

type Filter struct {
	Query *string `json:"query"`
}
