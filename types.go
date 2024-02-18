package main

type User struct {
	ID       int
	SocialId string
	Name     string
}

type Movie struct {
	ID       int
	IMDBLink string
	Genres   []string
	AddedBy  int
}

type MovieRating struct {
	ID      int
	MovieID int
	Rating  int
}

type Group struct {
	ID   int
	Name string
}
