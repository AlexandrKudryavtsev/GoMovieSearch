package entity

type Movie struct {
	RussianTitle  string `json:"russianTitle"`
	OriginalTitle string `json:"originalTitle"`
	Year          string `json:"year"`
	DetailsURL    string `json:"detailsUrl"`
	PosterURL     string `json:"posterUrl"`
}

type MovieSearchResult struct {
	Movies []Movie `json:"movies"`
	Total  int     `json:"total"`
}
