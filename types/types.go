package types

type User struct {
	ID        int    `json:"id"`
	SocialId  string `json:"-"`
	Name      string `json:"name"`
	AvatarURL string `json:"-"`
}

type Movie struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Genres      []string `json:"genres"`
	IMDBLink    string   `json:"imdbLink"`
	IMDBRating  float32  `json:"imdbRating"`
	AvatarLink  string   `json:"avatar"`
	AddedBy     int      `json:"addedBy"`
}

type MovieRating struct {
	ID      int `json:"id"`
	MovieID int `json:"movieId"`
	Rating  int `json:"rating"`
}

type Group struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedBy   int    `json:"createdBy"`
}
