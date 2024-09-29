package entities

type Song struct {
	ID          int    `json:"id" example:"1" description:"Song ID"`
	Group       string `json:"group" example:"Muse" description:"Group or band name"`
	Song        string `json:"song" example:"Supermassive Black Hole" description:"Song title"`
	ReleaseDate string `json:"release_date" example:"2006-06-19" description:"Song release date"`
	Text        string `json:"text" example:"Lyrics of the song" description:"Lyrics of the song"`
	Link        string `json:"link" example:"http://example.com/song" description:"Link to the song"`
}

type CreateSongRequest struct {
	Group string `json:"group" example:"Muse" description:"Название группы"`
	Song  string `json:"song" example:"Supermassive Black Hole" description:"Название песни"`
}

type SongsResponse struct {
	Data  []Song `json:"songs" example:"[{\"id\":1, \"group\":\"Muse\", \"song\":\"Supermassive Black Hole\"}]"`
	Total int    `json:"total" example:"100"`
}

type Error struct {
	Code string `json:"code" example:"400"`
	Msg  string `json:"msg" example:"Invalid request"`
}

type ErrorResponse struct {
	ErrorInfo Error `json:"error"`
}
